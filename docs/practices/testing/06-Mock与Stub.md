# Mock与Stub

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

---

## 📋 目录

- [Mock与Stub](#mock与stub)
  - [📋 目录](#-目录)
  - [1. 📖 概念介绍](#1--概念介绍)
  - [2. 🎯 核心知识点](#2--核心知识点)
    - [2.1 手动Mock](#21-手动mock)
    - [2.2 使用testify/mock](#22-使用testifymock)
    - [2.3 使用gomock](#23-使用gomock)
    - [2.4 HTTP Mock](#24-http-mock)
    - [2.5 数据库Mock](#25-数据库mock)
    - [2.6 Time Mock](#26-time-mock)
  - [3. 💡 最佳实践](#3--最佳实践)
    - [3.1 优先使用接口](#31-优先使用接口)
    - [3.2 Mock应该简单](#32-mock应该简单)
    - [3.3 只Mock外部依赖](#33-只mock外部依赖)
    - [3.4 验证交互](#34-验证交互)
  - [4. ⚠️ 常见问题](#4-️-常见问题)
  - [5. 📚 相关资源](#5--相关资源)

---

## 1. 📖 概念介绍

Mock和Stub是测试中隔离依赖的技术。Mock用于验证交互行为，Stub用于提供预定义的响应。Go通过接口和工具实现优雅的Mock。

---

## 2. 🎯 核心知识点

### 2.1 手动Mock

```go
// 定义接口
type UserRepository interface {
    GetByID(id int) (*User, error)
    Create(user *User) error
}

// Mock实现
type MockUserRepository struct {
    GetByIDFunc func(id int) (*User, error)
    CreateFunc  func(user *User) error
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
    if m.GetByIDFunc != nil {
        return m.GetByIDFunc(id)
    }
    return nil, nil
}

func (m *MockUserRepository) Create(user *User) error {
    if m.CreateFunc != nil {
        return m.CreateFunc(user)
    }
    return nil
}

// 使用Mock
func TestUserService(t *testing.T) {
    mock := &MockUserRepository{
        GetByIDFunc: func(id int) (*User, error) {
            return &User{ID: id, Name: "Alice"}, nil
        },
    }

    service := NewUserService(mock)
    user, err := service.GetUser(1)

    if err != nil {
        t.Fatal(err)
    }
    if user.Name != "Alice" {
        t.Errorf("got %s, want Alice", user.Name)
    }
}
```

---

### 2.2 使用testify/mock

```bash
go get github.com/stretchr/testify/mock
```

```go
package service

import (
    "testing"
    "github.com/stretchr/testify/mock"
)

// Mock结构
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Create(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

// 测试
func TestUserService_GetUser(t *testing.T) {
    mockRepo := new(MockUserRepository)

    // 设置期望
    expectedUser := &User{ID: 1, Name: "Alice"}
    mockRepo.On("GetByID", 1).Return(expectedUser, nil)

    service := NewUserService(mockRepo)
    user, err := service.GetUser(1)

    // 断言
    assert.NoError(t, err)
    assert.Equal(t, "Alice", user.Name)

    // 验证Mock被正确调用
    mockRepo.AssertExpectations(t)
}
```

---

### 2.3 使用gomock

```bash
go install github.com/golang/mock/mockgen@latest

# 生成mock
mockgen -source=repository.go -destination=mock_repository.go -package=service
```

```go
// repository.go
package service

type UserRepository interface {
    GetByID(id int) (*User, error)
    Create(user *User) error
}
```

生成的mock：

```go
// mock_repository.go (自动生成)
package service

import (
    gomock "github.com/golang/mock/gomock"
)

type MockUserRepository struct {
    ctrl     *gomock.Controller
    recorder *MockUserRepositoryMockRecorder
}

// ... 其他自动生成的代码
```

使用：

```go
func TestUserService_WithGomock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockUserRepository(ctrl)

    // 设置期望
    expectedUser := &User{ID: 1, Name: "Alice"}
    mockRepo.EXPECT().
        GetByID(1).
        Return(expectedUser, nil).
        Times(1)

    service := NewUserService(mockRepo)
    user, err := service.GetUser(1)

    if err != nil {
        t.Fatal(err)
    }
    if user.Name != "Alice" {
        t.Errorf("got %s, want Alice", user.Name)
    }
}
```

---

### 2.4 HTTP Mock

```go
package api

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestFetchUser(t *testing.T) {
    // 创建mock服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 验证请求
        if r.URL.Path != "/api/users/123" {
            t.Errorf("Expected path /api/users/123, got %s", r.URL.Path)
        }

        // 返回mock响应
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id":"123","name":"Alice"}`))
    }))
    defer server.Close()

    // 使用mock URL
    client := NewAPIClient(server.URL)
    user, err := client.FetchUser("123")

    if err != nil {
        t.Fatal(err)
    }
    if user.Name != "Alice" {
        t.Errorf("got %s, want Alice", user.Name)
    }
}
```

---

### 2.5 数据库Mock

```go
package repository

import (
    "database/sql"
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepository_GetByID(t *testing.T) {
    // 创建mock数据库
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create mock: %v", err)
    }
    defer db.Close()

    // 设置期望查询
    rows := sqlmock.NewRows([]string{"id", "name", "email"}).
        AddRow(1, "Alice", "alice@example.com")

    mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
        WithArgs(1).
        WillReturnRows(rows)

    // 测试
    repo := NewUserRepository(db)
    user, err := repo.GetByID(1)

    if err != nil {
        t.Fatal(err)
    }
    if user.Name != "Alice" {
        t.Errorf("got %s, want Alice", user.Name)
    }

    // 验证所有期望都被满足
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %s", err)
    }
}
```

---

### 2.6 Time Mock

```go
package service

import (
    "testing"
    "time"
)

// 定义时间接口
type Clock interface {
    Now() time.Time
}

// 真实实现
type RealClock struct{}

func (RealClock) Now() time.Time {
    return time.Now()
}

// Mock实现
type MockClock struct {
    current time.Time
}

func (m *MockClock) Now() time.Time {
    return m.current
}

// 使用Clock
type Service struct {
    clock Clock
}

func (s *Service) IsExpired(deadline time.Time) bool {
    return s.clock.Now().After(deadline)
}

// 测试
func TestService_IsExpired(t *testing.T) {
    mockClock := &MockClock{
        current: time.Date(2025, 10, 28, 12, 0, 0, 0, time.UTC),
    }

    service := &Service{clock: mockClock}

    deadline := time.Date(2025, 10, 28, 11, 0, 0, 0, time.UTC)

    if !service.IsExpired(deadline) {
        t.Error("Expected to be expired")
    }
}
```

---

## 3. 💡 最佳实践

### 3.1 优先使用接口

```go
// ✅ 好：基于接口
type Storage interface {
    Save(data []byte) error
}

// ❌ 差：依赖具体类型
type Service struct {
    storage *FileStorage
}
```

### 3.2 Mock应该简单

```go
// ✅ 好：简单的Mock
type MockStorage struct {
    SaveFunc func([]byte) error
}

// ❌ 差：过于复杂
type MockStorage struct {
    calls []Call
    expectations []Expectation
    // ... 太多状态
}
```

### 3.3 只Mock外部依赖

```go
// ✅ Mock：数据库、HTTP、文件系统
mockDB := NewMockDB()
mockHTTP := httptest.NewServer(...)

// ❌ 不要Mock：内部业务逻辑
```

### 3.4 验证交互

```go
func TestService_Delete(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("Delete", 1).Return(nil)

    service := NewUserService(mockRepo)
    service.DeleteUser(1)

    // 验证Delete被调用了一次
    mockRepo.AssertCalled(t, "Delete", 1)
    mockRepo.AssertNumberOfCalls(t, "Delete", 1)
}
```

---

## 4. ⚠️ 常见问题

- **Mock**: 验证交互行为（关注行为）

**Q2: 何时使用Mock？**

- 外部依赖（数据库、API）
- 难以构造的场景（错误情况）
- 昂贵的操作（网络、IO）

**Q3: Mock过多是坏味道吗？**

- 是的，可能表示设计问题
- 考虑重构，减少依赖
- 使用依赖注入

**Q4: 如何Mock私有函数？**

- 不应该Mock私有函数
- 通过公有接口测试
- 或者重构为可测试结构

---

## 5. 📚 相关资源

- [testify/mock](https://github.com/stretchr/testify)
- [gomock](https://github.com/golang/mock)
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- [httptest](https://pkg.go.dev/net/http/httptest)

**下一步**: [07-测试最佳实践](./07-测试最佳实践.md)

---

**最后更新**: 2025-10-29
