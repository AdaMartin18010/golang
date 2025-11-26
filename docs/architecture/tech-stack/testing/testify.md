# 1. 🧪 Testify 深度解析

> **简介**: 本文档详细阐述了 Testify 的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. 🧪 Testify 深度解析](#1--testify-深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 断言库](#131-断言库)
    - [1.3.2 Mock 框架](#132-mock-框架)
    - [1.3.3 Suite 测试套件](#133-suite-测试套件)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 测试最佳实践](#141-测试最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Testify 是什么？**

Testify 是一个 Go 测试工具包，提供了断言、Mock 和测试套件功能。

**核心特性**:

- ✅ **断言库**: 丰富的断言函数
- ✅ **Mock 框架**: 强大的 Mock 功能
- ✅ **测试套件**: 支持测试套件组织
- ✅ **易用性**: API 简洁，易于使用

---

## 1.2 选型论证

**为什么选择 Testify？**

**论证矩阵**:

| 评估维度 | 权重 | Testify | gocheck | GoConvey | 标准库 testing | 说明 |
|---------|------|---------|---------|----------|----------------|------|
| **功能完整性** | 30% | 10 | 8 | 9 | 5 | Testify 功能最完整 |
| **易用性** | 25% | 10 | 7 | 9 | 6 | Testify 易用性最好 |
| **Mock 支持** | 20% | 10 | 6 | 5 | 3 | Testify Mock 最强大 |
| **社区支持** | 15% | 10 | 6 | 7 | 10 | Testify 社区最活跃 |
| **性能** | 10% | 9 | 9 | 8 | 10 | Testify 性能优秀 |
| **加权总分** | - | **9.80** | 7.40 | 7.90 | 6.40 | Testify 得分最高 |

**核心优势**:

1. **功能完整性（权重 30%）**:
   - 断言、Mock、测试套件一体化
   - 功能丰富，覆盖全面
   - 适合各种测试场景

2. **易用性（权重 25%）**:
   - API 简洁，易于使用
   - 文档完善，示例丰富
   - 学习成本低

---

## 1.3 实际应用

### 1.3.1 断言库

**使用断言库**:

```go
// internal/domain/user/service_test.go
package user

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
    service := NewService(nil)

    user, err := service.CreateUser("test@example.com", "Test User")

    // 使用 assert（失败继续执行）
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
    assert.Equal(t, "Test User", user.Name)

    // 使用 require（失败立即停止）
    require.NotEmpty(t, user.ID)
}

func TestCreateUser_InvalidEmail(t *testing.T) {
    service := NewService(nil)

    user, err := service.CreateUser("invalid-email", "Test User")

    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "invalid email")
}
```

**常用断言函数**:

```go
// 相等性断言
assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)

// 布尔断言
assert.True(t, condition)
assert.False(t, condition)

// Nil 断言
assert.Nil(t, value)
assert.NotNil(t, value)

// 错误断言
assert.NoError(t, err)
assert.Error(t, err)

// 包含断言
assert.Contains(t, container, item)
assert.NotContains(t, container, item)

// 长度断言
assert.Len(t, slice, length)

// 类型断言
assert.IsType(t, expectedType, actual)
```

### 1.3.2 Mock 框架

**使用 Mock 框架**:

```go
// 定义接口
type Repository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
}

// 生成 Mock（使用 mockery）
//go:generate mockery --name=Repository

// 使用 Mock
func TestService_CreateUser(t *testing.T) {
    // 创建 Mock
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)

    // 设置期望
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).
        Return(nil).
        Once()

    // 执行测试
    user, err := service.CreateUser("test@example.com", "Test User")

    // 验证
    assert.NoError(t, err)
    assert.NotNil(t, user)
    mockRepo.AssertExpectations(t)
}

// 手动创建 Mock
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}
```

### 1.3.3 Suite 测试套件

**使用测试套件**:

```go
// 定义测试套件
type UserServiceTestSuite struct {
    suite.Suite
    service *Service
    repo    *MockRepository
}

// 设置测试套件
func (suite *UserServiceTestSuite) SetupTest() {
    suite.repo = new(MockRepository)
    suite.service = NewService(suite.repo)
}

// 清理测试套件
func (suite *UserServiceTestSuite) TearDownTest() {
    suite.repo.AssertExpectations(suite.T())
}

// 测试用例
func (suite *UserServiceTestSuite) TestCreateUser() {
    suite.repo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).
        Return(nil).
        Once()

    user, err := suite.service.CreateUser("test@example.com", "Test User")

    suite.NoError(err)
    suite.NotNil(user)
    suite.Equal("test@example.com", user.Email)
}

func (suite *UserServiceTestSuite) TestCreateUser_InvalidEmail() {
    user, err := suite.service.CreateUser("invalid-email", "Test User")

    suite.Error(err)
    suite.Nil(user)
}

// 运行测试套件
func TestUserServiceTestSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}
```

---

## 1.4 最佳实践

### 1.4.1 测试最佳实践

**为什么需要最佳实践？**

合理的测试实践可以提高测试的质量和可维护性。

**最佳实践原则**:

1. **测试组织**: 使用测试套件组织相关测试
2. **Mock 使用**: 合理使用 Mock，避免过度 Mock
3. **断言选择**: 根据场景选择合适的断言
4. **测试覆盖**: 追求合理的测试覆盖率

**实际应用示例**:

```go
// 测试最佳实践
type ServiceTestSuite struct {
    suite.Suite
    service *Service
    repo    *MockRepository
    ctx     context.Context
}

func (suite *ServiceTestSuite) SetupSuite() {
    // 整个套件的初始化
    suite.ctx = context.Background()
}

func (suite *ServiceTestSuite) SetupTest() {
    // 每个测试前的初始化
    suite.repo = new(MockRepository)
    suite.service = NewService(suite.repo)
}

func (suite *ServiceTestSuite) TearDownTest() {
    // 每个测试后的清理
    suite.repo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TearDownSuite() {
    // 整个套件的清理
}

// 表驱动测试
func (suite *ServiceTestSuite) TestCreateUser_Validation() {
    tests := []struct {
        name    string
        email   string
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid email",
            email:   "test@example.com",
            wantErr: false,
        },
        {
            name:    "invalid email",
            email:   "invalid-email",
            wantErr: true,
            errMsg:  "invalid email",
        },
        {
            name:    "empty email",
            email:   "",
            wantErr: true,
            errMsg:  "email required",
        },
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            suite.repo.On("Create", mock.Anything, mock.Anything).
                Return(nil).
                Maybe()

            user, err := suite.service.CreateUser(tt.email, "Test User")

            if tt.wantErr {
                suite.Error(err)
                suite.Nil(user)
                if tt.errMsg != "" {
                    suite.Contains(err.Error(), tt.errMsg)
                }
            } else {
                suite.NoError(err)
                suite.NotNil(user)
            }
        })
    }
}
```

**最佳实践要点**:

1. **测试组织**: 使用测试套件组织相关测试，提高可维护性
2. **Mock 使用**: 合理使用 Mock，避免过度 Mock 导致测试不真实
3. **断言选择**: 使用 `require` 进行关键断言，使用 `assert` 进行一般断言
4. **测试覆盖**: 追求合理的测试覆盖率，关注边界条件和错误处理

---

## 📚 扩展阅读

- [Testify 官方文档](https://github.com/stretchr/testify)
- [Mockery 文档](https://github.com/vektra/mockery)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Testify 的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
