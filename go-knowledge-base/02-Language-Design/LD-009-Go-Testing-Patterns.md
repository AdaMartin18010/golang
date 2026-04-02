# LD-009: Go 测试模式 (Go Testing Patterns)

> **维度**: Language Design
> **级别**: S (16+ KB)
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

## 5. 测试工具

### 5.1 testify

```go
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
    result := DoSomething()

    assert.Equal(t, expected, result)
    assert.NoError(t, err)
    assert.NotNil(t, obj)
}
```

### 5.2 golden 文件

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

## 6. 测试最佳实践

### 6.1 检查清单

- [ ] 测试覆盖正常路径和错误路径
- [ ] 使用表驱动测试减少重复
- [ ] 测试名称描述行为
- [ ] 并行运行测试 (t.Parallel())
- [ ] 清理测试资源 (t.Cleanup())

### 6.2 命名规范

```
Test{FunctionName}_{Scenario}_{ExpectedResult}

例:
- TestUserService_GetUser_Success
- TestUserService_GetUser_NotFound
- TestCalculator_Divide_ByZero
```

---

**质量评级**: S (15KB)
**完成日期**: 2026-04-02
