package env

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	value := Get("TEST_KEY", "default")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %s", value)
	}

	value = Get("NON_EXISTENT", "default")
	if value != "default" {
		t.Errorf("Expected 'default', got %s", value)
	}
}

func TestGetInt(t *testing.T) {
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	value := GetInt("TEST_INT", 0)
	if value != 123 {
		t.Errorf("Expected 123, got %d", value)
	}

	value = GetInt("NON_EXISTENT", 999)
	if value != 999 {
		t.Errorf("Expected 999, got %d", value)
	}
}

func TestGetInt64(t *testing.T) {
	os.Setenv("TEST_INT64", "123456789")
	defer os.Unsetenv("TEST_INT64")

	value := GetInt64("TEST_INT64", 0)
	if value != 123456789 {
		t.Errorf("Expected 123456789, got %d", value)
	}
}

func TestGetFloat64(t *testing.T) {
	os.Setenv("TEST_FLOAT", "123.456")
	defer os.Unsetenv("TEST_FLOAT")

	value := GetFloat64("TEST_FLOAT", 0.0)
	if value != 123.456 {
		t.Errorf("Expected 123.456, got %f", value)
	}
}

func TestGetBool(t *testing.T) {
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")

	value := GetBool("TEST_BOOL", false)
	if !value {
		t.Error("Expected true")
	}

	os.Setenv("TEST_BOOL", "false")
	value = GetBool("TEST_BOOL", true)
	if value {
		t.Error("Expected false")
	}
}

func TestGetSlice(t *testing.T) {
	os.Setenv("TEST_SLICE", "a,b,c")
	defer os.Unsetenv("TEST_SLICE")

	value := GetSlice("TEST_SLICE", []string{})
	if len(value) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(value))
	}
}

func TestSet(t *testing.T) {
	err := Set("TEST_SET", "value")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer os.Unsetenv("TEST_SET")

	if os.Getenv("TEST_SET") != "value" {
		t.Error("Expected environment variable to be set")
	}
}

func TestUnset(t *testing.T) {
	os.Setenv("TEST_UNSET", "value")
	err := Unset("TEST_UNSET")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if Has("TEST_UNSET") {
		t.Error("Expected environment variable to be unset")
	}
}

func TestHas(t *testing.T) {
	os.Setenv("TEST_HAS", "value")
	defer os.Unsetenv("TEST_HAS")

	if !Has("TEST_HAS") {
		t.Error("Expected environment variable to exist")
	}

	if Has("NON_EXISTENT") {
		t.Error("Expected environment variable not to exist")
	}
}

func TestGetWithPrefix(t *testing.T) {
	os.Setenv("TEST_PREFIX_1", "value1")
	os.Setenv("TEST_PREFIX_2", "value2")
	os.Setenv("OTHER_KEY", "value3")
	defer func() {
		os.Unsetenv("TEST_PREFIX_1")
		os.Unsetenv("TEST_PREFIX_2")
		os.Unsetenv("OTHER_KEY")
	}()

	result := GetWithPrefix("TEST_PREFIX_")
	if len(result) < 2 {
		t.Errorf("Expected at least 2 results, got %d", len(result))
	}
}

func TestExpand(t *testing.T) {
	os.Setenv("TEST_VAR", "world")
	defer os.Unsetenv("TEST_VAR")

	result := Expand("Hello ${TEST_VAR}")
	if result != "Hello world" {
		t.Errorf("Expected 'Hello world', got %s", result)
	}
}

func TestValidateRequired(t *testing.T) {
	os.Setenv("REQUIRED_1", "value1")
	defer os.Unsetenv("REQUIRED_1")

	err := ValidateRequired([]string{"REQUIRED_1"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = ValidateRequired([]string{"REQUIRED_1", "REQUIRED_2"})
	if err == nil {
		t.Error("Expected error for missing required variable")
	}
}

func TestGetSliceWithSeparator(t *testing.T) {
	os.Setenv("TEST_SLICE", "a|b|c")
	defer os.Unsetenv("TEST_SLICE")

	value := GetSliceWithSeparator("TEST_SLICE", "|", []string{})
	if len(value) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(value))
	}
}
