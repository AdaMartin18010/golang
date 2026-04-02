# TS-CL-009: Go regexp Package

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #regex #regular-expressions #pattern-matching
> **权威来源**:
>
> - [regexp Package](https://golang.org/pkg/regexp/) - Go standard library

---

## 1. Regular Expression Basics

```go
package main

import (
    "fmt"
    "regexp"
)

func basicRegex() {
    // Compile regex
    re := regexp.MustCompile(`\b\w+@\w+\.\w+\b`)

    // Match string
    text := "Contact us at support@example.com or sales@company.org"
    matches := re.FindAllString(text, -1)
    fmt.Println(matches) // [support@example.com sales@company.org]

    // Check if matches
    matched := re.MatchString("test@email.com")
    fmt.Println(matched) // true

    // Find first match
    first := re.FindString(text)
    fmt.Println(first) // support@example.com
}

// Common patterns
var patterns = struct {
    Email    *regexp.Regexp
    Phone    *regexp.Regexp
    URL      *regexp.Regexp
    IP       *regexp.Regexp
    Date     *regexp.Regexp
}{
    Email:    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
    Phone:    regexp.MustCompile(`^\+?\d{1,3}[-.\s]?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}$`),
    URL:      regexp.MustCompile(`https?://[^\s]+`),
    IP:       regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`),
    Date:     regexp.MustCompile(`\d{4}-\d{2}-\d{2}`),
}

func validationExamples() {
    // Email validation
    emails := []string{
        "user@example.com",
        "invalid.email",
        "test@domain.org",
    }

    for _, email := range emails {
        valid := patterns.Email.MatchString(email)
        fmt.Printf("%s: %v\n", email, valid)
    }
}
```

---

## 2. Advanced Patterns

```go
func advancedRegex() {
    // Capture groups
    re := regexp.MustCompile(`(\w+)\s+(\w+)`)
    matches := re.FindStringSubmatch("John Doe")
    fmt.Printf("Full: %s, First: %s, Last: %s\n",
        matches[0], matches[1], matches[2])

    // Named capture groups (Go 1.22+)
    re2 := regexp.MustCompile(`(?P<first>\w+)\s+(?P<last>\w+)`)
    matches2 := re2.FindStringSubmatch("Jane Smith")
    for i, name := range re2.SubexpNames() {
        if i != 0 && name != "" {
            fmt.Printf("%s: %s\n", name, matches2[i])
        }
    }

    // Find all with submatches
    re3 := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
    text := "Dates: 2024-01-15, 2024-02-20"
    allMatches := re3.FindAllStringSubmatch(text, -1)
    for _, match := range allMatches {
        fmt.Printf("Date: %s-%s-%s\n", match[1], match[2], match[3])
    }

    // Replace with groups
    re4 := regexp.MustCompile(`(\w+)\s+(\w+)`)
    result := re4.ReplaceAllString("John Doe", "$2, $1")
    fmt.Println(result) // Doe, John

    // Replace with function
    re5 := regexp.MustCompile(`\d+`)
    result2 := re5.ReplaceAllStringFunc("Age: 25, Score: 90", func(s string) string {
        return "[" + s + "]"
    })
    fmt.Println(result2) // Age: [25], Score: [90]
}

// Split with regex
func splitExample() {
    re := regexp.MustCompile(`\s+`)
    parts := re.Split("Hello   World\tTest", -1)
    fmt.Println(parts) // [Hello World Test]
}
```

---

## 3. Performance Considerations

```go
// Compile once, use many times
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isEmailValid(email string) bool {
    return emailRegex.MatchString(email)
}

// Avoid compiling in hot paths
func badExample(emails []string) {
    for _, email := range emails {
        // DON'T: Compile regex in loop
        re := regexp.MustCompile(`...`)
        re.MatchString(email)
    }
}

func goodExample(emails []string) {
    // DO: Compile once outside loop
    re := regexp.MustCompile(`...`)
    for _, email := range emails {
        re.MatchString(email)
    }
}
```

---

## 4. Best Practices

```
Regex Best Practices:
□ Compile regex at package level
□ Use raw strings for patterns
□ Test patterns thoroughly
□ Document what pattern matches
□ Consider readability vs complexity
□ Use capture groups wisely
□ Handle no-match cases
□ Escape special characters properly
```

---

## 5. Checklist

```
Regex Checklist:
□ Pattern tested with edge cases
□ Compiled once, reused
□ Proper escaping
□ Capture groups documented
□ Performance acceptable
□ Handles no-match case
□ Uses appropriate methods (Match/Find/Replace)
```
