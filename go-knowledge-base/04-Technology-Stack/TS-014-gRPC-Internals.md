# TS-014: gRPC 内部机制深度解析 (gRPC Internals)

> **维度**: Technology Stack
> **级别**: S (17+ KB)
> **标签**: #grpc #protobuf #http2 #rpc #streaming
> **权威来源**: [gRPC Documentation](https://grpc.io/docs/), [gRPC Core](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md)
> **版本**: gRPC 1.70+

---

## gRPC 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      gRPC Architecture                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Service Definition (Proto)                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ service UserService {                                               │    │
│  │   rpc GetUser(GetUserRequest) returns (User);                       │    │
│  │   rpc ListUsers(ListUsersRequest) returns (stream User);            │    │
│  │   rpc CreateUsers(stream CreateUserRequest) returns (UserList);     │    │
│  │   rpc Chat(stream Message) returns (stream Message);                │    │
│  │ }                                                                   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                              │                                               │
│                              ▼ protoc-gen-go-grpc                            │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      Generated Code                                 │    │
│  │  - Client Interface                                                 │    │
│  │  - Server Interface                                                 │    │
│  │  - Message Structs (protobuf)                                       │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│          │                                    │                              │
│          ▼                                    ▼                              │
│  ┌───────────────┐                    ┌───────────────┐                      │
│  │    Client     │◄─── HTTP/2 ───────►│    Server     │                      │
│  │               │    over TLS        │               │                      │
│  │ ┌───────────┐ │                    │ ┌───────────┐ │                      │
│  │ │ Channel   │ │                    │ │ Transport │ │                      │
│  │ │ Stub      │ │                    │ │ Handler   │ │                      │
│  │ │ Intercept │ │                    │ │ Service   │ │                      │
│  │ └───────────┘ │                    │ └───────────┘ │                      │
│  └───────────────┘                    └───────────────┘                      │
│                                                                              │
│  四种服务类型:                                                                │
│  1. Unary: 简单请求-响应                                                     │
│  2. Server Streaming: 服务端流                                               │
│  3. Client Streaming: 客户端流                                               │
│  4. Bidirectional Streaming: 双向流                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## HTTP/2 传输层

### gRPC over HTTP/2

```
请求头:
┌─────────────────────────────────────────────────────────────────────────────┐
│  :method = POST              ← 必须是 POST                                  │
│  :scheme = https             ← 或 http (开发环境)                           │
│  :authority = api.example.com                                               │
│  :path = /user.UserService/GetUser  ← 服务名/方法名                          │
│                                                                              │
│  content-type = application/grpc       ← 或 +proto, +json                    │
│  te = trailers                         ← 必需                                │
│  grpc-timeout = 10S                    ← 可选超时                            │
│  grpc-encoding = gzip                  ← 压缩算法                            │
│                                                                              │
│  custom-metadata-1 = value             ← 自定义元数据                         │
└─────────────────────────────────────────────────────────────────────────────┘

消息帧:
┌─────────────────────────────────────────────────────────────────────────────┐
│  LENGTH (4 bytes)  │  TYPE (1 byte)  │  FLAGS (1 byte)  │  DATA             │
│  0-16777215        │  0=DATA         │  END_STREAM      │  protobuf/JSON    │
│                    │  1=HEADERS      │  END_HEADERS     │                   │
└─────────────────────────────────────────────────────────────────────────────┘

响应尾:
┌─────────────────────────────────────────────────────────────────────────────┐
│  grpc-status = 0           ← 0=OK, 其他=错误码                               │
│  grpc-message =            ← 错误信息 (可选)                                 │
│  custom-trailer = value    ← 自定义 trailer                                  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 流控制

```
HTTP/2 流控制机制:

1. 连接级流控制
   - 整个 HTTP/2 连接的字节限制
   - 默认 65535 字节，可动态更新

2. 流级流控制
   - 每个 gRPC 流的字节限制
   - 独立窗口，防止单个流阻塞

3. gRPC 消息流控制
   - 应用层 ACK 机制
   - 背压传播到业务层
```

---

## Go gRPC 实现

### 服务端

```go
package server

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/health/grpc_health_v1"
    "google.golang.org/grpc/keepalive"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"

    pb "github.com/example/proto"
)

type UserServiceServer struct {
    pb.UnimplementedUserServiceServer
    repo UserRepository
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    // 元数据读取
    md, ok := metadata.FromIncomingContext(ctx)
    if ok {
        auth := md.Get("authorization")
        log.Printf("Auth: %v", auth)
    }

    // 获取用户
    user, err := s.repo.GetByID(ctx, req.Id)
    if err != nil {
        return nil, status.Error(codes.NotFound, "user not found")
    }

    return &pb.User{
        Id:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}

func (s *UserServiceServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    users, err := s.repo.List(stream.Context(), int(req.PageSize), int(req.PageToken))
    if err != nil {
        return err
    }

    for _, user := range users {
        if err := stream.Send(&pb.User{
            Id:   user.ID,
            Name: user.Name,
        }); err != nil {
            return err
        }
    }
    return nil
}

// 创建 gRPC 服务器
func NewGRPCServer() (*grpc.Server, error) {
    // TLS 凭证
    creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
    if err != nil {
        return nil, err
    }

    // 拦截器链
    unaryInterceptors := grpc.ChainUnaryInterceptor(
        loggingInterceptor,
        authInterceptor,
        recoveryInterceptor,
    )

    streamInterceptors := grpc.ChainStreamInterceptor(
        loggingStreamInterceptor,
    )

    // Keepalive 参数
    kaParams := keepalive.ServerParameters{
        MaxConnectionIdle:     5 * time.Minute,
        MaxConnectionAge:      2 * time.Hour,
        MaxConnectionAgeGrace: 5 * time.Minute,
        Time:                  1 * time.Minute,
        Timeout:               20 * time.Second,
    }

    server := grpc.NewServer(
        grpc.Creds(creds),
        unaryInterceptors,
        streamInterceptors,
        grpc.KeepaliveParams(kaParams),
        grpc.MaxConcurrentStreams(100),
        grpc.MaxRecvMsgSize(4*1024*1024),  // 4MB
        grpc.MaxSendMsgSize(4*1024*1024),
    )

    // 注册服务
    pb.RegisterUserServiceServer(server, &UserServiceServer{})

    // 健康检查
    grpc_health_v1.RegisterHealthServer(server, &healthServer{})

    return server, nil
}

// 拦截器示例
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    log.Printf("[%s] %s %v", time.Since(start), info.FullMethod, err)
    return resp, err
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }

    tokens := md.Get("authorization")
    if len(tokens) == 0 || !validateToken(tokens[0]) {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }

    return handler(ctx, req)
}
```

### 客户端

```go
package client

import (
    "context"
    "crypto/tls"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/backoff"
    "google.golang.org/grpc/connectivity"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/keepalive"

    pb "github.com/example/proto"
)

func NewUserClient(target string) (*grpc.ClientConn, pb.UserServiceClient, error) {
    // TLS 配置
    creds := credentials.NewClientTLSFromCert(nil, "")

    // 连接配置
    config := &grpc.Config{
        // 拦截器
        UnaryInterceptor:  grpc.ChainUnaryInterceptor(timeoutInterceptor, retryInterceptor),
        StreamInterceptor: grpc.ChainStreamInterceptor(loggingStreamInterceptor),

        // Keepalive
        KeepaliveParams: keepalive.ClientParameters{
            Time:                10 * time.Second,
            Timeout:             20 * time.Second,
            PermitWithoutStream: true,
        },

        // 连接超时
        DialTimeout: 10 * time.Second,

        // 退避策略
        ConnectParams: grpc.ConnectParams{
            Backoff: backoff.Config{
                BaseDelay:  1.0 * time.Second,
                Multiplier: 1.6,
                Jitter:     0.2,
                MaxDelay:   120 * time.Second,
            },
            MinConnectTimeout: 20 * time.Second,
        },
    }

    conn, err := grpc.Dial(target,
        grpc.WithTransportCredentials(creds),
        grpc.WithUnaryInterceptor(config.UnaryInterceptor),
        grpc.WithKeepaliveParams(config.KeepaliveParams),
    )
    if err != nil {
        return nil, nil, err
    }

    client := pb.NewUserServiceClient(conn)
    return conn, client, nil
}

// 使用示例
func exampleUsage() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    conn, client, err := NewUserClient("api.example.com:443")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Unary 调用
    user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "123"})
    if err != nil {
        st, ok := status.FromError(err)
        if ok {
            log.Printf("gRPC error: %s - %s", st.Code(), st.Message())
        }
        return
    }
    log.Printf("User: %v", user)
}
```

---

## 性能优化

| 优化项 | 配置 | 说明 |
|--------|------|------|
| 连接池 | 复用 ClientConn | 不要每次创建新连接 |
| 批处理 | 流式 RPC | 减少网络往返 |
| 压缩 | `grpc.UseCompressor("gzip")` | 大 payload 压缩 |
| 连接预热 | 提前建立连接 | 避免冷启动延迟 |
| 负载均衡 | `grpc.WithDefaultServiceConfig` | 客户端负载均衡 |

---

## 参考文献

1. [gRPC Documentation](https://grpc.io/docs/)
2. [gRPC Core Protocol](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md)
3. [Go gRPC Middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)
