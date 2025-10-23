package testutil

import (
"context"
"net/http"
"os"
"testing"
"time"

influxdb1 "github.com/influxdata/influxdb1-client/v2"
)

// SkipIfNoInfluxDB1 检查 InfluxDB 1.x 是否可用
func SkipIfNoInfluxDB1(t *testing.T, addr string) {
t.Helper()
if os.Getenv("SKIP_INTEGRATION") == "true" {
t.Skip("跳过集成测试 (SKIP_INTEGRATION=true)")
}

cli, err := influxdb1.NewHTTPClient(influxdb1.HTTPConfig{
Addr:    addr,
Timeout: 2 * time.Second,
})
if err != nil {
t.Skipf("无法创建 InfluxDB 1.x 客户端: %v", err)
return
}
defer cli.Close()

_, _, err = cli.Ping(2 * time.Second)
if err != nil {
t.Skipf("InfluxDB 1.x 不可用 (%s): %v", addr, err)
}
}

// SkipIfNoInfluxDB2 简单的健康检查 InfluxDB 2.x
func SkipIfNoInfluxDB2(t *testing.T, url string) {
t.Helper()
if os.Getenv("SKIP_INTEGRATION") == "true" {
t.Skip("跳过集成测试 (SKIP_INTEGRATION=true)")
return
}

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
req, err := http.NewRequestWithContext(ctx, "GET", url+"/health", nil)
if err != nil {
t.Skipf("无法创建请求: %v", err)
return
}
resp, err := http.DefaultClient.Do(req)
if err != nil || resp.StatusCode != http.StatusOK {
t.Skipf("InfluxDB 2.x 不可用 (%s): %v", url, err)
return
}
defer resp.Body.Close()
}

// SkipIfNoInfluxDB3 简单的健康检查 InfluxDB 3.x （与 2.x 类似）
func SkipIfNoInfluxDB3(t *testing.T, url string) {
SkipIfNoInfluxDB2(t, url)
}

// MustHaveDocker 确保 Docker 环境可用（简化实现）
func MustHaveDocker(t *testing.T) {
t.Helper()
if os.Getenv("SKIP_INTEGRATION") == "true" {
t.Skip("跳过集成测试 (SKIP_INTEGRATION=true)")
}
}
