# gRPC

> **分类**: 开源技术堆栈

---

## 什么是 gRPC

高性能、开源的 RPC 框架，基于 Protocol Buffers 和 HTTP/2。

---

## 定义服务

### proto/user.proto

```protobuf
syntax = "proto3";

package user;

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (stream User);
}

message GetUserRequest {
    int64 id = 1;
}

message User {
    int64 id = 1;
    string name = 2;
    string email = 3;
}
```

---

## 服务端

```go
import pb "github.com/user/project/proto"

type server struct {
    pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    return &pb.User{
        Id:    req.Id,
        Name:  "Alice",
        Email: "alice@example.com",
    }, nil
}

func main() {
    lis, _ := net.Listen("tcp", ":50051")
    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &server{})
    s.Serve(lis)
}
```

---

## 客户端

```go
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := pb.NewUserServiceClient(conn)

ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
if err != nil {
    log.Fatal(err)
}
fmt.Println(user.Name)
```

---

## 拦截器 (中间件)

```go
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    log.Printf("Method: %s, Duration: %v, Error: %v", info.FullMethod, time.Since(start), err)
    return resp, err
}

s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
```
