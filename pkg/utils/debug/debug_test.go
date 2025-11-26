package debug

import (
	"testing"
	"time"
)

func TestStack(t *testing.T) {
	stack := Stack()
	if len(stack) == 0 {
		t.Error("Expected non-empty stack")
	}
}

func TestCaller(t *testing.T) {
	file, line, function := Caller(0)
	if file == "" {
		t.Error("Expected non-empty file")
	}
	if line == 0 {
		t.Error("Expected non-zero line")
	}
	if function == "" {
		t.Error("Expected non-empty function")
	}
}

func TestCallers(t *testing.T) {
	callers := Callers(0, 5)
	if len(callers) == 0 {
		t.Error("Expected non-empty callers")
	}
}

func TestFuncName(t *testing.T) {
	name := FuncName(0)
	if name == "" {
		t.Error("Expected non-empty function name")
	}
}

func TestFileLine(t *testing.T) {
	file, line := FileLine(0)
	if file == "" {
		t.Error("Expected non-empty file")
	}
	if line == 0 {
		t.Error("Expected non-zero line")
	}
}

func TestTrace(t *testing.T) {
	defer Trace("test")()
	time.Sleep(10 * time.Millisecond)
}

func TestMeasure(t *testing.T) {
	duration := Measure(func() {
		time.Sleep(10 * time.Millisecond)
	})
	if duration < 10*time.Millisecond {
		t.Error("Expected duration >= 10ms")
	}
}

func TestMeasureWithResult(t *testing.T) {
	result, duration := MeasureWithResult(func() int {
		time.Sleep(10 * time.Millisecond)
		return 42
	})
	if result != 42 {
		t.Error("Expected result 42")
	}
	if duration < 10*time.Millisecond {
		t.Error("Expected duration >= 10ms")
	}
}

func TestAssert(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()
	Assert(false, "test")
}

func TestAssertEqual(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()
	AssertEqual(1, 2, "test")
}

func TestGetNumGoroutines(t *testing.T) {
	num := GetNumGoroutines()
	if num <= 0 {
		t.Error("Expected positive number of goroutines")
	}
}

func TestGetMemStats(t *testing.T) {
	stats := GetMemStats()
	if stats.Alloc == 0 {
		t.Error("Expected non-zero alloc")
	}
}

func TestSetDebug(t *testing.T) {
	SetDebug(true)
	if !IsDebug {
		t.Error("Expected debug mode to be true")
	}
	SetDebug(false)
	if IsDebug {
		t.Error("Expected debug mode to be false")
	}
}
