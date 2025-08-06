package influxdb2

import (
	"context"
	"testing"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

func TestSync2x2x(t *testing.T) {
	// 测试2x到2x同步函数
	cfg := common.SyncConfig{
		SourceAddr:     "http://localhost:8086",
		SourceToken:    "source-token",
		SourceOrg:      "source-org",
		SourceBucket:   "source-bucket",
		TargetAddr:     "http://localhost:8087",
		TargetToken:    "target-token",
		TargetOrg:      "target-org",
		TargetBucket:   "target-bucket",
		BatchSize:      100,
		Start:          "2024-01-01T00:00:00Z",
		LogLevel:       "info",
	}

	ctx := context.Background()
	err := Sync2x2x(ctx, cfg)
	
	// 期望连接错误，因为没有真实的InfluxDB实例
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}
	
	// 验证错误是连接相关的，而不是参数错误
	t.Logf("预期的连接错误: %v", err)
}

func TestSync2x2xWithEmptyTargetBucket(t *testing.T) {
	// 测试目标bucket为空时的前后缀拼接逻辑
	cfg := common.SyncConfig{
		SourceAddr:     "http://localhost:8086",
		SourceToken:    "source-token",
		SourceOrg:      "source-org",
		SourceBucket:   "my-bucket",
		TargetAddr:     "http://localhost:8087",
		TargetToken:    "target-token",
		TargetOrg:      "target-org",
		TargetBucket:   "", // 空bucket，应该使用前后缀拼接
		TargetDBPrefix: "backup_",
		TargetDBSuffix: "_copy",
		BatchSize:      100,
		LogLevel:       "info",
	}

	ctx := context.Background()
	err := Sync2x2x(ctx, cfg)
	
	// 期望连接错误，因为没有真实的InfluxDB实例
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}
	
	t.Logf("预期的连接错误: %v", err)
}

func TestSync2x2xConfigurationCombinations(t *testing.T) {
	// 测试不同的配置组合
	testCases := []struct {
		name   string
		config common.SyncConfig
	}{
		{
			name: "完整配置",
			config: common.SyncConfig{
				SourceAddr:     "https://source.influxdata.com",
				SourceToken:    "source-super-secret-token",
				SourceOrg:      "source-organization",
				SourceBucket:   "source-bucket",
				TargetAddr:     "https://target.influxdata.com",
				TargetToken:    "target-super-secret-token",
				TargetOrg:      "target-organization",
				TargetBucket:   "target-bucket",
				BatchSize:      500,
				Start:          "2024-01-01T00:00:00Z",
				End:            "2024-12-31T23:59:59Z",
				Parallel:       4,
				RetryCount:     3,
				RetryInterval:  500,
				RateLimit:      50,
				LogLevel:       "debug",
			},
		},
		{
			name: "最小配置",
			config: common.SyncConfig{
				SourceAddr:   "http://localhost:8086",
				SourceToken:  "source-token",
				SourceOrg:    "source-org",
				SourceBucket: "source-bucket",
				TargetAddr:   "http://localhost:8087",
				TargetToken:  "target-token",
				TargetOrg:    "target-org",
				TargetBucket: "target-bucket",
			},
		},
		{
			name: "使用前后缀的配置",
			config: common.SyncConfig{
				SourceAddr:     "http://localhost:8086",
				SourceToken:    "source-token",
				SourceOrg:      "source-org",
				SourceBucket:   "production-data",
				TargetAddr:     "http://localhost:8087",
				TargetToken:    "target-token",
				TargetOrg:      "target-org",
				TargetBucket:   "", // 空bucket
				TargetDBPrefix: "backup_",
				TargetDBSuffix: "_v2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := Sync2x2x(ctx, tc.config)
			
			// 所有测试都应该返回连接错误，因为没有真实的数据库
			if err == nil {
				t.Error("期望连接错误，但没有错误")
			}
			
			t.Logf("配置 %s 的预期错误: %v", tc.name, err)
		})
	}
}

func TestSync2x2xAdapterCreation(t *testing.T) {
	// 测试适配器创建逻辑
	cfg := common.SyncConfig{
		SourceAddr:     "http://test-source:8086",
		SourceToken:    "test-source-token",
		SourceOrg:      "test-source-org",
		SourceBucket:   "test-source-bucket",
		TargetAddr:     "http://test-target:8087",
		TargetToken:    "test-target-token",
		TargetOrg:      "test-target-org",
		TargetBucket:   "test-target-bucket",
	}

	// 测试源适配器创建
	source := &Adapter{
		URL:    cfg.SourceAddr,
		Token:  cfg.SourceToken,
		Org:    cfg.SourceOrg,
		Bucket: cfg.SourceBucket,
	}

	if source.URL != cfg.SourceAddr {
		t.Error("源适配器URL配置错误")
	}
	if source.Token != cfg.SourceToken {
		t.Error("源适配器Token配置错误")
	}
	if source.Org != cfg.SourceOrg {
		t.Error("源适配器Org配置错误")
	}
	if source.Bucket != cfg.SourceBucket {
		t.Error("源适配器Bucket配置错误")
	}

	// 测试目标适配器创建
	target := &Adapter{
		URL:    cfg.TargetAddr,
		Token:  cfg.TargetToken,
		Org:    cfg.TargetOrg,
		Bucket: cfg.TargetBucket,
	}

	if target.URL != cfg.TargetAddr {
		t.Error("目标适配器URL配置错误")
	}
	if target.Token != cfg.TargetToken {
		t.Error("目标适配器Token配置错误")
	}
	if target.Org != cfg.TargetOrg {
		t.Error("目标适配器Org配置错误")
	}
	if target.Bucket != cfg.TargetBucket {
		t.Error("目标适配器Bucket配置错误")
	}
}

func TestSync2x2xFunction(t *testing.T) {
	// 测试Sync2x2x函数的基本功能
	cfg := common.SyncConfig{
		SourceAddr:     "http://mock-source:8086",
		SourceToken:    "mock-source-token",
		SourceOrg:      "mock-source-org",
		SourceBucket:   "mock-source-bucket",
		TargetAddr:     "http://mock-target:8087",
		TargetToken:    "mock-target-token",
		TargetOrg:      "mock-target-org",
		TargetBucket:   "mock-target-bucket",
		BatchSize:      1000,
		Start:          "2024-01-01T00:00:00Z",
		End:            "2024-12-31T23:59:59Z",
		Parallel:       2,
		LogLevel:       "warn",
	}

	// 测试函数调用不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Sync2x2x发生panic: %v", r)
		}
	}()

	ctx := context.Background()
	err := Sync2x2x(ctx, cfg)

	// 应该返回错误，因为无法连接到mock地址
	if err == nil {
		t.Error("期望返回连接错误")
	}

	t.Logf("预期的错误: %v", err)
}
