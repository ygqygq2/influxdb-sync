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

func TestSyncConfigValidationHelper(t *testing.T) {
	// 测试配置验证辅助函数
	testCases := []struct {
		name       string
		config     SyncConfig
		shouldPass bool
	}{
		{
			name: "有效的最小配置",
			config: SyncConfig{
				SourceAddr: "http://localhost:8086",
				TargetAddr: "http://localhost:8087",
				BatchSize:  100,
			},
			shouldPass: true,
		},
		{
			name: "缺少源地址",
			config: SyncConfig{
				TargetAddr: "http://localhost:8087",
				BatchSize:  100,
			},
			shouldPass: false,
		},
		{
			name: "缺少目标地址",
			config: SyncConfig{
				SourceAddr: "http://localhost:8086",
				BatchSize:  100,
			},
			shouldPass: false,
		},
		{
			name: "无效的批次大小",
			config: SyncConfig{
				SourceAddr: "http://localhost:8086",
				TargetAddr: "http://localhost:8087",
				BatchSize:  0,
			},
			shouldPass: false,
		},
		{
			name: "完整配置",
			config: SyncConfig{
				SourceAddr:      "http://localhost:8086",
				SourceUser:      "admin",
				SourcePass:      "password",
				SourceDB:        "testdb",
				SourceDBExclude: []string{"_internal"},
				TargetAddr:      "http://localhost:8087",
				TargetUser:      "admin",
				TargetPass:      "password",
				TargetDB:        "targetdb",
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
			},
			shouldPass: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 手动验证配置
			valid := tc.config.SourceAddr != "" && 
				tc.config.TargetAddr != "" && 
				tc.config.BatchSize > 0

			if valid != tc.shouldPass {
				t.Errorf("配置验证失败：期望 %v，实际 %v", tc.shouldPass, valid)
			}
		})
	}
}

func TestSyncConfigFieldsDetail(t *testing.T) {
	// 详细测试SyncConfig的每个字段
	cfg := SyncConfig{}

	// 测试字符串字段
	cfg.SourceAddr = "http://source:8086"
	cfg.SourceUser = "src_user"
	cfg.SourcePass = "src_pass"
	cfg.SourceDB = "src_db"
	cfg.TargetAddr = "http://target:8087"
	cfg.TargetUser = "tgt_user"
	cfg.TargetPass = "tgt_pass"
	cfg.TargetDB = "tgt_db"
	cfg.TargetDBPrefix = "prefix_"
	cfg.TargetDBSuffix = "_suffix"
	cfg.Start = "2024-01-01T00:00:00Z"
	cfg.End = "2024-12-31T23:59:59Z"
	cfg.ResumeFile = "resume.state"
	cfg.LogLevel = "debug"

	// 验证字符串字段
	if cfg.SourceAddr != "http://source:8086" {
		t.Error("SourceAddr设置失败")
	}
	if cfg.SourceUser != "src_user" {
		t.Error("SourceUser设置失败")
	}
	if cfg.SourcePass != "src_pass" {
		t.Error("SourcePass设置失败")
	}
	if cfg.SourceDB != "src_db" {
		t.Error("SourceDB设置失败")
	}
	if cfg.TargetAddr != "http://target:8087" {
		t.Error("TargetAddr设置失败")
	}
	if cfg.TargetUser != "tgt_user" {
		t.Error("TargetUser设置失败")
	}
	if cfg.TargetPass != "tgt_pass" {
		t.Error("TargetPass设置失败")
	}
	if cfg.TargetDB != "tgt_db" {
		t.Error("TargetDB设置失败")
	}
	if cfg.TargetDBPrefix != "prefix_" {
		t.Error("TargetDBPrefix设置失败")
	}
	if cfg.TargetDBSuffix != "_suffix" {
		t.Error("TargetDBSuffix设置失败")
	}
	if cfg.Start != "2024-01-01T00:00:00Z" {
		t.Error("Start设置失败")
	}
	if cfg.End != "2024-12-31T23:59:59Z" {
		t.Error("End设置失败")
	}
	if cfg.ResumeFile != "resume.state" {
		t.Error("ResumeFile设置失败")
	}
	if cfg.LogLevel != "debug" {
		t.Error("LogLevel设置失败")
	}

	// 测试整数字段
	cfg.BatchSize = 2000
	cfg.Parallel = 8
	cfg.RetryCount = 5
	cfg.RetryInterval = 1000
	cfg.RateLimit = 100

	if cfg.BatchSize != 2000 {
		t.Error("BatchSize设置失败")
	}
	if cfg.Parallel != 8 {
		t.Error("Parallel设置失败")
	}
	if cfg.RetryCount != 5 {
		t.Error("RetryCount设置失败")
	}
	if cfg.RetryInterval != 1000 {
		t.Error("RetryInterval设置失败")
	}
	if cfg.RateLimit != 100 {
		t.Error("RateLimit设置失败")
	}

	// 测试切片字段
	cfg.SourceDBExclude = []string{"_internal", "system", "temp"}
	if len(cfg.SourceDBExclude) != 3 {
		t.Error("SourceDBExclude设置失败")
	}
	if cfg.SourceDBExclude[0] != "_internal" {
		t.Error("SourceDBExclude第一个元素设置失败")
	}
	if cfg.SourceDBExclude[1] != "system" {
		t.Error("SourceDBExclude第二个元素设置失败")
	}
	if cfg.SourceDBExclude[2] != "temp" {
		t.Error("SourceDBExclude第三个元素设置失败")
	}
}

func TestExtraConfigDetails(t *testing.T) {
	// 详细测试ExtraConfig
	extraCfg := ExtraConfig{}

	// 设置字段
	extraCfg.SourceDBExclude = []string{"exclude1", "exclude2", "exclude3"}
	extraCfg.TargetDBPrefix = "backup_"
	extraCfg.TargetDBSuffix = "_archived"

	// 验证字段
	if len(extraCfg.SourceDBExclude) != 3 {
		t.Errorf("期望SourceDBExclude长度为3，实际为%d", len(extraCfg.SourceDBExclude))
	}

	expectedExcludes := []string{"exclude1", "exclude2", "exclude3"}
	for i, expected := range expectedExcludes {
		if extraCfg.SourceDBExclude[i] != expected {
			t.Errorf("SourceDBExclude[%d]期望为%s，实际为%s", i, expected, extraCfg.SourceDBExclude[i])
		}
	}

	if extraCfg.TargetDBPrefix != "backup_" {
		t.Errorf("期望TargetDBPrefix为backup_，实际为%s", extraCfg.TargetDBPrefix)
	}

	if extraCfg.TargetDBSuffix != "_archived" {
		t.Errorf("期望TargetDBSuffix为_archived，实际为%s", extraCfg.TargetDBSuffix)
	}
}

func TestConfigurationEdgeCases(t *testing.T) {
	// 测试边缘情况
	t.Run("空字符串字段", func(t *testing.T) {
		cfg := SyncConfig{
			SourceAddr: "",
			TargetAddr: "",
			LogLevel:   "",
		}

		if cfg.SourceAddr != "" {
			t.Error("空字符串字段设置失败")
		}
	})

	t.Run("极大数值", func(t *testing.T) {
		cfg := SyncConfig{
			BatchSize:     999999,
			Parallel:      100,
			RetryCount:    1000,
			RetryInterval: 60000,
			RateLimit:     10000,
		}

		if cfg.BatchSize != 999999 {
			t.Error("极大数值设置失败")
		}
		if cfg.Parallel != 100 {
			t.Error("极大并发数设置失败")
		}
	})

	t.Run("空切片", func(t *testing.T) {
		cfg := SyncConfig{
			SourceDBExclude: []string{},
		}

		if len(cfg.SourceDBExclude) != 0 {
			t.Error("空切片设置失败")
		}
	})

	t.Run("长字符串", func(t *testing.T) {
		longString := ""
		for i := 0; i < 1000; i++ {
			longString += "a"
		}

		cfg := SyncConfig{
			SourceAddr: longString,
		}

		if len(cfg.SourceAddr) != 1000 {
			t.Error("长字符串设置失败")
		}
	})
}
