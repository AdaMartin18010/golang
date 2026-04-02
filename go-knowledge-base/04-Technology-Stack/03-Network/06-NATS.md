# TS-NET-006: NATS - Cloud Native Messaging

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #nats #messaging #pub-sub #streaming #cloud-native
> **权威来源**:
> - [NATS Documentation](https://docs.nats.io/) - NATS.io

---

## 1. NATS Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         NATS Messaging Architecture                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         NATS Server                                  │   │
│  │                                                                      │   │
│  │  ┌───────────────────────────────────────────────────────────────┐  │   │
│  │  │                     Subject Namespace                          │  │   │
│  │  │                                                                │  │   │
│  │  │  orders.create ──┐                                            │  │   │
│  │  │  orders.update   │──► ┌───────────┐                           │  │   │
│  │  │  orders.delete ──┘    │  Queue    │                           │  │   │
│  │  │                       │  Group    │                           │  │   │
│  │  │  payments.* ────────► └─────┬─────┘                           │  │   │
│  │  │  payments.process          │                                  │  │   │
│  │  │  payments.refund           ▼                                  │  │   │
│  │  │                       Subscribers                            │  │   │
│  │  │                                                                │  │   │
│  │  └───────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  │  Messaging Patterns:                                                 │   │
│  │  1. Pub-Sub (one-to-many)                                           │   │
│  │  2. Queue Groups (load-balanced)                                    │   │
│  │  3. Request-Reply (synchronous)                                     │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                              │                                               │
│         ┌────────────────────┼────────────────────┐                         │
│         │                    │                    │                         │
│         ▼                    ▼                    ▼                         │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐                │
│  │  Publisher  │      │  Publisher  │      │  Publisher  │                │
│  └─────────────┘      └─────────────┘      └─────────────┘                │
│                                                                              │
│         ┌────────────────────┼────────────────────┐                         │
│         │                    │                    │                         │
│         ▼                    ▼                    ▼                         │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐                │
│  │ Subscriber  │      │ Subscriber  │      │ Subscriber  │                │
│  │  (Queue)    │      │  (Queue)    │      │  (Direct)   │                │
│  └─────────────┘      └─────────────┘      └─────────────┘                │
│                                                                              │
│  Features:                                                                   │
│  - At-most-once delivery (fire-and-forget)                                  │
│  - At-least-once with JetStream                                             │
│  - Subject-based messaging                                                  │
│  - Wildcard subscriptions                                                   │
│  - Auto-discovery of servers                                                │
│  - Clustering support                                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go Client Usage

```go
package main

import (
    "log"
    "time"
    
    "github.com/nats-io/nats.go"
)

func main() {
    // Connect to NATS
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()
    
    // Simple Publisher
    nc.Publish("orders.create", []byte("New order #123"))
    
    // Simple Subscriber
    sub, err := nc.Subscribe("orders.create", func(msg *nats.Msg) {
        log.Printf("Received: %s", string(msg.Data))
    })
    if err != nil {
        log.Fatal(err)
    }
    defer sub.Unsubscribe()
    
    // Queue Group (load balanced)
    queueSub, _ := nc.QueueSubscribe("orders.*", "order-workers", func(msg *nats.Msg) {
        log.Printf("Worker received: %s", string(msg.Data))
    })
    defer queueSub.Unsubscribe()
    
    // Request-Reply
    msg, err := nc.Request("help.request", []byte("Need help"), 2*time.Second)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Response: %s", string(msg.Data))
    
    // Wildcard subscription
    wildcardSub, _ := nc.Subscribe("orders.>", func(msg *nats.Msg) {
        log.Printf("Order event: %s - %s", msg.Subject, string(msg.Data))
    })
    defer wildcardSub.Unsubscribe()
    
    // Keep running
    select {}
}
```

---

## 3. Subject Patterns

```
NATS Subject Patterns:

1. Dot-separated hierarchy
   - orders.create
   - orders.update
   - orders.us.create

2. Wildcards
   - *  : matches a single token
   - >  : matches one or more tokens
   
   Examples:
   - orders.*        : matches orders.create, orders.update
   - orders.>        : matches orders.us.create, orders.eu.update
   - *.create        : matches orders.create, users.create
   - orders.*.create : matches orders.us.create, orders.eu.create

3. Best practices
   - Use descriptive subjects
   - Keep hierarchy shallow (2-4 levels)
   - Use consistent naming
```

---

## 4. Checklist

```
NATS Checklist:
□ Subject naming convention defined
□ Queue groups for load balancing
□ Error handling implemented
□ Reconnection configured
□ Wildcard usage documented
□ JetStream for persistence (if needed)
□ Monitoring enabled
```
