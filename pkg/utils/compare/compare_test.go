package compare

import (
	"testing"
	"time"
)

func TestEqual(t *testing.T) {
	if !Equal(1, 1) {
		t.Error("Expected 1 == 1")
	}
	if Equal(1, 2) {
		t.Error("Expected 1 != 2")
	}
}

func TestCompareInt(t *testing.T) {
	if CompareInt(1, 2) != -1 {
		t.Error("Expected -1")
	}
	if CompareInt(2, 1) != 1 {
		t.Error("Expected 1")
	}
	if CompareInt(1, 1) != 0 {
		t.Error("Expected 0")
	}
}

func TestCompareFloat64(t *testing.T) {
	if CompareFloat64(1.0, 2.0) != -1 {
		t.Error("Expected -1")
	}
	if CompareFloat64(2.0, 1.0) != 1 {
		t.Error("Expected 1")
	}
	if CompareFloat64(1.0, 1.0) != 0 {
		t.Error("Expected 0")
	}
}

func TestCompareString(t *testing.T) {
	if CompareString("a", "b") != -1 {
		t.Error("Expected -1")
	}
	if CompareString("b", "a") != 1 {
		t.Error("Expected 1")
	}
	if CompareString("a", "a") != 0 {
		t.Error("Expected 0")
	}
}

func TestCompareTime(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	if CompareTime(t1, t2) != -1 {
		t.Error("Expected -1")
	}
	if CompareTime(t2, t1) != 1 {
		t.Error("Expected 1")
	}
	if CompareTime(t1, t1) != 0 {
		t.Error("Expected 0")
	}
}

func TestLess(t *testing.T) {
	if !Less(1, 2) {
		t.Error("Expected 1 < 2")
	}
	if Less(2, 1) {
		t.Error("Expected 2 not < 1")
	}
}

func TestGreater(t *testing.T) {
	if !Greater(2, 1) {
		t.Error("Expected 2 > 1")
	}
	if Greater(1, 2) {
		t.Error("Expected 1 not > 2")
	}
}

func TestMinInt(t *testing.T) {
	if MinInt(1, 2) != 1 {
		t.Error("Expected 1")
	}
}

func TestMaxInt(t *testing.T) {
	if MaxInt(1, 2) != 2 {
		t.Error("Expected 2")
	}
}

func TestInRangeInt(t *testing.T) {
	if !InRangeInt(5, 1, 10) {
		t.Error("Expected 5 in range [1, 10]")
	}
	if InRangeInt(15, 1, 10) {
		t.Error("Expected 15 not in range [1, 10]")
	}
}

func TestClampInt(t *testing.T) {
	if ClampInt(15, 1, 10) != 10 {
		t.Error("Expected 10")
	}
	if ClampInt(-5, 1, 10) != 1 {
		t.Error("Expected 1")
	}
	if ClampInt(5, 1, 10) != 5 {
		t.Error("Expected 5")
	}
}

func TestIsZero(t *testing.T) {
	if !IsZero(0) {
		t.Error("Expected 0 to be zero")
	}
	if !IsZero("") {
		t.Error("Expected empty string to be zero")
	}
	if IsZero(1) {
		t.Error("Expected 1 not to be zero")
	}
}

func TestIsNil(t *testing.T) {
	var s []int
	if !IsNil(s) {
		t.Error("Expected nil slice to be nil")
	}
	if IsNil(1) {
		t.Error("Expected 1 not to be nil")
	}
}

func TestIsEmpty(t *testing.T) {
	if !IsEmpty(0) {
		t.Error("Expected 0 to be empty")
	}
	if !IsEmpty("") {
		t.Error("Expected empty string to be empty")
	}
	var s []int
	if !IsEmpty(s) {
		t.Error("Expected nil slice to be empty")
	}
}

func TestCompareSlice(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{1, 2, 4}
	if CompareSlice(a, b) != -1 {
		t.Error("Expected -1")
	}
}

func TestCompareMap(t *testing.T) {
	a := map[string]int{"a": 1, "b": 2}
	b := map[string]int{"a": 1, "b": 2}
	if !CompareMap(a, b) {
		t.Error("Expected maps to be equal")
	}
	c := map[string]int{"a": 1, "b": 3}
	if CompareMap(a, c) {
		t.Error("Expected maps not to be equal")
	}
}
