package math

import (
	"testing"
)

func TestMax(t *testing.T) {
	if Max(3, 5) != 5 {
		t.Errorf("Expected 5, got %d", Max(3, 5))
	}
}

func TestMin(t *testing.T) {
	if Min(3, 5) != 3 {
		t.Errorf("Expected 3, got %d", Min(3, 5))
	}
}

func TestMaxInts(t *testing.T) {
	nums := []int{1, 5, 3, 9, 2}
	if MaxInts(nums) != 9 {
		t.Errorf("Expected 9, got %d", MaxInts(nums))
	}
}

func TestMinInts(t *testing.T) {
	nums := []int{1, 5, 3, 9, 2}
	if MinInts(nums) != 1 {
		t.Errorf("Expected 1, got %d", MinInts(nums))
	}
}

func TestSum(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	if Sum(nums) != 15 {
		t.Errorf("Expected 15, got %d", Sum(nums))
	}
}

func TestAverage(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	avg := Average(nums)
	if avg != 3.0 {
		t.Errorf("Expected 3.0, got %f", avg)
	}
}

func TestAbs(t *testing.T) {
	if Abs(-5) != 5 {
		t.Errorf("Expected 5, got %d", Abs(-5))
	}
	if Abs(5) != 5 {
		t.Errorf("Expected 5, got %d", Abs(5))
	}
}

func TestClamp(t *testing.T) {
	if Clamp(10, 0, 5) != 5 {
		t.Errorf("Expected 5, got %d", Clamp(10, 0, 5))
	}
	if Clamp(-5, 0, 5) != 0 {
		t.Errorf("Expected 0, got %d", Clamp(-5, 0, 5))
	}
	if Clamp(3, 0, 5) != 3 {
		t.Errorf("Expected 3, got %d", Clamp(3, 0, 5))
	}
}

func TestIsInRange(t *testing.T) {
	if !IsInRange(3, 0, 5) {
		t.Error("Expected true")
	}
	if IsInRange(10, 0, 5) {
		t.Error("Expected false")
	}
}

func TestGCD(t *testing.T) {
	if GCD(48, 18) != 6 {
		t.Errorf("Expected 6, got %d", GCD(48, 18))
	}
}

func TestLCM(t *testing.T) {
	if LCM(12, 18) != 36 {
		t.Errorf("Expected 36, got %d", LCM(12, 18))
	}
}

func TestFactorial(t *testing.T) {
	if Factorial(5) != 120 {
		t.Errorf("Expected 120, got %d", Factorial(5))
	}
}

func TestIsPrime(t *testing.T) {
	if !IsPrime(7) {
		t.Error("Expected 7 to be prime")
	}
	if IsPrime(8) {
		t.Error("Expected 8 not to be prime")
	}
}

func TestFibonacci(t *testing.T) {
	if Fibonacci(10) != 55 {
		t.Errorf("Expected 55, got %d", Fibonacci(10))
	}
}

func TestPercent(t *testing.T) {
	percent := Percent(25, 100)
	if percent != 25.0 {
		t.Errorf("Expected 25.0, got %f", percent)
	}
}

func TestLerp(t *testing.T) {
	result := Lerp(0, 10, 0.5)
	if result != 5.0 {
		t.Errorf("Expected 5.0, got %f", result)
	}
}

func TestDistance(t *testing.T) {
	dist := Distance(0, 0, 3, 4)
	if dist != 5.0 {
		t.Errorf("Expected 5.0, got %f", dist)
	}
}

func TestIsEven(t *testing.T) {
	if !IsEven(4) {
		t.Error("Expected 4 to be even")
	}
	if IsEven(5) {
		t.Error("Expected 5 not to be even")
	}
}

func TestIsOdd(t *testing.T) {
	if !IsOdd(5) {
		t.Error("Expected 5 to be odd")
	}
	if IsOdd(4) {
		t.Error("Expected 4 not to be odd")
	}
}
