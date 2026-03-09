// Package handlers provides tests for gRPC health handlers.
package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	healthpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/healthpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNewHealthHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	handler := NewHealthHandler(logger)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
	assert.NotNil(t, handler.checkers)
	assert.NotNil(t, handler.isReady)
	assert.True(t, handler.isReady())
}

func TestNewHealthHandler_NilLogger(t *testing.T) {
	handler := NewHealthHandler(nil)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
}

func TestNewHealthHandlerWithReadyFunc(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	isReady := func() bool { return false }

	handler := NewHealthHandlerWithReadyFunc(logger, isReady)

	assert.NotNil(t, handler)
	assert.False(t, handler.isReady())
}

func TestNewHealthHandlerWithReadyFunc_NilFunc(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	handler := NewHealthHandlerWithReadyFunc(logger, nil)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.isReady)
}

func TestHealthHandler_Check_Serving(t *testing.T) {
	handler := NewHealthHandler(nil)

	req := &emptypb.Empty{}
	resp, err := handler.Check(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, healthpb.HealthResponse_SERVING, resp.Status)
}

func TestHealthHandler_Check_NotReady(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	isReady := func() bool { return false }

	handler := NewHealthHandlerWithReadyFunc(logger, isReady)

	req := &emptypb.Empty{}
	resp, err := handler.Check(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, healthpb.HealthResponse_NOT_SERVING, resp.Status)
}

func TestHealthHandler_Check_WithSuccessfulChecker(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker := NewSimpleHealthChecker("test", func(ctx context.Context) error {
		return nil
	})

	handler := NewHealthHandler(logger, checker)

	req := &emptypb.Empty{}
	resp, err := handler.Check(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, healthpb.HealthResponse_SERVING, resp.Status)
}

func TestHealthHandler_Check_WithFailingChecker(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker := NewSimpleHealthChecker("test", func(ctx context.Context) error {
		return errors.New("check failed")
	})

	handler := NewHealthHandler(logger, checker)

	req := &emptypb.Empty{}
	resp, err := handler.Check(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, healthpb.HealthResponse_NOT_SERVING, resp.Status)
}

func TestHealthHandler_Check_MultipleCheckers(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker1 := NewSimpleHealthChecker("checker1", func(ctx context.Context) error {
		return nil
	})
	checker2 := NewSimpleHealthChecker("checker2", func(ctx context.Context) error {
		return nil
	})

	handler := NewHealthHandler(logger, checker1, checker2)

	req := &emptypb.Empty{}
	resp, err := handler.Check(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, healthpb.HealthResponse_SERVING, resp.Status)
}

func TestHealthHandler_Check_MultipleCheckersOneFailing(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker1 := NewSimpleHealthChecker("checker1", func(ctx context.Context) error {
		return nil
	})
	checker2 := NewSimpleHealthChecker("checker2", func(ctx context.Context) error {
		return errors.New("check failed")
	})

	handler := NewHealthHandler(logger, checker1, checker2)

	req := &emptypb.Empty{}
	resp, err := handler.Check(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, healthpb.HealthResponse_NOT_SERVING, resp.Status)
}

func TestHealthHandler_CheckWithDetails_Serving(t *testing.T) {
	handler := NewHealthHandler(nil)

	details, err := handler.CheckWithDetails(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, healthpb.HealthResponse_SERVING, details.Status)
	assert.Empty(t, details.Component)
}

func TestHealthHandler_CheckWithDetails_NotReady(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	isReady := func() bool { return false }

	handler := NewHealthHandlerWithReadyFunc(logger, isReady)

	details, err := handler.CheckWithDetails(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, healthpb.HealthResponse_NOT_SERVING, details.Status)
	assert.Equal(t, "service not ready", details.Message)
}

func TestHealthHandler_CheckWithDetails_WithSuccessfulChecker(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker := NewSimpleHealthChecker("test-checker", func(ctx context.Context) error {
		return nil
	})

	handler := NewHealthHandler(logger, checker)

	details, err := handler.CheckWithDetails(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, healthpb.HealthResponse_SERVING, details.Status)
	assert.Len(t, details.Component, 1)
	assert.Equal(t, healthpb.HealthResponse_SERVING, details.Component["checker_0"].Status)
}

func TestHealthHandler_CheckWithDetails_WithFailingChecker(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker := NewSimpleHealthChecker("test-checker", func(ctx context.Context) error {
		return errors.New("connection failed")
	})

	handler := NewHealthHandler(logger, checker)

	details, err := handler.CheckWithDetails(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, healthpb.HealthResponse_NOT_SERVING, details.Status)
	assert.Len(t, details.Component, 1)
	assert.Equal(t, healthpb.HealthResponse_NOT_SERVING, details.Component["checker_0"].Status)
	assert.Equal(t, "connection failed", details.Component["checker_0"].Message)
}

func TestHealthHandler_CheckWithDetails_MultipleCheckersMixed(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker1 := NewSimpleHealthChecker("db-checker", func(ctx context.Context) error {
		return nil
	})
	checker2 := NewSimpleHealthChecker("cache-checker", func(ctx context.Context) error {
		return errors.New("cache unavailable")
	})

	handler := NewHealthHandler(logger, checker1, checker2)

	details, err := handler.CheckWithDetails(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, healthpb.HealthResponse_NOT_SERVING, details.Status)
	assert.Len(t, details.Component, 2)
	assert.Equal(t, healthpb.HealthResponse_SERVING, details.Component["checker_0"].Status)
	assert.Equal(t, healthpb.HealthResponse_NOT_SERVING, details.Component["checker_1"].Status)
}

func TestGetCheckerName(t *testing.T) {
	checker := NewSimpleHealthChecker("my-checker", nil)

	name := getCheckerName(checker, 5)

	assert.Equal(t, "checker_5", name)
}

func TestSimpleHealthChecker_Check_Success(t *testing.T) {
	called := false
	checker := NewSimpleHealthChecker("test", func(ctx context.Context) error {
		called = true
		return nil
	})

	err := checker.Check(context.Background())

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestSimpleHealthChecker_Check_Error(t *testing.T) {
	testErr := errors.New("check error")
	checker := NewSimpleHealthChecker("test", func(ctx context.Context) error {
		return testErr
	})

	err := checker.Check(context.Background())

	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

func TestSimpleHealthChecker_Check_NilFunction(t *testing.T) {
	checker := NewSimpleHealthChecker("test", nil)

	err := checker.Check(context.Background())

	assert.NoError(t, err)
}

func TestNewSimpleHealthChecker(t *testing.T) {
	checkFunc := func(ctx context.Context) error { return nil }
	checker := NewSimpleHealthChecker("my-checker", checkFunc)

	assert.NotNil(t, checker)
	assert.Equal(t, "my-checker", checker.name)
	assert.NotNil(t, checker.check)
}

func TestDatabaseHealthChecker_Check_Success(t *testing.T) {
	called := false
	checker := NewDatabaseHealthChecker(func(ctx context.Context) error {
		called = true
		return nil
	})

	err := checker.Check(context.Background())

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestDatabaseHealthChecker_Check_Error(t *testing.T) {
	testErr := errors.New("database connection failed")
	checker := NewDatabaseHealthChecker(func(ctx context.Context) error {
		return testErr
	})

	err := checker.Check(context.Background())

	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

func TestDatabaseHealthChecker_Check_NilPing(t *testing.T) {
	checker := NewDatabaseHealthChecker(nil)

	err := checker.Check(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database ping function not set")
}

func TestNewDatabaseHealthChecker(t *testing.T) {
	pingFunc := func(ctx context.Context) error { return nil }
	checker := NewDatabaseHealthChecker(pingFunc)

	assert.NotNil(t, checker)
	assert.NotNil(t, checker.ping)
}

func TestHealthDetailsStruct(t *testing.T) {
	details := &HealthDetails{
		Status:    healthpb.HealthResponse_SERVING,
		Message:   "all systems operational",
		Component: make(map[string]HealthComponentStatus),
	}

	assert.Equal(t, healthpb.HealthResponse_SERVING, details.Status)
	assert.Equal(t, "all systems operational", details.Message)
	assert.NotNil(t, details.Component)
}

func TestHealthComponentStatusStruct(t *testing.T) {
	status := HealthComponentStatus{
		Status:  healthpb.HealthResponse_SERVING,
		Message: "operational",
	}

	assert.Equal(t, healthpb.HealthResponse_SERVING, status.Status)
	assert.Equal(t, "operational", status.Message)
}

// TestHealthCheckerInterface ensures our checkers implement the interface
func TestHealthCheckerInterface(t *testing.T) {
	// SimpleHealthChecker should implement HealthChecker
	var _ HealthChecker = (*SimpleHealthChecker)(nil)

	// DatabaseHealthChecker should implement HealthChecker
	var _ HealthChecker = (*DatabaseHealthChecker)(nil)
}

// TestHealthServiceServerInterface ensures HealthHandler implements the interface
func TestHealthServiceServerInterface(t *testing.T) {
	// This is verified by the var declaration at the end of health_handler.go
	// var _ healthpb.HealthServiceServer = (*HealthHandler)(nil)
	// If it doesn't compile, the test will fail
	assert.True(t, true)
}

func TestHealthHandler_WithMultipleCheckersAllFail(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker1 := NewSimpleHealthChecker("db", func(ctx context.Context) error {
		return errors.New("db down")
	})
	checker2 := NewSimpleHealthChecker("cache", func(ctx context.Context) error {
		return errors.New("cache down")
	})

	handler := NewHealthHandler(logger, checker1, checker2)

	details, err := handler.CheckWithDetails(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, healthpb.HealthResponse_NOT_SERVING, details.Status)
	assert.Len(t, details.Component, 2)
	assert.Equal(t, "db down", details.Component["checker_0"].Message)
	assert.Equal(t, "cache down", details.Component["checker_1"].Message)
}

func TestHealthHandler_Check_ContextCancellation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker := NewSimpleHealthChecker("slow-checker", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	handler := NewHealthHandler(logger, checker)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := &emptypb.Empty{}
	resp, err := handler.Check(ctx, req)

	assert.NoError(t, err) // Check doesn't propagate checker errors as handler errors
	assert.NotNil(t, resp)
	// Result depends on whether the context cancellation is detected
}

func TestHealthHandler_CheckWithDetails_ContextCancellation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	checker := NewSimpleHealthChecker("slow-checker", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	handler := NewHealthHandler(logger, checker)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	details, err := handler.CheckWithDetails(ctx)

	assert.NoError(t, err) // CheckWithDetails returns details, not error for checker failures
	assert.NotNil(t, details)
}
