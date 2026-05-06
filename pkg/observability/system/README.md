# 系统监控库

> **版本**: v1.0.0
> **Go版本**: 1.26+

完整的系统资源监控解决方案，提供 CPU、内存、IO、网络、磁盘等全方位的系统监控功能。

## 📋 功能特性

### 基础监控 ✅
- ✅ **CPU 监控**: Linux 精确实现（读取 `/proc/stat`），其他平台简化实现
- ✅ **内存监控**: 内存使用量、总内存、GC 统计、堆内存统计
- ✅ **IO 监控**: 读写字节数、操作数
- ✅ **网络监控**: 网络流量、连接数
- ✅ **磁盘监控**: Unix 完整实现，Windows 占位实现

### 高级功能 ✅
- ✅ **负载监控**: 系统负载平均值、请求速率、并发请求数、队列长度
- ✅ **APM 监控**: 请求持续时间、请求计数、错误计数、活跃连接数、吞吐量
- ✅ **限流器**: 可配置速率限制、时间窗口控制、请求统计
- ✅ **配置热重载**: 定期检查配置变化、动态更新配置
- ✅ **指标聚合**: 计数器、Gauge、直方图聚合、多维度聚合

### 智能分析 ✅
- ✅ **指标导出**: 指标快照导出、JSON 格式、历史记录管理、查询功能
- ✅ **告警系统**: 多级别告警、灵活的告警规则、冷却时间控制
- ✅ **诊断工具**: 系统信息收集、问题自动检测、建议生成
- ✅ **资源预测**: 线性预测算法、趋势分析、置信度计算

### 平台集成 ✅
- ✅ **平台检测**: 操作系统信息、容器检测、虚拟化检测
- ✅ **Kubernetes 集成**: Pod 信息自动检测、Labels 和 Annotations 提取
- ✅ **健康检查**: 可配置健康阈值、定期健康检查
- ✅ **错误处理**: 完善的错误类型、可配置重试机制

## 🚀 快速开始

### 基本使用

```go
import (
    "github.com/yourusername/golang/pkg/observability/system"
    "go.opentelemetry.io/otel/sdk/metric"
)

// 创建系统监控器
systemMonitor, err := system.NewSystemMonitor(system.SystemConfig{
    Meter:           mp.Meter("system"),
    Enabled:         true,
    CollectInterval: 5 * time.Second,
})
if err != nil {
    log.Fatal(err)
}

// 启动监控
if err := systemMonitor.Start(ctx); err != nil {
    log.Fatal(err)
}
defer systemMonitor.Stop(ctx)
```

### 完整功能使用

```go
// 创建系统监控器（启用所有功能）
systemMonitor, err := system.NewSystemMonitor(system.SystemConfig{
    Meter:            mp.Meter("system"),
    Tracer:           tp.Tracer("system"),
    Enabled:          true,
    CollectInterval:  5 * time.Second,
    EnableDiskMonitor: true,
    EnableLoadMonitor: true,
    EnableAPMMonitor:  true,
    RateLimitConfig: &system.RateLimiterConfig{
        Meter:   mp.Meter("ratelimit"),
        Enabled: true,
        Limit:   100,
        Window:  1 * time.Second,
    },
    HealthThresholds: system.DefaultHealthThresholds(),
})
```

## 📊 导出的指标

### 基础系统指标（21 个）
- **系统资源**: 7 个（CPU、内存、GC 等）
- **IO**: 4 个（读写字节数、操作数）
- **网络**: 5 个（流量、连接数等）
- **磁盘**: 5 个（使用量、总量、可用空间等）

### 高级功能指标（13 个）
- **负载监控**: 4 个
- **APM**: 5 个
- **限流器**: 3 个
- **Kubernetes**: 1 个

**总计**: **34 个监控指标**

## 📚 API 参考

### SystemMonitor

```go
type SystemMonitor struct {
    // ...
}

func NewSystemMonitor(cfg SystemConfig) (*SystemMonitor, error)
func (sm *SystemMonitor) Start(ctx context.Context) error
func (sm *SystemMonitor) Stop(ctx context.Context) error

// 获取各种监控器
func (sm *SystemMonitor) GetLoadMonitor() *LoadMonitor
func (sm *SystemMonitor) GetAPMMonitor() *APMMonitor
func (sm *SystemMonitor) GetRateLimiter() *RateLimiter
func (sm *SystemMonitor) GetMetricsExporter() *MetricsExporter
func (sm *SystemMonitor) GetAlertManager() *AlertManager
func (sm *SystemMonitor) GetDiagnostics() *Diagnostics
func (sm *SystemMonitor) GetPredictor() *ResourcePredictor
func (sm *SystemMonitor) GetDashboardExporter() *DashboardExporter
```

### 使用示例

#### APM 监控

```go
apmMonitor := systemMonitor.GetAPMMonitor()
apmMonitor.RecordRequest(ctx, duration, statusCode)
apmMonitor.RecordError(ctx, err)
```

#### 限流器

```go
rateLimiter := systemMonitor.GetRateLimiter()
if rateLimiter.Allow(ctx) {
    // 处理请求
}
```

#### 告警系统

```go
alertManager := systemMonitor.GetAlertManager()
alertManager.Check(ctx, "system.cpu.usage", 85.0, nil)
alerts := alertManager.GetAlertHistory(10)
```

#### 诊断工具

```go
diagnostics := systemMonitor.GetDiagnostics()
report, _ := diagnostics.GenerateReport(ctx)
jsonReport, _ := diagnostics.ExportJSON(ctx)
```

#### 资源预测

```go
predictor := systemMonitor.GetPredictor()
prediction, _ := predictor.Predict(ctx, "system.memory.usage", 1*time.Hour)
```

#### 仪表板导出

```go
dashboardExporter := systemMonitor.GetDashboardExporter()
jsonData, _ := dashboardExporter.ExportJSON(ctx)
promData, _ := dashboardExporter.ExportForPrometheus(ctx)
```

## 🎯 最佳实践

1. **监控间隔**: 根据需求设置合适的收集间隔（建议 5-10 秒）
2. **告警规则**: 根据实际负载设置合理的告警阈值
3. **资源预测**: 定期检查预测结果，提前规划资源
4. **健康检查**: 设置合理的健康阈值，及时发现问题
5. **指标导出**: 定期导出指标快照，用于分析和归档

## 📖 更多文档

- [系统监控实现报告](../../../docs/system-monitoring-implementation.md)
- [高级功能实现报告](../../../docs/ULTIMATE-ADVANCED-FEATURES.md)
- [最终完整实现报告](../../../docs/ULTIMATE-COMPLETE-IMPLEMENTATION-FINAL.md)

## 🔧 配置

### 健康检查阈值

```go
thresholds := system.HealthThresholds{
    MaxCPUUsage:    80.0,
    MaxMemoryUsage: 90.0,
    MaxDiskUsage:   85.0,
}
```

### 告警规则

```go
rule := system.AlertRule{
    ID:         "cpu-high",
    Name:       "CPU Usage High",
    MetricName: "system.cpu.usage",
    Condition:  "gt",
    Threshold:  80.0,
    Level:      system.AlertLevelWarning,
    Enabled:    true,
    Duration:   5 * time.Minute,
    Cooldown:   10 * time.Minute,
}
alertManager.AddRule(rule)
```

## ⚠️ 注意事项

1. **Linux CPU 监控**: 仅在 Linux 平台上提供精确实现，其他平台使用简化实现
2. **Windows 磁盘监控**: 当前为占位实现，需要后续完善
3. **eBPF 功能**: 需要编写 eBPF C 程序才能使用完整功能
4. **依赖**: 需要 OpenTelemetry SDK 依赖

## 🚀 路线图

- [ ] Windows 磁盘监控完整实现
- [ ] eBPF 程序实际实现
- [ ] 更多预测算法（ARIMA、LSTM 等）
- [ ] 告警通知集成（邮件、Slack、PagerDuty 等）

---

**版本**: v1.0.0
**最后更新**: 2025-01-XX
