# Go 1.23 分布式系统设计模型全面梳理

> 本文档系统性地梳理了 **Go 1.23** 语言在分布式系统中的核心设计模型，涵盖微服务架构、服务发现、负载均衡、熔断降级、分布式事务、一致性协议等11个关键领域。
>
> **Go 1.23 更新**：
>
> - 利用 `unique` 包优化服务实例标识的内存使用
> - 使用新的 `iter` 包简化服务列表遍历
> - 结合 `slices` / `maps` 迭代器优化配置管理

---

## 目录

- [Go 1.23 分布式系统设计模型全面梳理](#go-123-分布式系统设计模型全面梳理)
  - [目录](#目录)
  - [1. 微服务架构基础](#1-微服务架构基础)
    - [1.1 服务拆分原则](#11-服务拆分原则)
      - [概念定义](#概念定义)
      - [架构图](#架构图)
      - [拆分维度](#拆分维度)
      - [完整示例：电商系统拆分](#完整示例电商系统拆分)
      - [反例说明](#反例说明)
      - [注意事项](#注意事项)
      - [最佳实践](#最佳实践)
    - [1.2 服务通信方式](#12-服务通信方式)
      - [概念定义](#概念定义-1)
      - [架构图](#架构图-1)
      - [同步通信：gRPC实现](#同步通信grpc实现)
      - [异步通信：事件驱动实现](#异步通信事件驱动实现)
      - [反例说明](#反例说明-1)
      - [注意事项](#注意事项-1)
    - [1.3 API网关模式](#13-api网关模式)
      - [概念定义](#概念定义-2)
      - [架构图](#架构图-2)
      - [完整示例：Go实现API网关](#完整示例go实现api网关)
      - [高级网关：基于Kong/Envoy](#高级网关基于kongenvoy)
      - [反例说明](#反例说明-2)
      - [注意事项](#注意事项-2)
    - [1.4 BFF（Backend for Frontend）](#14-bffbackend-for-frontend)
      - [概念定义](#概念定义-3)
      - [架构图](#架构图-3)
      - [完整示例：BFF实现](#完整示例bff实现)
      - [移动端BFF优化](#移动端bff优化)
      - [反例说明](#反例说明-3)
      - [注意事项](#注意事项-3)
  - [2. 服务发现与注册](#2-服务发现与注册)
    - [2.1 客户端发现 vs 服务端发现](#21-客户端发现-vs-服务端发现)
      - [概念定义](#概念定义-4)
      - [架构图](#架构图-4)
      - [完整示例：客户端发现实现](#完整示例客户端发现实现)
      - [服务端发现：基于Nginx/Envoy](#服务端发现基于nginxenvoy)
      - [反例说明](#反例说明-4)
      - [注意事项](#注意事项-4)
    - [2.2 Consul集成](#22-consul集成)
      - [概念定义](#概念定义-5)
      - [架构图](#架构图-5)
      - [完整示例：Go集成Consul](#完整示例go集成consul)
      - [分布式锁实现](#分布式锁实现)
      - [反例说明](#反例说明-5)
      - [注意事项](#注意事项-5)
    - [2.3 etcd服务发现](#23-etcd服务发现)
      - [概念定义](#概念定义-6)
      - [架构图](#架构图-6)
      - [完整示例：Go集成etcd](#完整示例go集成etcd)
      - [分布式锁实现](#分布式锁实现-1)
      - [反例说明](#反例说明-6)
      - [注意事项](#注意事项-6)
    - [2.4 Kubernetes服务发现](#24-kubernetes服务发现)
      - [概念定义](#概念定义-7)
      - [架构图](#架构图-7)
      - [完整示例：Go应用集成K8s服务发现](#完整示例go应用集成k8s服务发现)
      - [K8s部署配置](#k8s部署配置)
      - [反例说明](#反例说明-7)
      - [注意事项](#注意事项-7)
  - [3. 负载均衡](#3-负载均衡)
    - [3.1 轮询、随机、加权负载均衡](#31-轮询随机加权负载均衡)
      - [概念定义](#概念定义-8)
      - [架构图](#架构图-8)
      - [完整示例：Go实现负载均衡算法](#完整示例go实现负载均衡算法)
      - [反例说明](#反例说明-8)
      - [注意事项](#注意事项-8)
    - [3.2 一致性哈希](#32-一致性哈希)
      - [概念定义](#概念定义-9)
      - [架构图](#架构图-9)
      - [完整示例：Go实现一致性哈希](#完整示例go实现一致性哈希)
      - [反例说明](#反例说明-9)
      - [注意事项](#注意事项-9)
    - [3.3 健康检查](#33-健康检查)
      - [概念定义](#概念定义-10)
      - [架构图](#架构图-10)
      - [完整示例：Go实现健康检查](#完整示例go实现健康检查)
      - [反例说明](#反例说明-10)
      - [注意事项](#注意事项-10)
    - [3.4 gRPC负载均衡](#34-grpc负载均衡)
      - [概念定义](#概念定义-11)
      - [架构图](#架构图-11)
      - [完整示例：Go实现gRPC负载均衡](#完整示例go实现grpc负载均衡)
      - [反例说明](#反例说明-11)
      - [注意事项](#注意事项-11)
  - [4. 熔断与降级](#4-熔断与降级)
    - [4.1 熔断器模式（Circuit Breaker）](#41-熔断器模式circuit-breaker)
      - [概念定义](#概念定义-12)
      - [架构图](#架构图-12)
      - [完整示例：Go实现熔断器](#完整示例go实现熔断器)
      - [反例说明](#反例说明-12)
      - [注意事项](#注意事项-12)
    - [4.2 降级策略](#42-降级策略)
      - [概念定义](#概念定义-13)
      - [架构图](#架构图-13)
      - [完整示例：Go实现降级策略](#完整示例go实现降级策略)
      - [反例说明](#反例说明-13)
      - [注意事项](#注意事项-13)
    - [4.3 hystrix-go实现](#43-hystrix-go实现)
      - [概念定义](#概念定义-14)
      - [完整示例：hystrix-go使用](#完整示例hystrix-go使用)
      - [反例说明](#反例说明-14)
    - [4.4 sentinel-go实现](#44-sentinel-go实现)
      - [概念定义](#概念定义-15)
      - [完整示例：sentinel-go使用](#完整示例sentinel-go使用)
      - [反例说明](#反例说明-15)
      - [注意事项](#注意事项-14)
  - [5. 限流与配额](#5-限流与配额)
    - [5.1 令牌桶算法](#51-令牌桶算法)
      - [概念定义](#概念定义-16)
      - [架构图](#架构图-14)
      - [完整示例：Go实现令牌桶](#完整示例go实现令牌桶)
      - [反例说明](#反例说明-16)
      - [注意事项](#注意事项-15)
    - [5.2 漏桶算法](#52-漏桶算法)
      - [概念定义](#概念定义-17)
      - [架构图](#架构图-15)
      - [完整示例：Go实现漏桶](#完整示例go实现漏桶)
      - [反例说明](#反例说明-17)
    - [5.3 滑动窗口](#53-滑动窗口)
      - [概念定义](#概念定义-18)
      - [架构图](#架构图-16)
      - [完整示例：Go实现滑动窗口](#完整示例go实现滑动窗口)
      - [反例说明](#反例说明-18)
    - [5.4 分布式限流](#54-分布式限流)
      - [概念定义](#概念定义-19)
      - [架构图](#架构图-17)
      - [完整示例：Go实现分布式限流](#完整示例go实现分布式限流)
      - [反例说明](#反例说明-19)
      - [注意事项](#注意事项-16)
  - [6. 重试与退避](#6-重试与退避)
    - [6.1 指数退避](#61-指数退避)
      - [概念定义](#概念定义-20)
      - [架构图](#架构图-18)
      - [完整示例：Go实现指数退避](#完整示例go实现指数退避)
      - [反例说明](#反例说明-20)
    - [6.2 抖动（Jitter）](#62-抖动jitter)
      - [概念定义](#概念定义-21)
      - [架构图](#架构图-19)
      - [完整示例：Go实现抖动](#完整示例go实现抖动)
      - [反例说明](#反例说明-21)
    - [6.3 重试策略](#63-重试策略)
      - [概念定义](#概念定义-22)
      - [完整示例：Go实现重试策略](#完整示例go实现重试策略)
      - [反例说明](#反例说明-22)
    - [6.4 幂等性设计](#64-幂等性设计)
      - [概念定义](#概念定义-23)
      - [架构图](#架构图-20)
      - [完整示例：Go实现幂等性](#完整示例go实现幂等性)
      - [反例说明](#反例说明-23)
      - [注意事项](#注意事项-17)
  - [7. 分布式事务](#7-分布式事务)
    - [7.1 两阶段提交（2PC）](#71-两阶段提交2pc)
      - [概念定义](#概念定义-24)
      - [架构图](#架构图-21)
      - [完整示例：Go实现2PC](#完整示例go实现2pc)
      - [反例说明](#反例说明-24)
      - [注意事项](#注意事项-18)
    - [7.2 三阶段提交（3PC）](#72-三阶段提交3pc)
      - [概念定义](#概念定义-25)
      - [架构图](#架构图-22)
      - [完整示例：Go实现3PC](#完整示例go实现3pc)
      - [反例说明](#反例说明-25)
    - [7.3 TCC（Try-Confirm-Cancel）](#73-tcctry-confirm-cancel)
      - [概念定义](#概念定义-26)
      - [架构图](#架构图-23)
      - [完整示例：Go实现TCC](#完整示例go实现tcc)
      - [反例说明](#反例说明-26)
      - [注意事项](#注意事项-19)
    - [7.4 Saga模式](#74-saga模式)
      - [概念定义](#概念定义-27)
      - [架构图](#架构图-24)
      - [完整示例：Go实现Saga](#完整示例go实现saga)
      - [反例说明](#反例说明-27)
    - [7.5 本地消息表](#75-本地消息表)
      - [概念定义](#概念定义-28)
      - [架构图](#架构图-25)
      - [完整示例：Go实现本地消息表](#完整示例go实现本地消息表)
      - [反例说明](#反例说明-28)
    - [7.6 最大努力通知](#76-最大努力通知)
      - [概念定义](#概念定义-29)
      - [完整示例：Go实现最大努力通知](#完整示例go实现最大努力通知)
      - [注意事项](#注意事项-20)
  - [8. 一致性协议](#8-一致性协议)
    - [8.1 Raft协议实现](#81-raft协议实现)
      - [概念定义](#概念定义-30)
      - [架构图](#架构图-26)
      - [完整示例：Go实现简化版Raft](#完整示例go实现简化版raft)
      - [反例说明](#反例说明-29)
    - [8.2 领导者选举](#82-领导者选举)
      - [概念定义](#概念定义-31)
      - [选举流程](#选举流程)
      - [注意事项](#注意事项-21)
    - [8.3 日志复制](#83-日志复制)
      - [概念定义](#概念定义-32)
      - [日志复制流程](#日志复制流程)
      - [日志匹配特性](#日志匹配特性)
    - [8.4 成员变更](#84-成员变更)
      - [概念定义](#概念定义-33)
      - [成员变更流程](#成员变更流程)
      - [注意事项](#注意事项-22)
  - [9. 分布式缓存](#9-分布式缓存)
    - [9.1 缓存穿透、击穿、雪崩](#91-缓存穿透击穿雪崩)
      - [概念定义](#概念定义-34)
      - [架构图](#架构图-27)
      - [完整示例：Go解决缓存三大问题](#完整示例go解决缓存三大问题)
      - [反例说明](#反例说明-30)
    - [9.2 缓存一致性](#92-缓存一致性)
      - [概念定义](#概念定义-35)
      - [架构图](#架构图-28)
      - [完整示例：Go实现缓存一致性](#完整示例go实现缓存一致性)
    - [9.3 本地缓存+分布式缓存](#93-本地缓存分布式缓存)
      - [完整示例：Go实现多级缓存](#完整示例go实现多级缓存)
    - [9.4 Cache-Aside模式](#94-cache-aside模式)
      - [完整示例：Go实现Cache-Aside](#完整示例go实现cache-aside)
  - [10. 分布式锁](#10-分布式锁)
    - [10.1 基于Redis的RedLock](#101-基于redis的redlock)
      - [概念定义](#概念定义-36)
      - [架构图](#架构图-29)
      - [完整示例：Go实现RedLock](#完整示例go实现redlock)
      - [反例说明](#反例说明-31)
    - [10.2 基于etcd的分布式锁](#102-基于etcd的分布式锁)
      - [概念定义](#概念定义-37)
      - [完整示例：Go实现etcd分布式锁](#完整示例go实现etcd分布式锁)
    - [10.3 基于ZooKeeper的锁](#103-基于zookeeper的锁)
      - [概念定义](#概念定义-38)
      - [完整示例：Go实现ZooKeeper分布式锁](#完整示例go实现zookeeper分布式锁)
    - [10.4 锁的续期与释放](#104-锁的续期与释放)
      - [完整示例：Go实现锁续期](#完整示例go实现锁续期)
  - [11. 消息队列](#11-消息队列)
    - [11.1 Kafka集成](#111-kafka集成)
      - [概念定义](#概念定义-39)
      - [架构图](#架构图-30)
      - [完整示例：Go集成Kafka](#完整示例go集成kafka)
      - [反例说明](#反例说明-32)
    - [11.2 RabbitMQ集成](#112-rabbitmq集成)
      - [概念定义](#概念定义-40)
      - [完整示例：Go集成RabbitMQ](#完整示例go集成rabbitmq)
    - [11.3 NATS集成](#113-nats集成)
      - [概念定义](#概念定义-41)
      - [完整示例：Go集成NATS](#完整示例go集成nats)
    - [11.4 消息可靠性保证](#114-消息可靠性保证)
      - [完整示例：Go实现消息可靠性](#完整示例go实现消息可靠性)
  - [总结](#总结)

---

## 1. 微服务架构基础

### 1.1 服务拆分原则

#### 概念定义

微服务拆分是将单体应用按照业务边界、数据边界或技术边界分解为多个独立部署的服务单元的过程。
核心原则包括：**单一职责原则(SRP)**、**高内聚低耦合**、**独立部署能力**。

#### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                      单体应用（拆分前）                        │
│  ┌─────────────┬─────────────┬─────────────┬─────────────┐  │
│  │   用户模块   │   订单模块   │   库存模块   │   支付模块   │  │
│  │  User Module│ Order Module│ Stock Module│  Pay Module │  │
│  └─────────────┴─────────────┴─────────────┴─────────────┘  │
│                    共享数据库 + 统一部署                        │
└─────────────────────────────────────────────────────────────┘
                              ↓ 拆分
┌─────────────────────────────────────────────────────────────┐
│                      微服务架构（拆分后）                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ 用户服务  │  │ 订单服务  │  │ 库存服务  │  │ 支付服务  │    │
│  │  User Svc│  │ Order Svc│  │ Stock Svc│  │  Pay Svc │    │
│  │  User DB │  │ Order DB │  │ Stock DB │  │  Pay DB  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
│       │             │             │             │            │
│       └─────────────┴──────┬──────┴─────────────┘            │
│                    轻量级通信（HTTP/gRPC）                      │
└─────────────────────────────────────────────────────────────┘
```

#### 拆分维度

| 维度 | 说明 | 示例 |
|------|------|------|
| **业务领域** | 按DDD领域边界拆分 | 用户域、订单域、商品域 |
| **数据边界** | 按数据访问模式拆分 | 读多写少分离、冷热数据分离 |
| **技术边界** | 按技术栈差异拆分 | AI服务(Python)、核心服务(Go) |
| **团队边界** | 按康威定律组织 | 2披萨团队对应2-3个服务 |

#### 完整示例：电商系统拆分

```go
// 用户服务 (user-service/main.go)
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

// User 用户实体
type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
}

// UserService 用户服务接口
type UserService interface {
    GetUser(ctx context.Context, id string) (*User, error)
    CreateUser(ctx context.Context, user *User) error
    UpdateUser(ctx context.Context, user *User) error
}

// UserHandler HTTP处理器
type UserHandler struct {
    service UserService
}

func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.service.GetUser(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, user)
}

func main() {
    r := gin.Default()
    handler := &UserHandler{service: NewUserServiceImpl()}

    r.GET("/users/:id", handler.GetUser)
    r.POST("/users", handler.CreateUser)
    r.PUT("/users/:id", handler.UpdateUser)

    log.Println("User Service starting on :8081")
    if err := r.Run(":8081"); err != nil {
        log.Fatal(err)
    }
}

// 订单服务 (order-service/main.go)
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// Order 订单实体
type Order struct {
    ID         string    `json:"id"`
    UserID     string    `json:"user_id"`
    ProductID  string    `json:"product_id"`
    Quantity   int       `json:"quantity"`
    TotalPrice float64   `json:"total_price"`
    Status     string    `json:"status"` // pending, paid, shipped, completed
    CreatedAt  time.Time `json:"created_at"`
}

// OrderService 订单服务
type OrderService struct {
    userServiceURL string // 用户服务地址
    stockServiceURL string // 库存服务地址
}

func (s *OrderService) CreateOrder(ctx context.Context, order *Order) error {
    // 1. 验证用户存在
    if err := s.validateUser(ctx, order.UserID); err != nil {
        return err
    }

    // 2. 检查库存
    if err := s.checkStock(ctx, order.ProductID, order.Quantity); err != nil {
        return err
    }

    // 3. 创建订单
    order.ID = generateOrderID()
    order.Status = "pending"
    order.CreatedAt = time.Now()

    // 保存到数据库...
    return nil
}

func (s *OrderService) validateUser(ctx context.Context, userID string) error {
    // 通过HTTP调用用户服务
    resp, err := http.Get(s.userServiceURL + "/users/" + userID)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return ErrUserNotFound
    }
    return nil
}

func (s *OrderService) checkStock(ctx context.Context, productID string, quantity int) error {
    // 调用库存服务...
    return nil
}

var ErrUserNotFound = errors.New("user not found")

func generateOrderID() string {
    return fmt.Sprintf("ORD%s", time.Now().Format("20060102150405"))
}

func main() {
    r := gin.Default()
    service := &OrderService{
        userServiceURL:  "http://localhost:8081",
        stockServiceURL: "http://localhost:8083",
    }

    r.POST("/orders", func(c *gin.Context) {
        var order Order
        if err := c.ShouldBindJSON(&order); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := service.CreateOrder(c.Request.Context(), &order); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, order)
    })

    log.Println("Order Service starting on :8082")
    r.Run(":8082")
}
```

#### 反例说明

```go
// ❌ 错误：过度拆分导致服务爆炸
// 将用户服务的每个CRUD操作都拆分成独立服务
// user-create-service, user-read-service, user-update-service, user-delete-service
// 问题：网络开销剧增、运维复杂度指数级增长、事务处理困难

// ❌ 错误：拆分粒度太粗，仍是分布式单体
// 所有业务逻辑仍耦合在一个服务中，只是部署成分布式
// 问题：没有获得微服务的好处，却承担了分布式系统的复杂度

// ❌ 错误：循环依赖
// 服务A依赖服务B，服务B又依赖服务A
// 问题：启动顺序依赖、级联故障风险、难以独立部署
```

#### 注意事项

1. **避免过早拆分**：单体应用<10万行代码时，拆分收益可能为负
2. **数据库独立**：每个微服务应有独立数据库，禁止直接访问其他服务的数据库
3. **API版本管理**：服务间API需要版本控制，保证向后兼容
4. **监控先行**：拆分前必须建立完善的监控和日志系统

#### 最佳实践

```go
// ✅ 合理的拆分粒度：每个服务可由2-3人维护
// ✅ 使用领域驱动设计(DDD)指导拆分
// ✅ 建立服务契约测试，保证接口兼容性
// ✅ 实施CI/CD，支持独立部署
```

---

### 1.2 服务通信方式

#### 概念定义

服务通信是微服务架构中服务间交互的核心机制，分为**同步通信**（请求-响应模式）和**异步通信**（事件驱动模式）两大类。

#### 架构图

```
同步通信（请求-响应）:
┌─────────┐     HTTP/gRPC      ┌─────────┐
│  服务A   │ ─────────────────> │  服务B   │
│         │ <───────────────── │         │
└─────────┘     阻塞等待响应    └─────────┘

异步通信（事件驱动）:
┌─────────┐    发布事件    ┌─────────┐
│  服务A   │ ────────────> │  消息队列  │
│         │               │ (Kafka)  │
└─────────┘               └────┬────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              ▼                ▼                ▼
         ┌─────────┐     ┌─────────┐     ┌─────────┐
         │  服务B   │     │  服务C   │     │  服务D   │
         │ (订阅者) │     │ (订阅者) │     │ (订阅者) │
         └─────────┘     └─────────┘     └─────────┘
```

#### 同步通信：gRPC实现

```go
// 定义protobuf (api/order.proto)
syntax = "proto3";
package order;
option go_package = "./orderpb";

service OrderService {
    rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
    rpc GetOrder(GetOrderRequest) returns (Order);
    rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (Order);
}

message Order {
    string id = 1;
    string user_id = 2;
    repeated OrderItem items = 3;
    double total_amount = 4;
    string status = 5;
}

message OrderItem {
    string product_id = 1;
    int32 quantity = 2;
    double price = 3;
}

message CreateOrderRequest {
    string user_id = 1;
    repeated OrderItem items = 2;
}

message CreateOrderResponse {
    string order_id = 1;
    string status = 2;
}

// 服务端实现 (order-service/grpc_server.go)
package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type OrderServer struct {
    orderpb.UnimplementedOrderServiceServer
    repo OrderRepository
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
    // 参数校验
    if req.UserId == "" {
        return nil, status.Error(codes.InvalidArgument, "user_id is required")
    }
    if len(req.Items) == 0 {
        return nil, status.Error(codes.InvalidArgument, "items cannot be empty")
    }

    // 计算总金额
    var totalAmount float64
    for _, item := range req.Items {
        totalAmount += float64(item.Quantity) * item.Price
    }

    // 创建订单
    order := &Order{
        ID:          generateID(),
        UserID:      req.UserId,
        Items:       convertItems(req.Items),
        TotalAmount: totalAmount,
        Status:      "pending",
    }

    if err := s.repo.Save(ctx, order); err != nil {
        log.Printf("failed to save order: %v", err)
        return nil, status.Error(codes.Internal, "failed to create order")
    }

    return &orderpb.CreateOrderResponse{
        OrderId: order.ID,
        Status:  order.Status,
    }, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.Order, error) {
    order, err := s.repo.GetByID(ctx, req.OrderId)
    if err != nil {
        if errors.Is(err, ErrOrderNotFound) {
            return nil, status.Error(codes.NotFound, "order not found")
        }
        return nil, status.Error(codes.Internal, "failed to get order")
    }
    return convertToProto(order), nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    server := grpc.NewServer(
        grpc.UnaryInterceptor(loggingInterceptor),
        grpc.UnaryInterceptor(recoveryInterceptor),
    )
    orderpb.RegisterOrderServiceServer(server, &OrderServer{repo: NewOrderRepository()})

    log.Println("gRPC server starting on :50051")
    if err := server.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

// 拦截器实现
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    log.Printf("method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
    return resp, err
}

func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("panic recovered: %v", r)
            err = status.Error(codes.Internal, "internal server error")
        }
    }()
    return handler(ctx, req)
}

// 客户端实现 (client/order_client.go)
package client

import (
    "context"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/keepalive"
)

type OrderClient struct {
    conn   *grpc.ClientConn
    client orderpb.OrderServiceClient
}

func NewOrderClient(addr string) (*OrderClient, error) {
    kacp := keepalive.ClientParameters{
        Time:                10 * time.Second,
        Timeout:             time.Second,
        PermitWithoutStream: true,
    }

    conn, err := grpc.Dial(addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithKeepaliveParams(kacp),
        grpc.WithDefaultServiceConfig(`{
            "loadBalancingConfig": [{"round_robin": {}}],
            "healthCheckConfig": {"serviceName": ""}
        }`),
    )
    if err != nil {
        return nil, err
    }

    return &OrderClient{
        conn:   conn,
        client: orderpb.NewOrderServiceClient(conn),
    }, nil
}

func (c *OrderClient) CreateOrder(ctx context.Context, userID string, items []*orderpb.OrderItem) (*orderpb.CreateOrderResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    return c.client.CreateOrder(ctx, &orderpb.CreateOrderRequest{
        UserId: userID,
        Items:  items,
    })
}

func (c *OrderClient) Close() error {
    return c.conn.Close()
}
```

#### 异步通信：事件驱动实现

```go
// 事件定义 (events/order_events.go)
package events

import (
    "encoding/json"
    "time"
)

// EventType 事件类型
type EventType string

const (
    OrderCreated   EventType = "order.created"
    OrderPaid      EventType = "order.paid"
    OrderShipped   EventType = "order.shipped"
    OrderCompleted EventType = "order.completed"
    OrderCancelled EventType = "order.cancelled"
)

// Event 领域事件
type Event struct {
    ID        string          `json:"id"`
    Type      EventType       `json:"type"`
    Timestamp time.Time       `json:"timestamp"`
    Payload   json.RawMessage `json:"payload"`
}

// OrderCreatedEvent 订单创建事件
type OrderCreatedEvent struct {
    OrderID     string    `json:"order_id"`
    UserID      string    `json:"user_id"`
    TotalAmount float64   `json:"total_amount"`
    CreatedAt   time.Time `json:"created_at"`
}

// EventPublisher 事件发布者接口
type EventPublisher interface {
    Publish(ctx context.Context, event *Event) error
    Close() error
}

// EventSubscriber 事件订阅者接口
type EventSubscriber interface {
    Subscribe(eventType EventType, handler EventHandler) error
    Start(ctx context.Context) error
    Close() error
}

// EventHandler 事件处理器类型
type EventHandler func(ctx context.Context, event *Event) error

// Kafka实现 (messaging/kafka_publisher.go)
package messaging

import (
    "context"
    "encoding/json"
    "log"

    "github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
    writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
    return &KafkaPublisher{
        writer: &kafka.Writer{
            Addr:     kafka.TCP(brokers...),
            Topic:    topic,
            Balancer: &kafka.LeastBytes{},
        },
    }
}

func (p *KafkaPublisher) Publish(ctx context.Context, event *events.Event) error {
    value, err := json.Marshal(event)
    if err != nil {
        return err
    }

    msg := kafka.Message{
        Key:   []byte(event.ID),
        Value: value,
        Headers: []kafka.Header{
            {Key: "event_type", Value: []byte(event.Type)},
        },
    }

    return p.writer.WriteMessages(ctx, msg)
}

func (p *KafkaPublisher) Close() error {
    return p.writer.Close()
}

// 订单服务发布事件 (order-service/event_publisher.go)
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
    // 1. 创建订单...
    order, err := s.createOrderInDB(ctx, req)
    if err != nil {
        return nil, err
    }

    // 2. 发布订单创建事件
    event := &events.OrderCreatedEvent{
        OrderID:     order.ID,
        UserID:      order.UserID,
        TotalAmount: order.TotalAmount,
        CreatedAt:   order.CreatedAt,
    }

    payload, _ := json.Marshal(event)
    domainEvent := &events.Event{
        ID:        generateEventID(),
        Type:      events.OrderCreated,
        Timestamp: time.Now(),
        Payload:   payload,
    }

    if err := s.eventPublisher.Publish(ctx, domainEvent); err != nil {
        // 记录日志，但不影响订单创建（最终一致性）
        log.Printf("failed to publish event: %v", err)
    }

    return order, nil
}

// 库存服务订阅事件 (inventory-service/event_handler.go)
func (s *InventoryService) HandleOrderCreated(ctx context.Context, event *events.Event) error {
    var orderEvent events.OrderCreatedEvent
    if err := json.Unmarshal(event.Payload, &orderEvent); err != nil {
        return err
    }

    // 预留库存
    if err := s.ReserveStock(ctx, orderEvent.OrderID, orderEvent.Items); err != nil {
        // 发布库存预留失败事件，触发补偿
        return s.publishStockReservationFailed(ctx, orderEvent.OrderID)
    }

    return nil
}

func main() {
    subscriber := messaging.NewKafkaSubscriber(kafkaBrokers, "order-events")
    inventoryService := NewInventoryService()

    // 注册事件处理器
    subscriber.Subscribe(events.OrderCreated, inventoryService.HandleOrderCreated)
    subscriber.Subscribe(events.OrderCancelled, inventoryService.HandleOrderCancelled)

    ctx := context.Background()
    subscriber.Start(ctx)
}
```

#### 反例说明

```go
// ❌ 错误：同步调用链过长
// A -> B -> C -> D -> E，任何一个服务故障都会导致整个链路失败
// 应该使用异步解耦或引入缓存

// ❌ 错误：混合使用多种通信协议无统一规划
// 服务间有的用HTTP、有的用gRPC、有的用消息队列
// 问题：增加学习和维护成本

// ❌ 错误：异步消息无顺序保证
// 订单创建事件和支付事件乱序到达
// 问题：库存状态不一致
```

#### 注意事项

1. **同步调用超时设置**：必须设置合理的超时时间，避免级联阻塞
2. **异步消息幂等性**：消费者必须实现幂等，防止重复处理
3. **消息顺序保证**：需要顺序的场景使用分区键或单分区
4. **死信队列**：处理失败的消息应进入死信队列，避免无限重试

---

### 1.3 API网关模式

#### 概念定义

API网关是微服务架构的统一入口，负责**请求路由**、**协议转换**、**认证授权**、**限流熔断**、**日志监控**等横切关注点。

#### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         API Gateway                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   路由模块   │  │   认证模块   │  │   限流模块   │         │
│  │   Routing   │  │    Auth     │  │ Rate Limiter│         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   熔断模块   │  │   日志模块   │  │   协议转换   │         │
│  │   Circuit   │  │   Logging   │  │  Protocol   │         │
│  │   Breaker   │  │             │  │   Convert   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
   ┌─────────┐          ┌─────────┐          ┌─────────┐
   │ 用户服务 │          │ 订单服务 │          │ 库存服务 │
   │ :8081   │          │ :8082   │          │ :8083   │
   └─────────┘          └─────────┘          └─────────┘
```

#### 完整示例：Go实现API网关

```go
// gateway/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

// Route 路由配置
type Route struct {
    Path         string
    ServiceURL   string
    StripPrefix  bool
    AuthRequired bool
    RateLimit    int // 每秒请求数
}

// Gateway API网关
type Gateway struct {
    routes      map[string]*Route
    proxies     map[string]*httputil.ReverseProxy
    limiters    map[string]*rate.Limiter
    authService AuthService
    mu          sync.RWMutex
}

func NewGateway() *Gateway {
    return &Gateway{
        routes:   make(map[string]*Route),
        proxies:  make(map[string]*httputil.ReverseProxy),
        limiters: make(map[string]*rate.Limiter),
    }
}

func (g *Gateway) RegisterRoute(route *Route) error {
    target, err := url.Parse(route.ServiceURL)
    if err != nil {
        return err
    }

    g.mu.Lock()
    defer g.mu.Unlock()

    g.routes[route.Path] = route
    g.proxies[route.Path] = httputil.NewSingleHostReverseProxy(target)

    if route.RateLimit > 0 {
        g.limiters[route.Path] = rate.NewLimiter(rate.Limit(route.RateLimit), route.RateLimit*2)
    }

    return nil
}

func (g *Gateway) Handler() gin.HandlerFunc {
    return func(c *gin.Context) {
        path := c.Request.URL.Path

        // 查找匹配的路由
        route, proxy := g.findRoute(path)
        if route == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
            c.Abort()
            return
        }

        // 认证检查
        if route.AuthRequired {
            if err := g.authenticate(c); err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
                c.Abort()
                return
            }
        }

        // 限流检查
        if limiter := g.limiters[route.Path]; limiter != nil {
            if !limiter.Allow() {
                c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
                c.Abort()
                return
            }
        }

        // 路径重写
        if route.StripPrefix {
            c.Request.URL.Path = strings.TrimPrefix(path, route.Path)
        }

        // 记录请求日志
        start := time.Now()
        defer func() {
            log.Printf("[%s] %s %s %d %v",
                c.Request.Method,
                path,
                c.ClientIP(),
                c.Writer.Status(),
                time.Since(start),
            )
        }()

        // 代理请求
        proxy.ServeHTTP(c.Writer, c.Request)
    }
}

func (g *Gateway) findRoute(path string) (*Route, *httputil.ReverseProxy) {
    g.mu.RLock()
    defer g.mu.RUnlock()

    for prefix, route := range g.routes {
        if strings.HasPrefix(path, prefix) {
            return route, g.proxies[prefix]
        }
    }
    return nil, nil
}

func (g *Gateway) authenticate(c *gin.Context) error {
    token := c.GetHeader("Authorization")
    if token == "" {
        return ErrMissingToken
    }

    // 调用认证服务验证token
    // 这里简化处理，实际应调用认证服务
    if !strings.HasPrefix(token, "Bearer ") {
        return ErrInvalidToken
    }

    return nil
}

var (
    ErrMissingToken = errors.New("missing authorization token")
    ErrInvalidToken = errors.New("invalid authorization token")
)

// 中间件实现

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        c.Next()

        if raw != "" {
            path = path + "?" + raw
        }

        log.Printf("[%s] %s | %d | %s | %s",
            c.Request.Method,
            path,
            c.Writer.Status(),
            time.Since(start),
            c.Errors.String(),
        )
    }
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "internal server error",
                })
            }
        }()
        c.Next()
    }
}

func main() {
    gateway := NewGateway()

    // 注册路由
    gateway.RegisterRoute(&Route{
        Path:         "/api/users",
        ServiceURL:   "http://localhost:8081",
        StripPrefix:  false,
        AuthRequired: true,
        RateLimit:    100,
    })

    gateway.RegisterRoute(&Route{
        Path:         "/api/orders",
        ServiceURL:   "http://localhost:8082",
        StripPrefix:  false,
        AuthRequired: true,
        RateLimit:    50,
    })

    gateway.RegisterRoute(&Route{
        Path:         "/api/public",
        ServiceURL:   "http://localhost:8084",
        StripPrefix:  false,
        AuthRequired: false,
        RateLimit:    200,
    })

    r := gin.New()
    r.Use(RecoveryMiddleware())
    r.Use(LoggingMiddleware())
    r.Use(CORSMiddleware())
    r.Use(gateway.Handler())

    // 健康检查
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    log.Println("API Gateway starting on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

#### 高级网关：基于Kong/Envoy

```go
// 使用go-control-plane与Envoy集成
package main

import (
    "context"
    "log"
    "net"

    cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
    listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
    route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
    "github.com/envoyproxy/go-control-plane/pkg/cache/types"
    "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
    "github.com/envoyproxy/go-control-plane/pkg/server/v3"
    "github.com/envoyproxy/go-control-plane/pkg/test/resource/v3"
)

// XDS服务器实现服务发现
func runXDSServer(ctx context.Context, snapshotCache cache.SnapshotCache) {
    grpcServer := grpc.NewServer()
    server := server.NewServer(ctx, snapshotCache, nil)

    // 注册xDS服务
    discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
    endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
    clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, server)
    routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, server)
    listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, server)

    lis, _ := net.Listen("tcp", ":18000")
    go func() {
        if err := grpcServer.Serve(lis); err != nil {
            log.Printf("xDS server error: %v", err)
        }
    }()

    <-ctx.Done()
    grpcServer.GracefulStop()
}
```

#### 反例说明

```go
// ❌ 错误：网关承担业务逻辑
// 在网关层实现订单计算、库存检查等业务逻辑
// 问题：网关臃肿、难以维护、违背单一职责

// ❌ 错误：无状态网关却存储会话
// 网关节点间会话不共享，导致用户请求漂移后认证失效
// 问题：需要使用Redis等共享存储会话

// ❌ 错误：网关层不做超时控制
// 后端服务响应慢导致网关连接耗尽
// 问题：雪崩效应
```

#### 注意事项

1. **网关无状态化**：网关应设计为无状态，支持水平扩展
2. **超时配置**：为每个路由配置合理的超时时间
3. **健康检查**：定期检测后端服务健康状态
4. **灰度发布**：支持基于Header或权重的流量分发

---

### 1.4 BFF（Backend for Frontend）

#### 概念定义

BFF模式为每种前端（Web、iOS、Android）提供定制化的后端服务，解决**多端数据需求差异**、**减少前端复杂度**、**优化网络请求**。

#### 架构图

```
传统模式（前端直接调用多个服务）:
┌─────────┐      ┌─────────┐
│  Web App │      │ Mobile  │
└────┬────┘      └────┬────┘
     │                │
     ├───── 用户服务 ──┤
     ├───── 订单服务 ──┤
     ├───── 库存服务 ──┤
     ├───── 支付服务 ──┤
     └───── 物流服务 ──┘

BFF模式（前端调用单一BFF）:
┌─────────┐      ┌─────────┐      ┌─────────┐
│  Web App │      │ iOS App │      │Android  │
└────┬────┘      └────┬────┘      └────┬────┘
     │                │                │
┌────┴────┐      ┌────┴────┐      ┌────┴────┐
│ Web BFF │      │ iOS BFF │      │Android  │
│         │      │         │      │  BFF    │
└────┬────┘      └────┬────┘      └────┬────┘
     │                │                │
     └────────────────┼────────────────┘
                      │
         ┌────────────┼────────────┐
         ▼            ▼            ▼
    ┌─────────┐  ┌─────────┐  ┌─────────┐
    │ 用户服务  │  │ 订单服务  │  │ 库存服务  │
    └─────────┘  └─────────┘  └─────────┘
```

#### 完整示例：BFF实现

```go
// web-bff/main.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

// WebBFF Web端BFF服务
type WebBFF struct {
    userClient   *UserClient
    orderClient  *OrderClient
    stockClient  *StockClient
    httpClient   *http.Client
}

// DashboardResponse 首页仪表盘数据（聚合多个服务）
type DashboardResponse struct {
    User      *UserInfo      `json:"user"`
    Orders    []OrderSummary `json:"orders"`
    CartCount int            `json:"cart_count"`
    Coupons   []Coupon       `json:"coupons"`
}

// UserInfo 用户信息
type UserInfo struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Avatar   string `json:"avatar"`
    Level    string `json:"level"`
    Points   int    `json:"points"`
}

// OrderSummary 订单摘要
type OrderSummary struct {
    ID          string  `json:"id"`
    Status      string  `json:"status"`
    TotalAmount float64 `json:"total_amount"`
    ItemCount   int     `json:"item_count"`
}

// GetDashboard 获取首页仪表盘数据（并行聚合）
func (b *WebBFF) GetDashboard(c *gin.Context) {
    userID := c.GetString("user_id") // 从JWT解析

    ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

    var (
        user      *UserInfo
        orders    []OrderSummary
        cartCount int
        coupons   []Coupon
        err       error
    )

    // 使用errgroup并行调用多个服务
    g, ctx := errgroup.WithContext(ctx)

    // 获取用户信息
    g.Go(func() error {
        user, err = b.userClient.GetUser(ctx, userID)
        return err
    })

    // 获取订单列表
    g.Go(func() error {
        orders, err = b.orderClient.GetRecentOrders(ctx, userID, 5)
        return err
    })

    // 获取购物车数量
    g.Go(func() error {
        cartCount, err = b.stockClient.GetCartItemCount(ctx, userID)
        return err
    })

    // 获取优惠券
    g.Go(func() error {
        coupons, err = b.userClient.GetCoupons(ctx, userID)
        return err
    })

    if err := g.Wait(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, DashboardResponse{
        User:      user,
        Orders:    orders,
        CartCount: cartCount,
        Coupons:   coupons,
    })
}

// ProductDetailResponse 商品详情页响应
type ProductDetailResponse struct {
    Product      *Product      `json:"product"`
    RelatedProducts []Product  `json:"related_products"`
    Reviews      []Review      `json:"reviews"`
    StockStatus  string        `json:"stock_status"`
    IsFavorite   bool          `json:"is_favorite"`
}

// GetProductDetail 获取商品详情（针对Web端优化）
func (b *WebBFF) GetProductDetail(c *gin.Context) {
    productID := c.Param("id")
    userID := c.GetString("user_id")

    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    resp := &ProductDetailResponse{}
    var wg sync.WaitGroup
    errChan := make(chan error, 4)

    // 并行获取商品信息
    wg.Add(1)
    go func() {
        defer wg.Done()
        product, err := b.stockClient.GetProduct(ctx, productID)
        if err != nil {
            errChan <- err
            return
        }
        resp.Product = product
    }()

    // 并行获取相关商品
    wg.Add(1)
    go func() {
        defer wg.Done()
        related, err := b.stockClient.GetRelatedProducts(ctx, productID, 4)
        if err != nil {
            errChan <- err
            return
        }
        resp.RelatedProducts = related
    }()

    // 并行获取评价
    wg.Add(1)
    go func() {
        defer wg.Done()
        reviews, err := b.stockClient.GetProductReviews(ctx, productID, 1, 10)
        if err != nil {
            errChan <- err
            return
        }
        resp.Reviews = reviews
    }()

    // 并行获取库存状态
    wg.Add(1)
    go func() {
        defer wg.Done()
        stock, err := b.stockClient.GetStockStatus(ctx, productID)
        if err != nil {
            errChan <- err
            return
        }
        resp.StockStatus = stock.Status
    }()

    // 检查是否收藏（仅登录用户）
    if userID != "" {
        wg.Add(1)
        go func() {
            defer wg.Done()
            isFav, err := b.userClient.IsFavorite(ctx, userID, productID)
            if err != nil {
                errChan <- err
                return
            }
            resp.IsFavorite = isFav
        }()
    }

    wg.Wait()
    close(errChan)

    // 检查是否有错误
    for err := range errChan {
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }

    c.JSON(http.StatusOK, resp)
}

// 客户端封装

type UserClient struct {
    baseURL string
    client  *http.Client
}

func NewUserClient(baseURL string) *UserClient {
    return &UserClient{
        baseURL: baseURL,
        client: &http.Client{
            Timeout: 2 * time.Second,
        },
    }
}

func (c *UserClient) GetUser(ctx context.Context, userID string) (*UserInfo, error) {
    req, err := http.NewRequestWithContext(ctx, "GET",
        fmt.Sprintf("%s/users/%s", c.baseURL, userID), nil)
    if err != nil {
        return nil, err
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("user service returned %d", resp.StatusCode)
    }

    var user UserInfo
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    return &user, nil
}

func (c *UserClient) GetCoupons(ctx context.Context, userID string) ([]Coupon, error) {
    // 实现...
    return nil, nil
}

func (c *UserClient) IsFavorite(ctx context.Context, userID, productID string) (bool, error) {
    // 实现...
    return false, nil
}

type OrderClient struct {
    baseURL string
    client  *http.Client
}

func (c *OrderClient) GetRecentOrders(ctx context.Context, userID string, limit int) ([]OrderSummary, error) {
    // 实现...
    return nil, nil
}

type StockClient struct {
    baseURL string
    client  *http.Client
}

func (c *StockClient) GetProduct(ctx context.Context, productID string) (*Product, error) {
    // 实现...
    return nil, nil
}

func (c *StockClient) GetRelatedProducts(ctx context.Context, productID string, limit int) ([]Product, error) {
    // 实现...
    return nil, nil
}

func (c *StockClient) GetProductReviews(ctx context.Context, productID string, page, size int) ([]Review, error) {
    // 实现...
    return nil, nil
}

func (c *StockClient) GetStockStatus(ctx context.Context, productID string) (*StockStatus, error) {
    // 实现...
    return nil, nil
}

func (c *StockClient) GetCartItemCount(ctx context.Context, userID string) (int, error) {
    // 实现...
    return 0, nil
}

// 数据结构定义
type Product struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Images      []string `json:"images"`
}

type Review struct {
    ID      string `json:"id"`
    UserID  string `json:"user_id"`
    Rating  int    `json:"rating"`
    Content string `json:"content"`
}

type StockStatus struct {
    ProductID string `json:"product_id"`
    Status    string `json:"status"`
    Quantity  int    `json:"quantity"`
}

type Coupon struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Discount float64 `json:"discount"`
}

func main() {
    bff := &WebBFF{
        userClient:  NewUserClient("http://user-service:8081"),
        orderClient: &OrderClient{baseURL: "http://order-service:8082"},
        stockClient: &StockClient{baseURL: "http://stock-service:8083"},
        httpClient:  &http.Client{Timeout: 5 * time.Second},
    }

    r := gin.Default()

    // 首页仪表盘
    r.GET("/api/dashboard", bff.GetDashboard)

    // 商品详情
    r.GET("/api/products/:id", bff.GetProductDetail)

    log.Println("Web BFF starting on :8090")
    r.Run(":8090")
}
```

#### 移动端BFF优化

```go
// mobile-bff/main.go
// 移动端BFF针对弱网环境优化

// 数据压缩和精简
type MobileProduct struct {
    ID       string  `json:"id"`
    Name     string  `json:"n"`  // 字段名缩短
    Price    float64 `json:"p"`
    Image    string  `json:"i"`  // 只返回一张主图
    InStock  bool    `json:"s"`
}

// 批量接口减少请求次数
func (b *MobileBFF) GetHomeData(c *gin.Context) {
    // 移动端首页需要的数据一次性返回
    // 避免多次网络请求
    type HomeData struct {
        Banners      []Banner      `json:"banners"`
        Categories   []Category    `json:"categories"`
        HotProducts  []MobileProduct `json:"hot_products"`
        FlashSale    *FlashSale    `json:"flash_sale,omitempty"`
    }

    // 并行获取所有数据
    // ...
}

// 图片适配
func adaptImageForMobile(imageURL string, deviceType string) string {
    // 根据设备类型返回不同尺寸的图片
    // 移动端使用压缩后的图片
    switch deviceType {
    case "ios":
        return imageURL + "?w=750&q=80"
    case "android":
        return imageURL + "?w=720&q=75"
    default:
        return imageURL
    }
}
```

#### 反例说明

```go
// ❌ 错误：BFF层实现业务逻辑
// BFF应该只做数据聚合和适配，不应该有业务规则

// ❌ 错误：所有前端共享一个BFF
// 失去了BFF的意义，前端仍需处理不必要的数据

// ❌ 错误：BFF直接访问数据库
// BFF应该通过API调用后端服务，禁止直接访问数据库
```

#### 注意事项

1. **BFF层保持薄**：只负责数据聚合、格式转换，不实现业务逻辑
2. **统一错误处理**：BFF应对下游服务的错误进行统一包装
3. **缓存策略**：对聚合结果进行合理缓存，减少下游调用
4. **超时控制**：聚合多个服务时，设置合理的超时时间

---

## 2. 服务发现与注册

### 2.1 客户端发现 vs 服务端发现

#### 概念定义

服务发现是分布式系统中服务实例动态定位的机制。**客户端发现**由客户端直接查询注册中心获取服务地址；**服务端发现**通过负载均衡器或Sidecar代理统一处理服务发现。

#### 架构图

```
客户端服务发现模式:
┌─────────┐     1.查询服务地址      ┌─────────────┐
│  客户端  │ ────────────────────> │  注册中心    │
│ (Client)│ <──────────────────── │ (Consul/etcd)│
└────┬────│     2.返回实例列表      └─────────────┘
     │
     │ 3.直接调用
     ▼
┌─────────┐
│ 服务实例A │
│ :8081   │
└─────────┘

服务端服务发现模式:
┌─────────┐     1.调用服务名        ┌─────────────┐
│  客户端  │ ────────────────────> │  负载均衡器  │
│ (Client)│                       │  (LB/Proxy) │
└─────────┘                       └──────┬──────┘
                                         │
                    2.查询并选择实例      │
                    ┌────────────────────┘
                    ▼
              ┌─────────────┐
              │  注册中心    │
              │ (Consul/etcd)│
              └─────────────┘
                    │
     ┌──────────────┼──────────────┐
     ▼              ▼              ▼
┌─────────┐   ┌─────────┐   ┌─────────┐
│ 服务实例A │   │ 服务实例B │   │ 服务实例C │
│ :8081   │   │ :8082   │   │ :8083   │
└─────────┘   └─────────┘   └─────────┘
```

#### 完整示例：客户端发现实现

```go
// discovery/client_discovery.go
package discovery

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// ServiceInstance 服务实例信息
type ServiceInstance struct {
    ID       string            `json:"id"`
    Name     string            `json:"name"`
    Host     string            `json:"host"`
    Port     int               `json:"port"`
    Metadata map[string]string `json:"metadata"`
    Healthy  bool              `json:"healthy"`
    Weight   int               `json:"weight"`
}

func (s *ServiceInstance) Address() string {
    return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Registry 服务注册中心接口
type Registry interface {
    Register(ctx context.Context, instance *ServiceInstance) error
    Deregister(ctx context.Context, instanceID string) error
    Discover(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
    Watch(ctx context.Context, serviceName string) (chan []*ServiceInstance, error)
    Close() error
}

// ClientDiscovery 客户端服务发现
type ClientDiscovery struct {
    registry   Registry
    cache      map[string][]*ServiceInstance
    watchers   map[string]chan struct{}
    mu         sync.RWMutex
    refreshInterval time.Duration
}

func NewClientDiscovery(registry Registry) *ClientDiscovery {
    cd := &ClientDiscovery{
        registry:        registry,
        cache:           make(map[string][]*ServiceInstance),
        watchers:        make(map[string]chan struct{}),
        refreshInterval: 10 * time.Second,
    }
    go cd.refreshLoop()
    return cd
}

// GetInstances 获取服务实例列表（优先从缓存读取）
func (cd *ClientDiscovery) GetInstances(ctx context.Context, serviceName string) ([]*ServiceInstance, error) {
    // 先查缓存
    cd.mu.RLock()
    if instances, ok := cd.cache[serviceName]; ok && len(instances) > 0 {
        cd.mu.RUnlock()
        return cd.filterHealthy(instances), nil
    }
    cd.mu.RUnlock()

    // 缓存未命中，从注册中心获取
    instances, err := cd.registry.Discover(ctx, serviceName)
    if err != nil {
        return nil, err
    }

    // 更新缓存
    cd.mu.Lock()
    cd.cache[serviceName] = instances
    cd.mu.Unlock()

    return cd.filterHealthy(instances), nil
}

// filterHealthy 过滤健康实例
func (cd *ClientDiscovery) filterHealthy(instances []*ServiceInstance) []*ServiceInstance {
    var healthy []*ServiceInstance
    for _, inst := range instances {
        if inst.Healthy {
            healthy = append(healthy, inst)
        }
    }
    return healthy
}

// refreshLoop 定期刷新缓存
func (cd *ClientDiscovery) refreshLoop() {
    ticker := time.NewTicker(cd.refreshInterval)
    defer ticker.Stop()

    for range ticker.C {
        cd.mu.RLock()
        services := make([]string, 0, len(cd.cache))
        for name := range cd.cache {
            services = append(services, name)
        }
        cd.mu.RUnlock()

        for _, serviceName := range services {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            instances, err := cd.registry.Discover(ctx, serviceName)
            cancel()

            if err != nil {
                continue
            }

            cd.mu.Lock()
            cd.cache[serviceName] = instances
            cd.mu.Unlock()
        }
    }
}

// WatchService 监听服务变化
func (cd *ClientDiscovery) WatchService(ctx context.Context, serviceName string) error {
    ch, err := cd.registry.Watch(ctx, serviceName)
    if err != nil {
        return err
    }

    go func() {
        for {
            select {
            case instances := <-ch:
                cd.mu.Lock()
                cd.cache[serviceName] = instances
                cd.mu.Unlock()
            case <-ctx.Done():
                return
            }
        }
    }()

    return nil
}

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
    Select(instances []*ServiceInstance) (*ServiceInstance, error)
}

// RoundRobinBalancer 轮询负载均衡器
type RoundRobinBalancer struct {
    counter uint64
    mu      sync.Mutex
}

func (r *RoundRobinBalancer) Select(instances []*ServiceInstance) (*ServiceInstance, error) {
    if len(instances) == 0 {
        return nil, fmt.Errorf("no available instances")
    }

    r.mu.Lock()
    idx := r.counter % uint64(len(instances))
    r.counter++
    r.mu.Unlock()

    return instances[idx], nil
}

// Client 服务发现客户端
type Client struct {
    discovery *ClientDiscovery
    balancer  LoadBalancer
}

func NewClient(registry Registry, balancer LoadBalancer) *Client {
    return &Client{
        discovery: NewClientDiscovery(registry),
        balancer:  balancer,
    }
}

// GetService 获取一个可用的服务实例
func (c *Client) GetService(ctx context.Context, serviceName string) (*ServiceInstance, error) {
    instances, err := c.discovery.GetInstances(ctx, serviceName)
    if err != nil {
        return nil, err
    }

    return c.balancer.Select(instances)
}
```

#### 服务端发现：基于Nginx/Envoy

```go
// discovery/server_discovery.go
package discovery

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// ServerSideRegistration 服务端注册（Sidecar模式）
type ServerSideRegistration struct {
    registry Registry
    instance *ServiceInstance
    interval time.Duration
    stopChan chan struct{}
}

func NewServerSideRegistration(registry Registry, instance *ServiceInstance, interval time.Duration) *ServerSideRegistration {
    return &ServerSideRegistration{
        registry: registry,
        instance: instance,
        interval: interval,
        stopChan: make(chan struct{}),
    }
}

// Start 启动服务注册和心跳
func (s *ServerSideRegistration) Start() error {
    // 初始注册
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := s.registry.Register(ctx, s.instance); err != nil {
        return fmt.Errorf("failed to register service: %w", err)
    }

    // 启动心跳
    go s.heartbeatLoop()

    // 监听退出信号
    go s.handleShutdown()

    return nil
}

func (s *ServerSideRegistration) heartbeatLoop() {
    ticker := time.NewTicker(s.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            // 更新健康状态
            s.instance.Healthy = true
            err := s.registry.Register(ctx, s.instance)
            cancel()

            if err != nil {
                fmt.Printf("heartbeat failed: %v\n", err)
            }
        case <-s.stopChan:
            return
        }
    }
}

func (s *ServerSideRegistration) handleShutdown() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan
    close(s.stopChan)

    // 注销服务
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := s.registry.Deregister(ctx, s.instance.ID); err != nil {
        fmt.Printf("failed to deregister: %v\n", err)
    }

    os.Exit(0)
}

// HealthCheckServer 健康检查服务器
type HealthCheckServer struct {
    port    int
    checks  map[string]HealthChecker
    mu      sync.RWMutex
}

type HealthChecker func() error

func NewHealthCheckServer(port int) *HealthCheckServer {
    return &HealthCheckServer{
        port:   port,
        checks: make(map[string]HealthChecker),
    }
}

func (h *HealthCheckServer) RegisterCheck(name string, checker HealthChecker) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.checks[name] = checker
}

func (h *HealthCheckServer) Start() {
    http.HandleFunc("/health", h.healthHandler)
    http.HandleFunc("/ready", h.readyHandler)

    go http.ListenAndServe(fmt.Sprintf(":%d", h.port), nil)
}

func (h *HealthCheckServer) healthHandler(w http.ResponseWriter, r *http.Request) {
    h.mu.RLock()
    checks := make(map[string]HealthChecker)
    for k, v := range h.checks {
        checks[k] = v
    }
    h.mu.RUnlock()

    status := make(map[string]string)
    healthy := true

    for name, checker := range checks {
        if err := checker(); err != nil {
            status[name] = "unhealthy: " + err.Error()
            healthy = false
        } else {
            status[name] = "healthy"
        }
    }

    response := map[string]interface{}{
        "status":  status,
        "healthy": healthy,
    }

    if !healthy {
        w.WriteHeader(http.StatusServiceUnavailable)
    }

    json.NewEncoder(w).Encode(response)
}

func (h *HealthCheckServer) readyHandler(w http.ResponseWriter, r *http.Request) {
    //  readiness检查，判断服务是否准备好接收流量
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]bool{"ready": true})
}
```

#### 反例说明

```go
// ❌ 错误：硬编码服务地址
const UserServiceURL = "http://192.168.1.100:8081"
// 问题：服务迁移时需要修改代码重新部署

// ❌ 错误：无健康检查
// 客户端继续向已故障的实例发送请求
// 问题：请求失败率高

// ❌ 错误：缓存不刷新
// 服务实例变化后客户端仍使用旧缓存
// 问题：请求发送到已下线的实例
```

#### 注意事项

1. **心跳机制**：服务需要定期发送心跳，注册中心超时剔除失效实例
2. **优雅退出**：服务关闭时应主动注销，避免请求丢失
3. **缓存一致性**：客户端缓存需要设置合理的过期时间或监听变更
4. **网络分区处理**：考虑脑裂情况下的服务发现策略

---

### 2.2 Consul集成

#### 概念定义

Consul是HashiCorp开源的服务网格解决方案，提供服务发现、健康检查、键值存储、安全服务通信等功能。

#### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        Consul集群                            │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                     │
│  │ Server  │  │ Server  │  │ Server  │  Raft共识            │
│  │ (Leader)│  │(Follower)│  │(Follower)│                     │
│  │ :8300  │  │ :8300  │  │ :8300  │                     │
│  └─────────┘  └─────────┘  └─────────┘                     │
└─────────────────────────────────────────────────────────────┘
         ▲              ▲              ▲
         │              │              │
    ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
    │  Client │    │  Client │    │  Client │
    │  (Agent)│    │  (Agent)│    │  (Agent)│
    └────┬────┘    └────┬────┘    └────┬────┘
         │              │              │
    ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
    │ 服务实例A │    │ 服务实例B │    │ 服务实例C │
    └─────────┘    └─────────┘    └─────────┘
```

#### 完整示例：Go集成Consul

```go
// consul/consul_client.go
package consul

import (
    "context"
    "fmt"
    "strconv"
    "time"

    "github.com/hashicorp/consul/api"
)

// ConsulClient Consul客户端封装
type ConsulClient struct {
    client    *api.Client
    config    *api.Config
    checkPort int
}

// NewConsulClient 创建Consul客户端
func NewConsulClient(addr string) (*ConsulClient, error) {
    config := api.DefaultConfig()
    config.Address = addr

    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create consul client: %w", err)
    }

    return &ConsulClient{
        client: client,
        config: config,
    }, nil
}

// RegisterService 注册服务
func (c *ConsulClient) RegisterService(ctx context.Context, svc *ServiceInstance) error {
    registration := &api.AgentServiceRegistration{
        ID:      svc.ID,
        Name:    svc.Name,
        Tags:    []string{"go", "microservice"},
        Port:    svc.Port,
        Address: svc.Host,
        Meta:    svc.Metadata,
        Weights: &api.AgentWeights{
            Passing: svc.Weight,
            Warning: 1,
        },
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s:%d/health", svc.Host, c.checkPort),
            Interval:                       "10s",
            Timeout:                        "5s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }

    // 支持gRPC健康检查
    if svc.Metadata["protocol"] == "grpc" {
        registration.Check = &api.AgentServiceCheck{
            GRPC:     fmt.Sprintf("%s:%d", svc.Host, svc.Port),
            Interval: "10s",
            Timeout:  "5s",
        }
    }

    return c.client.Agent().ServiceRegister(registration)
}

// DeregisterService 注销服务
func (c *ConsulClient) DeregisterService(ctx context.Context, serviceID string) error {
    return c.client.Agent().ServiceDeregister(serviceID)
}

// DiscoverService 发现服务
func (c *ConsulClient) DiscoverService(ctx context.Context, serviceName string) ([]*ServiceInstance, error) {
    entries, _, err := c.client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }

    var instances []*ServiceInstance
    for _, entry := range entries {
        instance := &ServiceInstance{
            ID:       entry.Service.ID,
            Name:     entry.Service.Service,
            Host:     entry.Service.Address,
            Port:     entry.Service.Port,
            Metadata: entry.Service.Meta,
            Healthy:  entry.Checks.AggregatedStatus() == api.HealthPassing,
        }

        // 解析权重
        if entry.Service.Weights != nil {
            instance.Weight = entry.Service.Weights.Passing
        }

        instances = append(instances, instance)
    }

    return instances, nil
}

// WatchService 监听服务变化
func (c *ConsulClient) WatchService(ctx context.Context, serviceName string) (chan []*ServiceInstance, error) {
    ch := make(chan []*ServiceInstance)

    go func() {
        defer close(ch)

        var lastIndex uint64

        for {
            select {
            case <-ctx.Done():
                return
            default:
            }

            entries, meta, err := c.client.Health().Service(serviceName, "", true, &api.QueryOptions{
                WaitIndex: lastIndex,
                WaitTime:  30 * time.Second,
            })

            if err != nil {
                time.Sleep(time.Second)
                continue
            }

            if meta.LastIndex == lastIndex {
                continue
            }

            lastIndex = meta.LastIndex

            var instances []*ServiceInstance
            for _, entry := range entries {
                instances = append(instances, &ServiceInstance{
                    ID:      entry.Service.ID,
                    Name:    entry.Service.Service,
                    Host:    entry.Service.Address,
                    Port:    entry.Service.Port,
                    Healthy: entry.Checks.AggregatedStatus() == api.HealthPassing,
                })
            }

            select {
            case ch <- instances:
            case <-ctx.Done():
                return
            }
        }
    }()

    return ch, nil
}

// Close 关闭连接
func (c *ConsulClient) Close() error {
    return nil
}

// KV操作封装
func (c *ConsulClient) GetKV(key string) (string, error) {
    kv, _, err := c.client.KV().Get(key, nil)
    if err != nil {
        return "", err
    }
    if kv == nil {
        return "", fmt.Errorf("key not found: %s", key)
    }
    return string(kv.Value), nil
}

func (c *ConsulClient) PutKV(key, value string) error {
    p := &api.KVPair{Key: key, Value: []byte(value)}
    _, err := c.client.KV().Put(p, nil)
    return err
}

// 服务注册助手函数
func RegisterServiceWithConsul(consulAddr string, svc *ServiceInstance) (*ConsulClient, error) {
    client, err := NewConsulClient(consulAddr)
    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.RegisterService(ctx, svc); err != nil {
        return nil, err
    }

    return client, nil
}

// 服务发现助手函数
func DiscoverServiceWithConsul(consulAddr, serviceName string) ([]*ServiceInstance, error) {
    client, err := NewConsulClient(consulAddr)
    if err != nil {
        return nil, err
    }
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return client.DiscoverService(ctx, serviceName)
}
```

#### 分布式锁实现

```go
// consul/distributed_lock.go
package consul

import (
    "context"
    "fmt"
    "time"

    "github.com/hashicorp/consul/api"
)

// DistributedLock 基于Consul的分布式锁
type DistributedLock struct {
    client    *api.Client
    sessionID string
    key       string
    ttl       time.Duration
}

// NewDistributedLock 创建分布式锁
func NewDistributedLock(consulAddr, key string, ttl time.Duration) (*DistributedLock, error) {
    config := api.DefaultConfig()
    config.Address = consulAddr

    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }

    return &DistributedLock{
        client: client,
        key:    key,
        ttl:    ttl,
    }, nil
}

// Lock 获取锁
func (l *DistributedLock) Lock(ctx context.Context) error {
    // 创建session
    session := l.client.Session()
    entry := &api.SessionEntry{
        TTL:      l.ttl.String(),
        Behavior: api.SessionBehaviorDelete,
    }

    sessionID, _, err := session.Create(entry, nil)
    if err != nil {
        return fmt.Errorf("failed to create session: %w", err)
    }

    l.sessionID = sessionID

    // 尝试获取锁
    kv := l.client.KV()
    p := &api.KVPair{
        Key:     l.key,
        Value:   []byte(l.sessionID),
        Session: l.sessionID,
    }

    acquired, _, err := kv.Acquire(p, nil)
    if err != nil {
        session.Destroy(l.sessionID, nil)
        return fmt.Errorf("failed to acquire lock: %w", err)
    }

    if !acquired {
        session.Destroy(l.sessionID, nil)
        return fmt.Errorf("lock already held by another session")
    }

    // 启动续约goroutine
    go l.renewSession()

    return nil
}

// Unlock 释放锁
func (l *DistributedLock) Unlock() error {
    if l.sessionID == "" {
        return nil
    }

    // 释放锁
    kv := l.client.KV()
    p := &api.KVPair{
        Key:     l.key,
        Value:   []byte(l.sessionID),
        Session: l.sessionID,
    }

    _, _, err := kv.Release(p, nil)
    if err != nil {
        return err
    }

    // 销毁session
    session := l.client.Session()
    _, err = session.Destroy(l.sessionID, nil)

    l.sessionID = ""
    return err
}

// renewSession 续约session
func (l *DistributedLock) renewSession() {
    session := l.client.Session()
    ticker := time.NewTicker(l.ttl / 2)
    defer ticker.Stop()

    for range ticker.C {
        if l.sessionID == "" {
            return
        }

        _, _, err := session.Renew(l.sessionID, nil)
        if err != nil {
            return
        }
    }
}

// TryLock 尝试获取锁（非阻塞）
func (l *DistributedLock) TryLock(ctx context.Context) (bool, error) {
    err := l.Lock(ctx)
    if err != nil {
        if err.Error() == "lock already held by another session" {
            return false, nil
        }
        return false, err
    }
    return true, nil
}
```

#### 反例说明

```go
// ❌ 错误：不使用健康检查
registration := &api.AgentServiceRegistration{
    ID:   svc.ID,
    Name: svc.Name,
    Port: svc.Port,
    // 缺少Check配置
}
// 问题：服务故障后Consul仍返回该实例

// ❌ 错误：session不续约
// 获取锁后不启动续约goroutine
// 问题：锁会自动过期，导致并发问题
```

#### 注意事项

1. **集群部署**：生产环境至少3个Server节点保证高可用
2. **健康检查配置**：合理设置检查间隔和超时时间
3. **Session TTL**：分布式锁的TTL应大于业务执行时间
4. **Watch性能**：大量服务监听时考虑使用长轮询优化

---

### 2.3 etcd服务发现

#### 概念定义

etcd是CoreOS开发的分布式键值存储系统，基于Raft协议保证一致性，广泛用于服务发现、配置管理、分布式锁等场景。

#### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        etcd集群                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                     │
│  │ Member  │  │ Member  │  │ Member  │                     │
│  │  :2379  │  │  :2379  │  │  :2379  │   Raft共识          │
│  │ (Leader)│  │(Follower)│  │(Follower)│                    │
│  └─────────┘  └─────────┘  └─────────┘                     │
│       ▲            ▲            ▲                          │
│       │            │            │                          │
│       └────────────┴────────────┘                          │
│              gRPC/HTTP通信                                  │
└─────────────────────────────────────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    ▼                   ▼
              ┌─────────┐         ┌─────────┐
              │ 服务注册  │         │ 服务发现  │
              │  (PUT)  │         │  (GET)  │
              └─────────┘         └─────────┘
```

#### 完整示例：Go集成etcd

```go
// etcd/etcd_client.go
package etcd

import (
    "context"
    "encoding/json"
    "fmt"
    "strings"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdClient etcd客户端封装
type EtcdClient struct {
    client *clientv3.Client
    prefix string
}

// NewEtcdClient 创建etcd客户端
func NewEtcdClient(endpoints []string, prefix string) (*EtcdClient, error) {
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to etcd: %w", err)
    }

    return &EtcdClient{
        client: cli,
        prefix: prefix,
    }, nil
}

// RegisterService 注册服务（带租约自动过期）
func (e *EtcdClient) RegisterService(ctx context.Context, svc *ServiceInstance, ttl int64) error {
    // 创建租约
    lease, err := e.client.Grant(ctx, ttl)
    if err != nil {
        return fmt.Errorf("failed to create lease: %w", err)
    }

    // 序列化服务信息
    value, err := json.Marshal(svc)
    if err != nil {
        return err
    }

    key := fmt.Sprintf("%s/services/%s/%s", e.prefix, svc.Name, svc.ID)

    // 写入etcd
    _, err = e.client.Put(ctx, key, string(value), clientv3.WithLease(lease.ID))
    if err != nil {
        return fmt.Errorf("failed to register service: %w", err)
    }

    // 启动自动续约
    go e.keepAlive(ctx, lease.ID)

    return nil
}

// keepAlive 自动续约
func (e *EtcdClient) keepAlive(ctx context.Context, leaseID clientv3.LeaseID) {
    ch, err := e.client.KeepAlive(ctx, leaseID)
    if err != nil {
        return
    }

    for {
        select {
        case _, ok := <-ch:
            if !ok {
                return
            }
        case <-ctx.Done():
            return
        }
    }
}

// DeregisterService 注销服务
func (e *EtcdClient) DeregisterService(ctx context.Context, serviceName, instanceID string) error {
    key := fmt.Sprintf("%s/services/%s/%s", e.prefix, serviceName, instanceID)
    _, err := e.client.Delete(ctx, key)
    return err
}

// DiscoverService 发现服务
func (e *EtcdClient) DiscoverService(ctx context.Context, serviceName string) ([]*ServiceInstance, error) {
    prefix := fmt.Sprintf("%s/services/%s/", e.prefix, serviceName)

    resp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix())
    if err != nil {
        return nil, err
    }

    var instances []*ServiceInstance
    for _, kv := range resp.Kvs {
        var svc ServiceInstance
        if err := json.Unmarshal(kv.Value, &svc); err != nil {
            continue
        }
        instances = append(instances, &svc)
    }

    return instances, nil
}

// WatchService 监听服务变化
func (e *EtcdClient) WatchService(ctx context.Context, serviceName string) (chan []*ServiceInstance, error) {
    prefix := fmt.Sprintf("%s/services/%s/", e.prefix, serviceName)
    ch := make(chan []*ServiceInstance)

    watchCh := e.client.Watch(ctx, prefix, clientv3.WithPrefix())

    go func() {
        defer close(ch)

        // 先获取当前所有实例
        instances, err := e.DiscoverService(ctx, serviceName)
        if err == nil {
            select {
            case ch <- instances:
            case <-ctx.Done():
                return
            }
        }

        for watchResp := range watchCh {
            if watchResp.Err() != nil {
                continue
            }

            // 重新获取所有实例
            instances, err := e.DiscoverService(ctx, serviceName)
            if err != nil {
                continue
            }

            select {
            case ch <- instances:
            case <-ctx.Done():
                return
            }
        }
    }()

    return ch, nil
}

// Close 关闭连接
func (e *EtcdClient) Close() error {
    return e.client.Close()
}

// 配置管理
func (e *EtcdClient) GetConfig(ctx context.Context, key string) (string, error) {
    fullKey := fmt.Sprintf("%s/config/%s", e.prefix, key)
    resp, err := e.client.Get(ctx, fullKey)
    if err != nil {
        return "", err
    }

    if len(resp.Kvs) == 0 {
        return "", fmt.Errorf("config not found: %s", key)
    }

    return string(resp.Kvs[0].Value), nil
}

func (e *EtcdClient) PutConfig(ctx context.Context, key, value string) error {
    fullKey := fmt.Sprintf("%s/config/%s", e.prefix, key)
    _, err := e.client.Put(ctx, fullKey, value)
    return err
}

func (e *EtcdClient) WatchConfig(ctx context.Context, key string) (chan string, error) {
    fullKey := fmt.Sprintf("%s/config/%s", e.prefix, key)
    ch := make(chan string)

    watchCh := e.client.Watch(ctx, fullKey)

    go func() {
        defer close(ch)
        for watchResp := range watchCh {
            for _, event := range watchResp.Events {
                if event.Type == clientv3.EventTypePut {
                    select {
                    case ch <- string(event.Kv.Value):
                    case <-ctx.Done():
                        return
                    }
                }
            }
        }
    }()

    return ch, nil
}

// 服务选举（基于etcd的分布式选举）
type Election struct {
    client    *clientv3.Client
    session   *concurrency.Session
    election  *concurrency.Election
    keyPrefix string
}

func NewElection(client *clientv3.Client, keyPrefix string) (*Election, error) {
    session, err := concurrency.NewSession(client)
    if err != nil {
        return nil, err
    }

    election := concurrency.NewElection(session, keyPrefix)

    return &Election{
        client:    client,
        session:   session,
        election:  election,
        keyPrefix: keyPrefix,
    }, nil
}

func (e *Election) Campaign(ctx context.Context, val string) error {
    return e.election.Campaign(ctx, val)
}

func (e *Election) Resign(ctx context.Context) error {
    return e.election.Resign(ctx)
}

func (e *Election) IsLeader(ctx context.Context) (bool, error) {
    leader, err := e.election.Leader(ctx)
    if err != nil {
        return false, err
    }

    // 比较当前节点是否是leader
    // 实际实现需要比较leader的值
    _ = leader
    return false, nil
}

func (e *Election) Close() error {
    return e.session.Close()
}
```

#### 分布式锁实现

```go
// etcd/distributed_lock.go
package etcd

import (
    "context"
    "fmt"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
    "go.etcd.io/etcd/client/v3/concurrency"
)

// EtcdLock 基于etcd的分布式锁
type EtcdLock struct {
    client *clientv3.Client
    mutex  *concurrency.Mutex
    session *concurrency.Session
}

// NewEtcdLock 创建分布式锁
func NewEtcdLock(client *clientv3.Client, key string, ttl int) (*EtcdLock, error) {
    session, err := concurrency.NewSession(client, concurrency.WithTTL(ttl))
    if err != nil {
        return nil, err
    }

    mutex := concurrency.NewMutex(session, key)

    return &EtcdLock{
        client:  client,
        mutex:   mutex,
        session: session,
    }, nil
}

// Lock 获取锁（阻塞）
func (l *EtcdLock) Lock(ctx context.Context) error {
    return l.mutex.Lock(ctx)
}

// Unlock 释放锁
func (l *EtcdLock) Unlock(ctx context.Context) error {
    return l.mutex.Unlock(ctx)
}

// TryLock 尝试获取锁（非阻塞）
func (l *EtcdLock) TryLock(ctx context.Context, timeout time.Duration) (bool, error) {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    err := l.mutex.Lock(ctx)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return false, nil
        }
        return false, err
    }
    return true, nil
}

// Close 关闭session
func (l *EtcdLock) Close() error {
    return l.session.Close()
}

// 使用示例
func ExampleEtcdLock() {
    cli, err := clientv3.New(clientv3.Config{
        Endpoints: []string{"localhost:2379"},
    })
    if err != nil {
        panic(err)
    }
    defer cli.Close()

    lock, err := NewEtcdLock(cli, "/locks/my-resource", 10)
    if err != nil {
        panic(err)
    }
    defer lock.Close()

    ctx := context.Background()

    // 获取锁
    if err := lock.Lock(ctx); err != nil {
        panic(err)
    }

    // 执行业务逻辑
    fmt.Println("Lock acquired, doing work...")
    time.Sleep(5 * time.Second)

    // 释放锁
    if err := lock.Unlock(ctx); err != nil {
        panic(err)
    }
    fmt.Println("Lock released")
}
```

#### 反例说明

```go
// ❌ 错误：不使用租约
_, err = client.Put(ctx, key, value)  // 没有使用WithLease
// 问题：服务崩溃后key永久存在，导致僵尸服务

// ❌ 错误：续约不处理channel关闭
for range ch {  // ch可能关闭
// 问题：goroutine泄漏

// ❌ 错误：锁不释放
mutex.Lock(ctx)
// 业务逻辑panic
// mutex.Unlock(ctx)  // 不会执行
// 问题：死锁
```

#### 注意事项

1. **租约TTL设置**：TTL应大于心跳间隔，避免服务正常但租约过期
2. **Watch性能**：大量Watch会增加etcd负载，考虑批量监听
3. **事务使用**：需要原子操作时，使用Txn事务API
4. **集群配置**：生产环境至少3节点，避免脑裂

---

### 2.4 Kubernetes服务发现

#### 概念定义

Kubernetes提供原生的服务发现机制，通过**Service**资源抽象服务访问，通过**DNS**或**环境变量**实现服务发现，支持**Headless Service**直接访问Pod。

#### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                      Kubernetes集群                          │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    Service (ClusterIP)               │   │
│  │  Selector: app=user-service                          │   │
│  │  ClusterIP: 10.96.0.1                                │   │
│  │  DNS: user-service.default.svc.cluster.local         │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                  │
│         ┌─────────────────┼─────────────────┐               │
│         ▼                 ▼                 ▼               │
│    ┌─────────┐      ┌─────────┐      ┌─────────┐           │
│    │  Pod A  │      │  Pod B  │      │  Pod C  │           │
│    │ :8080   │      │ :8080   │      │ :8080   │           │
│    │ Ready   │      │ Ready   │      │ NotReady│           │
│    └─────────┘      └─────────┘      └─────────┘           │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    CoreDNS                          │   │
│  │  user-service.default.svc.cluster.local -> 10.96.0.1│   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

#### 完整示例：Go应用集成K8s服务发现

```go
// k8s/k8s_discovery.go
package k8s

import (
    "context"
    "fmt"
    "os"
    "time"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/cache"
)

// K8sDiscovery Kubernetes服务发现
type K8sDiscovery struct {
    client     *kubernetes.Clientset
    namespace  string
    informer   cache.SharedInformer
    endpoints  map[string][]*Endpoint
    mu         sync.RWMutex
}

// Endpoint 服务端点
type Endpoint struct {
    IP       string
    Port     int32
    Ready    bool
    NodeName string
}

// NewK8sDiscovery 创建K8s服务发现（集群内运行）
func NewK8sDiscovery(namespace string) (*K8sDiscovery, error) {
    // 使用集群内配置
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }

    if namespace == "" {
        namespace = os.Getenv("NAMESPACE")
        if namespace == "" {
            namespace = "default"
        }
    }

    kd := &K8sDiscovery{
        client:    clientset,
        namespace: namespace,
        endpoints: make(map[string][]*Endpoint),
    }

    // 创建informer监听Endpoints变化
    factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
    kd.informer = factory.Core().V1().Endpoints().Informer()

    kd.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc:    kd.onEndpointsAdd,
        UpdateFunc: kd.onEndpointsUpdate,
        DeleteFunc: kd.onEndpointsDelete,
    })

    return kd, nil
}

// Start 启动监听
func (k *K8sDiscovery) Start(ctx context.Context) {
    go k.informer.Run(ctx.Done())

    // 等待同步完成
    if !cache.WaitForCacheSync(ctx.Done(), k.informer.HasSynced) {
        fmt.Println("failed to sync endpoints cache")
        return
    }

    fmt.Println("K8s endpoints cache synced")
}

func (k *K8sDiscovery) onEndpointsAdd(obj interface{}) {
    endpoints := obj.(*corev1.Endpoints)
    k.updateEndpoints(endpoints)
}

func (k *K8sDiscovery) onEndpointsUpdate(old, new interface{}) {
    endpoints := new.(*corev1.Endpoints)
    k.updateEndpoints(endpoints)
}

func (k *K8sDiscovery) onEndpointsDelete(obj interface{}) {
    endpoints := obj.(*corev1.Endpoints)
    k.mu.Lock()
    delete(k.endpoints, endpoints.Name)
    k.mu.Unlock()
}

func (k *K8sDiscovery) updateEndpoints(endpoints *corev1.Endpoints) {
    var eps []*Endpoint

    for _, subset := range endpoints.Subsets {
        for _, addr := range subset.Addresses {
            for _, port := range subset.Ports {
                eps = append(eps, &Endpoint{
                    IP:       addr.IP,
                    Port:     port.Port,
                    Ready:    true,
                    NodeName: *addr.NodeName,
                })
            }
        }

        // 未就绪的地址
        for _, addr := range subset.NotReadyAddresses {
            for _, port := range subset.Ports {
                eps = append(eps, &Endpoint{
                    IP:       addr.IP,
                    Port:     port.Port,
                    Ready:    false,
                    NodeName: *addr.NodeName,
                })
            }
        }
    }

    k.mu.Lock()
    k.endpoints[endpoints.Name] = eps
    k.mu.Unlock()

    fmt.Printf("Updated endpoints for %s: %d addresses\n", endpoints.Name, len(eps))
}

// GetEndpoints 获取服务端点
func (k *K8sDiscovery) GetEndpoints(serviceName string) ([]*Endpoint, error) {
    k.mu.RLock()
    defer k.mu.RUnlock()

    eps, ok := k.endpoints[serviceName]
    if !ok {
        return nil, fmt.Errorf("service not found: %s", serviceName)
    }

    // 只返回就绪的端点
    var ready []*Endpoint
    for _, ep := range eps {
        if ep.Ready {
            ready = append(ready, ep)
        }
    }

    return ready, nil
}

// 使用DNS进行服务发现
func DiscoverServiceByDNS(serviceName, namespace string) (string, error) {
    // K8s DNS格式: <service>.<namespace>.svc.cluster.local
    fqdn := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)

    // 使用标准net包解析
    addrs, err := net.LookupHost(fqdn)
    if err != nil {
        return "", err
    }

    if len(addrs) == 0 {
        return "", fmt.Errorf("no addresses found for %s", fqdn)
    }

    return addrs[0], nil
}

// 使用环境变量进行服务发现
func DiscoverServiceByEnv(serviceName string) (string, error) {
    // K8s会自动注入环境变量
    // <SERVICE_NAME>_SERVICE_HOST
    // <SERVICE_NAME>_SERVICE_PORT

    hostEnv := fmt.Sprintf("%s_SERVICE_HOST", strings.ToUpper(serviceName))
    portEnv := fmt.Sprintf("%s_SERVICE_PORT", strings.ToUpper(serviceName))

    host := os.Getenv(hostEnv)
    port := os.Getenv(portEnv)

    if host == "" || port == "" {
        return "", fmt.Errorf("service env not found: %s", serviceName)
    }

    return fmt.Sprintf("%s:%s", host, port), nil
}

// Headless Service直接发现Pod
func GetPodEndpoints(client *kubernetes.Clientset, namespace, serviceName string) ([]string, error) {
    // 获取Endpoints对象
    endpoints, err := client.CoreV1().Endpoints(namespace).Get(
        context.Background(),
        serviceName,
        metav1.GetOptions{},
    )
    if err != nil {
        return nil, err
    }

    var addresses []string
    for _, subset := range endpoints.Subsets {
        for _, addr := range subset.Addresses {
            for _, port := range subset.Ports {
                addresses = append(addresses, fmt.Sprintf("%s:%d", addr.IP, port.Port))
            }
        }
    }

    return addresses, nil
}
```

#### K8s部署配置

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: user-service:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
  labels:
    app: user-service
spec:
  selector:
    app: user-service
  ports:
  - port: 80
    targetPort: 8080
    name: http
  type: ClusterIP
---
# Headless Service（用于直接访问Pod）
apiVersion: v1
kind: Service
metadata:
  name: user-service-headless
spec:
  selector:
    app: user-service
  ports:
  - port: 8080
  clusterIP: None  # Headless
```

#### 反例说明

```go
// ❌ 错误：硬编码Pod IP
conn, err := grpc.Dial("10.244.1.5:8080")
// 问题：Pod重启后IP变化，连接失败

// ❌ 错误：不使用readiness探针
// 问题：Pod启动过程中被加入Service端点，导致请求失败

// ❌ 错误：忽略namespace
// 跨namespace访问服务时未指定namespace
// 问题：DNS解析失败
```

#### 注意事项

1. **探针配置**：合理配置liveness和readiness探针
2. **优雅关闭**：处理SIGTERM信号，完成当前请求后再退出
3. **资源限制**：设置合理的requests和limits
4. **DNS缓存**：应用内考虑DNS缓存，减少解析开销

---

## 3. 负载均衡

### 3.1 轮询、随机、加权负载均衡

#### 概念定义

负载均衡是将请求分发到多个服务实例的算法。**轮询(Round Robin)**按顺序分配；**随机(Random)**随机选择；**加权(Weighted)**根据实例权重分配。

#### 架构图

```
轮询算法:
请求1 ──> 实例A
请求2 ──> 实例B
请求3 ──> 实例C
请求4 ──> 实例A (循环)

随机算法:
请求1 ──> 实例B (随机)
请求2 ──> 实例A (随机)
请求3 ──> 实例B (随机)

加权轮询算法:
实例A: 权重5
实例B: 权重3
实例C: 权重2
分配序列: A,A,A,A,A,B,B,B,C,C
```

#### 完整示例：Go实现负载均衡算法

```go
// loadbalancer/loadbalancer.go
package loadbalancer

import (
    "errors"
    "math/rand"
    "sync"
    "sync/atomic"
    "time"
)

// Backend 后端服务实例
type Backend struct {
    ID      string
    Address string
    Weight  int
    Healthy bool
}

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
    Select(backends []*Backend) (*Backend, error)
}

// RoundRobinBalancer 轮询负载均衡器
type RoundRobinBalancer struct {
    counter uint64
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
    return &RoundRobinBalancer{}
}

func (r *RoundRobinBalancer) Select(backends []*Backend) (*Backend, error) {
    healthy := filterHealthy(backends)
    if len(healthy) == 0 {
        return nil, errors.New("no healthy backend available")
    }

    // 原子递增计数器
    idx := atomic.AddUint64(&r.counter, 1) % uint64(len(healthy))
    return healthy[idx], nil
}

// RandomBalancer 随机负载均衡器
type RandomBalancer struct {
    rnd *rand.Rand
    mu  sync.Mutex
}

func NewRandomBalancer() *RandomBalancer {
    return &RandomBalancer{
        rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (r *RandomBalancer) Select(backends []*Backend) (*Backend, error) {
    healthy := filterHealthy(backends)
    if len(healthy) == 0 {
        return nil, errors.New("no healthy backend available")
    }

    r.mu.Lock()
    idx := r.rnd.Intn(len(healthy))
    r.mu.Unlock()

    return healthy[idx], nil
}

// WeightedRoundRobinBalancer 加权轮询负载均衡器
type WeightedRoundRobinBalancer struct {
    mu           sync.Mutex
    currentIndex int
    currentWeight int
    maxWeight    int
    gcdWeight    int
    backendCount int
}

func NewWeightedRoundRobinBalancer() *WeightedRoundRobinBalancer {
    return &WeightedRoundRobinBalancer{}
}

func (w *WeightedRoundRobinBalancer) Select(backends []*Backend) (*Backend, error) {
    healthy := filterHealthy(backends)
    if len(healthy) == 0 {
        return nil, errors.New("no healthy backend available")
    }

    w.mu.Lock()
    defer w.mu.Unlock()

    // 计算最大权重和GCD
    w.maxWeight = getMaxWeight(healthy)
    w.gcdWeight = getGCDWeight(healthy)
    w.backendCount = len(healthy)

    for {
        w.currentIndex = (w.currentIndex + 1) % w.backendCount

        if w.currentIndex == 0 {
            w.currentWeight = w.currentWeight - w.gcdWeight
            if w.currentWeight <= 0 {
                w.currentWeight = w.maxWeight
            }
        }

        if healthy[w.currentIndex].Weight >= w.currentWeight {
            return healthy[w.currentIndex], nil
        }
    }
}

// WeightedRandomBalancer 加权随机负载均衡器
type WeightedRandomBalancer struct {
    rnd *rand.Rand
    mu  sync.Mutex
}

func NewWeightedRandomBalancer() *WeightedRandomBalancer {
    return &WeightedRandomBalancer{
        rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (w *WeightedRandomBalancer) Select(backends []*Backend) (*Backend, error) {
    healthy := filterHealthy(backends)
    if len(healthy) == 0 {
        return nil, errors.New("no healthy backend available")
    }

    // 计算总权重
    totalWeight := 0
    for _, b := range healthy {
        totalWeight += b.Weight
    }

    w.mu.Lock()
    randomWeight := w.rnd.Intn(totalWeight)
    w.mu.Unlock()

    // 根据权重选择
    currentWeight := 0
    for _, b := range healthy {
        currentWeight += b.Weight
        if randomWeight < currentWeight {
            return b, nil
        }
    }

    return healthy[len(healthy)-1], nil
}

// 辅助函数
func filterHealthy(backends []*Backend) []*Backend {
    var healthy []*Backend
    for _, b := range backends {
        if b.Healthy {
            healthy = append(healthy, b)
        }
    }
    return healthy
}

func getMaxWeight(backends []*Backend) int {
    max := 0
    for _, b := range backends {
        if b.Weight > max {
            max = b.Weight
        }
    }
    return max
}

func getGCDWeight(backends []*Backend) int {
    if len(backends) == 0 {
        return 0
    }

    gcd := backends[0].Weight
    for _, b := range backends[1:] {
        gcd = calculateGCD(gcd, b.Weight)
    }
    return gcd
}

func calculateGCD(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

// LoadBalancerManager 负载均衡管理器
type LoadBalancerManager struct {
    balancers map[string]LoadBalancer
    backends  map[string][]*Backend
    mu        sync.RWMutex
}

func NewLoadBalancerManager() *LoadBalancerManager {
    return &LoadBalancerManager{
        balancers: make(map[string]LoadBalancer),
        backends:  make(map[string][]*Backend),
    }
}

func (m *LoadBalancerManager) RegisterBalancer(serviceName string, lb LoadBalancer) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.balancers[serviceName] = lb
}

func (m *LoadBalancerManager) UpdateBackends(serviceName string, backends []*Backend) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.backends[serviceName] = backends
}

func (m *LoadBalancerManager) Select(serviceName string) (*Backend, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    lb, ok := m.balancers[serviceName]
    if !ok {
        return nil, errors.New("load balancer not found")
    }

    backends, ok := m.backends[serviceName]
    if !ok || len(backends) == 0 {
        return nil, errors.New("no backends available")
    }

    return lb.Select(backends)
}

// HTTP代理集成负载均衡
type LoadBalancedProxy struct {
    manager *LoadBalancerManager
    client  *http.Client
}

func NewLoadBalancedProxy(manager *LoadBalancerManager) *LoadBalancedProxy {
    return &LoadBalancedProxy{
        manager: manager,
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (p *LoadBalancedProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 从URL路径解析服务名
    serviceName := extractServiceName(r.URL.Path)

    backend, err := p.manager.Select(serviceName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }

    // 构建目标URL
    targetURL := &url.URL{
        Scheme: "http",
        Host:   backend.Address,
        Path:   r.URL.Path,
    }

    // 创建代理请求
    proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 复制Header
    proxyReq.Header = r.Header

    // 发送请求
    resp, err := p.client.Do(proxyReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // 复制响应
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

func extractServiceName(path string) string {
    parts := strings.Split(strings.Trim(path, "/"), "/")
    if len(parts) > 0 {
        return parts[0]
    }
    return ""
}
```

#### 反例说明

```go
// ❌ 错误：不考虑实例健康状态
func (r *RoundRobinBalancer) Select(backends []*Backend) (*Backend, error) {
    idx := atomic.AddUint64(&r.counter, 1) % uint64(len(backends))
    return backends[idx], nil  // 可能返回不健康实例
}

// ❌ 错误：加权算法权重为0时panic
if backend.Weight >= w.currentWeight {
    // 如果所有权重都是0，会无限循环
}

// ❌ 错误：随机算法使用全局rand
idx := rand.Intn(len(backends))  // 不是并发安全的
```

#### 注意事项

1. **健康检查过滤**：选择前必须过滤掉不健康实例
2. **权重配置**：权重为0的实例不应参与负载均衡
3. **并发安全**：计数器和随机数生成器需要并发安全
4. **空实例处理**：无可用实例时返回明确错误

---

### 3.2 一致性哈希

#### 概念定义

一致性哈希(Consistent Hashing)是一种特殊的哈希算法，在添加或删除节点时只需重新定位少量数据，广泛用于分布式缓存、负载均衡等场景。

#### 架构图

```
一致性哈希环:

         0°
         │
    315° ┼──────┐ 45°
         │      │
  270° ──┤  环  ├── 90°
         │      │
   225° └──────┘ 135°
         │
        180°

节点分布:
Node A: hash=30°
Node B: hash=120°
Node C: hash=240°

数据定位:
Key1(hash=50°) -> Node B (顺时针第一个)
Key2(hash=200°) -> Node C
Key3(hash=300°) -> Node A

虚拟节点:
Node A-1: hash=25°
Node A-2: hash=35°
Node A-3: hash=45°
(增加虚拟节点使分布更均匀)
```

#### 完整示例：Go实现一致性哈希

```go
// consistenthash/consistent_hash.go
package consistenthash

import (
    "hash/crc32"
    "sort"
    "strconv"
    "sync"
)

// Hash 哈希函数类型
type Hash func(data []byte) uint32

// ConsistentHash 一致性哈希
type ConsistentHash struct {
    hash     Hash           // 哈希函数
    replicas int            // 每个真实节点的虚拟节点数
    keys     []int          // 排序后的哈希环
    hashMap  map[int]string // 哈希值到真实节点的映射
    mu       sync.RWMutex
}

// New 创建一致性哈希
func New(replicas int, fn Hash) *ConsistentHash {
    if replicas <= 0 {
        replicas = 150  // 默认虚拟节点数
    }

    if fn == nil {
        fn = crc32.ChecksumIEEE  // 默认使用CRC32
    }

    return &ConsistentHash{
        replicas: replicas,
        hash:     fn,
        hashMap:  make(map[int]string),
    }
}

// Add 添加节点
func (c *ConsistentHash) Add(nodes ...string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    for _, node := range nodes {
        // 为每个真实节点创建虚拟节点
        for i := 0; i < c.replicas; i++ {
            // 虚拟节点key格式: node#i
            virtualKey := strconv.Itoa(i) + node
            hash := int(c.hash([]byte(virtualKey)))

            c.keys = append(c.keys, hash)
            c.hashMap[hash] = node
        }
    }

    // 排序哈希环
    sort.Ints(c.keys)
}

// Remove 移除节点
func (c *ConsistentHash) Remove(node string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    for i := 0; i < c.replicas; i++ {
        virtualKey := strconv.Itoa(i) + node
        hash := int(c.hash([]byte(virtualKey)))

        // 从keys中删除
        idx := sort.SearchInts(c.keys, hash)
        if idx < len(c.keys) && c.keys[idx] == hash {
            c.keys = append(c.keys[:idx], c.keys[idx+1:]...)
        }

        // 从hashMap中删除
        delete(c.hashMap, hash)
    }
}

// Get 获取key对应的节点
func (c *ConsistentHash) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if len(c.keys) == 0 {
        return "", false
    }

    hash := int(c.hash([]byte(key)))

    // 二分查找第一个大于等于hash的位置
    idx := sort.Search(len(c.keys), func(i int) bool {
        return c.keys[i] >= hash
    })

    // 如果超出范围，回到开头（环形）
    if idx == len(c.keys) {
        idx = 0
    }

    return c.hashMap[c.keys[idx]], true
}

// GetN 获取key对应的N个不同节点（用于副本）
func (c *ConsistentHash) GetN(key string, n int) []string {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if len(c.keys) == 0 || n <= 0 {
        return nil
    }

    if n > len(c.keys) {
        n = len(c.keys)
    }

    hash := int(c.hash([]byte(key)))
    idx := sort.Search(len(c.keys), func(i int) bool {
        return c.keys[i] >= hash
    })

    nodes := make([]string, 0, n)
    seen := make(map[string]bool)

    for len(nodes) < n {
        if idx >= len(c.keys) {
            idx = 0
        }

        node := c.hashMap[c.keys[idx]]
        if !seen[node] {
            nodes = append(nodes, node)
            seen[node] = true
        }

        idx++
    }

    return nodes
}

// GetAllNodes 获取所有真实节点
func (c *ConsistentHash) GetAllNodes() []string {
    c.mu.RLock()
    defer c.mu.RUnlock()

    seen := make(map[string]bool)
    var nodes []string

    for _, node := range c.hashMap {
        if !seen[node] {
            nodes = append(nodes, node)
            seen[node] = true
        }
    }

    return nodes
}

// 用于负载均衡的一致性哈希选择器
type ConsistentHashSelector struct {
    ch *ConsistentHash
}

func NewConsistentHashSelector(replicas int) *ConsistentHashSelector {
    return &ConsistentHashSelector{
        ch: New(replicas, nil),
    }
}

func (s *ConsistentHashSelector) AddBackends(backends []*Backend) {
    for _, b := range backends {
        if b.Healthy {
            s.ch.Add(b.Address)
        }
    }
}

func (s *ConsistentHashSelector) Select(key string) (string, bool) {
    return s.ch.Get(key)
}

// 使用示例：分布式缓存

type DistributedCache struct {
    ch      *ConsistentHash
    clients map[string]*CacheClient
    mu      sync.RWMutex
}

func NewDistributedCache(nodes []string) *DistributedCache {
    dc := &DistributedCache{
        ch:      New(150, nil),
        clients: make(map[string]*CacheClient),
    }

    for _, node := range nodes {
        dc.ch.Add(node)
        dc.clients[node] = NewCacheClient(node)
    }

    return dc
}

func (dc *DistributedCache) Get(ctx context.Context, key string) (string, error) {
    node, ok := dc.ch.Get(key)
    if !ok {
        return "", errors.New("no cache node available")
    }

    dc.mu.RLock()
    client := dc.clients[node]
    dc.mu.RUnlock()

    return client.Get(ctx, key)
}

func (dc *DistributedCache) Set(ctx context.Context, key, value string) error {
    node, ok := dc.ch.Get(key)
    if !ok {
        return errors.New("no cache node available")
    }

    dc.mu.RLock()
    client := dc.clients[node]
    dc.mu.RUnlock()

    return client.Set(ctx, key, value)
}

// 添加节点（仅影响部分key）
func (dc *DistributedCache) AddNode(node string) {
    dc.mu.Lock()
    defer dc.mu.Unlock()

    dc.ch.Add(node)
    dc.clients[node] = NewCacheClient(node)
}

// 移除节点（需要迁移数据）
func (dc *DistributedCache) RemoveNode(node string) {
    dc.mu.Lock()
    defer dc.mu.Unlock()

    // 获取该节点上的所有key
    keys := dc.clients[node].GetAllKeys()

    // 移除节点
    dc.ch.Remove(node)

    // 重新分布key
    for _, key := range keys {
        value, _ := dc.clients[node].Get(context.Background(), key)
        newNode, ok := dc.ch.Get(key)
        if ok && newNode != node {
            dc.clients[newNode].Set(context.Background(), key, value)
        }
    }

    delete(dc.clients, node)
}

// 统计分布均匀性
func (c *ConsistentHash) DistributionStats() map[string]int {
    c.mu.RLock()
    defer c.mu.RUnlock()

    stats := make(map[string]int)
    for _, node := range c.hashMap {
        stats[node]++
    }
    return stats
}
```

#### 反例说明

```go
// ❌ 错误：虚拟节点数太少
ch := New(3, nil)  // 只有3个虚拟节点
// 问题：分布不均匀，某些节点负载过高

// ❌ 错误：不使用排序
// 每次查找都遍历所有key
// 问题：时间复杂度O(n)，应为O(log n)

// ❌ 错误：并发不安全
func (c *ConsistentHash) Add(node string) {
    hash := int(c.hash([]byte(node)))
    c.keys = append(c.keys, hash)  // 并发写会panic
    c.hashMap[hash] = node
}
```

#### 注意事项

1. **虚拟节点数**：通常150-200个，平衡均匀性和内存占用
2. **哈希函数选择**：CRC32快速但分布一般，MurmurHash3更好
3. **节点变化处理**：添加/删除节点时需要考虑数据迁移
4. **并发安全**：读写操作需要加锁保护

---

### 3.3 健康检查

#### 概念定义

健康检查是检测服务实例可用性的机制，包括**主动检查**（定期探测）和**被动检查**（根据请求结果判断），用于及时剔除故障实例。

#### 架构图

```
主动健康检查:
┌─────────┐    定期HTTP/TCP检查    ┌─────────┐
│  检查器  │ ────────────────────> │ 服务实例 │
│         │ <──────────────────── │         │
└─────────┘   返回200 OK/超时     └─────────┘
         │
         │ 检查失败超过阈值
         ▼
    ┌─────────┐
    │ 标记不健康 │
    │ 停止转发  │
    └─────────┘

被动健康检查:
┌─────────┐     转发请求      ┌─────────┐
│ 负载均衡器│ ──────────────> │ 服务实例A │
│         │ <────────────── │  超时    │
└─────────┘   连续失败3次    └─────────┘
         │
         ▼
    标记A为不健康
    后续请求不再转发到A
```

#### 完整示例：Go实现健康检查

```go
// healthcheck/health_check.go
package healthcheck

import (
    "context"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// HealthStatus 健康状态
type HealthStatus int

const (
    StatusUnknown HealthStatus = iota
    StatusHealthy
    StatusUnhealthy
    StatusDegraded
)

func (s HealthStatus) String() string {
    switch s {
    case StatusHealthy:
        return "healthy"
    case StatusUnhealthy:
        return "unhealthy"
    case StatusDegraded:
        return "degraded"
    default:
        return "unknown"
    }
}

// HealthChecker 健康检查器接口
type HealthChecker interface {
    Check(ctx context.Context) HealthStatus
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
    Interval         time.Duration // 检查间隔
    Timeout          time.Duration // 超时时间
    FailureThreshold int           // 失败阈值
    SuccessThreshold int           // 成功阈值
}

// DefaultHealthCheckConfig 默认配置
func DefaultHealthCheckConfig() *HealthCheckConfig {
    return &HealthCheckConfig{
        Interval:         10 * time.Second,
        Timeout:          5 * time.Second,
        FailureThreshold: 3,
        SuccessThreshold: 2,
    }
}

// HTTPHealthChecker HTTP健康检查器
type HTTPHealthChecker struct {
    client  *http.Client
    url     string
    config  *HealthCheckConfig
}

func NewHTTPHealthChecker(url string, config *HealthCheckConfig) *HTTPHealthChecker {
    if config == nil {
        config = DefaultHealthCheckConfig()
    }

    return &HTTPHealthChecker{
        client: &http.Client{
            Timeout: config.Timeout,
        },
        url:    url,
        config: config,
    }
}

func (h *HTTPHealthChecker) Check(ctx context.Context) HealthStatus {
    req, err := http.NewRequestWithContext(ctx, "GET", h.url, nil)
    if err != nil {
        return StatusUnhealthy
    }

    resp, err := h.client.Do(req)
    if err != nil {
        return StatusUnhealthy
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        return StatusHealthy
    }

    if resp.StatusCode >= 500 {
        return StatusUnhealthy
    }

    return StatusDegraded
}

// TCPHealthChecker TCP健康检查器
type TCPHealthChecker struct {
    address string
    config  *HealthCheckConfig
}

func NewTCPHealthChecker(address string, config *HealthCheckConfig) *TCPHealthChecker {
    if config == nil {
        config = DefaultHealthCheckConfig()
    }

    return &TCPHealthChecker{
        address: address,
        config:  config,
    }
}

func (t *TCPHealthChecker) Check(ctx context.Context) HealthStatus {
    conn, err := net.DialTimeout("tcp", t.address, t.config.Timeout)
    if err != nil {
        return StatusUnhealthy
    }
    defer conn.Close()

    return StatusHealthy
}

// GRPCHealthChecker gRPC健康检查器
type GRPCHealthChecker struct {
    address string
    config  *HealthCheckConfig
}

func (g *GRPCHealthChecker) Check(ctx context.Context) HealthStatus {
    // 使用gRPC健康检查协议
    conn, err := grpc.Dial(g.address, grpc.WithInsecure(), grpc.WithBlock(),
        grpc.WithTimeout(g.config.Timeout))
    if err != nil {
        return StatusUnhealthy
    }
    defer conn.Close()

    healthClient := grpc_health_v1.NewHealthClient(conn)
    resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
    if err != nil {
        return StatusUnhealthy
    }

    switch resp.Status {
    case grpc_health_v1.HealthCheckResponse_SERVING:
        return StatusHealthy
    case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
        return StatusUnhealthy
    default:
        return StatusUnknown
    }
}

// HealthCheckManager 健康检查管理器
type HealthCheckManager struct {
    checkers   map[string]HealthChecker
    statuses   map[string]*HealthStatusRecord
    config     *HealthCheckConfig
    mu         sync.RWMutex
    stopChan   chan struct{}
}

type HealthStatusRecord struct {
    Status         HealthStatus
    LastCheck      time.Time
    FailureCount   int
    SuccessCount   int
    LastFailReason string
}

func NewHealthCheckManager(config *HealthCheckConfig) *HealthCheckManager {
    if config == nil {
        config = DefaultHealthCheckConfig()
    }

    return &HealthCheckManager{
        checkers: make(map[string]HealthChecker),
        statuses: make(map[string]*HealthStatusRecord),
        config:   config,
        stopChan: make(chan struct{}),
    }
}

func (m *HealthCheckManager) RegisterChecker(id string, checker HealthChecker) {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.checkers[id] = checker
    m.statuses[id] = &HealthStatusRecord{
        Status:    StatusUnknown,
        LastCheck: time.Time{},
    }
}

func (m *HealthCheckManager) Start() {
    ticker := time.NewTicker(m.config.Interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            m.runChecks()
        case <-m.stopChan:
            return
        }
    }
}

func (m *HealthCheckManager) Stop() {
    close(m.stopChan)
}

func (m *HealthCheckManager) runChecks() {
    m.mu.RLock()
    checkers := make(map[string]HealthChecker)
    for id, checker := range m.checkers {
        checkers[id] = checker
    }
    m.mu.RUnlock()

    for id, checker := range checkers {
        ctx, cancel := context.WithTimeout(context.Background(), m.config.Timeout)
        status := checker.Check(ctx)
        cancel()

        m.updateStatus(id, status)
    }
}

func (m *HealthCheckManager) updateStatus(id string, status HealthStatus) {
    m.mu.Lock()
    defer m.mu.Unlock()

    record, ok := m.statuses[id]
    if !ok {
        return
    }

    record.LastCheck = time.Now()

    if status == StatusHealthy {
        record.SuccessCount++
        record.FailureCount = 0

        if record.SuccessCount >= m.config.SuccessThreshold {
            record.Status = StatusHealthy
        }
    } else {
        record.FailureCount++
        record.SuccessCount = 0

        if record.FailureCount >= m.config.FailureThreshold {
            record.Status = StatusUnhealthy
            record.LastFailReason = fmt.Sprintf("Failed %d consecutive checks", record.FailureCount)
        }
    }
}

func (m *HealthCheckManager) GetStatus(id string) (HealthStatus, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    record, ok := m.statuses[id]
    if !ok {
        return StatusUnknown, false
    }

    return record.Status, true
}

func (m *HealthCheckManager) GetAllStatuses() map[string]HealthStatus {
    m.mu.RLock()
    defer m.mu.RUnlock()

    statuses := make(map[string]HealthStatus)
    for id, record := range m.statuses {
        statuses[id] = record.Status
    }

    return statuses
}

// 被动健康检查（基于请求结果）
type PassiveHealthChecker struct {
    failureThreshold int
    failureWindow    time.Duration
    records          map[string]*PassiveRecord
    mu               sync.RWMutex
}

type PassiveRecord struct {
    Failures   []time.Time
    LastAccess time.Time
}

func NewPassiveHealthChecker(failureThreshold int, failureWindow time.Duration) *PassiveHealthChecker {
    return &PassiveHealthChecker{
        failureThreshold: failureThreshold,
        failureWindow:    failureWindow,
        records:          make(map[string]*PassiveRecord),
    }
}

func (p *PassiveHealthChecker) RecordSuccess(backendID string) {
    p.mu.Lock()
    defer p.mu.Unlock()

    record, ok := p.records[backendID]
    if !ok {
        record = &PassiveRecord{}
        p.records[backendID] = record
    }

    record.LastAccess = time.Now()
    record.Failures = nil // 清空失败记录
}

func (p *PassiveHealthChecker) RecordFailure(backendID string) {
    p.mu.Lock()
    defer p.mu.Unlock()

    record, ok := p.records[backendID]
    if !ok {
        record = &PassiveRecord{}
        p.records[backendID] = record
    }

    record.LastAccess = time.Now()
    record.Failures = append(record.Failures, time.Now())

    // 清理过期失败记录
    cutoff := time.Now().Add(-p.failureWindow)
    var validFailures []time.Time
    for _, t := range record.Failures {
        if t.After(cutoff) {
            validFailures = append(validFailures, t)
        }
    }
    record.Failures = validFailures
}

func (p *PassiveHealthChecker) IsHealthy(backendID string) bool {
    p.mu.RLock()
    defer p.mu.RUnlock()

    record, ok := p.records[backendID]
    if !ok {
        return true // 无记录默认为健康
    }

    return len(record.Failures) < p.failureThreshold
}
```

#### 反例说明

```go
// ❌ 错误：检查间隔太短
Interval: 100 * time.Millisecond
// 问题：增加服务负载，可能被视为攻击

// ❌ 错误：不清理过期记录
// 被动检查失败记录无限累积
// 问题：内存泄漏

// ❌ 错误：单点检查器
// 只有一个检查器检查所有服务
// 问题：单点故障，检查器故障时无法检测服务状态
```

#### 注意事项

1. **检查间隔平衡**：太频繁增加负载，太少延迟发现故障
2. **阈值设置**：根据服务特性和SLA要求设置合理阈值
3. **多种检查方式**：结合主动和被动检查提高准确性
4. **优雅降级**：检查器本身故障时，默认认为服务健康

---

### 3.4 gRPC负载均衡

#### 概念定义

gRPC负载均衡支持**客户端负载均衡**（通过resolver和picker）和**服务端负载均衡**（通过L4/L7代理），支持多种负载均衡策略。

#### 架构图

```
gRPC客户端负载均衡:
┌─────────────────────────────────────────────────────────────┐
│                     gRPC Client                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  Resolver   │  │   Picker    │  │  Balancer   │         │
│  │  (服务发现)  │  │  (选择策略)  │  │  (负载均衡)  │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
         │                │                │
         ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────┐
│                     SubConn Pool                             │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                     │
│  │ SubConn │  │ SubConn │  │ SubConn │                     │
│  │ 实例A   │  │ 实例B   │  │ 实例C   │                     │
│  └────┬────┘  └────┬────┘  └────┬────┘                     │
└───────┼────────────┼────────────┼───────────────────────────┘
        │            │            │
        ▼            ▼            ▼
   ┌─────────┐  ┌─────────┐  ┌─────────┐
   │ 服务实例A │  │ 服务实例B │  │ 服务实例C │
   └─────────┘  └─────────┘  └─────────┘
```

#### 完整示例：Go实现gRPC负载均衡

```go
// grpc/lb_resolver.go
package grpc

import (
    "google.golang.org/grpc/resolver"
)

// ConsulResolverBuilder Consul解析器构建器
type ConsulResolverBuilder struct {
    consulAddr string
}

func NewConsulResolverBuilder(consulAddr string) *ConsulResolverBuilder {
    return &ConsulResolverBuilder{consulAddr: consulAddr}
}

func (b *ConsulResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
    r := &ConsulResolver{
        target: target,
        cc:     cc,
        consul: b.consulAddr,
        stopCh: make(chan struct{}),
    }

    go r.start()
    return r, nil
}

func (b *ConsulResolverBuilder) Scheme() string {
    return "consul"
}

func init() {
    resolver.Register(NewConsulResolverBuilder("localhost:8500"))
}

// ConsulResolver Consul解析器
type ConsulResolver struct {
    target resolver.Target
    cc     resolver.ClientConn
    consul string
    stopCh chan struct{}
}

func (r *ConsulResolver) start() {
    // 从Consul获取服务地址
    // 监听服务变化
    // 调用r.cc.UpdateState更新地址
}

func (r *ConsulResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (r *ConsulResolver) Close() {
    close(r.stopCh)
}

// 自定义Picker实现加权轮询

import (
    "google.golang.org/grpc/balancer"
    "google.golang.org/grpc/balancer/base"
)

const WeightedRoundRobin = "weighted_round_robin"

func init() {
    balancer.Register(base.NewBalancerBuilder(
        WeightedRoundRobin,
        &WeightedRRPickerBuilder{},
        base.Config{HealthCheck: true},
    ))
}

type WeightedRRPickerBuilder struct{}

func (b *WeightedRRPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
    if len(info.ReadySCs) == 0 {
        return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
    }

    var scs []balancer.SubConn
    var weights []int

    for sc, scInfo := range info.ReadySCs {
        scs = append(scs, sc)

        // 从地址属性获取权重
        weight := 1
        if scInfo.Address.Attributes != nil {
            if w, ok := scInfo.Address.Attributes.Value("weight").(int); ok {
                weight = w
            }
        }
        weights = append(weights, weight)
    }

    return &WeightedRRPicker{
        subConns: scs,
        weights:  weights,
    }
}

type WeightedRRPicker struct {
    subConns []balancer.SubConn
    weights  []int
    current  uint64
}

func (p *WeightedRRPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
    if len(p.subConns) == 0 {
        return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
    }

    // 加权轮询选择
    idx := p.weightedSelect()

    return balancer.PickResult{
        SubConn: p.subConns[idx],
        Done:    func(info balancer.DoneInfo) {},
    }, nil
}

func (p *WeightedRRPicker) weightedSelect() int {
    // 实现加权轮询算法
    totalWeight := 0
    for _, w := range p.weights {
        totalWeight += w
    }

    current := atomic.AddUint64(&p.current, 1) % uint64(totalWeight)

    cumulativeWeight := 0
    for i, w := range p.weights {
        cumulativeWeight += w
        if uint64(cumulativeWeight) > current {
            return i
        }
    }

    return len(p.subConns) - 1
}

// 使用示例
func ExampleGRPCClient() {
    // 使用Consul解析器和加权轮询负载均衡
    conn, err := grpc.Dial(
        "consul:///user-service",
        grpc.WithDefaultServiceConfig(`{
            "loadBalancingConfig": [{"weighted_round_robin": {}}]
        }`),
        grpc.WithInsecure(),
    )
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    client := pb.NewUserServiceClient(conn)
    // 使用client调用服务...
}

// 服务端实现健康检查
import (
    "google.golang.org/grpc/health"
    healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func startGRPCServer() {
    server := grpc.NewServer()

    // 注册健康检查服务
    healthServer := health.NewServer()
    healthpb.RegisterHealthServer(server, healthServer)

    // 设置服务状态
    healthServer.SetServingStatus("user.UserService", healthpb.HealthCheckResponse_SERVING)

    // 注册业务服务
    pb.RegisterUserServiceServer(server, &UserServiceImpl{})

    lis, _ := net.Listen("tcp", ":50051")
    server.Serve(lis)
}
```

#### 反例说明

```go
// ❌ 错误：不使用连接池
// 每次请求新建连接
// 问题：性能差，无法复用HTTP/2多路复用

// ❌ 错误：忽略gRPC状态码
// 不根据状态码调整负载均衡策略
// 问题：持续向故障节点发送请求

// ❌ 错误：不实现健康检查
// 服务端不实现Health服务
// 问题：客户端无法感知服务健康状态
```

#### 注意事项

1. **连接复用**：gRPC基于HTTP/2，应复用连接实现多路复用
2. **健康检查**：服务端必须实现Health服务
3. **优雅关闭**：服务端关闭前发送GOAWAY帧
4. **重试策略**：配置合理的重试和退避策略

---

## 4. 熔断与降级

### 4.1 熔断器模式（Circuit Breaker）

#### 概念定义

熔断器模式是一种容错设计，当服务故障率达到阈值时，**自动断开**对故障服务的调用，防止故障扩散。包含三种状态：**Closed（关闭）**、**Open（打开）**、**Half-Open（半开）**。

#### 架构图

```
熔断器状态转换:

                    失败率达到阈值
    ┌─────────────>─────────────┐
    │                           │
    ▼                           │
┌─────────┐              ┌─────────┐
│  CLOSED │              │  OPEN   │
│  (正常)  │              │ (熔断)  │
│ 允许请求 │              │拒绝请求 │
└────┬────┘              └────┬────┘
     │                        │
     │ 请求成功                │ 超时后
     │                         ▼
     │                   ┌─────────┐
     └───────────────────┤HALF-OPEN│
         测试请求成功     │(半开)   │
                         │允许有限 │
                         │请求测试 │
                         └─────────┘

工作流程:
1. CLOSED: 正常处理请求，统计失败率
2. OPEN: 失败率超过阈值，直接返回错误
3. HALF-OPEN: 超时后允许有限请求测试服务恢复
```

#### 完整示例：Go实现熔断器

```go
// circuitbreaker/circuit_breaker.go
package circuitbreaker

import (
    "context"
    "errors"
    "sync"
    "sync/atomic"
    "time"
)

// State 熔断器状态
type State int32

const (
    StateClosed State = iota    // 关闭状态，正常请求
    StateOpen                   // 打开状态，拒绝请求
    StateHalfOpen               // 半开状态，测试请求
)

func (s State) String() string {
    switch s {
    case StateClosed:
        return "closed"
    case StateOpen:
        return "open"
    case StateHalfOpen:
        return "half-open"
    default:
        return "unknown"
    }
}

// Config 熔断器配置
type Config struct {
    FailureThreshold    uint32        // 失败阈值（次数或百分比）
    SuccessThreshold    uint32        // 成功阈值（半开状态下）
    Timeout             time.Duration // 熔断持续时间
    HalfOpenMaxRequests uint32        // 半开状态最大请求数
    Interval            time.Duration // 统计窗口
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
    return &Config{
        FailureThreshold:    5,
        SuccessThreshold:    3,
        Timeout:             30 * time.Second,
        HalfOpenMaxRequests: 3,
        Interval:            10 * time.Second,
    }
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
    name           string
    config         *Config
    state          int32 // 原子操作
    failureCount   uint32
    successCount   uint32
    lastFailureTime int64 // UnixNano
    halfOpenCount  uint32
    mu             sync.Mutex
}

// New 创建熔断器
func New(name string, config *Config) *CircuitBreaker {
    if config == nil {
        config = DefaultConfig()
    }

    return &CircuitBreaker{
        name:   name,
        config: config,
        state:  int32(StateClosed),
    }
}

// State 获取当前状态
func (cb *CircuitBreaker) State() State {
    return State(atomic.LoadInt32(&cb.state))
}

// Execute 执行受保护的函数
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
    state := cb.State()

    switch state {
    case StateOpen:
        // 检查是否可以转换为半开状态
        if cb.canTransitionToHalfOpen() {
            cb.transitionTo(StateHalfOpen)
        } else {
            return ErrCircuitOpen
        }

    case StateHalfOpen:
        // 限制半开状态的并发请求数
        if !cb.allowHalfOpenRequest() {
            return ErrTooManyRequests
        }
    }

    // 执行函数
    err := fn()

    // 记录结果
    cb.recordResult(err)

    return err
}

// canTransitionToHalfOpen 检查是否可以转换到半开状态
func (cb *CircuitBreaker) canTransitionToHalfOpen() bool {
    lastFailure := atomic.LoadInt64(&cb.lastFailureTime)
    if lastFailure == 0 {
        return false
    }

    elapsed := time.Since(time.Unix(0, lastFailure))
    return elapsed >= cb.config.Timeout
}

// allowHalfOpenRequest 检查是否允许半开状态请求
func (cb *CircuitBreaker) allowHalfOpenRequest() bool {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.halfOpenCount < cb.config.HalfOpenMaxRequests {
        cb.halfOpenCount++
        return true
    }
    return false
}

// recordResult 记录执行结果
func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    state := cb.State()

    if err != nil {
        // 记录失败
        cb.failureCount++
        atomic.StoreInt64(&cb.lastFailureTime, time.Now().UnixNano())

        switch state {
        case StateClosed:
            if cb.failureCount >= cb.config.FailureThreshold {
                cb.transitionTo(StateOpen)
            }
        case StateHalfOpen:
            cb.transitionTo(StateOpen)
        }
    } else {
        // 记录成功
        cb.successCount++

        switch state {
        case StateHalfOpen:
            if cb.successCount >= cb.config.SuccessThreshold {
                cb.transitionTo(StateClosed)
            }
        case StateClosed:
            // 重置失败计数
            if cb.successCount >= cb.config.SuccessThreshold {
                cb.failureCount = 0
            }
        }
    }
}

// transitionTo 状态转换
func (cb *CircuitBreaker) transitionTo(state State) {
    oldState := cb.State()
    atomic.StoreInt32(&cb.state, int32(state))

    // 重置计数器
    cb.failureCount = 0
    cb.successCount = 0
    cb.halfOpenCount = 0

    // 可以在这里触发状态变更事件
    cb.onStateChange(oldState, state)
}

func (cb *CircuitBreaker) onStateChange(from, to State) {
    // 记录日志或发送指标
    // log.Printf("CircuitBreaker %s: %s -> %s", cb.name, from, to)
}

// 错误定义
var (
    ErrCircuitOpen   = errors.New("circuit breaker is open")
    ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// CircuitBreakerManager 熔断器管理器
type CircuitBreakerManager struct {
    breakers map[string]*CircuitBreaker
    config   *Config
    mu       sync.RWMutex
}

func NewCircuitBreakerManager(config *Config) *CircuitBreakerManager {
    return &CircuitBreakerManager{
        breakers: make(map[string]*CircuitBreaker),
        config:   config,
    }
}

func (m *CircuitBreakerManager) Get(name string) *CircuitBreaker {
    m.mu.RLock()
    cb, ok := m.breakers[name]
    m.mu.RUnlock()

    if ok {
        return cb
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    // 双重检查
    if cb, ok := m.breakers[name]; ok {
        return cb
    }

    cb = New(name, m.config)
    m.breakers[name] = cb
    return cb
}

// HTTP客户端集成熔断器
type CircuitBreakerClient struct {
    client *http.Client
    cb     *CircuitBreaker
}

func NewCircuitBreakerClient(cb *CircuitBreaker) *CircuitBreakerClient {
    return &CircuitBreakerClient{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        cb: cb,
    }
}

func (c *CircuitBreakerClient) Do(req *http.Request) (*http.Response, error) {
    var resp *http.Response
    var err error

    cbErr := c.cb.Execute(req.Context(), func() error {
        resp, err = c.client.Do(req)
        if err != nil {
            return err
        }

        // 5xx错误视为服务故障
        if resp.StatusCode >= 500 {
            return errors.New("server error")
        }

        return nil
    })

    if cbErr != nil {
        return nil, cbErr
    }

    return resp, err
}

// gRPC客户端集成熔断器
type CircuitBreakerGRPCClient struct {
    conn *grpc.ClientConn
    cb   *CircuitBreaker
}

func (c *CircuitBreakerGRPCClient) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
    return c.cb.Execute(ctx, func() error {
        return c.conn.Invoke(ctx, method, args, reply, opts...)
    })
}

// 统计信息
type CircuitBreakerStats struct {
    Name         string
    State        State
    FailureCount uint32
    SuccessCount uint32
}

func (cb *CircuitBreaker) Stats() *CircuitBreakerStats {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    return &CircuitBreakerStats{
        Name:         cb.name,
        State:        cb.State(),
        FailureCount: cb.failureCount,
        SuccessCount: cb.successCount,
    }
}
```

#### 反例说明

```go
// ❌ 错误：熔断阈值设置不合理
FailureThreshold: 1000,  // 太高，故障发现延迟
Timeout: 5 * time.Second, // 太短，服务未恢复就测试

// ❌ 错误：不考虑错误类型
// 将所有错误都视为服务故障
// 问题：客户端错误(4xx)不应触发熔断

// ❌ 错误：半开状态无限制
HalfOpenMaxRequests: 0  // 无限制
// 问题：服务刚恢复时可能再次被压垮
```

#### 注意事项

1. **阈值调优**：根据服务特性和SLA要求调整阈值
2. **错误分类**：区分服务端错误(5xx)和客户端错误(4xx)
3. **监控告警**：熔断状态变更应触发告警
4. **渐进恢复**：半开状态逐步增加流量

---

### 4.2 降级策略

#### 概念定义

降级是在系统资源不足或服务故障时，**主动牺牲非核心功能**，保证核心功能可用的策略。包括**功能降级**、**数据降级**、**页面降级**等。

#### 架构图

```
正常模式:
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ 核心功能 │    │推荐功能 │    │统计功能 │    │日志功能 │
│  订单   │    │ 商品推荐│    │ 数据报表│    │ 审计日志│
└────┬────┘    └────┬────┘    └────┬────┘    └────┬────┘
     │              │              │              │
     └──────────────┴──────────────┴──────────────┘
                        │
                   ┌─────────┐
                   │ 数据库  │
                   └─────────┘

降级模式:
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ 核心功能 │    │推荐功能 │    │统计功能 │    │日志功能 │
│  订单   │◄───│  关闭  │    │ 异步处理│    │ 采样记录│
│ (保留)  │    │ (降级) │    │ (降级) │    │ (降级) │
└────┬────┘    └─────────┘    └─────────┘    └─────────┘
     │
     ▼
┌─────────┐
│ 数据库  │  (减少连接池)
└─────────┘
```

#### 完整示例：Go实现降级策略

```go
// fallback/fallback.go
package fallback

import (
    "context"
    "encoding/json"
    "errors"
    "log"
    "time"
)

// FallbackFunc 降级函数类型
type FallbackFunc func(ctx context.Context, err error) (interface{}, error)

// Strategy 降级策略
type Strategy int

const (
    StrategyFailFast Strategy = iota     // 快速失败
    StrategyReturnDefault               // 返回默认值
    StrategyReturnCache                 // 返回缓存
    StrategyReturnEmpty                 // 返回空值
    StrategyAsyncExecute                // 异步执行
)

// Degrader 降级器
type Degrader struct {
    strategy     Strategy
    fallback     FallbackFunc
    defaultValue interface{}
    cache        Cache
}

// NewDegrader 创建降级器
func NewDegrader(strategy Strategy, opts ...DegraderOption) *Degrader {
    d := &Degrader{
        strategy: strategy,
    }

    for _, opt := range opts {
        opt(d)
    }

    return d
}

type DegraderOption func(*Degrader)

func WithFallback(fn FallbackFunc) DegraderOption {
    return func(d *Degrader) {
        d.fallback = fn
    }
}

func WithDefaultValue(v interface{}) DegraderOption {
    return func(d *Degrader) {
        d.defaultValue = v
    }
}

func WithCache(cache Cache) DegraderOption {
    return func(d *Degrader) {
        d.cache = cache
    }
}

// Execute 执行带降级的函数
func (d *Degrader) Execute(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
    result, err := fn()
    if err == nil {
        // 成功时更新缓存
        if d.cache != nil {
            d.cache.Set(ctx, "fallback_key", result, 5*time.Minute)
        }
        return result, nil
    }

    // 失败时执行降级
    log.Printf("Primary function failed: %v, executing fallback", err)

    switch d.strategy {
    case StrategyFailFast:
        return nil, err

    case StrategyReturnDefault:
        if d.defaultValue != nil {
            return d.defaultValue, nil
        }
        return nil, errors.New("no default value configured")

    case StrategyReturnCache:
        if d.cache != nil {
            cached, cacheErr := d.cache.Get(ctx, "fallback_key")
            if cacheErr == nil {
                return cached, nil
            }
        }
        return nil, err

    case StrategyReturnEmpty:
        return nil, nil

    case StrategyAsyncExecute:
        if d.fallback != nil {
            // 异步执行降级
            go func() {
                _, _ = d.fallback(ctx, err)
            }()
        }
        return nil, err

    default:
        if d.fallback != nil {
            return d.fallback(ctx, err)
        }
        return nil, err
    }
}

// 缓存接口
type Cache interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// 实际应用示例：电商系统降级

type ProductService struct {
    db       *sql.DB
    cache    Cache
    degrader *Degrader
}

func NewProductService(db *sql.DB, cache Cache) *ProductService {
    return &ProductService{
        db:    db,
        cache: cache,
        degrader: NewDegrader(
            StrategyReturnCache,
            WithCache(cache),
        ),
    }
}

// GetProduct 获取商品信息（带降级）
func (s *ProductService) GetProduct(ctx context.Context, productID string) (*Product, error) {
    result, err := s.degrader.Execute(ctx, func() (interface{}, error) {
        // 从数据库获取
        return s.getProductFromDB(ctx, productID)
    })

    if err != nil {
        return nil, err
    }

    return result.(*Product), nil
}

func (s *ProductService) getProductFromDB(ctx context.Context, productID string) (*Product, error) {
    var product Product
    err := s.db.QueryRowContext(ctx,
        "SELECT id, name, price, stock FROM products WHERE id = ?",
        productID,
    ).Scan(&product.ID, &product.Name, &product.Price, &product.Stock)

    if err != nil {
        return nil, err
    }

    return &product, nil
}

// GetRecommendations 获取推荐商品（带降级）
func (s *ProductService) GetRecommendations(ctx context.Context, userID string) ([]Product, error) {
    degrader := NewDegrader(
        StrategyReturnDefault,
        WithDefaultValue([]Product{
            {ID: "hot1", Name: "热销商品1", Price: 99.0},
            {ID: "hot2", Name: "热销商品2", Price: 199.0},
        }),
    )

    result, err := degrader.Execute(ctx, func() (interface{}, error) {
        // 调用推荐算法服务
        return s.getRecommendationsFromService(ctx, userID)
    })

    if err != nil {
        return nil, err
    }

    return result.([]Product), nil
}

// 多级降级策略

type MultiLevelDegrader struct {
    levels []DegraderLevel
}

type DegraderLevel struct {
    Condition func(error) bool
    Degrader  *Degrader
}

func (m *MultiLevelDegrader) Execute(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
    result, err := fn()
    if err == nil {
        return result, nil
    }

    // 根据错误类型选择降级级别
    for _, level := range m.levels {
        if level.Condition(err) {
            return level.Degrader.Execute(ctx, fn)
        }
    }

    return nil, err
}

// 使用示例
func ExampleMultiLevelDegrader() {
    mld := &MultiLevelDegrader{
        levels: []DegraderLevel{
            {
                // 超时错误：返回缓存
                Condition: func(err error) bool {
                    return errors.Is(err, context.DeadlineExceeded)
                },
                Degrader: NewDegrader(StrategyReturnCache, WithCache(nil)),
            },
            {
                // 连接错误：返回默认值
                Condition: func(err error) bool {
                    return errors.Is(err, ErrConnectionRefused)
                },
                Degrader: NewDegrader(StrategyReturnDefault, WithDefaultValue("default")),
            },
            {
                // 其他错误：快速失败
                Condition: func(err error) bool {
                    return true
                },
                Degrader: NewDegrader(StrategyFailFast),
            },
        },
    }

    result, err := mld.Execute(context.Background(), func() (interface{}, error) {
        // 业务逻辑
        return nil, nil
    })

    _ = result
    _ = err
}

var ErrConnectionRefused = errors.New("connection refused")
```

#### 反例说明

```go
// ❌ 错误：核心功能降级
// 订单创建功能设置降级返回默认值
// 问题：核心业务不能降级，否则影响数据一致性

// ❌ 错误：降级不记录日志
// 降级发生时没有任何记录
// 问题：无法追踪降级发生频率和原因

// ❌ 错误：缓存和主数据源不一致
// 降级返回的缓存数据已过期
// 问题：用户看到错误信息
```

#### 注意事项

1. **核心功能不降级**：订单创建、支付等核心功能不能降级
2. **降级数据一致性**：降级返回的数据应与正常数据格式一致
3. **降级监控**：记录降级发生次数和原因
4. **自动恢复**：主服务恢复后应自动切换回正常模式

---

### 4.3 hystrix-go实现

#### 概念定义

hystrix-go是Netflix Hystrix的Go语言实现，提供熔断、降级、监控等功能，通过**Command模式**封装受保护的调用。

#### 完整示例：hystrix-go使用

```go
// hystrix/hystrix_example.go
package hystrix

import (
    "context"
    "fmt"
    "time"

    "github.com/afex/hystrix-go/hystrix"
)

// 配置Hystrix
func init() {
    // 全局配置
    hystrix.ConfigureCommand("get_user", hystrix.CommandConfig{
        Timeout:                1000,  // 超时时间(ms)
        MaxConcurrentRequests:  100,   // 最大并发数
        RequestVolumeThreshold: 10,    // 触发熔断的最小请求数
        SleepWindow:            5000,  // 熔断后休眠时间(ms)
        ErrorPercentThreshold:  50,    // 错误率阈值(%)
    })
}

// UserService 用户服务
type UserService struct {
    client *http.Client
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    var user User

    // 使用Hystrix执行
    err := hystrix.Do("get_user", func() error {
        // 实际调用
        resp, err := s.client.Get(fmt.Sprintf("http://user-service/users/%s", userID))
        if err != nil {
            return err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("unexpected status: %d", resp.StatusCode)
        }

        return json.NewDecoder(resp.Body).Decode(&user)
    }, func(err error) error {
        // 降级函数
        fmt.Printf("GetUser fallback executed: %v\n", err)

        // 返回缓存数据或默认值
        user = User{
            ID:   userID,
            Name: "Unknown",
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    return &user, nil
}

// 使用Go方法异步执行
func (s *UserService) GetUserAsync(ctx context.Context, userID string) (*User, error) {
    output := make(chan *User, 1)
    errors := hystrix.Go("get_user", func() error {
        user, err := s.getUserFromRemote(ctx, userID)
        if err != nil {
            return err
        }
        output <- user
        return nil
    }, func(err error) error {
        // 降级
        output <- &User{ID: userID, Name: "Unknown"}
        return nil
    })

    select {
    case user := <-output:
        return user, nil
    case err := <-errors:
        return nil, err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// 批量请求使用舱壁隔离
func (s *UserService) GetUsers(ctx context.Context, userIDs []string) ([]*User, error) {
    var wg sync.WaitGroup
    users := make([]*User, len(userIDs))
    errChan := make(chan error, len(userIDs))

    // 使用信号量限制并发
    sem := make(chan struct{}, 10) // 最多10个并发

    for i, userID := range userIDs {
        wg.Add(1)
        go func(index int, id string) {
            defer wg.Done()

            sem <- struct{}{}        // 获取信号量
            defer func() { <-sem }() // 释放信号量

            user, err := s.GetUser(ctx, id)
            if err != nil {
                errChan <- err
                return
            }
            users[index] = user
        }(i, userID)
    }

    wg.Wait()
    close(errChan)

    // 检查是否有错误
    for err := range errChan {
        if err != nil {
            return users, err
        }
    }

    return users, nil
}

// 监控和统计
func StartHystrixStream() {
    // 启动Hystrix Dashboard数据流
    hystrixStreamHandler := hystrix.NewStreamHandler()
    hystrixStreamHandler.Start()

    go http.ListenAndServe(":8081", hystrixStreamHandler)
}

// 自定义指标收集
type MetricsCollector struct{}

func (m *MetricsCollector) Update(cmd string, eventType string, duration time.Duration) {
    // 发送到Prometheus或其他监控系统
    fmt.Printf("Command: %s, Event: %s, Duration: %v\n", cmd, eventType, duration)
}

func (m *MetricsCollector) Reset(cmd string) {
    fmt.Printf("Reset metrics for command: %s\n", cmd)
}

func init() {
    // 设置自定义指标收集器
    // hystrix.SetMetricsCollector(&MetricsCollector{})
}
```

#### 反例说明

```go
// ❌ 错误：所有命令使用相同配置
hystrix.Do("command1", fn, fallback)
hystrix.Do("command2", fn, fallback)  // 使用相同配置
// 问题：不同命令特性不同，应分别配置

// ❌ 错误：降级函数抛出panic
hystrix.Do("cmd", fn, func(err error) error {
    panic("fallback error")  // 会导致程序崩溃
    return nil
})

// ❌ 错误：不处理超时
Timeout: 100000  // 100秒，太长
// 问题：无法及时熔断
```

---

### 4.4 sentinel-go实现

#### 概念定义

sentinel-go是阿里巴巴开源的流量控制组件，提供**流量控制**、**熔断降级**、**系统保护**等功能，支持丰富的规则和灵活的扩展。

#### 完整示例：sentinel-go使用

```go
// sentinel/sentinel_example.go
package sentinel

import (
    "context"
    "fmt"

    sentinel "github.com/alibaba/sentinel-golang/api"
    "github.com/alibaba/sentinel-golang/core/circuitbreaker"
    "github.com/alibaba/sentinel-golang/core/flow"
    "github.com/alibaba/sentinel-golang/core/hotspot"
    "github.com/alibaba/sentinel-golang/core/system"
    "github.com/alibaba/sentinel-golang/logging"
)

// 初始化Sentinel
func InitSentinel() error {
    err := sentinel.InitDefault()
    if err != nil {
        return err
    }

    // 配置日志
    logging.ResetGlobalLoggerLevel(logging.InfoLevel)

    return nil
}

// 流量控制规则
func LoadFlowRules() error {
    _, err := flow.LoadRules([]*flow.Rule{
        {
            Resource:               "get_user",
            TokenCalculateStrategy: flow.Direct,
            ControlBehavior:        flow.Reject,
            Threshold:              100,  // QPS限制
            StatIntervalInMs:       1000,
        },
        {
            Resource:               "create_order",
            TokenCalculateStrategy: flow.WarmUp,  // 预热模式
            ControlBehavior:        flow.Reject,
            Threshold:              50,
            WarmUpPeriodSec:        10,
        },
    })
    return err
}

// 熔断规则
func LoadCircuitBreakerRules() error {
    _, err := circuitbreaker.LoadRules([]*circuitbreaker.Rule{
        {
            Resource:         "get_user",
            Strategy:         circuitbreaker.ErrorRatio,  // 错误率熔断
            RetryTimeoutMs:   3000,  // 熔断持续时间
            MinRequestAmount: 10,     // 最小请求数
            StatIntervalMs:   10000,  // 统计周期
            Threshold:        0.5,    // 错误率阈值50%
        },
        {
            Resource:         "slow_query",
            Strategy:         circuitbreaker.SlowRequestRatio,  // 慢调用熔断
            RetryTimeoutMs:   5000,
            MinRequestAmount: 10,
            StatIntervalMs:   10000,
            MaxAllowedRtMs:   500,    // 最大响应时间
            Threshold:        0.8,    // 慢调用比例阈值
        },
    })
    return err
}

// 热点参数限流
func LoadHotSpotRules() error {
    _, err := hotspot.LoadRules([]*hotspot.Rule{
        {
            Resource:        "get_product",
            MetricType:      hotspot.Concurrency,  // 并发数限流
            ParamIndex:      0,                    // 第一个参数(productID)
            Threshold:       10,
            DurationInSec:   1,
        },
    })
    return err
}

// 系统保护规则
func LoadSystemRules() error {
    _, err := system.LoadRules([]*system.Rule{
        {
            MetricType:   system.Load,
            TriggerCount: 8.0,  // 系统负载阈值
            Strategy:     system.BBR,
        },
        {
            MetricType:   system.CpuUsage,
            TriggerCount: 0.8,  // CPU使用率阈值
            Strategy:     system.BBR,
        },
    })
    return err
}

// 使用Entry进行流量控制
func GetUser(ctx context.Context, userID string) (*User, error) {
    // 获取Entry
    entry, err := sentinel.Entry("get_user", sentinel.WithTrafficType(base.Inbound))
    if err != nil {
        // 被限流或熔断
        if err == sentinel.ErrBlock {
            // 执行降级逻辑
            return getUserFromCache(userID)
        }
        return nil, err
    }
    defer entry.Exit()  // 确保Exit被调用

    // 执行业务逻辑
    user, err := getUserFromDB(userID)
    if err != nil {
        // 记录错误
        sentinel.TraceError(entry, err)
        return nil, err
    }

    return user, nil
}

// 带热点参数的Entry
func GetProduct(ctx context.Context, productID string) (*Product, error) {
    entry, err := sentinel.Entry(
        "get_product",
        sentinel.WithTrafficType(base.Inbound),
        sentinel.WithArgs(productID),  // 热点参数
    )
    if err != nil {
        if err == sentinel.ErrBlock {
            return getProductFromCache(productID)
        }
        return nil, err
    }
    defer entry.Exit()

    return getProductFromDB(productID)
}

// 自定义Slot扩展
type CustomSlot struct{}

func (s *CustomSlot) Order() uint32 {
    return 100
}

func (s *CustomSlot) Check(ctx *base.EntryContext) *base.TokenResult {
    // 自定义检查逻辑
    fmt.Printf("Custom check for resource: %s\n", ctx.Resource.Name())
    return nil
}

func init() {
    // 注册自定义Slot
    // slotChain.AddSlot(&CustomSlot{})
}

// 监控和指标
func StartMetricsExporter() {
    // 导出到Prometheus
    // sentinel_exporter.InitSentinelExporter()

    // 或使用自定义导出
    go func() {
        ticker := time.NewTicker(10 * time.Second)
        defer ticker.Stop()

        for range ticker.C {
            // 获取统计信息
            // metrics := sentinel.GetMetrics()
            // 发送到监控系统
        }
    }()
}
```

#### 反例说明

```go
// ❌ 错误：忘记调用Exit
entry, err := sentinel.Entry("resource")
if err != nil {
    return err
}
// defer entry.Exit()  // 忘记这行会导致统计错误

// ❌ 错误：不处理ErrBlock
entry, err := sentinel.Entry("resource")
if err != nil {
    return err  // 直接返回，没有降级处理
}

// ❌ 错误：规则配置不合理
Threshold: 0.01  // 错误率阈值1%，太敏感
RetryTimeoutMs: 100  // 熔断后100ms就恢复，太短
```

#### 注意事项

1. **规则热更新**：支持动态更新规则，无需重启
2. **多种熔断策略**：错误率、慢调用比例、错误数
3. **系统保护**：自动保护系统资源不被耗尽
4. **监控集成**：支持Prometheus等监控系统

---

## 5. 限流与配额

### 5.1 令牌桶算法

#### 概念定义

令牌桶算法以固定速率向桶中添加令牌，请求需要获取令牌才能执行。支持**突发流量**处理，是限流最常用的算法。

#### 架构图

```
令牌桶工作原理:

令牌生成速率: 10个/秒
桶容量: 100个令牌

时间线:
t=0:  桶满(100令牌) ──> 请求A消耗1令牌(剩99)
t=1s: 添加10令牌 ─────> 桶(100令牌，上限)
t=1s: 请求B消耗1令牌 ──> 桶(99令牌)
t=2s: 添加10令牌 ─────> 桶(100令牌)
...
突发流量:
同时100个请求 ──> 消耗100令牌 ──> 桶空
后续请求 ──────> 等待令牌生成

实现结构:
┌─────────────────────────────────┐
│           TokenBucket           │
│  ┌─────────────────────────┐   │
│  │  capacity: 100          │   │
│  │  tokens: 50             │   │
│  │  rate: 10/s             │   │
│  │  lastUpdate: timestamp  │   │
│  └─────────────────────────┘   │
│                                 │
│  Allow() bool                   │
│  AllowN(n int) bool             │
└─────────────────────────────────┘
```

#### 完整示例：Go实现令牌桶

```go
// ratelimit/token_bucket.go
package ratelimit

import (
    "context"
    "sync"
    "time"
)

// TokenBucket 令牌桶
type TokenBucket struct {
    capacity  int64         // 桶容量
    tokens    int64         // 当前令牌数
    rate      float64       // 令牌生成速率(个/秒)
    lastTime  time.Time     // 上次更新时间
    mu        sync.Mutex
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(capacity int64, rate float64) *TokenBucket {
    return &TokenBucket{
        capacity: capacity,
        tokens:   capacity, // 初始满桶
        rate:     rate,
        lastTime: time.Now(),
    }
}

// Allow 获取1个令牌
func (tb *TokenBucket) Allow() bool {
    return tb.AllowN(1)
}

// AllowN 获取N个令牌
func (tb *TokenBucket) AllowN(n int64) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastTime).Seconds()
    tb.lastTime = now

    // 添加新令牌
    newTokens := int64(elapsed * tb.rate)
    tb.tokens = min(tb.tokens+newTokens, tb.capacity)

    // 检查是否足够
    if tb.tokens >= n {
        tb.tokens -= n
        return true
    }

    return false
}

// Wait 等待直到获取令牌
func (tb *TokenBucket) Wait(ctx context.Context) error {
    return tb.WaitN(ctx, 1)
}

// WaitN 等待直到获取N个令牌
func (tb *TokenBucket) WaitN(ctx context.Context, n int64) error {
    for {
        if tb.AllowN(n) {
            return nil
        }

        // 计算需要等待的时间
        tb.mu.Lock()
        needed := n - tb.tokens
        waitTime := time.Duration(float64(needed) / tb.rate * float64(time.Second))
        tb.mu.Unlock()

        select {
        case <-time.After(waitTime):
            continue
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

// Reserve 预留令牌
func (tb *TokenBucket) Reserve() *Reservation {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastTime).Seconds()
    tb.lastTime = now

    // 添加新令牌
    newTokens := int64(elapsed * tb.rate)
    tb.tokens = min(tb.tokens+newTokens, tb.capacity)

    if tb.tokens >= 1 {
        tb.tokens--
        return &Reservation{OK: true, Delay: 0}
    }

    // 计算需要等待的时间
    needed := 1 - tb.tokens
    delay := time.Duration(float64(needed) / tb.rate * float64(time.Second))

    return &Reservation{OK: true, Delay: delay}
}

// Reservation 预留结果
type Reservation struct {
    OK    bool
    Delay time.Duration
}

func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

// RateLimiter 限流器接口
type RateLimiter interface {
    Allow() bool
    AllowN(n int64) bool
}

// MultiLevelRateLimiter 多级限流器
type MultiLevelRateLimiter struct {
    limiters []RateLimiter
}

func NewMultiLevelRateLimiter(limiters ...RateLimiter) *MultiLevelRateLimiter {
    return &MultiLevelRateLimiter{limiters: limiters}
}

func (m *MultiLevelRateLimiter) Allow() bool {
    for _, l := range m.limiters {
        if !l.Allow() {
            return false
        }
    }
    return true
}

func (m *MultiLevelRateLimiter) AllowN(n int64) bool {
    for _, l := range m.limiters {
        if !l.AllowN(n) {
            return false
        }
    }
    return true
}

// HTTP中间件集成
func TokenBucketMiddleware(tb *TokenBucket) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !tb.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// 基于用户ID的限流
func PerUserRateLimiter(rate float64, capacity int64) func(http.Handler) http.Handler {
    buckets := make(map[string]*TokenBucket)
    var mu sync.RWMutex

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userID := r.Header.Get("X-User-ID")
            if userID == "" {
                http.Error(w, "Missing user ID", http.StatusBadRequest)
                return
            }

            mu.RLock()
            bucket, ok := buckets[userID]
            mu.RUnlock()

            if !ok {
                mu.Lock()
                bucket = NewTokenBucket(capacity, rate)
                buckets[userID] = bucket
                mu.Unlock()
            }

            if !bucket.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

// 使用golang.org/x/time/rate
import "golang.org/x/time/rate"

func StandardTokenBucket() *rate.Limiter {
    // 每秒10个令牌，桶容量100
    return rate.NewLimiter(rate.Limit(10), 100)
}
```

#### 反例说明

```go
// ❌ 错误：不更新lastTime
func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    // 没有更新lastTime，导致令牌计算错误
    elapsed := time.Since(tb.lastTime).Seconds()
    newTokens := int64(elapsed * tb.rate)
    // ...
}

// ❌ 错误：桶容量为0
tb := NewTokenBucket(0, 10)
// 问题：永远无法获取令牌

// ❌ 错误：浮点数精度问题
newTokens := elapsed * tb.rate  // elapsed是float64
// 长时间运行后可能累积精度误差
```

#### 注意事项

1. **精度问题**：使用整数计数避免浮点数精度问题
2. **懒加载**：只在Allow时计算令牌数，避免定时任务
3. **并发安全**：使用锁或原子操作保证并发安全
4. **桶容量设置**：根据业务突发流量需求设置

---

### 5.2 漏桶算法

#### 概念定义

漏桶算法以固定速率处理请求，请求先进入桶中排队，桶满则拒绝新请求。保证**输出速率恒定**，适合需要严格平滑流量的场景。

#### 架构图

```
漏桶工作原理:

请求流入(可变速率):
  ↓ ↓↓ ↓  ↓↓↓   ↓
┌─────────────────┐
│     漏桶        │  桶容量: 100
│  ┌───────────┐  │  流出速率: 10个/秒
│  │ 请求队列   │  │
│  │ [req1]    │  │
│  │ [req2]    │  │
│  │ [req3]    │  │
│  │ ...       │  │
│  └─────┬─────┘  │
└────────┼────────┘
         │ 固定速率流出
         ▼
    请求处理

vs 令牌桶:
- 漏桶: 请求排队，流出速率固定，不支持突发
- 令牌桶: 请求不排队，支持突发，流入速率可变
```

#### 完整示例：Go实现漏桶

```go
// ratelimit/leaky_bucket.go
package ratelimit

import (
    "context"
    "sync"
    "time"
)

// LeakyBucket 漏桶
type LeakyBucket struct {
    capacity int           // 桶容量
    rate     time.Duration // 请求处理间隔
    queue    chan struct{} // 请求队列
    mu       sync.Mutex
    stopChan chan struct{}
}

// NewLeakyBucket 创建漏桶
func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
    lb := &LeakyBucket{
        capacity: capacity,
        rate:     rate,
        queue:    make(chan struct{}, capacity),
        stopChan: make(chan struct{}),
    }

    // 启动漏桶处理
    go lb.leak()

    return lb
}

// leak 漏桶核心：以固定速率处理请求
func (lb *LeakyBucket) leak() {
    ticker := time.NewTicker(lb.rate)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            select {
            case <-lb.queue:
                // 处理一个请求
            default:
                // 队列为空
            }
        case <-lb.stopChan:
            return
        }
    }
}

// Allow 尝试将请求加入桶中
func (lb *LeakyBucket) Allow() bool {
    select {
    case lb.queue <- struct{}{}:
        return true
    default:
        // 桶已满
        return false
    }
}

// Wait 等待直到请求可以加入桶中
func (lb *LeakyBucket) Wait(ctx context.Context) error {
    select {
    case lb.queue <- struct{}{}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Stop 停止漏桶
func (lb *LeakyBucket) Stop() {
    close(lb.stopChan)
}

// LeakyBucketV2 改进版漏桶（支持动态调整速率）
type LeakyBucketV2 struct {
    capacity   int
    rate       int64  // 每秒处理请求数
    water      int64  // 当前水位（队列中的请求数）
    lastLeak   int64  // 上次漏水时间戳(毫秒)
    mu         sync.Mutex
}

func NewLeakyBucketV2(capacity int, rate int64) *LeakyBucketV2 {
    return &LeakyBucketV2{
        capacity: capacity,
        rate:     rate,
        lastLeak: time.Now().UnixMilli(),
    }
}

func (lb *LeakyBucketV2) Allow() bool {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    now := time.Now().UnixMilli()
    elapsed := now - lb.lastLeak

    // 计算漏出的水量
    leaked := elapsed * lb.rate / 1000

    if leaked > 0 {
        lb.water = max(0, lb.water-leaked)
        lb.lastLeak = now
    }

    // 检查是否可以加入
    if lb.water < int64(lb.capacity) {
        lb.water++
        return true
    }

    return false
}

func max(a, b int64) int64 {
    if a > b {
        return a
    }
    return b
}

// 漏桶vs令牌桶选择
/*
使用漏桶的场景:
1. 需要严格平滑输出流量
2. 下游服务处理能力有限
3. 需要请求排队而非拒绝

使用令牌桶的场景:
1. 需要支持突发流量
2. 希望快速拒绝而非排队
3. 需要更灵活的限流策略
*/
```

#### 反例说明

```go
// ❌ 错误：漏桶速率设置不合理
rate: 1 * time.Millisecond  // 每秒1000个请求
// 问题：可能超过下游处理能力

// ❌ 错误：不处理队列积压
// 长时间高流量导致队列无限增长
// 问题：内存溢出

// ❌ 错误：漏桶goroutine泄漏
func Handler() {
    lb := NewLeakyBucket(100, time.Second)
    // 使用完不调用Stop()
}
// 问题：goroutine泄漏
```

---

### 5.3 滑动窗口

#### 概念定义

滑动窗口算法在固定时间窗口内统计请求数，相比固定窗口能更精确地控制流量，避免**窗口边界突发**问题。

#### 架构图

```
固定窗口问题:

窗口大小: 1分钟, 限制: 100请求

时间线:
|:00      |:30      |:00
|─────────|─────────|
    50请求    50请求

在:30和:00边界处:
  50请求(前窗口末尾) + 50请求(后窗口开头) = 100请求/秒
  实际在1秒内处理了100请求，远超限制

滑动窗口解决:

当前时间: 12:00:30
窗口: [11:59:30, 12:00:30]

统计:
11:59:30-11:59:45: 30请求
11:59:45-12:00:00: 40请求
12:00:00-12:00:15: 20请求
12:00:15-12:00:30: 10请求
─────────────────────────
总计: 100请求

精确控制，无边界问题
```

#### 完整示例：Go实现滑动窗口

```go
// ratelimit/sliding_window.go
package ratelimit

import (
    "sync"
    "time"
)

// SlidingWindow 滑动窗口
type SlidingWindow struct {
    windowSize time.Duration // 窗口大小
    limit      int           // 窗口内最大请求数
    buckets    []int         // 子窗口计数
    bucketSize time.Duration // 子窗口大小
    currentIdx int           // 当前子窗口索引
    lastUpdate time.Time     // 上次更新时间
    mu         sync.Mutex
}

// NewSlidingWindow 创建滑动窗口
func NewSlidingWindow(windowSize time.Duration, limit int, bucketCount int) *SlidingWindow {
    return &SlidingWindow{
        windowSize: windowSize,
        limit:      limit,
        buckets:    make([]int, bucketCount),
        bucketSize: windowSize / time.Duration(bucketCount),
        lastUpdate: time.Now(),
    }
}

// Allow 检查是否允许请求
func (sw *SlidingWindow) Allow() bool {
    sw.mu.Lock()
    defer sw.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(sw.lastUpdate)

    // 计算需要滑动的子窗口数
    bucketCount := len(sw.buckets)
    shiftBuckets := int(elapsed / sw.bucketSize)

    if shiftBuckets > 0 {
        // 滑动窗口：清空过期的子窗口
        for i := 0; i < shiftBuckets && i < bucketCount; i++ {
            sw.currentIdx = (sw.currentIdx + 1) % bucketCount
            sw.buckets[sw.currentIdx] = 0
        }
        sw.lastUpdate = now
    }

    // 计算当前窗口内的总请求数
    total := 0
    for _, count := range sw.buckets {
        total += count
    }

    // 检查是否超过限制
    if total >= sw.limit {
        return false
    }

    // 增加当前子窗口计数
    sw.buckets[sw.currentIdx]++
    return true
}

// GetCount 获取当前窗口内的请求数
func (sw *SlidingWindow) GetCount() int {
    sw.mu.Lock()
    defer sw.mu.Unlock()

    total := 0
    for _, count := range sw.buckets {
        total += count
    }
    return total
}

// SlidingWindowLog 滑动窗口日志实现（更精确但内存占用大）
type SlidingWindowLog struct {
    windowSize time.Duration
    limit      int
    timestamps []time.Time
    mu         sync.Mutex
}

func NewSlidingWindowLog(windowSize time.Duration, limit int) *SlidingWindowLog {
    return &SlidingWindowLog{
        windowSize: windowSize,
        limit:      limit,
        timestamps: make([]time.Time, 0, limit),
    }
}

func (sw *SlidingWindowLog) Allow() bool {
    sw.mu.Lock()
    defer sw.mu.Unlock()

    now := time.Now()
    windowStart := now.Add(-sw.windowSize)

    // 清理过期的时间戳
    validIdx := 0
    for i, ts := range sw.timestamps {
        if ts.After(windowStart) {
            validIdx = i
            break
        }
    }
    sw.timestamps = sw.timestamps[validIdx:]

    // 检查是否超过限制
    if len(sw.timestamps) >= sw.limit {
        return false
    }

    // 记录当前请求
    sw.timestamps = append(sw.timestamps, now)
    return true
}

// 滑动窗口计数器（Redis实现）
type RedisSlidingWindow struct {
    client *redis.Client
    window time.Duration
    limit  int
}

func (r *RedisSlidingWindow) Allow(ctx context.Context, key string) (bool, error) {
    now := time.Now().UnixMilli()
    windowStart := now - r.window.Milliseconds()

    pipe := r.client.Pipeline()

    // 移除窗口外的请求记录
    pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

    // 获取当前窗口内的请求数
    countCmd := pipe.ZCard(ctx, key)

    // 添加当前请求
    pipe.ZAdd(ctx, key, &redis.Z{
        Score:  float64(now),
        Member: now,
    })

    // 设置过期时间
    pipe.Expire(ctx, key, r.window)

    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }

    count := countCmd.Val()
    return count < int64(r.limit), nil
}
```

#### 反例说明

```go
// ❌ 错误：子窗口数太少
NewSlidingWindow(time.Minute, 100, 2)
// 问题：精度不够，仍有边界问题

// ❌ 错误：不清理过期数据
// 滑动窗口日志实现不清理过期时间戳
// 问题：内存无限增长

// ❌ 错误：并发不安全
func (sw *SlidingWindow) Allow() bool {
    // 没有加锁
    total := 0
    for _, count := range sw.buckets {
        total += count
    }
    // ...
}
```

---

### 5.4 分布式限流

#### 概念定义

分布式限流在集群层面统一控制流量，需要**集中式存储**（如Redis）协调各节点的限流状态，保证全局一致性。

#### 架构图

```
分布式令牌桶:

┌─────────┐     ┌─────────┐     ┌─────────┐
│ 服务节点A │     │ 服务节点B │     │ 服务节点C │
│ :8081   │     │ :8082   │     │ :8083   │
└────┬────┘     └────┬────┘     └────┬────┘
     │               │               │
     └───────────────┼───────────────┘
                     │
              ┌─────────────┐
              │    Redis    │
              │  令牌桶状态  │
              │  key: tokens│
              │  value: 50  │
              └─────────────┘

Redis Lua脚本保证原子性:
1. 获取当前令牌数
2. 计算新增令牌
3. 判断是否足够
4. 扣减令牌
5. 返回结果

所有操作在Redis单线程执行，保证原子性
```

#### 完整示例：Go实现分布式限流

```go
// ratelimit/distributed.go
package ratelimit

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedisTokenBucket 基于Redis的分布式令牌桶
type RedisTokenBucket struct {
    client   *redis.Client
    key      string
    capacity int64
    rate     float64
}

// Lua脚本：原子性获取令牌
const tokenBucketLua = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local requested = tonumber(ARGV[3])
local now = tonumber(ARGV[4])

local bucket = redis.call('HMGET', key, 'tokens', 'last_update')
local tokens = tonumber(bucket[1]) or capacity
local lastUpdate = tonumber(bucket[2]) or now

-- 计算新增令牌
local elapsed = now - lastUpdate
local newTokens = math.min(capacity, tokens + elapsed * rate / 1000)

-- 判断是否足够
if newTokens >= requested then
    newTokens = newTokens - requested
    redis.call('HMSET', key, 'tokens', newTokens, 'last_update', now)
    redis.call('EXPIRE', key, 60)
    return 1
else
    redis.call('HMSET', key, 'tokens', newTokens, 'last_update', now)
    redis.call('EXPIRE', key, 60)
    return 0
end
`

func NewRedisTokenBucket(client *redis.Client, key string, capacity int64, rate float64) *RedisTokenBucket {
    return &RedisTokenBucket{
        client:   client,
        key:      key,
        capacity: capacity,
        rate:     rate,
    }
}

func (r *RedisTokenBucket) Allow(ctx context.Context) (bool, error) {
    return r.AllowN(ctx, 1)
}

func (r *RedisTokenBucket) AllowN(ctx context.Context, n int64) (bool, error) {
    now := time.Now().UnixMilli()

    result, err := r.client.Eval(ctx, tokenBucketLua, []string{r.key},
        r.capacity, r.rate, n, now,
    ).Result()

    if err != nil {
        return false, err
    }

    allowed, ok := result.(int64)
    if !ok {
        return false, fmt.Errorf("unexpected result type")
    }

    return allowed == 1, nil
}

// Redis滑动窗口限流
const slidingWindowLua = `
local key = KEYS[1]
local window = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local windowStart = now - window

-- 移除窗口外的成员
redis.call('ZREMRANGEBYSCORE', key, 0, windowStart)

-- 获取当前窗口内的成员数
local count = redis.call('ZCARD', key)

-- 判断是否超过限制
if count >= limit then
    return 0
end

-- 添加当前请求
redis.call('ZADD', key, now, now)
redis.call('EXPIRE', key, math.ceil(window / 1000))

return 1
`

type RedisSlidingWindow struct {
    client *redis.Client
    window time.Duration
    limit  int
}

func (r *RedisSlidingWindow) Allow(ctx context.Context, key string) (bool, error) {
    now := time.Now().UnixMilli()

    result, err := r.client.Eval(ctx, slidingWindowLua, []string{key},
        r.window.Milliseconds(), r.limit, now,
    ).Result()

    if err != nil {
        return false, err
    }

    allowed, _ := result.(int64)
    return allowed == 1, nil
}

// 基于Redis Cell的限流（Redis 4.0+）
// 使用Redis Module: redis-cell
const redisCellLua = `
return redis.call('CL.THROTTLE', KEYS[1], ARGV[1], ARGV[2], ARGV[3], ARGV[4])
`

type RedisCellRateLimiter struct {
    client *redis.Client
}

func (r *RedisCellRateLimiter) Allow(ctx context.Context, key string, capacity, rate, period int64) (bool, error) {
    // CL.THROTTLE key max_burst count_per_period period_in_seconds [quantity]
    result, err := r.client.Eval(ctx, redisCellLua, []string{key},
        capacity-1, rate, period, 1,
    ).Result()

    if err != nil {
        return false, err
    }

    // 返回数组: [allowed, limit_remaining, limit_reset_after, retry_after]
    results := result.([]interface{})
    allowed := results[0].(int64)

    return allowed == 1, nil
}

// 分布式限流管理器
type DistributedRateLimiterManager struct {
    client *redis.Client
    prefix string
}

func NewDistributedRateLimiterManager(client *redis.Client, prefix string) *DistributedRateLimiterManager {
    return &DistributedRateLimiterManager{
        client: client,
        prefix: prefix,
    }
}

func (m *DistributedRateLimiterManager) GetLimiter(resource string, capacity int64, rate float64) *RedisTokenBucket {
    key := fmt.Sprintf("%s:ratelimit:%s", m.prefix, resource)
    return NewRedisTokenBucket(m.client, key, capacity, rate)
}

func (m *DistributedRateLimiterManager) Allow(ctx context.Context, resource string, capacity int64, rate float64) (bool, error) {
    limiter := m.GetLimiter(resource, capacity, rate)
    return limiter.Allow(ctx)
}

// 多级分布式限流
type MultiLevelDistributedLimiter struct {
    local  *TokenBucket       // 本地限流
    remote *RedisTokenBucket  // 分布式限流
}

func (m *MultiLevelDistributedLimiter) Allow(ctx context.Context) (bool, error) {
    // 先检查本地限流
    if !m.local.Allow() {
        return false, nil
    }

    // 再检查分布式限流
    return m.remote.Allow(ctx)
}

// 使用示例
func ExampleDistributedRateLimit() {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    manager := NewDistributedRateLimiterManager(client, "myapp")

    // API限流: 100请求/秒
    allowed, err := manager.Allow(context.Background(), "api", 100, 100)
    if err != nil {
        panic(err)
    }

    if !allowed {
        fmt.Println("Rate limit exceeded")
        return
    }

    // 处理请求
}
```

#### 反例说明

```go
// ❌ 错误：不使用Lua脚本
// 多个Redis命令之间有竞争条件
// 问题：限流不准确

// ❌ 错误：Redis故障时无降级
if err != nil {
    return false, err  // 直接返回错误
}
// 问题：Redis故障导致服务不可用
// 应降级为本地限流或直接放行

// ❌ 错误：key不设置过期时间
// 问题：Redis内存无限增长
```

#### 注意事项

1. **Lua脚本原子性**：使用Lua脚本保证操作原子性
2. **Redis故障降级**：Redis不可用时降级为本地限流或直接放行
3. **key过期设置**：设置合理的过期时间避免内存泄漏
4. **精度权衡**：分布式限流有网络延迟，精度略低于本地限流

---

## 6. 重试与退避

### 6.1 指数退避

#### 概念定义

指数退避(Exponential Backoff)是重试间隔按指数增长的策略，避免**重试风暴**，给故障服务恢复时间。

#### 架构图

```
指数退避时间计算:

初始间隔: 100ms
最大间隔: 30s
乘数: 2

重试次数    间隔时间
   1        100ms   = 100ms * 2^0
   2        200ms   = 100ms * 2^1
   3        400ms   = 100ms * 2^2
   4        800ms   = 100ms * 2^3
   5        1.6s    = 100ms * 2^4
   6        3.2s    = 100ms * 2^5
   ...
   n        min(100ms * 2^(n-1), 30s)

线性退避 vs 指数退避:

线性:  100ms, 200ms, 300ms, 400ms, 500ms...
指数:  100ms, 200ms, 400ms, 800ms, 1.6s...

指数退避更快增加间隔，更适合网络抖动场景
```

#### 完整示例：Go实现指数退避

```go
// retry/backoff.go
package retry

import (
    "context"
    "math"
    "math/rand"
    "time"
)

// Backoff 退避策略接口
type Backoff interface {
    Next() time.Duration
    Reset()
}

// ExponentialBackoff 指数退避
type ExponentialBackoff struct {
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Multiplier      float64
    currentInterval time.Duration
}

// NewExponentialBackoff 创建指数退避
func NewExponentialBackoff(initial, max time.Duration, multiplier float64) *ExponentialBackoff {
    if multiplier <= 1 {
        multiplier = 2
    }

    return &ExponentialBackoff{
        InitialInterval: initial,
        MaxInterval:     max,
        Multiplier:      multiplier,
        currentInterval: initial,
    }
}

// Next 获取下一次退避时间
func (e *ExponentialBackoff) Next() time.Duration {
    interval := e.currentInterval

    // 计算下一次间隔
    next := time.Duration(float64(e.currentInterval) * e.Multiplier)
    if next > e.MaxInterval {
        next = e.MaxInterval
    }
    e.currentInterval = next

    return interval
}

// Reset 重置退避
func (e *ExponentialBackoff) Reset() {
    e.currentInterval = e.InitialInterval
}

// LinearBackoff 线性退避
type LinearBackoff struct {
    InitialInterval time.Duration
    MaxInterval     time.Duration
    Increment       time.Duration
    currentInterval time.Duration
}

func NewLinearBackoff(initial, max, increment time.Duration) *LinearBackoff {
    return &LinearBackoff{
        InitialInterval: initial,
        MaxInterval:     max,
        Increment:       increment,
        currentInterval: initial,
    }
}

func (l *LinearBackoff) Next() time.Duration {
    interval := l.currentInterval

    next := l.currentInterval + l.Increment
    if next > l.MaxInterval {
        next = l.MaxInterval
    }
    l.currentInterval = next

    return interval
}

func (l *LinearBackoff) Reset() {
    l.currentInterval = l.InitialInterval
}

// FixedBackoff 固定间隔退避
type FixedBackoff struct {
    Interval time.Duration
}

func NewFixedBackoff(interval time.Duration) *FixedBackoff {
    return &FixedBackoff{Interval: interval}
}

func (f *FixedBackoff) Next() time.Duration {
    return f.Interval
}

func (f *FixedBackoff) Reset() {}

// NoBackoff 无退避（立即重试）
type NoBackoff struct{}

func (n *NoBackoff) Next() time.Duration {
    return 0
}

func (n *NoBackoff) Reset() {}
```

#### 反例说明

```go
// ❌ 错误：退避间隔无上限
for i := 0; i < retries; i++ {
    time.Sleep(time.Duration(math.Pow(2, float64(i))) * 100 * time.Millisecond)
}
// 问题：重试次数多时，间隔可能达到数小时

// ❌ 错误：乘数小于等于1
Multiplier: 1.0
// 问题：退避不增长，等同于固定间隔

// ❌ 错误：初始间隔太长
InitialInterval: 10 * time.Second
// 问题：轻微故障也需要等待10秒才重试
```

---

### 6.2 抖动（Jitter）

#### 概念定义

抖动是在退避时间上添加随机偏移，避免**同步重试**导致的服务压力突增。

#### 架构图

```
无抖动的重试风暴:

时间: 0ms   100ms   200ms   400ms   800ms
      │      │       │       │       │
      ▼      ▼       ▼       ▼       ▼
     失败   重试    重试    重试    重试

多个客户端同时重试，形成流量峰值

添加抖动后:

客户端A: 0ms   100ms   220ms   450ms   890ms
客户端B: 0ms   110ms   250ms   380ms   820ms
客户端C: 0ms   95ms    180ms   420ms   850ms

重试时间分散，避免流量峰值

抖动算法:
Full Jitter:  sleep = random(0, base)
Equal Jitter: sleep = base/2 + random(0, base/2)
Decorrelated Jitter: sleep = random(base, sleep*3)
```

#### 完整示例：Go实现抖动

```go
// retry/jitter.go
package retry

import (
    "math/rand"
    "sync"
    "time"
)

// JitterType 抖动类型
type JitterType int

const (
    NoJitter JitterType = iota
    FullJitter
    EqualJitter
    DecorrelatedJitter
)

// Jitter 抖动接口
type Jitter interface {
    Apply(base time.Duration) time.Duration
}

// FullJitterImpl 完全抖动 [0, base)
type FullJitterImpl struct {
    rnd *rand.Rand
    mu  sync.Mutex
}

func NewFullJitter() *FullJitterImpl {
    return &FullJitterImpl{
        rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (f *FullJitterImpl) Apply(base time.Duration) time.Duration {
    if base <= 0 {
        return 0
    }

    f.mu.Lock()
    n := f.rd.Int63n(int64(base))
    f.mu.Unlock()

    return time.Duration(n)
}

// EqualJitterImpl 等分抖动 [base/2, base)
type EqualJitterImpl struct {
    rnd *rand.Rand
    mu  sync.Mutex
}

func NewEqualJitter() *EqualJitterImpl {
    return &EqualJitterImpl{
        rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (e *EqualJitterImpl) Apply(base time.Duration) time.Duration {
    if base <= 0 {
        return 0
    }

    half := int64(base) / 2

    e.mu.Lock()
    n := half + e.rnd.Int63n(half)
    e.mu.Unlock()

    return time.Duration(n)
}

// DecorrelatedJitterImpl 去相关抖动
type DecorrelatedJitterImpl struct {
    rnd   *rand.Rand
    sleep time.Duration
    mu    sync.Mutex
}

func NewDecorrelatedJitter(initial time.Duration) *DecorrelatedJitterImpl {
    return &DecorrelatedJitterImpl{
        rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
        sleep: initial,
    }
}

func (d *DecorrelatedJitterImpl) Apply(base time.Duration) time.Duration {
    d.mu.Lock()
    defer d.mu.Unlock()

    // sleep = random(base, sleep*3)
    min := int64(base)
    max := int64(d.sleep) * 3

    if max <= min {
        max = min * 3
    }

    d.sleep = time.Duration(min + d.rnd.Int63n(max-min))
    return d.sleep
}

func (d *DecorrelatedJitterImpl) Reset() {
    d.mu.Lock()
    d.sleep = 0
    d.mu.Unlock()
}

// ExponentialBackoffWithJitter 带抖动的指数退避
type ExponentialBackoffWithJitter struct {
    *ExponentialBackoff
    jitter Jitter
}

func NewExponentialBackoffWithJitter(initial, max time.Duration, multiplier float64, jitterType JitterType) *ExponentialBackoffWithJitter {
    eb := NewExponentialBackoff(initial, max, multiplier)

    var jitter Jitter
    switch jitterType {
    case FullJitter:
        jitter = NewFullJitter()
    case EqualJitter:
        jitter = NewEqualJitter()
    case DecorrelatedJitter:
        jitter = NewDecorrelatedJitter(initial)
    default:
        jitter = nil
    }

    return &ExponentialBackoffWithJitter{
        ExponentialBackoff: eb,
        jitter:             jitter,
    }
}

func (e *ExponentialBackoffWithJitter) Next() time.Duration {
    base := e.ExponentialBackoff.Next()

    if e.jitter != nil {
        return e.jitter.Apply(base)
    }

    return base
}

// 使用AWS推荐的抖动策略
func AWSBackoff() *ExponentialBackoffWithJitter {
    return NewExponentialBackoffWithJitter(
        100*time.Millisecond,
        20*time.Second,
        2,
        FullJitter,
    )
}

// 使用Google推荐的抖动策略
func GoogleBackoff() *ExponentialBackoffWithJitter {
    return NewExponentialBackoffWithJitter(
        100*time.Millisecond,
        60*time.Second,
        2,
        EqualJitter,
    )
}
```

#### 反例说明

```go
// ❌ 错误：不使用抖动
for i := 0; i < retries; i++ {
    time.Sleep(time.Duration(math.Pow(2, float64(i))) * 100 * time.Millisecond)
    retry()
}
// 问题：多个客户端同时重试，形成重试风暴

// ❌ 错误：抖动范围太大
return time.Duration(rand.Int63n(int64(base * 10)))
// 问题：抖动过大，重试间隔不可预测

// ❌ 错误：使用全局随机数
rand.Int63n(int64(base))  // 使用全局rand
// 问题：不是并发安全的
```

---

### 6.3 重试策略

#### 概念定义

重试策略定义了何时重试、重试多少次、什么错误可重试等规则。合理的重试策略能**提高成功率**，不合理的策略会**放大故障**。

#### 完整示例：Go实现重试策略

```go
// retry/retry.go
package retry

import (
    "context"
    "errors"
    "fmt"
    "time"
)

// RetryableFunc 可重试函数类型
type RetryableFunc func() error

// Config 重试配置
type Config struct {
    MaxRetries  int           // 最大重试次数
    Backoff     Backoff       // 退避策略
    RetryIf     func(error) bool // 判断错误是否可重试
    OnRetry     func(uint, error) // 重试回调
    OnSuccess   func(uint)       // 成功回调
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
    return &Config{
        MaxRetries: 3,
        Backoff:    NewExponentialBackoff(100*time.Millisecond, 5*time.Second, 2),
        RetryIf:    IsRetryableError,
    }
}

// IsRetryableError 判断错误是否可重试
func IsRetryableError(err error) bool {
    if err == nil {
        return false
    }

    // 网络错误通常可重试
    var netErr net.Error
    if errors.As(err, &netErr) {
        return netErr.Temporary() || netErr.Timeout()
    }

    // 特定HTTP状态码可重试
    var httpErr *HTTPError
    if errors.As(err, &httpErr) {
        switch httpErr.StatusCode {
        case 408, 429, 500, 502, 503, 504:
            return true
        }
    }

    // 上下文取消不可重试
    if errors.Is(err, context.Canceled) {
        return false
    }

    // 超时错误可重试
    if errors.Is(err, context.DeadlineExceeded) {
        return true
    }

    return false
}

// HTTPError HTTP错误
type HTTPError struct {
    StatusCode int
    Message    string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// Do 执行重试
func Do(fn RetryableFunc, config *Config) error {
    if config == nil {
        config = DefaultConfig()
    }

    config.Backoff.Reset()

    var lastErr error

    for attempt := uint(0); attempt <= uint(config.MaxRetries); attempt++ {
        if attempt > 0 {
            // 等待退避时间
            time.Sleep(config.Backoff.Next())
        }

        err := fn()
        if err == nil {
            if config.OnSuccess != nil {
                config.OnSuccess(attempt)
            }
            return nil
        }

        lastErr = err

        // 检查是否可重试
        if config.RetryIf != nil && !config.RetryIf(err) {
            return err
        }

        // 达到最大重试次数
        if attempt >= uint(config.MaxRetries) {
            break
        }

        // 重试回调
        if config.OnRetry != nil {
            config.OnRetry(attempt+1, err)
        }
    }

    return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// DoWithContext 带上下文的执行重试
func DoWithContext(ctx context.Context, fn RetryableFunc, config *Config) error {
    if config == nil {
        config = DefaultConfig()
    }

    config.Backoff.Reset()

    var lastErr error

    for attempt := uint(0); attempt <= uint(config.MaxRetries); attempt++ {
        // 检查上下文是否取消
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        if attempt > 0 {
            // 带上下文的等待
            timer := time.NewTimer(config.Backoff.Next())
            select {
            case <-timer.C:
            case <-ctx.Done():
                timer.Stop()
                return ctx.Err()
            }
        }

        err := fn()
        if err == nil {
            if config.OnSuccess != nil {
                config.OnSuccess(attempt)
            }
            return nil
        }

        lastErr = err

        if config.RetryIf != nil && !config.RetryIf(err) {
            return err
        }

        if attempt >= uint(config.MaxRetries) {
            break
        }

        if config.OnRetry != nil {
            config.OnRetry(attempt+1, err)
        }
    }

    return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// HTTP客户端集成重试
type RetryableHTTPClient struct {
    client *http.Client
    config *Config
}

func NewRetryableHTTPClient(config *Config) *RetryableHTTPClient {
    return &RetryableHTTPClient{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        config: config,
    }
}

func (c *RetryableHTTPClient) Do(req *http.Request) (*http.Response, error) {
    var resp *http.Response
    var err error

    retryErr := DoWithContext(req.Context(), func() error {
        resp, err = c.client.Do(req)
        if err != nil {
            return err
        }

        // 5xx错误视为可重试
        if resp.StatusCode >= 500 {
            resp.Body.Close()
            return &HTTPError{StatusCode: resp.StatusCode}
        }

        return nil
    }, c.config)

    if retryErr != nil {
        return nil, retryErr
    }

    return resp, err
}

// gRPC客户端集成重试
func GRPCRetryInterceptor(config *Config) grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        var lastErr error

        config.Backoff.Reset()

        for attempt := uint(0); attempt <= uint(config.MaxRetries); attempt++ {
            if attempt > 0 {
                timer := time.NewTimer(config.Backoff.Next())
                select {
                case <-timer.C:
                case <-ctx.Done():
                    timer.Stop()
                    return ctx.Err()
                }
            }

            err := invoker(ctx, method, req, reply, cc, opts...)
            if err == nil {
                return nil
            }

            lastErr = err

            // gRPC状态码判断
            if st, ok := status.FromError(err); ok {
                switch st.Code() {
                case codes.DeadlineExceeded, codes.Unavailable, codes.ResourceExhausted:
                    // 可重试
                default:
                    return err
                }
            }

            if attempt >= uint(config.MaxRetries) {
                break
            }
        }

        return lastErr
    }
}
```

#### 反例说明

```go
// ❌ 错误：无限重试
for {
    err := doSomething()
    if err != nil {
        time.Sleep(time.Second)
        continue  // 无限循环
    }
    break
}

// ❌ 错误：重试所有错误
RetryIf: func(err error) bool {
    return true  // 所有错误都重试
}
// 问题：4xx错误重试无意义

// ❌ 错误：不检查上下文
for i := 0; i < retries; i++ {
    time.Sleep(backoff.Next())
    doSomething()  // 上下文已取消但仍继续
}
```

---

### 6.4 幂等性设计

#### 概念定义

幂等性是指**多次执行同一操作，结果相同**。在分布式系统中，由于重试、超时等原因，请求可能被处理多次，幂等性保证数据一致性。

#### 架构图

```
幂等性保证:

客户端            服务端
  │                │
  │── 请求A ──────>│
  │   (超时)       │ 处理请求A
  │                │ 生成ID: req-001
  │                │ 执行操作
  │                │ 返回结果
  │                │
  │ (未收到响应)   │
  │                │
  │── 重试A ──────>│
  │                │ 检查ID: req-001
  │                │ 已存在，返回缓存结果
  │<─ 响应 ────────│
  │                │

幂等性实现方式:
1. 唯一请求ID: 客户端生成，服务端去重
2. 乐观锁: 版本号控制
3. 数据库唯一约束: 业务字段唯一
4. Token机制: 预分配Token，执行时校验
```

#### 完整示例：Go实现幂等性

```go
// idempotency/idempotency.go
package idempotency

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// IdempotencyKey 幂等性Key生成器
type IdempotencyKey struct {
    Prefix    string
    Separator string
}

func NewIdempotencyKey(prefix string) *IdempotencyKey {
    return &IdempotencyKey{
        Prefix:    prefix,
        Separator: ":",
    }
}

func (k *IdempotencyKey) Generate(userID, action string) string {
    // 生成唯一ID
    b := make([]byte, 8)
    rand.Read(b)
    uniqueID := hex.EncodeToString(b)

    return fmt.Sprintf("%s%s%s%s%s%s%d",
        k.Prefix, k.Separator,
        userID, k.Separator,
        action, k.Separator,
        time.Now().UnixNano(),
    ) + uniqueID
}

// IdempotencyStore 幂等性存储接口
type IdempotencyStore interface {
    Get(ctx context.Context, key string) (*IdempotencyRecord, error)
    Set(ctx context.Context, key string, record *IdempotencyRecord, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

// IdempotencyRecord 幂等性记录
type IdempotencyRecord struct {
    Status   string          `json:"status"`    // pending, completed, failed
    Response json.RawMessage `json:"response"`
    CreatedAt time.Time      `json:"created_at"`
}

// RedisIdempotencyStore Redis实现
type RedisIdempotencyStore struct {
    client *redis.Client
}

func NewRedisIdempotencyStore(client *redis.Client) *RedisIdempotencyStore {
    return &RedisIdempotencyStore{client: client}
}

func (r *RedisIdempotencyStore) Get(ctx context.Context, key string) (*IdempotencyRecord, error) {
    data, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    var record IdempotencyRecord
    if err := json.Unmarshal([]byte(data), &record); err != nil {
        return nil, err
    }

    return &record, nil
}

func (r *RedisIdempotencyStore) Set(ctx context.Context, key string, record *IdempotencyRecord, ttl time.Duration) error {
    data, err := json.Marshal(record)
    if err != nil {
        return err
    }

    return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisIdempotencyStore) Delete(ctx context.Context, key string) error {
    return r.client.Del(ctx, key).Err()
}

// IdempotencyGuard 幂等性守卫
type IdempotencyGuard struct {
    store IdempotencyStore
}

func NewIdempotencyGuard(store IdempotencyStore) *IdempotencyGuard {
    return &IdempotencyGuard{store: store}
}

// Execute 执行幂等操作
func (g *IdempotencyGuard) Execute(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error)) (interface{}, error) {
    // 检查是否已存在
    record, err := g.store.Get(ctx, key)
    if err != nil {
        return nil, err
    }

    if record != nil {
        switch record.Status {
        case "completed":
            // 已完成，返回缓存结果
            var response interface{}
            if err := json.Unmarshal(record.Response, &response); err != nil {
                return nil, err
            }
            return response, nil
        case "pending":
            // 正在处理，返回冲突
            return nil, errors.New("request is being processed")
        case "failed":
            // 上次失败，可以重试
        }
    }

    // 设置pending状态
    pendingRecord := &IdempotencyRecord{
        Status:    "pending",
        CreatedAt: time.Now(),
    }
    if err := g.store.Set(ctx, key, pendingRecord, ttl); err != nil {
        return nil, err
    }

    // 执行业务逻辑
    result, err := fn()

    // 更新状态
    if err != nil {
        failedRecord := &IdempotencyRecord{
            Status:    "failed",
            CreatedAt: time.Now(),
        }
        g.store.Set(ctx, key, failedRecord, ttl)
        return nil, err
    }

    // 序列化结果
    responseData, _ := json.Marshal(result)
    completedRecord := &IdempotencyRecord{
        Status:    "completed",
        Response:  responseData,
        CreatedAt: time.Now(),
    }
    g.store.Set(ctx, key, completedRecord, ttl)

    return result, nil
}

// HTTP中间件集成幂等性
func IdempotencyMiddleware(store IdempotencyStore, headerName string) func(http.Handler) http.Handler {
    guard := NewIdempotencyGuard(store)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 只处理非幂等方法
            if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
                next.ServeHTTP(w, r)
                return
            }

            idempotencyKey := r.Header.Get(headerName)
            if idempotencyKey == "" {
                http.Error(w, "Missing idempotency key", http.StatusBadRequest)
                return
            }

            // 包装ResponseWriter捕获响应
            rw := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

            // 使用幂等性守卫执行
            _, err := guard.Execute(r.Context(), idempotencyKey, 24*time.Hour, func() (interface{}, error) {
                next.ServeHTTP(rw, r)
                if rw.statusCode >= 500 {
                    return nil, fmt.Errorf("server error: %d", rw.statusCode)
                }
                return rw.body, nil
            })

            if err != nil {
                if err.Error() == "request is being processed" {
                    http.Error(w, err.Error(), http.StatusConflict)
                    return
                }
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        })
    }
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
    body       []byte
}

func (r *responseRecorder) WriteHeader(code int) {
    r.statusCode = code
    r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(p []byte) (n int, err error) {
    r.body = append(r.body, p...)
    return r.ResponseWriter.Write(p)
}

// 数据库乐观锁实现幂等性
type Order struct {
    ID        string
    UserID    string
    Amount    float64
    Status    string
    Version   int  // 乐观锁版本号
    CreatedAt time.Time
}

func UpdateOrderStatus(ctx context.Context, db *sql.DB, orderID, status string, version int) error {
    result, err := db.ExecContext(ctx,
        "UPDATE orders SET status = ?, version = version + 1 WHERE id = ? AND version = ?",
        status, orderID, version,
    )
    if err != nil {
        return err
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rows == 0 {
        return errors.New("concurrent update detected")
    }

    return nil
}

// Token机制实现幂等性
type TokenService struct {
    store IdempotencyStore
}

func (s *TokenService) GenerateToken(ctx context.Context, userID string) (string, error) {
    token := generateUniqueToken()

    // 预分配Token
    record := &IdempotencyRecord{
        Status:    "pending",
        CreatedAt: time.Now(),
    }

    if err := s.store.Set(ctx, token, record, 5*time.Minute); err != nil {
        return "", err
    }

    return token, nil
}

func (s *TokenService) ConsumeToken(ctx context.Context, token string, fn func() error) error {
    record, err := s.store.Get(ctx, token)
    if err != nil {
        return err
    }

    if record == nil {
        return errors.New("invalid or expired token")
    }

    if record.Status != "pending" {
        return errors.New("token already consumed")
    }

    // 执行操作
    if err := fn(); err != nil {
        return err
    }

    // 标记Token已使用
    record.Status = "completed"
    return s.store.Set(ctx, token, record, 24*time.Hour)
}

func generateUniqueToken() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}
```

#### 反例说明

```go
// ❌ 错误：幂等Key不唯一
key := fmt.Sprintf("order:%s", userID)
// 问题：同一用户的不同订单使用相同key

// ❌ 错误：幂等记录不设置过期时间
store.Set(ctx, key, record, 0)  // 永不过期
// 问题：存储无限增长

// ❌ 错误：不处理pending状态
if record.Status == "pending" {
    return record.Response, nil  // 直接返回，可能还在处理中
}

// ❌ 错误：幂等性检查在事务外
record, _ := store.Get(ctx, key)
if record != nil {
    return record.Response, nil
}
// 此时另一个请求可能也在执行
// 问题：竞态条件
```

#### 注意事项

1. **Key设计**：Key应包含用户ID、操作类型、业务唯一标识
2. **过期时间**：设置合理的过期时间，避免存储无限增长
3. **pending处理**：pending状态需要处理，避免死锁
4. **原子性**：幂等性检查和业务操作应在同一事务中

---

## 7. 分布式事务

### 7.1 两阶段提交（2PC）

#### 概念定义

两阶段提交(Two-Phase Commit)是一种经典的分布式事务协议，通过**准备阶段**和**提交阶段**保证所有参与者要么全部提交，要么全部回滚。

#### 架构图

```
2PC执行流程:

阶段一: 准备阶段

┌─────────┐         ┌─────────┐         ┌─────────┐
│ 协调者   │         │ 参与者A  │         │ 参与者B  │
│(Coordinator)│      │         │         │         │
└────┬────┘         └────┬────┘         └────┬────┘
     │                   │                   │
     │── 1.CanCommit? ──>│                   │
     │                   │── 2.执行本地事务 ──>│
     │                   │   记录undo/redo   │
     │<─ 3.Yes/No ───────│                   │
     │                   │                   │
     │── 1.CanCommit? ──────────────────────>│
     │                   │                   │── 2.执行本地事务
     │                   │                   │   记录undo/redo
     │<─ 3.Yes/No ───────────────────────────│
     │                   │                   │

阶段二: 提交阶段(所有参与者返回Yes)

     │── 4.DoCommit ────>│                   │
     │                   │── 5.提交本地事务 ──>│
     │                   │                   │
     │<─ 6.Ack ──────────│                   │
     │                   │                   │
     │── 4.DoCommit ────────────────────────>│
     │                   │                   │── 5.提交本地事务
     │<─ 6.Ack ──────────────────────────────│
     │                   │                   │
     │ 完成               │                   │

回滚场景(任一参与者返回No):

     │── 4.DoAbort ─────>│                   │
     │                   │── 5.回滚本地事务 ──>│
     │<─ 6.Ack ──────────│                   │
```

#### 完整示例：Go实现2PC

```go
// twopc/twopc.go
package twopc

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"
)

// Vote 投票结果
type Vote int

const (
    VoteYes Vote = iota
    VoteNo
)

// Participant 参与者接口
type Participant interface {
    ID() string
    Prepare(ctx context.Context, txID string, data interface{}) (Vote, error)
    Commit(ctx context.Context, txID string) error
    Rollback(ctx context.Context, txID string) error
}

// Coordinator 协调者
type Coordinator struct {
    participants []Participant
    timeout      time.Duration
    log          TransactionLog
}

// TransactionLog 事务日志接口
type TransactionLog interface {
    Write(ctx context.Context, record *LogRecord) error
    Read(ctx context.Context, txID string) (*LogRecord, error)
}

// LogRecord 日志记录
type LogRecord struct {
    TxID        string
    Status      string // prepare, commit, abort
    Participants []string
    Data        interface{}
    CreatedAt   time.Time
}

// NewCoordinator 创建协调者
func NewCoordinator(log TransactionLog, timeout time.Duration, participants ...Participant) *Coordinator {
    return &Coordinator{
        participants: participants,
        timeout:      timeout,
        log:          log,
    }
}

// Execute 执行2PC事务
func (c *Coordinator) Execute(ctx context.Context, txID string, data interface{}) error {
    // 阶段一: 准备
    votes, err := c.preparePhase(ctx, txID, data)
    if err != nil {
        return fmt.Errorf("prepare phase failed: %w", err)
    }

    // 检查投票结果
    allYes := true
    for _, vote := range votes {
        if vote == VoteNo {
            allYes = false
            break
        }
    }

    // 阶段二: 提交或回滚
    if allYes {
        return c.commitPhase(ctx, txID)
    }
    return c.rollbackPhase(ctx, txID)
}

// preparePhase 准备阶段
func (c *Coordinator) preparePhase(ctx context.Context, txID string, data interface{}) (map[string]Vote, error) {
    votes := make(map[string]Vote)
    var mu sync.Mutex
    var wg sync.WaitGroup

    // 记录prepare日志
    participantIDs := make([]string, len(c.participants))
    for i, p := range c.participants {
        participantIDs[i] = p.ID()
    }

    c.log.Write(ctx, &LogRecord{
        TxID:         txID,
        Status:       "prepare",
        Participants: participantIDs,
        Data:         data,
        CreatedAt:    time.Now(),
    })

    // 并行发送prepare请求
    for _, p := range c.participants {
        wg.Add(1)
        go func(participant Participant) {
            defer wg.Done()

            vote, err := c.callPrepare(ctx, participant, txID, data)

            mu.Lock()
            if err != nil {
                votes[participant.ID()] = VoteNo
            } else {
                votes[participant.ID()] = vote
            }
            mu.Unlock()
        }(p)
    }

    wg.Wait()
    return votes, nil
}

func (c *Coordinator) callPrepare(ctx context.Context, p Participant, txID string, data interface{}) (Vote, error) {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    return p.Prepare(ctx, txID, data)
}

// commitPhase 提交阶段
func (c *Coordinator) commitPhase(ctx context.Context, txID string) error {
    // 记录commit日志
    c.log.Write(ctx, &LogRecord{
        TxID:      txID,
        Status:    "commit",
        CreatedAt: time.Now(),
    })

    var wg sync.WaitGroup
    errChan := make(chan error, len(c.participants))

    for _, p := range c.participants {
        wg.Add(1)
        go func(participant Participant) {
            defer wg.Done()

            if err := c.callCommit(ctx, participant, txID); err != nil {
                errChan <- fmt.Errorf("participant %s commit failed: %w", participant.ID(), err)
            }
        }(p)
    }

    wg.Wait()
    close(errChan)

    // 检查是否有错误
    var errs []error
    for err := range errChan {
        errs = append(errs, err)
    }

    if len(errs) > 0 {
        return fmt.Errorf("commit phase failed: %v", errs)
    }

    return nil
}

func (c *Coordinator) callCommit(ctx context.Context, p Participant, txID string) error {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    return p.Commit(ctx, txID)
}

// rollbackPhase 回滚阶段
func (c *Coordinator) rollbackPhase(ctx context.Context, txID string) error {
    // 记录abort日志
    c.log.Write(ctx, &LogRecord{
        TxID:      txID,
        Status:    "abort",
        CreatedAt: time.Now(),
    })

    var wg sync.WaitGroup

    for _, p := range c.participants {
        wg.Add(1)
        go func(participant Participant) {
            defer wg.Done()

            // 回滚失败需要记录，后续人工处理
            if err := c.callRollback(ctx, participant, txID); err != nil {
                fmt.Printf("participant %s rollback failed: %v\n", participant.ID(), err)
            }
        }(p)
    }

    wg.Wait()
    return errors.New("transaction aborted")
}

func (c *Coordinator) callRollback(ctx context.Context, p Participant, txID string) error {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    return p.Rollback(ctx, txID)
}

// 参与者实现示例
type DBParticipant struct {
    id string
    db *sql.DB
}

func (d *DBParticipant) ID() string {
    return d.id
}

func (d *DBParticipant) Prepare(ctx context.Context, txID string, data interface{}) (Vote, error) {
    // 开始本地事务
    tx, err := d.db.BeginTx(ctx, nil)
    if err != nil {
        return VoteNo, err
    }

    // 执行本地操作
    // 记录undo/redo日志

    // 不提交，等待协调者指令
    // 将tx保存在内存中，key为txID

    return VoteYes, nil
}

func (d *DBParticipant) Commit(ctx context.Context, txID string) error {
    // 获取之前保存的事务
    // tx := getTransaction(txID)
    // return tx.Commit()
    return nil
}

func (d *DBParticipant) Rollback(ctx context.Context, txID string) error {
    // 获取之前保存的事务
    // tx := getTransaction(txID)
    // return tx.Rollback()
    return nil
}
```

#### 反例说明

```go
// ❌ 错误：协调者单点故障
// 协调者宕机后，参与者无法知道应该提交还是回滚
// 问题：需要协调者恢复或人工干预

// ❌ 错误：不记录日志
// 协调者宕机重启后不知道事务状态
// 问题：数据不一致

// ❌ 错误：同步阻塞所有参与者
// 参与者A准备完成后等待协调者指令
// 参与者B准备失败，但A已锁定资源
// 问题：资源长时间锁定
```

#### 注意事项

1. **协调者高可用**：协调者需要集群部署，避免单点故障
2. **事务日志**：必须持久化事务日志，用于故障恢复
3. **超时处理**：设置合理的超时时间，避免无限等待
4. **性能影响**：2PC会阻塞资源，影响系统吞吐量

---

### 7.2 三阶段提交（3PC）

#### 概念定义

三阶段提交(Three-Phase Commit)是2PC的改进版，增加**CanCommit预检查阶段**，减少阻塞时间，协调者宕机时可以通过超时继续推进。

#### 架构图

```
3PC执行流程:

阶段一: CanCommit

协调者 ── CanCommit? ──> 参与者
参与者检查是否可以执行（不锁定资源）
参与者 <── Yes/No ──── 协调者

阶段二: PreCommit (所有参与者返回Yes)

协调者 ── PreCommit ──> 参与者
参与者锁定资源，记录undo/redo
参与者 <── Ack ─────── 协调者

阶段三: DoCommit

协调者 ── DoCommit ───> 参与者
参与者提交事务
参与者 <── Ack ─────── 协调者

超时处理:
- 参与者等待PreCommit超时: 直接中止
- 参与者等待DoCommit超时: 询问协调者或根据日志决定
- 协调者等待Ack超时: 继续重试

vs 2PC:
- 3PC增加CanCommit阶段，减少资源锁定时间
- 3PC协调者宕机时，参与者可以超时继续
- 3PC网络分区时仍可能不一致
```

#### 完整示例：Go实现3PC

```go
// threepc/threepc.go
package threepc

import (
    "context"
    "fmt"
    "time"
)

// Phase 阶段类型
type Phase int

const (
    PhaseCanCommit Phase = iota
    PhasePreCommit
    PhaseDoCommit
)

// Participant3PC 3PC参与者接口
type Participant3PC interface {
    ID() string
    CanCommit(ctx context.Context, txID string, data interface{}) (bool, error)
    PreCommit(ctx context.Context, txID string) error
    DoCommit(ctx context.Context, txID string) error
    Abort(ctx context.Context, txID string) error
    GetStatus(ctx context.Context, txID string) (string, error)
}

// Coordinator3PC 3PC协调者
type Coordinator3PC struct {
    participants []Participant3PC
    timeout      time.Duration
    log          TransactionLog
}

// Execute 执行3PC事务
func (c *Coordinator3PC) Execute(ctx context.Context, txID string, data interface{}) error {
    // 阶段一: CanCommit
    canCommits, err := c.canCommitPhase(ctx, txID, data)
    if err != nil {
        return fmt.Errorf("cancommit phase failed: %w", err)
    }

    // 检查是否所有参与者都可以执行
    allYes := true
    for _, can := range canCommits {
        if !can {
            allYes = false
            break
        }
    }

    if !allYes {
        c.abortPhase(ctx, txID)
        return fmt.Errorf("transaction aborted in cancommit phase")
    }

    // 阶段二: PreCommit
    if err := c.preCommitPhase(ctx, txID); err != nil {
        c.abortPhase(ctx, txID)
        return fmt.Errorf("precommit phase failed: %w", err)
    }

    // 阶段三: DoCommit
    return c.doCommitPhase(ctx, txID)
}

func (c *Coordinator3PC) canCommitPhase(ctx context.Context, txID string, data interface{}) (map[string]bool, error) {
    results := make(map[string]bool)

    for _, p := range c.participants {
        can, err := c.callCanCommit(ctx, p, txID, data)
        if err != nil {
            results[p.ID()] = false
        } else {
            results[p.ID()] = can
        }
    }

    return results, nil
}

func (c *Coordinator3PC) callCanCommit(ctx context.Context, p Participant3PC, txID string, data interface{}) (bool, error) {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    return p.CanCommit(ctx, txID, data)
}

func (c *Coordinator3PC) preCommitPhase(ctx context.Context, txID string) error {
    // 记录precommit日志
    c.log.Write(ctx, &LogRecord{
        TxID:      txID,
        Status:    "precommit",
        CreatedAt: time.Now(),
    })

    for _, p := range c.participants {
        if err := c.callPreCommit(ctx, p, txID); err != nil {
            return err
        }
    }

    return nil
}

func (c *Coordinator3PC) callPreCommit(ctx context.Context, p Participant3PC, txID string) error {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    return p.PreCommit(ctx, txID)
}

func (c *Coordinator3PC) doCommitPhase(ctx context.Context, txID string) error {
    // 记录docommit日志
    c.log.Write(ctx, &LogRecord{
        TxID:      txID,
        Status:    "docommit",
        CreatedAt: time.Now(),
    })

    for _, p := range c.participants {
        if err := c.callDoCommit(ctx, p, txID); err != nil {
            // 记录需要后续处理
            fmt.Printf("participant %s docommit failed: %v\n", p.ID(), err)
        }
    }

    return nil
}

func (c *Coordinator3PC) callDoCommit(ctx context.Context, p Participant3PC, txID string) error {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    return p.DoCommit(ctx, txID)
}

func (c *Coordinator3PC) abortPhase(ctx context.Context, txID string) {
    c.log.Write(ctx, &LogRecord{
        TxID:      txID,
        Status:    "abort",
        CreatedAt: time.Now(),
    })

    for _, p := range c.participants {
        ctx, cancel := context.WithTimeout(ctx, c.timeout)
        p.Abort(ctx, txID)
        cancel()
    }
}

// 参与者超时处理
type Participant3PCImpl struct {
    id      string
    timeout time.Duration
}

func (p *Participant3PCImpl) PreCommit(ctx context.Context, txID string) error {
    // 锁定资源
    // 记录状态为"prepared"

    // 启动超时定时器
    go func() {
        timer := time.NewTimer(p.timeout)
        defer timer.Stop()

        <-timer.C

        // 超时后询问协调者或根据日志决定
        status, _ := p.GetStatus(context.Background(), txID)
        if status == "prepared" {
            // 协调者可能宕机，询问其他参与者或等待
        }
    }()

    return nil
}
```

#### 反例说明

```go
// ❌ 错误：网络分区时仍可能不一致
// 协调者和部分参与者在一个分区，其他参与者在另一个分区
// 问题：部分提交部分回滚

// ❌ 错误：超时时间设置不合理
// 超时太短，正常事务被误判为失败
// 超时太长，阻塞时间增加
```

---

### 7.3 TCC（Try-Confirm-Cancel）

#### 概念定义

TCC是一种业务层面的分布式事务方案，将每个操作拆分为**Try（预留资源）**、**Confirm（确认执行）**、**Cancel（取消释放）**三个阶段。

#### 架构图

```
TCC执行流程:

订单服务                    库存服务                    账户服务
   │                         │                         │
   │── Try创建订单(待确认) ──>│                         │
   │   插入订单，状态=PENDING  │                         │
   │<─ 成功 ─────────────────│                         │
   │                         │                         │
   │── Try扣减库存 ──────────>│                         │
   │                         │  冻结库存               │
   │<─ 成功 ─────────────────│                         │
   │                         │                         │
   │── Try扣减余额 ────────────────────────────────────>│
   │                         │                         │  冻结金额
   │<─ 成功 ────────────────────────────────────────────│
   │                         │                         │
   │ 所有Try成功             │                         │
   │                         │                         │
   │── Confirm订单 ─────────>│                         │
   │   状态=CONFIRMED        │                         │
   │<─ 成功 ─────────────────│                         │
   │                         │                         │
   │── Confirm库存 ─────────>│                         │
   │                         │  确认扣减               │
   │<─ 成功 ─────────────────│                         │
   │                         │                         │
   │── Confirm账户 ────────────────────────────────────>│
   │                         │                         │  确认扣减
   │<─ 成功 ────────────────────────────────────────────│
   │                         │                         │
   │ 事务完成                │                         │

Cancel场景(Try失败):

   │── Try创建订单 ─────────>│                         │
   │<─ 成功 ─────────────────│                         │
   │                         │                         │
   │── Try扣减库存 ─────────>│                         │
   │<─ 失败 ─────────────────│                         │
   │                         │                         │
   │ 执行Cancel              │                         │
   │                         │                         │
   │── Cancel订单 ─────────>│                         │
   │   状态=CANCELLED        │                         │
   │<─ 成功 ─────────────────│                         │
```

#### 完整示例：Go实现TCC

```go
// tcc/tcc.go
package tcc

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// TCCAction TCC动作接口
type TCCAction interface {
    Try(ctx context.Context, bizID string, params interface{}) error
    Confirm(ctx context.Context, bizID string) error
    Cancel(ctx context.Context, bizID string) error
}

// TCCBranch TCC分支
type TCCBranch struct {
    Name   string
    Action TCCAction
}

// TCCTransaction TCC事务
type TCCTransaction struct {
    branches []TCCBranch
    timeout  time.Duration
}

// NewTCCTransaction 创建TCC事务
func NewTCCTransaction(timeout time.Duration, branches ...TCCBranch) *TCCTransaction {
    return &TCCTransaction{
        branches: branches,
        timeout:  timeout,
    }
}

// Execute 执行TCC事务
func (t *TCCTransaction) Execute(ctx context.Context, bizID string, params interface{}) error {
    // 阶段一: Try
    tryResults := make(map[string]error)

    for _, branch := range t.branches {
        ctx, cancel := context.WithTimeout(ctx, t.timeout)
        err := branch.Action.Try(ctx, bizID, params)
        cancel()

        tryResults[branch.Name] = err

        if err != nil {
            // Try失败，执行Cancel
            t.cancelAll(ctx, bizID)
            return fmt.Errorf("try failed for %s: %w", branch.Name, err)
        }
    }

    // 阶段二: Confirm
    var confirmErrs []error

    for _, branch := range t.branches {
        ctx, cancel := context.WithTimeout(ctx, t.timeout)
        err := branch.Action.Confirm(ctx, bizID)
        cancel()

        if err != nil {
            confirmErrs = append(confirmErrs, fmt.Errorf("confirm failed for %s: %w", branch.Name, err))
        }
    }

    if len(confirmErrs) > 0 {
        // Confirm失败需要人工处理或重试
        return fmt.Errorf("confirm errors: %v", confirmErrs)
    }

    return nil
}

func (t *TCCTransaction) cancelAll(ctx context.Context, bizID string) {
    var wg sync.WaitGroup

    for _, branch := range t.branches {
        wg.Add(1)
        go func(b TCCBranch) {
            defer wg.Done()

            ctx, cancel := context.WithTimeout(ctx, t.timeout)
            defer cancel()

            if err := b.Action.Cancel(ctx, bizID); err != nil {
                fmt.Printf("cancel failed for %s: %v\n", b.Name, err)
            }
        }(branch)
    }

    wg.Wait()
}

// 库存服务TCC实现
type InventoryTCC struct {
    db *sql.DB
}

func (i *InventoryTCC) Try(ctx context.Context, orderID string, params interface{}) error {
    p := params.(OrderParams)

    tx, err := i.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. 检查库存
    var available int
    err = tx.QueryRowContext(ctx,
        "SELECT stock FROM inventory WHERE product_id = ? FOR UPDATE",
        p.ProductID,
    ).Scan(&available)
    if err != nil {
        return err
    }

    if available < p.Quantity {
        return fmt.Errorf("insufficient stock")
    }

    // 2. 冻结库存
    _, err = tx.ExecContext(ctx,
        "UPDATE inventory SET stock = stock - ?, frozen = frozen + ? WHERE product_id = ?",
        p.Quantity, p.Quantity, p.ProductID,
    )
    if err != nil {
        return err
    }

    // 3. 记录冻结记录
    _, err = tx.ExecContext(ctx,
        "INSERT INTO inventory_freeze (order_id, product_id, quantity, status) VALUES (?, ?, ?, ?)",
        orderID, p.ProductID, p.Quantity, "PENDING",
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (i *InventoryTCC) Confirm(ctx context.Context, orderID string) error {
    tx, err := i.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. 更新冻结记录状态
    _, err = tx.ExecContext(ctx,
        "UPDATE inventory_freeze SET status = ? WHERE order_id = ?",
        "CONFIRMED", orderID,
    )
    if err != nil {
        return err
    }

    // 2. 扣减冻结库存
    _, err = tx.ExecContext(ctx,
        "UPDATE inventory SET frozen = frozen - (SELECT quantity FROM inventory_freeze WHERE order_id = ?) WHERE product_id = (SELECT product_id FROM inventory_freeze WHERE order_id = ?)",
        orderID, orderID,
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (i *InventoryTCC) Cancel(ctx context.Context, orderID string) error {
    tx, err := i.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. 获取冻结记录
    var productID string
    var quantity int
    err = tx.QueryRowContext(ctx,
        "SELECT product_id, quantity FROM inventory_freeze WHERE order_id = ? AND status = ?",
        orderID, "PENDING",
    ).Scan(&productID, &quantity)
    if err == sql.ErrNoRows {
        return nil // 已处理或不存在
    }
    if err != nil {
        return err
    }

    // 2. 释放冻结库存
    _, err = tx.ExecContext(ctx,
        "UPDATE inventory SET stock = stock + ?, frozen = frozen - ? WHERE product_id = ?",
        quantity, quantity, productID,
    )
    if err != nil {
        return err
    }

    // 3. 更新冻结记录状态
    _, err = tx.ExecContext(ctx,
        "UPDATE inventory_freeze SET status = ? WHERE order_id = ?",
        "CANCELLED", orderID,
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}

// 账户服务TCC实现
type AccountTCC struct {
    db *sql.DB
}

func (a *AccountTCC) Try(ctx context.Context, orderID string, params interface{}) error {
    p := params.(OrderParams)

    tx, err := a.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. 检查余额
    var balance float64
    err = tx.QueryRowContext(ctx,
        "SELECT balance FROM account WHERE user_id = ? FOR UPDATE",
        p.UserID,
    ).Scan(&balance)
    if err != nil {
        return err
    }

    if balance < p.Amount {
        return fmt.Errorf("insufficient balance")
    }

    // 2. 冻结金额
    _, err = tx.ExecContext(ctx,
        "UPDATE account SET balance = balance - ?, frozen = frozen + ? WHERE user_id = ?",
        p.Amount, p.Amount, p.UserID,
    )
    if err != nil {
        return err
    }

    // 3. 记录冻结记录
    _, err = tx.ExecContext(ctx,
        "INSERT INTO account_freeze (order_id, user_id, amount, status) VALUES (?, ?, ?, ?)",
        orderID, p.UserID, p.Amount, "PENDING",
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (a *AccountTCC) Confirm(ctx context.Context, orderID string) error {
    // 确认扣减，释放冻结金额
    // 类似库存服务的Confirm实现
    return nil
}

func (a *AccountTCC) Cancel(ctx context.Context, orderID string) error {
    // 释放冻结金额，恢复余额
    // 类似库存服务的Cancel实现
    return nil
}

// 使用示例
type OrderParams struct {
    UserID    string
    ProductID string
    Quantity  int
    Amount    float64
}

func CreateOrder(ctx context.Context, params OrderParams) error {
    orderID := generateOrderID()

    tcc := NewTCCTransaction(5*time.Second,
        TCCBranch{Name: "inventory", Action: &InventoryTCC{db: inventoryDB}},
        TCCBranch{Name: "account", Action: &AccountTCC{db: accountDB}},
    )

    return tcc.Execute(ctx, orderID, params)
}
```

#### 反例说明

```go
// ❌ 错误：Try不检查资源
func (i *InventoryTCC) Try(ctx context.Context, orderID string, params interface{}) error {
    // 直接冻结，不检查库存
    _, err := tx.Exec("UPDATE inventory SET frozen = frozen + ?", quantity)
    return err
}
// 问题：可能导致库存为负

// ❌ 错误：Confirm不幂等
func (i *InventoryTCC) Confirm(ctx context.Context, orderID string) error {
    // 直接扣减，不检查状态
    _, err := tx.Exec("UPDATE inventory SET stock = stock - ?", quantity)
    return err
}
// 问题：重复执行会导致超扣

// ❌ 错误：Cancel不处理空回滚
if err == sql.ErrNoRows {
    return err  // 应该返回nil
}
// 问题：Try失败时Cancel找不到记录会报错
```

#### 注意事项

1. **幂等性**：Confirm和Cancel必须幂等
2. **空回滚**：Cancel要处理Try未执行的情况
3. **悬挂问题**：Try超时后执行，但事务已结束
4. **幂等控制**：使用唯一键或状态机防止重复执行

---

### 7.4 Saga模式

#### 概念定义

Saga将长事务拆分为多个本地事务，每个本地事务提交后立即释放资源，通过**补偿事务**回滚已完成的操作。

#### 架构图

```
Saga执行流程:

正向流程:
T1 ──> T2 ──> T3 ──> T4
│       │       │       │
▼       ▼       ▼       ▼
成功    成功    成功    成功
完成

补偿流程(T3失败):
T1 ──> T2 ──> T3 (失败)
│       │       │
▼       ▼       ▼
成功    成功    失败
│       │
▼       ▼
C2      C1
(补偿T2) (补偿T1)

Saga编排方式:
1. 编排式(Choreography): 各服务监听事件，自主决定下一步
2. 协调式(Orchestration): 中央协调器控制流程
```

#### 完整示例：Go实现Saga

```go
// saga/saga.go
package saga

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// SagaStep Saga步骤
type SagaStep struct {
    Name       string
    Action     func(ctx context.Context) error
    Compensate func(ctx context.Context) error
}

// Saga Saga事务
type Saga struct {
    steps   []SagaStep
    history []string // 执行历史
    mu      sync.Mutex
}

// NewSaga 创建Saga
func NewSaga(steps ...SagaStep) *Saga {
    return &Saga{
        steps:   steps,
        history: make([]string, 0, len(steps)),
    }
}

// Execute 执行Saga
func (s *Saga) Execute(ctx context.Context) error {
    for i, step := range s.steps {
        // 执行正向操作
        if err := s.executeStep(ctx, step); err != nil {
            // 执行补偿
            s.compensate(ctx, i)
            return fmt.Errorf("step %s failed: %w", step.Name, err)
        }

        s.mu.Lock()
        s.history = append(s.history, step.Name)
        s.mu.Unlock()
    }

    return nil
}

func (s *Saga) executeStep(ctx context.Context, step SagaStep) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    return step.Action(ctx)
}

func (s *Saga) compensate(ctx context.Context, failedIndex int) {
    // 逆序执行补偿
    for i := failedIndex - 1; i >= 0; i-- {
        step := s.steps[i]

        if step.Compensate != nil {
            ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
            err := step.Compensate(ctx)
            cancel()

            if err != nil {
                fmt.Printf("compensation for %s failed: %v\n", step.Name, err)
                // 记录需要人工处理
            }
        }
    }
}

// 协调式Saga实现
type SagaOrchestrator struct {
    log SagaLog
}

type SagaLog interface {
    StartSaga(ctx context.Context, sagaID string, steps []string) error
    StepComplete(ctx context.Context, sagaID, step string) error
    StepFailed(ctx context.Context, sagaID, step string, err error) error
    CompensateStart(ctx context.Context, sagaID, step string) error
    CompensateComplete(ctx context.Context, sagaID, step string) error
    SagaComplete(ctx context.Context, sagaID string) error
}

func (o *SagaOrchestrator) Execute(ctx context.Context, sagaID string, steps []SagaStep) error {
    stepNames := make([]string, len(steps))
    for i, s := range steps {
        stepNames[i] = s.Name
    }

    // 记录Saga开始
    o.log.StartSaga(ctx, sagaID, stepNames)

    for i, step := range steps {
        // 执行步骤
        err := step.Action(ctx)

        if err != nil {
            o.log.StepFailed(ctx, sagaID, step.Name, err)

            // 执行补偿
            o.compensate(ctx, sagaID, steps[:i])
            return err
        }

        o.log.StepComplete(ctx, sagaID, step.Name)
    }

    o.log.SagaComplete(ctx, sagaID)
    return nil
}

func (o *SagaOrchestrator) compensate(ctx context.Context, sagaID string, steps []SagaStep) {
    for i := len(steps) - 1; i >= 0; i-- {
        step := steps[i]

        if step.Compensate != nil {
            o.log.CompensateStart(ctx, sagaID, step.Name)

            err := step.Compensate(ctx)
            if err != nil {
                fmt.Printf("compensation failed for %s: %v\n", step.Name, err)
            } else {
                o.log.CompensateComplete(ctx, sagaID, step.Name)
            }
        }
    }
}

// 订单创建Saga示例
func CreateOrderSaga(orderID string, userID string, productID string, quantity int, amount float64) *Saga {
    return NewSaga(
        SagaStep{
            Name: "create_order",
            Action: func(ctx context.Context) error {
                // 创建订单，状态=PENDING
                return createOrder(ctx, orderID, userID, amount)
            },
            Compensate: func(ctx context.Context) error {
                // 取消订单
                return cancelOrder(ctx, orderID)
            },
        },
        SagaStep{
            Name: "deduct_inventory",
            Action: func(ctx context.Context) error {
                // 扣减库存
                return deductInventory(ctx, productID, quantity)
            },
            Compensate: func(ctx context.Context) error {
                // 恢复库存
                return restoreInventory(ctx, productID, quantity)
            },
        },
        SagaStep{
            Name: "deduct_balance",
            Action: func(ctx context.Context) error {
                // 扣减余额
                return deductBalance(ctx, userID, amount)
            },
            Compensate: func(ctx context.Context) error {
                // 恢复余额
                return restoreBalance(ctx, userID, amount)
            },
        },
        SagaStep{
            Name: "confirm_order",
            Action: func(ctx context.Context) error {
                // 确认订单，状态=CONFIRMED
                return confirmOrder(ctx, orderID)
            },
        },
    )
}

// 辅助函数
func createOrder(ctx context.Context, orderID, userID string, amount float64) error {
    // 实现...
    return nil
}

func cancelOrder(ctx context.Context, orderID string) error {
    // 实现...
    return nil
}

func deductInventory(ctx context.Context, productID string, quantity int) error {
    // 实现...
    return nil
}

func restoreInventory(ctx context.Context, productID string, quantity int) error {
    // 实现...
    return nil
}

func deductBalance(ctx context.Context, userID string, amount float64) error {
    // 实现...
    return nil
}

func restoreBalance(ctx context.Context, userID string, amount float64) error {
    // 实现...
    return nil
}

func confirmOrder(ctx context.Context, orderID string) error {
    // 实现...
    return nil
}
```

#### 反例说明

```go
// ❌ 错误：补偿操作不幂等
func restoreInventory(ctx context.Context, productID string, quantity int) error {
    // 直接增加库存
    _, err := db.Exec("UPDATE inventory SET stock = stock + ?", quantity)
    return err
}
// 问题：重复补偿会导致库存超增

// ❌ 错误：补偿顺序错误
// 应该先补偿后执行的步骤
// 问题：数据不一致

// ❌ 错误：补偿失败不处理
if err := step.Compensate(ctx); err != nil {
    // 忽略错误
}
// 问题：部分补偿成功，数据不一致
```

---

### 7.5 本地消息表

#### 概念定义

本地消息表将分布式事务拆分为**本地事务**和**消息发送**，通过本地事务保证数据操作和消息记录的原子性，再通过定时任务发送消息。

#### 架构图

```
本地消息表流程:

┌─────────┐     1.开始本地事务    ┌─────────┐
│ 业务服务 │ ──────────────────> │  数据库  │
│         │                     │         │
│         │  2.执行业务操作      │         │
│         │  INSERT INTO orders │         │
│         │                     │         │
│         │  3.记录消息          │         │
│         │  INSERT INTO msg    │         │
│         │  (status=PENDING)   │         │
│         │                     │         │
│         │ <─ 4.提交事务 ──────│         │
└─────────┘                     └─────────┘
     │
     │ 5.发送消息到MQ
     ▼
┌─────────┐
│  消息队列 │
└─────────┘
     │
     ▼
┌─────────┐
│ 下游服务 │
└─────────┘

定时任务补偿:
┌─────────┐     查询PENDING消息    ┌─────────┐
│ 定时任务 │ ────────────────────> │ 消息表  │
│         │ <──────────────────── │         │
│         │     重新发送消息        │         │
│         │ ────────────────────> │  消息队列│
└─────────┘                       └─────────┘
```

#### 完整示例：Go实现本地消息表

```go
// msgtable/msg_table.go
package msgtable

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"
)

// Message 消息记录
type Message struct {
    ID        string          `json:"id"`
    Topic     string          `json:"topic"`
    Key       string          `json:"key"`
    Body      json.RawMessage `json:"body"`
    Status    string          `json:"status"` // pending, sent, failed
    RetryCount int            `json:"retry_count"`
    CreatedAt time.Time       `json:"created_at"`
    UpdatedAt time.Time       `json:"updated_at"`
}

// MessageTable 消息表管理
type MessageTable struct {
    db *sql.DB
}

func NewMessageTable(db *sql.DB) *MessageTable {
    return &MessageTable{db: db}
}

// CreateTable 创建消息表
func (m *MessageTable) CreateTable() error {
    _, err := m.db.Exec(`
        CREATE TABLE IF NOT EXISTS message_table (
            id VARCHAR(64) PRIMARY KEY,
            topic VARCHAR(128) NOT NULL,
            msg_key VARCHAR(128),
            body TEXT NOT NULL,
            status VARCHAR(20) NOT NULL DEFAULT 'pending',
            retry_count INT DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            INDEX idx_status (status),
            INDEX idx_created_at (created_at)
        )
    `)
    return err
}

// SaveMessage 保存消息（在业务事务中调用）
func (m *MessageTable) SaveMessage(ctx context.Context, tx *sql.Tx, msg *Message) error {
    msg.Status = "pending"
    msg.CreatedAt = time.Now()
    msg.UpdatedAt = time.Now()

    body, err := json.Marshal(msg.Body)
    if err != nil {
        return err
    }

    _, err = tx.ExecContext(ctx,
        "INSERT INTO message_table (id, topic, msg_key, body, status, retry_count, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        msg.ID, msg.Topic, msg.Key, body, msg.Status, msg.RetryCount, msg.CreatedAt, msg.UpdatedAt,
    )

    return err
}

// GetPendingMessages 获取待发送消息
func (m *MessageTable) GetPendingMessages(ctx context.Context, limit int) ([]*Message, error) {
    rows, err := m.db.QueryContext(ctx,
        "SELECT id, topic, msg_key, body, status, retry_count, created_at, updated_at FROM message_table WHERE status = ? AND retry_count < ? ORDER BY created_at LIMIT ?",
        "pending", 3, limit,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []*Message
    for rows.Next() {
        var msg Message
        var body []byte
        err := rows.Scan(&msg.ID, &msg.Topic, &msg.Key, &body, &msg.Status, &msg.RetryCount, &msg.CreatedAt, &msg.UpdatedAt)
        if err != nil {
            continue
        }
        msg.Body = body
        messages = append(messages, &msg)
    }

    return messages, nil
}

// UpdateStatus 更新消息状态
func (m *MessageTable) UpdateStatus(ctx context.Context, msgID, status string) error {
    _, err := m.db.ExecContext(ctx,
        "UPDATE message_table SET status = ?, updated_at = ? WHERE id = ?",
        status, time.Now(), msgID,
    )
    return err
}

// IncrementRetry 增加重试次数
func (m *MessageTable) IncrementRetry(ctx context.Context, msgID string) error {
    _, err := m.db.ExecContext(ctx,
        "UPDATE message_table SET retry_count = retry_count + 1, updated_at = ? WHERE id = ?",
        time.Now(), msgID,
    )
    return err
}

// MessagePublisher 消息发布者接口
type MessagePublisher interface {
    Publish(ctx context.Context, topic string, key string, body []byte) error
}

// MessageRelayService 消息投递服务
type MessageRelayService struct {
    msgTable  *MessageTable
    publisher MessagePublisher
    interval  time.Duration
}

func NewMessageRelayService(msgTable *MessageTable, publisher MessagePublisher, interval time.Duration) *MessageRelayService {
    return &MessageRelayService{
        msgTable:  msgTable,
        publisher: publisher,
        interval:  interval,
    }
}

// Start 启动消息投递服务
func (s *MessageRelayService) Start(ctx context.Context) {
    ticker := time.NewTicker(s.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            s.relayMessages(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (s *MessageRelayService) relayMessages(ctx context.Context) {
    messages, err := s.msgTable.GetPendingMessages(ctx, 100)
    if err != nil {
        fmt.Printf("get pending messages failed: %v\n", err)
        return
    }

    for _, msg := range messages {
        err := s.publisher.Publish(ctx, msg.Topic, msg.Key, msg.Body)

        if err != nil {
            // 发送失败，增加重试次数
            s.msgTable.IncrementRetry(ctx, msg.ID)
            fmt.Printf("publish message %s failed: %v\n", msg.ID, err)
        } else {
            // 发送成功，更新状态
            s.msgTable.UpdateStatus(ctx, msg.ID, "sent")
        }
    }
}

// 业务使用示例
type OrderService struct {
    db        *sql.DB
    msgTable  *MessageTable
    publisher MessagePublisher
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, amount float64) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    // 1. 创建订单
    orderID := generateOrderID()
    _, err = tx.ExecContext(ctx,
        "INSERT INTO orders (id, user_id, amount, status) VALUES (?, ?, ?, ?)",
        orderID, userID, amount, "created",
    )
    if err != nil {
        return err
    }

    // 2. 记录消息（同一事务）
    msg := &Message{
        ID:    generateMsgID(),
        Topic: "order_created",
        Key:   orderID,
        Body: mustJSON(map[string]interface{}{
            "order_id": orderID,
            "user_id":  userID,
            "amount":   amount,
        }),
    }

    if err := s.msgTable.SaveMessage(ctx, tx, msg); err != nil {
        return err
    }

    // 3. 提交事务
    if err := tx.Commit(); err != nil {
        return err
    }

    // 4. 同步发送消息（可选，提高实时性）
    go s.publisher.Publish(ctx, msg.Topic, msg.Key, msg.Body)

    return nil
}

func mustJSON(v interface{}) json.RawMessage {
    b, _ := json.Marshal(v)
    return b
}

func generateOrderID() string {
    return fmt.Sprintf("ORD%d", time.Now().UnixNano())
}

func generateMsgID() string {
    return fmt.Sprintf("MSG%d", time.Now().UnixNano())
}
```

#### 反例说明

```go
// ❌ 错误：消息和业务不在同一事务
func CreateOrder(ctx context.Context) error {
    // 先执行业务
    db.Exec("INSERT INTO orders ...")

    // 再记录消息（不同事务）
    msgTable.SaveMessage(ctx, msg)  // 使用新事务
}
// 问题：业务成功但消息记录失败，导致消息丢失

// ❌ 错误：不处理消息发送失败
func relayMessages() {
    for _, msg := range messages {
        publisher.Publish(msg)  // 忽略错误
        msgTable.UpdateStatus(msg.ID, "sent")
    }
}
// 问题：实际发送失败但标记为已发送
```

---

### 7.6 最大努力通知

#### 概念定义

最大努力通知是一种**柔性事务**方案，消息生产者尽力将消息送达消费者，消费者需要保证幂等性处理，适用于对一致性要求不高的场景。

#### 完整示例：Go实现最大努力通知

```go
// besteffort/best_effort.go
package besteffort

import (
    "context"
    "fmt"
    "time"
)

// Notifier 通知器
type Notifier struct {
    maxRetries  int
    retryDelays []time.Duration
    callback    NotifyCallback
}

// NotifyCallback 通知回调
type NotifyCallback func(ctx context.Context, notification *Notification) error

// Notification 通知内容
type Notification struct {
    ID        string
    Type      string
    Payload   []byte
    Attempts  int
    CreatedAt time.Time
}

func NewNotifier(callback NotifyCallback) *Notifier {
    return &Notifier{
        maxRetries: 5,
        retryDelays: []time.Duration{
            1 * time.Second,
            5 * time.Second,
            15 * time.Second,
            1 * time.Minute,
            5 * time.Minute,
        },
        callback: callback,
    }
}

// Notify 发送通知
func (n *Notifier) Notify(ctx context.Context, notification *Notification) {
    go n.doNotify(ctx, notification)
}

func (n *Notifier) doNotify(ctx context.Context, notification *Notification) {
    for attempt := 0; attempt < n.maxRetries; attempt++ {
        notification.Attempts = attempt + 1

        err := n.callback(ctx, notification)
        if err == nil {
            // 通知成功
            return
        }

        // 通知失败，等待后重试
        if attempt < len(n.retryDelays) {
            delay := n.retryDelays[attempt]

            select {
            case <-time.After(delay):
                continue
            case <-ctx.Done():
                return
            }
        }
    }

    // 达到最大重试次数，记录日志或进入死信队列
    fmt.Printf("notification %s failed after %d attempts\n", notification.ID, n.maxRetries)
}

// 使用示例：支付回调通知
func PaymentCallbackNotifier(paymentService *PaymentService) *Notifier {
    return NewNotifier(func(ctx context.Context, n *Notification) error {
        // 调用商户回调接口
        return paymentService.NotifyMerchant(ctx, n.Payload)
    })
}
```

#### 注意事项

1. **幂等性保证**：消费者必须实现幂等处理
2. **重试策略**：指数退避，避免对下游造成压力
3. **死信队列**：最终失败的消息进入死信队列人工处理
4. **监控告警**：通知失败率过高时触发告警

---

## 8. 一致性协议

### 8.1 Raft协议实现

#### 概念定义

Raft是一种**强一致性**的分布式共识算法，通过**领导者选举**、**日志复制**、**安全性**三个子问题保证集群状态一致。

#### 架构图

```
Raft集群架构:

┌─────────────────────────────────────────────────────────────┐
│                        Raft集群                              │
│                                                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐        │
│  │ Leader  │  │Follower │  │Follower │  │Follower │        │
│  │  (领导者)│  │ (跟随者) │  │ (跟随者) │  │ (跟随者) │        │
│  │         │  │         │  │         │  │         │        │
│  │ ┌─────┐ │  │ ┌─────┐ │  │ ┌─────┐ │  │ ┌─────┐ │        │
│  │ │Log 1│ │  │ │Log 1│ │  │ │Log 1│ │  │ │Log 1│ │        │
│  │ │Log 2│ │  │ │Log 2│ │  │ │Log 2│ │  │ │Log 2│ │        │
│  │ │Log 3│ │  │ │Log 3│ │  │ │Log 3│ │  │ │Log 3│ │        │
│  │ └─────┘ │  │ └─────┘ │  │ └─────┘ │  │ └─────┘ │        │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘        │
│       │            │            │            │              │
│       └────────────┴────────────┴────────────┘              │
│                    日志复制(Heartbeat)                        │
└─────────────────────────────────────────────────────────────┘

Raft状态机:
┌─────────┐  超时未收到心跳   ┌─────────┐  收到多数投票   ┌─────────┐
│Follower │ ───────────────> │Candidate│ ─────────────> │ Leader  │
│         │ <─────────────── │         │ <───────────── │         │
│         │   发现更高任期    │         │   发现更高任期  │         │
└─────────┘                  └─────────┘                └─────────┘
```

#### 完整示例：Go实现简化版Raft

```go
// raft/raft.go
package raft

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "sync/atomic"
    "time"
)

// NodeState 节点状态
type NodeState int32

const (
    StateFollower NodeState = iota
    StateCandidate
    StateLeader
)

func (s NodeState) String() string {
    switch s {
    case StateFollower:
        return "Follower"
    case StateCandidate:
        return "Candidate"
    case StateLeader:
        return "Leader"
    default:
        return "Unknown"
    }
}

// LogEntry 日志条目
type LogEntry struct {
    Index   uint64
    Term    uint64
    Command interface{}
}

// RaftNode Raft节点
type RaftNode struct {
    // 节点信息
    id      string
    peers   []string

    // 持久化状态
    currentTerm int64
    votedFor    string
    log         []LogEntry

    // 易失状态
    state        int32
    commitIndex  uint64
    lastApplied  uint64

    // Leader状态
    nextIndex   map[string]uint64
    matchIndex  map[string]uint64

    // 通道
    applyCh     chan ApplyMsg
    heartbeatCh chan struct{}

    // 定时器
    electionTimeout  time.Duration
    heartbeatInterval time.Duration

    // 网络层
    transport Transport

    // 停止信号
    stopCh chan struct{}
    wg     sync.WaitGroup

    mu sync.RWMutex
}

// ApplyMsg 应用消息
type ApplyMsg struct {
    Command interface{}
    Index   uint64
    Term    uint64
}

// Transport 网络传输接口
type Transport interface {
    SendRequestVote(peer string, req *RequestVoteRequest) (*RequestVoteResponse, error)
    SendAppendEntries(peer string, req *AppendEntriesRequest) (*AppendEntriesResponse, error)
}

// RequestVoteRequest 请求投票请求
type RequestVoteRequest struct {
    Term         int64
    CandidateID  string
    LastLogIndex uint64
    LastLogTerm  uint64
}

// RequestVoteResponse 请求投票响应
type RequestVoteResponse struct {
    Term        int64
    VoteGranted bool
}

// AppendEntriesRequest 追加日志请求
type AppendEntriesRequest struct {
    Term         int64
    LeaderID     string
    PrevLogIndex uint64
    PrevLogTerm  uint64
    Entries      []LogEntry
    LeaderCommit uint64
}

// AppendEntriesResponse 追加日志响应
type AppendEntriesResponse struct {
    Term    int64
    Success bool
}

// NewRaftNode 创建Raft节点
func NewRaftNode(id string, peers []string, transport Transport) *RaftNode {
    return &RaftNode{
        id:                id,
        peers:             peers,
        log:               make([]LogEntry, 1), // 索引从1开始
        state:             int32(StateFollower),
        nextIndex:         make(map[string]uint64),
        matchIndex:        make(map[string]uint64),
        applyCh:           make(chan ApplyMsg, 100),
        heartbeatCh:       make(chan struct{}),
        electionTimeout:   randomTimeout(150, 300),
        heartbeatInterval: 50 * time.Millisecond,
        transport:         transport,
        stopCh:            make(chan struct{}),
    }
}

// Start 启动节点
func (n *RaftNode) Start() {
    n.wg.Add(2)
    go n.electionLoop()
    go n.applyLoop()
}

// Stop 停止节点
func (n *RaftNode) Stop() {
    close(n.stopCh)
    n.wg.Wait()
}

// State 获取当前状态
func (n *RaftNode) State() NodeState {
    return NodeState(atomic.LoadInt32(&n.state))
}

// IsLeader 是否是Leader
func (n *RaftNode) IsLeader() bool {
    return n.State() == StateLeader
}

// GetTerm 获取当前任期
func (n *RaftNode) GetTerm() int64 {
    return atomic.LoadInt64(&n.currentTerm)
}

// electionLoop 选举循环
func (n *RaftNode) electionLoop() {
    defer n.wg.Done()

    timer := time.NewTimer(n.electionTimeout)
    defer timer.Stop()

    for {
        select {
        case <-n.stopCh:
            return

        case <-timer.C:
            n.mu.Lock()
            state := n.State()
            n.mu.Unlock()

            if state != StateLeader {
                n.startElection()
            }

            timer.Reset(n.electionTimeout)

        case <-n.heartbeatCh:
            // 收到心跳，重置选举定时器
            if !timer.Stop() {
                <-timer.C
            }
            timer.Reset(n.electionTimeout)
        }
    }
}

// startElection 开始选举
func (n *RaftNode) startElection() {
    n.mu.Lock()

    // 转换为Candidate
    atomic.StoreInt32(&n.state, int32(StateCandidate))
    atomic.AddInt64(&n.currentTerm, 1)
    n.votedFor = n.id

    term := n.currentTerm
    lastLogIndex := uint64(len(n.log) - 1)
    lastLogTerm := uint64(0)
    if lastLogIndex > 0 {
        lastLogTerm = n.log[lastLogIndex].Term
    }

    n.mu.Unlock()

    fmt.Printf("Node %s starting election for term %d\n", n.id, term)

    // 发送请求投票
    votes := 1 // 自己投自己
    var voteMu sync.Mutex

    var wg sync.WaitGroup
    for _, peer := range n.peers {
        if peer == n.id {
            continue
        }

        wg.Add(1)
        go func(peer string) {
            defer wg.Done()

            req := &RequestVoteRequest{
                Term:         term,
                CandidateID:  n.id,
                LastLogIndex: lastLogIndex,
                LastLogTerm:  lastLogTerm,
            }

            resp, err := n.transport.SendRequestVote(peer, req)
            if err != nil {
                return
            }

            // 检查任期
            if resp.Term > term {
                n.stepDown(resp.Term)
                return
            }

            if resp.VoteGranted {
                voteMu.Lock()
                votes++
                voteMu.Unlock()
            }
        }(peer)
    }

    wg.Wait()

    // 检查是否赢得选举
    if votes > len(n.peers)/2 {
        n.becomeLeader()
    }
}

// becomeLeader 成为Leader
func (n *RaftNode) becomeLeader() {
    n.mu.Lock()
    defer n.mu.Unlock()

    if n.State() != StateCandidate {
        return
    }

    atomic.StoreInt32(&n.state, int32(StateLeader))

    // 初始化Leader状态
    lastLogIndex := uint64(len(n.log))
    for _, peer := range n.peers {
        n.nextIndex[peer] = lastLogIndex
        n.matchIndex[peer] = 0
    }

    fmt.Printf("Node %s became leader for term %d\n", n.id, n.currentTerm)

    // 立即发送心跳
    go n.sendHeartbeats()

    // 启动心跳循环
    go n.heartbeatLoop()
}

// heartbeatLoop 心跳循环
func (n *RaftNode) heartbeatLoop() {
    ticker := time.NewTicker(n.heartbeatInterval)
    defer ticker.Stop()

    for {
        select {
        case <-n.stopCh:
            return
        case <-ticker.C:
            if n.State() != StateLeader {
                return
            }
            n.sendHeartbeats()
        }
    }
}

// sendHeartbeats 发送心跳
func (n *RaftNode) sendHeartbeats() {
    n.mu.RLock()
    term := n.currentTerm
    leaderID := n.id
    commitIndex := n.commitIndex
    n.mu.RUnlock()

    var wg sync.WaitGroup
    for _, peer := range n.peers {
        if peer == n.id {
            continue
        }

        wg.Add(1)
        go func(peer string) {
            defer wg.Done()

            req := &AppendEntriesRequest{
                Term:         term,
                LeaderID:     leaderID,
                LeaderCommit: commitIndex,
            }

            resp, err := n.transport.SendAppendEntries(peer, req)
            if err != nil {
                return
            }

            if resp.Term > term {
                n.stepDown(resp.Term)
            }
        }(peer)
    }

    wg.Wait()
}

// stepDown 降级为Follower
func (n *RaftNode) stepDown(term int64) {
    n.mu.Lock()
    defer n.mu.Unlock()

    atomic.StoreInt64(&n.currentTerm, term)
    atomic.StoreInt32(&n.state, int32(StateFollower))
    n.votedFor = ""
}

// HandleRequestVote 处理请求投票
func (n *RaftNode) HandleRequestVote(req *RequestVoteRequest) *RequestVoteResponse {
    n.mu.Lock()
    defer n.mu.Unlock()

    resp := &RequestVoteResponse{Term: n.currentTerm}

    // 任期检查
    if req.Term < n.currentTerm {
        resp.VoteGranted = false
        return resp
    }

    if req.Term > n.currentTerm {
        n.currentTerm = req.Term
        n.state = int32(StateFollower)
        n.votedFor = ""
    }

    // 投票检查
    if n.votedFor != "" && n.votedFor != req.CandidateID {
        resp.VoteGranted = false
        return resp
    }

    // 日志检查
    lastLogIndex := uint64(len(n.log) - 1)
    lastLogTerm := uint64(0)
    if lastLogIndex > 0 {
        lastLogTerm = n.log[lastLogIndex].Term
    }

    if req.LastLogTerm < lastLogTerm ||
        (req.LastLogTerm == lastLogTerm && req.LastLogIndex < lastLogIndex) {
        resp.VoteGranted = false
        return resp
    }

    n.votedFor = req.CandidateID
    resp.VoteGranted = true

    // 重置选举定时器
    select {
    case n.heartbeatCh <- struct{}{}:
    default:
    }

    return resp
}

// HandleAppendEntries 处理追加日志
func (n *RaftNode) HandleAppendEntries(req *AppendEntriesRequest) *AppendEntriesResponse {
    n.mu.Lock()
    defer n.mu.Unlock()

    resp := &AppendEntriesResponse{Term: n.currentTerm}

    // 任期检查
    if req.Term < n.currentTerm {
        resp.Success = false
        return resp
    }

    // 重置选举定时器（收到心跳）
    select {
    case n.heartbeatCh <- struct{}{}:
    default:
    }

    if req.Term > n.currentTerm {
        n.currentTerm = req.Term
        n.state = int32(StateFollower)
        n.votedFor = ""
    }

    // 日志一致性检查
    if req.PrevLogIndex > 0 {
        if req.PrevLogIndex >= uint64(len(n.log)) {
            resp.Success = false
            return resp
        }
        if n.log[req.PrevLogIndex].Term != req.PrevLogTerm {
            resp.Success = false
            return resp
        }
    }

    // 追加日志
    if len(req.Entries) > 0 {
        // 删除冲突的日志
        for i, entry := range req.Entries {
            idx := req.PrevLogIndex + 1 + uint64(i)
            if idx < uint64(len(n.log)) {
                if n.log[idx].Term != entry.Term {
                    n.log = n.log[:idx]
                    n.log = append(n.log, req.Entries[i:]...)
                    break
                }
            } else {
                n.log = append(n.log, req.Entries[i:]...)
                break
            }
        }
    }

    // 更新commitIndex
    if req.LeaderCommit > n.commitIndex {
        lastLogIndex := uint64(len(n.log) - 1)
        if req.LeaderCommit < lastLogIndex {
            n.commitIndex = req.LeaderCommit
        } else {
            n.commitIndex = lastLogIndex
        }
    }

    resp.Success = true
    return resp
}

// Propose 提议命令（Leader调用）
func (n *RaftNode) Propose(command interface{}) (uint64, error) {
    n.mu.Lock()
    defer n.mu.Unlock()

    if n.State() != StateLeader {
        return 0, fmt.Errorf("not leader")
    }

    entry := LogEntry{
        Index:   uint64(len(n.log)),
        Term:    uint64(n.currentTerm),
        Command: command,
    }

    n.log = append(n.log, entry)

    // 异步复制到其他节点
    go n.replicateLog()

    return entry.Index, nil
}

// replicateLog 复制日志到其他节点
func (n *RaftNode) replicateLog() {
    n.mu.RLock()
    term := n.currentTerm
    leaderID := n.id
    commitIndex := n.commitIndex
    n.mu.RUnlock()

    for _, peer := range n.peers {
        if peer == n.id {
            continue
        }

        go func(peer string) {
            n.mu.RLock()
            nextIdx := n.nextIndex[peer]
            entries := make([]LogEntry, 0)
            if nextIdx < uint64(len(n.log)) {
                entries = append(entries, n.log[nextIdx:]...)
            }
            prevLogIndex := nextIdx - 1
            prevLogTerm := uint64(0)
            if prevLogIndex > 0 && prevLogIndex < uint64(len(n.log)) {
                prevLogTerm = n.log[prevLogIndex].Term
            }
            n.mu.RUnlock()

            req := &AppendEntriesRequest{
                Term:         term,
                LeaderID:     leaderID,
                PrevLogIndex: prevLogIndex,
                PrevLogTerm:  prevLogTerm,
                Entries:      entries,
                LeaderCommit: commitIndex,
            }

            resp, err := n.transport.SendAppendEntries(peer, req)
            if err != nil {
                return
            }

            n.mu.Lock()
            defer n.mu.Unlock()

            if resp.Success {
                n.matchIndex[peer] = prevLogIndex + uint64(len(entries))
                n.nextIndex[peer] = n.matchIndex[peer] + 1
            } else {
                if n.nextIndex[peer] > 1 {
                    n.nextIndex[peer]--
                }
            }
        }(peer)
    }
}

// applyLoop 应用日志循环
func (n *RaftNode) applyLoop() {
    defer n.wg.Done()

    ticker := time.NewTicker(10 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-n.stopCh:
            return
        case <-ticker.C:
            n.mu.Lock()
            for n.lastApplied < n.commitIndex {
                n.lastApplied++
                entry := n.log[n.lastApplied]

                n.applyCh <- ApplyMsg{
                    Command: entry.Command,
                    Index:   entry.Index,
                    Term:    entry.Term,
                }
            }
            n.mu.Unlock()
        }
    }
}

// ApplyCh 获取应用通道
func (n *RaftNode) ApplyCh() <-chan ApplyMsg {
    return n.applyCh
}

// randomTimeout 生成随机超时时间
func randomTimeout(min, max int) time.Duration {
    return time.Duration(min+rand.Intn(max-min)) * time.Millisecond
}
```

#### 反例说明

```go
// ❌ 错误：不检查日志一致性就投票
if req.Term >= n.currentTerm {
    n.votedFor = req.CandidateID
    return &RequestVoteResponse{VoteGranted: true}
}
// 问题：可能选举出日志落后的Leader

// ❌ 错误：Leader不发送空心跳
// 只在有日志时发送AppendEntries
// 问题：Follower会频繁触发选举

// ❌ 错误：commitIndex更新不正确
n.commitIndex = req.LeaderCommit
// 问题：LeaderCommit可能大于本地日志长度
```

---

### 8.2 领导者选举

#### 概念定义

领导者选举是Raft的核心机制，通过**随机超时**和**多数投票**保证只有一个Leader，避免脑裂。

#### 选举流程

```go
// raft/election.go
package raft

// ElectionManager 选举管理器
type ElectionManager struct {
    node *RaftNode
}

// ElectionTimeout 选举超时处理
func (e *ElectionManager) OnElectionTimeout() {
    // 1. 增加任期
    e.node.incrementTerm()

    // 2. 转换为Candidate
    e.node.transitionTo(StateCandidate)

    // 3. 投票给自己
    e.node.voteFor(e.node.id)

    // 4. 并行发送RequestVote
    votes := e.node.requestVotes()

    // 5. 检查是否获得多数票
    if votes > len(e.node.peers)/2 {
        e.node.transitionTo(StateLeader)
    }
}

// 选举安全保证
// 1. 每个任期最多一个Leader
// 2. 只有包含全部已提交日志的节点才能当选
// 3. 通过随机超时避免活锁
```

#### 注意事项

1. **随机超时**：150-300ms随机，避免同时选举
2. **日志检查**：投票时检查候选者日志是否最新
3. **任期单调**：任期只增不减
4. **多数原则**：需要超过半数节点投票

---

### 8.3 日志复制

#### 概念定义

日志复制是Leader将日志条目复制到所有Follower，保证**日志一致性**和**状态机安全**。

#### 日志复制流程

```go
// raft/log_replication.go
package raft

// LogManager 日志管理器
type LogManager struct {
    log []LogEntry
}

// Append 追加日志
func (lm *LogManager) Append(entries []LogEntry) {
    lm.log = append(lm.log, entries...)
}

// Truncate 截断日志（用于回滚）
func (lm *LogManager) Truncate(fromIndex uint64) {
    if fromIndex < uint64(len(lm.log)) {
        lm.log = lm.log[:fromIndex]
    }
}

// GetEntries 获取日志条目
func (lm *LogManager) GetEntries(fromIndex uint64) []LogEntry {
    if fromIndex >= uint64(len(lm.log)) {
        return nil
    }
    return lm.log[fromIndex:]
}

// LastEntry 获取最后一条日志
func (lm *LogManager) LastEntry() *LogEntry {
    if len(lm.log) == 0 {
        return nil
    }
    return &lm.log[len(lm.log)-1]
}
```

#### 日志匹配特性

```
日志匹配特性:
如果两个日志条目有相同的索引和任期，则:
1. 它们存储相同的命令
2. 它们之前的所有日志都相同

冲突解决:
Leader发现Follower日志不一致时:
1. 递减nextIndex
2. 重试AppendEntries
3. 直到找到匹配点
4. 删除Follower冲突日志
5. 追加Leader日志
```

---

### 8.4 成员变更

#### 概念定义

成员变更是在不停止服务的情况下，动态增删集群节点。Raft使用**联合共识(Joint Consensus)**保证变更期间的安全性。

#### 成员变更流程

```go
// raft/membership.go
package raft

// Membership 集群成员
type Membership struct {
    nodes map[string]*NodeInfo
    mu    sync.RWMutex
}

type NodeInfo struct {
    ID      string
    Address string
}

// AddNode 添加节点
func (m *Membership) AddNode(node *NodeInfo) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if _, exists := m.nodes[node.ID]; exists {
        return fmt.Errorf("node %s already exists", node.ID)
    }

    m.nodes[node.ID] = node
    return nil
}

// RemoveNode 移除节点
func (m *Membership) RemoveNode(nodeID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if _, exists := m.nodes[nodeID]; !exists {
        return fmt.Errorf("node %s not found", nodeID)
    }

    delete(m.nodes, nodeID)
    return nil
}

// GetNodes 获取所有节点
func (m *Membership) GetNodes() []string {
    m.mu.RLock()
    defer m.mu.RUnlock()

    nodes := make([]string, 0, len(m.nodes))
    for id := range m.nodes {
        nodes = append(nodes, id)
    }
    return nodes
}

// JointConsensus 联合共识
type JointConsensus struct {
    oldConfig []string
    newConfig []string
}

// 联合共识阶段:
// 1. Cold,new: 需要Cold和Cnew都多数同意
// 2. Cnew: 只需要Cnew多数同意
```

#### 注意事项

1. **联合共识**：新旧配置同时生效期间，需要两者都多数同意
2. **Leader切换**：成员变更期间可能需要切换Leader
3. **配置持久化**：配置变更需要持久化，用于故障恢复
4. **单节点变更**：一次只变更一个节点，简化实现

---

## 9. 分布式缓存

### 9.1 缓存穿透、击穿、雪崩

#### 概念定义

缓存问题三件套：

- **穿透**：查询不存在的数据，绕过缓存直接查数据库
- **击穿**：热点key过期，大量请求同时查数据库
- **雪崩**：大量key同时过期，数据库压力剧增

#### 架构图

```
缓存穿透:
客户端 ──> 查询key="xxx" ──> 缓存未命中
                              │
                              ▼
                         查询数据库
                              │
                              ▼
                         数据库不存在
                              │
                              ▼
                         不写入缓存
                              │
客户端 ──> 查询key="xxx" ──> 缓存未命中 (循环)

解决方案: 布隆过滤器或缓存空值

缓存击穿:
                    key过期
                      │
客户端1 ──> 查询key ──> 缓存未命中
客户端2 ──> 查询key ──> 缓存未命中
客户端3 ──> 查询key ──> 缓存未命中
...N个并发请求
                      │
                      ▼
                  同时查数据库
                      │
                      ▼
                  数据库压力剧增

解决方案: 互斥锁或逻辑过期

缓存雪崩:
时间: t1    t2    t3    t4    t5
      │     │     │     │     │
      ▼     ▼     ▼     ▼     ▼
     key1  key2  key3  key4  key5
     过期   过期   过期   过期   过期
      │     │     │     │     │
      └─────┴─────┴─────┴─────┘
                  │
                  ▼
           大量请求到数据库

解决方案: 随机过期时间或热点key永不过期
```

#### 完整示例：Go解决缓存三大问题

```go
// cache/cache_problems.go
package cache

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "sync"
    "time"

    "github.com/redis/go-redis/v9"
)

// Cache 缓存接口
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
    SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
    Delete(ctx context.Context, key string) error
}

// RedisCache Redis缓存实现
type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(client *redis.Client) *Cache {
    return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
    return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
    return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
    return r.client.SetNX(ctx, key, value, ttl).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
    return r.client.Del(ctx, key).Err()
}

// ========== 缓存穿透解决方案 ==========

// PenetrationGuard 缓存穿透防护
type PenetrationGuard struct {
    cache       Cache
    bloomFilter *BloomFilter
}

// 方案1: 缓存空值
func (g *PenetrationGuard) GetWithNullCache(ctx context.Context, key string, fetch func() (string, error)) (string, error) {
    // 1. 查缓存
    value, err := g.cache.Get(ctx, key)
    if err == nil {
        // 检查是否是空值标记
        if value == "__NULL__" {
            return "", errors.New("key not found")
        }
        return value, nil
    }

    // 2. 查数据库
    value, err = fetch()
    if err != nil {
        // 3. 缓存空值（短时间）
        g.cache.Set(ctx, key, "__NULL__", 5*time.Minute)
        return "", err
    }

    // 4. 缓存结果
    g.cache.Set(ctx, key, value, 10*time.Minute)
    return value, nil
}

// 方案2: 布隆过滤器
func (g *PenetrationGuard) GetWithBloomFilter(ctx context.Context, key string, fetch func() (string, error)) (string, error) {
    // 1. 布隆过滤器检查
    if !g.bloomFilter.MayContain(key) {
        return "", errors.New("key not found")
    }

    // 2. 查缓存
    value, err := g.cache.Get(ctx, key)
    if err == nil {
        return value, nil
    }

    // 3. 查数据库
    value, err = fetch()
    if err != nil {
        return "", err
    }

    // 4. 缓存结果
    g.cache.Set(ctx, key, value, 10*time.Minute)
    return value, nil
}

// BloomFilter 布隆过滤器
type BloomFilter struct {
    bits   []bool
    hashes int
    size   int
}

func NewBloomFilter(size, hashes int) *BloomFilter {
    return &BloomFilter{
        bits:   make([]bool, size),
        hashes: hashes,
        size:   size,
    }
}

func (b *BloomFilter) Add(key string) {
    for i := 0; i < b.hashes; i++ {
        idx := b.hash(key, i) % b.size
        b.bits[idx] = true
    }
}

func (b *BloomFilter) MayContain(key string) bool {
    for i := 0; i < b.hashes; i++ {
        idx := b.hash(key, i) % b.size
        if !b.bits[idx] {
            return false
        }
    }
    return true
}

func (b *BloomFilter) hash(key string, seed int) int {
    h := 0
    for i, c := range key {
        h += int(c) * (seed + i + 1)
    }
    return h
}

// ========== 缓存击穿解决方案 ==========

// HotKeyGuard 热点key防护
type HotKeyGuard struct {
    cache Cache
    locks map[string]*sync.Mutex
    mu    sync.RWMutex
}

func NewHotKeyGuard(cache Cache) *HotKeyGuard {
    return &HotKeyGuard{
        cache: cache,
        locks: make(map[string]*sync.Mutex),
    }
}

// 方案1: 互斥锁
func (g *HotKeyGuard) GetWithMutex(ctx context.Context, key string, fetch func() (string, error)) (string, error) {
    // 1. 查缓存
    value, err := g.cache.Get(ctx, key)
    if err == nil {
        return value, nil
    }

    // 2. 获取锁
    lock := g.getLock(key)
    lock.Lock()
    defer lock.Unlock()

    // 3. 双重检查
    value, err = g.cache.Get(ctx, key)
    if err == nil {
        return value, nil
    }

    // 4. 查数据库
    value, err = fetch()
    if err != nil {
        return "", err
    }

    // 5. 缓存结果
    g.cache.Set(ctx, key, value, 10*time.Minute)
    return value, nil
}

func (g *HotKeyGuard) getLock(key string) *sync.Mutex {
    g.mu.RLock()
    lock, ok := g.locks[key]
    g.mu.RUnlock()

    if ok {
        return lock
    }

    g.mu.Lock()
    defer g.mu.Unlock()

    // 双重检查
    if lock, ok := g.locks[key]; ok {
        return lock
    }

    lock = &sync.Mutex{}
    g.locks[key] = lock
    return lock
}

// 方案2: 逻辑过期
func (g *HotKeyGuard) GetWithLogicalExpire(ctx context.Context, key string, fetch func() (string, error)) (string, error) {
    type CacheValue struct {
        Data       string    `json:"data"`
        ExpireTime time.Time `json:"expire_time"`
    }

    // 1. 查缓存
    valStr, err := g.cache.Get(ctx, key)
    if err != nil {
        // 缓存不存在，直接查数据库
        value, err := fetch()
        if err != nil {
            return "", err
        }

        cv := CacheValue{
            Data:       value,
            ExpireTime: time.Now().Add(10 * time.Minute),
        }
        data, _ := json.Marshal(cv)
        g.cache.Set(ctx, key, string(data), 0) // 永不过期
        return value, nil
    }

    var cv CacheValue
    if err := json.Unmarshal([]byte(valStr), &cv); err != nil {
        return "", err
    }

    // 2. 检查是否过期
    if cv.ExpireTime.After(time.Now()) {
        return cv.Data, nil
    }

    // 3. 已过期，尝试获取锁重建缓存
    lock := g.getLock(key + ":rebuild")
    if lock.TryLock() {
        defer lock.Unlock()

        // 异步重建缓存
        go func() {
            value, err := fetch()
            if err != nil {
                return
            }

            newCv := CacheValue{
                Data:       value,
                ExpireTime: time.Now().Add(10 * time.Minute),
            }
            data, _ := json.Marshal(newCv)
            g.cache.Set(ctx, key, string(data), 0)
        }()
    }

    // 4. 返回旧值
    return cv.Data, nil
}

// ========== 缓存雪崩解决方案 ==========

// AvalancheGuard 缓存雪崩防护
type AvalancheGuard struct {
    cache Cache
}

// 方案1: 随机过期时间
func (g *AvalancheGuard) SetWithRandomTTL(ctx context.Context, key string, value string, baseTTL time.Duration) error {
    // 添加随机偏移: 基础TTL ± 10%
    offset := time.Duration(rand.Intn(20)-10) * baseTTL / 100
    ttl := baseTTL + offset

    return g.cache.Set(ctx, key, value, ttl)
}

// 方案2: 热点key永不过期 + 定时刷新
func (g *AvalancheGuard) SetHotKey(ctx context.Context, key string, fetch func() (string, error), refreshInterval time.Duration) {
    // 初始加载
    value, _ := fetch()
    g.cache.Set(ctx, key, value, 0) // 永不过期

    // 定时刷新
    go func() {
        ticker := time.NewTicker(refreshInterval)
        defer ticker.Stop()

        for range ticker.C {
            value, err := fetch()
            if err != nil {
                continue
            }
            g.cache.Set(ctx, key, value, 0)
        }
    }()
}

// 方案3: 多级缓存 + 熔断
type MultiLevelCache struct {
    local  *LocalCache    // 本地缓存(Caffeine)
    remote Cache          // 远程缓存(Redis)
    breaker *CircuitBreaker
}

func (m *MultiLevelCache) Get(ctx context.Context, key string, fetch func() (string, error)) (string, error) {
    // 1. 查本地缓存
    if value, ok := m.local.Get(key); ok {
        return value, nil
    }

    // 2. 查远程缓存
    value, err := m.remote.Get(ctx, key)
    if err == nil {
        // 回填本地缓存
        m.local.Set(key, value, 1*time.Minute)
        return value, nil
    }

    // 3. 熔断检查
    if m.breaker.State() == StateOpen {
        return "", errors.New("circuit breaker open")
    }

    // 4. 查数据库
    value, err = fetch()
    if err != nil {
        m.breaker.RecordFailure()
        return "", err
    }

    m.breaker.RecordSuccess()

    // 5. 回填缓存
    m.remote.Set(ctx, key, value, 10*time.Minute)
    m.local.Set(key, value, 1*time.Minute)

    return value, nil
}

// LocalCache 本地缓存
type LocalCache struct {
    data map[string]cacheItem
    mu   sync.RWMutex
}

type cacheItem struct {
    value      string
    expireTime time.Time
}

func NewLocalCache() *LocalCache {
    lc := &LocalCache{
        data: make(map[string]cacheItem),
    }
    go lc.cleanup()
    return lc
}

func (lc *LocalCache) Get(key string) (string, bool) {
    lc.mu.RLock()
    item, ok := lc.data[key]
    lc.mu.RUnlock()

    if !ok || time.Now().After(item.expireTime) {
        return "", false
    }

    return item.value, true
}

func (lc *LocalCache) Set(key, value string, ttl time.Duration) {
    lc.mu.Lock()
    lc.data[key] = cacheItem{
        value:      value,
        expireTime: time.Now().Add(ttl),
    }
    lc.mu.Unlock()
}

func (lc *LocalCache) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        lc.mu.Lock()
        now := time.Now()
        for key, item := range lc.data {
            if now.After(item.expireTime) {
                delete(lc.data, key)
            }
        }
        lc.mu.Unlock()
    }
}
```

#### 反例说明

```go
// ❌ 错误：缓存空值不设置过期时间
g.cache.Set(ctx, key, "__NULL__", 0)  // 永不过期
// 问题：空值缓存永久占用内存

// ❌ 错误：互斥锁粒度太大
var globalLock sync.Mutex  // 全局锁
func Get(key string) {
    globalLock.Lock()  // 所有key串行
    defer globalLock.Unlock()
    // ...
}
// 问题：性能极差

// ❌ 错误：热点key设置过期时间
g.cache.Set(ctx, "hot_key", value, 10*time.Minute)
// 问题：过期时大量请求打到数据库
```

---

### 9.2 缓存一致性

#### 概念定义

缓存一致性保证缓存和数据库数据一致，主要策略有：**Cache-Aside**、**Read-Through**、**Write-Through**、**Write-Behind**。

#### 架构图

```
Cache-Aside模式:

读取:
客户端 ──> 查缓存 ──> 命中? ──> 是 ──> 返回数据
              │
              否
              │
              ▼
         查数据库
              │
              ▼
         写入缓存
              │
              ▼
         返回数据

写入:
客户端 ──> 更新数据库 ──> 删除缓存
                              │
                              ▼
                         下次读取时重建

Write-Through模式:

写入:
客户端 ──> 更新缓存 ──> 同步更新数据库
                              │
                              ▼
                         返回成功

Write-Behind模式:

写入:
客户端 ──> 更新缓存 ──> 异步写入数据库
                              │
                              ▼
                         立即返回
```

#### 完整示例：Go实现缓存一致性

```go
// cache/consistency.go
package cache

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "time"
)

// CacheStrategy 缓存策略接口
type CacheStrategy interface {
    Get(ctx context.Context, key string, fetch func() (interface{}, error)) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}) error
    Delete(ctx context.Context, key string) error
}

// CacheAside Cache-Aside策略
type CacheAside struct {
    cache Cache
    ttl   time.Duration
}

func NewCacheAside(cache Cache, ttl time.Duration) *CacheAside {
    return &CacheAside{
        cache: cache,
        ttl:   ttl,
    }
}

// Get 读取数据
func (c *CacheAside) Get(ctx context.Context, key string, fetch func() (interface{}, error)) (interface{}, error) {
    // 1. 查缓存
    valStr, err := c.cache.Get(ctx, key)
    if err == nil {
        var value interface{}
        if err := json.Unmarshal([]byte(valStr), &value); err == nil {
            return value, nil
        }
    }

    // 2. 查数据库
    value, err := fetch()
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    c.setCache(ctx, key, value)

    return value, nil
}

// Set 更新数据（先更新数据库，再删缓存）
func (c *CacheAside) Set(ctx context.Context, key string, value interface{}, updateDB func() error) error {
    // 1. 更新数据库
    if err := updateDB(); err != nil {
        return err
    }

    // 2. 删除缓存
    c.cache.Delete(ctx, key)

    return nil
}

func (c *CacheAside) setCache(ctx context.Context, key string, value interface{}) {
    data, _ := json.Marshal(value)
    c.cache.Set(ctx, key, string(data), c.ttl)
}

func (c *CacheAside) Delete(ctx context.Context, key string) error {
    return c.cache.Delete(ctx, key)
}

// WriteThrough Write-Through策略
type WriteThrough struct {
    cache Cache
    db    *sql.DB
    ttl   time.Duration
}

func (w *WriteThrough) Set(ctx context.Context, key string, value interface{}) error {
    // 1. 序列化
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    // 2. 开启事务
    tx, err := w.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 3. 更新数据库
    if err := w.updateDB(ctx, tx, key, value); err != nil {
        return err
    }

    // 4. 更新缓存
    if err := w.cache.Set(ctx, key, string(data), w.ttl); err != nil {
        return err
    }

    // 5. 提交事务
    return tx.Commit()
}

func (w *WriteThrough) updateDB(ctx context.Context, tx *sql.Tx, key string, value interface{}) error {
    // 实现数据库更新
    return nil
}

// WriteBehind Write-Behind策略
type WriteBehind struct {
    cache   Cache
    db      *sql.DB
    queue   chan *WriteRequest
    workers int
}

type WriteRequest struct {
    Key   string
    Value interface{}
}

func NewWriteBehind(cache Cache, db *sql.DB, workers int, queueSize int) *WriteBehind {
    wb := &WriteBehind{
        cache:   cache,
        db:      db,
        queue:   make(chan *WriteRequest, queueSize),
        workers: workers,
    }

    // 启动异步写入工作线程
    for i := 0; i < workers; i++ {
        go wb.worker()
    }

    return wb
}

func (w *WriteBehind) Set(ctx context.Context, key string, value interface{}) error {
    // 1. 更新缓存
    data, _ := json.Marshal(value)
    if err := w.cache.Set(ctx, key, string(data), 0); err != nil {
        return err
    }

    // 2. 异步写入数据库
    select {
    case w.queue <- &WriteRequest{Key: key, Value: value}:
        return nil
    default:
        return errors.New("write queue full")
    }
}

func (w *WriteBehind) worker() {
    for req := range w.queue {
        // 批量写入数据库
        w.batchWrite([]*WriteRequest{req})
    }
}

func (w *WriteBehind) batchWrite(requests []*WriteRequest) {
    // 批量写入实现
}

// 双删策略（解决缓存不一致）
type DoubleDeleteStrategy struct {
    cache Cache
}

func (d *DoubleDeleteStrategy) Update(ctx context.Context, key string, updateDB func() error) error {
    // 1. 先删缓存
    d.cache.Delete(ctx, key)

    // 2. 更新数据库
    if err := updateDB(); err != nil {
        return err
    }

    // 3. 延迟再删缓存
    go func() {
        time.Sleep(500 * time.Millisecond)
        d.cache.Delete(ctx, key)
    }()

    return nil
}

// Canal订阅MySQL binlog实现缓存同步
type CanalSync struct {
    cache Cache
}

func (c *CanalSync) OnRowChange(event *CanalEvent) {
    // 根据binlog事件更新缓存
    switch event.EventType {
    case "INSERT", "UPDATE":
        key := fmt.Sprintf("user:%d", event.Row["id"])
        data, _ := json.Marshal(event.Row)
        c.cache.Set(context.Background(), key, string(data), 10*time.Minute)
    case "DELETE":
        key := fmt.Sprintf("user:%d", event.Row["id"])
        c.cache.Delete(context.Background(), key)
    }
}
```

---

### 9.3 本地缓存+分布式缓存

#### 完整示例：Go实现多级缓存

```go
// cache/multilevel.go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// MultiLevelCache 多级缓存
type MultiLevelCache struct {
    l1      *LocalCache      // Caffeine/Go-Cache
    l2      Cache            // Redis
    l1TTL   time.Duration
    l2TTL   time.Duration
    stats   *CacheStats
}

// CacheStats 缓存统计
type CacheStats struct {
    L1Hits uint64
    L2Hits uint64
    Misses uint64
    mu     sync.RWMutex
}

func NewMultiLevelCache(l2 Cache, l1TTL, l2TTL time.Duration) *MultiLevelCache {
    return &MultiLevelCache{
        l1:    NewLocalCache(),
        l2:    l2,
        l1TTL: l1TTL,
        l2TTL: l2TTL,
        stats: &CacheStats{},
    }
}

// Get 读取数据（多级缓存）
func (m *MultiLevelCache) Get(ctx context.Context, key string, fetch func() (interface{}, error)) (interface{}, error) {
    // 1. L1缓存
    if val, ok := m.l1.Get(key); ok {
        m.stats.mu.Lock()
        m.stats.L1Hits++
        m.stats.mu.Unlock()
        return val, nil
    }

    // 2. L2缓存
    valStr, err := m.l2.Get(ctx, key)
    if err == nil {
        var value interface{}
        if err := json.Unmarshal([]byte(valStr), &value); err == nil {
            // 回填L1
            m.l1.Set(key, value, m.l1TTL)

            m.stats.mu.Lock()
            m.stats.L2Hits++
            m.stats.mu.Unlock()
            return value, nil
        }
    }

    // 3. 数据库
    value, err := fetch()
    if err != nil {
        m.stats.mu.Lock()
        m.stats.Misses++
        m.stats.mu.Unlock()
        return nil, err
    }

    // 4. 回填缓存
    m.set(ctx, key, value)

    return value, nil
}

// Set 设置缓存
func (m *MultiLevelCache) set(ctx context.Context, key string, value interface{}) {
    // L1
    m.l1.Set(key, value, m.l1TTL)

    // L2
    data, _ := json.Marshal(value)
    m.l2.Set(ctx, key, string(data), m.l2TTL)
}

// Delete 删除缓存
func (m *MultiLevelCache) Delete(ctx context.Context, key string) error {
    m.l1.Delete(key)
    return m.l2.Delete(ctx, key)
}

// Invalidate 使缓存失效（广播）
func (m *MultiLevelCache) Invalidate(ctx context.Context, pattern string) error {
    // 1. 删除L2缓存
    if err := m.l2.Delete(ctx, pattern); err != nil {
        return err
    }

    // 2. 发布失效消息（让所有节点删除L1）
    msg := &InvalidateMsg{
        Pattern: pattern,
        Time:    time.Now(),
    }
    data, _ := json.Marshal(msg)

    return m.l2.(*RedisCache).client.Publish(ctx, "cache:invalidate", data).Err()
}

// InvalidateMsg 失效消息
type InvalidateMsg struct {
    Pattern string    `json:"pattern"`
    Time    time.Time `json:"time"`
}

// 监听失效消息
func (m *MultiLevelCache) StartInvalidationListener(ctx context.Context) {
    pubsub := m.l2.(*RedisCache).client.Subscribe(ctx, "cache:invalidate")
    defer pubsub.Close()

    for msg := range pubsub.Channel() {
        var invalidateMsg InvalidateMsg
        if err := json.Unmarshal([]byte(msg.Payload), &invalidateMsg); err != nil {
            continue
        }

        // 删除匹配的L1缓存
        m.l1.DeleteByPattern(invalidateMsg.Pattern)
    }
}

// Stats 获取统计信息
func (m *MultiLevelCache) Stats() CacheStats {
    m.stats.mu.RLock()
    defer m.stats.mu.RUnlock()
    return *m.stats
}
```

---

### 9.4 Cache-Aside模式

#### 完整示例：Go实现Cache-Aside

```go
// cache/cacheaside.go
package cache

import (
    "context"
    "encoding/json"
    "errors"
    "sync"
    "time"
)

// CacheAsideManager Cache-Aside管理器
type CacheAsideManager struct {
    cache      Cache
    locks      map[string]*sync.Mutex
    mu         sync.RWMutex
    ttl        time.Duration
    nullTTL    time.Duration
}

func NewCacheAsideManager(cache Cache, ttl time.Duration) *CacheAsideManager {
    return &CacheAsideManager{
        cache:   cache,
        locks:   make(map[string]*sync.Mutex),
        ttl:     ttl,
        nullTTL: 5 * time.Minute,
    }
}

// Get 读取数据
func (c *CacheAsideManager) Get(ctx context.Context, key string, fetch func() (interface{}, error)) (interface{}, error) {
    // 1. 查缓存
    valStr, err := c.cache.Get(ctx, key)
    if err == nil {
        // 检查空值
        if valStr == "__NULL__" {
            return nil, errors.New("key not found")
        }

        var value interface{}
        if err := json.Unmarshal([]byte(valStr), &value); err == nil {
            return value, nil
        }
    }

    // 2. 获取锁（防止缓存击穿）
    lock := c.getLock(key)
    lock.Lock()
    defer lock.Unlock()

    // 3. 双重检查
    valStr, err = c.cache.Get(ctx, key)
    if err == nil {
        if valStr == "__NULL__" {
            return nil, errors.New("key not found")
        }

        var value interface{}
        if err := json.Unmarshal([]byte(valStr), &value); err == nil {
            return value, nil
        }
    }

    // 4. 查数据库
    value, err := fetch()
    if err != nil {
        // 缓存空值（防止穿透）
        c.cache.Set(ctx, key, "__NULL__", c.nullTTL)
        return nil, err
    }

    // 5. 写入缓存
    c.setCache(ctx, key, value)

    return value, nil
}

// Set 更新数据
func (c *CacheAsideManager) Set(ctx context.Context, key string, value interface{}, updateDB func() error) error {
    // 1. 更新数据库
    if err := updateDB(); err != nil {
        return err
    }

    // 2. 删除缓存
    c.cache.Delete(ctx, key)

    return nil
}

// Delete 删除数据
func (c *CacheAsideManager) Delete(ctx context.Context, key string, deleteDB func() error) error {
    // 1. 删除数据库
    if err := deleteDB(); err != nil {
        return err
    }

    // 2. 删除缓存
    c.cache.Delete(ctx, key)

    return nil
}

func (c *CacheAsideManager) getLock(key string) *sync.Mutex {
    c.mu.RLock()
    lock, ok := c.locks[key]
    c.mu.RUnlock()

    if ok {
        return lock
    }

    c.mu.Lock()
    defer c.mu.Unlock()

    if lock, ok := c.locks[key]; ok {
        return lock
    }

    lock = &sync.Mutex{}
    c.locks[key] = lock
    return lock
}

func (c *CacheAsideManager) setCache(ctx context.Context, key string, value interface{}) {
    data, _ := json.Marshal(value)
    c.cache.Set(ctx, key, string(data), c.ttl)
}

// BatchGet 批量读取
func (c *CacheAsideManager) BatchGet(ctx context.Context, keys []string, fetch func([]string) (map[string]interface{}, error)) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    missingKeys := make([]string, 0)

    // 1. 批量查缓存
    for _, key := range keys {
        valStr, err := c.cache.Get(ctx, key)
        if err == nil && valStr != "__NULL__" {
            var value interface{}
            if err := json.Unmarshal([]byte(valStr), &value); err == nil {
                result[key] = value
                continue
            }
        }
        missingKeys = append(missingKeys, key)
    }

    // 2. 批量查数据库
    if len(missingKeys) > 0 {
        dbResults, err := fetch(missingKeys)
        if err != nil {
            return result, err
        }

        // 3. 回填缓存
        for key, value := range dbResults {
            result[key] = value
            c.setCache(ctx, key, value)
        }
    }

    return result, nil
}
```

---

## 10. 分布式锁

### 10.1 基于Redis的RedLock

#### 概念定义

RedLock是Redis作者提出的分布式锁算法，在**多个独立Redis节点**上获取锁，保证锁的**安全性和可用性**。

#### 架构图

```
RedLock算法:

┌─────────┐     获取锁      ┌─────────┐
│ 客户端   │ ─────────────> │ Redis A │
│         │     获取锁      │         │
│         │ ─────────────> │ Redis B │
│         │     获取锁      │         │
│         │ ─────────────> │ Redis C │
│         │     获取锁      │         │
│         │ ─────────────> │ Redis D │
│         │     获取锁      │         │
│         │ ─────────────> │ Redis E │
└─────────┘                └─────────┘

算法步骤:
1. 获取当前时间戳
2. 依次向N个Redis节点获取锁
   - 使用SET key value NX PX timeout
   - 超时时间要远小于锁过期时间
3. 计算获取锁消耗的时间
4. 如果成功获取锁的节点数 > N/2 且 耗时 < 锁过期时间
   - 则认为获取锁成功
5. 否则，向所有节点释放锁

释放锁:
使用Lua脚本原子性删除:
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
```

#### 完整示例：Go实现RedLock

```go
// redlock/redlock.go
package redlock

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "fmt"
    "sync"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedLock RedLock分布式锁
type RedLock struct {
    clients     []*redis.Client
    quorum      int
    retryCount  int
    retryDelay  time.Duration
    driftFactor float64
}

// Lock 锁信息
type Lock struct {
    resource   string
    value      string
    validity   time.Duration
    until      time.Time
    redlock    *RedLock
}

// NewRedLock 创建RedLock
func NewRedLock(addrs []string, password string) (*RedLock, error) {
    clients := make([]*redis.Client, len(addrs))
    for i, addr := range addrs {
        clients[i] = redis.NewClient(&redis.Options{
            Addr:     addr,
            Password: password,
        })
    }

    return &RedLock{
        clients:     clients,
        quorum:      len(clients)/2 + 1,
        retryCount:  3,
        retryDelay:  200 * time.Millisecond,
        driftFactor: 0.01,
    }, nil
}

// Lock 获取锁
func (r *RedLock) Lock(ctx context.Context, resource string, ttl time.Duration) (*Lock, error) {
    value := generateUniqueValue()

    for i := 0; i < r.retryCount; i++ {
        start := time.Now()

        // 尝试在多数节点获取锁
        successes := r.lockOnMajority(ctx, resource, value, ttl)

        // 计算耗时
        elapsed := time.Since(start)
        validity := ttl - elapsed - time.Duration(float64(ttl)*r.driftFactor)

        // 检查是否成功
        if successes >= r.quorum && validity > 0 {
            return &Lock{
                resource: resource,
                value:    value,
                validity: validity,
                until:    time.Now().Add(validity),
                redlock:  r,
            }, nil
        }

        // 失败则释放所有锁
        r.unlockOnAll(ctx, resource, value)

        // 等待后重试
        if i < r.retryCount-1 {
            time.Sleep(r.retryDelay)
        }
    }

    return nil, errors.New("failed to acquire lock")
}

// lockOnMajority 在多数节点获取锁
func (r *RedLock) lockOnMajority(ctx context.Context, resource, value string, ttl time.Duration) int {
    successes := 0
    var mu sync.Mutex
    var wg sync.WaitGroup

    for _, client := range r.clients {
        wg.Add(1)
        go func(c *redis.Client) {
            defer wg.Done()

            // 使用SET NX PX获取锁
            ok, err := c.SetNX(ctx, resource, value, ttl).Result()
            if err == nil && ok {
                mu.Lock()
                successes++
                mu.Unlock()
            }
        }(client)
    }

    wg.Wait()
    return successes
}

// unlockOnAll 在所有节点释放锁
func (r *RedLock) unlockOnAll(ctx context.Context, resource, value string) {
    var wg sync.WaitGroup

    for _, client := range r.clients {
        wg.Add(1)
        go func(c *redis.Client) {
            defer wg.Done()
            r.unlockOnNode(ctx, c, resource, value)
        }(client)
    }

    wg.Wait()
}

// unlockOnNode 在单个节点释放锁
func (r *RedLock) unlockOnNode(ctx context.Context, client *redis.Client, resource, value string) {
    // 使用Lua脚本保证原子性
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
    client.Eval(ctx, script, []string{resource}, value)
}

// Unlock 释放锁
func (l *Lock) Unlock(ctx context.Context) error {
    l.redlock.unlockOnAll(ctx, l.resource, l.value)
    return nil
}

// Extend 延长锁有效期
func (l *Lock) Extend(ctx context.Context, ttl time.Duration) error {
    // 检查锁是否仍有效
    if time.Now().After(l.until) {
        return errors.New("lock has expired")
    }

    successes := 0
    var mu sync.Mutex
    var wg sync.WaitGroup

    for _, client := range l.redlock.clients {
        wg.Add(1)
        go func(c *redis.Client) {
            defer wg.Done()

            // 使用Lua脚本延长过期时间
            script := `
                if redis.call("get", KEYS[1]) == ARGV[1] then
                    return redis.call("pexpire", KEYS[1], ARGV[2])
                else
                    return 0
                end
            `
            ok, err := c.Eval(ctx, script, []string{l.resource}, l.value, ttl.Milliseconds()).Result()
            if err == nil && ok.(int64) == 1 {
                mu.Lock()
                successes++
                mu.Unlock()
            }
        }(client)
    }

    wg.Wait()

    if successes < l.redlock.quorum {
        return errors.New("failed to extend lock")
    }

    l.validity = ttl
    l.until = time.Now().Add(ttl)

    return nil
}

// IsValid 检查锁是否有效
func (l *Lock) IsValid() bool {
    return time.Now().Before(l.until)
}

func generateUniqueValue() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}

// 使用示例
func ExampleRedLock() {
    addrs := []string{
        "localhost:6379",
        "localhost:6380",
        "localhost:6381",
    }

    redlock, err := NewRedLock(addrs, "")
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // 获取锁
    lock, err := redlock.Lock(ctx, "my_resource", 10*time.Second)
    if err != nil {
        panic(err)
    }

    // 执行业务逻辑
    fmt.Println("Lock acquired, doing work...")

    // 延长锁（如果业务耗时较长）
    if err := lock.Extend(ctx, 10*time.Second); err != nil {
        fmt.Printf("Failed to extend lock: %v\n", err)
    }

    // 释放锁
    lock.Unlock(ctx)
}
```

#### 反例说明

```go
// ❌ 错误：单Redis节点
// 使用单个Redis节点实现分布式锁
// 问题：Redis故障时锁失效

// ❌ 错误：不使用Lua脚本释放锁
if client.Get(ctx, key).Val() == value {
    client.Del(ctx, key)  // 不是原子操作
}
// 问题：GET和DEL之间锁可能被其他客户端获取

// ❌ 错误：锁过期时间太短
Lock(ctx, "resource", 1*time.Second)
// 问题：业务未完成锁已过期，其他客户端可获取锁
```

---

### 10.2 基于etcd的分布式锁

#### 概念定义

etcd基于**Raft协议**保证一致性，提供原生的分布式锁支持，通过**租约(Lease)**和**前缀机制**实现锁。

#### 完整示例：Go实现etcd分布式锁

```go
// etcdlock/etcd_lock.go
package etcdlock

import (
    "context"
    "errors"
    "fmt"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
    "go.etcd.io/etcd/client/v3/concurrency"
)

// EtcdLock etcd分布式锁
type EtcdLock struct {
    client  *clientv3.Client
    session *concurrency.Session
    mutex   *concurrency.Mutex
    key     string
}

// NewEtcdLock 创建etcd锁
func NewEtcdLock(endpoints []string, key string, ttl int) (*EtcdLock, error) {
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        return nil, err
    }

    // 创建session
    session, err := concurrency.NewSession(cli, concurrency.WithTTL(ttl))
    if err != nil {
        cli.Close()
        return nil, err
    }

    // 创建互斥锁
    mutex := concurrency.NewMutex(session, key)

    return &EtcdLock{
        client:  cli,
        session: session,
        mutex:   mutex,
        key:     key,
    }, nil
}

// Lock 获取锁（阻塞）
func (l *EtcdLock) Lock(ctx context.Context) error {
    return l.mutex.Lock(ctx)
}

// TryLock 尝试获取锁（非阻塞）
func (l *EtcdLock) TryLock(ctx context.Context, timeout time.Duration) (bool, error) {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    err := l.mutex.Lock(ctx)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return false, nil
        }
        return false, err
    }
    return true, nil
}

// Unlock 释放锁
func (l *EtcdLock) Unlock(ctx context.Context) error {
    return l.mutex.Unlock(ctx)
}

// Close 关闭锁
func (l *EtcdLock) Close() error {
    return l.session.Close()
}

// 基于租约的锁实现
type LeaseLock struct {
    client *clientv3.Client
    lease  clientv3.LeaseID
    key    string
    value  string
}

// NewLeaseLock 创建租约锁
func NewLeaseLock(client *clientv3.Client, key string, ttl int64) (*LeaseLock, error) {
    // 创建租约
    lease, err := client.Grant(context.Background(), ttl)
    if err != nil {
        return nil, err
    }

    return &LeaseLock{
        client: client,
        lease:  lease.ID,
        key:    key,
        value:  generateUniqueValue(),
    }, nil
}

// Lock 获取锁
func (l *LeaseLock) Lock(ctx context.Context) error {
    // 使用事务原子性创建key
    txn := l.client.Txn(ctx)
    txn.If(clientv3.Compare(clientv3.CreateRevision(l.key), "=", 0)).
        Then(clientv3.OpPut(l.key, l.value, clientv3.WithLease(l.lease))).
        Else(clientv3.OpGet(l.key))

    resp, err := txn.Commit()
    if err != nil {
        return err
    }

    if !resp.Succeeded {
        return errors.New("lock already held")
    }

    // 启动自动续约
    go l.keepAlive()

    return nil
}

// keepAlive 自动续约
func (l *LeaseLock) keepAlive() {
    ch, err := l.client.KeepAlive(context.Background(), l.lease)
    if err != nil {
        return
    }

    for range ch {
        // 续约成功
    }
}

// Unlock 释放锁
func (l *LeaseLock) Unlock(ctx context.Context) error {
    // 使用事务原子性删除
    txn := l.client.Txn(ctx)
    txn.If(clientv3.Compare(clientv3.Value(l.key), "=", l.value)).
        Then(clientv3.OpDelete(l.key)).
        Else()

    _, err := txn.Commit()
    return err
}

// 基于前缀的分布式锁（支持公平锁）
type PrefixLock struct {
    client *clientv3.Client
    prefix string
    key    string
    lease  clientv3.LeaseID
}

// NewPrefixLock 创建前缀锁
func NewPrefixLock(client *clientv3.Client, prefix string, ttl int64) (*PrefixLock, error) {
    lease, err := client.Grant(context.Background(), ttl)
    if err != nil {
        return nil, err
    }

    return &PrefixLock{
        client: client,
        prefix: prefix,
        lease:  lease.ID,
    }, nil
}

// Lock 获取锁（公平锁）
func (l *PrefixLock) Lock(ctx context.Context) error {
    // 创建带序号的key
    key := fmt.Sprintf("%s/%v", l.prefix, l.lease)

    // 创建key
    _, err := l.client.Put(ctx, key, "", clientv3.WithLease(l.lease))
    if err != nil {
        return err
    }

    l.key = key

    // 获取当前所有key
    resp, err := l.client.Get(ctx, l.prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
    if err != nil {
        return err
    }

    // 如果不是第一个，等待
    if len(resp.Kvs) > 0 && string(resp.Kvs[0].Key) != key {
        // 监听前一个key的删除事件
        prevKey := string(resp.Kvs[0].Key)
        for _, kv := range resp.Kvs {
            if string(kv.Key) == key {
                break
            }
            prevKey = string(kv.Key)
        }

        watchCh := l.client.Watch(ctx, prevKey)
        <-watchCh
    }

    // 启动续约
    go l.keepAlive()

    return nil
}

func (l *PrefixLock) keepAlive() {
    ch, err := l.client.KeepAlive(context.Background(), l.lease)
    if err != nil {
        return
    }

    for range ch {
    }
}

// Unlock 释放锁
func (l *PrefixLock) Unlock(ctx context.Context) error {
    _, err := l.client.Delete(ctx, l.key)
    return err
}

func generateUniqueValue() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 使用示例
func ExampleEtcdLock() {
    endpoints := []string{"localhost:2379"}

    // 方式1: 使用concurrency包
    lock, err := NewEtcdLock(endpoints, "/locks/my-resource", 10)
    if err != nil {
        panic(err)
    }
    defer lock.Close()

    ctx := context.Background()

    if err := lock.Lock(ctx); err != nil {
        panic(err)
    }

    fmt.Println("Lock acquired")

    lock.Unlock(ctx)
}
```

---

### 10.3 基于ZooKeeper的锁

#### 概念定义

ZooKeeper通过**临时顺序节点**实现分布式锁，利用**Watcher机制**监听前一个节点的删除事件，实现公平锁。

#### 完整示例：Go实现ZooKeeper分布式锁

```go
// zklock/zk_lock.go
package zklock

import (
    "errors"
    "fmt"
    "path"
    "sort"
    "strings"
    "time"

    "github.com/go-zookeeper/zk"
)

// ZKLock ZooKeeper分布式锁
type ZKLock struct {
    conn   *zk.Conn
    path   string
    node   string
    prefix string
}

// NewZKLock 创建ZK锁
func NewZKLock(servers []string, path string, sessionTimeout time.Duration) (*ZKLock, error) {
    conn, _, err := zk.Connect(servers, sessionTimeout)
    if err != nil {
        return nil, err
    }

    // 创建锁目录
    exists, _, err := conn.Exists(path)
    if err != nil {
        return nil, err
    }

    if !exists {
        _, err = conn.Create(path, []byte{}, 0, zk.WorldACL(zk.PermAll))
        if err != nil && err != zk.ErrNodeExists {
            return nil, err
        }
    }

    return &ZKLock{
        conn:   conn,
        path:   path,
        prefix: "lock-",
    }, nil
}

// Lock 获取锁
func (l *ZKLock) Lock() error {
    // 1. 创建临时顺序节点
    nodePath, err := l.conn.Create(
        path.Join(l.path, l.prefix),
        []byte{},
        zk.FlagEphemeral|zk.FlagSequence,
        zk.WorldACL(zk.PermAll),
    )
    if err != nil {
        return err
    }

    l.node = path.Base(nodePath)

    // 2. 获取所有子节点
    for {
        children, _, err := l.conn.Children(l.path)
        if err != nil {
            return err
        }

        // 3. 排序
        sort.Strings(children)

        // 4. 找到当前节点位置
        idx := -1
        for i, child := range children {
            if child == l.node {
                idx = i
                break
            }
        }

        if idx == -1 {
            return errors.New("node not found")
        }

        // 5. 如果是第一个节点，获取锁成功
        if idx == 0 {
            return nil
        }

        // 6. 监听前一个节点
        prevNode := children[idx-1]
        exists, _, ch, err := l.conn.ExistsW(path.Join(l.path, prevNode))
        if err != nil {
            return err
        }

        // 7. 如果前一个节点已删除，重试
        if !exists {
            continue
        }

        // 8. 等待前一个节点删除
        <-ch
    }
}

// TryLock 尝试获取锁
func (l *ZKLock) TryLock(timeout time.Duration) (bool, error) {
    done := make(chan struct{})
    go func() {
        l.Lock()
        close(done)
    }()

    select {
    case <-done:
        return true, nil
    case <-time.After(timeout):
        return false, nil
    }
}

// Unlock 释放锁
func (l *ZKLock) Unlock() error {
    if l.node == "" {
        return nil
    }

    err := l.conn.Delete(path.Join(l.path, l.node), -1)
    l.node = ""
    return err
}

// Close 关闭连接
func (l *ZKLock) Close() {
    l.conn.Close()
}

// 可重入锁
type ReentrantZKLock struct {
    *ZKLock
    threadID string
    count    int
}

func NewReentrantZKLock(servers []string, path string, threadID string) (*ReentrantZKLock, error) {
    lock, err := NewZKLock(servers, path, 10*time.Second)
    if err != nil {
        return nil, err
    }

    return &ReentrantZKLock{
        ZKLock:   lock,
        threadID: threadID,
    }, nil
}

func (r *ReentrantZKLock) Lock() error {
    // 检查是否已持有锁
    if r.count > 0 {
        r.count++
        return nil
    }

    // 在节点数据中存储线程ID
    nodePath, err := r.conn.Create(
        path.Join(r.path, r.prefix),
        []byte(r.threadID),
        zk.FlagEphemeral|zk.FlagSequence,
        zk.WorldACL(zk.PermAll),
    )
    if err != nil {
        return err
    }

    r.node = path.Base(nodePath)

    // 获取锁逻辑...

    r.count = 1
    return nil
}

func (r *ReentrantZKLock) Unlock() error {
    r.count--
    if r.count > 0 {
        return nil
    }

    return r.ZKLock.Unlock()
}

// 读写锁
type ZKRWLock struct {
    conn     *zk.Conn
    path     string
    readPrefix  string
    writePrefix string
}

func NewZKRWLock(servers []string, path string) (*ZKRWLock, error) {
    conn, _, err := zk.Connect(servers, 10*time.Second)
    if err != nil {
        return nil, err
    }

    return &ZKRWLock{
        conn:        conn,
        path:        path,
        readPrefix:  "read-",
        writePrefix: "write-",
    }, nil
}

// RLock 获取读锁
func (r *ZKRWLock) RLock() error {
    // 创建读节点
    nodePath, err := r.conn.Create(
        path.Join(r.path, r.readPrefix),
        []byte{},
        zk.FlagEphemeral|zk.FlagSequence,
        zk.WorldACL(zk.PermAll),
    )
    if err != nil {
        return err
    }

    node := path.Base(nodePath)

    // 检查是否有写锁
    for {
        children, _, err := r.conn.Children(r.path)
        if err != nil {
            return err
        }

        sort.Strings(children)

        // 找到第一个写锁
        writeIdx := -1
        for i, child := range children {
            if strings.HasPrefix(child, r.writePrefix) {
                writeIdx = i
                break
            }
        }

        // 找到自己的位置
        myIdx := -1
        for i, child := range children {
            if child == node {
                myIdx = i
                break
            }
        }

        // 如果没有写锁或写锁在自己之后，获取读锁成功
        if writeIdx == -1 || writeIdx > myIdx {
            return nil
        }

        // 监听写锁删除
        exists, _, ch, err := r.conn.ExistsW(path.Join(r.path, children[writeIdx]))
        if err != nil {
            return err
        }

        if !exists {
            continue
        }

        <-ch
    }
}

// WLock 获取写锁
func (r *ZKRWLock) WLock() error {
    // 创建写节点
    nodePath, err := r.conn.Create(
        path.Join(r.path, r.writePrefix),
        []byte{},
        zk.FlagEphemeral|zk.FlagSequence,
        zk.WorldACL(zk.PermAll),
    )
    if err != nil {
        return err
    }

    node := path.Base(nodePath)

    // 等待成为第一个节点
    for {
        children, _, err := r.conn.Children(r.path)
        if err != nil {
            return err
        }

        sort.Strings(children)

        if children[0] == node {
            return nil
        }

        // 监听前一个节点
        idx := -1
        for i, child := range children {
            if child == node {
                idx = i
                break
            }
        }

        if idx <= 0 {
            continue
        }

        prevNode := children[idx-1]
        exists, _, ch, err := r.conn.ExistsW(path.Join(r.path, prevNode))
        if err != nil {
            return err
        }

        if !exists {
            continue
        }

        <-ch
    }
}

func (r *ZKRWLock) Unlock(node string) error {
    return r.conn.Delete(path.Join(r.path, node), -1)
}
```

---

### 10.4 锁的续期与释放

#### 完整示例：Go实现锁续期

```go
// lock/renewal.go
package lock

import (
    "context"
    "sync"
    "time"
)

// RenewableLock 可续期锁接口
type RenewableLock interface {
    Lock(ctx context.Context) error
    Unlock(ctx context.Context) error
    Renew(ctx context.Context, ttl time.Duration) error
}

// LockManager 锁管理器（自动续期）
type LockManager struct {
    lock      RenewableLock
    ttl       time.Duration
    renewInterval time.Duration
    stopChan  chan struct{}
    mu        sync.Mutex
    isLocked  bool
}

func NewLockManager(lock RenewableLock, ttl, renewInterval time.Duration) *LockManager {
    return &LockManager{
        lock:          lock,
        ttl:           ttl,
        renewInterval: renewInterval,
        stopChan:      make(chan struct{}),
    }
}

// Lock 获取锁并启动自动续期
func (m *LockManager) Lock(ctx context.Context) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if err := m.lock.Lock(ctx); err != nil {
        return err
    }

    m.isLocked = true

    // 启动自动续期
    go m.autoRenew()

    return nil
}

// Unlock 释放锁并停止续期
func (m *LockManager) Unlock(ctx context.Context) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if !m.isLocked {
        return nil
    }

    // 停止续期
    close(m.stopChan)

    m.isLocked = false

    return m.lock.Unlock(ctx)
}

// autoRenew 自动续期
func (m *LockManager) autoRenew() {
    ticker := time.NewTicker(m.renewInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            err := m.lock.Renew(ctx, m.ttl)
            cancel()

            if err != nil {
                // 续期失败，记录日志
                return
            }

        case <-m.stopChan:
            return
        }
    }
}

// 带超时的锁
type TimeoutLock struct {
    lock     RenewableLock
    maxHoldTime time.Duration
}

func (t *TimeoutLock) Lock(ctx context.Context) error {
    if err := t.lock.Lock(ctx); err != nil {
        return err
    }

    // 启动超时释放
    go func() {
        <-time.After(t.maxHoldTime)
        t.lock.Unlock(context.Background())
    }()

    return nil
}

// 锁状态监控
type LockMonitor struct {
    locks map[string]LockInfo
    mu    sync.RWMutex
}

type LockInfo struct {
    Resource   string
    Holder     string
    AcquiredAt time.Time
    ExpiresAt  time.Time
}

func (m *LockMonitor) RecordLock(info LockInfo) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.locks[info.Resource] = info
}

func (m *LockMonitor) RecordUnlock(resource string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.locks, resource)
}

func (m *LockMonitor) GetActiveLocks() []LockInfo {
    m.mu.RLock()
    defer m.mu.RUnlock()

    locks := make([]LockInfo, 0, len(m.locks))
    for _, info := range m.locks {
        locks = append(locks, info)
    }
    return locks
}
```

---

## 11. 消息队列

### 11.1 Kafka集成

#### 概念定义

Kafka是分布式流处理平台，基于**发布-订阅模型**，支持**高吞吐**、**持久化**、**分区**等特性。

#### 架构图

```
Kafka架构:

┌─────────────────────────────────────────────────────────────┐
│                      Kafka集群                               │
│                                                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐                     │
│  │ Broker  │  │ Broker  │  │ Broker  │                     │
│  │  :9092  │  │  :9092  │  │  :9092  │                     │
│  │         │  │         │  │         │                     │
│  │ TopicA  │  │ TopicA  │  │ TopicA  │   Partition 0       │
│  │   P0    │  │   P0    │  │   P0    │   (Leader)          │
│  │ TopicA  │  │ TopicA  │  │ TopicA  │                     │
│  │   P1    │  │   P1    │  │   P1    │   Partition 1       │
│  │         │  │         │  │         │   (Leader)          │
│  └─────────┘  └─────────┘  └─────────┘                     │
│       ▲            ▲            ▲                          │
│       │            │            │                          │
│       └────────────┴────────────┘                          │
│              ZooKeeper/KRaft                                │
└─────────────────────────────────────────────────────────────┘
       ▲                                    │
       │                                    │
┌──────┴──────┐                    ┌───────┴───────┐
│  Producer   │                    │   Consumer    │
│  (生产者)    │                    │   (消费者)     │
│             │                    │               │
│ 发送消息     │                    │ 消费消息       │
│ 指定Topic   │                    │ 指定Group     │
└─────────────┘                    └───────────────┘

消息可靠性:
- acks=0: 不等待确认
- acks=1: 等待Leader确认
- acks=all: 等待所有ISR确认
```

#### 完整示例：Go集成Kafka

```go
// kafka/kafka_client.go
package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/IBM/sarama"
)

// Producer Kafka生产者
type Producer struct {
    producer sarama.SyncProducer
    topic    string
}

// NewProducer 创建生产者
func NewProducer(brokers []string, topic string) (*Producer, error) {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll  // 等待所有副本确认
    config.Producer.Retry.Max = 3                      // 重试次数
    config.Producer.Return.Successes = true            // 返回成功
    config.Producer.Return.Errors = true               // 返回错误

    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }

    return &Producer{
        producer: producer,
        topic:    topic,
    }, nil
}

// SendMessage 发送消息
func (p *Producer) SendMessage(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    msg := &sarama.ProducerMessage{
        Topic: p.topic,
        Key:   sarama.StringEncoder(key),
        Value: sarama.ByteEncoder(data),
    }

    partition, offset, err := p.producer.SendMessage(msg)
    if err != nil {
        return err
    }

    log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
    return nil
}

// SendMessages 批量发送
func (p *Producer) SendMessages(ctx context.Context, messages []Message) error {
    var msgs []*sarama.ProducerMessage

    for _, m := range messages {
        data, _ := json.Marshal(m.Value)
        msgs = append(msgs, &sarama.ProducerMessage{
            Topic: p.topic,
            Key:   sarama.StringEncoder(m.Key),
            Value: sarama.ByteEncoder(data),
        })
    }

    return p.producer.SendMessages(msgs)
}

// Close 关闭生产者
func (p *Producer) Close() error {
    return p.producer.Close()
}

// Message 消息结构
type Message struct {
    Key   string
    Value interface{}
}

// Consumer Kafka消费者
type Consumer struct {
    consumer sarama.ConsumerGroup
    handler  sarama.ConsumerGroupHandler
    topics   []string
    groupID  string
}

// ConsumerHandler 消费者处理器
type ConsumerHandler struct {
    ready   chan bool
    handler func(message *sarama.ConsumerMessage) error
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
    close(h.ready)
    return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
    return nil
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for message := range claim.Messages() {
        log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s",
            string(message.Value), message.Timestamp, message.Topic)

        if err := h.handler(message); err != nil {
            log.Printf("Handle message error: %v", err)
            // 处理失败，不提交offset，稍后重试
            continue
        }

        // 处理成功，提交offset
        session.MarkMessage(message, "")
    }

    return nil
}

// NewConsumer 创建消费者
func NewConsumer(brokers []string, groupID string, topics []string, handler func(message *sarama.ConsumerMessage) error) (*Consumer, error) {
    config := sarama.NewConfig()
    config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    config.Consumer.Return.Errors = true

    consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
    if err != nil {
        return nil, err
    }

    return &Consumer{
        consumer: consumer,
        topics:   topics,
        groupID:  groupID,
        handler: &ConsumerHandler{
            ready:   make(chan bool),
            handler: handler,
        },
    }, nil
}

// Start 启动消费
func (c *Consumer) Start(ctx context.Context) error {
    for {
        if err := c.consumer.Consume(ctx, c.topics, c.handler); err != nil {
            log.Printf("Error from consumer: %v", err)
        }

        if ctx.Err() != nil {
            return ctx.Err()
        }
    }
}

// Close 关闭消费者
func (c *Consumer) Close() error {
    return c.consumer.Close()
}

// 事务生产者
type TransactionalProducer struct {
    producer sarama.AsyncProducer
}

func NewTransactionalProducer(brokers []string, transactionalID string) (*TransactionalProducer, error) {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Idempotent = true  // 幂等生产者
    config.Producer.Transaction.ID = transactionalID

    producer, err := sarama.NewAsyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }

    return &TransactionalProducer{producer: producer}, nil
}

func (p *TransactionalProducer) SendInTransaction(ctx context.Context, messages []Message, consumerOffsets map[string]int64) error {
    // 初始化事务
    if err := p.producer.BeginTxn(); err != nil {
        return err
    }

    // 发送消息
    for _, msg := range messages {
        data, _ := json.Marshal(msg.Value)
        p.producer.Input() <- &sarama.ProducerMessage{
            Topic: msg.Topic,
            Key:   sarama.StringEncoder(msg.Key),
            Value: sarama.ByteEncoder(data),
        }
    }

    // 提交消费位点
    for topicPartition, offset := range consumerOffsets {
        parts := strings.Split(topicPartition, "-")
        topic := parts[0]
        partition, _ := strconv.Atoi(parts[1])

        metadata := fmt.Sprintf("%d", offset)
        if err := p.producer.AddOffsetsToTxn(map[string][]*sarama.PartitionOffsetMetadata{
            topic: {{
                Partition: int32(partition),
                Offset:    offset,
                Metadata:  &metadata,
            }},
        }, groupID); err != nil {
            p.producer.AbortTxn()
            return err
        }
    }

    // 提交事务
    return p.producer.CommitTxn()
}

// 使用示例
func ExampleKafka() {
    brokers := []string{"localhost:9092"}
    topic := "orders"

    // 生产者
    producer, err := NewProducer(brokers, topic)
    if err != nil {
        panic(err)
    }
    defer producer.Close()

    // 发送消息
    order := map[string]interface{}{
        "order_id": "12345",
        "user_id":  "user001",
        "amount":   100.0,
    }

    ctx := context.Background()
    if err := producer.SendMessage(ctx, "order-12345", order); err != nil {
        panic(err)
    }

    // 消费者
    handler := func(msg *sarama.ConsumerMessage) error {
        fmt.Printf("Received: %s\n", string(msg.Value))
        return nil
    }

    consumer, err := NewConsumer(brokers, "order-consumer", []string{topic}, handler)
    if err != nil {
        panic(err)
    }
    defer consumer.Close()

    consumer.Start(ctx)
}
```

#### 反例说明

```go
// ❌ 错误：不处理消费失败
for msg := range claim.Messages() {
    process(msg)  // 可能panic
    session.MarkMessage(msg, "")  // 总是提交offset
}
// 问题：处理失败的消息也被标记为已消费

// ❌ 错误：不设置消费者组
config.Consumer.Group.Rebalance.Strategy = nil
// 问题：无法负载均衡，多个消费者消费相同消息

// ❌ 错误：自动提交offset
config.Consumer.Offsets.AutoCommit.Enable = true
// 问题：消息可能未处理完就提交offset，丢失消息
```

---

### 11.2 RabbitMQ集成

#### 概念定义

RabbitMQ是基于AMQP协议的消息队列，支持**路由**、**主题**、**扇出**等多种交换机类型。

#### 完整示例：Go集成RabbitMQ

```go
// rabbitmq/rabbitmq_client.go
package rabbitmq

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    amqp "github.com/rabbitmq/amqp091-go"
)

// Client RabbitMQ客户端
type Client struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

// NewClient 创建客户端
func NewClient(url string) (*Client, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }

    return &Client{
        conn:    conn,
        channel: ch,
    }, nil
}

// Close 关闭连接
func (c *Client) Close() {
    c.channel.Close()
    c.conn.Close()
}

// DeclareQueue 声明队列
func (c *Client) DeclareQueue(name string, durable, autoDelete, exclusive bool) (amqp.Queue, error) {
    return c.channel.QueueDeclare(
        name,
        durable,    // 持久化
        autoDelete, // 自动删除
        exclusive,  // 独占
        false,
        nil,
    )
}

// DeclareExchange 声明交换机
func (c *Client) DeclareExchange(name, kind string) error {
    return c.channel.ExchangeDeclare(
        name,
        kind,    // direct, topic, fanout, headers
        true,    // 持久化
        false,
        false,
        false,
        nil,
    )
}

// BindQueue 绑定队列
func (c *Client) BindQueue(queue, exchange, routingKey string) error {
    return c.channel.QueueBind(
        queue,
        routingKey,
        exchange,
        false,
        nil,
    )
}

// Producer 生产者
type Producer struct {
    client *Client
}

func NewProducer(client *Client) *Producer {
    return &Producer{client: client}
}

// Publish 发布消息
func (p *Producer) Publish(ctx context.Context, exchange, routingKey string, body interface{}) error {
    data, err := json.Marshal(body)
    if err != nil {
        return err
    }

    return p.client.channel.PublishWithContext(
        ctx,
        exchange,
        routingKey,
        false, // mandatory
        false, // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         data,
            DeliveryMode: amqp.Persistent, // 持久化消息
        },
    )
}

// PublishWithDelay 延迟消息
func (p *Producer) PublishWithDelay(ctx context.Context, exchange, routingKey string, body interface{}, delayMs int) error {
    data, err := json.Marshal(body)
    if err != nil {
        return err
    }

    headers := amqp.Table{
        "x-delay": delayMs,
    }

    return p.client.channel.PublishWithContext(
        ctx,
        exchange,
        routingKey,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        data,
            Headers:     headers,
        },
    )
}

// Consumer 消费者
type Consumer struct {
    client  *Client
    queue   string
    handler func(delivery amqp.Delivery) error
}

func NewConsumer(client *Client, queue string, handler func(delivery amqp.Delivery) error) *Consumer {
    return &Consumer{
        client:  client,
        queue:   queue,
        handler: handler,
    }
}

// Start 开始消费
func (c *Consumer) Start() error {
    msgs, err := c.client.channel.Consume(
        c.queue,
        "",    // consumer tag
        false, // auto-ack (手动确认)
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    go func() {
        for d := range msgs {
            if err := c.handler(d); err != nil {
                log.Printf("Handle message error: %v", err)
                // 拒绝消息，重新入队
                d.Nack(false, true)
            } else {
                // 确认消息
                d.Ack(false)
            }
        }
    }()

    return nil
}

// 使用示例
func ExampleRabbitMQ() {
    url := "amqp://guest:guest@localhost:5672/"

    client, err := NewClient(url)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // 声明交换机
    client.DeclareExchange("orders", "topic")

    // 声明队列
    queue, _ := client.DeclareQueue("order_queue", true, false, false)

    // 绑定队列
    client.BindQueue(queue.Name, "orders", "order.created")

    // 生产者
    producer := NewProducer(client)

    order := map[string]interface{}{
        "order_id": "12345",
        "user_id":  "user001",
    }

    ctx := context.Background()
    producer.Publish(ctx, "orders", "order.created", order)

    // 消费者
    handler := func(d amqp.Delivery) error {
        fmt.Printf("Received: %s\n", string(d.Body))
        return nil
    }

    consumer := NewConsumer(client, queue.Name, handler)
    consumer.Start()

    select {}
}

// 延迟队列实现
func SetupDelayQueue(client *Client) error {
    // 声明延迟交换机
    args := amqp.Table{
        "x-delayed-type": "direct",
    }

    err := client.channel.ExchangeDeclare(
        "delayed_exchange",
        "x-delayed-message",
        true,
        false,
        false,
        false,
        args,
    )
    if err != nil {
        return err
    }

    // 声明队列
    queue, err := client.DeclareQueue("delayed_queue", true, false, false)
    if err != nil {
        return err
    }

    // 绑定
    return client.BindQueue(queue.Name, "delayed_exchange", "")
}

// 死信队列实现
func SetupDLQ(client *Client) error {
    // 声明死信交换机
    client.DeclareExchange("dlx_exchange", "direct")

    // 声明死信队列
    dlq, _ := client.DeclareQueue("dlq", true, false, false)
    client.BindQueue(dlq.Name, "dlx_exchange", "dlq_routing_key")

    // 声明主队列（设置死信参数）
    args := amqp.Table{
        "x-dead-letter-exchange":    "dlx_exchange",
        "x-dead-letter-routing-key": "dlq_routing_key",
        "x-message-ttl":             30000, // 30秒过期
        "x-max-retries":             3,
    }

    _, err := client.channel.QueueDeclare(
        "main_queue",
        true,
        false,
        false,
        false,
        args,
    )

    return err
}
```

---

### 11.3 NATS集成

#### 概念定义

NATS是高性能、轻量级的消息系统，支持**发布-订阅**、**请求-回复**、**队列组**等模式。

#### 完整示例：Go集成NATS

```go
// nats/nats_client.go
package nats

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/nats-io/nats.go"
)

// Client NATS客户端
type Client struct {
    conn *nats.Conn
    js   nats.JetStreamContext
}

// NewClient 创建客户端
func NewClient(url string) (*Client, error) {
    conn, err := nats.Connect(url)
    if err != nil {
        return nil, err
    }

    js, err := conn.JetStream()
    if err != nil {
        conn.Close()
        return nil, err
    }

    return &Client{
        conn: conn,
        js:   js,
    }, nil
}

// Close 关闭连接
func (c *Client) Close() {
    c.conn.Close()
}

// Publisher 发布者
type Publisher struct {
    client *Client
}

func NewPublisher(client *Client) *Publisher {
    return &Publisher{client: client}
}

// Publish 发布消息
func (p *Publisher) Publish(subject string, data interface{}) error {
    payload, err := json.Marshal(data)
    if err != nil {
        return err
    }

    return p.client.conn.Publish(subject, payload)
}

// PublishWithContext 带上下文的发布
func (p *Publisher) PublishWithContext(ctx context.Context, subject string, data interface{}) error {
    payload, err := json.Marshal(data)
    if err != nil {
        return err
    }

    return p.client.conn.PublishMsg(&nats.Msg{
        Subject: subject,
        Data:    payload,
    })
}

// Request 请求-回复
func (p *Publisher) Request(subject string, data interface{}, timeout time.Duration) (*nats.Msg, error) {
    payload, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    return p.client.conn.Request(subject, payload, timeout)
}

// Subscriber 订阅者
type Subscriber struct {
    client *Client
    subs   []*nats.Subscription
}

func NewSubscriber(client *Client) *Subscriber {
    return &Subscriber{
        client: client,
        subs:   make([]*nats.Subscription, 0),
    }
}

// Subscribe 订阅
func (s *Subscriber) Subscribe(subject string, handler func(msg *nats.Msg)) error {
    sub, err := s.client.conn.Subscribe(subject, handler)
    if err != nil {
        return err
    }

    s.subs = append(s.subs, sub)
    return nil
}

// SubscribeQueue 队列组订阅
func (s *Subscriber) SubscribeQueue(subject, queue string, handler func(msg *nats.Msg)) error {
    sub, err := s.client.conn.QueueSubscribe(subject, queue, handler)
    if err != nil {
        return err
    }

    s.subs = append(s.subs, sub)
    return nil
}

// Unsubscribe 取消订阅
func (s *Subscriber) Unsubscribe() error {
    for _, sub := range s.subs {
        if err := sub.Unsubscribe(); err != nil {
            return err
        }
    }
    return nil
}

// JetStream发布者
type JetStreamPublisher struct {
    js nats.JetStreamContext
}

func NewJetStreamPublisher(client *Client) *JetStreamPublisher {
    return &JetStreamPublisher{js: client.js}
}

// Publish 发布持久化消息
func (p *JetStreamPublisher) Publish(subject string, data interface{}) (*nats.PubAck, error) {
    payload, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    return p.js.Publish(subject, payload)
}

// CreateStream 创建流
func (p *JetStreamPublisher) CreateStream(config *nats.StreamConfig) (*nats.StreamInfo, error) {
    return p.js.AddStream(config)
}

// JetStream消费者
type JetStreamConsumer struct {
    js nats.JetStreamContext
}

func NewJetStreamConsumer(client *Client) *JetStreamConsumer {
    return &JetStreamConsumer{js: client.js}
}

// CreateConsumer 创建消费者
func (c *JetStreamConsumer) CreateConsumer(stream string, config *nats.ConsumerConfig) (*nats.ConsumerInfo, error) {
    return c.js.AddConsumer(stream, config)
}

// Subscribe 订阅
func (c *JetStreamConsumer) Subscribe(subject string, handler func(msg *nats.Msg)) (nats.Subscription, error) {
    return c.js.Subscribe(subject, handler)
}

// PullSubscribe Pull订阅
func (c *JetStreamConsumer) PullSubscribe(subject, durable string) (nats.Subscription, error) {
    return c.js.PullSubscribe(subject, durable)
}

// 使用示例
func ExampleNATS() {
    url := "nats://localhost:4222"

    client, err := NewClient(url)
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // 发布-订阅
    publisher := NewPublisher(client)
    subscriber := NewSubscriber(client)

    // 订阅
    subscriber.Subscribe("orders.created", func(msg *nats.Msg) {
        fmt.Printf("Received: %s\n", string(msg.Data))
    })

    // 发布
    order := map[string]interface{}{
        "order_id": "12345",
        "amount":   100.0,
    }
    publisher.Publish("orders.created", order)

    // 队列组（负载均衡）
    subscriber.SubscribeQueue("orders.created", "order-processors", func(msg *nats.Msg) {
        fmt.Printf("Queue consumer received: %s\n", string(msg.Data))
    })

    // 请求-回复
    subscriber.Subscribe("orders.get", func(msg *nats.Msg) {
        // 处理请求
        response := map[string]string{"status": "ok"}
        data, _ := json.Marshal(response)
        msg.Respond(data)
    })

    resp, err := publisher.Request("orders.get", map[string]string{"order_id": "12345"}, 5*time.Second)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Response: %s\n", string(resp.Data))

    // JetStream（持久化）
    jsPub := NewJetStreamPublisher(client)
    jsSub := NewJetStreamConsumer(client)

    // 创建流
    jsPub.CreateStream(&nats.StreamConfig{
        Name:     "ORDERS",
        Subjects: []string{"orders.*"},
        Retention: nats.WorkQueuePolicy,
    })

    // 发布持久化消息
    jsPub.Publish("orders.processed", order)

    // 订阅持久化消息
    jsSub.Subscribe("orders.processed", func(msg *nats.Msg) {
        fmt.Printf("JetStream received: %s\n", string(msg.Data))
        msg.Ack()
    })

    select {}
}
```

---

### 11.4 消息可靠性保证

#### 完整示例：Go实现消息可靠性

```go
// messaging/reliability.go
package messaging

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
)

// ReliableProducer 可靠生产者
type ReliableProducer struct {
    producer  Producer
    outbox    OutboxStore
    publisher MessagePublisher
}

// OutboxStore 发件箱存储
type OutboxStore interface {
    Save(ctx context.Context, msg *OutboxMessage) error
    MarkSent(ctx context.Context, id string) error
    GetUnsent(ctx context.Context, limit int) ([]*OutboxMessage, error)
}

// OutboxMessage 发件箱消息
type OutboxMessage struct {
    ID        string
    Topic     string
    Key       string
    Payload   []byte
    Status    string // pending, sent
    CreatedAt time.Time
}

// Producer 生产者接口
type Producer interface {
    Send(ctx context.Context, topic string, key string, payload []byte) error
}

// MessagePublisher 消息发布接口
type MessagePublisher interface {
    Publish(ctx context.Context, topic string, key string, payload []byte) error
}

// Send 可靠发送
func (p *ReliableProducer) Send(ctx context.Context, topic, key string, payload interface{}) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    msg := &OutboxMessage{
        ID:        generateID(),
        Topic:     topic,
        Key:       key,
        Payload:   data,
        Status:    "pending",
        CreatedAt: time.Now(),
    }

    // 1. 保存到发件箱（在同一事务中）
    if err := p.outbox.Save(ctx, msg); err != nil {
        return err
    }

    // 2. 发送消息
    if err := p.publisher.Publish(ctx, topic, key, data); err != nil {
        return err
    }

    // 3. 标记为已发送
    return p.outbox.MarkSent(ctx, msg.ID)
}

// StartRelay 启动中继服务
func (p *ReliableProducer) StartRelay(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            p.relay(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (p *ReliableProducer) relay(ctx context.Context) {
    messages, err := p.outbox.GetUnsent(ctx, 100)
    if err != nil {
        return
    }

    for _, msg := range messages {
        err := p.publisher.Publish(ctx, msg.Topic, msg.Key, msg.Payload)
        if err != nil {
            continue
        }

        p.outbox.MarkSent(ctx, msg.ID)
    }
}

// ReliableConsumer 可靠消费者
type ReliableConsumer struct {
    consumer    Consumer
    idempotency IdempotencyStore
    handler     func(msg interface{}) error
}

// IdempotencyStore 幂等存储
type IdempotencyStore interface {
    IsProcessed(ctx context.Context, msgID string) (bool, error)
    MarkProcessed(ctx context.Context, msgID string) error
}

// Consume 可靠消费
func (c *ReliableConsumer) Consume(ctx context.Context, msgID string, payload []byte) error {
    // 1. 幂等检查
    processed, err := c.idempotency.IsProcessed(ctx, msgID)
    if err != nil {
        return err
    }

    if processed {
        return nil // 已处理，直接返回
    }

    // 2. 解析消息
    var msg interface{}
    if err := json.Unmarshal(payload, &msg); err != nil {
        return err
    }

    // 3. 处理消息
    if err := c.handler(msg); err != nil {
        return err
    }

    // 4. 标记为已处理
    return c.idempotency.MarkProcessed(ctx, msgID)
}

// DeadLetterQueue 死信队列
type DeadLetterQueue struct {
    producer Producer
    maxRetries int
}

func (d *DeadLetterQueue) Send(ctx context.Context, originalTopic string, msg interface{}, err error) error {
    dlqMsg := map[string]interface{}{
        "original_topic": originalTopic,
        "message":        msg,
        "error":          err.Error(),
        "timestamp":      time.Now(),
    }

    data, _ := json.Marshal(dlqMsg)
    return d.producer.Send(ctx, "DLQ", "", data)
}

// MessageTracing 消息追踪
type MessageTracing struct {
    tracer Tracer
}

type Tracer interface {
    StartSpan(ctx context.Context, operation string) (context.Context, Span)
}

type Span interface {
    SetTag(key string, value interface{})
    Finish()
}

func (m *MessageTracing) TracePublish(ctx context.Context, topic string, msg interface{}) context.Context {
    ctx, span := m.tracer.StartSpan(ctx, "publish")
    span.SetTag("topic", topic)
    defer span.Finish()

    return ctx
}

func (m *MessageTracing) TraceConsume(ctx context.Context, topic string, msg interface{}) context.Context {
    ctx, span := m.tracer.StartSpan(ctx, "consume")
    span.SetTag("topic", topic)
    defer span.Finish()

    return ctx
}

func generateID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
```

---

## 总结

本文档全面梳理了Go语言在分布式系统中的11个核心设计模型：

1. **微服务架构基础**：服务拆分原则、通信方式、API网关、BFF模式
2. **服务发现与注册**：客户端/服务端发现、Consul、etcd、Kubernetes
3. **负载均衡**：轮询、随机、加权、一致性哈希、健康检查
4. **熔断与降级**：Circuit Breaker、hystrix-go、sentinel-go
5. **限流与配额**：令牌桶、漏桶、滑动窗口、分布式限流
6. **重试与退避**：指数退避、抖动、重试策略、幂等性
7. **分布式事务**：2PC、3PC、TCC、Saga、本地消息表
8. **一致性协议**：Raft、领导者选举、日志复制、成员变更
9. **分布式缓存**：穿透、击穿、雪崩、一致性、Cache-Aside
10. **分布式锁**：RedLock、etcd锁、ZooKeeper锁
11. **消息队列**：Kafka、RabbitMQ、NATS集成

每个模型都包含概念定义、架构图、完整示例、反例说明和注意事项，为Go语言分布式系统开发提供全面的参考。

---

*文档生成时间: 2024年*
*作者: AI Assistant*
