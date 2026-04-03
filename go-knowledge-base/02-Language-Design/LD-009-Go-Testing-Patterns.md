# LD-009: Go 测试模式 (Go Testing Patterns)

> **维度**: Language Design
> **级别**: S (35+ KB)
> **标签**: #testing #patterns #table-driven #mock #benchmark
> **权威来源**:
>
> - [Testing in Go](https://go.dev/doc/tutorial/add-a-test) - Go Authors
> - [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests) - Go Wiki
> - [Advanced Testing in Go](https://speakerdeck.com/campoy/advanced-testing-in-go) - Francesc Campoy

---

## 1. 测试基础

### 1.1 测试函数签名

```go
// 单元测试
func TestXxx(t *testing.T)

// 基准测试
func BenchmarkXxx(b *testing.B)

// 模糊测试 (Go 1.18+)
func FuzzXxx(f *testing.F)

// 示例测试
func ExampleXxx()
```

### 1.2 测试结构

```
myproject/
├── foo.go
├── foo_test.go      // 白盒测试 (同包)
└── foo_blackbox_test.go  // 黑盒测试 (package_foo_test)
```

---

## 2. 表驱动测试

### 2.1 基本模式

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
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, got, tt.expected)
            }
        })
    }
}
```

### 2.2 带错误测试

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name     string
        a, b     float64
        expected float64
        wantErr  bool
    }{
        {"normal", 10, 2, 5, false},
        {"negative", -10, 2, -5, false},
        {"zero divisor", 10, 0, 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)
            if (err != nil) != tt.wantErr {
                t.Errorf("Divide() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && got != tt.expected {
                t.Errorf("Divide() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

---

## 3. Mock 和 Stub

### 3.1 接口 Mock

```go
// 定义接口
type DataStore interface {
    Get(id int) (*User, error)
    Save(u *User) error
}

// Mock 实现
type MockStore struct {
    GetFunc  func(id int) (*User, error)
    SaveFunc func(u *User) error
}

func (m *MockStore) Get(id int) (*User, error) {
    return m.GetFunc(id)
}

func (m *MockStore) Save(u *User) error {
    return m.SaveFunc(u)
}

// 测试使用
func TestServiceGetUser(t *testing.T) {
    mock := &MockStore{
        GetFunc: func(id int) (*User, error) {
            if id == 1 {
                return &User{ID: 1, Name: "Alice"}, nil
            }
            return nil, errors.New("not found")
        },
    }

    svc := NewService(mock)
    user, err := svc.GetUser(1)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("expected Alice, got %s", user.Name)
    }
}
```

### 3.2 HTTP 测试

```go
func TestHandler(t *testing.T) {
    tests := []struct {
        name       string
        method     string
        path       string
        wantStatus int
        wantBody   string
    }{
        {"get user", "GET", "/users/1", 200, `{"id":1}`},
        {"not found", "GET", "/users/999", 404, ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(tt.method, tt.path, nil)
            rec := httptest.NewRecorder()

            handler.ServeHTTP(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
            }

            if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
                t.Errorf("body = %s, want %s", rec.Body.String(), tt.wantBody)
            }
        })
    }
}
```

---

## 4. 基准测试

### 4.1 基本基准测试

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(1, 2)
    }
}
```

### 4.2 子基准测试

```go
func BenchmarkFib(b *testing.B) {
    benchmarks := []struct {
        name string
        n    int
    }{
        {"10", 10},
        {"20", 20},
        {"30", 30},
    }

    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Fib(bm.n)
            }
        })
    }
}
```

### 4.3 内存分配分析

```go
func BenchmarkAlloc(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        _ = make([]byte, 1024)
    }
}
```

---

## 5. 运行时行为分析

### 5.1 测试执行流程

```
┌─────────────────────────────────────────────────────────────────┐
│                    Test Execution Flow                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  go test                                                        │
│       │                                                          │
│       ▼                                                          │
│  ┌─────────────────────────────────────────────┐                │
│  │ 1. 编译测试二进制                            │                │
│  │    go test -c                               │                │
│  │    生成 .test 文件                          │                │
│  └─────────────────────────┬───────────────────┘                │
│                            │                                     │
│                            ▼                                     │
│  ┌─────────────────────────────────────────────┐                │
│  │ 2. 运行测试                                  │                │
│  │    - 执行 TestXxx 函数                      │                │
│  │    - 并行执行 (t.Parallel())                │                │
│  │    - 每个测试隔离                            │                │
│  └─────────────────────────┬───────────────────┘                │
│                            │                                     │
│                            ▼                                     │
│  ┌─────────────────────────────────────────────┐                │
│  │ 3. 收集结果                                  │                │
│  │    - 通过/失败计数                          │                │
│  │    - 覆盖率数据                             │                │
│  │    - 基准测试结果                           │                │
│  └─────────────────────────┬───────────────────┘                │
│                            │                                     │
│                            ▼                                     │
│  ┌─────────────────────────────────────────────┐                │
│  │ 4. 输出报告                                  │                │
│  │    - 测试摘要                               │                │
│  │    - 失败详情                               │                │
│  │    - 覆盖率报告                             │                │
│  └─────────────────────────────────────────────┘                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 测试并行执行模型

```
┌─────────────────────────────────────────────────────────────────┐
│                    Parallel Test Execution                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  默认模式:                                                       │
│  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐                               │
│  │ T1  │→│ T2  │→│ T3  │→│ T4  │  顺序执行                     │
│  └─────┘ └─────┘ └─────┘ └─────┘                               │
│                                                                  │
│  并行模式 (t.Parallel()):                                        │
│  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐                               │
│  │ T1  │ │ T2  │ │ T3  │ │ T4  │  并行执行                     │
│  └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘                               │
│     └───────┴───────┴───────┘                                   │
│            GOMAXPROCS goroutines                                │
│                                                                  │
│  控制并行度:                                                     │
│  go test -parallel=4                                            │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 5.3 基准测试统计

```
┌─────────────────────────────────────────────────────────────────┐
│                    Benchmark Statistics                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  BenchmarkName-8         1000000        1050 ns/op      48 B/op │
│  │                        │              │            │         │
│  │                        │              │            │         │
│  │                        │              │            │         │
│  │                        │              │            每次操作的 │
│  │                        │              │            内存分配    │
│  │                        │              │                       │
│  │                        │              每次操作的              │
│  │                        │              纳秒数                   │
│  │                        │                                     │
│  │                        迭代次数 (b.N)                         │
│  │                                                             │
│  测试名称，-8 表示 GOMAXPROCS=8                                │
│                                                                  │
│  统计方法:                                                       │
│  1. 自动调整 b.N 直到测试运行时间 >= 1s                         │
│  2. 计算平均耗时和内存分配                                       │
│  3. 多次运行取稳定结果                                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 6. 内存与性能特性

### 6.1 测试性能特征

| 操作 | 典型耗时 | 说明 |
|------|----------|------|
| 简单单元测试 | ~1ms | 无 I/O |
| 数据库测试 | ~10-100ms | 含连接开销 |
| HTTP 测试 | ~1-10ms | httptest.Server |
| 基准测试迭代 | ~1μs | 取决于被测函数 |

### 6.2 基准测试精度

```go
// 控制测量精度
func BenchmarkPrecise(b *testing.B) {
    // 重置计时器（排除初始化）
    b.ResetTimer()

    // 停止计时器
    b.StopTimer()
    setup()
    b.StartTimer()

    for i := 0; i < b.N; i++ {
        Operation()
    }
}
```

---

## 7. 多元表征

### 7.1 测试金字塔

```
                    ┌─────────┐
                    │   E2E   │  端到端测试 (少)
                    │  Tests  │
                   ┌┴─────────┴┐
                   │ Integration│ 集成测试 (中)
                   │   Tests   │
                  ┌┴───────────┴┐
                  │    Unit     │ 单元测试 (多)
                  │   Tests     │
                  └─────────────┘
```

### 7.2 测试决策树

```
写什么类型的测试?
│
├── 纯逻辑函数?
│   └── 单元测试 (table-driven)
│
├── 依赖外部服务?
│   ├── 使用 Mock/Stub
│   └── 集成测试
│
├── 性能关键代码?
│   └── 基准测试
│
├── 输入验证?
│   └── 模糊测试 (Go 1.18+)
│
└── 完整用户流程?
    └── 端到端测试
```

### 7.3 表驱动测试结构

```
┌─────────────────────────────────────────────────────────────────┐
│                    Table-Driven Test Structure                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  func TestXxx(t *testing.T) {                                   │
│       │                                                          │
│       ▼                                                          │
│  ┌─────────────────────────────────────────────┐                │
│  │ 1. 定义测试用例表                            │                │
│  │    tests := []struct {                      │                │
│  │        name     string                      │                │
│  │        input    InputType                   │                │
│  │        want     OutputType                  │                │
│  │        wantErr  bool                        │                │
│  │    }{...}                                   │                │
│  └─────────────────────────────────────────────┘                │
│       │                                                          │
│       ▼                                                          │
│  ┌─────────────────────────────────────────────┐                │
│  │ 2. 遍历执行                                  │                │
│  │    for _, tt := range tests {               │                │
│  │        t.Run(tt.name, func(t *testing.T) {  │                │
│  │            got, err := Xxx(tt.input)        │                │
│  │            // 断言...                        │                │
│  │        })                                   │                │
│  │    }                                        │                │
│  └─────────────────────────────────────────────┘                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 8. 完整代码示例

### 8.1 完整表驱动测试

```go
package calculator

import (
    "errors"
    "testing"
)

var (
    ErrDivideByZero = errors.New("divide by zero")
)

func Add(a, b int) int {
    return a + b
}

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, ErrDivideByZero
    }
    return a / b, nil
}

func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {
            name: "positive numbers",
            a:    1,
            b:    2,
            want: 3,
        },
        {
            name: "negative numbers",
            a:    -1,
            b:    -2,
            want: -3,
        },
        {
            name: "zero",
            a:    0,
            b:    0,
            want: 0,
        },
        {
            name: "mixed signs",
            a:    -1,
            b:    1,
            want: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, got, tt.want)
            }
        })
    }
}

func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        a       float64
        b       float64
        want    float64
        wantErr error
    }{
        {
            name:    "normal division",
            a:       10,
            b:       2,
            want:    5,
            wantErr: nil,
        },
        {
            name:    "divide by zero",
            a:       10,
            b:       0,
            want:    0,
            wantErr: ErrDivideByZero,
        },
        {
            name:    "negative divisor",
            a:       10,
            b:       -2,
            want:    -5,
            wantErr: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("Divide() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if err == nil && got != tt.want {
                t.Errorf("Divide() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 8.2 Mock 测试完整示例

```go
package service

import (
    "errors"
    "testing"
)

// 定义接口
type Repository interface {
    GetUser(id int) (*User, error)
    SaveUser(user *User) error
}

type User struct {
    ID   int
    Name string
}

// 实现
type UserService struct {
    repo Repository
}

func NewUserService(repo Repository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUserName(id int) (string, error) {
    user, err := s.repo.GetUser(id)
    if err != nil {
        return "", err
    }
    return user.Name, nil
}

// Mock 实现
type MockRepository struct {
    GetUserFunc  func(id int) (*User, error)
    SaveUserFunc func(user *User) error
}

func (m *MockRepository) GetUser(id int) (*User, error) {
    return m.GetUserFunc(id)
}

func (m *MockRepository) SaveUser(user *User) error {
    return m.SaveUserFunc(user)
}

// 测试
func TestUserService_GetUserName(t *testing.T) {
    tests := []struct {
        name    string
        userID  int
        mock    func() *MockRepository
        want    string
        wantErr bool
    }{
        {
            name:   "user found",
            userID: 1,
            mock: func() *MockRepository {
                return &MockRepository{
                    GetUserFunc: func(id int) (*User, error) {
                        if id == 1 {
                            return &User{ID: 1, Name: "Alice"}, nil
                        }
                        return nil, errors.New("not found")
                    },
                }
            },
            want:    "Alice",
            wantErr: false,
        },
        {
            name:   "user not found",
            userID: 999,
            mock: func() *MockRepository {
                return &MockRepository{
                    GetUserFunc: func(id int) (*User, error) {
                        return nil, errors.New("not found")
                    },
                }
            },
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := tt.mock()
            svc := NewUserService(mock)
            got, err := svc.GetUserName(tt.userID)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetUserName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("GetUserName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 8.3 HTTP 测试完整示例

```go
package handler

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

type UserHandler struct {
    service UserService
}

type UserService interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
}

type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/users/")
    user, err := h.service.GetUser(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(user)
}

func TestUserHandler_GetUser(t *testing.T) {
    tests := []struct {
        name       string
        path       string
        mockFunc   func(id string) (*User, error)
        wantStatus int
        wantBody   string
    }{
        {
            name: "user found",
            path: "/users/1",
            mockFunc: func(id string) (*User, error) {
                return &User{ID: "1", Name: "Alice"}, nil
            },
            wantStatus: http.StatusOK,
            wantBody:   `{"id":"1","name":"Alice"}`,
        },
        {
            name: "user not found",
            path: "/users/999",
            mockFunc: func(id string) (*User, error) {
                return nil, errors.New("not found")
            },
            wantStatus: http.StatusNotFound,
            wantBody:   "not found\n",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 创建 mock service
            mock := &MockUserService{
                GetUserFunc: tt.mockFunc,
            }
            handler := &UserHandler{service: mock}

            // 创建请求
            req := httptest.NewRequest(http.MethodGet, tt.path, nil)
            rec := httptest.NewRecorder()

            // 执行请求
            handler.GetUser(rec, req)

            // 验证结果
            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
            }
            if body := rec.Body.String(); body != tt.wantBody {
                t.Errorf("body = %s, want %s", body, tt.wantBody)
            }
        })
    }
}
```

### 8.4 基准测试完整示例

```go
package sort

import (
    "math/rand"
    "sort"
    "testing"
)

// BubbleSort 冒泡排序 (慢)
func BubbleSort(data []int) {
    n := len(data)
    for i := 0; i < n; i++ {
        for j := 0; j < n-i-1; j++ {
            if data[j] > data[j+1] {
                data[j], data[j+1] = data[j+1], data[j]
            }
        }
    }
}

// QuickSort 快速排序 (快)
func QuickSort(data []int) {
    sort.Ints(data)
}

// 基准测试
func BenchmarkBubbleSort_10(b *testing.B)   { benchmarkSort(b, 10, BubbleSort) }
func BenchmarkBubbleSort_100(b *testing.B)  { benchmarkSort(b, 100, BubbleSort) }
func BenchmarkBubbleSort_1000(b *testing.B) { benchmarkSort(b, 1000, BubbleSort) }

func BenchmarkQuickSort_10(b *testing.B)   { benchmarkSort(b, 10, QuickSort) }
func BenchmarkQuickSort_100(b *testing.B)  { benchmarkSort(b, 100, QuickSort) }
func BenchmarkQuickSort_1000(b *testing.B) { benchmarkSort(b, 1000, QuickSort) }

func benchmarkSort(b *testing.B, size int, sortFn func([]int)) {
    b.ReportAllocs()
    data := make([]int, size)

    for i := 0; i < b.N; i++ {
        // 每次重新生成随机数据
        for j := range data {
            data[j] = rand.Intn(size)
        }
        sortFn(data)
    }
}

// 比较基准测试
func BenchmarkSortComparison(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("bubble-%d", size), func(b *testing.B) {
            benchmarkSort(b, size, BubbleSort)
        })
        b.Run(fmt.Sprintf("quick-%d", size), func(b *testing.B) {
            benchmarkSort(b, size, QuickSort)
        })
    }
}
```

---

## 9. 测试工具

### 9.1 testify

```go
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
    result := DoSomething()

    assert.Equal(t, expected, result)
    assert.NoError(t, err)
    assert.NotNil(t, obj)
}
```

### 9.2 golden 文件

```go
func TestOutput(t *testing.T) {
    got := GenerateOutput()

    if *update {
        os.WriteFile("testdata/output.golden", got, 0644)
    }

    want, _ := os.ReadFile("testdata/output.golden")
    if !bytes.Equal(got, want) {
        t.Errorf("output mismatch")
    }
}
```

---

## 10. 最佳实践与反模式

### 10.1 ✅ 最佳实践

```go
// 1. 使用表驱动测试
func TestFunc(t *testing.T) {
    tests := []struct{...}{...}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}

// 2. 使用 t.Parallel() 加速测试
func TestParallel(t *testing.T) {
    t.Parallel()
    // 测试逻辑
}

// 3. 使用 t.Helper() 辅助函数
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

// 4. 使用 t.Cleanup() 清理资源
func TestWithCleanup(t *testing.T) {
    f, _ := os.CreateTemp("", "test")
    t.Cleanup(func() { os.Remove(f.Name()) })
    // 使用 f...
}

// 5. 命名规范: Test{FunctionName}_{Scenario}_{Expected}
func TestUserService_GetUser_Success(t *testing.T) {}
func TestUserService_GetUser_NotFound(t *testing.T) {}
```

### 10.2 ❌ 反模式

```go
// 1. 测试之间相互依赖
func TestA(t *testing.T) {
    globalVar = 1  // 修改全局状态
}

func TestB(t *testing.T) {
    if globalVar != 1 {  // 依赖 TestA
        t.Fail()
    }
}

// 2. 测试中使用 time.Sleep
func TestWithSleep(t *testing.T) {
    go doAsync()
    time.Sleep(time.Second)  // 不可靠！
}

// 3. 忽略错误返回值
func TestIgnoreError(t *testing.T) {
    result, _ := DoSomething()  // 忽略错误
    _ = result
}

// 4. 测试代码重复
func TestCase1(t *testing.T) {
    setup()
    // 测试1...
    cleanup()
}

func TestCase2(t *testing.T) {
    setup()  // 重复！
    // 测试2...
    cleanup()  // 重复！
}
```

---

## 11. 关系网络

```
Go Testing
├── Test Types
│   ├── Unit tests
│   ├── Integration tests
│   ├── Benchmark tests
│   ├── Fuzz tests (Go 1.18+)
│   └── Example tests
├── Patterns
│   ├── Table-driven tests
│   ├── Mock/Stub
│   ├── Golden files
│   └── Subtests (t.Run)
├── Tools
│   ├── Standard library (testing)
│   ├── Testify
│   ├── GoMock
│   └── httptest
└── Best Practices
    ├── Parallel execution
    ├── Test isolation
    ├── Coverage
    └── Performance
```

---

## 12. 参考文献

1. **Go Authors.** Testing Package Documentation.
2. **Campoy, F.** Advanced Testing in Go.
3. **Martin, R. C.** Clean Code: Unit Tests.

---

**质量评级**: S (35KB)
**完成日期**: 2026-04-02
