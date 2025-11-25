module github.com/yourusername/golang

go 1.25.3

require (
	// Web 框架
	github.com/go-chi/chi/v5 v5.0.10
	github.com/labstack/echo/v4 v4.11.3

	// ORM
	entgo.io/ent v0.12.5

	// 配置
	github.com/spf13/viper v1.17.0

	// 依赖注入
	github.com/google/wire v0.5.0

	// 数据库
	github.com/jackc/pgx/v5 v5.5.0

	// gRPC
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0

	// GraphQL
	github.com/99designs/gqlgen v0.17.40

	// MQTT
	github.com/eclipse/paho.mqtt.golang v1.4.3

	// Kafka
	github.com/IBM/sarama v1.42.1

	// OpenTelemetry
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
	go.opentelemetry.io/otel/metric v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.opentelemetry.io/otel/sdk/resource v1.21.0
	go.opentelemetry.io/otel/sdk/metric v1.21.0
	go.opentelemetry.io/otel/sdk/trace v1.21.0
	go.opentelemetry.io/otel/propagation v1.21.0
	go.opentelemetry.io/otel/semconv/v1.21.0 v1.21.0

	// eBPF
	github.com/cilium/ebpf v0.12.3

	// 测试
	github.com/stretchr/testify v1.8.4

	// UUID
	github.com/google/uuid v1.5.0
)
