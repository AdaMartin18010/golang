# Clean Architecture

本项目采用 Clean Architecture（整洁架构）设计。

## 架构层次

### 1. Domain Layer (领域层)

**位置**: `internal/domain/`

**职责**:

- 核心业务逻辑
- 领域实体和值对象
- 领域服务接口
- 仓储接口

**规则**:

- 不依赖任何外部框架
- 不依赖 Infrastructure 或 Interfaces 层
- 只包含业务逻辑

**示例**:

```go
// internal/domain/user/entity.go
type User struct {
    ID    string
    Email string
    Name  string
}
```

### 2. Application Layer (应用层)

**位置**: `internal/application/`

**职责**:

- 用例编排
- 协调领域对象
- 应用服务
- DTO（数据传输对象）

**规则**:

- 只能导入 Domain 层
- 不依赖 Infrastructure 或 Interfaces 层

**示例**:

```go
// internal/application/user/service.go
type Service struct {
    repo domain.UserRepository
}
```

### 3. Infrastructure Layer (基础设施层)

**位置**: `internal/infrastructure/`

**职责**:

- 技术实现细节
- 数据库访问
- 消息队列
- 外部服务集成
- 可观测性（OTLP）

**规则**:

- 实现 Domain 层定义的接口
- 可以依赖外部库

**示例**:

```go
// internal/infrastructure/database/postgres/user_repository.go
type UserRepository struct {
    db *ent.Client
}
```

### 4. Interfaces Layer (接口层)

**位置**: `internal/interfaces/`

**职责**:

- 外部接口适配
- HTTP 处理器
- gRPC 服务
- GraphQL 解析器
- MQTT 处理器

**规则**:

- 调用 Application 层
- 处理请求/响应转换

**示例**:

```go
// internal/interfaces/http/chi/handlers/user_handler.go
type UserHandler struct {
    service *application.UserService
}
```

## 依赖方向

```text
Interfaces → Application → Domain
     ↓            ↓
Infrastructure → Domain
```

## 优势

1. **独立性**: 业务逻辑不依赖框架
2. **可测试性**: 每层都可以独立测试
3. **可维护性**: 清晰的职责分离
4. **可扩展性**: 易于添加新功能
