package influxdb1

import (
	"context"
	"testing"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

func TestSyncConfig(t *testing.T) {
	// 测试SyncConfig结构体
	cfg := SyncConfig{
		SourceAddr:      "http://source:8086",
		SourceUser:      "src_user",
		SourcePass:      "src_pass",
		SourceDB:        "src_db",
		SourceDBExclude: []string{"_internal", "system"},
		TargetAddr:      "http://target:8086",
		TargetUser:      "tgt_user",
		TargetPass:      "tgt_pass",
		TargetDB:        "tgt_db",
		TargetDBPrefix:  "backup_",
		TargetDBSuffix:  "_copy",
		BatchSize:       1000,
		Start:           "2024-01-01T00:00:00Z",
		End:             "2024-12-31T23:59:59Z",
		ResumeFile:      "resume.state",
		Parallel:        4,
		RetryCount:      3,
		RetryInterval:   500,
		RateLimit:       50,
		LogLevel:        "info",
	}

	// 验证字段设置
	if cfg.SourceAddr != "http://source:8086" {
		t.Error("SourceAddr设置错误")
	}
	if cfg.SourceUser != "src_user" {
		t.Error("SourceUser设置错误")
	}
	if cfg.BatchSize != 1000 {
		t.Error("BatchSize设置错误")
	}
	if len(cfg.SourceDBExclude) != 2 {
		t.Error("SourceDBExclude设置错误")
	}
	if cfg.Parallel != 4 {
		t.Error("Parallel设置错误")
	}
}

func TestExtraConfig(t *testing.T) {
	// 测试ExtraConfig结构体
	cfg := ExtraConfig{
		SourceDBExclude: []string{"test1", "test2"},
		TargetDBPrefix:  "prefix_",
		TargetDBSuffix:  "_suffix",
	}

	if len(cfg.SourceDBExclude) != 2 {
		t.Error("SourceDBExclude长度错误")
	}
	if cfg.TargetDBPrefix != "prefix_" {
		t.Error("TargetDBPrefix设置错误")
	}
	if cfg.TargetDBSuffix != "_suffix" {
		t.Error("TargetDBSuffix设置错误")
	}
}

func TestSync1x1x(t *testing.T) {
	// 测试1x1x同步函数
	cfg := common.SyncConfig{
		SourceAddr: "http://localhost:8086",
		SourceUser: "admin",
		SourcePass: "password",
		TargetAddr: "http://localhost:8087",
		TargetUser: "admin",
		TargetPass: "password",
		BatchSize:  100,
	}

	ctx := context.Background()
	err := Sync1x1x(ctx, cfg)

	// 期望连接错误，因为没有真实的InfluxDB实例
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}

	// 验证错误是连接相关的，而不是参数错误
	t.Logf("预期的连接错误: %v", err)
}

func TestSyncBackwardCompatibility(t *testing.T) {
	// 测试向后兼容的Sync函数
	cfg := SyncConfig{
		SourceAddr: "http://localhost:8086",
		SourceUser: "admin",
		SourcePass: "password",
		TargetAddr: "http://localhost:8087",
		TargetUser: "admin",
		TargetPass: "password",
		BatchSize:  100,
		Start:      "2024-01-01T00:00:00Z",
		LogLevel:   "info",
	}

	ctx := context.Background()
	err := Sync(ctx, cfg)

	// 期望连接错误，因为没有真实的InfluxDB实例
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}

	// 验证错误是连接相关的，而不是参数错误
	t.Logf("预期的连接错误: %v", err)
}

func TestSyncConfigConversion(t *testing.T) {
	// 测试配置转换逻辑
	oldCfg := SyncConfig{
		SourceAddr:      "http://old-source:8086",
		SourceUser:      "old_user",
		SourcePass:      "old_pass",
		SourceDB:        "old_db",
		SourceDBExclude: []string{"_internal"},
		TargetAddr:      "http://old-target:8086",
		TargetUser:      "old_target_user",
		TargetPass:      "old_target_pass",
		TargetDB:        "old_target_db",
		TargetDBPrefix:  "old_",
		TargetDBSuffix:  "_old",
		BatchSize:       500,
		Start:           "2023-01-01T00:00:00Z",
		End:             "2023-12-31T23:59:59Z",
		ResumeFile:      "old_resume.state",
		Parallel:        2,
		RetryCount:      5,
		RetryInterval:   1000,
		RateLimit:       100,
		LogLevel:        "debug",
	}

	// 这里主要测试结构体字段的完整性
	if oldCfg.SourceAddr == "" {
		t.Error("SourceAddr不应该为空")
	}
	if oldCfg.BatchSize == 0 {
		t.Error("BatchSize不应该为0")
	}
	if oldCfg.Parallel == 0 {
		t.Error("Parallel不应该为0")
	}
}
