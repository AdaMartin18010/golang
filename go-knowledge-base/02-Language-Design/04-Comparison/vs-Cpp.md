# Go vs C++ 对比

> **分类**: 语言设计

---

## 概览

| 特性 | Go | C++ |
|------|-----|-----|
| 发布时间 | 2009 | 1985 |
| 抽象级别 | 高级 | 低到高级 |
| 内存管理 | GC | 手动/Smart Pointer |
| 编译 | 快 | 慢 |
| 运行时 | 需要 | 最小化 |
| 性能 | 好 | 极致 |
| 安全性 | 安全 | 需小心 |

---

## 代码对比

### 类/结构体

```go
// Go
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}
```

```cpp
// C++
class Rectangle {
public:
    double width, height;
    double area() const {
        return width * height;
    }
};
```

### 内存管理

```go
// Go: GC 管理
func process() {
    data := make([]byte, 1024)
    // 自动释放
}
```

```cpp
// C++: 手动管理
void process() {
    auto data = std::make_unique<std::byte[]>(1024);
    // RAII 自动释放
}
```

---

## 并发

### Go: CSP 模型

```go
ch := make(chan int)

go func() {
    ch <- 42
}()

v := <-ch
```

### C++: 线程 + 原子

```cpp
std::queue<int> q;
std::mutex m;
std::condition_variable cv;

// 生产者
std::thread([&] {
    std::lock_guard<std::mutex> lock(m);
    q.push(42);
    cv.notify_one();
}).detach();

// 消费者
std::unique_lock<std::mutex> lock(m);
cv.wait(lock, [&] { return !q.empty(); });
int v = q.front();
```

---

## 编译构建

| 特性 | Go | C++ |
|------|-----|-----|
| 编译速度 | 秒级 | 分钟级 |
| 依赖管理 | go mod | CMake/vcpkg/conan |
| 构建系统 | 内置 | 复杂 |
| 交叉编译 | 简单 | 复杂 |

---

## 适用场景

### 选择 Go

- 网络服务
- 云原生
- 快速迭代
- 大型团队协作
- DevOps 工具

### 选择 C++

- 游戏开发
- 嵌入式
- 高频交易
- 操作系统
- 性能关键

---

## 总结

Go 和 C++ 在抽象层次上差异巨大：

- **Go**: 现代、安全、高效开发
- **C++**: 底层控制、极致性能

两者互补，各有最佳适用领域。
