# Saga 分布式事务示例

电商订单 Saga 实现：订单 + 支付 + 库存

## 服务

- order-service: 订单服务
- payment-service: 支付服务
- inventory-service: 库存服务
- orchestrator: Saga 编排器

## 流程

1. 创建订单
2. 扣减库存
3. 处理支付
4. 完成订单

失败时按相反顺序补偿。

## 运行

```bash
docker-compose up -d
make run
```
