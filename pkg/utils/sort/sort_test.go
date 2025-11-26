package sort

import (
	"testing"
)

func TestInts(t *testing.T) {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	Ints(nums)
	if !IntsAreSorted(nums) {
		t.Error("Expected sorted slice")
	}
}

func TestFloat64s(t *testing.T) {
	nums := []float64{3.1, 1.4, 4.1, 1.5, 9.2, 6.5}
	Float64s(nums)
	if !Float64sAreSorted(nums) {
		t.Error("Expected sorted slice")
	}
}

func TestStrings(t *testing.T) {
	strs := []string{"banana", "apple", "cherry"}
	Strings(strs)
	if !StringsAreSorted(strs) {
		t.Error("Expected sorted slice")
	}
}

func TestIntsReverse(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	IntsReverse(nums)
	expected := []int{5, 4, 3, 2, 1}
	for i, v := range nums {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestSortByFunc(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	SortByFunc(people, func(a, b Person) bool {
		return a.Age < b.Age
	})
	if people[0].Name != "Bob" {
		t.Errorf("Expected Bob, got %s", people[0].Name)
	}
}

func TestReverse(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	Reverse(nums)
	expected := []int{5, 4, 3, 2, 1}
	for i, v := range nums {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestUniqueInts(t *testing.T) {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	unique := UniqueInts(nums)
	if len(unique) != 7 {
		t.Errorf("Expected 7 unique elements, got %d", len(unique))
	}
	if !IntsAreSorted(unique) {
		t.Error("Expected sorted unique slice")
	}
}

func TestTopNInts(t *testing.T) {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	top3 := TopNInts(nums, 3)
	if len(top3) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(top3))
	}
	if top3[0] != 9 {
		t.Errorf("Expected 9, got %d", top3[0])
	}
}

func TestBottomNInts(t *testing.T) {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	bottom3 := BottomNInts(nums, 3)
	if len(bottom3) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(bottom3))
	}
	if bottom3[0] != 1 {
		t.Errorf("Expected 1, got %d", bottom3[0])
	}
}

func TestSortByKeyInt(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	SortByKeyInt(people, func(p Person) int {
		return p.Age
	})
	if people[0].Name != "Bob" {
		t.Errorf("Expected Bob, got %s", people[0].Name)
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
