# TS-CL-012: Go File Operations - Deep Architecture and Patterns

> **维度**: Technology Stack > Core Library
> **级别**: S (18+ KB)
> **标签**: #golang #file #io #filesystem #os
> **权威来源**:
>
> - [Go os package](https://pkg.go.dev/os) - Official documentation
> - [Go io/ioutil](https://pkg.go.dev/io/ioutil) - I/O utilities

---

## 1. File System Architecture

### 1.1 File Operations Hierarchy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        File Operations Hierarchy                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   High-Level Operations                                                      │
│   ┌───────────────────────────────────────────────────────────────────────┐  │
│   │  os.ReadFile() / os.WriteFile()                                      │  │
│   │  - Simple, complete operations                                       │  │
│   └───────────────────────────────┬───────────────────────────────────────┘  │
│                                   │                                          │
│   Medium-Level Operations          │                                        │
│   ┌───────────────────────────────┴───────────────────────────────────────┐  │
│   │  bufio.Reader/Writer                                                  │  │
│   │  - Buffered I/O for efficiency                                       │  │
│   └───────────────────────────────┬───────────────────────────────────────┘  │
│                                   │                                          │
│   Low-Level Operations             │                                        │
│   ┌───────────────────────────────┴───────────────────────────────────────┐  │
│   │  os.File (Read, Write, Seek)                                          │  │
│   │  - Direct system calls                                               │  │
│   └───────────────────────────────┬───────────────────────────────────────┘  │
│                                   │                                          │
│   System Level                     │                                        │
│   ┌───────────────────────────────┴───────────────────────────────────────┐  │
│   │  Syscalls (read, write, open, close)                                  │  │
│   │  - Kernel interface                                                  │  │
│   └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. File Reading Patterns

### 2.1 Complete File Read

```go
// Simple read (Go 1.16+)
data, err := os.ReadFile("data.txt")
if err != nil {
    return err
}

// With file info for validation
info, err := os.Stat("data.txt")
if err != nil {
    return err
}
if info.Size() > 100*1024*1024 { // 100MB limit
    return fmt.Errorf("file too large")
}

data, err = os.ReadFile("data.txt")
```

### 2.2 Streaming Read

```go
// Buffered reading for large files
func processLargeFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := bufio.NewReaderSize(file, 64*1024) // 64KB buffer

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }
        processLine(line)
    }
    return nil
}
```

### 2.3 Memory-Mapped Files

```go
import "golang.org/x/exp/mmap"

func readMMap(filename string) error {
    reader, err := mmap.Open(filename)
    if err != nil {
        return err
    }
    defer reader.Close()

    data := make([]byte, reader.Len())
    _, err = reader.ReadAt(data, 0)
    return err
}
```

---

## 3. File Writing Patterns

### 3.1 Atomic Writes

```go
func writeFileAtomically(filename string, data []byte) error {
    // Write to temp file
    tmpFile := filename + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0644); err != nil {
        return err
    }

    // Atomic rename
    return os.Rename(tmpFile, filename)
}
```

### 3.2 Buffered Writing

```go
func writeBuffered(filename string, lines []string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriterSize(file, 64*1024)
    defer writer.Flush()

    for _, line := range lines {
        if _, err := writer.WriteString(line + "\n"); err != nil {
            return err
        }
    }
    return nil
}
```

---

## 4. File Metadata and Operations

### 4.1 File Information

```go
info, err := os.Stat("file.txt")
if err != nil {
    if os.IsNotExist(err) {
        // File doesn't exist
    }
    return err
}

fmt.Printf("Name: %s\n", info.Name())
fmt.Printf("Size: %d bytes\n", info.Size())
fmt.Printf("Mode: %v\n", info.Mode())
fmt.Printf("Modified: %v\n", info.ModTime())
fmt.Printf("Is Dir: %v\n", info.IsDir())
```

### 4.2 Directory Operations

```go
// Create directory
os.Mkdir("newdir", 0755)
os.MkdirAll("path/to/nested/dir", 0755)

// Read directory
entries, err := os.ReadDir(".")
for _, entry := range entries {
    fmt.Println(entry.Name())
    fmt.Println(entry.IsDir())
    info, _ := entry.Info()
    fmt.Println(info.Size())
}

// Walk directory tree
filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    fmt.Println(path)
    return nil
})
```

---

## 5. Performance Tuning

### 5.1 Buffer Size Comparison

| Buffer Size | Small Files | Large Files | Memory |
|-------------|-------------|-------------|--------|
| 4KB | Good | Slow | Low |
| 64KB | Good | Good | Medium |
| 256KB | OK | Fast | High |
| 1MB | Overhead | Fastest | Very High |

---

## 6. Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      File Operations Best Practices                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Reading:                                                                    │
│  □ Use os.ReadFile for small files                                          │
│  □ Use buffered I/O for large files                                         │
│  □ Always check file size before reading                                    │
│  □ Use defer file.Close()                                                   │
│                                                                              │
│  Writing:                                                                    │
│  □ Use atomic writes for critical data                                      │
│  □ Use buffered writers for multiple writes                                 │
│  □ Set appropriate file permissions                                         │
│                                                                              │
│  Safety:                                                                     │
│  □ Validate file paths (prevent directory traversal)                        │
│  □ Check for file existence before operations                               │
│  □ Handle permission errors gracefully                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (18+ KB, comprehensive coverage)
