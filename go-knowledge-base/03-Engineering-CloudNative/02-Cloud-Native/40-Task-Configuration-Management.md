# 任务配置管理 (Task Configuration Management)

> **分类**: 工程与云原生
> **标签**: #configuration #hot-reload #dynamic-config

---

## 动态配置

```go
type DynamicConfig struct {
    mu       sync.RWMutex
    config   Config
    watchers []func(Config)
}

type Config struct {
    MaxConcurrent    int
    DefaultTimeout   time.Duration
    RetryPolicy      RetryPolicy
    WorkerPoolSize   int
    QueueSize        int
}

func (dc *DynamicConfig) Load(source ConfigSource) error {
    cfg, err := source.Load()
    if err != nil {
        return err
    }

    dc.mu.Lock()
    dc.config = cfg
    dc.mu.Unlock()

    // 通知监听者
    dc.notifyWatchers(cfg)

    return nil
}

func (dc *DynamicConfig) Get() Config {
    dc.mu.RLock()
    defer dc.mu.RUnlock()
    return dc.config
}

func (dc *DynamicConfig) Watch(fn func(Config)) {
    dc.mu.Lock()
    dc.watchers = append(dc.watchers, fn)
    dc.mu.Unlock()
}

func (dc *DynamicConfig) notifyWatchers(cfg Config) {
    dc.mu.RLock()
    watchers := make([]func(Config), len(dc.watchers))
    copy(watchers, dc.watchers)
    dc.mu.RUnlock()

    for _, fn := range watchers {
        go fn(cfg)
    }
}
```

---

## 环境配置分离

```go
// 环境特定配置
type EnvironmentConfig struct {
    Development Config
    Staging     Config
    Production  Config
}

func LoadForEnvironment(env string) (Config, error) {
    // 加载基础配置
    baseConfig, _ := LoadFromFile("config/base.yaml")

    // 加载环境特定配置
    envConfig, _ := LoadFromFile(fmt.Sprintf("config/%s.yaml", env))

    // 合并配置（环境配置覆盖基础配置）
    merged := MergeConfigs(baseConfig, envConfig)

    // 环境变量覆盖
    merged = ApplyEnvOverrides(merged)

    return merged, nil
}

func ApplyEnvOverrides(cfg Config) Config {
    if v := os.Getenv("TASK_MAX_CONCURRENT"); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            cfg.MaxConcurrent = n
        }
    }

    if v := os.Getenv("TASK_TIMEOUT"); v != "" {
        if d, err := time.ParseDuration(v); err == nil {
            cfg.DefaultTimeout = d
        }
    }

    return cfg
}
```

---

## 配置验证

```go
type ConfigValidator struct {
    rules []ValidationRule
}

type ValidationRule struct {
    Field   string
    Check   func(interface{}) error
    Message string
}

func (cv *ConfigValidator) Validate(cfg Config) error {
    var errs []error

    if cfg.MaxConcurrent <= 0 {
        errs = append(errs, fmt.Errorf("max_concurrent must be > 0"))
    }

    if cfg.DefaultTimeout < time.Second {
        errs = append(errs, fmt.Errorf("default_timeout must be >= 1s"))
    }

    if cfg.WorkerPoolSize < cfg.MaxConcurrent {
        errs = append(errs, fmt.Errorf("worker_pool_size must be >= max_concurrent"))
    }

    if cfg.QueueSize <= 0 {
        errs = append(errs, fmt.Errorf("queue_size must be > 0"))
    }

    if len(errs) > 0 {
        return &ValidationError{Errors: errs}
    }

    return nil
}

func (cv *ConfigValidator) ValidateWithContext(ctx context.Context, cfg Config) error {
    // 检查依赖服务可用性
    if err := cv.checkDatabase(ctx, cfg.Database); err != nil {
        return fmt.Errorf("database check failed: %w", err)
    }

    if err := cv.checkRedis(ctx, cfg.Redis); err != nil {
        return fmt.Errorf("redis check failed: %w", err)
    }

    if err := cv.checkQueue(ctx, cfg.Queue); err != nil {
        return fmt.Errorf("queue check failed: %w", err)
    }

    return nil
}
```

---

## 配置热重载

```go
type HotReloader struct {
    config   *DynamicConfig
    sources  []ConfigSource
    interval time.Duration
}

func (hr *HotReloader) Start(ctx context.Context) {
    ticker := time.NewTicker(hr.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            hr.checkAndReload()
        }
    }
}

func (hr *HotReloader) checkAndReload() {
    // 检查配置是否变化
    checksum, err := hr.calculateChecksum()
    if err != nil {
        log.Printf("failed to calculate checksum: %v", err)
        return
    }

    if checksum == hr.lastChecksum {
        return  // 无变化
    }

    // 重新加载
    for _, source := range hr.sources {
        if err := hr.config.Load(source); err != nil {
            log.Printf("failed to reload config from %s: %v", source.Name(), err)
            continue
        }
    }

    hr.lastChecksum = checksum
    log.Println("config reloaded successfully")
}

func (hr *HotReloader) calculateChecksum() (string, error) {
    h := sha256.New()

    for _, source := range hr.sources {
        data, err := source.Raw()
        if err != nil {
            return "", err
        }
        h.Write(data)
    }

    return hex.EncodeToString(h.Sum(nil)), nil
}
```
