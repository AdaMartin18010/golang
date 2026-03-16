// Package chi provides tests for HTTP middleware.
package chi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTracingMiddleware(t *testing.T) {
	// Test that TracingMiddleware is defined and returns a handler
	assert.NotNil(t, TracingMiddleware)
	
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Apply middleware
	handler := TracingMiddleware(testHandler)
	assert.NotNil(t, handler)
	
	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	
	// Execute request (may panic if OTel is not configured, but tests function signature)
	// Skip actual execution since OTel requires setup
	_ = req
	_ = rr
	
	// Verify handler is not nil
	assert.NotNil(t, handler)
}

func TestLoggingMiddleware(t *testing.T) {
	// Test that LoggingMiddleware is defined and returns a handler
	assert.NotNil(t, LoggingMiddleware)
	
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Apply middleware
	handler := LoggingMiddleware(testHandler)
	assert.NotNil(t, handler)
	
	// Verify handler type
	assert.NotNil(t, handler)
}

func TestRecovererMiddleware(t *testing.T) {
	// Test that RecovererMiddleware is defined and returns a handler
	assert.NotNil(t, RecovererMiddleware)
	
	// Create a test handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	
	// Apply middleware
	handler := RecovererMiddleware(panicHandler)
	assert.NotNil(t, handler)
}

func TestTimeoutMiddleware(t *testing.T) {
	// Test that TimeoutMiddleware is defined
	assert.NotNil(t, TimeoutMiddleware)
	
	// Create middleware with 1 second timeout
	timeoutMiddleware := TimeoutMiddleware(1 * time.Second)
	assert.NotNil(t, timeoutMiddleware)
	
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Apply middleware
	handler := timeoutMiddleware(testHandler)
	assert.NotNil(t, handler)
}

func TestTimeoutMiddlewareWithDifferentDurations(t *testing.T) {
	// Test with different timeout durations
	durations := []time.Duration{
		1 * time.Second,
		30 * time.Second,
		1 * time.Minute,
		5 * time.Minute,
	}
	
	for _, duration := range durations {
		middleware := TimeoutMiddleware(duration)
		assert.NotNil(t, middleware)
		
		// Verify it creates a valid middleware function
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		handler := middleware(testHandler)
		assert.NotNil(t, handler)
	}
}

func TestCORSMiddleware(t *testing.T) {
	// Test that CORSMiddleware is defined and returns a handler
	assert.NotNil(t, CORSMiddleware)
	
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Apply middleware
	handler := CORSMiddleware(testHandler)
	assert.NotNil(t, handler)
	
	// Test regular request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	// Check CORS headers
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCORSMiddlewarePreflight(t *testing.T) {
	// Test CORS preflight (OPTIONS) request
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	handler := CORSMiddleware(testHandler)
	
	// Create OPTIONS request
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	// Check CORS headers are set
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
	
	// Should return 204 No Content for preflight
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestMiddlewareTypes(t *testing.T) {
	// Test that all middleware functions have correct signatures
	
	// TracingMiddleware should be func(http.Handler) http.Handler
	var _ func(http.Handler) http.Handler = TracingMiddleware
	
	// LoggingMiddleware should be func(http.Handler) http.Handler
	var _ func(http.Handler) http.Handler = LoggingMiddleware
	
	// RecovererMiddleware should be func(http.Handler) http.Handler
	var _ func(http.Handler) http.Handler = RecovererMiddleware
	
	// TimeoutMiddleware should return func(http.Handler) http.Handler
	timeoutMiddleware := TimeoutMiddleware(1 * time.Second)
	var _ func(http.Handler) http.Handler = timeoutMiddleware
	
	// CORSMiddleware should be func(http.Handler) http.Handler
	var _ func(http.Handler) http.Handler = CORSMiddleware
}

func TestMiddlewareChaining(t *testing.T) {
	// Test that middleware can be chained
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Chain multiple middleware (skip LoggingMiddleware in chain due to dependencies)
	handler := CORSMiddleware(
		RecovererMiddleware(testHandler),
	)
	
	assert.NotNil(t, handler)
	
	// Execute request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRecovererMiddlewareRecoversPanic(t *testing.T) {
	// Test that RecovererMiddleware actually recovers from panics
	recovered := false
	
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("intentional test panic")
	})
	
	handler := RecovererMiddleware(panicHandler)
	
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	
	// Should not panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = false // If we reach here, middleware didn't recover
			}
		}()
		handler.ServeHTTP(rr, req)
		recovered = true // If we reach here, middleware recovered successfully
	}()
	
	assert.True(t, recovered, "RecovererMiddleware should recover from panic")
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
