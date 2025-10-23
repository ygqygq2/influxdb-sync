package influxdb1

import (
	"testing"
	"time"

	"github.com/ygqygq2/influxdb-sync/internal/common"
	"github.com/ygqygq2/influxdb-sync/internal/testutil"
)

func TestEscapeMeasurement(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"simple", `"simple"`},
		{"measurement", `"measurement"`},
		{"with spaces", `"with spaces"`},
		{"with\"quotes", `"with\"quotes"`},
		{"with'single'quotes", `"with'single'quotes"`},
		{"", `""`},
		{"special-chars_123", `"special-chars_123"`},
		{"测量名称", `"测量名称"`}, // 中文测量名
	}

	for _, tc := range testCases {
		result := escapeMeasurement(tc.input)
		if result != tc.expected {
			t.Errorf("escapeMeasurement(%q) = %q, 期望 %q", tc.input, result, tc.expected)
		}
	}
}

func TestNewDataSource(t *testing.T) {
	config := DataSourceConfig{
		Addr: "http://localhost:8086",
		User: "admin",
		Pass: "password",
	}

	ds := NewDataSource(config)
	if ds == nil {
		t.Fatal("NewDataSource返回nil")
	}

	if ds.config.Addr != config.Addr {
		t.Errorf("期望地址为 %s, 实际为 %s", config.Addr, ds.config.Addr)
	}

	if ds.config.User != config.User {
		t.Errorf("期望用户为 %s, 实际为 %s", config.User, ds.config.User)
	}

	if ds.config.Pass != config.Pass {
		t.Errorf("期望密码为 %s, 实际为 %s", config.Pass, ds.config.Pass)
	}
}

func TestNewDataTarget(t *testing.T) {
	config := DataTargetConfig{
		Addr: "http://localhost:8087",
		User: "target_user",
		Pass: "target_pass",
	}

	dt := NewDataTarget(config)
	if dt == nil {
		t.Fatal("NewDataTarget返回nil")
	}

	if dt.config.Addr != config.Addr {
		t.Errorf("期望地址为 %s, 实际为 %s", config.Addr, dt.config.Addr)
	}

	if dt.config.User != config.User {
		t.Errorf("期望用户为 %s, 实际为 %s", config.User, dt.config.User)
	}

	if dt.config.Pass != config.Pass {
		t.Errorf("期望密码为 %s, 实际为 %s", config.Pass, dt.config.Pass)
	}
}

func TestDataSourceConfigStruct(t *testing.T) {
	// 测试DataSourceConfig结构体
	config := DataSourceConfig{
		Addr: "https://influxdb.example.com:8086",
		User: "test_user",
		Pass: "test_password",
	}

	if config.Addr != "https://influxdb.example.com:8086" {
		t.Errorf("地址设置不正确")
	}

	if config.User != "test_user" {
		t.Errorf("用户名设置不正确")
	}

	if config.Pass != "test_password" {
		t.Errorf("密码设置不正确")
	}
}

func TestDataTargetConfigStruct(t *testing.T) {
	// 测试DataTargetConfig结构体
	config := DataTargetConfig{
		Addr: "https://target.example.com:8086",
		User: "target_admin",
		Pass: "target_secret",
	}

	if config.Addr != "https://target.example.com:8086" {
		t.Errorf("目标地址设置不正确")
	}

	if config.User != "target_admin" {
		t.Errorf("目标用户名设置不正确")
	}

	if config.Pass != "target_secret" {
		t.Errorf("目标密码设置不正确")
	}
}

func TestDataSourceAndTargetSeparation(t *testing.T) {
	// 测试数据源和目标配置的独立性
	sourceConfig := DataSourceConfig{
		Addr: "http://source:8086",
		User: "source_user",
		Pass: "source_pass",
	}

	targetConfig := DataTargetConfig{
		Addr: "http://target:8086",
		User: "target_user",
		Pass: "target_pass",
	}

	source := NewDataSource(sourceConfig)
	target := NewDataTarget(targetConfig)

	// 验证配置不会互相影响
	if source.config.Addr == target.config.Addr {
		t.Error("源和目标地址不应该相同")
	}

	if source.config.User == target.config.User {
		t.Error("源和目标用户不应该相同")
	}

	if source.config.Pass == target.config.Pass {
		t.Error("源和目标密码不应该相同")
	}
}

func TestDataSource_GetTagKeys(t *testing.T) {
	testutil.SkipIfNoInfluxDB1(t, "http://localhost:18086")
	
	config := DataSourceConfig{
		Addr: "http://localhost:18086",
		User: "admin",
		Pass: "admin123",
	}
	
	ds := NewDataSource(config)
	err := ds.Connect()
	if err != nil {
		t.Fatalf("无法连接到数据库: %v", err)
	}
	defer ds.Close()
	
	tagKeys, err := ds.GetTagKeys("testdb", "cpu")
	if err != nil {
		t.Logf("GetTagKeys 错误: %v", err)
	}
	
	if len(tagKeys) > 0 {
		t.Logf("找到 %d 个标签", len(tagKeys))
	}
}

func TestDataSource_QueryData(t *testing.T) {
	testutil.SkipIfNoInfluxDB1(t, "http://localhost:18086")
	
	config := DataSourceConfig{
		Addr: "http://localhost:18086",
		User: "admin",
		Pass: "admin123",
	}
	
	ds := NewDataSource(config)
	err := ds.Connect()
	if err != nil {
		t.Fatalf("无法连接到数据库: %v", err)
	}
	defer ds.Close()
	
	points, maxTime, err := ds.QueryData("testdb", "cpu", 0, 1000)
	if err != nil {
		t.Logf("QueryData 错误: %v", err)
	}
	
	t.Logf("查询到 %d 个数据点，最大时间戳 %d", len(points), maxTime)
}

func TestDataTarget_WritePoints(t *testing.T) {
	testutil.SkipIfNoInfluxDB1(t, "http://localhost:18086")
	
	config := DataTargetConfig{
		Addr: "http://localhost:18086",
		User: "admin",
		Pass: "admin123",
	}
	
	dt := NewDataTarget(config)
	err := dt.Connect()
	if err != nil {
		t.Fatalf("无法连接到数据库: %v", err)
	}
	defer dt.Close()
	
	// 写入测试数据点
	now := time.Now()
	points := []common.DataPoint{
		{
			Measurement: "test_measurement",
			Tags: map[string]string{
				"host": "test_host",
				"region": "test_region",
			},
			Fields: map[string]interface{}{
				"value": 42.5,
				"count": 100,
			},
			Time: now,
		},
	}
	
	err = dt.WritePoints("testdb", points)
	if err != nil {
		t.Logf("WritePoints 错误: %v", err)
	} else {
		t.Logf("成功写入 %d 个数据点到 testdb.test_measurement", len(points))
	}
}
