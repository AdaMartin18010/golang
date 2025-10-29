# Go进阶主题

Go语言进阶主题，涵盖性能优化、架构设计、分布式系统、安全、AI/ML和现代Web技术。

---

## 📋 目录结构

### 核心模块

1. **[性能优化](./performance/README.md)** ⭐⭐⭐⭐⭐
   - 性能分析 (pprof)
   - 内存优化
   - 并发优化
   - GC调优

2. **[架构设计](./architecture/README.md)** ⭐⭐⭐⭐⭐
   - 设计模式
   - 架构模式
   - 并发型模式
   - 分布式型模式

3. **[架构实践](./architecture-practices/README.md)** ⭐⭐⭐⭐⭐
   - 微服务架构
   - 事件驱动
   - CQRS
   - 云原生架构

4. **[分布式系统](./distributed/README.md)** ⭐⭐⭐⭐⭐
   - CAP定理
   - 一致性协议 (Paxos/Raft)
   - 分布式锁
   - 分布式事务

5. **[安全](./security/README.md)** ⭐⭐⭐⭐
   - Web安全
   - 身份认证
   - 授权机制
   - 数据保护

6. **[AI/ML](./ai-ml/README.md)** ⭐⭐⭐⭐
   - Go与AI集成
   - 机器学习库
   - 模型推理
   - 数据处理

7. **[现代Web](./modern-web/README.md)** ⭐⭐⭐⭐
   - 现代Web框架
   - 实时通信 (WebSocket)
   - GraphQL
   - 微服务网关

---

## 🎯 学习路径

### 性能专家 (3-4周)
```
性能分析 → 内存优化 → 并发优化 → GC调优
```

### 架构师 (4-6周)
```
设计模式 → 架构模式 → 微服务 → 分布式系统
```

### 安全专家 (2-3周)
```
Web安全 → 认证授权 → 数据保护 → 审计
```

---

## 🚀 快速开始

### 性能分析

```go
import _ "net/http/pprof"

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    // 应用逻辑...
}
```

```bash
# CPU性能分析
go tool pprof http://localhost:6060/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap
```

### 分布式锁 (Redis)

```go
import "github.com/go-redis/redis/v8"

func AcquireLock(ctx context.Context, key string, ttl time.Duration) bool {
    result := rdb.SetNX(ctx, key, "locked", ttl)
    return result.Val()
}

func ReleaseLock(ctx context.Context, key string) {
    rdb.Del(ctx, key)
}
```

---

## 📖 系统文档

- **[知识图谱](./00-知识图谱.md)**: 进阶知识体系全景图
- **[对比矩阵](./00-对比矩阵.md)**: 技术方案对比
- **[概念定义体系](./00-概念定义体系.md)**: 核心概念详解

---

## 🛠️ 核心技术

### 性能工具
- pprof - 性能分析
- trace - 追踪分析
- benchstat - 基准测试对比

### 分布式工具
- etcd - 分布式KV存储
- Consul - 服务发现
- Kafka - 消息队列
- Redis - 缓存/锁

### 监控工具
- Prometheus - 监控
- Grafana - 可视化
- Jaeger - 分布式追踪

---

## 📚 推荐阅读顺序

1. **性能优化** → 分析 → 内存 → 并发 → GC
2. **架构设计** → 模式 → 微服务 → 分布式
3. **安全** → Web安全 → 认证 → 授权
4. **现代技术** → AI/ML → 现代Web

---

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3
