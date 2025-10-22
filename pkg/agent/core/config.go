package core

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// =============================================================================
// 配置管理 - Configuration Management
// =============================================================================

// ConfigManager 配置管理器
type ConfigManager struct {
	configs  map[string]interface{}
	mu       sync.RWMutex
	onChange []ConfigChangeHandler
	version  int
}

// ConfigChangeHandler 配置变更处理器
type ConfigChangeHandler func(key string, oldValue, newValue interface{})

// NewConfigManager 创建配置管理器
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs:  make(map[string]interface{}),
		onChange: make([]ConfigChangeHandler, 0),
		version:  0,
	}
}

// Set 设置配置项
func (cm *ConfigManager) Set(key string, value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	oldValue := cm.configs[key]
	cm.configs[key] = value
	cm.version++

	// 触发变更处理器
	for _, handler := range cm.onChange {
		go handler(key, oldValue, value)
	}
}

// Get 获取配置项
func (cm *ConfigManager) Get(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	value, exists := cm.configs[key]
	return value, exists
}

// GetString 获取字符串配置
func (cm *ConfigManager) GetString(key string) (string, error) {
	value, exists := cm.Get(key)
	if !exists {
		return "", fmt.Errorf("config key %s not found", key)
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("config key %s is not a string", key)
	}

	return str, nil
}

// GetInt 获取整数配置
func (cm *ConfigManager) GetInt(key string) (int, error) {
	value, exists := cm.Get(key)
	if !exists {
		return 0, fmt.Errorf("config key %s not found", key)
	}

	// 尝试不同的数字类型
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("config key %s is not an integer", key)
	}
}

// GetFloat 获取浮点数配置
func (cm *ConfigManager) GetFloat(key string) (float64, error) {
	value, exists := cm.Get(key)
	if !exists {
		return 0, fmt.Errorf("config key %s not found", key)
	}

	// 尝试不同的数字类型
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("config key %s is not a float", key)
	}
}

// GetBool 获取布尔配置
func (cm *ConfigManager) GetBool(key string) (bool, error) {
	value, exists := cm.Get(key)
	if !exists {
		return false, fmt.Errorf("config key %s not found", key)
	}

	b, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("config key %s is not a boolean", key)
	}

	return b, nil
}

// GetDuration 获取时间间隔配置
func (cm *ConfigManager) GetDuration(key string) (time.Duration, error) {
	value, exists := cm.Get(key)
	if !exists {
		return 0, fmt.Errorf("config key %s not found", key)
	}

	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case string:
		return time.ParseDuration(v)
	case int64:
		return time.Duration(v), nil
	default:
		return 0, fmt.Errorf("config key %s is not a duration", key)
	}
}

// GetOrDefault 获取配置项或返回默认值
func (cm *ConfigManager) GetOrDefault(key string, defaultValue interface{}) interface{} {
	value, exists := cm.Get(key)
	if !exists {
		return defaultValue
	}
	return value
}

// Delete 删除配置项
func (cm *ConfigManager) Delete(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.configs, key)
	cm.version++
}

// Has 检查配置项是否存在
func (cm *ConfigManager) Has(key string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	_, exists := cm.configs[key]
	return exists
}

// Keys 获取所有配置键
func (cm *ConfigManager) Keys() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	keys := make([]string, 0, len(cm.configs))
	for key := range cm.configs {
		keys = append(keys, key)
	}

	return keys
}

// SetMultiple 批量设置配置
func (cm *ConfigManager) SetMultiple(configs map[string]interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for key, value := range configs {
		oldValue := cm.configs[key]
		cm.configs[key] = value

		// 触发变更处理器
		for _, handler := range cm.onChange {
			go handler(key, oldValue, value)
		}
	}

	cm.version++
}

// LoadFromFile 从文件加载配置
func (cm *ConfigManager) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var configs map[string]interface{}
	if err := json.Unmarshal(data, &configs); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	cm.SetMultiple(configs)
	return nil
}

// SaveToFile 保存配置到文件
func (cm *ConfigManager) SaveToFile(filename string) error {
	cm.mu.RLock()
	data, err := json.MarshalIndent(cm.configs, "", "  ")
	cm.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal configs: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// OnChange 注册配置变更处理器
func (cm *ConfigManager) OnChange(handler ConfigChangeHandler) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.onChange = append(cm.onChange, handler)
}

// Clear 清空所有配置
func (cm *ConfigManager) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.configs = make(map[string]interface{})
	cm.version++
}

// GetVersion 获取配置版本号
func (cm *ConfigManager) GetVersion() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.version
}

// Clone 克隆配置
func (cm *ConfigManager) Clone() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	clone := make(map[string]interface{}, len(cm.configs))
	for key, value := range cm.configs {
		clone[key] = value
	}

	return clone
}

// =============================================================================
// 配置验证
// =============================================================================

// ConfigValidator 配置验证器
type ConfigValidator interface {
	Validate(key string, value interface{}) error
}

// ConfigValidatorFunc 配置验证器函数类型
type ConfigValidatorFunc func(key string, value interface{}) error

// Validate 实现ConfigValidator接口
func (f ConfigValidatorFunc) Validate(key string, value interface{}) error {
	return f(key, value)
}

// ValidatedConfigManager 带验证的配置管理器
type ValidatedConfigManager struct {
	*ConfigManager
	validators map[string]ConfigValidator
	mu         sync.RWMutex
}

// NewValidatedConfigManager 创建带验证的配置管理器
func NewValidatedConfigManager() *ValidatedConfigManager {
	return &ValidatedConfigManager{
		ConfigManager: NewConfigManager(),
		validators:    make(map[string]ConfigValidator),
	}
}

// RegisterValidator 注册验证器
func (vcm *ValidatedConfigManager) RegisterValidator(key string, validator ConfigValidator) {
	vcm.mu.Lock()
	defer vcm.mu.Unlock()

	vcm.validators[key] = validator
}

// Set 设置配置项（带验证）
func (vcm *ValidatedConfigManager) Set(key string, value interface{}) error {
	vcm.mu.RLock()
	validator, hasValidator := vcm.validators[key]
	vcm.mu.RUnlock()

	if hasValidator {
		if err := validator.Validate(key, value); err != nil {
			return fmt.Errorf("validation failed for key %s: %w", key, err)
		}
	}

	vcm.ConfigManager.Set(key, value)
	return nil
}
