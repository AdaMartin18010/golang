# TS-CL-006: Go time Package Deep Dive

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #time #timezone #duration #timing
> **权威来源**:
>
> - [time Package](https://golang.org/pkg/time/) - Go standard library
> - [Go Time Formatting](https://go.dev/src/time/format.go) - Source code

---

## 1. Time Representation

### 1.1 Time Structure

```go
// Time represents an instant in time
type Time struct {
    wall uint64    // 1 bit flag + 33 bits seconds + 30 bits nanoseconds
    ext  int64     // Monotonic reading for calculations
    loc  *Location // Timezone cache
}

// Wall time: Seconds since January 1, year 1, 00:00:00 UTC
// Monotonic: Nanoseconds since process start (for comparisons)
```

### 1.2 Creating Time Values

```go
package main

import (
    "fmt"
    "time"
)

func timeCreation() {
    // Current time
    now := time.Now()

    // Specific time
    t := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

    // From Unix timestamp
    unixTime := time.Unix(1705315800, 0)

    // From Unix milliseconds
    unixMilliTime := time.UnixMilli(1705315800000)

    // From Unix nanoseconds
    unixNanoTime := time.Unix(0, 1705315800000000000)

    // Zero time
    var zeroTime time.Time
    fmt.Println(zeroTime.IsZero()) // true
}
```

---

## 2. Time Formatting and Parsing

### 2.1 Format Reference Time

```go
func formatting() {
    t := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

    // Reference time: Mon Jan 2 15:04:05 MST 2006
    // Components:    01   2  3  4  5    6    7

    // Common formats
    fmt.Println(t.Format("2006-01-02"))                    // 2024-01-15
    fmt.Println(t.Format("2006-01-02 15:04:05"))          // 2024-01-15 10:30:00
    fmt.Println(t.Format(time.RFC3339))                    // 2024-01-15T10:30:00Z
    fmt.Println(t.Format(time.RFC3339Nano))                // 2024-01-15T10:30:00Z
    fmt.Println(t.Format(time.RFC1123))                    // Mon, 15 Jan 2024 10:30:00 UTC
    fmt.Println(t.Format(time.Kitchen))                    // 10:30AM
    fmt.Println(t.Format("January 2, 2006"))               // January 15, 2024
    fmt.Println(t.Format("Mon, 02 Jan"))                   // Mon, 15 Jan
    fmt.Println(t.Format("3:04 PM"))                       // 10:30 AM
}

func parsing() {
    // Parse time from string
    layouts := []string{
        time.RFC3339,
        "2006-01-02 15:04:05",
        "2006-01-02",
        "01/02/2006",
        "January 2, 2006",
    }

    input := "2024-01-15 10:30:00"

    for _, layout := range layouts {
        t, err := time.Parse(layout, input)
        if err == nil {
            fmt.Printf("Parsed with %s: %v\n", layout, t)
            break
        }
    }

    // Parse in specific location
    loc, _ := time.LoadLocation("America/New_York")
    t, _ := time.ParseInLocation("2006-01-02 15:04:05", input, loc)
    fmt.Println(t)
}
```

---

## 3. Duration

```go
func durationExamples() {
    // Creating durations
    d1 := 5 * time.Second
    d2 := 2 * time.Minute
    d3 := 1 * time.Hour
    d4 := 24 * time.Hour // 1 day

    // Parsing duration from string
    d, _ := time.ParseDuration("1h30m")
    d, _ = time.ParseDuration("2h45m30s")
    d, _ = time.ParseDuration("100ms")
    d, _ = time.ParseDuration("1.5h")

    // Duration arithmetic
    sum := d1 + d2
    diff := d2 - d1

    // Comparison
    if d1 < d2 {
        fmt.Println("d1 is shorter")
    }

    // Conversions
    nanos := d.Nanoseconds()
    millis := d.Milliseconds()
    seconds := d.Seconds()
    minutes := d.Minutes()
    hours := d.Hours()

    // Truncate and round
    truncated := d.Truncate(time.Hour)
    rounded := d.Round(time.Minute)

    // String representation
    fmt.Println(d.String()) // 1h30m0s
}
```

---

## 4. Time Zones

```go
func timezoneExamples() {
    // Load location
    nyc, _ := time.LoadLocation("America/New_York")
    tokyo, _ := time.LoadLocation("Asia/Tokyo")
    london, _ := time.LoadLocation("Europe/London")

    // Create time in specific timezone
    t := time.Date(2024, 1, 15, 10, 30, 0, 0, nyc)

    // Convert between timezones
    tInTokyo := t.In(tokyo)
    tInLondon := t.In(london)
    tInUTC := t.UTC()

    fmt.Println("NYC:", t)
    fmt.Println("Tokyo:", tInTokyo)
    fmt.Println("London:", tInLondon)
    fmt.Println("UTC:", tInUTC)

    // Location from time
    loc := t.Location()

    // Zone information
    name, offset := t.Zone()
    fmt.Printf("Zone: %s, Offset: %d\n", name, offset)
}
```

---

## 5. Timers and Tickers

```go
func timerExamples() {
    // One-shot timer
    timer := time.NewTimer(2 * time.Second)
    <-timer.C
    fmt.Println("Timer fired!")

    // Stop timer
    timer2 := time.NewTimer(5 * time.Second)
    go func() {
        <-timer2.C
        fmt.Println("This won't print")
    }()
    stopped := timer2.Stop()
    fmt.Println("Stopped:", stopped)

    // Reset timer
    timer3 := time.NewTimer(1 * time.Second)
    timer3.Reset(3 * time.Second)

    // After (convenience function)
    <-time.After(2 * time.Second)
}

func tickerExamples() {
    // Periodic ticker
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    done := make(chan bool)
    go func() {
        for {
            select {
            case t := <-ticker.C:
                fmt.Println("Tick at", t)
            case <-done:
                return
            }
        }
    }()

    time.Sleep(5 * time.Second)
    done <- true

    // Tick (convenience, leaks if not stopped)
    // Don't use: c := time.Tick(1 * time.Second)
}
```

---

## 6. Checklist

```
Time Package Checklist:
□ Use time.Time for instants
□ Use time.Duration for spans
□ Handle timezone properly
□ Use context for cancellation
□ Stop timers to free resources
□ Use monotonic time for comparisons
□ Format with reference time
□ Parse with correct layout
```
