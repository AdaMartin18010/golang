# LD-032-Go-Testing-Advanced

> **Dimension**: 02-Language-Design
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: Go 1.26 Testing Patterns
> **Size**: >20KB

---

## 1. Go测试基础

### 1.1 测试类型

```go
// 单元测试: *_test.go
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d, want 5", result)
    }
}

// 基准测试: Benchmark*
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}

// 示例测试: Example*
func ExampleAdd() {
    fmt.Println(Add(2, 3))
    // Output: 5
}

// Fuzz测试 (Go 1.18+): Fuzz*
func FuzzAdd(f *testing.F) {
    f.Add(2, 3)
    f.Fuzz(func(t *testing.T, a, b int) {
        result := Add(a, b)
        expected := a + b
        if result != expected {
            t.Errorf("Add(%d, %d) = %d, want %d", a, b, result, expected)
        }
    })
}
```

### 1.2 测试命令

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test -run TestAdd

# 详细输出
go test -v

# 代码覆盖率
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# 基准测试
go test -bench=.
go test -bench=BenchmarkAdd -benchmem

# 竞争检测
go test -race

# 模糊测试
go test -fuzz=FuzzAdd -fuzztime=30s
```

---

## 2. 表驱动测试

### 2.1 基本模式

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name     string
        a, b     float64
        expected float64
        wantErr  bool
    }{
        {"positive", 10, 2, 5, false},
        {"negative", -10, 2, -5, false},
        {"decimal", 7, 2, 3.5, false},
        {"by zero", 10, 0, 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.a, tt.b)

            if tt.wantErr {
                if err == nil {
                    t.Errorf("Divide(%v, %v) expected error", tt.a, tt.b)
                }
                return
            }

            if err != nil {
                t.Errorf("Divide(%v, %v) unexpected error: %v", tt.a, tt.b, err)
                return
            }

            if result != tt.expected {
                t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### 2.2 黄金文件测试

```go
// 用于复杂输出验证
func TestRenderer(t *testing.T) {
    tests := []string{"simple", "complex", "edge-cases"}

    for _, name := range tests {
        t.Run(name, func(t *testing.T) {
            input := readFile(t, filepath.Join("testdata", name+".input"))
            expected := readFile(t, filepath.Join("testdata", name+".golden"))

            result := Render(input)

            if *updateGolden {
                writeFile(t, filepath.Join("testdata", name+".golden"), result)
                return
            }

            if diff := cmp.Diff(string(expected), result); diff != "" {
                t.Errorf("Render() mismatch (-want +got):\n%s", diff)
            }
        })
    }
}

var updateGolden = flag.Bool("update-golden", false, "update golden files")
```

---

## 3. Mock与Stub

### 3.1 接口Mock

```go
// 定义接口
type Database interface {
    GetUser(id string) (*User, error)
    SaveUser(user *User) error
}

// 手动Mock
type MockDB struct {
    GetUserFunc  func(id string) (*User, error)
    SaveUserFunc func(user *User) error

    // 调用追踪
    GetUserCalls   []string
    SaveUserCalls  []*User
}

func (m *MockDB) GetUser(id string) (*User, error) {
    m.GetUserCalls = append(m.GetUserCalls, id)
    return m.GetUserFunc(id)
}

func (m *MockDB) SaveUser(user *User) error {
    m.SaveUserCalls = append(m.SaveUserCalls, user)
    return m.SaveUserFunc(user)
}

// 使用
func TestUserService(t *testing.T) {
    mockDB := &MockDB{
        GetUserFunc: func(id string) (*User, error) {
            if id == "123" {
                return &User{ID: "123", Name: "Alice"}, nil
            }
            return nil, ErrNotFound
        },
    }

    service := NewUserService(mockDB)
    user, err := service.GetUser("123")

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if user.Name != "Alice" {
        t.Errorf("expected Alice, got %s", user.Name)
    }

    // 验证调用
    if len(mockDB.GetUserCalls) != 1 {
        t.Errorf("expected 1 call, got %d", len(mockDB.GetUserCalls))
    }
}
```

### 3.2 使用Mockgen

```bash
# 安装
go install github.com/golang/mock/mockgen@latest

# 生成mock
mockgen -source=db.go -destination=db_mock.go -package=mocks
```

```go
// 使用生成的mock
func TestWithMockgen(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDB := mocks.NewMockDatabase(ctrl)

    // 设置期望
    mockDB.EXPECT().
        GetUser("123").
        Return(&User{ID: "123", Name: "Alice"}, nil).
        Times(1)

    mockDB.EXPECT().
        SaveUser(gomock.Any()).
        DoAndReturn(func(user *User) error {
            if user.Name == "" {
                return ErrInvalidUser
            }
            return nil
        })

    service := NewUserService(mockDB)
    service.ProcessUser("123")
}
```

### 3.3 HTTP Stub

```go
// 使用httptest
func TestAPIClient(t *testing.T) {
    // 创建测试服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/users/123" {
            t.Errorf("expected path /users/123, got %s", r.URL.Path)
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id":"123","name":"Alice"}`))
    }))
    defer server.Close()

    // 使用测试服务器URL
    client := NewAPIClient(server.URL)
    user, err := client.GetUser("123")

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if user.Name != "Alice" {
        t.Errorf("expected Alice, got %s", user.Name)
    }
}
```

---

## 4. 集成测试

### 4.1 测试容器

```go
// 使用testcontainers-go
import "github.com/testcontainers/testcontainers-go"

func TestWithPostgres(t *testing.T) {
    ctx := context.Background()

    // 启动PostgreSQL容器
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections"),
    }

    postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatal(err)
    }
    defer postgres.Terminate(ctx)

    // 获取连接信息
    host, _ := postgres.Host(ctx)
    port, _ := postgres.MappedPort(ctx, "5432")

    // 连接数据库
    dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb", host, port.Port())
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    // 运行迁移
    runMigrations(db)

    // 执行测试
    repo := NewUserRepository(db)
    user, err := repo.Create(&User{Name: "Alice"})

    if err != nil {
        t.Fatalf("create user: %v", err)
    }

    if user.ID == 0 {
        t.Error("expected user ID to be set")
    }
}
```

### 4.2 嵌入式测试

```go
// 使用嵌入式数据库
func TestWithSQLite(t *testing.T) {
    // 内存数据库
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    // 运行迁移
    runMigrations(db)

    // 测试...
}

// 使用embedded-postgres
import "github.com/fergusstrange/embedded-postgres"

func TestWithEmbeddedPostgres(t *testing.T) {
    postgres := embeddedpostgres.NewDatabase()
    if err := postgres.Start(); err != nil {
        t.Fatal(err)
    }
    defer postgres.Stop()

    db, _ := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres")
    // 测试...
}
```

---

## 5. 性能测试

### 5.1 基准测试进阶

```go
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}

// 子基准测试
func BenchmarkFibonacciMultiple(b *testing.B) {
    for _, n := range []int{10, 20, 30} {
        b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Fibonacci(n)
            }
        })
    }
}

// 比较基准测试
func BenchmarkFibonacciParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Fibonacci(20)
        }
    })
}

// 内存分析
func BenchmarkAllocate(b *testing.B) {
    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _ = make([]byte, 1024)
    }
}
```

### 5.2 Profiling

```go
// CPU Profiling
func TestWithCPUProfile(t *testing.T) {
    f, _ := os.Create("cpu.prof")
    defer f.Close()

    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // 执行测试逻辑
    HeavyOperation()
}

// Memory Profiling
func TestWithMemProfile(t *testing.T) {
    // 执行操作
    HeavyOperation()

    f, _ := os.Create("mem.prof")
    defer f.Close()

    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

---

## 6. 测试模式

### 6.1 Builder模式

```go
type UserBuilder struct {
    user *User
}

func NewUserBuilder() *UserBuilder {
    return &UserBuilder{
        user: &User{
            ID:        uuid.New().String(),
            CreatedAt: time.Now(),
        },
    }
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
    b.user.Name = name
    return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
    b.user.Email = email
    return b
}

func (b *UserBuilder) Build() *User {
    return b.user
}

// 测试中使用
func TestUserService(t *testing.T) {
    user := NewUserBuilder().
        WithName("Alice").
        WithEmail("alice@example.com").
        Build()

    service.CreateUser(user)
}
```

### 6.2 给定的-当-那么 (GWT)

```go
func TestOrderProcessing(t *testing.T) {
    // Given
    user := NewUserBuilder().WithBalance(100).Build()
    product := NewProductBuilder().WithPrice(50).WithStock(10).Build()

    // When
    order, err := service.CreateOrder(user, product, 1)

    // Then
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if order.Total != 50 {
        t.Errorf("expected total 50, got %d", order.Total)
    }

    if user.Balance != 50 {
        t.Errorf("expected balance 50, got %d", user.Balance)
    }
}
```

---

## 7. 持续测试

### 7.1 CI/CD集成

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

      redis:
        image: redis
        ports:
          - 6379:6379

    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-go@v4
        with:
          go-version: '1.26'

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
          REDIS_URL: localhost:6379

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

---

## 8. 参考文献

1. "The Go Programming Language" - Testing chapter
2. "Unit Testing Principles, Practices, and Patterns"
3. Google Testing Blog
4. Go Test Documentation

---

*Last Updated: 2026-04-03*
