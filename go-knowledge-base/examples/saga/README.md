# Saga 分布式事务示例

三服务分布式事务实现，展示 Saga 编排模式与补偿事务。

## 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Saga Transaction Flow                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  成功流程:                                                                   │
│                                                                              │
│  ┌─────────┐    Create Order    ┌───────────┐   Reserve    ┌───────────┐   │
│  │  Client │───────────────────►│   Order   │─────────────►│ Inventory │   │
│  └─────────┘                    │  Service  │              │  Service  │   │
│                                 │(Orchestrator)            └──────┬────┘   │
│                                 └──────┬───┘                       │        │
│                                        │                           │        │
│                                        │ Process                   │        │
│                                        │ Payment                   │        │
│                                        ▼                           │        │
│                                 ┌───────────┐                      │        │
│                                 │  Payment  │                      │        │
│                                 │  Service  │◄─────────────────────┘        │
│                                 └──────┬────┘   Confirm                     │
│                                        │                                    │
│                                        │ Complete                           │
│                                        ▼                                    │
│                                 ┌───────────┐                               │
│                                 │  Order    │                               │
│                                 │ Completed │                               │
│                                 └───────────┘                               │
│                                                                              │
│  失败补偿流程:                                                               │
│                                                                              │
│  ┌─────────┐    Create Order    ┌───────────┐   Reserve    ┌───────────┐   │
│  │  Client │───────────────────►│   Order   │─────────────►│ Inventory │   │
│  └─────────┘                    │  Service  │─────────────►│  Service  │   │
│                                 │           │              └──────┬────┘   │
│                                 └──────┬───┘                       │        │
│                                        │                           │        │
│                                        │ Process                   │        │
│                                        │ Payment                   │        │
│                   ┌────────────────────┼────────────────┐          │        │
│                   │       FAILED       │                │          │        │
│                   ▼                    │                │          │        │
│            ┌───────────┐               │                │          │        │
│            │ Refund    │◄──────────────┘                │          │        │
│            │ Payment   │         Compensate             │          │        │
│            └───────────┘                                │          │        │
│                   │                                     │          │        │
│                   │          Release Inventory          │          │        │
│                   └─────────────────────────────────────┴──────────┘        │
│                                                                              │
│  补偿顺序: 与执行顺序相反 (LIFO)                                              │
│  1. 退款 (Payment Compensation)                                              │
│  2. 释放库存 (Inventory Compensation)                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 服务说明

| 服务 | 职责 | 端口 | 补偿操作 |
|------|------|------|---------|
| Order Service | 编排器，协调 Saga | :8080 | 取消订单 |
| Payment Service | 支付处理 | :8081 | 退款 |
| Inventory Service | 库存管理 | :8082 | 释放库存 |

## 运行

```bash
# 启动所有服务
docker-compose up

# 或使用 air 热重载
cd cmd/order && air
cd cmd/payment && air
cd cmd/inventory && air
```

## 测试 Saga

```bash
# 成功场景
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "items": [
      {"product_id": "prod-1", "quantity": 2, "price": 50.00}
    ]
  }'

# 失败场景 (库存不足)
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-456",
    "items": [
      {"product_id": "prod-999", "quantity": 1000, "price": 10.00}
    ]
  }'
# 预期: 订单创建 → 库存预留失败 → 事务回滚
```

## 状态机

```
PENDING ──► RESERVED ──► PAID ──► COMPLETED
    │           │          │
    ▼           ▼          ▼
CANCELLED   RELEASED   REFUNDED

- PENDING: 初始状态
- RESERVED: 库存已预留
- PAID: 支付成功
- COMPLETED: 事务完成
- CANCELLED: 订单取消
- RELEASED: 库存释放 (补偿)
- REFUNDED: 退款完成 (补偿)
```

## 关键设计

1. **幂等性**: 每个补偿操作必须幂等
2. **持久化**: Saga 状态持久化到数据库
3. **重试**: 补偿失败时重试，最终一致性
4. **监控**: 记录 Saga 执行轨迹，便于审计
