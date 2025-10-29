# Go微服务开发

Go微服务架构完整指南，涵盖gRPC、服务发现、API网关和服务网格。

---

## 📚 核心内容

1. **微服务架构设计**
2. **gRPC与Protobuf**
3. **服务发现与注册**
4. **API网关**
5. **服务网格**
6. **监控与追踪**

---

## 🚀 gRPC示例

```go
// 服务端
type server struct {
    pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
```

---

## 📖 系统文档

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3
