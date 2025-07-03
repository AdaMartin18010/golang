# 策略模式 (Strategy Pattern)

## 目录

1. [概述](#1-概述)
2. [理论基础](#2-理论基础)
3. [Go语言实现](#3-go语言实现)
4. [工程案例](#4-工程案例)
5. [批判性分析](#5-批判性分析)
6. [面试题与考点](#6-面试题与考点)
7. [术语表](#7-术语表)
8. [常见陷阱](#8-常见陷阱)
9. [相关主题](#9-相关主题)
10. [学习路径](#10-学习路径)

## 1. 概述

### 1.1 定义

策略模式定义了一系列算法，并将每一个算法封装起来，使它们可以互相替换。策略模式让算法独立于使用它的客户而变化。

**形式化定义**:
$$Strategy = (Context, Strategy, ConcreteStrategy_1, ConcreteStrategy_2, ..., ConcreteStrategy_n)$$

其中：

- $Context$ 是上下文类
- $Strategy$ 是策略接口
- $ConcreteStrategy_i$ 是具体策略实现

### 1.2 核心特征

- **算法封装**: 每个算法都被封装在独立的类中
- **可替换性**: 算法可以互相替换
- **扩展性**: 新增算法不影响现有代码
- **单一职责**: 每个策略类只负责一个算法

## 2. 理论基础

### 2.1 数学形式化

**定义 2.1** (策略模式): 策略模式是一个四元组 $S = (C, \Sigma, A, E)$

其中：

- $C$ 是上下文集合
- $\Sigma$ 是策略集合
- $A$ 是算法函数，$A: C \times \Sigma \rightarrow Result$
- $E$ 是执行环境

**定理 2.1** (策略可替换性): 对于任意上下文 $c \in C$ 和策略 $\sigma_1, \sigma_2 \in \Sigma$，如果 $\sigma_1$ 和 $\sigma_2$ 具有相同的接口，则它们可以互相替换。

**证明**: 由策略接口的一致性保证。

### 2.2 范畴论视角

在范畴论中，策略模式可以表示为：

$$Strategy : Context \times Algorithm \rightarrow Result$$

其中 $Context$、$Algorithm$ 和 $Result$ 是相应的对象范畴。

## 3. Go语言实现

### 3.1 基础策略模式

```go
package strategy

import (
    "fmt"
    "sort"
)

// SortStrategy 排序策略接口
type SortStrategy interface {
    Sort(data []int) []int
    GetName() string
}

// BubbleSortStrategy 冒泡排序策略
type BubbleSortStrategy struct{}

func NewBubbleSortStrategy() *BubbleSortStrategy {
    return &BubbleSortStrategy{}
}

func (b *BubbleSortStrategy) Sort(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)
    
    n := len(result)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if result[j] > result[j+1] {
                result[j], result[j+1] = result[j+1], result[j]
            }
        }
    }
    
    return result
}

func (b *BubbleSortStrategy) GetName() string {
    return "Bubble Sort"
}

// QuickSortStrategy 快速排序策略
type QuickSortStrategy struct{}

func NewQuickSortStrategy() *QuickSortStrategy {
    return &QuickSortStrategy{}
}

func (q *QuickSortStrategy) Sort(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)
    
    q.quickSort(result, 0, len(result)-1)
    return result
}

func (q *QuickSortStrategy) quickSort(data []int, low, high int) {
    if low < high {
        pi := q.partition(data, low, high)
        q.quickSort(data, low, pi-1)
        q.quickSort(data, pi+1, high)
    }
}

func (q *QuickSortStrategy) partition(data []int, low, high int) int {
    pivot := data[high]
    i := low - 1
    
    for j := low; j < high; j++ {
        if data[j] <= pivot {
            i++
            data[i], data[j] = data[j], data[i]
        }
    }
    
    data[i+1], data[high] = data[high], data[i+1]
    return i + 1
}

func (q *QuickSortStrategy) GetName() string {
    return "Quick Sort"
}

// MergeSortStrategy 归并排序策略
type MergeSortStrategy struct{}

func NewMergeSortStrategy() *MergeSortStrategy {
    return &MergeSortStrategy{}
}

func (m *MergeSortStrategy) Sort(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)
    
    return m.mergeSort(result)
}

func (m *MergeSortStrategy) mergeSort(data []int) []int {
    if len(data) <= 1 {
        return data
    }
    
    mid := len(data) / 2
    left := m.mergeSort(data[:mid])
    right := m.mergeSort(data[mid:])
    
    return m.merge(left, right)
}

func (m *MergeSortStrategy) merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    
    return result
}

func (m *MergeSortStrategy) GetName() string {
    return "Merge Sort"
}

// BuiltinSortStrategy 内置排序策略
type BuiltinSortStrategy struct{}

func NewBuiltinSortStrategy() *BuiltinSortStrategy {
    return &BuiltinSortStrategy{}
}

func (b *BuiltinSortStrategy) Sort(data []int) []int {
    result := make([]int, len(data))
    copy(result, data)
    
    sort.Ints(result)
    return result
}

func (b *BuiltinSortStrategy) GetName() string {
    return "Built-in Sort"
}

// SortContext 排序上下文
type SortContext struct {
    strategy SortStrategy
}

func NewSortContext(strategy SortStrategy) *SortContext {
    return &SortContext{
        strategy: strategy,
    }
}

func (s *SortContext) SetStrategy(strategy SortStrategy) {
    s.strategy = strategy
}

func (s *SortContext) ExecuteSort(data []int) []int {
    if s.strategy == nil {
        return data
    }
    
    fmt.Printf("Using %s strategy\n", s.strategy.GetName())
    return s.strategy.Sort(data)
}

func (s *SortContext) GetCurrentStrategy() string {
    if s.strategy == nil {
        return "No strategy set"
    }
    return s.strategy.GetName()
}
```

### 3.2 支付策略模式

```go
package paymentstrategy

import (
    "fmt"
    "time"
)

// PaymentInfo 支付信息
type PaymentInfo struct {
    Amount      float64
    Currency    string
    OrderID     string
    CustomerID  string
    Timestamp   time.Time
}

func NewPaymentInfo(amount float64, currency, orderID, customerID string) *PaymentInfo {
    return &PaymentInfo{
        Amount:     amount,
        Currency:   currency,
        OrderID:    orderID,
        CustomerID: customerID,
        Timestamp:  time.Now(),
    }
}

// PaymentResult 支付结果
type PaymentResult struct {
    Success     bool
    TransactionID string
    Message     string
    Timestamp   time.Time
}

func NewPaymentResult(success bool, transactionID, message string) *PaymentResult {
    return &PaymentResult{
        Success:       success,
        TransactionID: transactionID,
        Message:       message,
        Timestamp:     time.Now(),
    }
}

// PaymentStrategy 支付策略接口
type PaymentStrategy interface {
    ProcessPayment(info *PaymentInfo) *PaymentResult
    GetName() string
    GetSupportedCurrencies() []string
    GetProcessingFee() float64
}

// CreditCardStrategy 信用卡支付策略
type CreditCardStrategy struct {
    processingFee float64
}

func NewCreditCardStrategy() *CreditCardStrategy {
    return &CreditCardStrategy{
        processingFee: 0.025, // 2.5%
    }
}

func (c *CreditCardStrategy) ProcessPayment(info *PaymentInfo) *PaymentResult {
    fmt.Printf("Processing credit card payment for order %s\n", info.OrderID)
    
    // 模拟支付处理
    if info.Amount > 0 {
        return NewPaymentResult(true, 
            fmt.Sprintf("CC_%s_%d", info.OrderID, time.Now().Unix()),
            "Credit card payment processed successfully")
    }
    
    return NewPaymentResult(false, "", "Invalid amount")
}

func (c *CreditCardStrategy) GetName() string {
    return "Credit Card"
}

func (c *CreditCardStrategy) GetSupportedCurrencies() []string {
    return []string{"USD", "EUR", "GBP", "JPY", "CNY"}
}

func (c *CreditCardStrategy) GetProcessingFee() float64 {
    return c.processingFee
}

// PayPalStrategy PayPal支付策略
type PayPalStrategy struct {
    processingFee float64
}

func NewPayPalStrategy() *PayPalStrategy {
    return &PayPalStrategy{
        processingFee: 0.029, // 2.9%
    }
}

func (p *PayPalStrategy) ProcessPayment(info *PaymentInfo) *PaymentResult {
    fmt.Printf("Processing PayPal payment for order %s\n", info.OrderID)
    
    // 模拟支付处理
    if info.Amount > 0 {
        return NewPaymentResult(true,
            fmt.Sprintf("PP_%s_%d", info.OrderID, time.Now().Unix()),
            "PayPal payment processed successfully")
    }
    
    return NewPaymentResult(false, "", "Invalid amount")
}

func (p *PayPalStrategy) GetName() string {
    return "PayPal"
}

func (p *PayPalStrategy) GetSupportedCurrencies() []string {
    return []string{"USD", "EUR", "GBP", "CAD", "AUD"}
}

func (p *PayPalStrategy) GetProcessingFee() float64 {
    return p.processingFee
}

// CryptoStrategy 加密货币支付策略
type CryptoStrategy struct {
    processingFee float64
}

func NewCryptoStrategy() *CryptoStrategy {
    return &CryptoStrategy{
        processingFee: 0.01, // 1%
    }
}

func (c *CryptoStrategy) ProcessPayment(info *PaymentInfo) *PaymentResult {
    fmt.Printf("Processing cryptocurrency payment for order %s\n", info.OrderID)
    
    // 模拟支付处理
    if info.Amount > 0 {
        return NewPaymentResult(true,
            fmt.Sprintf("CRYPTO_%s_%d", info.OrderID, time.Now().Unix()),
            "Cryptocurrency payment processed successfully")
    }
    
    return NewPaymentResult(false, "", "Invalid amount")
}

func (c *CryptoStrategy) GetName() string {
    return "Cryptocurrency"
}

func (c *CryptoStrategy) GetSupportedCurrencies() []string {
    return []string{"BTC", "ETH", "USDT", "USDC"}
}

func (c *CryptoStrategy) GetProcessingFee() float64 {
    return c.processingFee
}

// BankTransferStrategy 银行转账策略
type BankTransferStrategy struct {
    processingFee float64
}

func NewBankTransferStrategy() *BankTransferStrategy {
    return &BankTransferStrategy{
        processingFee: 0.005, // 0.5%
    }
}

func (b *BankTransferStrategy) ProcessPayment(info *PaymentInfo) *PaymentResult {
    fmt.Printf("Processing bank transfer for order %s\n", info.OrderID)
    
    // 模拟支付处理
    if info.Amount > 0 {
        return NewPaymentResult(true,
            fmt.Sprintf("BANK_%s_%d", info.OrderID, time.Now().Unix()),
            "Bank transfer initiated successfully")
    }
    
    return NewPaymentResult(false, "", "Invalid amount")
}

func (b *BankTransferStrategy) GetName() string {
    return "Bank Transfer"
}

func (b *BankTransferStrategy) GetSupportedCurrencies() []string {
    return []string{"USD", "EUR", "GBP", "JPY", "CNY", "CAD", "AUD"}
}

func (b *BankTransferStrategy) GetProcessingFee() float64 {
    return b.processingFee
}

// PaymentProcessor 支付处理器
type PaymentProcessor struct {
    strategy PaymentStrategy
}

func NewPaymentProcessor(strategy PaymentStrategy) *PaymentProcessor {
    return &PaymentProcessor{
        strategy: strategy,
    }
}

func (p *PaymentProcessor) SetStrategy(strategy PaymentStrategy) {
    p.strategy = strategy
}

func (p *PaymentProcessor) ProcessPayment(info *PaymentInfo) *PaymentResult {
    if p.strategy == nil {
        return NewPaymentResult(false, "", "No payment strategy set")
    }
    
    fmt.Printf("Using %s payment method\n", p.strategy.GetName())
    
    // 检查货币支持
    supported := false
    for _, currency := range p.strategy.GetSupportedCurrencies() {
        if currency == info.Currency {
            supported = true
            break
        }
    }
    
    if !supported {
        return NewPaymentResult(false, "", 
            fmt.Sprintf("Currency %s not supported by %s", 
                info.Currency, p.strategy.GetName()))
    }
    
    // 计算手续费
    fee := info.Amount * p.strategy.GetProcessingFee()
    totalAmount := info.Amount + fee
    
    fmt.Printf("Processing fee: %.2f %s\n", fee, info.Currency)
    fmt.Printf("Total amount: %.2f %s\n", totalAmount, info.Currency)
    
    return p.strategy.ProcessPayment(info)
}

func (p *PaymentProcessor) GetCurrentStrategy() string {
    if p.strategy == nil {
        return "No strategy set"
    }
    return p.strategy.GetName()
}

func (p *PaymentProcessor) GetSupportedCurrencies() []string {
    if p.strategy == nil {
        return []string{}
    }
    return p.strategy.GetSupportedCurrencies()
}
```

### 3.3 压缩策略模式

```go
package compressionstrategy

import (
    "bytes"
    "compress/gzip"
    "compress/zlib"
    "fmt"
    "io"
)

// CompressionResult 压缩结果
type CompressionResult struct {
    OriginalSize    int
    CompressedSize  int
    CompressionRatio float64
    Data            []byte
    Algorithm       string
}

func NewCompressionResult(originalSize, compressedSize int, data []byte, algorithm string) *CompressionResult {
    ratio := 0.0
    if originalSize > 0 {
        ratio = float64(compressedSize) / float64(originalSize)
    }
    
    return &CompressionResult{
        OriginalSize:    originalSize,
        CompressedSize:  compressedSize,
        CompressionRatio: ratio,
        Data:            data,
        Algorithm:       algorithm,
    }
}

func (c *CompressionResult) String() string {
    return fmt.Sprintf("%s: %d -> %d bytes (%.2f%% compression)", 
        c.Algorithm, c.OriginalSize, c.CompressedSize, (1-c.CompressionRatio)*100)
}

// CompressionStrategy 压缩策略接口
type CompressionStrategy interface {
    Compress(data []byte) *CompressionResult
    Decompress(data []byte) ([]byte, error)
    GetName() string
    GetExtension() string
}

// GzipStrategy Gzip压缩策略
type GzipStrategy struct{}

func NewGzipStrategy() *GzipStrategy {
    return &GzipStrategy{}
}

func (g *GzipStrategy) Compress(data []byte) *CompressionResult {
    var buf bytes.Buffer
    writer := gzip.NewWriter(&buf)
    
    _, err := writer.Write(data)
    if err != nil {
        return NewCompressionResult(len(data), 0, nil, "Gzip")
    }
    
    writer.Close()
    compressed := buf.Bytes()
    
    return NewCompressionResult(len(data), len(compressed), compressed, "Gzip")
}

func (g *GzipStrategy) Decompress(data []byte) ([]byte, error) {
    reader, err := gzip.NewReader(bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    defer reader.Close()
    
    var buf bytes.Buffer
    _, err = io.Copy(&buf, reader)
    if err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}

func (g *GzipStrategy) GetName() string {
    return "Gzip"
}

func (g *GzipStrategy) GetExtension() string {
    return ".gz"
}

// ZlibStrategy Zlib压缩策略
type ZlibStrategy struct{}

func NewZlibStrategy() *ZlibStrategy {
    return &ZlibStrategy{}
}

func (z *ZlibStrategy) Compress(data []byte) *CompressionResult {
    var buf bytes.Buffer
    writer := zlib.NewWriter(&buf)
    
    _, err := writer.Write(data)
    if err != nil {
        return NewCompressionResult(len(data), 0, nil, "Zlib")
    }
    
    writer.Close()
    compressed := buf.Bytes()
    
    return NewCompressionResult(len(data), len(compressed), compressed, "Zlib")
}

func (z *ZlibStrategy) Decompress(data []byte) ([]byte, error) {
    reader, err := zlib.NewReader(bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    defer reader.Close()
    
    var buf bytes.Buffer
    _, err = io.Copy(&buf, reader)
    if err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}

func (z *ZlibStrategy) GetName() string {
    return "Zlib"
}

func (z *ZlibStrategy) GetExtension() string {
    return ".zlib"
}

// NoCompressionStrategy 无压缩策略
type NoCompressionStrategy struct{}

func NewNoCompressionStrategy() *NoCompressionStrategy {
    return &NoCompressionStrategy{}
}

func (n *NoCompressionStrategy) Compress(data []byte) *CompressionResult {
    return NewCompressionResult(len(data), len(data), data, "No Compression")
}

func (n *NoCompressionStrategy) Decompress(data []byte) ([]byte, error) {
    return data, nil
}

func (n *NoCompressionStrategy) GetName() string {
    return "No Compression"
}

func (n *NoCompressionStrategy) GetExtension() string {
    return ""
}

// CompressionContext 压缩上下文
type CompressionContext struct {
    strategy CompressionStrategy
}

func NewCompressionContext(strategy CompressionStrategy) *CompressionContext {
    return &CompressionContext{
        strategy: strategy,
    }
}

func (c *CompressionContext) SetStrategy(strategy CompressionStrategy) {
    c.strategy = strategy
}

func (c *CompressionContext) Compress(data []byte) *CompressionResult {
    if c.strategy == nil {
        return NewCompressionResult(len(data), 0, nil, "No Strategy")
    }
    
    fmt.Printf("Using %s compression\n", c.strategy.GetName())
    return c.strategy.Compress(data)
}

func (c *CompressionContext) Decompress(data []byte) ([]byte, error) {
    if c.strategy == nil {
        return nil, fmt.Errorf("no strategy set")
    }
    
    return c.strategy.Decompress(data)
}

func (c *CompressionContext) GetCurrentStrategy() string {
    if c.strategy == nil {
        return "No strategy set"
    }
    return c.strategy.GetName()
}
```

## 4. 工程案例

### 4.1 路由策略模式

```go
package routingstrategy

import (
    "fmt"
    "math"
    "sort"
)

// Location 位置信息
type Location struct {
    ID       string
    Name     string
    Latitude float64
    Longitude float64
}

func NewLocation(id, name string, lat, lng float64) *Location {
    return &Location{
        ID:        id,
        Name:      name,
        Latitude:  lat,
        Longitude: lng,
    }
}

func (l *Location) DistanceTo(other *Location) float64 {
    const earthRadius = 6371 // 地球半径，单位：公里
    
    lat1 := l.Latitude * math.Pi / 180
    lat2 := other.Latitude * math.Pi / 180
    deltaLat := (other.Latitude - l.Latitude) * math.Pi / 180
    deltaLng := (other.Longitude - l.Longitude) * math.Pi / 180
    
    a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
        math.Cos(lat1)*math.Cos(lat2)*
            math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return earthRadius * c
}

// Route 路由信息
type Route struct {
    Start     *Location
    End       *Location
    Distance  float64
    Duration  float64
    Path      []*Location
}

func NewRoute(start, end *Location) *Route {
    return &Route{
        Start: start,
        End:   end,
        Path:  []*Location{start, end},
    }
}

func (r *Route) CalculateDistance() {
    r.Distance = r.Start.DistanceTo(r.End)
}

func (r *Route) String() string {
    return fmt.Sprintf("Route from %s to %s: %.2f km", 
        r.Start.Name, r.End.Name, r.Distance)
}

// RoutingStrategy 路由策略接口
type RoutingStrategy interface {
    FindRoute(start, end *Location, waypoints []*Location) *Route
    GetName() string
    GetDescription() string
}

// ShortestDistanceStrategy 最短距离策略
type ShortestDistanceStrategy struct{}

func NewShortestDistanceStrategy() *ShortestDistanceStrategy {
    return &ShortestDistanceStrategy{}
}

func (s *ShortestDistanceStrategy) FindRoute(start, end *Location, waypoints []*Location) *Route {
    route := NewRoute(start, end)
    
    if len(waypoints) == 0 {
        route.CalculateDistance()
        return route
    }
    
    // 使用最近邻算法找到最短路径
    current := start
    remaining := make([]*Location, len(waypoints))
    copy(remaining, waypoints)
    
    totalDistance := 0.0
    path := []*Location{start}
    
    for len(remaining) > 0 {
        nearest := s.findNearest(current, remaining)
        distance := current.DistanceTo(nearest)
        totalDistance += distance
        
        path = append(path, nearest)
        current = nearest
        
        // 从剩余点中移除最近的点
        for i, wp := range remaining {
            if wp.ID == nearest.ID {
                remaining = append(remaining[:i], remaining[i+1:]...)
                break
            }
        }
    }
    
    // 添加终点
    finalDistance := current.DistanceTo(end)
    totalDistance += finalDistance
    path = append(path, end)
    
    route.Distance = totalDistance
    route.Path = path
    
    return route
}

func (s *ShortestDistanceStrategy) findNearest(from *Location, to []*Location) *Location {
    if len(to) == 0 {
        return nil
    }
    
    nearest := to[0]
    minDistance := from.DistanceTo(nearest)
    
    for _, location := range to[1:] {
        distance := from.DistanceTo(location)
        if distance < minDistance {
            minDistance = distance
            nearest = location
        }
    }
    
    return nearest
}

func (s *ShortestDistanceStrategy) GetName() string {
    return "Shortest Distance"
}

func (s *ShortestDistanceStrategy) GetDescription() string {
    return "Finds the route with the shortest total distance"
}

// FastestRouteStrategy 最快路线策略
type FastestRouteStrategy struct {
    speedMap map[string]float64 // 不同道路类型的速度
}

func NewFastestRouteStrategy() *FastestRouteStrategy {
    return &FastestRouteStrategy{
        speedMap: map[string]float64{
            "highway": 120.0, // 高速公路 120 km/h
            "urban":   50.0,  // 城市道路 50 km/h
            "rural":   80.0,  // 乡村道路 80 km/h
        },
    }
}

func (f *FastestRouteStrategy) FindRoute(start, end *Location, waypoints []*Location) *Route {
    route := NewRoute(start, end)
    
    // 简化实现：假设所有路段都是高速公路
    avgSpeed := f.speedMap["highway"]
    
    if len(waypoints) == 0 {
        route.CalculateDistance()
        route.Duration = route.Distance / avgSpeed
        return route
    }
    
    // 使用最短距离策略，但计算时间而不是距离
    distanceStrategy := NewShortestDistanceStrategy()
    tempRoute := distanceStrategy.FindRoute(start, end, waypoints)
    
    route.Distance = tempRoute.Distance
    route.Path = tempRoute.Path
    route.Duration = route.Distance / avgSpeed
    
    return route
}

func (f *FastestRouteStrategy) GetName() string {
    return "Fastest Route"
}

func (f *FastestRouteStrategy) GetDescription() string {
    return "Finds the route with the shortest travel time"
}

// ScenicRouteStrategy 风景路线策略
type ScenicRouteStrategy struct {
    scenicPoints map[string]float64 // 风景点评分
}

func NewScenicRouteStrategy() *ScenicRouteStrategy {
    return &ScenicRouteStrategy{
        scenicPoints: map[string]float64{
            "mountain": 0.9,
            "lake":     0.8,
            "forest":   0.7,
            "city":     0.3,
        },
    }
}

func (s *ScenicRouteStrategy) FindRoute(start, end *Location, waypoints []*Location) *Route {
    route := NewRoute(start, end)
    
    // 简化实现：优先选择风景点
    if len(waypoints) == 0 {
        route.CalculateDistance()
        return route
    }
    
    // 按风景评分排序
    scenicWaypoints := make([]*Location, len(waypoints))
    copy(scenicWaypoints, waypoints)
    
    sort.Slice(scenicWaypoints, func(i, j int) bool {
        scoreI := s.getScenicScore(scenicWaypoints[i])
        scoreJ := s.getScenicScore(scenicWaypoints[j])
        return scoreI > scoreJ
    })
    
    // 使用排序后的路径
    distanceStrategy := NewShortestDistanceStrategy()
    tempRoute := distanceStrategy.FindRoute(start, end, scenicWaypoints)
    
    route.Distance = tempRoute.Distance
    route.Path = tempRoute.Path
    
    return route
}

func (s *ScenicRouteStrategy) getScenicScore(location *Location) float64 {
    // 简化实现：根据位置名称判断风景类型
    name := location.Name
    if contains(name, "mountain") {
        return s.scenicPoints["mountain"]
    } else if contains(name, "lake") {
        return s.scenicPoints["lake"]
    } else if contains(name, "forest") {
        return s.scenicPoints["forest"]
    }
    return s.scenicPoints["city"]
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) && (s == substr || 
        (len(s) > len(substr) && (s[:len(substr)] == substr || 
         s[len(s)-len(substr):] == substr || 
         contains(s[1:len(s)-1], substr))))
}

func (s *ScenicRouteStrategy) GetName() string {
    return "Scenic Route"
}

func (s *ScenicRouteStrategy) GetDescription() string {
    return "Finds the route with the most scenic views"
}

// RoutingContext 路由上下文
type RoutingContext struct {
    strategy RoutingStrategy
}

func NewRoutingContext(strategy RoutingStrategy) *RoutingContext {
    return &RoutingContext{
        strategy: strategy,
    }
}

func (r *RoutingContext) SetStrategy(strategy RoutingStrategy) {
    r.strategy = strategy
}

func (r *RoutingContext) FindRoute(start, end *Location, waypoints []*Location) *Route {
    if r.strategy == nil {
        return NewRoute(start, end)
    }
    
    fmt.Printf("Using %s routing strategy\n", r.strategy.GetName())
    return r.strategy.FindRoute(start, end, waypoints)
}

func (r *RoutingContext) GetCurrentStrategy() string {
    if r.strategy == nil {
        return "No strategy set"
    }
    return r.strategy.GetName()
}

func (r *RoutingContext) GetStrategyDescription() string {
    if r.strategy == nil {
        return "No strategy set"
    }
    return r.strategy.GetDescription()
}
```

## 5. 批判性分析

### 5.1 优势

1. **算法封装**: 每个算法独立封装
2. **可替换性**: 算法可以动态替换
3. **扩展性**: 新增算法不影响现有代码
4. **单一职责**: 每个策略只负责一个算法

### 5.2 劣势

1. **策略数量**: 策略过多会增加复杂度
2. **性能开销**: 策略切换可能有性能开销
3. **配置复杂**: 策略选择逻辑可能复杂
4. **测试困难**: 每个策略都需要单独测试

### 5.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 | 高 | 低 |
| Java | 接口 | 中 | 中 |
| C++ | 虚函数 | 高 | 中 |
| Python | 函数对象 | 中 | 低 |

### 5.4 最新趋势

1. **函数式策略**: 使用函数作为策略
2. **策略工厂**: 策略创建和管理
3. **策略组合**: 多个策略组合使用
4. **配置驱动**: 通过配置选择策略

## 6. 面试题与考点

### 6.1 基础考点

1. **Q**: 策略模式与状态模式的区别？
   **A**: 策略模式关注算法选择，状态模式关注状态转换

2. **Q**: 什么时候使用策略模式？
   **A**: 有多种算法选择、算法需要动态切换时

3. **Q**: 策略模式的优缺点？
   **A**: 优点：算法封装、可替换；缺点：策略数量多、配置复杂

### 6.2 进阶考点

1. **Q**: 如何避免策略模式中的if-else？
   **A**: 使用策略工厂、配置驱动

2. **Q**: 策略模式在微服务中的应用？
   **A**: 服务选择、负载均衡、路由策略

3. **Q**: 如何处理策略的性能差异？
   **A**: 性能监控、动态选择、缓存策略

## 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 策略模式 | 封装算法的设计模式 | Strategy Pattern |
| 策略 | 具体的算法实现 | Strategy |
| 上下文 | 使用策略的类 | Context |
| 算法族 | 相关的算法集合 | Algorithm Family |

## 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 策略过多 | 增加系统复杂度 | 策略分类、层次化 |
| 性能差异 | 不同策略性能不同 | 性能监控、动态选择 |
| 配置复杂 | 策略选择逻辑复杂 | 配置驱动、策略工厂 |
| 测试困难 | 每个策略都要测试 | 单元测试、集成测试 |

## 9. 相关主题

- [观察者模式](./01-Observer-Pattern.md)
- [命令模式](./03-Command-Pattern.md)
- [状态模式](./04-State-Pattern.md)
- [模板方法模式](./05-Template-Method-Pattern.md)
- [迭代器模式](./06-Iterator-Pattern.md)

## 10. 学习路径

### 10.1 新手路径

1. 理解策略模式的基本概念
2. 学习策略和上下文的关系
3. 实现简单的策略模式
4. 理解策略的可替换性

### 10.2 进阶路径

1. 学习复杂的策略实现
2. 理解策略的性能优化
3. 掌握策略的应用场景
4. 学习策略的最佳实践

### 10.3 高阶路径

1. 分析策略在大型项目中的应用
2. 理解策略与架构设计的关系
3. 掌握策略的性能调优
4. 学习策略的替代方案

---

**相关文档**: [行为型模式总览](./README.md) | [设计模式总览](../README.md)
