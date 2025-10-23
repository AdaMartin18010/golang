package patterns

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestGenerateContextCancellation(t *testing.T) {
	code := GenerateContextCancellation("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "context") {
		t.Error("Missing context import")
	}

	if !strings.Contains(code, "WithCancel") {
		t.Error("Missing WithCancel")
	}

	// 验证语法
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "context.go", code, parser.AllErrors)
	if err != nil {
		t.Fatalf("Syntax error: %v", err)
	}
}

func TestGenerateContextTimeout(t *testing.T) {
	code := GenerateContextTimeout("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "WithTimeout") {
		t.Error("Missing WithTimeout")
	}

	if !strings.Contains(code, "WithDeadline") {
		t.Error("Missing WithDeadline")
	}
}

func TestGenerateContextValue(t *testing.T) {
	code := GenerateContextValue("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "WithValue") {
		t.Error("Missing WithValue")
	}
}

func TestGenerateGracefulShutdown(t *testing.T) {
	code := GenerateGracefulShutdown("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "signal") {
		t.Error("Missing signal handling")
	}
}

func TestGenerateRateLimiting(t *testing.T) {
	code := GenerateRateLimiting("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "RateLimiter") {
		t.Error("Missing RateLimiter")
	}

	if !strings.Contains(code, "TokenBucket") {
		t.Error("Missing TokenBucket")
	}
}
