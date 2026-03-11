package tchttp

import (
	"net/http"
	"net/url"
	"time"
)

// Option 定义配置函数
type Option func(*http.Transport, *http.Client)

func WithRequestTimeout(timeout time.Duration) Option {
	return func(tr *http.Transport, c *http.Client) {
		c.Timeout = timeout
	}
}

func WithMaxIdleConnsPerHost(n int) Option {
	return func(tr *http.Transport, c *http.Client) {
		tr.MaxIdleConnsPerHost = n
	}
}

func WithMaxIdleConns(n int) Option {
	return func(tr *http.Transport, c *http.Client) {
		tr.MaxIdleConns = n
	}
}

func WithIdleConnTimeout(seconds int) Option {
	return func(tr *http.Transport, c *http.Client) {
		tr.IdleConnTimeout = time.Duration(seconds) * time.Second
	}
}

func WithProxy(proxy *url.URL) Option {
	return func(tr *http.Transport, c *http.Client) {
		tr.Proxy = http.ProxyURL(proxy)
	}
}
