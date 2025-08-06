package logx

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestSetLevel(t *testing.T) {
	// 保存原始级别
	originalLevel := currentLevel
	defer func() {
		currentLevel = originalLevel
	}()

	// 测试设置不同级别
	testCases := []string{"debug", "info", "warn", "error"}
	
	for _, level := range testCases {
		SetLevel(level)
		if currentLevel != level {
			t.Errorf("期望当前级别为 %s, 实际为 %s", level, currentLevel)
		}
	}

	// 测试大写级别
	SetLevel("INFO")
	if currentLevel != "info" {
		t.Errorf("期望当前级别为 info, 实际为 %s", currentLevel)
	}
}

func TestShouldLog(t *testing.T) {
	// 保存原始级别
	originalLevel := currentLevel
	defer func() {
		currentLevel = originalLevel
	}()

	testCases := []struct {
		currentLevel string
		targetLevel  string
		expected     bool
	}{
		{"debug", "debug", true},
		{"debug", "info", true},
		{"debug", "warn", true},
		{"debug", "error", true},
		{"info", "debug", false},
		{"info", "info", true},
		{"info", "warn", true},
		{"info", "error", true},
		{"warn", "debug", false},
		{"warn", "info", false},
		{"warn", "warn", true},
		{"warn", "error", true},
		{"error", "debug", false},
		{"error", "info", false},
		{"error", "warn", false},
		{"error", "error", true},
		{"invalid", "info", true}, // 无效级别应该默认输出
		{"info", "invalid", true}, // 无效级别应该默认输出
	}

	for _, tc := range testCases {
		currentLevel = tc.currentLevel
		result := shouldLog(tc.targetLevel)
		if result != tc.expected {
			t.Errorf("当前级别 %s, 目标级别 %s, 期望 %v, 实际 %v", 
				tc.currentLevel, tc.targetLevel, tc.expected, result)
		}
	}
}

func TestLogFunctions(t *testing.T) {
	// 保存原始logger和级别
	originalLogger := logger
	originalLevel := currentLevel
	defer func() {
		logger = originalLogger
		currentLevel = originalLevel
	}()

	// 创建缓冲区用于捕获日志输出
	var buf bytes.Buffer
	logger = log.New(&buf, "", 0)
	
	// 设置为debug级别，这样所有级别都会输出
	SetLevel("debug")

	// 测试Debug
	Debug("debug message")
	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Debug消息未正确输出")
	}
	if !strings.Contains(buf.String(), "[DEBUG]") {
		t.Error("Debug前缀未正确设置")
	}

	// 清空缓冲区
	buf.Reset()

	// 测试Info
	Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info消息未正确输出")
	}
	if !strings.Contains(buf.String(), "[INFO]") {
		t.Error("Info前缀未正确设置")
	}

	// 清空缓冲区
	buf.Reset()

	// 测试Warn
	Warn("warn message")
	if !strings.Contains(buf.String(), "warn message") {
		t.Error("Warn消息未正确输出")
	}
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Error("Warn前缀未正确设置")
	}

	// 清空缓冲区
	buf.Reset()

	// 测试Error
	Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("Error消息未正确输出")
	}
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Error("Error前缀未正确设置")
	}
}

func TestLogLevelFiltering(t *testing.T) {
	// 保存原始logger和级别
	originalLogger := logger
	originalLevel := currentLevel
	defer func() {
		logger = originalLogger
		currentLevel = originalLevel
	}()

	// 创建缓冲区用于捕获日志输出
	var buf bytes.Buffer
	logger = log.New(&buf, "", 0)
	
	// 设置为warn级别
	SetLevel("warn")

	// 测试Debug（不应该输出）
	Debug("debug message")
	if strings.Contains(buf.String(), "debug message") {
		t.Error("Debug消息不应该在warn级别输出")
	}

	// 测试Info（应该输出，因为Info函数没有级别检查）
	Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info消息应该输出（Info函数没有级别过滤）")
	}

	// 测试Warn（应该输出）
	Warn("warn message")
	if !strings.Contains(buf.String(), "warn message") {
		t.Error("Warn消息应该在warn级别输出")
	}

	// 测试Error（应该输出）
	Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("Error消息应该在warn级别输出")
	}
}

func TestSprint(t *testing.T) {
	result := sprint("test", " ", "message", " ", 123)
	expected := "test message 123"
	if result != expected {
		t.Errorf("期望 %s, 实际 %s", expected, result)
	}
}
