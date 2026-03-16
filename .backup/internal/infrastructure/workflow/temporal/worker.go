package temporal

import (
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/client"
)

// Worker 是 Temporal Worker 的封装，用于执行工作流和活动（Activities）。
//
// 功能说明：
// - 从 Temporal Server 接收任务（工作流任务和活动任务）
// - 执行已注册的工作流和活动
// - 将执行结果返回给 Temporal Server
//
// 设计说明：
// - 封装了 Temporal SDK 的 worker.Worker
// - 提供更简洁的 API 接口
// - 支持工作流和活动的注册和管理
//
// 工作流程：
// 1. 创建 Worker 并指定任务队列（Task Queue）
// 2. 注册工作流和活动函数
// 3. 启动 Worker，开始监听任务
// 4. Worker 从 Temporal Server 接收任务并执行
// 5. 将执行结果返回给 Temporal Server
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
//	// 创建 Worker
//	w := temporal.NewWorkerFromClient(temporalClient, "my-task-queue")
//
//	// 注册工作流
//	w.RegisterWorkflow(MyWorkflow)
//
//	// 注册活动
//	w.RegisterActivity(MyActivity)
//
//	// 启动 Worker（阻塞）
//	if err := w.Run(); err != nil {
//	    log.Fatal(err)
//	}
//
// 注意事项：
// - Worker 必须监听与客户端启动工作流时指定的相同任务队列
// - 工作流函数必须是确定性的（不能使用随机数、时间等）
// - 活动函数可以执行非确定性操作（网络请求、数据库操作等）
// - Worker 应该长时间运行，不要频繁启动和停止
type Worker struct {
	worker worker.Worker
}

// NewWorker 创建一个新的 Temporal Worker。
//
// 功能说明：
// - 使用底层 Temporal SDK 客户端创建 Worker
// - 指定任务队列名称
// - 配置 Worker 选项（并发数、活动重试策略等）
//
// 参数：
// - c: Temporal SDK 客户端实例
// - taskQueue: 任务队列名称
//   Worker 会监听此队列的任务
//   必须与客户端启动工作流时指定的队列名称一致
//
// 返回：
// - *Worker: 配置好的 Worker 实例
//
// 使用示例：
//
//	// 使用底层客户端创建 Worker
//	temporalClient, _ := client.Dial(client.Options{HostPort: "localhost:7233"})
//	w := temporal.NewWorker(temporalClient, "my-task-queue")
//
// 注意事项：
// - 任务队列名称应具有业务意义，便于管理
// - 多个 Worker 可以监听同一个任务队列（负载均衡）
// - Worker 选项可以通过 worker.Options 配置
func NewWorker(c client.Client, taskQueue string) *Worker {
	w := worker.New(c, taskQueue, worker.Options{
		// 可选配置：
		// MaxConcurrentActivityExecutionSize: 10,  // 最大并发活动数
		// MaxConcurrentWorkflowTaskExecutionSize: 10, // 最大并发工作流任务数
		// MaxConcurrentLocalActivityExecutionSize: 10, // 最大并发本地活动数
		// ActivityHeartbeatTimeout: time.Second * 30, // 活动心跳超时
		// WorkflowTaskHeartbeatTimeout: time.Second * 30, // 工作流任务心跳超时
	})
	return &Worker{worker: w}
}

// NewWorkerFromClient 从封装的 Client 创建 Worker。
//
// 功能说明：
// - 使用封装的 Client 创建 Worker
// - 更便捷的创建方式，无需直接访问底层客户端
//
// 参数：
// - c: 封装的 Temporal Client 实例
// - taskQueue: 任务队列名称
//
// 返回：
// - *Worker: 配置好的 Worker 实例
//
// 使用示例：
//
//	// 创建客户端
//	temporalClient, _ := temporal.NewClient("localhost:7233")
//	defer temporalClient.Close()
//
//	// 从客户端创建 Worker
//	w := temporal.NewWorkerFromClient(temporalClient, "my-task-queue")
//
// 注意事项：
// - 这是推荐的创建 Worker 的方式
// - Client 和 Worker 可以共享同一个连接
func NewWorkerFromClient(c *Client, taskQueue string) *Worker {
	return NewWorker(c.Client(), taskQueue)
}

// RegisterWorkflow 注册一个工作流函数。
//
// 功能说明：
// - 将工作流函数注册到 Worker
// - Worker 可以执行已注册的工作流
// - 工作流函数必须是已定义的函数
//
// 参数：
// - workflow: 工作流函数
//   工作流函数必须满足以下要求：
//   1. 函数签名：func(ctx workflow.Context, input InputType) (OutputType, error)
//   2. 函数必须是确定性的（不能使用随机数、时间、网络请求等）
//   3. 只能使用 workflow 包提供的 API
//
// 使用示例：
//
//	// 定义工作流函数
//	func UserCreationWorkflow(ctx workflow.Context, input UserInput) (UserOutput, error) {
//	    // 工作流逻辑
//	    return output, nil
//	}
//
//	// 注册工作流
//	w.RegisterWorkflow(UserCreationWorkflow)
//
// 注意事项：
// - 工作流函数必须是确定性的
// - 不能在工作流中直接调用外部服务（应使用 Activity）
// - 工作流函数会在 Worker 重启后继续执行
// - 同一个工作流函数可以注册多次（使用不同的名称）
func (w *Worker) RegisterWorkflow(workflow interface{}) {
	w.worker.RegisterWorkflow(workflow)
}

// RegisterActivity 注册一个活动（Activity）函数。
//
// 功能说明：
// - 将活动函数注册到 Worker
// - Worker 可以执行已注册的活动
// - 活动函数可以执行非确定性操作
//
// 参数：
// - activity: 活动函数
//   活动函数必须满足以下要求：
//   1. 函数签名：func(ctx context.Context, input InputType) (OutputType, error)
//   2. 可以使用标准库和第三方库
//   3. 可以执行网络请求、数据库操作等
//
// 使用示例：
//
//	// 定义活动函数
//	func CreateUserActivity(ctx context.Context, input UserInput) (UserOutput, error) {
//	    // 活动逻辑（可以执行数据库操作、网络请求等）
//	    user, err := userService.CreateUser(input)
//	    return UserOutput{User: user}, err
//	}
//
//	// 注册活动
//	w.RegisterActivity(CreateUserActivity)
//
// 注意事项：
// - 活动函数可以执行非确定性操作
// - 活动支持自动重试（可配置重试策略）
// - 活动支持超时控制（可配置超时时间）
// - 活动支持心跳机制（用于长时间运行的活动）
func (w *Worker) RegisterActivity(activity interface{}) {
	w.worker.RegisterActivity(activity)
}

// Start 启动 Worker 并开始处理任务。
//
// 功能说明：
// - 启动 Worker，开始从 Temporal Server 接收任务
// - 执行已注册的工作流和活动
// - 方法会阻塞，直到 Worker 停止
//
// 返回：
// - error: 如果启动失败，返回错误信息
//
// 使用示例：
//
//	// 启动 Worker（阻塞）
//	if err := w.Start(); err != nil {
//	    log.Fatal("Worker failed:", err)
//	}
//
// 注意事项：
// - 方法会阻塞当前 goroutine
// - 应在单独的 goroutine 中运行，或作为主程序的主循环
// - 使用 Stop() 方法可以停止 Worker
// - 停止后可以再次启动
func (w *Worker) Start() error {
	return w.worker.Run(worker.InterruptCh())
}

// Stop 停止 Worker。
//
// 功能说明：
// - 停止 Worker，不再接收新任务
// - 等待当前正在执行的任务完成
// - 释放 Worker 资源
//
// 使用示例：
//
//	// 在 goroutine 中启动 Worker
//	go func() {
//	    if err := w.Start(); err != nil {
//	        log.Fatal(err)
//	    }
//	}()
//
//	// 在需要时停止 Worker
//	w.Stop()
//
// 注意事项：
// - 停止是优雅的，会等待当前任务完成
// - 停止后可以再次启动
// - 应在应用程序退出前调用
func (w *Worker) Stop() {
	w.worker.Stop()
}

// Run 运行 Worker（阻塞）。
//
// 功能说明：
// - 启动 Worker 并开始处理任务
// - 方法会阻塞，直到收到中断信号或调用 Stop()
// - 与 Start() 方法功能相同，提供更语义化的命名
//
// 返回：
// - error: 如果运行失败，返回错误信息
//
// 使用示例：
//
//	// 运行 Worker（阻塞）
//	if err := w.Run(); err != nil {
//	    log.Fatal("Worker failed:", err)
//	}
//
// 注意事项：
// - 方法会阻塞当前 goroutine
// - 应在单独的 goroutine 中运行，或作为主程序的主循环
// - 使用 Stop() 方法可以停止 Worker
// - 与 Start() 方法功能相同，选择使用哪个取决于代码语义
func (w *Worker) Run() error {
	return w.worker.Run(worker.InterruptCh())
}
