package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	   Type   int    `yaml:"type"`   // 1: InfluxDB 1.x, 2: InfluxDB 2.x
	   URL    string `yaml:"url"`
	   User   string `yaml:"user"`
	   Pass   string `yaml:"pass"`
	   DB     string `yaml:"db"`
	   Token  string `yaml:"token"`
	   Org    string `yaml:"org"`
	   Bucket string `yaml:"bucket"`
}

type SyncConfig struct {
	   Start      string `yaml:"start"`
	   End        string `yaml:"end"`
	   BatchSize  int    `yaml:"batch_size"`
	   ResumeFile string `yaml:"resume_file"`
}

type Config struct {
	   Source DBConfig   `yaml:"source"`
	   Target DBConfig   `yaml:"target"`
	   Sync   SyncConfig `yaml:"sync"`
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
