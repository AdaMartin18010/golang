# 任务系统迁移指南 (Task System Migration Guide)

> **分类**: 工程与云原生
> **标签**: #migration #upgrade #backward-compatibility

---

## 版本迁移策略

```go
// 版本兼容性处理器
type VersionedTaskHandler struct {
    handlers map[int]TaskHandler  // version -> handler
    current  int
}

func (vth *VersionedTaskHandler) Handle(ctx context.Context, task *Task) error {
    version := task.Version
    if version == 0 {
        version = vth.current
    }

    handler, ok := vth.handlers[version]
    if !ok {
        return fmt.Errorf("unsupported task version: %d", version)
    }

    // 如果需要，升级到最新版本
    if version < vth.current {
        task = vth.migrateTask(task, version, vth.current)
    }

    return handler.Handle(ctx, task)
}

func (vth *VersionedTaskHandler) migrateTask(task *Task, from, to int) *Task {
    for v := from; v < to; v++ {
        task = vth.migrations[v](task)
        task.Version = v + 1
    }
    return task
}
```

---

## 数据结构迁移

```go
// V1 -> V2 迁移
type TaskMigrationV1ToV2 struct{}

func (m *TaskMigrationV1ToV2) Migrate(old TaskV1) TaskV2 {
    return TaskV2{
        ID:        old.ID,
        Name:      old.Name,
        Type:      old.Type,
        Payload:   m.migratePayload(old.Payload),
        Priority:  m.convertPriority(old.Priority),
        Schedule:  old.Schedule,
        CreatedAt: old.CreatedAt,
        // V2 新字段
        Metadata: map[string]string{
            "migrated_from": "v1",
            "migrated_at":   time.Now().Format(time.RFC3339),
        },
    }
}

func (m *TaskMigrationV1ToV2) migratePayload(old []byte) []byte {
    // 转换旧格式到新格式
    var oldPayload OldPayload
    json.Unmarshal(old, &oldPayload)

    newPayload := NewPayload{
        Data: oldPayload.Data,
        // 字段映射
    }

    new, _ := json.Marshal(newPayload)
    return new
}

func (m *TaskMigrationV1ToV2) convertPriority(old string) int {
    // 字符串优先级转数字
    switch old {
    case "high":
        return 10
    case "medium":
        return 5
    case "low":
        return 1
    default:
        return 5
    }
}
```

---

## 双写迁移

```go
type DualWriteMigration struct {
    oldStore TaskStore
    newStore TaskStore
    progress ProgressTracker
}

func (dwm *DualWriteMigration) Start() {
    // 阶段1: 双写（写入新旧两个系统）
    dwm.enableDualWrite()

    // 阶段2: 数据迁移
    dwm.migrateData()

    // 阶段3: 验证
    if err := dwm.verify(); err != nil {
        log.Fatal("verification failed:", err)
    }

    // 阶段4: 切换读
    dwm.switchReadToNew()

    // 阶段5: 停止写入旧系统
    dwm.disableOldWrite()

    // 阶段6: 清理
    dwm.cleanup()
}

func (dwm *DualWriteMigration) migrateData() error {
    cursor := ""

    for {
        tasks, nextCursor, err := dwm.oldStore.List(cursor, 100)
        if err != nil {
            return err
        }

        for _, task := range tasks {
            // 转换任务
            newTask := dwm.convertTask(task)

            // 写入新系统
            if err := dwm.newStore.Create(context.Background(), newTask); err != nil {
                return err
            }

            // 记录进度
            dwm.progress.Record(task.ID)
        }

        if nextCursor == "" {
            break
        }
        cursor = nextCursor
    }

    return nil
}
```

---

## 回滚策略

```go
type MigrationRollback struct {
    backup BackupStore
}

func (mr *MigrationRollback) CreateBackup(ctx context.Context) error {
    // 创建数据快照
    snapshot := &DataSnapshot{
        Timestamp: time.Now(),
        Tasks:     mr.dumpAllTasks(),
        Config:    mr.dumpConfig(),
    }

    return mr.backup.Save(ctx, snapshot)
}

func (mr *MigrationRollback) Rollback(ctx context.Context, snapshotID string) error {
    // 加载备份
    snapshot, err := mr.backup.Load(ctx, snapshotID)
    if err != nil {
        return err
    }

    // 停止当前系统
    mr.stopSystem()

    // 恢复数据
    if err := mr.restoreTasks(snapshot.Tasks); err != nil {
        return err
    }

    // 恢复配置
    mr.restoreConfig(snapshot.Config)

    // 重启系统
    mr.startSystem()

    return nil
}
```

---

## 兼容性测试

```go
func TestVersionCompatibility(t *testing.T) {
    // 测试旧客户端能否与新服务端通信
    oldClient := NewOldClient()
    newServer := NewNewServer()

    // 提交旧格式任务
    oldTask := OldTaskFormat{
        Name: "test",
        Data: "data",
    }

    taskID, err := oldClient.Submit(oldTask)
    if err != nil {
        t.Fatal(err)
    }

    // 验证新系统能正确处理
    status, err := newServer.GetStatus(taskID)
    if err != nil {
        t.Fatal(err)
    }

    if status != TaskStatusCompleted {
        t.Errorf("unexpected status: %s", status)
    }
}
```
