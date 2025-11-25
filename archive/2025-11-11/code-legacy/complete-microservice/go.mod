module github.com/yourusername/golang/examples/complete-microservice

go 1.25.3

require (
	github.com/yourusername/golang/pkg/agent v0.0.0
	github.com/yourusername/golang/pkg/concurrency v0.0.0
	github.com/yourusername/golang/pkg/http3 v0.0.0
	github.com/yourusername/golang/pkg/memory v0.0.0
	github.com/yourusername/golang/pkg/observability v0.0.0
	github.com/gorilla/websocket v1.5.3
)

replace (
	github.com/yourusername/golang/pkg/agent => ../../pkg/agent
	github.com/yourusername/golang/pkg/concurrency => ../../pkg/concurrency
	github.com/yourusername/golang/pkg/http3 => ../../pkg/http3
	github.com/yourusername/golang/pkg/memory => ../../pkg/memory
	github.com/yourusername/golang/pkg/observability => ../../pkg/observability
)

