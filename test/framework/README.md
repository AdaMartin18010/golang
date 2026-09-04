# 测试框架工具

> **状态**: ✅ 基础实现完成
> **版本**: v1.0.0
> **优先级**: P0 - 测试提升

---

## 📋 概述

本包提供了完整的测试辅助工具，用于简化测试编写和提高测试效率。

---

## 🚀 快速开始

### TestContext - 测试上下文

```go
func TestExample(t *testing.T) {
    tc := NewTestContext(t)
    defer tc.DeferCleanup()

    // 添加清理函数
    tc.AddCleanup(func() {
        // 清理资源
    })

    // 使用断言
    tc.AssertNoError(err, "should not have error")
    tc.AssertEqual(expected, actual, "should be equal")
}
```

### DatabaseHelper - 数据库测试辅助

```go
func TestDatabase(t *testing.T) {
    helper := NewDatabaseHelper("postgres", "postgres://user:pass@localhost/postgres", "test_db")

    err := helper.Setup(t)
    require.NoError(t, err)
    defer helper.Teardown(t)

    // 使用 helper.DB 进行测试
}
```

### HTTPTestHelper - HTTP 测试辅助

```go
func TestHTTP(t *testing.T) {
    helper := NewHTTPTestHelper("http://localhost:8080")
    helper.SetAuthToken("token123")

    // 使用 helper 进行 HTTP 测试
}
```

### RetryHelper - 重试辅助

```go
func TestRetry(t *testing.T) {
    helper := NewRetryHelper(3, 100*time.Millisecond)

    err := helper.Retry(func() error {
        // 可能失败的操作
        return doSomething()
    })

    require.NoError(t, err)
}
```

### EnvironmentHelper - 环境变量辅助

```go
func TestEnv(t *testing.T) {
    helper := NewEnvironmentHelper()
    defer helper.Restore()

    helper.SetEnv("TEST_VAR", "test_value")
    // 测试代码
    // 测试结束后自动恢复原始值
}
```

### TestDataHelper - 测试数据辅助

```go
func TestData(t *testing.T) {
    helper := NewTestDataHelper()

    helper.Set("user_id", "123")
    helper.Set("count", 42)

    userID, _ := helper.GetString("user_id")
    count, _ := helper.GetInt("count")
}
```

---

## 📚 API 文档

### TestContext

- `NewTestContext(t *testing.T) *TestContext` - 创建测试上下文
- `AddCleanup(fn func())` - 添加清理函数
- `CleanupAll()` - 执行所有清理函数
- `DeferCleanup()` - 延迟执行清理（在测试结束时）
- `AssertNoError(err error, ...)` - 断言没有错误
- `AssertError(err error, ...)` - 断言有错误
- `AssertEqual(expected, actual, ...)` - 断言相等
- `AssertNotNil(value interface{}, ...)` - 断言不为 nil
- `AssertTrue(condition bool, ...)` - 断言为 true
- `AssertFalse(condition bool, ...)` - 断言为 false

### DatabaseHelper

- `NewDatabaseHelper(driver, dsn, testDB string) *DatabaseHelper` - 创建数据库辅助工具
- `Setup(t *testing.T) error` - 设置测试数据库
- `Teardown(t *testing.T) error` - 清理测试数据库

### HTTPTestHelper

- `NewHTTPTestHelper(baseURL string) *HTTPTestHelper` - 创建 HTTP 测试辅助工具
- `SetHeader(key, value string)` - 设置请求头
- `SetAuthToken(token string)` - 设置认证令牌

### RetryHelper

- `NewRetryHelper(maxAttempts int, delay time.Duration) *RetryHelper` - 创建重试辅助工具
- `Retry(fn func() error) error` - 重试执行函数

### EnvironmentHelper

- `NewEnvironmentHelper() *EnvironmentHelper` - 创建环境变量辅助工具
- `SetEnv(key, value string)` - 设置环境变量（测试后自动恢复）
- `Restore()` - 恢复原始环境变量

### TestDataHelper

- `NewTestDataHelper() *TestDataHelper` - 创建测试数据辅助工具
- `Set(key string, value interface{})` - 设置测试数据
- `Get(key string) (interface{}, bool)` - 获取测试数据
- `GetString(key string) (string, bool)` - 获取字符串类型测试数据
- `GetInt(key string) (int, bool)` - 获取整数类型测试数据

---

## 🧪 测试

运行测试：

```bash
go test -v ./test/framework/...
```

---

## 📝 使用示例

### 完整示例

```go
package example_test

import (
    "testing"
    "github.com/yourusername/golang/test/framework"
)

func TestCompleteExample(t *testing.T) {
    // 创建测试上下文
    tc := NewTestContext(t)
    defer tc.DeferCleanup()

    // 设置环境变量
    envHelper := NewEnvironmentHelper()
    defer envHelper.Restore()
    envHelper.SetEnv("API_KEY", "test_key")

    // 创建测试数据
    dataHelper := NewTestDataHelper()
    dataHelper.Set("user_id", "123")

    // 添加清理函数
    tc.AddCleanup(func() {
        // 清理资源
    })

    // 执行测试
    userID, _ := dataHelper.GetString("user_id")
    tc.AssertEqual("123", userID, "user ID should match")

    // 使用重试
    retryHelper := NewRetryHelper(3, 100*time.Millisecond)
    err := retryHelper.Retry(func() error {
        // 可能失败的操作
        return nil
    })
    tc.AssertNoError(err, "retry should succeed")
}
```

---

## 🔗 相关文档

- 改进任务看板
- 改进路线图

---

**最后更新**: 2025-01-XX
