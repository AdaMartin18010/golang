# LD-024: Go 测试高级模式 (Go Testing Advanced Patterns)

> **维度**: Language Design
> **级别**: S (18+ KB)
> **标签**: #testing #tdd #mock #benchmark #fuzzing #table-driven #testify
> **权威来源**:
>
> - [Testing Package](https://github.com/golang/go/tree/master/src/testing) - Go Authors
> - [Go Test Patterns](https://go.dev/doc/code#Testing) - Go Authors
> - [Advanced Testing in Go](https://speakerdeck.com/campoy/advanced-testing-in-go) - Francesc Campoy

---

## 1. 测试基础架构

### 1.1 测试类型

```go
// 单元测试
func TestSomething(t *testing.T) {
    // 测试单个函数/方法
}

// 基准测试
func BenchmarkSomething(b *testing.B) {
    // 性能测试
}

// 模糊测试 (Go 1.18+)
func FuzzSomething(f *testing.F) {
    // 模糊测试
}

// 示例测试
func ExampleSomething() {
    // 文档示例 + 测试
}

// Main 测试
func TestMain(m *testing.M) {
    // 测试入口，设置/清理
    os.Exit(m.Run())
}
```

### 1.2 测试生命周期

```go
func TestMain(m *testing.M) {
    // 1. 全局设置
    setup()

    // 2. 运行所有测试
    code := m.Run()

    // 3. 全局清理
    teardown()

    os.Exit(code)
}

// 测试级别设置/清理
func TestSomething(t *testing.T) {
    // 每个测试的设置
    setup := func() {
        t.Helper()
        // ...
    }

    // 每个测试的清理
    t.Cleanup(func() {
        // 测试结束后执行
    })

    setup()
    // 测试代码...
}
```

---

## 2. 表驱动测试

### 2.1 基础模式

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -1, 1, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d, want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### 2.2 高级表驱动测试

```go
func TestUserService(t *testing.T) {
    type args struct {
        ctx    context.Context
        userID string
    }

    type mock struct {
        dbUsers    map[string]*User
        cacheUsers map[string]*User
    }

    tests := []struct {
        name      string
        args      args
        mock      mock
        want      *User
        wantErr   bool
        errType   error
        setup     func(*mock)
        teardown  func(*mock)
    }{
        {
            name: "cache hit",
            args: args{context.Background(), "1"},
            mock: mock{
                cacheUsers: map[string]*User{
                    "1": {ID: "1", Name: "Alice"},
                },
            },
            want: &User{ID: "1", Name: "Alice"},
        },
        {
            name: "cache miss db hit",
            args: args{context.Background(), "2"},
            mock: mock{
                dbUsers: map[string]*User{
                    "2": {ID: "2", Name: "Bob"},
                },
            },
            want: &User{ID: "2", Name: "Bob"},
            setup: func(m *mock) {
                // 验证缓存被填充
            },
        },
        {
            name:    "not found",
            args:    args{context.Background(), "3"},
            mock:    mock{},
            wantErr: true,
            errType: ErrNotFound,
        },
        {
            name:    "timeout",
            args:    args{ctx: timeoutCtx(0), userID: "1"},
            mock:    mock{},
            wantErr: true,
            errType: context.DeadlineExceeded,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 创建 mock 依赖
            cache := NewMockCache(tt.mock.cacheUsers)
            db := NewMockDB(tt.mock.dbUsers)
            svc := NewUserService(cache, db)

            if tt.setup != nil {
                tt.setup(&tt.mock)
            }
            if tt.teardown != nil {
                defer tt.teardown(&tt.mock)
            }

            got, err := svc.GetUser(tt.args.ctx, tt.args.userID)

            if (err != nil) != tt.wantErr {
                t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if tt.wantErr && tt.errType != nil && !errors.Is(err, tt.errType) {
                t.Errorf("GetUser() error type = %v, want %v", err, tt.errType)
                return
            }

            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetUser() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

## 3. Mock 与 Stub

### 3.1 接口 Mock

```go
// 定义接口
type Database interface {
    Get(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}

// 手动 Mock
type MockDB struct {
    GetFunc    func(ctx context.Context, id string) (*User, error)
    SaveFunc   func(ctx context.Context, user *User) error
    DeleteFunc func(ctx context.Context, id string) error

    // 调用追踪
    GetCalls    []GetCall
    SaveCalls   []SaveCall
    DeleteCalls []DeleteCall
}

type GetCall struct {
    Ctx context.Context
    ID  string
}

func (m *MockDB) Get(ctx context.Context, id string) (*User, error) {
    m.GetCalls = append(m.GetCalls, GetCall{Ctx: ctx, ID: id})
    if m.GetFunc != nil {
        return m.GetFunc(ctx, id)
    }
    return nil, errors.New("Get not implemented")
}

// 使用 mock
func TestUserService(t *testing.T) {
    mockDB := &MockDB{
        GetFunc: func(ctx context.Context, id string) (*User, error) {
            if id == "123" {
                return &User{ID: "123", Name: "Test"}, nil
            }
            return nil, ErrNotFound
        },
    }

    svc := NewUserService(mockDB)
    user, err := svc.GetUser(context.Background(), "123")

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    // 验证调用
    if len(mockDB.GetCalls) != 1 {
        t.Errorf("expected 1 Get call, got %d", len(mockDB.GetCalls))
    }

    if mockDB.GetCalls[0].ID != "123" {
        t.Errorf("expected ID 123, got %s", mockDB.GetCalls[0].ID)
    }
}
```

### 3.2 Mockgen 代码生成

```go
//go:generate mockgen -source=db.go -destination=mock_db.go -package=mocks

type Database interface {
    Get(ctx context.Context, id string) (*User, error)
}

// 生成代码使用
func TestWithMockgen(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDB := mocks.NewMockDatabase(ctrl)

    // 设置期望
    mockDB.EXPECT().
        Get(gomock.Any(), "123").
        Return(&User{ID: "123", Name: "Test"}, nil).
        Times(1)

    svc := NewUserService(mockDB)
    user, err := svc.GetUser(context.Background(), "123")

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if user.Name != "Test" {
        t.Errorf("expected name Test, got %s", user.Name)
    }
}
```

---

## 4. 基准测试

### 4.1 基础基准测试

```go
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(20)
    }
}

// 带并行测试
func BenchmarkFibonacciParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Fibonacci(20)
        }
    })
}

// 子基准测试
func BenchmarkFibonacciSizes(b *testing.B) {
    for _, n := range []int{10, 20, 30} {
        b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Fibonacci(n)
            }
        })
    }
}
```

### 4.2 内存分配分析

```go
func BenchmarkStringConcat(b *testing.B) {
    b.ReportAllocs() // 报告内存分配

    b.Run("naive", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var s string
            for j := 0; j < 100; j++ {
                s += "x" // 大量分配
            }
        }
    })

    b.Run("builder", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var builder strings.Builder
            builder.Grow(100)
            for j := 0; j < 100; j++ {
                builder.WriteString("x")
            }
            _ = builder.String()
        }
    })
}

// 结果:
// BenchmarkStringConcat/naive-8     100000    12345 ns/op    123456 B/op    100 allocs/op
// BenchmarkStringConcat/builder-8   5000000   234 ns/op      256 B/op      1 allocs/op
```

### 4.3 比较基准测试

```go
func BenchmarkSort(b *testing.B) {
    sizes := []int{100, 1000, 10000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            data := make([]int, size)

            b.Run("stdlib", func(b *testing.B) {
                for i := 0; i < b.N; i++ {
                    b.StopTimer()
                    rand.Shuffle(len(data), func(i, j int) {
                        data[i], data[j] = data[j], data[i]
                    })
                    b.StartTimer()

                    sort.Ints(data)
                }
            })

            b.Run("quicksort", func(b *testing.B) {
                for i := 0; i < b.N; i++ {
                    b.StopTimer()
                    rand.Shuffle(len(data), func(i, j int) {
                        data[i], data[j] = data[j], data[i]
                    })
                    b.StartTimer()

                    quickSort(data)
                }
            })
        })
    }
}
```

---

## 5. 模糊测试

### 5.1 基础模糊测试

```go
// Go 1.18+ 模糊测试
func FuzzParse(f *testing.F) {
    // 种子语料
    f.Add("hello world")
    f.Add("12345")
    f.Add("")

    f.Fuzz(func(t *testing.T, input string) {
        result, err := Parse(input)
        if err != nil {
            // 某些输入可能产生错误，这是正常的
            return
        }

        // 不变性检查
        if result == nil {
            t.Error("result should not be nil")
        }
    })
}

// 结构化模糊测试
func FuzzUser(f *testing.F) {
    f.Add("Alice", 25, "alice@example.com")
    f.Add("Bob", -1, "invalid")

    f.Fuzz(func(t *testing.T, name string, age int, email string) {
        user := &User{
            Name:  name,
            Age:   age,
            Email: email,
        }

        // 验证验证逻辑
        if err := user.Validate(); err != nil {
            return
        }

        // 序列化反序列化不变性
        data, _ := json.Marshal(user)
        var decoded User
        if err := json.Unmarshal(data, &decoded); err != nil {
            t.Fatalf("failed to unmarshal: %v", err)
        }

        if !reflect.DeepEqual(user, &decoded) {
            t.Errorf("roundtrip failed: %+v != %+v", user, decoded)
        }
    })
}
```

---

## 6. 测试辅助工具

### 6.1 testify 使用

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
)

// assert vs require
func TestWithTestify(t *testing.T) {
    result, err := DoSomething()

    // assert: 失败继续
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "expected", result.Value)

    // require: 失败立即停止
    require.NoError(t, err)
    require.NotNil(t, result)

    // 后续测试...
}

// 测试套件
type UserSuite struct {
    suite.Suite
    db     *sql.DB
    svc    *UserService
}

func (s *UserSuite) SetupSuite() {
    // 整个套件执行一次
    s.db = setupTestDB()
    s.svc = NewUserService(s.db)
}

func (s *UserSuite) TearDownSuite() {
    s.db.Close()
}

func (s *UserSuite) SetupTest() {
    // 每个测试前执行
    cleanDB(s.db)
}

func (s *UserSuite) TestCreate() {
    user, err := s.svc.Create("Alice")
    s.NoError(err)
    s.NotNil(user)
    s.Equal("Alice", user.Name)
}

func (s *UserSuite) TestGet() {
    s.svc.Create("Bob")

    user, err := s.svc.Get("Bob")
    s.NoError(err)
    s.Equal("Bob", user.Name)
}

func TestUserSuite(t *testing.T) {
    suite.Run(t, new(UserSuite))
}
```

### 6.2 Golden 文件测试

```go
// 用于大型输出（如 JSON、HTML）的测试
func TestGenerateHTML(t *testing.T) {
    data := LoadTestData()
    html := GenerateHTML(data)

    // 更新 golden 文件: go test -update
    if *update {
        os.WriteFile("testdata/output.golden", html, 0644)
        return
    }

    golden, err := os.ReadFile("testdata/output.golden")
    if err != nil {
        t.Fatalf("failed to read golden file: %v", err)
    }

    if !bytes.Equal(html, golden) {
        t.Errorf("output mismatch\nexpected:\n%s\nactual:\n%s", golden, html)
    }
}
```

---

## 7. 并发测试

### 7.1 竞态检测

```go
// 使用 -race 标志检测竞态
// go test -race ./...

func TestConcurrentAccess(t *testing.T) {
    cache := NewCache()

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            key := fmt.Sprintf("key%d", n%10)
            cache.Set(key, n)
            cache.Get(key)
        }(i)
    }
    wg.Wait()
}
```

### 7.2 超时测试

```go
func TestWithTimeout(t *testing.T) {
    done := make(chan bool)

    go func() {
        // 长时间操作
        time.Sleep(5 * time.Second)
        done <- true
    }()

    select {
    case <-done:
        // 成功
    case <-time.After(1 * time.Second):
        t.Fatal("operation timed out")
    }
}

// 使用 context
func TestWithContextTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    err := LongRunningOperation(ctx)
    if err != context.DeadlineExceeded {
        t.Fatalf("expected timeout, got %v", err)
    }
}
```

---

## 8. 视觉表征

### 8.1 测试金字塔

```
            /\
           /  \
          / E2E \          少量，慢，贵
         /________\
        /          \
       / Integration \     中等，API/DB 测试
      /________________\
     /                  \
    /     Unit Tests      \   大量，快，便宜
   /________________________\
```

### 8.2 测试组织结构

```
project/
├── cmd/
│   └── server/
├── internal/
│   └── service/
│       ├── user.go
│       └── user_test.go      # 同包白盒测试
├── pkg/
│   └── utils/
│       ├── helper.go
│       └── helper_test.go    # 同包测试
├── test/
│   ├── integration/          # 集成测试
│   │   └── api_test.go
│   ├── fixtures/             # 测试数据
│   │   └── users.json
│   └── mocks/                # 生成的 mocks
│       └── mock_db.go
└── go.mod
```

### 8.3 测试策略决策树

```
测试什么?
│
├── 纯函数/算法?
│   └── 表驱动单元测试
│
├── 外部依赖?
│   ├── HTTP API? → 使用 httptest
│   ├── 数据库? → 使用 testcontainers/sqlite
│   └── 缓存? → 使用 miniredis
│
├── 并发代码?
│   ├── 使用 -race 检测
│   └── 使用 sync.WaitGroup 协调
│
└── 性能关键?
    ├── 基准测试比较
    ├── 内存分配分析
    └── 使用 pprof 分析
```

---

**质量评级**: S (18KB)
**完成日期**: 2026-04-02
