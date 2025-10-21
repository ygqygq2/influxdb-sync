package config

import (
	"os"

	"github.com/ygqygq2/influxdb-sync/internal/influxdb3"
	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Type      int      `yaml:"type"` // 1: InfluxDB 1.x, 2: InfluxDB 2.x, 3: InfluxDB 3.x
	URL       string   `yaml:"url"`
	User      string   `yaml:"user"`
	Pass      string   `yaml:"pass"`
	DB        string   `yaml:"db"`
	DBExclude []string `yaml:"db_exclude"`
	DBPrefix  string   `yaml:"db_prefix"`
	DBSuffix  string   `yaml:"db_suffix"`
	Token     string   `yaml:"token"`
	Org       string   `yaml:"org"`
	Bucket    string   `yaml:"bucket"`
	// InfluxDB 3.x 特有配置
	CompatMode string `yaml:"compat_mode"` // "v1", "v2", "native"
	Namespace  string `yaml:"namespace"`   // 3.x 原生模式的命名空间
	Database   string `yaml:"database"`    // 3.x 数据库名称
}

type SyncConfig struct {
	Start         string `yaml:"start"`
	End           string `yaml:"end"`
	BatchSize     int    `yaml:"batch_size"`
	ResumeFile    string `yaml:"resume_file"`
	Parallel      int    `yaml:"parallel"`
	RetryCount    int    `yaml:"retry_count"`
	RetryInterval int    `yaml:"retry_interval"`
	RateLimit     int    `yaml:"rate_limit"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type Config struct {
	Source DBConfig   `yaml:"source"`
	Target DBConfig   `yaml:"target"`
	Sync   SyncConfig `yaml:"sync"`
	Log    LogConfig  `yaml:"log"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return &cfg, err
}

// ToV1CompatConfig 将 DBConfig 转换为 InfluxDB 3.x v1 兼容配置
func (db *DBConfig) ToV1CompatConfig() *influxdb3.V1CompatConfig {
	return &influxdb3.V1CompatConfig{
		Addr:     db.URL,
		User:     db.User,
		Pass:     db.Pass,
		Database: db.Database,
	}
}

// ToV2CompatConfig 将 DBConfig 转换为 InfluxDB 3.x v2 兼容配置
func (db *DBConfig) ToV2CompatConfig() *influxdb3.V2CompatConfig {
	return &influxdb3.V2CompatConfig{
		URL:      db.URL,
		Token:    db.Token,
		Org:      db.Org,
		Database: db.Database,
	}
}

// ToNativeConfig 将 DBConfig 转换为 InfluxDB 3.x 原生配置
func (db *DBConfig) ToNativeConfig() *influxdb3.NativeConfig {
	return &influxdb3.NativeConfig{
		URL:       db.URL,
		Token:     db.Token,
		Database:  db.Database,
		Namespace: db.Namespace,
	}
}
