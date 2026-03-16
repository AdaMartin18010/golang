package patterns

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockCommand 测试用的 Mock Command
type MockCommand struct {
	shouldError bool
	errMsg      string
}

func (m MockCommand) Execute(ctx context.Context) error {
	if m.shouldError {
		return errors.New(m.errMsg)
	}
	return nil
}

// MockCommandHandler 测试用的 Mock Command Handler
type MockCommandHandler struct {
	handleCalled bool
	lastCommand  MockCommand
	returnError  error
}

func (m *MockCommandHandler) Handle(ctx context.Context, cmd MockCommand) error {
	m.handleCalled = true
	m.lastCommand = cmd
	return m.returnError
}

// ==================== CommandResult 测试 ====================

// TestCommandResult 测试 CommandResult 结构
func TestCommandResult(t *testing.T) {
	tests := []struct {
		name     string
		result   CommandResult
		expected struct {
			success bool
			message string
			data    interface{}
		}
	}{
		{
			name: "成功结果",
			result: CommandResult{
				Success: true,
				Message: "操作成功",
				Data:    "created-id-123",
			},
			expected: struct {
				success bool
				message string
				data    interface{}
			}{
				success: true,
				message: "操作成功",
				data:    "created-id-123",
			},
		},
		{
			name: "失败结果",
			result: CommandResult{
				Success: false,
				Message: "操作失败",
				Data:    nil,
			},
			expected: struct {
				success bool
				message string
				data    interface{}
			}{
				success: false,
				message: "操作失败",
				data:    nil,
			},
		},
		{
			name: "复杂数据",
			result: CommandResult{
				Success: true,
				Message: "批量操作完成",
				Data: map[string]interface{}{
					"created": 5,
					"updated": 3,
					"failed":  0,
				},
			},
			expected: struct {
				success bool
				message string
				data    interface{}
			}{
				success: true,
				message: "批量操作完成",
				data: map[string]interface{}{
					"created": 5,
					"updated": 3,
					"failed":  0,
				},
			},
		},
		{
			name: "空消息",
			result: CommandResult{
				Success: true,
				Message: "",
				Data:    nil,
			},
			expected: struct {
				success bool
				message string
				data    interface{}
			}{
				success: true,
				message: "",
				data:    nil,
			},
		},
		{
			name: "结构体数据",
			result: CommandResult{
				Success: true,
				Message: "用户创建成功",
				Data: struct {
					ID    string `json:"id"`
					Email string `json:"email"`
				}{
					ID:    "user-123",
					Email: "test@example.com",
				},
			},
			expected: struct {
				success bool
				message string
				data    interface{}
			}{
				success: true,
				message: "用户创建成功",
				data: struct {
					ID    string `json:"id"`
					Email string `json:"email"`
				}{
					ID:    "user-123",
					Email: "test@example.com",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.success, tt.result.Success)
			assert.Equal(t, tt.expected.message, tt.result.Message)
			assert.Equal(t, tt.expected.data, tt.result.Data)
		})
	}
}

// TestCommandResultWithNilData 测试 CommandResult 的 nil 数据
func TestCommandResultWithNilData(t *testing.T) {
	result := CommandResult{
		Success: true,
		Message: "操作成功",
		Data:    nil,
	}

	assert.True(t, result.Success)
	assert.Equal(t, "操作成功", result.Message)
	assert.Nil(t, result.Data)
}

// TestCommandResultWithSliceData 测试 CommandResult 的切片数据
func TestCommandResultWithSliceData(t *testing.T) {
	result := CommandResult{
		Success: true,
		Message: "列表获取成功",
		Data:    []string{"item1", "item2", "item3"},
	}

	assert.True(t, result.Success)
	data, ok := result.Data.([]string)
	require.True(t, ok)
	assert.Len(t, data, 3)
}

// ==================== MockCommand 测试 ====================

// TestMockCommand_Execute 测试 Mock Command 执行
func TestMockCommand_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("执行成功", func(t *testing.T) {
		cmd := MockCommand{shouldError: false}
		err := cmd.Execute(ctx)
		require.NoError(t, err)
	})

	t.Run("执行失败", func(t *testing.T) {
		cmd := MockCommand{shouldError: true, errMsg: "执行错误"}
		err := cmd.Execute(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "执行错误")
	})

	t.Run("空错误消息", func(t *testing.T) {
		cmd := MockCommand{shouldError: true, errMsg: ""}
		err := cmd.Execute(ctx)
		require.Error(t, err)
	})
}

// TestMockCommand_ExecuteWithCancelledContext 测试带取消上下文的命令执行
func TestMockCommand_ExecuteWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	cmd := MockCommand{shouldError: false}
	err := cmd.Execute(ctx)
	// MockCommand 没有检查上下文，所以应该成功
	require.NoError(t, err)
}

// ==================== MockCommandHandler 测试 ====================

// TestMockCommandHandler_Handle 测试 Command Handler
func TestMockCommandHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("处理成功", func(t *testing.T) {
		handler := &MockCommandHandler{returnError: nil}
		cmd := MockCommand{shouldError: false}

		err := handler.Handle(ctx, cmd)

		require.NoError(t, err)
		assert.True(t, handler.handleCalled)
		assert.Equal(t, cmd, handler.lastCommand)
	})

	t.Run("处理失败", func(t *testing.T) {
		handler := &MockCommandHandler{returnError: errors.New("处理失败")}
		cmd := MockCommand{shouldError: false}

		err := handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.True(t, handler.handleCalled)
		assert.Contains(t, err.Error(), "处理失败")
	})

	t.Run("多次处理", func(t *testing.T) {
		handler := &MockCommandHandler{returnError: nil}
		
		for i := 0; i < 5; i++ {
			cmd := MockCommand{shouldError: false}
			err := handler.Handle(ctx, cmd)
			require.NoError(t, err)
		}

		assert.True(t, handler.handleCalled)
	})
}

// TestCommandHandlerWithContext 测试带上下文的 Command Handler
func TestCommandHandlerWithContext(t *testing.T) {
	t.Run("上下文取消", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		handler := &MockCommandHandler{returnError: nil}
		cmd := MockCommand{shouldError: false}

		err := handler.Handle(ctx, cmd)
		// 这里 handler 没有检查上下文，所以不会返回错误
		// 实际实现应该检查 ctx.Err()
		require.NoError(t, err)
	})

	t.Run("上下文超时", func(t *testing.T) {
		// 创建一个超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		handler := &MockCommandHandler{returnError: nil}
		cmd := MockCommand{shouldError: false}

		err := handler.Handle(ctx, cmd)
		// 同上，实际实现应该检查 ctx.Err()
		require.NoError(t, err)
	})

	t.Run("上下文值传递", func(t *testing.T) {
		type contextKey string
		key := contextKey("test-key")
		ctx := context.WithValue(context.Background(), key, "test-value")

		handler := &MockCommandHandler{returnError: nil}
		cmd := MockCommand{shouldError: false}

		err := handler.Handle(ctx, cmd)
		require.NoError(t, err)
		assert.True(t, handler.handleCalled)
	})
}

// ==================== 泛型 CommandHandler 接口测试 ====================

// TestCommandHandlerInterface 测试 CommandHandler 接口
func TestCommandHandlerInterface(t *testing.T) {
	t.Run("接口实现", func(t *testing.T) {
		// 验证 MockCommandHandler 实现了 CommandHandler 接口
		var _ CommandHandler[MockCommand] = &MockCommandHandler{}
	})

	t.Run("泛型处理", func(t *testing.T) {
		handler := &MockCommandHandler{returnError: nil}
		ctx := context.Background()
		cmd := MockCommand{shouldError: false}

		var h CommandHandler[MockCommand] = handler
		err := h.Handle(ctx, cmd)

		require.NoError(t, err)
	})
}

// ==================== 复杂场景测试 ====================

// TestCommandChaining 测试命令链式处理
func TestCommandChaining(t *testing.T) {
	ctx := context.Background()

	// 第一个命令处理器
	handler1 := &MockCommandHandler{returnError: nil}
	cmd1 := MockCommand{shouldError: false}

	err := handler1.Handle(ctx, cmd1)
	require.NoError(t, err)

	// 第二个命令处理器
	handler2 := &MockCommandHandler{returnError: nil}
	cmd2 := MockCommand{shouldError: false}

	err = handler2.Handle(ctx, cmd2)
	require.NoError(t, err)

	// 验证两个处理器都被调用
	assert.True(t, handler1.handleCalled)
	assert.True(t, handler2.handleCalled)
}

// TestCommandWithComplexData 测试带复杂数据的命令
func TestCommandWithComplexData(t *testing.T) {
	type CreateOrderCommand struct {
		UserID      string
		ProductIDs  []string
		TotalAmount float64
		Address     struct {
			Street  string
			City    string
			Country string
		}
	}

	cmd := CreateOrderCommand{
		UserID:      "user-123",
		ProductIDs:  []string{"prod-1", "prod-2"},
		TotalAmount: 199.99,
		Address: struct {
			Street  string
			City    string
			Country string
		}{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
	}

	assert.Equal(t, "user-123", cmd.UserID)
	assert.Len(t, cmd.ProductIDs, 2)
	assert.Equal(t, 199.99, cmd.TotalAmount)
	assert.Equal(t, "New York", cmd.Address.City)
}

// TestCommandResultEdgeCases 测试 CommandResult 边界条件
func TestCommandResultEdgeCases(t *testing.T) {
	t.Run("空结果", func(t *testing.T) {
		result := CommandResult{}
		assert.False(t, result.Success)
		assert.Empty(t, result.Message)
		assert.Nil(t, result.Data)
	})

	t.Run("布尔数据", func(t *testing.T) {
		result := CommandResult{
			Success: true,
			Message: "布尔测试",
			Data:    true,
		}
		data, ok := result.Data.(bool)
		assert.True(t, ok)
		assert.True(t, data)
	})

	t.Run("整数数据", func(t *testing.T) {
		result := CommandResult{
			Success: true,
			Message: "整数测试",
			Data:    42,
		}
		data, ok := result.Data.(int)
		assert.True(t, ok)
		assert.Equal(t, 42, data)
	})

	t.Run("Map数据", func(t *testing.T) {
		result := CommandResult{
			Success: true,
			Message: "Map测试",
			Data: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		}
		data, ok := result.Data.(map[string]string)
		assert.True(t, ok)
		assert.Equal(t, "value1", data["key1"])
	})
}

// ==================== 性能测试 ====================

// BenchmarkCommand_Execute Command 执行性能测试
func BenchmarkCommand_Execute(b *testing.B) {
	ctx := context.Background()
	cmd := MockCommand{shouldError: false}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cmd.Execute(ctx)
	}
}

// BenchmarkCommandHandler_Handle Command Handler 性能测试
func BenchmarkCommandHandler_Handle(b *testing.B) {
	ctx := context.Background()
	handler := &MockCommandHandler{returnError: nil}
	cmd := MockCommand{shouldError: false}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.Handle(ctx, cmd)
	}
}

// BenchmarkCommandResult_Create CommandResult 创建性能测试
func BenchmarkCommandResult_Create(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CommandResult{
			Success: true,
			Message: "操作成功",
			Data:    "test-data",
		}
	}
}

// ==================== 示例函数 ====================

// ExampleCommandResult CommandResult 使用示例
func ExampleCommandResult() {
	// 创建成功结果
	result := CommandResult{
		Success: true,
		Message: "用户创建成功",
		Data:    "user-id-123",
	}

	_ = result
	// Output:
}

// ExampleMockCommand MockCommand 使用示例
func ExampleMockCommand() {
	ctx := context.Background()
	cmd := MockCommand{shouldError: false}
	_ = cmd.Execute(ctx)
	// Output:
}
