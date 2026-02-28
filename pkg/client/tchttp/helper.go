package tchttp

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

// Do 发起请求，封装了通用的逻辑
func Do(req *http.Request) ([]byte, error) {
	response, err := DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() // 必须 Close，否则连接无法回收到池中复用

	respBytes, _ := io.ReadAll(response.Body)

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
