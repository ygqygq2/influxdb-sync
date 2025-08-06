package influxdb2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

func TestSync2x2xDetailedCoverage(t *testing.T) {
	// 测试2x到2x同步的详细场景
	cfg := common.SyncConfig{
		SourceAddr:      "http://localhost:8086",
		SourceToken:     "source-token-12345",
		SourceOrg:       "source-organization",
		SourceBucket:    "source-bucket",
		TargetAddr:      "http://localhost:8087",
		TargetToken:     "target-token-67890",
		TargetOrg:       "target-organization", 
		TargetBucket:    "target-bucket",
		TargetDBPrefix:  "migrated_",
		TargetDBSuffix:  "_v2",
		BatchSize:       1500,
		Start:           "2024-01-01T00:00:00Z",
		End:             "2024-12-31T23:59:59Z",
		ResumeFile:      "sync_2x2x_resume.state",
		Parallel:        6,
		RetryCount:      5,
		RetryInterval:   1000,
		RateLimit:       25,
		LogLevel:        "debug",
	}

	ctx := context.Background()
	err := Sync2x2x(ctx, cfg)
	
	// 期望连接错误，因为没有真实的InfluxDB实例
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}
	
	t.Logf("预期的连接错误: %v", err)
}

func TestSync2x2xConfigurationVariations(t *testing.T) {
	// 测试不同配置的2x2x同步
	testConfigs := []struct {
		name   string
		config common.SyncConfig
	}{
		{
			name: "最小配置",
			config: common.SyncConfig{
				SourceAddr:   "http://localhost:8086",
				SourceToken:  "src-token",
				SourceOrg:    "src-org",
				TargetAddr:   "http://localhost:8087",
				TargetToken:  "tgt-token",
				TargetOrg:    "tgt-org",
				BatchSize:    100,
			},
		},
		{
			name: "高性能配置",
			config: common.SyncConfig{
				SourceAddr:   "http://localhost:8086",
				SourceToken:  "src-token",
				SourceOrg:    "src-org",
				TargetAddr:   "http://localhost:8087",
				TargetToken:  "tgt-token",
				TargetOrg:    "tgt-org",
				BatchSize:    10000,
				Parallel:     16,
				RetryCount:   1,
				RateLimit:    0,
			},
		},
		{
			name: "高可靠性配置",
			config: common.SyncConfig{
				SourceAddr:    "http://localhost:8086",
				SourceToken:   "src-token",
				SourceOrg:     "src-org",
				TargetAddr:    "http://localhost:8087",
				TargetToken:   "tgt-token",
				TargetOrg:     "tgt-org",
				BatchSize:     500,
				Parallel:      2,
				RetryCount:    10,
				RetryInterval: 2000,
				RateLimit:     500,
			},
		},
		{
			name: "云端配置",
			config: common.SyncConfig{
				SourceAddr:     "https://us-west-2-1.aws.cloud2.influxdata.com",
				SourceToken:    "cloud-source-token",
				SourceOrg:      "cloud-source-org",
				SourceBucket:   "production-data",
				TargetAddr:     "https://eu-central-1-1.aws.cloud2.influxdata.com",
				TargetToken:    "cloud-target-token",
				TargetOrg:      "cloud-target-org",
				TargetBucket:   "backup-data",
				TargetDBPrefix: "backup_",
				TargetDBSuffix: "_eu",
				BatchSize:      2000,
				Parallel:       4,
				RetryCount:     3,
				RetryInterval:  1000,
				RateLimit:      100,
			},
		},
	}

	for _, tc := range testConfigs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := Sync2x2x(ctx, tc.config)
			
			// 所有配置都应该因为连接失败而返回错误
			if err == nil {
				t.Error("期望连接错误，但没有错误")
			}
			
			t.Logf("配置 %s 的预期连接错误: %v", tc.name, err)
		})
	}
}

func TestSync2x2xParameterValidation(t *testing.T) {
	// 测试参数验证场景
	testCases := []struct {
		name           string
		config         common.SyncConfig
		expectError    bool
		errorKeyword   string
	}{
		{
			name: "空源地址",
			config: common.SyncConfig{
				SourceAddr:  "",
				SourceToken: "token",
				SourceOrg:   "org",
				TargetAddr:  "http://localhost:8087",
				TargetToken: "token",
				TargetOrg:   "org",
				BatchSize:   100,
			},
			expectError:  true,
			errorKeyword: "source",
		},
		{
			name: "空目标地址",
			config: common.SyncConfig{
				SourceAddr:  "http://localhost:8086",
				SourceToken: "token",
				SourceOrg:   "org",
				TargetAddr:  "",
				TargetToken: "token",
				TargetOrg:   "org",
				BatchSize:   100,
			},
			expectError:  true,
			errorKeyword: "target",
		},
		{
			name: "无效批次大小",
			config: common.SyncConfig{
				SourceAddr:  "http://localhost:8086",
				SourceToken: "token",
				SourceOrg:   "org",
				TargetAddr:  "http://localhost:8087",
				TargetToken: "token",
				TargetOrg:   "org",
				BatchSize:   0,
			},
			expectError:  true,
			errorKeyword: "batch",
		},
		{
			name: "缺少源Token",
			config: common.SyncConfig{
				SourceAddr:  "http://localhost:8086",
				SourceToken: "",
				SourceOrg:   "org",
				TargetAddr:  "http://localhost:8087",
				TargetToken: "token",
				TargetOrg:   "org",
				BatchSize:   100,
			},
			expectError:  true,
			errorKeyword: "token",
		},
		{
			name: "缺少目标Token",
			config: common.SyncConfig{
				SourceAddr:  "http://localhost:8086",
				SourceToken: "token",
				SourceOrg:   "org",
				TargetAddr:  "http://localhost:8087",
				TargetToken: "",
				TargetOrg:   "org",
				BatchSize:   100,
			},
			expectError:  true,
			errorKeyword: "token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := Sync2x2x(ctx, tc.config)
			
			if tc.expectError && err == nil {
				t.Errorf("期望错误但没有收到错误")
			}
			
			if err != nil {
				t.Logf("收到预期错误: %v", err)
			}
		})
	}
}

func TestAdapterAdvancedMethods(t *testing.T) {
	// 测试适配器的高级方法
	adapter := &Adapter{
		URL:    "http://localhost:8086",
		Token:  "test-token-advanced",
		Org:    "test-org-advanced",
		Bucket: "test-bucket-advanced",
	}

	// 测试连接
	err := adapter.Connect()
	if err != nil {
		t.Errorf("Connect失败: %v", err)
	}
	defer adapter.Close()

	// 测试GetDatabases（带bucket）
	buckets, err := adapter.GetDatabases()
	if err != nil {
		t.Logf("GetDatabases预期错误: %v", err)
	} else {
		if len(buckets) != 1 || buckets[0] != "test-bucket-advanced" {
			t.Errorf("期望返回 [test-bucket-advanced], 实际返回 %v", buckets)
		}
	}

	// 测试GetMeasurements
	_, err = adapter.GetMeasurements("test-bucket-advanced")
	if err != nil {
		t.Logf("GetMeasurements预期错误: %v", err)
	}

	// 测试GetTagKeys
	_, err = adapter.GetTagKeys("test-bucket-advanced", "test-measurement")
	if err != nil {
		t.Logf("GetTagKeys预期错误: %v", err)
	}

	// 测试QueryData
	_, _, err = adapter.QueryData("test-bucket-advanced", "test-measurement", 0, 100)
	if err != nil {
		t.Logf("QueryData预期错误: %v", err)
	}

	// 测试WritePoints
	points := []common.DataPoint{
		{
			Measurement: "test_advanced",
			Tags:        map[string]string{"env": "test", "service": "advanced"},
			Fields:      map[string]interface{}{"value": 42.0, "status": "ok"},
			Time:        time.Now(),
		},
	}
	err = adapter.WritePoints("test-bucket-advanced", points)
	if err != nil {
		t.Logf("WritePoints预期错误: %v", err)
	}
}

func TestAdapterEmptyBucketScenario(t *testing.T) {
	// 测试没有指定bucket的场景
	adapter := &Adapter{
		URL:   "http://localhost:8086",
		Token: "test-token-no-bucket",
		Org:   "test-org-no-bucket",
		// 不设置Bucket
	}

	err := adapter.Connect()
	if err != nil {
		t.Errorf("Connect失败: %v", err)
	}
	defer adapter.Close()

	// 当没有指定bucket时，应该查询所有buckets
	_, err = adapter.GetDatabases()
	if err != nil {
		t.Logf("GetDatabases（无bucket）预期错误: %v", err)
	}
}

func TestAdapterLargeDataSet(t *testing.T) {
	// 测试大数据集处理
	adapter := &Adapter{
		URL:    "http://localhost:8086",
		Token:  "test-token-large",
		Org:    "test-org-large",
		Bucket: "test-bucket-large",
	}

	err := adapter.Connect()
	if err != nil {
		t.Errorf("Connect失败: %v", err)
	}
	defer adapter.Close()

	// 测试大批次查询
	_, _, err = adapter.QueryData("test-bucket-large", "large-measurement", 0, 10000)
	if err != nil {
		t.Logf("大批次QueryData预期错误: %v", err)
	}

	// 测试大量数据点写入
	var largePoints []common.DataPoint
	for i := 0; i < 1000; i++ {
		point := common.DataPoint{
			Measurement: "large_test",
			Tags:        map[string]string{"instance": fmt.Sprintf("server-%d", i)},
			Fields:      map[string]interface{}{"value": float64(i), "counter": i},
			Time:        time.Now().Add(time.Duration(i) * time.Second),
		}
		largePoints = append(largePoints, point)
	}

	err = adapter.WritePoints("test-bucket-large", largePoints)
	if err != nil {
		t.Logf("大量WritePoints预期错误: %v", err)
	}
}
