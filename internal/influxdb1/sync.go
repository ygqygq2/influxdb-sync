package influxdb1

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/ygqygq2/influxdb-sync/internal/logx"
)

// measurement 名称转义，双引号包裹并转义内部双引号
func escapeMeasurement(m string) string {
	return "\"" + strings.ReplaceAll(m, "\"", "\\\"") + "\""
}



func Sync(ctx context.Context, cfg SyncConfig) error {
	src, err := NewClient(cfg.SourceAddr, cfg.SourceUser, cfg.SourcePass, 0)
	if err != nil {
		logx.Error("源库连接失败:", err)
		return err
	}
	defer src.Close()
	tgt, err := NewClient(cfg.TargetAddr, cfg.TargetUser, cfg.TargetPass, 0)
	if err != nil {
		logx.Error("目标库连接失败:", err)
		return err
	}
	defer tgt.Close()

   // 断点续传：读取配置起始时间和 resume 文件
   startTime := cfg.Start
   if startTime == "" {
	   startTime = "1970-01-01T00:00:00Z"
   }
   if cfg.ResumeFile != "" {
	   if data, err := os.ReadFile(cfg.ResumeFile); err == nil {
		   last := string(data)
		   // 取配置 start 与 resume 文件中时间的较大者
		   if t0, err0 := time.Parse(time.RFC3339Nano, startTime); err0 == nil {
			   if t1, err1 := time.Parse(time.RFC3339Nano, last); err1 == nil && t1.After(t0) {
				   startTime = last
			   }
		   }
	   }
   }

	// 查询所有数据库
	cli := src.cli
   dbs := []string{}
   if cfg.SourceDB == "" {
	   dbRes, err := cli.Query(client.NewQuery("SHOW DATABASES", "", ""))
	   if err != nil {
		   logx.Error("SHOW DATABASES 查询失败:", err)
		   return err
	   }
	   if dbRes.Error() != nil {
		   logx.Error("SHOW DATABASES 响应错误:", dbRes.Error())
		   return dbRes.Error()
	   }
	   for _, result := range dbRes.Results {
		   for _, series := range result.Series {
			   for _, v := range series.Values {
				   if len(v) > 0 {
					   if name, ok := v[0].(string); ok {
						   if name == "_internal" {
							   continue
						   }
						   // 过滤
						   if containsString(cfg.SourceDBExclude, name) {
							   continue
						   }
						   dbs = append(dbs, name)
					   }
				   }
			   }
		   }
	   }
   } else {
	   dbs = append(dbs, cfg.SourceDB)
   }

	for _, db := range dbs {
		logx.Info("同步数据库:", db)
		// 查询所有 measurement
		showRes, err := cli.Query(client.NewQuery("SHOW MEASUREMENTS", db, ""))
		if err != nil {
			logx.Error("SHOW MEASUREMENTS 查询失败:", err)
			return err
		}
		if showRes.Error() != nil {
			logx.Error("SHOW MEASUREMENTS 响应错误:", showRes.Error())
			return showRes.Error()
		}
		var measurements []string
		for _, result := range showRes.Results {
			for _, series := range result.Series {
				for _, v := range series.Values {
					if len(v) > 0 {
						if name, ok := v[0].(string); ok {
							measurements = append(measurements, name)
						}
					}
				}
			}
		}
		if len(measurements) == 0 {
			logx.Warn("库", db, "无数据表")
			continue
		}
		batchSize := cfg.BatchSize
		if batchSize <= 0 {
			batchSize = 1000
		}
		parallel := cfg.Parallel
		if parallel <= 0 {
			parallel = 4
		}
		type syncResult struct {
			m   string
			err error
		}
	   jobs := make(chan string)
	   results := make(chan syncResult)

	  targetDB := cfg.TargetDBPrefix + db + cfg.TargetDBSuffix

	   // worker 函数定义
   // worker 函数定义：一次性查询全部数据并同步，适用于数据量适中场景
   worker := func(idx, total int) {
	   for m := range jobs {
		   // 构建查询SQL，不使用分页
		   em := escapeMeasurement(m)
		   q := fmt.Sprintf("SELECT * FROM %s", em)
		   res, err := cli.Query(client.NewQuery(q, db, "ns"))
		   if err != nil {
			   logx.Error("查询", m, "失败:", err)
			   results <- syncResult{m, err}
			   return
		   }
		   if res.Error() != nil {
			   logx.Error("响应错误:", res.Error())
			   results <- syncResult{m, res.Error()}
			   return
		   }
		   moved := 0
		   // 处理返回结果，写入所有点
		   for _, result := range res.Results {
			   for _, series := range result.Series {
				   bp, _ := client.NewBatchPoints(client.BatchPointsConfig{Database: targetDB, Precision: "ns"})
				   colIdx := map[string]int{}
				   for i, col := range series.Columns {
					   colIdx[col] = i
				   }
				   for _, row := range series.Values {
					   tags := map[string]string{}
					   fields := map[string]interface{}{}
					   var t time.Time
					   for col, idx := range colIdx {
						   switch col {
						   case "time":
							   switch v := row[idx].(type) {
							   case string:
								   t, _ = time.Parse(time.RFC3339Nano, v)
							   case time.Time:
								   t = v
							   case int64:
								   t = time.Unix(0, v)
							   case float64:
								   t = time.Unix(0, int64(v))
							   case json.Number:
								   if ns, err := v.Int64(); err == nil {
									   t = time.Unix(0, ns)
								   }
							   }
						   case "host", "region":
							   if s, ok := row[idx].(string); ok {
								   tags[col] = s
							   }
						   default:
							   if val := row[idx]; val != nil {
								   if sv, ok := val.(string); !ok || sv != "" {
									   fields[col] = val
								   }
							   }
						   }
					   }
					   if pt, err := client.NewPoint(series.Name, tags, fields, t); err == nil {
						   bp.AddPoint(pt)
						   moved++
					   }
				   }
				   if err := tgt.cli.Write(bp); err != nil {
					   logx.Error("写入目标库失败:", err)
					   results <- syncResult{m, err}
					   return
				   }
			   }
		   }
		   results <- syncResult{m, nil}
	   }
   }

		// 启动 worker
		for i := 0; i < parallel; i++ {
			go worker(i, len(measurements))
		}
		// 分发任务
		for _, m := range measurements {
			logx.Info("分发 measurement:", m)
			jobs <- m
		}
		close(jobs)
		// 收集结果
		var firstErr error
		for i := 0; i < len(measurements); i++ {
			r := <-results
			if r.err != nil && firstErr == nil {
				firstErr = r.err
			}
		}
		if firstErr != nil {
			return firstErr
		}
	}
	return nil
}
