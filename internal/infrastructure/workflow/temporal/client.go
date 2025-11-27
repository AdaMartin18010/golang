// Package temporal provides Temporal workflow orchestration client and worker implementations.
//
// Temporal 是一个开源的工作流编排引擎，用于构建可靠的分布式应用。
//
// 设计原则：
// 1. 可靠性：工作流状态持久化，支持故障恢复
// 2. 可扩展性：支持大规模并发工作流执行
// 3. 可观测性：提供完整的工作流执行历史和状态查询
// 4. 灵活性：支持复杂的工作流模式（顺序、并行、条件、循环等）
//
// 核心组件：
// - Client: 用于启动、查询、信号化工作流
// - Worker: 用于执行工作流和活动（Activities）
// - Workflow: 定义业务流程逻辑
// - Activity: 执行具体的业务操作
//
// 使用场景：
// - 长时间运行的业务流程（订单处理、数据迁移等）
// - 需要可靠执行的关键业务逻辑
// - 复杂的多步骤业务流程编排
// - 需要支持重试、超时、取消的业务操作
//
// 架构说明：
// - Temporal Server: 管理工作流状态、调度任务
// - Worker: 执行工作流和活动
// - Client: 与 Temporal Server 通信，启动和查询工作流
package temporal

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
)

// Client 是 Temporal 客户端的封装，用于与 Temporal Server 通信。
//
// 功能说明：
// - 启动工作流（ExecuteWorkflow）
// - 查询工作流状态（GetWorkflow、QueryWorkflow）
// - 向工作流发送信号（SignalWorkflow）
// - 管理工作流生命周期
//
// 设计说明：
// - 封装了 Temporal SDK 的 client.Client
// - 提供更简洁的 API 接口
// - 支持依赖注入和测试
//
// 使用示例：
//
//	// 创建客户端
//	temporalClient, err := temporal.NewClient("localhost:7233")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer temporalClient.Close()
//
//	// 启动工作流
//	workflowOptions := client.StartWorkflowOptions{
//	    ID:        "workflow-id",
//	    TaskQueue: "my-task-queue",
//	}
//	workflowRun, err := temporalClient.ExecuteWorkflow(ctx, workflowOptions, MyWorkflow, input)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 查询工作流结果
//	var result MyResult
//	err = workflowRun.Get(ctx, &result)
//
// 注意事项：
// - 客户端是线程安全的，可以在多个 goroutine 中使用
// - 应在应用程序生命周期中复用客户端实例
// - 退出前应调用 Close() 关闭客户端连接
type Client struct {
	client client.Client
}

// NewClient 创建并连接到 Temporal Server 的客户端。
//
// 功能说明：
// - 建立与 Temporal Server 的 gRPC 连接
// - 配置客户端选项（地址、命名空间等）
// - 返回配置好的客户端实例
//
// 参数：
// - address: Temporal Server 的地址和端口
//   格式：host:port（例如：localhost:7233）
//   默认端口：7233（gRPC）
//
// 返回：
// - *Client: 配置好的客户端实例
// - error: 如果连接失败，返回错误信息
//
// 配置选项：
// - HostPort: Temporal Server 地址
// - Namespace: 命名空间（默认：default）
// - ConnectionOptions: 连接选项（TLS、重试等）
//
// 使用示例：
//
//	// 基本用法
//	client, err := temporal.NewClient("localhost:7233")
//	if err != nil {
//	    log.Fatal("Failed to create Temporal client:", err)
//	}
//	defer client.Close()
//
//	// 高级配置（如果需要）
//	// 可以通过修改 client.Options 添加更多配置：
//	// - Namespace: 指定命名空间
//	// - ConnectionOptions: 配置 TLS、重试策略等
//
// 注意事项：
// - 确保 Temporal Server 已启动并监听指定地址
// - 客户端连接是长连接，应复用客户端实例
// - 生产环境建议配置 TLS 和命名空间隔离
func NewClient(address string) (*Client, error) {
	c, err := client.Dial(client.Options{
		HostPort: address,
		// 其他可选配置：
		// Namespace: "production",
		// ConnectionOptions: client.ConnectionOptions{
		//     TLS: &tls.Config{...},
		// },
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client: %w", err)
	}

	return &Client{client: c}, nil
}

// ExecuteWorkflow 启动并执行一个工作流。
//
// 功能说明：
// - 向 Temporal Server 发送工作流启动请求
// - 返回 WorkflowRun 用于查询工作流状态和结果
// - 工作流会异步执行，不阻塞调用
//
// 参数：
// - ctx: 上下文，用于控制请求超时和取消
// - options: 工作流启动选项
//   - ID: 工作流 ID（唯一标识符）
//   - TaskQueue: 任务队列名称（Worker 监听此队列）
//   - WorkflowExecutionTimeout: 工作流执行超时时间
//   - WorkflowTaskTimeout: 工作流任务超时时间
//   - WorkflowIDReusePolicy: 工作流 ID 重用策略
// - workflow: 工作流函数（必须是已注册的工作流）
// - args: 传递给工作流的参数
//
// 返回：
// - client.WorkflowRun: 工作流运行实例，用于查询状态和结果
// - error: 如果启动失败，返回错误信息
//
// 使用示例：
//
//	workflowOptions := client.StartWorkflowOptions{
//	    ID:        "user-creation-workflow-123",
//	    TaskQueue: "user-task-queue",
//	    WorkflowExecutionTimeout: time.Hour,
//	}
//	workflowRun, err := client.ExecuteWorkflow(ctx, workflowOptions, UserCreationWorkflow, userInput)
//	if err != nil {
//	    log.Fatal("Failed to start workflow:", err)
//	}
//
//	// 等待工作流完成并获取结果
//	var result UserCreationResult
//	err = workflowRun.Get(ctx, &result)
//
// 注意事项：
// - 工作流 ID 应具有唯一性，避免冲突
// - TaskQueue 必须与 Worker 监听的队列名称一致
// - 工作流函数必须是确定性的（不能使用随机数、时间等）
func (c *Client) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return c.client.ExecuteWorkflow(ctx, options, workflow, args...)
}

// GetWorkflow 根据工作流 ID 和运行 ID 获取工作流运行实例。
//
// 功能说明：
// - 获取已存在的工作流运行实例
// - 用于查询工作流状态、结果或发送信号
// - 如果工作流不存在，返回的 WorkflowRun 在调用时会返回错误
//
// 参数：
// - ctx: 上下文
// - workflowID: 工作流 ID（唯一标识符）
// - runID: 工作流运行 ID（可选，如果为空则获取最新运行）
//
// 返回：
// - client.WorkflowRun: 工作流运行实例
//
// 使用示例：
//
//	// 获取特定运行的工作流
//	workflowRun := client.GetWorkflow(ctx, "workflow-id", "run-id")
//	var result MyResult
//	err := workflowRun.Get(ctx, &result)
//
//	// 获取最新运行的工作流（runID 为空）
//	workflowRun := client.GetWorkflow(ctx, "workflow-id", "")
//
// 注意事项：
// - 如果工作流不存在，WorkflowRun 的方法调用会返回错误
// - runID 为空时，获取该工作流 ID 的最新运行实例
func (c *Client) GetWorkflow(ctx context.Context, workflowID, runID string) client.WorkflowRun {
	return c.client.GetWorkflow(ctx, workflowID, runID)
}

// SignalWorkflow 向工作流发送信号（Signal）。
//
// 功能说明：
// - 信号是工作流外部与工作流通信的机制
// - 工作流可以通过 workflow.GetSignalChannel() 接收信号
// - 信号是异步的，不会阻塞工作流执行
//
// 参数：
// - ctx: 上下文，用于控制请求超时
// - workflowID: 工作流 ID
// - runID: 工作流运行 ID（可选，如果为空则发送给最新运行）
// - signalName: 信号名称（工作流中定义的信号名称）
// - arg: 信号参数（可以是任意类型）
//
// 返回：
// - error: 如果发送失败，返回错误信息
//
// 使用场景：
// - 取消工作流执行
// - 更新工作流配置
// - 通知工作流外部事件
// - 用户交互（如审批、确认等）
//
// 使用示例：
//
//	// 发送取消信号
//	err := client.SignalWorkflow(ctx, "workflow-id", "", "cancel", nil)
//
//	// 发送更新信号
//	updateData := UpdateData{Status: "approved"}
//	err := client.SignalWorkflow(ctx, "workflow-id", "", "update", updateData)
//
// 注意事项：
// - 信号是异步的，工作流可能不会立即处理
// - 如果工作流不存在或已结束，信号会被忽略
// - 信号参数会被序列化，确保类型可序列化
func (c *Client) SignalWorkflow(ctx context.Context, workflowID, runID, signalName string, arg interface{}) error {
	return c.client.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}

// QueryWorkflow 查询工作流状态（Query）。
//
// 功能说明：
// - 查询是同步的，立即返回工作流的当前状态
// - 工作流必须实现查询处理器（Query Handler）
// - 查询不会影响工作流的执行状态
//
// 参数：
// - ctx: 上下文，用于控制请求超时
// - workflowID: 工作流 ID
// - runID: 工作流运行 ID（可选）
// - queryType: 查询类型（工作流中定义的查询名称）
// - args: 查询参数（可选）
//
// 返回：
// - interface{}: 查询结果（需要类型断言）
// - error: 如果查询失败，返回错误信息
//
// 使用场景：
// - 查询工作流当前状态
// - 获取工作流中间结果
// - 检查工作流进度
// - 获取工作流配置信息
//
// 使用示例：
//
//	// 查询工作流状态
//	result, err := client.QueryWorkflow(ctx, "workflow-id", "", "getStatus", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	status := result.(string)
//
//	// 查询工作流进度
//	result, err := client.QueryWorkflow(ctx, "workflow-id", "", "getProgress", nil)
//	progress := result.(int)
//
// 注意事项：
// - 查询是同步的，会等待工作流处理查询
// - 工作流必须实现对应的查询处理器
// - 查询结果需要类型断言才能使用
func (c *Client) QueryWorkflow(ctx context.Context, workflowID, runID, queryType string, args ...interface{}) (interface{}, error) {
	return c.client.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
}

// Close 关闭客户端连接。
//
// 功能说明：
// - 关闭与 Temporal Server 的 gRPC 连接
// - 释放客户端资源
// - 停止所有后台任务
//
// 使用示例：
//
//	defer client.Close()
//
// 注意事项：
// - 应在应用程序退出前调用
// - 关闭后不应再使用该客户端
// - 关闭是同步的，会等待所有待处理的请求完成
func (c *Client) Close() {
	c.client.Close()
}

// Client 返回底层的 Temporal SDK 客户端。
//
// 功能说明：
// - 提供对底层 client.Client 的访问
// - 用于需要直接使用 SDK 功能的场景（如创建 Worker）
//
// 返回：
// - client.Client: 底层的 Temporal SDK 客户端
//
// 使用场景：
// - 创建 Worker 时需要底层客户端
// - 需要使用 SDK 的高级功能
// - 与第三方库集成
//
// 使用示例：
//
//	// 使用底层客户端创建 Worker
//	temporalClient := client.Client()
//	worker := worker.New(temporalClient, "task-queue", worker.Options{})
//
// 注意事项：
// - 返回的客户端与封装客户端共享连接
// - 不应直接关闭返回的客户端（应使用 Close() 方法）
func (c *Client) Client() client.Client {
	return c.client
}
