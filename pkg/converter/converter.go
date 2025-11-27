package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Converter 数据转换器接口
// 提供各种数据格式和类型之间的转换能力
type Converter interface {
	// ToString 转换为字符串
	ToString(v interface{}) string

	// ToInt 转换为整数
	ToInt(v interface{}) (int, error)

	// ToInt64 转换为 int64
	ToInt64(v interface{}) (int64, error)

	// ToFloat64 转换为 float64
	ToFloat64(v interface{}) (float64, error)

	// ToBool 转换为布尔值
	ToBool(v interface{}) (bool, error)

	// ToTime 转换为时间
	ToTime(v interface{}) (time.Time, error)

	// ToJSON 转换为 JSON 字符串
	ToJSON(v interface{}) (string, error)

	// FromJSON 从 JSON 字符串解析
	FromJSON(data string, v interface{}) error

	// ToMap 转换为 map
	ToMap(v interface{}) (map[string]interface{}, error)

	// ToSlice 转换为切片
	ToSlice(v interface{}) ([]interface{}, error)

	// Convert 通用转换方法
	Convert(v interface{}, targetType reflect.Type) (interface{}, error)
}

// DefaultConverter 默认转换器实现
type DefaultConverter struct{}

// NewConverter 创建转换器
func NewConverter() Converter {
	return &DefaultConverter{}
}

func (c *DefaultConverter) ToString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%g", val)
	case bool:
		return strconv.FormatBool(val)
	case time.Time:
		return val.Format(time.RFC3339)
	case time.Duration:
		return val.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (c *DefaultConverter) ToInt(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case uint:
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		return int(val), nil
	case uint64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

func (c *DefaultConverter) ToInt64(v interface{}) (int64, error) {
	switch val := v.(type) {
	case int64:
		return val, nil
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

func (c *DefaultConverter) ToFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case bool:
		if val {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

func (c *DefaultConverter) ToBool(v interface{}) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		return strconv.ParseBool(val)
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(val).Int() != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(val).Uint() != 0, nil
	case float32, float64:
		return reflect.ValueOf(val).Float() != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

func (c *DefaultConverter) ToTime(v interface{}) (time.Time, error) {
	switch val := v.(type) {
	case time.Time:
		return val, nil
	case string:
		// 尝试多种时间格式
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, val); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot parse time: %s", val)
	case int64:
		return time.Unix(val, 0), nil
	default:
		return time.Time{}, fmt.Errorf("cannot convert %T to time.Time", v)
	}
}

func (c *DefaultConverter) ToJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return string(data), nil
}

func (c *DefaultConverter) FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func (c *DefaultConverter) ToMap(v interface{}) (map[string]interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot convert nil to map")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct && rv.Kind() != reflect.Map {
		return nil, fmt.Errorf("cannot convert %T to map", v)
	}

	result := make(map[string]interface{})

	if rv.Kind() == reflect.Map {
		for _, key := range rv.MapKeys() {
			result[c.ToString(key.Interface())] = rv.MapIndex(key).Interface()
		}
		return result, nil
	}

	// 处理结构体
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if jsonTag != "" {
				fieldName = jsonTag
			}
		}

		result[fieldName] = rv.Field(i).Interface()
	}

	return result, nil
}

func (c *DefaultConverter) ToSlice(v interface{}) ([]interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot convert nil to slice")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("cannot convert %T to slice", v)
	}

	result := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}

	return result, nil
}

func (c *DefaultConverter) Convert(v interface{}, targetType reflect.Type) (interface{}, error) {
	if v == nil {
		return reflect.Zero(targetType).Interface(), nil
	}

	sourceType := reflect.TypeOf(v)
	if sourceType == targetType {
		return v, nil
	}

	// 处理指针类型
	if targetType.Kind() == reflect.Ptr {
		elemType := targetType.Elem()
		converted, err := c.Convert(v, elemType)
		if err != nil {
			return nil, err
		}
		ptr := reflect.New(elemType)
		ptr.Elem().Set(reflect.ValueOf(converted))
		return ptr.Interface(), nil
	}

	// 根据目标类型进行转换
	switch targetType.Kind() {
	case reflect.String:
		return c.ToString(v), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := c.ToInt64(v)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(val).Convert(targetType).Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := c.ToInt64(v)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(uint64(val)).Convert(targetType).Interface(), nil
	case reflect.Float32, reflect.Float64:
		val, err := c.ToFloat64(v)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(val).Convert(targetType).Interface(), nil
	case reflect.Bool:
		return c.ToBool(v)
	default:
		return nil, fmt.Errorf("unsupported target type: %v", targetType)
	}
}
