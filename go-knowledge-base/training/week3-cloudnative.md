# Week 3: Cloud-Native Patterns

## Module Overview

**Duration:** 40 hours (5 days)
**Prerequisites:** Week 2 completion (Concurrency)
**Learning Goal:** Build production-ready cloud-native applications with resilience and observability

---

## Learning Objectives

By the end of this week, you will be able to:

1. **Microservices Architecture**
   - Design microservices following Domain-Driven Design (DDD)
   - Implement proper service boundaries and APIs
   - Choose between REST and gRPC appropriately

2. **Resilience Patterns**
   - Implement circuit breakers for fault isolation
   - Design retry mechanisms with exponential backoff
   - Apply timeout and deadline patterns
   - Use bulkhead pattern for resource isolation

3. **Observability**
   - Implement structured logging
   - Add metrics with Prometheus
   - Configure distributed tracing with OpenTelemetry
   - Create health checks and readiness probes

4. **Context Management**
   - Propagate context through service calls
   - Implement request cancellation chains
   - Add context-aware logging and tracing

5. **Deployment**
   - Containerize Go applications
   - Write Kubernetes manifests
   - Configure graceful shutdown
   - Implement rolling updates

---

## Reading Assignments

### Required Reading (Complete by Day 3)

1. **[Microservices Patterns](../03-Engineering-CloudNative/EC-001-Microservices.md)**
   - Study: Service boundaries and decomposition
   - Learn: API Gateway and BFF patterns
   - Understand: Database per service pattern

2. **[Circuit Breaker Pattern](../03-Engineering-CloudNative/EC-001-Circuit-Breaker-Pattern.md)**
   - Master: Circuit breaker states and transitions
   - Learn: Half-open state testing
   - Study: Integration with observability

3. **[Context Management](../03-Engineering-CloudNative/EC-005-Context-Management.md)**
   - Understand: Request-scoped values
   - Learn: Cancellation propagation
   - Study: Context with timeout/deadline

4. **[Distributed Tracing](../03-Engineering-CloudNative/EC-006-Distributed-Tracing.md)**
   - Learn: Trace and span concepts
   - Study: Trace context propagation
   - Understand: Sampling strategies

5. **[Graceful Shutdown](../03-Engineering-CloudNative/EC-007-Graceful-Shutdown-Complete.md)**
   - Master: Signal handling
   - Learn: Drain in-flight requests
   - Study: Connection closure order

### Supplementary Reading (Complete by Day 5)

1. **[Retry Pattern](../03-Engineering-CloudNative/EC-002-Retry-Pattern.md)**
   - Understand: Exponential backoff
   - Learn: Jitter strategies
   - Study: Idempotency requirements

2. **[Rate Limiting Pattern](../03-Engineering-CloudNative/EC-005-Rate-Limiting-Pattern.md)**
   - Learn: Token bucket and leaky bucket
   - Study: Distributed rate limiting
   - Understand: Rate limit headers

3. **[OpenTelemetry Production](../03-Engineering-CloudNative/EC-060-OpenTelemetry-Distributed-Tracing-Production.md)**
   - Master: Instrumentation libraries
   - Learn: Collector configuration
   - Study: Production best practices

---

## Hands-on Exercises

### Day 1: RESTful API with Framework

#### Exercise 1.1: Echo Framework Mastery (3 hours)

Build a production-ready API with Echo:

```go
package main

import (
    "context"
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

// Server encapsulates the HTTP server
type Server struct {
    echo    *echo.Echo
    service *UserService
    config  Config
}

type Config struct {
    Port            string
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    ShutdownTimeout time.Duration
}

func NewServer(config Config, service *UserService) *Server {
    e := echo.New()

    s := &Server{
        echo:    e,
        service: service,
        config:  config,
    }

    s.setupMiddleware()
    s.setupRoutes()

    return s
}

func (s *Server) setupMiddleware() {
    // Recovery middleware
    s.echo.Use(middleware.Recover())

    // Request logging with structured format
    s.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: `{"time":"${time_rfc3339_nano}",` +
            `"level":"info",` +
            `"method":"${method}",` +
            `"uri":"${uri}",` +
            `"status":${status},` +
            `"latency":"${latency_human}",` +
            `"bytes_in":${bytes_in},` +
            `"bytes_out":${bytes_out}}` + "\n",
    }))

    // Request ID
    s.echo.Use(middleware.RequestID())

    // CORS
    s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
    }))

    // Timeout
    s.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
        Timeout: 30 * time.Second,
    }))

    // Custom context middleware
    s.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            ctx, cancel := context.WithCancel(c.Request().Context())
            defer cancel()

            c.SetRequest(c.Request().WithContext(ctx))
            return next(c)
        }
    })
}

func (s *Server) setupRoutes() {
    api := s.echo.Group("/api/v1")

    // Health check
    api.GET("/health", s.handleHealth)

    // User routes
    users := api.Group("/users")
    users.GET("", s.listUsers)
    users.POST("", s.createUser)
    users.GET("/:id", s.getUser)
    users.PUT("/:id", s.updateUser)
    users.DELETE("/:id", s.deleteUser)
}

type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name" validate:"required"`
    Email     string    `json:"email" validate:"required,email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func (s *Server) handleHealth(c echo.Context) error {
    return c.JSON(http.StatusOK, APIResponse{
        Success: true,
        Data: map[string]string{
            "status": "healthy",
            "time":   time.Now().Format(time.RFC3339),
        },
    })
}

func (s *Server) createUser(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, APIResponse{
            Success: false,
            Error: &APIError{
                Code:    "INVALID_REQUEST",
                Message: err.Error(),
            },
        })
    }

    if err := c.Validate(req); err != nil {
        return c.JSON(http.StatusBadRequest, APIResponse{
            Success: false,
            Error: &APIError{
                Code:    "VALIDATION_ERROR",
                Message: err.Error(),
            },
        })
    }

    user, err := s.service.Create(c.Request().Context(), req)
    if err != nil {
        return s.handleError(c, err)
    }

    return c.JSON(http.StatusCreated, APIResponse{
        Success: true,
        Data:    user,
    })
}

func (s *Server) handleError(c echo.Context, err error) error {
    // Map domain errors to HTTP status codes
    var status int
    var code string

    switch {
    case errors.Is(err, ErrNotFound):
        status = http.StatusNotFound
        code = "NOT_FOUND"
    case errors.Is(err, ErrConflict):
        status = http.StatusConflict
        code = "CONFLICT"
    case errors.Is(err, ErrUnauthorized):
        status = http.StatusUnauthorized
        code = "UNAUTHORIZED"
    default:
        status = http.StatusInternalServerError
        code = "INTERNAL_ERROR"
    }

    return c.JSON(status, APIResponse{
        Success: false,
        Error: &APIError{
            Code:    code,
            Message: err.Error(),
        },
    })
}

func (s *Server) Start() error {
    s.echo.Server.ReadTimeout = s.config.ReadTimeout
    s.echo.Server.WriteTimeout = s.config.WriteTimeout

    return s.echo.Start(s.config.Port)
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.echo.Shutdown(ctx)
}
```

**Tasks:**

1. Add JWT authentication middleware
2. Implement request validation
3. Add rate limiting middleware
4. Create OpenAPI/Swagger documentation
5. Implement pagination for list endpoints

**Deliverable:** Production-ready REST API with tests

#### Exercise 1.2: gRPC Service (2 hours)

Implement a gRPC service:

```protobuf
syntax = "proto3";
package users;
option go_package = "github.com/example/users";

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc UpdateUser(UpdateUserRequest) returns (User);
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);

    // Streaming example
    rpc StreamUsers(StreamUsersRequest) returns (stream User);
}

message User {
    string id = 1;
    string name = 2;
    string email = 3;
    string created_at = 4;
    string updated_at = 5;
}

message GetUserRequest {
    string id = 1;
}

message ListUsersRequest {
    int32 page_size = 1;
    string page_token = 2;
}

message ListUsersResponse {
    repeated User users = 1;
    string next_page_token = 2;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
}

message UpdateUserRequest {
    string id = 1;
    string name = 2;
    string email = 3;
}

message DeleteUserRequest {
    string id = 1;
}

message DeleteUserResponse {
    bool success = 1;
}

message StreamUsersRequest {
    string filter = 1;
}
```

**Implementation:**

```go
package grpc

import (
    "context"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type Server struct {
    pb.UnimplementedUserServiceServer
    service *UserService
}

func NewServer(service *UserService) *Server {
    return &Server{service: service}
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, err := s.service.Get(ctx, req.Id)
    if err != nil {
        return nil, s.mapError(err)
    }
    return toProtoUser(user), nil
}

func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    users, nextToken, err := s.service.List(ctx, int(req.PageSize), req.PageToken)
    if err != nil {
        return nil, s.mapError(err)
    }

    return &pb.ListUsersResponse{
        Users:         toProtoUsers(users),
        NextPageToken: nextToken,
    }, nil
}

func (s *Server) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
    ctx := stream.Context()

    users, err := s.service.Stream(ctx, req.Filter)
    if err != nil {
        return s.mapError(err)
    }

    for user := range users {
        if err := stream.Send(toProtoUser(user)); err != nil {
            return err
        }
    }

    return nil
}

func (s *Server) mapError(err error) error {
    switch {
    case errors.Is(err, ErrNotFound):
        return status.Error(codes.NotFound, err.Error())
    case errors.Is(err, ErrInvalidInput):
        return status.Error(codes.InvalidArgument, err.Error())
    case errors.Is(err, ErrUnauthorized):
        return status.Error(codes.Unauthenticated, err.Error())
    default:
        return status.Error(codes.Internal, "internal error")
    }
}
```

**Deliverable:** gRPC server with reflection and health checking

---

### Day 2: Resilience Patterns

#### Exercise 2.1: Circuit Breaker Implementation (3 hours)

Build a robust circuit breaker:

```go
package resilience

import (
    "context"
    "errors"
    "sync"
    "sync/atomic"
    "time"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type State int32

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

type Config struct {
    FailureThreshold    int
    SuccessThreshold    int
    Timeout             time.Duration
    HalfOpenMaxRequests int
}

type CircuitBreaker struct {
    config     Config
    state      int32 // atomic
    failures   int32 // atomic
    successes  int32 // atomic
    lastFailureTime int64 // atomic (UnixNano)

    mu        sync.RWMutex
    halfOpenCount int32 // atomic
}

func NewCircuitBreaker(config Config) *CircuitBreaker {
    return &CircuitBreaker{
        config: config,
    }
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
    state := cb.currentState()

    switch state {
    case StateOpen:
        if cb.shouldAttemptReset() {
            cb.setState(StateHalfOpen)
        } else {
            return ErrCircuitOpen
        }
    case StateHalfOpen:
        if !cb.canAttempt() {
            return ErrCircuitOpen
        }
    }

    err := fn()
    cb.recordResult(err)

    return err
}

func (cb *CircuitBreaker) currentState() State {
    return State(atomic.LoadInt32(&cb.state))
}

func (cb *CircuitBreaker) setState(state State) {
    atomic.StoreInt32(&cb.state, int32(state))

    switch state {
    case StateClosed:
        atomic.StoreInt32(&cb.failures, 0)
        atomic.StoreInt32(&cb.successes, 0)
    case StateOpen:
        atomic.StoreInt32(&cb.halfOpenCount, 0)
    case StateHalfOpen:
        atomic.StoreInt32(&cb.halfOpenCount, 0)
    }
}

func (cb *CircuitBreaker) shouldAttemptReset() bool {
    lastFailure := atomic.LoadInt64(&cb.lastFailureTime)
    return time.Since(time.Unix(0, lastFailure)) > cb.config.Timeout
}

func (cb *CircuitBreaker) canAttempt() bool {
    count := atomic.AddInt32(&cb.halfOpenCount, 1)
    return count <= int32(cb.config.HalfOpenMaxRequests)
}

func (cb *CircuitBreaker) recordResult(err error) {
    state := cb.currentState()

    if err != nil {
        atomic.StoreInt64(&cb.lastFailureTime, time.Now().UnixNano())

        switch state {
        case StateClosed:
            failures := atomic.AddInt32(&cb.failures, 1)
            if int(failures) >= cb.config.FailureThreshold {
                cb.setState(StateOpen)
            }
        case StateHalfOpen:
            cb.setState(StateOpen)
        }
    } else {
        switch state {
        case StateHalfOpen:
            successes := atomic.AddInt32(&cb.successes, 1)
            if int(successes) >= cb.config.SuccessThreshold {
                cb.setState(StateClosed)
            }
        case StateClosed:
            atomic.StoreInt32(&cb.failures, 0)
        }
    }
}

func (cb *CircuitBreaker) State() State {
    return cb.currentState()
}

func (cb *CircuitBreaker) Metrics() Metrics {
    return Metrics{
        State:     cb.currentState(),
        Failures:  int(atomic.LoadInt32(&cb.failures)),
        Successes: int(atomic.LoadInt32(&cb.successes)),
    }
}

type Metrics struct {
    State     State
    Failures  int
    Successes int
}
```

**Integration example:**

```go
// HTTP client with circuit breaker
type ResilientClient struct {
    client         *http.Client
    circuitBreaker *CircuitBreaker
}

func (c *ResilientClient) Do(req *http.Request) (*http.Response, error) {
    var resp *http.Response
    var err error

    cbErr := c.circuitBreaker.Execute(req.Context(), func() error {
        resp, err = c.client.Do(req)
        if err != nil {
            return err
        }

        // Treat 5xx as failures
        if resp.StatusCode >= 500 {
            return fmt.Errorf("server error: %d", resp.StatusCode)
        }

        return nil
    })

    if cbErr == ErrCircuitOpen {
        return nil, fmt.Errorf("service temporarily unavailable")
    }

    return resp, err
}
```

**Deliverable:** Circuit breaker with metrics and configuration

#### Exercise 2.2: Retry with Exponential Backoff (2 hours)

Implement retry mechanism:

```go
package resilience

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
    RetryableFn func(error) bool
}

func DefaultRetryConfig() RetryConfig {
    return RetryConfig{
        MaxRetries:  3,
        BaseDelay:   100 * time.Millisecond,
        MaxDelay:    30 * time.Second,
        Multiplier:  2.0,
        Jitter:      0.1,
        RetryableFn: IsRetryableError,
    }
}

func IsRetryableError(err error) bool {
    if err == nil {
        return false
    }

    // Retry on specific errors
    var retryable []error
    // Add specific retryable errors

    for _, r := range retryable {
        if errors.Is(err, r) {
            return true
        }
    }

    return false
}

func Retry(ctx context.Context, config RetryConfig, fn func() error) error {
    var err error

    for attempt := 0; attempt <= config.MaxRetries; attempt++ {
        err = fn()
        if err == nil {
            return nil
        }

        if attempt == config.MaxRetries {
            break
        }

        if config.RetryableFn != nil && !config.RetryableFn(err) {
            return err
        }

        delay := calculateDelay(attempt, config)

        select {
        case <-time.After(delay):
            continue
        case <-ctx.Done():
            return ctx.Err()
        }
    }

    return fmt.Errorf("max retries exceeded: %w", err)
}

func calculateDelay(attempt int, config RetryConfig) time.Duration {
    // Exponential backoff
    delay := float64(config.BaseDelay) * math.Pow(config.Multiplier, float64(attempt))

    // Cap at max delay
    if delay > float64(config.MaxDelay) {
        delay = float64(config.MaxDelay)
    }

    // Add jitter
    if config.Jitter > 0 {
        jitter := delay * config.Jitter * (rand.Float64()*2 - 1)
        delay += jitter
    }

    return time.Duration(delay)
}
```

**Deliverable:** Retry decorator for HTTP clients

---

### Day 3: Observability

#### Exercise 3.1: Structured Logging (2 hours)

Implement structured logging:

```go
package observability

import (
    "context"
    "encoding/json"
    "io"
    "os"
    "time"
)

type Level int

const (
    Debug Level = iota
    Info
    Warn
    Error
    Fatal
)

type Logger struct {
    writer io.Writer
    level  Level
    fields map[string]interface{}
}

func NewLogger(w io.Writer, level Level) *Logger {
    if w == nil {
        w = os.Stdout
    }
    return &Logger{
        writer: w,
        level:  level,
        fields: make(map[string]interface{}),
    }
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
    newFields := make(map[string]interface{})
    for k, v := range l.fields {
        newFields[k] = v
    }
    newFields[key] = value
    return &Logger{writer: l.writer, level: l.level, fields: newFields}
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
    newFields := make(map[string]interface{})
    for k, v := range l.fields {
        newFields[k] = v
    }
    for k, v := range fields {
        newFields[k] = v
    }
    return &Logger{writer: l.writer, level: l.level, fields: newFields}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
    // Extract trace ID, request ID, etc. from context
    newLogger := l

    if traceID := GetTraceID(ctx); traceID != "" {
        newLogger = newLogger.WithField("trace_id", traceID)
    }

    if requestID := GetRequestID(ctx); requestID != "" {
        newLogger = newLogger.WithField("request_id", requestID)
    }

    return newLogger
}

func (l *Logger) log(level Level, msg string, fields map[string]interface{}) {
    if level < l.level {
        return
    }

    entry := LogEntry{
        Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
        Level:     level.String(),
        Message:   msg,
        Fields:    make(map[string]interface{}),
    }

    // Add logger fields
    for k, v := range l.fields {
        entry.Fields[k] = v
    }

    // Add event fields
    for k, v := range fields {
        entry.Fields[k] = v
    }

    encoder := json.NewEncoder(l.writer)
    encoder.Encode(entry)
}

type LogEntry struct {
    Timestamp string                 `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Fields    map[string]interface{} `json:"fields,omitempty"`
}

func (l Level) String() string {
    switch l {
    case Debug:
        return "debug"
    case Info:
        return "info"
    case Warn:
        return "warn"
    case Error:
        return "error"
    case Fatal:
        return "fatal"
    default:
        return "unknown"
    }
}

func (l *Logger) Debug(msg string) { l.log(Debug, msg, nil) }
func (l *Logger) Info(msg string)  { l.log(Info, msg, nil) }
func (l *Logger) Warn(msg string)  { l.log(Warn, msg, nil) }
func (l *Logger) Error(msg string) { l.log(Error, msg, nil) }
func (l *Logger) Fatal(msg string) { l.log(Fatal, msg, nil) }
```

**Deliverable:** Logging package with context integration

#### Exercise 3.2: Metrics with Prometheus (2 hours)

Add Prometheus metrics:

```go
package observability

import (
    "context"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint", "status"},
    )

    requestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    activeRequests = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "http_active_requests",
            Help: "Number of active HTTP requests",
        },
        []string{"method", "endpoint"},
    )

    businessMetrics = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "business_events_total",
            Help: "Business events counter",
        },
        []string{"event_type"},
    )
)

func init() {
    prometheus.MustRegister(requestDuration)
    prometheus.MustRegister(requestCount)
    prometheus.MustRegister(activeRequests)
    prometheus.MustRegister(businessMetrics)
}

// MetricsMiddleware records HTTP metrics
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

        endpoint := r.URL.Path
        method := r.Method

        activeRequests.WithLabelValues(method, endpoint).Inc()
        defer activeRequests.WithLabelValues(method, endpoint).Dec()

        next.ServeHTTP(wrapped, r)

        duration := time.Since(start).Seconds()
        status := http.StatusText(wrapped.statusCode)

        requestDuration.WithLabelValues(method, endpoint, status).Observe(duration)
        requestCount.WithLabelValues(method, endpoint, status).Inc()
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

// MetricsHandler exposes Prometheus metrics
func MetricsHandler() http.Handler {
    return promhttp.Handler()
}
```

**Deliverable:** Instrumented service with metrics endpoint

---

### Day 4: Distributed Tracing

#### Exercise 4.1: OpenTelemetry Integration (3 hours)

Implement distributed tracing:

```go
package observability

import (
    "context"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
    "go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func InitTracer(serviceName, serviceVersion string) (*sdktrace.TracerProvider, error) {
    ctx := context.Background()

    // Create exporter
    client := otlptracegrpc.NewClient()
    exporter, err := otlptrace.New(ctx, client)
    if err != nil {
        return nil, err
    }

    // Create resource
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String(serviceName),
            semconv.ServiceVersionKey.String(serviceVersion),
        ),
    )
    if err != nil {
        return nil, err
    }

    // Create tracer provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)),
    )

    otel.SetTracerProvider(tp)
    tracer = tp.Tracer(serviceName)

    return tp, nil
}

// SpanFromContext creates a child span
type Span struct {
    trace.Span
    ctx context.Context
}

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, *Span) {
    ctx, span := tracer.Start(ctx, name, opts...)
    return ctx, &Span{Span: span, ctx: ctx}
}

func (s *Span) SetError(err error) {
    s.RecordError(err)
    s.SetStatus(codes.Error, err.Error())
}

func (s *Span) SetAttributes(attrs ...attribute.KeyValue) {
    s.Span.SetAttributes(attrs...)
}

func (s *Span) End() {
    s.Span.End()
}

// HTTP middleware for trace context propagation
func TracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, span := StartSpan(r.Context(), r.URL.Path,
            trace.WithAttributes(
                attribute.String("http.method", r.Method),
                attribute.String("http.url", r.URL.String()),
                attribute.String("http.user_agent", r.UserAgent()),
            ),
        )
        defer span.End()

        // Add trace context to response headers
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Deliverable:** Traced service with Jaeger visualization

---

### Day 5: Deployment

#### Exercise 5.1: Containerization (2 hours)

Create production Dockerfile:

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Copy config if needed
COPY --from=builder /app/configs ./configs

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run as non-root user
RUN adduser -D -s /bin/sh appuser
USER appuser

CMD ["./server"]
```

**Deliverable:** Multi-stage Dockerfile with security best practices

#### Exercise 5.2: Kubernetes Manifests (3 hours)

Create K8s deployment:

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: user-service:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: PORT
          value: "8080"
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        lifecycle:
          preStop:
            exec:
              command: ["/bin/sh", "-c", "sleep 10"]
      terminationGracePeriodSeconds: 30

---
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: user-service
  annotations:
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /users
        pathType: Prefix
        backend:
          service:
            name: user-service
            port:
              number: 80
```

**Deliverable:** Complete K8s manifests with HPA

---

## Code Review Checklist

### Cloud-Native Best Practices

- [ ] Context propagation throughout call chain
- [ ] Graceful shutdown handling
- [ ] Health checks implemented
- [ ] Structured logging with correlation IDs
- [ ] Metrics exposed in Prometheus format
- [ ] Tracing spans properly nested

### Resilience

- [ ] Circuit breaker on external calls
- [ ] Retry with exponential backoff
- [ ] Request timeouts configured
- [ ] Rate limiting applied
- [ ] Bulkhead pattern where needed

### Security

- [ ] No secrets in code or images
- [ ] Non-root container user
- [ ] Security headers set
- [ ] Input validation on all endpoints
- [ ] TLS configured

### Observability

- [ ] Structured logging (JSON format)
- [ ] Appropriate log levels
- [ ] Key metrics instrumented
- [ ] Distributed tracing enabled
- [ ] Error tracking integrated

---

## Assessment Criteria

### Knowledge Assessment (30%)

- Microservices design principles
- Resilience patterns and trade-offs
- Observability concepts (metrics, logs, traces)
- Kubernetes deployment concepts

**Passing Score:** 80%

### Coding Challenge (50%)

**Problem:** Build a resilient payment service

**Requirements:**

- REST API with validation
- Circuit breaker on payment gateway
- Retry with idempotency keys
- Prometheus metrics
- OpenTelemetry tracing
- Graceful shutdown
- Kubernetes deployment

### System Design (20%)

**Problem:** Design observability for a microservices platform

---

*Next: [Week 4: System Design](week4-systemdesign.md)*
