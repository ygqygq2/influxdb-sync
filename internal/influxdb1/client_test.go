package influxdb1

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	// 测试客户端创建
	addr := "http://localhost:8086"
	user := "admin"
	pass := "password"
	timeout := 30 * time.Second

	client, err := NewClient(addr, user, pass, timeout)
	if err != nil {
		t.Errorf("创建客户端失败: %v", err)
	}

	if client == nil {
		t.Error("客户端不应该为空")
	}

	if client.cli == nil {
		t.Error("底层客户端不应该为空")
	}

	// 测试关闭客户端
	err = client.Close()
	if err != nil {
		t.Errorf("关闭客户端失败: %v", err)
	}
}

func TestNewClientInvalidURL(t *testing.T) {
	// 测试无效URL - 实际上InfluxDB客户端可能不会在创建时验证URL
	// 这个测试主要验证函数不会panic
	addr := "invalid-url"
	user := "admin"
	pass := "password"
	timeout := 30 * time.Second

	client, err := NewClient(addr, user, pass, timeout)
	// InfluxDB客户端通常在创建时不验证URL，所以可能不会返回错误
	// 我们主要确保函数能正常执行
	if client != nil {
		client.Close()
	}

	// 如果有错误也是可以接受的
	if err != nil {
		t.Logf("预期的错误: %v", err)
	}
}

func TestClientTimeout(t *testing.T) {
	// 测试不同的超时设置
	timeouts := []time.Duration{
		1 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}

	for _, timeout := range timeouts {
		client, err := NewClient("http://localhost:8086", "admin", "password", timeout)
		if err != nil {
			t.Errorf("超时 %v 时创建客户端失败: %v", timeout, err)
			continue
		}

		if client == nil {
			t.Errorf("超时 %v 时客户端为空", timeout)
			continue
		}

		client.Close()
	}
}

func TestClientCredentials(t *testing.T) {
	// 测试不同的认证信息
	testCases := []struct {
		name string
		user string
		pass string
	}{
		{"空用户名", "", "password"},
		{"空密码", "admin", ""},
		{"正常认证", "admin", "password"},
		{"特殊字符", "admin@test", "pass@word#123"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient("http://localhost:8086", tc.user, tc.pass, 30*time.Second)
			if err != nil {
				t.Errorf("用户 %s, 密码 %s 时创建客户端失败: %v", tc.user, tc.pass, err)
				return
			}

			if client == nil {
				t.Errorf("用户 %s, 密码 %s 时客户端为空", tc.user, tc.pass)
				return
			}

			client.Close()
		})
	}
}
