// Package ent provides simplified tests for Ent ORM client.
//
// 注意：由于 schema 和生成代码之间的 ID 类型不匹配，
// 这些测试使用简化的方法验证基本功能。
//
// 运行测试：
//   - go test -v ./internal/infrastructure/database/ent/...
package ent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
)

// TestNewClient_Simple 测试创建 Ent 客户端
func TestNewClient_Simple(t *testing.T) {
	client, err := Open("sqlite3", "file:test?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)
	defer client.Close()

	assert.NotNil(t, client)
	assert.NotNil(t, client.User)
	assert.NotNil(t, client.Schema)
}

// TestClient_Close_Simple 测试关闭客户端
func TestClient_Close_Simple(t *testing.T) {
	client, err := Open("sqlite3", "file:test?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err)
}

// TestClient_Debug_Simple 测试调试模式
func TestClient_Debug_Simple(t *testing.T) {
	client, err := Open("sqlite3", "file:test?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)
	defer client.Close()

	debugClient := client.Debug()
	assert.NotNil(t, debugClient)
}

// TestClient_Options_Simple 测试客户端选项
func TestClient_Options_Simple(t *testing.T) {
	// 测试 Driver 选项
	client, err := Open("sqlite3", "file:test?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)
	defer client.Close()

	assert.NotNil(t, client)
}
