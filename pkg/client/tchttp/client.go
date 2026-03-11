package tchttp

import (
	"net/http"
	"sync"
	"time"
)

var clientOnce sync.Once
var httpClient *http.Client

func GetClient() *http.Client {
	clientOnce.Do(func() {
		httpClient = NewClient()
	})
	return httpClient
}

func SetClient(c *http.Client) {
	httpClient = c
}

// NewClient 创建一个新的、经过配置的 http.Client
func NewClient(opts ...Option) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second, // 整个请求的绝对超时时间（含读取 Body）
	}

	// 应用自定义配置
	for _, opt := range opts {
		if opt != nil {
			opt(tr, c)
		}
	}

	return c
}
