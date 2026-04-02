# 任务版本管理 (Task Versioning)

> **分类**: 工程与云原生  
> **标签**: #versioning #migration #compatibility

---

## 任务版本控制

```go
type TaskVersion struct {
    Version     string
    Schema      TaskSchema
    Handler     TaskHandler
    Migrate     func(oldData []byte) ([]byte, error)  // 数据迁移函数
    Deprecated  bool
    SupportedUntil time.Time
}

type VersionRegistry struct {
    versions map[string]*TaskVersion
    current  string
}

func (vr *VersionRegistry) Register(v *TaskVersion) {
    vr.versions[v.Version] = v
}

func (vr *VersionRegistry) Get(version string) (*TaskVersion, error) {
    v, ok := vr.versions[version]
    if !ok {
        return nil, fmt.Errorf("unknown task version: %s", version)
    }
    
    if v.Deprecated && time.Now().After(v.SupportedUntil) {
        return nil, fmt.Errorf("task version %s is no longer supported", version)
    }
    
    return v, nil
}

func (vr *VersionRegistry) GetCurrent() *TaskVersion {
    return vr.versions[vr.current]
}

// 版本迁移
func (vr *VersionRegistry) Migrate(oldVersion string, data []byte) ([]byte, string, error) {
    current := vr.GetCurrent()
    
    // 已经是当前版本
    if oldVersion == current.Version {
        return data, oldVersion, nil
    }
    
    // 逐版本迁移
    versions := vr.getMigrationPath(oldVersion, current.Version)
    
    for _, v := range versions {
        if v.Migrate != nil {
            migrated, err := v.Migrate(data)
            if err != nil {
                return nil, "", fmt.Errorf("migration to %s failed: %w", v.Version, err)
            }
            data = migrated
        }
    }
    
    return data, current.Version, nil
}

func (vr *VersionRegistry) getMigrationPath(from, to string) []*TaskVersion {
    // 返回从 from 到 to 的迁移路径
    // 简化版本：假设版本是线性的 v1 -> v2 -> v3
    var path []*TaskVersion
    started := false
    
    for _, v := range vr.getSortedVersions() {
        if v.Version == from {
            started = true
        }
        if started {
            path = append(path, v)
        }
        if v.Version == to {
            break
        }
    }
    
    return path
}
```

---

## 版本兼容性

```go
// 向前兼容性：新版本能处理旧数据
type ForwardCompatibleTask struct {
    version string
    data    map[string]interface{}
}

func (fct *ForwardCompatibleTask) UnmarshalJSON(data []byte) error {
    var raw map[string]interface{}
    if err := json.Unmarshal(data, &raw); err != nil {
        return err
    }
    
    // 处理版本字段
    if v, ok := raw["version"].(string); ok {
        fct.version = v
    } else {
        fct.version = "1.0"  // 默认版本
    }
    
    // 设置默认值
    fct.data = fct.applyDefaults(raw, fct.version)
    
    return nil
}

func (fct *ForwardCompatibleTask) applyDefaults(data map[string]interface{}, version string) map[string]interface{} {
    defaults := map[string]interface{}{
        "timeout": 30,
        "retry":   3,
    }
    
    // 不同版本默认值不同
    switch version {
    case "1.0":
        defaults["timeout"] = 60
    case "2.0":
        defaults["timeout"] = 30
        defaults["priority"] = "normal"
    }
    
    // 合并
    for k, v := range defaults {
        if _, ok := data[k]; !ok {
            data[k] = v
        }
    }
    
    return data
}
```

---

## 灰度发布

```go
type CanaryDeployment struct {
    newVersion string
    oldVersion string
    percentage float64  // 0-100
}

func (cd *CanaryDeployment) ShouldUseNewVersion(taskID string) bool {
    // 基于任务ID哈希决定使用哪个版本
    hash := hashString(taskID)
    return float64(hash%100) < cd.percentage
}

func (cd *CanaryDeployment) RouteTask(ctx context.Context, task *Task) error {
    if cd.ShouldUseNewVersion(task.ID) {
        return executeWithVersion(ctx, task, cd.newVersion)
    }
    return executeWithVersion(ctx, task, cd.oldVersion)
}

// 逐步放量
func (cd *CanaryDeployment) GradualRollout(targetPercentage float64, duration time.Duration) {
    steps := 10
    stepDuration := duration / time.Duration(steps)
    stepSize := (targetPercentage - cd.percentage) / float64(steps)
    
    for i := 0; i < steps; i++ {
        time.Sleep(stepDuration)
        cd.percentage += stepSize
        
        // 监控新版本的错误率
        if cd.checkErrorRate() > 0.05 {
            // 错误率过高，回滚
            cd.rollback()
            return
        }
    }
}
```
