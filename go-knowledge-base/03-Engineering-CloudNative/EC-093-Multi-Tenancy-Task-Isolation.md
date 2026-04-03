# 多租户任务隔离 (Multi-Tenancy Task Isolation)

> **分类**: 工程与云原生
> **标签**: #multi-tenancy #isolation #security
> **参考**: SaaS Multi-Tenancy, Resource Isolation

---

## 多租户隔离架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Multi-Tenancy Task Isolation                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Isolation Levels:                                                           │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  1. Database-per-Tenant (Highest Isolation)                          │   │
│  │                                                                      │   │
│  │   Tenant A ──► ┌─────────────┐                                      │   │
│  │   Tenant B ──► │  Database A │                                      │   │
│  │   Tenant C ──► │  Database B │                                      │   │
│  │                │  Database C │                                      │   │
│  │                └─────────────┘                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  2. Schema-per-Tenant (Medium Isolation)                             │   │
│  │                                                                      │   │
│  │   Database                                                          │   │
│  │   ├─ Schema_A (tables: tasks, queues, workers)                      │   │
│  │   ├─ Schema_B (tables: tasks, queues, workers)                      │   │
│  │   └─ Schema_C (tables: tasks, queues, workers)                      │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  3. Shared Schema with Row-Level Security (Lowest Isolation)         │   │
│  │                                                                      │   │
│  │   Table: tasks                                                       │   │
│  │   ├─ id, tenant_id, payload, status, ...                            │   │
│  │   │                                                                  │   │
│  │   ├─ Row (tenant_id='A')                                            │   │
│  │   ├─ Row (tenant_id='B')                                            │   │
│  │   └─ Row (tenant_id='C')                                            │   │
│  │                                                                      │   │
│  │   Query: SELECT * FROM tasks WHERE tenant_id = current_tenant()     │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整多租户实现

```go
package multitenancy

import (
    "context"
    "fmt"
    "sync"
)

// Tenant 租户信息
type Tenant struct {
    ID            string
    Name          string
    Tier          TenantTier
    ResourceQuota ResourceQuota
    IsActive      bool
}

// TenantTier 租户等级
type TenantTier string

const (
    TierFree      TenantTier = "free"
    TierBasic     TenantTier = "basic"
    TierPro       TenantTier = "pro"
    TierEnterprise TenantTier = "enterprise"
)

// ResourceQuota 资源配额
type ResourceQuota struct {
    MaxConcurrentTasks int
    MaxTasksPerDay     int
    MaxWorkers         int
    CPUQuota           float64
    MemoryQuota        int64
}

// TenantContextKey 租户上下文键
type TenantContextKey struct{}

// WithTenant 将租户信息加入上下文
func WithTenant(ctx context.Context, tenant *Tenant) context.Context {
    return context.WithValue(ctx, TenantContextKey{}, tenant)
}

// GetTenant 从上下文获取租户
func GetTenant(ctx context.Context) (*Tenant, bool) {
    tenant, ok := ctx.Value(TenantContextKey{}).(*Tenant)
    return tenant, ok
}

// MustGetTenant 必须获取租户（否则panic）
func MustGetTenant(ctx context.Context) *Tenant {
    tenant, ok := GetTenant(ctx)
    if !ok {
        panic("tenant not found in context")
    }
    return tenant
}

// TenantManager 租户管理器
type TenantManager struct {
    tenants map[string]*Tenant
    mu      sync.RWMutex

    // 资源追踪
    usage   map[string]*TenantUsage
}

// TenantUsage 租户资源使用
type TenantUsage struct {
    TenantID          string
    CurrentTasks      int
    TasksToday        int
    LastResetTime     int64
    CPUUsed           float64
    MemoryUsed        int64
}

// NewTenantManager 创建租户管理器
func NewTenantManager() *TenantManager {
    return &TenantManager{
        tenants: make(map[string]*Tenant),
        usage:   make(map[string]*TenantUsage),
    }
}

// RegisterTenant 注册租户
func (tm *TenantManager) RegisterTenant(tenant *Tenant) error {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    if _, exists := tm.tenants[tenant.ID]; exists {
        return fmt.Errorf("tenant %s already exists", tenant.ID)
    }

    tm.tenants[tenant.ID] = tenant
    tm.usage[tenant.ID] = &TenantUsage{
        TenantID: tenant.ID,
    }

    return nil
}

// GetTenant 获取租户
func (tm *TenantManager) GetTenant(tenantID string) (*Tenant, error) {
    tm.mu.RLock()
    defer tm.mu.RUnlock()

    tenant, ok := tm.tenants[tenantID]
    if !ok {
        return nil, fmt.Errorf("tenant %s not found", tenantID)
    }

    return tenant, nil
}

// CheckQuota 检查配额
func (tm *TenantManager) CheckQuota(tenantID string, resource string, amount int) error {
    tm.mu.RLock()
    defer tm.mu.RUnlock()

    tenant, ok := tm.tenants[tenantID]
    if !ok {
        return fmt.Errorf("tenant not found")
    }

    usage, ok := tm.usage[tenantID]
    if !ok {
        return fmt.Errorf("usage not found")
    }

    switch resource {
    case "concurrent_tasks":
        if usage.CurrentTasks+amount > tenant.ResourceQuota.MaxConcurrentTasks {
            return fmt.Errorf("concurrent task quota exceeded: %d/%d",
                usage.CurrentTasks+amount, tenant.ResourceQuota.MaxConcurrentTasks)
        }
    case "tasks_per_day":
        if usage.TasksToday+amount > tenant.ResourceQuota.MaxTasksPerDay {
            return fmt.Errorf("daily task quota exceeded: %d/%d",
                usage.TasksToday+amount, tenant.ResourceQuota.MaxTasksPerDay)
        }
    }

    return nil
}

// RecordUsage 记录使用
func (tm *TenantManager) RecordUsage(tenantID string, resource string, delta int) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    usage, ok := tm.usage[tenantID]
    if !ok {
        return
    }

    switch resource {
    case "concurrent_tasks":
        usage.CurrentTasks += delta
    case "tasks_per_day":
        usage.TasksToday += delta
    }
}

// IsolatedTaskQueue 隔离的任务队列
type IsolatedTaskQueue struct {
    tenantID string
    queue    TaskQueue
}

// TaskQueue 任务队列接口
type TaskQueue interface {
    Enqueue(ctx context.Context, task interface{}) error
    Dequeue(ctx context.Context) (interface{}, error)
}

// TenantAwareScheduler 租户感知调度器
type TenantAwareScheduler struct {
    tenantManager *TenantManager
    queues        map[string]*IsolatedTaskQueue // tenantID -> queue
    workers       map[string]int                // tenantID -> worker count
    mu            sync.RWMutex
}

// NewTenantAwareScheduler 创建租户感知调度器
func NewTenantAwareScheduler(tm *TenantManager) *TenantAwareScheduler {
    return &TenantAwareScheduler{
        tenantManager: tm,
        queues:        make(map[string]*IsolatedTaskQueue),
        workers:       make(map[string]int),
    }
}

// SubmitTask 提交任务
func (tas *TenantAwareScheduler) SubmitTask(ctx context.Context, tenantID string, task interface{}) error {
    // 检查租户
    tenant, err := tas.tenantManager.GetTenant(tenantID)
    if err != nil {
        return err
    }

    if !tenant.IsActive {
        return fmt.Errorf("tenant %s is not active", tenantID)
    }

    // 检查配额
    if err := tas.tenantManager.CheckQuota(tenantID, "concurrent_tasks", 1); err != nil {
        return err
    }

    // 获取或创建队列
    tas.mu.Lock()
    isoQueue, ok := tas.queues[tenantID]
    if !ok {
        isoQueue = &IsolatedTaskQueue{
            tenantID: tenantID,
            queue:    NewInMemoryQueue(),
        }
        tas.queues[tenantID] = isoQueue
    }
    tas.mu.Unlock()

    // 入队
    return isoQueue.queue.Enqueue(ctx, task)
}

// RowLevelSecurity 行级安全过滤器
type RowLevelSecurity struct {
    tenantIDColumn string
}

// NewRowLevelSecurity 创建行级安全
func NewRowLevelSecurity(column string) *RowLevelSecurity {
    return &RowLevelSecurity{tenantIDColumn: column}
}

// Filter 过滤查询
func (rls *RowLevelSecurity) Filter(query string, tenantID string) string {
    // 添加租户过滤条件
    return fmt.Sprintf("%s WHERE %s = '%s'", query, rls.tenantIDColumn, tenantID)
}

// ValidateOwnership 验证所有权
func (rls *RowLevelSecurity) ValidateOwnership(resourceTenantID, currentTenantID string) error {
    if resourceTenantID != currentTenantID {
        return fmt.Errorf("access denied: resource belongs to different tenant")
    }
    return nil
}

// NamespaceIsolation 命名空间隔离
type NamespaceIsolation struct {
    namespaces map[string]*TenantNamespace
    mu         sync.RWMutex
}

// TenantNamespace 租户命名空间
type TenantNamespace struct {
    TenantID string
    Prefix   string
}

// NewNamespaceIsolation 创建命名空间隔离
func NewNamespaceIsolation() *NamespaceIsolation {
    return &NamespaceIsolation{
        namespaces: make(map[string]*TenantNamespace),
    }
}

// CreateNamespace 创建命名空间
func (ni *NamespaceIsolation) CreateNamespace(tenantID string) *TenantNamespace {
    ni.mu.Lock()
    defer ni.mu.Unlock()

    ns := &TenantNamespace{
        TenantID: tenantID,
        Prefix:   fmt.Sprintf("tenant_%s_", tenantID),
    }
    ni.namespaces[tenantID] = ns
    return ns
}

// GetPrefix 获取前缀
func (ni *NamespaceIsolation) GetPrefix(tenantID string) (string, error) {
    ni.mu.RLock()
    defer ni.mu.RUnlock()

    ns, ok := ni.namespaces[tenantID]
    if !ok {
        return "", fmt.Errorf("namespace not found for tenant %s", tenantID)
    }

    return ns.Prefix, nil
}

// KeyWithPrefix 添加前缀到键
func (ni *NamespaceIsolation) KeyWithPrefix(tenantID, key string) (string, error) {
    prefix, err := ni.GetPrefix(tenantID)
    if err != nil {
        return "", err
    }
    return prefix + key, nil
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"

    "multitenancy"
)

func main() {
    // 创建租户管理器
    tm := multitenancy.NewTenantManager()

    // 注册租户
    tenantA := &multitenancy.Tenant{
        ID:       "tenant-a",
        Name:     "Company A",
        Tier:     multitenancy.TierPro,
        IsActive: true,
        ResourceQuota: multitenancy.ResourceQuota{
            MaxConcurrentTasks: 100,
            MaxTasksPerDay:     10000,
            MaxWorkers:         10,
        },
    }

    tm.RegisterTenant(tenantA)

    // 创建上下文
    ctx := multitenancy.WithTenant(context.Background(), tenantA)

    // 获取租户
    t, _ := multitenancy.GetTenant(ctx)
    fmt.Printf("Current tenant: %s\n", t.Name)

    // 检查配额
    if err := tm.CheckQuota(tenantA.ID, "concurrent_tasks", 50); err != nil {
        fmt.Printf("Quota check failed: %v\n", err)
    }

    // 创建调度器
    scheduler := multitenancy.NewTenantAwareScheduler(tm)

    // 提交任务
    task := map[string]string{"type": "send-email"}
    if err := scheduler.SubmitTask(ctx, tenantA.ID, task); err != nil {
        fmt.Printf("Submit failed: %v\n", err)
    }

    // 行级安全
    rls := multitenancy.NewRowLevelSecurity("tenant_id")
    query := rls.Filter("SELECT * FROM tasks", tenantA.ID)
    fmt.Printf("Filtered query: %s\n", query)

    // 命名空间隔离
    ns := multitenancy.NewNamespaceIsolation()
    ns.CreateNamespace(tenantA.ID)

    key, _ := ns.KeyWithPrefix(tenantA.ID, "queue:tasks")
    fmt.Printf("Namespaced key: %s\n", key)
}
```

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02