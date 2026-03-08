# Go 1.26 语言变化详解

> 深入解析 Go 1.26 的两大语言层面改进：`new()` 表达式支持和递归泛型约束。

---

## 目录

- [Go 1.26 语言变化详解](#go-126-语言变化详解)
  - [目录](#目录)
  - [1. new() 支持表达式操作数](#1-new-支持表达式操作数)
    - [1.1 变化概述](#11-变化概述)
    - [1.2 语法形式](#12-语法形式)
    - [1.3 语义等价性](#13-语义等价性)
    - [1.4 实际应用场景](#14-实际应用场景)
      - [场景1: 可选字段序列化](#场景1-可选字段序列化)
      - [场景2: API 请求构建](#场景2-api-请求构建)
      - [场景3: 配置对象初始化](#场景3-配置对象初始化)
    - [1.5 反例与注意事项](#15-反例与注意事项)
    - [1.6 性能考量](#16-性能考量)
  - [2. 递归泛型约束](#2-递归泛型约束)
    - [2.1 变化概述](#21-变化概述)
    - [2.2 理论基础](#22-理论基础)
    - [2.3 核心应用场景](#23-核心应用场景)
      - [场景1: 通用比较接口](#场景1-通用比较接口)
      - [场景2: 通用树遍历](#场景2-通用树遍历)
      - [场景3: 图算法抽象](#场景3-图算法抽象)
    - [2.4 约束规则详解](#24-约束规则详解)
    - [2.5 类型推断改进](#25-类型推断改进)
  - [3. 形式化语义分析](#3-形式化语义分析)
    - [3.1 new() 的形式化语义](#31-new-的形式化语义)
    - [3.2 递归约束的形式化语义](#32-递归约束的形式化语义)
  - [4. 最佳实践](#4-最佳实践)
    - [4.1 new() 使用指南](#41-new-使用指南)
    - [4.2 递归约束使用指南](#42-递归约束使用指南)
    - [4.3 迁移建议](#43-迁移建议)
  - [参考](#参考)

---

## 1. new() 支持表达式操作数

### 1.1 变化概述

Go 1.26 放宽了内置 `new` 函数的限制，允许其操作数为**表达式**，而不仅仅是类型名。

```go
// Go 1.25 及之前 - 操作数必须是类型
ptr := new(int)           // ✓ 合法
// ptr := new(int(42))    // ✗ 非法

// Go 1.26 - 操作数可以是表达式
ptr := new(int(42))       // ✓ 合法！
ptr := new(MyStruct{...}) // ✓ 合法！
```

### 1.2 语法形式

```ebnf
// 新的语法（Go 1.26）
NewExpr = "new" Expression

// 表达式可以是：
// - 类型转换: new(T(value))
// - 复合字面量: new(T{fields...})
// - 函数调用: new(funcReturningT())
```

### 1.3 语义等价性

```go
// 以下三种形式在语义上等价：

// 形式1: 传统方式
x := MyType{Field: value}
ptr1 := &x

// 形式2: 复合字面量直接取地址
ptr2 := &MyType{Field: value}

// 形式3: Go 1.26 new() 表达式
ptr3 := new(MyType{Field: value})
```

**等价定理**:

```text
new(T{...}) ≡ &T{...}
new(T(v)) ≡ &T(v) ≡ &v where v is of type T
```

### 1.4 实际应用场景

#### 场景1: 可选字段序列化

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
)

type Person struct {
    Name    string  `json:"name"`
    Age     *int    `json:"age,omitempty"`     // 可选字段
    Email   *string `json:"email,omitempty"`   // 可选字段
}

func yearsSince(t time.Time) int {
    return int(time.Since(t).Hours() / (365.25 * 24))
}

// Go 1.25 方式 - 需要多行代码
func personJSONOld(name string, born time.Time, email string) ([]byte, error) {
    p := Person{Name: name}
    if !born.IsZero() {
        age := yearsSince(born)
        p.Age = &age  // 必须声明中间变量
    }
    if email != "" {
        p.Email = &email  // 可以复用参数
    }
    return json.Marshal(p)
}

// Go 1.26 方式 - 简洁直观
func personJSONNew(name string, born time.Time, email string) ([]byte, error) {
    return json.Marshal(Person{
        Name:  name,
        Age:   new(yearsSince(born)),  // 直接内联！
        Email: new(email),             // 简洁！
    })
}

func main() {
    data, _ := personJSONNew("Alice", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), "alice@example.com")
    fmt.Println(string(data))
    // {"name":"Alice","age":26,"email":"alice@example.com"}
}
```

#### 场景2: API 请求构建

```go
type CreateUserRequest struct {
    Name     string
    Email    string
    Phone    *string  // 可选
    Nickname *string  // 可选
}

// 构建请求时简洁处理可选字段
func buildRequest(name, email string, opts ...string) *CreateUserRequest {
    req := &CreateUserRequest{
        Name:  name,
        Email: email,
    }

    if len(opts) > 0 {
        req.Phone = new(opts[0])  // Go 1.26 简洁语法
    }
    if len(opts) > 1 {
        req.Nickname = new(opts[1])
    }

    return req
}
```

#### 场景3: 配置对象初始化

```go
type Config struct {
    Host string
    Port int
    TLS  *TLSConfig  // 可选
}

type TLSConfig struct {
    Cert string
    Key  string
}

// Go 1.26 快速创建带嵌套可选配置的对象
func createConfig(host string, port int, useTLS bool) *Config {
    cfg := &Config{
        Host: host,
        Port: port,
    }

    if useTLS {
        cfg.TLS = new(TLSConfig{  // 一行初始化
            Cert: "/path/to/cert",
            Key:  "/path/to/key",
        })
    }

    return cfg
}
```

### 1.5 反例与注意事项

```go
// ❌ 错误: 表达式必须返回可寻址的值类型
ptr := new(someFunc())  // 如果 someFunc() 返回非指针值，合法
                        // 但如果返回指针，可能不是你想要的

// ✅ 正确: 明确意图
result := someFunc()
ptr := new(*result)     // 如果你想创建指针的指针

// ❌ 错误: new 不能用于 nil
// ptr := new(nil)      // 编译错误

// ❌ 错误: 不能用于接口类型（除非是类型断言）
var i interface{} = 42
// ptr := new(i)        // 编译错误
ptr := new(i.(int))     // ✅ 正确: 类型断言后使用
```

### 1.6 性能考量

```go
// 三种方式的汇编输出几乎相同
// 编译器会将它们优化为相同的机器码

// &T{...} 和 new(T{...}) 性能完全一致
// 选择取决于代码可读性
```

---

## 2. 递归泛型约束

### 2.1 变化概述

Go 1.26 移除了"泛型类型不能在类型参数列表中引用自身"的限制，允许定义**自引用类型约束**。

```go
// Go 1.25 及之前 - 编译错误
// type Adder[A Adder[A]] interface {
//     ^^^^^ 循环类型约束错误
//     Add(A) A
// }

// Go 1.26 - 合法！
type Adder[A Adder[A]] interface {
    Add(A) A
}
```

### 2.2 理论基础

递归类型约束对应于类型论中的**最小不动点 (Least Fixed Point)**：

```text
类型方程: F(X) = Base ∪ { Op: X → X }
解: μX.F(X)  (最小不动点)

其中 μ 是最小不动点算子，满足:
1. F(μX.F(X)) = μX.F(X)  (不动点性质)
2. ∀Y. F(Y) ⊆ Y ⟹ μX.F(X) ⊆ Y  (最小性)
```

### 2.3 核心应用场景

#### 场景1: 通用比较接口

```go
package main

import "fmt"

// 自引用约束: Ordered 类型的元素必须是 Ordered
// 这允许我们在约束中使用类型自身的方法
type Ordered[T Ordered[T]] interface {
    comparable
    Less(other T) bool
    Equal(other T) bool
}

// 二叉树节点实现 Ordered
type IntNode struct {
    value int
}

func (n IntNode) Less(other IntNode) bool {
    return n.value < other.value
}

func (n IntNode) Equal(other IntNode) bool {
    return n.value == other.value
}

// 通用二叉搜索树
type BST[T Ordered[T]] struct {
    root *BSTNode[T]
}

type BSTNode[T Ordered[T]] struct {
    value T
    left  *BSTNode[T]
    right *BSTNode[T]
}

func (t *BST[T]) Insert(v T) {
    t.root = insertNode(t.root, v)
}

func insertNode[T Ordered[T]](n *BSTNode[T], v T) *BSTNode[T] {
    if n == nil {
        return &BSTNode[T]{value: v}
    }
    if v.Less(n.value) {
        n.left = insertNode(n.left, v)
    } else if n.value.Less(v) {
        n.right = insertNode(n.right, v)
    }
    return n
}

func (t *BST[T]) Search(v T) bool {
    return searchNode(t.root, v)
}

func searchNode[T Ordered[T]](n *BSTNode[T], v T) bool {
    if n == nil {
        return false
    }
    if v.Equal(n.value) {
        return true
    }
    if v.Less(n.value) {
        return searchNode(n.left, v)
    }
    return searchNode(n.right, v)
}

func main() {
    tree := &BST[IntNode]{}
    tree.Insert(IntNode{5})
    tree.Insert(IntNode{3})
    tree.Insert(IntNode{7})

    fmt.Println(tree.Search(IntNode{3}))  // true
    fmt.Println(tree.Search(IntNode{10})) // false
}
```

#### 场景2: 通用树遍历

```go
package main

import "fmt"

// 树节点接口 - 自引用约束
type TreeNode[T TreeNode[T]] interface {
    Value() int
    Children() []T  // 返回相同类型的子节点
}

// 二叉树实现
type BinaryTree struct {
    val   int
    left  *BinaryTree
    right *BinaryTree
}

func (b *BinaryTree) Value() int {
    return b.val
}

func (b *BinaryTree) Children() []*BinaryTree {
    children := make([]*BinaryTree, 0, 2)
    if b.left != nil {
        children = append(children, b.left)
    }
    if b.right != nil {
        children = append(children, b.right)
    }
    return children
}

// 多叉树实现
type NaryTree struct {
    val      int
    children []*NaryTree
}

func (n *NaryTree) Value() int {
    return n.val
}

func (n *NaryTree) Children() []*NaryTree {
    return n.children
}

// 通用遍历函数 - 适用于任何实现 TreeNode 的类型
func PreOrder[T TreeNode[T]](root T, visit func(int)) {
    if root == nil {
        return
    }
    visit(root.Value())
    for _, child := range root.Children() {
        PreOrder(child, visit)
    }
}

func PostOrder[T TreeNode[T]](root T, visit func(int)) {
    if root == nil {
        return
    }
    for _, child := range root.Children() {
        PostOrder(child, visit)
    }
    visit(root.Value())
}

func main() {
    // 测试二叉树
    binary := &BinaryTree{
        val: 1,
        left: &BinaryTree{
            val:   2,
            left:  &BinaryTree{val: 4},
            right: &BinaryTree{val: 5},
        },
        right: &BinaryTree{
            val: 3,
        },
    }

    fmt.Print("Binary PreOrder: ")
    PreOrder(binary, func(v int) { fmt.Printf("%d ", v) })
    fmt.Println() // 1 2 4 5 3

    // 测试多叉树
    nary := &NaryTree{
        val: 1,
        children: []*NaryTree{
            {val: 2, children: []*NaryTree{{val: 4}, {val: 5}}},
            {val: 3},
        },
    }

    fmt.Print("Nary PreOrder: ")
    PreOrder(nary, func(v int) { fmt.Printf("%d ", v) })
    fmt.Println() // 1 2 4 5 3
}
```

#### 场景3: 图算法抽象

```go
package main

import (
    "container/heap"
    "fmt"
)

// 图节点接口
type GraphNode[N GraphNode[N]] interface {
    comparable
    Neighbors() []N           // 相邻节点
    EdgeCost(to N) int        // 到邻居的边成本
    Heuristic(target N) int   // 启发式函数 (A*)
}

// Dijkstra 最短路径 - 通用实现
func Dijkstra[N GraphNode[N]](start, target N) ([]N, int) {
    dist := make(map[N]int)
    prev := make(map[N]N)
    visited := make(map[N]bool)

    dist[start] = 0

    for len(dist) > 0 {
        // 找到未访问的最小距离节点
        var u N
        minDist := int(^uint(0) >> 1) // MaxInt

        for node, d := range dist {
            if !visited[node] && d < minDist {
                minDist = d
                u = node
            }
        }

        if visited[u] {
            break
        }
        visited[u] = true

        if u == target {
            break
        }

        // 更新邻居距离
        for _, v := range u.Neighbors() {
            if visited[v] {
                continue
            }
            cost := dist[u] + u.EdgeCost(v)
            if d, ok := dist[v]; !ok || cost < d {
                dist[v] = cost
                prev[v] = u
            }
        }
    }

    // 重建路径
    if _, ok := dist[target]; !ok {
        return nil, -1
    }

    path := []N{target}
    for node := target; node != start; node = prev[node] {
        path = append([]N{prev[node]}, path...)
    }

    return path, dist[target]
}

// 简单图节点实现
type SimpleNode struct {
    name      string
    neighbors map[*SimpleNode]int
}

func (n *SimpleNode) Neighbors() []*SimpleNode {
    result := make([]*SimpleNode, 0, len(n.neighbors))
    for neighbor := range n.neighbors {
        result = append(result, neighbor)
    }
    return result
}

func (n *SimpleNode) EdgeCost(to *SimpleNode) int {
    return n.neighbors[to]
}

func (n *SimpleNode) Heuristic(target *SimpleNode) int {
    return 0 // Dijkstra 不需要启发式
}

func main() {
    // 构建图
    a := &SimpleNode{name: "A", neighbors: make(map[*SimpleNode]int)}
    b := &SimpleNode{name: "B", neighbors: make(map[*SimpleNode]int)}
    c := &SimpleNode{name: "C", neighbors: make(map[*SimpleNode]int)}
    d := &SimpleNode{name: "D", neighbors: make(map[*SimpleNode]int)}

    a.neighbors[b] = 1
    a.neighbors[c] = 4
    b.neighbors[c] = 2
    b.neighbors[d] = 5
    c.neighbors[d] = 1

    path, cost := Dijkstra(a, d)

    fmt.Printf("Path: ")
    for _, node := range path {
        fmt.Printf("%s ", node.name)
    }
    fmt.Printf("\nCost: %d\n", cost) // A B C D, Cost: 4
}
```

### 2.4 约束规则详解

```go
// 合法的自引用约束模式

// 模式1: 接口自引用
type Constraint1[T Constraint1[T]] interface {
    Method(T)
}

// 模式2: 结构体自引用
type Node[T Node[T]] struct {
    children []T
}

// 模式3: 多个类型参数
type Pair[A Pair[A, B], B any] interface {
    First() A
    Second() B
}

// 非法模式 (会导致无限展开)

// 模式1: 直接类型递归
type Bad1[T Bad1[T]] T  // 编译错误: 无限类型展开

// 模式2: 间接无限递归
type Bad2[T Bad3[T]] interface{}
type Bad3[T Bad2[T]] interface{}  // 编译错误
```

### 2.5 类型推断改进

```go
type Container[C Container[C]] interface {
    Add(item any)
    Get() any
}

// Go 1.26 改进了递归约束的类型推断
func Process[C Container[C]](c C) {
    // 编译器能更好地推断 C 的具体类型
}

// 使用示例
type MyContainer struct {
    items []any
}

func (m *MyContainer) Add(item any) {
    m.items = append(m.items, item)
}

func (m *MyContainer) Get() any {
    if len(m.items) == 0 {
        return nil
    }
    return m.items[0]
}

func main() {
    c := &MyContainer{}
    Process(c)  // 类型推断: C = *MyContainer
}
```

---

## 3. 形式化语义分析

### 3.1 new() 的形式化语义

```text
new 表达式的操作语义:
────────────────────────────────────────

[new-expr]
──────────
E ⊢ e : T    T 是值类型
─────────────────────────────
E ⊢ new(e) : *T

内存效果:
─────────
1. 分配 sizeof(T) 字节的内存
2. 将 e 的值复制到该内存
3. 返回指向该内存的指针

等价变换:
─────────
new(T(v)) ⟺ &T(v) ⟺ tmp := T(v); &tmp
```

### 3.2 递归约束的形式化语义

```text
递归类型约束的判定规则:
────────────────────────────────────────

给定约束: type C[T C[T]] interface { ... }

类型满足性判定 (Type Satisfaction):
────────────────────────────────────────
T satisfies C[T] ⟺
    1. T 实现了 C[T] 的所有方法
    2. 对于 C[T] 中的自引用方法 M(T):
       - T.M 的参数类型必须与 T 匹配

终止性保证 (Termination):
────────────────────────────────────────
递归约束必须满足:
1. 递归必须通过接口类型 (不是具体类型)
2. 不能出现无限制的展开: T → C[T] → C[C[T]] → ...
3. 实现类型必须是具体可实例化的
```

---

## 4. 最佳实践

### 4.1 new() 使用指南

| 场景 | 推荐方式 | 理由 |
|------|---------|------|
| 简单取地址 | `&T{...}` | 最清晰、最常用 |
| 可选字段 | `new(func())` | 简洁内联 |
| 复杂计算 | 先计算再取地址 | 可读性更好 |
| 性能敏感 | `&T{...}` 或 `new(T{...})` | 两者等价 |

### 4.2 递归约束使用指南

```go
// ✅ 使用递归约束的场景

// 1. 需要自引用方法的数据结构
type TreeNode[T TreeNode[T]] interface {
    Children() []T
}

// 2. 可比较的自引用类型
type Ordered[T Ordered[T]] interface {
    comparable
    Less(T) bool
}

// 3. 抽象算法需要统一接口
type GraphNode[N GraphNode[N]] interface {
    Neighbors() []N
}

// ❌ 避免使用递归约束的场景

// 1. 简单的泛型函数 - 不需要自引用
func Max[T comparable](a, b T) T  // 使用 comparable 即可

// 2. 容器类型 - 通常不需要
type Stack[T any] struct { ... }  // T any 足够

// 3. 类型转换函数
func Convert[T, U any](v T) U     // 无自引用需要
```

### 4.3 迁移建议

```go
// 旧代码: 使用 &T{...}
cfg := &Config{
    Host: "localhost",
    Port: 8080,
}

// 新代码: 可以保持不变，或在需要时使用 new()
// 两种写法都是合法且推荐的

// 推荐的迁移策略:
// 1. 保持现有代码不变
// 2. 新代码根据场景选择
// 3. 处理可选字段时优先考虑 new(expr)
```

---

## 参考

- [Go 1.26 Release Notes - Language Changes](https://go.dev/doc/go1.26#language)
- [Go Proposal: Recursive Generic Types](<https://github.com/golang/go/issues/>...
- [System F with Recursive Types](https://en.wikipedia.org/wiki/System_F)
