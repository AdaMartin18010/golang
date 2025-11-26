package httpclient

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
)

// Client HTTP客户端
type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string
	timeout    time.Duration
}

// Config 客户端配置
type Config struct {
	BaseURL    string
	Timeout    time.Duration
	Headers    map[string]string
	Transport  *http.Transport
}

// NewClient 创建HTTP客户端
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	if config.Transport != nil {
		httpClient.Transport = config.Transport
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		headers:    config.Headers,
		timeout:    config.Timeout,
	}
}

// SetHeader 设置请求头
func (c *Client) SetHeader(key, value string) {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[key] = value
}

// SetHeaders 批量设置请求头
func (c *Client) SetHeaders(headers map[string]string) {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	for k, v := range headers {
		c.headers[k] = v
	}
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, path string, params map[string]string) (*Response, error) {
	return c.Request(ctx, "GET", path, params, nil, nil)
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Request(ctx, "POST", path, nil, body, headers)
}

// Put 发送PUT请求
func (c *Client) Put(ctx context.Context, path string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Request(ctx, "PUT", path, nil, body, headers)
}

// Delete 发送DELETE请求
func (c *Client) Delete(ctx context.Context, path string, params map[string]string) (*Response, error) {
	return c.Request(ctx, "DELETE", path, params, nil, nil)
}

// Patch 发送PATCH请求
func (c *Client) Patch(ctx context.Context, path string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Request(ctx, "PATCH", path, nil, body, headers)
}

// Request 发送HTTP请求
func (c *Client) Request(
	ctx context.Context,
	method, path string,
	params map[string]string,
	body interface{},
	headers map[string]string,
) (*Response, error) {
	// 构建URL
	fullURL := c.buildURL(path, params)

	// 构建请求体
	var reqBody io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 设置请求头
	c.setHeaders(req, headers)

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
	}, nil
}

// buildURL 构建完整URL
func (c *Client) buildURL(path string, params map[string]string) string {
	var fullURL string
	if c.baseURL != "" {
		fullURL = strings.TrimSuffix(c.baseURL, "/") + "/" + strings.TrimPrefix(path, "/")
	} else {
		fullURL = path
	}

	if len(params) > 0 {
		u, err := url.Parse(fullURL)
		if err == nil {
			q := u.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			u.RawQuery = q.Encode()
			fullURL = u.String()
		}
	}

	return fullURL
}

// setHeaders 设置请求头
func (c *Client) setHeaders(req *http.Request, extraHeaders map[string]string) {
	// 设置默认请求头
	if c.headers != nil {
		for k, v := range c.headers {
			req.Header.Set(k, v)
		}
	}

	// 设置额外请求头
	if extraHeaders != nil {
		for k, v := range extraHeaders {
			req.Header.Set(k, v)
		}
	}

	// 如果请求体是JSON，设置Content-Type
	if req.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
}

// Response HTTP响应
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// JSON 将响应体解析为JSON
func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

// String 返回响应体字符串
func (r *Response) String() string {
	return string(r.Body)
}

// IsSuccess 判断是否成功（2xx状态码）
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError 判断是否错误（4xx或5xx状态码）
func (r *Response) IsError() bool {
	return r.StatusCode >= 400
}

// Get 获取响应头值
func (r *Response) GetHeader(key string) string {
	return r.Headers.Get(key)
}

// DefaultClient 默认HTTP客户端
var DefaultClient = NewClient(Config{
	Timeout: 30 * time.Second,
})

// Get 使用默认客户端发送GET请求
func Get(ctx context.Context, url string, params map[string]string) (*Response, error) {
	return DefaultClient.Get(ctx, url, params)
}

// Post 使用默认客户端发送POST请求
func Post(ctx context.Context, url string, body interface{}) (*Response, error) {
	return DefaultClient.Post(ctx, url, body, nil)
}

// Put 使用默认客户端发送PUT请求
func Put(ctx context.Context, url string, body interface{}) (*Response, error) {
	return DefaultClient.Put(ctx, url, body, nil)
}

// Delete 使用默认客户端发送DELETE请求
func Delete(ctx context.Context, url string, params map[string]string) (*Response, error) {
	return DefaultClient.Delete(ctx, url, params)
}
