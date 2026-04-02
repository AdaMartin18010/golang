# TS-NET-007: etcd - Distributed Key-Value Store

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #etcd #distributed-systems #key-value #consensus #raft
> **权威来源**:
>
> - [etcd Documentation](https://etcd.io/docs/) - etcd project

---

## 1. etcd Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         etcd Cluster Architecture                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        etcd Cluster (3+ nodes)                       │   │
│  │                                                                      │   │
│  │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │   │
│  │  │   Node 1    │◄──►│   Node 2    │◄──►│   Node 3    │             │   │
│  │  │  (Leader)   │    │  (Follower) │    │  (Follower) │             │   │
│  │  │             │    │             │    │             │             │   │
│  │  │  ┌───────┐  │    │  ┌───────┐  │    │  ┌───────┐  │             │   │
│  │  │  │  Raft │  │    │  │  Raft │  │    │  │  Raft │  │             │   │
│  │  │  │ State │  │    │  │ State │  │    │  │ State │  │             │   │
│  │  │  │Machine│  │    │  │Machine│  │    │  │Machine│  │             │   │
│  │  │  └───┬───┘  │    │  └───┬───┘  │    │  └───┬───┘  │             │   │
│  │  │      │      │    │      │      │    │      │      │             │   │
│  │  │  ┌───▼───┐  │    │  ┌───▼───┐  │    │  ┌───▼───┐  │             │   │
│  │  │  │  WAL  │  │    │  │  WAL  │  │    │  │  WAL  │  │             │   │
│  │  │  │(Write│  │    │  │(Write│  │    │  │(Write│  │             │   │
│  │  │  │ Ahead│  │    │  │ Ahead│  │    │  │ Ahead│  │             │   │
│  │  │  │ Log) │  │    │  │ Log) │  │    │  │ Log) │  │             │   │
│  │  │  └───┬───┘  │    │  └───┬───┘  │    │  └───┬───┘  │             │   │
│  │  │      │      │    │      │      │    │      │      │             │   │
│  │  │  ┌───▼───┐  │    │  ┌───▼───┐  │    │  ┌───▼───┐  │             │   │
│  │  │  │ BoltDB│  │    │  │ BoltDB│  │    │  │ BoltDB│  │             │   │
│  │  │  │(Store)│  │    │  │(Store)│  │    │  │(Store)│  │             │   │
│  │  │  └───────┘  │    │  └───────┘  │    │  └───────┘  │             │   │
│  │  └─────────────┘    └─────────────┘    └─────────────┘             │   │
│  │          │                │                │                        │   │
│  │          └────────────────┴────────────────┘                        │   │
│  │                           │                                         │   │
│  │                      Consensus (Raft)                               │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Key Characteristics:                                                        │
│  - Linearizable reads/writes                                                │
│  - Strong consistency (not eventually consistent)                           │
│  - Leader election with Raft                                                │
│  - Watch for changes                                                        │
│  - TTL for keys                                                             │
│  - Transactions (compare-and-swap)                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go Client Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
    // Create client
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer cli.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Put key-value
    _, err = cli.Put(ctx, "key", "value")
    if err != nil {
        log.Fatal(err)
    }

    // Get value
    resp, err := cli.Get(ctx, "key")
    if err != nil {
        log.Fatal(err)
    }

    for _, ev := range resp.Kvs {
        fmt.Printf("%s : %s\n", ev.Key, ev.Value)
    }

    // Watch for changes
    watchChan := cli.Watch(context.Background(), "key")
    go func() {
        for wresp := range watchChan {
            for _, ev := range wresp.Events {
                fmt.Printf("Watch: %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
            }
        }
    }()

    // Put with TTL
    lease, err := cli.Grant(ctx, 60) // 60 seconds
    if err != nil {
        log.Fatal(err)
    }

    _, err = cli.Put(ctx, "temp", "data", clientv3.WithLease(lease.ID))
    if err != nil {
        log.Fatal(err)
    }

    // Transaction (compare-and-swap)
    txn := cli.Txn(ctx).
        If(clientv3.Compare(clientv3.Value("key"), "=", "value")).
        Then(clientv3.OpPut("key", "new_value")).
        Else(clientv3.OpGet("key"))

    txnResp, err := txn.Commit()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Transaction succeeded: %v\n", txnResp.Succeeded)
}
```

---

## 3. Use Cases

```
etcd Use Cases:

1. Service Discovery
   - Register service endpoints
   - Health checking
   - Load balancing

2. Configuration Management
   - Centralized config store
   - Dynamic configuration
   - Config versioning

3. Distributed Coordination
   - Leader election
   - Distributed locks
   - Barriers

4. Kubernetes
   - Cluster state storage
   - Custom resources
   - Controllers
```

---

## 4. Checklist

```
etcd Checklist:
□ Cluster size odd number (3, 5, 7)
□ Proper endpoints configured
□ TLS enabled for production
□ Regular backups
□ Monitoring in place
□ Watch for critical keys
□ TTL for ephemeral data
□ Transactions for atomic operations
```
