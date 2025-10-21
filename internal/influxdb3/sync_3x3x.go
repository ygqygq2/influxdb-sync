package influxdb3

import (
	"context"
	"fmt"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

// Sync3x3x 执行 InfluxDB 3.x 到 3.x 的同步
func Sync3x3x(ctx context.Context, cfg common.SyncConfig) error {
	// 创建源适配器 (InfluxDB 3.x, v1 兼容模式)
	sourceConfig := V1CompatConfig{
		Addr:     cfg.SourceAddr,
		User:     cfg.SourceUser,
		Pass:     cfg.SourcePass,
		Database: cfg.SourceDB,
	}
	source := &DataSource3x{
		client: nil, // 将在连接时初始化
		config: sourceConfig,
	}

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
