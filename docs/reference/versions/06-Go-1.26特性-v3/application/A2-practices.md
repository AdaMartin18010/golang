# A2: 工程实践

> **层级**: 应用层 (Application)
> **地位**: 基于模式的工程实践
> **依赖**: A1

---

## 实践 1: 配置管理

```go
package config

// AppConfig 应用配置
type AppConfig struct {
    // 服务器配置
    Server struct {
        Host    *string
        Port    *int
        Timeout *time.Duration
    }

    // 数据库配置
    Database struct {
        URL      *string
        MaxConns *int
        Timeout  *time.Duration
    }

    // 缓存配置
    Cache struct {
        Enabled  *bool
        Size     *int
        TTL      *time.Duration
    }
}

// LoadDefaults 加载默认值
func (c *AppConfig) LoadDefaults() {
    // 服务器默认值
    if c.Server.Host == nil {
        c.Server.Host = new("0.0.0.0")
    }
    if c.Server.Port == nil {
        c.Server.Port = new(8080)
    }
    if c.Server.Timeout == nil {
        c.Server.Timeout = new(30 * time.Second)
    }

    // 数据库默认值
    if c.Database.URL == nil {
        c.Database.URL = new("postgres://localhost/myapp")
    }
    if c.Database.MaxConns == nil {
        c.Database.MaxConns = new(100)
    }
    if c.Database.Timeout == nil {
        c.Database.Timeout = new(10 * time.Second)
    }

    // 缓存默认值
    if c.Cache.Enabled == nil {
        c.Cache.Enabled = new(true)
    }
    if c.Cache.Size == nil {
        c.Cache.Size = new(1000)
    }
    if c.Cache.TTL == nil {
        c.Cache.TTL = new(5 * time.Minute)
    }
}

// 从环境变量加载
func (c *AppConfig) LoadFromEnv() {
    if host := os.Getenv("SERVER_HOST"); host != "" {
        c.Server.Host = new(host)
    }
    if port := os.Getenv("SERVER_PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            c.Server.Port = new(p)
        }
    }
    // ... 其他环境变量
}

// 从配置文件加载
func LoadConfig(path string) (*AppConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config AppConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    config.LoadDefaults()
    config.LoadFromEnv()

    return &config, nil
}
```

---

## 实践 2: API 设计

```go
package api

// Client 配置选项
type ClientOptions struct {
    BaseURL    *string
    Timeout    *time.Duration
    RetryCount *int
    APIKey     *string
}

type Option func(*ClientOptions)

// 选项函数使用 new 表达式
func WithBaseURL(url string) Option {
    return func(o *ClientOptions) {
        o.BaseURL = new(url)
    }
}

func WithTimeout(d time.Duration) Option {
    return func(o *ClientOptions) {
        o.Timeout = new(d)
    }
}

func WithRetryCount(n int) Option {
    return func(o *ClientOptions) {
        o.RetryCount = new(n)
    }
}

func WithAPIKey(key string) Option {
    return func(o *ClientOptions) {
        o.APIKey = new(key)
    }
}

type Client struct {
    opts ClientOptions
}

func NewClient(options ...Option) *Client {
    c := &Client{
        opts: ClientOptions{
            BaseURL:    new("https://api.example.com"),
            Timeout:    new(30 * time.Second),
            RetryCount: new(3),
        },
    }

    for _, opt := range options {
        opt(&c.opts)
    }

    return c
}

// 使用
client := api.NewClient(
    api.WithBaseURL("https://api.service.com"),
    api.WithTimeout(10*time.Second),
    api.WithRetryCount(5),
)
```

---

## 实践 3: 延迟初始化

```go
package lazy

// Value 延迟加载的值
type Value[T any] struct {
    once  sync.Once
    value T
    fn    func() T
}

func NewValue[T any](fn func() T) *Value[T] {
    return &Value[T]{fn: fn}
}

func (v *Value[T]) Get() T {
    v.once.Do(func() {
        v.value = v.fn()
    })
    return v.value
}

// 使用示例
type Database struct {
    connection *Value[*sql.DB]
}

func NewDatabase(connString string) *Database {
    return &Database{
        connection: NewValue(func() *sql.DB {
            db, err := sql.Open("postgres", connString)
            if err != nil {
                panic(err)
            }
            return db
        }),
    }
}

func (d *Database) Query(ctx context.Context, query string) (*sql.Rows, error) {
    return d.connection.Get().QueryContext(ctx, query)
}
```

---

**下一章**: [A3-性能优化](A3-optimization.md)
