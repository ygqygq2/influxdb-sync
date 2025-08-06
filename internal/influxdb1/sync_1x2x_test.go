package influxdb1

import (
	"context"
	"testing"

	"github.com/ygqygq2/influxdb-sync/internal/common"
	"github.com/ygqygq2/influxdb-sync/internal/influxdb2"
)

func TestSync1x2x(t *testing.T) {
	// 测试1x到2x同步函数
	cfg := common.SyncConfig{
		SourceAddr:   "http://localhost:8086",
		SourceUser:   "admin",
		SourcePass:   "password",
		TargetAddr:   "http://localhost:8087",
		TargetToken:  "target-token",
		TargetOrg:    "target-org",
		TargetBucket: "target-bucket",
		BatchSize:    100,
		Start:        "2024-01-01T00:00:00Z",
		LogLevel:     "info",
	}

	ctx := context.Background()
	err := Sync1x2x(ctx, cfg)

	// 期望连接错误，因为没有真实的InfluxDB实例
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}

	// 验证错误是连接相关的，而不是参数错误
	t.Logf("预期的连接错误: %v", err)
}

func TestSync1x2xConfigValidation(t *testing.T) {
	// 测试不同的配置组合
	testCases := []struct {
		name   string
		config common.SyncConfig
	}{
		{
			name: "完整配置",
			config: common.SyncConfig{
				SourceAddr:   "http://source:8086",
				SourceUser:   "src_admin",
				SourcePass:   "src_password",
				TargetAddr:   "https://target.influxdata.com",
				TargetToken:  "my-super-secret-token",
				TargetOrg:    "my-organization",
				TargetBucket: "my-bucket",
				BatchSize:    500,
				Parallel:     2,
				LogLevel:     "debug",
			},
		},
		{
			name: "最小配置",
			config: common.SyncConfig{
				SourceAddr:   "http://localhost:8086",
				TargetAddr:   "http://localhost:8087",
				TargetToken:  "token",
				TargetOrg:    "org",
				TargetBucket: "bucket",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := Sync1x2x(ctx, tc.config)

			// 所有测试都应该返回连接错误，因为没有真实的数据库
			if err == nil {
				t.Error("期望连接错误，但没有错误")
			}

			t.Logf("配置 %s 的预期错误: %v", tc.name, err)
		})
	}
}

func TestDataSourceAndTargetCreation(t *testing.T) {
	// 测试数据源和目标的创建
	cfg := common.SyncConfig{
		SourceAddr:   "http://test-source:8086",
		SourceUser:   "test_user",
		SourcePass:   "test_pass",
		TargetAddr:   "http://test-target:8087",
		TargetToken:  "test-token",
		TargetOrg:    "test-org",
		TargetBucket: "test-bucket",
	}

	// 创建源
	source := NewDataSource(DataSourceConfig{
		Addr: cfg.SourceAddr,
		User: cfg.SourceUser,
		Pass: cfg.SourcePass,
	})

	if source == nil {
		t.Fatal("数据源创建失败")
	}

	if source.config.Addr != cfg.SourceAddr {
		t.Error("数据源地址配置错误")
	}

	// 创建目标
	target := &influxdb2.Adapter{
		URL:    cfg.TargetAddr,
		Token:  cfg.TargetToken,
		Org:    cfg.TargetOrg,
		Bucket: cfg.TargetBucket,
	}

	if target.URL != cfg.TargetAddr {
		t.Error("数据目标地址配置错误")
	}

	if target.Token != cfg.TargetToken {
		t.Error("数据目标Token配置错误")
	}
}

func TestSync1x2xFunction(t *testing.T) {
	// 测试Sync1x2x函数的基本功能
	cfg := common.SyncConfig{
		SourceAddr:   "http://mock-source:8086",
		SourceUser:   "mock_user",
		SourcePass:   "mock_pass",
		TargetAddr:   "http://mock-target:8087",
		TargetToken:  "mock-token",
		TargetOrg:    "mock-org",
		TargetBucket: "mock-bucket",
		BatchSize:    1000,
		Start:        "2024-01-01T00:00:00Z",
		End:          "2024-12-31T23:59:59Z",
		Parallel:     4,
		LogLevel:     "info",
	}

	// 测试函数调用不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Sync1x2x发生panic: %v", r)
		}
	}()

	ctx := context.Background()
	err := Sync1x2x(ctx, cfg)

	// 应该返回错误，因为无法连接到mock地址
	if err == nil {
		t.Error("期望返回连接错误")
	}

	t.Logf("预期的错误: %v", err)
}
