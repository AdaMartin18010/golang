# Go vs Java 对比

> **分类**: 语言设计

---

## 概览

| 特性 | Go | Java |
|------|-----|------|
| 发布时间 | 2009 | 1995 |
| 运行时 | 编译为机器码 + 轻量运行时 | JVM 字节码 + 重型运行时 |
| 内存管理 | GC | GC (更复杂) |
| 并发 | Goroutine | 线程/虚拟线程 |
| 编译速度 | 快 | 慢 |
| 启动速度 | 快 | 慢 |
| 生态 | 现代云原生 | 成熟企业级 |

---

## 代码对比

### Hello World

```go
// Go
package main
import "fmt"
func main() {
    fmt.Println("Hello")
}
```

```java
// Java
public class Hello {
    public static void main(String[] args) {
        System.out.println("Hello");
    }
}
```

### 并发

```go
// Go: 轻量 goroutine
func main() {
    go func() {
        fmt.Println("async")
    }()
    time.Sleep(time.Second)
}
```

```java
// Java 21: 虚拟线程
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
    executor.submit(() -> System.out.println("async"));
}
```

---

## 类型系统

### Go: 结构子类型

```go
type Reader interface { Read() }

type File struct{}
func (f File) Read() {}  // 隐式实现
```

### Java: 名义子类型

```java
interface Reader { void read(); }

class File implements Reader {  // 显式实现
    public void read() {}
}
```

---

## 生态对比

| 领域 | Go | Java |
|------|-----|------|
| Web 框架 | Gin, Echo | Spring Boot |
| ORM | GORM, Ent | Hibernate |
| 微服务 | 标准库为主 | Spring Cloud |
| 云原生 | 主导 | 追赶中 |
| 企业级 | 增长中 | 成熟 |

---

## 适用场景

### 选择 Go

- 微服务
- 云原生
- CLI 工具
- 网络服务
- 容器/Docker/K8s

### 选择 Java

- 大型企业应用
- Android
- 大数据
- 长期维护项目
- 成熟生态系统

---

## 趋势

- **Go**: 云原生时代主流选择
- **Java**: 虚拟线程 (Project Loom) 缩小差距
- **交集**: 两者都在进化，互相学习
