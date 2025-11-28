package patterns

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestGenerateActorModel(t *testing.T) {
	code := GenerateActorModel("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "Actor") {
		t.Error("Missing Actor")
	}

	if !strings.Contains(code, "Message") {
		t.Error("Missing Message")
	}
}

func TestGenerateFuturePromise(t *testing.T) {
	code := GenerateFuturePromise("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "Future") {
		t.Error("Missing Future")
	}
}

func TestGenerateMapReduce(t *testing.T) {
	code := GenerateMapReduce("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "MapReduce") {
		t.Error("Missing MapReduce")
	}
}

func TestGeneratePubSub(t *testing.T) {
	code := GeneratePubSub("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "PubSub") {
		t.Error("Missing PubSub")
	}

	if !strings.Contains(code, "Subscribe") {
		t.Error("Missing Subscribe")
	}

	if !strings.Contains(code, "Publish") {
		t.Error("Missing Publish")
	}
}

func TestGenerateSessionTypes(t *testing.T) {
	code := GenerateSessionTypes("main")

	if code == "" {
		t.Fatal("Generated code is empty")
	}

	if !strings.Contains(code, "Session") {
		t.Error("Missing Session")
	}
}

func TestAllAdvancedPatterns(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
	}{
		{"Actor", GenerateActorModel},
		{"FuturePromise", GenerateFuturePromise},
		{"MapReduce", GenerateMapReduce},
		{"PubSub", GeneratePubSub},
		{"SessionTypes", GenerateSessionTypes},
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
