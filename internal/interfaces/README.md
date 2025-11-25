# Interfaces Layer (接口层)

Clean Architecture 的接口层，包含外部接口适配。

## 结构

```
interfaces/
├── http/          # HTTP 接口
│   ├── chi/       # Chi 路由
│   ├── echo/      # Echo 路由
│   └── openapi/   # OpenAPI 规范
├── grpc/          # gRPC 接口
│   └── proto/     # Protocol Buffers
├── graphql/       # GraphQL 接口
└── asyncapi/      # AsyncAPI 规范
```

## 规则

- ✅ 调用 application 层
- ✅ 处理外部请求
- ✅ 适配不同协议
