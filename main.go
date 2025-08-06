package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ygqygq2/influxdb-sync/internal/common"
	"github.com/ygqygq2/influxdb-sync/internal/config"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb1"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb2"
	"github.com/ygqygq2/influxdb-sync/internal/logx"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: influxdb-sync config.yaml")
		os.Exit(1)
	}
	cfgPath := os.Args[1]
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		fmt.Println("配置加载失败:", err)
		os.Exit(2)
	}

	// 设置日志级别
	logLevel := cfg.Log.Level
	if logLevel == "" {
		logLevel = "info" // 默认级别
	}
	logx.SetLevel(logLevel)

	// 根据 type 字段自动判断同步模式
	if cfg.Source.Type == 1 && cfg.Target.Type == 1 {
		if err := runSync1x1x(cfg); err != nil {
			fmt.Println("同步失败:", err)
			os.Exit(2)
		}
		fmt.Println("同步完成")
	} else if cfg.Source.Type == 1 && cfg.Target.Type == 2 {
		if err := runSync1x2x(cfg); err != nil {
			fmt.Println("同步失败:", err)
			os.Exit(2)
		}
		fmt.Println("同步完成")
	} else if cfg.Source.Type == 2 && cfg.Target.Type == 2 {
		if err := runSync2x2x(cfg); err != nil {
			fmt.Println("同步失败:", err)
			os.Exit(2)
		}
		fmt.Println("同步完成")
	} else {
		fmt.Println("暂不支持该类型同步: source.type=", cfg.Source.Type, ", target.type=", cfg.Target.Type)
		os.Exit(1)
	}
}

func runSync1x1x(cfg *config.Config) error {
	c := influxdb1.SyncConfig{
		SourceAddr:      cfg.Source.URL,
		SourceUser:      cfg.Source.User,
		SourcePass:      cfg.Source.Pass,
		SourceDB:        cfg.Source.DB,
		SourceDBExclude: cfg.Source.DBExclude,
		TargetAddr:      cfg.Target.URL,
		TargetUser:      cfg.Target.User,
		TargetPass:      cfg.Target.Pass,
		TargetDB:        cfg.Target.DB,
		TargetDBPrefix:  cfg.Target.DBPrefix,
		TargetDBSuffix:  cfg.Target.DBSuffix,
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
	return influxdb1.Sync(context.Background(), c)
}

func runSync1x2x(cfg *config.Config) error {
	c := common.SyncConfig{
		SourceAddr:      cfg.Source.URL,
		SourceUser:      cfg.Source.User,
		SourcePass:      cfg.Source.Pass,
		SourceDB:        cfg.Source.DB,
		SourceDBExclude: cfg.Source.DBExclude,
		TargetAddr:      cfg.Target.URL,
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
	return influxdb1.Sync1x2x(context.Background(), c)
}

func runSync2x2x(cfg *config.Config) error {
	c := common.SyncConfig{
		SourceAddr:     cfg.Source.URL,
		SourceToken:    cfg.Source.Token,
		SourceOrg:      cfg.Source.Org,
		SourceBucket:   cfg.Source.Bucket,
		TargetAddr:     cfg.Target.URL,
		TargetDBPrefix: cfg.Target.DBPrefix,
		TargetDBSuffix: cfg.Target.DBSuffix,
		TargetToken:    cfg.Target.Token,
		TargetOrg:      cfg.Target.Org,
		TargetBucket:   cfg.Target.Bucket,
		BatchSize:      cfg.Sync.BatchSize,
		Start:          cfg.Sync.Start,
		End:            cfg.Sync.End,
		ResumeFile:     cfg.Sync.ResumeFile,
		Parallel:       cfg.Sync.Parallel,
		RetryCount:     cfg.Sync.RetryCount,
		RetryInterval:  cfg.Sync.RetryInterval,
		RateLimit:      cfg.Sync.RateLimit,
		LogLevel:       cfg.Log.Level,
	}
	return influxdb2.Sync2x2x(context.Background(), c)
}
