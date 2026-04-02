# API 客户端设计 (API Client Design)

> **分类**: 开源技术堆栈  
> **标签**: #api-client #http #sdk

---

## 客户端结构

```go
type Client struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
    retryPolicy RetryPolicy
    
    // 子服务客户端
    Users  *UserService
    Orders *OrderService
}

func NewClient(baseURL, apiKey string, opts ...ClientOption) *Client {
    c := &Client{
        baseURL: baseURL,
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        retryPolicy: DefaultRetryPolicy,
    }
    
    // 应用选项
    for _, opt := range opts {
        opt(c)
    }
    
    // 初始化子服务
    c.Users = &UserService{client: c}
    c.Orders = &OrderService{client: c}
    
    return c
}

type ClientOption func(*Client)

func WithHTTPClient(hc *http.Client) ClientOption {
    return func(c *Client) {
        c.httpClient = hc
    }
}

func WithRetryPolicy(rp RetryPolicy) ClientOption {
    return func(c *Client) {
        c.retryPolicy = rp
    }
}
```

---

## 请求构建

```go
type RequestBuilder struct {
    method  string
    path    string
    headers map[string]string
    query   url.Values
    body    interface{}
}

func (c *Client) NewRequest(method, path string) *RequestBuilder {
    return &RequestBuilder{
        method:  method,
        path:    path,
        headers: make(map[string]string),
        query:   make(url.Values),
    }
}

func (rb *RequestBuilder) WithQuery(key, value string) *RequestBuilder {
    rb.query.Set(key, value)
    return rb
}

func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
    rb.headers[key] = value
    return rb
}

func (rb *RequestBuilder) WithBody(body interface{}) *RequestBuilder {
    rb.body = body
    return rb
}

func (rb *RequestBuilder) Execute(ctx context.Context, result interface{}) error {
    // 构建 URL
    u, _ := url.Parse(rb.path)
    u.RawQuery = rb.query.Encode()
    
    // 序列化 body
    var bodyReader io.Reader
    if rb.body != nil {
        data, _ := json.Marshal(rb.body)
        bodyReader = bytes.NewReader(data)
    }
    
    // 创建请求
    req, err := http.NewRequestWithContext(ctx, rb.method, u.String(), bodyReader)
    if err != nil {
        return err
    }
    
    // 设置 headers
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")
    for k, v := range rb.headers {
        req.Header.Set(k, v)
    }
    
    // 执行
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // 处理响应
    if resp.StatusCode >= 400 {
        return c.handleError(resp)
    }
    
    if result != nil {
        return json.NewDecoder(resp.Body).Decode(result)
    }
    
    return nil
}
```

---

## 使用示例

```go
client := NewClient("https://api.example.com", "api-key",
    WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
)

// 使用
user, err := client.Users.Get(ctx, "user-123")

// 创建
newUser, err := client.Users.Create(ctx, &CreateUserRequest{
    Name:  "John",
    Email: "john@example.com",
})

// 子服务
orders, err := client.Orders.List(ctx, ListOrdersRequest{
    UserID: "user-123",
    Status: "pending",
})
```

---

## 错误处理

```go
type APIError struct {
    StatusCode int
    Code       string
    Message    string
    Details    map[string]interface{}
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
}

func (c *Client) handleError(resp *http.Response) error {
    var apiErr APIError
    apiErr.StatusCode = resp.StatusCode
    
    // 尝试解析错误响应
    if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
        apiErr.Message = http.StatusText(resp.StatusCode)
    }
    
    return &apiErr
}
```

---

## 分页处理

```go
type PaginatedResponse struct {
    Data       []User
    NextCursor string
    HasMore    bool
}

func (s *UserService) ListAll(ctx context.Context) ([]User, error) {
    var allUsers []User
    cursor := ""
    
    for {
        resp, err := s.List(ctx, ListRequest{Cursor: cursor, Limit: 100})
        if err != nil {
            return nil, err
        }
        
        allUsers = append(allUsers, resp.Data...)
        
        if !resp.HasMore {
            break
        }
        cursor = resp.NextCursor
    }
    
    return allUsers, nil
}
```
