# TS-CL-004: Go Context Package - Deep Architecture and Cancellation Patterns

> **维度**: Technology Stack > Core Library
> **级别**: S (20+ KB)
> **标签**: #golang #context #cancellation #deadline #timeout #tracing
> **权威来源**:
>
> - [Go context package](https://pkg.go.dev/context) - Official documentation
> - [Go Concurrency Patterns: Context](https://go.dev/blog/context) - Go Blog
> - [Understanding Context](https://medium.com/@cep21/go-contexts-3-examples-4e63725f31f2) - Practical examples

---

## 1. Context Architecture Deep Dive

### 1.1 The Context Tree

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Context Tree Structure                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                              Background()                                    │
│                                   │                                          │
│                    ┌──────────────┼──────────────┐                          │
│                    │              │              │                          │
│               ┌────▼────┐   ┌────▼────┐   ┌────▼────┐                      │
│               │WithValue│   │WithCancel│  │WithTimeout│                     │
│               │ (key=1) │   │         │   │ (30s)     │                     │
│               └────┬────┘   └────┬────┘   └────┬────┘                      │
│                    │             │             │                            │
│              ┌─────▼─────┐  ┌────▼────┐  ┌────▼────┐                        │
│              │WithValue  │  │WithValue│  │WithCancel│                       │
│              │ (key=2)   │  │ (key=3) │  │         │                        │
│              └─────┬─────┘  └────┬────┘  └────┬────┘                        │
│                    │             │             │                             │
│                    └─────────────┴─────────────┘                             │
│                                  │                                           │
│                           ┌──────▼──────┐                                    │
│                           │   Request   │                                    │
│                           │   Handler   │                                    │
│                           └─────────────┘                                    │
│                                                                              │
│  Key Properties:                                                             │
│  - Immutable: Each With* creates a new context                              │
│  - Hierarchical: Children inherit from parents                               │
│  - Cancellation propagates down the tree                                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Context Interface

```go
// Context is the core interface
type Context interface {
    // Deadline returns when this context will be canceled
    Deadline() (deadline time.Time, ok bool)

    // Done returns a channel that's closed when context is canceled
    Done() <-chan struct{}

    // Err returns why the context was canceled
    Err() error

    // Value retrieves data stored in the context
    Value(key interface{}) interface{}
}
```

### 1.3 Context Types

| Type | Purpose | Use Case |
|------|---------|----------|
| `context.Background()` | Root context | Main function, initialization |
| `context.TODO()` | Placeholder | Temporary, when unsure |
| `WithCancel(parent)` | Manual cancellation | User cancels operation |
| `WithDeadline(parent, time)` | Absolute deadline | Scheduled tasks |
| `WithTimeout(parent, duration)` | Relative deadline | API calls, DB queries |
| `WithValue(parent, key, val)` | Request-scoped data | Authentication, tracing |

---

## 2. Cancellation Patterns

### 2.1 Basic Cancellation

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // Always call cancel to release resources

    go worker(ctx)

    // Cancel after 5 seconds
    time.Sleep(5 * time.Second)
    cancel() // Signals all goroutines to stop
}

func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Worker shutting down:", ctx.Err())
            return
        default:
            // Do work
            time.Sleep(100 * time.Millisecond)
        }
    }
}
```

### 2.2 Timeout Pattern

```go
func queryDatabase(ctx context.Context, query string) error {
    // Create timeout context for this specific operation
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    result := make(chan error, 1)
    go func() {
        result <- performQuery(query)
    }()

    select {
    case err := <-result:
        return err
    case <-ctx.Done():
        return fmt.Errorf("query timeout: %w", ctx.Err())
    }
}
```

### 2.3 Deadline Pattern

```go
func scheduleTask(ctx context.Context, executeAt time.Time) error {
    ctx, cancel := context.WithDeadline(ctx, executeAt)
    defer cancel()

    // Check if deadline is already passed
    if deadline, ok := ctx.Deadline(); ok {
        if time.Until(deadline) < 0 {
            return errors.New("deadline already passed")
        }
    }

    <-ctx.Done()
    return ctx.Err()
}
```

---

## 3. Context Value Storage

### 3.1 Value Storage Best Practices

```go
// Define custom key types to avoid collisions
type contextKey string

const (
    userIDKey     contextKey = "userID"
    requestIDKey  contextKey = "requestID"
    traceIDKey    contextKey = "traceID"
)

// Setter functions
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(userIDKey).(string)
    return userID, ok
}

// Usage
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    ctx = WithUserID(ctx, "user-123")
    ctx = WithRequestID(ctx, generateRequestID())

    processRequest(ctx)
}
```

### 3.2 What to Store in Context

| ✅ Store | ❌ Don't Store |
|---------|---------------|
| Request ID / Trace ID | Database connections |
| User ID / Session info | Large objects |
| Authentication tokens | Business logic data |
| Deadline/timeout info | Logger instances |
| Flags for behavior | Configuration |

---

## 4. Go Client Integration

### 4.1 HTTP Client with Context

```go
func makeRequest(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}

// Usage with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

data, err := makeRequest(ctx, "https://api.example.com/data")
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Request timed out")
    }
}
```

### 4.2 Database Operations with Context

```go
func queryWithContext(ctx context.Context, db *sql.DB) error {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    rows, err := db.QueryContext(ctx, "SELECT * FROM users WHERE active = ?", true)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Process row
        }
    }
    return rows.Err()
}
```

### 4.3 gRPC with Context

```go
func grpcCall(ctx context.Context, client pb.ServiceClient) (*pb.Response, error) {
    // Set deadline
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // Add metadata
    md := metadata.New(map[string]string{
        "request-id": "abc-123",
        "user-id": "user-456",
    })
    ctx = metadata.NewOutgoingContext(ctx, md)

    return client.GetData(ctx, &pb.Request{})
}
```

---

## 5. Performance Tuning Guidelines

### 5.1 Context Overhead

| Operation | Approximate Cost |
|-----------|-----------------|
| context.Background() | ~0 (singleton) |
| WithCancel() | ~50-100ns + allocation |
| WithTimeout() | ~100-200ns + allocation |
| WithValue() | ~50-100ns + allocation |
| ctx.Done() check | ~5-10ns |

### 5.2 Best Practices for Performance

```go
// 1. Pass context as first parameter
func DoSomething(ctx context.Context, arg string) error

// 2. Don't store context in structs
type BadService struct {
    ctx context.Context // DON'T DO THIS
}

// 3. Accept context at method call time
type GoodService struct{}
func (s *GoodService) DoWork(ctx context.Context) error

// 4. Create child contexts for specific operations
func (s *Service) HandleRequest(ctx context.Context) error {
    // Short timeout for auth check
    authCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
    defer cancel()
    if err := s.authenticate(authCtx); err != nil {
        return err
    }

    // Longer timeout for main processing
    processCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    return s.process(processCtx)
}
```

---

## 6. Visual Representations

### 6.1 Cancellation Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Cancellation Propagation                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Parent Context                    Child Contexts                          │
│   ┌──────────────┐                  ┌─────────────┐  ┌─────────────┐        │
│   │   cancel()   │─────────────────>│  <-Done()   │  │  <-Done()   │        │
│   └──────────────┘                  └──────┬──────┘  └──────┬──────┘        │
│        │                                    │                │              │
│        │                           ┌────────▼────────┐       │              │
│        │                           │  Grandchild 1   │       │              │
│        │                           │   <-Done()      │       │              │
│        │                           └────────┬────────┘       │              │
│        │                                    │                │              │
│        └────────────────────────────────────┴────────────────┘              │
│                                             │                               │
│                                    ┌────────▼────────┐                      │
│                                    │  Grandchild 2   │                      │
│                                    │   <-Done()      │                      │
│                                    └─────────────────┘                      │
│                                                                              │
│   All descendants receive cancellation signal simultaneously                │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Context Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Context Lifecycle                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Create ──> Use ──> Propagate ──> Cancel ──> Cleanup                        │
│     │        │         │           │          │                             │
│     │        │         │           │          └── Release resources         │
│     │        │         │           └── Signal done channel                  │
│     │        │         └── Pass to child contexts                          │
│     │        └── Store values, check deadlines                              │
│     └── Background(), TODO(), WithCancel(), WithTimeout()                   │
│                                                                              │
│  Critical Points:                                                            │
│  - Always call cancel() to release resources                                │
│  - Check ctx.Err() to determine cancellation reason                         │
│  - Don't pass nil context, use context.TODO()                               │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Comparison with Alternatives

| Approach | Pros | Cons | When to Use |
|----------|------|------|-------------|
| **Context** | Standard, composable, cancellation support | Can be misused for data | Go 1.7+, all new code |
| **Channels** | Direct control, no magic | Verbose, manual propagation | Complex coordination |
| **sync.WaitGroup** | Simple wait semantics | No cancellation | Wait for completion |
| **Atomic flags** | Fast, simple | No composition | Simple stop signals |
| **Error return** | Explicit, testable | No async cancellation | Synchronous code |

---

## 8. Configuration Best Practices

```go
// Timeout configuration constants
const (
    // External service timeouts
    DefaultHTTPTimeout    = 10 * time.Second
    DefaultDBTimeout      = 5 * time.Second
    DefaultGRPCTimeout    = 10 * time.Second

    // Internal operation timeouts
    DefaultCacheTimeout   = 100 * time.Millisecond
    DefaultAuthTimeout    = 500 * time.Millisecond
)

// Context configuration struct
type ContextConfig struct {
    RequestTimeout  time.Duration
    ShutdownTimeout time.Duration
    EnableTracing   bool
}

func NewServerContext(cfg ContextConfig) (context.Context, context.CancelFunc) {
    ctx := context.Background()

    if cfg.EnableTracing {
        ctx = withTracing(ctx)
    }

    return context.WithTimeout(ctx, cfg.RequestTimeout)
}
```

---

## 9. Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Context Best Practices                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  API Design:                                                                 │
│  □ Context as first parameter of functions/methods                          │
│  □ Store only request-scoped values (IDs, metadata)                         │
│  □ Never store context in struct fields                                     │
│  □ Propagate context through call chain                                     │
│                                                                              │
│  Cancellation:                                                               │
│  □ Always call cancel() (defer immediately after creation)                  │
│  □ Check ctx.Done() in long-running operations                              │
│  □ Handle context errors appropriately                                      │
│  □ Use specific timeouts for different operations                           │
│                                                                              │
│  Values:                                                                     │
│  □ Use typed keys to avoid collisions                                       │
│  □ Provide getter/setter functions                                          │
│  □ Document what values are available                                       │
│  □ Don't store large objects or connections                                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (20+ KB, comprehensive coverage)
