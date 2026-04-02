# 项目结构 (Project Structure)

> **分类**: 工程与云原生
> **标签**: #project-structure #layout

---

## 标准布局

```
myapp/
├── cmd/                    # 可执行程序入口
│   ├── server/
│   │   └── main.go        # HTTP 服务
│   └── worker/
│       └── main.go        # 后台任务
│
├── internal/              # 私有代码
│   ├── domain/            # 领域模型
│   │   ├── user.go
│   │   └── order.go
│   ├── repository/        # 数据访问
│   ├── service/           # 业务逻辑
│   └── handler/           # HTTP 处理
│
├── pkg/                   # 公开库（可被外部使用）
│   ├── logger/
│   ├── validator/
│   └── errors/
│
├── api/                   # API 定义
│   ├── proto/             # Protocol Buffers
│   └── openapi/           # OpenAPI/Swagger
│
├── web/                   # 前端静态文件
│
├── configs/               # 配置文件
│   ├── config.yaml
│   └── config.prod.yaml
│
├── scripts/               # 脚本
│   ├── build.sh
│   └── migrate.sh
│
├── deployments/           # 部署配置
│   ├── docker/
│   └── k8s/
│
├── docs/                  # 文档
│
├── test/                  # 测试数据和集成测试
│
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

---

## 按功能分层

```
internal/
├── user/                  # 用户模块
│   ├── domain.go          # 实体、值对象
│   ├── repository.go      # 接口定义
│   ├── repository_pg.go   # PostgreSQL 实现
│   ├── service.go         # 业务逻辑
│   ├── handler.go         # HTTP 处理
│   └── dto.go             # 数据传输对象
│
├── order/                 # 订单模块
│   ├── domain.go
│   ├── repository.go
│   └── service.go
```

---

## Clean Architecture

```
internal/
├── entities/              # 核心领域模型
│   └── user.go
│
├── usecases/              # 业务用例
│   └── user_usecase.go
│
├── interface/             # 适配器接口
│   ├── repository/
│   │   └── user_repo.go
│   └── presenter/
│       └── user_presenter.go
│
├── infrastructure/        # 具体实现
│   ├── database/
│   │   └── pg_user_repo.go
│   └── web/
│       └── user_handler.go
```

---

## 包命名规范

| 包名 | 用途 | 示例 |
|------|------|------|
| `cmd/` | 程序入口 | `cmd/server/main.go` |
| `internal/` | 私有代码 | `internal/service/` |
| `pkg/` | 公开库 | `pkg/logger/` |
| `api/` | API 定义 | `api/v1/` |
| `web/` | 前端资源 | `web/static/` |
| `configs/` | 配置 | `configs/app.yaml` |

---

## 导入组织

```go
import (
    // 标准库
    "context"
    "fmt"
    "time"

    // 第三方库
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
    "go.uber.org/zap"

    // 内部包
    "myapp/internal/domain"
    "myapp/internal/service"
    "myapp/pkg/logger"
)
```
