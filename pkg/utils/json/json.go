package json

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

// Marshal 序列化为JSON
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent 序列化为格式化的JSON
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal 反序列化JSON
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// UnmarshalString 从字符串反序列化JSON
func UnmarshalString(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// MarshalString 序列化为JSON字符串
func MarshalString(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MarshalIndentString 序列化为格式化的JSON字符串
func MarshalIndentString(v interface{}, prefix, indent string) (string, error) {
	data, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// PrettyPrint 美化打印JSON
func PrettyPrint(v interface{}) (string, error) {
	return MarshalIndentString(v, "", "  ")
}

// IsValidJSON 检查字符串是否为有效的JSON
func IsValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// Get 从JSON对象中获取值（使用点号分隔的路径）
func Get(data []byte, path string) (interface{}, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	keys := strings.Split(path, ".")
	current := interface{}(obj)

	for _, key := range keys {
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[key]; exists {
				current = val
			} else {
				return nil, errors.New("key not found: " + key)
			}
		} else {
			return nil, errors.New("invalid path: " + path)
		}
	}

	return current, nil
}

// Set 设置JSON对象中的值（使用点号分隔的路径）
func Set(data []byte, path string, value interface{}) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	keys := strings.Split(path, ".")
	current := interface{}(obj)

	for i, key := range keys {
		if i == len(keys)-1 {
			// 最后一个键，设置值
			if m, ok := current.(map[string]interface{}); ok {
				m[key] = value
			} else {
				return nil, errors.New("invalid path: " + path)
			}
		} else {
			// 中间键，创建或获取嵌套对象
			if m, ok := current.(map[string]interface{}); ok {
				if val, exists := m[key]; exists {
					if _, ok := val.(map[string]interface{}); !ok {
						m[key] = make(map[string]interface{})
					}
					current = m[key]
				} else {
					newMap := make(map[string]interface{})
					m[key] = newMap
					current = newMap
				}
			} else {
				return nil, errors.New("invalid path: " + path)
			}
		}
	}

	return json.Marshal(obj)
}

// Merge 合并多个JSON对象
func Merge(jsons ...[]byte) ([]byte, error) {
	result := make(map[string]interface{})

	for _, data := range jsons {
		var obj map[string]interface{}
		if err := json.Unmarshal(data, &obj); err != nil {
			return nil, err
		}

		for k, v := range obj {
			result[k] = v
		}
	}

	return json.Marshal(result)
}

// ReadFile 从文件读取JSON
func ReadFile(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteFile 将JSON写入文件
func WriteFile(filename string, v interface{}, indent bool) error {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Decode 从Reader解码JSON
func Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// Encode 编码JSON到Writer
func Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// Transform 转换JSON结构
func Transform(data []byte, transformer func(map[string]interface{}) map[string]interface{}) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	transformed := transformer(obj)
	return json.Marshal(transformed)
}

// Filter 过滤JSON对象
func Filter(data []byte, filter func(string, interface{}) bool) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	filtered := make(map[string]interface{})
	for k, v := range obj {
		if filter(k, v) {
			filtered[k] = v
		}
	}

	return json.Marshal(filtered)
}

// Flatten 扁平化嵌套JSON对象
func Flatten(data []byte, separator string) ([]byte, error) {
	if separator == "" {
		separator = "."
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	flattened := make(map[string]interface{})
	flattenMap(obj, "", separator, flattened)

	return json.Marshal(flattened)
}

func flattenMap(m map[string]interface{}, prefix, separator string, result map[string]interface{}) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + separator + k
		}

		if nested, ok := v.(map[string]interface{}); ok {
			flattenMap(nested, key, separator, result)
		} else {
			result[key] = v
		}
	}
}

// Unflatten 反扁平化JSON对象
func Unflatten(data []byte, separator string) ([]byte, error) {
	if separator == "" {
		separator = "."
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	unflattened := make(map[string]interface{})
	for k, v := range obj {
		keys := strings.Split(k, separator)
		setNestedValue(unflattened, keys, v)
	}

	return json.Marshal(unflattened)
}

func setNestedValue(m map[string]interface{}, keys []string, value interface{}) {
	if len(keys) == 1 {
		m[keys[0]] = value
		return
	}

	key := keys[0]
	if _, exists := m[key]; !exists {
		m[key] = make(map[string]interface{})
	}

	if nested, ok := m[key].(map[string]interface{}); ok {
		setNestedValue(nested, keys[1:], value)
	}
}
