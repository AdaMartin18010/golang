# TS-DT-010: Go Fuzzing

> **维度**: Technology Stack > Development Tools
> **级别**: S (16+ KB)
> **标签**: #fuzzing #testing #golang #security #fuzz-testing
> **权威来源**:
>
> - [Go Fuzzing Tutorial](https://go.dev/doc/security/fuzz/) - Go team
> - [Native Go Fuzzing](https://go.dev/doc/fuzz/) - Go documentation

---

## 1. Fuzzing Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Go Fuzzing Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Fuzzing Process:                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  1. Seed Corpus                                                      │   │
│  │     ├── Valid inputs to start with                                   │   │
│  │     ├── Example: "hello", "12345", "test@example.com"               │   │
│  │     └── Stored in testdata/fuzz/FuzzName/*                          │   │
│  │                                                                      │   │
│  │  2. Fuzzer generates mutations                                       │   │
│  │     ├── Bit flipping                                                 │   │
│  │     ├── Byte insertion/deletion                                      │   │
│  │     ├── Interesting values (0, -1, MAX_INT)                         │   │
│  │     └── Dictionary words                                             │   │
│  │                                                                      │   │
│  │  3. Test function executes                                           │   │
│  │     └── func FuzzName(f *testing.F)                                 │   │
│  │                                                                      │   │
│  │  4. Coverage guidance                                                │   │
│  │     ├── Track which code paths are executed                          │   │
│  │     ├── Prioritize inputs that find new paths                        │   │
│  │     └── Continue until crash or timeout                              │   │
│  │                                                                      │   │
│  │  5. Findings                                                         │   │
│  │     ├── Crashes (panics, errors)                                     │   │
│  │     ├── Hangs (infinite loops)                                       │   │
│  │     └── OOM (memory exhaustion)                                      │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Benefits:                                                                   │
│  - Find edge cases and bugs automatically                                    │
│  - Discover security vulnerabilities                                         │
│  - Test with inputs you wouldn't think of                                    │
│  - Continuous improvement with coverage guidance                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Writing Fuzz Tests

```go
package parser

import (
    "testing"
)

// FuzzParseJSON tests JSON parsing with random inputs
func FuzzParseJSON(f *testing.F) {
    // Seed corpus - valid inputs to start
    testcases := []string{
        `{"name": "John", "age": 30}`,
        `[]`,
        `{}`,
        `null`,
        `true`,
        `false`,
        `123`,
        `"string"`,
        `[1, 2, 3]`,
    }

    for _, tc := range testcases {
        f.Add(tc) // Add to seed corpus
    }

    // Fuzz target
    f.Fuzz(func(t *testing.T, input string) {
        // Function under test
        result, err := ParseJSON(input)

        // Check for panics
        // (Fuzzing automatically catches panics)

        // Validate: if no error, result should be usable
        if err == nil && result == nil {
            t.Error("nil result without error")
        }

        // Validate: certain errors are expected
        if err != nil {
            // Error should be one of known types
            if !IsValidJSONError(err) {
                t.Errorf("unexpected error type: %v", err)
            }
        }
    })
}

// FuzzStringReverse tests string reversal
func FuzzStringReverse(f *testing.F) {
    f.Add("hello", "world")
    f.Add("", "empty")
    f.Add("1234567890", "numbers")

    f.Fuzz(func(t *testing.T, input string, expected string) {
        reversed := ReverseString(input)

        // Property: reverse twice should equal original
        doubleReversed := ReverseString(reversed)
        if doubleReversed != input {
            t.Errorf("Reverse(Reverse(%q)) = %q, want %q", input, doubleReversed, input)
        }
    })
}

// FuzzCalculate tests mathematical calculation
func FuzzCalculate(f *testing.F) {
    // Seed with boundary values
    f.Add(int64(0), int64(0))
    f.Add(int64(1), int64(1))
    f.Add(int64(-1), int64(-1))
    f.Add(int64(9223372036854775807), int64(1))  // Max int64
    f.Add(int64(-9223372036854775808), int64(1)) // Min int64

    f.Fuzz(func(t *testing.T, a, b int64) {
        result, err := SafeAdd(a, b)

        // Check for overflow/underflow
        if err != nil {
            // Should return error on overflow
            if (b > 0 && a > 0 && result < 0) ||
               (b < 0 && a < 0 && result > 0) {
                // Expected overflow
                return
            }
            t.Errorf("unexpected error: %v", err)
        }

        // Property: a + b = b + a
        result2, _ := SafeAdd(b, a)
        if result != result2 {
            t.Errorf("SafeAdd(%d, %d) != SafeAdd(%d, %d)", a, b, b, a)
        }
    })
}
```

---

## 3. Running Fuzz Tests

```bash
# Run fuzzing for 10 seconds
go test -fuzz=FuzzParseJSON -fuzztime=10s

# Run fuzzing until crash or manual stop
go test -fuzz=FuzzParseJSON

# Run with verbose output
go test -v -fuzz=FuzzParseJSON

# Run with multiple workers
go test -fuzz=FuzzParseJSON -parallel=4

# Run specific fuzz test with corpus
go test -run=FuzzParseJSON ./testdata/fuzz/FuzzParseJSON/...

# Minimize crash input
go test -fuzz=FuzzParseJSON -minimize=1000

# View fuzzing coverage
go test -fuzz=FuzzParseJSON -cover

# Fuzz with memory limit
go test -fuzz=FuzzParseJSON -fuzzminimizetime=30s
```

---

## 4. Best Practices

```
Fuzzing Best Practices:
□ Provide good seed corpus
□ Check properties, not specific outputs
□ Handle expected errors gracefully
□ Use appropriate types (string, []byte, int, etc.)
□ Run fuzzing in CI/CD
□ Save crashers for regression testing
□ Minimize crash inputs
□ Document found bugs
```
