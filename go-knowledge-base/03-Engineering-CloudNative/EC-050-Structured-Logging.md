# EC-050: Structured Logging Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #structured-logging #observability #json-logging #correlation-id #centralized-logging
> **Authoritative Sources**:
> - [The Art of Logging](https://www.oreilly.com/library/view/the-art-of/9781492081626/) - Phil Zona (2020)
> - [12 Factor App - Logs](https://12factor.net/logs) - Adam Wiggins (2011)
> - [Splunk Logging Best Practices](https://docs.splunk.com/Documentation/Splunk/latest/Data/BestPractices) - Splunk (2024)
> - [Fluentd Architecture](https://docs.fluentd.org/quickstart/life-of-a-fluentd-event) - Fluentd (2024)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Log Data Domain)**
Let $\mathcal{L}$ be the log output of application $A$ executing on distributed nodes $\mathcal{N} = \{N_1, N_2, ..., N_n\}$ where:
- Each log entry $l \in \mathcal{L}$ has timestamp $t(l)$, level $v(l)$, message $m(l)$
- Traditional format: $l_{text} = "[timestamp] [level] message text"
- Structured format: $l_{struct} = \{k_1: v_1, k_2: v_2, ..., k_m: v_m\}$

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Volume Scale** | $|\mathcal{L}| \to 10^{9+} entries/day$ | Parsing overhead must be minimal |
| **Query Velocity** | $T_{query}(\mathcal{L}, pattern) < 1s$ | Index-friendly format required |
| **Correlation Need** | $\forall l_i, l_j: correlate(l_i, l_j)$ | Standardized correlation fields |
| **Schema Evolution** | $\exists t: schema(\mathcal{L}_t) \neq schema(\mathcal{L}_{t+1})$ | Backward compatibility |
| **Multi-source Aggregation** | $\mathcal{L} = \bigcup_{i=1}^{n} \mathcal{L}_i$ | Common format across sources |

### 1.2 Problem Statement

**Problem 1.1 (Log Queryability Problem)**
Given a query predicate $\phi$ over log corpus $\mathcal{L}$, find all entries satisfying $\phi$:

$$\mathcal{L}_\phi = \{l \in \mathcal{L} \mid \phi(l) = true\}$$

With time complexity $O(\log |\mathcal{L}|)$ or better.

**Key Challenges:**

1. **Unstructured Parsing**: Text logs require regex parsing ($O(n)$ per entry)
2. **Schema Variance**: Different log formats prevent unified querying
3. **Context Loss**: Flat text loses hierarchical context
4. **Cardinality Issues**: High-cardinality fields affect storage
5. **Timestamp Ambiguity**: Timezone and format inconsistencies

### 1.3 Formal Requirements Specification

**Requirement 1.1 (Structured Format)**
$$\forall l \in \mathcal{L}: l \in \mathbb{JSON} \lor l \in \mathbb{KV}$$

**Requirement 1.2 (Correlation Support)**
$$\forall l: trace\_id(l) \neq \emptyset \land span\_id(l) \neq \emptyset$$

**Requirement 1.3 (Schema Stability)**
$$\forall k \in keys(l_t): k \in keys(l_{t+1}) \lor deprecated(k)$$

---

## 2. Solution Architecture

### 2.1 Structured Log Model

**Definition 2.1 (Structured Log Entry)**
A structured log entry $l$ is a tuple $\langle Timestamp, Level, Message, Context, Metadata \rangle$:

- $Timestamp$: ISO 8601 with millisecond precision
- $Level \in \{DEBUG, INFO, WARN, ERROR, FATAL\}$
- $Message$: Human-readable description
- $Context = \{trace\_id, span\_id, service, operation\}$: Distributed tracing context
- $Metadata$: Application-specific key-value pairs

**Log Level Semantics:**

| Level | Usage | Production Volume |
|-------|-------|-------------------|
| DEBUG | Detailed diagnostics | Disabled |
| INFO | Normal operations | High |
| WARN | Anomalies, recoverable | Medium |
| ERROR | Failures, handled | Low |
| FATAL | System failure, exit | Minimal |

---

## 3. Visual Representations

### 3.1 Structured Logging Pipeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    STRUCTURED LOGGING PIPELINE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  APPLICATION LAYER                                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                     │   │
│  │  ┌─────────────────────────────────────────────────────────────┐   │   │
│  │  │                      Logger (Application)                    │   │   │
│  │  │                                                               │   │   │
│  │  │   logger.Info("Order processed",                             │   │   │
│  │  │     zap.String("order_id", "123"),                           │   │   │
│  │  │     zap.String("customer", "456"),                           │   │   │
│  │  │     zap.Float64("amount", 99.99),                            │   │   │
│  │  │     zap.String("trace_id", traceID),                         │   │   │
│  │  │   )                                                          │   │   │
│  │  │                                                               │   │   │
│  │  │   Output:                                                     │   │   │
│  │  │   {                                                           │   │   │
│  │  │     "timestamp": "2024-01-15T09:30:00.123Z",                 │   │   │
│  │  │     "level": "INFO",                                          │   │   │
│  │  │     "message": "Order processed",                             │   │   │
│  │  │     "service": "order-service",                               │   │   │
│  │  │     "trace_id": "abc123def456",                               │   │   │
│  │  │     "span_id": "span789",                                     │   │   │
│  │  │     "order_id": "123",                                        │   │   │
│  │  │     "customer": "456",                                        │   │   │
│  │  │     "amount": 99.99                                           │   │   │
│  │  │   }                                                           │   │   │
│  │  └─────────────────────────────────────────────────────────────┘   │   │
│  │                              │                                      │   │
│  └──────────────────────────────┼──────────────────────────────────────┘   │
│                                 │                                           │
│                                 ▼                                           │
│  AGGREGATION LAYER                                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      Log Shipper (Fluentd/Fluent Bit)                │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│  │  │   Input     │  │   Parser    │  │   Filter    │  │   Output    │ │   │
│  │  │             │──►│             │──►│             │──►│             │ │   │
│  │  │• Tail files │  │• JSON parse │  │• Add host   │  │• Forward    │ │   │
│  │  │• Journald   │  │• Grok       │  │• Add k8s    │  │• Elasticsearch│  │   │
│  │  │• TCP/UDP    │  │• Regex      │  │  metadata   │  │• Kafka      │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                 │                                           │
│                                 ▼                                           │
│  STORAGE & ANALYSIS                                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────────────┐    │   │
│  │  │Elasticsearch/ │  │    Kafka      │  │    Object Storage     │    │   │
│  │  │   OpenSearch  │  │   (Buffer)    │  │    (S3/GCS/Azure)     │    │   │
│  │  │               │  │               │  │                       │    │   │
│  │  │• Full-text    │  │• Stream       │  │• Long-term retention  │    │   │
│  │  │  search       │  │  processing   │  │• Compliance archive   │    │   │
│  │  │• Aggregations │  │• Replay       │  │• Cost-optimized       │    │   │
│  │  └───────┬───────┘  └───────────────┘  └───────────────────────┘    │   │
│  │          │                                                          │   │
│  │          ▼                                                          │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                     Visualization (Kibana/Grafana)           │    │   │
│  │  │                                                               │    │   │
│  │  │  [Dashboards]  [Alerts]  [Search]  [Anomaly Detection]       │    │   │
│  │  │                                                               │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Correlation and Context Propagation

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    LOG CORRELATION ACROSS SERVICES                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Request Flow:                                                              │
│                                                                             │
│  Client ──────► API Gateway ──────► Order Service ──────► Payment Service   │
│                                                                             │
│  Correlation Context:                                                       │
│                                                                             │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐│
│  │  Request ID   │──│  Request ID   │──│  Request ID   │──│  Request ID   ││
│  │  trace_id:    │  │  trace_id:    │  │  trace_id:    │  │  trace_id:    ││
│  │  "abc123"     │  │  "abc123"     │  │  "abc123"     │  │  "abc123"     ││
│  │               │  │               │  │               │  │               ││
│  │  Span ID      │  │  Span ID      │  │  Span ID      │  │  Span ID      ││
│  │  span_id:     │  │  span_id:     │  │  span_id:     │  │  span_id:     ││
│  │  "gateway-1"  │  │  "api-1"      │  │  "order-1"    │  │  "pay-1"      ││
│  └───────────────┘  └───────────────┘  └───────────────┘  └───────────────┘│
│                                                                             │
│  Generated Logs (searchable by trace_id="abc123"):                          │
│                                                                             │
│  TIME          SERVICE        LEVEL  MESSAGE                        SPAN    │
│  ─────────────────────────────────────────────────────────────────────────  │
│  09:30:00.123  api-gateway    INFO   Request received               gateway-1│
│  09:30:00.145  api-gateway    INFO   Routing to order service       gateway-1│
│  09:30:00.156  order-service  INFO   Order validation started       order-1  │
│  09:30:00.178  order-service  INFO   Order validation completed     order-1  │
│  09:30:00.189  order-service  INFO   Calling payment service        order-1  │
│  09:30:00.201  payment-svc    INFO   Payment processing started     pay-1    │
│  09:30:00.245  payment-svc    INFO   Payment authorized              pay-1    │
│  09:30:00.267  order-service  INFO   Order confirmed                 order-1  │
│  09:30:00.289  api-gateway    INFO   Response sent                   gateway-1│
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.3 Log Level and Sampling Strategy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    LOG LEVEL HIERARCHY AND SAMPLING                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Volume:    ████████████████████▌                                           │
│             DEBUG (Disabled in prod)                                        │
│             Detailed diagnostics, development only                          │
│                                                                             │
│  Volume:    █████████████████▌                                              │
│             INFO (Sampled 100%)                                             │
│             Normal operations, business events                              │
│                                                                             │
│  Volume:    █████████▌                                                      │
│             WARN (Sampled 100%)                                             │
│             Anomalies, recoverable errors                                   │
│                                                                             │
│  Volume:    ███▌                                                            │
│             ERROR (Sampled 100%)                                            │
│             Failures, handled exceptions                                    │
│                                                                             │
│  Volume:    █▌                                                              │
│             FATAL (Sampled 100%)                                            │
│             System failure, immediate attention                             │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  Sampling Rules:                                                            │
│  • ERROR/FATAL: Always keep (0% sampling)                                   │
│  • WARN: Keep 100% for 1 hour, then sample 10%                              │
│  • INFO: Keep 100% for 24 hours, then sample 1%                             │
│  • DEBUG: Never in production                                               │
│                                                                             │
│  Retention Policy:                                                          │
│  • Hot storage (SSD): 7 days, all levels                                    │
│  • Warm storage (HDD): 30 days, WARN+                                       │
│  • Cold storage (S3): 1 year, ERROR+                                        │
│  • Archive: 7 years, FATAL only (compliance)                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 Structured Logger with zap

```go
package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional context support
type Logger struct {
	*zap.Logger
	service string
	version string
}

// Config holds logger configuration
type Config struct {
	Service     string
	Version     string
	Environment string
	Level       string
	Format      string // json or console
	Output      string // stdout, file path
}

// New creates a new structured logger
func New(cfg Config) (*Logger, error) {
	// Parse log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Configure encoder
	var encoder zapcore.Encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Configure output
	var writeSyncer zapcore.WriteSyncer
	if cfg.Output == "stdout" || cfg.Output == "" {
		writeSyncer = zapcore.AddSync(os.Stdout)
	} else {
		file, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writeSyncer = zapcore.AddSync(file)
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger with default fields
	zapLogger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("service", cfg.Service),
			zap.String("version", cfg.Version),
			zap.String("environment", cfg.Environment),
			zap.String("hostname", getHostname()),
			zap.Int("pid", os.Getpid()),
		),
	)

	return &Logger{
		Logger:  zapLogger,
		service: cfg.Service,
		version: cfg.Version,
	}, nil
}

// WithContext returns a logger with trace context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return l
	}

	spanContext := span.SpanContext()
	if !spanContext.IsValid() {
		return l
	}

	return &Logger{
		Logger: l.With(
			zap.String("trace_id", spanContext.TraceID().String()),
			zap.String("span_id", spanContext.SpanID().String()),
			zap.Bool("trace_sampled", spanContext.IsSampled()),
		),
		service: l.service,
		version: l.version,
	}
}

// WithField returns a logger with additional field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Logger: l.With(zap.Any(key, value)),
	}
}

// WithFields returns a logger with multiple fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{
		Logger: l.With(zapFields...),
	}
}

// WithError returns a logger with error field
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.With(zap.Error(err)),
	}
}

// Debug logs debug message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Info logs info message
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Warn logs warning message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Error logs error message
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Fatal logs fatal message and exits
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// LogRequest logs HTTP request details
func (l *Logger) LogRequest(method, path string, status int, duration time.Duration, fields ...zap.Field) {
	fields = append(fields,
		zap.String("http.method", method),
		zap.String("http.path", path),
		zap.Int("http.status_code", status),
		zap.Duration("http.duration", duration),
	)

	if status >= 500 {
		l.Error("HTTP request failed", fields...)
	} else if status >= 400 {
		l.Warn("HTTP request warning", fields...)
	} else {
		l.Info("HTTP request completed", fields...)
	}
}

// LogDBQuery logs database query details
func (l *Logger) LogDBQuery(operation, table string, duration time.Duration, rows int64, fields ...zap.Field) {
	fields = append(fields,
		zap.String("db.operation", operation),
		zap.String("db.table", table),
		zap.Duration("db.duration", duration),
		zap.Int64("db.rows_affected", rows),
	)
	l.Debug("Database query executed", fields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// getHostname returns the hostname
func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
```

### 4.2 HTTP Middleware for Request Logging

```go
package logging

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.written = true
	return rw.ResponseWriter.Write(b)
}

// Middleware creates HTTP logging middleware
func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create response writer wrapper
		wrapped := newResponseWriter(w)
		
		// Extract request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		// Add request ID to response header
		wrapped.Header().Set("X-Request-ID", requestID)
		
		// Create request-scoped logger
		requestLogger := l.With(
			zap.String("request_id", requestID),
			zap.String("trace_id", r.Header.Get("X-Trace-ID")),
		)
		
		// Log request start
		requestLogger.Debug("HTTP request started",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)
		
		// Process request
		next.ServeHTTP(wrapped, r)
		
		// Log request completion
		duration := time.Since(start)
		requestLogger.LogRequest(
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			zap.Int64("bytes_written", r.ContentLength),
		)
	})
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}
```

### 4.3 gRPC Interceptor for Structured Logging

```go
package logging

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor creates a gRPC unary server interceptor with logging
func (l *Logger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		
		// Extract metadata
		md, _ := metadata.FromIncomingContext(ctx)
		
		// Get peer info
		p, _ := peer.FromContext(ctx)
		
		// Create request-scoped logger
		requestLogger := l.With(
			zap.String("grpc.method", info.FullMethod),
			zap.String("grpc.peer", p.Addr.String()),
			zap.Strings("grpc.metadata.request_id", md.Get("x-request-id")),
		)
		
		// Log request
		requestLogger.Debug("gRPC request started",
			zap.Any("request", req),
		)
		
		// Handle request
		resp, err := handler(ctx, req)
		
		// Log response
		duration := time.Since(start)
		st, _ := status.FromError(err)
		
		fields := []zap.Field{
			zap.Duration("grpc.duration", duration),
			zap.String("grpc.code", st.Code().String()),
		}
		
		if err != nil {
			fields = append(fields, zap.Error(err))
			requestLogger.Error("gRPC request failed", fields...)
		} else {
			requestLogger.Info("gRPC request completed", fields...)
		}
		
		return resp, err
	}
}

// StreamServerInterceptor creates a gRPC stream server interceptor with logging
func (l *Logger) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		
		ctx := ss.Context()
		md, _ := metadata.FromIncomingContext(ctx)
		p, _ := peer.FromContext(ctx)
		
		requestLogger := l.With(
			zap.String("grpc.method", info.FullMethod),
			zap.String("grpc.peer", p.Addr.String()),
			zap.Bool("grpc.client_stream", info.IsClientStream),
			zap.Bool("grpc.server_stream", info.IsServerStream),
		)
		
		requestLogger.Info("gRPC stream started",
			zap.Strings("grpc.metadata.request_id", md.Get("x-request-id")),
		)
		
		// Wrap stream to log messages
		wrapped := &loggedServerStream{
			ServerStream: ss,
			logger:       requestLogger,
		}
		
		err := handler(srv, wrapped)
		
		duration := time.Since(start)
		st, _ := status.FromError(err)
		
		fields := []zap.Field{
			zap.Duration("grpc.duration", duration),
			zap.String("grpc.code", st.Code().String()),
			zap.Int64("grpc.messages_sent", wrapped.messagesSent),
			zap.Int64("grpc.messages_received", wrapped.messagesReceived),
		}
		
		if err != nil {
			requestLogger.Error("gRPC stream failed", append(fields, zap.Error(err))...)
		} else {
			requestLogger.Info("gRPC stream completed", fields...)
		}
		
		return err
	}
}

// loggedServerStream wraps grpc.ServerStream to count messages
type loggedServerStream struct {
	grpc.ServerStream
	logger            *Logger
	messagesSent      int64
	messagesReceived  int64
}

func (s *loggedServerStream) SendMsg(m interface{}) error {
	s.messagesSent++
	return s.ServerStream.SendMsg(m)
}

func (s *loggedServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.messagesReceived++
	}
	return err
}
```

---

## 5. Failure Scenarios and Mitigations

### 5.1 Failure Taxonomy

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Log Flooding** | Storage exhaustion | Rate limit exceeded | Circuit breaker + Sampling |
| **Sensitive Data Leak** | Security incident | Pattern matching | Redaction filters |
| **Disk Full** | Application crash | Disk usage monitoring | Log rotation + Archival |
| **Network Partition** | Log loss | Export failure | Local buffering + Retry |
| **Timestamp Drift** | Ordering errors | NTP monitoring | UTC timestamps + Monotonic clocks |
| **Hot Shard** | Query degradation | Index monitoring | Shard rebalancing |

---

## 6. Semantic Trade-off Analysis

### 6.1 Format Comparison

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    LOG FORMAT COMPARISON                                     │
├─────────────────────┬─────────────────┬─────────────────┬───────────────────┤
│     Dimension       │     Plain Text  │      JSON       │    Binary (Proto) │
├─────────────────────┼─────────────────┼─────────────────┼────────────────────┤
│ Human Readability   │ ✅ Excellent    │ ⚠️  Moderate    │ ❌ Poor           │
│ Machine Parsing     │ ❌ Regex needed │ ✅ Native       │ ✅ Native         │
│ Query Performance   │ ❌ Slow         │ ✅ Fast         │ ✅ Fast           │
│ Storage Efficiency  │ ⚠️  Moderate    │ ⚠️  Moderate    │ ✅ Compact        │
│ Schema Flexibility  │ ❌ None         │ ✅ High         │ ⚠️  Versioned     │
│ Tooling Support     │ ⚠️  Limited     │ ✅ Excellent    │ ⚠️  Specialized   │
└─────────────────────┴─────────────────┴─────────────────┴────────────────────┘

Recommendation: JSON for most applications; Protocol Buffers for high-throughput,
resource-constrained environments.
```

---

## 7. References

1. Zona, P. (2020). *The Art of Logging*. O'Reilly Media.
2. Wiggins, A. (2011). *The Twelve-Factor App*. 12factor.net.
3. Splunk. (2024). *Logging Best Practices*. docs.splunk.com.
4. Fluentd Project. (2024). *Fluentd Architecture*. fluentd.org.
5. Uber Engineering. (2017). *Logging at Uber with zap*. uber.com/blog.
