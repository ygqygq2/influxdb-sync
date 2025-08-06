package common

import (
	"testing"
	"time"
)

func TestSyncConfigStruct(t *testing.T) {
	// 测试SyncConfig结构体的完整性
	cfg := SyncConfig{
		SourceAddr:      "http://source.example.com:8086",
		SourceUser:      "source_admin",
		SourcePass:      "source_password",
		SourceDB:        "source_database",
		SourceDBExclude: []string{"_internal", "system", "temp"},
		SourceToken:     "source-api-token",
		SourceOrg:       "source-organization",
		SourceBucket:    "source-bucket",
		TargetAddr:      "http://target.example.com:8087",
		TargetUser:      "target_admin",
		TargetPass:      "target_password",
		TargetDB:        "target_database",
		TargetDBPrefix:  "migrated_",
		TargetDBSuffix:  "_backup",
		TargetToken:     "target-api-token",
		TargetOrg:       "target-organization",
		TargetBucket:    "target-bucket",
		BatchSize:       2000,
		Start:           "2024-01-01T00:00:00Z",
		End:             "2024-12-31T23:59:59Z",
		ResumeFile:      "/tmp/sync_resume.state",
		Parallel:        8,
		RetryCount:      5,
		RetryInterval:   1000,
		RateLimit:       200,
		LogLevel:        "debug",
	}

	// 验证所有字段
	if cfg.SourceAddr != "http://source.example.com:8086" {
		t.Error("SourceAddr字段设置错误")
	}
	if cfg.SourceUser != "source_admin" {
		t.Error("SourceUser字段设置错误")
	}
	if cfg.SourcePass != "source_password" {
		t.Error("SourcePass字段设置错误")
	}
	if cfg.SourceDB != "source_database" {
		t.Error("SourceDB字段设置错误")
	}
	if len(cfg.SourceDBExclude) != 3 {
		t.Error("SourceDBExclude字段设置错误")
	}
	if cfg.SourceToken != "source-api-token" {
		t.Error("SourceToken字段设置错误")
	}
	if cfg.SourceOrg != "source-organization" {
		t.Error("SourceOrg字段设置错误")
	}
	if cfg.SourceBucket != "source-bucket" {
		t.Error("SourceBucket字段设置错误")
	}

	// 验证目标配置
	if cfg.TargetAddr != "http://target.example.com:8087" {
		t.Error("TargetAddr字段设置错误")
	}
	if cfg.TargetUser != "target_admin" {
		t.Error("TargetUser字段设置错误")
	}
	if cfg.TargetPass != "target_password" {
		t.Error("TargetPass字段设置错误")
	}
	if cfg.TargetDB != "target_database" {
		t.Error("TargetDB字段设置错误")
	}
	if cfg.TargetDBPrefix != "migrated_" {
		t.Error("TargetDBPrefix字段设置错误")
	}
	if cfg.TargetDBSuffix != "_backup" {
		t.Error("TargetDBSuffix字段设置错误")
	}
	if cfg.TargetToken != "target-api-token" {
		t.Error("TargetToken字段设置错误")
	}
	if cfg.TargetOrg != "target-organization" {
		t.Error("TargetOrg字段设置错误")
	}
	if cfg.TargetBucket != "target-bucket" {
		t.Error("TargetBucket字段设置错误")
	}

	// 验证同步配置
	if cfg.BatchSize != 2000 {
		t.Error("BatchSize字段设置错误")
	}
	if cfg.Start != "2024-01-01T00:00:00Z" {
		t.Error("Start字段设置错误")
	}
	if cfg.End != "2024-12-31T23:59:59Z" {
		t.Error("End字段设置错误")
	}
	if cfg.ResumeFile != "/tmp/sync_resume.state" {
		t.Error("ResumeFile字段设置错误")
	}
	if cfg.Parallel != 8 {
		t.Error("Parallel字段设置错误")
	}
	if cfg.RetryCount != 5 {
		t.Error("RetryCount字段设置错误")
	}
	if cfg.RetryInterval != 1000 {
		t.Error("RetryInterval字段设置错误")
	}
	if cfg.RateLimit != 200 {
		t.Error("RateLimit字段设置错误")
	}
	if cfg.LogLevel != "debug" {
		t.Error("LogLevel字段设置错误")
	}
}

func TestDataPointStructFields(t *testing.T) {
	// 测试DataPoint结构体的所有字段
	now := time.Now()
	point := DataPoint{
		Measurement: "system_metrics",
		Tags: map[string]string{
			"host":        "server-01",
			"region":      "us-west-2",
			"environment": "production",
			"service":     "api-gateway",
		},
		Fields: map[string]interface{}{
			"cpu_usage":          85.5,
			"memory_usage":       72.3,
			"disk_usage":         45.1,
			"load_avg":           2.8,
			"active_connections": 150,
			"response_time":      250.5,
			"error_rate":         0.02,
			"throughput":         1250.7,
		},
		Time: now,
	}

	// 验证测量名
	if point.Measurement != "system_metrics" {
		t.Error("Measurement字段设置错误")
	}

	// 验证标签
	if len(point.Tags) != 4 {
		t.Error("Tags数量错误")
	}
	if point.Tags["host"] != "server-01" {
		t.Error("host标签设置错误")
	}
	if point.Tags["region"] != "us-west-2" {
		t.Error("region标签设置错误")
	}
	if point.Tags["environment"] != "production" {
		t.Error("environment标签设置错误")
	}
	if point.Tags["service"] != "api-gateway" {
		t.Error("service标签设置错误")
	}

	// 验证字段
	if len(point.Fields) != 8 {
		t.Error("Fields数量错误")
	}
	if point.Fields["cpu_usage"] != 85.5 {
		t.Error("cpu_usage字段设置错误")
	}
	if point.Fields["memory_usage"] != 72.3 {
		t.Error("memory_usage字段设置错误")
	}
	if point.Fields["active_connections"] != 150 {
		t.Error("active_connections字段设置错误")
	}

	// 验证时间戳
	if !point.Time.Equal(now) {
		t.Error("Time字段设置错误")
	}
}

func TestSyncResultStructVariations(t *testing.T) {
	// 测试SyncResult结构体的不同情况
	testCases := []struct {
		name        string
		measurement string
		error       error
		expectError bool
	}{
		{
			name:        "成功结果",
			measurement: "cpu_metrics",
			error:       nil,
			expectError: false,
		},
		{
			name:        "失败结果",
			measurement: "memory_metrics",
			error:       &mockError{"同步失败"},
			expectError: true,
		},
		{
			name:        "空测量名成功",
			measurement: "",
			error:       nil,
			expectError: false,
		},
		{
			name:        "长测量名失败",
			measurement: "very_long_measurement_name_that_might_cause_issues_in_some_systems",
			error:       &mockError{"测量名过长"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SyncResult{
				Measurement: tc.measurement,
				Error:       tc.error,
			}

			if result.Measurement != tc.measurement {
				t.Errorf("Measurement期望 %s, 实际 %s", tc.measurement, result.Measurement)
			}

			hasError := result.Error != nil
			if hasError != tc.expectError {
				t.Errorf("Error状态期望 %v, 实际 %v", tc.expectError, hasError)
			}

			if tc.expectError && result.Error.Error() == "" {
				t.Error("错误消息不应该为空")
			}
		})
	}
}

func TestDataPointWithDifferentFieldTypes(t *testing.T) {
	// 测试DataPoint支持不同类型的字段值
	point := DataPoint{
		Measurement: "mixed_types",
		Tags: map[string]string{
			"location": "datacenter-1",
		},
		Fields: map[string]interface{}{
			"string_field":  "text_value",
			"int_field":     42,
			"int64_field":   int64(9223372036854775807),
			"float32_field": float32(3.14),
			"float64_field": 3.141592653589793,
			"bool_field":    true,
			"nil_field":     nil,
		},
		Time: time.Now(),
	}

	// 验证不同类型的字段
	if point.Fields["string_field"] != "text_value" {
		t.Error("string_field类型错误")
	}
	if point.Fields["int_field"] != 42 {
		t.Error("int_field类型错误")
	}
	if point.Fields["int64_field"] != int64(9223372036854775807) {
		t.Error("int64_field类型错误")
	}
	if point.Fields["bool_field"] != true {
		t.Error("bool_field类型错误")
	}
	if point.Fields["nil_field"] != nil {
		t.Error("nil_field应该为nil")
	}

	// 验证类型断言
	if _, ok := point.Fields["string_field"].(string); !ok {
		t.Error("string_field应该是string类型")
	}
	if _, ok := point.Fields["int_field"].(int); !ok {
		t.Error("int_field应该是int类型")
	}
	if _, ok := point.Fields["bool_field"].(bool); !ok {
		t.Error("bool_field应该是bool类型")
	}
}
