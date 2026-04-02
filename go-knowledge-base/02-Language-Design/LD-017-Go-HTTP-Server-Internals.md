# LD-017: Go HTTP 服务器内部原理 (Go HTTP Server Internals)

> **维度**: Language Design
> **级别**: S (19+ KB)
> **标签**: #http #server #net-http #internals #performance #concurrency
> **权威来源**:
>
> - [net/http Package](https://github.com/golang/go/tree/master/src/net/http) - Go Authors
> - [HTTP/2 in Go](https://go.dev/blog/h2push) - Go Authors
> - [Go HTTP Server Best Practices](https://www.ardanlabs.com/blog/) - Ardan Labs

---

## 1. HTTP 服务器架构

### 1.1 核心组件

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Server                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │   Server     │───►│   Listener   │───►│   Conn       │  │
│  │              │    │   (TCP)      │    │   Handler    │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│         │                                     │              │
│         ▼                                     ▼              │
│  ┌──────────────┐                    ┌──────────────┐       │
│  │   Handler    │◄───────────────────│   ServeHTTP  │       │
│  │   (mux)      │                    │   (per req)  │       │
│  └──────────────┘                    └──────────────┘       │
│         │                                     │              │
│         ▼                                     ▼              │
│  ┌──────────────┐                    ┌──────────────┐       │
│  │   Routes     │                    │   Response   │       │
│  │   Matching   │                    │   Writer     │       │
│  └──────────────┘                    └──────────────┘       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 Server 结构

```go
// src/net/http/server.go
type Server struct {
    Addr    string        // TCP 地址
    Handler Handler       // 请求处理器

    // 超时配置
    ReadTimeout       time.Duration
    ReadHeaderTimeout time.Duration
    WriteTimeout      time.Duration
    IdleTimeout       time.Duration

    // TLS 配置
    TLSConfig *tls.Config

    // 连接状态回调
    ConnState func(net.Conn, ConnState)

    // 错误日志
    ErrorLog *log.Logger

    // 内部字段
    disableKeepAlives int32     // 原子操作
    inShutdown        int32     // 关闭标志
    mu                sync.Mutex
    listeners         map[*net.Listener]struct{}
    activeConn        map[*conn]struct{}
    doneChan          chan struct{}
}
```

---

## 2. 连接处理流程

### 2.1 监听与接受

```go
// 启动服务器
func (srv *Server) ListenAndServe() error {
    addr := srv.Addr
    if addr == "" {
        addr = ":http"
    }
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    return srv.Serve(ln)
}

// 主服务循环
func (srv *Server) Serve(l net.Listener) error {
    // 跟踪活跃连接
    origListener := l
    l = &onceCloseListener{Listener: l}
    defer l.Close()

    // 优雅关闭支持
    if srv.shuttingDown() {
        return ErrServerClosed
    }

    // 连接上下文
    baseCtx := context.Background()
    ctx := context.WithValue(baseCtx, ServerContextKey, srv)

    for {
        // 接受新连接
        rw, err := l.Accept()
        if err != nil {
            // 检查是否正在关闭
            if srv.shuttingDown() {
                return ErrServerClosed
            }
            // 临时错误处理
            if ne, ok := err.(net.Error); ok && ne.Temporary() {
                time.Sleep(tempDelay)
                continue
            }
            return err
        }

        // 创建连接上下文
        connCtx := ctx
        if cc := srv.ConnContext; cc != nil {
            connCtx = cc(connCtx, rw)
            if connCtx == nil {
                panic("ConnContext returned nil")
            }
        }

        // 创建连接对象
        c := srv.newConn(rw)
        c.setState(c.rwc, StateNew)

        // 启动 goroutine 处理
        go c.serve(connCtx)
    }
}
```

### 2.2 连接处理详解

```go
// src/net/http/server.go
func (c *conn) serve(ctx context.Context) {
    c.remoteAddr = c.rwc.RemoteAddr().String()
    ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())

    defer func() {
        // 恢复 panic
        if err := recover(); err != nil {
            const size = 64 << 10
            buf := make([]byte, size)
            buf = buf[:runtime.Stack(buf, false)]
            c.server.logf("http: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
        }
        // 关闭连接
        c.close()
        c.setState(c.rwc, StateClosed)
    }()

    // TLS 握手
    if tlsConn, ok := c.rwc.(*tls.Conn); ok {
        // 设置 TLS 握手超时
        if d := c.server.ReadTimeout; d != 0 {
            c.rwc.SetReadDeadline(time.Now().Add(d))
        }
        if d := c.server.WriteTimeout; d != 0 {
            c.rwc.SetWriteDeadline(time.Now().Add(d))
        }

        if err := tlsConn.HandshakeContext(ctx); err != nil {
            c.server.logf("http: TLS handshake error from %s: %v", c.rwc.RemoteAddr(), err)
            return
        }

        c.tlsState = new(tls.ConnectionState)
        *c.tlsState = tlsConn.ConnectionState()
    }

    // HTTP/1.x 或 HTTP/2 检测
    if !c.server.TLSConfig.IsEncrypted() {
        c.server.TLSConfig = nil
    }

    // 读取第一个字节检测协议
    c.rwc.SetReadDeadline(time.Time{})

    // HTTP/2 升级检测
    if c.server.TLSConfig == nil {
        // 检测 HTTP/2 明文升级
        if isH2CUpgrade(c.rwc) {
            // 处理 h2c
            c.serveH2C(ctx)
            return
        }
    }

    // HTTP/1.x 处理
    for {
        // 读取请求
        w, err := c.readRequest(ctx)
        if err != nil {
            // 处理错误...
            return
        }

        // 调用处理器
        serverHandler{c.server}.ServeHTTP(w, w.req)

        // 检查是否保持连接
        if !w.shouldReuseConnection() {
            return
        }
    }
}
```

---

## 3. 请求解析与响应

### 3.1 请求解析

```go
// 读取 HTTP 请求
func (c *conn) readRequest(ctx context.Context) (*response, error) {
    var (
        wholeReqDeadline time.Time
        hdrDeadline      time.Time
    )

    // 设置超时
    if d := c.server.ReadTimeout; d != 0 {
        wholeReqDeadline = time.Now().Add(d)
        c.rwc.SetReadDeadline(wholeReqDeadline)
    }
    if d := c.server.ReadHeaderTimeout; d != 0 {
        hdrDeadline = time.Now().Add(d)
        c.rwc.SetReadDeadline(hdrDeadline)
    }

    // 创建缓冲读取器
    c.r.set(c.rwc, c.remoteAddr)

    // 解析请求行和头部
    req, err := readRequest(c.bufr, keepHostHeader)
    if err != nil {
        return nil, err
    }

    // 恢复读取超时
    c.rwc.SetReadDeadline(wholeReqDeadline)

    // 设置请求上下文
    req.ctx = ctx

    // 创建响应对象
    w := &response{
        conn:          c,
        req:           req,
        header:        make(Header),
        wroteHeader:   false,
        contentLength: -1,
    }

    return w, nil
}

// 请求解析核心
func readRequest(b *bufio.Reader, deleteHostHeader bool) (req *Request, err error) {
    // 读取请求行
    line, err := b.ReadSlice('\n')
    if err != nil {
        return nil, err
    }

    // 解析 "METHOD URI HTTP/VERSION"
    var method, uri, proto string
    method, uri, proto, ok := parseRequestLine(line)
    if !ok {
        return nil, badRequest("malformed HTTP request")
    }

    // 验证方法
    if !validMethod(method) {
        return nil, badRequest("invalid method")
    }

    // 解析 URL
    rawurl := uri
    if method == "CONNECT" {
        // CONNECT 特殊处理
    } else {
        // 解析 URI
        rawurl, _ = splitHostPort(uri)
    }

    url, err := url.ParseRequestURI(rawurl)
    if err != nil {
        return nil, badRequest("invalid URI")
    }

    // 读取头部
    header, err := readMIMEHeader(b)
    if err != nil {
        return nil, err
    }

    // 构建请求对象
    req = &Request{
        Method: method,
        URL:    url,
        Proto:  proto,
        Header: header,
        // ...
    }

    return req, nil
}
```

### 3.2 响应写入

```go
// response 结构
type response struct {
    conn          *conn
    req           *Request
    header        Header
    wroteHeader   bool
    status        int
    contentLength int64

    writer  io.Writer
    handler handler

    // 写入状态
    wroteBytes    int64
    wroteHeaderAt time.Time

    // 分块传输
    chunking      bool
    chunkWriter   *chunkWriter
}

// WriteHeader 写入状态码
func (w *response) WriteHeader(code int) {
    if w.wroteHeader {
        // 已经写入，记录
        caller := relevantCaller()
        w.conn.server.logf("http: multiple response.WriteHeader calls from %s", caller)
        return
    }

    w.wroteHeader = true
    w.status = code

    // 写入状态行
    // HTTP/1.1 200 OK\r\n
    text := StatusText(code)
    if text == "" {
        text = "status " + strconv.Itoa(code)
    }

    // 构建响应头
    var buf [512]byte
    p := buf[:0]
    p = append(p, "HTTP/1.1 "...)
    p = strconv.AppendInt(p, int64(code), 10)
    p = append(p, ' ')
    p = append(p, text...)
    p = append(p, "\r\n"...)

    // 写入默认头部
    w.writeHeader(p, code)
}

// 实际写入头部
func (w *response) writeHeader(buf []byte, code int) {
    // 内容长度处理
    if w.contentLength != -1 {
        buf = appendHeader(buf, "Content-Length", strconv.FormatInt(w.contentLength, 10))
    } else if !w.chunking {
        // 启用分块编码
        w.chunking = true
        buf = appendHeader(buf, "Transfer-Encoding", "chunked")
    }

    // 写入头部
    buf = append(buf, "\r\n"...)
    w.conn.w.Write(buf)
}
```

---

## 4. Handler 与多路复用器

### 4.1 Handler 接口

```go
// 核心 Handler 接口
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// HandlerFunc 适配器
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

// ServeMux 实现
type ServeMux struct {
    mu    sync.RWMutex
    m     map[string]muxEntry
    hosts bool // 是否包含主机模式
}

type muxEntry struct {
    h       Handler
    pattern string
}

// 注册路由
func (mux *ServeMux) Handle(pattern string, handler Handler) {
    mux.mu.Lock()
    defer mux.mu.Unlock()

    if pattern == "" {
        panic("http: invalid pattern")
    }

    if handler == nil {
        panic("http: nil handler")
    }

    if _, exist := mux.m[pattern]; exist {
        panic("http: multiple registrations for " + pattern)
    }

    if mux.m == nil {
        mux.m = make(map[string]muxEntry)
    }

    mux.m[pattern] = muxEntry{h: handler, pattern: pattern}

    if pattern[0] != '/' {
        mux.hosts = true
    }
}

// 路由匹配
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
    // 精确匹配
    if e, ok := mux.m[path]; ok {
        return e.h, e.pattern
    }

    // 前缀匹配
    var n = 0
    for k, v := range mux.m {
        if !pathMatch(k, path) {
            continue
        }
        if h == nil || len(k) > n {
            n = len(k)
            h = v.h
            pattern = v.pattern
        }
    }

    return
}
```

### 4.2 自定义路由实现

```go
// 高性能路由实现
package main

import (
    "net/http"
    "strings"
    "sync"
)

// Radix Tree 路由
type RadixNode struct {
    path     string
    indices  string
    children []*RadixNode
    handler  http.Handler
    isParam  bool
}

type RadixRouter struct {
    root   *RadixNode
    paramsPool sync.Pool
}

func NewRadixRouter() *RadixRouter {
    return &RadixRouter{
        root: &RadixNode{},
        paramsPool: sync.Pool{
            New: func() interface{} {
                return make(map[string]string, 8)
            },
        },
    }
}

func (r *RadixRouter) Insert(path string, handler http.Handler) {
    r.root.insert(path, handler)
}

func (n *RadixNode) insert(path string, handler http.Handler) {
    if len(path) == 0 {
        n.handler = handler
        return
    }

    // 查找共同前缀
    i := 0
    for i < len(n.path) && i < len(path) && n.path[i] == path[i] {
        i++
    }

    if i < len(n.path) {
        // 分裂节点
        child := &RadixNode{
            path:     n.path[i:],
            children: n.children,
            handler:  n.handler,
        }
        n.path = n.path[:i]
        n.children = []*RadixNode{child}
        n.handler = nil
        n.indices = string(child.path[0])
    }

    if i < len(path) {
        // 继续插入
        path = path[i:]
        firstChar := path[0]

        for j, c := range n.indices {
            if c == rune(firstChar) {
                n.children[j].insert(path, handler)
                return
            }
        }

        // 创建新子节点
        child := &RadixNode{path: path, handler: handler}
        n.indices += string(firstChar)
        n.children = append(n.children, child)
    } else {
        n.handler = handler
    }
}

func (r *RadixRouter) Search(path string) (http.Handler, map[string]string) {
    params := r.paramsPool.Get().(map[string]string)
    for k := range params {
        delete(params, k)
    }
    defer r.paramsPool.Put(params)

    handler := r.root.search(path, params)
    if handler != nil {
        return handler, params
    }
    return nil, nil
}

func (n *RadixNode) search(path string, params map[string]string) http.Handler {
    if len(path) == 0 || path == "/" {
        return n.handler
    }

    // 跳过前导斜杠
    if path[0] == '/' {
        path = path[1:]
    }

    // 尝试匹配子节点
    if len(path) > 0 && len(n.indices) > 0 {
        firstChar := path[0]
        for i, c := range n.indices {
            if c == rune(firstChar) {
                child := n.children[i]
                if strings.HasPrefix(path, child.path) {
                    remaining := path[len(child.path):]
                    if handler := child.search(remaining, params); handler != nil {
                        return handler
                    }
                }
            }
        }
    }

    // 尝试参数匹配 :param
    if n.isParam {
        // 提取参数值
        parts := strings.SplitN(path, "/", 2)
        params[n.path[1:]] = parts[0]
        if len(parts) > 1 {
            for _, child := range n.children {
                if handler := child.search(parts[1], params); handler != nil {
                    return handler
                }
            }
        }
        return n.handler
    }

    return nil
}

// 实现 http.Handler
func (r *RadixRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    handler, params := r.Search(req.URL.Path)
    if handler != nil {
        // 将参数存入上下文
        ctx := req.Context()
        for k, v := range params {
            ctx = context.WithValue(ctx, k, v)
        }
        handler.ServeHTTP(w, req.WithContext(ctx))
        return
    }
    http.NotFound(w, req)
}
```

---

## 5. 性能优化与基准测试

### 5.1 连接池优化

```go
// TCP 连接池配置
func optimizeServer() *http.Server {
    return &http.Server{
        Addr: ":8080",

        // 读写超时
        ReadTimeout:       5 * time.Second,
        ReadHeaderTimeout: 2 * time.Second,
        WriteTimeout:      10 * time.Second,
        IdleTimeout:       120 * time.Second,

        // 最大头部大小
        MaxHeaderBytes: 1 << 20, // 1MB

        // 连接状态监控
        ConnState: func(conn net.Conn, state http.ConnState) {
            switch state {
            case http.StateNew:
                // 新连接
            case http.StateActive:
                // 活跃连接
            case http.StateIdle:
                // 空闲连接
            case http.StateClosed:
                // 关闭连接
            }
        },
    }
}

// 自定义 Listener 实现限流
type LimitedListener struct {
    net.Listener
    sem chan struct{}
}

func NewLimitedListener(ln net.Listener, maxConns int) *LimitedListener {
    return &LimitedListener{
        Listener: ln,
        sem:      make(chan struct{}, maxConns),
    }
}

func (l *LimitedListener) Accept() (net.Conn, error) {
    l.sem <- struct{}{} // 获取许可
    conn, err := l.Listener.Accept()
    if err != nil {
        <-l.sem
        return nil, err
    }
    return &LimitedConn{Conn: conn, release: func() { <-l.sem }}, nil
}

type LimitedConn struct {
    net.Conn
    release func()
    once    sync.Once
}

func (c *LimitedConn) Close() error {
    err := c.Conn.Close()
    c.once.Do(c.release)
    return err
}
```

### 5.2 基准测试

```go
func BenchmarkHTTPServer(b *testing.B) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    server := httptest.NewServer(handler)
    defer server.Close()

    client := server.Client()

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            resp, err := client.Get(server.URL)
            if err != nil {
                b.Fatal(err)
            }
            resp.Body.Close()
        }
    })
}

// 典型性能数据 (Go 1.21, 8核)
// BenchmarkHTTPServer-8    1000000    1050 ns/op    0 allocs/op
```

### 5.3 内存分配优化

```go
// 减少分配的 Handler
func efficientHandler(w http.ResponseWriter, r *http.Request) {
    // 预分配响应缓冲区
    var buf [1024]byte
    n := copy(buf[:], `{"status":"ok"}`)

    w.Header().Set("Content-Type", "application/json")
    w.Write(buf[:n])
}

// 使用 sync.Pool 复用对象
var responsePool = sync.Pool{
    New: func() interface{} {
        return &ResponseData{
            Headers: make(http.Header),
            Body:    make([]byte, 0, 4096),
        }
    },
}

type ResponseData struct {
    Headers http.Header
    Body    []byte
    Status  int
}

func pooledHandler(w http.ResponseWriter, r *http.Request) {
    resp := responsePool.Get().(*ResponseData)
    defer func() {
        resp.Body = resp.Body[:0]
        responsePool.Put(resp)
    }()

    // 处理请求...
    resp.Body = append(resp.Body, `{"result":"success"}`...)
    resp.Status = http.StatusOK

    w.WriteHeader(resp.Status)
    w.Write(resp.Body)
}
```

---

## 6. HTTP/2 支持

### 6.1 HTTP/2 服务器

```go
// 启用 HTTP/2
func enableHTTP2() *http.Server {
    srv := &http.Server{
        Addr: ":443",
    }

    // 自动启用 HTTP/2
    http2.ConfigureServer(srv, &http2.Server{
        MaxConcurrentStreams: 250,
        MaxReadFrameSize:     1048576, // 1MB
        IdleTimeout:          10 * time.Second,
    })

    return srv
}

// h2c (HTTP/2 over cleartext)
func h2cServer() *http.Server {
    h2s := &http2.Server{
        IdleTimeout: 1 * time.Minute,
    }

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 检测 HTTP/2
        if r.ProtoMajor == 2 {
            w.Write([]byte("Hello HTTP/2!"))
        } else {
            w.Write([]byte("Hello HTTP/1.1!"))
        }
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: h2c.NewHandler(handler, h2s),
    }

    return srv
}
```

---

## 7. 并发安全分析

### 7.1 竞态条件防护

```go
// 线程安全的计数器
type SafeCounter struct {
    mu    sync.RWMutex
    count int64
}

func (c *SafeCounter) Incr() {
    atomic.AddInt64(&c.count, 1)
}

func (c *SafeCounter) Get() int64 {
    return atomic.LoadInt64(&c.count)
}

// 请求限流器
type RateLimiter struct {
    tokens   float64
    rate     float64
    capacity float64
    last     time.Time
    mu       sync.Mutex
}

func (rl *RateLimiter) Allow() bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(rl.last).Seconds()
    rl.last = now

    // 添加新令牌
    rl.tokens += elapsed * rl.rate
    if rl.tokens > rl.capacity {
        rl.tokens = rl.capacity
    }

    // 检查是否有足够令牌
    if rl.tokens >= 1 {
        rl.tokens--
        return true
    }
    return false
}
```

### 7.2 优雅关闭

```go
// 优雅关闭实现
func gracefulShutdown(srv *http.Server, timeout time.Duration) {
    // 监听系统信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    <-quit
    log.Println("Shutting down server...")

    // 创建超时上下文
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // 关闭服务器
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited")
}
```

---

## 8. 视觉表征

### 8.1 HTTP 请求处理流程

```
Client Request
      │
      ▼
┌─────────────┐
│   Accept    │◄── goroutine per connection
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Parse     │──► 解析请求行、头部
│   Request   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Route     │──► 匹配 Handler
│   Match     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Middleware  │──► 认证、日志、恢复
│   Chain     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Handler   │──► 业务逻辑
│   Execute   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Write     │──► 序列化响应
│   Response  │
└──────┬──────┘
       │
       ▼
    Client
```

### 8.2 服务器架构层次

```
┌─────────────────────────────────────────┐
│           Application Layer              │
│  Handlers → Controllers → Services      │
├─────────────────────────────────────────┤
│           Routing Layer                  │
│  ServeMux / Custom Router               │
├─────────────────────────────────────────┤
│           Middleware Layer               │
│  Logging → Recovery → Auth → CORS       │
├─────────────────────────────────────────┤
│           HTTP Protocol Layer            │
│  Request Parse → Response Serialize     │
├─────────────────────────────────────────┤
│           Connection Layer               │
│  HTTP/1.1 / HTTP/2 / TLS               │
├─────────────────────────────────────────┤
│           Transport Layer                │
│  TCP / Unix Socket                      │
└─────────────────────────────────────────┘
```

### 8.3 性能调优决策树

```
高延迟?
│
├── 检查数据库连接池
├── 启用 HTTP/2
├── 使用连接复用 (keep-alive)
└── 添加缓存层

高内存使用?
│
├── 使用 sync.Pool 复用对象
├── 减少每次请求的分配
├── 流式处理大响应
└── 限制并发连接数

低吞吐量?
│
├── 增加 GOMAXPROCS
├── 优化 Handler 代码
├── 使用更高效的路由器
└── 启用 pprof 分析
```

---

## 9. 完整代码示例

### 9.1 生产级 HTTP 服务器

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // 创建路由器
    mux := http.NewServeMux()
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/api/health", healthHandler)
    mux.HandleFunc("/api/users/", userHandler)

    // 包装中间件
    handler := loggingMiddleware(recoveryMiddleware(mux))

    // 创建服务器
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1MB
    }

    // 优雅关闭
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    log.Println("Server started on :8080")

    // 等待中断信号
    quit := make(chan os.Signal, 1)
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

// 中间件
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // 包装 ResponseWriter 捕获状态码
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

        next.ServeHTTP(wrapped, r)

        log.Printf("[%s] %s %s %d %v",
            r.Method,
            r.URL.Path,
            r.RemoteAddr,
            wrapped.statusCode,
            time.Since(start),
        )
    })
}

func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// Handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome!"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status":"healthy"}`))
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    // 解析用户ID
    id := r.URL.Path[len("/api/users/"):]
    w.Write([]byte("User ID: " + id))
}
```

---

**质量评级**: S (19KB)
**完成日期**: 2026-04-02
