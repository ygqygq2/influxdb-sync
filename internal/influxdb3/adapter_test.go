package influxdb3

import (
	"testing"

	"github.com/ygqygq2/influxdb-sync/internal/common"
)

func TestNewDataSource3x(t *testing.T) {
	tests := []struct {
		name    string
		config  interface{}
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
			name:    "invalid config type",
			config:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds, err := NewDataSource3x(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDataSource3x() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ds == nil && !tt.wantErr {
				t.Errorf("NewDataSource3x() returned nil")
			}
			if ds != nil {
				defer ds.Close()
			}
		})
	}
}

func TestNewDataTarget3x(t *testing.T) {
	tests := []struct {
		name    string
		config  interface{}
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
			name:    "invalid config type",
			config:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt, err := NewDataTarget3x(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDataTarget3x() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if dt == nil && !tt.wantErr {
				t.Errorf("NewDataTarget3x() returned nil")
			}
			if dt != nil {
				defer dt.Close()
			}
		})
	}
}

func TestDataSource3x_GetMeasurements(t *testing.T) {
	// 需要实际的 InfluxDB 3.x 实例
	t.Skip("需要实际的 InfluxDB 3.x 实例")

	config := V1CompatConfig{
		Addr:     "http://localhost:8086",
		User:     "admin",
		Pass:     "password",
		Database: "test-db",
	}

	ds, err := NewDataSource3x(config)
	if err != nil {
		t.Fatalf("NewDataSource3x() error = %v", err)
	}
	defer ds.Close()

	err = ds.Connect()
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	measurements, err := ds.GetMeasurements("test-db")
	if err != nil {
		t.Errorf("GetMeasurements() error = %v", err)
	}

	t.Logf("Found %d measurements", len(measurements))
}

func TestDataTarget3x_WritePoints(t *testing.T) {
	// 需要实际的 InfluxDB 3.x 实例
	t.Skip("需要实际的 InfluxDB 3.x 实例")

	config := V2CompatConfig{
		URL:      "http://localhost:8086",
		Token:    "test-token",
		Org:      "test-org",
		Database: "test-db",
	}

	dt, err := NewDataTarget3x(config)
	if err != nil {
		t.Fatalf("NewDataTarget3x() error = %v", err)
	}
	defer dt.Close()

	err = dt.Connect()
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	// 创建测试数据点
	points := []common.DataPoint{
		{
			Measurement: "test_measurement",
			Tags: map[string]string{
				"tag1": "value1",
			},
			Fields: map[string]interface{}{
				"field1": 1.0,
			},
		},
	}

	err = dt.WritePoints("test-db", points)
	if err != nil {
		t.Errorf("WritePoints() error = %v", err)
	}
}
