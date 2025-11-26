package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/golang/pkg/logger"
	"log/slog"
)

func TestRecoveryMiddleware(t *testing.T) {
	log := logger.NewLogger(slog.LevelError)
	config := RecoveryConfig{
		Logger: log,
	}

	r := chi.NewRouter()
	r.Use(RecoveryMiddleware(config))
	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	// 应该不会导致测试panic
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestRecoveryMiddleware_NoPanic(t *testing.T) {
	log := logger.NewLogger(slog.LevelError)
	config := RecoveryConfig{
		Logger: log,
	}

	r := chi.NewRouter()
	r.Use(RecoveryMiddleware(config))
	r.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/ok", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
