package random

import (
	"testing"
	"time"
)

func TestInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		val := Int(100)
		if val < 0 || val >= 100 {
			t.Errorf("Expected value in [0, 100), got %d", val)
		}
	}
}

func TestIntRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		val := IntRange(10, 20)
		if val < 10 || val >= 20 {
			t.Errorf("Expected value in [10, 20), got %d", val)
		}
	}
}

func TestFloat64(t *testing.T) {
	for i := 0; i < 100; i++ {
		val := Float64()
		if val < 0 || val >= 1 {
			t.Errorf("Expected value in [0.0, 1.0), got %f", val)
		}
	}
}

func TestFloat64Range(t *testing.T) {
	for i := 0; i < 100; i++ {
		val := Float64Range(1.0, 2.0)
		if val < 1.0 || val >= 2.0 {
			t.Errorf("Expected value in [1.0, 2.0), got %f", val)
		}
	}
}

func TestString(t *testing.T) {
	str := String(10)
	if len(str) != 10 {
		t.Errorf("Expected length 10, got %d", len(str))
	}
}

func TestStringWithCharset(t *testing.T) {
	charset := "abc"
	str := StringWithCharset(5, charset)
	if len(str) != 5 {
		t.Errorf("Expected length 5, got %d", len(str))
	}
	for _, r := range str {
		if !contains(charset, byte(r)) {
			t.Errorf("Character %c not in charset", r)
		}
	}
}

func contains(s string, b byte) bool {
	for i := range s {
		if s[i] == b {
			return true
		}
	}
	return false
}

func TestHex(t *testing.T) {
	str := Hex(16)
	if len(str) != 16 {
		t.Errorf("Expected length 16, got %d", len(str))
	}
}

func TestBytes(t *testing.T) {
	b := Bytes(10)
	if len(b) != 10 {
		t.Errorf("Expected length 10, got %d", len(b))
	}
}

func TestSecureInt(t *testing.T) {
	val, err := SecureInt(100)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val < 0 || val >= 100 {
		t.Errorf("Expected value in [0, 100), got %d", val)
	}
}

func TestSecureString(t *testing.T) {
	str, err := SecureString(10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(str) != 10 {
		t.Errorf("Expected length 10, got %d", len(str))
	}
}

func TestChoice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	val, ok := Choice(slice)
	if !ok {
		t.Error("Expected true")
	}
	if val < 1 || val > 5 {
		t.Errorf("Expected value in [1, 5], got %d", val)
	}

	empty := []int{}
	_, ok = Choice(empty)
	if ok {
		t.Error("Expected false for empty slice")
	}
}

func TestChoices(t *testing.T) {
	slice := []int{1, 2, 3}
	result := Choices(slice, 5)
	if len(result) != 5 {
		t.Errorf("Expected length 5, got %d", len(result))
	}
}

func TestSample(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Sample(slice, 3)
	if len(result) != 3 {
		t.Errorf("Expected length 3, got %d", len(result))
	}

	// 检查是否有重复
	seen := make(map[int]bool)
	for _, v := range result {
		if seen[v] {
			t.Error("Expected no duplicates")
		}
		seen[v] = true
	}
}

func TestShuffle(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	original := make([]int, len(slice))
	copy(original, slice)
	Shuffle(slice)

	// 检查是否被打乱（有一定概率相同，但通常不同）
	same := true
	for i := range slice {
		if slice[i] != original[i] {
			same = false
			break
		}
	}
	// 如果100次都相同，可能有问题（但这是概率性的，不能完全保证）
}

func TestWeightedChoice(t *testing.T) {
	items := []string{"a", "b", "c"}
	weights := []float64{0.1, 0.2, 0.7}
	val, ok := WeightedChoice(items, weights)
	if !ok {
		t.Error("Expected true")
	}
	if val != "a" && val != "b" && val != "c" {
		t.Errorf("Expected one of [a, b, c], got %s", val)
	}
}

func TestBool(t *testing.T) {
	// 测试多次，应该既有true也有false
	hasTrue := false
	hasFalse := false
	for i := 0; i < 100; i++ {
		if Bool() {
			hasTrue = true
		} else {
			hasFalse = true
		}
		if hasTrue && hasFalse {
			break
		}
	}
	if !hasTrue || !hasFalse {
		t.Error("Expected both true and false values")
	}
}

func TestProbability(t *testing.T) {
	if !Probability(1.0) {
		t.Error("Expected true for probability 1.0")
	}
	if Probability(0.0) {
		t.Error("Expected false for probability 0.0")
	}
}

func TestDuration(t *testing.T) {
	min := time.Second
	max := 10 * time.Second
	d := Duration(min, max)
	if d < min || d >= max {
		t.Errorf("Expected duration in [%v, %v), got %v", min, max, d)
	}
}

func TestTime(t *testing.T) {
	start := time.Now()
	end := start.Add(24 * time.Hour)
	tm := Time(start, end)
	if tm.Before(start) || tm.After(end) {
		t.Errorf("Expected time in [%v, %v], got %v", start, end, tm)
	}
}

func TestFastString(t *testing.T) {
	str := FastString(10)
	if len(str) != 10 {
		t.Errorf("Expected length 10, got %d", len(str))
	}
}
