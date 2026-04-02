# API 文档 (API Documentation)

> **分类**: 开源技术堆栈  
> **标签**: #api-docs #openapi #swagger

---

## OpenAPI 生成

### 从代码生成

```go
// 使用 swaggo
type CreateUserRequest struct {
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john@example.com"`
}

type UserResponse struct {
    ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// @Summary      Create user
// @Description  Create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      CreateUserRequest  true  "User info"
// @Success      201      {object}  UserResponse
// @Failure      400      {object}  ErrorResponse
// @Router       /users [post]
func CreateUser(c *gin.Context) {
    // ...
}
```

```bash
# 生成文档
swag init

# 运行
swag init -g cmd/main.go -o docs
```

---

## 自定义文档

```go
type APIDocumentation struct {
    OpenAPI string                   `json:"openapi"`
    Info    Info                     `json:"info"`
    Paths   map[string]PathItem      `json:"paths"`
}

type Info struct {
    Title       string `json:"title"`
    Version     string `json:"version"`
    Description string `json:"description"`
}

func GenerateAPIDoc(routes []Route) *APIDocumentation {
    doc := &APIDocumentation{
        OpenAPI: "3.0.0",
        Info: Info{
            Title:       "My API",
            Version:     "1.0.0",
            Description: "API documentation",
        },
        Paths: make(map[string]PathItem),
    }
    
    for _, route := range routes {
        doc.Paths[route.Path] = PathItem{
            route.Method: Operation{
                Summary:     route.Summary,
                Description: route.Description,
                Parameters:  route.Parameters,
                Responses:   route.Responses,
            },
        }
    }
    
    return doc
}
```

---

## 文档 UI

```go
import swaggerFiles "github.com/swaggo/files"
import ginSwagger "github.com/swaggo/gin-swagger"

func main() {
    r := gin.Default()
    
    // Swagger UI
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // ReDoc
    r.GET("/docs", func(c *gin.Context) {
        c.Data(200, "text/html", []byte(redocHTML))
    })
    
    r.Run()
}

const redocHTML = `
<!DOCTYPE html>
<html>
  <head>
    <title>API Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
      body { margin: 0; padding: 0; }
    </style>
  </head>
  <body>
    <redoc spec-url="/swagger/doc.json"></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
  </body>
</html>
`
```

---

## 文档测试

```go
// 使用文档测试 API
func TestAPIWithDocs(t *testing.T) {
    // 加载 OpenAPI 规范
    doc, _ := openapi3.NewLoader().LoadFromFile("docs/swagger.json")
    
    // 验证路由
    router := setupRouter()
    
    for path, pathItem := range doc.Paths {
        for method, operation := range pathItem.Operations() {
            t.Run(fmt.Sprintf("%s %s", method, path), func(t *testing.T) {
                // 构造请求
                req := httptest.NewRequest(method, path, nil)
                w := httptest.NewRecorder()
                
                router.ServeHTTP(w, req)
                
                // 验证响应码
                validCodes := make(map[int]bool)
                for code := range operation.Responses {
                    if c, err := strconv.Atoi(code); err == nil {
                        validCodes[c] = true
                    }
                }
                
                if !validCodes[w.Code] {
                    t.Errorf("Invalid status code: %d", w.Code)
                }
            })
        }
    }
}
```
