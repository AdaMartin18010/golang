# 采样机制

框架级别的采样机制，提供可配置的采样策略，用于控制数据收集和处理的频率。

## 📋 功能特性

- ✅ **多种采样策略**: 总是采样、从不采样、概率采样、速率限制采样、自适应采样
- ✅ **动态调整**: 支持运行时动态调整采样率
- ✅ **线程安全**: 所有采样器都是线程安全的
- ✅ **上下文感知**: 支持基于上下文的采样决策

## 🚀 快速开始

### 概率采样

```go
import "github.com/yourusername/golang/pkg/sampling"

// 创建概率采样器（50% 采样率）
sampler, err := sampling.NewProbabilisticSampler(0.5)
if err != nil {
    log.Fatal(err)
}

// 判断是否采样
if sampler.ShouldSample(ctx) {
    // 执行采样操作
    collectData()
}
```

### 速率限制采样

```go
// 创建速率限制采样器（每秒最多 100 次）
sampler, err := sampling.NewRateLimitingSampler(100.0)
if err != nil {
    log.Fatal(err)
}

// 判断是否采样
if sampler.ShouldSample(ctx) {
    // 执行采样操作
}
```

### 自适应采样

```go
// 创建自适应采样器
sampler, err := sampling.NewAdaptiveSampler(0.5, 0.1, 1.0)
if err != nil {
    log.Fatal(err)
}

// 根据系统负载调整采样率
adaptiveSampler := sampler.(*sampling.AdaptiveSampler)
adaptiveSampler.AdjustForLoad(0.9) // 负载 90%，降低采样率

// 判断是否采样
if sampler.ShouldSample(ctx) {
    // 执行采样操作
}
```

## 📚 API 参考

### Sampler 接口

```go
type Sampler interface {
    ShouldSample(ctx context.Context) bool
    SampleRate() float64
    UpdateRate(rate float64) error
}
```

### 采样器类型

- **AlwaysSampler**: 总是采样（100%）
- **NeverSampler**: 从不采样（0%）
- **ProbabilisticSampler**: 概率采样
- **RateLimitingSampler**: 速率限制采样
- **AdaptiveSampler**: 自适应采样

## 🎯 使用场景

1. **追踪采样**: 控制分布式追踪的采样率
2. **指标收集**: 控制指标收集的频率
3. **日志采样**: 控制日志记录的频率
4. **性能分析**: 控制性能分析数据的收集频率

## 🔗 相关文档

- [OTLP 集成](../observability/README.md)
- 追踪和定位
