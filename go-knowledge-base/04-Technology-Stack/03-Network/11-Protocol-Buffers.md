# TS-NET-011: Protocol Buffers in Go

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #protobuf #serialization #grpc #golang #protocol-buffers
> **权威来源**:
>
> - [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers) - Google
> - [Go Protocol Buffers](https://pkg.go.dev/google.golang.org/protobuf) - Go package

---

## 1. Protocol Buffers Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Protocol Buffers Architecture                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Protocol Buffers vs JSON:                                                   │
│                                                                              │
│  JSON:                                        Protocol Buffers:             │
│  {                                            message Person {              │
│    "id": 123,                                   int32 id = 1;               │
│    "name": "John Doe",                          string name = 2;            │
│    "email": "john@example.com",                 string email = 3;           │
│    "phones": [                                repeated Phone phones = 4;    │
│      {"number": "555-1234",                   }                             │
│       "type": "HOME"                          message Phone {               │
│      }                                          string number = 1;          │
│    ]                                            PhoneType type = 2;         │
│  }                                              }                           │
│                                               enum PhoneType {              │
│  Size: ~80 bytes                              MOBILE = 0;                   │
│  Text format                                  HOME = 1;                     │
│  No schema validation                         WORK = 2;                     │
│  Slower parsing                               }                             │
│                                               }                             │
│                                                                              │
│                                               Binary size: ~20 bytes        │
│                                               Type safe                     │
│                                               Schema evolution              │
│                                               Fast parsing                  │
│                                                                              │
│  Use Cases:                                                                  │
│  - gRPC services                                                             │
│  - Data storage                                                              │
│  - Microservice communication                                                │
│  - Configuration files                                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Defining Messages

```protobuf
// user.proto
syntax = "proto3";

package user;
option go_package = "github.com/example/proto/user";

import "google/protobuf/timestamp.proto";
import "google/protobuf/any.proto";

// User message
message User {
    // Field numbers are used for binary encoding
    int32 id = 1;
    string username = 2;
    string email = 3;

    // Nested message
    Profile profile = 4;

    // Repeated field (array/slice)
    repeated string roles = 5;

    // Enum
    Status status = 6;

    // Timestamp
    google.protobuf.Timestamp created_at = 7;
    google.protobuf.Timestamp updated_at = 8;

    // Oneof - mutually exclusive fields
    oneof contact_method {
        string phone = 9;
        string email_secondary = 10;
    }

    // Map
    map<string, string> metadata = 11;

    // Optional (proto3)
    optional string nickname = 12;

    enum Status {
        UNKNOWN = 0;
        ACTIVE = 1;
        INACTIVE = 2;
        SUSPENDED = 3;
    }
}

message Profile {
    string first_name = 1;
    string last_name = 2;
    string bio = 3;
    string avatar_url = 4;
}

// Service definition (for gRPC)
service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    rpc UpdateUser(UpdateUserRequest) returns (User);
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
}

message GetUserRequest {
    int32 id = 1;
}

message CreateUserRequest {
    User user = 1;
}

message ListUsersRequest {
    int32 page_size = 1;
    string page_token = 2;
}

message ListUsersResponse {
    repeated User users = 1;
    string next_page_token = 2;
}

message UpdateUserRequest {
    User user = 1;
    // Field mask for partial updates
    google.protobuf.FieldMask update_mask = 2;
}

message DeleteUserRequest {
    int32 id = 1;
}

message DeleteUserResponse {
    bool success = 1;
}
```

---

## 3. Code Generation and Usage

```bash
# Install protoc compiler
# Download from https://github.com/protocolbuffers/protobuf/releases

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       user.proto

# Generated files:
# - user.pb.go: Message types and accessors
# - user_grpc.pb.go: gRPC client and server interfaces
```

```go
package main

import (
    "fmt"
    "log"
    "time"

    "google.golang.org/protobuf/types/known/timestamppb"

    pb "github.com/example/proto/user"
)

func main() {
    // Create a new user
    user := &pb.User{
        Id:       1,
        Username: "johndoe",
        Email:    "john@example.com",
        Profile: &pb.Profile{
            FirstName: "John",
            LastName:  "Doe",
            Bio:       "Software engineer",
        },
        Roles:  []string{"user", "admin"},
        Status: pb.User_ACTIVE,
        Metadata: map[string]string{
            "department": "engineering",
            "location":   "sf",
        },
        CreatedAt: timestamppb.Now(),
        UpdatedAt: timestamppb.Now(),
    }

    // Serialize to binary
    data, err := proto.Marshal(user)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Serialized size: %d bytes\n", len(data))

    // Deserialize
    newUser := &pb.User{}
    if err := proto.Unmarshal(data, newUser); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("User: %s (%s)\n", newUser.Username, newUser.Email)

    // JSON serialization (for debugging)
    jsonData, err := protojson.Marshal(user)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("JSON: %s\n", jsonData)
}
```

---

## 4. Best Practices

```
Protocol Buffers Best Practices:
□ Use proto3 for new projects
□ Reserve field numbers when removing fields
□ Use appropriate field types
□ Avoid changing field numbers
□ Use meaningful message and field names
□ Document with comments
□ Version your proto files
□ Use packages for namespacing
```
