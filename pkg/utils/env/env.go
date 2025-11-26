package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Get 获取环境变量，如果不存在返回默认值
func Get(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetRequired 获取必需的环境变量，如果不存在则panic
func GetRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

// GetInt 获取整数类型的环境变量
func GetInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetIntRequired 获取必需的整数类型的环境变量
func GetIntRequired(key string) int {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not a valid integer: %s", key, value))
	}
	return intValue
}

// GetInt64 获取64位整数类型的环境变量
func GetInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetInt64Required 获取必需的64位整数类型的环境变量
func GetInt64Required(key string) int64 {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not a valid int64: %s", key, value))
	}
	return intValue
}

// GetFloat64 获取浮点数类型的环境变量
func GetFloat64(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return floatValue
}

// GetFloat64Required 获取必需的浮点数类型的环境变量
func GetFloat64Required(key string) float64 {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not a valid float64: %s", key, value))
	}
	return floatValue
}

// GetBool 获取布尔类型的环境变量
func GetBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

// GetBoolRequired 获取必需的布尔类型的环境变量
func GetBoolRequired(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s is not a valid boolean: %s", key, value))
	}
	return boolValue
}

// GetSlice 获取字符串切片类型的环境变量（使用逗号分隔）
func GetSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}

// GetSliceRequired 获取必需的字符串切片类型的环境变量
func GetSliceRequired(key string) []string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return strings.Split(value, ",")
}

// GetSliceWithSeparator 获取字符串切片类型的环境变量（使用指定分隔符）
func GetSliceWithSeparator(key string, separator string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, separator)
}

// Set 设置环境变量
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset 删除环境变量
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Has 检查环境变量是否存在
func Has(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

// GetAll 获取所有环境变量
func GetAll() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			env[pair[0]] = pair[1]
		}
	}
	return env
}

// GetWithPrefix 获取所有以指定前缀开头的环境变量
func GetWithPrefix(prefix string) map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 && strings.HasPrefix(pair[0], prefix) {
			env[pair[0]] = pair[1]
		}
	}
	return env
}

// Expand 展开环境变量（支持 ${VAR} 或 $VAR 格式）
func Expand(s string) string {
	return os.ExpandEnv(s)
}

// ExpandMap 展开map中的环境变量
func ExpandMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = os.ExpandEnv(v)
	}
	return result
}

// LoadFromFile 从文件加载环境变量（.env格式）
func LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// 移除引号
			if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}
			if err := os.Setenv(key, value); err != nil {
				return err
			}
		}
	}
	return nil
}

// MustLoadFromFile 从文件加载环境变量，如果失败则panic
func MustLoadFromFile(filename string) {
	if err := LoadFromFile(filename); err != nil {
		panic(fmt.Sprintf("failed to load environment file %s: %v", filename, err))
	}
}

// ValidateRequired 验证必需的环境变量是否都已设置
func ValidateRequired(keys []string) error {
	var missing []string
	for _, key := range keys {
		if !Has(key) {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}

// GetOrDefault 获取环境变量，如果不存在或为空则返回默认值（别名）
func GetOrDefault(key, defaultValue string) string {
	return Get(key, defaultValue)
}

// IsSet 检查环境变量是否已设置（别名）
func IsSet(key string) bool {
	return Has(key)
}

// Clear 清除所有环境变量（仅清除当前进程的，不影响系统环境变量）
func Clear() {
	env := GetAll()
	for key := range env {
		_ = os.Unsetenv(key)
	}
}

// Copy 复制环境变量到新map
func Copy() map[string]string {
	return GetAll()
}

// Merge 合并环境变量（后面的会覆盖前面的）
func Merge(envs ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, env := range envs {
		for k, v := range env {
			result[k] = v
		}
	}
	return result
}

// Filter 过滤环境变量
func Filter(fn func(key, value string) bool) map[string]string {
	result := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 && fn(pair[0], pair[1]) {
			result[pair[0]] = pair[1]
		}
	}
	return result
}
