# TS-CL-004: Go Context Package - Deep Dive

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #context #cancellation #deadline #request-scoping
> **权威来源**:
> - [Go Context Package](https://golang.org/pkg/context/) - Go standard library
> - [Go Concurrency Patterns: Context](https://blog.golang.org/context) - Go Blog
> - [Context and structs](https://go.dev/blog/context-and-structs) - Go Blog

---

## 1. Context Architecture

### 1.1 Context Tree Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Context Tree Structure                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                         Background() / TODO()                                │
│                                (root)                                        │
│                                  │                                           │
│           ┌──────────────────────┼──────────────────────┐                   │
│           │                      │                      │                   │
│           ▼                      ▼                      ▼                   │
│    ┌─────────────┐        ┌─────────────┐        ┌─────────────┐           │
│    │ WithCancel  │        │ WithTimeout │        │ WithValue   │           │
│    │  (ctx1)     │        │  (ctx2)     │        │  (ctx3)     │           │
│    └──────┬──────┘        └──────┬──────┘        └──────┬──────┘           │
│           │                      │                      │                   │
│    ┌──────┴──────┐               │                      │                   │
│    ▼             ▼               ▼                      ▼                   │
│ ┌────────┐  ┌────────┐    ┌─────────────┐        ┌─────────────┐           │
│ │ ctx1.1 │  │ ctx1.2 │    │ WithCancel  │        │ WithCancel  │           │
│ │        │  │        │    │  (ctx2.1)   │        │  (ctx3.1)   │           │
│ └────────┘  └────────┘    └─────────────┘        └──────┬──────┘           │
│                                                         │                   │
│                                                         ▼                   │
│                                                  ┌─────────────┐           │
│                                                  │ WithValue   │           │
│                                                  │  (ctx3.1.1) │           │
│                                                  └─────────────┘           │
│                                                                              │
│  Propagation Rules:                                                          │
│  - Cancellation: Propagates DOWN the tree (parent cancels → children cancel)│
│  - Deadline: Propagates DOWN (child deadline ≤ parent deadline)             │
│  - Values: Propagates UP (lookup goes from child to parent)                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Context Interface

```go
// Context interface - the foundation of request scoping
type Context interface {
    // Deadline returns the time when this Context will be cancelled
    // ok=false if no deadline is set
    Deadline() (deadline time.Time, ok bool)
    
    // Done returns a channel that closes when Context is cancelled
    // Used for cancellation detection
    Done() <-chan struct{}
    
    // Err returns why the Context was cancelled
    // nil if not cancelled, Canceled or DeadlineExceeded if cancelled
    Err() error
    
    // Value returns the value for the key, or nil if not found
    // Searches up the context tree
    Value(key interface{}) interface{}
}

// Context errors
var Canceled = errors.New("context canceled")
var DeadlineExceeded error = deadlineExceededError{}
```

---

## 2. Context Creation Patterns

### 2.1 Root Contexts

```go
// Background() - empty root context, never cancelled, no values
// Used for: main function, initialization, tests
ctx := context.Background()

// TODO() - placeholder for when context is unclear
// Used for: refactoring, unclear propagation paths
ctx := context.TODO()

// Both return emptyCtx{} - never cancelled, no deadline, no values
```

### 2.2 Derived Contexts

```go
// WithCancel - returns ctx + cancel function
// Use for: explicit cancellation control
ctx, cancel := context.WithCancel(parent)
defer cancel() // Always call to prevent goroutine leak

// WithDeadline - cancels at specific time
// Use for: absolute time limits
ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Minute))
defer cancel()

// WithTimeout - cancels after duration
// Use for: relative time limits
ctx, cancel := context.WithTimeout(parent, 30*time.Second)
defer cancel()

// WithValue - carries request-scoped values
// Use for: request IDs, user info, tracing data
ctx := context.WithValue(parent, key, value)
```

---

## 3. Cancellation Patterns

### 3.1 Basic Cancellation

```go
func longRunningOperation(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err() // Return when cancelled
        default:
            // Do work...
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        time.Sleep(5 * time.Second)
        cancel() // Cancel after 5 seconds
    }()
    
    if err := longRunningOperation(ctx); err != nil {
        log.Println("Operation cancelled:", err)
    }
}
```

### 3.2 Timeout Pattern

```go
func callExternalAPI(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    req, err := http.NewRequestWithContext(ctx, "GET", "https://api.example.com", nil)
    if err != nil {
        return err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("API call timeout: %w", err)
        }
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

### 3.3 Propagating Cancellation

```go
func processRequest(ctx context.Context, req Request) error {
    // Pass context to all downstream calls
    
    user, err := getUser(ctx, req.UserID)
    if err != nil {
        return err
    }
    
    orders, err := getOrders(ctx, user.ID)
    if err != nil {
        return err
    }
    
    return sendNotification(ctx, user, orders)
}

func getUser(ctx context.Context, userID string) (*User, error) {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    // Database call with context
    row := db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", userID)
    // ...
}
```

---

## 4. Value Storage Patterns

### 4.1 Typed Keys

```go
// Private type prevents collisions
type contextKey string

const (
    userKey     contextKey = "user"
    requestIDKey contextKey = "requestID"
)

func WithUser(ctx context.Context, user *User) context.Context {
    return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (*User, bool) {
    user, ok := ctx.Value(userKey).(*User)
    return user, ok
}

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(requestIDKey).(string)
    return id, ok
}
```

### 4.2 HTTP Middleware Pattern

```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        ctx := WithRequestID(r.Context(), requestID)
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        next.ServeHTTP(w, r)
        
        if requestID, ok := RequestIDFromContext(r.Context()); ok {
            log.Printf("[%s] %s %s %v", requestID, r.Method, r.URL.Path, time.Since(start))
        }
    })
}
```

---

## 5. Best Practices

```go
// DO: Pass context as first parameter
func DoSomething(ctx context.Context, arg string) error

// DON'T: Store context in structs
type BadService struct {
    ctx context.Context // BAD!
}

// DO: Pass explicitly
type GoodService struct{}
func (s *GoodService) DoSomething(ctx context.Context) error

// DO: Always call cancel
defer cancel()

// DO: Check ctx.Err() before expensive operations
if err := ctx.Err(); err != nil {
    return err
}

// DO: Use typed keys for values
// DON'T: Use string keys (collision risk)

// DO: Respect cancellation in loops
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case item := <-items:
        process(item)
    }
}
```

---

## 6. Common Pitfalls

```
Pitfalls to Avoid:
□ Storing context in structs
□ Not checking ctx.Done() in long operations
□ Not calling cancel() (goroutine leak)
□ Using string keys (collision risk)
□ Storing optional parameters in context
□ Passing nil context (use context.Background())
□ Not propagating context to child goroutines
```

---

## 7. Checklist

```
Context Usage Checklist:
□ Context is first parameter in functions
□ cancel() is always called (defer)
□ Typed keys for values
□ Context propagated to all downstream calls
□ ctx.Err() checked in loops
□ Timeouts configured appropriately
□ Request-scoped data only in context
□ No context stored in structs
```
