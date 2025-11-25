package temporal

import (
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/client"
)

// Worker Temporal Worker
type Worker struct {
	worker worker.Worker
}

// NewWorker 创建 Worker
func NewWorker(c client.Client, taskQueue string) *Worker {
	w := worker.New(c, taskQueue, worker.Options{})
	return &Worker{worker: w}
}

// NewWorkerFromClient 从 Client 包装器创建 Worker
func NewWorkerFromClient(c *Client, taskQueue string) *Worker {
	return NewWorker(c.Client(), taskQueue)
}

// RegisterWorkflow 注册工作流
func (w *Worker) RegisterWorkflow(workflow interface{}) {
	w.worker.RegisterWorkflow(workflow)
}

// RegisterActivity 注册活动
func (w *Worker) RegisterActivity(activity interface{}) {
	w.worker.RegisterActivity(activity)
}

// Start 启动 Worker
func (w *Worker) Start() error {
	return w.worker.Run(worker.InterruptCh())
}

// Stop 停止 Worker
func (w *Worker) Stop() {
	w.worker.Stop()
}

// Run 运行 Worker（阻塞）
func (w *Worker) Run() error {
	return w.worker.Run(worker.InterruptCh())
}
