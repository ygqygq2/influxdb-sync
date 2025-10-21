package influxdb3

import (
	"context"
	"fmt"

	"github.com/ygqygq2/influxdb-sync/internal/common"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb1"
)

// Sync1x3x 执行 InfluxDB 1.x 到 3.x 的同步
func Sync1x3x(ctx context.Context, cfg common.SyncConfig) error {
	// 创建源适配器 (InfluxDB 1.x)
	source := influxdb1.NewDataSource(influxdb1.DataSourceConfig{
		Addr: cfg.SourceAddr,
		User: cfg.SourceUser,
		Pass: cfg.SourcePass,
	})

	// 创建目标适配器 (InfluxDB 3.x, v2 兼容模式)
	var targetBucket string
	if cfg.TargetBucket != "" {
		targetBucket = cfg.TargetBucket
	} else {
		targetBucket = fmt.Sprintf("%s%s%s", cfg.TargetDBPrefix, cfg.SourceDB, cfg.TargetDBSuffix)
	}

	targetConfig := V2CompatConfig{
		URL:      cfg.TargetAddr,
		Token:    cfg.TargetToken,
		Org:      cfg.TargetOrg,
		Database: targetBucket,
	}
	target := &DataTarget3x{
		client: nil, // 将在连接时初始化
		config: targetConfig,
	}

	// 创建同步器
	syncer := common.NewSyncer(cfg, source, target)

	// 执行同步
	return syncer.Sync(ctx)
}
