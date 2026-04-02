# LD-009: Go 测试模式与最佳实践 (Go Testing Patterns & Best Practices)

> **维度**: Language Design
> **级别**: S (16+ KB)
> **标签**: #go-testing #tdd #benchmark #mocking #test-coverage
> **权威来源**: [Go Testing](https://go.dev/doc/tutorial/add-a-test), [Testify](https://github.com/stretchr/testify)

---

## 测试金字塔

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Go Testing Pyramid                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                                    ▲                                        │
│                                   /_\                                       │
│                                  / E2E \    Integration tests (少量)         │
│                                 /_______\   (API/Service)                    │
│                                /_________\                                  │
│                               /             \                               │
│                              /  Integration  \  HTTP/DB 集成测试              │
│                             /_________________\ (中等数量)                     │
│                            /                     \                          │
│                           /                       \                         │
│                          /       Unit Tests        \   纯函数/结构体测试       │
│                         /___________________________\  (大量)                 │
│                                                                              │
│  比例建议: Unit (70%) : Integration (20%) : E2E (10%)                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 单元测试

### 表驱动测试

```go
package service

import "testing"

func TestCalculatePrice(t *testing.T) {
    tests := []struct {
        name     string
        qty      int
        price    float64
        discount float64
        want     float64
        wantErr  bool
    }{
        {
            name:     "正常计算",
            qty:      2,
            price:    100,
            discount: 0.1,
            want:     180,
            wantErr:  false,
        },
        {
            name:     "数量为零",
            qty:      0,
            price:    100,
            discount: 0,
            want:     0,
            wantErr:  true,
        },
        {
            name:     "负价格",
            qty:      1,
            price:    -10,
            discount: 0,
            want:     0,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := CalculatePrice(tt.qty, tt.price, tt.discount)
            if (err != nil) != tt.wantErr {
                t.Errorf("CalculatePrice() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("CalculatePrice() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 测试 Fixtures

```go
// testdata/ 目录存放测试数据
// - testdata/users.json
// - testdata/orders/

func TestLoadUsers(t *testing.T) {
    data, err := os.ReadFile("testdata/users.json")
    if err != nil {
        t.Fatalf("failed to load test data: %v", err)
    }

    var users []User
    if err := json.Unmarshal(data, &users); err != nil {
        t.Fatalf("failed to unmarshal: %v", err)
    }

    // 测试逻辑...
}
```

---

## Mock 与依赖注入

### 接口 Mock

```go
// 定义接口
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

// 真实实现
type PostgresUserRepo struct { db *sql.DB }

// Mock 实现
type MockUserRepo struct {
    mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// 测试使用
func TestUserService_GetUser(t *testing.T) {
    mockRepo := new(MockUserRepo)
    service := NewUserService(mockRepo)

    mockRepo.On("GetByID", mock.Anything, "123").
        Return(&User{ID: "123", Name: "Alice"}, nil)

    user, err := service.GetUser(context.Background(), "123")

    assert.NoError(t, err)
    assert.Equal(t, "Alice", user.Name)
    mockRepo.AssertExpectations(t)
}
```

### 测试容器 (Integration)

```go
func TestPostgresRepo(t *testing.T) {
    ctx := context.Background()

    // 启动 PostgreSQL 容器
    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:18-alpine"),
        postgres.WithDatabase("test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    if err != nil {
        t.Fatalf("failed to start container: %v", err)
    }
    defer container.Terminate(ctx)

    connStr, _ := container.ConnectionString(ctx)
    db, _ := sql.Open("postgres", connStr)

    repo := NewPostgresUserRepo(db)

    // 运行测试...
}
```

---

## Benchmark

```go
func BenchmarkCalculatePrice(b *testing.B) {
    for i := 0; i < b.N; i++ {
        CalculatePrice(10, 99.99, 0.15)
    }
}

// 带子基准测试
func BenchmarkCalculatePriceParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            CalculatePrice(10, 99.99, 0.15)
        }
    })
}

// 内存分配分析
func BenchmarkWithAlloc(b *testing.B) {
    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        heavyAllocation()
    }
}
```

---

## 测试覆盖率

```bash
# 运行测试
make test

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看 HTML 报告
go tool cover -html=coverage.out

# 设置覆盖率阈值
make test-coverage  # 需要 > 80%
```

---

## 参考文献

1. [Go Testing](https://go.dev/doc/tutorial/add-a-test)
2. [Testify](https://github.com/stretchr/testify)
3. [Go Test Containers](https://golang.testcontainers.org/)
