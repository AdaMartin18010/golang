package patterns

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestGenerateWorkerPool(t *testing.T) {
	data := map[string]interface{}{
		"PackageName": "main",
		"NumWorkers":  5,
		"BufferSize":  10,
	}

	code := GenerateWorkerPool(data)

	// 验证代码非空
	if code == "" {
		t.Fatal("Generated code is empty")
	}

	// 验证包含关键字
	if !strings.Contains(code, "package main") {
		t.Error("Missing package declaration")
	}

	if !strings.Contains(code, "WorkerPool") {
		t.Error("Missing WorkerPool struct")
	}

	// 验证代码可以解析
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "generated.go", code, parser.AllErrors)
	if err != nil {
		t.Fatalf("Generated code has syntax errors: %v", err)
	}
}

func TestGenerateFanIn(t *testing.T) {
	data := map[string]interface{}{
		"PackageName": "test",
	}

	code := GenerateFanIn(data)

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "FanIn") {
		t.Error("Missing FanIn function")
	}

	// 验证语法
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "fanin.go", code, parser.AllErrors)
	if err != nil {
		t.Fatalf("Syntax error: %v", err)
	}
}

func TestGenerateFanOut(t *testing.T) {
	data := map[string]interface{}{
		"PackageName": "test",
	}

	code := GenerateFanOut(data)

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "FanOut") {
		t.Error("Missing FanOut function")
	}
}

func TestGeneratePipeline(t *testing.T) {
	data := map[string]interface{}{
		"PackageName": "test",
		"Stages":      3,
	}

	code := GeneratePipeline(data)

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "Pipeline") {
		t.Error("Missing Pipeline function")
	}
}

func TestGenerateGenerator(t *testing.T) {
	data := map[string]interface{}{
		"PackageName": "test",
	}

	code := GenerateGenerator(data)

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "Generator") {
		t.Error("Missing Generator function")
	}
}

// Benchmark tests
func BenchmarkGenerateWorkerPool(b *testing.B) {
	data := map[string]interface{}{
		"PackageName": "main",
		"NumWorkers":  5,
		"BufferSize":  10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateWorkerPool(data)
	}
}
