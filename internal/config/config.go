package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Type      int      `yaml:"type"` // 1: InfluxDB 1.x, 2: InfluxDB 2.x
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
