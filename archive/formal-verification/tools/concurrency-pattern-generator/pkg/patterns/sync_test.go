package patterns

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestGenerateMutexSimple(t *testing.T) {
	code := GenerateMutexSimple("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "SafeCounter") {
		t.Error("Missing SafeCounter")
	}

	if !strings.Contains(code, "sync.Mutex") {
		t.Error("Missing sync.Mutex")
	}

	// 验证语法
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "mutex.go", code, parser.AllErrors)
	if err != nil {
		t.Fatalf("Syntax error: %v", err)
	}
}

func TestGenerateRWMutexSimple(t *testing.T) {
	code := GenerateRWMutexSimple("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "Cache") {
		t.Error("Missing Cache")
	}

	if !strings.Contains(code, "sync.RWMutex") {
		t.Error("Missing sync.RWMutex")
	}
}

func TestGenerateWaitGroupSimple(t *testing.T) {
	code := GenerateWaitGroupSimple("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "ParallelTasks") {
		t.Error("Missing ParallelTasks")
	}

	if !strings.Contains(code, "sync.WaitGroup") {
		t.Error("Missing sync.WaitGroup")
	}
}

func TestGenerateOnceSimple(t *testing.T) {
	code := GenerateOnceSimple("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "Singleton") {
		t.Error("Missing Singleton")
	}

	if !strings.Contains(code, "sync.Once") {
		t.Error("Missing sync.Once")
	}
}

func TestGenerateSemaphoreSimple(t *testing.T) {
	code := GenerateSemaphoreSimple("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "Semaphore") {
		t.Error("Missing Semaphore")
	}

	if !strings.Contains(code, "Acquire") {
		t.Error("Missing Acquire method")
	}

	if !strings.Contains(code, "Release") {
		t.Error("Missing Release method")
	}
}

func TestAllSyncPatterns(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
	}{
		{"Mutex", GenerateMutexSimple},
		{"RWMutex", GenerateRWMutexSimple},
		{"WaitGroup", GenerateWaitGroupSimple},
		{"Once", GenerateOnceSimple},
		{"Semaphore", GenerateSemaphoreSimple},
		{"Barrier", GenerateBarrierSimple},
		{"CountDownLatch", GenerateCountDownLatchSimple},
		{"Cond", GenerateCondSimple},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := tt.fn("test")
			if code == "" {
				t.Errorf("%s generated empty code", tt.name)
			}

			// 验证语法
			fset := token.NewFileSet()
			_, err := parser.ParseFile(fset, tt.name+".go", code, parser.AllErrors)
			if err != nil {
				t.Errorf("%s has syntax errors: %v", tt.name, err)
			}
		})
	}
}
