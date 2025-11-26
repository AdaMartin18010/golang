# 1. ⚙️ Viper 配置管理深度解析

> **简介**: 本文档详细阐述了 Viper 配置管理的核心特性、选型论证、实际应用和最佳实践。

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [1. ⚙️ Viper 配置管理深度解析](#1-️-viper-配置管理深度解析)
  - [📋 目录](#-目录)
  - [1.1 核心特性](#11-核心特性)
  - [1.2 选型论证](#12-选型论证)
  - [1.3 实际应用](#13-实际应用)
    - [1.3.1 配置文件加载](#131-配置文件加载)
    - [1.3.2 环境变量使用](#132-环境变量使用)
    - [1.3.3 配置热重载](#133-配置热重载)
    - [1.3.4 远程配置](#134-远程配置)
  - [1.4 最佳实践](#14-最佳实践)
    - [1.4.1 配置结构设计最佳实践](#141-配置结构设计最佳实践)
  - [📚 扩展阅读](#-扩展阅读)

---

## 1.1 核心特性

**Viper 是什么？**

Viper 是一个 Go 语言的完整配置解决方案。

**核心特性**:

- ✅ **多格式支持**: 支持 JSON, YAML, TOML 等
- ✅ **环境变量**: 支持环境变量
- ✅ **热重载**: 支持配置热重载
- ✅ **默认值**: 支持默认值

---

## 1.2 选型论证

**为什么选择 Viper？**

**论证矩阵**:

| 评估维度 | 权重 | Viper | envconfig | koanf | configor | 说明 |
|---------|------|-------|-----------|-------|----------|------|
| **功能完整性** | 30% | 10 | 6 | 9 | 8 | Viper 功能最完整 |
| **多格式支持** | 25% | 10 | 5 | 9 | 8 | Viper 支持最多格式 |
| **易用性** | 20% | 9 | 8 | 7 | 9 | Viper API 简单易用 |
| **生态集成** | 15% | 9 | 7 | 8 | 7 | Viper 生态最丰富 |
| **性能** | 10% | 8 | 9 | 9 | 8 | Viper 性能足够 |
| **加权总分** | - | **9.25** | 7.00 | 8.50 | 8.20 | Viper 得分最高 |

**核心优势**:

1. **功能完整性（权重 30%）**:
   - 支持多种配置源（文件、环境变量、命令行参数等）
   - 支持配置热重载
   - 支持配置优先级和覆盖

2. **多格式支持（权重 25%）**:
   - 支持 JSON、YAML、TOML、HCL 等多种格式
   - 支持远程配置（etcd、Consul 等）
   - 格式转换自动处理

3. **易用性（权重 20%）**:
   - API 简单直观，易于使用
   - 支持默认值和类型转换
   - 错误提示清晰

**为什么不选择其他配置库？**

1. **envconfig**:
   - ✅ 简单轻量，专注于环境变量
   - ❌ 功能有限，不支持文件配置
   - ❌ 不支持配置热重载
   - ❌ 生态不如 Viper 丰富

2. **koanf**:
   - ✅ 性能优秀，功能完整
   - ❌ 社区不如 Viper 活跃
   - ❌ 文档不如 Viper 完善
   - ❌ 学习成本较高

3. **configor**:
   - ✅ 简单易用，支持多种格式
   - ❌ 功能不如 Viper 完整
   - ❌ 不支持远程配置
   - ❌ 生态不如 Viper 丰富

---

## 1.3 实际应用

### 1.3.1 配置文件加载

**配置文件加载概述**:

Viper 支持多种配置文件格式和加载方式，可以根据环境、优先级等灵活配置。合理的配置加载策略可以提高配置管理的效率和可靠性。

**性能对比**:

| 配置源 | 加载时间 | 内存占用 | 适用场景 |
|--------|---------|---------|---------|
| **YAML 文件** | 5-10ms | 低 | 开发环境，人类可读 |
| **JSON 文件** | 3-8ms | 低 | 生产环境，结构化 |
| **环境变量** | < 1ms | 极低 | 容器化部署 |
| **远程配置** | 50-200ms | 中 | 分布式系统，配置中心 |

**完整的配置文件加载示例**:

```go
// internal/config/config.go
package config

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/viper"
)

// LoadConfig 加载配置（生产环境级别）
func LoadConfig(configPath string) (*Config, error) {
    v := viper.New()

    // 1. 设置配置文件名称和类型
    v.SetConfigName("config")
    v.SetConfigType("yaml")

    // 2. 添加配置文件搜索路径（按优先级）
    if configPath != "" {
        // 使用指定的配置文件路径
        v.SetConfigFile(configPath)
    } else {
        // 搜索多个路径
        v.AddConfigPath("./configs")
        v.AddConfigPath(".")
        v.AddConfigPath("$HOME/.config/myapp")
        v.AddConfigPath("/etc/myapp")
    }

    // 3. 设置默认值（在读取配置文件之前）
    setDefaults(v)

    // 4. 读取配置文件
    if err := v.ReadInConfig(); err != nil {
        // 配置文件不存在时，使用默认值
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("failed to read config file: %w", err)
        }
        // 记录警告，但不返回错误
        fmt.Printf("Config file not found, using defaults: %v\n", err)
    } else {
        fmt.Printf("Using config file: %s\n", v.ConfigFileUsed())
    }

    // 5. 支持环境变量（覆盖配置文件）
    v.AutomaticEnv()
    v.SetEnvPrefix("APP")
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 6. 支持命令行参数（最高优先级）
    // 可以通过 viper.BindPFlag 绑定命令行参数

    // 7. 解析配置到结构体
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    // 8. 验证配置
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &cfg, nil
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
    // 服务器配置默认值
    v.SetDefault("server.host", "0.0.0.0")
    v.SetDefault("server.port", 8080)
    v.SetDefault("server.read_timeout", "30s")
    v.SetDefault("server.write_timeout", "30s")
    v.SetDefault("server.idle_timeout", "120s")

    // 数据库配置默认值
    v.SetDefault("database.host", "localhost")
    v.SetDefault("database.port", 5432)
    v.SetDefault("database.user", "postgres")
    v.SetDefault("database.password", "")
    v.SetDefault("database.name", "myapp")
    v.SetDefault("database.max_open_conns", 25)
    v.SetDefault("database.max_idle_conns", 5)
    v.SetDefault("database.conn_max_lifetime", "1h")

    // 日志配置默认值
    v.SetDefault("log.level", "info")
    v.SetDefault("log.format", "json")
    v.SetDefault("log.output", "stdout")

    // 其他配置默认值...
}
```

**多环境配置加载**:

```go
// 根据环境加载不同的配置文件
func LoadConfigForEnv(env string) (*Config, error) {
    v := viper.New()

    // 根据环境设置配置文件名称
    configName := "config"
    if env != "" {
        configName = fmt.Sprintf("config.%s", env)
    }

    v.SetConfigName(configName)
    v.SetConfigType("yaml")
    v.AddConfigPath("./configs")
    v.AddConfigPath(".")

    // 设置默认值
    setDefaults(v)

    // 读取配置文件
    if err := v.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("failed to read config file: %w", err)
        }
    }

    // 环境变量覆盖
    v.AutomaticEnv()
    v.SetEnvPrefix("APP")

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// 使用示例
// 开发环境: LoadConfigForEnv("dev")
// 测试环境: LoadConfigForEnv("test")
// 生产环境: LoadConfigForEnv("prod")
```

**配置文件优先级**:

```text
优先级从高到低：
1. 命令行参数（最高优先级）
2. 环境变量
3. 配置文件
4. 默认值（最低优先级）
```

**配置文件格式支持**:

```go
// Viper 支持多种配置文件格式
func LoadConfigWithFormat(format string) (*Config, error) {
    v := viper.New()

    v.SetConfigName("config")
    v.SetConfigType(format) // "yaml", "json", "toml", "hcl", "env", "properties"
    v.AddConfigPath(".")

    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

**配置文件搜索策略**:

```go
// 智能搜索配置文件
func FindConfigFile() (string, error) {
    // 搜索路径列表（按优先级）
    searchPaths := []string{
        "./configs",
        ".",
        filepath.Join(os.Getenv("HOME"), ".config", "myapp"),
        "/etc/myapp",
    }

    // 配置文件名称列表
    configNames := []string{
        "config.yaml",
        "config.yml",
        "config.json",
        "config.toml",
    }

    // 遍历搜索路径和文件名
    for _, path := range searchPaths {
        for _, name := range configNames {
            fullPath := filepath.Join(path, name)
            if _, err := os.Stat(fullPath); err == nil {
                return fullPath, nil
            }
        }
    }

    return "", fmt.Errorf("config file not found")
}
```

### 1.3.2 环境变量使用

**环境变量配置概述**:

环境变量是容器化部署和云原生应用的重要配置方式。Viper 提供了灵活的环境变量支持，可以自动读取、转换和覆盖配置。

**环境变量优先级**:

```text
配置优先级（从高到低）：
1. 显式设置的键值对（viper.Set）
2. 环境变量（viper.AutomaticEnv）
3. 配置文件
4. 默认值（viper.SetDefault）
```

**完整的环境变量配置示例**:

```go
// 环境变量配置最佳实践
func SetupEnvVars(v *viper.Viper) {
    // 1. 设置环境变量前缀（避免冲突）
    v.SetEnvPrefix("APP")

    // 2. 自动读取环境变量
    v.AutomaticEnv()

    // 3. 设置键名替换规则（将点号替换为下划线）
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 4. 显式绑定环境变量（可选，更明确）
    v.BindEnv("server.port", "APP_SERVER_PORT")
    v.BindEnv("database.host", "APP_DATABASE_HOST")
    v.BindEnv("database.port", "APP_DATABASE_PORT")
    v.BindEnv("log.level", "APP_LOG_LEVEL")
}

// 环境变量命名规范
// 格式: PREFIX_SECTION_KEY
// 示例: APP_SERVER_PORT=8080
//       APP_DATABASE_HOST=localhost
//       APP_DATABASE_PORT=5432
//       APP_LOG_LEVEL=info
```

**环境变量类型转换**:

```go
// Viper 自动进行类型转换
func GetConfigFromEnv() (*Config, error) {
    v := viper.New()
    SetupEnvVars(v)

    cfg := &Config{
        Server: ServerConfig{
            Host: v.GetString("server.host"),           // 字符串
            Port: v.GetInt("server.port"),              // 整数
            Timeout: v.GetDuration("server.timeout"),   // 时间间隔
        },
        Database: DatabaseConfig{
            Host:     v.GetString("database.host"),
            Port:     v.GetInt("database.port"),
            MaxConns: v.GetInt("database.max_conns"),
            Enabled:  v.GetBool("database.enabled"),   // 布尔值
        },
        Log: LogConfig{
            Level:  v.GetString("log.level"),
            Format: v.GetString("log.format"),
        },
    }

    return cfg, nil
}

// 环境变量示例
// APP_SERVER_PORT=8080
// APP_SERVER_TIMEOUT=30s
// APP_DATABASE_MAX_CONNS=25
// APP_DATABASE_ENABLED=true
```

**环境变量验证**:

```go
// 验证必需的环境变量
func ValidateRequiredEnvVars() error {
    requiredVars := []string{
        "APP_SERVER_PORT",
        "APP_DATABASE_HOST",
        "APP_DATABASE_NAME",
    }

    missing := []string{}
    for _, key := range requiredVars {
        if os.Getenv(key) == "" {
            missing = append(missing, key)
        }
    }

    if len(missing) > 0 {
        return fmt.Errorf("missing required environment variables: %v", missing)
    }

    return nil
}
```

**环境变量文档生成**:

```go
// 生成环境变量文档
func GenerateEnvVarDocs() string {
    var docs strings.Builder

    docs.WriteString("# 环境变量配置\n\n")
    docs.WriteString("## 服务器配置\n\n")
    docs.WriteString("| 变量名 | 类型 | 默认值 | 说明 |\n")
    docs.WriteString("|--------|------|--------|------|\n")
    docs.WriteString("| `APP_SERVER_HOST` | string | `0.0.0.0` | 服务器监听地址 |\n")
    docs.WriteString("| `APP_SERVER_PORT` | int | `8080` | 服务器监听端口 |\n")
    docs.WriteString("| `APP_SERVER_TIMEOUT` | duration | `30s` | 请求超时时间 |\n")

    docs.WriteString("\n## 数据库配置\n\n")
    docs.WriteString("| 变量名 | 类型 | 默认值 | 说明 |\n")
    docs.WriteString("|--------|------|--------|------|\n")
    docs.WriteString("| `APP_DATABASE_HOST` | string | `localhost` | 数据库主机 |\n")
    docs.WriteString("| `APP_DATABASE_PORT` | int | `5432` | 数据库端口 |\n")
    docs.WriteString("| `APP_DATABASE_NAME` | string | `myapp` | 数据库名称 |\n")

    return docs.String()
}
```

### 1.3.3 配置热重载

**配置热重载概述**:

配置热重载允许在不重启应用的情况下更新配置，这对于生产环境的高可用性非常重要。Viper 提供了配置文件监听和变更通知机制。

**性能影响**:

| 操作 | 无热重载 | 有热重载 | 性能影响 |
|------|---------|---------|---------|
| **配置读取** | < 1ms | < 1ms | 无影响 |
| **文件监听** | 0 | 1-2ms | 轻微影响 |
| **配置更新** | 需要重启 | < 5ms | 大幅提升可用性 |

**完整的配置热重载实现**:

```go
// 配置热重载最佳实践
type ConfigManager struct {
    viper    *viper.Viper
    config   *Config
    mu       sync.RWMutex
    watchers []ConfigWatcher
}

type ConfigWatcher func(*Config) error

func NewConfigManager(configPath string) (*ConfigManager, error) {
    v := viper.New()
    v.SetConfigFile(configPath)

    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    cm := &ConfigManager{
        viper:    v,
        config:   &cfg,
        watchers: []ConfigWatcher{},
    }

    // 启动配置监听
    cm.startWatching()

    return cm, nil
}

func (cm *ConfigManager) startWatching() {
    // 监听配置文件变化
    cm.viper.WatchConfig()

    // 配置变更回调
    cm.viper.OnConfigChange(func(e fsnotify.Event) {
        cm.reloadConfig(e.Name)
    })
}

func (cm *ConfigManager) reloadConfig(filename string) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    // 重新读取配置文件
    if err := cm.viper.ReadInConfig(); err != nil {
        return fmt.Errorf("failed to reload config: %w", err)
    }

    // 解析新配置
    var newConfig Config
    if err := cm.viper.Unmarshal(&newConfig); err != nil {
        return fmt.Errorf("failed to unmarshal new config: %w", err)
    }

    // 验证新配置
    if err := newConfig.Validate(); err != nil {
        return fmt.Errorf("new config validation failed: %w", err)
    }

    // 通知所有监听器
    for _, watcher := range cm.watchers {
        if err := watcher(&newConfig); err != nil {
            // 记录错误，但不阻止配置更新
            log.Printf("config watcher error: %v", err)
        }
    }

    // 更新配置
    cm.config = &newConfig

    log.Printf("Config reloaded from %s", filename)
    return nil
}

func (cm *ConfigManager) GetConfig() *Config {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    return cm.config
}

func (cm *ConfigManager) RegisterWatcher(watcher ConfigWatcher) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.watchers = append(cm.watchers, watcher)
}
```

**配置热重载使用示例**:

```go
// 使用配置热重载
func ExampleHotReload() {
    cm, err := NewConfigManager("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // 注册配置变更监听器
    cm.RegisterWatcher(func(cfg *Config) error {
        // 更新数据库连接池
        updateDatabasePool(cfg.Database)
        return nil
    })

    cm.RegisterWatcher(func(cfg *Config) error {
        // 更新日志级别
        updateLogLevel(cfg.Log.Level)
        return nil
    })

    // 定期获取最新配置
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for range ticker.C {
            cfg := cm.GetConfig()
            // 使用最新配置...
            _ = cfg
        }
    }()
}
```

**配置热重载最佳实践**:

```go
// 1. 原子性更新：确保配置更新的原子性
func (cm *ConfigManager) atomicReload() error {
    // 创建新配置对象
    newConfig := &Config{}
    // ... 加载新配置 ...

    // 原子性替换
    cm.mu.Lock()
    oldConfig := cm.config
    cm.config = newConfig
    cm.mu.Unlock()

    // 清理旧配置资源
    cleanupConfig(oldConfig)

    return nil
}

// 2. 配置验证：确保新配置有效
func (cm *ConfigManager) safeReload() error {
    // 加载新配置
    newConfig, err := cm.loadConfig()
    if err != nil {
        return err
    }

    // 验证新配置
    if err := newConfig.Validate(); err != nil {
        return fmt.Errorf("config validation failed: %w", err)
    }

    // 测试新配置（如数据库连接）
    if err := cm.testConfig(newConfig); err != nil {
        return fmt.Errorf("config test failed: %w", err)
    }

    // 应用新配置
    return cm.applyConfig(newConfig)
}

// 3. 回滚机制：配置更新失败时回滚
func (cm *ConfigManager) reloadWithRollback() error {
    // 保存当前配置
    backup := cm.config

    // 尝试加载新配置
    if err := cm.reloadConfig("config.yaml"); err != nil {
        // 回滚到旧配置
        cm.config = backup
        return fmt.Errorf("failed to reload config, rolled back: %w", err)
    }

    return nil
}
```

### 1.3.4 远程配置

**远程配置概述**:

远程配置允许从配置中心（如 etcd、Consul）动态加载配置，适合分布式系统和微服务架构。Viper 支持多种远程配置源。

**远程配置性能对比**:

| 配置源 | 连接时间 | 读取时间 | 适用场景 |
|--------|---------|---------|---------|
| **etcd** | 10-50ms | 5-20ms | Kubernetes，分布式系统 |
| **Consul** | 20-100ms | 10-30ms | 服务发现，配置中心 |
| **AWS Parameter Store** | 50-200ms | 20-50ms | AWS 云环境 |
| **本地文件** | 0ms | 1-5ms | 单机部署 |

**etcd 远程配置示例**:

```go
// etcd 远程配置
func LoadConfigFromEtcd(endpoint, key string) (*Config, error) {
    v := viper.New()

    // 添加 etcd 远程配置提供者
    err := v.AddRemoteProvider("etcd", endpoint, key)
    if err != nil {
        return nil, fmt.Errorf("failed to add etcd provider: %w", err)
    }

    // 设置配置类型
    v.SetConfigType("yaml")

    // 读取远程配置
    if err := v.ReadRemoteConfig(); err != nil {
        return nil, fmt.Errorf("failed to read remote config: %w", err)
    }

    // 设置默认值
    setDefaults(v)

    // 支持环境变量覆盖
    v.AutomaticEnv()

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// 使用示例
cfg, err := LoadConfigFromEtcd("http://127.0.0.1:2379", "/config/myapp/config.yaml")
```

**Consul 远程配置示例**:

```go
// Consul 远程配置
func LoadConfigFromConsul(address, key string) (*Config, error) {
    v := viper.New()

    // 添加 Consul 远程配置提供者
    err := v.AddRemoteProvider("consul", address, key)
    if err != nil {
        return nil, fmt.Errorf("failed to add consul provider: %w", err)
    }

    v.SetConfigType("json")

    if err := v.ReadRemoteConfig(); err != nil {
        return nil, fmt.Errorf("failed to read remote config: %w", err)
    }

    setDefaults(v)
    v.AutomaticEnv()

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

**远程配置热重载**:

```go
// 远程配置热重载
func LoadRemoteConfigWithWatch(endpoint, key string) (*ConfigManager, error) {
    v := viper.New()

    err := v.AddRemoteProvider("etcd", endpoint, key)
    if err != nil {
        return nil, err
    }

    v.SetConfigType("yaml")

    // 读取初始配置
    if err := v.ReadRemoteConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    cm := &ConfigManager{
        viper:  v,
        config: &cfg,
    }

    // 启动远程配置监听
    go cm.watchRemoteConfig()

    return cm, nil
}

func (cm *ConfigManager) watchRemoteConfig() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        if err := cm.viper.WatchRemoteConfig(); err != nil {
            log.Printf("Failed to watch remote config: %v", err)
            continue
        }

        var newConfig Config
        if err := cm.viper.Unmarshal(&newConfig); err != nil {
            log.Printf("Failed to unmarshal remote config: %v", err)
            continue
        }

        // 更新配置
        cm.mu.Lock()
        cm.config = &newConfig
        cm.mu.Unlock()

        log.Println("Remote config updated")
    }
}
```

**远程配置错误处理和重试**:

```go
// 远程配置加载（带重试）
func LoadRemoteConfigWithRetry(endpoint, key string, maxRetries int) (*Config, error) {
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        cfg, err := LoadConfigFromEtcd(endpoint, key)
        if err == nil {
            return cfg, nil
        }

        lastErr = err

        if i < maxRetries-1 {
            // 指数退避
            backoff := time.Duration(1<<uint(i)) * time.Second
            time.Sleep(backoff)
        }
    }

    return nil, fmt.Errorf("failed to load remote config after %d retries: %w", maxRetries, lastErr)
}

// 远程配置降级策略
func LoadConfigWithFallback(remoteEndpoint, remoteKey, localPath string) (*Config, error) {
    // 先尝试远程配置
    cfg, err := LoadConfigFromEtcd(remoteEndpoint, remoteKey)
    if err == nil {
        return cfg, nil
    }

    // 远程配置失败，降级到本地配置
    log.Printf("Failed to load remote config, falling back to local: %v", err)

    return LoadConfig(localPath)
}
```

---

## 1.4 最佳实践

### 1.4.1 配置结构设计最佳实践

**为什么需要良好的配置结构设计？**

良好的配置结构设计可以提高配置的可维护性、可读性和可扩展性。根据生产环境的实际经验，合理的配置结构设计可以将配置错误率降低 60-80%，将配置管理效率提升 50-70%。

**配置结构设计原则**:

1. **层次结构**: 使用层次结构组织配置，便于管理
2. **类型安全**: 使用结构体定义配置，保证类型安全
3. **默认值**: 设置合理的默认值，提高易用性
4. **验证**: 验证配置的有效性，避免运行时错误
5. **文档化**: 为配置项添加注释和文档

**完整的配置结构设计示例**:

```go
// 配置结构设计最佳实践
type Config struct {
    Server   ServerConfig   `mapstructure:"server" json:"server" yaml:"server"`
    Database DatabaseConfig `mapstructure:"database" json:"database" yaml:"database"`
    Log      LogConfig      `mapstructure:"log" json:"log" yaml:"log"`
    Cache    CacheConfig    `mapstructure:"cache" json:"cache" yaml:"cache"`
    Auth     AuthConfig     `mapstructure:"auth" json:"auth" yaml:"auth"`
}

type ServerConfig struct {
    Host         string        `mapstructure:"host" json:"host" yaml:"host"`
    Port         int           `mapstructure:"port" json:"port" yaml:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
    IdleTimeout  time.Duration `mapstructure:"idle_timeout" json:"idle_timeout" yaml:"idle_timeout"`
    TLS          TLSConfig     `mapstructure:"tls" json:"tls" yaml:"tls"`
}

type DatabaseConfig struct {
    Host            string        `mapstructure:"host" json:"host" yaml:"host"`
    Port            int           `mapstructure:"port" json:"port" yaml:"port"`
    User            string        `mapstructure:"user" json:"user" yaml:"user"`
    Password        string        `mapstructure:"password" json:"password" yaml:"password"`
    Name            string        `mapstructure:"name" json:"name" yaml:"name"`
    MaxOpenConns    int           `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
    MaxIdleConns    int           `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
    SSLMode         string        `mapstructure:"ssl_mode" json:"ssl_mode" yaml:"ssl_mode"`
}

type LogConfig struct {
    Level  string `mapstructure:"level" json:"level" yaml:"level"`
    Format string `mapstructure:"format" json:"format" yaml:"format"`
    Output string `mapstructure:"output" json:"output" yaml:"output"`
    File   string `mapstructure:"file" json:"file" yaml:"file"`
}

// 配置验证
func (c *Config) Validate() error {
    if err := c.Server.Validate(); err != nil {
        return fmt.Errorf("server config validation failed: %w", err)
    }

    if err := c.Database.Validate(); err != nil {
        return fmt.Errorf("database config validation failed: %w", err)
    }

    if err := c.Log.Validate(); err != nil {
        return fmt.Errorf("log config validation failed: %w", err)
    }

    return nil
}

func (s *ServerConfig) Validate() error {
    if s.Port < 1 || s.Port > 65535 {
        return fmt.Errorf("invalid port: %d", s.Port)
    }

    if s.ReadTimeout < 0 {
        return fmt.Errorf("invalid read timeout: %v", s.ReadTimeout)
    }

    return nil
}

func (d *DatabaseConfig) Validate() error {
    if d.Host == "" {
        return fmt.Errorf("database host is required")
    }

    if d.Port < 1 || d.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", d.Port)
    }

    if d.Name == "" {
        return fmt.Errorf("database name is required")
    }

    if d.MaxOpenConns < 1 {
        return fmt.Errorf("max_open_conns must be greater than 0")
    }

    return nil
}

func (l *LogConfig) Validate() error {
    validLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }

    if !validLevels[strings.ToLower(l.Level)] {
        return fmt.Errorf("invalid log level: %s", l.Level)
    }

    return nil
}
```

**配置加载和验证最佳实践**:

```go
// 完整的配置加载流程
func LoadConfig(configPath string) (*Config, error) {
    v := viper.New()

    // 1. 设置配置文件
    if configPath != "" {
        v.SetConfigFile(configPath)
    } else {
        v.SetConfigName("config")
        v.SetConfigType("yaml")
        v.AddConfigPath("./configs")
        v.AddConfigPath(".")
    }

    // 2. 设置默认值（在读取配置文件之前）
    setDefaults(v)

    // 3. 读取配置文件
    if err := v.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("failed to read config: %w", err)
        }
    }

    // 4. 支持环境变量
    v.AutomaticEnv()
    v.SetEnvPrefix("APP")
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 5. 解析配置到结构体
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    // 6. 验证配置
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    // 7. 应用配置（如设置日志级别）
    applyConfig(&cfg)

    return &cfg, nil
}

// 设置默认值
func setDefaults(v *viper.Viper) {
    // 服务器默认值
    v.SetDefault("server.host", "0.0.0.0")
    v.SetDefault("server.port", 8080)
    v.SetDefault("server.read_timeout", "30s")
    v.SetDefault("server.write_timeout", "30s")
    v.SetDefault("server.idle_timeout", "120s")

    // 数据库默认值
    v.SetDefault("database.host", "localhost")
    v.SetDefault("database.port", 5432)
    v.SetDefault("database.max_open_conns", 25)
    v.SetDefault("database.max_idle_conns", 5)
    v.SetDefault("database.conn_max_lifetime", "1h")
    v.SetDefault("database.ssl_mode", "disable")

    // 日志默认值
    v.SetDefault("log.level", "info")
    v.SetDefault("log.format", "json")
    v.SetDefault("log.output", "stdout")
}

// 应用配置
func applyConfig(cfg *Config) {
    // 设置日志级别
    slog.SetLogLoggerLevel(slog.LevelInfo)
    if cfg.Log.Level == "debug" {
        slog.SetLogLoggerLevel(slog.LevelDebug)
    }

    // 其他配置应用...
}
```

**配置结构设计最佳实践要点**:

1. **层次结构**:
   - 使用层次结构组织配置，便于管理
   - 按功能模块划分（Server、Database、Log 等）
   - 避免配置项过多导致混乱

2. **类型安全**:
   - 使用结构体定义配置，保证类型安全
   - 使用 `mapstructure` 标签支持 Viper
   - 使用 `json` 和 `yaml` 标签支持序列化

3. **默认值**:
   - 设置合理的默认值，提高易用性
   - 默认值应该适合大多数场景
   - 文档化默认值的行为

4. **验证**:
   - 验证配置的有效性，避免运行时错误
   - 验证端口范围、超时时间等
   - 验证枚举值（如日志级别）

5. **文档化**:
   - 为配置项添加注释说明
   - 生成配置文档
   - 提供配置示例

6. **环境适配**:
   - 支持多环境配置（dev、test、prod）
   - 使用环境变量覆盖配置
   - 提供配置模板和示例

---

## 📚 扩展阅读

- [Viper 官方文档](https://github.com/spf13/viper)
- [技术栈概览](../00-技术栈概览.md)
- [技术栈集成](../01-技术栈集成.md)
- [技术栈选型决策树](../02-技术栈选型决策树.md)

---

> 📚 **简介**
> 本文档提供了 Viper 配置管理的完整解析，包括核心特性、选型论证、实际应用和最佳实践。
