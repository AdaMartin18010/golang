package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/yourusername/golang/internal/infrastructure/messaging/nats"
)

func main() {
	// 创建客户端
	client, err := nats.NewClient(nats.DefaultConfig())
	if err != nil {
		log.Fatal("Failed to create NATS client:", err)
	}
	defer client.Close()

	// 订阅主题
	sub, err := client.Subscribe("user.created", func(msg *nats.Msg) {
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return
		}
		log.Printf("Received message: %+v", data)
	})
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}
	defer sub.Unsubscribe()

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发布消息
	message := map[string]interface{}{
		"user_id": 123,
		"name":    "Alice",
		"email":   "alice@example.com",
	}

	if err := client.Publish("user.created", message); err != nil {
		log.Fatal("Failed to publish:", err)
	}

	log.Println("Message published successfully")

	// 等待消息处理
	time.Sleep(1 * time.Second)
}
