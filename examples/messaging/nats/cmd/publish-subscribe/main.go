package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// 连接到 NATS 服务器（默认本地）
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	// 订阅主题
	sub, err := nc.Subscribe("user.created", func(msg *nats.Msg) {
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

	data, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Failed to marshal message:", err)
	}

	if err := nc.Publish("user.created", data); err != nil {
		log.Fatal("Failed to publish:", err)
	}

	log.Println("Message published successfully")

	// 等待消息处理
	time.Sleep(1 * time.Second)
}
