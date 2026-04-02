# 任务多租户隔离 (Task Multi-Tenancy)

> **分类**: 工程与云原生  
> **标签**: #multi-tenancy #isolation #security

---

## 租户上下文

```go
// 租户上下文键
type tenantKey struct{}

type TenantContext struct {
    TenantID    string
    OrgID       string
    Plan        string  // free, basic, enterprise
    Quotas      Quotas
}

type Quotas struct {
    MaxConcurrentTasks int
    MaxTasksPerHour    int
    MaxTaskDuration    time.Duration
    PriorityLevels     []int
}

func WithTenant(ctx context.Context, tenant TenantContext) context.Context {
    return context.WithValue(ctx, tenantKey{}, tenant)
}

func TenantFromContext(ctx context.Context) (TenantContext, bool) {
    t, ok := ctx.Value(tenantKey{}).(TenantContext)
    return t, ok
}

// HTTP 中间件
func TenantMiddleware(tenantResolver TenantResolver) gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetHeader("X-Tenant-ID")
        if tenantID == "" {
            c.AbortWithStatusJSON(400, gin.H{"error": "missing tenant"})
            return
        }
        
        tenant, err := tenantResolver.Resolve(c.Request.Context(), tenantID)
        if err != nil {
            c.AbortWithStatusJSON(404, gin.H{"error": "tenant not found"})
            return
        }
        
        ctx := WithTenant(c.Request.Context(), tenant)
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

---

## 租户隔离执行器

```go
type TenantAwareExecutor struct {
    executors     map[string]*DedicatedExecutor  // tenantID -> executor
    sharedPool    *SharedExecutorPool
    isolationMode IsolationMode
}

type IsolationMode int

const (
    SharedPool IsolationMode = iota  // 共享资源，逻辑隔离
    DedicatedPool                     // 专用资源池
    DedicatedWorker                   // 专用工作节点
)

func (tae *TenantAwareExecutor) Execute(ctx context.Context, task *Task) error {
    tenant, ok := TenantFromContext(ctx)
    if !ok {
        return errors.New("no tenant in context")
    }
    
    // 检查配额
    if err := tae.checkQuota(ctx, tenant, task); err != nil {
        return err
    }
    
    // 根据隔离模式选择执行器
    switch tae.isolationMode {
    case DedicatedPool, DedicatedWorker:
        return tae.executeDedicated(ctx, tenant.TenantID, task)
    default:
        return tae.executeShared(ctx, tenant, task)
    }
}

func (tae *TenantAwareExecutor) checkQuota(ctx context.Context, tenant TenantContext, task *Task) error {
    quotas := tenant.Quotas
    
    // 检查并发数
    current := tae.getTenantConcurrentTasks(tenant.TenantID)
    if current >= quotas.MaxConcurrentTasks {
        return fmt.Errorf("concurrent task quota exceeded: %d/%d", 
            current, quotas.MaxConcurrentTasks)
    }
    
    // 检查每小时任务数
    hourly := tae.getTenantHourlyTasks(tenant.TenantID)
    if hourly >= quotas.MaxTasksPerHour {
        return fmt.Errorf("hourly task quota exceeded: %d/%d",
            hourly, quotas.MaxTasksPerHour)
    }
    
    return nil
}

func (tae *TenantAwareExecutor) executeDedicated(ctx context.Context, tenantID string, task *Task) error {
    executor, exists := tae.executors[tenantID]
    if !exists {
        // 动态创建专用执行器
        executor = tae.createDedicatedExecutor(tenantID)
        tae.executors[tenantID] = executor
    }
    
    return executor.Execute(ctx, task)
}
```

---

## 租户资源配额

```go
type TenantQuotaManager struct {
    store QuotaStore
}

func (tqm *TenantQuotaManager) EnforceQuota(ctx context.Context, tenantID string, task *Task) error {
    quotas, _ := tqm.store.GetQuota(ctx, tenantID)
    
    // CPU/内存限制
    if task.ResourceRequest.CPU > quotas.MaxCPU {
        return fmt.Errorf("CPU request %f exceeds quota %f", 
            task.ResourceRequest.CPU, quotas.MaxCPU)
    }
    
    if task.ResourceRequest.Memory > quotas.MaxMemory {
        return fmt.Errorf("memory request %d exceeds quota %d",
            task.ResourceRequest.Memory, quotas.MaxMemory)
    }
    
    // 优先级限制
    if !tqm.isPriorityAllowed(quotas, task.Priority) {
        return fmt.Errorf("priority %d not allowed for tenant", task.Priority)
    }
    
    return nil
}

func (tqm *TenantQuotaManager) GetUsage(ctx context.Context, tenantID string) (QuotaUsage, error) {
    tasks, _ := tqm.store.ListTenantTasks(ctx, tenantID)
    
    usage := QuotaUsage{
        ActiveTasks: len(tasks),
        CPUUsed:     0,
        MemoryUsed:  0,
    }
    
    for _, task := range tasks {
        usage.CPUUsed += task.ResourceUsage.CPU
        usage.MemoryUsed += task.ResourceUsage.Memory
    }
    
    return usage, nil
}
```

---

## 租户数据隔离

```go
// 数据库行级安全
type TenantIsolation struct {
    db *sql.DB
}

func (ti *TenantIsolation) QueryTenantTasks(ctx context.Context, tenantID string) ([]Task, error) {
    // 自动添加租户过滤
    rows, err := ti.db.QueryContext(ctx, 
        "SELECT * FROM tasks WHERE tenant_id = $1", tenantID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    // ...
}

// Redis 键隔离
func TenantKey(tenantID, key string) string {
    return fmt.Sprintf("tenant:%s:%s", tenantID, key)
}

func (ti *TenantIsolation) GetTaskQueue(ctx context.Context, tenantID string) (*TaskQueue, error) {
    key := TenantKey(tenantID, "task_queue")
    // 使用隔离的键访问 Redis
    return ti.redisClient.GetQueue(key)
}
```
