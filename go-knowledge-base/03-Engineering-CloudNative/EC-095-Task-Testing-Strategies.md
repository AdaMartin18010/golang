# 任务测试策略 (Task Testing Strategies)

> **分类**: 工程与云原生
> **标签**: #testing #unit-test #integration-test #mock
> **参考**: Go Testing, Testify, Testing Patterns

---

## 测试策略架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Task Testing Strategy Pyramid                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                    E2E Tests (Top)                                   │   │
│   │   - Full workflow testing                                            │   │
│   │   - Integration with real dependencies                               │   │
│   │   - Slow, comprehensive                                              │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                     ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                    Integration Tests                                 │   │
│   │   - Component interaction testing                                    │   │
│   │   - Database, queue, cache integration                               │   │
│   │   - Medium speed                                                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                     ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                    Unit Tests (Base)                                 │   │
│   │   - Function-level testing                                           │   │
│   │   - Mocked dependencies                                              │   │
│   │   - Fast, isolated                                                   │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整测试实现

```go
package tasktest

import (
    "context"
    "errors"
    "testing"
    "time"
)

// Task 任务定义
type Task struct {
    ID       string
    Type     string
    Payload  interface{}
    Handler  TaskHandler
}

// TaskHandler 任务处理器
type TaskHandler func(ctx context.Context, payload interface{}) (interface{}, error)

// TaskExecutor 任务执行器
type TaskExecutor struct {
    handlers map[string]TaskHandler
}

// NewTaskExecutor 创建执行器
func NewTaskExecutor() *TaskExecutor {
    return &TaskExecutor{
        handlers: make(map[string]TaskHandler),
    }
}

// Register 注册处理器
func (te *TaskExecutor) Register(taskType string, handler TaskHandler) {
    te.handlers[taskType] = handler
}

// Execute 执行任务
func (te *TaskExecutor) Execute(ctx context.Context, task *Task) (interface{}, error) {
    handler, ok := te.handlers[task.Type]
    if !ok {
        return nil, errors.New("handler not found")
    }

    return handler(ctx, task.Payload)
}

// MockTaskQueue 模拟任务队列
type MockTaskQueue struct {
    tasks    []*Task
    dequeued []*Task
    mu       sync.Mutex
}

// NewMockTaskQueue 创建模拟队列
func NewMockTaskQueue() *MockTaskQueue {
    return &MockTaskQueue{
        tasks:    make([]*Task, 0),
        dequeued: make([]*Task, 0),
    }
}

// Enqueue 入队
func (mq *MockTaskQueue) Enqueue(task *Task) error {
    mq.mu.Lock()
    defer mq.mu.Unlock()
    mq.tasks = append(mq.tasks, task)
    return nil
}

// Dequeue 出队
func (mq *MockTaskQueue) Dequeue() (*Task, error) {
    mq.mu.Lock()
    defer mq.mu.Unlock()

    if len(mq.tasks) == 0 {
        return nil, nil
    }

    task := mq.tasks[0]
    mq.tasks = mq.tasks[1:]
    mq.dequeued = append(mq.dequeued, task)

    return task, nil
}

// GetDequeued 获取已出队任务
func (mq *MockTaskQueue) GetDequeued() []*Task {
    mq.mu.Lock()
    defer mq.mu.Unlock()
    return mq.dequeued
}

// Reset 重置
func (mq *MockTaskQueue) Reset() {
    mq.mu.Lock()
    defer mq.mu.Unlock()
    mq.tasks = make([]*Task, 0)
    mq.dequeued = make([]*Task, 0)
}

// MockTaskStore 模拟任务存储
type MockTaskStore struct {
    tasks map[string]*Task
    mu    sync.RWMutex
}

// NewMockTaskStore 创建模拟存储
func NewMockTaskStore() *MockTaskStore {
    return &MockTaskStore{
        tasks: make(map[string]*Task),
    }
}

// Save 保存任务
func (ms *MockTaskStore) Save(task *Task) error {
    ms.mu.Lock()
    defer ms.mu.Unlock()
    ms.tasks[task.ID] = task
    return nil
}

// Get 获取任务
func (ms *MockTaskStore) Get(id string) (*Task, error) {
    ms.mu.RLock()
    defer ms.mu.RUnlock()

    task, ok := ms.tasks[id]
    if !ok {
        return nil, errors.New("task not found")
    }

    return task, nil
}

// TestTaskExecutor_Execute_Success 测试成功执行
func TestTaskExecutor_Execute_Success(t *testing.T) {
    // Arrange
    executor := NewTaskExecutor()
    executor.Register("test-task", func(ctx context.Context, payload interface{}) (interface{}, error) {
        return "success", nil
    })

    task := &Task{
        ID:      "task-1",
        Type:    "test-task",
        Payload: map[string]string{"key": "value"},
    }

    // Act
    result, err := executor.Execute(context.Background(), task)

    // Assert
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if result != "success" {
        t.Errorf("expected 'success', got %v", result)
    }
}

// TestTaskExecutor_Execute_HandlerNotFound 测试处理器不存在
func TestTaskExecutor_Execute_HandlerNotFound(t *testing.T) {
    // Arrange
    executor := NewTaskExecutor()

    task := &Task{
        ID:   "task-1",
        Type: "unknown-task",
    }

    // Act
    _, err := executor.Execute(context.Background(), task)

    // Assert
    if err == nil {
        t.Error("expected error, got nil")
    }
}

// TestTaskExecutor_Execute_Timeout 测试超时
func TestTaskExecutor_Execute_Timeout(t *testing.T) {
    // Arrange
    executor := NewTaskExecutor()
    executor.Register("slow-task", func(ctx context.Context, payload interface{}) (interface{}, error) {
        select {
        case <-time.After(5 * time.Second):
            return "done", nil
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    })

    task := &Task{
        ID:   "task-1",
        Type: "slow-task",
    }

    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    // Act
    _, err := executor.Execute(ctx, task)

    // Assert
    if err != context.DeadlineExceeded {
        t.Errorf("expected DeadlineExceeded, got %v", err)
    }
}

// TestMockTaskQueue_EnqueueDequeue 测试队列操作
func TestMockTaskQueue_EnqueueDequeue(t *testing.T) {
    // Arrange
    queue := NewMockTaskQueue()
    task := &Task{ID: "task-1", Type: "test"}

    // Act
    err := queue.Enqueue(task)
    if err != nil {
        t.Fatalf("enqueue failed: %v", err)
    }

    dequeued, err := queue.Dequeue()

    // Assert
    if err != nil {
        t.Errorf("dequeue failed: %v", err)
    }
    if dequeued.ID != task.ID {
        t.Errorf("expected task %s, got %s", task.ID, dequeued.ID)
    }
}

// IntegrationTest 集成测试
type IntegrationTest struct {
    executor *TaskExecutor
    queue    *MockTaskQueue
    store    *MockTaskStore
}

// SetupIntegrationTest 设置集成测试
func SetupIntegrationTest() *IntegrationTest {
    return &IntegrationTest{
        executor: NewTaskExecutor(),
        queue:    NewMockTaskQueue(),
        store:    NewMockTaskStore(),
    }
}

// TestIntegration_FullWorkflow 测试完整工作流
func TestIntegration_FullWorkflow(t *testing.T) {
    // Setup
    test := SetupIntegrationTest()

    test.executor.Register("process-order", func(ctx context.Context, payload interface{}) (interface{}, error) {
        data := payload.(map[string]string)
        return map[string]string{
            "order_id": data["order_id"],
            "status":   "processed",
        }, nil
    })

    // Execute workflow
    task := &Task{
        ID:      "order-123",
        Type:    "process-order",
        Payload: map[string]string{"order_id": "ORD-456"},
    }

    // Save task
    if err := test.store.Save(task); err != nil {
        t.Fatalf("save failed: %v", err)
    }

    // Enqueue
    if err := test.queue.Enqueue(task); err != nil {
        t.Fatalf("enqueue failed: %v", err)
    }

    // Dequeue and execute
    dequeued, _ := test.queue.Dequeue()
    result, err := test.executor.Execute(context.Background(), dequeued)

    // Assert
    if err != nil {
        t.Errorf("execution failed: %v", err)
    }

    resultMap := result.(map[string]string)
    if resultMap["status"] != "processed" {
        t.Errorf("expected status 'processed', got %s", resultMap["status"])
    }
}

// BenchmarkTaskExecutor_Execute 基准测试
func BenchmarkTaskExecutor_Execute(b *testing.B) {
    executor := NewTaskExecutor()
    executor.Register("benchmark-task", func(ctx context.Context, payload interface{}) (interface{}, error) {
        return payload, nil
    })

    task := &Task{
        ID:      "task-1",
        Type:    "benchmark-task",
        Payload: "test-payload",
    }

    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = executor.Execute(ctx, task)
    }
}

// ParallelBenchmark 并行基准测试
func BenchmarkTaskExecutor_Execute_Parallel(b *testing.B) {
    executor := NewTaskExecutor()
    executor.Register("parallel-task", func(ctx context.Context, payload interface{}) (interface{}, error) {
        time.Sleep(1 * time.Millisecond)
        return payload, nil
    })

    task := &Task{
        Type:    "parallel-task",
        Payload: "test",
    }

    ctx := context.Background()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, _ = executor.Execute(ctx, task)
        }
    })
}

// TableDrivenTest 表驱动测试
func TestTaskExecutor_Execute_TableDriven(t *testing.T) {
    executor := NewTaskExecutor()
    executor.Register("add", func(ctx context.Context, payload interface{}) (interface{}, error) {
        nums := payload.([]int)
        return nums[0] + nums[1], nil
    })
    executor.Register("error", func(ctx context.Context, payload interface{}) (interface{}, error) {
        return nil, errors.New("intentional error")
    })

    tests := []struct {
        name     string
        taskType string
        payload  interface{}
        want     interface{}
        wantErr  bool
    }{
        {"add two numbers", "add", []int{1, 2}, 3, false},
        {"add negatives", "add", []int{-1, -2}, -3, false},
        {"handler error", "error", nil, nil, true},
        {"unknown handler", "unknown", nil, nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            task := &Task{
                Type:    tt.taskType,
                Payload: tt.payload,
            }

            got, err := executor.Execute(context.Background(), task)

            if (err != nil) != tt.wantErr {
                t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Execute() = %v, want %v", got, tt.want)
            }
        })
    }
}

import "sync"
