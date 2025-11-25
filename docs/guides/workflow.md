# 工作流指南

## Temporal 工作流

本项目使用 [Temporal](https://temporal.io/) 作为工作流编排引擎，用于构建可靠的分布式应用。

## 架构

### 组件

1. **Temporal Server** - 工作流引擎
2. **Worker** - 执行工作流和活动
3. **Client** - 启动和查询工作流

### 工作流类型

- **UserWorkflow** - 用户相关操作工作流
  - 创建用户
  - 更新用户
  - 删除用户

## 快速开始

### 1. 启动 Temporal Server

使用 Docker Compose：

```bash
cd deployments/docker
docker-compose up -d temporal temporal-db temporal-ui
```

### 2. 启动 Worker

```bash
go run ./cmd/temporal-worker
```

### 3. 启动应用

```bash
go run ./cmd/server
```

### 4. 访问 Temporal UI

打开浏览器访问：http://localhost:8088

## 使用工作流

### 启动工作流

```bash
curl -X POST http://localhost:8080/api/v1/workflows/user \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "email": "test@example.com",
    "name": "Test User",
    "action": "create"
  }'
```

### 查询工作流结果

```bash
curl http://localhost:8080/api/v1/workflows/user/{workflow_id}?run_id={run_id}
```

## 工作流定义

### UserWorkflow

用户工作流处理用户相关的操作，包括：

1. **创建用户**
   - 验证用户信息
   - 创建用户
   - 发送通知

2. **更新用户**
   - 更新用户信息
   - 发送通知

3. **删除用户**
   - 删除用户
   - 发送通知

### 活动（Activities）

- `ValidateUserActivity` - 验证用户信息
- `CreateUserActivity` - 创建用户
- `UpdateUserActivity` - 更新用户
- `DeleteUserActivity` - 删除用户
- `SendNotificationActivity` - 发送通知

## 配置

在 `configs/config.yaml` 中配置：

```yaml
workflow:
  temporal:
    address: "localhost:7233"
    task_queue: "user-task-queue"
    namespace: "default"
```

## 开发工作流

### 1. 定义工作流

在 `internal/application/workflow/` 中定义工作流函数：

```go
func MyWorkflow(ctx workflow.Context, input MyInput) (MyOutput, error) {
    // 工作流逻辑
    return output, nil
}
```

### 2. 定义活动

在 `internal/application/workflow/` 中定义活动函数：

```go
func MyActivity(ctx context.Context, input string) (string, error) {
    // 活动逻辑
    return result, nil
}
```

### 3. 注册工作流和活动

在 `cmd/temporal-worker/main.go` 中注册：

```go
w.RegisterWorkflow(MyWorkflow)
w.RegisterActivity(MyActivity)
```

## 最佳实践

1. **幂等性** - 确保活动和查询是幂等的
2. **超时设置** - 为活动设置合理的超时时间
3. **重试策略** - 配置适当的重试策略
4. **错误处理** - 正确处理和传播错误
5. **信号和查询** - 使用信号和查询进行工作流交互

## 监控

### Temporal UI

访问 http://localhost:8088 查看：
- 工作流执行历史
- 活动执行状态
- 工作流查询和信号

### 指标

Temporal 提供丰富的指标，可以集成到 Prometheus 和 Grafana。

## 故障排除

### Worker 无法连接

- 检查 Temporal Server 是否运行
- 验证连接地址配置
- 检查网络连接

### 工作流执行失败

- 查看 Temporal UI 中的错误信息
- 检查活动日志
- 验证输入参数

## 参考资源

- [Temporal 文档](https://docs.temporal.io/)
- [Temporal Go SDK](https://docs.temporal.io/dev-guide/go)
- [工作流模式](https://docs.temporal.io/workflows)
