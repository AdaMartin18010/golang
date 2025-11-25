# Go分布式系统

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go分布式系统](#go分布式系统)
  - [📋 目录](#-目录)
  - [📚 核心内容](#-核心内容)
  - [🚀 Redis分布式锁](#-redis分布式锁)
  - [📖 系统文档](#-系统文档)

---

## 📚 核心内容

1. **[分布式基础概念](./01-分布式系统基础.md)** ⭐⭐⭐⭐⭐
   - CAP定理
   - BASE理论
   - 一致性级别

2. **[RPC与消息队列](./02-RPC与消息队列.md)** ⭐⭐⭐⭐⭐
   - gRPC
   - Kafka, RabbitMQ, NATS

3. **[分布式一致性](./03-分布式一致性.md)** ⭐⭐⭐⭐⭐
   - Paxos, Raft, ZAB
   - 一致性哈希

4. **[分布式锁](./04-分布式锁.md)** ⭐⭐⭐⭐⭐
   - Redis锁
   - etcd锁
   - ZooKeeper锁

5. **[分布式事务](./05-分布式事务.md)** ⭐⭐⭐⭐⭐
   - 2PC/3PC
   - TCC
   - Saga

6. **[负载均衡](./06-负载均衡.md)** ⭐⭐⭐⭐
   - 算法
   - 实现

---

## 🚀 Redis分布式锁

```go
func AcquireLock(ctx Context.Context, key string, ttl time.Duration) bool {
    result := rdb.SetNX(ctx, key, "locked", ttl)
    return result.Val()
}

func ReleaseLock(ctx Context.Context, key string) {
    rdb.Del(ctx, key)
}
```

---

## 📖 系统文档
