package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func initTracer() (func(context.Context) error, error) {
	res, _ := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("example-observability-app"),
		semconv.DeploymentEnvironment("dev"),
	))
	exp, err := otlptracegrpc.New(context.Background())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

func main() {
	shutdown, err := initTracer()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = shutdown(context.Background())
	}()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("example").Start(r.Context(), "hello")
		defer span.End()
		time.Sleep(50 * time.Millisecond)
		traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		logger.Info("hello called", "trace_id", traceID)
		_, _ = fmt.Fprintln(w, "hello")
	})

	addr := ":8080"
	logger.Info("starting server", "addr", addr)
	_ = http.ListenAndServe(addr, nil)
}
