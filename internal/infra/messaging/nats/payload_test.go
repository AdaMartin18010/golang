// Package nats provides unit tests for payload marshaling logic.
package nats

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMarshalPayload_ByteSlice 测试字节数组类型
func TestMarshalPayload_ByteSlice(t *testing.T) {
	input := []byte{0x01, 0x02, 0x03}
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, input, result)
}

// TestMarshalPayload_String 测试字符串类型
func TestMarshalPayload_String(t *testing.T) {
	input := "test message"
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte(input), result)
}

// TestMarshalPayload_EmptyString 测试空字符串
func TestMarshalPayload_EmptyString(t *testing.T) {
	input := ""
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, result)
}

// TestMarshalPayload_Integer 测试整数类型
func TestMarshalPayload_Integer(t *testing.T) {
	input := 42
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("42"), result)
}

// TestMarshalPayload_Float 测试浮点数类型
func TestMarshalPayload_Float(t *testing.T) {
	input := 3.14159
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte("3.14159"), result)
}

// TestMarshalPayload_Boolean 测试布尔类型
func TestMarshalPayload_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		result, err := marshalPayload(true)
		require.NoError(t, err)
		assert.Equal(t, []byte("true"), result)
	})

	t.Run("false", func(t *testing.T) {
		result, err := marshalPayload(false)
		require.NoError(t, err)
		assert.Equal(t, []byte("false"), result)
	})
}

// TestMarshalPayload_Map 测试 Map 类型
func TestMarshalPayload_Map(t *testing.T) {
	input := map[string]interface{}{
		"event":   "user.created",
		"user_id": 123,
	}
	result, err := marshalPayload(input)
	require.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(result, &decoded)
	require.NoError(t, err)
	assert.Equal(t, "user.created", decoded["event"])
}

// TestMarshalPayload_Struct 测试结构体类型
func TestMarshalPayload_Struct(t *testing.T) {
	type Event struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
	}

	input := Event{
		Type:    "order.created",
		Payload: "data",
	}

	result, err := marshalPayload(input)
	require.NoError(t, err)

	var decoded Event
	err = json.Unmarshal(result, &decoded)
	require.NoError(t, err)
	assert.Equal(t, input, decoded)
}

// TestMarshalPayload_Slice 测试切片类型
func TestMarshalPayload_Slice(t *testing.T) {
	input := []string{"a", "b", "c"}
	result, err := marshalPayload(input)
	require.NoError(t, err)
	assert.Equal(t, []byte(`["a","b","c"]`), result)
}

// TestMarshalPayload_Nil 测试 nil 值
func TestMarshalPayload_Nil(t *testing.T) {
	result, err := marshalPayload(nil)
	require.NoError(t, err)
	assert.Equal(t, []byte("null"), result)
}

// TestMarshalPayload_InvalidJSON 测试无法序列化的类型
func TestMarshalPayload_InvalidJSON(t *testing.T) {
	type Circular struct {
		Self *Circular
	}
	circular := &Circular{}
	circular.Self = circular

	_, err := marshalPayload(circular)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal message")
}

// TestMarshalPayload_ComplexData 测试复杂数据结构
func TestMarshalPayload_ComplexData(t *testing.T) {
	type Nested struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type Message struct {
		Action  string   `json:"action"`
		Targets []string `json:"targets"`
		Data    Nested   `json:"data"`
	}

	input := Message{
		Action:  "notify",
		Targets: []string{"user1", "user2"},
		Data: Nested{
			ID:   "123",
			Name: "test",
		},
	}

	result, err := marshalPayload(input)
	require.NoError(t, err)

	var decoded Message
	err = json.Unmarshal(result, &decoded)
	require.NoError(t, err)
	assert.Equal(t, input, decoded)
}

// BenchmarkMarshalPayload_ByteSlice 性能测试
func BenchmarkMarshalPayload_ByteSlice(b *testing.B) {
	input := []byte("benchmark data")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = marshalPayload(input)
	}
}

// BenchmarkMarshalPayload_Struct 性能测试
func BenchmarkMarshalPayload_Struct(b *testing.B) {
	type Data struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	input := Data{ID: "123", Name: "test"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = marshalPayload(input)
	}
}
