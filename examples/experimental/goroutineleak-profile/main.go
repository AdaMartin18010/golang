// Package main demonstrates Go 1.26 experimental goroutine leak profile
//
// Enable with: GOEXPERIMENT=goroutineleakprofile go run main.go
//
// This feature adds a new "goroutineleak" profile to runtime/pprof.
// Unlike the regular goroutine profile which lists all alive goroutines,
// the goroutineleak profile identifies goroutines that are blocked on
// concurrency primitives (channels, mutexes) and are unreachable from
// any runnable goroutine — i.e., they will never unblock.
//
// The feature is zero-overhead when not in use.
// It is expected to be enabled by default in Go 1.27.
package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	fmt.Println("Go 1.26 Experimental: goroutineleak profile")
	fmt.Println("============================================")
	fmt.Println("Run with: GOEXPERIMENT=goroutineleakprofile go run main.go")
	fmt.Println()

	// Create a deliberate goroutine leak
	leakChannel()

	// Give the runtime time to detect the unreachable goroutines
	time.Sleep(100 * time.Millisecond)

	// Try to write the goroutineleak profile
	prof := pprof.Lookup("goroutineleak")
	if prof == nil {
		fmt.Println("goroutineleak profile not available.")
		fmt.Println("Please rebuild with GOEXPERIMENT=goroutineleakprofile")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  GOEXPERIMENT=goroutineleakprofile go run main.go")
		return
	}

	fmt.Println("=== goroutineleak profile ===")
	if err := prof.WriteTo(os.Stdout, 1); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing profile: %v\n", err)
		os.Exit(1)
	}
}

// leakChannel creates a classic goroutine leak:
// goroutines blocked on an unbuffered channel that becomes unreachable.
func leakChannel() {
	ch := make(chan int)

	// These goroutines will block forever on channel receive
	for i := 0; i < 3; i++ {
		go func(id int) {
			<-ch // Block forever - no sender
			fmt.Printf("This will never print: %d\n", id)
		}(i)
	}

	// ch is local; after leakChannel returns, ch is unreachable.
	// The blocked goroutines are also unreachable => leak.
}
