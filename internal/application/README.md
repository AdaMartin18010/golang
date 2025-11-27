# Application Layer (应用层)

Clean Architecture 的应用层，包含用例编排。

## 结构

```text
application/
├── user/          # 用户用例
│   ├── dto.go         # 用户 DTO
│   └── service.go     # 用户应用服务
├── order/         # 订单用例
│   ├── dto.go         # 订单 DTO
│   └── service.go     # 订单应用服务
└── product/       # 产品用例
    ├── dto.go         # 产品 DTO
    └── service.go     # 产品应用服务
```

## 规则

- ✅ 只能导入 domain 层
- ✅ 协调领域对象
- ✅ 包含用例逻辑
- ❌ 不能导入 infrastructure 或 interfaces 层

## 应用服务

### User Service (用户服务)

**功能**:
- `CreateUser()` - 创建用户
- `GetUser()` - 获取用户
- `UpdateUser()` - 更新用户
- `DeleteUser()` - 删除用户
- `ListUsers()` - 列出用户

**DTO**:
- `CreateUserRequest` - 创建用户请求
- `UpdateUserRequest` - 更新用户请求
- `UserDTO` - 用户数据传输对象

### Order Service (订单服务)

**功能**:
- `CreateOrder()` - 创建订单
- `GetOrder()` - 获取订单
- `GetUserOrders()` - 获取用户订单列表
- `PayOrder()` - 支付订单
- `ShipOrder()` - 发货
- `DeliverOrder()` - 送达
- `CancelOrder()` - 取消订单
- `RefundOrder()` - 退款
- `UpdateOrder()` - 更新订单

**DTO**:
- `CreateOrderRequest` - 创建订单请求
- `CreateOrderItemRequest` - 创建订单项请求
- `UpdateOrderRequest` - 更新订单请求
- `OrderDTO` - 订单数据传输对象
- `OrderItemDTO` - 订单项数据传输对象

### Product Service (产品服务)

**功能**:
- `CreateProduct()` - 创建产品
- `GetProduct()` - 获取产品
- `GetProductBySKU()` - 根据 SKU 获取产品
- `UpdateProduct()` - 更新产品
- `DeleteProduct()` - 删除产品
- `ListProducts()` - 列出产品
- `SearchProducts()` - 搜索产品
- `GetProductsByCategory()` - 根据分类获取产品
- `ActivateProduct()` - 上架产品
- `DeactivateProduct()` - 下架产品
- `UpdateStock()` - 更新库存
- `IncreaseStock()` - 增加库存
- `DecreaseStock()` - 减少库存

**DTO**:
- `CreateProductRequest` - 创建产品请求
- `UpdateProductRequest` - 更新产品请求
- `ProductDTO` - 产品数据传输对象

## 设计原则

1. **用例编排**: 应用服务负责协调领域对象完成业务用例
2. **DTO 转换**: 应用层负责领域对象和 DTO 之间的转换
3. **事务边界**: 应用服务方法通常是一个事务边界
4. **依赖注入**: 通过构造函数注入领域仓储和领域服务
