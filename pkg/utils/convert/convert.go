package convert

import (
	"fmt"
	"strconv"
	"time"
)

// ToString 转换为字符串
func ToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
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

// ToInt 转换为int
func ToInt(v interface{}) (int, error) {
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

// ToInt64 转换为int64
func ToInt64(v interface{}) (int64, error) {
	switch val := v.(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
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

// ToFloat64 转换为float64
func ToFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
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
	case float32:
		return float64(val), nil
	case float64:
		return val, nil
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

// ToBool 转换为bool
func ToBool(v interface{}) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case int:
		return val != 0, nil
	case int8:
		return val != 0, nil
	case int16:
		return val != 0, nil
	case int32:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case uint:
		return val != 0, nil
	case uint8:
		return val != 0, nil
	case uint16:
		return val != 0, nil
	case uint32:
		return val != 0, nil
	case uint64:
		return val != 0, nil
	case float32:
		return val != 0, nil
	case float64:
		return val != 0, nil
	case string:
		return strconv.ParseBool(val)
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

// ToBytes 转换为[]byte
func ToBytes(v interface{}) []byte {
	switch val := v.(type) {
	case []byte:
		return val
	case string:
		return []byte(val)
	default:
		return []byte(ToString(v))
	}
}

// MustInt 转换为int，失败时panic
func MustInt(v interface{}) int {
	val, err := ToInt(v)
	if err != nil {
		panic(err)
	}
	return val
}

// MustInt64 转换为int64，失败时panic
func MustInt64(v interface{}) int64 {
	val, err := ToInt64(v)
	if err != nil {
		panic(err)
	}
	return val
}

// MustFloat64 转换为float64，失败时panic
func MustFloat64(v interface{}) float64 {
	val, err := ToFloat64(v)
	if err != nil {
		panic(err)
	}
	return val
}

// MustBool 转换为bool，失败时panic
func MustBool(v interface{}) bool {
	val, err := ToBool(v)
	if err != nil {
		panic(err)
	}
	return val
}

// ToIntDefault 转换为int，失败时返回默认值
func ToIntDefault(v interface{}, defaultValue int) int {
	val, err := ToInt(v)
	if err != nil {
		return defaultValue
	}
	return val
}

// ToInt64Default 转换为int64，失败时返回默认值
func ToInt64Default(v interface{}, defaultValue int64) int64 {
	val, err := ToInt64(v)
	if err != nil {
		return defaultValue
	}
	return val
}

// ToFloat64Default 转换为float64，失败时返回默认值
func ToFloat64Default(v interface{}, defaultValue float64) float64 {
	val, err := ToFloat64(v)
	if err != nil {
		return defaultValue
	}
	return val
}

// ToBoolDefault 转换为bool，失败时返回默认值
func ToBoolDefault(v interface{}, defaultValue bool) bool {
	val, err := ToBool(v)
	if err != nil {
		return defaultValue
	}
	return val
}

// ToStringSlice 转换为[]string
func ToStringSlice(v interface{}) []string {
	switch val := v.(type) {
	case []string:
		return val
	case []interface{}:
		result := make([]string, len(val))
		for i, item := range val {
			result[i] = ToString(item)
		}
		return result
	case []int:
		result := make([]string, len(val))
		for i, item := range val {
			result[i] = ToString(item)
		}
		return result
	case []int64:
		result := make([]string, len(val))
		for i, item := range val {
			result[i] = ToString(item)
		}
		return result
	case []float64:
		result := make([]string, len(val))
		for i, item := range val {
			result[i] = ToString(item)
		}
		return result
	case []bool:
		result := make([]string, len(val))
		for i, item := range val {
			result[i] = ToString(item)
		}
		return result
	default:
		return []string{ToString(v)}
	}
}

// ToIntSlice 转换为[]int
func ToIntSlice(v interface{}) ([]int, error) {
	switch val := v.(type) {
	case []int:
		return val, nil
	case []interface{}:
		result := make([]int, len(val))
		for i, item := range val {
			num, err := ToInt(item)
			if err != nil {
				return nil, err
			}
			result[i] = num
		}
		return result, nil
	case []string:
		result := make([]int, len(val))
		for i, item := range val {
			num, err := ToInt(item)
			if err != nil {
				return nil, err
			}
			result[i] = num
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to []int", v)
	}
}

// ToInt64Slice 转换为[]int64
func ToInt64Slice(v interface{}) ([]int64, error) {
	switch val := v.(type) {
	case []int64:
		return val, nil
	case []int:
		result := make([]int64, len(val))
		for i, item := range val {
			result[i] = int64(item)
		}
		return result, nil
	case []interface{}:
		result := make([]int64, len(val))
		for i, item := range val {
			num, err := ToInt64(item)
			if err != nil {
				return nil, err
			}
			result[i] = num
		}
		return result, nil
	case []string:
		result := make([]int64, len(val))
		for i, item := range val {
			num, err := ToInt64(item)
			if err != nil {
				return nil, err
			}
			result[i] = num
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to []int64", v)
	}
}

// ToFloat64Slice 转换为[]float64
func ToFloat64Slice(v interface{}) ([]float64, error) {
	switch val := v.(type) {
	case []float64:
		return val, nil
	case []float32:
		result := make([]float64, len(val))
		for i, item := range val {
			result[i] = float64(item)
		}
		return result, nil
	case []int:
		result := make([]float64, len(val))
		for i, item := range val {
			result[i] = float64(item)
		}
		return result, nil
	case []interface{}:
		result := make([]float64, len(val))
		for i, item := range val {
			num, err := ToFloat64(item)
			if err != nil {
				return nil, err
			}
			result[i] = num
		}
		return result, nil
	case []string:
		result := make([]float64, len(val))
		for i, item := range val {
			num, err := ToFloat64(item)
			if err != nil {
				return nil, err
			}
			result[i] = num
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to []float64", v)
	}
}

// ToBoolSlice 转换为[]bool
func ToBoolSlice(v interface{}) ([]bool, error) {
	switch val := v.(type) {
	case []bool:
		return val, nil
	case []interface{}:
		result := make([]bool, len(val))
		for i, item := range val {
			b, err := ToBool(item)
			if err != nil {
				return nil, err
			}
			result[i] = b
		}
		return result, nil
	case []string:
		result := make([]bool, len(val))
		for i, item := range val {
			b, err := ToBool(item)
			if err != nil {
				return nil, err
			}
			result[i] = b
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to []bool", v)
	}
}

// ToMapStringInterface 转换为map[string]interface{}
func ToMapStringInterface(v interface{}) (map[string]interface{}, error) {
	switch val := v.(type) {
	case map[string]interface{}:
		return val, nil
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			result[ToString(k)] = v
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to map[string]interface{}", v)
	}
}

// ToMapStringString 转换为map[string]string
func ToMapStringString(v interface{}) (map[string]string, error) {
	switch val := v.(type) {
	case map[string]string:
		return val, nil
	case map[string]interface{}:
		result := make(map[string]string)
		for k, v := range val {
			result[k] = ToString(v)
		}
		return result, nil
	case map[interface{}]interface{}:
		result := make(map[string]string)
		for k, v := range val {
			result[ToString(k)] = ToString(v)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to map[string]string", v)
	}
}

// IsNumeric 检查是否为数字类型
func IsNumeric(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case string:
		_, err := strconv.ParseFloat(v.(string), 64)
		return err == nil
	default:
		return false
	}
}

// IsInteger 检查是否为整数类型
func IsInteger(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case string:
		_, err := strconv.ParseInt(v.(string), 10, 64)
		return err == nil
	default:
		return false
	}
}

// IsFloat 检查是否为浮点数类型
func IsFloat(v interface{}) bool {
	switch v.(type) {
	case float32, float64:
		return true
	case string:
		_, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return false
		}
		// 检查是否包含小数点
		return contains(v.(string), ".")
	default:
		return false
	}
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || indexOf(s, substr) >= 0)
}

// indexOf 查找子字符串的位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

