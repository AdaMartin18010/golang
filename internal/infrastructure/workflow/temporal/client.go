package temporal

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
)

// Client Temporal 客户端
type Client struct {
	client client.Client
}

// NewClient 创建 Temporal 客户端
func NewClient(address string) (*Client, error) {
	c, err := client.Dial(client.Options{
		HostPort: address,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client: %w", err)
	}

	return &Client{client: c}, nil
}

// ExecuteWorkflow 执行工作流
func (c *Client) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return c.client.ExecuteWorkflow(ctx, options, workflow, args...)
}

// GetWorkflow 获取工作流
func (c *Client) GetWorkflow(ctx context.Context, workflowID, runID string) client.WorkflowRun {
	return c.client.GetWorkflow(ctx, workflowID, runID)
}

// SignalWorkflow 发送信号给工作流
func (c *Client) SignalWorkflow(ctx context.Context, workflowID, runID, signalName string, arg interface{}) error {
	return c.client.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}

// QueryWorkflow 查询工作流
func (c *Client) QueryWorkflow(ctx context.Context, workflowID, runID, queryType string, args ...interface{}) (interface{}, error) {
	return c.client.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
}

// Close 关闭客户端
func (c *Client) Close() {
	c.client.Close()
}

// Client 返回底层客户端（用于 Worker）
func (c *Client) Client() client.Client {
	return c.client
}
