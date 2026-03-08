# 配置集成指南

## 📋 概述

本文档说明如何将可观测性配置集成到应用的统一配置系统中。

## 🚀 快速开始

### 1. 配置文件方式

```yaml
# configs/config.yaml
observability:
  otlp:
    endpoint: "localhost:4317"
    insecure: true
    service_name: "my-service"
    service_version: "v1.0.0"
  system:
    enabled: true
    collect_interval: 5s
    enable_disk_monitor: true
    enable_load_monitor: true
    enable_apm_monitor: true
    rate_limit:
      enabled: true
      limit: 100
      window: 1s
    health_thresholds:
      max_memory_usage: 90.0
      max_cpu_usage: 95.0
      max_goroutines: 10000
    alerts:
      - id: "cpu-high"
        name: "CPU Usage High"
        metric_name: "system.cpu.usage"
        condition: "gt"
        threshold: 80.0
        level: "warning"
        enabled: true
        duration: 5m
        cooldown: 10m
```

### 2. 代码中使用

```go
import (
    "github.com/yourusername/golang/internal/config"
    "github.com/yourusername/golang/pkg/observability"
)

// 加载配置
appConfig, _ := config.LoadConfig()

// 从应用配置创建可观测性配置
obsConfig := observability.ConfigFromAppConfig(appConfig)

// 创建可观测性集成
obs, _ := observability.NewObservability(obsConfig)

// 应用告警规则
observability.ApplyAlertRules(obs, appConfig.Observability.System.Alerts)

// 启动
obs.Start()
defer obs.Stop(ctx)
```

### 3. 环境变量方式

```bash
# OTLP 配置
export APP_OTLP_ENDPOINT=localhost:4317
export APP_OTLP_SERVICE_NAME=my-service

# 系统监控配置
export APP_OBSERVABILITY_SYSTEM_ENABLED=true
export APP_OBSERVABILITY_SYSTEM_COLLECT_INTERVAL=5s
export APP_OBSERVABILITY_SYSTEM_ENABLE_DISK_MONITOR=true
export APP_OBSERVABILITY_SYSTEM_ENABLE_LOAD_MONITOR=true
export APP_OBSERVABILITY_SYSTEM_ENABLE_APM_MONITOR=true
```

## 📚 配置项说明

### OTLP 配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `endpoint` | string | `localhost:4317` | OTLP 端点地址 |
| `insecure` | bool | `true` | 是否使用不安全连接 |
| `service_name` | string | `app` | 服务名称 |
| `service_version` | string | `1.0.0` | 服务版本 |

### 系统监控配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `enabled` | bool | `false` | 是否启用系统监控 |
| `collect_interval` | string | `5s` | 收集间隔 |
| `enable_disk_monitor` | bool | `false` | 是否启用磁盘监控 |
| `enable_load_monitor` | bool | `false` | 是否启用负载监控 |
| `enable_apm_monitor` | bool | `false` | 是否启用 APM 监控 |

### 限流器配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `enabled` | bool | `false` | 是否启用限流器 |
| `limit` | int64 | `100` | 每秒限制 |
| `window` | string | `1s` | 时间窗口 |

### 健康检查阈值

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `max_memory_usage` | float64 | `90.0` | 最大内存使用率（%） |
| `max_cpu_usage` | float64 | `95.0` | 最大 CPU 使用率（%） |
| `max_goroutines` | int | `10000` | 最大 Goroutine 数量 |

### 告警规则配置

| 配置项 | 类型 | 说明 |
|--------|------|------|
| `id` | string | 告警规则 ID |
| `name` | string | 告警规则名称 |
| `metric_name` | string | 指标名称 |
| `condition` | string | 条件（gt, lt, eq, gte, lte） |
| `threshold` | float64 | 阈值 |
| `level` | string | 级别（info, warning, critical） |
| `enabled` | bool | 是否启用 |
| `duration` | string | 持续时间（如 "5m"） |
| `cooldown` | string | 冷却时间（如 "10m"） |

## 🔧 配置优先级

1. **环境变量** - 最高优先级
2. **配置文件** - 中等优先级
3. **默认值** - 最低优先级

## 📖 更多信息

- [快速开始指南](./QUICK-START.md)
- [完整使用指南](./OBSERVABILITY-COMPLETE-GUIDE.md)
- [配置示例](../configs/observability.yaml)

---

**版本**: v1.0.0
