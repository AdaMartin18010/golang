package errors

import (
	"errors"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *AppError
		want    string
		wantErr bool
	}{
		{
			name: "error without cause",
			err:  NewNotFoundError("user", "123"),
			want: "NOT_FOUND: user with id 123 not found",
		},
		{
			name: "error with cause",
			err: &AppError{
				Code:    ErrCodeInternal,
				Message: "test error",
				Cause:   errors.New("underlying error"),
			},
			want: "INTERNAL_ERROR: test error: underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &AppError{
		Code:    ErrCodeInternal,
		Message: "test error",
		Cause:   cause,
	}

	if got := err.Unwrap(); got != cause {
		t.Errorf("AppError.Unwrap() = %v, want %v", got, cause)
	}
}

func TestAppError_WithDetails(t *testing.T) {
	err := NewNotFoundError("user", "123")
	err.WithDetails("field", "value")

	if err.Details == nil {
		t.Error("Details should not be nil")
	}
	if err.Details["field"] != "value" {
		t.Errorf("Details[field] = %v, want value", err.Details["field"])
	}
}

func TestAppError_WithTraceID(t *testing.T) {
	err := NewNotFoundError("user", "123")
	traceID := "trace-123"
	err.WithTraceID(traceID)

	if err.TraceID != traceID {
		t.Errorf("TraceID = %v, want %v", err.TraceID, traceID)
	}
}

func TestAppError_HTTPStatusCode(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		want int
	}{
		{
			name: "not found",
			err:  NewNotFoundError("user", "123"),
			want: 404,
		},
		{
			name: "invalid input",
			err:  NewInvalidInputError("invalid"),
			want: 400,
		},
		{
			name: "internal error",
			err:  NewInternalError("error", nil),
			want: 500,
		},
		{
			name: "unauthorized",
			err:  NewUnauthorizedError("unauthorized"),
			want: 401,
		},
		{
			name: "forbidden",
			err:  NewForbiddenError("forbidden"),
			want: 403,
		},
		{
			name: "conflict",
			err:  NewConflictError("conflict"),
			want: 409,
		},
		{
			name: "timeout",
			err:  NewTimeoutError("timeout"),
			want: 504,
		},
		{
			name: "rate limit",
			err:  NewRateLimitError("rate limit"),
			want: 429,
		},
		{
			name: "service unavailable",
			err:  NewServiceUnavailableError("unavailable"),
			want: 503,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.HTTPStatusCode(); got != tt.want {
				t.Errorf("AppError.HTTPStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppError_IsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		want bool
	}{
		{
			name: "not retryable",
			err:  NewNotFoundError("user", "123"),
			want: false,
		},
		{
			name: "retryable timeout",
			err:  NewTimeoutError("timeout"),
			want: true,
		},
		{
			name: "retryable rate limit",
			err:  NewRateLimitError("rate limit"),
			want: true,
		},
		{
			name: "retryable service unavailable",
			err:  NewServiceUnavailableError("unavailable"),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsRetryable(); got != tt.want {
				t.Errorf("AppError.IsRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromDomainError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want *AppError
	}{
		{
			name: "nil error",
			err:  nil,
			want: nil,
		},
		{
			name: "AppError",
			err:  NewNotFoundError("user", "123"),
			want: NewNotFoundError("user", "123"),
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: &AppError{
				Code:       ErrCodeInternal,
				Message:    "Internal error",
				Category:   CategoryServer,
				HTTPStatus: 500,
				Retryable:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromDomainError(tt.err)
			if tt.want == nil {
				if got != nil {
					t.Errorf("FromDomainError() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Error("FromDomainError() = nil, want AppError")
				return
			}
			if got.Code != tt.want.Code {
				t.Errorf("FromDomainError().Code = %v, want %v", got.Code, tt.want.Code)
			}
		})
	}
}

func TestErrorCategory_HTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		category ErrorCategory
		want     int
	}{
		{
			name:     "client error",
			category: CategoryClient,
			want:     400,
		},
		{
			name:     "server error",
			category: CategoryServer,
			want:     500,
		},
		{
			name:     "business error",
			category: CategoryBusiness,
			want:     400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.HTTPStatus(); got != tt.want {
				t.Errorf("ErrorCategory.HTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
