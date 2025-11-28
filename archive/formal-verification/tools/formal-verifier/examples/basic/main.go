// Package main demonstrates common code quality issues that FV can detect
package main

import (
	"fmt"
	"sync"
	"time"
)

// Example 1: High Complexity Function
// This function has too many nested conditionals and branches
func processOrder(orderID int, items []string, discount float64,
	shipping bool, priority int, customer map[string]string) error {

	if orderID <= 0 {
		return fmt.Errorf("invalid order ID")
	}

	if len(items) == 0 {
		return fmt.Errorf("no items")
	}

	total := 0.0
	for _, item := range items {
		if item == "book" {
			total += 10.0
		} else if item == "pen" {
			total += 2.0
		} else if item == "notebook" {
			total += 5.0
		} else if item == "laptop" {
			total += 1000.0
		} else if item == "mouse" {
			total += 20.0
		} else {
			total += 1.0
		}
	}

	if discount > 0 {
		if discount < 0.1 {
			total *= 0.95
		} else if discount < 0.2 {
			total *= 0.9
		} else if discount < 0.3 {
			total *= 0.85
		} else {
			total *= 0.8
		}
	}

	if shipping {
		if priority == 1 {
			total += 20.0
		} else if priority == 2 {
			total += 10.0
		} else {
			total += 5.0
		}
	}

	fmt.Printf("Order %d total: $%.2f\n", orderID, total)
	return nil
}

// Example 2: Goroutine Leak
// The goroutine will never exit because the channel is never closed
func startWorker() {
	ch := make(chan int)

	go func() {
		for v := range ch {
			fmt.Println("Processing:", v)
		}
	}()

	ch <- 1
	ch <- 2
	// Missing: close(ch)
}

// Example 3: Data Race
// Multiple goroutines accessing shared variable without synchronization
var counter int

func incrementCounter() {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // Data race here!
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}

// Example 4: Unsafe Type Assertion
// Type assertion without checking can panic
func processValue(v interface{}) {
	str := v.(string) // Unsafe! What if v is not a string?
	fmt.Println("Value:", str)
}

// Example 5: Unbuffered Channel Misuse
// Can cause deadlock if not handled properly
func sendData() {
	ch := make(chan string)

	ch <- "data" // This will block forever if no receiver
	fmt.Println("Sent data")
}

func main() {
	fmt.Println("Running examples with code quality issues...")
	fmt.Println("Use 'fv analyze' to detect these issues!")

	// Run examples
	processOrder(1, []string{"book", "pen"}, 0.1, true, 1, nil)
	startWorker()
	incrementCounter()
	processValue("hello")

	time.Sleep(100 * time.Millisecond)
}
