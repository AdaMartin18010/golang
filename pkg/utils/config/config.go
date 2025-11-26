package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Loader 配置加载器接口
type Loader interface {
	Load(config interface{}) error
}

// FileLoader 文件配置加载器
type FileLoader struct {
	filename string
}

// NewFileLoader 创建文件配置加载器
func NewFileLoader(filename string) *FileLoader {
	return &FileLoader{filename: filename}
}

// Load 从文件加载配置
func (l *FileLoader) Load(config interface{}) error {
	data, err := os.ReadFile(l.filename)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	ext := getFileExt(l.filename)
	switch ext {
	case ".json":
		return json.Unmarshal(data, config)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}
}

// EnvLoader 环境变量配置加载器
type EnvLoader struct {
	prefix string
}

// NewEnvLoader 创建环境变量配置加载器
func NewEnvLoader(prefix string) *EnvLoader {
	return &EnvLoader{prefix: prefix}
}

// Load 从环境变量加载配置
func (l *EnvLoader) Load(config interface{}) error {
	return loadFromEnv(config, l.prefix)
}

// MapLoader Map配置加载器
type MapLoader struct {
	data map[string]interface{}
}

// NewMapLoader 创建Map配置加载器
func NewMapLoader(data map[string]interface{}) *MapLoader {
	return &MapLoader{data: data}
}

// Load 从Map加载配置
func (l *MapLoader) Load(config interface{}) error {
	return loadFromMap(config, l.data)
}

// MultiLoader 多源配置加载器
type MultiLoader struct {
	loaders []Loader
}

// NewMultiLoader 创建多源配置加载器
func NewMultiLoader(loaders ...Loader) *MultiLoader {
	return &MultiLoader{loaders: loaders}
}

// Load 从多个源加载配置（后面的会覆盖前面的）
func (l *MultiLoader) Load(config interface{}) error {
	for _, loader := range l.loaders {
		if err := loader.Load(config); err != nil {
			return err
		}
	}
	return nil
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(config interface{}, prefix string) error {
	val := reflect.ValueOf(config)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("config must be a pointer to struct")
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !field.CanSet() {
			continue
		}

		envKey := getEnvKey(fieldType, prefix)
		envValue := os.Getenv(envKey)

		if envValue == "" {
			continue
		}

		if err := setFieldValue(field, envValue); err != nil {
			return fmt.Errorf("set field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// loadFromMap 从Map加载配置
func loadFromMap(config interface{}, data map[string]interface{}) error {
	val := reflect.ValueOf(config)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("config must be a pointer to struct")
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !field.CanSet() {
			continue
		}

		key := getMapKey(fieldType)
		value, exists := data[key]
		if !exists {
			continue
		}

		if err := setFieldValueFromInterface(field, value); err != nil {
			return fmt.Errorf("set field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// getEnvKey 获取环境变量键名
func getEnvKey(field reflect.StructField, prefix string) string {
	// 检查是否有env标签
	if tag := field.Tag.Get("env"); tag != "" {
		if prefix != "" {
			return prefix + "_" + tag
		}
		return tag
	}

	// 使用字段名
	name := field.Name
	if prefix != "" {
		return prefix + "_" + strings.ToUpper(name)
	}
	return strings.ToUpper(name)
}

// getMapKey 获取Map键名
func getMapKey(field reflect.StructField) string {
	// 检查是否有map标签
	if tag := field.Tag.Get("map"); tag != "" {
		return tag
	}

	// 使用字段名（转换为小写）
	return strings.ToLower(field.Name)
}

// setFieldValue 设置字段值（从字符串）
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

// setFieldValueFromInterface 设置字段值（从interface{}）
func setFieldValueFromInterface(field reflect.Value, value interface{}) error {
	val := reflect.ValueOf(value)

	if val.Type().AssignableTo(field.Type()) {
		field.Set(val)
		return nil
	}

	if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))
		return nil
	}

	// 尝试转换为字符串再设置
	if str, ok := value.(string); ok {
		return setFieldValue(field, str)
	}

	return fmt.Errorf("cannot convert %T to %s", value, field.Type())
}

// getFileExt 获取文件扩展名
func getFileExt(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return ""
	}
	return strings.ToLower(filename[idx:])
}

// Load 从文件加载配置（便捷函数）
func Load(filename string, config interface{}) error {
	loader := NewFileLoader(filename)
	return loader.Load(config)
}

// LoadFromEnv 从环境变量加载配置（便捷函数）
func LoadFromEnv(prefix string, config interface{}) error {
	loader := NewEnvLoader(prefix)
	return loader.Load(config)
}

// LoadFromMap 从Map加载配置（便捷函数）
func LoadFromMap(data map[string]interface{}, config interface{}) error {
	loader := NewMapLoader(data)
	return loader.Load(config)
}
