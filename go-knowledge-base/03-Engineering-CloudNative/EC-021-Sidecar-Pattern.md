# EC-021: Sidecar Pattern

## Problem Formalization

### The Cross-Cutting Concern Challenge

Microservices require common functionality (logging, monitoring, configuration, networking) that should not pollute the core business logic. Traditional approaches either duplicate code across services or create tight coupling through shared libraries.

#### Problem Statement

Given:
- A set of microservices S = {s₁, s₂, ..., sₙ}
- Cross-cutting concerns C = {c₁, c₂, ..., cₘ} where each cᵢ includes:
  - Observability (metrics, logging, tracing)
  - Security (mTLS, authentication)
  - Communication (service discovery, load balancing)
  - Resilience (circuit breaking, retries)

Find a deployment model D such that:
```
Minimize: CodeDuplication(S, C)
Minimize: Coupling(S, C)
Maximize: Independence(sᵢ) for each sᵢ ∈ S
Subject to:
  - Each sᵢ has access to all c ∈ C
  - Upgrades to C don't require redeploying S
  - Failure in c doesn't cascade to s
```

### Monolithic vs Sidecar Approach

```
Without Sidecar (Shared Library Approach):
┌─────────────────────────────────────────────────────────────────────────┐
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                  │
│  │  Service A   │  │  Service B   │  │  Service C   │                  │
│  │  ┌────────┐  │  │  ┌────────┐  │  │  ┌────────┐  │                  │
│  │  │Business│  │  │  │Business│  │  │  │Business│  │                  │
│  │  │ Logic  │  │  │  │ Logic  │  │  │  │ Logic  │  │                  │
│  │  └────────┘  │  │  └────────┘  │  │  └────────┘  │                  │
│  │  ┌────────┐  │  │  ┌────────┐  │  │  ┌────────┐  │                  │
│  │  │Metrics │  │  │  │Metrics │  │  │  │Metrics │  │                  │
│  │  │Logging │  │  │  │Logging │  │  │  │Logging │  │                  │
│  │  │Tracing │  │  │  │Tracing │  │  │  │Tracing │  │                  │
│  │  │SSL/TLS │  │  │  │SSL/TLS │  │  │  │SSL/TLS │  │                  │
│  │  │Retries │  │  │  │Retries │  │  │  │Retries │  │                  │
│  │  └────────┘  │  │  └────────┘  │  │  └────────┘  │                  │
│  └──────────────┘  └──────────────┘  └──────────────┘                  │
│                                                                         │
│  Problems:                                                              │
│  • Library version conflicts                                            │
│  • Language-specific implementations                                    │
│  • Upgrades require redeploying all services                            │
│  • Different configurations per service                                 │
└─────────────────────────────────────────────────────────────────────────┘

With Sidecar Pattern:
┌─────────────────────────────────────────────────────────────────────────┐
│  ┌────────────────────┐  ┌────────────────────┐  ┌────────────────────┐│
│  │  Pod/Container     │  │  Pod/Container     │  │  Pod/Container     ││
│  │  ┌──────────────┐  │  │  ┌──────────────┐  │  │  ┌──────────────┐  ││
│  │  │   Service A  │  │  │  │   Service B  │  │  │  │   Service C  │  ││
│  │  │   (Business) │  │  │  │   (Business) │  │  │  │   (Business) │  ││
│  │  │   localhost  │  │  │  │   localhost  │  │  │  │   localhost  │  ││
│  │  └──────┬───────┘  │  │  └──────┬───────┘  │  │  └──────┬───────┘  ││
│  │         │          │  │         │          │  │         │          ││
│  │  ┌──────┴───────┐  │  │  ┌──────┴───────┐  │  │  ┌──────┴───────┐  ││
│  │  │   Sidecar    │  │  │  │   Sidecar    │  │  │  │   Sidecar    │  ││
│  │  │   (Envoy/    │  │  │  │   (Envoy/    │  │  │  │   (Envoy/    │  ││
│  │  │   Linkerd)   │  │  │  │   Linkerd)   │  │  │  │   Linkerd)   │  ││
│  │  │   mTLS       │  │  │  │   mTLS       │  │  │  │   mTLS       │  ││
│  │  │   Metrics    │  │  │  │   Metrics    │  │  │  │   Metrics    │  ││
│  │  │   Retry      │  │  │  │   Retry      │  │  │  │   Retry      │  ││
│  │  └──────────────┘  │  │  └──────────────┘  │  │  └──────────────┘  ││
│  └────────────────────┘  └────────────────────┘  └────────────────────┘│
│                                                                         │
│  Benefits:                                                              │
│  • Language-agnostic                                                    │
│  • Independent upgrades                                                 │
│  • Consistent configuration                                             │
│  • Failure isolation                                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Sidecar Deployment Model

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     Sidecar Communication Patterns                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Pattern 1: Localhost Networking (Most Common)                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                            │   │
│  │  ┌──────────────┐         ┌──────────────┐                     │   │
│  │  │   Service    │◄───────►│   Sidecar    │                     │   │
│  │  │              │  :8080  │              │                     │   │
│  │  │ (App Port)   │         │ (Proxy Port) │                     │   │
│  │  └──────────────┘         └──────┬───────┘                     │   │
│  │                                  │                              │   │
│  │                          Network Interface                      │   │
│  │                                  │                              │   │
│  │                                  ▼                              │   │
│  │                         Other Services                         │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Pattern 2: Unix Domain Sockets (Higher Performance)                    │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                            │   │
│  │  ┌──────────────┐         ┌──────────────┐                     │   │
│  │  │   Service    │◄───────►│   Sidecar    │                     │   │
│  │  │              │/tmp/app.sock                                │   │
│  │  └──────────────┘         └──────┬───────┘                     │   │
│  │                                  │                              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Pattern 3: Shared Memory (Lowest Latency)                              │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Pod                                                            │   │
│  │  ┌──────────────┐         ┌──────────────┐                     │   │
│  │  │   Service    │◄───────►│   Sidecar    │                     │   │
│  │  │              │  shm    │              │                     │   │
│  │  └──────────────┘         └──────────────┘                     │   │
│  │         ▲                        ▲                              │   │
│  │         └──────────/dev/shm──────┘                              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Service Mesh with Sidecars

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Service Mesh Architecture                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Control Plane (Istio/Linkerd)                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │   │
│  │  │   Pilot      │  │  Citadel     │  │   Galley             │  │   │
│  │  │   (Config)   │  │  (Certs)     │  │   (Validation)       │  │   │
│  │  └──────────────┘  └──────────────┘  └──────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│         │                          │                          │        │
│         │ xDS Protocol             │ Certificate Distribution  │        │
│         │ (Envoy API)              │ (SPIFFE/SPIRE)            │        │
│         ▼                          ▼                          ▼        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Data Plane (Envoy Sidecars)                  │   │
│  │                                                                  │   │
│  │  ┌──────────┐      ┌──────────┐      ┌──────────┐              │   │
│  │  │ Service A│◄────►│ Service B│◄────►│ Service C│              │   │
│  │  │  + Envoy │ mTLS │  + Envoy │ mTLS │  + Envoy │              │   │
│  │  └──────────┘      └──────────┘      └──────────┘              │   │
│  │       │                 │                 │                     │   │
│  │       └─────────────────┼─────────────────┘                     │   │
│  │                         │                                       │   │
│  │                    Service Discovery                            │   │
│  │                    (Kubernetes DNS)                             │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Sidecar Capabilities:                                                  │
│  • Dynamic service discovery                                            │
│  • Load balancing (least_request, ring_hash, random)                    │
│  • Health checks                                                        │
│  • mTLS encryption                                                      │
│  • Circuit breaking                                                     │
│  • Fault injection (chaos testing)                                      │
│  • Rich metrics (Prometheus)                                            │
│  • Distributed tracing (Zipkin/Jaeger)                                  │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Custom Sidecar Implementation

```go
// cmd/sidecar/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "go.uber.org/zap"
)

// SidecarConfig holds all sidecar configuration
type SidecarConfig struct {
    // Application connection
    AppHost string
    AppPort int
    
    // Sidecar listener
    ListenPort    int
    AdminPort     int
    MetricsPort   int
    
    // TLS configuration
    TLSCertPath   string
    TLSKeyPath    string
    CAPath        string
    
    // Resilience settings
    CircuitBreaker CircuitBreakerConfig
    Retry          RetryConfig
    Timeout        TimeoutConfig
    
    // Observability
    LogLevel      string
    TraceEndpoint string
}

type CircuitBreakerConfig struct {
    Enabled           bool
    FailureThreshold  int
    SuccessThreshold  int
    Timeout           time.Duration
}

type RetryConfig struct {
    Enabled      bool
    MaxAttempts  int
    BaseDelay    time.Duration
    MaxDelay     time.Duration
}

type TimeoutConfig struct {
    Request  time.Duration
    Idle     time.Duration
}

// Sidecar implements the sidecar proxy functionality
type Sidecar struct {
    config      *SidecarConfig
    logger      *zap.Logger
    
    // Proxies
    appProxy    *httputil.ReverseProxy
    
    // Components
    metrics     *Metrics
    tracer      *Tracer
    breaker     *CircuitBreaker
    retrier     *Retrier
    authenticator *Authenticator
    
    // Servers
    proxyServer  *http.Server
    adminServer  *http.Server
    metricsServer *http.Server
    
    ctx         context.Context
    cancel      context.CancelFunc
}

func NewSidecar(cfg *SidecareConfig) (*Sidecar, error) {
    logger, _ := zap.NewProduction()
    
    ctx, cancel := context.WithCancel(context.Background())
    
    s := &Sidecar{
        config: cfg,
        logger: logger,
        ctx:    ctx,
        cancel: cancel,
    }
    
    // Setup application proxy
    appURL := &url.URL{
        Scheme: "http",
        Host:   fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort),
    }
    s.appProxy = httputil.NewSingleHostReverseProxy(appURL)
    s.appProxy.ErrorHandler = s.handleProxyError
    
    // Initialize components
    s.metrics = NewMetrics()
    s.tracer = NewTracer(cfg.TraceEndpoint)
    
    if cfg.CircuitBreaker.Enabled {
        s.breaker = NewCircuitBreaker(cfg.CircuitBreaker)
    }
    
    if cfg.Retry.Enabled {
        s.retrier = NewRetrier(cfg.Retry)
    }
    
    s.setupServers()
    
    return s, nil
}

func (s *Sidecar) setupServers() {
    // Main proxy server
    proxyMux := http.NewServeMux()
    proxyMux.Handle("/", s.wrapMiddleware(s.appProxy))
    
    s.proxyServer = &http.Server{
        Addr:         fmt.Sprintf(":%d", s.config.ListenPort),
        Handler:      proxyMux,
        ReadTimeout:  s.config.Timeout.Request,
        WriteTimeout: s.config.Timeout.Request,
        IdleTimeout:  s.config.Timeout.Idle,
    }
    
    // Admin server (health, config)
    adminMux := http.NewServeMux()
    adminMux.HandleFunc("/health", s.healthHandler)
    adminMux.HandleFunc("/ready", s.readyHandler)
    adminMux.HandleFunc("/config", s.configHandler)
    
    s.adminServer = &http.Server{
        Addr:    fmt.Sprintf(":%d", s.config.AdminPort),
        Handler: adminMux,
    }
    
    // Metrics server (Prometheus)
    metricsMux := http.NewServeMux()
    metricsMux.Handle("/metrics", promhttp.Handler())
    
    s.metricsServer = &http.Server{
        Addr:    fmt.Sprintf(":%d", s.config.MetricsPort),
        Handler: metricsMux,
    }
}

// wrapMiddleware applies all sidecar middleware
func (s *Sidecar) wrapMiddleware(handler http.Handler) http.Handler {
    h := handler
    
    // Order matters: outside-in for request, inside-out for response
    h = s.loggingMiddleware(h)
    h = s.tracingMiddleware(h)
    h = s.metricsMiddleware(h)
    h = s.authMiddleware(h)
    h = s.circuitBreakerMiddleware(h)
    h = s.retryMiddleware(h)
    h = s.timeoutMiddleware(h)
    
    return h
}

func (s *Sidecar) loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status
        wrapped := &responseRecorder{ResponseWriter: w}
        
        next.ServeHTTP(wrapped, r)
        
        s.logger.Info("request",
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.Int("status", wrapped.statusCode),
            zap.Duration("duration", time.Since(start)),
            zap.String("trace_id", r.Header.Get("X-Trace-ID")),
        )
    })
}

func (s *Sidecar) tracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract or create trace context
        ctx := s.tracer.Extract(r.Context(), r.Header)
        
        // Start span
        span, ctx := s.tracer.StartSpan(ctx, "sidecar_proxy")
        defer span.Finish()
        
        // Add trace headers for application
        s.tracer.Inject(ctx, r.Header)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (s *Sidecar) circuitBreakerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if s.breaker == nil {
            next.ServeHTTP(w, r)
            return
        }
        
        if !s.breaker.Allow() {
            s.metrics.RecordCircuitBreakerOpen()
            http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
            return
        }
        
        wrapped := &responseRecorder{ResponseWriter: w}
        next.ServeHTTP(wrapped, r)
        
        // Record result
        if wrapped.statusCode >= 500 {
            s.breaker.RecordFailure()
        } else {
            s.breaker.RecordSuccess()
        }
    })
}

func (s *Sidecar) retryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if s.retrier == nil {
            next.ServeHTTP(w, r)
            return
        }
        
        var lastErr error
        
        for attempt := 0; attempt < s.retrier.MaxAttempts; attempt++ {
            if attempt > 0 {
                // Calculate backoff
                delay := s.retrier.CalculateBackoff(attempt)
                time.Sleep(delay)
                
                // Clone request for retry
                r = r.Clone(r.Context())
            }
            
            wrapped := &responseRecorder{ResponseWriter: w}
            next.ServeHTTP(wrapped, r)
            
            // Success
            if wrapped.statusCode < 500 {
                if attempt > 0 {
                    s.metrics.RecordRetrySuccess(attempt)
                }
                return
            }
            
            // Check if retryable
            if !s.retrier.IsRetryable(wrapped.statusCode) {
                return
            }
            
            lastErr = fmt.Errorf("attempt %d failed with status %d", attempt+1, wrapped.statusCode)
        }
        
        s.metrics.RecordRetryExhausted()
        s.logger.Warn("retry exhausted",
            zap.Error(lastErr),
            zap.String("path", r.URL.Path),
        )
    })
}

func (s *Sidecar) Run() error {
    errCh := make(chan error, 3)
    
    // Start servers
    go func() {
        s.logger.Info("starting proxy server", zap.Int("port", s.config.ListenPort))
        errCh <- s.proxyServer.ListenAndServe()
    }()
    
    go func() {
        s.logger.Info("starting admin server", zap.Int("port", s.config.AdminPort))
        errCh <- s.adminServer.ListenAndServe()
    }()
    
    go func() {
        s.logger.Info("starting metrics server", zap.Int("port", s.config.MetricsPort))
        errCh <- s.metricsServer.ListenAndServe()
    }()
    
    // Wait for shutdown signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    select {
    case err := <-errCh:
        return err
    case <-sigCh:
        return s.Shutdown()
    }
}

func (s *Sidecar) Shutdown() error {
    s.logger.Info("shutting down sidecar...")
    s.cancel()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    var errs []error
    
    if err := s.proxyServer.Shutdown(ctx); err != nil {
        errs = append(errs, fmt.Errorf("proxy shutdown: %w", err))
    }
    
    if err := s.adminServer.Shutdown(ctx); err != nil {
        errs = append(errs, fmt.Errorf("admin shutdown: %w", err))
    }
    
    if err := s.metricsServer.Shutdown(ctx); err != nil {
        errs = append(errs, fmt.Errorf("metrics shutdown: %w", err))
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("shutdown errors: %v", errs)
    }
    
    return nil
}

func (s *Sidecar) healthHandler(w http.ResponseWriter, r *http.Request) {
    // Check application health
    resp, err := http.Get(fmt.Sprintf("http://%s:%d/health", s.config.AppHost, s.config.AppPort))
    if err != nil || resp.StatusCode != http.StatusOK {
        http.Error(w, "unhealthy", http.StatusServiceUnavailable)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"healthy"}`))
}

// Metrics collection
type Metrics struct {
    requestDuration   *prometheus.HistogramVec
    requestsTotal     *prometheus.CounterVec
    circuitBreakerOpens prometheus.Counter
    retrySuccess      *prometheus.CounterVec
    retryExhausted    prometheus.Counter
}

func NewMetrics() *Metrics {
    m := &Metrics{
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "sidecar_request_duration_seconds",
                Help: "Request duration",
                Buckets: prometheus.DefBuckets,
            },
            []string{"method", "status"},
        ),
        requestsTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "sidecar_requests_total",
                Help: "Total requests",
            },
            []string{"method", "status"},
        ),
        circuitBreakerOpens: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "sidecar_circuit_breaker_opens_total",
                Help: "Circuit breaker opens",
            },
        ),
        retrySuccess: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "sidecar_retry_success_total",
                Help: "Successful retries",
            },
            []string{"attempt"},
        ),
        retryExhausted: prometheus.NewCounter(
            prometheus.CounterOpts{
                Name: "sidecar_retry_exhausted_total",
                Help: "Retry attempts exhausted",
            },
        ),
    }
    
    prometheus.MustRegister(
        m.requestDuration,
        m.requestsTotal,
        m.circuitBreakerOpens,
        m.retrySuccess,
        m.retryExhausted,
    )
    
    return m
}
```

### Sidecar Configuration

```yaml
# config/sidecar.yaml
sidecar:
  # Application connection
  app:
    host: "127.0.0.1"
    port: 8080
  
  # Listener configuration
  proxy:
    port: 8000
    read_timeout: 30s
    write_timeout: 30s
    idle_timeout: 120s
  
  # Admin endpoints
  admin:
    port: 8001
  
  # Prometheus metrics
  metrics:
    port: 8002
    path: "/metrics"
  
  # TLS configuration (mTLS)
  tls:
    enabled: true
    cert_file: "/etc/certs/sidecar.crt"
    key_file: "/etc/certs/sidecar.key"
    ca_file: "/etc/certs/ca.crt"
    client_auth: "require_and_verify"
  
  # Circuit breaker
  circuit_breaker:
    enabled: true
    failure_threshold: 5
    success_threshold: 2
    timeout: 30s
    half_open_max_calls: 3
  
  # Retry policy
  retry:
    enabled: true
    max_attempts: 3
    base_delay: 100ms
    max_delay: 2s
    retryable_status_codes: [502, 503, 504]
  
  # Observability
  tracing:
    enabled: true
    endpoint: "http://jaeger:14268/api/traces"
    sample_rate: 0.1
  
  logging:
    level: "info"
    format: "json"
```

## Trade-off Analysis

### Sidecar vs In-Process

| Aspect | Sidecar | In-Process Library | Notes |
|--------|---------|-------------------|-------|
| **Resource Usage** | Higher (extra container) | Lower | Sidecar needs CPU/memory |
| **Latency** | Higher (~1-3ms) | Lower | Network hop vs function call |
| **Language Support** | Universal | Language-specific | Sidecar works with any language |
| **Upgrade Independence** | Yes | No | Sidecar upgrades without app changes |
| **Configuration Consistency** | High | Low | Centralized sidecar config |
| **Debugging Complexity** | Higher | Lower | Distributed across containers |
| **Security Isolation** | Strong | Weak | Sidecar as security boundary |

### Sidecar vs Shared Node Agent

```
Sidecar Pattern:
┌─────────────────────────────────────────────────────────────────────────┐
│  Pros:                                                                  │
│  • Resource isolation (CPU/memory limits per service)                   │
│  • Independent lifecycle (upgrade one without affecting others)         │
│  • Blast radius containment (failure isolated to one pod)               │
│                                                                         │
│  Cons:                                                                  │
│  • Higher resource overhead (per-pod sidecar)                           │
│  • More connections (N sidecars vs 1 agent)                             │
│  • Configuration duplication                                            │
└─────────────────────────────────────────────────────────────────────────┘

Shared Node Agent:
┌─────────────────────────────────────────────────────────────────────────┐
│  Pros:                                                                  │
│  • Lower resource overhead (single instance per node)                   │
│  • Shared connection pooling                                            │
│  • Centralized configuration                                            │
│                                                                         │
│  Cons:                                                                  │
│  • No resource isolation between services                               │
│  • Upgrade affects all services on node                                 │
│  • Single point of failure (if agent fails, all services affected)      │
│  • Security concerns (shared process space)                             │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### Sidecar Testing

```go
// test/sidecar/integration_test.go
package sidecar

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSidecarProxy(t *testing.T) {
    // Start mock application
    app := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-App-Header", "value")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    }))
    defer app.Close()
    
    // Create sidecar
    cfg := &SidecarConfig{
        AppHost: "localhost",
        AppPort: parsePort(app.URL),
        ListenPort: 18000,
        AdminPort: 18001,
        MetricsPort: 18002,
        CircuitBreaker: CircuitBreakerConfig{Enabled: false},
        Retry: RetryConfig{Enabled: false},
    }
    
    sidecar, err := NewSidecar(cfg)
    require.NoError(t, err)
    
    go sidecar.Run()
    defer sidecar.Shutdown()
    
    // Wait for startup
    time.Sleep(100 * time.Millisecond)
    
    // Test proxy functionality
    t.Run("proxies request to app", func(t *testing.T) {
        resp, err := http.Get("http://localhost:18000/test")
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        assert.Equal(t, "value", resp.Header.Get("X-App-Header"))
    })
    
    t.Run("adds tracing headers", func(t *testing.T) {
        req, _ := http.NewRequest("GET", "http://localhost:18000/test", nil)
        resp, err := http.DefaultClient.Do(req)
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.NotEmpty(t, resp.Header.Get("X-Trace-ID"))
    })
}

func TestSidecarCircuitBreaker(t *testing.T) {
    // Start failing application
    failCount := 0
    app := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        failCount++
        if failCount < 10 {
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }
        w.WriteHeader(http.StatusOK)
    }))
    defer app.Close()
    
    cfg := &SidecarConfig{
        AppHost: "localhost",
        AppPort: parsePort(app.URL),
        ListenPort: 18010,
        CircuitBreaker: CircuitBreakerConfig{
            Enabled:          true,
            FailureThreshold: 3,
            SuccessThreshold: 2,
            Timeout:          1 * time.Second,
        },
        Retry: RetryConfig{Enabled: false},
    }
    
    sidecar, err := NewSidecar(cfg)
    require.NoError(t, err)
    
    go sidecar.Run()
    defer sidecar.Shutdown()
    
    time.Sleep(100 * time.Millisecond)
    
    // Trigger failures to open circuit
    for i := 0; i < 5; i++ {
        http.Get("http://localhost:18010/test")
    }
    
    // Circuit should be open now
    resp, err := http.Get("http://localhost:18010/test")
    require.NoError(t, err)
    defer resp.Body.Close()
    
    assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func TestSidecarRetry(t *testing.T) {
    attemptCount := 0
    app := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        attemptCount++
        if attemptCount < 3 {
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }
        w.WriteHeader(http.StatusOK)
    }))
    defer app.Close()
    
    cfg := &SidecarConfig{
        AppHost: "localhost",
        AppPort: parsePort(app.URL),
        ListenPort: 18020,
        CircuitBreaker: CircuitBreakerConfig{Enabled: false},
        Retry: RetryConfig{
            Enabled:     true,
            MaxAttempts: 3,
            BaseDelay:   10 * time.Millisecond,
        },
    }
    
    sidecar, err := NewSidecar(cfg)
    require.NoError(t, err)
    
    go sidecar.Run()
    defer sidecar.Shutdown()
    
    time.Sleep(100 * time.Millisecond)
    
    resp, err := http.Get("http://localhost:18020/test")
    require.NoError(t, err)
    defer resp.Body.Close()
    
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    assert.Equal(t, 3, attemptCount) // Retried until success
}
```

## Summary

The Sidecar Pattern provides:

1. **Language Agnostic**: Works with any application language
2. **Independent Lifecycle**: Upgrade infrastructure without touching apps
3. **Consistent Operations**: Same observability/security everywhere
4. **Resource Isolation**: Sidecar failures don't crash the app
5. **Simplified Applications**: Developers focus on business logic

Key considerations:
- Resource overhead (memory/CPU per sidecar)
- Latency impact (minimal with localhost)
- Debugging complexity (distributed logs)
- Configuration management at scale
