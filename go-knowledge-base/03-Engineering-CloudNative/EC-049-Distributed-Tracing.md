# EC-049: Distributed Tracing Pattern

> **Dimension**: Engineering-CloudNative
> **Level**: S (>15KB)
> **Tags**: #distributed-tracing #observability #opentelemetry #jaeger #zipkin #microservices
> **Authoritative Sources**:
> - [Dapper, a Large-Scale Distributed Systems Tracing Infrastructure](https://research.google/pubs/dapper-a-large-scale-distributed-systems-tracing-infrastructure/) - Sigelman et al. (2010)
> - [OpenTelemetry Specification](https://opentelemetry.io/docs/specs/otel/) - CNCF (2024)
> - [W3C Trace Context](https://www.w3.org/TR/trace-context/) - W3C (2024)
> - [Mastering Distributed Tracing](https://www.packtpub.com/product/mastering-distributed-tracing/9781788628464) - Juraci (2018)

---

## 1. Problem Formalization

### 1.1 System Context and Constraints

**Definition 1.1 (Distributed Request Topology)**
Let $\mathcal{R}$ be a request traversing services $\mathcal{S} = \{S_1, S_2, ..., S_n\}$ where:
- Each service $S_i$ executes on node $N_i$ with local clock $C_i$
- Services communicate via asynchronous message passing $M_{i,j}$
- Execution produces log entries $L_i = \{l_{i,1}, l_{i,2}, ...\}$ per service

**System Constraints:**

| Constraint | Formal Definition | Impact |
|------------|-------------------|--------|
| **Clock Skew** | $\forall i, j: |C_i(t) - C_j(t)| > \epsilon$ | Cannot rely on absolute timestamps |
| **Asynchronous Execution** | $\forall S_i, S_j: async(S_i) \land async(S_j)$ | Call stacks don't span services |
| **Log Distribution** | $\forall L_i: storage(L_i) \neq storage(L_j)$ | Correlation requires external mechanism |
| **Sampling Necessity** | $\forall R: P(trace(R)) \ll 1$ | Cannot trace all requests at scale |
| **Cardinality Explosion** | $|\{trace\_ids\}| \to \infty$ | Storage and query optimization required |

### 1.2 Problem Statement

**Problem 1.1 (Causal Relationship Reconstruction)**
Given a set of distributed log entries $\mathcal{L} = \bigcup_{i=1}^{n} L_i$ produced by request $R$, reconstruct the causal execution graph $G_R = (V, E)$ where:
- $V = \{operations\_of(R)\}$: All operations comprising $R$
- $E = \{(u, v) \mid u \prec v\}$: Happens-before relationships

**Key Challenges:**

1. **Context Propagation**: Maintain correlation across process boundaries
2. **Clock Synchronization**: Establish temporal ordering despite clock skew
3. **Sampling Strategy**: Balance observability with overhead
4. **Storage Efficiency**: Handle high-cardinality trace data
5. **Query Performance**: Enable sub-second trace retrieval

### 1.3 Formal Requirements Specification

**Requirement 1.1 (Trace Completeness)**
$$\forall R: traced(R) \Rightarrow \forall op \in R: op \in trace(R)$$

**Requirement 1.2 (Causal Accuracy)**
$$\forall u, v \in trace(R): u \prec v \Leftrightarrow (u, v) \in edges(trace(R))$$

**Requirement 1.3 (Overhead Bound)**
$$\forall S_i: overhead(tracing) < \theta_{max}\% \cdot throughput(S_i)$$

---

## 2. Solution Architecture

### 2.1 Formal Tracing Model

**Definition 2.1 (Distributed Trace)**
A Distributed Trace $T$ is a 5-tuple $\langle ID, Spans, Context, References, Attributes \rangle$:

- $ID$: Globally unique 16-byte identifier
- $Spans = \{span_1, span_2, ..., span_m\}$: Timed operations
- $Context = \{trace\_id, span\_id, flags, state\}$: Propagation context
- $References \subseteq Spans \times Spans$: Parent-child and follow-from relationships
- $Attributes$: Key-value metadata annotations

**Definition 2.2 (Span)**
A Span $s$ is $\langle span\_id, operation\_name, start\_time, duration, parent\_id, tags, logs \rangle$

**Span Relationships:**
- **ChildOf**: Synchronous call ($caller \to callee$)
- **FollowsFrom**: Asynchronous call ($producer \to consumer$)

### 2.2 Sampling Strategies

| Strategy | Formula | Use Case |
|----------|---------|----------|
| **Head-based** | $P(sample) = r$ at trace start | Consistent sampling decisions |
| **Tail-based** | $P(sample \mid attributes)$ at end | Error-focused retention |
| **Probabilistic** | $hash(trace\_id) < threshold$ | Uniform distribution |
| **Adaptive** | $r = f(current\_load)$ | Load-aware sampling |

---

## 3. Visual Representations

### 3.1 Distributed Tracing Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DISTRIBUTED TRACING SYSTEM                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        APPLICATION LAYER                             │   │
│  │                                                                      │   │
│  │  Service A          Service B          Service C          Service D  │   │
│  │  ┌─────────┐       ┌─────────┐       ┌─────────┐       ┌─────────┐  │   │
│  │  │ Handler │──────►│ Handler │──────►│ Handler │──────►│ Handler │  │   │
│  │  └────┬────┘       └────┬────┘       └────┬────┘       └────┬────┘  │   │
│  │       │                 │                 │                 │       │   │
│  │  ┌────┴────┐       ┌────┴────┐       ┌────┴────┐       ┌────┴────┐  │   │
│  │  │  OTel   │       │  OTel   │       │  OTel   │       │  OTel   │  │   │
│  │  │  SDK    │◄─────►│  SDK    │◄─────►│  SDK    │◄─────►│  SDK    │  │   │
│  │  │         │       │         │       │         │       │         │  │   │
│  │  │• Auto   │       │• Auto   │       │• Auto   │       │• Auto   │  │   │
│  │  │  instr. │       │  instr. │       │  instr. │       │  instr. │  │   │
│  │  │• Context│       │• Context│       │• Context│       │• Context│  │   │
│  │  │  propag.│       │  propag.│       │  propag.│       │  propag.│  │   │
│  │  │• Span   │       │• Span   │       │• Span   │       │• Span   │  │   │
│  │  │  creation│      │  creation│      │  creation│      │  creation│  │   │
│  │  └────┬────┘       └────┬────┘       └────┬────┘       └────┬────┘  │   │
│  │       │                 │                 │                 │       │   │
│  └───────┼─────────────────┼─────────────────┼─────────────────┼───────┘   │
│          │                 │                 │                 │            │
│          └─────────────────┴─────────────────┴─────────────────┘            │
│                                   │                                        │
│                                   ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      COLLECTION LAYER                                │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                    OpenTelemetry Collector                   │    │   │
│  │  │                                                               │    │   │
│  │  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │    │   │
│  │  │  │   Receivers  │  │  Processors  │  │  Exporters   │        │    │   │
│  │  │  │              │  │              │  │              │        │    │   │
│  │  │  │• OTLP/gRPC   │  │• Batch       │  │• OTLP        │        │    │   │
│  │  │  │• OTLP/HTTP   │──►│• Memory Limit│──►│• Jaeger      │        │    │   │
│  │  │  │• Zipkin      │  │• Attributes  │  │• Zipkin      │        │    │   │
│  │  │  │• Prometheus  │  │• Resource    │  │• Prometheus  │        │    │   │
│  │  │  └──────────────┘  └──────────────┘  └──────────────┘        │    │   │
│  │  │                                                               │    │   │
│  │  │  Sampling: [Head-based ▼] [Tail-based ▼] [Adaptive ▼]        │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                   │                                        │
│                                   ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        STORAGE LAYER                                 │   │
│  │                                                                      │   │
│  │  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐            │   │
│  │  │    Jaeger     │  │    Tempo      │  │   ClickHouse  │            │   │
│  │  │  (Badger/     │  │   (Object     │  │  (Columnar    │            │   │
│  │  │   Cassandra)  │  │    Storage)   │  │    Storage)   │            │   │
│  │  └───────────────┘  └───────────────┘  └───────────────┘            │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                   │                                        │
│                                   ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        VISUALIZATION                                 │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                     Jaeger UI / Grafana                      │    │   │
│  │  │                                                               │    │   │
│  │  │  [Timeline View]  [Service Graph]  [Trace Comparison]        │    │   │
│  │  │                                                               │    │   │
│  │  │  ┌─────────────────────────────────────────────────────┐     │    │   │
│  │  │  │  Trace: a1b2c3d4...                                  │     │    │   │
│  │  │  │  ┌──────────┐┌─────┐┌──────────┐┌─────┐              │     │    │   │
│  │  │  │  │Service A ││     ││Service B ││     │              │     │    │   │
│  │  │  │  │[══════]  │────►│[════]    │────►│...             │     │    │   │
│  │  │  │  └──────────┘└─────┘└──────────┘└─────┘              │     │    │   │
│  │  │  └─────────────────────────────────────────────────────┘     │    │   │
│  │  │                                                               │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Trace Context Propagation Flow

```
Request Flow with Trace Context:

Client          Service A          Service B          Service C
  │               │                  │                  │
  │  Request      │                  │                  │
  ├──────────────►│                  │                  │
  │  Headers:     │                  │                  │
  │  traceparent: │                  │                  │
  │  00-abc123-   │                  │                  │
  │  def456-01    │                  │                  │
  │               │                  │                  │
  │               │  Create Span A   │                  │
  │               │  span_id: spanA  │                  │
  │               │                  │                  │
  │               │  Outgoing Request│                  │
  │               ├─────────────────►│                  │
  │               │  traceparent:    │                  │
  │               │  00-abc123-      │                  │
  │               │  spanB-01        │                  │
  │               │                  │                  │
  │               │                  │  Create Span B   │
  │               │                  │  parent: spanA   │
  │               │                  │                  │
  │               │                  │  Outgoing Request│
  │               │                  ├─────────────────►│
  │               │                  │  traceparent:    │
  │               │                  │  00-abc123-      │
  │               │                  │  spanC-01        │
  │               │                  │                  │
  │               │                  │                  │  Create Span C
  │               │                  │                  │  parent: spanB
  │               │                  │                  │
  │               │                  │  Response        │
  │               │                  │◄─────────────────┤
  │               │  Response        │                  │
  │               │◄─────────────────┤                  │
  │  Response     │                  │                  │
  │◄──────────────┤                  │                  │
  │               │                  │                  │

Resulting Trace Structure:
┌─────────────────────────────────────────────────────────────────┐
│ Trace: abc123                                                   │
│                                                                 │
│ ┌─────────────────────────────────────────────────────────────┐ │
│ │ Span A (Service A: handle_request)                          │ │
│ │ [████████══════] 150ms                                      │ │
│ │                                                             │ │
│ │  ┌─────────────────────────────────────────────────────┐   │ │
│ │  │ Span B (Service B: process)                         │   │ │
│ │  │ [██████] 80ms                                       │   │ │
│ │  │                                                     │   │ │
│ │  │  ┌─────────────────────────┐                       │   │ │
│ │  │  │ Span C (Service C: db)  │                       │   │ │
│ │  │  │ [███] 30ms              │                       │   │ │
│ │  │  └─────────────────────────┘                       │   │ │
│ │  └─────────────────────────────────────────────────────┘   │ │
│ └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### 3.3 Sampling Decision Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      SAMPLING DECISION FLOW                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  HEAD-BASED SAMPLING (at request start)                                     │
│  ═══════════════════════════════════════                                    │
│                                                                             │
│  Request ──► [Sampler] ──Yes──► [Trace] ──► Propagate: sampled=1            │
│                │                                               │            │
│                No                                              ▼            │
│                │                                   All spans recorded       │
│                ▼                                               │            │
│           [No Trace]                                           ▼            │
│                                                Export to backend            │
│                                                                             │
│  PROS: Simple, consistent decision, low overhead                            │
│  CONS: May miss important traces that start normally but end with error     │
│                                                                             │
│  ═══════════════════════════════════════════════════════════════════════   │
│                                                                             │
│  TAIL-BASED SAMPLING (at request end)                                       │
│  ═════════════════════════════════════                                      │
│                                                                             │
│  Request ──► [Buffer Spans] ──► [Evaluate at End]                           │
│                                     │                                       │
│                     ┌───────────────┼───────────────┐                       │
│                     │               │               │                       │
│                   Error?        Latency > T?     Custom Rule?               │
│                     │               │               │                       │
│                     ▼               ▼               ▼                       │
│                   [Yes]           [Yes]           [Yes]                     │
│                     │               │               │                       │
│                     └───────────────┴───────────────┘                       │
│                                     │                                       │
│                                     ▼                                       │
│                               [Export Full Trace]                           │
│                                                                             │
│  [Discard] ◄───No──┘                                                       │
│                                                                             │
│  PROS: Never miss important traces based on outcome                         │
│  CONS: Requires buffering, higher memory overhead                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Production Go Implementation

### 4.1 OpenTelemetry SDK Configuration

```go
package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TracerConfig holds tracing configuration
type TracerConfig struct {
	ServiceName        string
	ServiceVersion     string
	Environment        string
	OTLPAddr           string
	Insecure           bool
	SampleRate         float64
	MaxQueueSize       int
	BatchTimeout       time.Duration
	ExportTimeout      time.Duration
	MaxExportBatchSize int
	Headers            map[string]string
}

// DefaultConfig returns default tracing configuration
func DefaultConfig() TracerConfig {
	return TracerConfig{
		ServiceName:        "unknown-service",
		ServiceVersion:     "1.0.0",
		Environment:        "development",
		OTLPAddr:           "localhost:4317",
		Insecure:           true,
		SampleRate:         1.0,
		MaxQueueSize:       2048,
		BatchTimeout:       100 * time.Millisecond,
		ExportTimeout:      30 * time.Second,
		MaxExportBatchSize: 512,
		Headers:            make(map[string]string),
	}
}

// TracerProvider manages the OpenTelemetry tracer provider
type TracerProvider struct {
	provider *sdktrace.TracerProvider
	cfg      TracerConfig
}

// NewTracerProvider creates and configures a new tracer provider
func NewTracerProvider(cfg TracerConfig) (*TracerProvider, error) {
	ctx := context.Background()

	// Create OTLP exporter
	exporter, err := createOTLPExporter(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
			attribute.String("host.id", getHostID()),
			attribute.String("host.name", getHostname()),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Configure sampler
	sampler := createSampler(cfg.SampleRate)

	// Configure batch span processor
	batchProcessor := sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithMaxQueueSize(cfg.MaxQueueSize),
		sdktrace.WithBatchTimeout(cfg.BatchTimeout),
		sdktrace.WithExportTimeout(cfg.ExportTimeout),
		sdktrace.WithMaxExportBatchSize(cfg.MaxExportBatchSize),
	)

	// Create tracer provider
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
		sdktrace.WithSpanProcessor(batchProcessor),
	)

	// Set as global provider
	otel.SetTracerProvider(provider)

	// Configure propagators (W3C Trace Context + Baggage)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &TracerProvider{
		provider: provider,
		cfg:      cfg,
	}, nil
}

func createOTLPExporter(ctx context.Context, cfg TracerConfig) (sdktrace.SpanExporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.OTLPAddr),
		otlptracegrpc.WithTimeout(cfg.ExportTimeout),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	} else {
		creds := credentials.NewClientTLSFromCert(nil, "")
		opts = append(opts, otlptracegrpc.WithTLSCredentials(creds))
	}

	if len(cfg.Headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(cfg.Headers))
	}

	// Add gRPC dial options for reliability
	opts = append(opts, otlptracegrpc.WithDialOption(
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`),
	))

	return otlptracegrpc.New(ctx, opts...)
}

func createSampler(rate float64) sdktrace.Sampler {
	// Parent-based sampler with trace ID ratio
	return sdktrace.ParentBased(
		sdktrace.TraceIDRatioBased(rate),
		sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
		sdktrace.WithRemoteParentNotSampled(sdktrace.NeverSample()),
		sdktrace.WithLocalParentSampled(sdktrace.AlwaysSample()),
		sdktrace.WithLocalParentNotSampled(sdktrace.TraceIDRatioBased(rate)),
	)
}

// Shutdown gracefully shuts down the tracer provider
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	return tp.provider.Shutdown(ctx)
}

// Tracer returns a tracer for the given instrumentation scope
func (tp *TracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return tp.provider.Tracer(name, opts...)
}

// Helper functions
func getHostID() string {
	// Implementation to get unique host identifier
	return ""
}

func getHostname() string {
	// Implementation to get hostname
	return ""
}
```

### 4.2 HTTP Middleware and gRPC Interceptors

```go
package tracing

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	grpctrace "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// HTTPMiddleware creates tracing middleware for HTTP handlers
func HTTPMiddleware(operation string, opts ...otelhttp.Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, operation,
			append([]otelhttp.Option{
				otelhttp.WithPublicEndpoint(),
				otelhttp.WithSpanOptions(
					trace.WithAttributes(
						attribute.String("http.server_name", operation),
					),
				),
			}, opts...)...,
		)
	}
}

// TracedHTTPClient returns an HTTP client with tracing enabled
func TracedHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttp.DefaultClientTrace(ctx)
			}),
		),
		Timeout: timeout,
	}
}

// UnaryClientInterceptor creates a gRPC unary client interceptor with tracing
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return grpctrace.UnaryClientInterceptor(
		grpctrace.WithPropagators(propagation.TraceContext{}),
	)
}

// UnaryServerInterceptor creates a gRPC unary server interceptor with tracing
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpctrace.UnaryServerInterceptor(
		grpctrace.WithPropagators(propagation.TraceContext{}),
	)
}

// StreamClientInterceptor creates a gRPC stream client interceptor with tracing
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return grpctrace.StreamClientInterceptor(
		grpctrace.WithPropagators(propagation.TraceContext{}),
	)
}

// StreamServerInterceptor creates a gRPC stream server interceptor with tracing
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpctrace.StreamServerInterceptor(
		grpctrace.WithPropagators(propagation.TraceContext{}),
	)
}

// ExtractContext extracts trace context from gRPC metadata
func ExtractContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	return propagation.TraceContext{}.Extract(ctx, GRPCHeaderCarrier(md))
}

// GRPCHeaderCarrier adapts gRPC metadata for propagation.TextMapCarrier
type GRPCHeaderCarrier metadata.MD

func (m GRPCHeaderCarrier) Get(key string) string {
	values := metadata.MD(m).Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (m GRPCHeaderCarrier) Set(key string, value string) {
	metadata.MD(m).Set(key, value)
}

func (m GRPCHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range metadata.MD(m) {
		keys = append(keys, k)
	}
	return keys
}
```

### 4.3 Database and Message Queue Instrumentation

```go
package tracing

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TracedDB wraps sql.DB with tracing
type TracedDB struct {
	*sql.DB
	tracer trace.Tracer
	dbName string
}

// NewTracedDB creates a traced database connection
func NewTracedDB(db *sql.DB, tracer trace.Tracer, dbName string) *TracedDB {
	return &TracedDB{
		DB:     db,
		tracer: tracer,
		dbName: dbName,
	}
}

// QueryContext traces SQL queries
func (db *TracedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx, span := db.tracer.Start(ctx, "db.query",
		trace.WithAttributes(
			attribute.String("db.system", db.detectDBSystem()),
			attribute.String("db.name", db.dbName),
			attribute.String("db.statement", query),
			attribute.String("db.operation", "SELECT"),
		),
	)
	defer span.End()

	start := time.Now()
	rows, err := db.DB.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	span.SetAttributes(attribute.Int64("db.duration_ms", duration.Milliseconds()))

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("error", true))
	}

	return rows, err
}

// ExecContext traces SQL executions
func (db *TracedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	op := detectOperation(query)
	ctx, span := db.tracer.Start(ctx, fmt.Sprintf("db.%s", op),
		trace.WithAttributes(
			attribute.String("db.system", db.detectDBSystem()),
			attribute.String("db.name", db.dbName),
			attribute.String("db.statement", query),
			attribute.String("db.operation", op),
		),
	)
	defer span.End()

	start := time.Now()
	result, err := db.DB.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	span.SetAttributes(attribute.Int64("db.duration_ms", duration.Milliseconds()))

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("error", true))
	} else {
		if rowsAffected, err := result.RowsAffected(); err == nil {
			span.SetAttributes(attribute.Int64("db.rows_affected", rowsAffected))
		}
	}

	return result, err
}

func (db *TracedDB) detectDBSystem() string {
	// Detect database type from driver
	return "postgresql"
}

func detectOperation(query string) string {
	// Simple detection based on query prefix
	if len(query) < 6 {
		return "UNKNOWN"
	}
	switch query[:6] {
	case "SELECT":
		return "SELECT"
	case "INSERT":
		return "INSERT"
	case "UPDATE":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

// TracedRedis wraps redis.Client with tracing
type TracedRedis struct {
	client *redis.Client
	tracer trace.Tracer
}

// NewTracedRedis creates a traced Redis client
func NewTracedRedis(client *redis.Client, tracer trace.Tracer) *TracedRedis {
	return &TracedRedis{
		client: client,
		tracer: tracer,
	}
}

// Get traces Redis GET operation
func (r *TracedRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	ctx, span := r.tracer.Start(ctx, "redis.get",
		trace.WithAttributes(
			attribute.String("db.system", "redis"),
			attribute.String("db.operation", "GET"),
			attribute.String("db.redis.key", key),
		),
	)
	defer span.End()

	cmd := r.client.Get(ctx, key)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("error", true))
	} else {
		span.SetAttributes(attribute.Bool("cache.hit", err != redis.Nil))
	}

	return cmd
}

// Set traces Redis SET operation
func (r *TracedRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	ctx, span := r.tracer.Start(ctx, "redis.set",
		trace.WithAttributes(
			attribute.String("db.system", "redis"),
			attribute.String("db.operation", "SET"),
			attribute.String("db.redis.key", key),
			attribute.String("db.redis.ttl", expiration.String()),
		),
	)
	defer span.End()

	cmd := r.client.Set(ctx, key, value, expiration)
	if err := cmd.Err(); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("error", true))
	}

	return cmd
}
```

---

## 5. Failure Scenarios and Mitigations

### 5.1 Tracing Failure Taxonomy

| Scenario | Impact | Detection | Mitigation |
|----------|--------|-----------|------------|
| **Exporter Timeout** | Trace data loss | Export latency | Queue + Retry + Drop policy |
| **Collector Overload** | Sampling rate degradation | Queue depth | Backpressure + Adaptive sampling |
| **Context Loss** | Broken traces | Orphan spans | Context validation + Logging |
| **High Cardinality** | Storage exhaustion | Cardinality metrics | Attribute limits + Cardinality capping |
| **Clock Skew** | Incorrect span ordering | Span timestamp validation | Clock sync + Logical timestamps |
| **Trace ID Collision** | Merged traces | ID uniqueness check | Sufficient ID entropy |

### 5.2 Resilience Patterns

```go
// ResilientExporter wraps an exporter with retry logic
type ResilientExporter struct {
	exporter sdktrace.SpanExporter
	retries  int
	backoff  time.Duration
}

func (re *ResilientExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	var lastErr error
	for i := 0; i <= re.retries; i++ {
		if i > 0 {
			time.Sleep(re.backoff * time.Duration(i))
		}
		err := re.exporter.ExportSpans(ctx, spans)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("export failed after %d retries: %w", re.retries, lastErr)
}
```

---

## 6. Semantic Trade-off Analysis

### 6.1 Sampling Strategy Comparison

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SAMPLING STRATEGY COMPARISON                              │
├─────────────────────┬───────────────┬───────────────┬───────────────────────┤
│     Dimension       │  Head-Based   │  Tail-Based   │  Adaptive             │
├─────────────────────┼───────────────┼───────────────┼───────────────────────┤
│ Decision Point      │ Request start │ Request end   │ Continuous            │
│ Memory Usage        │ Low           │ High (buffer) │ Medium                │
│ CPU Overhead        │ Minimal       │ Moderate      │ Higher                │
│ Important Traces    │ May miss      │ Never miss    │ Configurable          │
│ Implementation      │ Simple        │ Complex       │ Most complex          │
│ Use Case            │ General       │ Error-focused │ Cost-sensitive        │
└─────────────────────┴───────────────┴───────────────┴───────────────────────┘
```

### 6.2 Backend Storage Trade-offs

| Backend | Retention | Query Speed | Cost | Scale |
|---------|-----------|-------------|------|-------|
| **Jaeger/Cassandra** | Medium | Fast | Medium | Medium |
| **Tempo/S3** | Long | Medium | Low | High |
| **ClickHouse** | Long | Very Fast | Medium | High |
| **Elasticsearch** | Configurable | Medium | High | High |

---

## 7. References

1. Sigelman, B. H., et al. (2010). Dapper, a Large-Scale Distributed Systems Tracing Infrastructure. *Google Research*.
2. OpenTelemetry Project. (2024). *OpenTelemetry Specification*. opentelemetry.io.
3. W3C. (2024). *Trace Context*. W3C Recommendation.
4. Juraci, Y. (2018). *Mastering Distributed Tracing*. Packt Publishing.
5. Bryan, L. (2022). *Cloud-Native Observability with OpenTelemetry*. O'Reilly Media.
