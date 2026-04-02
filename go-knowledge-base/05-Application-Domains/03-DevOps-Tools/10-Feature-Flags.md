# 特性开关 (Feature Flags)

> **分类**: 成熟应用领域

---

## 基础实现

```go
type FeatureFlags struct {
    flags map[string]bool
    mu    sync.RWMutex
}

func New() *FeatureFlags {
    return &FeatureFlags{
        flags: make(map[string]bool),
    }
}

func (f *FeatureFlags) Enable(name string) {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.flags[name] = true
}

func (f *FeatureFlags) IsEnabled(name string) bool {
    f.mu.RLock()
    defer f.mu.RUnlock()
    return f.flags[name]
}

// 使用
if flags.IsEnabled("new-checkout") {
    newCheckout.Process()
} else {
    oldCheckout.Process()
}
```

---

## LaunchDarkly

```go
import "github.com/launchdarkly/go-server-sdk/v6"

client, _ := ld.MakeClient("sdk-key", 5*time.Second)

flag, _ := client.BoolVariation("new-feature", user, false)
if flag {
    // 新功能
}
```

---

## 基于百分比的推出

```go
type GradualRollout struct {
    percentage int
}

func (r *GradualRollout) IsEnabled(userID string) bool {
    hash := fnv32(userID)
    return hash%100 < uint32(r.percentage)
}

func fnv32(s string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(s))
    return h.Sum32()
}
```

---

## 最佳实践

1. **可观测性**: 记录特性开关使用情况
2. **清理**: 功能稳定后移除开关
3. **测试**: 测试两种代码路径
4. **回滚**: 快速禁用问题功能
