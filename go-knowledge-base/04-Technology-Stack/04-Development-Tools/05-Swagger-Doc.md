# Swagger 文档

> **分类**: 开源技术堆栈

---

## swaggo

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

---

## 注释格式

```go
// GetUser godoc
// @Summary      Get a user
// @Description  get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  User
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Router       /users/{id} [get]
func GetUser(c *gin.Context) {
    id := c.Param("id")
    // ...
}
```

---

## 生成文档

```bash
# 生成 docs 目录
swag init

# 指定主文件
swag init -g cmd/main.go
```

---

## Gin 集成

```go
import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "project/docs"
)

func main() {
    r := gin.Default()
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    r.Run()
}
```

---

## 访问

```
http://localhost:8080/swagger/index.html
```
