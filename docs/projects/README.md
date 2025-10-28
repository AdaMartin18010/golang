# Go项目实战

Go项目实战，包含示例项目、教程和项目模板。

---

## 📚 目录结构

### 核心模块

1. **[示例项目](./examples/README.md)** ⭐⭐⭐⭐⭐
   - Web应用示例
   - 微服务示例
   - CLI工具示例
   - 库项目示例
   - 实战案例

2. **[教程项目](./tutorials/README.md)** ⭐⭐⭐⭐⭐
   - 入门教程
   - 进阶教程
   - 专家教程
   - 实战演练

3. **[项目模板](./templates/README.md)** ⭐⭐⭐⭐⭐
   - 项目结构模板
   - 微服务模板
   - Web应用模板
   - CLI工具模板
   - 库项目模板

---

## 🎯 项目类型

### Web应用 (1-2周)
```
博客系统 → API服务 → 电商后端 → 实时聊天
```

### 微服务 (2-3周)
```
用户服务 → 订单服务 → 支付服务 → 网关
```

### CLI工具 (1周)
```
文件处理 → 系统监控 → 代码生成 → 部署工具
```

---

## 🚀 快速开始

### 使用项目模板

```bash
# 克隆模板
git clone https://github.com/your-org/go-web-template.git myproject
cd myproject

# 初始化
go mod init github.com/username/myproject
go mod tidy

# 运行
go run main.go
```

### 项目结构

```
myproject/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   └── repository/
├── pkg/
│   └── util/
├── api/
│   └── proto/
├── configs/
├── scripts/
├── docs/
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── README.md
```

---

## 📖 示例项目

### 1. RESTful API

```go
// main.go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    // 路由
    r.GET("/users", getUsers)
    r.POST("/users", createUser)
    r.GET("/users/:id", getUser)
    r.PUT("/users/:id", updateUser)
    r.DELETE("/users/:id", deleteUser)
    
    r.Run(":8080")
}
```

### 2. 微服务

```go
// user-service/main.go
package main

import (
    "google.golang.org/grpc"
    pb "myapp/proto/user"
)

func main() {
    lis, _ := net.Listen("tcp", ":50051")
    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &server{})
    s.Serve(lis)
}
```

### 3. CLI工具

```go
// main.go
package main

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "mytool",
    Short: "A brief description",
    Run: func(cmd *cobra.Command, args []string) {
        // 主逻辑
    },
}

func main() {
    rootCmd.Execute()
}
```

---

## 🛠️ 开发工具

### 项目管理
- Makefile - 构建自动化
- Docker Compose - 本地开发
- air - 热重载

### 代码质量
- golangci-lint - 代码检查
- gofmt - 格式化
- go vet - 代码分析

---

## 📚 学习资源

- [项目示例索引](./examples/00-示例索引.md)
- [教程项目索引](./tutorials/00-教程索引.md)
- [模板使用指南](./templates/06-快速开始指南.md)

---

## 🔗 相关资源

- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [Awesome Go Projects](https://github.com/avelino/awesome-go#project-layout)

---

**最后更新**: 2025-10-28  
**Go版本**: 1.25.3

