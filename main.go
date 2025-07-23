package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ygqygq2/influxdb-sync/internal/config"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb1"
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
	// 根据 type 字段自动判断同步模式
	if cfg.Source.Type == 1 && cfg.Target.Type == 1 {
		if err := runSync1x1x(cfg); err != nil {
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
		SourceAddr: cfg.Source.URL,
		SourceUser: cfg.Source.User,
		SourcePass: cfg.Source.Pass,
		SourceDB:   cfg.Source.DB,
		TargetAddr: cfg.Target.URL,
		TargetUser: cfg.Target.User,
		TargetPass: cfg.Target.Pass,
		TargetDB:   cfg.Target.DB,
		BatchSize:  cfg.Sync.BatchSize,
	   Start:      cfg.Sync.Start,
	   End:        cfg.Sync.End,
	   ResumeFile: cfg.Sync.ResumeFile,
	}
	return influxdb1.Sync(context.Background(), c)
}
