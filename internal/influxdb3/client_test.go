package influxdb3

import (
	"fmt"
	"testing"
	"time"
)

func TestNewClient3x(t *testing.T) {
	tests := []struct {
		name    string
		config  NativeConfig
		wantErr bool
	}{
		{
			name: "valid native config",
			config: NativeConfig{
				URL:       "http://localhost:8086",
				Token:     "test-token",
				Database:  "test-db",
				Namespace: "test-ns",
			},
			wantErr: false,
		},
		{
			name:    "empty config",
			config:  NativeConfig{},
			wantErr: false, // 客户端创建不会失败，但连接会失败
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient3x(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient3x() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if client == nil && !tt.wantErr {
				t.Errorf("NewClient3x() returned nil client")
			}
			if client != nil {
				defer client.Close()
			}
		})
	}
}

func TestNewV1CompatClient(t *testing.T) {
	tests := []struct {
		name    string
		config  V1CompatConfig
		wantErr bool
	}{
		{
			name: "valid v1 compat config",
			config: V1CompatConfig{
				Addr:     "http://localhost:8086",
				User:     "admin",
				Pass:     "password",
				Database: "test-db",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewV1CompatClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewV1CompatClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if client == nil && !tt.wantErr {
				t.Errorf("NewV1CompatClient() returned nil client")
			}
			if client != nil {
				defer client.Close()
			}
		})
	}
}

func TestNewV2CompatClient(t *testing.T) {
	tests := []struct {
		name    string
		config  V2CompatConfig
		wantErr bool
	}{
		{
			name: "valid v2 compat config",
			config: V2CompatConfig{
				URL:      "http://localhost:8086",
				Token:    "test-token",
				Org:      "test-org",
				Database: "test-db",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewV2CompatClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewV2CompatClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if client == nil && !tt.wantErr {
				t.Errorf("NewV2CompatClient() returned nil client")
			}
			if client != nil {
				defer client.Close()
			}
		})
	}
}

func TestClient3x_WriteLineProtocol(t *testing.T) {
	// 这个测试需要实际的 InfluxDB 3.x 实例，所以跳过
	t.Skip("需要实际的 InfluxDB 3.x 实例")

	config := NativeConfig{
		URL:       "http://localhost:8086",
		Token:     "test-token",
		Database:  "test-db",
		Namespace: "test-ns",
	}

	client, err := NewClient3x(config)
	if err != nil {
		t.Fatalf("NewClient3x() error = %v", err)
	}
	defer client.Close()

	// 测试数据
	data := fmt.Sprintf("measurement,tag1=value1 field1=1.0 %d", time.Now().UnixNano())

	err = client.WriteLineProtocol(data)
	if err != nil {
		t.Errorf("WriteLineProtocol() error = %v", err)
	}
}
