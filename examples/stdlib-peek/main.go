// Package main demonstrates Go 1.26 bytes.Buffer.Peek
//
// Peek returns the next n bytes without advancing the buffer.
// This is useful for parsing protocols where you need to inspect
// data before committing to consume it.
package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println("Go 1.26 bytes.Buffer.Peek Demo")
	fmt.Println("================================")

	buf := bytes.NewBufferString("Hello, World!")

	// Peek at the first 5 bytes without consuming them
	peeked, err := buf.Peek(5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Peek(5): %q\n", peeked)

	// The buffer position is unchanged
	remaining := buf.Bytes()
	fmt.Printf("Buffer after Peek: %q\n", remaining)

	// Now actually consume some bytes
	consumed := make([]byte, 5)
	n, err := buf.Read(consumed)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read(%d): %q\n", n, consumed[:n])

	// Buffer is now shorter
	remaining = buf.Bytes()
	fmt.Printf("Buffer after Read: %q\n", remaining)

	// Peek beyond buffer length returns available bytes + error
	smallBuf := bytes.NewBufferString("Hi")
	peeked, err = smallBuf.Peek(10)
	fmt.Printf("\nPeek(10) from 2-byte buffer:\n")
	fmt.Printf("  Data: %q\n", peeked)
	fmt.Printf("  Error: %v\n", err)
}
