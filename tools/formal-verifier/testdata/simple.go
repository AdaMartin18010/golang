package main

import "fmt"

// Simple function to test CFG construction
func process(x int) int {
	if x > 0 {
		fmt.Println("Positive")
		return x * 2
	} else {
		fmt.Println("Non-positive")
		return x / 2
	}
}

// Function with loop to test CFG construction
func sum(n int) int {
	total := 0
	for i := 0; i < n; i++ {
		total += i
	}
	return total
}

// Function with range loop
func sumSlice(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Function with switch statement
func classify(x int) string {
	switch {
	case x < 0:
		return "negative"
	case x == 0:
		return "zero"
	case x > 0:
		return "positive"
	default:
		return "unknown"
	}
}

func main() {
	x := 10
	result := process(x)
	fmt.Printf("Result: %d\n", result)

	sum10 := sum(10)
	fmt.Printf("Sum: %d\n", sum10)

	nums := []int{1, 2, 3, 4, 5}
	total := sumSlice(nums)
	fmt.Printf("Total: %d\n", total)

	class := classify(5)
	fmt.Printf("Class: %s\n", class)
}
