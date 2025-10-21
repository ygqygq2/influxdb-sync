package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/ygqygq2/influxdb-sync/internal/common"
	"github.com/ygqygq2/influxdb-sync/internal/config"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb1"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb2"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb3"
)

// detectSyncMode 根据配置自动识别同步模式
func detectSyncMode(cfg *config.Config) string {
	// 检查源端版本
	sourceVersion := ""
	switch cfg.Source.Type {
	case 1:
		sourceVersion = "1x"
	case 2:
		sourceVersion = "2x"
	case 3:
		sourceVersion = "3x"
	default:
		// 兼容旧配置，通过字段判断
		if cfg.Source.Token != "" {
			if cfg.Source.Database != "" {
				sourceVersion = "3x"
			} else {
				sourceVersion = "2x"
			}
		} else {
			sourceVersion = "1x"
		}
	}

	// 检查目标端版本
	targetVersion := ""
	switch cfg.Target.Type {
	case 1:
		targetVersion = "1x"
	case 2:
		targetVersion = "2x"
	case 3:
		targetVersion = "3x"
	default:
		// 兼容旧配置，通过字段判断
		if cfg.Target.Token != "" {
			if cfg.Target.Database != "" {
				targetVersion = "3x"
			} else {
				targetVersion = "2x"
			}
		} else {
			targetVersion = "1x"
		}
	}

	return sourceVersion + targetVersion
}

// Run 执行同步命令，自动识别版本
func Run(cfgPath string) error {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	// 自动识别同步模式
	mode := detectSyncMode(cfg)

	// 转换配置为通用格式
	// 对于3.x版本，优先使用Database字段；对于1.x/2.x版本，使用DB字段
	sourceDB := cfg.Source.DB
	if cfg.Source.Type == 3 && cfg.Source.Database != "" {
		sourceDB = cfg.Source.Database
	}

	targetDB := cfg.Target.DB
	if cfg.Target.Type == 3 && cfg.Target.Database != "" {
		targetDB = cfg.Target.Database
	}

	syncConfig := common.SyncConfig{
		SourceAddr:      cfg.Source.URL,
		SourceUser:      cfg.Source.User,
		SourcePass:      cfg.Source.Pass,
		SourceDB:        sourceDB,
		SourceDBExclude: cfg.Source.DBExclude,
		SourceToken:     cfg.Source.Token,
		SourceOrg:       cfg.Source.Org,
		SourceBucket:    cfg.Source.Bucket,
		TargetAddr:      cfg.Target.URL,
		TargetUser:      cfg.Target.User,
		TargetPass:      cfg.Target.Pass,
		TargetDB:        targetDB,
		TargetDBPrefix:  cfg.Target.DBPrefix,
		TargetDBSuffix:  cfg.Target.DBSuffix,
		TargetToken:     cfg.Target.Token,
		TargetOrg:       cfg.Target.Org,
		TargetBucket:    cfg.Target.Bucket,
		BatchSize:       cfg.Sync.BatchSize,
		Start:           cfg.Sync.Start,
		End:             cfg.Sync.End,
		ResumeFile:      cfg.Sync.ResumeFile,
		Parallel:        cfg.Sync.Parallel,
		RetryCount:      cfg.Sync.RetryCount,
		RetryInterval:   cfg.Sync.RetryInterval,
		RateLimit:       cfg.Sync.RateLimit,
		LogLevel:        cfg.Log.Level,
	}

	// 根据模式选择同步方式
	switch strings.ToLower(mode) {
	case "1x1x", "1x-1x", "":
		// 兼容原有的1x-1x同步
		return runInfluxdb1Sync(syncConfig)
	case "1x2x", "1x-2x":
		// 1x到2x同步
		return influxdb1.Sync1x2x(context.Background(), syncConfig)
	case "2x2x", "2x-2x":
		// 2x到2x同步
		return influxdb2.Sync2x2x(context.Background(), syncConfig)
	case "1x3x", "1x-3x":
		// 1x到3x同步
		return influxdb3.Sync1x3x(context.Background(), syncConfig)
	case "2x3x", "2x-3x":
		// 2x到3x同步
		return influxdb3.Sync2x3x(context.Background(), syncConfig)
	case "3x3x", "3x-3x":
		// 3x到3x同步
		return influxdb3.Sync3x3x(context.Background(), syncConfig)
	default:
		return fmt.Errorf("不支持的同步模式: %s，支持的模式: 1x1x, 1x2x, 2x2x, 1x3x, 2x3x, 3x3x", mode)
	}
}

// runInfluxdb1Sync 执行原有的1x-1x同步
func runInfluxdb1Sync(cfg common.SyncConfig) error {
	// 转换为influxdb1的配置格式
	c := influxdb1.SyncConfig{
		SourceAddr: cfg.SourceAddr,
		SourceUser: cfg.SourceUser,
		SourcePass: cfg.SourcePass,
		SourceDB:   cfg.SourceDB,
		TargetAddr: cfg.TargetAddr,
		TargetUser: cfg.TargetUser,
		TargetPass: cfg.TargetPass,
		TargetDB:   cfg.TargetDB,
		BatchSize:  cfg.BatchSize,
		ResumeFile: cfg.ResumeFile,
	}
	return influxdb1.Sync(context.Background(), c)
}

// ShowUsage 显示使用说明
func ShowUsage() {
	fmt.Println("InfluxDB 同步工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  influxdb-sync <config.yaml>")
	fmt.Println("")
	fmt.Println("参数:")
	fmt.Println("  config.yaml  配置文件路径")
	fmt.Println("")
	fmt.Println("支持的同步场景 (自动识别):")
	fmt.Println("  InfluxDB 1.x 到 1.x")
	fmt.Println("  InfluxDB 1.x 到 2.x")
	fmt.Println("  InfluxDB 2.x 到 2.x")
	fmt.Println("  InfluxDB 1.x 到 3.x")
	fmt.Println("  InfluxDB 2.x 到 3.x")
	fmt.Println("  InfluxDB 3.x 到 3.x")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  influxdb-sync config.yaml")
	fmt.Println("  influxdb-sync config_1x2x.yaml")
	fmt.Println("  influxdb-sync config_2x2x.yaml")
	fmt.Println("  influxdb-sync config_1x3x.yaml")
	fmt.Println("  influxdb-sync config_2x3x.yaml")
	fmt.Println("  influxdb-sync config_3x3x.yaml")
}
