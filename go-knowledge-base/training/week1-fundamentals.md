# Week 1: Go Fundamentals and Tooling

## Module Overview

**Duration:** 40 hours (5 days)  
**Prerequisites:** Basic programming experience in any language  
**Learning Goal:** Write idiomatic Go code with proper testing and tooling

---

## Learning Objectives

By the end of this week, you will be able to:

1. **Language Fundamentals**
   - Explain Go's design philosophy (simplicity, explicitness, composition)
   - Use Go's type system effectively (structs, interfaces, type assertions)
   - Write clean, idiomatic Go code following team standards

2. **Memory and Runtime**
   - Understand Go's memory model and garbage collection
   - Explain the difference between value and pointer receivers
   - Use escape analysis to optimize memory allocation

3. **Error Handling**
   - Implement proper error handling patterns
   - Create custom error types with context
   - Use error wrapping and unwrapping effectively

4. **Testing Fundamentals**
   - Write comprehensive unit tests
   - Implement table-driven tests
   - Create benchmarks for performance measurement
   - Use fuzzing for input validation

5. **Development Tooling**
   - Configure a professional development environment
   - Use Go modules for dependency management
   - Debug Go applications with Delve
   - Profile applications for performance optimization

---

## Reading Assignments

### Required Reading (Complete by Day 3)

1. **[Go Design Philosophy](../02-Language-Design/01-Design-Philosophy/README.md)**
   - Read: Simplicity, Composition, Explicitness, Orthogonality
   - Focus on: How these principles influence code structure
   - Reflection: Write 3 examples of explicit vs implicit code

2. **[Go Type System](../02-Language-Design/02-Language-Features/01-Type-System.md)**
   - Study: Structural typing, interface satisfaction
   - Practice: Define interfaces based on behavior, not types
   - Key concept: "Accept interfaces, return structs"

3. **[Go Interfaces](../02-Language-Design/02-Language-Features/02-Interfaces.md)**
   - Master: Implicit interface implementation
   - Understand: Empty interface (any) usage and pitfalls
   - Learn: Interface composition patterns

4. **[Error Handling Patterns](../02-Language-Design/LD-023-Go-Error-Handling-Patterns.md)**
   - Learn: Sentinel errors vs custom error types
   - Practice: Error wrapping with fmt.Errorf and %w
   - Understand: When to panic vs return error

5. **[Testing Patterns](../02-Language-Design/LD-009-Go-Testing-Patterns.md)**
   - Master: Table-driven tests
   - Learn: Test organization and naming conventions
   - Understand: Mocking strategies in Go

### Supplementary Reading (Complete by Day 5)

6. **[Go Memory Management](../02-Language-Design/02-Language-Features/09-Memory-Management.md)**
   - Understand: Stack vs heap allocation
   - Learn: Escape analysis basics
   - Study: Memory layout of structs

7. **[Go Garbage Collection](../02-Language-Design/02-Language-Features/10-GC.md)**
   - Learn: Tri-color mark and sweep algorithm
   - Understand: GC tuning parameters
   - Study: Memory leak patterns in Go

8. **[Go Modules](../04-Technology-Stack/04-Development-Tools/01-Go-Modules.md)**
   - Master: Module initialization and versioning
   - Learn: Semantic import versioning
   - Understand: Private module proxy usage

---

## Hands-on Exercises

### Day 1: Environment Setup and Basics

#### Exercise 1.1: Development Environment (2 hours)

Set up your development environment:

```bash
# 1. Install Go 1.22+
# Verify installation
go version

# 2. Configure workspace
mkdir -p ~/go/src/github.com/yourusername
cd ~/go/src/github.com/yourusername

# 3. Initialize your first module
go mod init onboarding-week1

# 4. Install essential tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/go-delve/delve/cmd/dlv@latest
```

**Deliverable:** Screenshot of successful `go version` and `go env` output

#### Exercise 1.2: Hello Universe (1 hour)

Create a command-line greeting program:

```go
package main

import (
    "flag"
    "fmt"
    "strings"
    "time"
)

type Greeting struct {
    Name      string
    Message   string
    Language  string
    Timestamp time.Time
}

func (g Greeting) String() string {
    return fmt.Sprintf("[%s] %s: %s %s",
        g.Timestamp.Format(time.RFC3339),
        g.Language,
        g.Message,
        g.Name)
}

func main() {
    var name = flag.String("name", "World", "Name to greet")
    var lang = flag.String("lang", "en", "Language code (en, es, fr, de)")
    var repeat = flag.Int("repeat", 1, "Number of times to repeat")
    
    flag.Parse()
    
    messages := map[string]string{
        "en": "Hello",
        "es": "Hola",
        "fr": "Bonjour",
        "de": "Hallo",
    }
    
    msg, ok := messages[*lang]
    if !ok {
        msg = messages["en"]
    }
    
    g := Greeting{
        Name:      *name,
        Message:   msg,
        Language:  *lang,
        Timestamp: time.Now(),
    }
    
    fmt.Println(strings.Repeat(g.String()+"\n", *repeat))
}
```

**Extensions:**
- Add more languages
- Support custom greeting templates
- Add JSON output format
- Add file logging

**Deliverable:** Working program with at least 3 extensions

#### Exercise 1.3: Data Types Deep Dive (2 hours)

Explore Go's type system:

```go
package types

// Exercise 1.3.1: Create a Person type with embedded types
type Address struct {
    Street  string
    City    string
    Country string
    ZIP     string
}

type ContactInfo struct {
    Email   string
    Phone   string
    Address Address
}

type Person struct {
    ID        string
    FirstName string
    LastName  string
    Age       int
    ContactInfo  // Embedded struct
}

// Exercise 1.3.2: Implement Stringer interface
func (p Person) String() string {
    return fmt.Sprintf("Person[%s]: %s %s, %d years old",
        p.ID, p.FirstName, p.LastName, p.Age)
}

// Exercise 1.3.3: Create a validation interface
type Validator interface {
    Validate() error
}

func (p Person) Validate() error {
    if p.ID == "" {
        return fmt.Errorf("person ID cannot be empty")
    }
    if p.FirstName == "" || p.LastName == "" {
        return fmt.Errorf("first and last name are required")
    }
    if p.Age < 0 || p.Age > 150 {
        return fmt.Errorf("invalid age: %d", p.Age)
    }
    return nil
}
```

**Tasks:**
1. Create different struct types (Product, Order, Company)
2. Implement common interfaces (Stringer, Validator)
3. Use struct tags for JSON marshaling
4. Practice type conversions and assertions

**Deliverable:** Package with at least 5 types and complete test coverage

---

### Day 2: Interfaces and Methods

#### Exercise 2.1: Interface Design (3 hours)

Design interfaces for a notification system:

```go
package notification

import "context"

// Notifier is the core interface
type Notifier interface {
    Send(ctx context.Context, msg Message) error
    Close() error
}

// Message represents a notification
type Message struct {
    ID      string
    To      string
    Subject string
    Body    string
    Metadata map[string]string
}

// PriorityNotifier adds priority capabilities
type PriorityNotifier interface {
    Notifier
    SendPriority(ctx context.Context, msg Message, priority Priority) error
}

type Priority int

const (
    PriorityLow Priority = iota
    PriorityNormal
    PriorityHigh
    PriorityCritical
)

// Implementations to create:
// - EmailNotifier
// - SMSNotifier
// - SlackNotifier
// - PushNotifier
```

**Requirements:**
1. Each implementation must have its own configuration
2. Support for retry logic with exponential backoff
3. Metrics collection (success/failure counts)
4. Context cancellation support

**Deliverable:** Working notification package with 4 implementations

#### Exercise 2.2: Interface Composition (2 hours)

Build composable interfaces for a storage system:

```go
package storage

// Basic operations
type Reader interface {
    Read(ctx context.Context, key string) ([]byte, error)
}

type Writer interface {
    Write(ctx context.Context, key string, value []byte) error
}

type Deleter interface {
    Delete(ctx context.Context, key string) error
}

type Lister interface {
    List(ctx context.Context, prefix string) ([]string, error)
}

// Composed interfaces
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteDeleter interface {
    Reader
    Writer
    Deleter
}

type FullStorage interface {
    Reader
    Writer
    Deleter
    Lister
    Closer
}

type Closer interface {
    Close() error
}
```

**Implement:**
1. MemoryStorage (in-memory map-based)
2. FileStorage (filesystem-based)
3. CachedStorage (decorator with caching)
4. MeteredStorage (decorator with metrics)

**Deliverable:** Storage package with decorator pattern implementations

---

### Day 3: Error Handling Mastery

#### Exercise 3.1: Error Types Hierarchy (3 hours)

Create a comprehensive error handling system:

```go
package errors

import (
    "errors"
    "fmt"
)

// Domain errors
type ErrorCode string

const (
    ErrCodeNotFound     ErrorCode = "NOT_FOUND"
    ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
    ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
    ErrCodeForbidden    ErrorCode = "FORBIDDEN"
    ErrCodeConflict     ErrorCode = "CONFLICT"
    ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
)

// ApplicationError is the base error type
type ApplicationError struct {
    Code       ErrorCode
    Message    string
    Operation  string
    Resource   string
    Cause      error
    Context    map[string]interface{}
}

func (e *ApplicationError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("[%s] %s: %s (caused by: %v)",
            e.Code, e.Operation, e.Message, e.Cause)
    }
    return fmt.Sprintf("[%s] %s: %s", e.Code, e.Operation, e.Message)
}

func (e *ApplicationError) Unwrap() error {
    return e.Cause
}

// Helper constructors
func NewNotFound(resource, id string) error {
    return &ApplicationError{
        Code:     ErrCodeNotFound,
        Message:  fmt.Sprintf("%s with id '%s' not found", resource, id),
        Resource: resource,
    }
}

func NewInvalidInput(field, reason string) error {
    return &ApplicationError{
        Code:    ErrCodeInvalidInput,
        Message: fmt.Sprintf("invalid %s: %s", field, reason),
        Context: map[string]interface{}{"field": field},
    }
}

func Wrap(cause error, code ErrorCode, operation string) error {
    return &ApplicationError{
        Code:      code,
        Message:   cause.Error(),
        Operation: operation,
        Cause:     cause,
    }
}
```

**Tasks:**
1. Add stack trace capture
2. Implement error filtering by code
3. Create error logging formatter
4. Add retryable error detection

**Deliverable:** Production-ready error package with tests

#### Exercise 3.2: Error Handling Patterns (2 hours)

Implement common error handling patterns:

```go
package patterns

// Pattern 1: Error accumulation for batch operations
type BatchResult struct {
    Successes int
    Failures  int
    Errors    []error
}

func (br *BatchResult) Error() error {
    if len(br.Errors) == 0 {
        return nil
    }
    return fmt.Errorf("batch operation completed with %d failures: %v",
        len(br.Errors), br.Errors)
}

// Pattern 2: Sentinel errors
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidToken    = errors.New("invalid authentication token")
    ErrRateLimited     = errors.New("rate limit exceeded")
    ErrServiceNotReady = errors.New("service not ready")
)

// Pattern 3: Error with retry info
type RetryableError struct {
    Err        error
    RetryAfter time.Duration
    MaxRetries int
}

func (e *RetryableError) Error() string {
    return fmt.Sprintf("retryable error (retry after %v): %v",
        e.RetryAfter, e.Err)
}

func (e *RetryableError) Unwrap() error {
    return e.Err
}
```

**Deliverable:** Document with 5+ error patterns and examples

---

### Day 4: Testing Excellence

#### Exercise 4.1: Table-Driven Tests (3 hours)

Master table-driven test patterns:

```go
package calculator

import (
    "errors"
    "testing"
)

// Calculator operations
type Calculator struct {
    precision int
}

func NewCalculator(precision int) *Calculator {
    return &Calculator{precision: precision}
}

func (c *Calculator) Add(a, b float64) float64 {
    return c.round(a + b)
}

func (c *Calculator) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return c.round(a / b), nil
}

func (c *Calculator) round(n float64) float64 {
    factor := math.Pow(10, float64(c.precision))
    return math.Round(n*factor) / factor
}

// Table-driven tests
func TestCalculator_Add(t *testing.T) {
    tests := []struct {
        name      string
        precision int
        a, b      float64
        want      float64
    }{
        {
            name:      "add positive numbers",
            precision: 2,
            a:         1.5,
            b:         2.5,
            want:      4.0,
        },
        {
            name:      "add negative numbers",
            precision: 2,
            a:         -1.5,
            b:         -2.5,
            want:      -4.0,
        },
        {
            name:      "rounding with precision",
            precision: 2,
            a:         1.005,
            b:         2.005,
            want:      3.01,
        },
        {
            name:      "zero addition",
            precision: 2,
            a:         0,
            b:         5.5,
            want:      5.5,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := NewCalculator(tt.precision)
            got := c.Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%v, %v) = %v, want %v",
                    tt.a, tt.b, got, tt.want)
            }
        })
    }
}

func TestCalculator_Divide(t *testing.T) {
    tests := []struct {
        name    string
        a, b    float64
        want    float64
        wantErr bool
    }{
        {
            name: "divide positive numbers",
            a:    10,
            b:    2,
            want: 5,
        },
        {
            name:    "divide by zero",
            a:       10,
            b:       0,
            wantErr: true,
        },
        {
            name: "divide negative numbers",
            a:    -10,
            b:    -2,
            want: 5,
        },
        {
            name: "divide with decimal result",
            a:    10,
            b:    3,
            want: 3.33,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := NewCalculator(2)
            got, err := c.Divide(tt.a, tt.b)
            if (err != nil) != tt.wantErr {
                t.Errorf("Divide() error = %v, wantErr %v",
                    err, tt.wantErr)
                return
            }
            if !tt.wantErr && got != tt.want {
                t.Errorf("Divide() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Tasks:**
1. Create subtests for organization
2. Add test helpers and fixtures
3. Implement golden file testing
4. Add parallel test execution

**Deliverable:** Calculator package with comprehensive tests (>90% coverage)

#### Exercise 4.2: Benchmarks and Fuzzing (2 hours)

Learn performance testing:

```go
package calculator

import (
    "testing"
    "time"
)

// Benchmarks
func BenchmarkCalculator_Add(b *testing.B) {
    c := NewCalculator(2)
    for i := 0; i < b.N; i++ {
        c.Add(1.5, 2.5)
    }
}

func BenchmarkCalculator_AddParallel(b *testing.B) {
    c := NewCalculator(2)
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            c.Add(1.5, 2.5)
        }
    })
}

// Comparison benchmark
func BenchmarkStringConcatenation(b *testing.B) {
    b.Run("Plus", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = "Hello" + " " + "World"
        }
    })
    
    b.Run("Sprintf", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = fmt.Sprintf("%s %s", "Hello", "World")
        }
    })
    
    b.Run("Builder", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var sb strings.Builder
            sb.WriteString("Hello")
            sb.WriteString(" ")
            sb.WriteString("World")
            _ = sb.String()
        }
    })
}

// Fuzzing
func FuzzCalculator_Add(f *testing.F) {
    f.Add(1.5, 2.5)
    f.Add(-1.0, 1.0)
    f.Add(0.0, 0.0)
    
    f.Fuzz(func(t *testing.T, a, b float64) {
        c := NewCalculator(10)
        result := c.Add(a, b)
        
        // Properties to verify
        // Commutative: a + b == b + a
        result2 := c.Add(b, a)
        if result != result2 {
            t.Errorf("Add not commutative: %v + %v = %v, but %v + %v = %v",
                a, b, result, b, a, result2)
        }
        
        // Identity: a + 0 == a
        if b == 0 && result != a {
            t.Errorf("Add identity failed: %v + 0 = %v, want %v",
                a, result, a)
        }
    })
}
```

**Deliverable:** Benchmark results and fuzzing corpus

---

### Day 5: Project Structure and Modules

#### Exercise 5.1: Project Layout (3 hours)

Create a well-structured Go project:

```
myproject/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── cli/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── user.go
│   │   └── order.go
│   ├── repository/
│   │   ├── user_repo.go
│   │   └── order_repo.go
│   ├── service/
│   │   ├── user_service.go
│   │   └── order_service.go
│   └── infrastructure/
│       ├── database/
│       └── http/
├── pkg/
│   ├── validator/
│   └── logger/
├── api/
│   └── openapi.yaml
├── configs/
│   └── config.yaml
├── scripts/
│   └── migrate.sh
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

**Requirements:**
1. Follow standard Go project layout
2. Use internal for application-specific code
3. Create reusable packages in pkg/
4. Implement clean architecture separation

**Deliverable:** Complete project structure with working code

#### Exercise 5.2: Module Management (2 hours)

Practice Go module operations:

```bash
# 1. Create a module
go mod init github.com/yourusername/myproject

# 2. Add dependencies
go get github.com/gin-gonic/gin
go get github.com/stretchr/testify

# 3. Update dependencies
go get -u ./...
go mod tidy

# 4. Vendor dependencies
go mod vendor

# 5. Verify dependencies
go mod verify

# 6. View dependency graph
go mod graph

# 7. Clean unused dependencies
go mod tidy

# 8. Work with multiple modules
go work init
go work use ./module1
go work use ./module2
```

**Deliverable:** Multiple-module workspace demonstration

---

## Code Review Checklist

### Style and Idioms

- [ ] Code follows `gofmt` formatting
- [ ] `go vet` passes without warnings
- [ ] `golangci-lint` passes
- [ ] Variable names are clear and idiomatic
- [ ] Error messages are lowercase (no punctuation at end)
- [ ] Comments are complete sentences
- [ ] Exported symbols have documentation comments

### Error Handling

- [ ] All errors are checked
- [ ] Errors provide context with wrapping
- [ ] Sentinel errors are used appropriately
- [ ] Panic is only used for unrecoverable errors
- [ ] Error types implement `Unwrap()` when wrapping

### Interfaces

- [ ] Interfaces are defined where they are used (consumer side)
- [ ] Interfaces are small and focused
- [ ] Empty interfaces (`any`) are used sparingly
- [ ] Interface naming follows Go conventions (`-er` suffix)

### Concurrency (Basic)

- [ ] Shared state is avoided where possible
- [ ] Mutex fields are not copied
- [ ] Channels are closed by sender, not receiver

### Testing

- [ ] All new code has unit tests
- [ ] Tests use table-driven patterns where appropriate
- [ ] Test coverage is at least 80%
- [ ] Benchmarks exist for performance-critical code
- [ ] Tests clean up resources properly

### Performance

- [ ] Preallocate slices with known capacity
- [ ] Use `strings.Builder` for string concatenation in loops
- [ ] Avoid unnecessary allocations
- [ ] Use value receivers when no mutation needed

---

## Assessment Criteria

### Knowledge Assessment (40%)

**Quiz Topics:**
1. Go type system (structural typing, interfaces)
2. Memory management (stack vs heap, escape analysis)
3. Error handling patterns
4. Testing methodology
5. Module system

**Passing Score:** 85%

### Coding Challenge (40%)

**Problem:** Build a URL shortener service

**Requirements:**
- HTTP API with CRUD operations
- In-memory storage (map-based)
- Custom error types
- Comprehensive tests (>85% coverage)
- Benchmarks for key operations
- Race-condition free (verified with -race)

**Evaluation Criteria:**
- Correctness: 30%
- Code quality: 25%
- Test coverage: 20%
- Performance: 15%
- Documentation: 10%

### Code Review Participation (20%)

**Requirements:**
- Review at least 3 peer submissions
- Provide constructive feedback
- Identify potential issues
- Suggest improvements

---

## Resources and References

### Official Documentation
- [Go Language Specification](https://golang.org/ref/spec)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go FAQ](https://golang.org/doc/faq)

### Recommended Tools
- **IDE:** VS Code with Go extension or GoLand
- **Linting:** golangci-lint
- **Formatting:** gofmt, goimports
- **Debugging:** Delve
- **Testing:** go test, gotestsum

### Books
- "The Go Programming Language" (Donovan & Kernighan)
- "Learning Go" (Jon Bodner)
- "100 Go Mistakes and How to Avoid Them" (Teiva Harsanyi)

---

## Week 1 Completion Checklist

- [ ] Environment set up and verified
- [ ] All exercises completed
- [ ] Code review checklist understood
- [ ] Week 1 assessment passed (85%+)
- [ ] Capstone project submitted
- [ ] Peer reviews completed

---

*Next: [Week 2: Concurrency Deep Dive](week2-concurrency.md)*
