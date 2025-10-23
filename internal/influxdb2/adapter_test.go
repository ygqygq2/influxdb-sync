package influxdb2

import (
	"testing"

	"github.com/ygqygq2/influxdb-sync/internal/testutil"
)

func TestNewAdapter(t *testing.T) {
	// 测试适配器创建
	adapter := &Adapter{
		URL:    "http://localhost:8086",
		Token:  "test-token",
		Org:    "test-org",
		Bucket: "test-bucket",
	}

	if adapter.URL != "http://localhost:8086" {
		t.Errorf("期望URL为 http://localhost:8086, 实际为 %s", adapter.URL)
	}

	if adapter.Token != "test-token" {
		t.Errorf("期望Token为 test-token, 实际为 %s", adapter.Token)
	}

	if adapter.Org != "test-org" {
		t.Errorf("期望Org为 test-org, 实际为 %s", adapter.Org)
	}

	if adapter.Bucket != "test-bucket" {
		t.Errorf("期望Bucket为 test-bucket, 实际为 %s", adapter.Bucket)
	}
}

func TestAdapterStructFields(t *testing.T) {
	// 测试适配器结构体字段
	adapter := &Adapter{}

	// 设置各个字段
	adapter.URL = "https://cloud.influxdata.com"
	adapter.Token = "my-super-secret-token"
	adapter.Org = "my-organization"
	adapter.Bucket = "my-bucket"

	// 验证字段设置
	if adapter.URL != "https://cloud.influxdata.com" {
		t.Error("URL字段设置失败")
	}

	if adapter.Token != "my-super-secret-token" {
		t.Error("Token字段设置失败")
	}

	if adapter.Org != "my-organization" {
		t.Error("Org字段设置失败")
	}

	if adapter.Bucket != "my-bucket" {
		t.Error("Bucket字段设置失败")
	}
}

func TestAdapterConfiguration(t *testing.T) {
	// 测试不同的配置组合
	testCases := []struct {
		name   string
		url    string
		token  string
		org    string
		bucket string
	}{
		{"本地实例", "http://localhost:8086", "local-token", "local-org", "local-bucket"},
		{"云实例", "https://cloud.influxdata.com", "cloud-token", "cloud-org", "cloud-bucket"},
		{"自定义端口", "http://custom-host:9999", "custom-token", "custom-org", "custom-bucket"},
		{"空bucket", "http://localhost:8086", "token", "org", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adapter := &Adapter{
				URL:    tc.url,
				Token:  tc.token,
				Org:    tc.org,
				Bucket: tc.bucket,
			}

			if adapter.URL != tc.url {
				t.Errorf("URL不匹配: 期望 %s, 实际 %s", tc.url, adapter.URL)
			}

			if adapter.Token != tc.token {
				t.Errorf("Token不匹配: 期望 %s, 实际 %s", tc.token, adapter.Token)
			}

			if adapter.Org != tc.org {
				t.Errorf("Org不匹配: 期望 %s, 实际 %s", tc.org, adapter.Org)
			}

			if adapter.Bucket != tc.bucket {
				t.Errorf("Bucket不匹配: 期望 %s, 实际 %s", tc.bucket, adapter.Bucket)
			}
		})
	}
}

func TestAdapterEmptyConfiguration(t *testing.T) {
	// 测试空配置
	adapter := &Adapter{}

	if adapter.URL != "" {
		t.Error("空适配器的URL应该为空字符串")
	}

	if adapter.Token != "" {
		t.Error("空适配器的Token应该为空字符串")
	}

	if adapter.Org != "" {
		t.Error("空适配器的Org应该为空字符串")
	}

	if adapter.Bucket != "" {
		t.Error("空适配器的Bucket应该为空字符串")
	}

	if adapter.client != nil {
		t.Error("空适配器的client应该为nil")
	}
}

func TestAdapterMultipleInstances(t *testing.T) {
	// 测试多个适配器实例的独立性
	adapter1 := &Adapter{
		URL:    "http://server1:8086",
		Token:  "token1",
		Org:    "org1",
		Bucket: "bucket1",
	}

	adapter2 := &Adapter{
		URL:    "http://server2:8086",
		Token:  "token2",
		Org:    "org2",
		Bucket: "bucket2",
	}

	// 验证两个适配器实例互不影响
	if adapter1.URL == adapter2.URL {
		t.Error("两个适配器实例的URL不应该相同")
	}

	if adapter1.Token == adapter2.Token {
		t.Error("两个适配器实例的Token不应该相同")
	}

	if adapter1.Org == adapter2.Org {
		t.Error("两个适配器实例的Org不应该相同")
	}

	if adapter1.Bucket == adapter2.Bucket {
		t.Error("两个适配器实例的Bucket不应该相同")
	}

	// 修改adapter1不应该影响adapter2
	adapter1.URL = "http://modified:8086"
	if adapter2.URL == adapter1.URL {
		t.Error("修改adapter1不应该影响adapter2")
	}
}

func TestAdapter_Connect(t *testing.T) {
	// 需要实际的数据库实例，跳过
	t.Skip("需要实际的 InfluxDB 2.x 实例")
	
	adapter := &Adapter{
		URL:    "http://localhost:8086",
		Token:  "test-token",
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	
	err := adapter.Connect()
	if err != nil {
		t.Logf("Connect 错误: %v", err)
	}
	defer adapter.Close()
}

func TestAdapter_GetDatabases(t *testing.T) {
	// 需要实际的数据库实例，跳过
	t.Skip("需要实际的 InfluxDB 2.x 实例")
	
	adapter := &Adapter{
		URL:    "http://localhost:8086",
		Token:  "test-token",
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	
	err := adapter.Connect()
	if err != nil {
		t.Skip("无法连接到数据库")
	}
	defer adapter.Close()
	
	buckets, err := adapter.GetDatabases()
	if err != nil {
		t.Logf("GetDatabases 错误: %v", err)
	}
	
	t.Logf("找到 %d 个buckets", len(buckets))
}

func TestAdapter_GetMeasurements(t *testing.T) {
	// 需要实际的数据库实例，跳过
	t.Skip("需要实际的 InfluxDB 2.x 实例")
	
	adapter := &Adapter{
		URL:    "http://localhost:8086",
		Token:  "test-token",
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	
	err := adapter.Connect()
	if err != nil {
		t.Skip("无法连接到数据库")
	}
	defer adapter.Close()
	
	measurements, err := adapter.GetMeasurements("test-bucket")
	if err != nil {
		t.Logf("GetMeasurements 错误: %v", err)
	}
	
	t.Logf("找到 %d 个measurements", len(measurements))
}

func TestAdapter_WritePoints(t *testing.T) {
	// 需要实际的数据库实例，跳过
	t.Skip("需要实际的 InfluxDB 2.x 实例")
	
	adapter := &Adapter{
		URL:    "http://localhost:8086",
		Token:  "test-token",
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	
	err := adapter.Connect()
	if err != nil {
		t.Skip("无法连接到数据库")
	}
	defer adapter.Close()
	
	// 模拟写入数据点
	// 实际测试需要有效的数据点
}

// Integration tests with real InfluxDB 2.x (docker)
func TestAdapter_WritePoints_Integration(t *testing.T) {
	testutil.SkipIfNoInfluxDB3(t, "http://localhost:18088")

	adapter := &Adapter{
		URL:    "http://localhost:18088",
		Token:  "test3xtoken",
		Org:    "testorg",
		Bucket: "testbucket",
	}

	err := adapter.Connect()
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer adapter.Close()

	// TODO: 添加真实写入测试
	t.Log("InfluxDB 2.x WritePoints integration test placeholder")
}

func TestAdapter_GetMeasurements_Integration(t *testing.T) {
	testutil.SkipIfNoInfluxDB3(t, "http://localhost:18088")

	adapter := &Adapter{
		URL:    "http://localhost:18088",
		Token:  "test3xtoken",
		Org:    "testorg",
		Bucket: "testbucket",
	}

	err := adapter.Connect()
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer adapter.Close()

	measurements, err := adapter.GetMeasurements("testbucket")
	if err != nil {
		t.Fatalf("GetMeasurements 失败: %v", err)
	}

	t.Logf("成功获取 %d 个 measurements", len(measurements))
	if len(measurements) > 0 {
		t.Logf("第一个 measurement: %s", measurements[0])
	}
}

func TestAdapter_GetTagKeys_Integration(t *testing.T) {
	testutil.SkipIfNoInfluxDB3(t, "http://localhost:18088")

	adapter := &Adapter{
		URL:    "http://localhost:18088",
		Token:  "test3xtoken",
		Org:    "testorg",
		Bucket: "testbucket",
	}

	err := adapter.Connect()
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer adapter.Close()

	// 先获取 measurements
	measurements, err := adapter.GetMeasurements("testbucket")
	if err != nil || len(measurements) == 0 {
		t.Skip("没有可用的 measurements 进行测试")
	}

	// 测试获取第一个 measurement 的 tag keys
	tagKeys, err := adapter.GetTagKeys("testbucket", measurements[0])
	if err != nil {
		t.Fatalf("GetTagKeys 失败: %v", err)
	}

	t.Logf("成功获取 %d 个 tag keys", len(tagKeys))
	for key := range tagKeys {
		t.Logf("  - %s", key)
	}
}

func TestAdapter_QueryData_Integration(t *testing.T) {
	testutil.SkipIfNoInfluxDB3(t, "http://localhost:18088")

	adapter := &Adapter{
		URL:    "http://localhost:18088",
		Token:  "test3xtoken",
		Org:    "testorg",
		Bucket: "testbucket",
	}

	err := adapter.Connect()
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer adapter.Close()

	// 先获取 measurements
	measurements, err := adapter.GetMeasurements("testbucket")
	if err != nil || len(measurements) == 0 {
		t.Skip("没有可用的 measurements 进行测试")
	}

	// 查询数据
	points, lastTime, err := adapter.QueryData("testbucket", measurements[0], 0, 100)
	if err != nil {
		t.Fatalf("QueryData 失败: %v", err)
	}

	t.Logf("成功查询 %d 个数据点, lastTime=%d", len(points), lastTime)
	if len(points) > 0 {
		t.Logf("第一个点: measurement=%s, time=%v, tags=%v, fields=%v",
			points[0].Measurement, points[0].Time, points[0].Tags, points[0].Fields)
	}
}
