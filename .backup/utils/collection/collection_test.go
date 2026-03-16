package collection

import (
	"testing"
)

func TestContains(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	if !Contains(slice, 3) {
		t.Error("Expected slice to contain 3")
	}
	if Contains(slice, 6) {
		t.Error("Expected slice not to contain 6")
	}
}

func TestIndex(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	if Index(slice, 3) != 2 {
		t.Error("Expected index 2 for value 3")
	}
	if Index(slice, 6) != -1 {
		t.Error("Expected index -1 for value 6")
	}
}

func TestRemove(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Remove(slice, 3)
	if Contains(result, 3) {
		t.Error("Expected result not to contain 3")
	}
	if len(result) != 4 {
		t.Error("Expected result length to be 4")
	}
}

func TestUnique(t *testing.T) {
	slice := []int{1, 2, 2, 3, 3, 3, 4}
	result := Unique(slice)
	if len(result) != 4 {
		t.Errorf("Expected 4 unique elements, got %d", len(result))
	}
}

func TestFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Filter(slice, func(x int) bool {
		return x%2 == 0
	})
	if len(result) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(result))
	}
}

func TestMap(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Map(slice, func(x int) int {
		return x * 2
	})
	if len(result) != 5 {
		t.Errorf("Expected 5 elements, got %d", len(result))
	}
	if result[0] != 2 {
		t.Errorf("Expected 2, got %d", result[0])
	}
}

func TestReduce(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Reduce(slice, 0, func(acc, x int) int {
		return acc + x
	})
	if result != 15 {
		t.Errorf("Expected 15, got %d", result)
	}
}

func TestChunk(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7}
	result := Chunk(slice, 3)
	if len(result) != 3 {
		t.Errorf("Expected 3 chunks, got %d", len(result))
	}
	if len(result[0]) != 3 {
		t.Errorf("Expected first chunk size 3, got %d", len(result[0]))
	}
}

func TestReverse(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Reverse(slice)
	if result[0] != 5 {
		t.Errorf("Expected 5, got %d", result[0])
	}
}

func TestFirst(t *testing.T) {
	slice := []int{1, 2, 3}
	first, ok := First(slice)
	if !ok || first != 1 {
		t.Errorf("Expected (1, true), got (%d, %v)", first, ok)
	}

	empty := []int{}
	_, ok = First(empty)
	if ok {
		t.Error("Expected false for empty slice")
	}
}

func TestLast(t *testing.T) {
	slice := []int{1, 2, 3}
	last, ok := Last(slice)
	if !ok || last != 3 {
		t.Errorf("Expected (3, true), got (%d, %v)", last, ok)
	}
}

func TestTake(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Take(slice, 3)
	if len(result) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(result))
	}
}

func TestDrop(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Drop(slice, 2)
	if len(result) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(result))
	}
	if result[0] != 3 {
		t.Errorf("Expected 3, got %d", result[0])
	}
}

func TestPartition(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	even, odd := Partition(slice, func(x int) bool {
		return x%2 == 0
	})
	if len(even) != 2 {
		t.Errorf("Expected 2 even numbers, got %d", len(even))
	}
	if len(odd) != 3 {
		t.Errorf("Expected 3 odd numbers, got %d", len(odd))
	}
}

func TestGroupBy(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := GroupBy(slice, func(x int) string {
		if x%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if len(result["even"]) != 2 {
		t.Errorf("Expected 2 even numbers, got %d", len(result["even"]))
	}
}

func TestCount(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	count := Count(slice, func(x int) bool {
		return x%2 == 0
	})
	if count != 2 {
		t.Errorf("Expected 2, got %d", count)
	}
}

func TestAny(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	if !Any(slice, func(x int) bool {
		return x > 3
	}) {
		t.Error("Expected true")
	}
	if Any(slice, func(x int) bool {
		return x > 10
	}) {
		t.Error("Expected false")
	}
}

func TestAll(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	if !All(slice, func(x int) bool {
		return x > 0
	}) {
		t.Error("Expected true")
	}
	if All(slice, func(x int) bool {
		return x > 3
	}) {
		t.Error("Expected false")
	}
}

func TestSum(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	sum := Sum(slice)
	if sum != 15 {
		t.Errorf("Expected 15, got %d", sum)
	}
}

func TestMax(t *testing.T) {
	slice := []int{1, 5, 3, 2, 4}
	max, ok := Max(slice)
	if !ok || max != 5 {
		t.Errorf("Expected (5, true), got (%d, %v)", max, ok)
	}
}

func TestMin(t *testing.T) {
	slice := []int{5, 1, 3, 2, 4}
	min, ok := Min(slice)
	if !ok || min != 1 {
		t.Errorf("Expected (1, true), got (%d, %v)", min, ok)
	}
}

func TestIntersect(t *testing.T) {
	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4, 5, 6}
	result := Intersect(slice1, slice2)
	if len(result) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(result))
	}
}

func TestUnion(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{3, 4, 5}
	result := Union(slice1, slice2)
	if len(result) != 5 {
		t.Errorf("Expected 5 elements, got %d", len(result))
	}
}

func TestDifference(t *testing.T) {
	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4}
	result := Difference(slice1, slice2)
	if len(result) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(result))
	}
}

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := MapKeys(m)
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
}

func TestMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	values := MapValues(m)
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}
}

func TestMapContains(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	if !MapContains(m, "a") {
		t.Error("Expected map to contain key 'a'")
	}
	if MapContains(m, "c") {
		t.Error("Expected map not to contain key 'c'")
	}
}

func TestMapGet(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	if MapGet(m, "a", 0) != 1 {
		t.Error("Expected 1")
	}
	if MapGet(m, "c", 99) != 99 {
		t.Error("Expected 99")
	}
}
