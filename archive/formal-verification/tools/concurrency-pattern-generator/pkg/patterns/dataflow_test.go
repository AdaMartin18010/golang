package patterns

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestGenerateProducerConsumer(t *testing.T) {
	code := GenerateProducerConsumer("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "Producer") {
		t.Error("Missing Producer")
	}

	if !strings.Contains(code, "Consumer") {
		t.Error("Missing Consumer")
	}
}

func TestGenerateBufferedChannel(t *testing.T) {
	code := GenerateBufferedChannel("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "make(chan") {
		t.Error("Missing channel creation")
	}
}

func TestGenerateSelectPattern(t *testing.T) {
	code := GenerateSelectPattern("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "select") {
		t.Error("Missing select statement")
	}
}

func TestAllDataFlowPatterns(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
	}{
		{"ProducerConsumer", GenerateProducerConsumer},
		{"BufferedChannel", GenerateBufferedChannel},
		{"UnbufferedChannel", GenerateUnbufferedChannel},
		{"Select", GenerateSelectPattern},
		{"ForSelectLoop", GenerateForSelectLoop},
		{"DoneChannel", GenerateDoneChannel},
		{"ErrorChannel", GenerateErrorChannel},
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
