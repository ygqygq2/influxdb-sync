package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ygqygq2/influxdb-sync/internal/common"
	"github.com/ygqygq2/influxdb-sync/internal/config"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb1"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb2"
)

// Run 执行同步命令
func Run(cfgPath, mode string) error {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	// 转换配置为通用格式
	syncConfig := common.SyncConfig{
		SourceAddr:      cfg.Source.URL,
		SourceUser:      cfg.Source.User,
		SourcePass:      cfg.Source.Pass,
		SourceDB:        cfg.Source.DB,
		SourceDBExclude: cfg.Source.DBExclude,
		SourceToken:     cfg.Source.Token,
		SourceOrg:       cfg.Source.Org,
		SourceBucket:    cfg.Source.Bucket,
		TargetAddr:      cfg.Target.URL,
		TargetUser:      cfg.Target.User,
		TargetPass:      cfg.Target.Pass,
		TargetDB:        cfg.Target.DB,
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
	default:
		return fmt.Errorf("不支持的同步模式: %s，支持的模式: 1x1x, 1x2x, 2x2x", mode)
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
	fmt.Println("  influxdb-sync <config.yaml> [mode]")
	fmt.Println("")
	fmt.Println("参数:")
	fmt.Println("  config.yaml  配置文件路径")
	fmt.Println("  mode         同步模式 (可选)")
	fmt.Println("")
	fmt.Println("支持的同步模式:")
	fmt.Println("  1x1x, 1x-1x  InfluxDB 1.x 到 1.x (默认)")
	fmt.Println("  1x2x, 1x-2x  InfluxDB 1.x 到 2.x")
	fmt.Println("  2x2x, 2x-2x  InfluxDB 2.x 到 2.x")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  influxdb-sync config.yaml")
	fmt.Println("  influxdb-sync config_1x2x.yaml 1x2x")
	fmt.Println("  influxdb-sync config_2x2x.yaml 2x2x")
}

func main() {
	if len(os.Args) < 2 {
		ShowUsage()
		os.Exit(1)
	}

	cfgPath := os.Args[1]
	mode := ""
	if len(os.Args) >= 3 {
		mode = os.Args[2]
	}

	if err := Run(cfgPath, mode); err != nil {
		fmt.Printf("同步失败: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("同步完成")
}
