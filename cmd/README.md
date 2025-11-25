# Command Applications

主程序入口目录。

## 结构

```text
cmd/
├── server/        # HTTP 服务器
├── grpc-server/   # gRPC 服务器
├── graphql-server/# GraphQL 服务器
├── mqtt-client/   # MQTT 客户端
└── cli/           # CLI 工具
```

## 说明

每个子目录包含一个独立的可执行程序。
