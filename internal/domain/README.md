# Domain Layer (领域层)

Clean Architecture 的领域层，包含核心业务逻辑。

## 结构

```text
domain/
├── user/          # 用户领域
│   ├── entity.go      # 用户实体
│   ├── repository.go  # 用户仓储接口
│   ├── service.go     # 用户领域服务接口
│   └── errors.go      # 用户领域错误
├── order/         # 订单领域
│   ├── entity.go      # 订单实体
│   ├── repository.go  # 订单仓储接口
│   ├── service.go     # 订单领域服务接口
│   └── errors.go      # 订单领域错误
└── product/       # 产品领域
    ├── entity.go      # 产品实体
    ├── repository.go  # 产品仓储接口
    ├── service.go     # 产品领域服务接口
    └── errors.go      # 产品领域错误
```

## 规则

- ✅ 不依赖任何外部框架
- ✅ 只包含业务逻辑
- ✅ 定义接口，不包含实现
- ❌ 不能导入 infrastructure 或 interfaces 层

## 领域模型

### User (用户领域)

**实体**: `User`

- ID、Email、Name
- 创建时间、更新时间

**业务方法**:

- `NewUser()` - 创建新用户
- `UpdateName()` - 更新用户名
- `UpdateEmail()` - 更新邮箱
- `IsValid()` - 验证用户有效性

**仓储接口**: `Repository`

- `Create()`, `FindByID()`, `FindByEmail()`, `Update()`, `Delete()`, `List()`

### Order (订单领域)

**实体**: `Order`

- ID、UserID、Items、TotalAmount、Status
- 订单状态流转：Pending → Paid → Shipped → Delivered
- 支持取消和退款

**订单状态**:

- `pending` - 待支付
- `paid` - 已支付
- `shipped` - 已发货
- `delivered` - 已送达
- `cancelled` - 已取消
- `refunded` - 已退款

**业务方法**:

- `Pay()` - 支付订单
- `Ship()` - 发货
- `Deliver()` - 送达
- `Cancel()` - 取消订单
- `Refund()` - 退款
- `AddItem()` - 添加订单项
- `RemoveItem()` - 移除订单项
- `CanBeCancelled()` - 检查是否可以取消
- `CanBeRefunded()` - 检查是否可以退款

**仓储接口**: `Repository`

- `Create()`, `FindByID()`, `FindByUserID()`, `Update()`, `Delete()`, `FindByStatus()`, `CountByUserID()`

### Product (产品领域)

**实体**: `Product`

- ID、Name、Description、Price、Stock、Status、CategoryID、SKU
- 产品状态：Active、Inactive、Deleted

**业务方法**:

- `UpdatePrice()` - 更新价格
- `UpdateStock()` - 更新库存
- `IncreaseStock()` - 增加库存
- `DecreaseStock()` - 减少库存
- `Activate()` - 上架产品
- `Deactivate()` - 下架产品
- `Delete()` - 删除产品（软删除）
- `IsAvailable()` - 检查产品是否可用
- `HasStock()` - 检查是否有库存

**仓储接口**: `Repository`

- `Create()`, `FindByID()`, `FindBySKU()`, `Update()`, `Delete()`, `List()`, `FindByCategory()`, `FindByStatus()`, `Search()`

## 设计原则

1. **领域驱动设计 (DDD)**: 每个领域都是独立的，包含完整的业务逻辑
2. **接口隔离**: 仓储接口和领域服务接口在领域层定义，实现层实现
3. **业务规则封装**: 所有业务规则都在实体方法中，保证业务逻辑的一致性
4. **状态管理**: 通过状态机管理订单状态流转，确保状态转换的合法性
