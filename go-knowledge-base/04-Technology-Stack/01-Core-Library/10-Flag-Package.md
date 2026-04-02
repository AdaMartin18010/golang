# TS-CL-010: Go flag Package

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #cli #flags #command-line #arguments
> **权威来源**:
>
> - [flag Package](https://golang.org/pkg/flag/) - Go standard library

---

## 1. Basic Flag Usage

```go
package main

import (
    "flag"
    "fmt"
    "os"
    "time"
)

func basicFlags() {
    // Define flags
    var (
        host    = flag.String("host", "localhost", "Server host")
        port    = flag.Int("port", 8080, "Server port")
        debug   = flag.Bool("debug", false, "Enable debug mode")
        timeout = flag.Duration("timeout", 30*time.Second, "Request timeout")
    )

    // Parse flags
    flag.Parse()

    // Use flags
    fmt.Printf("Server: %s:%d\n", *host, *port)
    fmt.Printf("Debug: %v\n", *debug)
    fmt.Printf("Timeout: %v\n", *timeout)
}

// Custom flag type
type arrayFlags []string

func (a *arrayFlags) String() string {
    return fmt.Sprintf("%v", *a)
}

func (a *arrayFlags) Set(value string) error {
    *a = append(*a, value)
    return nil
}

func customFlagType() {
    var tags arrayFlags
    flag.Var(&tags, "tag", "Tags (can be specified multiple times)")

    flag.Parse()

    fmt.Printf("Tags: %v\n", tags)
    // Usage: ./app -tag=go -tag=backend -tag=api
}
```

---

## 2. Advanced Flag Usage

```go
func advancedFlags() {
    // Flag with environment variable fallback
    host := flag.String("host", getEnv("HOST", "localhost"), "Server host")

    flag.Parse()

    fmt.Printf("Host: %s\n", *host)
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// Subcommands with flag.NewFlagSet
func subcommands() {
    // Define subcommands
    serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
    serveHost := serveCmd.String("host", "localhost", "Server host")
    servePort := serveCmd.Int("port", 8080, "Server port")

    migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
    migrateDirection := migrateCmd.String("direction", "up", "Migration direction (up/down)")

    // Check subcommand
    if len(os.Args) < 2 {
        fmt.Println("Expected 'serve' or 'migrate' subcommand")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "serve":
        serveCmd.Parse(os.Args[2:])
        fmt.Printf("Serving on %s:%d\n", *serveHost, *servePort)

    case "migrate":
        migrateCmd.Parse(os.Args[2:])
        fmt.Printf("Migrating %s\n", *migrateDirection)

    default:
        fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
        os.Exit(1)
    }
}
```

---

## 3. Best Practices

```go
// Configuration struct with flags
type Config struct {
    Host    string
    Port    int
    Debug   bool
    Timeout time.Duration
}

func parseFlags() Config {
    var cfg Config

    flag.StringVar(&cfg.Host, "host", "localhost", "Server host")
    flag.IntVar(&cfg.Port, "port", 8080, "Server port")
    flag.BoolVar(&cfg.Debug, "debug", false, "Enable debug mode")
    flag.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "Request timeout")

    flag.Parse()

    return cfg
}

// Help text customization
func customHelp() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
        fmt.Fprintln(os.Stderr, "My Application - A sample Go application")
        fmt.Fprintln(os.Stderr, "\nOptions:")
        flag.PrintDefaults()
        fmt.Fprintln(os.Stderr, "\nExamples:")
        fmt.Fprintln(os.Stderr, "  ./app -host=0.0.0.0 -port=9090")
        fmt.Fprintln(os.Stderr, "  ./app -debug -timeout=1m")
    }

    flag.Parse()
}
```

---

## 4. Checklist

```
Flag Package Checklist:
□ Meaningful flag names
□ Sensible defaults
□ Clear help text
□ Environment variable fallback
□ Required flags validated
□ Custom flag types documented
□ Subcommands if needed
□ Version flag
```
