package core

import (
	"os"
	"testing"
	"time"
)

// TestConfigManagerBasic 测试配置管理器基础功能
func TestConfigManagerBasic(t *testing.T) {
	cm := NewConfigManager()

	// Set and Get
	cm.Set("key1", "value1")
	value, exists := cm.Get("key1")

	if !exists {
		t.Error("Expected key to exist")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}
}

// TestConfigManagerGetString 测试获取字符串配置
func TestConfigManagerGetString(t *testing.T) {
	cm := NewConfigManager()

	cm.Set("str_key", "test_value")
	value, err := cm.GetString("str_key")

	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}
}

// TestConfigManagerGetStringError 测试获取字符串错误情况
func TestConfigManagerGetStringError(t *testing.T) {
	cm := NewConfigManager()

	// 不存在的key
	_, err := cm.GetString("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent key")
	}

	// 错误类型
	cm.Set("int_key", 123)
	_, err = cm.GetString("int_key")
	if err == nil {
		t.Error("Expected error for wrong type")
	}
}

// TestConfigManagerGetInt 测试获取整数配置
func TestConfigManagerGetInt(t *testing.T) {
	cm := NewConfigManager()

	// 测试int类型
	cm.Set("int_key", 42)
	value, err := cm.GetInt("int_key")
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	// 测试int64类型
	cm.Set("int64_key", int64(100))
	value, err = cm.GetInt("int64_key")
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	if value != 100 {
		t.Errorf("Expected 100, got %d", value)
	}

	// 测试float64类型转换
	cm.Set("float_key", 99.0)
	value, err = cm.GetInt("float_key")
	if err != nil {
		t.Fatalf("GetInt failed: %v", err)
	}
	if value != 99 {
		t.Errorf("Expected 99, got %d", value)
	}
}

// TestConfigManagerGetFloat 测试获取浮点数配置
func TestConfigManagerGetFloat(t *testing.T) {
	cm := NewConfigManager()

	// 测试float64类型
	cm.Set("float_key", 3.14)
	value, err := cm.GetFloat("float_key")
	if err != nil {
		t.Fatalf("GetFloat failed: %v", err)
	}
	if value != 3.14 {
		t.Errorf("Expected 3.14, got %f", value)
	}

	// 测试int类型转换
	cm.Set("int_key", 42)
	value, err = cm.GetFloat("int_key")
	if err != nil {
		t.Fatalf("GetFloat failed: %v", err)
	}
	if value != 42.0 {
		t.Errorf("Expected 42.0, got %f", value)
	}
}

// TestConfigManagerGetBool 测试获取布尔配置
func TestConfigManagerGetBool(t *testing.T) {
	cm := NewConfigManager()

	cm.Set("bool_key", true)
	value, err := cm.GetBool("bool_key")

	if err != nil {
		t.Fatalf("GetBool failed: %v", err)
	}
	if value != true {
		t.Error("Expected true")
	}
}

// TestConfigManagerGetDuration 测试获取时间间隔配置
func TestConfigManagerGetDuration(t *testing.T) {
	cm := NewConfigManager()

	// 测试time.Duration类型
	cm.Set("duration_key", 5*time.Second)
	value, err := cm.GetDuration("duration_key")
	if err != nil {
		t.Fatalf("GetDuration failed: %v", err)
	}
	if value != 5*time.Second {
		t.Errorf("Expected 5s, got %v", value)
	}

	// 测试字符串类型
	cm.Set("duration_str", "10m")
	value, err = cm.GetDuration("duration_str")
	if err != nil {
		t.Fatalf("GetDuration failed: %v", err)
	}
	if value != 10*time.Minute {
		t.Errorf("Expected 10m, got %v", value)
	}
}

// TestConfigManagerGetOrDefault 测试获取或默认值
func TestConfigManagerGetOrDefault(t *testing.T) {
	cm := NewConfigManager()

	cm.Set("key1", "value1")

	// 存在的key
	value := cm.GetOrDefault("key1", "default")
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}

	// 不存在的key
	value = cm.GetOrDefault("nonexistent", "default")
	if value != "default" {
		t.Errorf("Expected 'default', got '%v'", value)
	}
}

// TestConfigManagerDelete 测试删除配置
func TestConfigManagerDelete(t *testing.T) {
	cm := NewConfigManager()

	cm.Set("key1", "value1")
	cm.Delete("key1")

	_, exists := cm.Get("key1")
	if exists {
		t.Error("Key should be deleted")
	}
}

// TestConfigManagerHas 测试检查配置存在
func TestConfigManagerHas(t *testing.T) {
	cm := NewConfigManager()

	cm.Set("key1", "value1")

	if !cm.Has("key1") {
		t.Error("Key should exist")
	}
	if cm.Has("nonexistent") {
		t.Error("Key should not exist")
	}
}

// TestConfigManagerKeys 测试获取所有键
func TestConfigManagerKeys(t *testing.T) {
	cm := NewConfigManager()

	cm.Set("key1", "value1")
	cm.Set("key2", "value2")
	cm.Set("key3", "value3")

	keys := cm.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
}

// TestConfigManagerSetMultiple 测试批量设置
func TestConfigManagerSetMultiple(t *testing.T) {
	cm := NewConfigManager()

	configs := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	cm.SetMultiple(configs)

	if !cm.Has("key1") || !cm.Has("key2") || !cm.Has("key3") {
		t.Error("All keys should be set")
	}
}

// TestConfigManagerClear 测试清空配置
func TestConfigManagerClear(t *testing.T) {
	cm := NewConfigManager()

	cm.Set("key1", "value1")
	cm.Set("key2", "value2")

	cm.Clear()

	keys := cm.Keys()
	if len(keys) != 0 {
		t.Error("All keys should be cleared")
	}
}

// TestConfigManagerVersion 测试版本号
func TestConfigManagerVersion(t *testing.T) {
	cm := NewConfigManager()

	v1 := cm.GetVersion()
	cm.Set("key1", "value1")
	v2 := cm.GetVersion()

	if v2 <= v1 {
		t.Error("Version should increase after modification")
	}
}

// TestConfigManagerClone 测试克隆配置
func TestConfigManagerClone(t *testing.T) {
	cm := NewConfigManager()

	cm.Set("key1", "value1")
	cm.Set("key2", 42)

	clone := cm.Clone()

	if len(clone) != 2 {
		t.Errorf("Expected 2 items in clone, got %d", len(clone))
	}
	if clone["key1"] != "value1" {
		t.Error("Clone should have same values")
	}

	// 修改克隆不应该影响原始
	clone["key3"] = "value3"
	if cm.Has("key3") {
		t.Error("Modifying clone should not affect original")
	}
}

// TestConfigManagerLoadFromFile 测试从文件加载
func TestConfigManagerLoadFromFile(t *testing.T) {
	// 创建临时配置文件
	tempFile := "test_config.json"
	content := `{"key1": "value1", "key2": 42, "key3": true}`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)

	cm := NewConfigManager()
	err = cm.LoadFromFile(tempFile)

	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	if !cm.Has("key1") || !cm.Has("key2") || !cm.Has("key3") {
		t.Error("All keys should be loaded")
	}
}

// TestConfigManagerSaveToFile 测试保存到文件
func TestConfigManagerSaveToFile(t *testing.T) {
	cm := NewConfigManager()
	cm.Set("key1", "value1")
	cm.Set("key2", 42)

	tempFile := "test_save.json"
	defer os.Remove(tempFile)

	err := cm.SaveToFile(tempFile)
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("File should be created")
	}

	// 加载并验证
	cm2 := NewConfigManager()
	err = cm2.LoadFromFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to load saved file: %v", err)
	}

	if !cm2.Has("key1") || !cm2.Has("key2") {
		t.Error("Saved keys should be loadable")
	}
}

// TestConfigManagerOnChange 测试变更处理器
func TestConfigManagerOnChange(t *testing.T) {
	cm := NewConfigManager()

	changed := false
	var changedKey string
	var newValue interface{}

	cm.OnChange(func(key string, oldValue, newVal interface{}) {
		changed = true
		changedKey = key
		newValue = newVal
	})

	cm.Set("key1", "value1")

	// 等待异步处理器执行
	time.Sleep(100 * time.Millisecond)

	if !changed {
		t.Error("Change handler should be called")
	}
	if changedKey != "key1" {
		t.Errorf("Expected key 'key1', got '%s'", changedKey)
	}
	if newValue != "value1" {
		t.Errorf("Expected value 'value1', got '%v'", newValue)
	}
}

// TestValidatedConfigManager 测试带验证的配置管理器
func TestValidatedConfigManager(t *testing.T) {
	vcm := NewValidatedConfigManager()

	// 注册验证器
	validator := ConfigValidatorFunc(func(key string, value interface{}) error {
		if value == "invalid" {
			return NewError(ErrorCodeInvalidConfig, "invalid value")
		}
		return nil
	})

	vcm.RegisterValidator("key1", validator)

	// 有效值应该成功
	err := vcm.Set("key1", "valid")
	if err != nil {
		t.Fatalf("Valid value should be accepted: %v", err)
	}

	// 无效值应该失败
	err = vcm.Set("key1", "invalid")
	if err == nil {
		t.Error("Invalid value should be rejected")
	}
}

// TestConfigManagerConcurrent 测试并发访问
func TestConfigManagerConcurrent(t *testing.T) {
	cm := NewConfigManager()

	done := make(chan bool)
	numGoroutines := 10

	// 并发写入
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				cm.Set("key", id)
				cm.Get("key")
				cm.Has("key")
			}
			done <- true
		}(i)
	}

	// 等待完成
	timeout := time.After(5 * time.Second)
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}
}

// BenchmarkConfigManagerSet 基准测试：设置配置
func BenchmarkConfigManagerSet(b *testing.B) {
	cm := NewConfigManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Set("key", "value")
	}
}

// BenchmarkConfigManagerGet 基准测试：获取配置
func BenchmarkConfigManagerGet(b *testing.B) {
	cm := NewConfigManager()
	cm.Set("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Get("key")
	}
}

// BenchmarkConfigManagerGetString 基准测试：获取字符串配置
func BenchmarkConfigManagerGetString(b *testing.B) {
	cm := NewConfigManager()
	cm.Set("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.GetString("key")
	}
}

// TestConfigManagerLoadFromFileError 测试文件加载错误
func TestConfigManagerLoadFromFileError(t *testing.T) {
	cm := NewConfigManager()

	// 不存在的文件
	err := cm.LoadFromFile("nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}

	// 无效的JSON
	tempFile := "invalid.json"
	os.WriteFile(tempFile, []byte("invalid json"), 0644)
	defer os.Remove(tempFile)

	err = cm.LoadFromFile(tempFile)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

// TestConfigManagerMultipleChangeListen 测试多个变更监听器
func TestConfigManagerMultipleChangeListeners(t *testing.T) {
	cm := NewConfigManager()

	count1 := 0
	count2 := 0

	cm.OnChange(func(key string, oldValue, newValue interface{}) {
		count1++
	})

	cm.OnChange(func(key string, oldValue, newValue interface{}) {
		count2++
	})

	cm.Set("key1", "value1")

	// 等待异步处理器执行
	time.Sleep(100 * time.Millisecond)

	if count1 == 0 || count2 == 0 {
		t.Error("All change handlers should be called")
	}
}

// TestValidatedConfigManagerWithoutValidator 测试没有验证器的情况
func TestValidatedConfigManagerWithoutValidator(t *testing.T) {
	vcm := NewValidatedConfigManager()

	// 没有注册验证器的key应该直接设置成功
	err := vcm.Set("key1", "value1")
	if err != nil {
		t.Errorf("Set without validator should succeed: %v", err)
	}

	if !vcm.Has("key1") {
		t.Error("Key should be set")
	}
}
