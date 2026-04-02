# 设计模式 (Design Patterns)

> **分类**: 工程与云原生

---

## 创建型模式

### 单例

```go
type Database struct {
    conn *sql.DB
}

var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        instance = &Database{}
        instance.conn = createConnection()
    })
    return instance
}
```

### 工厂

```go
func NewLogger(logType string) (Logger, error) {
    switch logType {
    case "file":
        return &FileLogger{}, nil
    case "console":
        return &ConsoleLogger{}, nil
    default:
        return nil, errors.New("unknown type")
    }
}
```

---

## 结构型模式

### 适配器

```go
type Adapter struct {
    adaptee *Adaptee
}

func (a *Adapter) Request() string {
    return a.adaptee.SpecificRequest()
}
```

---

## Go 惯用模式

### 函数选项

```go
type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{host: "localhost", port: 8080}
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

### 管道

```go
func Generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}
```
