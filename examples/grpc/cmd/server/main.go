package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 定义简单的 gRPC 服务（简化示例，不依赖 internal 包）

// LoggingInterceptor 日志拦截器
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	log.Printf("gRPC request: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	if err != nil {
		log.Printf("gRPC error: %s, duration: %v, error: %v", info.FullMethod, duration, err)
	} else {
		log.Printf("gRPC success: %s, duration: %v", info.FullMethod, duration)
	}

	return resp, err
}

// RecoveryInterceptor 恢复拦截器
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("gRPC panic recovered: %v", r)
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

func main() {
	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	)

	// 注意：这里需要生成 gRPC 代码后注册实际的服务
	// 当前为简化示例，不包含具体服务实现

	// 启动服务器
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	log.Println("gRPC server listening on :8081")
	log.Println("Note: This is a simplified example. Generate protobuf code to add actual services.")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
