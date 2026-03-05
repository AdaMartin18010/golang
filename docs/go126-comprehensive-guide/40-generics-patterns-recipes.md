# Go泛型模式与实战食谱

> 从基础用法到高级模式的完整泛型指南

---

## 一、泛型设计决策

### 1.1 什么时候使用泛型

```text
适合使用泛型的场景：
────────────────────────────────────────

1. 数据结构：

// 通用的Stack实现
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

价值：
- 一次实现，支持任意类型
- 类型安全，无需类型断言
- 编译时类型检查

2. 通用算法：

// 通用的Map函数
func Map[T, R any](input []T, fn func(T) R) []R {
    result := make([]R, len(input))
    for i, v := range input {
        result[i] = fn(v)
    }
    return result
}

// 使用
ints := []int{1, 2, 3}
strings := Map(ints, func(i int) string {
    return strconv.Itoa(i)
})

3. 类型约束接口：

type Number interface {
    ~int | ~int64 | ~float64
}

func Sum[T Number](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}

不适合使用泛型的场景：
────────────────────────────────────────

1. 简单的类型转换：

// 过度设计
func ToString[T fmt.Stringer](v T) string {
    return v.String()
}

// 简单直接更好
func ToString(v fmt.Stringer) string {
    return v.String()
}

2. 只有interface就能解决：

// 不需要泛型
type Processor[T any] interface {
    Process(T) T
}

// interface就够了
type Processor interface {
    Process(interface{}) interface{}
}

3. 过早抽象：

// 先写具体实现
func ProcessInts([]int) { ... }

// 确实需要多种类型时再泛化
func Process[T any]([]T) { ... }
```

### 1.2 泛型的性能考量

```text
编译期实例化：
────────────────────────────────────────

Go泛型采用编译期实例化：
- 每个类型组合生成独立代码
- 无运行时类型擦除开销
- 与手写特定类型代码性能相同

对比Java泛型：

Java（类型擦除）：
List<Integer> 和 List<String>
→ 运行时都是 List<Object>
→ 装箱/拆箱开销

Go（编译期实例化）：
Stack[int] 和 Stack[string]
→ 两个独立的类型
→ 无额外开销

二进制大小影响：
────────────────────────────────────────

每个实例化增加代码：
- Stack[int] 生成一份代码
- Stack[string] 生成另一份代码
- N 个类型 → N 份代码

实践建议：
- 避免过多类型参数组合
- 常用类型可以单独优化
- 权衡通用性和二进制大小

性能测试：
────────────────────────────────────────

基准测试泛型 vs 非泛型：

func BenchmarkGeneric(b *testing.B) {
    s := NewStack[int]()
    for i := 0; i < b.N; i++ {
        s.Push(i)
        s.Pop()
    }
}

func BenchmarkNonGeneric(b *testing.B) {
    s := NewIntStack()
    for i := 0; i < b.N; i++ {
        s.Push(i)
        s.Pop()
    }
}

// 两者性能相同
```

---

## 二、类型约束设计

### 2.1 约束接口模式

```text
基础约束：
────────────────────────────────────────

any 约束（无限制）：
func Print[T any](v T) {
    fmt.Println(v)
}

comparable 约束：
func Contains[T comparable](slice []T, item T) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

方法约束：
type Stringer interface {
    String() string
}

func Join[T Stringer](items []T, sep string) string {
    // 使用 items[i].String()
}

复杂约束设计：
────────────────────────────────────────

数值类型约束：
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64 |
    ~complex64 | ~complex128
}

type Signed interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Float interface {
    ~float32 | ~float64
}

带方法的联合类型：
────────────────────────────────────────

type Ordered interface {
    constraints.Ordered  // int, uint, float, string
}

// 自定义Ordered类型
type MyInt int

func Max[T Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// 使用
m := Max(MyInt(10), MyInt(20))  // 正常工作

约束组合：
────────────────────────────────────────

type Addable interface {
    constraints.Integer | constraints.Float | ~string
}

func Add[T Addable](a, b T) T {
    return a + b
}

type Container[T any] interface {
    Len() int
    Get(int) T
    Set(int, T)
}
```

### 2.2 Go 1.26 递归类型约束

```text
递归约束的定义：
────────────────────────────────────────

Go 1.26 允许类型参数引用自身：

type Ordered[T interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64 |
    ~string |
    Ordered[T]  // 递归引用！
}] interface {
    Less(T) bool
}

应用场景：
────────────────────────────────────────

1. 树结构：

type Node[T interface {
    *TreeNode[T]  // 递归
}] struct {
    Value int
    Left  T
    Right T
}

type TreeNode[T any] struct {
    Value T
    Left  *TreeNode[T]
    Right *TreeNode[T]
}

// 通用树遍历
func InOrder[T interface{ *TreeNode[T] }](root T, visit func(T)) {
    if root == nil {
        return
    }
    InOrder(root.Left, visit)
    visit(root)
    InOrder(root.Right, visit)
}

2. 链表：

type ListNode[T any] struct {
    Value T
    Next  *ListNode[T]
}

// 递归处理链表
type List[T interface{ *ListNode[T] }] struct {
    Head T
}
```

---

## 三、泛型实战模式

### 3.1 数据结构实现

```text
通用缓存：
────────────────────────────────────────

type Cache[K comparable, V any] struct {
    data map[K]V
    mu   sync.RWMutex
    ttl  time.Duration
}

func NewCache[K comparable, V any](ttl time.Duration) *Cache[K, V] {
    return &Cache[K, V]{
        data: make(map[K]V),
        ttl:  ttl,
    }
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.data[key]
    return v, ok
}

func (c *Cache[K, V]) Set(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}

// 使用
cache := NewCache[string, User](time.Hour)
cache.Set("user:123", user)
user, ok := cache.Get("user:123")

通用结果类型：
────────────────────────────────────────

type Result[T any] struct {
    Value T
    Error error
}

func (r Result[T]) Ok() bool {
    return r.Error == nil
}

func (r Result[T]) Unwrap() (T, error) {
    return r.Value, r.Error
}

// 使用
func FetchData() Result[[]byte] {
    data, err := http.Get(url)
    if err != nil {
        return Result[[]byte]{Error: err}
    }
    return Result[[]byte]{Value: data}
}

func main() {
    result := FetchData()
    if !result.Ok() {
        log.Fatal(result.Error)
    }
    data, _ := result.Unwrap()
    // 使用 data
}

函数式编程工具：
────────────────────────────────────────

// Filter
func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce
func Reduce[T, R any](slice []T, initial R, fn func(R, T) R) R {
    result := initial
    for _, v := range slice {
        result = fn(result, v)
    }
    return result
}

// Find
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
    var zero T
    for _, v := range slice {
        if predicate(v) {
            return v, true
        }
    }
    return zero, false
}

// 使用示例
numbers := []int{1, 2, 3, 4, 5}

// 过滤偶数
evens := Filter(numbers, func(n int) bool {
    return n%2 == 0
})

// 求和
sum := Reduce(numbers, 0, func(acc, n int) int {
    return acc + n
})

// 查找第一个大于3的数
first, found := Find(numbers, func(n int) bool {
    return n > 3
})
```

### 3.2 设计模式实现

```text
通用工厂模式：
────────────────────────────────────────

type Creator[T any] interface {
    Create() T
}

type Factory[T any, C Creator[T]] struct {
    creators map[string]C
}

func (f *Factory[T, C]) Register(name string, creator C) {
    f.creators[name] = creator
}

func (f *Factory[T, C]) Create(name string) (T, error) {
    var zero T
    creator, ok := f.creators[name]
    if !ok {
        return zero, fmt.Errorf("unknown type: %s", name)
    }
    return creator.Create(), nil
}

通用构建器：
────────────────────────────────────────

type Builder[T any] struct {
    target T
    steps  []func(T) T
}

func NewBuilder[T any](initial T) *Builder[T] {
    return &Builder[T]{target: initial}
}

func (b *Builder[T]) Step(fn func(T) T) *Builder[T] {
    b.steps = append(b.steps, fn)
    return b
}

func (b *Builder[T]) Build() T {
    result := b.target
    for _, step := range b.steps {
        result = step(result)
    }
    return result
}

// 使用
type Config struct {
    Host string
    Port int
    SSL  bool
}

config := NewBuilder(Config{}).
    Step(func(c Config) Config {
        c.Host = "localhost"
        return c
    }).
    Step(func(c Config) Config {
        c.Port = 8080
        return c
    }).
    Step(func(c Config) Config {
        c.SSL = true
        return c
    }).
    Build()
```

---

*本章提供了Go泛型的完整实战指南，从基础到高级模式。*
