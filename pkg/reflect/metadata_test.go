package reflect

import (
	"testing"
)

func TestInspector_InspectType(t *testing.T) {
	inspector := NewInspector()

	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	user := User{}
	metadata := inspector.InspectType(user)

	if metadata.Name != "User" {
		t.Errorf("Expected name User, got %s", metadata.Name)
	}

	if metadata.Kind != "struct" {
		t.Errorf("Expected kind struct, got %s", metadata.Kind)
	}

	if len(metadata.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(metadata.Fields))
	}
}

func TestInspector_InspectStruct(t *testing.T) {
	inspector := NewInspector()

	type User struct {
		ID    int    `json:"id" db:"id"`
		Name  string `json:"name" db:"name"`
		Email string `json:"email" db:"email"`
	}

	user := User{}
	metadata := inspector.InspectStruct(user)

	if metadata.Name != "User" {
		t.Errorf("Expected name User, got %s", metadata.Name)
	}

	if len(metadata.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(metadata.Fields))
	}

	// 检查标签
	if metadata.Tags["ID"]["json"] != "id" {
		t.Errorf("Expected json tag 'id' for ID field")
	}
}

func TestInspector_InspectFunction(t *testing.T) {
	inspector := NewInspector()

	testFunc := func(a, b int) int {
		return a + b
	}

	metadata := inspector.InspectFunction(testFunc)

	if metadata.Name == "" {
		t.Error("Expected non-empty function name")
	}

	if len(metadata.Inputs) != 2 {
		t.Errorf("Expected 2 inputs, got %d", len(metadata.Inputs))
	}

	if len(metadata.Outputs) != 1 {
		t.Errorf("Expected 1 output, got %d", len(metadata.Outputs))
	}
}

func TestInspector_Describe(t *testing.T) {
	inspector := NewInspector()

	type User struct {
		ID   int
		Name string
	}

	user := User{}
	description := inspector.Describe(user)

	if description == "" {
		t.Error("Expected non-empty description")
	}

	if len(description) < 10 {
		t.Error("Expected detailed description")
	}
}

