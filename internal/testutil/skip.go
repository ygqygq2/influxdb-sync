package testutil

import (
	"os"
	"testing"
)

// SkipIfNoDatabase 如果没有设置 INTEGRATION_TEST 环境变量，则跳过测试
// 用法: testutil.SkipIfNoDatabase(t)
func SkipIfNoDatabase(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("跳过集成测试，设置 INTEGRATION_TEST=1 来运行")
	}
}

// SkipIfNoDatabaseWithMsg 跳过测试并显示自定义消息
func SkipIfNoDatabaseWithMsg(t *testing.T, msg string) {
	t.Helper()
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip(msg)
	}
}
