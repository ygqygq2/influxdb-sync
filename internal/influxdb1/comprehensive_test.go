package influxdb1

import (
	"context"
	"testing"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

func TestSync1x2xAdditional(t *testing.T) {
	// 测试1x到2x同步的更多场景
	cfg := common.SyncConfig{
		SourceAddr:      "http://localhost:8086",
		SourceUser:      "admin",
		SourcePass:      "password",
		SourceDB:        "testdb",
		SourceDBExclude: []string{"_internal", "system"},
		TargetAddr:      "http://localhost:8087",
		TargetToken:     "target-token",
		TargetOrg:       "target-org",
		TargetBucket:    "target-bucket",
		TargetDBPrefix:  "backup_",
		TargetDBSuffix:  "_v2",
		BatchSize:       1000,
		Start:           "2024-01-01T00:00:00Z",
		End:             "2024-12-31T23:59:59Z",
		ResumeFile:      "test_resume.state",
		Parallel:        4,
		RetryCount:      3,
		RetryInterval:   500,
		RateLimit:       50,
		LogLevel:        "debug",
	}

	ctx := context.Background()
	err := Sync1x2x(ctx, cfg)

	// 期望连接错误，因为没有真实的InfluxDB实例
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}

	t.Logf("预期的连接错误: %v", err)
}

func TestInfluxDB1xSpecificFunctionality(t *testing.T) {
	// 测试 InfluxDB 1.x 特有的功能
	t.Run("数据源和目标配置", func(t *testing.T) {
		sourceConfig := DataSourceConfig{
			Addr: "http://localhost:8086",
			User: "admin", 
			Pass: "password",
		}

		targetConfig := DataTargetConfig{
			Addr: "http://localhost:8087",
			User: "admin",
			Pass: "password", 
		}

		source := NewDataSource(sourceConfig)
		target := NewDataTarget(targetConfig)

		if source == nil {
			t.Error("数据源创建失败")
		}
		if target == nil {
			t.Error("数据目标创建失败")
		}
	})

	t.Run("ExtraConfig配置", func(t *testing.T) {
		extraCfg := ExtraConfig{
			SourceDBExclude: []string{"_internal", "system"},
			TargetDBPrefix:  "backup_",
			TargetDBSuffix:  "_archived",
		}

		if len(extraCfg.SourceDBExclude) != 2 {
			t.Error("SourceDBExclude设置失败")
		}
		if extraCfg.TargetDBPrefix != "backup_" {
			t.Error("TargetDBPrefix设置失败")
		}
		if extraCfg.TargetDBSuffix != "_archived" {
			t.Error("TargetDBSuffix设置失败")
		}
	})

	t.Run("Measurement转义功能", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"cpu", `"cpu"`},
			{"cpu usage", `"cpu usage"`},
			{"cpu-memory_disk", `"cpu-memory_disk"`},
			{`cpu"memory`, `"cpu\"memory"`},
			{"1cpu", `"1cpu"`},
			{"处理器", `"处理器"`},
		}

		for _, tc := range testCases {
			result := escapeMeasurement(tc.input)
			if result != tc.expected {
				t.Errorf("escapeMeasurement(%q) = %q, 期望 %q",
					tc.input, result, tc.expected)
			}
		}
	})
}
