package core

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestPluginManager 测试插件管理器基础功能
func TestPluginManager(t *testing.T) {
	pm := NewPluginManager()

	// 注册插件
	plugin := NewLoggingPlugin()
	info := PluginInfo{
		Name:        plugin.Name(),
		Version:     plugin.Version(),
		Type:        plugin.Type(),
		Description: "Test logging plugin",
		Author:      "Test Author",
	}

	err := pm.Register(plugin, info)
	if err != nil {
		t.Fatalf("Failed to register plugin: %v", err)
	}

	// 验证插件已注册
	plugins := pm.List()
	if len(plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(plugins))
	}
}

// TestPluginManagerDuplicateRegistration 测试重复注册
func TestPluginManagerDuplicateRegistration(t *testing.T) {
	pm := NewPluginManager()
	plugin := NewLoggingPlugin()
	info := PluginInfo{Name: plugin.Name(), Version: plugin.Version(), Type: plugin.Type()}

	// 首次注册
	err := pm.Register(plugin, info)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// 重复注册应该失败
	err = pm.Register(plugin, info)
	if err == nil {
		t.Error("Expected error for duplicate registration, got nil")
	}
}

// TestPluginManagerGet 测试获取插件
func TestPluginManagerGet(t *testing.T) {
	pm := NewPluginManager()
	plugin := NewLoggingPlugin()
	info := PluginInfo{Name: plugin.Name(), Version: plugin.Version(), Type: plugin.Type()}

	pm.Register(plugin, info)

	// 获取存在的插件
	retrieved, err := pm.Get("logging")
	if err != nil {
		t.Fatalf("Failed to get plugin: %v", err)
	}
	if retrieved.Name() != "logging" {
		t.Errorf("Expected plugin name 'logging', got '%s'", retrieved.Name())
	}

	// 获取不存在的插件
	_, err = pm.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent plugin, got nil")
	}
}

// TestPluginManagerUnregister 测试注销插件
func TestPluginManagerUnregister(t *testing.T) {
	pm := NewPluginManager()
	plugin := NewLoggingPlugin()
	info := PluginInfo{Name: plugin.Name(), Version: plugin.Version(), Type: plugin.Type()}

	pm.Register(plugin, info)

	// 注销插件
	err := pm.Unregister("logging")
	if err != nil {
		t.Fatalf("Failed to unregister plugin: %v", err)
	}

	// 验证插件已删除
	plugins := pm.List()
	if len(plugins) != 0 {
		t.Errorf("Expected 0 plugins, got %d", len(plugins))
	}
}

// TestPluginManagerListByType 测试按类型列出插件
func TestPluginManagerListByType(t *testing.T) {
	pm := NewPluginManager()

	// 注册不同类型的插件
	logging := NewLoggingPlugin()
	loggingInfo := PluginInfo{Name: logging.Name(), Version: logging.Version(), Type: logging.Type()}
	pm.Register(logging, loggingInfo)

	validation := NewValidationPlugin()
	validationInfo := PluginInfo{Name: validation.Name(), Version: validation.Version(), Type: validation.Type()}
	pm.Register(validation, validationInfo)

	// 按类型获取
	middlewarePlugins := pm.ListByType(PluginTypeMiddleware)
	if len(middlewarePlugins) != 1 {
		t.Errorf("Expected 1 middleware plugin, got %d", len(middlewarePlugins))
	}

	preprocessorPlugins := pm.ListByType(PluginTypePreProcessor)
	if len(preprocessorPlugins) != 1 {
		t.Errorf("Expected 1 preprocessor plugin, got %d", len(preprocessorPlugins))
	}
}

// TestPluginManagerExecute 测试执行插件
func TestPluginManagerExecute(t *testing.T) {
	pm := NewPluginManager()
	plugin := NewValidationPlugin()
	info := PluginInfo{Name: plugin.Name(), Version: plugin.Version(), Type: plugin.Type()}

	pm.Register(plugin, info)

	// 执行插件
	ctx := context.Background()
	data := map[string]interface{}{"test": "data"}

	result, err := pm.Execute(ctx, "validation", data)
	if err != nil {
		t.Fatalf("Plugin execution failed: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result")
	}
}

// TestPluginManagerExecuteChain 测试插件链执行
func TestPluginManagerExecuteChain(t *testing.T) {
	pm := NewPluginManager()

	// 注册多个插件
	logging := NewLoggingPlugin()
	loggingInfo := PluginInfo{Name: logging.Name(), Version: logging.Version(), Type: logging.Type()}
	pm.Register(logging, loggingInfo)

	validation := NewValidationPlugin()
	validationInfo := PluginInfo{Name: validation.Name(), Version: validation.Version(), Type: validation.Type()}
	pm.Register(validation, validationInfo)

	// 执行插件链
	ctx := context.Background()
	data := map[string]interface{}{"test": "data"}

	result, err := pm.ExecuteChain(ctx, []string{"validation", "logging"}, data)
	if err != nil {
		t.Fatalf("Plugin chain execution failed: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result")
	}
}

// TestPluginManagerCleanupAll 测试清理所有插件
func TestPluginManagerCleanupAll(t *testing.T) {
	pm := NewPluginManager()

	// 注册多个插件
	logging := NewLoggingPlugin()
	loggingInfo := PluginInfo{Name: logging.Name(), Version: logging.Version(), Type: logging.Type()}
	pm.Register(logging, loggingInfo)

	validation := NewValidationPlugin()
	validationInfo := PluginInfo{Name: validation.Name(), Version: validation.Version(), Type: validation.Type()}
	pm.Register(validation, validationInfo)

	// 清理所有插件
	err := pm.CleanupAll()
	if err != nil {
		t.Fatalf("CleanupAll failed: %v", err)
	}

	// 验证所有插件已删除
	plugins := pm.List()
	if len(plugins) != 0 {
		t.Errorf("Expected 0 plugins after cleanup, got %d", len(plugins))
	}
}

// TestValidationPluginNilData 测试验证插件处理nil数据
func TestValidationPluginNilData(t *testing.T) {
	plugin := NewValidationPlugin()
	plugin.Initialize(nil)

	ctx := context.Background()
	_, err := plugin.Execute(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil data, got nil")
	}
}

// BenchmarkPluginManagerRegister 基准测试：插件注册
func BenchmarkPluginManagerRegister(b *testing.B) {
	pm := NewPluginManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		plugin := &LoggingPlugin{
			name:    fmt.Sprintf("plugin_%d", i),
			version: "1.0.0",
		}
		info := PluginInfo{
			Name:    plugin.Name(),
			Version: plugin.Version(),
			Type:    PluginTypeMiddleware,
		}
		pm.Register(plugin, info)
	}
}

// BenchmarkPluginManagerExecute 基准测试：插件执行
func BenchmarkPluginManagerExecute(b *testing.B) {
	pm := NewPluginManager()
	plugin := NewLoggingPlugin()
	info := PluginInfo{Name: plugin.Name(), Version: plugin.Version(), Type: plugin.Type()}
	pm.Register(plugin, info)

	ctx := context.Background()
	data := map[string]interface{}{"test": "data"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pm.Execute(ctx, "logging", data)
	}
}

// TestPluginInfo 测试插件信息
func TestPluginInfo(t *testing.T) {
	info := PluginInfo{
		Name:        "test-plugin",
		Version:     "1.0.0",
		Type:        PluginTypeMiddleware,
		Description: "A test plugin",
		Author:      "Test Author",
		Config:      map[string]interface{}{"key": "value"},
	}

	if info.Name != "test-plugin" {
		t.Errorf("Expected name 'test-plugin', got '%s'", info.Name)
	}
	if info.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", info.Version)
	}
	if info.Type != PluginTypeMiddleware {
		t.Errorf("Expected type PluginTypeMiddleware, got %v", info.Type)
	}
}

// TestPluginTypes 测试所有插件类型
func TestPluginTypes(t *testing.T) {
	types := []PluginType{
		PluginTypePreProcessor,
		PluginTypePostProcessor,
		PluginTypeMiddleware,
		PluginTypeExtension,
	}

	for _, pluginType := range types {
		t.Run(string(pluginType), func(t *testing.T) {
			if string(pluginType) == "" {
				t.Error("Plugin type should not be empty")
			}
		})
	}
}

// MockPlugin 测试用模拟插件
type MockPlugin struct {
	name        string
	version     string
	executeFunc func(ctx context.Context, data interface{}) (interface{}, error)
}

func (m *MockPlugin) Name() string                                   { return m.name }
func (m *MockPlugin) Version() string                                { return m.version }
func (m *MockPlugin) Type() PluginType                               { return PluginTypeMiddleware }
func (m *MockPlugin) Initialize(config map[string]interface{}) error { return nil }
func (m *MockPlugin) Execute(ctx context.Context, data interface{}) (interface{}, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, data)
	}
	return data, nil
}
func (m *MockPlugin) Cleanup() error { return nil }

// TestPluginManagerWithMockPlugin 测试使用模拟插件
func TestPluginManagerWithMockPlugin(t *testing.T) {
	pm := NewPluginManager()

	mock := &MockPlugin{
		name:    "mock",
		version: "1.0.0",
		executeFunc: func(ctx context.Context, data interface{}) (interface{}, error) {
			return "processed", nil
		},
	}

	info := PluginInfo{Name: mock.Name(), Version: mock.Version(), Type: mock.Type()}
	pm.Register(mock, info)

	ctx := context.Background()
	result, err := pm.Execute(ctx, "mock", "test")
	if err != nil {
		t.Fatalf("Mock plugin execution failed: %v", err)
	}

	if result != "processed" {
		t.Errorf("Expected 'processed', got '%v'", result)
	}
}

// TestPluginManagerConcurrent 测试并发访问
func TestPluginManagerConcurrent(t *testing.T) {
	pm := NewPluginManager()

	// 注册插件
	plugin := NewLoggingPlugin()
	info := PluginInfo{Name: plugin.Name(), Version: plugin.Version(), Type: plugin.Type()}
	pm.Register(plugin, info)

	ctx := context.Background()
	data := "test"

	// 并发执行
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				pm.Execute(ctx, "logging", data)
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	timeout := time.After(5 * time.Second)
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}
}
