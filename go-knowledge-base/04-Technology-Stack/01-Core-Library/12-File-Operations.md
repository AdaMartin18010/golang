# TS-CL-012: Go File Operations Deep Dive

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #file-io #os #filesystem #buffering
> **权威来源**:
>
> - [os Package](https://golang.org/pkg/os/) - Go standard library
> - [Effective Go: Data](https://golang.org/doc/effective_go.html#data) - Go team

---

## 1. File System Operations

### 1.1 Basic File Operations

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

// Reading files
func readFile(filename string) ([]byte, error) {
    // Read entire file
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return data, nil
}

// Reading large files efficiently
func readLargeFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Use buffered reader for efficiency
    reader := bufio.NewReaderSize(file, 64*1024) // 64KB buffer

    buf := make([]byte, 4096)
    for {
        n, err := reader.Read(buf)
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        // Process buf[:n]
        process(buf[:n])
    }

    return nil
}

// Writing files
func writeFile(filename string, data []byte) error {
    // Write entire file (truncates existing)
    return os.WriteFile(filename, data, 0644)
}

// Append to file
func appendToFile(filename string, data []byte) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.Write(data)
    return err
}

// Line-by-line reading
func readLines(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    return lines, scanner.Err()
}
```

### 1.2 File Information and Permissions

```go
// File info and operations
func fileOperations(path string) error {
    // Get file info
    info, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("File does not exist")
        }
        return err
    }

    fmt.Printf("Name: %s\n", info.Name())
    fmt.Printf("Size: %d bytes\n", info.Size())
    fmt.Printf("Mode: %s\n", info.Mode())
    fmt.Printf("ModTime: %v\n", info.ModTime())
    fmt.Printf("IsDir: %v\n", info.IsDir())

    // Check permissions
    mode := info.Mode()
    if mode&0400 != 0 {
        fmt.Println("Readable by owner")
    }
    if mode&0200 != 0 {
        fmt.Println("Writable by owner")
    }
    if mode&0100 != 0 {
        fmt.Println("Executable by owner")
    }

    // Change permissions
    if err := os.Chmod(path, 0755); err != nil {
        return err
    }

    // Change owner (Unix only)
    // os.Chown(path, uid, gid)

    return nil
}
```

---

## 2. Directory Operations

```go
// Directory operations
func directoryOperations() error {
    // Create directory
    if err := os.Mkdir("newdir", 0755); err != nil {
        return err
    }

    // Create directory tree
    if err := os.MkdirAll("parent/child/grandchild", 0755); err != nil {
        return err
    }

    // List directory contents
    entries, err := os.ReadDir(".")
    if err != nil {
        return err
    }

    for _, entry := range entries {
        info, err := entry.Info()
        if err != nil {
            continue
        }
        fmt.Printf("%s (size: %d, dir: %v)\n", entry.Name(), info.Size(), entry.IsDir())
    }

    // Walk directory tree
    err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        fmt.Printf("%s (size: %d)\n", path, info.Size())
        return nil
    })

    // Remove directory
    os.Remove("emptydir")           // Remove empty directory
    os.RemoveAll("parent")          // Remove directory and contents

    return err
}
```

---

## 3. Path Manipulation

```go
// Path operations
func pathOperations() {
    // Join paths (cross-platform)
    path := filepath.Join("home", "user", "documents", "file.txt")
    // Windows: home\user\documents\file.txt
    // Unix: home/user/documents/file.txt

    // Split path
    dir := filepath.Dir(path)
    base := filepath.Base(path)
    ext := filepath.Ext(path)

    fmt.Printf("Dir: %s\n", dir)
    fmt.Printf("Base: %s\n", base)
    fmt.Printf("Ext: %s\n", ext)

    // Clean path
    dirty := "/home//user/../user/./documents/file.txt"
    clean := filepath.Clean(dirty)
    fmt.Printf("Clean: %s\n", clean) // /home/user/documents/file.txt

    // Absolute path
    abs, _ := filepath.Abs("relative/path")

    // Relative path
    rel, _ := filepath.Rel("/home/user", "/home/user/documents/file.txt")

    // Check if absolute
    isAbs := filepath.IsAbs("/absolute/path")
}
```

---

## 4. Best Practices

```go
// Best practices for file operations

// 1. Always defer Close()
func readFileProperly(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Process file...
    return nil
}

// 2. Use temporary files for atomic writes
func atomicWrite(filename string, data []byte) error {
    // Create temp file in same directory
    tmpFile, err := os.CreateTemp(filepath.Dir(filename), "*.tmp")
    if err != nil {
        return err
    }
    tmpName := tmpFile.Name()

    // Cleanup on failure
    defer func() {
        if err != nil {
            os.Remove(tmpName)
        }
    }()

    // Write data
    if _, err = tmpFile.Write(data); err != nil {
        tmpFile.Close()
        return err
    }

    if err = tmpFile.Close(); err != nil {
        return err
    }

    // Atomic rename
    return os.Rename(tmpName, filename)
}

// 3. Handle large files with streaming
func copyLargeFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()

    // Copy with buffer
    buf := make([]byte, 1024*1024) // 1MB buffer
    _, err = io.CopyBuffer(destination, source, buf)
    return err
}
```

---

## 5. Checklist

```
File Operations Checklist:
□ Always defer file.Close()
□ Check errors for all operations
□ Use buffered I/O for performance
□ Use atomic writes for critical data
□ Handle large files with streaming
□ Use proper file permissions
□ Cross-platform path handling
□ Clean up temporary files
```
