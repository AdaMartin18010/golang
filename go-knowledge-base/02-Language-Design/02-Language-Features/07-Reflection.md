# 反射 (Reflection)

> **分类**: 语言设计

---

## reflect 包

Go 通过 `reflect` 包提供运行时类型检查和操作。

```go
import "reflect"
```

---

## 基本操作

### 获取类型和值

```go
x := 42
v := reflect.ValueOf(x)    // Value
t := reflect.TypeOf(x)     // Type

fmt.Println(v.Kind())      // int
fmt.Println(t.Name())      // int
```

### 修改值

```go
x := 42
v := reflect.ValueOf(&x)   // 必须传指针
v.Elem().SetInt(100)
fmt.Println(x)             // 100
```

---

## 结构体反射

### 遍历字段

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

u := User{Name: "Alice", Age: 30}
t := reflect.TypeOf(u)

for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    fmt.Println(field.Name, field.Tag.Get("json"))
}
```

### 创建实例

```go
t := reflect.TypeOf(User{})
v := reflect.New(t)  // 创建指针

user := v.Interface().(*User)
```

---

## 方法调用

```go
type Calculator struct{}

func (c Calculator) Add(a, b int) int {
    return a + b
}

c := Calculator{}
v := reflect.ValueOf(c)

method := v.MethodByName("Add")
args := []reflect.Value{
    reflect.ValueOf(1),
    reflect.ValueOf(2),
}
result := method.Call(args)
fmt.Println(result[0].Int())  // 3
```

---

## 常见用途

### 1. JSON 序列化

```go
// encoding/json 内部使用反射
json.Marshal(user)
```

### 2. 深拷贝

```go
func deepCopy(dst, src interface{}) error {
    // 使用反射递归复制
}
```

### 3. 依赖注入

```go
// 根据类型自动注入依赖
```

---

## 性能注意

**反射比直接调用慢 10-100 倍**

```go
// 直接调用: 快
x := user.Name

// 反射: 慢
v := reflect.ValueOf(user)
x := v.FieldByName("Name").String()
```

---

## 最佳实践

1. **避免热路径使用反射**
2. **优先使用类型断言**
3. **缓存反射结果**

```go
// 缓存类型信息
var userType = reflect.TypeOf(User{})

func processMany(users []User) {
    for _, u := range users {
        v := reflect.ValueOf(u)
        // 使用 userType 信息
    }
}
```
