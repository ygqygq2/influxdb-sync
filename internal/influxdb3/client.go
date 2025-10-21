package influxdb3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	client "github.com/influxdata/influxdb1-client/v2"
)

// Client3x InfluxDB 3.x 客户端，支持多种兼容模式
type Client3x struct {
	baseURL    string
	token      string
	database   string
	namespace  string
	httpClient *http.Client
	v1Client   client.Client    // v1 兼容客户端
	v2Client   influxdb2.Client // v2 兼容客户端
	compatMode string           // "v1", "v2", "native"
}

// NewClient3x 创建新的 3.x 客户端
func NewClient3x(config NativeConfig) (*Client3x, error) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	c := &Client3x{
		baseURL:    strings.TrimSuffix(config.URL, "/"),
		token:      config.Token,
		database:   config.Database,
		namespace:  config.Namespace,
		httpClient: httpClient,
		compatMode: "native",
	}

	return c, nil
}

// NewV1CompatClient 创建 v1 兼容模式客户端
func NewV1CompatClient(config V1CompatConfig) (*Client3x, error) {
	// 创建 v1 客户端
	v1Client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Addr,
		Username: config.User,
		Password: config.Pass,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create v1 client: %w", err)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	c := &Client3x{
		baseURL:    strings.TrimSuffix(config.Addr, "/"),
		database:   config.Database,
		httpClient: httpClient,
		v1Client:   v1Client,
		compatMode: "v1",
	}

	return c, nil
}

// NewV2CompatClient 创建 v2 兼容模式客户端
func NewV2CompatClient(config V2CompatConfig) (*Client3x, error) {
	// 创建 v2 客户端
	v2Client := influxdb2.NewClient(config.URL, config.Token)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	c := &Client3x{
		baseURL:    strings.TrimSuffix(config.URL, "/"),
		token:      config.Token,
		database:   config.Database,
		httpClient: httpClient,
		v2Client:   v2Client,
		compatMode: "v2",
	}

	return c, nil
}

// Close 关闭客户端连接
func (c *Client3x) Close() error {
	if c.v1Client != nil {
		return c.v1Client.Close()
	}
	if c.v2Client != nil {
		c.v2Client.Close()
	}
	return nil
}

// Ping 测试连接
func (c *Client3x) Ping(timeout time.Duration) error {
	switch c.compatMode {
	case "v1":
		if c.v1Client != nil {
			_, _, err := c.v1Client.Ping(timeout)
			return err
		}
	case "v2":
		if c.v2Client != nil {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			_, err := c.v2Client.Ping(ctx)
			return err
		}
	case "native":
		// 对于原生 3.x，尝试访问 health 端点
		healthURL := c.baseURL + "/health"
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
		if err != nil {
			return err
		}

		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("ping failed with status: %d", resp.StatusCode)
		}
	}
	return nil
}

// QuerySQL 执行 SQL 查询 (原生 3.x 功能)
func (c *Client3x) QuerySQL(query string) ([]byte, error) {
	if c.compatMode != "native" {
		return nil, fmt.Errorf("SQL queries only supported in native mode")
	}

	queryURL := c.baseURL + "/v1/sql"

	// 构建请求体
	reqBody := map[string]interface{}{
		"query":    query,
		"database": c.database,
	}
	if c.namespace != "" {
		reqBody["namespace"] = c.namespace
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("SQL query failed with status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// QueryInfluxQL 执行 InfluxQL 查询 (v1 兼容)
func (c *Client3x) QueryInfluxQL(query, database string) (*client.Response, error) {
	if c.compatMode != "v1" || c.v1Client == nil {
		return nil, fmt.Errorf("InfluxQL queries only supported in v1 compatibility mode")
	}

	q := client.NewQuery(query, database, "ns")
	return c.v1Client.Query(q)
}

// QueryFlux 执行 Flux 查询 (v2 兼容)
func (c *Client3x) QueryFlux(query, org string) (*api.QueryTableResult, error) {
	if c.compatMode != "v2" || c.v2Client == nil {
		return nil, fmt.Errorf("Flux queries only supported in v2 compatibility mode")
	}

	queryAPI := c.v2Client.QueryAPI(org)
	result, err := queryAPI.Query(context.Background(), query)
	return result, err
}

// WriteLineProtocol 写入 Line Protocol 数据
func (c *Client3x) WriteLineProtocol(data string) error {
	var writeURL string

	switch c.compatMode {
	case "v1":
		// v1 兼容模式
		writeURL = c.baseURL + "/write"
		params := url.Values{}
		params.Set("db", c.database)
		writeURL += "?" + params.Encode()

	case "v2":
		// v2 兼容模式
		writeURL = c.baseURL + "/api/v2/write"
		params := url.Values{}
		params.Set("bucket", c.database)
		params.Set("precision", "ns")
		writeURL += "?" + params.Encode()

	case "native":
		// 原生 3.x 模式
		writeURL = c.baseURL + "/v1/write"
		params := url.Values{}
		params.Set("database", c.database)
		if c.namespace != "" {
			params.Set("namespace", c.namespace)
		}
		params.Set("precision", "ns")
		writeURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("POST", writeURL, strings.NewReader(data))
	if err != nil {
		return err
	}

	// 设置认证头
	switch c.compatMode {
	case "v2", "native":
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("write failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetDatabases 获取数据库列表
func (c *Client3x) GetDatabases() ([]string, error) {
	switch c.compatMode {
	case "v1":
		if c.v1Client == nil {
			return nil, fmt.Errorf("v1 client not available")
		}
		q := client.NewQuery("SHOW DATABASES", "", "")
		resp, err := c.v1Client.Query(q)
		if err != nil {
			return nil, err
		}
		if resp.Error() != nil {
			return nil, resp.Error()
		}

		var databases []string
		if len(resp.Results) > 0 && len(resp.Results[0].Series) > 0 {
			for _, row := range resp.Results[0].Series[0].Values {
				if len(row) > 0 {
					if dbName, ok := row[0].(string); ok {
						databases = append(databases, dbName)
					}
				}
			}
		}
		return databases, nil

	case "v2":
		// v2 兼容模式返回当前 bucket
		if c.database != "" {
			return []string{c.database}, nil
		}
		return []string{}, nil

	case "native":
		// 原生 3.x 模式，返回当前数据库
		if c.database != "" {
			return []string{c.database}, nil
		}
		return []string{}, nil
	}

	return []string{}, nil
}
