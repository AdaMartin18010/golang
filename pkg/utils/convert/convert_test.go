package convert

import (
	"testing"
	"time"
)

func TestToString(t *testing.T) {
	if ToString(42) != "42" {
		t.Errorf("Expected '42', got %s", ToString(42))
	}
	if ToString("hello") != "hello" {
		t.Errorf("Expected 'hello', got %s", ToString("hello"))
	}
	if ToString(true) != "true" {
		t.Errorf("Expected 'true', got %s", ToString(true))
	}
}

func TestToInt(t *testing.T) {
	val, err := ToInt("42")
	if err != nil || val != 42 {
		t.Errorf("Expected (42, nil), got (%d, %v)", val, err)
	}

	val, err = ToInt(42)
	if err != nil || val != 42 {
		t.Errorf("Expected (42, nil), got (%d, %v)", val, err)
	}

	val, err = ToInt(42.5)
	if err != nil || val != 42 {
		t.Errorf("Expected (42, nil), got (%d, %v)", val, err)
	}
}

func TestToInt64(t *testing.T) {
	val, err := ToInt64("42")
	if err != nil || val != 42 {
		t.Errorf("Expected (42, nil), got (%d, %v)", val, err)
	}

	val, err = ToInt64(int64(42))
	if err != nil || val != 42 {
		t.Errorf("Expected (42, nil), got (%d, %v)", val, err)
	}
}

func TestToFloat64(t *testing.T) {
	val, err := ToFloat64("42.5")
	if err != nil || val != 42.5 {
		t.Errorf("Expected (42.5, nil), got (%f, %v)", val, err)
	}

	val, err = ToFloat64(42.5)
	if err != nil || val != 42.5 {
		t.Errorf("Expected (42.5, nil), got (%f, %v)", val, err)
	}
}

func TestToBool(t *testing.T) {
	val, err := ToBool("true")
	if err != nil || !val {
		t.Errorf("Expected (true, nil), got (%v, %v)", val, err)
	}

	val, err = ToBool(1)
	if err != nil || !val {
		t.Errorf("Expected (true, nil), got (%v, %v)", val, err)
	}

	val, err = ToBool(0)
	if err != nil || val {
		t.Errorf("Expected (false, nil), got (%v, %v)", val, err)
	}
}

func TestToBytes(t *testing.T) {
	bytes := ToBytes("hello")
	if string(bytes) != "hello" {
		t.Errorf("Expected 'hello', got %s", string(bytes))
	}
}

func TestMustInt(t *testing.T) {
	val := MustInt("42")
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}
}

func TestToIntDefault(t *testing.T) {
	val := ToIntDefault("42", 0)
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}

	val = ToIntDefault("invalid", 100)
	if val != 100 {
		t.Errorf("Expected 100, got %d", val)
	}
}

func TestToStringSlice(t *testing.T) {
	slice := ToStringSlice([]int{1, 2, 3})
	if len(slice) != 3 || slice[0] != "1" {
		t.Errorf("Expected ['1', '2', '3'], got %v", slice)
	}
}

func TestToIntSlice(t *testing.T) {
	slice, err := ToIntSlice([]string{"1", "2", "3"})
	if err != nil || len(slice) != 3 || slice[0] != 1 {
		t.Errorf("Expected ([1, 2, 3], nil), got (%v, %v)", slice, err)
	}
}

func TestIsNumeric(t *testing.T) {
	if !IsNumeric(42) {
		t.Error("Expected 42 to be numeric")
	}
	if !IsNumeric("42") {
		t.Error("Expected '42' to be numeric")
	}
	if IsNumeric("hello") {
		t.Error("Expected 'hello' not to be numeric")
	}
}

func TestIsInteger(t *testing.T) {
	if !IsInteger(42) {
		t.Error("Expected 42 to be integer")
	}
	if !IsInteger("42") {
		t.Error("Expected '42' to be integer")
	}
	if IsInteger("42.5") {
		t.Error("Expected '42.5' not to be integer")
	}
}

func TestToMapStringString(t *testing.T) {
	m := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	result, err := ToMapStringString(m)
	if err != nil || result["key1"] != "value1" || result["key2"] != "42" {
		t.Errorf("Expected map with string values, got %v, %v", result, err)
	}
}

func TestToStringWithTime(t *testing.T) {
	now := time.Now()
	str := ToString(now)
	if str == "" {
		t.Error("Expected non-empty string for time.Time")
	}
}
