package tchttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// Do 发起请求，封装了通用的逻辑
func Do(req *http.Request) ([]byte, error) {
	response, err := GetClient().Do(req)
	if err != nil {
		return nil, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// 防止OOM(内存溢出)
	const maxBodySize = 32 << 20 // 32 MB
	respBytes, err := io.ReadAll(io.LimitReader(response.Body, maxBodySize))
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	return respBytes, nil
}

// PostJSON 快速发起 JSON POST 请求
// body: 序列化后的请求体
func PostJSON(ctx context.Context, url string, body []byte) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return Do(req)
}

// Get 快速发起 GET 请求
func Get(ctx context.Context, url string) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return Do(req)
}
