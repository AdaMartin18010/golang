package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"rate-limiter/pkg/limiter"
)

// Limiter interface for rate limiting
type Limiter interface {
	Allow() bool
}

// Result represents a rate limit check result
type Result struct {
	Allowed    bool
	Limit      int
	Remaining  int
	ResetAt    time.Time
	RetryAfter time.Duration
}

// KeyFunc generates a rate limit key from a request
type KeyFunc func(r *http.Request) string

// RateLimiterMiddleware creates a rate limiting middleware
func RateLimiterMiddleware(limiters map[string]Limiter, keyFunc KeyFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			
			limiter, ok := limiters[key]
			if !ok {
				// Try default limiter
				limiter, ok = limiters["default"]
				if !ok {
					next.ServeHTTP(w, r)
					return
				}
			}
			
			allowed := limiter.Allow()
			
			// Add rate limit headers
			headers := w.Header()
			headers.Set("X-RateLimit-Limit", "100")
			
			if !allowed {
				headers.Set("Retry-After", "60")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// KeyByIP generates key by client IP
func KeyByIP(r *http.Request) string {
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}
	return "ip:" + ip
}

// KeyByUser generates key by user ID
func KeyByUser(headerName string) KeyFunc {
	return func(r *http.Request) string {
		userID := r.Header.Get(headerName)
		if userID == "" {
			return KeyByIP(r)
		}
		return "user:" + userID
	}
}

// KeyByAPIKey generates key by API key
func KeyByAPIKey(r *http.Request) string {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		return KeyByIP(r)
	}
	return "apikey:" + apiKey
}

// AdaptiveRateLimiter adapts rate limits based on behavior
type AdaptiveRateLimiter struct {
	baseLimiter *limiter.TokenBucket
	violations  int
	lastViolation time.Time
	mu          sync.Mutex
}

// NewAdaptiveRateLimiter creates an adaptive rate limiter
func NewAdaptiveRateLimiter(capacity, refillRate float64) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		baseLimiter: limiter.NewTokenBucket(capacity, refillRate),
	}
}

// Allow checks if request is allowed with adaptive behavior
func (a *AdaptiveRateLimiter) Allow() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	// Reduce limit if too many violations
	if a.violations > 10 && time.Since(a.lastViolation) < time.Hour {
		// Temporary reduction - implement with a reduced limiter
		return a.baseLimiter.Allow()
	}
	
	allowed := a.baseLimiter.Allow()
	if !allowed {
		a.violations++
		a.lastViolation = time.Now()
	}
	
	return allowed
}
