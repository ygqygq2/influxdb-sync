package influxdb2

import (
	"context"
	"fmt"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

// Sync2x2x 执行 InfluxDB 2.x 到 2.x 的同步
func Sync2x2x(ctx context.Context, cfg common.SyncConfig) error {
	// 创建源适配器
	source := &Adapter{
		URL:    cfg.SourceAddr,
		Token:  cfg.SourceToken,
		Org:    cfg.SourceOrg,
		Bucket: cfg.SourceBucket,
	}

	// 创建目标适配器
	target := &Adapter{
		URL:    cfg.TargetAddr,
		Token:  cfg.TargetToken,
		Org:    cfg.TargetOrg,
		Bucket: cfg.TargetBucket,
	}

	// 如果目标 bucket 为空，使用前后缀拼接逻辑
	if target.Bucket == "" {
		target.Bucket = fmt.Sprintf("%s%s%s", cfg.TargetDBPrefix, source.Bucket, cfg.TargetDBSuffix)
	}

	// 创建同步器
	syncCfg := common.SyncConfig{
		BatchSize:     cfg.BatchSize,
		Start:         cfg.Start,
		End:           cfg.End,
		ResumeFile:    cfg.ResumeFile,
		Parallel:      cfg.Parallel,
		RetryCount:    cfg.RetryCount,
		RetryInterval: cfg.RetryInterval,
		RateLimit:     cfg.RateLimit,
		LogLevel:      cfg.LogLevel,
	}
	syncer := common.NewSyncer(syncCfg, source, target)

	// 执行同步
	return syncer.Sync(ctx)
}
