# API 版本控制 (API Versioning)

> **分类**: 成熟应用领域  
> **标签**: #api #versioning #backward-compatibility

---

## 版本策略

### URL 路径版本

```
/api/v1/users
/api/v2/users
```

```go
func SetupRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", userHandlerV1.List)
        v1.POST("/users", userHandlerV1.Create)
    }
    
    v2 := r.Group("/api/v2")
    {
        v2.GET("/users", userHandlerV2.List)
        v2.POST("/users", userHandlerV2.Create)
    }
}
```

### Header 版本

```
Accept: application/vnd.api+json;version=2
X-API-Version: 2
```

```go
func VersionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        version := c.GetHeader("X-API-Version")
        if version == "" {
            version = "1"
        }
        
        c.Set("api_version", version)
        c.Next()
    }
}

func GetUser(c *gin.Context) {
    version := c.GetString("api_version")
    
    switch version {
    case "1":
        getUserV1(c)
    case "2":
        getUserV2(c)
    default:
        c.JSON(400, gin.H{"error": "unsupported version"})
    }
}
```

---

## 版本适配器

```go
// v1 响应
type UserV1 struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// v2 响应
type UserV2 struct {
    ID        string    `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// 适配器
func AdaptToV1(u *UserV2) *UserV1 {
    return &UserV1{
        ID:    parseID(u.ID),
        Name:  u.FirstName + " " + u.LastName,
        Email: u.Email,
    }
}
```

---

## 弃用策略

```go
func DeprecatedMiddleware(version, sunset string) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Deprecation", "true")
        c.Header("Sunset", sunset)
        c.Header("Link", fmt.Sprintf("</api/v%d%s>; rel=successor-version", 
            parseVersion(version)+1, c.Request.URL.Path))
        
        c.Next()
    }
}

// 使用
v1.GET("/legacy", DeprecatedMiddleware("1", "Sat, 01 Jun 2024 00:00:00 GMT"), legacyHandler)
```

---

## 版本协商

```go
func ContentNegotiation() gin.HandlerFunc {
    return func(c *gin.Context) {
        accept := c.GetHeader("Accept")
        
        // 解析 Accept 头
        mediaType, params, _ := mime.ParseMediaType(accept)
        
        // 检查版本
        if version := params["version"]; version != "" {
            c.Set("api_version", version)
        }
        
        c.Next()
    }
}
```

---

## 最佳实践

1. **主版本变化**: 破坏性变更
2. **次版本变化**: 向后兼容的新功能
3. **保持 N-1 兼容**: 至少支持两个版本
4. **文档**: 清晰的版本变更日志
5. **监控**: 各版本使用量统计
