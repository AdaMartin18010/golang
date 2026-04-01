// Package mqtt provides unit tests for payload conversion logic.
package mqtt

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConvertPayload_ByteSlice 测试字节数组类型的 payload 转换
func TestConvertPayload_ByteSlice(t *testing.T) {
	input := []byte{0x01, 0x02, 0x03}
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, input, result)
}

// TestConvertPayload_String 测试字符串类型的 payload 转换
func TestConvertPayload_String(t *testing.T) {
	input := "test message"
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte(input), result)
}

// TestConvertPayload_EmptyString 测试空字符串
func TestConvertPayload_EmptyString(t *testing.T) {
	input := ""
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, result)
}

// TestConvertPayload_Integer 测试整数类型
func TestConvertPayload_Integer(t *testing.T) {
	input := 42
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("42"), result)
}

// TestConvertPayload_Float 测试浮点数类型
func TestConvertPayload_Float(t *testing.T) {
	input := 3.14159
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("3.14159"), result)
}

// TestConvertPayload_Boolean 测试布尔类型
func TestConvertPayload_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		result, err := convertPayload(true)
		require.NoError(t, err)
		assert.Equal(t, []byte("true"), result)
	})

	t.Run("false", func(t *testing.T) {
		result, err := convertPayload(false)
		require.NoError(t, err)
		assert.Equal(t, []byte("false"), result)
	})
}

// TestConvertPayload_Map 测试 Map 类型
func TestConvertPayload_Map(t *testing.T) {
	input := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	result, err := convertPayload(input)
	require.NoError(t, err)

	// 验证结果是有效的 JSON
	var decoded map[string]interface{}
	err = json.Unmarshal(result, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "value1", decoded["key1"])
	assert.Equal(t, float64(123), decoded["key2"])
}

// TestConvertPayload_Struct 测试结构体类型
func TestConvertPayload_Struct(t *testing.T) {
	type SensorData struct {
		SensorID    string  `json:"sensor_id"`
		Temperature float64 `json:"temperature"`
		Active      bool    `json:"active"`
	}

	input := SensorData{
		SensorID:    "temp-001",
		Temperature: 25.5,
		Active:      true,
	}

	result, err := convertPayload(input)
	require.NoError(t, err)

	// 验证结果是有效的 JSON
	var decoded SensorData
	err = json.Unmarshal(result, &decoded)
	require.NoError(t, err)
	assert.Equal(t, input, decoded)
}

// TestConvertPayload_Slice 测试切片类型
func TestConvertPayload_Slice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result, err := convertPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("[1,2,3,4,5]"), result)
}

// TestConvertPayload_Nil 测试 nil 值
func TestConvertPayload_Nil(t *testing.T) {
	result, err := convertPayload(nil)
	require.NoError(t, err)
	assert.Equal(t, []byte("null"), result)
}

// TestConvertPayload_InvalidJSON 测试无法序列化的类型
func TestConvertPayload_InvalidJSON(t *testing.T) {
	// 创建一个循环引用，无法序列化
	type Circular struct {
		Self *Circular
	}
	circular := &Circular{}
	circular.Self = circular

	_, err := convertPayload(circular)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal payload")
}

// TestConvertPayload_ComplexStruct 测试复杂嵌套结构体
func TestConvertPayload_ComplexStruct(t *testing.T) {
	type Nested struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	type ComplexData struct {
		ID       string   `json:"id"`
		Tags     []string `json:"tags"`
		Nested   Nested   `json:"nested"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	input := ComplexData{
		ID:   "complex-001",
		Tags: []string{"tag1", "tag2"},
		Nested: Nested{
			Name:  "nested-name",
			Value: 100,
		},
		Metadata: map[string]interface{}{
			"key": "value",
		},
	}

	result, err := convertPayload(input)
	require.NoError(t, err)

	// 验证结果是有效的 JSON
	var decoded ComplexData
	err = json.Unmarshal(result, &decoded)
	require.NoError(t, err)
	assert.Equal(t, input.ID, decoded.ID)
	assert.Equal(t, input.Tags, decoded.Tags)
	assert.Equal(t, input.Nested, decoded.Nested)
}

// BenchmarkConvertPayload_ByteSlice 性能测试：字节数组
func BenchmarkConvertPayload_ByteSlice(b *testing.B) {
	input := []byte("benchmark test data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = convertPayload(input)
	}
}

// BenchmarkConvertPayload_String 性能测试：字符串
func BenchmarkConvertPayload_String(b *testing.B) {
	input := "benchmark test data"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = convertPayload(input)
	}
}

// BenchmarkConvertPayload_Struct 性能测试：结构体
func BenchmarkConvertPayload_Struct(b *testing.B) {
	type Data struct {
		ID    string `json:"id"`
		Value int    `json:"value"`
	}
	input := Data{ID: "bench-001", Value: 42}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = convertPayload(input)
	}
}
