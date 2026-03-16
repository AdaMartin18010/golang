package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// 连接到 NATS 服务器
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	// 服务端：订阅并回复
	sub, err := nc.Subscribe("user.get", func(msg *nats.Msg) {
		userID := string(msg.Data)
		log.Printf("Received request for user: %s", userID)

		// 处理请求并回复
		response := "user:" + userID + ":found"
		if err := msg.Respond([]byte(response)); err != nil {
			log.Printf("Failed to respond: %v", err)
		}
	})
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 客户端：发送请求
	reply, err := nc.Request("user.get", []byte("123"), 5*time.Second)
	if err != nil {
		log.Fatal("Failed to send request:", err)
	}

	log.Printf("Response: %s", string(reply.Data))
}
