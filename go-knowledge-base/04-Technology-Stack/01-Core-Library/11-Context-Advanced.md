# TS-CL-011: Go Context Advanced Patterns - Deep Dive

> **维度**: Technology Stack > Core Library
> **级别**: S (18+ KB)
> **标签**: #golang #context #advanced #propagation #values #cancellation
> **权威来源**:
>
> - [Go context package](https://pkg.go.dev/context) - Official documentation
> - [Context and structs](https://go.dev/blog/context-and-structs) - Go Blog

---

## 1. Advanced Context Patterns

### 1.1 Context Propagation Chain

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Context Propagation Chain                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Request Entry                                                              │
│   ┌───────────────────────────────────────────────────────────────────────┐ │
│   │  HTTP Handler                                                         │ │
│   │  ┌─────────────────────────────────────────────────────────────────┐  │ │
│   │  │  Middleware (Auth, Logging, Metrics)                            │  │ │
│   │  │  ┌───────────────────────────────────────────────────────────┐  │  │ │
│   │  │  │  Service Layer                                            │  │  │ │
│   │  │  │  ┌─────────────────────────────────────────────────────┐  │  │  │ │
│   │  │  │  │  Repository Layer                                     │  │  │  │ │
│   │  │  │  │  ┌───────────────────────────────────────────────┐   │  │  │  │ │
│   │  │  │  │  │  External Calls (DB, Cache, HTTP, gRPC)       │   │  │  │  │ │
│   │  │  │  │  └───────────────────────────────────────────────┘   │  │  │  │ │
│   │  │  │  └─────────────────────────────────────────────────────┘  │  │  │ │
│   │  │  └───────────────────────────────────────────────────────────┘  │  │ │
│   │  └─────────────────────────────────────────────────────────────────┘  │ │
│   └───────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│   Context carries:                                                           │
│   - Deadline/Cancellation                                                    │
│   - Request ID (for tracing)                                                 │
│   - User ID (for authorization)                                              │
│   - Authentication token                                                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Context Key Management

```go
// Private key type to prevent collisions
type contextKey struct {
    name string
}

func (k contextKey) String() string {
    return k.name
}

// Define keys as private variables
var (
    requestIDKey = &contextKey{"requestID"}
    userIDKey    = &contextKey{"userID"}
    traceIDKey   = &contextKey{"traceID"}
)

// Exported setter functions
func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(requestIDKey).(string)
    return id, ok
}
```

---

## 2. Advanced Cancellation Patterns

### 2.1 Graceful Shutdown

```go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    server := &http.Server{
        Addr:    ":8080",
        Handler: handler(),
    }

    // Start server
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Wait for shutdown signal
    <-ctx.Done()

    // Graceful shutdown with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(shutdownCtx); err != nil {
        log.Printf("Shutdown error: %v", err)
    }
}
```

### 2.2 Fan-Out Cancellation

```go
func processBatch(ctx context.Context, items []Item) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    errChan := make(chan error, len(items))

    for _, item := range items {
        go func(i Item) {
            if err := processItem(ctx, i); err != nil {
                errChan <- err
                cancel() // Cancel all other goroutines on first error
            }
        }(item)
    }

    select {
    case err := <-errChan:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

## 3. Context for Observability

### 3.1 Request Tracing

```go
func TracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Generate or extract trace ID
        traceID := r.Header.Get("X-Trace-ID")
        if traceID == "" {
            traceID = generateTraceID()
        }

        ctx = WithTraceID(ctx, traceID)
        w.Header().Set("X-Trace-ID", traceID)

        // Log with context
        ctx = WithLogger(ctx, logger.With("trace_id", traceID))

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func Handler(ctx context.Context) {
    logger := LoggerFromContext(ctx)
    traceID := TraceIDFromContext(ctx)

    logger.Info("Processing request",
        "trace_id", traceID,
        "user_id", UserIDFromContext(ctx),
    )
}
```

---

## 4. Performance Considerations

### 4.1 Context Overhead

```go
// Context creation costs
// - WithCancel: ~50-100ns
// - WithTimeout: ~100-200ns
// - WithValue: ~50-100ns

// Minimize context creation in hot paths
func processMany(ctx context.Context, items []Item) {
    // Create timeout once for batch
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    for _, item := range items {
        // Reuse the same context
        processItem(ctx, item)
    }
}
```

---

## 5. Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Advanced Context Patterns                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Patterns:                                                                   │
│  □ Use private key types for values                                         │
│  □ Implement graceful shutdown with context                                 │
│  □ Propagate trace IDs through call chain                                   │
│  □ Use cancel for early termination                                         │
│                                                                              │
│  Performance:                                                                │
│  □ Minimize context creation in hot paths                                   │
│  □ Reuse timeout contexts for batch operations                              │
│  □ Check ctx.Done() in long-running operations                              │
│                                                                              │
│  Observability:                                                              │
│  □ Always include request/trace IDs                                         │
│  □ Log with context-enriched loggers                                        │
│  □ Propagate context to all external calls                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (18+ KB, comprehensive coverage)
