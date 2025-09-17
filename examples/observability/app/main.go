package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func initProviders() (func(context.Context) error, func(context.Context) error, error) {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "example-observability-app"
	}
	env := os.Getenv("OTEL_ENV")
	if env == "" {
		env = "dev"
	}
	res, _ := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.DeploymentEnvironment(env),
	))

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true"

	if endpoint == "" {
		// 无导出端点：静默降级为本地 Provider（不创建导出器，不打印连接错误）
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
		)
		otel.SetTracerProvider(tp)
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
		)
		otel.SetMeterProvider(mp)
		return tp.Shutdown, mp.Shutdown, nil
	}

	// 已配置端点：按配置初始化 gRPC 导出器
	traceOpts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(endpoint)}
	metricOpts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(endpoint)}
	if insecure {
		traceOpts = append(traceOpts, otlptracegrpc.WithInsecure())
		metricOpts = append(metricOpts, otlpmetricgrpc.WithInsecure())
	}

	texp, err := otlptracegrpc.New(context.Background(), traceOpts...)
	if err != nil {
		return nil, nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(texp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	mexp, err := otlpmetricgrpc.New(context.Background(), metricOpts...)
	if err != nil {
		return tp.Shutdown, nil, err
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(mexp)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	return tp.Shutdown, mp.Shutdown, nil
}

func main() {
	tShutdown, mShutdown, err := initProviders()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = tShutdown(context.Background())
		_ = mShutdown(context.Background())
	}()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	meter := otel.Meter("example-observability-app")
	reqDur, _ := meter.Float64Histogram("http_server_request_duration_seconds")
	reqTotal, _ := meter.Int64Counter("http_server_requests_total")

	// HTTP 路由
	mux := http.NewServeMux()

	// 健康检查：同时支持 /healthz 与 /health，返回 JSON
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"status":"ok"}`)
	}
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("example").Start(r.Context(), "hello")
		defer span.End()
		time.Sleep(50 * time.Millisecond)
		traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		logger.Info("hello called", "trace_id", traceID)
		w.Header().Set("Trace-ID", traceID)
		_, _ = fmt.Fprintln(w, "hello")
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("example").Start(r.Context(), "error")
		defer span.End()
		time.Sleep(20 * time.Millisecond)
		traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		w.Header().Set("Trace-ID", traceID)
		http.Error(w, "internal error", http.StatusInternalServerError)
	})

	// 记录时延与计数（含标签）的中间件

	metricsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			duration := time.Since(start).Seconds()
			attrs := []attribute.KeyValue{
				attribute.String("method", r.Method),
				attribute.String("path", r.URL.Path),
				attribute.String("code", strconv.Itoa(rec.status)),
			}
			reqDur.Record(r.Context(), duration, metric.WithAttributes(attrs...))
			reqTotal.Add(r.Context(), 1, metric.WithAttributes(attrs...))
		})
	}

	addr := ":8080"
	logger.Info("starting server", "addr", addr)
	root := metricsMiddleware(otelhttp.NewHandler(mux, "http.server"))
	srv := &http.Server{Addr: addr, Handler: root}

	go func() {
		_ = srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
