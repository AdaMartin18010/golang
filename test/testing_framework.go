package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// BaseTestSuite 基础测试套件
// 提供通用的测试工具和辅助方法
type BaseTestSuite struct {
	suite.Suite
	ctx context.Context
}

// SetupSuite 在整个测试套件开始前执行一次
func (s *BaseTestSuite) SetupSuite() {
	s.ctx = context.Background()
}

// SetupTest 在每个测试方法执行前执行
func (s *BaseTestSuite) SetupTest() {
	// 每个测试都有独立的 context
	s.ctx = context.Background()
}

// TearDownTest 在每个测试方法执行后执行
func (s *BaseTestSuite) TearDownTest() {
	// 清理资源
}

// TearDownSuite 在整个测试套件结束后执行一次
func (s *BaseTestSuite) TearDownSuite() {
	// 清理全局资源
}

// AssertNoError 断言没有错误（使用 require，失败时立即停止）
func (s *BaseTestSuite) AssertNoError(err error, msgAndArgs ...interface{}) {
	require.NoError(s.T(), err, msgAndArgs...)
}

// AssertError 断言有错误
func (s *BaseTestSuite) AssertError(err error, msgAndArgs ...interface{}) {
	require.Error(s.T(), err, msgAndArgs...)
}

// AssertEqual 断言相等
func (s *BaseTestSuite) AssertEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	assert.Equal(s.T(), expected, actual, msgAndArgs...)
}

// AssertNotNil 断言非空
func (s *BaseTestSuite) AssertNotNil(object interface{}, msgAndArgs ...interface{}) {
	assert.NotNil(s.T(), object, msgAndArgs...)
}

// Context 获取测试上下文
func (s *BaseTestSuite) Context() context.Context {
	return s.ctx
}

// TableTest 表格驱动测试辅助函数
type TableTest[T any] struct {
	Name     string
	Input    T
	Expected interface{}
	WantErr  bool
	Setup    func(*testing.T)
	Teardown func(*testing.T)
}

// RunTableTests 运行表格驱动测试
func RunTableTests[T any](t *testing.T, tests []TableTest[T], testFunc func(*testing.T, T) (interface{}, error)) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Setup
			if tt.Setup != nil {
				tt.Setup(t)
			}

			// 执行测试
			result, err := testFunc(t, tt.Input)

			// 验证错误
			if tt.WantErr {
				require.Error(t, err, "Expected error but got nil")
			} else {
				require.NoError(t, err, "Unexpected error: %v", err)
				if tt.Expected != nil {
					assert.Equal(t, tt.Expected, result, "Result mismatch")
				}
			}

			// Teardown
			if tt.Teardown != nil {
				tt.Teardown(t)
			}
		})
	}
}

// MockContext 创建 mock context（用于测试）
func MockContext() context.Context {
	return context.Background()
}

// MockContextWithTimeout 创建带超时的 mock context
func MockContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
