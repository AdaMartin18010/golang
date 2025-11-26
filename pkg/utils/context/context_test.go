package context

import (
	"context"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	ctx, cancel := WithTimeout(context.Background(), time.Second)
	defer cancel()

	if ctx == nil {
		t.Error("Expected non-nil context")
	}
}

func TestWithTimeoutSeconds(t *testing.T) {
	ctx, cancel := WithTimeoutSeconds(context.Background(), 1)
	defer cancel()

	if ctx == nil {
		t.Error("Expected non-nil context")
	}
}

func TestIsDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	if IsDone(ctx) {
		t.Error("Expected context not done")
	}
	cancel()
	if !IsDone(ctx) {
		t.Error("Expected context done")
	}
}

func TestGetValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")
	value := GetValue(ctx, "key")
	if value != "value" {
		t.Errorf("Expected 'value', got %v", value)
	}
}

func TestGetStringValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")
	value, ok := GetStringValue(ctx, "key")
	if !ok || value != "value" {
		t.Errorf("Expected ('value', true), got (%s, %v)", value, ok)
	}
}

func TestGetIntValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", 42)
	value, ok := GetIntValue(ctx, "key")
	if !ok || value != 42 {
		t.Errorf("Expected (42, true), got (%d, %v)", value, ok)
	}
}

func TestGetBoolValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", true)
	value, ok := GetBoolValue(ctx, "key")
	if !ok || !value {
		t.Errorf("Expected (true, true), got (%v, %v)", value, ok)
	}
}

func TestSleep(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := Sleep(ctx, 100*time.Millisecond)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDoWithTimeout(t *testing.T) {
	err := DoWithTimeout(context.Background(), time.Second, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestWithTraceID(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace123")
	traceID, ok := GetTraceID(ctx)
	if !ok || traceID != "trace123" {
		t.Errorf("Expected ('trace123', true), got (%s, %v)", traceID, ok)
	}
}

func TestWithUserID(t *testing.T) {
	ctx := WithUserID(context.Background(), "user123")
	userID, ok := GetUserID(ctx)
	if !ok || userID != "user123" {
		t.Errorf("Expected ('user123', true), got (%s, %v)", userID, ok)
	}
}

func TestContextBuilder(t *testing.T) {
	ctx := Chain(context.Background()).
		WithTraceID("trace123").
		WithUserID("user123").
		WithRequestID("req123").
		Build()

	traceID, _ := GetTraceID(ctx)
	if traceID != "trace123" {
		t.Errorf("Expected 'trace123', got %s", traceID)
	}

	userID, _ := GetUserID(ctx)
	if userID != "user123" {
		t.Errorf("Expected 'user123', got %s", userID)
	}
}

func TestWithStringValue(t *testing.T) {
	ctx := WithStringValue(context.Background(), "key", "value")
	value := GetStringKeyValue(ctx, "key")
	if value != "value" {
		t.Errorf("Expected 'value', got %v", value)
	}
}

func TestRetryWithContext(t *testing.T) {
	count := 0
	err := RetryWithContext(context.Background(), 3, func(ctx context.Context) error {
		count++
		if count < 3 {
			return context.DeadlineExceeded
		}
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 retries, got %d", count)
	}
}
