# encoding/json 详解

> **分类**: 开源技术堆栈

---

## 序列化

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

p := Person{Name: "Alice", Age: 30}
data, _ := json.Marshal(p)
// {"name":"Alice","age":30}
```

---

## 反序列化

```go
var p Person
json.Unmarshal(data, &p)
```

---

## Tag 选项

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name,omitempty"`  // 空值省略
    Email     string    `json:"-"`               // 忽略
    CreatedAt time.Time `json:"created_at"`
}
```

---

## 流处理

```go
// Encoder
encoder := json.NewEncoder(writer)
encoder.Encode(value)

// Decoder
decoder := json.NewDecoder(reader)
decoder.Decode(&value)
```

---

## 自定义类型

```go
type MyTime time.Time

func (t MyTime) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Time(t).Format("2006-01-02"))
}
```
