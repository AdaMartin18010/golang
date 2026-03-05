# Go测试技术完全指南

> 从单元测试到集成测试的全面实践

---

## 一、测试基础哲学

### 1.1 为什么需要测试

```text
测试的价值：
────────────────────────────────────────

1. 验证正确性：
   - 代码是否按预期工作
   - 边界条件是否处理正确
   - 异常情况是否有处理

2. 防止回归：
   - 新功能不破坏旧功能
   - 重构时有安全网
   - CI/CD自动化验证

3. 文档作用：
   - 展示API使用方式
   - 说明预期行为
   - 作为使用示例

4. 设计反馈：
   - 难测试的代码往往设计有问题
   - 推动更好的代码结构
   - 促进关注点分离

测试金字塔：
────────────────────────────────────────

        /
       / \     E2E测试 (少量)
      /   \    - 用户场景
     /─────\
    /       \   集成测试 (中等)
   /         \  - 组件协作
  /───────────\
 /             \ 单元测试 (大量)
/               \ - 函数/方法

比例建议：
- 单元测试：70%
- 集成测试：20%
- E2E测试：10%

不写测试的代价：
────────────────────────────────────────

场景：紧急修复生产bug

有测试：
1. 复现bug → 编写测试 → 修复 → 验证
2. 10分钟完成
3. 高置信度不会引入新问题

无测试：
1. 尝试复现bug
2. 手动验证修复
3. 担心破坏其他功能
4. 花费数小时
5. 低置信度，可能引入回归

长期成本：
- 手动测试时间 >> 写测试时间
- bug修复成本 >> 预防成本
- 技术债务累积
```

### 1.2 好的测试特征

```text
FIRST原则：
────────────────────────────────────────

F - Fast (快速)：
- 测试应该在毫秒级完成
- 慢测试不会被频繁运行
- 使用mock避免真实依赖

I - Isolated (独立)：
- 测试之间不应相互依赖
- 可以单独运行任意测试
- 执行顺序不影响结果

R - Repeatable (可重复)：
- 每次运行结果相同
- 不依赖外部状态
- 不依赖随机数据

S - Self-validating (自验证)：
- 测试结果应该是布尔值
- 无需人工检查输出
- CI/CD可以自动判断

T - Timely (及时)：
- 测试应该与代码一起编写
- 理想情况：先写测试后写代码
- 最迟：代码提交前完成

AAA模式：
────────────────────────────────────────

Arrange (准备)：
- 设置测试环境
- 创建测试数据
- 初始化依赖

Act (执行)：
- 调用被测代码
- 通常只有一行

Assert (验证)：
- 验证结果
- 检查副作用
- 确认预期行为

示例：
func TestAdd(t *testing.T) {
    // Arrange
    a, b := 2, 3
    expected := 5

    // Act
    result := Add(a, b)

    // Assert
    if result != expected {
        t.Errorf("Add(%d, %d) = %d, want %d", a, b, result, expected)
    }
}
```

---

## 二、表格驱动测试

### 2.1 基础表格测试

```text
什么是表格驱动测试：
────────────────────────────────────────

将多个测试用例组织在一个表格中，
循环执行每个用例。

优势：
- 易于添加新用例
- 结构清晰
- 避免重复代码

基础示例：
────────────────────────────────────────

func Add(a, b int) int {
    return a + b
}

func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
        {"mixed", -2, 3, 1},
        {"zero", 0, 5, 5},
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

运行结果：
$ go test -v
=== RUN   TestAdd
=== RUN   TestAdd/positive
=== RUN   TestAdd/negative
=== RUN   TestAdd/mixed
=== RUN   TestAdd/zero
--- PASS: TestAdd (0.00s)
    --- PASS: TestAdd/positive (0.00s)
    --- PASS: TestAdd/negative (0.00s)
    --- PASS: TestAdd/mixed (0.00s)
    --- PASS: TestAdd/zero (0.00s)
```

### 2.2 高级表格测试

```text
包含错误测试：
────────────────────────────────────────

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func TestDivide(t *testing.T) {
    tests := []struct {
        name        string
        a, b        float64
        expected    float64
        expectError bool
    }{
        {"normal", 10, 2, 5, false},
        {"negative", -10, 2, -5, false},
        {"by_zero", 10, 0, 0, true},
        {"decimal", 7, 2, 3.5, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.a, tt.b)

            if tt.expectError {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }

            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }

            if result != tt.expected {
                t.Errorf("Divide(%f, %f) = %f, want %f",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

使用辅助函数：
────────────────────────────────────────

// 断言助手
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %+v, want %+v", got, want)
    }
}

func assertError(t *testing.T, got error, wantErr bool) {
    t.Helper()
    if (got != nil) != wantErr {
        t.Errorf("error = %v, wantErr %v", got, wantErr)
    }
}

func TestWithHelpers(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    int
        wantErr bool
    }{
        {"valid", "42", 42, false},
        {"invalid", "abc", 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := strconv.Atoi(tt.input)
            assertEqual(t, got, tt.want)
            assertError(t, err, tt.wantErr)
        })
    }
}
```

---

## 三、Mock与依赖注入

### 3.1 接口与测试

```text
为什么接口利于测试：
────────────────────────────────────────

// 依赖具体实现 - 难测试
type UserService struct {
    db *sql.DB  // 真实数据库连接
}

func (s *UserService) GetUser(id int) (*User, error) {
    row := s.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    // ...
}

测试问题：
- 需要真实数据库
- 测试慢
- 依赖外部状态

使用接口改进：
────────────────────────────────────────

// 定义接口
type UserRepository interface {
    GetUser(id int) (*User, error)
}

// 生产实现
type DBRepository struct {
    db *sql.DB
}

func (r *DBRepository) GetUser(id int) (*User, error) {
    // 真实数据库操作
}

// 服务依赖接口
type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.GetUser(id)
}

// Mock实现
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetUser(id int) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

// 测试
func TestUserService(t *testing.T) {
    mockRepo := new(MockRepository)
    service := &UserService{repo: mockRepo}

    expectedUser := &User{ID: 1, Name: "Test"}
    mockRepo.On("GetUser", 1).Return(expectedUser, nil)

    user, err := service.GetUser(1)

    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}
```

### 3.2 手动Mock

```text
简单Mock实现：
────────────────────────────────────────

type MockUserRepository struct {
    users map[int]*User
    err   error
}

func (m *MockUserRepository) GetUser(id int) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("user not found")
    }
    return user, nil
}

// 测试
func TestUserService(t *testing.T) {
    mockRepo := &MockUserRepository{
        users: map[int]*User{
            1: {ID: 1, Name: "Test"},
        },
    }

    service := &UserService{repo: mockRepo}

    user, err := service.GetUser(1)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if user.Name != "Test" {
        t.Errorf("got %s, want Test", user.Name)
    }
}

使用gomock：
────────────────────────────────────────

// 生成mock代码
//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go

func TestWithGomock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)

    service := &UserService{repo: mockRepo}

    mockRepo.EXPECT().GetUser(1).Return(&User{ID: 1}, nil)

    service.GetUser(1)
}
```

---

## 四、基准测试与分析

### 4.1 编写基准测试

```text
基础基准测试：
────────────────────────────────────────

func Fib(n int) int {
    if n < 2 {
        return n
    }
    return Fib(n-1) + Fib(n-2)
}

func BenchmarkFib(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fib(10)
    }
}

运行：
go test -bench=Fib -benchmem

输出：
BenchmarkFib-8    2865334    416 ns/op    0 B/op    0 allocs/op

解读：
- 2865334：运行次数
- 416 ns/op：每次操作耗时
- 0 B/op：每次操作分配内存
- 0 allocs/op：每次操作分配次数

带参数的基准测试：
────────────────────────────────────────

func BenchmarkFibSizes(b *testing.B) {
    sizes := []int{10, 20, 30}
    for _, n := range sizes {
        b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Fib(n)
            }
        })
    }
}

输出：
BenchmarkFibSizes/n=10-8    2983456    402 ns/op
BenchmarkFibSizes/n=20-8      30956  38792 ns/op
BenchmarkFibSizes/n=30-8         32   37429512 ns/op

对比基准：
────────────────────────────────────────

func BenchmarkConcat(b *testing.B) {
    b.Run("strings+", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = "Hello" + " " + "World"
        }
    })

    b.Run("fmt.Sprintf", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = fmt.Sprintf("%s %s", "Hello", "World")
        }
    })

    b.Run("strings.Builder", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var b strings.Builder
            b.WriteString("Hello")
            b.WriteString(" ")
            b.WriteString("World")
            _ = b.String()
        }
    })
}
```

### 4.2 内存分析

```text
内存分配追踪：
────────────────────────────────────────

func BenchmarkAlloc(b *testing.B) {
    b.ReportAllocs()  // 显示分配信息

    for i := 0; i < b.N; i++ {
        _ = make([]int, 100)
    }
}

重置计时器：
────────────────────────────────────────

func BenchmarkComplex(b *testing.B) {
    // 准备数据，不计入基准
    data := generateLargeData()

    b.ResetTimer()  // 重置计时器

    for i := 0; i < b.N; i++ {
        process(data)
    }

    b.StopTimer()  // 停止计时
    cleanup()
}

并行基准测试：
────────────────────────────────────────

func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // 并发执行的代码
            doWork()
        }
    })
}

这个会启动GOMAXPROCS个goroutine并发执行，
用于测试并发安全性和扩展性。
```

---

## 五、测试覆盖率

### 5.1 覆盖率分析

```text
运行覆盖率测试：
────────────────────────────────────────

生成覆盖率报告：
go test -cover

输出：
PASS
coverage: 78.3% of statements
ok      mypackage    0.123s

详细覆盖率：
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

查看函数覆盖率：
go tool cover -func=coverage.out

设置覆盖率目标：
────────────────────────────────────────

CI/CD中强制执行：
go test -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | awk '{if($1 < 80) exit 1}'

覆盖率类型：
────────────────────────────────────────

-set：语句是否执行
-count：执行次数
-atomic：原子计数（并发安全）

边界情况：
────────────────────────────────────────

100%覆盖率 ≠ 无bug

// 有100%覆盖率但有问题的代码
func Divide(a, b int) int {
    if b != 0 {  // 覆盖：true 和 false
        return a / b
    }
    return 0  // 隐藏了错误处理
}

测试：
func TestDivide(t *testing.T) {
    assert.Equal(t, 2, Divide(10, 5))  // b != 0
    assert.Equal(t, 0, Divide(10, 0))  // b == 0
}

覆盖率100%，但：
- 除0返回0可能不是预期行为
- 没有测试负数
- 没有测试最大值

结论：
覆盖率是有用指标，但不是唯一指标。
质量测试 > 高覆盖率
```

---

*本章提供了Go测试的全面指南，从基础到高级技术。*
