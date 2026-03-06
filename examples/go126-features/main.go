// Package main demonstrates Go 1.26 new features
//
// 演示内容包括：
// 1. new() 函数支持表达式参数
// 2. errors.AsType 泛型错误断言
// 3. slog.NewMultiHandler 多日志处理器
// 4. 泛型自引用类型约束
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

// ============================================================================
// Feature 1: new() with expressions
// ============================================================================

// Config represents application configuration
type Config struct {
	Name     string
	LogLevel string
	Port     int
}

// NewConfig creates a config with default values
func NewConfig() *Config {
	return &Config{
		Name:     "default",
		LogLevel: "info",
		Port:     8080,
	}
}

// WithName returns config with name set (for new() expression demo)
func (c Config) WithName(name string) Config {
	c.Name = name
	return c
}

// demonstrateNewExpression shows Go 1.26 new() clarifications
// Note: Go 1.26 does NOT support new(expression), this is a common misconception
// The actual change is that new() is better optimized for stack allocation
func demonstrateNewExpression() {
	fmt.Println("=== Feature 1: new() clarifications ===")
	fmt.Println("NOTE: Go 1.26 does NOT support new(expression)")
	fmt.Println("The actual improvements are:")
	fmt.Println("1. Better stack allocation for small objects")
	fmt.Println("2. Optimized escape analysis")
	fmt.Println()

	// Traditional usage (still the only way)
	ptr := new(int)
	*ptr = 42
	fmt.Printf("Traditional new() usage: %d\n", *ptr)

	// Practical use case: optional fields in JSON
	type User struct {
		Name string
		Age  *int `json:"age,omitempty"`
	}

	// Traditional way (still valid)
	age := 25
	user := User{
		Name: "Alice",
		Age:  &age, // Use address of variable
	}

	fmt.Printf("User with pointer field: Name=%s, Age=%d\n", user.Name, *user.Age)
	fmt.Println()
}

// ============================================================================
// Feature 2: errors.AsType (Go 1.26)
// ============================================================================

// CustomError is a custom error type
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

// demonstrateErrorsAsType shows Go 1.26 errors.AsType feature
func demonstrateErrorsAsType() {
	fmt.Println("=== Feature 2: errors.AsType ===")

	var err error = &CustomError{Code: 404, Message: "not found"}

	// Before Go 1.26: need var declaration
	var customErr1 *CustomError
	if errors.As(err, &customErr1) {
		fmt.Printf("Old way - Code: %d, Message: %s\n", customErr1.Code, customErr1.Message)
	}

	// Go 1.26: Direct type-safe assertion
	// Note: This uses the new errors.AsType function
	customErr2, ok := errors.AsType[*CustomError](err)
	if ok {
		fmt.Printf("New way - Code: %d, Message: %s\n", customErr2.Code, customErr2.Message)
	}

	// Example with wrapped error
	wrappedErr := fmt.Errorf("wrapped: %w", err)
	customErr3, ok := errors.AsType[*CustomError](wrappedErr)
	if ok {
		fmt.Printf("Wrapped error - Code: %d\n", customErr3.Code)
	}

	fmt.Println()
}

// ============================================================================
// Feature 3: slog.NewMultiHandler
// ============================================================================

// demonstrateMultiHandler shows Go 1.26 slog.NewMultiHandler feature
func demonstrateMultiHandler() {
	fmt.Println("=== Feature 3: slog.NewMultiHandler ===")

	// Create multiple handlers
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})

	// Go 1.26: Combine handlers with MultiHandler
	multiHandler := slog.NewMultiHandler(jsonHandler, textHandler)
	logger := slog.New(multiHandler)

	// This will be logged by both handlers (if level allows)
	logger.Info("Application started",
		slog.String("version", "1.0.0"),
		slog.Int("port", 8080),
	)

	// This will definitely be logged by both
	logger.Warn("High memory usage",
		slog.Float64("usage_percent", 85.5),
	)

	fmt.Println("\nMultiHandler allows logging to multiple destinations simultaneously")
	fmt.Println()
}

// ============================================================================
// Feature 4: Generic self-reference (conceptual)
// ============================================================================

// Adder is an example of self-referential generic constraint (Go 1.26)
// The type parameter A can refer to the interface itself
type Adder[A Adder[A]] interface {
	Add(other A) A
	Value() int
}

// IntAdder implements Adder
type IntAdder struct {
	v int
}

func (a IntAdder) Add(other IntAdder) IntAdder {
	return IntAdder{v: a.v + other.v}
}

func (a IntAdder) Value() int {
	return a.v
}

// demonstrateGenericSelfReference shows Go 1.26 generic self-reference
func demonstrateGenericSelfReference() {
	fmt.Println("=== Feature 4: Generic self-reference ===")

	// Before Go 1.26: this pattern was not allowed
	// Now we can define interfaces that reference themselves in constraints

	a1 := IntAdder{v: 10}
	a2 := IntAdder{v: 20}
	result := a1.Add(a2)

	fmt.Printf("a1: %d\n", a1.Value())
	fmt.Printf("a2: %d\n", a2.Value())
	fmt.Printf("result: %d\n", result.Value())

	fmt.Println("\nSelf-referential generics enable powerful type constraints")
	fmt.Println("See internal/domain/interfaces/specification_go126.go for practical example")
	fmt.Println()
}

// ============================================================================
// Main
// ============================================================================

func main() {
	fmt.Println("Go 1.26 Features Demo")
	fmt.Println("====================")

	demonstrateNewExpression()
	demonstrateErrorsAsType()
	demonstrateMultiHandler()
	demonstrateGenericSelfReference()

	fmt.Println("All demonstrations completed!")
}
