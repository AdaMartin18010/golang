# 3.2.1 代理模式 (Proxy Pattern)

## 3.2.1.1 目录

## 3.2.1.2 1. 概述

### 3.2.1.2.1 定义

代理模式为其他对象提供一种代理以控制对这个对象的访问。

**形式化定义**:
$$Proxy = (Subject, RealSubject, Proxy, Client, AccessControl)$$

其中：

- $Subject$ 是主题接口
- $RealSubject$ 是真实主题
- $Proxy$ 是代理类
- $Client$ 是客户端
- $AccessControl$ 是访问控制

### 3.2.1.2.2 核心特征

- **访问控制**: 控制对真实对象的访问
- **延迟加载**: 延迟创建真实对象
- **透明性**: 客户端无需知道代理存在
- **功能增强**: 在访问前后添加额外功能

## 3.2.1.3 2. 理论基础

### 3.2.1.3.1 数学形式化

**定义 2.1** (代理模式): 代理模式是一个六元组 $P = (S, R, P, A, C, V)$

其中：

- $S$ 是主题集合
- $R$ 是真实主题集合
- $P$ 是代理集合
- $A$ 是访问控制函数，$A: P \times S \rightarrow \{allow, deny\}$
- $C$ 是控制策略
- $V$ 是验证规则

**定理 2.1** (访问控制一致性): 对于任意代理 $p \in P$ 和主题 $s \in S$，$A(p, s) = allow$ 当且仅当满足控制策略 $C$

**证明**: 由访问控制函数的实现保证。

### 3.2.1.3.2 范畴论视角

在范畴论中，代理模式可以表示为：

$$Proxy : Subject \rightarrow Subject$$

其中 $Subject$ 是对象范畴。

## 3.2.1.4 3. Go语言实现

### 3.2.1.4.1 基础代理模式

```go
package proxy

import (
    "fmt"
    "time"
)

// Subject 主题接口
type Subject interface {
    Request() string
}

// RealSubject 真实主题
type RealSubject struct {
    name string
}

func (r *RealSubject) Request() string {
    // 模拟耗时操作
    time.Sleep(100 * time.Millisecond)
    return fmt.Sprintf("RealSubject(%s): Request processed", r.name)
}

// Proxy 代理
type Proxy struct {
    realSubject *RealSubject
    name        string
}

func NewProxy(name string) *Proxy {
    return &Proxy{name: name}
}

func (p *Proxy) Request() string {
    // 延迟初始化
    if p.realSubject == nil {
        p.realSubject = &RealSubject{name: p.name}
        fmt.Printf("Proxy: Created RealSubject for %s\n", p.name)
    }
    
    // 访问控制
    if !p.checkAccess() {
        return "Proxy: Access denied"
    }
    
    // 前置处理
    fmt.Printf("Proxy: Before request for %s\n", p.name)
    
    // 调用真实对象
    result := p.realSubject.Request()
    
    // 后置处理
    fmt.Printf("Proxy: After request for %s\n", p.name)
    
    return result
}

func (p *Proxy) checkAccess() bool {
    // 简单的访问控制逻辑
    return p.name != "forbidden"
}

```

### 3.2.1.4.2 虚拟代理模式

```go
package virtualproxy

import (
    "fmt"
    "image"
    "image/jpeg"
    "os"
    "sync"
)

// Image 图像接口
type Image interface {
    Display()
    GetWidth() int
    GetHeight() int
}

// RealImage 真实图像
type RealImage struct {
    filename string
    image    image.Image
    loaded   bool
}

func NewRealImage(filename string) *RealImage {
    return &RealImage{filename: filename}
}

func (r *RealImage) loadFromDisk() {
    if r.loaded {
        return
    }
    
    file, err := os.Open(r.filename)
    if err != nil {
        fmt.Printf("Error loading image: %v\n", err)
        return
    }
    defer file.Close()
    
    img, err := jpeg.Decode(file)
    if err != nil {
        fmt.Printf("Error decoding image: %v\n", err)
        return
    }
    
    r.image = img
    r.loaded = true
    fmt.Printf("Loaded image: %s\n", r.filename)
}

func (r *RealImage) Display() {
    r.loadFromDisk()
    if r.loaded {
        fmt.Printf("Displaying image: %s (%dx%d)\n", r.filename, r.GetWidth(), r.GetHeight())
    }
}

func (r *RealImage) GetWidth() int {
    r.loadFromDisk()
    if r.loaded {
        return r.image.Bounds().Dx()
    }
    return 0
}

func (r *RealImage) GetHeight() int {
    r.loadFromDisk()
    if r.loaded {
        return r.image.Bounds().Dy()
    }
    return 0
}

// VirtualProxy 虚拟代理
type VirtualProxy struct {
    realImage *RealImage
    filename  string
    mu        sync.Mutex
}

func NewVirtualProxy(filename string) *VirtualProxy {
    return &VirtualProxy{filename: filename}
}

func (v *VirtualProxy) Display() {
    v.mu.Lock()
    defer v.mu.Unlock()
    
    if v.realImage == nil {
        v.realImage = NewRealImage(v.filename)
    }
    v.realImage.Display()
}

func (v *VirtualProxy) GetWidth() int {
    v.mu.Lock()
    defer v.mu.Unlock()
    
    if v.realImage == nil {
        v.realImage = NewRealImage(v.filename)
    }
    return v.realImage.GetWidth()
}

func (v *VirtualProxy) GetHeight() int {
    v.mu.Lock()
    defer v.mu.Unlock()
    
    if v.realImage == nil {
        v.realImage = NewRealImage(v.filename)
    }
    return v.realImage.GetHeight()
}

```

### 3.2.1.4.3 保护代理模式

```go
package protectionproxy

import (
    "fmt"
    "time"
)

// User 用户
type User struct {
    Name     string
    Role     string
    IsActive bool
}

// Document 文档
type Document struct {
    ID      string
    Title   string
    Content string
    Owner   string
}

// DocumentService 文档服务接口
type DocumentService interface {
    ReadDocument(user User, docID string) (string, error)
    WriteDocument(user User, docID string, content string) error
    DeleteDocument(user User, docID string) error
}

// RealDocumentService 真实文档服务
type RealDocumentService struct {
    documents map[string]*Document
}

func NewRealDocumentService() *RealDocumentService {
    return &RealDocumentService{
        documents: make(map[string]*Document),
    }
}

func (r *RealDocumentService) ReadDocument(user User, docID string) (string, error) {
    doc, exists := r.documents[docID]
    if !exists {
        return "", fmt.Errorf("document not found: %s", docID)
    }
    
    // 模拟读取延迟
    time.Sleep(50 * time.Millisecond)
    
    return doc.Content, nil
}

func (r *RealDocumentService) WriteDocument(user User, docID string, content string) error {
    doc, exists := r.documents[docID]
    if !exists {
        doc = &Document{ID: docID, Owner: user.Name}
        r.documents[docID] = doc
    }
    
    doc.Content = content
    doc.Title = fmt.Sprintf("Document %s", docID)
    
    // 模拟写入延迟
    time.Sleep(100 * time.Millisecond)
    
    return nil
}

func (r *RealDocumentService) DeleteDocument(user User, docID string) error {
    _, exists := r.documents[docID]
    if !exists {
        return fmt.Errorf("document not found: %s", docID)
    }
    
    delete(r.documents, docID)
    
    // 模拟删除延迟
    time.Sleep(50 * time.Millisecond)
    
    return nil
}

// ProtectionProxy 保护代理
type ProtectionProxy struct {
    realService DocumentService
}

func NewProtectionProxy() *ProtectionProxy {
    return &ProtectionProxy{
        realService: NewRealDocumentService(),
    }
}

func (p *ProtectionProxy) ReadDocument(user User, docID string) (string, error) {
    // 检查用户权限
    if !p.canRead(user, docID) {
        return "", fmt.Errorf("access denied: user %s cannot read document %s", user.Name, docID)
    }
    
    return p.realService.ReadDocument(user, docID)
}

func (p *ProtectionProxy) WriteDocument(user User, docID string, content string) error {
    // 检查用户权限
    if !p.canWrite(user, docID) {
        return fmt.Errorf("access denied: user %s cannot write document %s", user.Name, docID)
    }
    
    return p.realService.WriteDocument(user, docID, content)
}

func (p *ProtectionProxy) DeleteDocument(user User, docID string) error {
    // 检查用户权限
    if !p.canDelete(user, docID) {
        return fmt.Errorf("access denied: user %s cannot delete document %s", user.Name, docID)
    }
    
    return p.realService.DeleteDocument(user, docID)
}

func (p *ProtectionProxy) canRead(user User, docID string) bool {
    // 简单的权限检查逻辑
    return user.IsActive && (user.Role == "admin" || user.Role == "user")
}

func (p *ProtectionProxy) canWrite(user User, docID string) bool {
    // 简单的权限检查逻辑
    return user.IsActive && (user.Role == "admin" || user.Role == "editor")
}

func (p *ProtectionProxy) canDelete(user User, docID string) bool {
    // 简单的权限检查逻辑
    return user.IsActive && user.Role == "admin"
}

```

## 3.2.1.5 4. 工程案例

### 3.2.1.5.1 缓存代理

```go
package cacheproxy

import (
    "fmt"
    "sync"
    "time"
)

// DataService 数据服务接口
type DataService interface {
    GetData(key string) (interface{}, error)
    SetData(key string, value interface{}) error
}

// RealDataService 真实数据服务
type RealDataService struct {
    data map[string]interface{}
}

func NewRealDataService() *RealDataService {
    return &RealDataService{
        data: make(map[string]interface{}),
    }
}

func (r *RealDataService) GetData(key string) (interface{}, error) {
    // 模拟数据库查询延迟
    time.Sleep(200 * time.Millisecond)
    
    value, exists := r.data[key]
    if !exists {
        return nil, fmt.Errorf("key not found: %s", key)
    }
    
    return value, nil
}

func (r *RealDataService) SetData(key string, value interface{}) error {
    // 模拟数据库写入延迟
    time.Sleep(100 * time.Millisecond)
    
    r.data[key] = value
    return nil
}

// CacheItem 缓存项
type CacheItem struct {
    Value      interface{}
    Expiry     time.Time
    LastAccess time.Time
}

// CacheProxy 缓存代理
type CacheProxy struct {
    realService DataService
    cache       map[string]*CacheItem
    mu          sync.RWMutex
    maxSize     int
    ttl         time.Duration
}

func NewCacheProxy(realService DataService, maxSize int, ttl time.Duration) *CacheProxy {
    proxy := &CacheProxy{
        realService: realService,
        cache:       make(map[string]*CacheItem),
        maxSize:     maxSize,
        ttl:         ttl,
    }
    
    // 启动清理协程
    go proxy.cleanupRoutine()
    
    return proxy
}

func (c *CacheProxy) GetData(key string) (interface{}, error) {
    c.mu.RLock()
    if item, exists := c.cache[key]; exists && time.Now().Before(item.Expiry) {
        item.LastAccess = time.Now()
        c.mu.RUnlock()
        fmt.Printf("Cache hit for key: %s\n", key)
        return item.Value, nil
    }
    c.mu.RUnlock()
    
    // 缓存未命中，从真实服务获取
    fmt.Printf("Cache miss for key: %s\n", key)
    value, err := c.realService.GetData(key)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // 检查缓存大小
    if len(c.cache) >= c.maxSize {
        c.evictLRU()
    }
    
    c.cache[key] = &CacheItem{
        Value:      value,
        Expiry:     time.Now().Add(c.ttl),
        LastAccess: time.Now(),
    }
    
    return value, nil
}

func (c *CacheProxy) SetData(key string, value interface{}) error {
    // 直接写入真实服务
    err := c.realService.SetData(key, value)
    if err != nil {
        return err
    }
    
    // 更新缓存
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.cache[key] = &CacheItem{
        Value:      value,
        Expiry:     time.Now().Add(c.ttl),
        LastAccess: time.Now(),
    }
    
    return nil
}

func (c *CacheProxy) evictLRU() {
    var oldestKey string
    var oldestTime time.Time
    
    for key, item := range c.cache {
        if oldestKey == "" || item.LastAccess.Before(oldestTime) {
            oldestKey = key
            oldestTime = item.LastAccess
        }
    }
    
    if oldestKey != "" {
        delete(c.cache, oldestKey)
        fmt.Printf("Evicted key from cache: %s\n", oldestKey)
    }
}

func (c *CacheProxy) cleanupRoutine() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, item := range c.cache {
            if now.After(item.Expiry) {
                delete(c.cache, key)
                fmt.Printf("Expired key from cache: %s\n", key)
            }
        }
        c.mu.Unlock()
    }
}

```

### 3.2.1.5.2 远程代理

```go
package remoteproxy

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// UserService 用户服务接口
type UserService interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
    UpdateUser(user *User) error
    DeleteUser(id string) error
}

// User 用户模型
type User struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    IsActive bool   `json:"is_active"`
}

// RemoteUserService 远程用户服务
type RemoteUserService struct {
    baseURL string
    client  *http.Client
}

func NewRemoteUserService(baseURL string) *RemoteUserService {
    return &RemoteUserService{
        baseURL: baseURL,
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (r *RemoteUserService) GetUser(id string) (*User, error) {
    url := fmt.Sprintf("%s/users/%s", r.baseURL, id)
    
    resp, err := r.client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get user: status %d", resp.StatusCode)
    }
    
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, fmt.Errorf("failed to decode user: %w", err)
    }
    
    return &user, nil
}

func (r *RemoteUserService) CreateUser(user *User) error {
    url := fmt.Sprintf("%s/users", r.baseURL)
    
    data, err := json.Marshal(user)
    if err != nil {
        return fmt.Errorf("failed to marshal user: %w", err)
    }
    
    resp, err := r.client.Post(url, "application/json", bytes.NewBuffer(data))
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("failed to create user: status %d", resp.StatusCode)
    }
    
    return nil
}

func (r *RemoteUserService) UpdateUser(user *User) error {
    url := fmt.Sprintf("%s/users/%s", r.baseURL, user.ID)
    
    data, err := json.Marshal(user)
    if err != nil {
        return fmt.Errorf("failed to marshal user: %w", err)
    }
    
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := r.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to update user: status %d", resp.StatusCode)
    }
    
    return nil
}

func (r *RemoteUserService) DeleteUser(id string) error {
    url := fmt.Sprintf("%s/users/%s", r.baseURL, id)
    
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    
    resp, err := r.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("failed to delete user: status %d", resp.StatusCode)
    }
    
    return nil
}

// RemoteProxy 远程代理
type RemoteProxy struct {
    realService UserService
    cache       map[string]*CachedUser
    mu          sync.RWMutex
    ttl         time.Duration
}

type CachedUser struct {
    User      *User
    Expiry    time.Time
}

func NewRemoteProxy(realService UserService, ttl time.Duration) *RemoteProxy {
    return &RemoteProxy{
        realService: realService,
        cache:       make(map[string]*CachedUser),
        ttl:         ttl,
    }
}

func (r *RemoteProxy) GetUser(id string) (*User, error) {
    // 检查缓存
    r.mu.RLock()
    if cached, exists := r.cache[id]; exists && time.Now().Before(cached.Expiry) {
        r.mu.RUnlock()
        fmt.Printf("Cache hit for user: %s\n", id)
        return cached.User, nil
    }
    r.mu.RUnlock()
    
    // 从远程服务获取
    fmt.Printf("Cache miss for user: %s\n", id)
    user, err := r.realService.GetUser(id)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    r.mu.Lock()
    r.cache[id] = &CachedUser{
        User:   user,
        Expiry: time.Now().Add(r.ttl),
    }
    r.mu.Unlock()
    
    return user, nil
}

func (r *RemoteProxy) CreateUser(user *User) error {
    err := r.realService.CreateUser(user)
    if err != nil {
        return err
    }
    
    // 清除相关缓存
    r.mu.Lock()
    delete(r.cache, user.ID)
    r.mu.Unlock()
    
    return nil
}

func (r *RemoteProxy) UpdateUser(user *User) error {
    err := r.realService.UpdateUser(user)
    if err != nil {
        return err
    }
    
    // 更新缓存
    r.mu.Lock()
    r.cache[user.ID] = &CachedUser{
        User:   user,
        Expiry: time.Now().Add(r.ttl),
    }
    r.mu.Unlock()
    
    return nil
}

func (r *RemoteProxy) DeleteUser(id string) error {
    err := r.realService.DeleteUser(id)
    if err != nil {
        return err
    }
    
    // 清除缓存
    r.mu.Lock()
    delete(r.cache, id)
    r.mu.Unlock()
    
    return nil
}

```

### 3.2.1.5.3 智能引用代理

```go
package smartreferenceproxy

import (
    "fmt"
    "sync"
    "time"
)

// Resource 资源接口
type Resource interface {
    Use() string
    GetID() string
    GetSize() int64
}

// ExpensiveResource 昂贵资源
type ExpensiveResource struct {
    id       string
    size     int64
    refCount int
    mu       sync.Mutex
}

func NewExpensiveResource(id string, size int64) *ExpensiveResource {
    return &ExpensiveResource{
        id:   id,
        size: size,
    }
}

func (e *ExpensiveResource) Use() string {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    e.refCount++
    fmt.Printf("Resource %s used, ref count: %d\n", e.id, e.refCount)
    
    // 模拟资源使用
    time.Sleep(50 * time.Millisecond)
    
    return fmt.Sprintf("Resource %s is being used", e.id)
}

func (e *ExpensiveResource) GetID() string {
    return e.id
}

func (e *ExpensiveResource) GetSize() int64 {
    return e.size
}

func (e *ExpensiveResource) Release() {
    e.mu.Lock()
    defer e.mu.Unlock()
    
    e.refCount--
    fmt.Printf("Resource %s released, ref count: %d\n", e.id, e.refCount)
}

// SmartReferenceProxy 智能引用代理
type SmartReferenceProxy struct {
    resource     *ExpensiveResource
    refCount     int
    lastAccess   time.Time
    mu           sync.Mutex
    resourcePool map[string]*ExpensiveResource
}

func NewSmartReferenceProxy() *SmartReferenceProxy {
    return &SmartReferenceProxy{
        resourcePool: make(map[string]*ExpensiveResource),
    }
}

func (s *SmartReferenceProxy) GetResource(id string) Resource {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // 检查资源池
    if resource, exists := s.resourcePool[id]; exists {
        s.refCount++
        s.lastAccess = time.Now()
        fmt.Printf("Reusing resource: %s\n", id)
        return s.createProxy(resource)
    }
    
    // 创建新资源
    resource := NewExpensiveResource(id, 1024)
    s.resourcePool[id] = resource
    s.refCount++
    s.lastAccess = time.Now()
    
    fmt.Printf("Created new resource: %s\n", id)
    return s.createProxy(resource)
}

func (s *SmartReferenceProxy) createProxy(resource *ExpensiveResource) Resource {
    return &ResourceProxy{
        resource: resource,
        parent:   s,
    }
}

func (s *SmartReferenceProxy) ReleaseResource(id string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if resource, exists := s.resourcePool[id]; exists {
        resource.Release()
        s.refCount--
        
        // 如果引用计数为0，从池中移除
        if resource.refCount == 0 {
            delete(s.resourcePool, id)
            fmt.Printf("Removed resource from pool: %s\n", id)
        }
    }
}

// ResourceProxy 资源代理
type ResourceProxy struct {
    resource *ExpensiveResource
    parent   *SmartReferenceProxy
}

func (r *ResourceProxy) Use() string {
    return r.resource.Use()
}

func (r *ResourceProxy) GetID() string {
    return r.resource.GetID()
}

func (r *ResourceProxy) GetSize() int64 {
    return r.resource.GetSize()
}

func (r *ResourceProxy) Release() {
    r.parent.ReleaseResource(r.resource.GetID())
}

```

## 3.2.1.6 5. 批判性分析

### 3.2.1.6.1 优势

1. **访问控制**: 提供细粒度的访问控制
2. **延迟加载**: 延迟创建昂贵对象
3. **缓存**: 提供缓存功能
4. **透明性**: 客户端无需知道代理存在

### 3.2.1.6.2 劣势

1. **复杂性**: 增加系统复杂性
2. **性能开销**: 代理层可能影响性能
3. **调试困难**: 代理链调试复杂
4. **过度使用**: 可能导致过度设计

### 3.2.1.6.3 行业对比

| 语言 | 实现方式 | 性能 | 复杂度 |
|------|----------|------|--------|
| Go | 接口 + 结构体 | 高 | 中 |
| Java | 动态代理 | 中 | 中 |
| C++ | 虚函数 | 中 | 中 |
| Python | 描述符 | 高 | 低 |

### 3.2.1.6.4 最新趋势

1. **智能代理**: 基于AI的访问控制
2. **分布式代理**: 微服务架构中的代理
3. **缓存代理**: 多级缓存策略
4. **安全代理**: 零信任架构

## 3.2.1.7 6. 面试题与考点

### 3.2.1.7.1 基础考点

1. **Q**: 代理模式与装饰器模式的区别？
   **A**: 代理控制访问，装饰器增强功能

2. **Q**: 什么时候使用代理模式？
   **A**: 需要访问控制、延迟加载、缓存时

3. **Q**: 代理模式的类型？
   **A**: 虚拟代理、保护代理、远程代理、缓存代理

### 3.2.1.7.2 进阶考点

1. **Q**: 如何设计高性能的缓存代理？
   **A**: 使用LRU、TTL、分布式缓存

2. **Q**: 代理模式在微服务中的应用？
   **A**: API网关、服务网格、负载均衡

3. **Q**: 如何处理代理的性能问题？
   **A**: 异步处理、连接池、批量操作

## 3.2.1.8 7. 术语表

| 术语 | 定义 | 英文 |
|------|------|------|
| 代理模式 | 控制对象访问的设计模式 | Proxy Pattern |
| 虚拟代理 | 延迟创建对象的代理 | Virtual Proxy |
| 保护代理 | 控制访问权限的代理 | Protection Proxy |
| 远程代理 | 控制远程对象访问的代理 | Remote Proxy |
| 缓存代理 | 提供缓存功能的代理 | Cache Proxy |

## 3.2.1.9 8. 常见陷阱

| 陷阱 | 问题 | 解决方案 |
|------|------|----------|
| 代理链过长 | 性能下降 | 优化代理层次 |
| 缓存一致性问题 | 数据不一致 | 使用缓存策略 |
| 访问控制复杂 | 权限管理困难 | 简化权限模型 |
| 过度代理 | 功能过度复杂 | 评估实际需求 |

## 3.2.1.10 9. 相关主题

- [适配器模式](./01-Adapter-Pattern.md)
- [装饰器模式](./02-Decorator-Pattern.md)
- [外观模式](./04-Facade-Pattern.md)
- [桥接模式](./05-Bridge-Pattern.md)
- [组合模式](./06-Composite-Pattern.md)

## 3.2.1.11 10. 学习路径

### 3.2.1.11.1 新手路径

1. 理解代理模式的基本概念
2. 学习不同类型的代理
3. 实现简单的代理
4. 理解访问控制的重要性

### 3.2.1.11.2 进阶路径

1. 学习缓存代理和远程代理
2. 理解代理的性能优化
3. 掌握代理的应用场景
4. 学习代理的最佳实践

### 3.2.1.11.3 高阶路径

1. 分析代理在大型项目中的应用
2. 理解代理与架构设计的关系
3. 掌握代理的性能调优
4. 学习代理的替代方案

---

**相关文档**: [结构型模式总览](./README.md) | [设计模式总览](../README.md)
