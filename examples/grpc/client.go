package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// TODO: 生成 gRPC 代码后取消注释
	// userpb "github.com/yourusername/golang/internal/interfaces/grpc/proto/userpb"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	// TODO: 生成 gRPC 代码后取消注释
	// client := userpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 示例：获取用户
	// req := &userpb.GetUserRequest{Id: "123"}
	// resp, err := client.GetUser(ctx, req)
	// if err != nil {
	//     log.Fatal("Failed to get user:", err)
	// }
	// log.Printf("User: %+v", resp.User)

	log.Println("gRPC client example (requires generated code)")
}
