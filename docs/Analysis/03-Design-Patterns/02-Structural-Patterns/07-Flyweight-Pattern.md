# 3.2.1 享元模式 (Flyweight Pattern)

<!-- TOC START -->
- [3.2.1 享元模式 (Flyweight Pattern)](#321-享元模式-flyweight-pattern)
  - [3.2.1.1 目录](#3211-目录)
  - [3.2.1.2 1. 概述](#3212-1-概述)
    - [3.2.1.2.1 定义](#32121-定义)
    - [3.2.1.2.2 核心特征](#32122-核心特征)
  - [3.2.1.3 2. 理论基础](#3213-2-理论基础)
    - [3.2.1.3.1 数学形式化](#32131-数学形式化)
    - [3.2.1.3.2 范畴论视角](#32132-范畴论视角)
  - [3.2.1.4 3. Go语言实现](#3214-3-go语言实现)
    - [3.2.1.4.1 基础享元模式](#32141-基础享元模式)
    - [3.2.1.4.2 字符享元模式](#32142-字符享元模式)
    - [3.2.1.4.3 图形享元模式](#32143-图形享元模式)
  - [3.2.1.5 4. 工程案例](#3215-4-工程案例)
    - [3.2.1.5.1 数据库连接池享元模式](#32151-数据库连接池享元模式)
    - [3.2.1.5.2 缓存享元模式](#32152-缓存享元模式)
  - [3.2.1.6 5. 批判性分析](#3216-5-批判性分析)
    - [3.2.1.6.1 优势](#32161-优势)
    - [3.2.1.6.2 劣势](#32162-劣势)
    - [3.2.1.6.3 行业对比](#32163-行业对比)
    - [3.2.1.6.4 最新趋势](#32164-最新趋势)
  - [3.2.1.7 6. 面试题与考点](#3217-6-面试题与考点)
    - [3.2.1.7.1 基础考点](#32171-基础考点)
    - [3.2.1.7.2 进阶考点](#32172-进阶考点)
  - [3.2.1.8 7. 术语表](#3218-7-术语表)
  - [3.2.1.9 8. 常见陷阱](#3219-8-常见陷阱)
  - [3.2.1.10 9. 相关主题](#32110-9-相关主题)
  - [3.2.1.11 10. 学习路径](#32111-10-学习路径)
    - [3.2.1.11.1 新手路径](#321111-新手路径)
    - [3.2.1.11.2 进阶路径](#321112-进阶路径)
    - [3.2.1.11.3 高阶路径](#321113-高阶路径)
<!-- TOC END -->

## 3.2.1.1 目录

## 3.2.1.2 1. 概述

### 3.2.1.2.1 定义

享元模式通过共享技术有效地支持大量细粒度对象的复用。

**形式化定义**:
$$Flyweight = (Flyweight, ConcreteFlyweight, FlyweightFactory, Client, IntrinsicState, ExtrinsicState)$$

其中：

- $Flyweight$ 是享元接口
- $ConcreteFlyweight$ 是具体享元
- $FlyweightFactory$ 是享元工厂
- $Client$ 是客户端
- $IntrinsicState$ 是内部状态
- $ExtrinsicState$ 是外部状态

### 3.2.1.2.2 核心特征

- **共享对象**: 共享相同内部状态的对象
- **状态分离**: 内部状态与外部状态分离
- **对象池**: 使用对象池管理享元对象
- **内存优化**: 减少内存使用和对象创建开销

## 3.2.1.3 2. 理论基础

### 3.2.1.3.1 数学形式化

**定义 2.1** (享元模式): 享元模式是一个七元组 $F = (I, E, P, F, C, M, V)$

其中：

- $I$ 是内部状态集合
- $E$ 是外部状态集合
- $P$ 是享元池
- $F$ 是享元工厂函数，$F: I \rightarrow P$
- $C$ 是客户端集合
- $M$ 是内存映射函数
- $V$ 是验证规则

**定理 2.1** (共享性): 对于任意内部状态 $i \in I$，存在唯一的享元对象 $p \in P$ 使得 $F(i) = p$

**证明**: 由享元工厂的单例性质保证。

### 3.2.1.3.2 范畴论视角

在范畴论中，享元模式可以表示为：

$$Flyweight : IntrinsicState \rightarrow SharedObject$$

其中 $IntrinsicState$ 是内部状态范畴，$SharedObject$ 是共享对象范畴。

## 3.2.1.4 3. Go语言实现

### 3.2.1.4.1 基础享元模式

```go
package flyweight

import (
    "fmt"
    "sync"
)

// Flyweight 享元接口
type Flyweight interface {
    Operation(extrinsicState string) string
    GetIntrinsicState() string
}

// ConcreteFlyweight 具体享元
type ConcreteFlyweight struct {
    intrinsicState string
}

func NewConcreteFlyweight(intrinsicState string) *ConcreteFlyweight {
    return &ConcreteFlyweight{intrinsicState: intrinsicState}
}

func (c *ConcreteFlyweight) Operation(extrinsicState string) string {
    return fmt.Sprintf("ConcreteFlyweight: intrinsic=%s, extrinsic=%s", 
        c.intrinsicState, extrinsicState)
}

func (c *ConcreteFlyweight) GetIntrinsicState() string {
    return c.intrinsicState
}

// UnsharedConcreteFlyweight 非共享具体享元
type UnsharedConcreteFlyweight struct {
    allState string
}

func NewUnsharedConcreteFlyweight(allState string) *UnsharedConcreteFlyweight {
    return &UnsharedConcreteFlyweight{allState: allState}
}

func (u *UnsharedConcreteFlyweight) Operation(extrinsicState string) string {
    return fmt.Sprintf("UnsharedConcreteFlyweight: allState=%s, extrinsic=%s", 
        u.allState, extrinsicState)
}

func (u *UnsharedConcreteFlyweight) GetIntrinsicState() string {
    return u.allState
}

// FlyweightFactory 享元工厂
type FlyweightFactory struct {
    flyweights map[string]Flyweight
    mu         sync.RWMutex
}

func NewFlyweightFactory() *FlyweightFactory {
    return &FlyweightFactory{
        flyweights: make(map[string]Flyweight),
    }
}

func (f *FlyweightFactory) GetFlyweight(key string) Flyweight {
    f.mu.RLock()
    if flyweight, exists := f.flyweights[key]; exists {
        f.mu.RUnlock()
        return flyweight
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    // 双重检查
    if flyweight, exists := f.flyweights[key]; exists {
        return flyweight
    }
    
    // 创建新的享元对象
    flyweight := NewConcreteFlyweight(key)
    f.flyweights[key] = flyweight
    
    return flyweight
}

func (f *FlyweightFactory) GetFlyweightCount() int {
    f.mu.RLock()
    defer f.mu.RUnlock()
    return len(f.flyweights)
}

func (f *FlyweightFactory) ListFlyweights() []string {
    f.mu.RLock()
    defer f.mu.RUnlock()
    
    keys := make([]string, 0, len(f.flyweights))
    for key := range f.flyweights {
        keys = append(keys, key)
    }
    return keys
}
```

### 3.2.1.4.2 字符享元模式

```go
package characterflyweight

import (
    "fmt"
    "sync"
)

// Character 字符接口
type Character interface {
    Display(font string, size int, color string) string
    GetChar() rune
}

// CharacterFlyweight 字符享元
type CharacterFlyweight struct {
    char rune
}

func NewCharacterFlyweight(char rune) *CharacterFlyweight {
    return &CharacterFlyweight{char: char}
}

func (c *CharacterFlyweight) Display(font string, size int, color string) string {
    return fmt.Sprintf("Character '%c' with font=%s, size=%d, color=%s", 
        c.char, font, size, color)
}

func (c *CharacterFlyweight) GetChar() rune {
    return c.char
}

// CharacterFactory 字符工厂
type CharacterFactory struct {
    characters map[rune]*CharacterFlyweight
    mu         sync.RWMutex
}

func NewCharacterFactory() *CharacterFactory {
    return &CharacterFactory{
        characters: make(map[rune]*CharacterFlyweight),
    }
}

func (f *CharacterFactory) GetCharacter(char rune) *CharacterFlyweight {
    f.mu.RLock()
    if character, exists := f.characters[char]; exists {
        f.mu.RUnlock()
        return character
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    // 双重检查
    if character, exists := f.characters[char]; exists {
        return character
    }
    
    // 创建新的字符享元
    character := NewCharacterFlyweight(char)
    f.characters[char] = character
    
    return character
}

func (f *CharacterFactory) GetCharacterCount() int {
    f.mu.RLock()
    defer f.mu.RUnlock()
    return len(f.characters)
}

// TextEditor 文本编辑器
type TextEditor struct {
    factory    *CharacterFactory
    characters []*CharacterFlyweight
    positions  []Position
}

type Position struct {
    x, y int
}

func NewTextEditor() *TextEditor {
    return &TextEditor{
        factory:    NewCharacterFactory(),
        characters: make([]*CharacterFlyweight, 0),
        positions:  make([]Position, 0),
    }
}

func (t *TextEditor) AddCharacter(char rune, x, y int) {
    character := t.factory.GetCharacter(char)
    t.characters = append(t.characters, character)
    t.positions = append(t.positions, Position{x: x, y: y})
}

func (t *TextEditor) Display(font string, size int, color string) string {
    result := fmt.Sprintf("TextEditor with %d characters:\n", len(t.characters))
    
    for i, character := range t.characters {
        pos := t.positions[i]
        display := character.Display(font, size, color)
        result += fmt.Sprintf("  [%d,%d]: %s\n", pos.x, pos.y, display)
    }
    
    return result
}

func (t *TextEditor) GetUniqueCharacterCount() int {
    return t.factory.GetCharacterCount()
}

func (t *TextEditor) GetTotalCharacterCount() int {
    return len(t.characters)
}
```

### 3.2.1.4.3 图形享元模式

```go
package shapeflyweight

import (
    "fmt"
    "sync"
)

// Shape 图形接口
type Shape interface {
    Draw(x, y int, color string) string
    GetType() string
}

// CircleFlyweight 圆形享元
type CircleFlyweight struct {
    radius int
}

func NewCircleFlyweight(radius int) *CircleFlyweight {
    return &CircleFlyweight{radius: radius}
}

func (c *CircleFlyweight) Draw(x, y int, color string) string {
    return fmt.Sprintf("Circle at (%d,%d) with radius %d, color %s", 
        x, y, c.radius, color)
}

func (c *CircleFlyweight) GetType() string {
    return "circle"
}

// RectangleFlyweight 矩形享元
type RectangleFlyweight struct {
    width  int
    height int
}

func NewRectangleFlyweight(width, height int) *RectangleFlyweight {
    return &RectangleFlyweight{
        width:  width,
        height: height,
    }
}

func (r *RectangleFlyweight) Draw(x, y int, color string) string {
    return fmt.Sprintf("Rectangle at (%d,%d) with size %dx%d, color %s", 
        x, y, r.width, r.height, color)
}

func (r *RectangleFlyweight) GetType() string {
    return "rectangle"
}

// TriangleFlyweight 三角形享元
type TriangleFlyweight struct {
    base   int
    height int
}

func NewTriangleFlyweight(base, height int) *TriangleFlyweight {
    return &TriangleFlyweight{
        base:   base,
        height: height,
    }
}

func (t *TriangleFlyweight) Draw(x, y int, color string) string {
    return fmt.Sprintf("Triangle at (%d,%d) with base %d, height %d, color %s", 
        x, y, t.base, t.height, color)
}

func (t *TriangleFlyweight) GetType() string {
    return "triangle"
}

// ShapeFactory 图形工厂
type ShapeFactory struct {
    circles    map[int]*CircleFlyweight
    rectangles map[string]*RectangleFlyweight
    triangles  map[string]*TriangleFlyweight
    mu         sync.RWMutex
}

func NewShapeFactory() *ShapeFactory {
    return &ShapeFactory{
        circles:    make(map[int]*CircleFlyweight),
        rectangles: make(map[string]*RectangleFlyweight),
        triangles:  make(map[string]*TriangleFlyweight),
    }
}

func (f *ShapeFactory) GetCircle(radius int) *CircleFlyweight {
    f.mu.RLock()
    if circle, exists := f.circles[radius]; exists {
        f.mu.RUnlock()
        return circle
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    if circle, exists := f.circles[radius]; exists {
        return circle
    }
    
    circle := NewCircleFlyweight(radius)
    f.circles[radius] = circle
    return circle
}

func (f *ShapeFactory) GetRectangle(width, height int) *RectangleFlyweight {
    key := fmt.Sprintf("%dx%d", width, height)
    
    f.mu.RLock()
    if rect, exists := f.rectangles[key]; exists {
        f.mu.RUnlock()
        return rect
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    if rect, exists := f.rectangles[key]; exists {
        return rect
    }
    
    rect := NewRectangleFlyweight(width, height)
    f.rectangles[key] = rect
    return rect
}

func (f *ShapeFactory) GetTriangle(base, height int) *TriangleFlyweight {
    key := fmt.Sprintf("%dx%d", base, height)
    
    f.mu.RLock()
    if triangle, exists := f.triangles[key]; exists {
        f.mu.RUnlock()
        return triangle
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    if triangle, exists := f.triangles[key]; exists {
        return triangle
    }
    
    triangle := NewTriangleFlyweight(base, height)
    f.triangles[key] = triangle
    return triangle
}

func (f *ShapeFactory) GetFlyweightCount() int {
    f.mu.RLock()
    defer f.mu.RUnlock()
    return len(f.circles) + len(f.rectangles) + len(f.triangles)
}

// DrawingCanvas 绘图画布
type DrawingCanvas struct {
    factory *ShapeFactory
    shapes  []ShapeInstance
}

type ShapeInstance struct {
    shape Shape
    x, y  int
    color string
}

func NewDrawingCanvas() *DrawingCanvas {
    return &DrawingCanvas{
        factory: NewShapeFactory(),
        shapes:  make([]ShapeInstance, 0),
    }
}

func (d *DrawingCanvas) AddCircle(radius, x, y int, color string) {
    circle := d.factory.GetCircle(radius)
    d.shapes = append(d.shapes, ShapeInstance{
        shape: circle,
        x:     x,
        y:     y,
        color: color,
    })
}

func (d *DrawingCanvas) AddRectangle(width, height, x, y int, color string) {
    rect := d.factory.GetRectangle(width, height)
    d.shapes = append(d.shapes, ShapeInstance{
        shape: rect,
        x:     x,
        y:     y,
        color: color,
    })
}

func (d *DrawingCanvas) AddTriangle(base, height, x, y int, color string) {
    triangle := d.factory.GetTriangle(base, height)
    d.shapes = append(d.shapes, ShapeInstance{
        shape: triangle,
        x:     x,
        y:     y,
        color: color,
    })
}

func (d *DrawingCanvas) Draw() string {
    result := fmt.Sprintf("DrawingCanvas with %d shapes:\n", len(d.shapes))
    
    for i, instance := range d.shapes {
        result += fmt.Sprintf("  %d: %s\n", i+1, instance.shape.Draw(instance.x, instance.y, instance.color))
    }
    
    return result
}

func (d *DrawingCanvas) GetUniqueShapeCount() int {
    return d.factory.GetFlyweightCount()
}

func (d *DrawingCanvas) GetTotalShapeCount() int {
    return len(d.shapes)
}
```

## 3.2.1.5 4. 工程案例

### 3.2.1.5.1 数据库连接池享元模式

```go
package connectionflyweight

import (
    "database/sql"
    "fmt"
    "sync"
    "time"
)

// ConnectionConfig 连接配置
type ConnectionConfig struct {
    Driver   string
    Host     string
    Port     int
    Database string
    Username string
    Password string
}

func (c *ConnectionConfig) GetKey() string {
    return fmt.Sprintf("%s://%s:%s@%s:%d/%s", 
        c.Driver, c.Username, c.Password, c.Host, c.Port, c.Database)
}

// ConnectionFlyweight 连接享元
type ConnectionFlyweight struct {
    config *ConnectionConfig
    db     *sql.DB
    mu     sync.Mutex
    inUse  bool
    lastUsed time.Time
}

func NewConnectionFlyweight(config *ConnectionConfig) (*ConnectionFlyweight, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", 
        config.Username, config.Password, config.Host, config.Port, config.Database)
    
    db, err := sql.Open(config.Driver, dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    return &ConnectionFlyweight{
        config:   config,
        db:       db,
        lastUsed: time.Now(),
    }, nil
}

func (c *ConnectionFlyweight) GetDB() *sql.DB {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.inUse = true
    c.lastUsed = time.Now()
    return c.db
}

func (c *ConnectionFlyweight) Release() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.inUse = false
    c.lastUsed = time.Now()
}

func (c *ConnectionFlyweight) IsInUse() bool {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.inUse
}

func (c *ConnectionFlyweight) GetLastUsed() time.Time {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.lastUsed
}

func (c *ConnectionFlyweight) Close() error {
    return c.db.Close()
}

// ConnectionPool 连接池
type ConnectionPool struct {
    configs    map[string]*ConnectionConfig
    connections map[string][]*ConnectionFlyweight
    maxConnections int
    mu          sync.RWMutex
}

func NewConnectionPool(maxConnections int) *ConnectionPool {
    return &ConnectionPool{
        configs:       make(map[string]*ConnectionConfig),
        connections:   make(map[string][]*ConnectionFlyweight),
        maxConnections: maxConnections,
    }
}

func (p *ConnectionPool) RegisterConfig(config *ConnectionConfig) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    key := config.GetKey()
    p.configs[key] = config
    p.connections[key] = make([]*ConnectionFlyweight, 0)
}

func (p *ConnectionPool) GetConnection(configKey string) (*ConnectionFlyweight, error) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    connections, exists := p.connections[configKey]
    if !exists {
        return nil, fmt.Errorf("config not found: %s", configKey)
    }
    
    // 查找可用的连接
    for _, conn := range connections {
        if !conn.IsInUse() {
            conn.GetDB() // 标记为使用中
            return conn, nil
        }
    }
    
    // 创建新连接
    if len(connections) < p.maxConnections {
        config := p.configs[configKey]
        conn, err := NewConnectionFlyweight(config)
        if err != nil {
            return nil, err
        }
        
        conn.GetDB() // 标记为使用中
        p.connections[configKey] = append(p.connections[configKey], conn)
        return conn, nil
    }
    
    return nil, fmt.Errorf("connection pool is full")
}

func (p *ConnectionPool) ReleaseConnection(conn *ConnectionFlyweight) {
    conn.Release()
}

func (p *ConnectionPool) GetConnectionCount(configKey string) int {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    if connections, exists := p.connections[configKey]; exists {
        return len(connections)
    }
    return 0
}

func (p *ConnectionPool) GetActiveConnectionCount(configKey string) int {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    if connections, exists := p.connections[configKey]; exists {
        count := 0
        for _, conn := range connections {
            if conn.IsInUse() {
                count++
            }
        }
        return count
    }
    return 0
}

func (p *ConnectionPool) CloseAll() error {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    for _, connections := range p.connections {
        for _, conn := range connections {
            if err := conn.Close(); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

### 3.2.1.5.2 缓存享元模式

```go
package cacheflyweight

import (
    "fmt"
    "sync"
    "time"
)

// CacheItem 缓存项
type CacheItem struct {
    key       string
    value     interface{}
    createdAt time.Time
    lastUsed  time.Time
    useCount  int
}

func NewCacheItem(key string, value interface{}) *CacheItem {
    now := time.Now()
    return &CacheItem{
        key:       key,
        value:     value,
        createdAt: now,
        lastUsed:  now,
        useCount:  1,
    }
}

func (c *CacheItem) GetValue() interface{} {
    c.lastUsed = time.Now()
    c.useCount++
    return c.value
}

func (c *CacheItem) GetKey() string {
    return c.key
}

func (c *CacheItem) GetLastUsed() time.Time {
    return c.lastUsed
}

func (c *CacheItem) GetUseCount() int {
    return c.useCount
}

func (c *CacheItem) GetAge() time.Duration {
    return time.Since(c.createdAt)
}

// CacheFlyweight 缓存享元
type CacheFlyweight struct {
    items map[string]*CacheItem
    mu    sync.RWMutex
    maxItems int
    ttl      time.Duration
}

func NewCacheFlyweight(maxItems int, ttl time.Duration) *CacheFlyweight {
    cache := &CacheFlyweight{
        items:    make(map[string]*CacheItem),
        maxItems: maxItems,
        ttl:      ttl,
    }
    
    // 启动清理协程
    go cache.cleanupRoutine()
    
    return cache
}

func (c *CacheFlyweight) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    item, exists := c.items[key]
    c.mu.RUnlock()
    
    if !exists {
        return nil, false
    }
    
    // 检查TTL
    if c.ttl > 0 && time.Since(item.lastUsed) > c.ttl {
        c.mu.Lock()
        delete(c.items, key)
        c.mu.Unlock()
        return nil, false
    }
    
    return item.GetValue(), true
}

func (c *CacheFlyweight) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // 检查容量限制
    if len(c.items) >= c.maxItems {
        c.evictLRU()
    }
    
    c.items[key] = NewCacheItem(key, value)
}

func (c *CacheFlyweight) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

func (c *CacheFlyweight) evictLRU() {
    var oldestKey string
    var oldestTime time.Time
    
    for key, item := range c.items {
        if oldestKey == "" || item.lastUsed.Before(oldestTime) {
            oldestKey = key
            oldestTime = item.lastUsed
        }
    }
    
    if oldestKey != "" {
        delete(c.items, oldestKey)
    }
}

func (c *CacheFlyweight) cleanupRoutine() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, item := range c.items {
            if c.ttl > 0 && now.Sub(item.lastUsed) > c.ttl {
                delete(c.items, key)
            }
        }
        c.mu.Unlock()
    }
}

func (c *CacheFlyweight) GetStats() map[string]interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    stats := make(map[string]interface{})
    stats["total_items"] = len(c.items)
    stats["max_items"] = c.maxItems
    
    if len(c.items) > 0 {
        var totalUses int
        var oldestItem *CacheItem
        var newestItem *CacheItem
        
        for _, item := range c.items {
            totalUses += item.useCount
            if oldestItem == nil || item.createdAt.Before(oldestItem.createdAt) {
                oldestItem = item
            }
            if newestItem == nil || item.createdAt.After(newestItem.createdAt) {
                newestItem = item
            }
        }
        
        stats["total_uses"] = totalUses
        stats["avg_uses_per_item"] = float64(totalUses) / float64(len(c.items))
        stats["oldest_item_age"] = oldestItem.GetAge().String()
        stats["newest_item_age"] = newestItem.GetAge().String()
    }
    
    return stats
}

// CacheFactory 缓存工厂
type CacheFactory struct {
    caches map[string]*CacheFlyweight
    mu     sync.RWMutex
}

func NewCacheFactory() *CacheFactory {
    return &CacheFactory{
        caches: make(map[string]*CacheFlyweight),
    }
}

func (f *CacheFactory) GetCache(name string, maxItems int, ttl time.Duration) *CacheFlyweight {
    f.mu.RLock()
    if cache, exists := f.caches[name]; exists {
        f.mu.RUnlock()
        return cache
    }
    f.mu.RUnlock()
    
    f.mu.Lock()
    defer f.mu.Unlock()
    
    if cache, exists := f.caches[name]; exists {
        return cache
    }
    
    cache := NewCacheFlyweight(maxItems, ttl)
    f.caches[name] = cache
    return cache
}

func (f *CacheFactory) GetCacheCount() int {
    f.mu.RLock()
    defer f.mu.RUnlock()
    return len(f.caches)
}

func (f *CacheFactory) ListCaches() []string {
    f.mu.RLock()
    defer f.mu.RUnlock()
    
    names := make([]string, 0, len(f.caches))
    for name := range f.caches {
        names = append(names, name)
    }
    return names
}
```

## 3.2.1.6 5. 批判性分析

### 3.2.1.6.1 优势

1. **内存优化**: 减少内存使用和对象创建开销
2. **性能提升**: 共享对象减少创建和销毁开销
3. **统一管理**: 集中管理共享对象
4. **扩展性**: 易于添加新的享元类型

### 3.2.1.6.2 劣势

1. **复杂性**: 增加系统复杂性
2. **线程安全**: 需要处理并发访问
3. **状态管理**: 内部状态和外部状态分离复杂
4. **调试困难**: 共享对象调试困难

### 3.2.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 + 工厂 | 高 | 中 |
| Java | 享元类 + 工厂 | 中 | 中 |
| C++ | 享元类 | 中 | 中 |
| Python | 对象池 | 高 | 低 |

### 3.2.1.6.4 最新趋势

1. **对象池模式**: 更通用的对象复用
2. **内存池**: 专门的内存管理
3. **连接池**: 数据库连接复用
4. **缓存策略**: 多级缓存

## 3.2.1.7 6. 面试题与考点

### 3.2.1.7.1 基础考点

1. **Q**: 享元模式与单例模式的区别？
   **A**: 享元共享多个对象，单例只有一个对象

2. **Q**: 什么时候使用享元模式？
   **A**: 需要大量相似对象、内存受限时

3. **Q**: 享元模式的优缺点？
   **A**: 优点：内存优化、性能提升；缺点：复杂性、线程安全

### 3.2.1.7.2 进阶考点

1. **Q**: 如何设计线程安全的享元？
   **A**: 使用互斥锁、读写锁、原子操作

2. **Q**: 享元模式在大型项目中的应用？
   **A**: 连接池、缓存、对象池

3. **Q**: 如何处理享元的内存管理？
   **A**: 引用计数、垃圾回收、手动清理

## 3.2.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 享元模式 | 共享细粒度对象的设计模式 | Flyweight Pattern |
| 内部状态 | 共享的不可变状态 | Intrinsic State |
| 外部状态 | 不共享的可变状态 | Extrinsic State |
| 享元工厂 | 管理享元对象的工厂 | Flyweight Factory |
| 对象池 | 对象复用的容器 | Object Pool |

## 3.2.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 线程安全 | 并发访问导致问题 | 使用锁机制 |
| 内存泄漏 | 享元对象无法释放 | 引用计数、清理机制 |
| 状态混乱 | 内部状态被意外修改 | 不可变设计 |
| 过度共享 | 过度共享导致复杂性 | 适度共享 |

## 3.2.1.10 9. 相关主题

- [适配器模式](./01-Adapter-Pattern.md)
- [装饰器模式](./02-Decorator-Pattern.md)
- [代理模式](./03-Proxy-Pattern.md)
- [外观模式](./04-Facade-Pattern.md)
- [桥接模式](./05-Bridge-Pattern.md)

## 3.2.1.11 10. 学习路径

### 3.2.1.11.1 新手路径

1. 理解享元模式的基本概念
2. 学习内部状态和外部状态
3. 实现简单的享元模式
4. 理解对象共享的原理

### 3.2.1.11.2 进阶路径

1. 学习复杂的享元实现
2. 理解享元的性能优化
3. 掌握享元的应用场景
4. 学习享元的最佳实践

### 3.2.1.11.3 高阶路径

1. 分析享元在大型项目中的应用
2. 理解享元与架构设计的关系
3. 掌握享元的性能调优
4. 学习享元的替代方案

---

**相关文档**: [结构型模式总览](./README.md) | [设计模式总览](../README.md)
