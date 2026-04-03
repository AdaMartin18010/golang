package limiter

import (
	"sync"
	"time"
)

// SlidingWindow implements the sliding window rate limiting algorithm
type SlidingWindow struct {
	limit      int
	windowSize time.Duration

	previousWindow struct {
		count     int
		timestamp time.Time
	}
	currentWindow struct {
		count     int
		timestamp time.Time
	}

	mu sync.Mutex
}

// NewSlidingWindow creates a new sliding window rate limiter
func NewSlidingWindow(limit int, windowSize time.Duration) *SlidingWindow {
	now := time.Now()
	return &SlidingWindow{
		limit:      limit,
		windowSize: windowSize,
		currentWindow: struct {
			count     int
			timestamp time.Time
		}{timestamp: now},
		previousWindow: struct {
			count     int
			timestamp time.Time
		}{timestamp: now.Add(-windowSize)},
	}
}

// Allow checks if a request should be allowed
func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(sw.currentWindow.timestamp)

	// Check if we need to move to a new window
	if elapsed >= sw.windowSize {
		sw.previousWindow = sw.currentWindow
		sw.currentWindow = struct {
			count     int
			timestamp time.Time
		}{timestamp: now}
		elapsed = 0
	}

	// Calculate weighted count
	weight := 1.0 - (elapsed.Seconds() / sw.windowSize.Seconds())
	weightedCount := float64(sw.previousWindow.count)*weight + float64(sw.currentWindow.count)

	if int(weightedCount) < sw.limit {
		sw.currentWindow.count++
		return true
	}

	return false
}

// Stats returns current window statistics
func (sw *SlidingWindow) Stats() (current, previous, weighted int) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(sw.currentWindow.timestamp)
	weight := 1.0 - (elapsed.Seconds() / sw.windowSize.Seconds())

	return sw.currentWindow.count, sw.previousWindow.count,
		int(float64(sw.previousWindow.count)*weight + float64(sw.currentWindow.count))
}

// Reset resets the window counters
func (sw *SlidingWindow) Reset() {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	sw.currentWindow = struct {
		count     int
		timestamp time.Time
	}{timestamp: now}
	sw.previousWindow = struct {
		count     int
		timestamp time.Time
	}{timestamp: now.Add(-sw.windowSize)}
}

// Limit returns the rate limit
func (sw *SlidingWindow) Limit() int {
	return sw.limit
}

// WindowSize returns the window size
func (sw *SlidingWindow) WindowSize() time.Duration {
	return sw.windowSize
}
