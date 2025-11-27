# 精细控制机制

框架级别的精细控制机制，提供功能开关、速率控制、熔断器等细粒度控制能力。

## 📋 功能特性

- ✅ **功能开关**: 动态启用/禁用功能
- ✅ **配置管理**: 动态更新配置
- ✅ **配置监听**: 监听配置变化
- ✅ **速率控制**: 细粒度的速率限制
- ✅ **熔断器**: 自动熔断和恢复

## 🚀 快速开始

### 功能控制器

```go
import "github.com/yourusername/golang/pkg/control"

controller := control.NewFeatureController()

// 注册功能
controller.Register("feature-a", "Feature A description", true, map[string]interface{}{
    "max_requests": 100,
})

// 启用/禁用功能
controller.Enable("feature-a")
controller.Disable("feature-a")

// 检查功能状态
if controller.IsEnabled("feature-a") {
    // 执行功能
}

// 更新配置
controller.SetConfig("feature-a", map[string]interface{}{
    "max_requests": 200,
})

// 监听配置变化
controller.Watch("feature-a", func(config interface{}) {
    fmt.Printf("Config updated: %v\n", config)
})
```

### 速率控制器

```go
rateController := control.NewRateController()

// 设置速率限制（每秒最多 100 次）
rateController.SetRateLimit("api-calls", 100.0, time.Second)

// 检查是否允许
if rateController.Allow("api-calls") {
    // 执行操作
}
```

### 熔断器控制器

```go
circuitController := control.NewCircuitController()

// 注册熔断器
circuitController.RegisterCircuit("external-api", 10, 5, 30*time.Second)

// 记录成功/失败
circuitController.RecordSuccess("external-api")
circuitController.RecordFailure("external-api")

// 检查是否允许
if circuitController.Allow("external-api") {
    // 执行操作
}
```

## 📚 API 参考

### Controller 接口

```go
type Controller interface {
    Enable(name string) error
    Disable(name string) error
    IsEnabled(name string) bool
    SetConfig(name string, config interface{}) error
    GetConfig(name string) (interface{}, error)
    Watch(name string, callback func(interface{})) error
    Unwatch(name string) error
}
```

## 🎯 使用场景

1. **功能开关**: 灰度发布、A/B 测试
2. **速率控制**: API 限流、资源保护
3. **熔断器**: 防止级联故障
4. **动态配置**: 运行时配置更新

## 🔗 相关文档

- [采样机制](../sampling/README.md)
- [追踪和定位](../tracing/README.md)
