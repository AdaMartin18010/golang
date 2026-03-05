# Go应用部署与运维指南

> 从构建到生产环境的完整部署实践

---

## 一、构建优化

### 1.1 编译优化

```text
构建标志：
────────────────────────────────────────

生产构建常用标志：

# 基础优化
go build -ldflags="-s -w" -o app

-s：去除符号表（减小体积）
-w：去除调试信息（减小体积）

# 完全优化
go build \
  -ldflags="-s -w -X main.version=1.0.0" \
  -trimpath \
  -o app

-trimpath：去除构建路径信息
-X：注入版本信息

CGO考虑：
────────────────────────────────────────

禁用CGO（静态链接）：
CGO_ENABLED=0 go build -o app

优势：
- 静态链接，无依赖
- 可以在scratch镜像中运行
- 更好的可移植性

需要CGO时：
- 使用sqlite等C库
- 需要系统库功能
- 使用-musl工具链

多平台构建：
────────────────────────────────────────

交叉编译：

# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o app-linux

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o app-linux-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o app.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o app-darwin
GOOS=darwin GOARCH=arm64 go build -o app-darwin-arm64

使用GoReleaser自动化：
# .goreleaser.yml
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
```

### 1.2 容器化

```text
多阶段构建：
────────────────────────────────────────

# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]

scratch镜像：
────────────────────────────────────────

# 最小镜像
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main .

FROM scratch
COPY --from=builder /app/main /main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["/main"]

大小对比：
┌────────────────┬─────────────┐
│    镜像类型     │    大小     │
├────────────────┼─────────────┤
│ golang:alpine  │   ~350MB    │
│ alpine         │   ~15MB     │
│ scratch        │   ~5MB      │
└────────────────┴─────────────┘

健康检查：
────────────────────────────────────────

Dockerfile：
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

应用实现：
func healthHandler(w http.ResponseWriter, r *http.Request) {
    // 检查依赖
    if err := checkDatabase(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

---

## 二、配置管理

### 2.1 配置来源优先级

```text
12-Factor App配置原则：
────────────────────────────────────────

配置与代码分离：
- 不在代码中硬编码配置
- 通过环境变量传递配置
- 不同环境不同配置

配置优先级（从高到低）：
1. 环境变量
2. 配置文件
3. 默认值

实现模式：
────────────────────────────────────────

type Config struct {
    ServerPort  string `env:"PORT" default:"8080"`
    DatabaseURL string `env:"DATABASE_URL" required:"true"`
    LogLevel    string `env:"LOG_LEVEL" default:"info"`
    Debug       bool   `env:"DEBUG" default:"false"`
}

func LoadConfig() (*Config, error) {
    var cfg Config

    // 设置默认值
    if err := defaults.Set(&cfg); err != nil {
        return nil, err
    }

    // 从环境变量读取
    if err := env.Parse(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

使用：
func main() {
    cfg, err := LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    // 使用配置
    srv := &http.Server{
        Addr: ":" + cfg.ServerPort,
    }
}

敏感信息处理：
────────────────────────────────────────

不要：
- 将密码提交到代码仓库
- 在日志中打印敏感信息
- 通过命令行参数传递密码

应该：
- 使用环境变量或secrets管理
- 日志脱敏
- 使用secret管理工具

Kubernetes secrets：
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  DATABASE_URL: <base64-encoded>
  API_KEY: <base64-encoded>
```

### 2.2 运行时配置

```text
热加载配置：
────────────────────────────────────────

使用fsnotify监听配置变化：

type ConfigManager struct {
    config atomic.Value
}

func (cm *ConfigManager) Load(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return err
    }

    cm.config.Store(&cfg)
    return nil
}

func (cm *ConfigManager) Watch(path string) {
    watcher, _ := fsnotify.NewWatcher()
    watcher.Add(path)

    go func() {
        for event := range watcher.Events {
            if event.Op&fsnotify.Write == fsnotify.Write {
                cm.Load(path)
            }
        }
    }()
}

func (cm *ConfigManager) Get() *Config {
    return cm.config.Load().(*Config)
}

功能开关：
────────────────────────────────────────

使用feature flag：

type FeatureFlags struct {
    NewAlgorithm bool `env:"FEATURE_NEW_ALGORITHM"`
    BetaFeature  bool `env:"FEATURE_BETA"`
}

func processData(data []byte) {
    if flags.NewAlgorithm {
        processWithNewAlgorithm(data)
    } else {
        processWithOldAlgorithm(data)
    }
}

好处：
- 灰度发布
- 快速回滚
- A/B测试
```

---

## 三、可观测性

### 3.1 日志

```text
结构化日志：
────────────────────────────────────────

使用标准库log/slog：

import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// 结构化日志
logger.Info("request processed",
    "method", r.Method,
    "path", r.URL.Path,
    "duration", time.Since(start),
    "status", statusCode,
)

输出：
{
    "time": "2024-01-15T10:30:00Z",
    "level": "INFO",
    "msg": "request processed",
    "method": "GET",
    "path": "/api/users",
    "duration": "45ms",
    "status": 200
}

日志级别：
────────────────────────────────────────

slog.Debug("debug information")      // 开发调试
slog.Info("normal operation")         // 正常运行
slog.Warn("something unexpected")     // 需要注意
slog.Error("operation failed")        // 错误发生

生产环境配置：
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,  // 只记录Info及以上
}))

上下文日志：
────────────────────────────────────────

添加请求ID：
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := generateRequestID()
        ctx := context.WithValue(r.Context(), "requestID", requestID)

        logger := slog.With("requestID", requestID)
        ctx = WithLogger(ctx, logger)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    logger := LoggerFromContext(r.Context())
    logger.Info("handling request")
}
```

### 3.2 指标

```text
使用Prometheus：
────────────────────────────────────────

import "github.com/prometheus/client_golang/prometheus"

var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

func init() {
    prometheus.MustRegister(requestsTotal)
    prometheus.MustRegister(requestDuration)
}

func instrumentHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        next.ServeHTTP(wrapped, r)

        duration := time.Since(start).Seconds()
        requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
        requestsTotal.WithLabelValues(r.Method, r.URL.Path,
            strconv.Itoa(wrapped.statusCode)).Inc()
    })
}

暴露端点：
http.Handle("/metrics", promhttp.Handler())

业务指标：
────────────────────────────────────────

var (
    ordersProcessed = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "orders_processed_total",
        Help: "Total orders processed",
    })

    orderValue = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name:    "order_value_usd",
        Help:    "Order value distribution",
        Buckets: []float64{10, 50, 100, 500, 1000},
    })
)

func ProcessOrder(order Order) error {
    // 处理订单

    ordersProcessed.Inc()
    orderValue.Observe(order.Total)

    return nil
}
```

### 3.3 分布式追踪

```text
使用OpenTelemetry：
────────────────────────────────────────

import "go.opentelemetry.io/otel"

func initTracer() func() {
    exp, _ := jaeger.New(jaeger.WithAgentEndpoint())

    tp := trace.NewTracerProvider(
        trace.WithBatcher(exp),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("my-service"),
        )),
    )

    otel.SetTracerProvider(tp)

    return func() {
        tp.Shutdown(context.Background())
    }
}

创建span：
────────────────────────────────────────

func processOrder(ctx context.Context, order Order) error {
    ctx, span := otel.Tracer("order-service").Start(ctx, "processOrder")
    defer span.End()

    span.SetAttributes(
        attribute.Int("order.id", order.ID),
        attribute.Float64("order.total", order.Total),
    )

    if err := validateOrder(ctx, order); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return err
    }

    return saveOrder(ctx, order)
}

传播追踪上下文：
────────────────────────────────────────

HTTP客户端：
func callService(ctx context.Context, url string) error {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    // 注入追踪信息
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

    resp, err := http.DefaultClient.Do(req)
    // ...
}

HTTP服务端：
func handler(w http.ResponseWriter, r *http.Request) {
    // 提取追踪信息
    ctx := otel.GetTextMapPropagator().Extract(r.Context(),
        propagation.HeaderCarrier(r.Header))

    ctx, span := tracer.Start(ctx, "handler")
    defer span.End()

    // 处理请求...
}
```

---

*本章提供了Go应用从构建到生产部署的完整指南。*
