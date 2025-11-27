# 数据转换工具

框架级别的数据转换工具，提供各种数据格式和类型之间的转换能力。

## 📋 功能特性

- ✅ **类型转换**: 字符串、整数、浮点数、布尔值、时间等
- ✅ **JSON 转换**: JSON 序列化和反序列化
- ✅ **Map 转换**: 结构体到 Map 的转换
- ✅ **Slice 转换**: 数组/切片转换
- ✅ **通用转换**: 基于反射的通用类型转换

## 🚀 快速开始

### 基本类型转换

```go
import "github.com/yourusername/golang/pkg/converter"

conv := converter.NewConverter()

// 转换为字符串
str := conv.ToString(123)        // "123"
str = conv.ToString(true)        // "true"
str = conv.ToString(time.Now())  // "2025-01-01T00:00:00Z"

// 转换为整数
num, _ := conv.ToInt("123")      // 123
num, _ = conv.ToInt(123.45)      // 123

// 转换为浮点数
f, _ := conv.ToFloat64("123.45") // 123.45

// 转换为布尔值
b, _ := conv.ToBool("true")      // true
b, _ = conv.ToBool(1)            // true
```

### JSON 转换

```go
// 转换为 JSON
data := map[string]interface{}{
    "name": "John",
    "age":  30,
}
jsonStr, _ := conv.ToJSON(data)

// 从 JSON 解析
var result map[string]interface{}
conv.FromJSON(jsonStr, &result)
```

### Map 和 Slice 转换

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{ID: 1, Name: "John", Email: "john@example.com"}

// 转换为 Map
m, _ := conv.ToMap(user)
// map[string]interface{}{
//     "id": 1,
//     "name": "John",
//     "email": "john@example.com",
// }

// 转换为 Slice
slice := []int{1, 2, 3}
s, _ := conv.ToSlice(slice)
```

### 通用转换

```go
// 使用反射进行类型转换
targetType := reflect.TypeOf(int64(0))
result, _ := conv.Convert("123", targetType)
// result 是 int64 类型的 123
```

## 📚 API 参考

### Converter 接口

```go
type Converter interface {
    ToString(v interface{}) string
    ToInt(v interface{}) (int, error)
    ToInt64(v interface{}) (int64, error)
    ToFloat64(v interface{}) (float64, error)
    ToBool(v interface{}) (bool, error)
    ToTime(v interface{}) (time.Time, error)
    ToJSON(v interface{}) (string, error)
    FromJSON(data string, v interface{}) error
    ToMap(v interface{}) (map[string]interface{}, error)
    ToSlice(v interface{}) ([]interface{}, error)
    Convert(v interface{}, targetType reflect.Type) (interface{}, error)
}
```

## 🎯 使用场景

1. **API 数据转换**: 请求/响应数据转换
2. **配置解析**: 配置文件数据转换
3. **数据验证**: 类型转换和验证
4. **序列化/反序列化**: 数据格式转换

## 🔗 相关文档

- [反射/自解释能力](../reflect/README.md)
