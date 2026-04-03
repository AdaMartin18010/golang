# AD-030-E-commerce-Systems-2026

> **Dimension**: 05-Application-Domains
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: E-commerce 2026 (Marketplace, OMS, Inventory, Payments)
> **Size**: >20KB

---

## 1. 电商系统概览

### 1.1 核心领域

| 领域 | 关键功能 | 技术挑战 |
|------|---------|---------|
| 商品 | SKU管理、类目、搜索 | 大规模数据、实时同步 |
| 库存 | 分配、预留、同步 | 超卖、分布式一致性 |
| 订单 | 生命周期、状态机 | 幂等性、事务 |
| 支付 | 多渠道、对账 | 安全、合规 |
| 物流 | 履约、追踪 | 多方集成 |
| 营销 | 促销、优惠券 | 并发控制 |
| 用户 | 画像、推荐 | 隐私合规 |

### 1.2 系统架构

```
┌─────────────────────────────────────────┐
│           电商系统架构                  │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐ │
│  │   Web   │  │ Mobile  │  │  Admin  │ │
│  │   App   │  │   App   │  │  Portal │ │
│  └────┬────┘  └────┬────┘  └────┬────┘ │
│       │            │            │      │
│       └────────────┼────────────┘      │
│                    │                    │
│              API Gateway               │
│            (Kong/AWS API GW)           │
│                    │                   │
│  ┌─────────────────┼─────────────────┐ │
│  │           BFF层                  │ │
│  │  (Next.js / React Query)         │ │
│  └─────────────────┼─────────────────┘ │
│                    │                   │
│  ┌─────────────────┼─────────────────┐ │
│  │           微服务层                │ │
│  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐│ │
│  │  │商品 │ │库存 │ │订单 │ │支付 ││ │
│  │  │服务 │ │服务 │ │服务 │ │服务 ││ │
│  │  └─────┘ └─────┘ └─────┘ └─────┘│ │
│  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐│ │
│  │  │用户 │ │搜索 │ │推荐 │ │物流 ││ │
│  │  │服务 │ │服务 │ │服务 │ │服务 ││ │
│  │  └─────┘ └─────┘ └─────┘ └─────┘│ │
│  └─────────────────┼─────────────────┘ │
│                    │                   │
│  ┌─────────────────┼─────────────────┐ │
│  │           数据层                 │ │
│  │  MySQL  Redis  ES  Kafka  OSS   │ │
│  └───────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

---

## 2. 商品系统

### 2.1 SKU数据模型

```go
// SKU (库存单位)
type SKU struct {
    ID          string            `json:"id"`           // SKU ID
    SPUID       string            `json:"spu_id"`       // 产品ID
    Attributes  map[string]string `json:"attributes"`   // {颜色: 红色, 尺码: XL}
    Price       decimal.Decimal   `json:"price"`        // 售价
    Cost        decimal.Decimal   `json:"cost"`         // 成本
    Barcode     string            `json:"barcode"`      // 条形码
    Weight      float64           `json:"weight"`       // 重量(kg)
    Images      []string          `json:"images"`       // 图片URL
    Status      SKUStatus         `json:"status"`       // 上架/下架
    CreatedAt   time.Time         `json:"created_at"`
}

// SPU (标准化产品单元)
type SPU struct {
    ID           string     `json:"id"`
    Name         string     `json:"name"`
    Description  string     `json:"description"`
    CategoryID   string     `json:"category_id"`
    BrandID      string     `json:"brand_id"`
    MainImage    string     `json:"main_image"`
    AttributeDefs []AttributeDef `json:"attribute_defs"`
}

// 类目
type Category struct {
    ID       string   `json:"id"`
    ParentID string   `json:"parent_id"`
    Name     string   `json:"name"`
    Level    int      `json:"level"`
    Path     string   `json:"path"`      // 1/2/3
    IsLeaf   bool     `json:"is_leaf"`
    Attributes []CategoryAttribute `json:"attributes"`
}
```

### 2.2 搜索系统

```go
// Elasticsearch商品索引
type ProductDocument struct {
    SKUId       string            `json:"sku_id"`
    SPUId       string            `json:"spu_id"`
    Title       string            `json:"title"`
    Description string            `json:"description"`
    CategoryID  string            `json:"category_id"`
    CategoryPath string           `json:"category_path"`
    Price       float64           `json:"price"`
    Attributes  map[string]string `json:"attributes"`
    Tags        []string          `json:"tags"`
    Sales       int64             `json:"sales"`
    Rating      float64           `json:"rating"`
    InStock     bool              `json:"in_stock"`
    CreatedAt   time.Time         `json:"created_at"`

    // 向量字段 (语义搜索)
    TitleVector []float32 `json:"title_vector"`
}

// 搜索服务
type SearchService struct {
    es *elasticsearch.Client
}

func (s *SearchService) Search(ctx context.Context, req SearchRequest) (*SearchResult, error) {
    query := map[string]interface{}{
        "bool": map[string]interface{}{
            "must": []map[string]interface{}{
                {
                    "multi_match": map[string]interface{}{
                        "query":  req.Keyword,
                        "fields": []string{"title^3", "description", "brand"},
                    },
                },
            },
            "filter": []map[string]interface{}{},
        },
    }

    // 价格过滤
    if req.MinPrice > 0 || req.MaxPrice > 0 {
        priceRange := map[string]interface{}{}
        if req.MinPrice > 0 {
            priceRange["gte"] = req.MinPrice
        }
        if req.MaxPrice > 0 {
            priceRange["lte"] = req.MaxPrice
        }
        query["bool"]["filter"] = append(
            query["bool"]["filter"].([]map[string]interface{}),
            map[string]interface{}{"range": map[string]interface{}{"price": priceRange}},
        )
    }

    // 类目过滤
    if req.CategoryID != "" {
        query["bool"]["filter"] = append(
            query["bool"]["filter"].([]map[string]interface{}),
            map[string]interface{}{"term": map[string]interface{}{"category_id": req.CategoryID}},
        )
    }

    // 属性过滤
    for attr, value := range req.Filters {
        query["bool"]["filter"] = append(
            query["bool"]["filter"].([]map[string]interface{}),
            map[string]interface{}{"term": map[string]interface{}{fmt.Sprintf("attributes.%s", attr): value}},
        )
    }

    // 排序
    sort := []map[string]interface{}{}
    switch req.SortBy {
    case "price_asc":
        sort = append(sort, map[string]interface{}{"price": "asc"})
    case "price_desc":
        sort = append(sort, map[string]interface{}{"price": "desc"})
    case "sales":
        sort = append(sort, map[string]interface{}{"sales": "desc"})
    case "rating":
        sort = append(sort, map[string]interface{}{"rating": "desc"})
    default:
        sort = append(sort, map[string]interface{}{"_score": "desc"})
    }

    searchBody := map[string]interface{}{
        "query": query,
        "sort":  sort,
        "from":  (req.Page - 1) * req.PageSize,
        "size":  req.PageSize,
        "aggs": map[string]interface{}{
            "categories": map[string]interface{}{
                "terms": map[string]interface{}{"field": "category_id"},
            },
            "price_stats": map[string]interface{}{
                "stats": map[string]interface{}{"field": "price"},
            },
        },
    }

    return s.executeSearch(ctx, searchBody)
}

// 语义搜索 (向量相似度)
func (s *SearchService) SemanticSearch(ctx context.Context, query string) (*SearchResult, error) {
    // 1. 查询向量化
    vector := s.embeddingModel.Encode(query)

    // 2. KNN搜索
    knnQuery := map[string]interface{}{
        "knn": map[string]interface{}{
            "field": "title_vector",
            "query_vector": vector,
            "k": 100,
            "num_candidates": 1000,
        },
    }

    return s.executeSearch(ctx, knnQuery)
}
```

---

## 3. 库存系统

### 3.1 库存模型

```go
// 库存分层
type Inventory struct {
    SKUId         string          `json:"sku_id"`
    WarehouseID   string          `json:"warehouse_id"`

    // 物理库存
    TotalQty      int64           `json:"total_qty"`       // 总库存
    AvailableQty  int64           `json:"available_qty"`   // 可用库存
    ReservedQty   int64           `json:"reserved_qty"`    // 已预留
    LockedQty     int64           `json:"locked_qty"`      // 锁定(售后等)

    // 渠道库存
    Channels      map[string]ChannelInventory `json:"channels"`

    Version       int64           `json:"version"`         // 乐观锁版本
    UpdatedAt     time.Time       `json:"updated_at"`
}

type ChannelInventory struct {
    ChannelID    string `json:"channel_id"`
    AvailableQty int64  `json:"available_qty"`
    ReservedQty  int64  `json:"reserved_qty"`
}

// 库存操作
type InventoryOperation struct {
    ID          string    `json:"id"`
    SKUId       string    `json:"sku_id"`
    WarehouseID string    `json:"warehouse_id"`
    Type        OpType    `json:"type"`       // RESERVE, RELEASE, DEDUCT, INCREASE
    Quantity    int64     `json:"quantity"`
    OrderID     string    `json:"order_id"`
    Reason      string    `json:"reason"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### 3.2 库存预留模式

```go
// 库存服务
type InventoryService struct {
    redis    *redis.Client
    mysql    *sql.DB
    kafka    sarama.AsyncProducer
}

// 预留库存 (下单时)
func (s *InventoryService) Reserve(ctx context.Context, req ReserveRequest) error {
    key := fmt.Sprintf("inventory:%s:%s", req.SKUId, req.WarehouseID)

    // 1. Redis Lua脚本原子操作
    script := `
        local available = tonumber(redis.call('hget', KEYS[1], 'available'))
        local reserve = tonumber(redis.call('hget', KEYS[1], 'reserved'))
        local qty = tonumber(ARGV[1])

        if available >= qty then
            redis.call('hincrby', KEYS[1], 'available', -qty)
            redis.call('hincrby', KEYS[1], 'reserved', qty)
            redis.call('hset', KEYS[1], 'reserve_' .. ARGV[2], qty)
            return 1
        else
            return 0
        end
    `

    result, err := s.redis.Eval(ctx, script, []string{key}, req.Quantity, req.OrderID).Result()
    if err != nil || result.(int64) == 0 {
        return ErrInsufficientInventory
    }

    // 2. 发送库存变更事件
    s.kafka.Send(&sarama.ProducerMessage{
        Topic: "inventory_changes",
        Value: sarama.StringEncoder(fmt.Sprintf(`{"sku":"%s","op":"reserve","qty":%d}`, req.SKUId, req.Quantity)),
    })

    // 3. 异步同步到MySQL
    return nil
}

// 扣减库存 (支付成功)
func (s *InventoryService) Deduct(ctx context.Context, req DeductRequest) error {
    key := fmt.Sprintf("inventory:%s:%s", req.SKUId, req.WarehouseID)

    script := `
        local reserved = tonumber(redis.call('hget', KEYS[1], 'reserved'))
        local reserve_key = 'reserve_' .. ARGV[1]
        local reserved_qty = tonumber(redis.call('hget', KEYS[1], reserve_key) or 0)

        if reserved_qty > 0 then
            redis.call('hincrby', KEYS[1], 'reserved', -reserved_qty)
            redis.call('hincrby', KEYS[1], 'total', -reserved_qty)
            redis.call('hdel', KEYS[1], reserve_key)
            return 1
        end
        return 0
    `

    _, err := s.redis.Eval(ctx, script, []string{key}, req.OrderID).Result()
    return err
}

// 释放库存 (取消订单)
func (s *InventoryService) Release(ctx context.Context, req ReleaseRequest) error {
    key := fmt.Sprintf("inventory:%s:%s", req.SKUId, req.WarehouseID)

    script := `
        local reserve_key = 'reserve_' .. ARGV[1]
        local reserved_qty = tonumber(redis.call('hget', KEYS[1], reserve_key) or 0)

        if reserved_qty > 0 then
            redis.call('hincrby', KEYS[1], 'available', reserved_qty)
            redis.call('hincrby', KEYS[1], 'reserved', -reserved_qty)
            redis.call('hdel', KEYS[1], reserve_key)
            return reserved_qty
        end
        return 0
    `

    _, err := s.redis.Eval(ctx, script, []string{key}, req.OrderID).Result()
    return err
}

// 定时释放过期预留
func (s *InventoryService) ExpireReservations(ctx context.Context) {
    // 扫描所有预留
    // 如果预留超过15分钟未支付，自动释放
}
```

---

## 4. 订单系统

### 4.1 订单状态机

```
┌─────────┐   创建   ┌─────────┐   支付   ┌─────────┐
│  INIT   │ ───────►│ CREATED │────────►│  PAID   │
└─────────┘         └────┬────┘         └────┬────┘
                         │                   │
                    取消  ▼              发货  ▼
                   ┌─────────┐          ┌─────────┐
                   │CANCELLED│          │SHIPPED  │
                   └─────────┘          └────┬────┘
                                             │
                                        签收  ▼
                                       ┌─────────┐
                                       │COMPLETED│
                                       └─────────┘
```

### 4.2 订单模型

```go
type Order struct {
    ID              string      `json:"id"`
    UserID          string      `json:"user_id"`
    Status          OrderStatus `json:"status"`

    // 金额
    TotalAmount     decimal.Decimal `json:"total_amount"`
    DiscountAmount  decimal.Decimal `json:"discount_amount"`
    ShippingFee     decimal.Decimal `json:"shipping_fee"`
    PayAmount       decimal.Decimal `json:"pay_amount"`

    // 商品
    Items           []OrderItem `json:"items"`

    // 地址
    ShippingAddress Address     `json:"shipping_address"`

    // 支付
    Payment         PaymentInfo `json:"payment"`

    // 物流
    Shipment        *ShipmentInfo `json:"shipment,omitempty"`

    // 时间戳
    CreatedAt       time.Time   `json:"created_at"`
    PaidAt          *time.Time  `json:"paid_at,omitempty"`
    ShippedAt       *time.Time  `json:"shipped_at,omitempty"`
    CompletedAt     *time.Time  `json:"completed_at,omitempty"`
}

type OrderItem struct {
    SKUId       string          `json:"sku_id"`
    SPUName     string          `json:"spu_name"`
    SKUAttrs    map[string]string `json:"sku_attrs"`
    Image       string          `json:"image"`
    Price       decimal.Decimal `json:"price"`
    Quantity    int             `json:"quantity"`
    TotalPrice  decimal.Decimal `json:"total_price"`
}
```

### 4.3 订单创建流程

```go
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    // 1. 参数校验
    if len(req.Items) == 0 {
        return nil, ErrEmptyCart
    }

    // 2. 获取商品信息并计算价格
    items := make([]OrderItem, 0, len(req.Items))
    var totalAmount decimal.Decimal

    for _, cartItem := range req.Items {
        sku, err := s.productService.GetSKU(ctx, cartItem.SKUId)
        if err != nil {
            return nil, err
        }

        if sku.Status != SKUStatusActive {
            return nil, ErrSKUUnavailable
        }

        item := OrderItem{
            SKUId:      cartItem.SKUId,
            SPUName:    sku.SPUName,
            SKUAttrs:   sku.Attributes,
            Image:      sku.Images[0],
            Price:      sku.Price,
            Quantity:   cartItem.Quantity,
            TotalPrice: sku.Price.Mul(decimal.NewFromInt(int64(cartItem.Quantity))),
        }
        items = append(items, item)
        totalAmount = totalAmount.Add(item.TotalPrice)
    }

    // 3. 应用优惠
    discountAmount := s.calculateDiscount(ctx, req, totalAmount)

    // 4. 计算运费
    shippingFee := s.calculateShipping(ctx, req.ShippingAddress, items)

    // 5. 创建订单(预扣库存)
    order := &Order{
        ID:             generateOrderID(),
        UserID:         req.UserID,
        Status:         OrderStatusCreated,
        TotalAmount:    totalAmount,
        DiscountAmount: discountAmount,
        ShippingFee:    shippingFee,
        PayAmount:      totalAmount.Sub(discountAmount).Add(shippingFee),
        Items:          items,
        ShippingAddress: req.ShippingAddress,
        CreatedAt:      time.Now(),
    }

    // 6. 预留库存
    for _, item := range items {
        if err := s.inventoryService.Reserve(ctx, ReserveRequest{
            SKUId:       item.SKUId,
            Quantity:    item.Quantity,
            OrderID:     order.ID,
        }); err != nil {
            // 回滚已预留的库存
            s.rollbackReservation(ctx, order.ID, items[:len(items)-1])
            return nil, err
        }
    }

    // 7. 保存订单
    if err := s.orderRepo.Create(ctx, order); err != nil {
        s.rollbackReservation(ctx, order.ID, items)
        return nil, err
    }

    // 8. 发送延迟消息(15分钟未支付自动取消)
    s.delayQueue.Publish(DelayMessage{
        Topic:     "order_auto_cancel",
        Key:       order.ID,
        Payload:   order.ID,
        Delay:     15 * time.Minute,
    })

    return order, nil
}
```

---

## 5. 促销系统

### 5.1 促销类型

```go
type PromotionType int

const (
    PromotionTypeDirectDiscount PromotionType = iota // 直降
    PromotionTypePercentage                           // 折扣率
    PromotionTypeFixedAmount                          // 固定金额
    PromotionTypeBuyXGetY                             // 买X送Y
    PromotionTypeBundle                               // 捆绑销售
    PromotionTypeFlashSale                            // 秒杀
)

type Promotion struct {
    ID          string        `json:"id"`
    Name        string        `json:"name"`
    Type        PromotionType `json:"type"`
    Rules       PromotionRules `json:"rules"`

    // 时间
    StartTime   time.Time     `json:"start_time"`
    EndTime     time.Time     `json:"end_time"`

    // 范围
    Scope       PromotionScope `json:"scope"`

    // 限制
    UserLimit   int           `json:"user_limit"`   // 每用户限领
    TotalLimit  int           `json:"total_limit"`  // 总限量
    UsedCount   int           `json:"used_count"`

    Status      PromotionStatus `json:"status"`
}

// 秒杀特殊处理
type FlashSale struct {
    Promotion
    SKUId       string    `json:"sku_id"`
    Stock       int64     `json:"stock"`
    SalePrice   decimal.Decimal `json:"sale_price"`

    // 限流
    RateLimit   int       `json:"rate_limit"`  // 每秒请求数
}
```

### 5.2 秒杀实现

```go
func (s *FlashSaleService) SecKill(ctx context.Context, req SecKillRequest) (*Order, error) {
    // 1. 限流检查
    if !s.limiter.Allow(req.UserID) {
        return nil, ErrRateLimited
    }

    // 2. 用户资格检查 (是否已购买)
    if s.hasBought(ctx, req.UserID, req.FlashSaleID) {
        return nil, ErrAlreadyPurchased
    }

    // 3. Redis原子扣减库存
    key := fmt.Sprintf("flashsale:%s:stock", req.FlashSaleID)
    remaining, err := s.redis.Decr(ctx, key).Result()
    if err != nil || remaining < 0 {
        // 库存不足，恢复库存计数
        s.redis.Incr(ctx, key)
        return nil, ErrOutOfStock
    }

    // 4. 发送订单创建消息到队列 (异步处理)
    msg := SecKillMessage{
        FlashSaleID: req.FlashSaleID,
        UserID:      req.UserID,
        SKUId:       req.SKUId,
        Quantity:    1,
    }

    if err := s.orderQueue.Publish(msg); err != nil {
        // 恢复库存
        s.redis.Incr(ctx, key)
        return nil, err
    }

    // 5. 标记用户已参与
    s.markBought(ctx, req.UserID, req.FlashSaleID)

    // 6. 返回排队中状态
    return &Order{Status: OrderStatusQueued}, nil
}
```

---

## 6. 推荐系统

### 6.1 推荐架构

```go
// 多路召回 + 排序
type RecommendationService struct {
    // 召回层
    recallers []Recaller

    // 排序模型
    ranker *Ranker
}

type Recaller interface {
    Recall(ctx context.Context, userID string, count int) ([]Item, error)
}

// 实现多种召回策略
type CFRecaller struct { /* 协同过滤 */ }
type ContentRecaller struct { /* 内容相似 */ }
type HotRecaller struct { /* 热门推荐 */ }
type NewArrivalRecaller struct { /* 新品推荐 */ }

func (s *RecommendationService) Recommend(ctx context.Context, req RecommendRequest) ([]Item, error) {
    // 1. 多路召回
    recallCh := make(chan []Item, len(s.recallers))

    for _, recaller := range s.recallers {
        go func(r Recaller) {
            items, _ := r.Recall(ctx, req.UserID, req.Count*3)
            recallCh <- items
        }(recaller)
    }

    // 合并召回结果
    candidateMap := make(map[string]Item)
    for i := 0; i < len(s.recallers); i++ {
        items := <-recallCh
        for _, item := range items {
            candidateMap[item.ID] = item
        }
    }

    candidates := make([]Item, 0, len(candidateMap))
    for _, item := range candidateMap {
        candidates = append(candidates, item)
    }

    // 2. 粗排 (轻量级模型，快速过滤)
    candidates = s.coarseRank(ctx, candidates, req.Count*2)

    // 3. 精排 (深度模型)
    ranked := s.ranker.Rank(ctx, req.UserID, candidates)

    // 4. 重排 (多样性、业务规则)
    ranked = s.reRank(ctx, ranked, req.Count)

    return ranked, nil
}
```

---

## 7. 最佳实践

### 7.1 性能优化

| 优化点 | 方案 |
|--------|------|
| 商品详情 | CDN + 本地缓存 |
| 搜索 | Elasticsearch + Redis |
| 库存 | Redis + Lua原子操作 |
| 订单 | 分库分表 + 读写分离 |
| 秒杀 | 令牌桶限流 + 消息队列 |

### 7.2 数据一致性

```
最终一致性场景:
- 库存同步: 异步消息
- 搜索索引: 监听Binlog
- 缓存更新: Cache Aside模式

强一致性场景:
- 订单状态: 分布式事务
- 库存扣减: Redis原子操作
- 支付状态: 两阶段提交
```

---

## 8. 参考文献

1. "Designing Data-Intensive Applications" - Martin Kleppmann
2. "Building Microservices" - Sam Newman
3. "E-commerce System Design"
4. Alibaba E-commerce Architecture
5. JD E-commerce Technical Whitepaper

---

*Last Updated: 2026-04-03*
