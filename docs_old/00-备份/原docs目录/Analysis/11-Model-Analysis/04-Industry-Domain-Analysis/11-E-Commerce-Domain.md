# 电子商务领域分析

## 目录

1. [概述](#概述)
2. [形式化定义](#形式化定义)
3. [电商架构](#电商架构)
4. [订单管理](#订单管理)
5. [支付系统](#支付系统)
6. [最佳实践](#最佳实践)

## 概述

电子商务是现代商业的重要形式，涉及商品管理、订单处理、支付结算等多个技术领域。本文档从电商架构、订单管理、支付系统等维度深入分析电子商务领域的Golang实现方案。

### 核心特征

- **商品管理**: 产品目录和库存
- **订单处理**: 购物车和订单流程
- **支付结算**: 多种支付方式
- **用户管理**: 客户账户和权限
- **数据分析**: 销售和用户行为分析

## 形式化定义

### 电商系统定义

**定义 15.1** (电商系统)
电商系统是一个八元组 $\mathcal{ECS} = (U, P, O, C, I, S, A, T)$，其中：

- $U$ 是用户集合 (Users)
- $P$ 是产品集合 (Products)
- $O$ 是订单集合 (Orders)
- $C$ 是购物车 (Cart)
- $I$ 是库存系统 (Inventory)
- $S$ 是销售系统 (Sales)
- $A$ 是分析系统 (Analytics)
- $T$ 是交易系统 (Transaction)

**定义 15.2** (订单)
订单是一个六元组 $\mathcal{Order} = (I, U, P, Q, S, T)$，其中：

- $I$ 是订单信息 (Order Information)
- $U$ 是用户信息 (User Information)
- $P$ 是产品信息 (Product Information)
- $Q$ 是数量 (Quantity)
- $S$ 是状态 (Status)
- $T$ 是时间戳 (Timestamp)

### 库存管理模型

**定义 15.3** (库存管理)
库存管理是一个四元组 $\mathcal{IM} = (S, R, A, L)$，其中：

- $S$ 是库存状态 (Stock Status)
- $R$ 是补货规则 (Replenishment Rules)
- $A$ 是可用性 (Availability)
- $L$ 是库存水平 (Stock Level)

**性质 15.1** (库存约束)
对于任意产品 $p$，必须满足：
$\text{stock}(p) \geq 0$

## 电商架构

### 产品管理系统

```go
// 产品
type Product struct {
    ID          string
    Name        string
    Description string
    Category    string
    Brand       string
    SKU         string
    Price       *Price
    Inventory   *Inventory
    Images      []string
    Attributes  map[string]interface{}
    Status      ProductStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
    mu          sync.RWMutex
}

// 价格
type Price struct {
    RegularPrice float64
    SalePrice    *float64
    Currency     string
    Discount     float64
}

// 库存
type Inventory struct {
    Quantity    int
    Reserved    int
    Available   int
    LowStock    int
    OutOfStock  bool
}

// 产品状态
type ProductStatus string

const (
    ProductStatusActive   ProductStatus = "active"
    ProductStatusInactive ProductStatus = "inactive"
    ProductStatusDraft    ProductStatus = "draft"
    ProductStatusDeleted  ProductStatus = "deleted"
)

// 产品管理器
type ProductManager struct {
    products map[string]*Product
    categories map[string]*Category
    mu       sync.RWMutex
}

// 分类
type Category struct {
    ID          string
    Name        string
    Description string
    ParentID    *string
    Level       int
    Products    []string
}

// 创建产品
func (pm *ProductManager) CreateProduct(product *Product) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()
    
    if _, exists := pm.products[product.ID]; exists {
        return fmt.Errorf("product %s already exists", product.ID)
    }
    
    // 验证产品信息
    if err := pm.validateProduct(product); err != nil {
        return fmt.Errorf("product validation failed: %w", err)
    }
    
    // 设置默认值
    product.CreatedAt = time.Now()
    product.UpdatedAt = time.Now()
    product.Status = ProductStatusDraft
    
    // 计算可用库存
    if product.Inventory != nil {
        product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved
        product.Inventory.OutOfStock = product.Inventory.Available <= 0
    }
    
    // 创建产品
    pm.products[product.ID] = product
    
    return nil
}

// 验证产品信息
func (pm *ProductManager) validateProduct(product *Product) error {
    if product.ID == "" {
        return fmt.Errorf("product ID is required")
    }
    
    if product.Name == "" {
        return fmt.Errorf("product name is required")
    }
    
    if product.SKU == "" {
        return fmt.Errorf("SKU is required")
    }
    
    if product.Price == nil {
        return fmt.Errorf("price is required")
    }
    
    if product.Price.RegularPrice <= 0 {
        return fmt.Errorf("regular price must be positive")
    }
    
    return nil
}

// 获取产品
func (pm *ProductManager) GetProduct(productID string) (*Product, error) {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    
    product, exists := pm.products[productID]
    if !exists {
        return nil, fmt.Errorf("product %s not found", productID)
    }
    
    return product, nil
}

// 更新库存
func (pm *ProductManager) UpdateInventory(productID string, quantity int, operation InventoryOperation) error {
    product, err := pm.GetProduct(productID)
    if err != nil {
        return err
    }
    
    product.mu.Lock()
    defer product.mu.Unlock()
    
    switch operation {
    case InventoryOperationAdd:
        product.Inventory.Quantity += quantity
    case InventoryOperationSubtract:
        if product.Inventory.Quantity < quantity {
            return fmt.Errorf("insufficient inventory")
        }
        product.Inventory.Quantity -= quantity
    case InventoryOperationReserve:
        if product.Inventory.Available < quantity {
            return fmt.Errorf("insufficient available inventory")
        }
        product.Inventory.Reserved += quantity
    case InventoryOperationRelease:
        if product.Inventory.Reserved < quantity {
            return fmt.Errorf("insufficient reserved inventory")
        }
        product.Inventory.Reserved -= quantity
    }
    
    // 更新可用库存
    product.Inventory.Available = product.Inventory.Quantity - product.Inventory.Reserved
    product.Inventory.OutOfStock = product.Inventory.Available <= 0
    
    product.UpdatedAt = time.Now()
    
    return nil
}

// 库存操作类型
type InventoryOperation string

const (
    InventoryOperationAdd     InventoryOperation = "add"
    InventoryOperationSubtract InventoryOperation = "subtract"
    InventoryOperationReserve InventoryOperation = "reserve"
    InventoryOperationRelease InventoryOperation = "release"
)

// 搜索产品
func (pm *ProductManager) SearchProducts(query string, filters map[string]interface{}) ([]*Product, error) {
    pm.mu.RLock()
    defer pm.mu.RUnlock()
    
    var results []*Product
    
    for _, product := range pm.products {
        if product.Status != ProductStatusActive {
            continue
        }
        
        // 检查查询匹配
        if query != "" {
            if !pm.matchesQuery(product, query) {
                continue
            }
        }
        
        // 检查过滤器
        if !pm.matchesFilters(product, filters) {
            continue
        }
        
        results = append(results, product)
    }
    
    return results, nil
}

// 检查查询匹配
func (pm *ProductManager) matchesQuery(product *Product, query string) bool {
    query = strings.ToLower(query)
    
    if strings.Contains(strings.ToLower(product.Name), query) {
        return true
    }
    
    if strings.Contains(strings.ToLower(product.Description), query) {
        return true
    }
    
    if strings.Contains(strings.ToLower(product.Brand), query) {
        return true
    }
    
    return false
}

// 检查过滤器匹配
func (pm *ProductManager) matchesFilters(product *Product, filters map[string]interface{}) bool {
    for key, value := range filters {
        switch key {
        case "category":
            if category, ok := value.(string); ok {
                if product.Category != category {
                    return false
                }
            }
        case "brand":
            if brand, ok := value.(string); ok {
                if product.Brand != brand {
                    return false
                }
            }
        case "min_price":
            if minPrice, ok := value.(float64); ok {
                if product.Price.RegularPrice < minPrice {
                    return false
                }
            }
        case "max_price":
            if maxPrice, ok := value.(float64); ok {
                if product.Price.RegularPrice > maxPrice {
                    return false
                }
            }
        case "in_stock":
            if inStock, ok := value.(bool); ok {
                if inStock && product.Inventory.OutOfStock {
                    return false
                }
            }
        }
    }
    
    return true
}

```

### 购物车系统

```go
// 购物车
type Cart struct {
    ID          string
    UserID      string
    Items       []*CartItem
    Subtotal    float64
    Tax         float64
    Shipping    float64
    Total       float64
    CreatedAt   time.Time
    UpdatedAt   time.Time
    mu          sync.RWMutex
}

// 购物车项目
type CartItem struct {
    ID          string
    ProductID   string
    ProductName string
    Quantity    int
    Price       float64
    Total       float64
    AddedAt     time.Time
}

// 购物车管理器
type CartManager struct {
    carts map[string]*Cart
    mu    sync.RWMutex
}

// 创建购物车
func (cm *CartManager) CreateCart(userID string) (*Cart, error) {
    cart := &Cart{
        ID:        uuid.New().String(),
        UserID:    userID,
        Items:     make([]*CartItem, 0),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    cm.mu.Lock()
    cm.carts[cart.ID] = cart
    cm.mu.Unlock()
    
    return cart, nil
}

// 获取购物车
func (cm *CartManager) GetCart(cartID string) (*Cart, error) {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    
    cart, exists := cm.carts[cartID]
    if !exists {
        return nil, fmt.Errorf("cart %s not found", cartID)
    }
    
    return cart, nil
}

// 添加商品到购物车
func (cm *CartManager) AddToCart(cartID, productID string, quantity int) error {
    cart, err := cm.GetCart(cartID)
    if err != nil {
        return err
    }
    
    cart.mu.Lock()
    defer cart.mu.Unlock()
    
    // 检查商品是否已在购物车中
    for _, item := range cart.Items {
        if item.ProductID == productID {
            item.Quantity += quantity
            item.Total = float64(item.Quantity) * item.Price
            cm.updateCartTotals(cart)
            return nil
        }
    }
    
    // 获取产品信息
    product, err := cm.getProduct(productID)
    if err != nil {
        return err
    }
    
    // 检查库存
    if product.Inventory.Available < quantity {
        return fmt.Errorf("insufficient inventory")
    }
    
    // 创建购物车项目
    item := &CartItem{
        ID:          uuid.New().String(),
        ProductID:   productID,
        ProductName: product.Name,
        Quantity:    quantity,
        Price:       product.Price.RegularPrice,
        Total:       product.Price.RegularPrice * float64(quantity),
        AddedAt:     time.Now(),
    }
    
    cart.Items = append(cart.Items, item)
    cm.updateCartTotals(cart)
    
    return nil
}

// 更新购物车项目数量
func (cm *CartManager) UpdateCartItemQuantity(cartID, itemID string, quantity int) error {
    cart, err := cm.GetCart(cartID)
    if err != nil {
        return err
    }
    
    cart.mu.Lock()
    defer cart.mu.Unlock()
    
    for _, item := range cart.Items {
        if item.ID == itemID {
            if quantity <= 0 {
                // 移除项目
                cm.removeCartItem(cart, itemID)
            } else {
                // 检查库存
                product, err := cm.getProduct(item.ProductID)
                if err != nil {
                    return err
                }
                
                if product.Inventory.Available < quantity {
                    return fmt.Errorf("insufficient inventory")
                }
                
                item.Quantity = quantity
                item.Total = float64(quantity) * item.Price
            }
            
            cm.updateCartTotals(cart)
            return nil
        }
    }
    
    return fmt.Errorf("cart item %s not found", itemID)
}

// 移除购物车项目
func (cm *CartManager) RemoveFromCart(cartID, itemID string) error {
    cart, err := cm.GetCart(cartID)
    if err != nil {
        return err
    }
    
    cart.mu.Lock()
    defer cart.mu.Unlock()
    
    cm.removeCartItem(cart, itemID)
    cm.updateCartTotals(cart)
    
    return nil
}

// 移除购物车项目
func (cm *CartManager) removeCartItem(cart *Cart, itemID string) {
    for i, item := range cart.Items {
        if item.ID == itemID {
            cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
            break
        }
    }
}

// 更新购物车总计
func (cm *CartManager) updateCartTotals(cart *Cart) {
    subtotal := 0.0
    for _, item := range cart.Items {
        subtotal += item.Total
    }
    
    cart.Subtotal = subtotal
    cart.Tax = subtotal * 0.1 // 10% 税率
    cart.Shipping = 0.0 // 免费配送
    cart.Total = cart.Subtotal + cart.Tax + cart.Shipping
    cart.UpdatedAt = time.Now()
}

// 获取产品信息
func (cm *CartManager) getProduct(productID string) (*Product, error) {
    // 这里应该从产品管理器获取产品信息
    // 简化实现
    return &Product{
        ID:   productID,
        Name: "Sample Product",
        Price: &Price{
            RegularPrice: 100.0,
        },
        Inventory: &Inventory{
            Available: 100,
        },
    }, nil
}

```

## 订单管理

### 订单系统

```go
// 订单
type Order struct {
    ID          string
    UserID      string
    Items       []*OrderItem
    Status      OrderStatus
    Subtotal    float64
    Tax         float64
    Shipping    float64
    Discount    float64
    Total       float64
    ShippingAddress *Address
    BillingAddress  *Address
    PaymentMethod   string
    CreatedAt       time.Time
    UpdatedAt       time.Time
    mu              sync.RWMutex
}

// 订单项目
type OrderItem struct {
    ID          string
    ProductID   string
    ProductName string
    Quantity    int
    Price       float64
    Total       float64
}

// 订单状态
type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusConfirmed OrderStatus = "confirmed"
    OrderStatusProcessing OrderStatus = "processing"
    OrderStatusShipped   OrderStatus = "shipped"
    OrderStatusDelivered OrderStatus = "delivered"
    OrderStatusCancelled OrderStatus = "cancelled"
    OrderStatusRefunded  OrderStatus = "refunded"
)

// 地址
type Address struct {
    FirstName string
    LastName  string
    Street    string
    City      string
    State     string
    ZipCode   string
    Country   string
    Phone     string
}

// 订单管理器
type OrderManager struct {
    orders map[string]*Order
    mu     sync.RWMutex
}

// 创建订单
func (om *OrderManager) CreateOrder(userID string, cart *Cart, shippingAddress, billingAddress *Address, paymentMethod string) (*Order, error) {
    // 验证购物车
    if len(cart.Items) == 0 {
        return nil, fmt.Errorf("cart is empty")
    }
    
    // 检查库存
    if err := om.checkInventory(cart.Items); err != nil {
        return nil, err
    }
    
    // 创建订单
    order := &Order{
        ID:              uuid.New().String(),
        UserID:          userID,
        Status:          OrderStatusPending,
        ShippingAddress: shippingAddress,
        BillingAddress:  billingAddress,
        PaymentMethod:   paymentMethod,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }
    
    // 转换购物车项目为订单项目
    for _, cartItem := range cart.Items {
        orderItem := &OrderItem{
            ID:          uuid.New().String(),
            ProductID:   cartItem.ProductID,
            ProductName: cartItem.ProductName,
            Quantity:    cartItem.Quantity,
            Price:       cartItem.Price,
            Total:       cartItem.Total,
        }
        order.Items = append(order.Items, orderItem)
    }
    
    // 计算订单总计
    om.calculateOrderTotals(order)
    
    // 保存订单
    om.mu.Lock()
    om.orders[order.ID] = order
    om.mu.Unlock()
    
    // 预留库存
    om.reserveInventory(order.Items)
    
    return order, nil
}

// 检查库存
func (om *OrderManager) checkInventory(items []*CartItem) error {
    for _, item := range items {
        product, err := om.getProduct(item.ProductID)
        if err != nil {
            return err
        }
        
        if product.Inventory.Available < item.Quantity {
            return fmt.Errorf("insufficient inventory for product %s", product.Name)
        }
    }
    
    return nil
}

// 预留库存
func (om *OrderManager) reserveInventory(items []*OrderItem) {
    for _, item := range items {
        // 这里应该调用产品管理器预留库存
        // 简化实现
        log.Printf("Reserving %d units of product %s", item.Quantity, item.ProductID)
    }
}

// 计算订单总计
func (om *OrderManager) calculateOrderTotals(order *Order) {
    subtotal := 0.0
    for _, item := range order.Items {
        subtotal += item.Total
    }
    
    order.Subtotal = subtotal
    order.Tax = subtotal * 0.1 // 10% 税率
    order.Shipping = 0.0 // 免费配送
    order.Discount = 0.0 // 无折扣
    order.Total = order.Subtotal + order.Tax + order.Shipping - order.Discount
}

// 获取订单
func (om *OrderManager) GetOrder(orderID string) (*Order, error) {
    om.mu.RLock()
    defer om.mu.RUnlock()
    
    order, exists := om.orders[orderID]
    if !exists {
        return nil, fmt.Errorf("order %s not found", orderID)
    }
    
    return order, nil
}

// 更新订单状态
func (om *OrderManager) UpdateOrderStatus(orderID string, status OrderStatus) error {
    order, err := om.GetOrder(orderID)
    if err != nil {
        return err
    }
    
    order.mu.Lock()
    order.Status = status
    order.UpdatedAt = time.Now()
    order.mu.Unlock()
    
    // 处理状态变更
    om.handleStatusChange(order, status)
    
    return nil
}

// 处理状态变更
func (om *OrderManager) handleStatusChange(order *Order, status OrderStatus) {
    switch status {
    case OrderStatusConfirmed:
        // 确认订单，开始处理
        log.Printf("Order %s confirmed", order.ID)
    case OrderStatusProcessing:
        // 开始处理订单
        log.Printf("Order %s processing", order.ID)
    case OrderStatusShipped:
        // 订单已发货
        log.Printf("Order %s shipped", order.ID)
    case OrderStatusDelivered:
        // 订单已送达
        log.Printf("Order %s delivered", order.ID)
    case OrderStatusCancelled:
        // 取消订单，释放库存
        om.releaseInventory(order.Items)
        log.Printf("Order %s cancelled", order.ID)
    case OrderStatusRefunded:
        // 退款处理
        log.Printf("Order %s refunded", order.ID)
    }
}

// 释放库存
func (om *OrderManager) releaseInventory(items []*OrderItem) {
    for _, item := range items {
        // 这里应该调用产品管理器释放库存
        // 简化实现
        log.Printf("Releasing %d units of product %s", item.Quantity, item.ProductID)
    }
}

// 获取产品信息
func (om *OrderManager) getProduct(productID string) (*Product, error) {
    // 这里应该从产品管理器获取产品信息
    // 简化实现
    return &Product{
        ID:   productID,
        Name: "Sample Product",
        Inventory: &Inventory{
            Available: 100,
        },
    }, nil
}

// 获取用户订单
func (om *OrderManager) GetUserOrders(userID string) ([]*Order, error) {
    om.mu.RLock()
    defer om.mu.RUnlock()
    
    var userOrders []*Order
    for _, order := range om.orders {
        if order.UserID == userID {
            userOrders = append(userOrders, order)
        }
    }
    
    return userOrders, nil
}

```

## 支付系统

### 支付处理

```go
// 支付
type Payment struct {
    ID          string
    OrderID     string
    UserID      string
    Amount      float64
    Currency    string
    Method      PaymentMethod
    Status      PaymentStatus
    TransactionID string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    mu          sync.RWMutex
}

// 支付方式
type PaymentMethod string

const (
    PaymentMethodCreditCard PaymentMethod = "credit_card"
    PaymentMethodDebitCard  PaymentMethod = "debit_card"
    PaymentMethodPayPal     PaymentMethod = "paypal"
    PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
    PaymentMethodCrypto     PaymentMethod = "crypto"
)

// 支付状态
type PaymentStatus string

const (
    PaymentStatusPending   PaymentStatus = "pending"
    PaymentStatusProcessing PaymentStatus = "processing"
    PaymentStatusCompleted PaymentStatus = "completed"
    PaymentStatusFailed    PaymentStatus = "failed"
    PaymentStatusRefunded  PaymentStatus = "refunded"
)

// 支付处理器
type PaymentProcessor struct {
    payments map[string]*Payment
    gateways map[string]PaymentGateway
    mu       sync.RWMutex
}

// 支付网关接口
type PaymentGateway interface {
    ProcessPayment(payment *Payment) error
    RefundPayment(paymentID string, amount float64) error
    Name() string
}

// 信用卡支付网关
type CreditCardGateway struct{}

func (ccg *CreditCardGateway) Name() string {
    return "credit_card_gateway"
}

func (ccg *CreditCardGateway) ProcessPayment(payment *Payment) error {
    // 模拟支付处理
    payment.mu.Lock()
    payment.Status = PaymentStatusProcessing
    payment.UpdatedAt = time.Now()
    payment.mu.Unlock()
    
    // 模拟处理延迟
    time.Sleep(time.Second * 2)
    
    // 模拟支付结果
    success := rand.Float64() > 0.1 // 90% 成功率
    
    payment.mu.Lock()
    if success {
        payment.Status = PaymentStatusCompleted
        payment.TransactionID = uuid.New().String()
    } else {
        payment.Status = PaymentStatusFailed
    }
    payment.UpdatedAt = time.Now()
    payment.mu.Unlock()
    
    if !success {
        return fmt.Errorf("payment processing failed")
    }
    
    return nil
}

func (ccg *CreditCardGateway) RefundPayment(paymentID string, amount float64) error {
    // 模拟退款处理
    log.Printf("Processing refund for payment %s: %.2f", paymentID, amount)
    return nil
}

// PayPal支付网关
type PayPalGateway struct{}

func (ppg *PayPalGateway) Name() string {
    return "paypal_gateway"
}

func (ppg *PayPalGateway) ProcessPayment(payment *Payment) error {
    // 模拟PayPal支付处理
    payment.mu.Lock()
    payment.Status = PaymentStatusProcessing
    payment.UpdatedAt = time.Now()
    payment.mu.Unlock()
    
    // 模拟处理延迟
    time.Sleep(time.Second * 1)
    
    // 模拟支付结果
    success := rand.Float64() > 0.05 // 95% 成功率
    
    payment.mu.Lock()
    if success {
        payment.Status = PaymentStatusCompleted
        payment.TransactionID = uuid.New().String()
    } else {
        payment.Status = PaymentStatusFailed
    }
    payment.UpdatedAt = time.Now()
    payment.mu.Unlock()
    
    if !success {
        return fmt.Errorf("PayPal payment processing failed")
    }
    
    return nil
}

func (ppg *PayPalGateway) RefundPayment(paymentID string, amount float64) error {
    // 模拟PayPal退款处理
    log.Printf("Processing PayPal refund for payment %s: %.2f", paymentID, amount)
    return nil
}

// 处理支付
func (pp *PaymentProcessor) ProcessPayment(payment *Payment) error {
    // 验证支付信息
    if err := pp.validatePayment(payment); err != nil {
        return err
    }
    
    // 获取支付网关
    gateway, exists := pp.gateways[string(payment.Method)]
    if !exists {
        return fmt.Errorf("payment gateway not found for method %s", payment.Method)
    }
    
    // 保存支付记录
    pp.mu.Lock()
    pp.payments[payment.ID] = payment
    pp.mu.Unlock()
    
    // 处理支付
    err := gateway.ProcessPayment(payment)
    if err != nil {
        payment.mu.Lock()
        payment.Status = PaymentStatusFailed
        payment.UpdatedAt = time.Now()
        payment.mu.Unlock()
        return err
    }
    
    return nil
}

// 验证支付信息
func (pp *PaymentProcessor) validatePayment(payment *Payment) error {
    if payment.ID == "" {
        return fmt.Errorf("payment ID is required")
    }
    
    if payment.OrderID == "" {
        return fmt.Errorf("order ID is required")
    }
    
    if payment.UserID == "" {
        return fmt.Errorf("user ID is required")
    }
    
    if payment.Amount <= 0 {
        return fmt.Errorf("payment amount must be positive")
    }
    
    if payment.Method == "" {
        return fmt.Errorf("payment method is required")
    }
    
    return nil
}

// 获取支付
func (pp *PaymentProcessor) GetPayment(paymentID string) (*Payment, error) {
    pp.mu.RLock()
    defer pp.mu.RUnlock()
    
    payment, exists := pp.payments[paymentID]
    if !exists {
        return nil, fmt.Errorf("payment %s not found", paymentID)
    }
    
    return payment, nil
}

// 退款
func (pp *PaymentProcessor) RefundPayment(paymentID string, amount float64) error {
    payment, err := pp.GetPayment(paymentID)
    if err != nil {
        return err
    }
    
    if payment.Status != PaymentStatusCompleted {
        return fmt.Errorf("payment is not completed")
    }
    
    if amount > payment.Amount {
        return fmt.Errorf("refund amount cannot exceed payment amount")
    }
    
    // 获取支付网关
    gateway, exists := pp.gateways[string(payment.Method)]
    if !exists {
        return fmt.Errorf("payment gateway not found")
    }
    
    // 处理退款
    err = gateway.RefundPayment(paymentID, amount)
    if err != nil {
        return err
    }
    
    payment.mu.Lock()
    payment.Status = PaymentStatusRefunded
    payment.UpdatedAt = time.Now()
    payment.mu.Unlock()
    
    return nil
}

// 添加支付网关
func (pp *PaymentProcessor) AddGateway(method PaymentMethod, gateway PaymentGateway) {
    pp.mu.Lock()
    pp.gateways[string(method)] = gateway
    pp.mu.Unlock()
}

```

## 最佳实践

### 1. 错误处理

```go
// 电子商务错误类型
type ECommerceError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    UserID  string `json:"user_id,omitempty"`
    OrderID string `json:"order_id,omitempty"`
    Details string `json:"details,omitempty"`
}

func (e *ECommerceError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// 错误码常量
const (
    ErrCodeProductNotFound  = "PRODUCT_NOT_FOUND"
    ErrCodeInsufficientInventory = "INSUFFICIENT_INVENTORY"
    ErrCodeOrderNotFound    = "ORDER_NOT_FOUND"
    ErrCodePaymentFailed    = "PAYMENT_FAILED"
    ErrCodeCartEmpty        = "CART_EMPTY"
)

// 统一错误处理
func HandleECommerceError(err error, userID, orderID string) *ECommerceError {
    switch {
    case errors.Is(err, ErrProductNotFound):
        return &ECommerceError{
            Code:   ErrCodeProductNotFound,
            Message: "Product not found",
            UserID: userID,
        }
    case errors.Is(err, ErrInsufficientInventory):
        return &ECommerceError{
            Code:   ErrCodeInsufficientInventory,
            Message: "Insufficient inventory",
            UserID: userID,
        }
    case errors.Is(err, ErrPaymentFailed):
        return &ECommerceError{
            Code:    ErrCodePaymentFailed,
            Message: "Payment processing failed",
            OrderID: orderID,
        }
    default:
        return &ECommerceError{
            Code: ErrCodeOrderNotFound,
            Message: "Order not found",
        }
    }
}

```

### 2. 监控和日志

```go
// 电子商务指标
type ECommerceMetrics struct {
    productCount    prometheus.Gauge
    orderCount      prometheus.Counter
    revenueTotal    prometheus.Counter
    cartAbandonment prometheus.Counter
    errorCount      prometheus.Counter
}

func NewECommerceMetrics() *ECommerceMetrics {
    return &ECommerceMetrics{
        productCount: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "ecommerce_products_total",
            Help: "Total number of products",
        }),
        orderCount: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "ecommerce_orders_total",
            Help: "Total number of orders",
        }),
        revenueTotal: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "ecommerce_revenue_total",
            Help: "Total revenue",
        }),
        cartAbandonment: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "ecommerce_cart_abandonment_total",
            Help: "Total cart abandonments",
        }),
        errorCount: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "ecommerce_errors_total",
            Help: "Total number of e-commerce errors",
        }),
    }
}

// 电子商务日志
type ECommerceLogger struct {
    logger *zap.Logger
}

func (l *ECommerceLogger) LogProductCreated(product *Product) {
    l.logger.Info("product created",
        zap.String("product_id", product.ID),
        zap.String("product_name", product.Name),
        zap.String("category", product.Category),
        zap.Float64("price", product.Price.RegularPrice),
    )
}

func (l *ECommerceLogger) LogOrderCreated(order *Order) {
    l.logger.Info("order created",
        zap.String("order_id", order.ID),
        zap.String("user_id", order.UserID),
        zap.Float64("total", order.Total),
        zap.String("status", string(order.Status)),
    )
}

func (l *ECommerceLogger) LogPaymentProcessed(payment *Payment) {
    l.logger.Info("payment processed",
        zap.String("payment_id", payment.ID),
        zap.String("order_id", payment.OrderID),
        zap.Float64("amount", payment.Amount),
        zap.String("method", string(payment.Method)),
        zap.String("status", string(payment.Status)),
    )
}

```

### 3. 测试策略

```go
// 单元测试
func TestProductManager_CreateProduct(t *testing.T) {
    manager := &ProductManager{
        products:   make(map[string]*Product),
        categories: make(map[string]*Category),
    }
    
    product := &Product{
        ID:   "product1",
        Name: "Test Product",
        SKU:  "TEST001",
        Price: &Price{
            RegularPrice: 100.0,
            Currency:     "USD",
        },
        Inventory: &Inventory{
            Quantity: 50,
        },
    }
    
    // 测试创建产品
    err := manager.CreateProduct(product)
    if err != nil {
        t.Errorf("Failed to create product: %v", err)
    }
    
    if len(manager.products) != 1 {
        t.Errorf("Expected 1 product, got %d", len(manager.products))
    }
}

// 集成测试
func TestOrderManager_CreateOrder(t *testing.T) {
    // 创建订单管理器
    orderManager := &OrderManager{
        orders: make(map[string]*Order),
    }
    
    // 创建购物车
    cart := &Cart{
        ID:     "cart1",
        UserID: "user1",
        Items: []*CartItem{
            {
                ID:          "item1",
                ProductID:   "product1",
                ProductName: "Test Product",
                Quantity:    2,
                Price:       100.0,
                Total:       200.0,
            },
        },
        Total: 200.0,
    }
    
    // 创建地址
    address := &Address{
        FirstName: "John",
        LastName:  "Doe",
        Street:    "123 Main St",
        City:      "New York",
        State:     "NY",
        ZipCode:   "10001",
    }
    
    // 测试创建订单
    order, err := orderManager.CreateOrder("user1", cart, address, address, "credit_card")
    if err != nil {
        t.Errorf("Failed to create order: %v", err)
    }
    
    if order.UserID != "user1" {
        t.Errorf("Expected user ID 'user1', got '%s'", order.UserID)
    }
    
    if order.Total != 220.0 { // 200 + 20 tax
        t.Errorf("Expected total 220.0, got %.2f", order.Total)
    }
}

// 性能测试
func BenchmarkProductManager_SearchProducts(b *testing.B) {
    manager := &ProductManager{
        products: make(map[string]*Product),
    }
    
    // 创建测试产品
    for i := 0; i < 1000; i++ {
        product := &Product{
            ID:   fmt.Sprintf("product%d", i),
            Name: fmt.Sprintf("Product %d", i),
            Category: "electronics",
            Price: &Price{
                RegularPrice: float64(i * 10),
            },
            Status: ProductStatusActive,
        }
        manager.products[product.ID] = product
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := manager.SearchProducts("Product", map[string]interface{}{
            "category": "electronics",
        })
        if err != nil {
            b.Fatalf("Search failed: %v", err)
        }
    }
}

```

---

## 总结

本文档深入分析了电子商务领域的核心概念、技术架构和实现方案，包括：

1. **形式化定义**: 电商系统、订单、库存管理的数学建模
2. **电商架构**: 产品管理、购物车系统的设计
3. **订单管理**: 订单系统、状态管理的实现
4. **支付系统**: 支付处理、网关集成的实现
5. **最佳实践**: 错误处理、监控、测试策略

电子商务系统需要在商品管理、订单处理、支付结算等多个方面找到平衡，通过合理的架构设计和实现方案，可以构建出高效、可靠、用户友好的电子商务系统。

---

**最后更新**: 2024-12-19  
**当前状态**: ✅ 电子商务领域分析完成  
**下一步**: 汽车技术领域分析
