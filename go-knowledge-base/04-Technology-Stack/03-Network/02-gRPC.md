# TS-NET-002: gRPC Architecture and Go Implementation

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #grpc #protobuf #rpc #microservices #streaming
> **权威来源**:
>
> - [gRPC Documentation](https://grpc.io/docs/) - CNCF
> - [Protocol Buffers](https://developers.google.com/protocol-buffers) - Google
> - [gRPC-Go](https://github.com/grpc/grpc-go) - Go implementation

---

## 1. gRPC Architecture

### 1.1 Core Concepts

```
┌─────────────────────────────────────────────────────────────────┐
│                        gRPC Architecture                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Service Definition (.proto)                                     │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ service UserService {                                    │   │
│  │   rpc GetUser(GetUserReq) returns (User);               │   │
│  │   rpc ListUsers(ListReq) returns (stream User);         │   │
│  │   rpc CreateUsers(stream CreateReq) returns (UserList); │   │
│  │   rpc Chat(stream Msg) returns (stream Msg);            │   │
│  │ }                                                        │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              protoc (Protocol Compiler)                  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                   │
│              ┌───────────────┴───────────────┐                   │
│              ▼                               ▼                   │
│  ┌─────────────────────┐         ┌─────────────────────┐       │
│  │   Client Code Gen   │         │   Server Code Gen   │       │
│  └──────────┬──────────┘         └──────────┬──────────┘       │
│             │                               │                   │
│  ┌──────────▼──────────┐         ┌──────────▼──────────┐       │
│  │    Client Stub      │◄───────►│    Server Stub      │       │
│  │  - Marshal request  │  HTTP/2 │  - Unmarshal request│       │
│  │  - Send RPC         │         │  - Invoke handler   │       │
│  │  - Unmarshal resp   │         │  - Marshal response │       │
│  └──────────┬──────────┘         └──────────┬──────────┘       │
│             │                               │                   │
│  ┌──────────▼──────────┐         ┌──────────▼──────────┐       │
│  │    Channel (HTTP/2) │◄───────►│     Transport       │       │
│  │  - Connection mgmt  │         │  - HTTP/2 frames    │       │
│  │  - Flow control     │         │  - TLS              │       │
│  │  - Load balancing   │         │  - Compression      │       │
│  └─────────────────────┘         └─────────────────────┘       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 RPC Types

```
┌─────────────────────────────────────────────────────────────────┐
│                      gRPC RPC Types                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. UNARY (Request-Response)                                     │
│  ┌────────┐   request    ┌────────┐   response   ┌────────┐     │
│  │ Client │─────────────►│ Server │─────────────►│ Client │     │
│  └────────┘              └────────┘              └────────┘     │
│                                                                  │
│  2. SERVER STREAMING                                             │
│  ┌────────┐   request    ┌────────┐  response1   ┌────────┐     │
│  │ Client │─────────────►│ Server │─────────────►│        │     │
│  └────────┘              └────────┤  response2   │ Client │     │
│                                   ├─────────────►│        │     │
│                                   │  response3   └────────┘     │
│                                   └─────────────►               │
│                                                                  │
│  3. CLIENT STREAMING                                             │
│  ┌────────┐   request1   ┌────────┐                             │
│  │        │─────────────►│        │                             │
│  │ Client │   request2   │ Server │   response   ┌────────┐     │
│  │        │─────────────►│        │─────────────►│ Client │     │
│  │        │   request3   │        │              └────────┘     │
│  └────────┴─────────────►└────────┘                             │
│                                                                  │
│  4. BIDIRECTIONAL STREAMING                                      │
│  ┌────────┐   msg1       ┌────────┐   msg1       ┌────────┐     │
│  │        │─────────────►│        │─────────────►│        │     │
│  │ Client │              │ Server │              │ Client │     │
│  │        │◄─────────────│        │◄─────────────│        │     │
│  └────────┘   msg2       └────────┘   msg2       └────────┘     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Protocol Buffers

### 2.1 Message Definition

```protobuf
syntax = "proto3";
package user;
option go_package = "github.com/example/proto/user";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";

// User message
message User {
    string id = 1;
    string email = 2;
    string name = 3;
    int32 age = 4;
    Status status = 5;
    google.protobuf.Timestamp created_at = 6;
    repeated string tags = 7;
    Address address = 8;

    enum Status {
        UNKNOWN = 0;
        ACTIVE = 1;
        INACTIVE = 2;
        BANNED = 3;
    }
}

message Address {
    string street = 1;
    string city = 2;
    string country = 3;
    string zip_code = 4;
}

// Request/Response messages
message GetUserRequest {
    string user_id = 1;
}

message ListUsersRequest {
    int32 page_size = 1;
    string page_token = 2;
    string filter = 3;
}

message ListUsersResponse {
    repeated User users = 1;
    string next_page_token = 2;
    int32 total_count = 3;
}

message CreateUserRequest {
    User user = 1;
}

message UpdateUserRequest {
    User user = 1;
    google.protobuf.FieldMask update_mask = 2;
}

message DeleteUserRequest {
    string user_id = 1;
}

// Service definition
service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc UpdateUser(UpdateUserRequest) returns (User);
    rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty);

    // Streaming methods
    rpc StreamUsers(ListUsersRequest) returns (stream User);
    rpc BatchCreateUsers(stream CreateUserRequest) returns (ListUsersResponse);
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}

message ChatMessage {
    string user_id = 1;
    string content = 2;
    google.protobuf.Timestamp timestamp = 3;
}
```

### 2.2 Code Generation

```bash
# Install protoc-gen-go
 go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       user.proto

# Generated files:
# - user.pb.go: Message types
# - user_grpc.pb.go: Client and server interfaces
```

---

## 3. Go Implementation

### 3.1 Server Implementation

```go
package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"

    pb "github.com/example/proto/user"
)

// Server struct implements generated UserServiceServer interface
type userServer struct {
    pb.UnimplementedUserServiceServer
    users map[string]*pb.User
}

// Unary RPC
func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    // Authentication check
    if err := s.authenticate(ctx); err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "auth failed: %v", err)
    }

    user, exists := s.users[req.UserId]
    if !exists {
        return nil, status.Errorf(codes.NotFound, "user not found: %s", req.UserId)
    }

    return user, nil
}

// Server streaming
func (s *userServer) StreamUsers(req *pb.ListUsersRequest, stream pb.UserService_StreamUsersServer) error {
    for _, user := range s.users {
        if err := stream.Send(user); err != nil {
            return err
        }
    }
    return nil
}

// Client streaming
func (s *userServer) BatchCreateUsers(stream pb.UserService_BatchCreateUsersServer) error {
    var users []*pb.User

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.ListUsersResponse{
                Users: users,
            })
        }
        if err != nil {
            return err
        }

        user := req.User
        s.users[user.Id] = user
        users = append(users, user)
    }
}

// Bidirectional streaming
func (s *userServer) Chat(stream pb.UserService_ChatServer) error {
    ctx := stream.Context()

    // Send welcome message
    if err := stream.Send(&pb.ChatMessage{
        UserId:    "system",
        Content:   "Welcome!",
        Timestamp: timestamppb.Now(),
    }); err != nil {
        return err
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }

        // Echo back
        if err := stream.Send(msg); err != nil {
            return err
        }
    }
}

// Authentication helper
func (s *userServer) authenticate(ctx context.Context) error {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return status.Error(codes.Unauthenticated, "missing metadata")
    }

    tokens := md.Get("authorization")
    if len(tokens) == 0 {
        return status.Error(codes.Unauthenticated, "missing token")
    }

    // Validate token...
    return nil
}

func main() {
    // Create listener
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // TLS credentials
    creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
    if err != nil {
        log.Fatalf("failed to load TLS: %v", err)
    }

    // Create gRPC server with interceptors
    srv := grpc.NewServer(
        grpc.Creds(creds),
        grpc.UnaryInterceptor(unaryInterceptor),
        grpc.StreamInterceptor(streamInterceptor),
        grpc.MaxRecvMsgSize(1024*1024*10), // 10MB
        grpc.MaxSendMsgSize(1024*1024*10),
    )

    // Register service
    pb.RegisterUserServiceServer(srv, &userServer{
        users: make(map[string]*pb.User),
    })

    log.Println("Server starting on :50051")
    if err := srv.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

// Unary interceptor
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    start := time.Now()

    // Add request ID to context
    requestID := uuid.New().String()
    ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", requestID)

    // Call handler
    resp, err := handler(ctx, req)

    // Log
    log.Printf("[%s] Method: %s, Duration: %v, Error: %v",
        requestID, info.FullMethod, time.Since(start), err)

    return resp, err
}

// Stream interceptor
func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    wrapped := &wrappedStream{ServerStream: ss}
    return handler(srv, wrapped)
}

type wrappedStream struct {
    grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
    log.Printf("Receive message: %T", m)
    return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
    log.Printf("Send message: %T", m)
    return w.ServerStream.SendMsg(m)
}
```

### 3.2 Client Implementation

```go
package main

import (
    "context"
    "io"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/metadata"

    pb "github.com/example/proto/user"
)

func createClient() (pb.UserServiceClient, *grpc.ClientConn, error) {
    // TLS credentials
    creds, err := credentials.NewClientTLSFromFile("ca.crt", "localhost")
    if err != nil {
        return nil, nil, err
    }

    // Connection options
    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(creds),
        grpc.WithBlock(),
        grpc.WithTimeout(5 * time.Second),
        grpc.WithUnaryInterceptor(clientUnaryInterceptor),
        grpc.WithStreamInterceptor(clientStreamInterceptor),
    }

    // Connect
    conn, err := grpc.Dial("localhost:50051", opts...)
    if err != nil {
        return nil, nil, err
    }

    return pb.NewUserServiceClient(conn), conn, nil
}

func unaryCall(client pb.UserServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Add metadata
    ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer token123")

    resp, err := client.GetUser(ctx, &pb.GetUserRequest{UserId: "user-1"})
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }

    log.Printf("User: %v", resp)
}

func serverStreaming(client pb.UserServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    stream, err := client.StreamUsers(ctx, &pb.ListUsersRequest{PageSize: 100})
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }

    for {
        user, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Error: %v", err)
            return
        }
        log.Printf("User: %v", user.Name)
    }
}

func clientStreaming(client pb.UserServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    stream, err := client.BatchCreateUsers(ctx)
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }

    // Send multiple requests
    for i := 0; i < 10; i++ {
        if err := stream.Send(&pb.CreateUserRequest{
            User: &pb.User{
                Id:    fmt.Sprintf("user-%d", i),
                Name:  fmt.Sprintf("User %d", i),
                Email: fmt.Sprintf("user%d@example.com", i),
            },
        }); err != nil {
            log.Printf("Send error: %v", err)
            return
        }
    }

    resp, err := stream.CloseAndRecv()
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }

    log.Printf("Created %d users", len(resp.Users))
}

func bidirectionalStreaming(client pb.UserServiceClient) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    stream, err := client.Chat(ctx)
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }

    // Send messages
    go func() {
        for i := 0; i < 10; i++ {
            if err := stream.Send(&pb.ChatMessage{
                UserId:    "client",
                Content:   fmt.Sprintf("Message %d", i),
                Timestamp: timestamppb.Now(),
            }); err != nil {
                log.Printf("Send error: %v", err)
                return
            }
            time.Sleep(time.Second)
        }
        stream.CloseSend()
    }()

    // Receive messages
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Recv error: %v", err)
            return
        }
        log.Printf("Received: %s", msg.Content)
    }
}

func clientUnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
    start := time.Now()
    err := invoker(ctx, method, req, reply, cc, opts...)
    log.Printf("Call: %s, Duration: %v, Error: %v", method, time.Since(start), err)
    return err
}

func clientStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
    log.Printf("Streaming: %s", method)
    return streamer(ctx, desc, cc, method, opts...)
}
```

---

## 4. Advanced Features

### 4.1 Load Balancing

```go
import "google.golang.org/grpc/balancer/roundrobin"

// Round-robin load balancing
conn, err := grpc.Dial(
    "dns:///service.example.com",
    grpc.WithDefaultServiceConfig(`{
        "loadBalancingPolicy": "round_robin"
    }`),
    grpc.WithTransportCredentials(creds),
)

// Custom resolver
import "google.golang.org/grpc/resolver"

func init() {
    resolver.Register(&customResolverBuilder{})
}

type customResolverBuilder struct{}

func (b *customResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
    r := &customResolver{
        target: target,
        cc:     cc,
    }
    r.start()
    return r, nil
}

func (b *customResolverBuilder) Scheme() string { return "custom" }
```

### 4.2 Health Checking

```go
import "google.golang.org/grpc/health"
import healthpb "google.golang.org/grpc/health/grpc_health_v1"

// Server-side health check
healthServer := health.NewServer()
healthServer.SetServingStatus("UserService", healthpb.HealthCheckResponse_SERVING)
healthpb.RegisterHealthServer(grpcServer, healthServer)

// Client-side health check
conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
healthClient := healthpb.NewHealthClient(conn)

resp, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{Service: "UserService"})
if err != nil || resp.Status != healthpb.HealthCheckResponse_SERVING {
    // Handle unhealthy
}

// Streaming health check
stream, _ := healthClient.Watch(ctx, &healthpb.HealthCheckRequest{Service: "UserService"})
for {
    resp, err := stream.Recv()
    if err != nil {
        return
    }
    log.Printf("Health status: %v", resp.Status)
}
```

---

## 5. Performance Optimization

```
gRPC Performance Checklist:
□ Use connection pooling (single connection is multiplexed via HTTP/2)
□ Enable compression (gzip or snappy)
□ Set appropriate message size limits
□ Use streaming for large datasets
□ Enable keepalive for connection health
□ Use protobuf instead of JSON (smaller payload)
□ Implement caching for idempotent requests
□ Use load balancing for multiple backends
```

---

## 6. Checklist

```
gRPC Development Checklist:
□ Define clear service contracts in .proto
□ Version your APIs (v1, v2)
□ Use appropriate RPC types
□ Implement proper error handling with status codes
□ Add request ID for tracing
□ Implement authentication/authorization
□ Configure TLS for production
□ Set timeouts and deadlines
□ Implement health checks
□ Add metrics and logging
□ Use interceptors for cross-cutting concerns
```
