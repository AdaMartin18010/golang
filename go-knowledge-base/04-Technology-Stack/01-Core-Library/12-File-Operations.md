# 文件操作 (File Operations)

> **分类**: 开源技术堆栈

---

## 基本操作

### 读写文件

```go
// 读取整个文件
data, err := os.ReadFile("config.json")
if err != nil {
    log.Fatal(err)
}

// 写入文件
err = os.WriteFile("output.txt", []byte("hello"), 0644)

// 追加模式
f, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
defer f.Close()
f.WriteString("new line\n")
```

---

## 遍历目录

```go
// 遍历目录树
err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    
    fmt.Printf("%s (size: %d)\n", path, info.Size())
    return nil
})

// Go 1.16+ 的 WalkDir (更高效)
err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }
    
    info, _ := d.Info()
    fmt.Printf("%s (size: %d)\n", path, info.Size())
    return nil
})
```

---

## 文件监控

```go
import "github.com/fsnotify/fsnotify"

watcher, err := fsnotify.NewWatcher()
if err != nil {
    log.Fatal(err)
}
defer watcher.Close()

go func() {
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            
            if event.Op&fsnotify.Write == fsnotify.Write {
                log.Println("Modified file:", event.Name)
            }
            
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Println("Error:", err)
        }
    }
}()

// 添加监控路径
err = watcher.Add("/tmp/foo")
```

---

## 内存映射

```go
import "golang.org/x/exp/mmap"

// 只读内存映射
reader, err := mmap.Open("largefile.dat")
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

// 读取数据
data := make([]byte, 1024)
n, err := reader.ReadAt(data, 0)
```

---

## 临时文件

```go
// 创建临时文件
tmpFile, err := os.CreateTemp("", "example-*.txt")
if err != nil {
    log.Fatal(err)
}
defer os.Remove(tmpFile.Name())  // 清理
defer tmpFile.Close()

// 创建临时目录
tmpDir, err := os.MkdirTemp("", "example-*")
if err != nil {
    log.Fatal(err)
}
defer os.RemoveAll(tmpDir)
```

---

## 文件锁定

```go
import "github.com/gofrs/flock"

// 创建文件锁
fileLock := flock.New("/var/lock/my.lock")

// 获取锁
locked, err := fileLock.TryLock()
if err != nil {
    log.Fatal(err)
}

if !locked {
    log.Println("Could not acquire lock")
    return
}
defer fileLock.Unlock()

// 执行需要互斥的操作
```
