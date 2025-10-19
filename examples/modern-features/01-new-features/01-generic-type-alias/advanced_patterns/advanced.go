package advanced_patterns

import (
	"context"
	"encoding/json"
	"fmt"
)

// --- 场景1: 简化Web框架中的Handler和Middleware签名 ---

// HandlerFunc 定义了一个通用的HTTP处理器函数签名。
// 它接受一个请求类型 Req，并返回一个响应类型 Res 或一个错误。
type HandlerFunc[Req, Res any] func(ctx context.Context, req Req) (Res, error)

// Middleware 定义了一个通用的中间件函数签名。
// 它接收一个处理器，并返回一个新的处理器，实现了装饰器模式。
type Middleware[Req, Res any] func(next HandlerFunc[Req, Res]) HandlerFunc[Req, Res]

// --- 示例实现 ---

// 1. 定义具体的请求和响应结构体
type GetUserRequest struct {
	UserID string `json:"user_id"`
}

type GetUserResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// 2. 实现一个具体的处理器
func GetUserHandler(ctx context.Context, req GetUserRequest) (GetUserResponse, error) {
	fmt.Printf("Handling request for user: %s\n", req.UserID)
	// 在实际应用中，这里会查询数据库
	if req.UserID == "123" {
		return GetUserResponse{UserID: "123", Name: "John Doe"}, nil
	}
	return GetUserResponse{}, fmt.Errorf("user not found")
}

// 3. 实现一个日志中间件
func LoggingMiddleware[Req, Res any](next HandlerFunc[Req, Res]) HandlerFunc[Req, Res] {
	return func(ctx context.Context, req Req) (Res, error) {
		// 将请求结构体转为JSON字符串以便记录
		reqBytes, _ := json.Marshal(req)
		fmt.Printf("Request received: %s\n", string(reqBytes))

		// 调用下一个处理器
		res, err := next(ctx, req)
		if err != nil {
			fmt.Printf("Handler error: %v\n", err)
		}

		// 将响应结构体转为JSON字符串
		resBytes, _ := json.Marshal(res)
		fmt.Printf("Response sent: %s\n", string(resBytes))
		return res, err
	}
}

// --- 场景2: 组合约束和别名 ---

// Hasher 约束，要求类型必须有一个 Hash() 方法返回字符串。
type Hasher interface {
	Hash() string
}

// Serializable 约束，要求类型可以被序列化为 []byte。
type Serializable interface {
	Serialize() ([]byte, error)
}

// Storable 约束组合了 Hasher 和 Serializable，代表可存储的对象。
type Storable interface {
	Hasher
	Serializable
}

// Cache 是一个泛型缓存，其键必须是可比较的，值必须是可存储的。
type Cache[K comparable, V Storable] struct {
	data map[K]V
}

// 使用别名简化Cache的定义
type UserCache = Cache[string, UserData]

// UserData 是一个符合Storable约束的示例类型。
type UserData struct {
	ID   string
	Data string
}

func (u UserData) Hash() string {
	return u.ID // 简单示例
}

func (u UserData) Serialize() ([]byte, error) {
	return json.Marshal(u)
}

// PrintAdvancedUsage 演示了高级模式的用法。
func PrintAdvancedUsage() {
	fmt.Println("--- Advanced Pattern 1: HTTP Middleware ---")
	// 定义一个处理器别名
	type GetUserHandlerFunc = HandlerFunc[GetUserRequest, GetUserResponse]
	// 定义一个中间件别名
	type GetUserMiddleware = Middleware[GetUserRequest, GetUserResponse]

	// 组合处理器和中间件
	var h GetUserHandlerFunc = GetUserHandler
	var m GetUserMiddleware = LoggingMiddleware[GetUserRequest, GetUserResponse]

	chainedHandler := m(h)

	// 模拟HTTP请求
	req := GetUserRequest{UserID: "123"}
	_, _ = chainedHandler(context.Background(), req)

	fmt.Println("\n--- Advanced Pattern 2: Constraint Combination ---")
	// 使用别名创建缓存实例
	userCache := UserCache{data: make(map[string]UserData)}
	user := UserData{ID: "user-456", Data: "Some important data"}

	userCache.data[user.Hash()] = user
	fmt.Printf("UserCache contains: %+v\n", userCache.data)
}
