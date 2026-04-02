# 结构体嵌入 (Struct Embedding)

> **分类**: 语言设计

---

## 基本嵌入

```go
type Reader struct{}
func (r Reader) Read() {}

type Writer struct{}
func (w Writer) Write() {}

// 嵌入
type ReadWriter struct {
    Reader
    Writer
}

// 自动拥有 Read() 和 Write() 方法
var rw ReadWriter
rw.Read()
rw.Write()
```

---

## 嵌入 vs 组合

```go
// 嵌入 - 方法提升到外层
type Engine struct{}
func (e Engine) Start() {}

type Car struct {
    Engine  // 嵌入
}

car := Car{}
car.Start()  // 直接调用

// 组合 - 需要间接访问
type Car2 struct {
    engine Engine  // 组合
}

car2 := Car2{}
car2.engine.Start()  // 通过字段访问
```

---

## 嵌入指针

```go
type Widget struct {
    X, Y int
}

// 嵌入指针
type Label struct {
    *Widget
    Text string
}

label := Label{
    Widget: &Widget{X: 10, Y: 20},
    Text:   "Hello",
}
```

---

## 方法覆盖

```go
type Base struct{}
func (b Base) Method() string { return "base" }

type Derived struct {
    Base
}

// 覆盖方法
func (d Derived) Method() string { return "derived" }

// 访问被覆盖的方法
type Derived2 struct {
    Base
}

func (d Derived2) Method() string {
    return d.Base.Method() + " extended"
}
```

---

## 实际应用

### io.ReadWriter

```go
type ReadWriter struct {
    *PipeReader
    *PipeWriter
}
```

### HTTP Handler

```go
type MyHandler struct {
    http.ServeMux  // 嵌入路由
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 预处理
    log.Println("Request:", r.URL)

    // 调用嵌入的 ServeMux
    h.ServeMux.ServeHTTP(w, r)
}
```
