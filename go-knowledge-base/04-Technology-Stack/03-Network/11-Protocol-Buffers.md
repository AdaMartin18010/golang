# Protocol Buffers

> **分类**: 开源技术堆栈  
> **标签**: #protobuf #serialization #grpc

---

## 定义消息

```protobuf
syntax = "proto3";

package user;

option go_package = "github.com/example/api/user";

message User {
    int64 id = 1;
    string email = 2;
    string name = 3;
    Status status = 4;
    google.protobuf.Timestamp created_at = 5;
    
    repeated string tags = 6;
    map<string, string> metadata = 7;
}

enum Status {
    UNKNOWN = 0;
    ACTIVE = 1;
    INACTIVE = 2;
}

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (stream User);
    rpc CreateUser(CreateUserRequest) returns (User);
}
```

---

## 生成代码

```bash
# 安装 protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go-grpc@latest

# 生成
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    user.proto
```

---

## 使用生成的代码

```go
import pb "github.com/example/api/user"

// 创建消息
user := &pb.User{
    Id:        1,
    Email:     "user@example.com",
    Name:      "John Doe",
    Status:    pb.Status_ACTIVE,
    CreatedAt: timestamppb.Now(),
    Tags:      []string{"vip", "beta"},
    Metadata: map[string]string{
        "source": "web",
    },
}

// 序列化
data, err := proto.Marshal(user)
if err != nil {
    log.Fatal(err)
}

// 反序列化
newUser := &pb.User{}
if err := proto.Unmarshal(data, newUser); err != nil {
    log.Fatal(err)
}
```

---

## 与 JSON 对比

| 特性 | Protobuf | JSON |
|------|----------|------|
| 大小 | 小 (二进制) | 大 (文本) |
| 速度 | 快 | 慢 |
| 可读性 | 差 | 好 |
| 模式 | 强类型 | 无模式 |
| 适用 | 服务间通信 | Web API |

---

## Oneof 类型

```protobuf
message Notification {
    int64 id = 1;
    
    oneof content {
        string text = 2;
        bytes image = 3;
        Video video = 4;
    }
}
```

```go
notif := &pb.Notification{
    Id: 1,
    Content: &pb.Notification_Text{
        Text: "Hello",
    },
}

switch c := notif.Content.(type) {
case *pb.Notification_Text:
    fmt.Println(c.Text)
case *pb.Notification_Image:
    fmt.Println(len(c.Image))
}
```

---

## 版本兼容

```protobuf
// v1
message User {
    int64 id = 1;
    string name = 2;
}

// v2 - 向后兼容
message User {
    int64 id = 1;
    string name = 2;
    string email = 3;  // 新增字段
    reserved 4, 5;      // 保留字段号
    reserved "phone";   // 保留字段名
}
```

---

## 最佳实践

1. **使用 proto3** - 更简洁的语法
2. **字段编号不要复用** - 向后兼容
3. **避免使用 required** - 灵活性
4. **使用枚举** - 类型安全
5. **使用时间戳类型** - 标准格式
