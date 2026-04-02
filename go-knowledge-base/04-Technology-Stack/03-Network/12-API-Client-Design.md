# TS-NET-012: API Client Design Patterns in Go

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #api-client #http-client #golang #resilience #patterns #circuit-breaker
> **权威来源**:
> - [Go HTTP Client Best Practices](https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779) - Medium
> - [Resilience Patterns](https://docs.microsoft.com/en-us/azure/architecture/patterns/category/resiliency) - Microsoft Azure

---

## 1. API Client Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         API Client Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      API Client                                     │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐ │   │
│  │  │  Request    │  │  Circuit    │  │   Retry     │  │   Timeout  │ │   │
│  │  │  Builder    │──►│  Breaker    │──►│   Logic     │──►│   Handler  │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────┬──────┘ │   │
│  │                                                          │        │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │        │   │
│  │  │  Auth       │  │  Logging    │  │  Metrics    │      │        │   │
│  │  │  Handler    │  │  Handler    │  │  Handler    │      │        │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘      │        │   │
│  │                                                          │        │   │
│  └──────────────────────────────────────────────────────────┼────────┘   │
│                                                             │              │
│  ┌──────────────────────────────────────────────────────────┼────────┐   │
│  │                      HTTP Client                          │        │   │
│  │  ┌───────────────────────────────────────────────────────┼──────┐ │   │
│  │  │                   Connection Pool                      │      │ │   │
│  │  │  ┌──────────┐  ┌──────────┐  ┌──────────┐             │      │ │   │
│  │  │  │ Conn 1   │  │ Conn 2   │  │ Conn N   │             │      │ │   │
│  │  │  │ (Active) │  │ (Idle)   │  │ (Active) │             │      │ │   │
│  │  │  └──────────┘  └──────────┘  └──────────┘             │      │ │   │
│  │  └───────────────────────────────────────────────────────┼──────┘ │   │
│  └──────────────────────────────────────────────────────────┼────────┘   │
│                                                             │              │
│                                                        ┌────┴────┐        │
│                                                        │   API   │        │
│                                                        └─────────┘        │
│                                                                              │
│  Resilience Patterns:                                                        │
│  - Circuit Breaker: Fail fast when service is unhealthy                     │
│  - Retry: Exponential backoff for transient failures                        │
│  - Timeout: Limit request duration                                          │
│  - Bulkhead: Isolate failures                                               │
│  - Rate Limiting: Prevent overwhelming the service                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Complete Client Implementation

```go
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"
)

// Config for API client
type Config struct {
    BaseURL            string
    Timeout            time.Duration
    MaxRetries         int
    RetryDelay         time.Duration
    MaxConnsPerHost    int
    InsecureSkipVerify bool
}

// Client is the API client
type Client struct {
    httpClient *http.Client
    baseURL    string
    config     Config
    auth       AuthProvider
}

// AuthProvider handles authentication
type AuthProvider interface {
    Apply(req *http.Request) error
    Refresh(ctx context.Context) error
}

// New creates a new API client
func New(config Config, auth AuthProvider) (*Client, error) {
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: config.MaxConnsPerHost,
        MaxConnsPerHost:     config.MaxConnsPerHost,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    }
    
    return &Client{
        httpClient: &http.Client{
            Timeout:   config.Timeout,
            Transport: transport,
        },
        baseURL: config.BaseURL,
        config:  config,
        auth:    auth,
    }, nil
}

// Request builder
type RequestBuilder struct {
    client  *Client
    method  string
    path    string
    query   url.Values
    headers http.Header
    body    interface{}
}

func (c *Client) Request(method, path string) *RequestBuilder {
    return &RequestBuilder{
        client:  c,
        method:  method,
        path:    path,
        query:   make(url.Values),
        headers: make(http.Header),
    }
}

func (rb *RequestBuilder) WithQuery(key, value string) *RequestBuilder {
    rb.query.Add(key, value)
    return rb
}

func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
    rb.headers.Set(key, value)
    return rb
}

func (rb *RequestBuilder) WithBody(body interface{}) *RequestBuilder {
    rb.body = body
    return rb
}

func (rb *RequestBuilder) Execute(ctx context.Context, result interface{}) error {
    return rb.client.do(ctx, rb, result)
}

// Do the actual request with retries
func (c *Client) do(ctx context.Context, rb *RequestBuilder, result interface{}) error {
    var lastErr error
    
    for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
        if attempt > 0 {
            time.Sleep(c.config.RetryDelay * time.Duration(attempt))
        }
        
        err := c.executeOnce(ctx, rb, result)
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // Don't retry on client errors (4xx)
        if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode >= 400 && apiErr.StatusCode < 500 {
            return err
        }
    }
    
    return fmt.Errorf("max retries exceeded: %w", lastErr)
}

func (c *Client) executeOnce(ctx context.Context, rb *RequestBuilder, result interface{}) error {
    // Build URL
    u, err := url.Parse(c.baseURL + rb.path)
    if err != nil {
        return err
    }
    u.RawQuery = rb.query.Encode()
    
    // Build body
    var bodyReader io.Reader
    if rb.body != nil {
        bodyBytes, err := json.Marshal(rb.body)
        if err != nil {
            return err
        }
        bodyReader = bytes.NewReader(bodyBytes)
    }
    
    // Create request
    req, err := http.NewRequestWithContext(ctx, rb.method, u.String(), bodyReader)
    if err != nil {
        return err
    }
    
    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    for key, values := range rb.headers {
        for _, value := range values {
            req.Header.Add(key, value)
        }
    }
    
    // Apply authentication
    if c.auth != nil {
        if err := c.auth.Apply(req); err != nil {
            return err
        }
    }
    
    // Execute
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Read response body
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    // Check status
    if resp.StatusCode >= 400 {
        return &APIError{
            StatusCode: resp.StatusCode,
            Body:       respBody,
            Message:    fmt.Sprintf("API error: %s", resp.Status),
        }
    }
    
    // Parse result
    if result != nil && len(respBody) > 0 {
        if err := json.Unmarshal(respBody, result); err != nil {
            return err
        }
    }
    
    return nil
}

// APIError represents an API error
type APIError struct {
    StatusCode int
    Body       []byte
    Message    string
}

func (e *APIError) Error() string {
    return e.Message
}
```

---

## 3. Circuit Breaker Pattern

```go
package client

import (
    "errors"
    "sync"
    "time"
)

type CircuitState int

const (
    StateClosed CircuitState = iota    // Normal operation
    StateOpen                          // Failing, reject requests
    StateHalfOpen                      // Testing if recovered
)

type CircuitBreaker struct {
    failureThreshold int
    successThreshold int
    timeout          time.Duration
    
    state          CircuitState
    failures       int
    successes      int
    lastFailure    time.Time
    mu             sync.Mutex
}

func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        failureThreshold: failureThreshold,
        successThreshold: successThreshold,
        timeout:          timeout,
        state:            StateClosed,
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    
    switch cb.state {
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
            cb.failures = 0
            cb.successes = 0
        } else {
            cb.mu.Unlock()
            return errors.New("circuit breaker open")
        }
    case StateHalfOpen:
        // Continue to execute
    }
    
    cb.mu.Unlock()
    
    err := fn()
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.recordFailure()
    } else {
        cb.recordSuccess()
    }
    
    return err
}

func (cb *CircuitBreaker) recordFailure() {
    cb.failures++
    cb.lastFailure = time.Now()
    
    switch cb.state {
    case StateClosed:
        if cb.failures >= cb.failureThreshold {
            cb.state = StateOpen
        }
    case StateHalfOpen:
        cb.state = StateOpen
    }
}

func (cb *CircuitBreaker) recordSuccess() {
    cb.successes++
    
    switch cb.state {
    case StateHalfOpen:
        if cb.successes >= cb.successThreshold {
            cb.state = StateClosed
            cb.failures = 0
        }
    case StateClosed:
        cb.failures = 0
    }
}

func (cb *CircuitBreaker) State() CircuitState {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    return cb.state
}
```

---

## 4. Retry with Exponential Backoff

```go
package client

import (
    "context"
    "math"
    "math/rand"
    "time"
)

type RetryConfig struct {
    MaxRetries  int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Multiplier  float64
    Jitter      float64
}

func DefaultRetryConfig() RetryConfig {
    return RetryConfig{
        MaxRetries: 3,
        BaseDelay:  100 * time.Millisecond,
        MaxDelay:   30 * time.Second,
        Multiplier: 2.0,
        Jitter:     0.1,
    }
}

func (rc *RetryConfig) CalculateDelay(attempt int) time.Duration {
    if attempt <= 0 {
        return 0
    }
    
    // Exponential backoff
    delay := float64(rc.BaseDelay) * math.Pow(rc.Multiplier, float64(attempt-1))
    
    // Cap at max delay
    if delay > float64(rc.MaxDelay) {
        delay = float64(rc.MaxDelay)
    }
    
    // Add jitter
    if rc.Jitter > 0 {
        jitter := delay * rc.Jitter * (rand.Float64()*2 - 1)
        delay += jitter
    }
    
    return time.Duration(delay)
}

func RetryWithConfig(ctx context.Context, config RetryConfig, fn func() error) error {
    var lastErr error
    
    for attempt := 0; attempt <= config.MaxRetries; attempt++ {
        if attempt > 0 {
            delay := config.CalculateDelay(attempt)
            
            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return ctx.Err()
            }
        }
        
        err := fn()
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // Check if error is retryable
        if !isRetryableError(err) {
            return err
        }
    }
    
    return lastErr
}

func isRetryableError(err error) bool {
    if err == nil {
        return false
    }
    
    // Check for specific retryable errors
    errStr := err.Error()
    retryableErrors := []string{
        "connection refused",
        "connection reset",
        "broken pipe",
        "timeout",
        "temporary",
        "too many requests",
        "service unavailable",
    }
    
    for _, retryable := range retryableErrors {
        if contains(errStr, retryable) {
            return true
        }
    }
    
    return false
}
```

---

## 5. Checklist

```
API Client Checklist:
□ Connection pooling configured
□ Timeout set appropriately
□ Retry logic with exponential backoff
□ Circuit breaker for resilience
□ Authentication handling
□ Request/response logging
□ Metrics collection
□ Proper error handling
□ Context support for cancellation
□ User-Agent header set
□ Compression enabled
□ Request/response interceptors
```
