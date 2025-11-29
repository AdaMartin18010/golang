package framework

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContext 测试上下文
type TestContext struct {
	Context context.Context
	Cancel  context.CancelFunc
	T       *testing.T
	Cleanup []func()
}

// NewTestContext 创建测试上下文
func NewTestContext(t *testing.T) *TestContext {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	return &TestContext{
		Context: ctx,
		Cancel:  cancel,
		T:       t,
		Cleanup: make([]func(), 0),
	}
}

// AddCleanup 添加清理函数
func (tc *TestContext) AddCleanup(fn func()) {
	tc.Cleanup = append(tc.Cleanup, fn)
}

// CleanupAll 执行所有清理函数
func (tc *TestContext) CleanupAll() {
	for i := len(tc.Cleanup) - 1; i >= 0; i-- {
		tc.Cleanup[i]()
	}
}

// DeferCleanup 延迟执行清理（在测试结束时）
func (tc *TestContext) DeferCleanup() {
	tc.T.Cleanup(func() {
		tc.CleanupAll()
		tc.Cancel()
	})
}

// AssertNoError 断言没有错误
func (tc *TestContext) AssertNoError(err error, msgAndArgs ...interface{}) {
	require.NoError(tc.T, err, msgAndArgs...)
}

// AssertError 断言有错误
func (tc *TestContext) AssertError(err error, msgAndArgs ...interface{}) {
	require.Error(tc.T, err, msgAndArgs...)
}

// AssertEqual 断言相等
func (tc *TestContext) AssertEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	assert.Equal(tc.T, expected, actual, msgAndArgs...)
}

// AssertNotNil 断言不为 nil
func (tc *TestContext) AssertNotNil(value interface{}, msgAndArgs ...interface{}) {
	require.NotNil(tc.T, value, msgAndArgs...)
}

// AssertTrue 断言为 true
func (tc *TestContext) AssertTrue(condition bool, msgAndArgs ...interface{}) {
	assert.True(tc.T, condition, msgAndArgs...)
}

// AssertFalse 断言为 false
func (tc *TestContext) AssertFalse(condition bool, msgAndArgs ...interface{}) {
	assert.False(tc.T, condition, msgAndArgs...)
}

// DatabaseHelper 数据库测试辅助工具
type DatabaseHelper struct {
	DB       *sql.DB
	Driver   string
	DSN      string
	TestDB   string
	OriginalDB string
}

// NewDatabaseHelper 创建数据库辅助工具
func NewDatabaseHelper(driver, dsn, testDB string) *DatabaseHelper {
	return &DatabaseHelper{
		Driver:   driver,
		DSN:      dsn,
		TestDB:   testDB,
	}
}

// Setup 设置测试数据库
func (h *DatabaseHelper) Setup(t *testing.T) error {
	// 连接到主数据库
	db, err := sql.Open(h.Driver, h.DSN)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// 创建测试数据库
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", h.TestDB))
	if err != nil {
		return fmt.Errorf("failed to create test database: %w", err)
	}

	// 连接到测试数据库
	testDSN := fmt.Sprintf("%s/%s", h.DSN, h.TestDB)
	h.DB, err = sql.Open(h.Driver, testDSN)
	if err != nil {
		return fmt.Errorf("failed to open test database: %w", err)
	}

	return nil
}

// Teardown 清理测试数据库
func (h *DatabaseHelper) Teardown(t *testing.T) error {
	if h.DB != nil {
		h.DB.Close()
	}

	// 连接到主数据库
	db, err := sql.Open(h.Driver, h.DSN)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// 删除测试数据库
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", h.TestDB))
	if err != nil {
		return fmt.Errorf("failed to drop test database: %w", err)
	}

	return nil
}

// HTTPTestHelper HTTP 测试辅助工具
type HTTPTestHelper struct {
	BaseURL string
	Headers map[string]string
}

// NewHTTPTestHelper 创建 HTTP 测试辅助工具
func NewHTTPTestHelper(baseURL string) *HTTPTestHelper {
	return &HTTPTestHelper{
		BaseURL: baseURL,
		Headers: make(map[string]string),
	}
}

// SetHeader 设置请求头
func (h *HTTPTestHelper) SetHeader(key, value string) {
	h.Headers[key] = value
}

// SetAuthToken 设置认证令牌
func (h *HTTPTestHelper) SetAuthToken(token string) {
	h.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
}

// MockHelper Mock 辅助工具
type MockHelper struct {
	Mocks map[string]interface{}
}

// NewMockHelper 创建 Mock 辅助工具
func NewMockHelper() *MockHelper {
	return &MockHelper{
		Mocks: make(map[string]interface{}),
	}
}

// RegisterMock 注册 Mock
func (h *MockHelper) RegisterMock(name string, mock interface{}) {
	h.Mocks[name] = mock
}

// GetMock 获取 Mock
func (h *MockHelper) GetMock(name string) (interface{}, bool) {
	mock, ok := h.Mocks[name]
	return mock, ok
}

// TestDataHelper 测试数据辅助工具
type TestDataHelper struct {
	Data map[string]interface{}
}

// NewTestDataHelper 创建测试数据辅助工具
func NewTestDataHelper() *TestDataHelper {
	return &TestDataHelper{
		Data: make(map[string]interface{}),
	}
}

// Set 设置测试数据
func (h *TestDataHelper) Set(key string, value interface{}) {
	h.Data[key] = value
}

// Get 获取测试数据
func (h *TestDataHelper) Get(key string) (interface{}, bool) {
	value, ok := h.Data[key]
	return value, ok
}

// GetString 获取字符串类型测试数据
func (h *TestDataHelper) GetString(key string) (string, bool) {
	value, ok := h.Get(key)
	if !ok {
		return "", false
	}
	str, ok := value.(string)
	return str, ok
}

// GetInt 获取整数类型测试数据
func (h *TestDataHelper) GetInt(key string) (int, bool) {
	value, ok := h.Get(key)
	if !ok {
		return 0, false
	}
	i, ok := value.(int)
	return i, ok
}

// EnvironmentHelper 环境变量辅助工具
type EnvironmentHelper struct {
	OriginalEnv map[string]string
}

// NewEnvironmentHelper 创建环境变量辅助工具
func NewEnvironmentHelper() *EnvironmentHelper {
	return &EnvironmentHelper{
		OriginalEnv: make(map[string]string),
	}
}

// SetEnv 设置环境变量（测试后自动恢复）
func (h *EnvironmentHelper) SetEnv(key, value string) {
	// 保存原始值
	if original, exists := os.LookupEnv(key); exists {
		h.OriginalEnv[key] = original
	} else {
		h.OriginalEnv[key] = "" // 标记为不存在
	}
	os.Setenv(key, value)
}

// Restore 恢复原始环境变量
func (h *EnvironmentHelper) Restore() {
	for key, value := range h.OriginalEnv {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
	h.OriginalEnv = make(map[string]string)
}

// RetryHelper 重试辅助工具
type RetryHelper struct {
	MaxAttempts int
	Delay       time.Duration
}

// NewRetryHelper 创建重试辅助工具
func NewRetryHelper(maxAttempts int, delay time.Duration) *RetryHelper {
	return &RetryHelper{
		MaxAttempts: maxAttempts,
		Delay:       delay,
	}
}

// Retry 重试执行函数
func (h *RetryHelper) Retry(fn func() error) error {
	var lastErr error
	for i := 0; i < h.MaxAttempts; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
			if i < h.MaxAttempts-1 {
				time.Sleep(h.Delay)
			}
		}
	}
	return fmt.Errorf("retry failed after %d attempts: %w", h.MaxAttempts, lastErr)
}

// CoverageHelper 覆盖率辅助工具
type CoverageHelper struct {
	Package string
}

// NewCoverageHelper 创建覆盖率辅助工具
func NewCoverageHelper(packagePath string) *CoverageHelper {
	return &CoverageHelper{
		Package: packagePath,
	}
}

// CheckCoverage 检查覆盖率（需要在测试中调用）
func (h *CoverageHelper) CheckCoverage(t *testing.T, minCoverage float64) {
	// 这个功能需要在测试外部通过 go test -cover 来检查
	// 这里只是提供一个接口
	t.Logf("Coverage check for package %s (minimum: %.2f%%)", h.Package, minCoverage*100)
}
