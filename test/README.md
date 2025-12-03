# 测试框架

**版本**: v2.0
**更新日期**: 2025-12-03
**目标**: 测试覆盖率 > 80%

---

## 🎯 测试策略

### 测试金字塔

```text
       ┌─────────┐
       │   E2E   │  10% - 端到端测试
       ├─────────┤
       │  集成测试 │  20% - 集成测试
       ├─────────┤
       │  单元测试 │  70% - 单元测试
       └─────────┘
```

### 测试类型

| 类型 | 占比 | 工具 | 目标 |
|------|------|------|------|
| **单元测试** | 70% | testify | 覆盖率 > 80% |
| **集成测试** | 20% | testcontainers | 关键路径 |
| **E2E 测试** | 10% | httptest | 核心场景 |

---

## 🛠️ 测试工具

### 核心库

```go
// go.mod
require (
    github.com/stretchr/testify v1.9.0      // 断言和 Mock
    github.com/testcontainers/testcontainers-go v0.33.0  // 容器测试
    github.com/golang/mock v1.6.0           // Mock 生成
)
```

### 工具链

- **testify** - 断言、Mock、Suite
- **testcontainers** - 集成测试容器
- **gomock** - Mock 生成（可选）
- **httptest** - HTTP 测试
- **sqlmock** - 数据库 Mock

---

## 📁 测试结构

```text
test/
├── testing_framework.go      # 测试框架基础 ✅
├── mocks/
│   ├── repository_mock.go   # Repository Mock ✅
│   └── service_mock.go      # Service Mock
├── fixtures/
│   ├── users.json           # 测试数据
│   └── config.yaml          # 测试配置
├── integration/
│   ├── database_test.go     # 数据库集成测试
│   ├── api_test.go          # API 集成测试
│   └── messaging_test.go    # 消息队列集成测试
└── e2e/
    ├── user_flow_test.go    # 用户流程 E2E
    └── api_flow_test.go     # API 流程 E2E

internal/
├── domain/
│   └── user/
│       ├── entity_test.go   # 实体测试 ✅
│       └── specifications/
│           └── user_spec_test.go
├── application/
│   └── user/
│       └── service_test.go  # 服务测试 ✅
└── infrastructure/
    └── database/
        └── ent/
            └── repository/
                └── user_repository_test.go
```

---

## 🚀 快速开始

### 1. 运行所有测试

```bash
# 运行所有测试
make test

# 运行单元测试
go test ./internal/... -v

# 运行集成测试
go test ./test/integration/... -v

# 运行 E2E 测试
go test ./test/e2e/... -v
```

### 2. 查看覆盖率

```bash
# 生成覆盖率报告
make coverage

# 查看 HTML 报告
go tool cover -html=coverage.out
```

### 3. 运行特定测试

```bash
# 运行特定包的测试
go test ./internal/domain/user/... -v

# 运行特定测试函数
go test ./internal/domain/user/... -run TestUser_Validate -v

# 运行性能测试
go test ./internal/domain/user/... -bench=. -benchmem
```

---

## 📖 测试示例

### 1. 单元测试（使用 testify）

```go
func TestUser_Validate(t *testing.T) {
    user := &User{
        Email: "test@example.com",
        Name:  "Test User",
    }

    err := user.Validate()
    assert.NoError(t, err)
}
```

### 2. 表格驱动测试

```go
func TestUser_Validate_TableDriven(t *testing.T) {
    tests := []struct {
        name    string
        user    *User
        wantErr bool
    }{
        {"valid user", &User{Email: "test@example.com"}, false},
        {"empty email", &User{Email: ""}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 3. Mock 测试

```go
func TestService_CreateUser(t *testing.T) {
    repo := &mocks.MockRepository[user.User]{}
    service := &Service{repo: repo}

    // 设置期望
    repo.On("Create", mock.Anything, mock.Anything).Return(nil)

    // 执行测试
    err := service.CreateUser(context.Background(), &user.User{
        Email: "test@example.com",
    })

    // 验证
    assert.NoError(t, err)
    repo.AssertExpectations(t)
}
```

### 4. 测试套件

```go
type UserServiceSuite struct {
    suite.Suite
    repo    *mocks.MockRepository[user.User]
    service *Service
}

func (s *UserServiceSuite) SetupTest() {
    s.repo = &mocks.MockRepository[user.User]{}
    s.service = &Service{repo: s.repo}
}

func (s *UserServiceSuite) TestCreateUser() {
    s.repo.On("Create", mock.Anything, mock.Anything).Return(nil)
    err := s.service.CreateUser(context.Background(), &user.User{})
    s.NoError(err)
}

func TestUserServiceSuite(t *testing.T) {
    suite.Run(t, new(UserServiceSuite))
}
```

---

## 🎯 测试覆盖率目标

### 当前状态

| 层次 | 当前覆盖率 | 目标覆盖率 | 状态 |
|------|-----------|-----------|------|
| Domain | < 30% | > 90% | ⚠️ 需提升 |
| Application | < 40% | > 85% | ⚠️ 需提升 |
| Infrastructure | < 20% | > 70% | ⚠️ 需提升 |
| Interfaces | < 30% | > 75% | ⚠️ 需提升 |
| **总体** | **< 30%** | **> 80%** | ⚠️ 需提升 |

### 提升计划

**Week 1**: Domain Layer
- ✅ Entity 测试
- ⏳ Specification 测试
- ⏳ Repository 接口测试

**Week 2**: Application Layer
- ✅ Service 测试框架
- ⏳ Command/Query 测试
- ⏳ Event 测试

**Week 3**: Infrastructure Layer
- ⏳ Repository 实现测试
- ⏳ Cache 测试
- ⏳ Messaging 测试

**Week 4**: Integration & E2E
- ⏳ API 集成测试
- ⏳ 数据库集成测试
- ⏳ E2E 流程测试

---

## 📚 最佳实践

### 1. 测试命名

```go
// 格式: Test<Function>_<Scenario>
func TestCreateUser_Success(t *testing.T) {}
func TestCreateUser_ValidationError(t *testing.T) {}
func TestCreateUser_RepositoryError(t *testing.T) {}
```

### 2. AAA 模式

```go
func TestExample(t *testing.T) {
    // Arrange - 准备
    user := &User{Email: "test@example.com"}

    // Act - 执行
    err := user.Validate()

    // Assert - 断言
    assert.NoError(t, err)
}
```

### 3. 表格驱动测试

- ✅ 覆盖多种场景
- ✅ 易于添加新用例
- ✅ 清晰的测试结构

### 4. Mock 使用原则

- ✅ Mock 外部依赖
- ✅ 不 Mock 被测试对象
- ✅ 验证 Mock 调用

---

## 🔧 Makefile 命令

```makefile
test: ## 运行所有测试
	go test ./... -v -race

test-unit: ## 运行单元测试
	go test ./internal/... -v -short

test-integration: ## 运行集成测试
	go test ./test/integration/... -v

test-e2e: ## 运行 E2E 测试
	go test ./test/e2e/... -v

coverage: ## 生成覆盖率报告
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

bench: ## 运行性能测试
	go test ./... -bench=. -benchmem
```

---

## 📊 CI/CD 集成

### GitHub Actions

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25.3'

      - name: Run tests
        run: make test

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
```

---

## 🎯 下一步

### 立即任务

1. **完善 Domain 测试**
   - Specification 测试
   - Entity 边界测试

2. **完善 Application 测试**
   - 所有 Service 方法
   - Command/Query 测试

3. **添加集成测试**
   - 数据库集成
   - API 集成

### 本周目标

- Domain Layer: > 90%
- Application Layer: > 85%
- 总体覆盖率: > 50%

---

**状态**: 🔄 建设中
**目标**: 覆盖率 > 80%
**优先级**: P0 (最高)
