package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ygqygq2/influxdb-sync/internal/config"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb1"
)

func RunInfluxdb1Sync(cfgPath string) error {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}
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
	ResumeFile: cfg.Sync.ResumeFile,
}
	return influxdb1.Sync(context.Background(), c)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: influxdb-sync <config.yaml>")
		os.Exit(1)
	}
	if err := RunInfluxdb1Sync(os.Args[1]); err != nil {
		fmt.Println("同步失败:", err)
		os.Exit(2)
	}
	fmt.Println("同步完成")
}
