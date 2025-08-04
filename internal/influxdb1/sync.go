package influxdb1

import (
	"context"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

// 1.x 到 1.x 同步
func Sync1x1x(ctx context.Context, cfg common.SyncConfig) error {
	// 创建源和目标
	source := NewDataSource(DataSourceConfig{
		Addr: cfg.SourceAddr,
		User: cfg.SourceUser,
		Pass: cfg.SourcePass,
	})

	target := NewDataTarget(DataTargetConfig{
		Addr: cfg.TargetAddr,
		User: cfg.TargetUser,
		Pass: cfg.TargetPass,
	})

	// 创建同步器
	syncer := common.NewSyncer(cfg, source, target)

	// 执行同步
	return syncer.Sync(ctx)
}

// 保持向后兼容的旧接口
func Sync(ctx context.Context, cfg SyncConfig) error {
	// 转换为新的配置格式
	newCfg := common.SyncConfig{
		SourceAddr:      cfg.SourceAddr,
		SourceUser:      cfg.SourceUser,
		SourcePass:      cfg.SourcePass,
		SourceDB:        cfg.SourceDB,
		SourceDBExclude: cfg.SourceDBExclude,
		TargetAddr:      cfg.TargetAddr,
		TargetUser:      cfg.TargetUser,
		TargetPass:      cfg.TargetPass,
		TargetDB:        cfg.TargetDB,
		TargetDBPrefix:  cfg.TargetDBPrefix,
		TargetDBSuffix:  cfg.TargetDBSuffix,
		BatchSize:       cfg.BatchSize,
		Start:           cfg.Start,
		End:             cfg.End,
		ResumeFile:      cfg.ResumeFile,
		Parallel:        cfg.Parallel,
		RetryCount:      cfg.RetryCount,
		RetryInterval:   cfg.RetryInterval,
		RateLimit:       cfg.RateLimit,
		LogLevel:        cfg.LogLevel,
	}

	return Sync1x1x(ctx, newCfg)
}
