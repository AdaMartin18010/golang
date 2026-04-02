# TS-CL-002: Go I/O Package - Deep Architecture and Patterns

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #io #streaming #interfaces #zero-copy #buffering

---

## 1. I/O Architecture Fundamentals

### 1.1 The Universal Interface Philosophy

The `io` package defines the fundamental abstractions that power Go's composable I/O ecosystem.

**Core Principle:** Every I/O source implements `io.Reader`. Every I/O destination implements `io.Writer`. This enables universal composability.

### 1.2 Core Interfaces

```go
// Reader is the interface that wraps the basic Read method.
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer is the interface that wraps the basic Write method.
type Writer interface {
    Write(p []byte) (n int, err error)
}

// Closer - Resource cleanup
type Closer interface {
    Close() error
}

// Seeker - Random access
type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}
```

**EOF Handling Pattern:**

```go
// CORRECT EOF handling pattern
for {
    n, err := r.Read(buf)
    if n > 0 {
        process(buf[:n])
    }
    if err != nil {
        if err == io.EOF {
            break
        }
        return err
    }
}
```

---

## 2. Buffered I/O Architecture

### 2.1 Why Buffering Matters

System Call Overhead Comparison:

- Unbuffered: 1000 reads = 1000 syscalls
- Buffered (64KB): 1000 x 1KB reads = ~16 syscalls
- Typical speedup: 10-100x for small reads

### 2.2 bufio.Reader Implementation

```go
type Reader struct {
    buf          []byte
    rd           io.Reader
    r, w         int
    err          error
}
```

**Usage:**

```go
file, _ := os.Open("data.txt")
defer file.Close()

reader := bufio.NewReaderSize(file, 64*1024) // 64KB buffer

line, err := reader.ReadString('\n')
peek, err := reader.Peek(10)
reader.Discard(100)
```

### 2.3 bufio.Scanner - Token-Based Reading

```go
// Efficient for large files - minimal memory allocation
scanner := bufio.NewScanner(file)
scanner.Buffer(make([]byte, 4096), 1024*1024) // Custom buffer, max token size

for scanner.Scan() {
    line := scanner.Text() // Returns string (copy)
    bytes := scanner.Bytes() // Returns slice (shares buffer)
}

if err := scanner.Err(); err != nil {
    return err
}
```

**Custom Split Functions:**

```go
// Split by words
scanner.Split(bufio.ScanWords)

// Split by bytes
scanner.Split(bufio.ScanBytes)

// Custom split: records separated by "\n\n"
onDoubleNewline := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
    if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
        return i + 2, data[:i], nil
    }
    if atEOF && len(data) > 0 {
        return len(data), data, nil
    }
    return 0, nil, nil
}
scanner.Split(onDoubleNewline)
```

---

## 3. Utility Functions and Patterns

### 3.1 io.Copy - Zero-Copy Transfer

```go
// Copy uses buffer pool and may use sendfile/splice
func Copy(dst io.Writer, src io.Reader) (written int64, err error)

// Copy with custom buffer
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error)

// Copy at most N bytes
func CopyN(dst io.Writer, src io.Reader, n int64) (written int64, err error)
```

**Internal Optimization:**

```go
// io.Copy checks for specific interfaces to optimize:
// - src.(WriterTo) - src writes directly to dst
// - dst.(ReaderFrom) - dst reads directly from src
// - Both os.File - may use sendfile/splice syscall

// Example: sendfile optimization
copyBuffer(dst, src, buf) {
    if wt, ok := src.(WriterTo); ok {
        return wt.WriteTo(dst)  // Optimized path
    }
    if rt, ok := dst.(ReaderFrom); ok {
        return rt.ReadFrom(src) // Optimized path
    }
    // Fallback to buffered copy
}
```

### 3.2 Pipe - Synchronous Communication

```go
// io.Pipe creates a connected reader/writer pair
pr, pw := io.Pipe()

// Writer blocks until reader reads
go func() {
    pw.Write([]byte("data"))
    pw.Close()
}()

data, _ := io.ReadAll(pr)
```

### 3.3 MultiReader and MultiWriter

```go
// Combine multiple readers sequentially
combined := io.MultiReader(reader1, reader2, reader3)
data, _ := io.ReadAll(combined) // Reads all in sequence

// Broadcast to multiple writers
destinations := []io.Writer{file1, file2, file3}
multi := io.MultiWriter(destinations...)
fmt.Fprint(multi, "broadcast") // Writes to all
```

### 3.4 TeeReader

```go
// Reads from src and writes to dst simultaneously
// Useful for calculating checksums while reading

hash := sha256.New()
tee := io.TeeReader(file, hash)

// Data flows to both file processing and hash calculation
processFile(tee)
checksum := hash.Sum(nil)
```

---

## 4. Advanced Patterns

### 4.1 Reader/Writer Wrappers

```go
// CountingReader - track bytes read
type CountingReader struct {
    Reader io.Reader
    N      int64
}

func (c *CountingReader) Read(p []byte) (n int, err error) {
    n, err = c.Reader.Read(p)
    c.N += int64(n)
    return
}

// LimitingReader - restrict read amount (prevents OOM)
type LimitingReader struct {
    R io.Reader
    N int64 // max bytes remaining
}

func (l *LimitingReader) Read(p []byte) (n int, err error) {
    if l.N <= 0 {
        return 0, io.EOF
    }
    if int64(len(p)) > l.N {
        p = p[:l.N]
    }
    n, err = l.R.Read(p)
    l.N -= int64(n)
    return
}
```

### 4.2 SectionReader - Random Access

```go
// Access a section of a file without seeking
section := io.NewSectionReader(file, offset, length)

// Implements: io.Reader, io.Seeker, io.ReaderAt
// Safe for concurrent use
```

---

## 5. Performance Benchmarks

```go
// Benchmark: Buffered vs Unbuffered

func BenchmarkUnbufferedRead(b *testing.B) {
    data := make([]byte, 1024)
    for i := 0; i < b.N; i++ {
        file, _ := os.Open("test.txt")
        for {
            _, err := file.Read(data)
            if err == io.EOF {
                break
            }
        }
        file.Close()
    }
}

func BenchmarkBufferedRead(b *testing.B) {
    data := make([]byte, 1024)
    for i := 0; i < b.N; i++ {
        file, _ := os.Open("test.txt")
        reader := bufio.NewReaderSize(file, 64*1024)
        for {
            _, err := reader.Read(data)
            if err == io.EOF {
                break
            }
        }
        file.Close()
    }
}

// Results (typical):
// Unbuffered: ~50 MB/s, many syscalls
// Buffered:   ~500 MB/s, fewer syscalls
```

---

## 6. Best Practices

```
Always use bufio for small reads/writes
Always defer reader/writer Close()
Always handle io.EOF correctly
Use io.Copy for large transfers
Use scanner for line-by-line processing
Use pipes for goroutine communication
```
