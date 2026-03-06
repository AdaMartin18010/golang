# 示例项目索引与指南

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.26

---

## 📋 目录

- [示例项目索引与指南](#示例项目索引与指南)
  - [📋 目录](#-目录)
  - [1. 概述](#1-概述)
    - [1.1 示例项目列表](#11-示例项目列表)
  - [2. Web服务示例](#2-web服务示例)
    - [2.1 RESTful API服务](#21-restful-api服务)
  - [3. 微服务示例](#3-微服务示例)
    - [3.1 gRPC微服务](#31-grpc微服务)
  - [4. CLI工具示例](#4-cli工具示例)
    - [4.1 数据分析CLI](#41-数据分析cli)
  - [5. 数据处理示例](#5-数据处理示例)
    - [5.1 实时日志处理](#51-实时日志处理)
  - [6. 实时系统示例](#6-实时系统示例)
    - [6.1 Web爬虫系统](#61-web爬虫系统)
  - [7. 使用指南](#7-使用指南)
    - [7.1 快速开始](#71-快速开始)
    - [7.2 学习路径](#72-学习路径)

---

---

## 1. 概述

### 1.1 示例项目列表

```text
5个完整示例项目:

┌─────────────────────────────────────┐
│         示例项目                    │
├─────────────────────────────────────┤
│                                     │
│  1. RESTful API服务                 │
│     └─ HTTP/3, 认证, CRUD, 测试     │
│                                     │
│  2. 微服务架构                      │
│     └─ gRPC, 服务发现, 追踪, 监控   │
│                                     │
│  3. CLI数据分析工具                 │
│     └─ REPL, 可视化, 插件, 导出     │
│                                     │
│  4. 实时日志处理                    │
│     └─ 流式, 聚合, 告警, 存储       │
│                                     │
│  5. Web爬虫系统                     │
│     └─ 并发, 去重, 存储, 限速       │
│                                     │
└─────────────────────────────────────┘
```

---

## 2. Web服务示例

### 2.1 RESTful API服务

**项目结构**:

```text
restful-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   ├── user.go
│   │   │   └── product.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── logger.go
│   │   │   └── ratelimit.go
│   │   └── router.go
│   ├── model/
│   │   ├── user.go
│   │   └── product.go
│   ├── repository/
│   │   ├── user_repo.go
│   │   └── product_repo.go
│   └── service/
│       ├── user_service.go
│       └── product_service.go
├── pkg/
│   ├── database/
│   │   └── postgres.go
│   └── jwt/
│       └── token.go
├── config/
│   └── config.yaml
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

**核心代码**:

```go
// cmd/server/main.go

package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "restful-api/internal/api"
    "restful-api/pkg/database"
)

func main() {
    // 初始化数据库
    db, err := database.NewPostgres()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 创建路由
    router := api.NewRouter(db)

    // 创建服务器
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    // 启动服务器
    go func() {
        log.Println("Starting server on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // 优雅关闭
    quit := make(Channel os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited")
}
```

```go
// internal/api/handler/user.go

package handler

import (
    "encoding/json"
    "net/http"

    "restful-api/internal/service"
)

type UserHandler struct {
    service *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
    return &UserHandler{service: svc}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user, err := h.service.Create(r.Context(), req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    user, err := h.service.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

type CreateUserRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

**运行说明**:

```bash
# 1. 启动依赖
docker-compose up -d

# 2. 运行迁移
go run cmd/migrate/main.go

# 3. 启动服务
go run cmd/server/main.go

# 4. 测试API
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret"}'
```

---

## 3. 微服务示例

### 3.1 gRPC微服务

```text
microservices/
├── services/
│   ├── user/
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── handler/
│   │   │   └── repository/
│   │   └── proto/
│   │       └── user.proto
│   ├── order/
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── handler/
│   │   │   └── repository/
│   │   └── proto/
│   │       └── order.proto
│   └── product/
│       ├── cmd/main.go
│       ├── internal/
│       └── proto/
│           └── product.proto
├── pkg/
│   ├── discovery/
│   │   └── consul.go
│   ├── tracing/
│   │   └── jaeger.go
│   └── metrics/
│       └── prometheus.go
├── gateway/
│   └── main.go
└── docker-compose.yml
```

**核心代码**:

```go
// services/user/cmd/main.go

package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/health"
    "google.golang.org/grpc/health/grpc_health_v1"
    "google.golang.org/grpc/reflection"

    pb "microservices/services/user/proto"
    "microservices/pkg/discovery"
    "microservices/pkg/tracing"
)

type server struct {
    pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx Context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    // 实现获取用户逻辑
    return &pb.User{
        Id:    req.Id,
        Email: "user@example.com",
        Name:  "John Doe",
    }, nil
}

func main() {
    // 初始化追踪
    tracer, err := tracing.NewJaegerTracer("user-service")
    if err != nil {
        log.Fatal(err)
    }
    defer tracer.Close()

    // 注册服务发现
    consul := discovery.NewConsulClient()
    if err := consul.Register("user-service", "localhost", 50051); err != nil {
        log.Fatal(err)
    }
    defer consul.Deregister("user-service")

    // 创建gRPC服务器
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal(err)
    }

    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &server{})

    // 注册健康检查
    healthServer := health.NewServer()
    grpc_health_v1.RegisterHealthServer(s, healthServer)
    healthServer.SetServingStatus("user-service", grpc_health_v1.HealthCheckResponse_SERVING)

    // 注册反射
    reflection.Register(s)

    log.Println("User service started on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatal(err)
    }
}
```

---

## 4. CLI工具示例

### 4.1 数据分析CLI

**项目结构**:

```text
data-cli/
├── cmd/
│   └── datacli/
│       └── main.go
├── internal/
│   ├── analyzer/
│   │   ├── stats.go
│   │   └── trends.go
│   ├── export/
│   │   ├── csv.go
│   │   └── json.go
│   ├── repl/
│   │   └── shell.go
│   └── visualize/
│       ├── chart.go
│       └── table.go
├── plugins/
│   └── example/
│       └── plugin.go
└── go.mod
```

**核心代码**:

```go
// cmd/datacli/main.go

package main

import (
    "fmt"
    "os"

    "data-cli/internal/repl"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "datacli",
    Short: "A powerful data analysis CLI tool",
}

var analyzeCmd = &cobra.Command{
    Use:   "analyze [file]",
    Short: "Analyze data file",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        filename := args[0]
        fmt.Printf("Analyzing %s...\n", filename)
        // 分析逻辑
    },
}

var replCmd = &cobra.Command{
    Use:   "repl",
    Short: "Start interactive REPL",
    Run: func(cmd *cobra.Command, args []string) {
        shell := repl.NewShell()
        shell.Run()
    },
}

func main() {
    rootCmd.AddCommand(analyzeCmd)
    rootCmd.AddCommand(replCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

---

## 5. 数据处理示例

### 5.1 实时日志处理

**项目结构**:

```text
log-processor/
├── cmd/
│   └── processor/
│       └── main.go
├── internal/
│   ├── aggregator/
│   │   └── metrics.go
│   ├── alert/
│   │   └── rules.go
│   ├── parser/
│   │   └── log_parser.go
│   └── storage/
│       ├── elasticsearch.go
│       └── timeseries.go
├── config/
│   └── config.yaml
└── docker-compose.yml
```

**核心代码**:

```go
// cmd/processor/main.go

package main

import (
    "bufio"
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "log-processor/internal/aggregator"
    "log-processor/internal/parser"
    "log-processor/internal/storage"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 初始化组件
    parser := parser.NewLogParser()
    agg := aggregator.NewMetricsAggregator()
    store := storage.NewElasticsearch()

    // 处理日志流
    scanner := bufio.NewScanner(os.Stdin)

    go func() {
        for scanner.Scan() {
            line := scanner.Text()

            // 解析日志
            entry, err := parser.Parse(line)
            if err != nil {
                continue
            }

            // 聚合指标
            agg.Add(entry)

            // 存储
            if err := store.Index(ctx, entry); err != nil {
                log.Printf("Failed to store: %v", err)
            }
        }
    }()

    // 等待退出信号
    quit := make(Channel os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down...")
}
```

---

## 6. 实时系统示例

### 6.1 Web爬虫系统

**项目结构**:

```text
web-crawler/
├── cmd/
│   └── crawler/
│       └── main.go
├── internal/
│   ├── crawler/
│   │   ├── crawler.go
│   │   └── worker_pool.go
│   ├── dedup/
│   │   └── bloom_filter.go
│   ├── parser/
│   │   └── html_parser.go
│   ├── ratelimit/
│   │   └── limiter.go
│   └── storage/
│       ├── mongo.go
│       └── cache.go
├── config/
│   └── config.yaml
└── docker-compose.yml
```

**核心代码**:

```go
// internal/crawler/crawler.go

package crawler

import (
    "context"
    "net/http"
    "sync"
    "time"

    "web-crawler/internal/dedup"
    "web-crawler/internal/parser"
    "web-crawler/internal/ratelimit"
)

type Crawler struct {
    client    *http.Client
    parser    *parser.HTMLParser
    dedup     *dedup.BloomFilter
    limiter   *ratelimit.Limiter
    workerPool *WorkerPool
}

func NewCrawler(workers int) *Crawler {
    return &Crawler{
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
        parser:    parser.NewHTMLParser(),
        dedup:     dedup.NewBloomFilter(1000000),
        limiter:   ratelimit.NewLimiter(100), // 100 req/s
        workerPool: NewWorkerPool(workers),
    }
}

func (c *Crawler) Crawl(ctx Context.Context, urls []string) error {
    c.workerPool.Start()
    defer c.workerPool.Stop()

    var wg sync.WaitGroup

    for _, url := range urls {
        if c.dedup.Contains(url) {
            continue
        }

        c.dedup.Add(url)
        wg.Add(1)

        c.workerPool.Submit(&CrawlTask{
            URL:     url,
            Crawler: c,
            WG:      &wg,
        })
    }

    wg.Wait()
    return nil
}

type CrawlTask struct {
    URL     string
    Crawler *Crawler
    WG      *sync.WaitGroup
}

func (t *CrawlTask) Execute(ctx Context.Context) error {
    defer t.WG.Done()

    // 限速
    t.Crawler.limiter.Wait()

    // 获取页面
    resp, err := t.Crawler.client.Get(t.URL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 解析页面
    doc, err := t.Crawler.parser.Parse(resp.Body)
    if err != nil {
        return err
    }

    // 提取链接
    links := t.Crawler.parser.ExtractLinks(doc)

    // 递归爬取
    for _, link := range links {
        if !t.Crawler.dedup.Contains(link) {
            t.Crawler.dedup.Add(link)
            t.WG.Add(1)

            t.Crawler.workerPool.Submit(&CrawlTask{
                URL:     link,
                Crawler: t.Crawler,
                WG:      t.WG,
            })
        }
    }

    return nil
}
```

---

## 7. 使用指南

### 7.1 快速开始

**1. RESTful API服务**:

```bash
cd examples/restful-api
docker-compose up -d
go run cmd/server/main.go
```

**2. 微服务架构**:

```bash
cd examples/microservices
docker-compose up -d
# 启动各个服务
go run services/user/cmd/main.go
go run services/order/cmd/main.go
go run services/product/cmd/main.go
```

```bash
cd examples/data-cli
go build -o datacli cmd/datacli/main.go
./datacli repl
```

**4. 日志处理**:

```bash
cd examples/log-processor
tail -f /var/log/app.log | go run cmd/processor/main.go
```

**5. Web爬虫**:

```bash
cd examples/web-crawler
go run cmd/crawler/main.go --start-url https://example.com
```

---

### 7.2 学习路径
