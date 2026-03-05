# Go测试策略深度解析

> 全面的Go测试方法论：单元测试、集成测试、契约测试

---

## 一、测试金字塔

```text
测试金字塔:
────────────────────────────────────────

    /\
   /E2E\          端到端测试 (少量)
  /─────\
 /Integration\    集成测试 (中等)
/─────────────\
/  Unit Tests   \  单元测试 (大量)
─────────────────

比例建议:
├─ 单元测试: 70%
├─ 集成测试: 20%
└─ E2E测试: 10%
```

---

## 二、单元测试

### 2.1 表驱动测试

```text
表驱动测试模式:
────────────────────────────────────────

代码示例:
func Calculate(a, b int, op string) (int, error) {
    switch op {
    case "+":
        return a + b, nil
    case "-":
        return a - b, nil
    case "*":
        return a * b, nil
    case "/":
        if b == 0 {
            return 0, errors.New("division by zero")
        }
        return a / b, nil
    default:
        return 0, fmt.Errorf("unknown operator: %s", op)
    }
}

func TestCalculate(t *testing.T) {
    tests := []struct {
        name    string
        a, b    int
        op      string
        want    int
        wantErr bool
    }{
        {"add", 1, 2, "+", 3, false},
        {"subtract", 5, 3, "-", 2, false},
        {"multiply", 4, 5, "*", 20, false},
        {"divide", 10, 2, "/", 5, false},
        {"divide by zero", 10, 0, "/", 0, true},
        {"invalid op", 1, 2, "^", 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Calculate(tt.a, tt.b, tt.op)
            if (err != nil) != tt.wantErr {
                t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Calculate() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 2.2 Mock测试

```text
使用gomock:
────────────────────────────────────────

//go:generate mockgen -source=service.go -destination=mocks/mock_service.go -package=mocks

// 接口
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

// 测试
func TestUserService_GetUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)

    // 设置期望
    mockRepo.EXPECT().
        GetByID(gomock.Any(), "123").
        Return(&User{ID: "123", Name: "John"}, nil)

    service := NewUserService(mockRepo)
    user, err := service.GetUser(context.Background(), "123")

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if user.Name != "John" {
        t.Errorf("unexpected user: %v", user)
    }
}
```

---

## 三、集成测试

### 3.1 数据库集成测试

```text
数据库测试:
────────────────────────────────────────

使用testcontainers:
import "github.com/testcontainers/testcontainers-go"

func TestUserRepository(t *testing.T) {
    ctx := context.Background()

    // 启动PostgreSQL容器
    req := testcontainers.ContainerRequest{
        Image:        "postgres:13",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
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

    // 连接数据库并测试
    dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb", host, port)
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        t.Fatal(err)
    }

    // 执行测试...
}
```

### 3.2 HTTP集成测试

```text
HTTP测试:
────────────────────────────────────────

func TestAPI(t *testing.T) {
    // 创建测试服务器
    handler := setupRouter()
    server := httptest.NewServer(handler)
    defer server.Close()

    client := &http.Client{Timeout: 5 * time.Second}

    // 测试创建用户
    resp, err := client.Post(
        server.URL+"/api/users",
        "application/json",
        strings.NewReader(`{"name":"John","email":"john@example.com"}`),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        t.Errorf("expected 201, got %d", resp.StatusCode)
    }

    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        t.Fatal(err)
    }

    if user.Name != "John" {
        t.Errorf("expected John, got %s", user.Name)
    }
}
```

---

## 四、契约测试

```text
Pact契约测试:
────────────────────────────────────────

消费者测试:
func TestConsumer(t *testing.T) {
    pact := dsl.Pact{
        Consumer: "user-service-client",
        Provider: "user-service",
    }

    pact.AddInteraction().
        Given("User exists").
        UponReceiving("A request for user").
        WithRequest(dsl.Request{
            Method: "GET",
            Path:   dsl.String("/users/123"),
        }).
        WillRespondWith(dsl.Response{
            Status: 200,
            Body: dsl.MapMatcher{
                "id":   dsl.String("123"),
                "name": dsl.String("John"),
            },
        })

    if err := pact.Verify(testUserClient); err != nil {
        t.Fatal(err)
    }
}

提供者验证:
func TestProvider(t *testing.T) {
    pact := dsl.Pact{
        Provider: "user-service",
    }

    _, err := pact.VerifyProvider(t, types.VerifyRequest{
        ProviderBaseURL: "http://localhost:8080",
        PactURLs:        []string{"./pacts/user-service-client-user-service.json"},
    })

    if err != nil {
        t.Fatal(err)
    }
}
```

---

## 五、测试覆盖率

```text
覆盖率目标:
────────────────────────────────────────

目标:
├─ 单元测试覆盖率: > 80%
├─ 核心逻辑覆盖率: > 90%
└─ 集成测试: 覆盖主要场景

生成报告:
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

排除文件:
//go:build !test

代码示例:
// 不可测试的代码分离
//go:build test

package main

var timeNow = time.Now

//go:build !test

package main

var timeNow = func() time.Time { return time.Now() }
```

---

*本章提供了全面的Go测试策略，涵盖单元测试、集成测试和契约测试。*
