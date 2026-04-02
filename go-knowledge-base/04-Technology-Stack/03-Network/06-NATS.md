# NATS

> **分类**: 开源技术堆栈

---

## nats.go

```go
import "github.com/nats-io/nats.go"

nc, _ := nats.Connect(nats.DefaultURL)
defer nc.Close()
```

---

## Pub/Sub

```go
// 订阅
sub, _ := nc.Subscribe("foo", func(m *nats.Msg) {
    fmt.Printf("Received: %s\n", string(m.Data))
})

// 发布
nc.Publish("foo", []byte("Hello World"))
```

---

## Request/Reply

```go
// 请求
msg, _ := nc.Request("help", []byte("help me"), 100*time.Millisecond)

// 回复
nc.Subscribe("help", func(m *nats.Msg) {
    nc.Publish(m.Reply, []byte("I can help!"))
})
```
