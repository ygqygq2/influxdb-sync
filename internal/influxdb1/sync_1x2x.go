package influxdb1

import (
	"context"

	"github.com/ygqygq2/influxdb-sync/internal/common"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb2"
)

// 1.x 到 2.x 同步
func Sync1x2x(ctx context.Context, cfg common.SyncConfig) error {
	// 创建源和目标
	source := NewDataSource(DataSourceConfig{
		Addr: cfg.SourceAddr,
		User: cfg.SourceUser,
		Pass: cfg.SourcePass,
	})

	target := &influxdb2.Adapter{
		URL:    cfg.TargetAddr,
		Token:  cfg.TargetToken,
		Org:    cfg.TargetOrg,
		Bucket: cfg.TargetBucket,
	}

	// 创建同步器
	syncer := common.NewSyncer(cfg, source, target)

	// 执行同步
	return syncer.Sync(ctx)
}
