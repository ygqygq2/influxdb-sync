package influxdb1

import (
	"testing"
)

// 测试 NewDataSource 函数
func TestNewDataSource_Creation(t *testing.T) {
	cfg := DataSourceConfig{
		Addr: "http://localhost:8086",
		User: "admin",
		Pass: "password",
	}

	ds := NewDataSource(cfg)
	if ds == nil {
		t.Error("NewDataSource 返回 nil")
	}
}

// 测试 NewDataTarget 函数
func TestNewDataTarget_Creation(t *testing.T) {
	cfg := DataTargetConfig{
		Addr: "http://localhost:8086",
		User: "admin",
		Pass: "password",
	}

	dt := NewDataTarget(cfg)
	if dt == nil {
		t.Error("NewDataTarget 返回 nil")
	}
}
