package influxdb1

import (
	"context"
	"log"
)

type SyncConfig struct {
	SourceAddr string
	SourceUser string
	SourcePass string
	SourceDB   string
	TargetAddr string
	TargetUser string
	TargetPass string
	TargetDB   string
	BatchSize  int
	ResumeFile string
}

func Sync(ctx context.Context, cfg SyncConfig) error {
	src, err := NewClient(cfg.SourceAddr, cfg.SourceUser, cfg.SourcePass, 0)
	if err != nil {
		return err
	}
	defer src.Close()
	tgt, err := NewClient(cfg.TargetAddr, cfg.TargetUser, cfg.TargetPass, 0)
	if err != nil {
		return err
	}
	defer tgt.Close()

	// TODO: 查询源库数据，批量写入目标库，实现断点续传
	log.Println("[TODO] influxdb1.8 同步逻辑待实现")
	return nil
}
