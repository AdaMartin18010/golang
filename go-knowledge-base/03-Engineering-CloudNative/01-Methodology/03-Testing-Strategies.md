# EC-M03: Testing Strategies in Go (S-Level)

> **维度**: Engineering-CloudNative / Methodology
> **级别**: S (15+ KB)
> **标签**: #testing #go #unit-test #integration-test #tdd #benchmark #mock
> **权威来源**:
>
> - [Test-Driven Development by Example](https://en.wikipedia.org/wiki/Test-Driven_Development_by_Example) - Kent Beck (2002)
> - [Unit Testing Principles, Practices, and Patterns](https://www.manning.com/books/unit-testing) - Vladimir Khorikov (2020)
> - [Go Testing](https://go.dev/doc/testing) - The Go Authors

---

## 1. 测试金字塔

```
                    /\
                   /  \
                  / E2E \          <- 少量 (5%)
                 /________\
                /          \
               / Integration \     <- 中等 (15%)
              /______________\
             /                \
            /     Unit Test     \   <- 大量 (80%)
           /______________________\
```

---

## 2. 单元测试

### 2.1 基本测试结构

```go
package service

import "testing"

func TestCalculateTotal(t *testing.T) {
    // Arrange
    items := []Item{
        {Price: 10.0, Quantity: 2},
        {Price: 20.0, Quantity: 1},
    }

    // Act
    total := CalculateTotal(items)

    // Assert
    expected := 40.0
    if total != expected {
        t.Errorf("expected %f, got %f", expected, total)
    }
}
```

### 2.2 表驱动测试

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name     string
        dividend float64
        divisor  float64
        want     float64
        wantErr  bool
    }{
        {"normal", 10, 2, 5, false},
        {"negative", -10, 2, -5, false},
        {"by zero", 10, 0, 0, true},
        {"fraction", 5, 2, 2.5, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.dividend, tt.divisor)
            if (err != nil) != tt.wantErr {
                t.Errorf("Divide() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Divide() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 2.3 Mock 测试

```go
// 使用 testify/mock
import "github.com/stretchr/testify/mock"

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Get(id string) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

func TestService_GetUser(t *testing.T) {
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)

    expectedUser := &User{ID: "1", Name: "John"}
    mockRepo.On("Get", "1").Return(expectedUser, nil)

    user, err := service.GetUser("1")

    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}
```

---

## 3. 集成测试

```go
// +build integration

func TestDatabaseIntegration(t *testing.T) {
    db, err := sql.Open("postgres", testDSN)
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    repo := NewRepository(db)

    // 使用事务回滚保持测试隔离
    tx, _ := db.Begin()
    defer tx.Rollback()

    // 测试操作...
}
```

---

## 4. 基准测试

```go
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}

func BenchmarkFibonacciParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Fibonacci(20)
        }
    })
}
```

---

## 5. 模糊测试

```go
func FuzzParseEmail(f *testing.F) {
    f.Add("test@example.com")
    f.Fuzz(func(t *testing.T, email string) {
        _, err := ParseEmail(email)
        // 不应该 panic
    })
}
```

---

## 6. 测试最佳实践

### 6.1 AAA 模式

- Arrange: 准备测试数据
- Act: 执行被测代码
- Assert: 验证结果

### 6.2 测试命名规范

```
Test{FunctionName}_{Scenario}_{ExpectedResult}
TestDivide_ByZero_ReturnsError
TestDivide_Normal_ReturnsQuotient
```

### 6.3 覆盖率目标

- 单元测试: >80%
- 关键路径: >90%
- 集成测试: 核心流程全覆盖

---

## 7. 测试工具

| 工具 | 用途 |
|------|------|
| testing | Go 标准测试库 |
| testify | 断言和 mock |
| gomock | 代码生成 mock |
| httptest | HTTP 测试 |
| dockertest | 容器集成测试 |

---

## 8. 生产检查清单

```
Testing Checklist:
□ 所有导出函数都有测试
□ 表驱动测试覆盖边界情况
□ 使用 t.Parallel() 加速测试
□ Mock 外部依赖
□ 集成测试标记 +build integration
□ 基准测试验证性能
□ CI 中运行测试并检查覆盖率
```

---

**质量评级**: S (15+ KB)
