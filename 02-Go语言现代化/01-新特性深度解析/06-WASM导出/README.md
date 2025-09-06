# Go WebAssembly导出功能深度解析

<!-- TOC START -->
- [Go WebAssembly导出功能深度解析](#go-webassembly导出功能深度解析)
  - [1.1 概述](#11-概述)
  - [1.2 WASM导出基础](#12-wasm导出基础)
  - [1.3 函数导出](#13-函数导出)
  - [1.4 内存管理](#14-内存管理)
  - [1.5 类型转换](#15-类型转换)
  - [1.6 高级特性](#16-高级特性)
  - [1.7 性能优化](#17-性能优化)
  - [1.8 实际应用](#18-实际应用)
<!-- TOC END -->

## 1.1 概述

Go 1.21+ 引入了强大的WebAssembly导出功能，允许Go程序编译为WASM模块并在浏览器或其他WASM运行时中执行。这个功能为Go语言开辟了新的应用场景，特别是在前端开发、边缘计算和跨平台部署方面。

## 1.2 WASM导出基础

### 1.2.1 基本编译配置

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "fmt"
)

func main() {
    // 注册全局函数
    js.Global().Set("goAdd", js.FuncOf(add))
    js.Global().Set("goMultiply", js.FuncOf(multiply))
    
    // 保持程序运行
    select {}
}

// add 加法函数
func add(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Expected 2 arguments")
    }
    
    a := args[0].Float()
    b := args[1].Float()
    
    return js.ValueOf(a + b)
}

// multiply 乘法函数
func multiply(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Expected 2 arguments")
    }
    
    a := args[0].Float()
    b := args[1].Float()
    
    return js.ValueOf(a * b)
}
```

### 1.2.2 编译命令

```bash
# 编译为WASM
GOOS=js GOARCH=wasm go build -o main.wasm main.go

# 使用TinyGo编译（更小的文件大小）
tinygo build -target wasm -o main.wasm main.go
```

## 1.3 函数导出

### 1.3.1 基础函数导出

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "encoding/json"
)

// 导出数学运算函数
func init() {
    js.Global().Set("math", js.ValueOf(map[string]interface{}{
        "add":      js.FuncOf(add),
        "subtract": js.FuncOf(subtract),
        "multiply": js.FuncOf(multiply),
        "divide":   js.FuncOf(divide),
        "power":    js.FuncOf(power),
    }))
}

func add(this js.Value, args []js.Value) interface{} {
    if len(args) < 2 {
        return js.ValueOf("Error: At least 2 arguments required")
    }
    
    result := 0.0
    for _, arg := range args {
        result += arg.Float()
    }
    
    return js.ValueOf(result)
}

func subtract(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Exactly 2 arguments required")
    }
    
    return js.ValueOf(args[0].Float() - args[1].Float())
}

func multiply(this js.Value, args []js.Value) interface{} {
    if len(args) < 2 {
        return js.ValueOf("Error: At least 2 arguments required")
    }
    
    result := 1.0
    for _, arg := range args {
        result *= arg.Float()
    }
    
    return js.ValueOf(result)
}

func divide(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Exactly 2 arguments required")
    }
    
    divisor := args[1].Float()
    if divisor == 0 {
        return js.ValueOf("Error: Division by zero")
    }
    
    return js.ValueOf(args[0].Float() / divisor)
}

func power(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Exactly 2 arguments required")
    }
    
    base := args[0].Float()
    exponent := args[1].Float()
    
    result := 1.0
    for i := 0; i < int(exponent); i++ {
        result *= base
    }
    
    return js.ValueOf(result)
}
```

### 1.3.2 复杂数据结构导出

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "encoding/json"
    "fmt"
)

// User 用户结构体
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    IsActive bool   `json:"isActive"`
}

// UserManager 用户管理器
type UserManager struct {
    users []User
}

// NewUserManager 创建用户管理器
func NewUserManager() *UserManager {
    return &UserManager{
        users: make([]User, 0),
    }
}

// AddUser 添加用户
func (um *UserManager) AddUser(name, email string, age int) User {
    user := User{
        ID:       len(um.users) + 1,
        Name:     name,
        Email:    email,
        Age:      age,
        IsActive: true,
    }
    um.users = append(um.users, user)
    return user
}

// GetUser 获取用户
func (um *UserManager) GetUser(id int) (User, bool) {
    for _, user := range um.users {
        if user.ID == id {
            return user, true
        }
    }
    return User{}, false
}

// GetAllUsers 获取所有用户
func (um *UserManager) GetAllUsers() []User {
    return um.users
}

// 全局用户管理器
var userManager = NewUserManager()

func init() {
    // 导出用户管理函数
    js.Global().Set("userManager", js.ValueOf(map[string]interface{}{
        "addUser":    js.FuncOf(addUser),
        "getUser":    js.FuncOf(getUser),
        "getAllUsers": js.FuncOf(getAllUsers),
    }))
}

func addUser(this js.Value, args []js.Value) interface{} {
    if len(args) != 3 {
        return js.ValueOf("Error: Expected 3 arguments (name, email, age)")
    }
    
    name := args[0].String()
    email := args[1].String()
    age := args[2].Int()
    
    user := userManager.AddUser(name, email, age)
    
    // 转换为JSON
    userJSON, err := json.Marshal(user)
    if err != nil {
        return js.ValueOf("Error: Failed to serialize user")
    }
    
    return js.ValueOf(string(userJSON))
}

func getUser(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (id)")
    }
    
    id := args[0].Int()
    user, exists := userManager.GetUser(id)
    
    if !exists {
        return js.ValueOf("Error: User not found")
    }
    
    userJSON, err := json.Marshal(user)
    if err != nil {
        return js.ValueOf("Error: Failed to serialize user")
    }
    
    return js.ValueOf(string(userJSON))
}

func getAllUsers(this js.Value, args []js.Value) interface{} {
    users := userManager.GetAllUsers()
    
    usersJSON, err := json.Marshal(users)
    if err != nil {
        return js.ValueOf("Error: Failed to serialize users")
    }
    
    return js.ValueOf(string(usersJSON))
}
```

## 1.4 内存管理

### 1.4.1 内存分配和释放

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "unsafe"
)

// MemoryManager 内存管理器
type MemoryManager struct {
    allocated map[uintptr]int // 记录分配的内存大小
}

// NewMemoryManager 创建内存管理器
func NewMemoryManager() *MemoryManager {
    return &MemoryManager{
        allocated: make(map[uintptr]int),
    }
}

// Allocate 分配内存
func (mm *MemoryManager) Allocate(size int) uintptr {
    // 在Go中分配内存
    data := make([]byte, size)
    
    // 获取内存地址
    ptr := unsafe.Pointer(&data[0])
    addr := uintptr(ptr)
    
    // 记录分配
    mm.allocated[addr] = size
    
    return addr
}

// Free 释放内存
func (mm *MemoryManager) Free(addr uintptr) {
    delete(mm.allocated, addr)
}

// GetSize 获取分配的内存大小
func (mm *MemoryManager) GetSize(addr uintptr) int {
    return mm.allocated[addr]
}

// 全局内存管理器
var memoryManager = NewMemoryManager()

func init() {
    js.Global().Set("memory", js.ValueOf(map[string]interface{}{
        "allocate": js.FuncOf(allocateMemory),
        "free":     js.FuncOf(freeMemory),
        "getSize":  js.FuncOf(getMemorySize),
    }))
}

func allocateMemory(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (size)")
    }
    
    size := args[0].Int()
    addr := memoryManager.Allocate(size)
    
    return js.ValueOf(addr)
}

func freeMemory(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (address)")
    }
    
    addr := uintptr(args[0].Int())
    memoryManager.Free(addr)
    
    return js.ValueOf("OK")
}

func getMemorySize(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (address)")
    }
    
    addr := uintptr(args[0].Int())
    size := memoryManager.GetSize(addr)
    
    return js.ValueOf(size)
}
```

### 1.4.2 字符串内存管理

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "unsafe"
)

// StringManager 字符串管理器
type StringManager struct {
    strings map[uintptr]string
}

// NewStringManager 创建字符串管理器
func NewStringManager() *StringManager {
    return &StringManager{
        strings: make(map[uintptr]string),
    }
}

// StoreString 存储字符串
func (sm *StringManager) StoreString(s string) uintptr {
    // 将字符串转换为字节数组
    bytes := []byte(s)
    
    // 获取内存地址
    ptr := unsafe.Pointer(&bytes[0])
    addr := uintptr(ptr)
    
    // 存储字符串引用
    sm.strings[addr] = s
    
    return addr
}

// GetString 获取字符串
func (sm *StringManager) GetString(addr uintptr) (string, bool) {
    s, exists := sm.strings[addr]
    return s, exists
}

// FreeString 释放字符串
func (sm *StringManager) FreeString(addr uintptr) {
    delete(sm.strings, addr)
}

// 全局字符串管理器
var stringManager = NewStringManager()

func init() {
    js.Global().Set("stringManager", js.ValueOf(map[string]interface{}{
        "store": js.FuncOf(storeString),
        "get":   js.FuncOf(getString),
        "free":  js.FuncOf(freeString),
    }))
}

func storeString(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (string)")
    }
    
    s := args[0].String()
    addr := stringManager.StoreString(s)
    
    return js.ValueOf(addr)
}

func getString(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (address)")
    }
    
    addr := uintptr(args[0].Int())
    s, exists := stringManager.GetString(addr)
    
    if !exists {
        return js.ValueOf("Error: String not found")
    }
    
    return js.ValueOf(s)
}

func freeString(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (address)")
    }
    
    addr := uintptr(args[0].Int())
    stringManager.FreeString(addr)
    
    return js.ValueOf("OK")
}
```

## 1.5 类型转换

### 1.5.1 Go类型到JavaScript类型

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "encoding/json"
    "reflect"
)

// TypeConverter 类型转换器
type TypeConverter struct{}

// ConvertToJS 将Go值转换为JavaScript值
func (tc *TypeConverter) ConvertToJS(value interface{}) js.Value {
    if value == nil {
        return js.Null()
    }
    
    switch v := value.(type) {
    case bool:
        return js.ValueOf(v)
    case int:
        return js.ValueOf(v)
    case int8:
        return js.ValueOf(v)
    case int16:
        return js.ValueOf(v)
    case int32:
        return js.ValueOf(v)
    case int64:
        return js.ValueOf(v)
    case uint:
        return js.ValueOf(v)
    case uint8:
        return js.ValueOf(v)
    case uint16:
        return js.ValueOf(v)
    case uint32:
        return js.ValueOf(v)
    case uint64:
        return js.ValueOf(v)
    case float32:
        return js.ValueOf(v)
    case float64:
        return js.ValueOf(v)
    case string:
        return js.ValueOf(v)
    case []interface{}:
        return tc.convertSliceToJS(v)
    case map[string]interface{}:
        return tc.convertMapToJS(v)
    default:
        return tc.convertStructToJS(value)
    }
}

// convertSliceToJS 转换切片
func (tc *TypeConverter) convertSliceToJS(slice []interface{}) js.Value {
    jsArray := js.Global().Get("Array").New(len(slice))
    for i, item := range slice {
        jsArray.SetIndex(i, tc.ConvertToJS(item))
    }
    return jsArray
}

// convertMapToJS 转换映射
func (tc *TypeConverter) convertMapToJS(m map[string]interface{}) js.Value {
    jsObject := js.Global().Get("Object").New()
    for key, value := range m {
        jsObject.Set(key, tc.ConvertToJS(value))
    }
    return jsObject
}

// convertStructToJS 转换结构体
func (tc *TypeConverter) convertStructToJS(value interface{}) js.Value {
    // 使用反射获取结构体字段
    v := reflect.ValueOf(value)
    t := reflect.TypeOf(value)
    
    jsObject := js.Global().Get("Object").New()
    
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        fieldValue := v.Field(i)
        
        // 获取JSON标签
        jsonTag := field.Tag.Get("json")
        if jsonTag == "" {
            jsonTag = field.Name
        }
        
        jsObject.Set(jsonTag, tc.ConvertToJS(fieldValue.Interface()))
    }
    
    return jsObject
}

// 全局类型转换器
var typeConverter = &TypeConverter{}

func init() {
    js.Global().Set("typeConverter", js.ValueOf(map[string]interface{}{
        "convert": js.FuncOf(convertToJS),
    }))
}

func convertToJS(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument")
    }
    
    // 这里需要从JavaScript传递Go值，实际使用中可能需要序列化
    return js.ValueOf("Type conversion not implemented for this example")
}
```

### 1.5.2 JavaScript类型到Go类型

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "strconv"
)

// JSValueConverter JavaScript值转换器
type JSValueConverter struct{}

// ConvertFromJS 将JavaScript值转换为Go值
func (jvc *JSValueConverter) ConvertFromJS(jsValue js.Value) interface{} {
    switch jsValue.Type() {
    case js.TypeUndefined:
        return nil
    case js.TypeNull:
        return nil
    case js.TypeBoolean:
        return jsValue.Bool()
    case js.TypeNumber:
        return jsValue.Float()
    case js.TypeString:
        return jsValue.String()
    case js.TypeObject:
        return jvc.convertObjectFromJS(jsValue)
    case js.TypeFunction:
        return jsValue
    default:
        return jsValue
    }
}

// convertObjectFromJS 转换JavaScript对象
func (jvc *JSValueConverter) convertObjectFromJS(jsValue js.Value) interface{} {
    // 检查是否为数组
    if jsValue.Get("length").Type() != js.TypeUndefined {
        return jvc.convertArrayFromJS(jsValue)
    }
    
    // 转换为map
    result := make(map[string]interface{})
    
    // 获取对象的所有属性
    keys := js.Global().Get("Object").Call("keys", jsValue)
    for i := 0; i < keys.Length(); i++ {
        key := keys.Index(i).String()
        value := jsValue.Get(key)
        result[key] = jvc.ConvertFromJS(value)
    }
    
    return result
}

// convertArrayFromJS 转换JavaScript数组
func (jvc *JSValueConverter) convertArrayFromJS(jsValue js.Value) []interface{} {
    length := jsValue.Length()
    result := make([]interface{}, length)
    
    for i := 0; i < length; i++ {
        result[i] = jvc.ConvertFromJS(jsValue.Index(i))
    }
    
    return result
}

// ConvertToInt 转换为整数
func (jvc *JSValueConverter) ConvertToInt(jsValue js.Value) (int, error) {
    if jsValue.Type() == js.TypeNumber {
        return int(jsValue.Int()), nil
    }
    
    if jsValue.Type() == js.TypeString {
        return strconv.Atoi(jsValue.String())
    }
    
    return 0, js.ValueOf("Cannot convert to int")
}

// ConvertToFloat 转换为浮点数
func (jvc *JSValueConverter) ConvertToFloat(jsValue js.Value) (float64, error) {
    if jsValue.Type() == js.TypeNumber {
        return jsValue.Float(), nil
    }
    
    if jsValue.Type() == js.TypeString {
        return strconv.ParseFloat(jsValue.String(), 64)
    }
    
    return 0, js.ValueOf("Cannot convert to float")
}

// ConvertToBool 转换为布尔值
func (jvc *JSValueConverter) ConvertToBool(jsValue js.Value) bool {
    if jsValue.Type() == js.TypeBoolean {
        return jsValue.Bool()
    }
    
    if jsValue.Type() == js.TypeString {
        return jsValue.String() == "true"
    }
    
    if jsValue.Type() == js.TypeNumber {
        return jsValue.Int() != 0
    }
    
    return false
}

// 全局转换器
var jsValueConverter = &JSValueConverter{}

func init() {
    js.Global().Set("jsConverter", js.ValueOf(map[string]interface{}{
        "convert":    js.FuncOf(convertFromJS),
        "toInt":      js.FuncOf(convertToInt),
        "toFloat":    js.FuncOf(convertToFloat),
        "toBool":     js.FuncOf(convertToBool),
    }))
}

func convertFromJS(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument")
    }
    
    result := jsValueConverter.ConvertFromJS(args[0])
    return js.ValueOf(result)
}

func convertToInt(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument")
    }
    
    result, err := jsValueConverter.ConvertToInt(args[0])
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    return js.ValueOf(result)
}

func convertToFloat(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument")
    }
    
    result, err := jsValueConverter.ConvertToFloat(args[0])
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    return js.ValueOf(result)
}

func convertToBool(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument")
    }
    
    result := jsValueConverter.ConvertToBool(args[0])
    return js.ValueOf(result)
}
```

## 1.6 高级特性

### 1.6.1 异步操作支持

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "time"
)

// AsyncManager 异步管理器
type AsyncManager struct {
    promises map[string]js.Value
}

// NewAsyncManager 创建异步管理器
func NewAsyncManager() *AsyncManager {
    return &AsyncManager{
        promises: make(map[string]js.Value),
    }
}

// CreatePromise 创建Promise
func (am *AsyncManager) CreatePromise(id string) js.Value {
    promise := js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        resolve := args[0]
        reject := args[1]
        
        // 存储resolve和reject函数
        am.promises[id] = js.ValueOf(map[string]interface{}{
            "resolve": resolve,
            "reject":  reject,
        })
        
        return nil
    }))
    
    return promise
}

// ResolvePromise 解决Promise
func (am *AsyncManager) ResolvePromise(id string, value interface{}) {
    if promise, exists := am.promises[id]; exists {
        resolve := promise.Get("resolve")
        resolve.Invoke(js.ValueOf(value))
        delete(am.promises, id)
    }
}

// RejectPromise 拒绝Promise
func (am *AsyncManager) RejectPromise(id string, reason string) {
    if promise, exists := am.promises[id]; exists {
        reject := promise.Get("reject")
        reject.Invoke(js.ValueOf(reason))
        delete(am.promises, id)
    }
}

// 全局异步管理器
var asyncManager = NewAsyncManager()

func init() {
    js.Global().Set("asyncManager", js.ValueOf(map[string]interface{}{
        "createPromise": js.FuncOf(createPromise),
        "resolvePromise": js.FuncOf(resolvePromise),
        "rejectPromise": js.FuncOf(rejectPromise),
        "asyncOperation": js.FuncOf(asyncOperation),
    }))
}

func createPromise(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (id)")
    }
    
    id := args[0].String()
    promise := asyncManager.CreatePromise(id)
    
    return promise
}

func resolvePromise(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Expected 2 arguments (id, value)")
    }
    
    id := args[0].String()
    value := args[1]
    
    asyncManager.ResolvePromise(id, value)
    
    return js.ValueOf("OK")
}

func rejectPromise(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Expected 2 arguments (id, reason)")
    }
    
    id := args[0].String()
    reason := args[1].String()
    
    asyncManager.RejectPromise(id, reason)
    
    return js.ValueOf("OK")
}

func asyncOperation(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (id)")
    }
    
    id := args[0].String()
    
    // 创建Promise
    promise := asyncManager.CreatePromise(id)
    
    // 模拟异步操作
    go func() {
        time.Sleep(2 * time.Second)
        
        // 模拟操作结果
        result := map[string]interface{}{
            "id":      id,
            "result":  "Operation completed",
            "timestamp": time.Now().Unix(),
        }
        
        asyncManager.ResolvePromise(id, result)
    }()
    
    return promise
}
```

### 1.6.2 事件系统

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "sync"
)

// EventEmitter 事件发射器
type EventEmitter struct {
    listeners map[string][]js.Value
    mutex     sync.RWMutex
}

// NewEventEmitter 创建事件发射器
func NewEventEmitter() *EventEmitter {
    return &EventEmitter{
        listeners: make(map[string][]js.Value),
    }
}

// On 添加事件监听器
func (ee *EventEmitter) On(event string, listener js.Value) {
    ee.mutex.Lock()
    defer ee.mutex.Unlock()
    
    ee.listeners[event] = append(ee.listeners[event], listener)
}

// Off 移除事件监听器
func (ee *EventEmitter) Off(event string, listener js.Value) {
    ee.mutex.Lock()
    defer ee.mutex.Unlock()
    
    listeners := ee.listeners[event]
    for i, l := range listeners {
        if l.Equal(listener) {
            ee.listeners[event] = append(listeners[:i], listeners[i+1:]...)
            break
        }
    }
}

// Emit 发射事件
func (ee *EventEmitter) Emit(event string, args ...interface{}) {
    ee.mutex.RLock()
    listeners := ee.listeners[event]
    ee.mutex.RUnlock()
    
    for _, listener := range listeners {
        // 转换参数
        jsArgs := make([]interface{}, len(args))
        for i, arg := range args {
            jsArgs[i] = js.ValueOf(arg)
        }
        
        // 调用监听器
        listener.Invoke(jsArgs...)
    }
}

// 全局事件发射器
var eventEmitter = NewEventEmitter()

func init() {
    js.Global().Set("eventEmitter", js.ValueOf(map[string]interface{}{
        "on":   js.FuncOf(onEvent),
        "off":  js.FuncOf(offEvent),
        "emit": js.FuncOf(emitEvent),
    }))
}

func onEvent(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Expected 2 arguments (event, listener)")
    }
    
    event := args[0].String()
    listener := args[1]
    
    eventEmitter.On(event, listener)
    
    return js.ValueOf("OK")
}

func offEvent(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Expected 2 arguments (event, listener)")
    }
    
    event := args[0].String()
    listener := args[1]
    
    eventEmitter.Off(event, listener)
    
    return js.ValueOf("OK")
}

func emitEvent(this js.Value, args []js.Value) interface{} {
    if len(args) < 1 {
        return js.ValueOf("Error: Expected at least 1 argument (event)")
    }
    
    event := args[0].String()
    
    // 转换剩余参数
    eventArgs := make([]interface{}, len(args)-1)
    for i, arg := range args[1:] {
        eventArgs[i] = arg
    }
    
    eventEmitter.Emit(event, eventArgs...)
    
    return js.ValueOf("OK")
}
```

## 1.7 性能优化

### 1.7.1 内存优化

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "sync"
)

// MemoryPool 内存池
type MemoryPool struct {
    pools map[int]*sync.Pool
    mutex sync.RWMutex
}

// NewMemoryPool 创建内存池
func NewMemoryPool() *MemoryPool {
    return &MemoryPool{
        pools: make(map[int]*sync.Pool),
    }
}

// Get 获取内存块
func (mp *MemoryPool) Get(size int) []byte {
    mp.mutex.RLock()
    pool, exists := mp.pools[size]
    mp.mutex.RUnlock()
    
    if !exists {
        mp.mutex.Lock()
        pool = &sync.Pool{
            New: func() interface{} {
                return make([]byte, size)
            },
        }
        mp.pools[size] = pool
        mp.mutex.Unlock()
    }
    
    return pool.Get().([]byte)
}

// Put 归还内存块
func (mp *MemoryPool) Put(buf []byte) {
    size := cap(buf)
    mp.mutex.RLock()
    pool, exists := mp.pools[size]
    mp.mutex.RUnlock()
    
    if exists {
        // 重置切片长度
        buf = buf[:0]
        pool.Put(buf)
    }
}

// 全局内存池
var memoryPool = NewMemoryPool()

func init() {
    js.Global().Set("memoryPool", js.ValueOf(map[string]interface{}{
        "get": js.FuncOf(getMemory),
        "put": js.FuncOf(putMemory),
    }))
}

func getMemory(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (size)")
    }
    
    size := args[0].Int()
    buf := memoryPool.Get(size)
    
    // 返回内存地址（实际使用中需要更复杂的处理）
    return js.ValueOf(len(buf))
}

func putMemory(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument (buffer)")
    }
    
    // 实际使用中需要从JavaScript传递缓冲区
    return js.ValueOf("OK")
}
```

### 1.7.2 批量操作优化

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "sync"
)

// BatchProcessor 批量处理器
type BatchProcessor struct {
    batchSize int
    buffer    []interface{}
    mutex     sync.Mutex
    processor func([]interface{}) interface{}
}

// NewBatchProcessor 创建批量处理器
func NewBatchProcessor(batchSize int, processor func([]interface{}) interface{}) *BatchProcessor {
    return &BatchProcessor{
        batchSize: batchSize,
        buffer:    make([]interface{}, 0, batchSize),
        processor: processor,
    }
}

// Add 添加项目
func (bp *BatchProcessor) Add(item interface{}) interface{} {
    bp.mutex.Lock()
    defer bp.mutex.Unlock()
    
    bp.buffer = append(bp.buffer, item)
    
    if len(bp.buffer) >= bp.batchSize {
        return bp.flush()
    }
    
    return nil
}

// Flush 刷新缓冲区
func (bp *BatchProcessor) Flush() interface{} {
    bp.mutex.Lock()
    defer bp.mutex.Unlock()
    return bp.flush()
}

// flush 内部刷新方法
func (bp *BatchProcessor) flush() interface{} {
    if len(bp.buffer) == 0 {
        return nil
    }
    
    batch := make([]interface{}, len(bp.buffer))
    copy(batch, bp.buffer)
    bp.buffer = bp.buffer[:0]
    
    return bp.processor(batch)
}

// 全局批量处理器
var batchProcessor = NewBatchProcessor(10, func(items []interface{}) interface{} {
    // 处理批量数据
    result := make([]interface{}, len(items))
    for i, item := range items {
        result[i] = map[string]interface{}{
            "processed": true,
            "data":      item,
            "index":     i,
        }
    }
    return result
})

func init() {
    js.Global().Set("batchProcessor", js.ValueOf(map[string]interface{}{
        "add":   js.FuncOf(addToBatch),
        "flush": js.FuncOf(flushBatch),
    }))
}

func addToBatch(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: Expected 1 argument")
    }
    
    result := batchProcessor.Add(args[0])
    if result != nil {
        return js.ValueOf(result)
    }
    
    return js.ValueOf("Added to batch")
}

func flushBatch(this js.Value, args []js.Value) interface{} {
    result := batchProcessor.Flush()
    if result != nil {
        return js.ValueOf(result)
    }
    
    return js.ValueOf("Batch is empty")
}
```

## 1.8 实际应用

### 1.8.1 图像处理

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "image"
    "image/color"
)

// ImageProcessor 图像处理器
type ImageProcessor struct{}

// ProcessImage 处理图像
func (ip *ImageProcessor) ProcessImage(imgData []byte, width, height int) []byte {
    // 创建图像
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    
    // 处理像素数据
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            // 获取像素索引
            idx := (y*width + x) * 4
            
            if idx+3 < len(imgData) {
                // 读取RGBA值
                r := imgData[idx]
                g := imgData[idx+1]
                b := imgData[idx+2]
                a := imgData[idx+3]
                
                // 应用滤镜（例如：灰度化）
                gray := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
                
                // 设置像素
                img.Set(x, y, color.RGBA{gray, gray, gray, a})
            }
        }
    }
    
    // 返回处理后的数据
    return img.Pix
}

// 全局图像处理器
var imageProcessor = &ImageProcessor{}

func init() {
    js.Global().Set("imageProcessor", js.ValueOf(map[string]interface{}{
        "process": js.FuncOf(processImage),
    }))
}

func processImage(this js.Value, args []js.Value) interface{} {
    if len(args) != 3 {
        return js.ValueOf("Error: Expected 3 arguments (data, width, height)")
    }
    
    // 获取参数
    data := args[0]
    width := args[1].Int()
    height := args[2].Int()
    
    // 转换数据
    imgData := make([]byte, data.Length())
    for i := 0; i < data.Length(); i++ {
        imgData[i] = byte(data.Index(i).Int())
    }
    
    // 处理图像
    result := imageProcessor.ProcessImage(imgData, width, height)
    
    // 返回结果
    jsResult := js.Global().Get("Array").New(len(result))
    for i, b := range result {
        jsResult.SetIndex(i, js.ValueOf(b))
    }
    
    return jsResult
}
```

### 1.8.2 数据加密

```go
//go:build js,wasm

package main

import (
    "syscall/js"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

// CryptoManager 加密管理器
type CryptoManager struct{}

// Encrypt 加密数据
func (cm *CryptoManager) Encrypt(data []byte, key []byte) ([]byte, error) {
    // 创建AES cipher
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    // 创建GCM
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // 生成随机nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    // 加密数据
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    
    return ciphertext, nil
}

// Decrypt 解密数据
func (cm *CryptoManager) Decrypt(data []byte, key []byte) ([]byte, error) {
    // 创建AES cipher
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    // 创建GCM
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    // 提取nonce
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, js.ValueOf("ciphertext too short")
    }
    
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    
    // 解密数据
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }
    
    return plaintext, nil
}

// 全局加密管理器
var cryptoManager = &CryptoManager{}

func init() {
    js.Global().Set("cryptoManager", js.ValueOf(map[string]interface{}{
        "encrypt": js.FuncOf(encryptData),
        "decrypt": js.FuncOf(decryptData),
    }))
}

func encryptData(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Expected 2 arguments (data, key)")
    }
    
    // 获取参数
    data := args[0]
    key := args[1]
    
    // 转换数据
    dataBytes := make([]byte, data.Length())
    for i := 0; i < data.Length(); i++ {
        dataBytes[i] = byte(data.Index(i).Int())
    }
    
    keyBytes := make([]byte, key.Length())
    for i := 0; i < key.Length(); i++ {
        keyBytes[i] = byte(key.Index(i).Int())
    }
    
    // 加密
    result, err := cryptoManager.Encrypt(dataBytes, keyBytes)
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    // 返回base64编码的结果
    return js.ValueOf(base64.StdEncoding.EncodeToString(result))
}

func decryptData(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: Expected 2 arguments (data, key)")
    }
    
    // 获取参数
    data := args[0].String()
    key := args[1]
    
    // 解码base64
    dataBytes, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    // 转换密钥
    keyBytes := make([]byte, key.Length())
    for i := 0; i < key.Length(); i++ {
        keyBytes[i] = byte(key.Index(i).Int())
    }
    
    // 解密
    result, err := cryptoManager.Decrypt(dataBytes, keyBytes)
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    // 返回结果
    jsResult := js.Global().Get("Array").New(len(result))
    for i, b := range result {
        jsResult.SetIndex(i, js.ValueOf(b))
    }
    
    return jsResult
}
```

---

**总结**: Go WebAssembly导出功能为Go语言提供了强大的跨平台能力，特别是在前端开发和边缘计算领域。通过合理的类型转换、内存管理和性能优化，可以实现高效的WASM模块，为现代Web应用提供强大的后端计算能力。
