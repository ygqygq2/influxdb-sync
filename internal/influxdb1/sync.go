package influxdb1

import (
	"context"
	"fmt"
	"os"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/ygqygq2/influxdb-sync/internal/logx"
)

type SyncConfig struct {
	SourceAddr string
	SourceUser string
	SourcePass string
	SourceDB   string
	TargetAddr string
	TargetUser string
	TargetPass string
	TargetDB   string
	BatchSize  int
	Start      string // 起始时间
	End        string // 结束时间
	ResumeFile string
	Parallel   int    // 并发同步表数，默认4
}

// 兼容新配置结构体
type ExtraConfig struct {
	SourceDBExclude []string
	TargetDBPrefix  string
	TargetDBSuffix  string
}

func getExtraConfig() ExtraConfig {
	// 这里假设通过环境变量传递，实际项目应从 config 结构体传递
	// 你可根据自己的 config 结构体调整此处
	// 这里只做演示，实际应从 main.go 传递参数
	return ExtraConfig{
		SourceDBExclude: []string{},
		TargetDBPrefix:  "",
		TargetDBSuffix:  "",
	}
}

func containsString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
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
	extra := getExtraConfig()
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
							if containsString(extra.SourceDBExclude, name) {
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

		// worker goroutine
		worker := func() {
			// 目标库名加前后缀
			targetDB := extra.TargetDBPrefix + db + extra.TargetDBSuffix
			for m := range jobs {
				logx.Info("同步 measurement:", m)
				var lastTime = startTime
				for {
					if lastTime == "" {
						lastTime = "1970-01-01T00:00:00Z"
					}
					var q string
					if cfg.End != "" {
						q = fmt.Sprintf("SELECT * FROM %s WHERE time > '%s' AND time < '%s' LIMIT %d", m, lastTime, cfg.End, batchSize)
					} else {
						q = fmt.Sprintf("SELECT * FROM %s WHERE time > '%s' LIMIT %d", m, lastTime, batchSize)
					}
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
					for _, result := range res.Results {
						for _, series := range result.Series {
							bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
								Database: targetDB,
								Precision: "ns",
							})
							colIdx := map[string]int{}
							for idx, col := range series.Columns {
								colIdx[col] = idx
							}
							for _, row := range series.Values {
								tags := map[string]string{}
								fields := map[string]interface{}{}
								var t time.Time
								for col, idx := range colIdx {
									if col == "time" {
										switch v := row[idx].(type) {
										case string:
											tt, _ := time.Parse(time.RFC3339Nano, v)
											t = tt
										case time.Time:
											t = v
										}
									} else if col == "host" || col == "region" {
										if s, ok := row[idx].(string); ok {
											tags[col] = s
										}
									} else {
										fields[col] = row[idx]
									}
								}
								pt, err := client.NewPoint(series.Name, tags, fields, t)
								if err == nil {
									bp.AddPoint(pt)
									moved++
									lastTime = t.Format(time.RFC3339Nano)
								}
							}
							if err := tgt.cli.Write(bp); err != nil {
								logx.Error("写入目标库失败:", err)
								results <- syncResult{m, err}
								return
							}
						}
					}
					// 记录断点
					if moved > 0 && cfg.ResumeFile != "" {
						os.WriteFile(cfg.ResumeFile, []byte(lastTime), 0644)
					}
					if moved < batchSize {
						break // 本表已同步完毕
					}
				}
				results <- syncResult{m, nil}
			}
		}

		// 启动 worker
		for i := 0; i < parallel; i++ {
			go worker()
		}
		// 分发任务
		go func() {
			for _, m := range measurements {
				jobs <- m
			}
			close(jobs)
		}()
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
	logx.Info("全部同步完成")
	return nil
}
