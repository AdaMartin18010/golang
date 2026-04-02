# RESTful API 开发

> **分类**: 成熟应用领域

---

## 设计原则

### URL 设计

```
GET    /users          # 列表
GET    /users/:id      # 详情
POST   /users          # 创建
PUT    /users/:id      # 更新
PATCH  /users/:id      # 部分更新
DELETE /users/:id      # 删除
```

### 状态码

| 状态码 | 用途 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 不存在 |
| 500 | 服务器错误 |

---

## Gin 实现

```go
func main() {
    r := gin.Default()

    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", listUsers)
        v1.GET("/users/:id", getUser)
        v1.POST("/users", createUser)
        v1.PUT("/users/:id", updateUser)
        v1.DELETE("/users/:id", deleteUser)
    }

    r.Run()
}
```

---

## 统一响应

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, code int, data interface{}) {
    c.JSON(200, Response{
        Code:    code,
        Message: getMessage(code),
        Data:    data,
    })
}
```
