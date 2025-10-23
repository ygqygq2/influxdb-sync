package influxdb3

import (
	"testing"
	"time"

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

func TestEscapeMeasurement(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal measurement",
			input: "cpu",
			want:  `"cpu"`,
		},
		{
			name:  "measurement with quotes",
			input: `my"measurement`,
			want:  `"my\"measurement"`,
		},
		{
			name:  "measurement with spaces",
			input: "cpu usage",
			want:  `"cpu usage"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeMeasurement(tt.input)
			if got != tt.want {
				t.Errorf("escapeMeasurement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatLineProtocol(t *testing.T) {
	tests := []struct {
		name  string
		point common.DataPoint
		want  string
	}{
		{
			name: "simple point",
			point: common.DataPoint{
				Measurement: "cpu",
				Tags: map[string]string{
					"host": "server01",
				},
				Fields: map[string]interface{}{
					"value": 0.64,
				},
				Time: time.Unix(0, 1609459200000000000),
			},
			want: "cpu,host=server01 value=0.64 1609459200000000000",
		},
		{
			name: "point with multiple fields",
			point: common.DataPoint{
				Measurement: "memory",
				Tags: map[string]string{
					"host": "server01",
				},
				Fields: map[string]interface{}{
					"used":      int64(1024),
					"available": 2048.5,
					"cached":    "string_value",
				},
				Time: time.Unix(0, 1609459200000000000),
			},
			want: "", // 因为字段顺序不固定，只验证不为空
		},
		{
			name: "point with no fields",
			point: common.DataPoint{
				Measurement: "cpu",
				Tags: map[string]string{
					"host": "server01",
				},
				Fields: map[string]interface{}{},
				Time:   time.Unix(0, 1609459200000000000),
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatLineProtocol(tt.point)
			if tt.want == "" {
				if tt.name == "point with no fields" {
					if got != "" {
						t.Errorf("formatLineProtocol() should return empty string for no fields")
					}
				} else {
					if got == "" {
						t.Errorf("formatLineProtocol() returned empty string unexpectedly")
					}
				}
			} else if got != tt.want {
				t.Errorf("formatLineProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
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

func TestDataSource3x_ConnectError(t *testing.T) {
	config := V1CompatConfig{
		Addr:     "http://invalid-host:9999",
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
	// 应该返回连接错误
	if err == nil {
		t.Error("应该返回连接错误")
	}
}

func TestDataTarget3x_ConnectError(t *testing.T) {
	config := V2CompatConfig{
		URL:      "http://invalid-host:9999",
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
	// 应该返回连接错误
	if err == nil {
		t.Error("应该返回连接错误")
	}
}

func TestDataSource3x_GetDatabasesNotConnected(t *testing.T) {
	config := V1CompatConfig{
		Addr:     "http://localhost:8086",
		User:     "admin",
		Pass:     "password",
		Database: "test-db",
	}

	ds := NewV1CompatDataSource(config)
	
	_, err := ds.GetDatabases()
	if err == nil {
		t.Error("未连接时应该返回错误")
	}
}

func TestDataSource3x_GetMeasurementsNotConnected(t *testing.T) {
	config := V1CompatConfig{
		Addr:     "http://localhost:8086",
		User:     "admin",
		Pass:     "password",
		Database: "test-db",
	}

	ds := NewV1CompatDataSource(config)
	
	_, err := ds.GetMeasurements("test-db")
	if err == nil {
		t.Error("未连接时应该返回错误")
	}
}

func TestDataSource3x_QueryDataNotConnected(t *testing.T) {
	config := V1CompatConfig{
		Addr:     "http://localhost:8086",
		User:     "admin",
		Pass:     "password",
		Database: "test-db",
	}

	ds := NewV1CompatDataSource(config)
	
	_, _, err := ds.QueryData("test-db", "test_measurement", 0, 1000)
	if err == nil {
		t.Error("未连接时应该返回错误")
	}
}

func TestDataTarget3x_WritePointsNotConnected(t *testing.T) {
	config := V2CompatConfig{
		URL:      "http://localhost:8086",
		Token:    "test-token",
		Org:      "test-org",
		Database: "test-db",
	}

	dt := NewV2CompatDataTarget(config)
	
	points := []common.DataPoint{
		{
			Measurement: "test",
			Tags:        map[string]string{"host": "server1"},
			Fields:      map[string]interface{}{"value": 1.0},
			Time:        time.Now(),
		},
	}
	
	err := dt.WritePoints("test-db", points)
	if err == nil {
		t.Error("未连接时应该返回错误")
	}
}

func TestDataSource3x_Close(t *testing.T) {
	config := V1CompatConfig{
		Addr:     "http://localhost:8086",
		User:     "admin",
		Pass:     "password",
		Database: "test-db",
	}

	ds := NewV1CompatDataSource(config)
	err := ds.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestDataTarget3x_Close(t *testing.T) {
	config := V2CompatConfig{
		URL:      "http://localhost:8086",
		Token:    "test-token",
		Org:      "test-org",
		Database: "test-db",
	}

	dt := NewV2CompatDataTarget(config)
	err := dt.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

// Integration tests with real InfluxDB (docker)
func TestDataTarget3x_WriteLineProtocol_Integration(t *testing.T) {
	// 使用 docker 容器 influxdb:2.7 在 18088 端口
	config := V2CompatConfig{
		URL:      "http://localhost:18088",
		Token:    "test3xtoken",
		Org:      "testorg",
		Database: "testbucket",
	}

	dt := NewV2CompatDataTarget(config)
	if dt == nil {
		t.Fatal("NewV2CompatDataTarget returned nil")
	}

	// 尝试连接
	err := dt.Connect()
	if err != nil {
		t.Skipf("无法连接到 InfluxDB 3.x (v2 compat) at localhost:18088: %v", err)
		return
	}
	defer dt.Close()

	// 测试写入数据
	points := []common.DataPoint{
		{
			Measurement: "test_measurement_3x",
			Tags: map[string]string{
				"host":   "server01",
				"region": "us-west",
			},
			Fields: map[string]interface{}{
				"cpu":    float64(80.5),
				"memory": float64(1024),
				"status": "online",
			},
			Time: time.Now(),
		},
		{
			Measurement: "test_measurement_3x",
			Tags: map[string]string{
				"host":   "server02",
				"region": "us-east",
			},
			Fields: map[string]interface{}{
				"cpu":    float64(60.3),
				"memory": float64(2048),
			},
			Time: time.Now().Add(-time.Minute),
		},
	}

	err = dt.WritePoints("testbucket", points)
	if err != nil {
		t.Fatalf("WritePoints() error = %v", err)
	}

	t.Log("Successfully wrote points to InfluxDB 3.x (v2 compat)")
}

func TestDataSource3x_QueryData_V2Compat_Integration(t *testing.T) {
	// 使用 docker 容器 influxdb:2.7 在 18088 端口
	config := V2CompatConfig{
		URL:      "http://localhost:18088",
		Token:    "test3xtoken",
		Org:      "testorg",
		Database: "testbucket",
	}

	ds := NewV2CompatDataSource(config)
	if ds == nil {
		t.Fatal("NewV2CompatDataSource returned nil")
	}

	err := ds.Connect()
	if err != nil {
		t.Skipf("无法连接到 InfluxDB 3.x (v2 compat) at localhost:18088: %v", err)
		return
	}
	defer ds.Close()

	// 先写入测试数据
	dt := NewV2CompatDataTarget(config)
	_ = dt.Connect()
	testPoints := []common.DataPoint{
		{
			Measurement: "query_test_3x",
			Tags:        map[string]string{"sensor": "temp01"},
			Fields:      map[string]interface{}{"value": float64(25.5)},
			Time:        time.Now(),
		},
	}
	_ = dt.WritePoints("testbucket", testPoints)
	dt.Close()

	// 等待数据写入
	time.Sleep(2 * time.Second)

	// 测试查询 (startTime 使用 0 表示从最早开始，batchSize 使用 1000)
	points, lastTime, err := ds.QueryData("testbucket", "query_test_3x", 0, 1000)
	if err != nil {
		t.Fatalf("QueryData() error = %v", err)
	}

	if len(points) == 0 {
		t.Log("Warning: No points returned, but no error")
	} else {
		t.Logf("Successfully queried %d points from InfluxDB 3.x (v2 compat), lastTime=%d", len(points), lastTime)
	}
}

func TestDataSource3x_GetTagKeys_V2Compat_Integration(t *testing.T) {
	config := V2CompatConfig{
		URL:      "http://localhost:18088",
		Token:    "test3xtoken",
		Org:      "testorg",
		Database: "testbucket",
	}

	ds := NewV2CompatDataSource(config)
	if ds == nil {
		t.Fatal("NewV2CompatDataSource returned nil")
	}

	err := ds.Connect()
	if err != nil {
		t.Skipf("无法连接到 InfluxDB 3.x (v2 compat) at localhost:18088: %v", err)
		return
	}
	defer ds.Close()

	// 先写入带标签的测试数据
	dt := NewV2CompatDataTarget(config)
	_ = dt.Connect()
	testPoints := []common.DataPoint{
		{
			Measurement: "tagkeys_test_3x",
			Tags: map[string]string{
				"location": "beijing",
				"device":   "sensor01",
			},
			Fields: map[string]interface{}{"temp": float64(23.5)},
			Time:   time.Now(),
		},
	}
	_ = dt.WritePoints("testbucket", testPoints)
	dt.Close()

	time.Sleep(2 * time.Second)

	// 测试获取 tag keys
	tagKeys, err := ds.GetTagKeys("testbucket", "tagkeys_test_3x")
	if err != nil {
		t.Fatalf("GetTagKeys() error = %v", err)
	}

	if len(tagKeys) == 0 {
		t.Log("Warning: No tag keys returned")
	} else {
		t.Logf("Successfully got %d tag keys: %v", len(tagKeys), tagKeys)
	}
}

func TestDataSource3x_GetMeasurements_V2Compat_Integration(t *testing.T) {
	config := V2CompatConfig{
		URL:      "http://localhost:18088",
		Token:    "test3xtoken",
		Org:      "testorg",
		Database: "testbucket",
	}

	ds := NewV2CompatDataSource(config)
	if ds == nil {
		t.Fatal("NewV2CompatDataSource returned nil")
	}

	err := ds.Connect()
	if err != nil {
		t.Skipf("无法连接到 InfluxDB 3.x (v2 compat) at localhost:18088: %v", err)
		return
	}
	defer ds.Close()

	// 测试获取 measurements
	measurements, err := ds.GetMeasurements("testbucket")
	if err != nil {
		t.Fatalf("GetMeasurements() error = %v", err)
	}

	t.Logf("Successfully got %d measurements", len(measurements))
	if len(measurements) > 0 {
		t.Logf("First measurement: %s", measurements[0])
	}
}

func TestDataSource3x_GetDatabases_V2Compat_Integration(t *testing.T) {
	config := V2CompatConfig{
		URL:      "http://localhost:18088",
		Token:    "test3xtoken",
		Org:      "testorg",
		Database: "testbucket",
	}

	ds := NewV2CompatDataSource(config)
	if ds == nil {
		t.Fatal("NewV2CompatDataSource returned nil")
	}

	err := ds.Connect()
	if err != nil {
		t.Skipf("无法连接到 InfluxDB 3.x (v2 compat) at localhost:18088: %v", err)
		return
	}
	defer ds.Close()

	// 测试获取 databases（在 v2 compat 模式下返回 buckets）
	databases, err := ds.GetDatabases()
	if err != nil {
		t.Fatalf("GetDatabases() error = %v", err)
	}

	t.Logf("Successfully got %d databases/buckets", len(databases))
	for _, db := range databases {
		t.Logf("  - %s", db)
	}
}
