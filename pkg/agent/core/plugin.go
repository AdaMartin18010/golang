package core

import (
	"context"
	"fmt"
	"sync"
)

// =============================================================================
// 插件系统 - Plugin System
// =============================================================================

// PluginType 插件类型
type PluginType string

const (
	PluginTypePreProcessor  PluginType = "preprocessor"  // 预处理器
	PluginTypePostProcessor PluginType = "postprocessor" // 后处理器
	PluginTypeMiddleware    PluginType = "middleware"    // 中间件
	PluginTypeExtension     PluginType = "extension"     // 扩展功能
)

// Plugin 插件接口
type Plugin interface {
	// Name 返回插件名称
	Name() string

	// Version 返回插件版本
	Version() string

	// Type 返回插件类型
	Type() PluginType

	// Initialize 初始化插件
	Initialize(config map[string]interface{}) error

	// Execute 执行插件逻辑
	Execute(ctx context.Context, data interface{}) (interface{}, error)

	// Cleanup 清理插件资源
	Cleanup() error
}

// PluginInfo 插件信息
type PluginInfo struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Type        PluginType             `json:"type"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	Config      map[string]interface{} `json:"config"`
}

// PluginManager 插件管理器
type PluginManager struct {
	plugins map[string]Plugin
	info    map[string]PluginInfo
	mu      sync.RWMutex
}

// NewPluginManager 创建插件管理器
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
		info:    make(map[string]PluginInfo),
	}
}

// Register 注册插件
func (pm *PluginManager) Register(plugin Plugin, info PluginInfo) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	name := plugin.Name()
	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	// 初始化插件
	if err := plugin.Initialize(info.Config); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", name, err)
	}

	pm.plugins[name] = plugin
	pm.info[name] = info

	return nil
}

// Unregister 注销插件
func (pm *PluginManager) Unregister(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// 清理插件资源
	if err := plugin.Cleanup(); err != nil {
		return fmt.Errorf("failed to cleanup plugin %s: %w", name, err)
	}

	delete(pm.plugins, name)
	delete(pm.info, name)

	return nil
}

// Get 获取插件
func (pm *PluginManager) Get(name string) (Plugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// List 列出所有插件
func (pm *PluginManager) List() []PluginInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make([]PluginInfo, 0, len(pm.info))
	for _, info := range pm.info {
		result = append(result, info)
	}

	return result
}

// ListByType 按类型列出插件
func (pm *PluginManager) ListByType(pluginType PluginType) []PluginInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make([]PluginInfo, 0)
	for _, info := range pm.info {
		if info.Type == pluginType {
			result = append(result, info)
		}
	}

	return result
}

// Execute 执行插件
func (pm *PluginManager) Execute(ctx context.Context, name string, data interface{}) (interface{}, error) {
	plugin, err := pm.Get(name)
	if err != nil {
		return nil, err
	}

	return plugin.Execute(ctx, data)
}

// ExecuteChain 按顺序执行多个插件
func (pm *PluginManager) ExecuteChain(ctx context.Context, names []string, data interface{}) (interface{}, error) {
	result := data

	for _, name := range names {
		plugin, err := pm.Get(name)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugin %s: %w", name, err)
		}

		result, err = plugin.Execute(ctx, result)
		if err != nil {
			return nil, fmt.Errorf("plugin %s execution failed: %w", name, err)
		}
	}

	return result, nil
}

// CleanupAll 清理所有插件
func (pm *PluginManager) CleanupAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var lastErr error
	for name, plugin := range pm.plugins {
		if err := plugin.Cleanup(); err != nil {
			lastErr = fmt.Errorf("failed to cleanup plugin %s: %w", name, err)
		}
	}

	pm.plugins = make(map[string]Plugin)
	pm.info = make(map[string]PluginInfo)

	return lastErr
}

// =============================================================================
// 内置插件示例
// =============================================================================

// LoggingPlugin 日志插件示例
type LoggingPlugin struct {
	name    string
	version string
	config  map[string]interface{}
}

// NewLoggingPlugin 创建日志插件
func NewLoggingPlugin() *LoggingPlugin {
	return &LoggingPlugin{
		name:    "logging",
		version: "1.0.0",
	}
}

func (p *LoggingPlugin) Name() string {
	return p.name
}

func (p *LoggingPlugin) Version() string {
	return p.version
}

func (p *LoggingPlugin) Type() PluginType {
	return PluginTypeMiddleware
}

func (p *LoggingPlugin) Initialize(config map[string]interface{}) error {
	p.config = config
	return nil
}

func (p *LoggingPlugin) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	// 这里可以添加日志记录逻辑
	// log.Printf("[LoggingPlugin] Processing data: %+v", data)
	return data, nil
}

func (p *LoggingPlugin) Cleanup() error {
	return nil
}

// ValidationPlugin 验证插件示例
type ValidationPlugin struct {
	name    string
	version string
	config  map[string]interface{}
}

// NewValidationPlugin 创建验证插件
func NewValidationPlugin() *ValidationPlugin {
	return &ValidationPlugin{
		name:    "validation",
		version: "1.0.0",
	}
}

func (p *ValidationPlugin) Name() string {
	return p.name
}

func (p *ValidationPlugin) Version() string {
	return p.version
}

func (p *ValidationPlugin) Type() PluginType {
	return PluginTypePreProcessor
}

func (p *ValidationPlugin) Initialize(config map[string]interface{}) error {
	p.config = config
	return nil
}

func (p *ValidationPlugin) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	// 这里可以添加数据验证逻辑
	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}
	return data, nil
}

func (p *ValidationPlugin) Cleanup() error {
	return nil
}
