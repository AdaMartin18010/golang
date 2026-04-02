# 测试工具链

> **分类**: 成熟应用领域

---

## testify

```go
import "github.com/stretchr/testify/assert"
import "github.com/stretchr/testify/mock"

func TestSomething(t *testing.T) {
    assert := assert.New(t)

    assert.Equal(123, 123, "they should be equal")
    assert.NotNil(t, object)
    assert.True(t, result)
}
```

---

## gomock

```go
// 生成 mock
//go:generate mockgen -source=store.go -destination=mock_store.go -package=db

type MockStore struct {
    mock.Mock
}

func (m *MockStore) GetUser(id int) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

// 使用
mockStore := new(MockStore)
mockStore.On("GetUser", 1).Return(&User{ID: 1}, nil)
```

---

## Testcontainers

```go
import "github.com/testcontainers/testcontainers-go"

ctx := context.Background()
req := testcontainers.ContainerRequest{
    Image:        "redis:latest",
    ExposedPorts: []string{"6379/tcp"},
    WaitingFor:   wait.ForLog("Ready to accept connections"),
}
redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
    ContainerRequest: req,
    Started:          true,
})
```
