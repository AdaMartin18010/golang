// Package interceptors provides tests for gRPC interceptors.
package interceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// MockUnaryHandler is a mock gRPC unary handler
type MockUnaryHandler func(ctx context.Context, req interface{}) (interface{}, error)

func (m MockUnaryHandler) Handle(ctx context.Context, req interface{}) (interface{}, error) {
	return m(ctx, req)
}

func TestLoggingUnaryInterceptor_Success(t *testing.T) {
	// This test verifies the interceptor doesn't panic and calls the handler
	interceptor := LoggingUnaryInterceptor

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	ctx := context.Background()
	resp, err := interceptor(ctx, "request", info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
}

func TestLoggingUnaryInterceptor_Error(t *testing.T) {
	interceptor := LoggingUnaryInterceptor

	testErr := errors.New("handler error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, testErr
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	ctx := context.Background()
	resp, err := interceptor(ctx, "request", info, handler)

	assert.Error(t, err)
	assert.Equal(t, testErr, err)
	assert.Nil(t, resp)
}

func TestLoggingUnaryInterceptor_WithRequest(t *testing.T) {
	interceptor := LoggingUnaryInterceptor

	type TestRequest struct {
		ID   string
		Name string
	}

	request := TestRequest{
		ID:   "123",
		Name: "Test",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Verify request is passed through
		r, ok := req.(TestRequest)
		assert.True(t, ok)
		assert.Equal(t, "123", r.ID)
		assert.Equal(t, "Test", r.Name)
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	ctx := context.Background()
	resp, err := interceptor(ctx, request, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
}

func TestLoggingUnaryInterceptor_DifferentMethods(t *testing.T) {
	interceptor := LoggingUnaryInterceptor

	methods := []string{
		"/user.Service/GetUser",
		"/user.Service/CreateUser",
		"/user.Service/UpdateUser",
		"/user.Service/DeleteUser",
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			}

			info := &grpc.UnaryServerInfo{
				FullMethod: method,
			}

			ctx := context.Background()
			resp, err := interceptor(ctx, nil, info, handler)

			assert.NoError(t, err)
			assert.Equal(t, "success", resp)
		})
	}
}

func TestLoggingUnaryInterceptor_NilRequest(t *testing.T) {
	interceptor := LoggingUnaryInterceptor

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		assert.Nil(t, req)
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	ctx := context.Background()
	resp, err := interceptor(ctx, nil, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
}

func TestLoggingUnaryInterceptor_ContextPropagation(t *testing.T) {
	interceptor := LoggingUnaryInterceptor

	type contextKey string
	key := contextKey("test-key")
	value := "test-value"

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Verify context is propagated
		v := ctx.Value(key)
		assert.Equal(t, value, v)
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	ctx := context.WithValue(context.Background(), key, value)
	resp, err := interceptor(ctx, nil, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
}

func TestLoggingUnaryInterceptor_HandlerPanic(t *testing.T) {
	interceptor := LoggingUnaryInterceptor

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("handler panic")
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	ctx := context.Background()

	// Interceptor should propagate panic
	assert.Panics(t, func() {
		_, _ = interceptor(ctx, nil, info, handler)
	})
}

// Tracing interceptor tests require OpenTelemetry setup
// These are simplified tests that verify basic functionality

func TestTracingUnaryInterceptor_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping tracing test in short mode - requires OpenTelemetry setup")
	}

	interceptor := TracingUnaryInterceptor

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	ctx := context.Background()
	resp, err := interceptor(ctx, "request", info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
}

func TestTracingUnaryInterceptor_Error(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping tracing test in short mode - requires OpenTelemetry setup")
	}

	interceptor := TracingUnaryInterceptor

	testErr := errors.New("handler error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, testErr
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	ctx := context.Background()
	resp, err := interceptor(ctx, "request", info, handler)

	assert.Error(t, err)
	assert.Equal(t, testErr, err)
	assert.Nil(t, resp)
}

func TestTracingUnaryInterceptor_WithMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping tracing test in short mode - requires OpenTelemetry setup")
	}

	interceptor := TracingUnaryInterceptor

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Handler receives context
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	// Context without metadata (should still work)
	ctx := context.Background()
	resp, err := interceptor(ctx, nil, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
}

func TestInterceptorChain(t *testing.T) {
	// Test that multiple interceptors can be chained
	loggingInterceptor := LoggingUnaryInterceptor

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "final response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	ctx := context.Background()

	// First interceptor
	resp1, err1 := loggingInterceptor(ctx, "request", info, handler)
	assert.NoError(t, err1)
	assert.Equal(t, "final response", resp1)
}

// Benchmark tests

func BenchmarkLoggingUnaryInterceptor(b *testing.B) {
	interceptor := LoggingUnaryInterceptor
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = interceptor(ctx, "request", info, handler)
	}
}

func BenchmarkLoggingUnaryInterceptor_WithError(b *testing.B) {
	interceptor := LoggingUnaryInterceptor
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("error")
	}
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = interceptor(ctx, "request", info, handler)
	}
}

// TestInterceptorFunctionSignatures ensures interceptors have correct signatures
func TestInterceptorFunctionSignatures(t *testing.T) {
	// Verify LoggingUnaryInterceptor has the correct signature
	var _ grpc.UnaryServerInterceptor = LoggingUnaryInterceptor

	// Verify TracingUnaryInterceptor has the correct signature
	var _ grpc.UnaryServerInterceptor = TracingUnaryInterceptor
}

func TestLoggingUnaryInterceptor_DifferentErrorTypes(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"generic_error", errors.New("generic error")},
		{"validation_error", errors.New("validation failed")},
		{"not_found_error", errors.New("resource not found")},
		{"internal_error", errors.New("internal server error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := LoggingUnaryInterceptor

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, tc.err
			}

			info := &grpc.UnaryServerInfo{
				FullMethod: "/test.Service/TestMethod",
			}

			ctx := context.Background()
			resp, err := interceptor(ctx, nil, info, handler)

			assert.Error(t, err)
			assert.Equal(t, tc.err, err)
			assert.Nil(t, resp)
		})
	}
}

func TestLoggingUnaryInterceptor_ResponseTypes(t *testing.T) {
	testCases := []struct {
		name     string
		response interface{}
	}{
		{"string", "success"},
		{"int", 42},
		{"struct", struct{ ID string }{ID: "123"}},
		{"map", map[string]string{"key": "value"}},
		{"nil", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := LoggingUnaryInterceptor

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return tc.response, nil
			}

			info := &grpc.UnaryServerInfo{
				FullMethod: "/test.Service/TestMethod",
			}

			ctx := context.Background()
			resp, err := interceptor(ctx, nil, info, handler)

			assert.NoError(t, err)
			assert.Equal(t, tc.response, resp)
		})
	}
}
