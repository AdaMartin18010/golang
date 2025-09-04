# 13.1 微服务案例分析

<!-- TOC START -->
- [13.1 微服务案例分析](#微服务案例分析)
  - [13.1.1 目录](#目录)
  - [13.1.2 概述](#概述)
    - [13.1.2.1 核心目标](#核心目标)
  - [13.1.3 电商微服务架构](#电商微服务架构)
    - [13.1.3.1 架构概述](#架构概述)
    - [13.1.3.2 服务划分](#服务划分)
  - [13.1.4 金融微服务架构](#金融微服务架构)
    - [13.1.4.1 架构概述](#架构概述)
    - [13.1.4.2 服务划分](#服务划分)
  - [13.1.5 社交网络微服务架构](#社交网络微服务架构)
    - [13.1.5.1 架构概述](#架构概述)
    - [13.1.5.2 服务划分](#服务划分)
  - [13.1.6 总结](#总结)
    - [13.1.6.1 关键要点](#关键要点)
    - [13.1.6.2 技术优势](#技术优势)
    - [13.1.6.3 应用场景](#应用场景)
<!-- TOC END -->

## 13.1.1 目录

1. [概述](#概述)
2. [电商微服务架构](#电商微服务架构)
3. [金融微服务架构](#金融微服务架构)
4. [社交网络微服务架构](#社交网络微服务架构)
5. [物联网微服务架构](#物联网微服务架构)
6. [游戏微服务架构](#游戏微服务架构)
7. [医疗健康微服务架构](#医疗健康微服务架构)
8. [教育微服务架构](#教育微服务架构)
9. [总结](#总结)

## 13.1.2 概述

微服务案例分析通过实际的应用场景，展示微服务架构在不同行业中的实现方法和最佳实践。本分析基于Golang技术栈，提供详细的架构设计、实现方案和性能优化策略。

### 13.1.2.1 核心目标

- **实际应用**: 展示微服务在真实场景中的应用
- **架构设计**: 提供完整的架构设计方案
- **技术实现**: 基于Golang的具体实现方案
- **性能优化**: 针对特定场景的性能优化策略

## 13.1.3 电商微服务架构

### 13.1.3.1 架构概述

电商系统是微服务架构的典型应用场景，需要处理高并发、高可用的业务需求。

**系统特点**:

- 高并发访问
- 复杂的业务流程
- 数据一致性要求
- 实时库存管理
- 支付安全要求

### 13.1.3.2 服务划分

```go
// 电商微服务系统
type EcommerceSystem struct {
    UserService      *UserService
    ProductService   *ProductService
    OrderService     *OrderService
    PaymentService   *PaymentService
    InventoryService *InventoryService
    NotificationService *NotificationService
    SearchService    *SearchService
    RecommendationService *RecommendationService
    APIGateway       *APIGateway
}

// 用户服务
type UserService struct {
    db          *gorm.DB
    cache       *redis.Client
    auth        *AuthService
    logger      *zap.Logger
}

// 用户实体
type User struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    Email       string    `json:"email" gorm:"uniqueIndex"`
    Name        string    `json:"name"`
    Password    string    `json:"-" gorm:"column:password"`
    Phone       string    `json:"phone"`
    Address     string    `json:"address"`
    Status      UserStatus `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// 用户状态
type UserStatus int

const (
    Active UserStatus = iota
    Inactive
    Suspended
)

// 创建用户
func (us *UserService) CreateUser(user *User) error {
    // 验证用户数据
    if err := us.validateUser(user); err != nil {
        return fmt.Errorf("user validation failed: %v", err)
    }
    
    // 检查用户是否已存在
    if exists, _ := us.userExists(user.Email); exists {
        return fmt.Errorf("user already exists")
    }
    
    // 加密密码
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("password encryption failed: %v", err)
    }
    user.Password = string(hashedPassword)
    
    // 保存用户
    if err := us.db.Create(user).Error; err != nil {
        return fmt.Errorf("failed to create user: %v", err)
    }
    
    // 缓存用户信息
    us.cacheUser(user)
    
    // 发布用户创建事件
    us.publishUserCreatedEvent(user)
    
    return nil
}

// 用户认证
func (us *UserService) Authenticate(email, password string) (*User, error) {
    user, err := us.GetUserByEmail(email)
    if err != nil {
        return nil, fmt.Errorf("authentication failed")
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, fmt.Errorf("authentication failed")
    }
    
    return user, nil
}

// 产品服务
type ProductService struct {
    db          *gorm.DB
    cache       *redis.Client
    search      *SearchService
    logger      *zap.Logger
}

// 产品实体
type Product struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Category    string    `json:"category"`
    Brand       string    `json:"brand"`
    Images      []string  `json:"images" gorm:"serializer:json"`
    Attributes  map[string]interface{} `json:"attributes" gorm:"serializer:json"`
    Status      ProductStatus `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// 产品状态
type ProductStatus int

const (
    Available ProductStatus = iota
    OutOfStock
    Discontinued
)

// 创建产品
func (ps *ProductService) CreateProduct(product *Product) error {
    // 验证产品数据
    if err := ps.validateProduct(product); err != nil {
        return fmt.Errorf("product validation failed: %v", err)
    }
    
    // 保存产品
    if err := ps.db.Create(product).Error; err != nil {
        return fmt.Errorf("failed to create product: %v", err)
    }
    
    // 缓存产品信息
    ps.cacheProduct(product)
    
    // 索引产品到搜索引擎
    ps.search.IndexProduct(product)
    
    // 发布产品创建事件
    ps.publishProductCreatedEvent(product)
    
    return nil
}

// 搜索产品
func (ps *ProductService) SearchProducts(query string, filters map[string]interface{}) ([]*Product, error) {
    // 从搜索引擎搜索
    productIDs, err := ps.search.SearchProducts(query, filters)
    if err != nil {
        return nil, fmt.Errorf("search failed: %v", err)
    }
    
    // 从数据库获取产品详情
    products := make([]*Product, 0)
    for _, productID := range productIDs {
        if product, err := ps.GetProductByID(productID); err == nil {
            products = append(products, product)
        }
    }
    
    return products, nil
}

// 订单服务
type OrderService struct {
    db          *gorm.DB
    cache       *redis.Client
    productService *ProductService
    inventoryService *InventoryService
    paymentService *PaymentService
    notificationService *NotificationService
    logger      *zap.Logger
}

// 订单实体
type Order struct {
    ID              string    `json:"id" gorm:"primaryKey"`
    UserID          string    `json:"user_id"`
    Items           []OrderItem `json:"items" gorm:"serializer:json"`
    Total           float64   `json:"total"`
    Status          OrderStatus `json:"status"`
    PaymentMethod   string    `json:"payment_method"`
    PaymentID       string    `json:"payment_id"`
    ShippingAddress string    `json:"shipping_address"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// 订单项
type OrderItem struct {
    ProductID   string  `json:"product_id"`
    Quantity    int     `json:"quantity"`
    Price       float64 `json:"price"`
    Subtotal    float64 `json:"subtotal"`
}

// 订单状态
type OrderStatus int

const (
    Pending OrderStatus = iota
    Confirmed
    Paid
    Shipped
    Delivered
    Cancelled
)

// 创建订单
func (os *OrderService) CreateOrder(order *Order) error {
    // 开始事务
    tx := os.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 验证商品
    for _, item := range order.Items {
        product, err := os.productService.GetProduct(item.ProductID)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("product not found: %v", err)
        }
        
        // 检查库存
        if err := os.inventoryService.ReserveStock(item.ProductID, item.Quantity); err != nil {
            tx.Rollback()
            return fmt.Errorf("insufficient stock: %v", err)
        }
        
        // 计算小计
        item.Price = product.Price
        item.Subtotal = product.Price * float64(item.Quantity)
    }
    
    // 计算订单总额
    total := 0.0
    for _, item := range order.Items {
        total += item.Subtotal
    }
    order.Total = total
    
    // 保存订单
    if err := tx.Create(order).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create order: %v", err)
    }
    
    // 提交事务
    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }
    
    // 异步处理支付
    go os.processPayment(order)
    
    // 发送通知
    go os.notificationService.SendOrderConfirmation(order)
    
    return nil
}

// 处理支付
func (os *OrderService) processPayment(order *Order) {
    // 创建支付请求
    paymentRequest := &PaymentRequest{
        OrderID: order.ID,
        Amount:  order.Total,
        Method:  order.PaymentMethod,
    }
    
    // 调用支付服务
    payment, err := os.paymentService.ProcessPayment(paymentRequest)
    if err != nil {
        // 支付失败，释放库存
        os.releaseStock(order.Items)
        os.notificationService.SendPaymentFailedNotification(order)
        return
    }
    
    // 更新订单状态
    order.Status = Paid
    order.PaymentID = payment.ID
    
    if err := os.db.Save(order).Error; err != nil {
        os.logger.Error("Failed to update order status", zap.Error(err))
    }
    
    // 发送支付成功通知
    os.notificationService.SendPaymentSuccessNotification(order)
}

// 库存服务
type InventoryService struct {
    db          *gorm.DB
    cache       *redis.Client
    logger      *zap.Logger
}

// 库存实体
type Inventory struct {
    ProductID   string    `json:"product_id" gorm:"primaryKey"`
    Quantity    int       `json:"quantity"`
    Reserved    int       `json:"reserved"`
    Available   int       `json:"available"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// 预留库存
func (is *InventoryService) ReserveStock(productID string, quantity int) error {
    // 使用Redis分布式锁
    lock := is.cache.NewMutex(fmt.Sprintf("inventory_lock:%s", productID))
    if err := lock.Lock(); err != nil {
        return fmt.Errorf("failed to acquire lock: %v", err)
    }
    defer lock.Unlock()
    
    // 检查库存
    inventory, err := is.GetInventory(productID)
    if err != nil {
        return fmt.Errorf("inventory not found: %v", err)
    }
    
    if inventory.Available < quantity {
        return fmt.Errorf("insufficient stock")
    }
    
    // 更新库存
    inventory.Reserved += quantity
    inventory.Available -= quantity
    
    if err := is.db.Save(inventory).Error; err != nil {
        return fmt.Errorf("failed to update inventory: %v", err)
    }
    
    // 更新缓存
    is.cacheInventory(inventory)
    
    return nil
}

// 释放库存
func (is *InventoryService) ReleaseStock(productID string, quantity int) error {
    lock := is.cache.NewMutex(fmt.Sprintf("inventory_lock:%s", productID))
    if err := lock.Lock(); err != nil {
        return fmt.Errorf("failed to acquire lock: %v", err)
    }
    defer lock.Unlock()
    
    inventory, err := is.GetInventory(productID)
    if err != nil {
        return fmt.Errorf("inventory not found: %v", err)
    }
    
    inventory.Reserved -= quantity
    inventory.Available += quantity
    
    if err := is.db.Save(inventory).Error; err != nil {
        return fmt.Errorf("failed to update inventory: %v", err)
    }
    
    is.cacheInventory(inventory)
    
    return nil
}

```

## 13.1.4 金融微服务架构

### 13.1.4.1 架构概述

金融系统对安全性、一致性和可靠性有极高要求，微服务架构需要特别考虑这些因素。

**系统特点**:

- 高安全性要求
- 强一致性保证
- 实时交易处理
- 合规性要求
- 审计追踪

### 13.1.4.2 服务划分

```go
// 金融微服务系统
type FinancialSystem struct {
    AccountService     *AccountService
    TransactionService *TransactionService
    PaymentService     *PaymentService
    RiskService        *RiskService
    ComplianceService  *ComplianceService
    AuditService       *AuditService
    NotificationService *NotificationService
    APIGateway         *APIGateway
}

// 账户服务
type AccountService struct {
    db          *gorm.DB
    cache       *redis.Client
    encryption  *EncryptionService
    audit       *AuditService
    logger      *zap.Logger
}

// 账户实体
type Account struct {
    ID              string    `json:"id" gorm:"primaryKey"`
    UserID          string    `json:"user_id"`
    AccountNumber   string    `json:"account_number" gorm:"uniqueIndex"`
    AccountType     AccountType `json:"account_type"`
    Balance         float64   `json:"balance"`
    Currency        string    `json:"currency"`
    Status          AccountStatus `json:"status"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// 账户类型
type AccountType int

const (
    Savings AccountType = iota
    Checking
    Credit
    Investment
)

// 账户状态
type AccountStatus int

const (
    Active AccountStatus = iota
    Frozen
    Closed
)

// 创建账户
func (as *AccountService) CreateAccount(account *Account) error {
    // 验证账户数据
    if err := as.validateAccount(account); err != nil {
        return fmt.Errorf("account validation failed: %v", err)
    }
    
    // 生成账户号码
    account.AccountNumber = as.generateAccountNumber()
    
    // 加密敏感数据
    if err := as.encryptAccountData(account); err != nil {
        return fmt.Errorf("encryption failed: %v", err)
    }
    
    // 保存账户
    if err := as.db.Create(account).Error; err != nil {
        return fmt.Errorf("failed to create account: %v", err)
    }
    
    // 审计日志
    as.audit.LogAccountCreation(account)
    
    // 缓存账户信息
    as.cacheAccount(account)
    
    return nil
}

// 更新余额
func (as *AccountService) UpdateBalance(accountID string, amount float64) error {
    // 使用数据库事务
    tx := as.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 锁定账户记录
    var account Account
    if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", accountID).First(&account).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("account not found: %v", err)
    }
    
    // 更新余额
    account.Balance += amount
    
    if err := tx.Save(&account).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update balance: %v", err)
    }
    
    // 提交事务
    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }
    
    // 审计日志
    as.audit.LogBalanceUpdate(accountID, amount)
    
    // 更新缓存
    as.cacheAccount(&account)
    
    return nil
}

// 交易服务
type TransactionService struct {
    db          *gorm.DB
    cache       *redis.Client
    accountService *AccountService
    riskService *RiskService
    complianceService *ComplianceService
    auditService *AuditService
    logger      *zap.Logger
}

// 交易实体
type Transaction struct {
    ID              string    `json:"id" gorm:"primaryKey"`
    FromAccountID   string    `json:"from_account_id"`
    ToAccountID     string    `json:"to_account_id"`
    Amount          float64   `json:"amount"`
    Currency        string    `json:"currency"`
    Type            TransactionType `json:"type"`
    Status          TransactionStatus `json:"status"`
    Description     string    `json:"description"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// 交易类型
type TransactionType int

const (
    Transfer TransactionType = iota
    Deposit
    Withdrawal
    Payment
)

// 交易状态
type TransactionStatus int

const (
    Pending TransactionStatus = iota
    Processing
    Completed
    Failed
    Cancelled
)

// 执行交易
func (ts *TransactionService) ExecuteTransaction(transaction *Transaction) error {
    // 风险检查
    if err := ts.riskService.CheckTransactionRisk(transaction); err != nil {
        return fmt.Errorf("risk check failed: %v", err)
    }
    
    // 合规检查
    if err := ts.complianceService.CheckCompliance(transaction); err != nil {
        return fmt.Errorf("compliance check failed: %v", err)
    }
    
    // 开始分布式事务
    tx := ts.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 锁定源账户
    var fromAccount Account
    if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", transaction.FromAccountID).First(&fromAccount).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("source account not found: %v", err)
    }
    
    // 检查余额
    if fromAccount.Balance < transaction.Amount {
        tx.Rollback()
        return fmt.Errorf("insufficient balance")
    }
    
    // 锁定目标账户
    var toAccount Account
    if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", transaction.ToAccountID).First(&toAccount).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("target account not found: %v", err)
    }
    
    // 更新余额
    fromAccount.Balance -= transaction.Amount
    toAccount.Balance += transaction.Amount
    
    if err := tx.Save(&fromAccount).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update source account: %v", err)
    }
    
    if err := tx.Save(&toAccount).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update target account: %v", err)
    }
    
    // 保存交易记录
    transaction.Status = Completed
    if err := tx.Create(transaction).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create transaction: %v", err)
    }
    
    // 提交事务
    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }
    
    // 审计日志
    ts.auditService.LogTransaction(transaction)
    
    // 更新缓存
    ts.cacheTransaction(transaction)
    
    return nil
}

// 风险服务
type RiskService struct {
    rules       map[string]*RiskRule
    thresholds  map[string]float64
    logger      *zap.Logger
}

// 风险规则
type RiskRule struct {
    ID          string
    Name        string
    Condition   func(*Transaction) bool
    Action      func(*Transaction) error
    Priority    int
}

// 检查交易风险
func (rs *RiskService) CheckTransactionRisk(transaction *Transaction) error {
    // 按优先级检查风险规则
    rules := rs.getRulesByPriority()
    
    for _, rule := range rules {
        if rule.Condition(transaction) {
            if err := rule.Action(transaction); err != nil {
                return fmt.Errorf("risk rule %s failed: %v", rule.Name, err)
            }
        }
    }
    
    return nil
}

// 大额交易检查
func (rs *RiskService) CheckLargeTransaction(transaction *Transaction) error {
    threshold := rs.thresholds["large_transaction"]
    
    if transaction.Amount > threshold {
        // 触发大额交易审核
        rs.logger.Warn("Large transaction detected", 
            zap.String("transaction_id", transaction.ID),
            zap.Float64("amount", transaction.Amount))
        
        return fmt.Errorf("large transaction requires manual review")
    }
    
    return nil
}

// 合规服务
type ComplianceService struct {
    rules       map[string]*ComplianceRule
    logger      *zap.Logger
}

// 合规规则
type ComplianceRule struct {
    ID          string
    Name        string
    Check       func(*Transaction) error
}

// 检查合规性
func (cs *ComplianceService) CheckCompliance(transaction *Transaction) error {
    for _, rule := range cs.rules {
        if err := rule.Check(transaction); err != nil {
            return fmt.Errorf("compliance rule %s failed: %v", rule.Name, err)
        }
    }
    
    return nil
}

// 反洗钱检查
func (cs *ComplianceService) CheckAML(transaction *Transaction) error {
    // 实现反洗钱检查逻辑
    // 检查交易频率、金额模式等
    
    return nil
}

// 审计服务
type AuditService struct {
    db          *gorm.DB
    logger      *zap.Logger
}

// 审计日志
type AuditLog struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    UserID      string    `json:"user_id"`
    Action      string    `json:"action"`
    Resource    string    `json:"resource"`
    ResourceID  string    `json:"resource_id"`
    Details     map[string]interface{} `json:"details" gorm:"serializer:json"`
    IPAddress   string    `json:"ip_address"`
    UserAgent   string    `json:"user_agent"`
    CreatedAt   time.Time `json:"created_at"`
}

// 记录审计日志
func (as *AuditService) LogAction(userID, action, resource, resourceID string, details map[string]interface{}) error {
    log := &AuditLog{
        ID:         uuid.New().String(),
        UserID:     userID,
        Action:     action,
        Resource:   resource,
        ResourceID: resourceID,
        Details:    details,
        CreatedAt:  time.Now(),
    }
    
    if err := as.db.Create(log).Error; err != nil {
        return fmt.Errorf("failed to create audit log: %v", err)
    }
    
    return nil
}

```

## 13.1.5 社交网络微服务架构

### 13.1.5.1 架构概述

社交网络系统需要处理大量的用户交互、内容分发和实时通信，微服务架构能够很好地支持这些需求。

**系统特点**:

- 高并发用户访问
- 实时消息推送
- 内容推荐算法
- 用户关系管理
- 多媒体内容处理

### 13.1.5.2 服务划分

```go
// 社交网络微服务系统
type SocialNetworkSystem struct {
    UserService       *UserService
    PostService       *PostService
    CommentService    *CommentService
    LikeService       *LikeService
    FollowService     *FollowService
    MessageService    *MessageService
    NotificationService *NotificationService
    RecommendationService *RecommendationService
    MediaService      *MediaService
    APIGateway        *APIGateway
}

// 用户服务
type UserService struct {
    db          *gorm.DB
    cache       *redis.Client
    search      *SearchService
    logger      *zap.Logger
}

// 用户实体
type User struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    Username    string    `json:"username" gorm:"uniqueIndex"`
    Email       string    `json:"email" gorm:"uniqueIndex"`
    Name        string    `json:"name"`
    Bio         string    `json:"bio"`
    Avatar      string    `json:"avatar"`
    Followers   int       `json:"followers"`
    Following   int       `json:"following"`
    Posts       int       `json:"posts"`
    Status      UserStatus `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// 创建用户
func (us *UserService) CreateUser(user *User) error {
    // 验证用户数据
    if err := us.validateUser(user); err != nil {
        return fmt.Errorf("user validation failed: %v", err)
    }
    
    // 检查用户名是否已存在
    if exists, _ := us.usernameExists(user.Username); exists {
        return fmt.Errorf("username already exists")
    }
    
    // 保存用户
    if err := us.db.Create(user).Error; err != nil {
        return fmt.Errorf("failed to create user: %v", err)
    }
    
    // 索引用户到搜索引擎
    us.search.IndexUser(user)
    
    // 缓存用户信息
    us.cacheUser(user)
    
    return nil
}

// 帖子服务
type PostService struct {
    db          *gorm.DB
    cache       *redis.Client
    search      *SearchService
    media       *MediaService
    logger      *zap.Logger
}

// 帖子实体
type Post struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    UserID      string    `json:"user_id"`
    Content     string    `json:"content"`
    Media       []string  `json:"media" gorm:"serializer:json"`
    Likes       int       `json:"likes"`
    Comments    int       `json:"comments"`
    Shares      int       `json:"shares"`
    Visibility  Visibility `json:"visibility"`
    Status      PostStatus `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// 可见性
type Visibility int

const (
    Public Visibility = iota
    Friends
    Private
)

// 帖子状态
type PostStatus int

const (
    Published PostStatus = iota
    Draft
    Deleted
)

// 创建帖子
func (ps *PostService) CreatePost(post *Post) error {
    // 验证帖子数据
    if err := ps.validatePost(post); err != nil {
        return fmt.Errorf("post validation failed: %v", err)
    }
    
    // 处理媒体文件
    if len(post.Media) > 0 {
        processedMedia := make([]string, 0)
        for _, mediaURL := range post.Media {
            if processedURL, err := ps.media.ProcessMedia(mediaURL); err == nil {
                processedMedia = append(processedMedia, processedURL)
            }
        }
        post.Media = processedMedia
    }
    
    // 保存帖子
    if err := ps.db.Create(post).Error; err != nil {
        return fmt.Errorf("failed to create post: %v", err)
    }
    
    // 索引帖子到搜索引擎
    ps.search.IndexPost(post)
    
    // 缓存帖子信息
    ps.cachePost(post)
    
    // 更新用户帖子计数
    ps.updateUserPostCount(post.UserID)
    
    return nil
}

// 获取用户时间线
func (ps *PostService) GetUserTimeline(userID string, page, size int) ([]*Post, error) {
    // 获取用户关注的人
    following, err := ps.getFollowingUsers(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get following users: %v", err)
    }
    
    // 查询帖子
    var posts []*Post
    offset := (page - 1) * size
    
    query := ps.db.Where("user_id IN (?) AND visibility = ?", following, Public)
    query = query.Order("created_at DESC").Offset(offset).Limit(size)
    
    if err := query.Find(&posts).Error; err != nil {
        return nil, fmt.Errorf("failed to get posts: %v", err)
    }
    
    return posts, nil
}

// 关注服务
type FollowService struct {
    db          *gorm.DB
    cache       *redis.Client
    notification *NotificationService
    logger      *zap.Logger
}

// 关注关系
type Follow struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    FollowerID  string    `json:"follower_id"`
    FollowingID string    `json:"following_id"`
    CreatedAt   time.Time `json:"created_at"`
}

// 关注用户
func (fs *FollowService) FollowUser(followerID, followingID string) error {
    // 检查是否已经关注
    if exists, _ := fs.isFollowing(followerID, followingID); exists {
        return fmt.Errorf("already following")
    }
    
    // 创建关注关系
    follow := &Follow{
        ID:          uuid.New().String(),
        FollowerID:  followerID,
        FollowingID: followingID,
        CreatedAt:   time.Now(),
    }
    
    if err := fs.db.Create(follow).Error; err != nil {
        return fmt.Errorf("failed to create follow: %v", err)
    }
    
    // 更新关注者计数
    fs.updateFollowerCount(followingID, 1)
    fs.updateFollowingCount(followerID, 1)
    
    // 发送通知
    fs.notification.SendFollowNotification(followerID, followingID)
    
    return nil
}

// 取消关注
func (fs *FollowService) UnfollowUser(followerID, followingID string) error {
    // 删除关注关系
    if err := fs.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&Follow{}).Error; err != nil {
        return fmt.Errorf("failed to delete follow: %v", err)
    }
    
    // 更新关注者计数
    fs.updateFollowerCount(followingID, -1)
    fs.updateFollowingCount(followerID, -1)
    
    return nil
}

// 消息服务
type MessageService struct {
    db          *gorm.DB
    cache       *redis.Client
    websocket   *WebSocketManager
    logger      *zap.Logger
}

// 消息实体
type Message struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    SenderID    string    `json:"sender_id"`
    ReceiverID  string    `json:"receiver_id"`
    Content     string    `json:"content"`
    Type        MessageType `json:"type"`
    Status      MessageStatus `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

// 消息类型
type MessageType int

const (
    Text MessageType = iota
    Image
    Video
    File
)

// 消息状态
type MessageStatus int

const (
    Sent MessageStatus = iota
    Delivered
    Read
)

// 发送消息
func (ms *MessageService) SendMessage(message *Message) error {
    // 验证消息数据
    if err := ms.validateMessage(message); err != nil {
        return fmt.Errorf("message validation failed: %v", err)
    }
    
    // 保存消息
    if err := ms.db.Create(message).Error; err != nil {
        return fmt.Errorf("failed to create message: %v", err)
    }
    
    // 缓存消息
    ms.cacheMessage(message)
    
    // 通过WebSocket发送实时消息
    ms.websocket.SendToUser(message.ReceiverID, "new_message", message)
    
    return nil
}

// 获取对话历史
func (ms *MessageService) GetConversation(user1ID, user2ID string, page, size int) ([]*Message, error) {
    var messages []*Message
    offset := (page - 1) * size
    
    query := ms.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", 
        user1ID, user2ID, user2ID, user1ID)
    query = query.Order("created_at DESC").Offset(offset).Limit(size)
    
    if err := query.Find(&messages).Error; err != nil {
        return nil, fmt.Errorf("failed to get messages: %v", err)
    }
    
    return messages, nil
}

// 推荐服务
type RecommendationService struct {
    db          *gorm.DB
    cache       *redis.Client
    algorithm   *RecommendationAlgorithm
    logger      *zap.Logger
}

// 推荐算法
type RecommendationAlgorithm struct {
    models      map[string]*RecommendationModel
}

// 推荐模型
type RecommendationModel struct {
    Name        string
    Algorithm   func(string) ([]string, error)
    Weight      float64
}

// 获取推荐用户
func (rs *RecommendationService) GetRecommendedUsers(userID string, limit int) ([]*User, error) {
    // 使用多种推荐算法
    recommendations := make(map[string]float64)
    
    for _, model := range rs.algorithm.models {
        if userIDs, err := model.Algorithm(userID); err == nil {
            for _, id := range userIDs {
                recommendations[id] += model.Weight
            }
        }
    }
    
    // 排序并获取前N个推荐
    sortedUsers := rs.sortByScore(recommendations)
    
    // 获取用户详情
    users := make([]*User, 0)
    for i := 0; i < limit && i < len(sortedUsers); i++ {
        if user, err := rs.getUserByID(sortedUsers[i]); err == nil {
            users = append(users, user)
        }
    }
    
    return users, nil
}

// 协同过滤算法
func (rs *RecommendationService) CollaborativeFiltering(userID string) ([]string, error) {
    // 获取用户的关注列表
    following, err := rs.getFollowingUsers(userID)
    if err != nil {
        return nil, err
    }
    
    // 获取关注用户的关注列表
    recommendations := make(map[string]int)
    for _, followingID := range following {
        if followingFollowing, err := rs.getFollowingUsers(followingID); err == nil {
            for _, id := range followingFollowing {
                if id != userID {
                    recommendations[id]++
                }
            }
        }
    }
    
    // 排序并返回推荐用户ID
    return rs.sortByCount(recommendations), nil
}

```

## 13.1.6 总结

微服务案例分析展示了不同行业中微服务架构的实际应用，每个案例都有其特定的需求和挑战。

### 13.1.6.1 关键要点

1. **电商系统**: 高并发、数据一致性、库存管理
2. **金融系统**: 安全性、强一致性、合规性
3. **社交网络**: 实时通信、内容分发、推荐算法

### 13.1.6.2 技术优势

- **可扩展性**: 支持业务快速增长
- **高可用性**: 通过服务隔离提高系统稳定性
- **技术多样性**: 支持不同技术栈的选择
- **团队自治**: 支持团队独立开发和部署

### 13.1.6.3 应用场景

- **大型系统**: 复杂业务系统的模块化
- **高并发系统**: 需要高并发处理的应用
- **实时系统**: 需要实时响应的系统
- **多团队开发**: 支持团队独立开发

通过合理应用微服务架构，可以构建出更加灵活、可扩展和可维护的软件系统。
