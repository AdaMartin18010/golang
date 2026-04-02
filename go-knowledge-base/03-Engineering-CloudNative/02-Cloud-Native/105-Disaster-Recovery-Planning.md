# 灾难恢复规划 (Disaster Recovery Planning)

> **分类**: 工程与云原生
> **标签**: #disaster-recovery #business-continuity #backup
> **参考**: AWS DR Strategies, Azure Site Recovery

---

## 灾难恢复架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Disaster Recovery Architecture                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  RPO (Recovery Point Objective): 0-5 minutes                                │
│  RTO (Recovery Time Objective): 15 minutes                                  │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Multi-Region Deployment                           │   │
│  │                                                                      │   │
│  │   Region A (Primary)          Region B (Standby)                    │   │
│  │   ┌─────────────────┐         ┌─────────────────┐                  │   │
│  │   │  Active Cluster │         │ Standby Cluster │                  │   │
│  │   │  ┌───────────┐  │◄───────►│  ┌───────────┐  │                  │   │
│  │   │  │ Master    │  │   Sync  │  │ Replica   │  │                  │   │
│  │   │  │ (Leader)  │  │────────►│  │ (Follower)│  │                  │   │
│  │   │  └───────────┘  │         │  └───────────┘  │                  │   │
│  │   └─────────────────┘         └─────────────────┘                  │   │
│  │                                                                      │   │
│  │   Data Replication: Async/Sync                                       │   │
│  │   Failover: Automatic with health checks                            │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整灾难恢复实现

```go
package dr

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// DRConfig 灾难恢复配置
type DRConfig struct {
    // RPO/RTO 目标
    RPO time.Duration // 恢复点目标
    RTO time.Duration // 恢复时间目标

    // 复制配置
    ReplicationMode string // "sync" | "async"
    SyncTimeout     time.Duration

    // 故障转移配置
    HealthCheckInterval time.Duration
    FailoverThreshold   int // 连续失败次数

    // 备份配置
    BackupInterval    time.Duration
    BackupRetention   int // 保留份数
}

// Site 站点
type Site struct {
    ID       string
    Region   string
    Endpoint string
    Status   SiteStatus
    Role     SiteRole

    // 健康状态
    LastHealthCheck time.Time
    ConsecutiveFails int
}

type SiteStatus string
const (
    SiteHealthy   SiteStatus = "healthy"
    SiteDegraded  SiteStatus = "degraded"
    SiteUnhealthy SiteStatus = "unhealthy"
)

type SiteRole string
const (
    RolePrimary  SiteRole = "primary"
    RoleStandby  SiteRole = "standby"
    RoleReplica  SiteRole = "replica"
)

// DisasterRecoveryManager 灾难恢复管理器
type DisasterRecoveryManager struct {
    config DRConfig

    // 站点管理
    sites       map[string]*Site
    primarySite string

    // 复制管理
    replicator DataReplicator

    // 故障检测
    healthChecker *HealthChecker

    // 故障转移控制
    failoverLock sync.RWMutex
    isFailoverInProgress bool

    // 备份管理
    backupManager *BackupManager
}

// DataReplicator 数据复制器接口
type DataReplicator interface {
    StartReplication(ctx context.Context, source, target string) error
    StopReplication(ctx context.Context) error
    GetReplicationLag(ctx context.Context) (time.Duration, error)
    Sync(ctx context.Context) error
}

// HealthChecker 健康检查器
type HealthChecker struct {
    config     DRConfig
    sites      map[string]*Site
    checkFunc  func(site *Site) (bool, error)
    stopCh     chan struct{}
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(config DRConfig, sites map[string]*Site) *HealthChecker {
    return &HealthChecker{
        config:    config,
        sites:     sites,
        checkFunc: defaultHealthCheck,
        stopCh:    make(chan struct{}),
    }
}

// Start 启动健康检查
func (hc *HealthChecker) Start() {
    ticker := time.NewTicker(hc.config.HealthCheckInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            hc.checkAllSites()
        case <-hc.stopCh:
            return
        }
    }
}

func (hc *HealthChecker) checkAllSites() {
    for _, site := range hc.sites {
        healthy, err := hc.checkFunc(site)
        site.LastHealthCheck = time.Now()

        if err != nil || !healthy {
            site.ConsecutiveFails++
            if site.ConsecutiveFails >= hc.config.FailoverThreshold {
                site.Status = SiteUnhealthy
            } else {
                site.Status = SiteDegraded
            }
        } else {
            site.ConsecutiveFails = 0
            site.Status = SiteHealthy
        }
    }
}

func defaultHealthCheck(site *Site) (bool, error) {
    // 实现健康检查逻辑
    return true, nil
}

// Failover 执行故障转移
func (drm *DisasterRecoveryManager) Failover(ctx context.Context, targetSiteID string) error {
    drm.failoverLock.Lock()
    defer drm.failoverLock.Unlock()

    if drm.isFailoverInProgress {
        return fmt.Errorf("failover already in progress")
    }

    drm.isFailoverInProgress = true
    defer func() { drm.isFailoverInProgress = false }()

    targetSite, ok := drm.sites[targetSiteID]
    if !ok {
        return fmt.Errorf("target site not found: %s", targetSiteID)
    }

    if targetSite.Status != SiteHealthy {
        return fmt.Errorf("target site is not healthy: %s", targetSite.Status)
    }

    // 1. 停止复制
    if err := drm.replicator.StopReplication(ctx); err != nil {
        return fmt.Errorf("failed to stop replication: %w", err)
    }

    // 2. 最终同步
    if err := drm.replicator.Sync(ctx); err != nil {
        return fmt.Errorf("final sync failed: %w", err)
    }

    // 3. 提升目标站点为主站点
    oldPrimary := drm.sites[drm.primarySite]
    oldPrimary.Role = RoleReplica

    targetSite.Role = RolePrimary
    drm.primarySite = targetSiteID

    // 4. 启动新的复制
    for _, site := range drm.sites {
        if site.ID != targetSiteID && site.Status == SiteHealthy {
            go drm.replicator.StartReplication(ctx, targetSiteID, site.ID)
        }
    }

    return nil
}

// BackupManager 备份管理器
type BackupManager struct {
    config      DRConfig
    storage     BackupStorage
    scheduler   *BackupScheduler
    retention   *RetentionPolicy
}

// BackupStorage 备份存储接口
type BackupStorage interface {
    Store(ctx context.Context, backup Backup) error
    Retrieve(ctx context.Context, backupID string) (*Backup, error)
    List(ctx context.Context) ([]Backup, error)
    Delete(ctx context.Context, backupID string) error
}

// Backup 备份
type Backup struct {
    ID          string
    Timestamp   time.Time
    Type        string // "full" | "incremental"
    Size        int64
    Checksum    string
    Metadata    map[string]string
    Status      string // "completed" | "failed" | "in_progress"
}

// BackupScheduler 备份调度器
type BackupScheduler struct {
    interval time.Duration
    executor func() error
    stopCh   chan struct{}
}

// Start 启动备份调度
func (bs *BackupScheduler) Start() {
    ticker := time.NewTicker(bs.interval)
    defer ticker.Stop()

    // 立即执行一次
    bs.executor()

    for {
        select {
        case <-ticker.C:
            if err := bs.executor(); err != nil {
                // 记录错误
            }
        case <-bs.stopCh:
            return
        }
    }
}

// RetentionPolicy 保留策略
type RetentionPolicy struct {
    maxBackups    int
    maxAge        time.Duration
    minBackups    int
}

// Apply 应用保留策略
func (rp *RetentionPolicy) Apply(backups []Backup) []string {
    if len(backups) <= rp.minBackups {
        return nil
    }

    var toDelete []string
    now := time.Now()

    // 按时间排序（旧的在前）
    sortBackupsByTime(backups)

    for i, backup := range backups {
        // 保留最小数量
        if len(backups)-len(toDelete) <= rp.minBackups {
            break
        }

        // 检查最大数量
        if len(backups)-len(toDelete) > rp.maxBackups {
            toDelete = append(toDelete, backup.ID)
            continue
        }

        // 检查最大年龄
        if now.Sub(backup.Timestamp) > rp.maxAge {
            toDelete = append(toDelete, backup.ID)
        }
    }

    return toDelete
}

func sortBackupsByTime(backups []Backup) {
    // 冒泡排序示例
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if backups[i].Timestamp.After(backups[j].Timestamp) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

// PointInTimeRecovery 时间点恢复
func (bm *BackupManager) PointInTimeRecovery(ctx context.Context, targetTime time.Time) error {
    // 1. 找到目标时间之前的最新完整备份
    fullBackup := bm.findLatestFullBackupBefore(targetTime)
    if fullBackup == nil {
        return fmt.Errorf("no full backup found before %v", targetTime)
    }

    // 2. 恢复完整备份
    if err := bm.restoreBackup(ctx, fullBackup); err != nil {
        return err
    }

    // 3. 应用增量备份直到目标时间
    incrementals := bm.findIncrementalBackups(fullBackup.Timestamp, targetTime)
    for _, inc := range incrementals {
        if err := bm.applyIncrementalBackup(ctx, inc); err != nil {
            return err
        }
    }

    return nil
}

func (bm *BackupManager) findLatestFullBackupBefore(t time.Time) *Backup {
    // 实现查找逻辑
    return nil
}

func (bm *BackupManager) restoreBackup(ctx context.Context, backup *Backup) error {
    // 实现恢复逻辑
    return nil
}

func (bm *BackupManager) findIncrementalBackups(start, end time.Time) []Backup {
    // 实现查找逻辑
    return nil
}

func (bm *BackupManager) applyIncrementalBackup(ctx context.Context, backup Backup) error {
    // 实现应用逻辑
    return nil
}
```

---

## 灾难恢复演练

```go
// DRDrill 灾难恢复演练
type DRDrill struct {
    manager *DisasterRecoveryManager
    scenario DrillScenario
}

type DrillScenario string
const (
    ScenarioPrimaryFailure DrillScenario = "primary_failure"
    ScenarioNetworkPartition DrillScenario = "network_partition"
    ScenarioDataCorruption DrillScenario = "data_corruption"
    ScenarioRegionOutage DrillScenario = "region_outage"
)

// RunDrill 执行演练
func (drill *DRDrill) RunDrill(ctx context.Context) (*DrillResult, error) {
    startTime := time.Now()

    // 1. 注入故障
    if err := drill.injectFault(ctx); err != nil {
        return nil, err
    }

    // 2. 检测故障
    detectionTime := time.Now()

    // 3. 执行故障转移
    if err := drill.manager.Failover(ctx, "standby-site"); err != nil {
        return nil, err
    }

    failoverTime := time.Now()

    // 4. 验证恢复
    if err := drill.verifyRecovery(ctx); err != nil {
        return nil, err
    }

    recoveryTime := time.Now()

    return &DrillResult{
        DetectionTime:   detectionTime.Sub(startTime),
        FailoverTime:    failoverTime.Sub(detectionTime),
        RecoveryTime:    recoveryTime.Sub(failoverTime),
        TotalRTO:        recoveryTime.Sub(startTime),
        DataLoss:        drill.calculateDataLoss(),
        Success:         true,
    }, nil
}

// DrillResult 演练结果
type DrillResult struct {
    DetectionTime time.Duration
    FailoverTime  time.Duration
    RecoveryTime  time.Duration
    TotalRTO      time.Duration
    DataLoss      time.Duration // RPO
    Success       bool
}
```

---

## 检查清单

- [ ] **多区域部署**: 主站点 + 至少1个备用站点
- [ ] **自动故障转移**: 健康检查 + 自动切换
- [ ] **数据复制**: RPO < 5分钟
- [ ] **定期备份**: 完整备份 + 增量备份
- [ ] **备份验证**: 定期测试备份可恢复性
- [ ] **演练计划**: 每季度至少一次DR演练
- [ ] **文档更新**: 恢复流程文档保持最新
- [ ] **RTO/RPO监控**: 实时监控恢复目标达成情况
