# Go常见陷阱与解决方案

> 收集Go开发中的常见陷阱、错误模式及最佳解决方案

---

## 一、并发陷阱

### 1.1 Goroutine泄漏

```text
泄漏场景:
────────────────────────────────────────

1. 无退出条件的goroutine
2. 阻塞在channel发送/接收
3. 无限循环

代码示例:
// 陷阱: goroutine泄漏
func leakyChannel() {
    ch := make(chan int)
    go func() {
        // 永远阻塞，如果无人接收
        ch <- 1
    }()
    // 函数返回，但goroutine仍在等待
}

// 解决方案: 使用buffered channel或select
type SafeChannel struct {
    ch      chan int
    done    chan struct{}
    timeout time.Duration
}

func (sc *SafeChannel) Send(v int) bool {
    select {
    case sc.ch <- v:
        return true
    case <-sc.done:
        return false
    case <-time.After(sc.timeout):
        return false
    }
}

// 解决方案: 使用context
type Worker struct {
    ctx    context.Context
    cancel context.CancelFunc
}

func NewWorker() *Worker {
    ctx, cancel := context.WithCancel(context.Background())
    return &Worker{ctx: ctx, cancel: cancel}
}

func (w *Worker) Start() {
    go func() {
        for {
            select {
            case <-w.ctx.Done():
                return
            default:
                // 工作
            }
        }
    }()
}

func (w *Worker) Stop() {
    w.cancel()
}
```

### 1.2 竞态条件

```text
竞态场景:
────────────────────────────────────────

1. 非原子操作
2. 非并发安全的map
3. 闭包变量捕获

代码示例:
// 陷阱: 竞态条件
func raceCondition() {
    var counter int
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++  // 竞态!
        }()
    }
    wg.Wait()
    fmt.Println(counter) // 不确定
}

// 解决方案1: Mutex
func fixedWithMutex() {
    var counter int
    var mu sync.Mutex
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }
    wg.Wait()
}

// 解决方案2: atomic
func fixedWithAtomic() {
    var counter int64
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1)
        }()
    }
    wg.Wait()
}

// 陷阱: 闭包变量竞态
func closureRace() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            fmt.Println(i)  // 竞态! 所有goroutine共享i
        }()
    }
    wg.Wait()
}

// 解决方案: 参数传递
func fixedClosure() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {  // 每个goroutine有自己的n
            defer wg.Done()
            fmt.Println(n)
        }(i)
    }
    wg.Wait()
}
```

### 1.3 死锁

```text
死锁场景:
────────────────────────────────────────

1. 循环等待
2. 锁顺序不一致
3. 无缓冲channel阻塞

代码示例:
// 陷阱: 锁顺序死锁
func deadlockExample() {
    var mu1, mu2 sync.Mutex

    go func() {
        mu1.Lock()
        time.Sleep(time.Millisecond)
        mu2.Lock()  // 等待mu2
        // ...
        mu2.Unlock()
        mu1.Unlock()
    }()

    mu2.Lock()
    time.Sleep(time.Millisecond)
    mu1.Lock()  // 等待mu1 - 死锁!
    // ...
    mu1.Unlock()
    mu2.Unlock()
}

// 解决方案: 统一锁顺序
type Resource struct {
    id   int
    mu   sync.Mutex
    data int
}

func acquireInOrder(r1, r2 *Resource) (func(), func()) {
    if r1.id < r2.id {
        r1.mu.Lock()
        r2.mu.Lock()
        return func() { r2.mu.Unlock() }, func() { r1.mu.Unlock() }
    }
    r2.mu.Lock()
    r1.mu.Lock()
    return func() { r1.mu.Unlock() }, func() { r2.mu.Unlock() }
}

// 陷阱: channel死锁
func channelDeadlock() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        ch1 <- 1
        <-ch2
    }()

    ch2 <- 2  // 等待ch2的接收者
    <-ch1     // 永远不会执行
}

// 解决方案: 使用select
type SafeChannelComm struct {
    ch1 chan int
    ch2 chan int
}

func (s *SafeChannelComm) Exchange(v1, v2 int) bool {
    select {
    case s.ch1 <- v1:
        select {
        case v := <-s.ch2:
            _ = v
            return true
        case <-time.After(time.Second):
            return false
        }
    case <-time.After(time.Second):
        return false
    }
}
```

---

## 二、类型系统陷阱

### 2.1 接口陷阱

```text
常见陷阱:
────────────────────────────────────────

1. nil接口不等于nil
2. 方法集不匹配
3. 类型断言panic

代码示例:
// 陷阱1: nil接口
func nilInterfaceTrap() {
    var p *int = nil
    var i interface{} = p

    fmt.Println(p == nil) // true
    fmt.Println(i == nil) // false!

    // 原因: i包含 (*int, nil)，类型信息非空

    // 正确检查方式
    if p, ok := i.(*int); ok && p == nil {
        fmt.Println("nil *int")
    }
}

// 陷阱2: 值接收者vs指针接收者
type MyInt int

func (m MyInt) ValueMethod() {}
func (m *MyInt) PointerMethod() {}

type Interface interface {
    ValueMethod()
    PointerMethod()
}

func receiverTrap() {
    var v MyInt
    // var i Interface = v  // 错误: v没有PointerMethod
    var i Interface = &v  // 正确: *MyInt有所有方法
    _ = i
}

// 陷阱3: 接口比较panic
func interfaceCompareTrap() {
    var a interface{} = []int{1, 2, 3}
    var b interface{} = []int{1, 2, 3}
    // fmt.Println(a == b)  // panic: 切片不可比较

    // 解决方案
    equal := reflect.DeepEqual(a, b)
    fmt.Println(equal)
}
```

### 2.2 泛型陷阱

```text
泛型陷阱:
────────────────────────────────────────

1. 类型约束不满足
2. 类型推断失败
3. ~ (近似约束) 误解

代码示例:
// 陷阱1: 约束不满足
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64 | ~string
}

func Min[T Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

type MyInt int

func constraintTrap() {
    // Min(MyInt(1), MyInt(2))  // 错误: MyInt不满足Ordered

    // 需要显式转换
    Min(int(MyInt(1)), int(MyInt(2)))
}

// 陷阱2: 类型推断失败
func inferenceTrap() {
    // result := Map([]int{1, 2}, func(n int) string { return fmt.Sprint(n) })
    // 需要显式指定类型或确保编译器能推断
}
```

---

## 三、内存陷阱

### 3.1 内存泄漏

```text
泄漏场景:
────────────────────────────────────────

1. 全局变量累积
2. 未关闭的资源
3. 切片引用

代码示例:
// 陷阱: 切片引用导致内存不释放
func sliceReferenceLeak() {
    big := make([]byte, 100<<20)  // 100MB
    small := big[:10]

    // big可以被GC，但底层数组不能
    // 因为small引用它

    // 解决方案: 复制
    smallCopy := make([]byte, 10)
    copy(smallCopy, big[:10])
    // 现在big的底层数组可以被GC
    _ = smallCopy
}

// 陷阱: 全局map累积
var globalCache = make(map[string][]byte)

func globalMapLeak() {
    // 无限制增长
    data := make([]byte, 1000)
    globalCache[time.Now().String()] = data
    // 从不清理...
}

// 解决方案: LRU缓存
type LRUCache struct {
    mu       sync.Mutex
    cache    map[string]*list.Element
    lru      *list.List
    maxSize  int
}

type entry struct {
    key   string
    value []byte
}

func (c *LRUCache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if elem, ok := c.cache[key]; ok {
        c.lru.MoveToFront(elem)
        return elem.Value.(*entry).value, true
    }
    return nil, false
}

func (c *LRUCache) Set(key string, value []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if elem, ok := c.cache[key]; ok {
        c.lru.MoveToFront(elem)
        elem.Value.(*entry).value = value
        return
    }

    if c.lru.Len() >= c.maxSize {
        // 淘汰最久未使用
        back := c.lru.Back()
        if back != nil {
            c.lru.Remove(back)
            delete(c.cache, back.Value.(*entry).key)
        }
    }

    elem := c.lru.PushFront(&entry{key, value})
    c.cache[key] = elem
}
```

### 3.2 逃逸分析陷阱

```text
逃逸陷阱:
────────────────────────────────────────

1. 意外的堆分配
2. 接口装箱
3. 闭包捕获

代码示例:
// 陷阱: 意外的逃逸
func escapeTrap() *int {
    x := 10
    return &x  // x逃逸到堆
}

// 陷阱: 接口装箱
func interfaceBoxing() {
    var x int = 10
    var i interface{} = x  // x被装箱到堆
    _ = i
}

// 陷阱: 闭包捕获
func closureCapture() func() int {
    x := 10
    return func() int {  // x逃逸到堆
        return x
    }
}

// 优化: 使用值返回
func noEscape() int {
    x := 10
    return x  // 栈分配
}
```

---

## 四、标准库陷阱

### 4.1 JSON处理陷阱

```text
JSON陷阱:
────────────────────────────────────────

1. 大小写敏感
2. Number处理
3. 空值处理

代码示例:
// 陷阱1: 大小写
func jsonCaseTrap() {
    type Person struct {
        Name string  // 导出字段
        age  int     // 未导出，不会被序列化
    }

    p := Person{Name: "John", age: 30}
    b, _ := json.Marshal(p)
    fmt.Println(string(b)) // {"Name":"John"}
}

// 解决方案
func jsonCaseFixed() {
    type Person struct {
        Name string `json:"name"`  // 指定JSON字段名
        Age  int    `json:"age"`
    }

    p := Person{Name: "John", Age: 30}
    b, _ := json.Marshal(p)
    fmt.Println(string(b)) // {"name":"John","age":30}
}

// 陷阱2: Number精度
func jsonNumberTrap() {
    data := `{"id": 9007199254740993}`  // 超过float53精度
    var result map[string]interface{}
    json.Unmarshal([]byte(data), &result)

    // id被解析为float64，精度丢失
    fmt.Println(result["id"]) // 9.007199254740992e+15
}

// 解决方案
func jsonNumberFixed() {
    data := `{"id": 9007199254740993}`
    var result map[string]json.Number
    json.Unmarshal([]byte(data), &result)

    // 使用json.Number保持精度
    id, _ := result["id"].Int64()
    fmt.Println(id) // 9007199254740993
}
```

### 4.2 Time处理陷阱

```text
Time陷阱:
────────────────────────────────────────

1. 时区问题
2. 单调时钟vs墙钟
3. JSON序列化

代码示例:
// 陷阱: 时区问题
func timezoneTrap() {
    t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    fmt.Println(t.Local())  // 转换到本地时区，可能不是同一天
}

// 解决方案: 统一使用UTC
func timezoneFixed() {
    t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    fmt.Println(t.UTC())
}

// 陷阱: time.After内存泄漏
func afterLeak() {
    for {
        select {
        case <-time.After(5 * time.Minute):  // 每次创建新timer
            // 处理
        }
    }
}

// 解决方案: 复用timer
func afterFixed() {
    timer := time.NewTimer(5 * time.Minute)
    defer timer.Stop()

    for {
        timer.Reset(5 * time.Minute)
        select {
        case <-timer.C:
            // 处理
        }
    }
}
```

---

*本章收集了Go开发中的常见陷阱和解决方案，帮助开发者避免常见错误。*
