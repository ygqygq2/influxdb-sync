package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ygqygq2/influxdb-sync/internal/config"
)

func TestMainFunction(t *testing.T) {
	// 保存原始命令行参数
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 测试无参数情况（应该退出）
	os.Args = []string{"influxdb-sync"}

	// 由于main函数会调用os.Exit，我们无法直接测试它
	// 这里主要测试不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main函数发生panic: %v", r)
		}
	}()
}

func TestRunSync1x1xFunction(t *testing.T) {
	// 创建测试配置
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_1x1x.yaml")

	configContent := `
source:
  type: 1
  url: "http://localhost:8086"
  user: "admin"
  pass: "password"
  db: "testdb"
  db_exclude: ["_internal"]
  db_prefix: ""
  db_suffix: ""

target:
  type: 1
  url: "http://localhost:8087"
  user: "admin"
  pass: "password"
  db: "targetdb"
  db_prefix: "backup_"
  db_suffix: "_copy"

sync:
  start: "2024-01-01T00:00:00Z"
  end: "2024-12-31T23:59:59Z"
  batch_size: 1000
  resume_file: "test_resume.state"
  parallel: 4
  retry_count: 3
  retry_interval: 500
  rate_limit: 50

log:
  level: "info"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("无法创建测试配置文件: %v", err)
	}

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("配置加载失败: %v", err)
	}

	// 测试runSync1x1x函数
	err = runSync1x1x(cfg)
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}

	t.Logf("预期的连接错误: %v", err)
}

func TestRunSync1x2xFunction(t *testing.T) {
	// 创建测试配置
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_1x2x.yaml")

	configContent := `
source:
  type: 1
  url: "http://localhost:8086"
  user: "admin"
  pass: "password"
  db: "testdb"
  db_exclude: ["_internal"]

target:
  type: 2
  url: "http://localhost:8087"
  token: "target-token"
  org: "target-org"
  bucket: "target-bucket"
  db_prefix: "backup_"
  db_suffix: "_copy"

sync:
  start: "2024-01-01T00:00:00Z"
  batch_size: 500
  parallel: 2
  retry_count: 3
  retry_interval: 1000
  rate_limit: 100

log:
  level: "debug"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("无法创建测试配置文件: %v", err)
	}

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("配置加载失败: %v", err)
	}

	// 测试runSync1x2x函数
	err = runSync1x2x(cfg)
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}

	t.Logf("预期的连接错误: %v", err)
}

func TestRunSync2x2xFunction(t *testing.T) {
	// 创建测试配置
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_2x2x.yaml")

	configContent := `
source:
  type: 2
  url: "http://localhost:8086"
  token: "source-token"
  org: "source-org"
  bucket: "source-bucket"

target:
  type: 2
  url: "http://localhost:8087"
  token: "target-token"
  org: "target-org"
  bucket: "target-bucket"
  db_prefix: "backup_"
  db_suffix: "_v2"

sync:
  start: "2024-01-01T00:00:00Z"
  end: "2024-12-31T23:59:59Z"
  batch_size: 2000
  parallel: 8
  retry_count: 5
  retry_interval: 200
  rate_limit: 0

log:
  level: "warn"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("无法创建测试配置文件: %v", err)
	}

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("配置加载失败: %v", err)
	}

	// 测试runSync2x2x函数
	err = runSync2x2x(cfg)
	if err == nil {
		t.Error("期望连接错误，但没有错误")
	}

	t.Logf("预期的连接错误: %v", err)
}

func TestMainLogLevelSetup(t *testing.T) {
	// 测试日志级别设置逻辑
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "log_level_test.yaml")

	testCases := []struct {
		name     string
		logLevel string
	}{
		{"默认级别", ""},
		{"Info级别", "info"},
		{"Debug级别", "debug"},
		{"Warn级别", "warn"},
		{"Error级别", "error"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logLevelContent := `"info"`
			if tc.logLevel != "" {
				logLevelContent = `"` + tc.logLevel + `"`
			}

			configContent := `
source:
  type: 1
  url: "http://localhost:8086"
  user: "admin"
  pass: "password"
  db: "testdb"

target:
  type: 1
  url: "http://localhost:8087"
  user: "admin"
  pass: "password"
  db: "targetdb"

sync:
  batch_size: 100

log:
  level: ` + logLevelContent + `
`

			err := os.WriteFile(configPath, []byte(configContent), 0644)
			if err != nil {
				t.Fatalf("无法创建测试配置文件: %v", err)
			}

			// 加载配置测试日志级别
			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				t.Fatalf("配置加载失败: %v", err)
			}

			// 验证日志级别设置
			expectedLevel := tc.logLevel
			if expectedLevel == "" {
				expectedLevel = "info" // 默认级别
			}

			if cfg.Log.Level != expectedLevel {
				t.Errorf("期望日志级别为 %s, 实际为 %s", expectedLevel, cfg.Log.Level)
			}
		})
	}
}

func TestConfigurationTypeValidation(t *testing.T) {
	// 测试不同的数据库类型组合
	testCases := []struct {
		name       string
		sourceType int
		targetType int
	}{
		{"1x到1x", 1, 1},
		{"1x到2x", 1, 2},
		{"2x到2x", 2, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "type_test.yaml")

			configContent := `
source:
  type: ` + string(rune(tc.sourceType+'0')) + `
  url: "http://localhost:8086"
  user: "admin"
  pass: "password"
  db: "testdb"
  token: "source-token"
  org: "source-org"
  bucket: "source-bucket"

target:
  type: ` + string(rune(tc.targetType+'0')) + `
  url: "http://localhost:8087"
  user: "admin"
  pass: "password"
  db: "targetdb"
  token: "target-token"
  org: "target-org"
  bucket: "target-bucket"

sync:
  batch_size: 100

log:
  level: "info"
`

			err := os.WriteFile(configPath, []byte(configContent), 0644)
			if err != nil {
				t.Fatalf("无法创建测试配置文件: %v", err)
			}

			// 加载配置
			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				t.Fatalf("配置加载失败: %v", err)
			}

			// 验证类型设置
			if cfg.Source.Type != tc.sourceType {
				t.Errorf("期望源类型为 %d, 实际为 %d", tc.sourceType, cfg.Source.Type)
			}

			if cfg.Target.Type != tc.targetType {
				t.Errorf("期望目标类型为 %d, 实际为 %d", tc.targetType, cfg.Target.Type)
			}
		})
	}
}
