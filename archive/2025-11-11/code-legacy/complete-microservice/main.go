// Package main demonstrates a complete microservice application
// integrating all core modules: Agent, Concurrency, HTTP/3, Memory, and Observability.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/golang/pkg/agent/core"
	"github.com/yourusername/golang/pkg/concurrency/patterns"
	"github.com/yourusername/golang/pkg/memory"
	"github.com/yourusername/golang/pkg/observability"
)

// Config holds application configuration
type Config struct {
	ServerAddr     string
	WorkerPoolSize int
	LogLevel       observability.LogLevel
}

// Application represents the microservice application
type Application struct {
	config  *Config
	agent   *core.BaseAgent
	logger  *observability.Logger
	memPool *memory.GenericPool[*Request]
	server  *http.Server
}

// Request represents a business request
type Request struct {
	ID        string
	Timestamp time.Time
	Data      map[string]interface{}
}

// Reset resets the request for reuse
func (r *Request) Reset() {
	r.ID = ""
	r.Timestamp = time.Time{}
	r.Data = nil
}

// NewApplication creates a new application instance
func NewApplication(config *Config) *Application {
	// Initialize observability
	logger := observability.NewLogger(config.LogLevel, os.Stdout)
	logger.AddHook(observability.NewMetricsHook())
	observability.SetDefaultLogger(logger)

	// Create AI-Agent
	agent := core.NewBaseAgent("microservice-agent")

	// Create memory pool for requests
	memPool := memory.NewGenericPool(
		func() *Request { return &Request{Data: make(map[string]interface{})} },
		func(r *Request) { r.Reset() },
		1000,
	)

	return &Application{
		config:  config,
		agent:   agent,
		logger:  logger,
		memPool: memPool,
	}
}

// Start starts the application
func (app *Application) Start(ctx context.Context) error {
	span, ctx := observability.StartSpan(ctx, "app-start")
	defer span.Finish()

	app.logger.Info("Starting microservice application...")

	// Register metrics
	requestCounter := observability.RegisterCounter("requests_total", "Total number of requests", nil)
	requestDuration := observability.RegisterHistogram("request_duration_seconds", "Request duration", nil)

	// Setup HTTP handlers
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy","version":"v2.0.0"}`))
	})

	// Process endpoint with full integration
	mux.HandleFunc("/api/process", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestCounter.Inc()

		// Create span for tracing
		span, reqCtx := observability.StartSpan(r.Context(), "process-request")
		defer span.Finish()

		// Get request from pool
		req := app.memPool.Get()
		defer app.memPool.Put(req)

		req.ID = fmt.Sprintf("req-%d", time.Now().UnixNano())
		req.Timestamp = time.Now()

		// Log with context
		observability.WithContext(reqCtx).Info("Processing request", "request_id", req.ID)

		// Use rate limiter
		limiter := patterns.NewTokenBucket(100, time.Second)
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Simulate processing with worker pool
		jobs := make(chan int, 10)
		results := patterns.WorkerPool(reqCtx, app.config.WorkerPoolSize, jobs)

		// Send jobs
		go func() {
			for i := 0; i < 10; i++ {
				jobs <- i
			}
			close(jobs)
		}()

		// Collect results
		var processed int
		for range results {
			processed++
		}

		// Record duration
		duration := time.Since(start).Seconds()
		requestDuration.Observe(duration)

		// Response
		w.Header().Set("Content-Type", "application/json")
		response := fmt.Sprintf(`{"request_id":"%s","processed":%d,"duration":%.3f}`,
			req.ID, processed, duration)
		_, _ = w.Write([]byte(response))

		observability.WithContext(reqCtx).Info("Request completed",
			"request_id", req.ID,
			"duration", duration)
	})

	// Metrics endpoint
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		metrics := observability.ExportPrometheusMetrics()
		_, _ = w.Write([]byte(metrics))
	})

	// Create and start server
	app.server = &http.Server{
		Addr:    app.config.ServerAddr,
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		app.logger.Info("Server listening", "addr", app.config.ServerAddr)
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("Server failed", "error", err)
		}
	}()

	app.logger.Info("Microservice started successfully")
	return nil
}

// Stop gracefully stops the application
func (app *Application) Stop(ctx context.Context) error {
	span, ctx := observability.StartSpan(ctx, "app-stop")
	defer span.Finish()

	app.logger.Info("Stopping microservice application...")

	if app.server != nil {
		if err := app.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
	}

	app.logger.Info("Microservice stopped successfully")
	return nil
}

func main() {
	// Configuration
	config := &Config{
		ServerAddr:     ":8080",
		WorkerPoolSize: 5,
		LogLevel:       observability.InfoLevel,
	}

	// Create application
	app := NewApplication(config)

	// Start application
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nReceived interrupt signal, shutting down...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Stop(shutdownCtx); err != nil {
		log.Fatalf("Failed to stop application: %v", err)
	}

	fmt.Println("Application stopped successfully")
}
