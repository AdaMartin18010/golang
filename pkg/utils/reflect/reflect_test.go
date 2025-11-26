package reflect

import (
	"testing"
)

type TestStruct struct {
	Name  string
	Age   int
	Email string `json:"email" db:"user_email"`
}

func (t *TestStruct) GetName() string {
	return t.Name
}

func (t *TestStruct) SetName(name string) {
	t.Name = name
}

func TestGetType(t *testing.T) {
	if GetType(42) != "int" {
		t.Errorf("Expected 'int', got %s", GetType(42))
	}
}

func TestGetKind(t *testing.T) {
	if GetKind(42) != reflect.Int {
		t.Error("Expected Int kind")
	}
}

func TestIsNil(t *testing.T) {
	var s *string
	if !IsNil(s) {
		t.Error("Expected nil")
	}
	if IsNil(42) {
		t.Error("Expected not nil")
	}
}

func TestIsZero(t *testing.T) {
	if !IsZero(0) {
		t.Error("Expected zero value")
	}
	if IsZero(42) {
		t.Error("Expected not zero value")
	}
}

func TestIsPointer(t *testing.T) {
	var s *string
	if !IsPointer(s) {
		t.Error("Expected pointer")
	}
	if IsPointer(42) {
		t.Error("Expected not pointer")
	}
}

func TestIsSlice(t *testing.T) {
	if !IsSlice([]int{1, 2, 3}) {
		t.Error("Expected slice")
	}
	if IsSlice(42) {
		t.Error("Expected not slice")
	}
}

func TestDereference(t *testing.T) {
	value := 42
	ptr := &value
	result := Dereference(ptr)
	if result != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestGetField(t *testing.T) {
	s := TestStruct{Name: "test", Age: 30}
	value, err := GetField(s, "Name")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != "test" {
		t.Errorf("Expected 'test', got %v", value)
	}
}

func TestSetField(t *testing.T) {
	s := &TestStruct{Name: "test"}
	err := SetField(s, "Name", "new")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if s.Name != "new" {
		t.Errorf("Expected 'new', got %s", s.Name)
	}
}

func TestHasField(t *testing.T) {
	s := TestStruct{}
	if !HasField(s, "Name") {
		t.Error("Expected field to exist")
	}
	if HasField(s, "NonExistent") {
		t.Error("Expected field not to exist")
	}
}

func TestGetFieldNames(t *testing.T) {
	s := TestStruct{}
	fields := GetFieldNames(s)
	if len(fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(fields))
	}
}

func TestCallMethod(t *testing.T) {
	s := &TestStruct{Name: "test"}
	results, err := CallMethod(s, "GetName")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if results[0] != "test" {
		t.Errorf("Expected 'test', got %v", results[0])
	}
}

func TestHasMethod(t *testing.T) {
	s := &TestStruct{}
	if !HasMethod(s, "GetName") {
		t.Error("Expected method to exist")
	}
	if HasMethod(s, "NonExistent") {
		t.Error("Expected method not to exist")
	}
}

func TestNewInstance(t *testing.T) {
	s := TestStruct{}
	newS := NewInstance(s)
	if newS == nil {
		t.Error("Expected new instance")
	}
}

func TestDeepEqual(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	if !DeepEqual(a, b) {
		t.Error("Expected equal")
	}
}

func TestGetSliceElement(t *testing.T) {
	slice := []int{1, 2, 3}
	value, err := GetSliceElement(slice, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != 1 {
		t.Errorf("Expected 1, got %v", value)
	}
}

func TestGetMapValue(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	value, ok := GetMapValue(m, "a")
	if !ok {
		t.Error("Expected value to exist")
	}
	if value != 1 {
		t.Errorf("Expected 1, got %v", value)
	}
}

func TestGetLength(t *testing.T) {
	slice := []int{1, 2, 3}
	length, err := GetLength(slice)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if length != 3 {
		t.Errorf("Expected 3, got %d", length)
	}
}

func TestIsAssignable(t *testing.T) {
	if !IsAssignable(42, 0) {
		t.Error("Expected assignable")
	}
}

func TestIsConvertible(t *testing.T) {
	if !IsConvertible(42, int64(0)) {
		t.Error("Expected convertible")
	}
}
