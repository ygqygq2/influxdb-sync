package common

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ygqygq2/influxdb-sync/internal/logx"
)

// 通用同步器
type Syncer struct {
	cfg    SyncConfig
	source DataSource
	target DataTarget
}

// 创建新的同步器
func NewSyncer(cfg SyncConfig, source DataSource, target DataTarget) *Syncer {
	return &Syncer{
		cfg:    cfg,
		source: source,
		target: target,
	}
}

// 执行同步
func (s *Syncer) Sync(ctx context.Context) error {
	// 连接源和目标
	if err := s.source.Connect(); err != nil {
		logx.Error("源库连接失败:", err)
		return err
	}
	defer s.source.Close()

	if err := s.target.Connect(); err != nil {
		logx.Error("目标库连接失败:", err)
		return err
	}
	defer s.target.Close()

	// 获取起始时间
	startTimeNano, err := s.getStartTime()
	if err != nil {
		return err
	}

	// 获取数据库列表
	dbs, err := s.getDatabases()
	if err != nil {
		return err
	}

	// 同步每个数据库
	for _, db := range dbs {
		if err := s.syncDatabase(ctx, db, startTimeNano); err != nil {
			return err
		}
	}

	return nil
}

// 获取起始时间
func (s *Syncer) getStartTime() (int64, error) {
	startTime := s.cfg.Start
	if startTime == "" {
		startTime = "1970-01-01T00:00:00Z"
	}

	var startTimeNano int64 = 0

	if s.cfg.ResumeFile != "" {
		if data, err := os.ReadFile(s.cfg.ResumeFile); err == nil {
			last := string(data)
			// 取配置 start 与 resume 文件中时间的较大者
			if t0, err0 := time.Parse(time.RFC3339Nano, startTime); err0 == nil {
				startTimeNano = t0.UnixNano()
				if t1, err1 := time.Parse(time.RFC3339Nano, last); err1 == nil && t1.After(t0) {
					startTimeNano = t1.UnixNano()
				}
			}
		} else {
			if t0, err0 := time.Parse(time.RFC3339Nano, startTime); err0 == nil {
				startTimeNano = t0.UnixNano()
			}
		}
	} else {
		if t0, err0 := time.Parse(time.RFC3339Nano, startTime); err0 == nil {
			startTimeNano = t0.UnixNano()
		}
	}

	return startTimeNano, nil
}

// 获取数据库列表
func (s *Syncer) getDatabases() ([]string, error) {
	if s.cfg.SourceDB != "" {
		return []string{s.cfg.SourceDB}, nil
	}

	dbs, err := s.source.GetDatabases()
	if err != nil {
		return nil, err
	}

	// 过滤排除的数据库
	var filteredDBs []string
	for _, db := range dbs {
		if db == "_internal" {
			continue
		}
		if containsString(s.cfg.SourceDBExclude, db) {
			continue
		}
		filteredDBs = append(filteredDBs, db)
	}

	return filteredDBs, nil
}

// 同步单个数据库
func (s *Syncer) syncDatabase(ctx context.Context, db string, startTimeNano int64) error {
	logx.Info("同步数据库:", db)

	// 获取 measurements
	measurements, err := s.source.GetMeasurements(db)
	if err != nil {
		return err
	}

	if len(measurements) == 0 {
		logx.Warn("库", db, "无数据表")
		return nil
	}

	// 设置默认值
	batchSize := s.cfg.BatchSize
	if batchSize <= 0 {
		batchSize = 1000
	}
	parallel := s.cfg.Parallel
	if parallel <= 0 {
		parallel = 4
	}

	// 创建任务通道
	jobs := make(chan string, len(measurements))
	results := make(chan SyncResult, len(measurements))

	// 启动 worker
	for i := 0; i < parallel; i++ {
		go s.worker(ctx, db, startTimeNano, batchSize, jobs, results)
	}

	// 分发任务
	for _, m := range measurements {
		logx.Info("分发 measurement:", m)
		jobs <- m
	}
	close(jobs)

	// 收集结果
	var allErrors []error
	for i := 0; i < len(measurements); i++ {
		r := <-results
		if r.Error != nil {
			allErrors = append(allErrors, fmt.Errorf("measurement %s: %v", r.Measurement, r.Error))
		}
	}

	if len(allErrors) > 0 {
		for _, err := range allErrors {
			logx.Error("同步错误:", err)
		}
		return fmt.Errorf("同步失败，共%d个错误", len(allErrors))
	}

	return nil
}

// 工作协程
func (s *Syncer) worker(ctx context.Context, db string, startTimeNano int64, batchSize int, jobs <-chan string, results chan<- SyncResult) {
	for measurement := range jobs {
		logx.Info(fmt.Sprintf("开始处理 measurement: %s", measurement))
		start := time.Now()
		if err := s.syncMeasurement(ctx, db, measurement, startTimeNano, batchSize); err != nil {
			logx.Error(fmt.Sprintf("处理 measurement %s 失败，耗时: %v，错误: %v", measurement, time.Since(start), err))
			results <- SyncResult{Measurement: measurement, Error: err}
		} else {
			logx.Info(fmt.Sprintf("处理 measurement %s 成功，耗时: %v", measurement, time.Since(start)))
			results <- SyncResult{Measurement: measurement, Error: nil}
		}
	}
}

// 同步单个 measurement
func (s *Syncer) syncMeasurement(ctx context.Context, db, measurement string, startTimeNano int64, batchSize int) error {
	// 获取标签字段
	tagKeys, err := s.source.GetTagKeys(db, measurement)
	if err != nil {
		logx.Error("获取", measurement, "标签字段失败:", err)
		return err
	}
	logx.Debug("获取到", measurement, "的标签字段:", tagKeys)

	// 设置重试和限流参数
	retryCount := s.cfg.RetryCount
	if retryCount <= 0 {
		retryCount = 3
	}
	retryInterval := s.cfg.RetryInterval
	if retryInterval <= 0 {
		retryInterval = 500
	}
	rateLimit := s.cfg.RateLimit
	if rateLimit < 0 {
		rateLimit = 50
	}

	var lastTime int64 = startTimeNano

	// 确定目标数据库/bucket名称
	var targetName string
	if s.cfg.TargetBucket != "" {
		// 如果明确配置了 TargetBucket，使用它（适用于固定 bucket 名称）
		targetName = s.cfg.TargetBucket
	} else {
		// 否则使用前后缀拼接源数据库名（适用于动态命名）
		// 对于 1x->1x: targetName = prefix + db + suffix
		// 对于 1x->2x: targetName = prefix + db + suffix (作为 bucket 名)
		targetName = s.cfg.TargetDBPrefix + db + s.cfg.TargetDBSuffix
	}

	for {
		// 查询数据
		logx.Info(fmt.Sprintf("开始查询 %s，起始时间: %d", measurement, lastTime))
		queryStart := time.Now()
		points, maxTime, err := s.source.QueryData(db, measurement, lastTime, batchSize)
		queryDuration := time.Since(queryStart)
		if err != nil {
			logx.Error(fmt.Sprintf("查询 %s 失败，耗时: %v，错误: %v", measurement, queryDuration, err))
			return err
		}
		logx.Info(fmt.Sprintf("查询 %s 完成，耗时: %v，返回 %d 个点", measurement, queryDuration, len(points)))

		if len(points) == 0 {
			logx.Info(fmt.Sprintf("measurement %s 没有更多数据", measurement))
			break // 没有更多数据
		}

		logx.Debug(fmt.Sprintf("处理 %s: %d 个点，时间范围: %d -> %d", measurement, len(points), lastTime, maxTime))

		// 写入目标库，重试机制
		var writeErr error
		for i := 0; i < retryCount; i++ {
			writeErr = s.target.WritePoints(targetName, points)
			if writeErr == nil {
				break
			}
			logx.Warn(fmt.Sprintf("写入目标库失败，第%d次重试: %v", i+1, writeErr))
			time.Sleep(time.Duration(retryInterval) * time.Millisecond)
		}

		if writeErr != nil {
			logx.Error("写入目标库失败，重试失败:", writeErr)
			return writeErr
		}

		logx.Debug(fmt.Sprintf("成功写入 %s: %d 个点", measurement, len(points)))

		// 更新断点续传文件
		if s.cfg.ResumeFile != "" && maxTime > lastTime {
			resumeTime := time.Unix(0, maxTime).Format(time.RFC3339Nano)
			if err := os.WriteFile(s.cfg.ResumeFile, []byte(resumeTime), 0644); err != nil {
				logx.Warn("更新断点续传文件失败:", err)
			}
		}

		lastTime = maxTime

		// 限流控制
		if rateLimit > 0 {
			time.Sleep(time.Duration(rateLimit) * time.Millisecond)
		}

		// 如果本批不足batchSize，说明拉完了
		if len(points) < batchSize {
			break
		}
	}

	return nil
}

// 工具函数：检查字符串数组是否包含指定字符串
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
