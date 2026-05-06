// Package main demonstrates Go 1.26 experimental runtime/secret package
//
// Enable with: GOEXPERIMENT=runtimesecret go run main.go
//
// runtime/secret provides secret.Do, which guarantees that after the
// function returns, registers and stack memory used during execution
// are cleared. This is critical for handling ephemeral keys, passwords,
// and other sensitive data that must not survive in memory.
//
// Note: This is an experimental feature in Go 1.26. API may change.
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Go 1.26 Experimental: runtime/secret")
	fmt.Println("=====================================")
	fmt.Println("Run with: GOEXPERIMENT=runtimesecret go run main.go")
	fmt.Println()

	// Since runtime/secret requires GOEXPERIMENT, we provide a conceptual
	// demonstration that compiles without the experiment flag.
	fmt.Println("Conceptual demonstration (compile with GOEXPERIMENT=runtimesecret):")
	fmt.Println()
	fmt.Println("import \"runtime/secret\"")
	fmt.Println()
	fmt.Println("func handleSensitiveData() {")
	fmt.Println("    secret.Do(func() {")
	fmt.Println("        // All memory used here (registers, stack, heap)")
	fmt.Println("        // is guaranteed to be cleared after return.")
	fmt.Println("        key := generateEphemeralKey()")
	fmt.Println("        ciphertext := encrypt(data, key)")
	fmt.Println("        _ = ciphertext")
	fmt.Println("    })")
	fmt.Println("}")
	fmt.Println()

	if os.Getenv("GOEXPERIMENT") == "runtimesecret" {
		fmt.Println("✅ GOEXPERIMENT=runtimesecret is active!")
		fmt.Println("   In a real build, secret.Do would clear sensitive state.")
	} else {
		fmt.Println("⚠️  GOEXPERIMENT not set. This is a conceptual demo.")
		fmt.Println("   To run with actual secret clearing:")
		fmt.Println("   GOEXPERIMENT=runtimesecret go run main.go")
	}
}
