module example.com/golang-examples

go 1.27

// Go 1.26.2 推荐的模块配置
// 使用 workspace 模式时，本模块可以直接引用其他本地模块

require (
	entgo.io/ent v0.14.6
	github.com/go-chi/chi/v5 v5.2.5
	github.com/google/wire v0.7.0
	github.com/nats-io/nats.go v1.49.0
	github.com/yourusername/golang v0.0.0-00010101000000-000000000000
	github.com/yourusername/golang/pkg/observability v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.42.0
	google.golang.org/grpc v1.79.3
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cilium/ebpf v0.20.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/lib/pq v1.12.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.28 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.42.0 // indirect
	go.opentelemetry.io/otel/sdk v1.42.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.42.0 // indirect
	go.opentelemetry.io/otel/trace v1.42.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.1 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/yourusername/golang => ../

replace github.com/yourusername/golang/pkg/observability => ../pkg/observability

replace github.com/yourusername/golang/pkg/concurrency => ../pkg/concurrency

replace github.com/yourusername/golang/pkg/http3 => ../pkg/http3

replace github.com/yourusername/golang/pkg/memory => ../pkg/memory
