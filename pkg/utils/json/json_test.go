package json

import (
	"strings"
	"testing"
)

func TestMarshalString(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"age":  30,
	}

	result, err := MarshalString(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "test") {
		t.Error("Expected result to contain 'test'")
	}
}

func TestUnmarshalString(t *testing.T) {
	jsonStr := `{"name":"test","age":30}`
	var result map[string]interface{}

	err := UnmarshalString(jsonStr, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected 'test', got %v", result["name"])
	}
}

func TestPrettyPrint(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"age":  30,
	}

	result, err := PrettyPrint(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "\n") {
		t.Error("Expected formatted JSON with newlines")
	}
}

func TestIsValidJSON(t *testing.T) {
	valid := `{"name":"test"}`
	if !IsValidJSON(valid) {
		t.Error("Expected valid JSON")
	}

	invalid := `{name:test}`
	if IsValidJSON(invalid) {
		t.Error("Expected invalid JSON")
	}
}

func TestGet(t *testing.T) {
	data := []byte(`{"user":{"name":"test","age":30}}`)

	value, err := Get(data, "user.name")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if value != "test" {
		t.Errorf("Expected 'test', got %v", value)
	}
}

func TestSet(t *testing.T) {
	data := []byte(`{"user":{"name":"old"}}`)

	result, err := Set(data, "user.name", "new")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var obj map[string]interface{}
	Unmarshal(result, &obj)

	user := obj["user"].(map[string]interface{})
	if user["name"] != "new" {
		t.Errorf("Expected 'new', got %v", user["name"])
	}
}

func TestMerge(t *testing.T) {
	json1 := []byte(`{"a":1,"b":2}`)
	json2 := []byte(`{"b":3,"c":4}`)

	result, err := Merge(json1, json2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var obj map[string]interface{}
	Unmarshal(result, &obj)

	if obj["a"] != float64(1) {
		t.Errorf("Expected 1, got %v", obj["a"])
	}
	if obj["b"] != float64(3) {
		t.Errorf("Expected 3, got %v", obj["b"])
	}
	if obj["c"] != float64(4) {
		t.Errorf("Expected 4, got %v", obj["c"])
	}
}

func TestFlatten(t *testing.T) {
	data := []byte(`{"user":{"name":"test","age":30}}`)

	result, err := Flatten(data, ".")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var obj map[string]interface{}
	Unmarshal(result, &obj)

	if obj["user.name"] != "test" {
		t.Errorf("Expected 'test', got %v", obj["user.name"])
	}
}

func TestUnflatten(t *testing.T) {
	data := []byte(`{"user.name":"test","user.age":30}`)

	result, err := Unflatten(data, ".")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var obj map[string]interface{}
	Unmarshal(result, &obj)

	user := obj["user"].(map[string]interface{})
	if user["name"] != "test" {
		t.Errorf("Expected 'test', got %v", user["name"])
	}
}

func TestFilter(t *testing.T) {
	data := []byte(`{"a":1,"b":2,"c":3}`)

	result, err := Filter(data, func(k string, v interface{}) bool {
		return k != "b"
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var obj map[string]interface{}
	Unmarshal(result, &obj)

	if _, exists := obj["b"]; exists {
		t.Error("Expected 'b' to be filtered out")
	}
}
