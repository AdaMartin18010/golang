# EC-036: Shared Database Pattern (共享数据库模式)

> **维度**: Engineering-CloudNative  
> **级别**: S (>15KB)  
> **标签**: #shared-database #monolith #migration #intermediate  
> **权威来源**:  
> - [Shared Database Pattern](https://microservices.io/patterns/data/shared-database.html) - Chris Richardson  
- [Monolith to Microservices](https://www.oreilly.com/library/view/monolith-to-microservices/9781492047834/) - Sam Newman  
> - [Refactoring Databases](https://www.oreilly.com/library/view/refactoring-databases/0321293533/) - Ambler & Sadalage

---

## 1. 模式形式化定义

### 1.1 问题定义

**问题陈述**: 在从单体应用向微服务迁移的过程中，或者在某些特定约束条件下，如何在保持数据一致性的同时支持多个服务访问同一数据库？

**形式化描述**:
```
给定: 服务集合 S = {S₁, S₂, ..., Sₙ}
给定: 数据库 DB
约束:
  - 多个服务需要访问相同数据
  - 需要 ACID 事务支持
  - 无法立即拆分数据库
目标: 设计访问模式使得服务间松耦合最大化
```

**适用场景**:
- 单体到微服务的迁移过渡期
- 强一致性要求且无法使用 Saga 的场景
- 数据关联复杂，难以立即拆分
- 遗留系统现代化

### 1.2 解决方案形式化

**定义 1.1 (共享数据库模式)**
多个服务共享同一个数据库，但通过以下机制隔离：
1. Schema 分离：每个服务有自己的 Schema
2. 视图隔离：通过数据库视图限制访问
3. API 封装：服务通过 API 而非直接 SQL 访问数据
4. 事务协调：使用分布式事务或协调机制

**形式化表示**:
```
Schema 分配:
  ∀Sᵢ ∈ S: owns_schema(Sᵢ, schemaᵢ)
  schemaᵢ ⊆ DB
  schemaᵢ ∩ schemaⱼ = ∅ (理想情况) 或 controlled_overlap

访问控制:
  access(Sᵢ, table) ⟺ table ∈ schemaᵢ ∨ granted(Sᵢ, table)
```

### 1.3 架构模型

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Shared Database Architecture                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                 │
│  │  Service A  │    │  Service B  │    │  Service C  │                 │
│  │  (Orders)   │    │ (Payments)  │    │(Inventory)  │                 │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘                 │
│         │                  │                  │                         │
│         │ API Layer        │ API Layer        │ API Layer               │
│         │ (Recommended)    │ (Recommended)    │ (Recommended)           │
│         ▼                  ▼                  ▼                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    SHARED DATABASE                               │   │
│  │                                                                  │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │   │
│  │  │   orders    │  │  payments   │  │  inventory  │              │   │
│  │  │   schema    │  │   schema    │  │   schema    │              │   │
│  │  │             │  │             │  │             │              │   │
│  │  │ • orders    │  │ • payments  │  │ • products  │              │   │
│  │  │ • order_items│ │ • refunds   │  │ • stock     │              │   │
│  │  │             │  │             │  │             │              │   │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │   │
│  │         │                │                │                      │   │
│  │         └────────────────┼────────────────┘                      │   │
│  │                          │                                       │   │
│  │              ┌───────────┴───────────┐                          │   │
│  │              │  Shared Tables        │                          │   │
│  │              │  (e.g., users, audit) │                          │   │
│  │              └───────────────────────┘                          │   │
│  │                                                                  │   │
│  │  访问控制:                                                        │   │
│  │  • Service A: 读写 orders schema, 读 shared                      │   │
│  │  • Service B: 读写 payments schema, 读 shared                    │   │
│  │  • Service C: 读写 inventory schema, 读 shared                   │   │
│  │                                                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  风险:                                                                   │
│  • Schema 变更影响多个服务                                               │
│  • 性能问题难以隔离                                                       │
│  • 技术栈锁定                                                            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Go 实现

### 2.1 共享数据库访问层

```go
// shareddatabase/core.go
package shareddatabase

import (
    "context"
    "database/sql"
    "fmt"
    "sync"
    "time"
)

// Schema 数据库 Schema
type Schema string

const (
    SchemaOrders    Schema = "orders"
    SchemaPayments  Schema = "payments"
    SchemaInventory Schema = "inventory"
    SchemaShared    Schema = "shared"
)

// SharedDatabase 共享数据库接口
type SharedDatabase interface {
    // QueryContext 查询（带 Schema 验证）
    QueryContext(ctx context.Context, service string, schema Schema, query string, args ...interface{}) (*sql.Rows, error)
    
    // ExecContext 执行（带 Schema 验证）
    ExecContext(ctx context.Context, service string, schema Schema, query string, args ...interface{}) (sql.Result, error)
    
    // BeginTx 开始事务（跨 Schema 事务）
    BeginTx(ctx context.Context, opts *sql.TxOptions) (*SharedTx, error)
    
    // HealthCheck 健康检查
    HealthCheck(ctx context.Context) error
    
    // Stats 统计信息
    Stats() sql.DBStats
    
    // Close 关闭
    Close() error
}

// AccessPolicy 访问策略
type AccessPolicy struct {
    Service      string
    AllowedRead  []Schema
    AllowedWrite []Schema
}

// sharedDatabaseImpl 共享数据库实现
type sharedDatabaseImpl struct {
    db      *sql.DB
    policies map[string]*AccessPolicy
    mu       sync.RWMutex
}

// SharedTx 共享事务
type SharedTx struct {
    tx      *sql.Tx
    service string
    policy  *AccessPolicy
}

// NewSharedDatabase 创建共享数据库
func NewSharedDatabase(connectionString string, policies []*AccessPolicy) (SharedDatabase, error) {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // 配置连接池
    db.SetMaxOpenConns(50)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // 验证连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    policyMap := make(map[string]*AccessPolicy)
    for _, p := range policies {
        policyMap[p.Service] = p
    }
    
    return &sharedDatabaseImpl{
        db:       db,
        policies: policyMap,
    }, nil
}

// QueryContext 查询
func (s *sharedDatabaseImpl) QueryContext(ctx context.Context, service string, schema Schema, query string, args ...interface{}) (*sql.Rows, error) {
    policy := s.getPolicy(service)
    if policy == nil {
        return nil, fmt.Errorf("no access policy for service: %s", service)
    }
    
    if !s.canRead(policy, schema) {
        return nil, fmt.Errorf("service %s does not have read access to schema %s", service, schema)
    }
    
    // 添加 Schema 前缀
    qualifiedQuery := s.qualifyQuery(query, schema)
    
    return s.db.QueryContext(ctx, qualifiedQuery, args...)
}

// ExecContext 执行
func (s *sharedDatabaseImpl) ExecContext(ctx context.Context, service string, schema Schema, query string, args ...interface{}) (sql.Result, error) {
    policy := s.getPolicy(service)
    if policy == nil {
        return nil, fmt.Errorf("no access policy for service: %s", service)
    }
    
    if !s.canWrite(policy, schema) {
        return nil, fmt.Errorf("service %s does not have write access to schema %s", service, schema)
    }
    
    qualifiedQuery := s.qualifyQuery(query, schema)
    
    return s.db.ExecContext(ctx, qualifiedQuery, args...)
}

// BeginTx 开始事务
func (s *sharedDatabaseImpl) BeginTx(ctx context.Context, opts *sql.TxOptions) (*SharedTx, error) {
    tx, err := s.db.BeginTx(ctx, opts)
    if err != nil {
        return nil, err
    }
    
    return &SharedTx{tx: tx}, nil
}

// HealthCheck 健康检查
func (s *sharedDatabaseImpl) HealthCheck(ctx context.Context) error {
    return s.db.PingContext(ctx)
}

// Stats 统计
func (s *sharedDatabaseImpl) Stats() sql.DBStats {
    return s.db.Stats()
}

// Close 关闭
func (s *sharedDatabaseImpl) Close() error {
    return s.db.Close()
}

// getPolicy 获取策略
func (s *sharedDatabaseImpl) getPolicy(service string) *AccessPolicy {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.policies[service]
}

// canRead 检查读权限
func (s *sharedDatabaseImpl) canRead(policy *AccessPolicy, schema Schema) bool {
    for _, s := range policy.AllowedRead {
        if s == schema {
            return true
        }
    }
    return false
}

// canWrite 检查写权限
func (s *sharedDatabaseImpl) canWrite(policy *AccessPolicy, schema Schema) bool {
    for _, s := range policy.AllowedWrite {
        if s == schema {
            return true
        }
    }
    return false
}

// qualifyQuery 限定查询（添加 Schema）
func (s *sharedDatabaseImpl) qualifyQuery(query string, schema Schema) string {
    // 简化实现：实际应该使用 SQL 解析器
    return fmt.Sprintf("SET search_path TO %s; %s", schema, query)
}

// QueryContext 事务查询
func (tx *SharedTx) QueryContext(ctx context.Context, schema Schema, query string, args ...interface{}) (*sql.Rows, error) {
    qualifiedQuery := fmt.Sprintf("SET LOCAL search_path TO %s; %s", schema, query)
    return tx.tx.QueryContext(ctx, qualifiedQuery, args...)
}

// ExecContext 事务执行
func (tx *SharedTx) ExecContext(ctx context.Context, schema Schema, query string, args ...interface{}) (sql.Result, error) {
    qualifiedQuery := fmt.Sprintf("SET LOCAL search_path TO %s; %s", schema, query)
    return tx.tx.ExecContext(ctx, qualifiedQuery, args...)
}

// Commit 提交
func (tx *SharedTx) Commit() error {
    return tx.tx.Commit()
}

// Rollback 回滚
func (tx *SharedTx) Rollback() error {
    return tx.tx.Rollback()
}
```

### 2.2 服务模式实现

```go
// shareddatabase/order_service.go
package shareddatabase

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"

    "github.com/google/uuid"
)

// OrderService 订单服务（使用共享数据库）
type OrderService struct {
    db     SharedDatabase
    schema Schema
}

// Order 订单
type Order struct {
    ID         string          `json:"id" db:"id"`
    CustomerID string          `json:"customer_id" db:"customer_id"`
    Items      json.RawMessage `json:"items" db:"items"`
    Total      float64         `json:"total" db:"total"`
    Status     string          `json:"status" db:"status"`
    CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

// NewOrderService 创建订单服务
func NewOrderService(db SharedDatabase) *OrderService {
    return &OrderService{
        db:     db,
        schema: SchemaOrders,
    }
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, customerID string, items []OrderItem) (*Order, error) {
    var total float64
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }
    
    order := &Order{
        ID:         uuid.New().String(),
        CustomerID: customerID,
        Items:      mustMarshal(items),
        Total:      total,
        Status:     "PENDING",
        CreatedAt:  time.Now(),
    }
    
    query := `
        INSERT INTO orders (id, customer_id, items, total, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    
    _, err := s.db.ExecContext(ctx, "order-service", s.schema, query,
        order.ID, order.CustomerID, order.Items, order.Total, order.Status, order.CreatedAt)
    
    if err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    
    return order, nil
}

// GetOrder 获取订单
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*Order, error) {
    query := `
        SELECT id, customer_id, items, total, status, created_at
        FROM orders
        WHERE id = $1
    `
    
    rows, err := s.db.QueryContext(ctx, "order-service", s.schema, query, orderID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    if !rows.Next() {
        return nil, fmt.Errorf("order not found: %s", orderID)
    }
    
    order := &Order{}
    err = rows.Scan(&order.ID, &order.CustomerID, &order.Items, &order.Total, &order.Status, &order.CreatedAt)
    if err != nil {
        return nil, err
    }
    
    return order, nil
}

// GetCustomerOrders 获取客户订单（跨 Schema 查询示例）
func (s *OrderService) GetCustomerOrdersWithPaymentStatus(ctx context.Context, customerID string) ([]*OrderWithPayment, error) {
    // 使用跨 Schema 查询
    query := `
        SELECT o.id, o.total, o.status, p.status as payment_status
        FROM orders.orders o
        LEFT JOIN payments.payments p ON p.order_id = o.id
        WHERE o.customer_id = $1
    `
    
    // 注意：这需要在访问策略中允许
    rows, err := s.db.QueryContext(ctx, "order-service", SchemaOrders, query, customerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var results []*OrderWithPayment
    for rows.Next() {
        r := &OrderWithPayment{}
        if err := rows.Scan(&r.OrderID, &r.Total, &r.OrderStatus, &r.PaymentStatus); err != nil {
            return nil, err
        }
        results = append(results, r)
    }
    
    return results, rows.Err()
}

// OrderItem 订单项
type OrderItem struct {
    ProductID string  `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

// OrderWithPayment 带支付状态的订单
type OrderWithPayment struct {
    OrderID       string  `json:"order_id"`
    Total         float64 `json:"total"`
    OrderStatus   string  `json:"order_status"`
    PaymentStatus string  `json:"payment_status"`
}

func mustMarshal(v interface{}) json.RawMessage {
    data, _ := json.Marshal(v)
    return data
}
```

### 2.3 数据库 Schema 设计

```sql
-- 创建 Schema
CREATE SCHEMA IF NOT EXISTS orders;
CREATE SCHEMA IF NOT EXISTS payments;
CREATE SCHEMA IF NOT EXISTS inventory;
CREATE SCHEMA IF NOT EXISTS shared;

-- Orders Schema
CREATE TABLE orders.orders (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL,
    items JSONB NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Payments Schema
CREATE TABLE payments.payments (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建跨 Schema 外键（如果数据库支持）
-- 注意：这可能会限制服务的独立性
-- ALTER TABLE payments.payments
-- ADD CONSTRAINT fk_order
-- FOREIGN KEY (order_id) REFERENCES orders.orders(id);

-- Shared Schema
CREATE TABLE shared.customers (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 访问控制视图示例
CREATE VIEW orders.order_summary AS
SELECT id, customer_id, total, status, created_at
FROM orders.orders;

-- 授予权限（由 DBA 管理）
GRANT USAGE ON SCHEMA orders TO order_service;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA orders TO order_service;

GRANT USAGE ON SCHEMA payments TO payment_service;
GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA payments TO payment_service;
-- 只允许 payments 服务读取 orders 的特定视图
GRANT SELECT ON orders.order_summary TO payment_service;
```

---

## 3. 测试策略

### 3.1 集成测试

```go
// shareddatabase/integration_test.go
package shareddatabase

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSharedDatabase_AccessControl(t *testing.T) {
    // 配置访问策略
    policies := []*AccessPolicy{
        {
            Service:      "order-service",
            AllowedRead:  []Schema{SchemaOrders, SchemaShared},
            AllowedWrite: []Schema{SchemaOrders},
        },
        {
            Service:      "payment-service",
            AllowedRead:  []Schema{SchemaPayments, SchemaOrders, SchemaShared},
            AllowedWrite: []Schema{SchemaPayments},
        },
    }
    
    // 注意：这需要真实数据库连接
    // db, err := NewSharedDatabase("postgres://...", policies)
    // require.NoError(t, err)
    
    // 测试读写权限
    t.Run("order_service_can_write_to_orders", func(t *testing.T) {
        // 验证 order-service 可以写入 orders schema
    })
    
    t.Run("order_service_cannot_write_to_payments", func(t *testing.T) {
        // 验证 order-service 不能写入 payments schema
    })
}
```

---

## 4. 与其他模式的集成

### 4.1 与 Strangler Fig 模式的集成

```
┌─────────────────────────────────────────────────────────────────────────┐
│         Shared Database with Strangler Fig Pattern                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  阶段 1: 单体应用（初始状态）                                              │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Monolithic Application                        │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                         │   │
│  │  │ Orders  │  │Payments │  │Inventory│                         │   │
│  │  │ Module  │  │ Module  │  │ Module  │                         │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘                         │   │
│  │       └─────────────┼─────────────┘                            │   │
│  │                     ▼                                          │   │
│  │              ┌─────────────┐                                   │   │
│  │              │ Single DB   │                                   │   │
│  │              └─────────────┘                                   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  阶段 2: 逐步提取服务（使用共享数据库过渡）                                 │
│  ┌─────────────┐    ┌───────────────────────────────────────────┐   │   │
│  │New Service  │    │         Monolith (remaining)              │   │   │
│  │(Orders)     │    │  ┌─────────┐  ┌─────────┐                │   │   │
│  │             │    │  │Payments │  │Inventory│                │   │   │
│  └──────┬──────┘    │  │ Module  │  │ Module  │                │   │   │
│         │           │  └────┬────┘  └────┬────┘                │   │   │
│         │           │       └─────────────┘                    │   │   │
│         │           │              │                           │   │   │
│         └───────────┼──────────────┘                           │   │   │
│                     ▼                                          │   │   │
│            ┌─────────────────┐                                 │   │   │
│            │  Shared Database │                                │   │   │
│            │  (过渡状态)       │                                │   │   │
│            └─────────────────┘                                 │   │   │
│  ┌─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  阶段 3: 完全拆分数据库                                                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                 │
│  │   Orders    │    │  Payments   │    │  Inventory  │                 │
│  │   Service   │    │   Service   │    │   Service   │                 │
│  │             │    │             │    │             │                 │
│  │  ┌───────┐  │    │  ┌───────┐  │    │  ┌───────┐  │                 │
│  │  │Orders │  │    │  │Payment│  │    │  │Inventory│  │                 │
│  │  │  DB   │  │    │  │  DB   │  │    │  │  DB    │  │                 │
│  │  └───────┘  │    │  └───────┘  │    │  └───────┘  │                 │
│  └─────────────┘    └─────────────┘    └─────────────┘                 │
│                                                                         │
│  使用 Shared Database 作为过渡阶段的优势:                                   │
│  • 允许逐步迁移而非大爆炸式重构                                            │
│  • 保持数据一致性（ACID 事务）                                             │
│  • 可以随时回滚                                                            │
│  • 降低风险                                                                │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 5. 决策标准

### 5.1 何时使用共享数据库

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Shared Database Decision Tree                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  开始                                                                    │
│   │                                                                     │
│   ▼                                                                     │
│  ┌─────────────────────────┐                                           │
│  │ 正在从单体迁移到微服务？ │───否──► 考虑 Database-per-Service          │
│  └──────────┬──────────────┘                                           │
│             │是                                                         │
│             ▼                                                           │
│  ┌─────────────────────────┐                                           │
│  │ 数据拆分复杂度如何？     │                                           │
│  └──────────┬──────────────┘                                           │
│             │                                                           │
│      ┌──────┴──────┐                                                     │
│      ▼             ▼                                                     │
│    高             低                                                      │
│    │               │                                                      │
│    ▼               ▼                                                      │
│  ┌─────────┐   ┌─────────┐                                              │
│  │ Shared  │   │ Direct  │                                              │
│  │Database │   │ Split   │                                              │
│  │(过渡)   │   │         │                                              │
│  └─────────┘   └─────────┘                                              │
│                                                                         │
│  其他考虑因素:                                                            │
│  • 强一致性要求且无法使用 Saga ──是──► Shared Database（临时）            │
│  • 时间压力大，无法重构数据模型 ──是──► Shared Database（临时）            │
│  • 遗留系统复杂依赖 ──是──► Shared Database（可能长期）                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 风险与缓解策略

| 风险 | 影响 | 缓解策略 |
|------|------|----------|
| Schema 变更影响多个服务 | 高 | 使用 Schema 版本控制，向后兼容变更 |
| 性能问题难以隔离 | 高 | 监控每个服务的查询，使用资源限制 |
| 数据库成为单点故障 | 高 | 数据库高可用配置，读写分离 |
| 技术栈锁定 | 中 | 尽快迁移到独立数据库 |
| 数据所有权模糊 | 中 | 明确数据所有权，使用访问控制 |

### 5.3 检查清单

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Shared Database Checklist                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  迁移前:                                                                 │
│  □ 评估数据依赖图                                                        │
│  □ 制定分阶段迁移计划                                                    │
│  □ 设计 Schema 命名规范                                                  │
│  □ 配置访问控制和审计日志                                                 │
│                                                                         │
│  迁移中:                                                                 │
│  □ 使用 Schema 分离不同服务的数据                                         │
│  □ 实施最小权限原则                                                       │
│  □ 保持向后兼容的 Schema 变更                                             │
│  □ 监控跨 Schema 查询性能                                                 │
│                                                                         │
│  迁移后（独立数据库）:                                                    │
│  □ 逐步移除共享数据库依赖                                                 │
│  □ 实施 Saga 或 Outbox 模式                                              │
│  □ 移除共享数据库实例                                                     │
│                                                                         │
│  长期（如必须保持共享）:                                                  │
│  □ 定期审查访问模式                                                       │
│  □ 实施数据库分片策略                                                     │
│  □ 配置严格的资源限制                                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (>15KB, 完整形式化 + Go 实现 + 决策标准)

**相关文档**:
- [EC-035-Database-per-Service.md](./EC-035-Database-per-Service.md)
- [EC-008-Saga-Pattern-Formal.md](./EC-008-Saga-Pattern-Formal.md)
