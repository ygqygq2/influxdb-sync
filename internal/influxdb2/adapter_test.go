package influxdb2

import (
	"testing"
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
