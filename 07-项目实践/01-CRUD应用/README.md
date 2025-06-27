# CRUD应用项目模板

## 📚 **项目概述**

这是一个完整的CRUD（Create, Read, Update, Delete）应用项目模板，展示了如何使用Go语言构建一个功能完整的Web应用。项目采用现代化的技术栈和最佳实践，适合作为学习和实际开发的参考。

## 🎯 **项目特色**

### **技术栈**

- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **认证**: JWT
- **日志**: Zap
- **配置**: Viper

### **架构特点**

- **分层架构**: 清晰的分层设计
- **依赖注入**: 使用接口和依赖注入
- **错误处理**: 统一的错误处理机制
- **中间件**: 可复用的中间件组件
- **测试**: 完整的单元测试和集成测试

## 📋 **项目结构**

```text
crud-app/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── database/
│   │   └── database.go          # 数据库连接
│   ├── models/
│   │   └── user.go              # 数据模型
│   ├── handlers/
│   │   └── user_handler.go      # HTTP处理器
│   ├── services/
│   │   └── user_service.go      # 业务逻辑
│   ├── repositories/
│   │   └── user_repository.go   # 数据访问层
│   └── middleware/
│       ├── auth.go              # 认证中间件
│       ├── cors.go              # CORS中间件
│       └── logging.go           # 日志中间件
├── pkg/
│   ├── errors/
│   │   └── errors.go            # 错误定义
│   ├── utils/
│   │   └── utils.go             # 工具函数
│   └── validators/
│       └── validators.go        # 数据验证
├── api/
│   └── docs/                    # API文档
├── scripts/
│   ├── build.sh                 # 构建脚本
│   └── deploy.sh                # 部署脚本
├── tests/
│   ├── unit/                    # 单元测试
│   └── integration/             # 集成测试
├── docker/
│   ├── Dockerfile               # Docker镜像
│   └── docker-compose.yml       # 容器编排
├── configs/
│   ├── config.yaml              # 配置文件
│   └── config.prod.yaml         # 生产配置
├── go.mod                       # Go模块文件
├── go.sum                       # 依赖校验
└── README.md                    # 项目文档
```

## 🚀 **快速开始**

### **环境要求**

- Go 1.21+
- PostgreSQL 13+
- Redis 6+
- Docker (可选)

### **安装依赖**

```bash
# 克隆项目
git clone <repository-url>
cd crud-app

# 安装依赖
go mod download

# 设置环境变量
export DATABASE_URL="postgres://user:password@localhost:5432/crud_app"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="your-secret-key"
```

### **运行项目**

```bash
# 开发模式
go run cmd/server/main.go

# 构建项目
go build -o bin/server cmd/server/main.go

# 运行测试
go test ./...

# 使用Docker
docker-compose up -d
```

## 📊 **API接口**

### **用户管理接口**

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/api/v1/users` | 创建用户 | 否 |
| GET | `/api/v1/users` | 获取用户列表 | 是 |
| GET | `/api/v1/users/:id` | 获取用户详情 | 是 |
| PUT | `/api/v1/users/:id` | 更新用户 | 是 |
| DELETE | `/api/v1/users/:id` | 删除用户 | 是 |

### **认证接口**

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/refresh` | 刷新令牌 |

## 💻 **代码示例**

### **主程序入口**

```go
// cmd/server/main.go
package main

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/your-username/crud-app/internal/config"
    "github.com/your-username/crud-app/internal/database"
    "github.com/your-username/crud-app/internal/handlers"
    "github.com/your-username/crud-app/internal/middleware"
)

func main() {
    // 加载配置
    cfg := config.Load()
    
    // 初始化数据库
    db := database.Init(cfg.DatabaseURL)
    
    // 创建Gin引擎
    r := gin.Default()
    
    // 添加中间件
    r.Use(middleware.CORS())
    r.Use(middleware.Logging())
    
    // 设置路由
    setupRoutes(r, db)
    
    // 启动服务器
    log.Printf("Server starting on %s", cfg.ServerPort)
    http.ListenAndServe(":"+cfg.ServerPort, r)
}

func setupRoutes(r *gin.Engine, db *gorm.DB) {
    // API v1 路由组
    v1 := r.Group("/api/v1")
    {
        // 认证路由
        auth := v1.Group("/auth")
        {
            auth.POST("/register", handlers.Register)
            auth.POST("/login", handlers.Login)
            auth.POST("/refresh", handlers.RefreshToken)
        }
        
        // 用户路由 (需要认证)
        users := v1.Group("/users")
        users.Use(middleware.Auth())
        {
            users.GET("", handlers.GetUsers)
            users.GET("/:id", handlers.GetUser)
            users.POST("", handlers.CreateUser)
            users.PUT("/:id", handlers.UpdateUser)
            users.DELETE("/:id", handlers.DeleteUser)
        }
    }
}
```

## 🧪 **测试示例**

### **单元测试**

```go
// tests/unit/user_service_test.go
package unit

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/your-username/crud-app/internal/models"
    "github.com/your-username/crud-app/internal/services"
)

func TestUserService_CreateUser(t *testing.T) {
    // 测试用例
    req := &models.CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
        Age:      25,
    }
    
    // 执行测试
    user, err := service.CreateUser(req)
    
    // 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, req.Username, user.Username)
}
```

## 🚀 **部署指南**

### **Docker部署**

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./main"]
```

## 📊 **性能优化**

### **数据库优化**

- 使用连接池
- 添加适当的索引
- 优化查询语句

### **缓存策略**

- Redis缓存热点数据
- 实现缓存失效机制
- 使用缓存预热

## 🛡️ **安全考虑**

### **认证授权**

- JWT令牌认证
- 密码加密存储
- 权限控制

### **数据验证**

- 输入数据验证
- SQL注入防护
- XSS攻击防护

---

**项目维护者**: AI Assistant  
**最后更新**: 2024年6月27日  
**项目状态**: 开发中
