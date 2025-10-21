package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")

	// 创建测试配置文件
	configContent := `
source:
  type: 1
  url: "http://localhost:8086"
  user: "admin"
  pass: "password"
  db: "testdb"
  db_exclude: ["_internal"]
  db_prefix: "src_"
  db_suffix: "_backup"
  token: "test-token"
  org: "test-org"
  bucket: "test-bucket"

target:
  type: 2
  url: "http://localhost:8087"
  user: "target_user"
  pass: "target_pass"
  db: "targetdb"
  db_exclude: ["system"]
  db_prefix: "tgt_"
  db_suffix: "_copy"
  token: "target-token"
  org: "target-org"
  bucket: "target-bucket"

sync:
  start: "2024-01-01T00:00:00Z"
  end: "2024-12-31T23:59:59Z"
  batch_size: 1000
  resume_file: "resume.state"
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

	// 测试配置加载
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("配置加载失败: %v", err)
	}

	// 验证源配置
	if cfg.Source.Type != 1 {
		t.Errorf("期望源类型为 1, 实际为 %d", cfg.Source.Type)
	}
	if cfg.Source.URL != "http://localhost:8086" {
		t.Errorf("期望源URL为 http://localhost:8086, 实际为 %s", cfg.Source.URL)
	}
	if cfg.Source.User != "admin" {
		t.Errorf("期望源用户为 admin, 实际为 %s", cfg.Source.User)
	}
	if cfg.Source.Pass != "password" {
		t.Errorf("期望源密码为 password, 实际为 %s", cfg.Source.Pass)
	}
	if cfg.Source.DB != "testdb" {
		t.Errorf("期望源数据库为 testdb, 实际为 %s", cfg.Source.DB)
	}
	if len(cfg.Source.DBExclude) != 1 || cfg.Source.DBExclude[0] != "_internal" {
		t.Errorf("期望排除数据库为 [_internal], 实际为 %v", cfg.Source.DBExclude)
	}

	// 验证目标配置
	if cfg.Target.Type != 2 {
		t.Errorf("期望目标类型为 2, 实际为 %d", cfg.Target.Type)
	}
	if cfg.Target.URL != "http://localhost:8087" {
		t.Errorf("期望目标URL为 http://localhost:8087, 实际为 %s", cfg.Target.URL)
	}
	if cfg.Target.Token != "target-token" {
		t.Errorf("期望目标token为 target-token, 实际为 %s", cfg.Target.Token)
	}

	// 验证同步配置
	if cfg.Sync.BatchSize != 1000 {
		t.Errorf("期望批次大小为 1000, 实际为 %d", cfg.Sync.BatchSize)
	}
	if cfg.Sync.Parallel != 4 {
		t.Errorf("期望并发数为 4, 实际为 %d", cfg.Sync.Parallel)
	}
	if cfg.Sync.RetryCount != 3 {
		t.Errorf("期望重试次数为 3, 实际为 %d", cfg.Sync.RetryCount)
	}

	// 验证日志配置
	if cfg.Log.Level != "info" {
		t.Errorf("期望日志级别为 info, 实际为 %s", cfg.Log.Level)
	}
}

func TestLoadConfigFileNotExists(t *testing.T) {
	// 测试文件不存在的情况
	_, err := LoadConfig("non_existent_file.yaml")
	if err == nil {
		t.Error("期望返回错误，但没有错误")
	}
}

func TestLoadConfigInvalidYaml(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid_config.yaml")

	// 创建无效的YAML内容
	invalidContent := `
source:
  type: 1
  url: "http://localhost:8086"
target:
  invalid yaml structure
    missing proper indentation
`

	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("无法创建测试配置文件: %v", err)
	}

	// 测试加载无效配置
	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("期望返回YAML解析错误，但没有错误")
	}
}

func TestDBConfig_ToV1CompatConfig(t *testing.T) {
	dbConfig := DBConfig{
		Type:     3,
		URL:      "http://localhost:8086",
		User:     "admin",
		Pass:     "password",
		Database: "testdb",
	}

	v1Config := dbConfig.ToV1CompatConfig()
	if v1Config.Addr != dbConfig.URL {
		t.Errorf("Addr = %v, want %v", v1Config.Addr, dbConfig.URL)
	}
	if v1Config.User != dbConfig.User {
		t.Errorf("User = %v, want %v", v1Config.User, dbConfig.User)
	}
	if v1Config.Pass != dbConfig.Pass {
		t.Errorf("Pass = %v, want %v", v1Config.Pass, dbConfig.Pass)
	}
	if v1Config.Database != dbConfig.Database {
		t.Errorf("Database = %v, want %v", v1Config.Database, dbConfig.Database)
	}
}

func TestDBConfig_ToV2CompatConfig(t *testing.T) {
	dbConfig := DBConfig{
		Type:     3,
		URL:      "http://localhost:8086",
		Token:    "test-token",
		Org:      "test-org",
		Database: "testdb",
	}

	v2Config := dbConfig.ToV2CompatConfig()
	if v2Config.URL != dbConfig.URL {
		t.Errorf("URL = %v, want %v", v2Config.URL, dbConfig.URL)
	}
	if v2Config.Token != dbConfig.Token {
		t.Errorf("Token = %v, want %v", v2Config.Token, dbConfig.Token)
	}
	if v2Config.Org != dbConfig.Org {
		t.Errorf("Org = %v, want %v", v2Config.Org, dbConfig.Org)
	}
	if v2Config.Database != dbConfig.Database {
		t.Errorf("Database = %v, want %v", v2Config.Database, dbConfig.Database)
	}
}

func TestDBConfig_ToNativeConfig(t *testing.T) {
	dbConfig := DBConfig{
		Type:      3,
		URL:       "http://localhost:8086",
		Token:     "test-token",
		Database:  "testdb",
		Namespace: "test-ns",
	}

	nativeConfig := dbConfig.ToNativeConfig()
	if nativeConfig.URL != dbConfig.URL {
		t.Errorf("URL = %v, want %v", nativeConfig.URL, dbConfig.URL)
	}
	if nativeConfig.Token != dbConfig.Token {
		t.Errorf("Token = %v, want %v", nativeConfig.Token, dbConfig.Token)
	}
	if nativeConfig.Database != dbConfig.Database {
		t.Errorf("Database = %v, want %v", nativeConfig.Database, dbConfig.Database)
	}
	if nativeConfig.Namespace != dbConfig.Namespace {
		t.Errorf("Namespace = %v, want %v", nativeConfig.Namespace, dbConfig.Namespace)
	}
}
