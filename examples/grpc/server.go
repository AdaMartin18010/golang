package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/yourusername/golang/internal/interfaces/grpc/handlers"
	"github.com/yourusername/golang/internal/interfaces/grpc/interceptors"
	// TODO: 生成 gRPC 代码后取消注释
	// userpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/userpb"
	// healthpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/healthpb"
)

func main() {
	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingUnaryInterceptor,
			interceptors.TracingUnaryInterceptor,
		),
	)

	// TODO: 生成 gRPC 代码后取消注释
	// 注册服务
	// userHandler := handlers.NewUserHandler(nil) // 需要传入实际的 user.Service
	// userpb.RegisterUserServiceServer(grpcServer, userHandler)

	// healthHandler := handlers.NewHealthHandler()
	// healthpb.RegisterHealthServiceServer(grpcServer, healthHandler)

	// 启动服务器
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	log.Println("gRPC server listening on :8081")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
