// Package main demonstrates Go 1.26 new features
//
// 演示内容包括：
// 1. new() 函数支持表达式参数 (Go 1.26)
// 2. errors.AsType 泛型错误断言 (Go 1.26)
// 3. slog.NewMultiHandler 多日志处理器 (Go 1.26)
// 4. 泛型自引用类型约束 (Go 1.26)
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// ============================================================================
// Feature 1: new() with expressions (Go 1.26)
// ============================================================================

// yearsSince calculates approximate years since a given time
func yearsSince(t time.Time) int {
	return int(time.Since(t).Hours() / (365.25 * 24)) // approximately
}

// demonstrateNewExpression shows Go 1.26 new() with expression operand
// Go 1.26 语言变更：内置 new 函数允许操作数为表达式，指定变量的初始值
func demonstrateNewExpression() {
	fmt.Println("=== Feature 1: new() with expressions ===")
	fmt.Println("Go 1.26 Change: The built-in new function now allows its operand to be an expression.")
	fmt.Println()

	// Official Go 1.26 example: new(expression)
	// 特别适用于序列化包中表示可选值的指针字段
	type Person struct {
		Name string `json:"name"`
		Age  *int   `json:"age"` // age if known; nil otherwise
	}

	born := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	person := Person{
		Name: "Alice",
		Age:  new(yearsSince(born)), // Go 1.26: new accepts expression
	}
	fmt.Printf("Person: Name=%s, Age=%d\n", person.Name, *person.Age)

	// 更简洁的用法：不再需要临时变量
	enabled := new(true)    // *bool pointing to true
	count := new(42)        // *int pointing to 42
	message := new("hello") // *string pointing to "hello"

	fmt.Printf("new(true) = %t\n", *enabled)
	fmt.Printf("new(42) = %d\n", *count)
	fmt.Printf("new(\"hello\") = %s\n", *message)
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

	// Before Go 1.26: need var declaration and pointer-to-pointer
	if customErr1, ok := errors.AsType[*CustomError](err); ok {
		fmt.Printf("Old way - Code: %d, Message: %s\n", customErr1.Code, customErr1.Message)
	}

	// Go 1.26: Direct type-safe assertion with generics
	// AsType is type-safe at compile time, no reflection overhead
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
// Feature 3: slog.NewMultiHandler (Go 1.26)
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
// Feature 4: Generic self-reference (Go 1.26)
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
