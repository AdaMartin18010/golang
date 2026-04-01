module example.com/golang-examples

go 1.26

// Go 1.25.3 推荐的模块配置
// 使用 workspace 模式时，本模块可以直接引用其他本地模块

require (
	entgo.io/ent v0.14.6
	github.com/go-chi/chi/v5 v5.2.5
	github.com/nats-io/nats.go v1.49.0
	go.opentelemetry.io/otel v1.42.0
	google.golang.org/grpc v1.79.3
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.42.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.42.0 // indirect
	go.opentelemetry.io/otel/trace v1.42.0 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)


replace github.com/yourusername/golang => ../
replace github.com/yourusername/golang/pkg/observability => ../pkg/observability
replace github.com/yourusername/golang/pkg/concurrency => ../pkg/concurrency
replace github.com/yourusername/golang/pkg/http3 => ../pkg/http3
replace github.com/yourusername/golang/pkg/memory => ../pkg/memory
