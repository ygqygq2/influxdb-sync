package influxdb1

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type Client struct {
	cli client.Client
}

func NewClient(addr, user, pass string, timeout time.Duration) (*Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: pass,
		Timeout:  timeout,
	})
	if err != nil {
		return nil, err
	}
	return &Client{cli: c}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}
