# A1: 代码模式

> **层级**: 应用层 (Application)
> **地位**: 基于形式化理论的工程模式
> **依赖**: D3 (定理证明)

---

## 模式 P1: 可选字段 (基于 Th1.1)

### 理论基础

```
由 Th1.1: new(v) ≡ &v'

应用:
  可选字段需要: *T 类型的指针
  new(v) 提供: 简洁的 *T 创建方式
```

### 实现

```go
// 可选配置结构
type Config struct {
    Timeout  *time.Duration
    MaxConns *int
    Retries  *int
}

// 使用 new 表达式创建
func NewConfig() Config {
    return Config{
        Timeout:  new(30 * time.Second),
        MaxConns: new(100),
        Retries:  new(3),
    }
}

// 应用默认值
func (c *Config) ApplyDefaults() {
    if c.Timeout == nil {
        c.Timeout = new(30 * time.Second)
    }
    if c.MaxConns == nil {
        c.MaxConns = new(100)
    }
    if c.Retries == nil {
        c.Retries = new(3)
    }
}
```

### 正确性论证

```
正确性: 基于 Th1.1

new(30 * time.Second) 创建指向 30s 的指针
&time.Duration(30 * time.Second) 同样创建指向 30s 的指针
两者语义等价，但 new 更简洁
```

---

## 模式 P2: 构造者模式 (基于 Th1.1)

### 实现

```go
type ServerConfig struct {
    Host *string
    Port *int
    TLS  *bool
}

type ServerConfigBuilder struct {
    config ServerConfig
}

func NewServerConfigBuilder() *ServerConfigBuilder {
    return &ServerConfigBuilder{}
}

func (b *ServerConfigBuilder) WithHost(host string) *ServerConfigBuilder {
    b.config.Host = new(host)  // 简洁
    return b
}

func (b *ServerConfigBuilder) WithPort(port int) *ServerConfigBuilder {
    b.config.Port = new(port)
    return b
}

func (b *ServerConfigBuilder) Build() ServerConfig {
    return b.config
}

// 使用
config := NewServerConfigBuilder().
    WithHost("localhost").
    WithPort(8080).
    Build()
```

---

## 模式 P3: 树遍历 (基于 Th1.2)

### 理论基础

```
由 Th1.2: wellformed(C) → terminates(unfold(C))

应用:
  定义递归约束 Node[T Node[T]]
  保证树遍历算法终止
```

### 实现

```go
// 通用树节点约束
type Node[T Node[T]] interface {
    Children() []T
}

// 通用遍历 - 保证终止 (由Th1.2)
func Walk[T Node[T]](node T, fn func(T)) {
    fn(node)
    for _, child := range node.Children() {
        Walk(child, fn)  // 递归调用，类型安全
    }
}

// 具体实现
type FileNode struct {
    Name     string
    Children []*FileNode
}

func (f *FileNode) Children() []*FileNode {
    return f.Children
}

// 使用
func main() {
    root := buildFileTree("/home/user")

    // 遍历所有文件
    Walk(root, func(node *FileNode) {
        if !strings.HasSuffix(node.Name, ".go") {
            return
        }
        fmt.Println(node.Name)
    })
}
```

### 终止性保证

```
由 Th1.2，Walk 函数保证终止:
  1. FileNode 满足结构递归（子节点 < 父节点）
  2. 每次递归调用处理子节点
  3. 有限树深度保证有限递归深度
```

---

## 模式 P4: 递归算法抽象

```go
// 可比较接口
type Comparable[T Comparable[T]] interface {
    Compare(T) int
    LessThan(T) bool
}

// 通用二叉搜索树
type BST[T Comparable[T]] struct {
    Value T
    Left  *BST[T]
    Right *BST[T]
}

func (t *BST[T]) Insert(v T) *BST[T] {
    if t == nil {
        return &BST[T]{Value: v}
    }
    if v.LessThan(t.Value) {
        t.Left = t.Left.Insert(v)
    } else {
        t.Right = t.Right.Insert(v)
    }
    return t
}

func (t *BST[T]) Search(v T) (T, bool) {
    var zero T
    if t == nil {
        return zero, false
    }
    cmp := v.Compare(t.Value)
    if cmp == 0 {
        return t.Value, true
    }
    if cmp < 0 {
        return t.Left.Search(v)
    }
    return t.Right.Search(v)
}
```

---

**下一章**: [A2-工程实践](A2-practices.md)
