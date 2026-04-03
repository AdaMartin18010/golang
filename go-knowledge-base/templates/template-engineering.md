# EC-XXX: [Engineering/Cloud-Native Topic] - Quick Contribution Template

> **Dimension**: Engineering & Cloud-Native (EC)
> **Level**: S/A/B - Target >[TODO: 15KB/10KB/5KB]
> **Status**: [TODO: Draft/Review/Complete]
> **Tags**: #[TODO: cloud-native] #[TODO: pattern] #[TODO: go]
> **Author**: [TODO: Your Name]
> **Created**: [TODO: YYYY-MM-DD]
> **Estimated Reading Time**: [TODO: XX minutes]

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Introduction](#introduction)
3. [Pattern Overview](#pattern-overview)
4. [Implementation](#implementation)
5. [Production Considerations](#production-considerations)
6. [Visual Representations](#visual-representations)
7. [Code Examples](#code-examples)
8. [Best Practices](#best-practices)
9. [Cross-References](#cross-references)
10. [References](#references)

---

## Executive Summary

[TODO: 2-3 paragraph overview]

**Pattern at a Glance**:
- **Category**: [TODO: Creational/Structural/Behavioral/Architectural]
- **Difficulty**: [TODO: Beginner/Intermediate/Advanced]
- **Use When**: [TODO: Key condition]
- **Trade-offs**: [TODO: Key trade-off]

---

## Introduction

### What is [Pattern Name]?

[TODO: One-paragraph definition]

**Also Known As**: [TODO: Alternative names]

### Why Use This Pattern?

**Benefits**:
- [TODO: Benefit 1]
- [TODO: Benefit 2]
- [TODO: Benefit 3]

**Challenges**:
- [TODO: Challenge 1]
- [TODO: Challenge 2]

### When to Use

**Appropriate For**:
- [TODO: Use case 1]
- [TODO: Use case 2]
- [TODO: Use case 3]

**Not Appropriate For**:
- [TODO: Anti-use case 1]
- [TODO: Anti-use case 2]

### Prerequisites

- [TODO: [Go Concurrency](../02-Language-Design/02-Language-Features/03-Goroutines.md)]
- [TODO: [Context Package](../04-Technology-Stack/01-Core-Library/04-Context-Package.md)]

---

## Pattern Overview

### Problem Statement

[TODO: Describe the problem this pattern solves]

### Solution

[TODO: How this pattern solves the problem]

### Structure

#### Participants

| Role | Responsibility | In Go |
|------|----------------|-------|
| [Role 1] | [TODO: What it does] | [TODO: Interface/Struct] |
| [Role 2] | [TODO: What it does] | [TODO: Interface/Struct] |
| [Role 3] | [TODO: What it does] | [TODO: Interface/Struct] |

#### Collaboration Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    COLLABORATION SEQUENCE                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   Client                                                        │
│     │                                                           │
│     │  1. [Operation]                                           │
│     ▼                                                           │
│   ┌──────────┐     2. [Call]     ┌──────────┐                  │
│   │ [Role 1] │ ────────────────▶ │ [Role 2] │                  │
│   └──────────┘                   └────┬─────┘                  │
│                                       │                         │
│                                       │ 3. [Process]            │
│                                       ▼                         │
│                                 ┌──────────┐                   │
│                                 │ [Role 3] │                   │
│                                 └────┬─────┘                   │
│                                       │                         │
│     │  4. [Result]                   │ 4. [Result]             │
│     ▼                                ▼                         │
│   [Complete]                                                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Implementation

### Basic Implementation

```go
// file: basic.go
// description: Minimal working implementation
package pattern

import (
    "context"
    "fmt"
)

// [Interface] defines the contract for [role].
type Interface interface {
    Method(ctx context.Context, input Input) (Output, error)
}

// [Implementation] implements [Interface].
type Implementation struct {
    // [TODO: Dependencies]
}

// New creates a new [Implementation].
func New() *Implementation {
    return &Implementation{}
}

// Method implements [Interface].
func (i *Implementation) Method(ctx context.Context, input Input) (Output, error) {
    // [TODO: Basic implementation]
    return Output{}, nil
}
```

### Production Implementation

```go
// file: production.go
// description: Production-ready implementation with all features
package pattern

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"

    "go.uber.org/zap"
)

var (
    ErrInvalidInput = errors.New("invalid input")
    ErrTimeout      = errors.New("operation timed out")
    ErrCircuitOpen  = errors.New("circuit breaker is open")
)

// Config holds configuration options.
type Config struct {
    Timeout       time.Duration
    RetryAttempts int
    RetryDelay    time.Duration
    Logger        *zap.Logger
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
    return Config{
        Timeout:       30 * time.Second,
        RetryAttempts: 3,
        RetryDelay:    100 * time.Millisecond,
    }
}

func (c *Config) validate() error {
    if c.Timeout <= 0 {
        return fmt.Errorf("%w: timeout must be positive", ErrInvalidInput)
    }
    if c.RetryAttempts < 0 {
        return fmt.Errorf("%w: retry attempts cannot be negative", ErrInvalidInput)
    }
    return nil
}

// Component is a production-ready [component] implementation.
type Component struct {
    config Config
    mu     sync.RWMutex
    state  State
    logger *zap.Logger
}

// State represents the component state.
type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

// New creates a new Component with the given configuration.
func New(cfg Config) (*Component, error) {
    if err := cfg.validate(); err != nil {
        return nil, err
    }
    
    if cfg.Logger == nil {
        cfg.Logger = zap.NewNop()
    }
    
    return &Component{
        config: cfg,
        state:  StateClosed,
        logger: cfg.Logger,
    }, nil
}

// Execute performs the main operation with full error handling.
func (c *Component) Execute(ctx context.Context, input Input) (Output, error) {
    c.mu.RLock()
    state := c.state
    c.mu.RUnlock()
    
    if state == StateOpen {
        return Output{}, ErrCircuitOpen
    }
    
    ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
    defer cancel()
    
    var result Output
    var err error
    
    // Retry logic
    for attempt := 0; attempt < c.config.RetryAttempts; attempt++ {
        result, err = c.executeOnce(ctx, input)
        if err == nil {
            c.recordSuccess()
            return result, nil
        }
        
        if !isRetryable(err) {
            break
        }
        
        c.logger.Warn("attempt failed, retrying",
            zap.Int("attempt", attempt+1),
            zap.Error(err),
        )
        
        select {
        case <-time.After(c.config.RetryDelay):
            continue
        case <-ctx.Done():
            return Output{}, ctx.Err()
        }
    }
    
    c.recordFailure()
    return Output{}, fmt.Errorf("operation failed after %d attempts: %w", 
        c.config.RetryAttempts, err)
}

func (c *Component) executeOnce(ctx context.Context, input Input) (Output, error) {
    // [TODO: Actual implementation]
    return Output{}, nil
}

func (c *Component) recordSuccess() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if c.state == StateHalfOpen {
        c.state = StateClosed
        c.logger.Info("circuit closed")
    }
}

func (c *Component) recordFailure() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.state = StateOpen
    c.logger.Warn("circuit opened")
}

func isRetryable(err error) bool {
    // [TODO: Implement retryable error detection]
    return true
}
```

### Variations

#### Variation 1: [Name]

[TODO: Description and when to use]

```go
// file: variation1.go
// description: [Variation description]
func Variation1() {
    // [TODO: Implementation]
}
```

#### Variation 2: [Name]

[TODO: Description and when to use]

```go
// file: variation2.go
// description: [Variation description]
func Variation2() {
    // [TODO: Implementation]
}
```

---

## Production Considerations

### Performance Characteristics

| Metric | Value | Notes |
|--------|-------|-------|
| Time Complexity | O([TODO]) | [TODO: Explanation] |
| Space Complexity | O([TODO]) | [TODO: Explanation] |
| Latency (p50) | [TODO] ms | Under normal load |
| Latency (p99) | [TODO] ms | Under normal load |
| Throughput | [TODO] ops/sec | Single instance |

### Scalability

**Horizontal Scaling**:
- [TODO: How the pattern scales horizontally]
- [TODO: Limitations]

**Vertical Scaling**:
- [TODO: Resource requirements]
- [TODO: Bottlenecks]

### Reliability

| Scenario | Behavior | Mitigation |
|----------|----------|------------|
| Network partition | [TODO] | [TODO] |
| Node failure | [TODO] | [TODO] |
| High load | [TODO] | [TODO] |

### Monitoring

**Key Metrics**:
- `operation_latency_seconds` - Histogram of operation latency
- `operation_total` - Counter of operations by status
- `state_changes_total` - Counter of state transitions

**Health Checks**:

```go
// file: health.go
// description: Health check implementation
package pattern

import (
    "context"
    "net/http"
)

// HealthChecker implements health check interface.
type HealthChecker struct {
    component *Component
}

// Check performs health check.
func (h *HealthChecker) Check(ctx context.Context) error {
    return h.component.Health(ctx)
}

// HTTPHandler returns HTTP handler for health checks.
func (h *HealthChecker) HTTPHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := h.Check(r.Context()); err != nil {
            http.Error(w, err.Error(), http.StatusServiceUnavailable)
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("healthy"))
    }
}
```

### Security Considerations

- [TODO: Security concern 1 and mitigation]
- [TODO: Security concern 2 and mitigation]

---

## Visual Representations

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    SYSTEM ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌──────────────┐                                              │
│   │    Client    │                                              │
│   └──────┬───────┘                                              │
│          │                                                       │
│          ▼                                                       │
│   ┌──────────────┐                                              │
│   │   API Layer  │                                              │
│   └──────┬───────┘                                              │
│          │                                                       │
│          ▼                                                       │
│   ┌──────────────────────────────────────────┐                  │
│   │           [Pattern Implementation]       │                  │
│   │  ┌──────────┐  ┌──────────┐  ┌─────────┐ │                  │
│   │  │Component1│  │Component2│  │Component3│ │                  │
│   │  └──────────┘  └──────────┘  └─────────┘ │                  │
│   └────────┬─────────────────────────────────┘                  │
│            │                                                     │
│            ▼                                                     │
│   ┌──────────────────────────────────────────┐                  │
│   │         Infrastructure Layer             │                  │
│   │  ┌────────┐ ┌────────┐ ┌────────┐       │                  │
│   │  │Database│ │ Cache  │ │ Message│       │                  │
│   │  └────────┘ └────────┘ └────────┘       │                  │
│   └──────────────────────────────────────────┘                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### State Machine

```
                    ┌──────────────┐
         ┌─────────│    CLOSED    │◀────────┐
         │         │  (Normal)    │         │
         │         └──────┬───────┘         │
    Success│              │Failure          │Success
         │                ▼                │
         │         ┌──────────────┐        │
         └────────▶│    OPEN      │────────┘
                   │ (Rejecting)  │
                   └──────┬───────┘
                          │
                          │Timeout
                          ▼
                   ┌──────────────┐
                   │  HALF-OPEN   │
                   │ (Testing)    │
                   └──────────────┘
```

### Decision Tree

```
Should you use [Pattern]?
│
├── Do you have [condition 1]?
│   ├── Yes → Continue
│   └── No → Consider [Alternative A]
│
├── Is [condition 2] critical?
│   ├── Yes → [Pattern] is recommended
│   └── No → Continue
│
└── Do you need [feature]?
    ├── Yes → [Pattern] with [Variation]
    └── No → [Alternative B] might be simpler
```

### Comparison Matrix

| Pattern | Latency | Complexity | Resilience | Use Case |
|---------|---------|------------|------------|----------|
| [Pattern A] | Low | Low | Basic | [Scenario] |
| [Pattern B] | Med | Med | Good | [Scenario] |
| **[This Pattern]** | Med | High | Excellent | [Scenario] |

---

## Code Examples

### Example 1: [Scenario Name]

**Context**: [TODO: When to use this example]

```go
// file: example1.go
// description: [Scenario] implementation
package main

import (
    "context"
    "fmt"
    "log"
    "time"
)

func main() {
    cfg := pattern.DefaultConfig()
    component, err := pattern.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    result, err := component.Execute(ctx, pattern.Input{
        // [TODO: Input fields]
    })
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    
    fmt.Printf("Result: %+v\n", result)
}
```

### Example 2: Integration with [Technology]

```go
// file: integration.go
// description: Integration with [technology]
package main

import (
    "context"
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewDevelopment()
    cfg := pattern.Config{
        Timeout:       30 * time.Second,
        RetryAttempts: 3,
        Logger:        logger,
    }
    
    component, err := pattern.New(cfg)
    if err != nil {
        logger.Fatal("failed to create component", zap.Error(err))
    }
    
    r := gin.Default()
    
    r.POST("/api/operation", func(c *gin.Context) {
        var req Request
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        
        result, err := component.Execute(c.Request.Context(), req.ToInput())
        if err != nil {
            logger.Error("operation failed", zap.Error(err))
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(200, result)
    })
    
    r.Run(":8080")
}
```

### Example 3: Testing

```go
// file: example_test.go
// description: Integration test example
package main

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Setup
    cfg := pattern.DefaultConfig()
    cfg.Timeout = 5 * time.Second
    
    component, err := pattern.New(cfg)
    require.NoError(t, err)
    
    ctx := context.Background()
    
    t.Run("successful operation", func(t *testing.T) {
        result, err := component.Execute(ctx, pattern.Input{
            // [TODO: Valid input]
        })
        
        assert.NoError(t, err)
        assert.NotNil(t, result)
    })
    
    t.Run("timeout handling", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
        defer cancel()
        
        _, err := component.Execute(ctx, pattern.Input{
            // [TODO: Input]
        })
        
        assert.ErrorIs(t, err, context.DeadlineExceeded)
    })
}
```

---

## Best Practices

### Design Principles

1. **Fail Fast**
   ```go
   // ✅ Good: Validate early
   func (c *Component) Execute(ctx context.Context, input Input) (Output, error) {
       if err := input.Validate(); err != nil {
           return Output{}, fmt.Errorf("invalid input: %w", err)
       }
       // ...
   }
   ```

2. **Observability**
   ```go
   // ✅ Good: Add structured logging
   logger.Info("operation started",
       zap.String("operation_id", opID),
       zap.Duration("timeout", cfg.Timeout),
   )
   ```

3. **Graceful Degradation**
   ```go
   // ✅ Good: Have fallback strategies
   result, err := primary.Execute(ctx, input)
   if err != nil {
       logger.Warn("primary failed, using fallback", zap.Error(err))
       result, err = fallback.Execute(ctx, input)
   }
   ```

### Configuration Guidelines

| Parameter | Default | Range | Recommendation |
|-----------|---------|-------|----------------|
| Timeout | 30s | 1s - 300s | Start with 30s, adjust based on p99 |
| Retries | 3 | 0 - 10 | 3 for idempotent ops, 0 for non-idempotent |
| [TODO: More] | [TODO] | [TODO] | [TODO] |

### Common Pitfalls

| Pitfall | Problem | Solution |
|---------|---------|----------|
| [Pitfall 1] | [TODO: Description] | [TODO: Solution] |
| [Pitfall 2] | [TODO: Description] | [TODO: Solution] |

---

## Cross-References

### Prerequisites

- [TODO: [Go Concurrency](../02-Language-Design/02-Language-Features/03-Goroutines.md)]
- [TODO: [Context Package](../04-Technology-Stack/01-Core-Library/04-Context-Package.md)]

### Related Patterns

| Pattern | Relationship | When to Use |
|---------|--------------|-------------|
| [Pattern A](../03-Engineering-CloudNative/EC-XXX-A.md) | Alternative | [TODO] |
| [Pattern B](../03-Engineering-CloudNative/EC-XXX-B.md) | Complementary | [TODO] |

### Other Dimensions

- **Formal Theory**: [TODO: [FT-XXX](../01-Formal-Theory/FT-XXX-Name.md)]
- **Language Design**: [TODO: [LD-XXX](../02-Language-Design/LD-XXX-Name.md)]
- **Technology**: [TODO: [TS-XXX](../04-Technology-Stack/TS-XXX-Name.md)]
- **Application**: [TODO: [AD-XXX](../05-Application-Domains/AD-XXX-Name.md)]

---

## References

### Books

[1] [TODO: Book Title](https://) - [TODO: Author]

### Articles

[2] [TODO: Article Title](https://) - [TODO: Author/Source]

### Go Libraries

[3] [TODO: Library](https://pkg.go.dev/...) - [TODO: Description]

### Case Studies

[4] [TODO: Company/Project] - [TODO: How they use this pattern]

---

## Document History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | [TODO: YYYY-MM-DD] | Initial pattern documentation | [TODO: Name] |

---

*Template: EC-XXX - Engineering/Cloud-Native Pattern (S/A-Level)*
*For contribution guidelines, see [CONTRIBUTING.md](../CONTRIBUTING.md)*
