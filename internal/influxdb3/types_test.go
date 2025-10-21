package influxdb3

import (
	"testing"
)

func TestV1CompatConfig(t *testing.T) {
	config := V1CompatConfig{
		Addr:     "http://localhost:8086",
		User:     "admin",
		Pass:     "password",
		Database: "testdb",
	}

	if config.Addr == "" {
		t.Error("Addr should not be empty")
	}
	if config.User != "admin" {
		t.Errorf("User = %v, want admin", config.User)
	}
	if config.Database != "testdb" {
		t.Errorf("Database = %v, want testdb", config.Database)
	}
}

func TestV2CompatConfig(t *testing.T) {
	config := V2CompatConfig{
		URL:      "http://localhost:8086",
		Token:    "test-token",
		Org:      "test-org",
		Bucket:   "test-bucket",
		Database: "testdb",
	}

	if config.URL == "" {
		t.Error("URL should not be empty")
	}
	if config.Token != "test-token" {
		t.Errorf("Token = %v, want test-token", config.Token)
	}
	if config.Org != "test-org" {
		t.Errorf("Org = %v, want test-org", config.Org)
	}
}

func TestNativeConfig(t *testing.T) {
	config := NativeConfig{
		URL:       "http://localhost:8086",
		Token:     "test-token",
		Database:  "testdb",
		Namespace: "test-ns",
		UseSQL:    true,
	}

	if config.URL == "" {
		t.Error("URL should not be empty")
	}
	if config.Token != "test-token" {
		t.Errorf("Token = %v, want test-token", config.Token)
	}
	if !config.UseSQL {
		t.Error("UseSQL should be true")
	}
}
