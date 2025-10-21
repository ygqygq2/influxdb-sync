package influxdb3

import (
	"context"
	"testing"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

func TestSync1x3x(t *testing.T) {
	// 需要实际的数据库实例
	t.Skip("需要实际的 InfluxDB 实例")

	cfg := common.SyncConfig{
		SourceAddr:   "http://localhost:8086",
		SourceUser:   "admin",
		SourcePass:   "password",
		SourceDB:     "testdb",
		TargetAddr:   "http://localhost:8087",
		TargetToken:  "test-token",
		TargetOrg:    "test-org",
		TargetBucket: "test-bucket",
		BatchSize:    1000,
	}

	err := Sync1x3x(context.Background(), cfg)
	if err != nil {
		t.Logf("预期错误（无实际数据库）: %v", err)
	}
}

func TestSync2x3x(t *testing.T) {
	// 需要实际的数据库实例
	t.Skip("需要实际的 InfluxDB 实例")

	cfg := common.SyncConfig{
		SourceAddr:   "http://localhost:8086",
		SourceToken:  "source-token",
		SourceOrg:    "source-org",
		SourceBucket: "source-bucket",
		TargetAddr:   "http://localhost:8087",
		TargetToken:  "target-token",
		TargetOrg:    "target-org",
		TargetBucket: "target-bucket",
		BatchSize:    1000,
	}

	err := Sync2x3x(context.Background(), cfg)
	if err != nil {
		t.Logf("预期错误（无实际数据库）: %v", err)
	}
}

func TestSync3x3x(t *testing.T) {
	// 需要实际的数据库实例
	t.Skip("需要实际的 InfluxDB 实例")

	cfg := common.SyncConfig{
		SourceAddr:   "http://localhost:8086",
		SourceUser:   "admin",
		SourcePass:   "password",
		SourceDB:     "testdb",
		TargetAddr:   "http://localhost:8087",
		TargetToken:  "target-token",
		TargetOrg:    "target-org",
		TargetBucket: "target-bucket",
		BatchSize:    1000,
	}

	err := Sync3x3x(context.Background(), cfg)
	if err != nil {
		t.Logf("预期错误（无实际数据库）: %v", err)
	}
}

func TestSync1x3x_InvalidConfig(t *testing.T) {
	cfg := common.SyncConfig{
		SourceAddr: "",
		TargetAddr: "",
	}

	err := Sync1x3x(context.Background(), cfg)
	if err == nil {
		t.Error("应该返回配置错误")
	}
}

func TestSync2x3x_InvalidConfig(t *testing.T) {
	cfg := common.SyncConfig{
		SourceAddr: "",
		TargetAddr: "",
	}

	err := Sync2x3x(context.Background(), cfg)
	if err == nil {
		t.Error("应该返回配置错误")
	}
}

func TestSync3x3x_InvalidConfig(t *testing.T) {
	cfg := common.SyncConfig{
		SourceAddr: "",
		TargetAddr: "",
	}

	err := Sync3x3x(context.Background(), cfg)
	if err == nil {
		t.Error("应该返回配置错误")
	}
}
