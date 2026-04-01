// Package migrate provides tests for ent schema definitions.
package migrate

import (
	"testing"

	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
	"github.com/stretchr/testify/assert"
)

func TestUsersColumns(t *testing.T) {
	// Test that UsersColumns is defined and has expected columns
	assert.NotNil(t, UsersColumns)
	assert.Equal(t, 5, len(UsersColumns))

	// Test column definitions
	expectedColumns := map[string]struct {
		fieldType field.Type
		unique    bool
		size      int64
	}{
		"id":         {field.TypeString, false, 0},  // v0.14.6: 主键不再显式标记 Unique
		"email":      {field.TypeString, true, 0},   // v0.14.6: 默认不设置 Size
		"name":       {field.TypeString, false, 50}, // 与 schema 定义一致
		"created_at": {field.TypeTime, false, 0},
		"updated_at": {field.TypeTime, false, 0},
	}

	for _, col := range UsersColumns {
		expected, ok := expectedColumns[col.Name]
		assert.True(t, ok, "unexpected column: %s", col.Name)
		assert.Equal(t, expected.fieldType, col.Type, "column %s type mismatch", col.Name)
		assert.Equal(t, expected.unique, col.Unique, "column %s unique mismatch", col.Name)
		if expected.size > 0 {
			assert.Equal(t, expected.size, col.Size, "column %s size mismatch", col.Name)
		}
	}
}

func TestUsersTable(t *testing.T) {
	// Test that UsersTable is properly defined
	assert.NotNil(t, UsersTable)
	assert.Equal(t, "users", UsersTable.Name)
	assert.Equal(t, UsersColumns, UsersTable.Columns)

	// Test primary key
	assert.Equal(t, 1, len(UsersTable.PrimaryKey))
	assert.Equal(t, "id", UsersTable.PrimaryKey[0].Name)

	// Test indexes (v0.14.6: 索引从 schema 生成，如未定义则为空)
	// 注意：当前 schema 未定义显式索引，因此 Indexes 为空
	// 如需索引，请在 schema.User 中添加 Indexes() 方法
}

func TestTables(t *testing.T) {
	// Test that Tables contains all expected tables
	assert.NotNil(t, Tables)
	assert.Equal(t, 1, len(Tables))
	assert.Equal(t, UsersTable, Tables[0])
}

func TestSchemaColumnTypes(t *testing.T) {
	// Test that all column types are correctly defined
	type testCase struct {
		columnName string
		wantType   field.Type
	}

	tests := []testCase{
		{"id", field.TypeString},
		{"email", field.TypeString},
		{"name", field.TypeString},
		{"created_at", field.TypeTime},
		{"updated_at", field.TypeTime},
	}

	columnMap := make(map[string]*schema.Column)
	for _, col := range UsersColumns {
		columnMap[col.Name] = col
	}

	for _, tc := range tests {
		t.Run(tc.columnName, func(t *testing.T) {
			col, ok := columnMap[tc.columnName]
			assert.True(t, ok, "column %s not found", tc.columnName)
			assert.Equal(t, tc.wantType, col.Type, "column %s type mismatch", tc.columnName)
		})
	}
}

func TestTableReferences(t *testing.T) {
	// Verify that Tables slice correctly references UsersTable
	found := false
	for _, table := range Tables {
		if table == UsersTable {
			found = true
			break
		}
	}
	assert.True(t, found, "UsersTable should be in Tables slice")
}

func TestPrimaryKeyConfiguration(t *testing.T) {
	// Verify primary key configuration
	assert.NotNil(t, UsersTable.PrimaryKey)
	assert.Equal(t, 1, len(UsersTable.PrimaryKey))
	assert.Equal(t, UsersColumns[0], UsersTable.PrimaryKey[0])
	assert.Equal(t, "id", UsersTable.PrimaryKey[0].Name)
}
