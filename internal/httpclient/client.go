package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/damola-oladipo/probe/internal/request"
)

const maxBodyBytes = 8 << 20 // 8MiB

type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       string
	Duration   time.Duration
}

type Client struct {
	client *http.Client
}

func New() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Send(ctx context.Context, req request.Request) (Response, error) {
	start := time.Now()

	var body io.Reader
	if req.Body != "" && req.Method != http.MethodHead {
		body = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, body)
	if err != nil {
		return Response{}, err
	}

	for key, values := range req.Headers {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return Response{}, err
	}
	defer httpResp.Body.Close()

	limited := io.LimitReader(httpResp.Body, int64(maxBodyBytes)+1)
	raw, err := io.ReadAll(limited)
	if err != nil {
		return Response{
			StatusCode: httpResp.StatusCode,
			Status:     httpResp.Status,
			Headers:    httpResp.Header,
			Duration:   time.Since(start),
		}, err
	}

	out := Response{
		StatusCode: httpResp.StatusCode,
		Status:     httpResp.Status,
		Headers:    httpResp.Header,
		Duration:   time.Since(start),
	}
	if len(raw) > maxBodyBytes {
		return out, fmt.Errorf("response body exceeds 8MiB")
	}
	out.Body = string(raw)
	return out, nil
}
