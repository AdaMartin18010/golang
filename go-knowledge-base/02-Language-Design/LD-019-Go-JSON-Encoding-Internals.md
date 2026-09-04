# LD-019: Go JSON 编码内部原理 (Go JSON Encoding Internals)

> **维度**: Language Design
> **级别**: S (22 KB)
> **标签**: #json #encoding #reflection #performance #codegen #serialization
> **权威来源**:
>
> - [encoding/json Package](https://github.com/golang/go/tree/master/src/encoding/json) - Go Authors
> - [JSON and Go](https://go.dev/blog/json) - Go Authors
> - [High Performance JSON](https://github.com/json-iterator/go-benchmark) - JSON Benchmarks

> **Go 版本**: 1.27+
---

## 1. JSON 包架构

### 1.1 核心组件

```
┌─────────────────────────────────────────────────────────────┐
│                   encoding/json                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │  Marshal    │───►│  encodeState│───►│  encode     │     │
│  │             │    │  (buffer)   │    │  (types)    │     │
│  └─────────────┘    └─────────────┘    └─────────────┘     │
│                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │  Unmarshal  │───►│  Decoder    │───►│  decode     │     │
│  │             │    │  (scanner)  │    │  (types)    │     │
│  └─────────────┘    └─────────────┘    └─────────────┘     │
│                                                              │
│  ┌─────────────┐    ┌─────────────┐                        │
│  │  Scanner    │    │  reflect    │                        │
│  │  (lexer)    │    │  (types)    │                        │
│  └─────────────┘    └─────────────┘                        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 关键数据结构

```go
// src/encoding/json/encode.go

// encodeState 编码状态
type encodeState struct {
    bytes.Buffer           // 输出缓冲
    scratch      [64]byte // 临时缓冲区

    // 递归深度限制
    ptrLevel     uint
    ptrSeen      map[any]struct{} // 循环引用检测
}

// encOpts 编码选项
type encOpts struct {
    quoted bool // 字符串引号
    escape bool // HTML 转义
}

// 类型编码器缓存
type encoderFunc func(e *encodeState, v reflect.Value, opts encOpts)

var encoderCache sync.Map // map[reflect.Type]encoderFunc
```

---

## 2. Marshal 实现原理

### 2.1 编码流程

```go
func Marshal(v any) ([]byte, error) {
    e := newEncodeState()
    defer encodeStatePool.Put(e)

    err := e.marshal(v, encOpts{escape: true})
    if err != nil {
        return nil, err
    }

    buf := append([]byte(nil), e.Bytes()...)
    return buf, nil
}

func (e *encodeState) marshal(v any, opts encOpts) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
        }
    }()

    e.reflectValue(reflect.ValueOf(v), opts)
    return nil
}

func (e *encodeState) reflectValue(v reflect.Value, opts encOpts) {
    valueEncoder(v)(e, v, opts)
}
```

### 2.2 类型编码器选择

```go
func valueEncoder(v reflect.Value) encoderFunc {
    if !v.IsValid() {
        return invalidValueEncoder
    }
    return typeEncoder(v.Type())
}

func typeEncoder(t reflect.Type) encoderFunc {
    // 检查缓存
    if fi, ok := encoderCache.Load(t); ok {
        return fi.(encoderFunc)
    }

    // 加锁创建
    encoderCacheMu.Lock()
    defer encoderCacheMu.Unlock()

    if fi, ok := encoderCache.Load(t); ok {
        return fi.(encoderFunc)
    }

    // 创建编码器
    var f encoderFunc
    wg := &sync.WaitGroup{}
    wg.Add(1)
    encoderCache.Store(t, func(e *encodeState, v reflect.Value, opts encOpts) {
        wg.Wait()
        f(e, v, opts)
    })

    f = newTypeEncoder(t, true)
    wg.Done()
    encoderCache.Store(t, f)
    return f
}
```

### 2.3 具体类型编码

```go
// 字符串编码
func stringEncoder(e *encodeState, v reflect.Value, opts encOpts) {
    if v.Type() == numberType {
        // json.Number 特殊处理
        numStr := v.String()
        if numStr == "" {
            numStr = "0"
        }
        e.WriteString(numStr)
        return
    }

    s := v.String()
    if opts.quoted {
        e.WriteByte('"')
    }
    e.WriteString(strconv.Quote(s))
    if opts.quoted {
        e.WriteByte('"')
    }
}

// 结构体编码
func newStructEncoder(t reflect.Type) encoderFunc {
    // 收集字段信息
    fields := typeFields(t) // 解析 json tag
    se := &structEncoder{
        fields:    fields,
        fieldEncs: make([]encoderFunc, len(fields)),
    }

    for i, f := range fields {
        se.fieldEncs[i] = typeEncoder(f.typ)
    }

    return se.encode
}

func (se *structEncoder) encode(e *encodeState, v reflect.Value, opts encOpts) {
    e.WriteByte('{')
    first := true

    for i, f := range se.fields {
        fv := fieldByIndex(v, f.index)
        if !fv.IsValid() || f.omitEmpty && isEmptyValue(fv) {
            continue
        }

        if first {
            first = false
        } else {
            e.WriteByte(',')
        }

        // 字段名
        e.WriteString(f.nameQuoted)
        e.WriteByte(':')

        // 字段值
        opts.quoted = f.quoted
        se.fieldEncs[i](e, fv, opts)
    }

    e.WriteByte('}')
}
```

### 2.4 字段标签解析

```go
// 解析 json tag
type field struct {
    name      string      // JSON 字段名
    nameQuoted string     // 带引号的字段名
    tag       bool        // 是否有 tag
    index     []int       // 字段索引
    typ       reflect.Type
    omitEmpty bool
    quoted    bool
}

func typeFields(t reflect.Type) []field {
    // 使用缓存
    if f, ok := fieldCache.Load(t); ok {
        return f.([]field)
    }

    var fields []field

    for i := 0; i < t.NumField(); i++ {
        sf := t.Field(i)

        // 跳过未导出字段
        if !sf.IsExported() {
            continue
        }

        tag := sf.Tag.Get("json")
        if tag == "-" {
            continue
        }

        name, opts := parseTag(tag)
        if name == "" {
            name = sf.Name
        }

        field := field{
            name:       name,
            nameQuoted: strconv.Quote(name),
            tag:        tag != "",
            index:      []int{i},
            typ:        sf.Type,
            omitEmpty:  opts.Contains("omitempty"),
            quoted:     opts.Contains("string"),
        }

        fields = append(fields, field)
    }

    fieldCache.Store(t, fields)
    return fields
}
```

---

## 3. Unmarshal 实现原理

### 3.1 解码流程

```go
func Unmarshal(data []byte, v any) error {
    var d decodeState
    d.init(data)
    return d.unmarshal(v)
}

func (d *decodeState) unmarshal(v any) error {
    rv := reflect.ValueOf(v)
    if rv.Kind() != reflect.Pointer || rv.IsNil() {
        return &InvalidUnmarshalError{reflect.TypeOf(v)}
    }

    // 扫描 JSON
    d.scanWhile(scanSkipSpace)

    // 解码到目标类型
    d.value(rv)

    return d.savedError
}
```

### 3.2 Scanner 实现

```go
// src/encoding/json/scanner.go

// Scanner 状态机
type scanner struct {
    step       func(*scanner, byte) int
    endTop     bool     // 顶层结束
    parseState []int    // 解析栈
    err        error
}

// 扫描状态
const (
    scanContinue     = iota // 继续
    scanBeginLiteral        // 开始字面量
    scanBeginObject         // 开始对象
    scanObjectKey           // 对象键
    scanObjectColon         // 对象冒号
    scanObjectValue         // 对象值
    scanBeginArray          // 开始数组
    scanArrayValue          // 数组值
    scanArrayComma          // 数组逗号
    scanEndObject           // 结束对象
    scanEndArray            // 结束数组
    scanError               // 错误
)

func (s *scanner) scan(next byte) int {
    return s.step(s, next)
}

// 状态转换
func stateBegin(s *scanner, c byte) int {
    switch c {
    case '{':
        s.step = stateBeginObject
        return scanBeginObject
    case '[':
        s.step = stateBeginArray
        return scanBeginArray
    case '"':
        s.step = stateInString
        return scanBeginLiteral
    case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
        s.step = stateBeginNumber
        return scanBeginLiteral
    case 't', 'f', 'n': // true, false, null
        s.step = stateBeginLiteral
        return scanBeginLiteral
    }
    return scanError
}
```

### 3.3 解码类型匹配

```go
func (d *decodeState) value(v reflect.Value) {
    // 检查 UnmarshalJSON 接口
    if u, ok := v.Interface().(Unmarshaler); ok {
        d.literalStore(u)
        return
    }

    // 根据 JSON 类型解码
    switch d.opcode {
    case scanBeginArray:
        d.array(v)
    case scanBeginObject:
        d.object(v)
    case scanBeginLiteral:
        d.literal(v)
    }
}

// 对象解码
func (d *decodeState) object(v reflect.Value) {
    // 创建字段映射
    fields := cachedTypeFields(v.Type())

    // 读取 {
    d.scanWhile(scanSkipSpace)

    // 空对象
    if d.opcode == scanEndObject {
        return
    }

    for {
        // 读取键
        start := d.readIndex()
        d.scanWhile(scanContinue)
        item := d.data[start:d.readIndex()]
        key := string(item)

        // 查找字段
        var subv reflect.Value
        if f, ok := fields[key]; ok {
            subv = v.FieldByIndex(f.index)
        }

        // 读取值
        d.scanWhile(scanSkipSpace)
        d.value(subv)

        // 检查 }
        d.scanWhile(scanSkipSpace)
        if d.opcode == scanEndObject {
            break
        }

        // 期望逗号
        if d.opcode != scanObjectComma {
            d.error("expected comma after object element")
        }
    }
}
```

---

## 4. 性能优化

### 4.1 编码优化策略

```go
// 1. 使用 sync.Pool 复用 encoder
var encodeStatePool = sync.Pool{
    New: func() interface{} {
        return &encodeState{}
    },
}

func newEncodeState() *encodeState {
    e := encodeStatePool.Get().(*encodeState)
    e.Reset()
    return e
}

// 2. 避免反射 - 代码生成
go:generate go run gen.go

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// 生成代码 (gen.go 输出)
func (u User) MarshalJSON() ([]byte, error) {
    var buf [64]byte
    b := buf[:0]
    b = append(b, '{')
    b = append(b, `"id":`...)
    b = strconv.AppendInt(b, int64(u.ID), 10)
    b = append(b, ',')
    b = append(b, `"name":`...)
    b = append(b, strconv.Quote(u.Name)...)
    b = append(b, '}')
    return b, nil
}

// 3. 使用 json.RawMessage 延迟解析
type Event struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func processEvent(data []byte) error {
    var event Event
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    // 根据类型解析 Data
    switch event.Type {
    case "user":
        var user User
        return json.Unmarshal(event.Data, &user)
    case "order":
        var order Order
        return json.Unmarshal(event.Data, &order)
    }
    return nil
}
```

### 4.2 解码优化策略

```go
// 1. 使用 Decoder 流式处理
func processLargeJSON(r io.Reader) error {
    dec := json.NewDecoder(r)

    // 读取 [
    _, err := dec.Token()
    if err != nil {
        return err
    }

    // 逐个解码
    for dec.More() {
        var item Item
        if err := dec.Decode(&item); err != nil {
            return err
        }
        // 处理 item...
    }

    // 读取 ]
    _, err = dec.Token()
    return err
}

// 2. 使用 Number 避免精度丢失
func decodeNumber(data []byte) error {
    var result map[string]json.Number
    if err := json.Unmarshal(data, &result); err != nil {
        return err
    }

    n := result["big_number"]
    // 保留精度
    i, err := n.Int64()
    f, err := n.Float64()
    s := n.String()
}

// 3. 预分配切片容量
type Response struct {
    Items []Item `json:"items"`
}

func decodeWithCapacity(data []byte) (*Response, error) {
    // 使用自定义类型控制解码
    var raw struct {
        Items json.RawMessage `json:"items"`
    }
    if err := json.Unmarshal(data, &raw); err != nil {
        return nil, err
    }

    // 先获取数组长度
    var items []json.RawMessage
    if err := json.Unmarshal(raw.Items, &items); err != nil {
        return nil, err
    }

    // 预分配
    result := &Response{
        Items: make([]Item, 0, len(items)),
    }

    for _, item := range items {
        var i Item
        if err := json.Unmarshal(item, &i); err != nil {
            return nil, err
        }
        result.Items = append(result.Items, i)
    }

    return result, nil
}
```

### 4.3 基准测试

```go
func BenchmarkMarshal(b *testing.B) {
    type User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    user := User{ID: 1, Name: "John", Age: 30}

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _, err := json.Marshal(user)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkUnmarshal(b *testing.B) {
    data := []byte(`{"id":1,"name":"John","age":30}`)

    type User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
        Age  int    `json:"age"`
    }

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        var user User
        if err := json.Unmarshal(data, &user); err != nil {
            b.Fatal(err)
        }
    }
}

// 典型结果 (Go 1.21)
// BenchmarkMarshal-8      5000000    285 ns/op    128 B/op    2 allocs/op
// BenchmarkUnmarshal-8    3000000    425 ns/op    192 B/op    4 allocs/op
```

---

## 5. 并发安全分析

### 5.1 线程安全保证

```go
// 类型编码器缓存是并发安全的
var encoderCache sync.Map

// Encoder/Decoder 不是并发安全的
type Encoder struct {
    w       io.Writer
    encodeState
}

type Decoder struct {
    r       io.Reader
    scanner
}

// 正确用法
func threadSafeEncoding() {
    // 每个 goroutine 独立的 Encoder
    var buf bytes.Buffer
    enc := json.NewEncoder(&buf)
    enc.Encode(data)

    // 或者使用 Marshal（并发安全）
    data, _ := json.Marshal(obj)
}
```

### 5.2 并发编码模式

```go
// 并行编码大量对象
func parallelEncode(items []Item) [][]byte {
    results := make([][]byte, len(items))

    var wg sync.WaitGroup
    sem := make(chan struct{}, runtime.GOMAXPROCS(0))

    for i, item := range items {
        wg.Add(1)
        sem <- struct{}{}

        go func(idx int, it Item) {
            defer wg.Done()
            defer func() { <-sem }()

            data, err := json.Marshal(it)
            if err != nil {
                return
            }
            results[idx] = data
        }(i, item)
    }

    wg.Wait()
    return results
}
```

---

## 6. 视觉表征

### 6.1 Marshal 流程图

```
Input: Go Value
      │
      ▼
┌─────────────┐
│  Marshal()  │
└──────┬──────┘
       │
       ▼
┌─────────────┐     缓存命中    ┌─────────────┐
│  typeEncoder│────────────────►│  返回缓存    │
│  (查缓存)   │                │  encoder    │
└──────┬──────┘                └─────────────┘
       │ 缓存未命中
       ▼
┌─────────────┐
│ newTypeEncoder│
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 类型分类    │───► 基本类型 (int, string, bool)
│            │───► 复合类型 (struct, slice, map)
│            │───► 接口类型 (json.Marshaler)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  递归编码   │───► 写入 encodeState.Buffer
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 返回 []byte │
└─────────────┘
```

### 6.2 Unmarshal 状态机

```
    ┌─────────┐
    │  Start  │
    └────┬────┘
         │
    ┌────┴────┬────────┬────────┐
    │         │        │        │
    ▼         ▼        ▼        ▼
┌──────┐  ┌──────┐ ┌──────┐ ┌──────┐
│  {   │  │  [   │ │  "   │ │number│
└──┬───┘  └──┬───┘ └──┬───┘ └──┬───┘
   │         │        │        │
   ▼         ▼        ▼        ▼
┌──────┐  ┌──────┐ ┌──────┐ ┌──────┐
│Object│  │Array │ │String│ │Literal
│State │  │State │ │State │ │State │
└──────┘  └──────┘ └──────┘ └──────┘
```

### 6.3 性能优化决策树

```
JSON 性能问题?
│
├── 编码慢?
│   ├── 使用代码生成替代反射
│   ├── 实现 json.Marshaler 接口
│   ├── 使用 json.RawMessage 延迟编码
│   └── 复用 encodeState (sync.Pool)
│
├── 解码慢?
│   ├── 使用 Decoder 流式处理
│   ├── 实现 json.Unmarshaler 接口
│   ├── 预分配目标切片容量
│   └── 使用 json.Number 避免二次解析
│
└── 内存占用高?
    ├── 减少大对象嵌套
    ├── 使用流式 API (Encoder/Decoder)
    └── 复用缓冲区
```

---

## 7. 完整代码示例

### 7.1 高性能 JSON 处理

```go
package main

import (
    "bytes"
    "encoding/json"
    "strconv"
    "sync"
)

// 自定义 Marshal 减少分配
type User struct {
    ID   int
    Name string
    Age  int
}

func (u User) MarshalJSON() ([]byte, error) {
    // 预计算容量避免扩容
    var buf [128]byte
    b := buf[:0]

    b = append(b, '{')

    // ID
    b = append(b, `"id":`...)
    b = strconv.AppendInt(b, int64(u.ID), 10)
    b = append(b, ',')

    // Name
    b = append(b, `"name":`...)
    b = strconv.AppendQuote(b, u.Name)
    b = append(b, ',')

    // Age
    b = append(b, `"age":`...)
    b = strconv.AppendInt(b, int64(u.Age), 10)

    b = append(b, '}')

    result := make([]byte, len(b))
    copy(result, b)
    return result, nil
}

// 缓冲池
type bufferPool struct {
    pool sync.Pool
}

func newBufferPool() *bufferPool {
    return &bufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                return new(bytes.Buffer)
            },
        },
    }
}

func (p *bufferPool) Get() *bytes.Buffer {
    return p.pool.Get().(*bytes.Buffer)
}

func (p *bufferPool) Put(b *bytes.Buffer) {
    b.Reset()
    p.pool.Put(b)
}

// 批量编码器
type BatchEncoder struct {
    pool *bufferPool
}

func NewBatchEncoder() *BatchEncoder {
    return &BatchEncoder{
        pool: newBufferPool(),
    }
}

func (e *BatchEncoder) EncodeUsers(users []User) []byte {
    buf := e.pool.Get()
    defer e.pool.Put(buf)

    buf.WriteByte('[')
    for i, u := range users {
        if i > 0 {
            buf.WriteByte(',')
        }
        data, _ := u.MarshalJSON()
        buf.Write(data)
    }
    buf.WriteByte(']')

    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result
}

// 流式解码器
func StreamDecode(r io.Reader, callback func(User) error) error {
    dec := json.NewDecoder(r)

    // 读取 [
    _, err := dec.Token()
    if err != nil {
        return err
    }

    for dec.More() {
        var user User
        if err := dec.Decode(&user); err != nil {
            return err
        }
        if err := callback(user); err != nil {
            return err
        }
    }

    // 读取 ]
    _, err = dec.Token()
    return err
}

func main() {
    users := []User{
        {ID: 1, Name: "Alice", Age: 30},
        {ID: 2, Name: "Bob", Age: 25},
    }

    encoder := NewBatchEncoder()
    data := encoder.EncodeUsers(users)
    println(string(data))
}
```

---

**质量评级**: S (17KB)
**完成日期**: 2026-04-02
