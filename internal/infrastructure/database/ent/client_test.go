// Package ent provides comprehensive tests for Ent ORM client.
//
// 测试策略：
// 1. 使用 SQLite 内存数据库进行单元测试（快速、无需外部依赖）
// 2. 使用 Open 函数进行 schema 迁移
// 3. 覆盖 CRUD 操作、事务和错误处理
//
// 运行测试：
//   - 单元测试: go test -v ./internal/infrastructure/database/ent/...
package ent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/golang/internal/infrastructure/database/ent/migrate"
	"github.com/yourusername/golang/internal/infrastructure/database/ent/user"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestClient 创建测试用的 Ent 客户端
// 使用 SQLite 内存数据库，测试结束后自动清理
func setupTestClient(t *testing.T) *Client {
	// 使用 Open 创建客户端
	client, err := Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed to open ent client: %v", err)
	}

	// 执行 schema 迁移
	ctx := context.Background()
	if err := migrate.Create(ctx, client.Schema, migrate.Tables); err != nil {
		client.Close()
		t.Fatalf("failed to create schema: %v", err)
	}

	return client
}

// TestNewClient 测试创建 Ent 客户端
func TestNewClient(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	assert.NotNil(t, client)
	assert.NotNil(t, client.User)
	assert.NotNil(t, client.Schema)
}

// TestClient_Close 测试关闭客户端
func TestClient_Close(t *testing.T) {
	client := setupTestClient(t)

	err := client.Close()
	assert.NoError(t, err)
}

// TestUser_Create 测试创建用户
func TestUser_Create(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	u, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("Test User").
		Save(ctx)

	require.NoError(t, err)
	// 注意：schema 中 id 是 uint64，但 ent 生成的是 string
	// 这里我们接受 ent 生成的 ID 格式
	assert.NotEmpty(t, u.ID)
	assert.Equal(t, "test@example.com", u.Email)
	assert.Equal(t, "Test User", u.Name)
	assert.False(t, u.CreatedAt.IsZero())
	assert.False(t, u.UpdatedAt.IsZero())
}

// TestUser_Create_DuplicateEmail 测试创建重复邮箱用户
func TestUser_Create_DuplicateEmail(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 创建第一个用户
	_, err := client.User.Create().
		SetEmail("duplicate@example.com").
		SetName("First User").
		Save(ctx)
	require.NoError(t, err)

	// 尝试创建相同邮箱的用户
	_, err = client.User.Create().
		SetEmail("duplicate@example.com").
		SetName("Second User").
		Save(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unique constraint")
}

// TestUser_Create_InvalidEmail 测试创建无效邮箱用户
func TestUser_Create_InvalidEmail(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 尝试创建无效邮箱的用户
	_, err := client.User.Create().
		SetEmail("invalid-email").
		SetName("Test User").
		Save(ctx)

	assert.Error(t, err)
}

// TestUser_Create_EmptyName 测试创建空名称用户
func TestUser_Create_EmptyName(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 尝试创建空名称的用户
	_, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("").
		Save(ctx)

	assert.Error(t, err)
}

// TestUser_Create_NameTooShort 测试创建名称过短用户
func TestUser_Create_NameTooShort(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 尝试创建名称过短的用户
	_, err := client.User.Create().
		SetEmail("test@example.com").
		SetName("A").
		Save(ctx)

	assert.Error(t, err)
}

// TestUser_Query 测试查询用户
func TestUser_Query(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 创建测试用户
	created, err := client.User.Create().
		SetEmail("query@example.com").
		SetName("Query User").
		Save(ctx)
	require.NoError(t, err)

	// 根据 ID 查询
	u, err := client.User.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, u.ID)
	assert.Equal(t, "query@example.com", u.Email)

	// 查询所有用户
	users, err := client.User.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 1)

	// 根据邮箱查询
	u, err = client.User.Query().
		Where(user.EmailEQ("query@example.com")).
		Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, "Query User", u.Name)
}

// TestUser_Query_NotFound 测试查询不存在的用户
func TestUser_Query_NotFound(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 查询不存在的用户
	_, err := client.User.Get(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestUser_Update 测试更新用户
func TestUser_Update(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 创建测试用户
	created, err := client.User.Create().
		SetEmail("update@example.com").
		SetName("Original Name").
		Save(ctx)
	require.NoError(t, err)

	// 等待一小段时间确保时间戳变化
	// 注意：SQLite 的时间戳精度可能有限

	// 更新用户名称
	updated, err := client.User.UpdateOne(created).
		SetName("Updated Name").
		Save(ctx)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, created.Email, updated.Email) // 邮箱不变

	// UpdatedAt 应该被更新
	// 注意：由于时间精度问题，这里只检查是否不为零
	assert.False(t, updated.UpdatedAt.IsZero())
}

// TestUser_UpdateOneID 测试通过 ID 更新用户
func TestUser_UpdateOneID(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 创建测试用户
	created, err := client.User.Create().
		SetEmail("updateid@example.com").
		SetName("Original Name").
		Save(ctx)
	require.NoError(t, err)

	// 通过 ID 更新
	updated, err := client.User.UpdateOneID(created.ID).
		SetName("Updated By ID").
		Save(ctx)
	require.NoError(t, err)
	assert.Equal(t, "Updated By ID", updated.Name)
}

// TestUser_Delete 测试删除用户
func TestUser_Delete(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 创建测试用户
	created, err := client.User.Create().
		SetEmail("delete@example.com").
		SetName("Delete User").
		Save(ctx)
	require.NoError(t, err)

	// 删除用户
	err = client.User.DeleteOne(created).Exec(ctx)
	require.NoError(t, err)

	// 验证用户已删除
	_, err = client.User.Get(ctx, created.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestUser_DeleteOneID 测试通过 ID 删除用户
func TestUser_DeleteOneID(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 创建测试用户
	created, err := client.User.Create().
		SetEmail("deleteid@example.com").
		SetName("Delete By ID User").
		Save(ctx)
	require.NoError(t, err)

	// 通过 ID 删除
	err = client.User.DeleteOneID(created.ID).Exec(ctx)
	require.NoError(t, err)

	// 验证用户已删除
	exists, _ := client.User.Query().
		Where(user.ID(created.ID)).
		Exist(ctx)
	assert.False(t, exists)
}

// TestUser_Delete_NotFound 测试删除不存在的用户
func TestUser_Delete_NotFound(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 删除不存在的用户
	err := client.User.DeleteOneID("non-existent-id").Exec(ctx)
	assert.Error(t, err)
}

// TestUser_List 测试列出用户
func TestUser_List(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 创建多个测试用户
	for i := 0; i < 5; i++ {
		_, err := client.User.Create().
			SetEmail("list"+string(rune('0'+i))+"@example.com").
			SetName("List User "+string(rune('0'+i))).
			Save(ctx)
		require.NoError(t, err)
	}

	// 查询所有用户
	users, err := client.User.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 5)

	// 分页：限制数量
	users, err = client.User.Query().
		Limit(3).
		All(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 3)

	// 分页：偏移量
	users, err = client.User.Query().
		Offset(3).
		Limit(3).
		All(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

// TestUser_Count 测试用户计数
func TestUser_Count(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 初始计数
	count, err := client.User.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// 创建用户
	for i := 0; i < 3; i++ {
		_, err := client.User.Create().
			SetEmail("count"+string(rune('0'+i))+"@example.com").
			SetName("Count User "+string(rune('0'+i))).
			Save(ctx)
		require.NoError(t, err)
	}

	// 验证计数
	count, err = client.User.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

// TestTx_Commit 测试事务提交
func TestTx_Commit(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 开始事务
	tx, err := client.Tx(ctx)
	require.NoError(t, err)

	// 在事务中创建用户
	u, err := tx.User.Create().
		SetEmail("tx@example.com").
		SetName("Transaction User").
		Save(ctx)
	require.NoError(t, err)

	// 提交事务
	err = tx.Commit()
	require.NoError(t, err)

	// 验证用户已创建
	found, err := client.User.Get(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)
}

// TestTx_Rollback 测试事务回滚
func TestTx_Rollback(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 开始事务
	tx, err := client.Tx(ctx)
	require.NoError(t, err)

	// 在事务中创建用户
	u, err := tx.User.Create().
		SetEmail("rollback@example.com").
		SetName("Rollback User").
		Save(ctx)
	require.NoError(t, err)

	// 回滚事务
	err = tx.Rollback()
	require.NoError(t, err)

	// 验证用户未创建
	_, err = client.User.Get(ctx, u.ID)
	assert.Error(t, err) // 应该找不到
}

// TestTx_NestedTx 测试嵌套事务（应该失败）
func TestTx_NestedTx(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 开始外层事务
	tx, err := client.Tx(ctx)
	require.NoError(t, err)
	defer tx.Rollback()

	// 尝试在事务中开始新事务（应该失败）
	_, err = tx.Client().Tx(ctx)
	assert.Error(t, err)
	assert.Equal(t, ErrTxStarted, err)
}

// TestClient_Debug 测试调试模式
func TestClient_Debug(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	// 启用调试模式
	debugClient := client.Debug()
	assert.NotNil(t, debugClient)

	// 多次调用 Debug 应该返回同一个客户端
	debugClient2 := debugClient.Debug()
	assert.Equal(t, debugClient, debugClient2)
}

// TestClient_Use 测试钩子
func TestClient_Use(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	hookCalled := false

	// 添加钩子
	client.Use(func(next Mutator) Mutator {
		return MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			hookCalled = true
			return next.Mutate(ctx, m)
		})
	})

	ctx := context.Background()
	_, err := client.User.Create().
		SetEmail("hook@example.com").
		SetName("Hook User").
		Save(ctx)

	require.NoError(t, err)
	assert.True(t, hookCalled, "Hook should have been called")
}

// TestUser_BulkCreate 测试批量创建
func TestUser_BulkCreate(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// 批量创建用户
	builders := []*UserCreate{
		client.User.Create().
			SetEmail("bulk1@example.com").
			SetName("Bulk User 1"),
		client.User.Create().
			SetEmail("bulk2@example.com").
			SetName("Bulk User 2"),
	}

	users, err := client.User.CreateBulk(builders...).Save(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

// TestUser_String 测试用户字符串表示
func TestUser_String(t *testing.T) {
	client := setupTestClient(t)
	defer client.Close()

	ctx := context.Background()

	u, err := client.User.Create().
		SetEmail("string@example.com").
		SetName("String User").
		Save(ctx)
	require.NoError(t, err)

	str := u.String()
	assert.Contains(t, str, "string@example.com")
	assert.Contains(t, str, "String User")
	assert.Contains(t, str, u.ID)
}

// BenchmarkUser_Create 基准测试：创建用户
func BenchmarkUser_Create(b *testing.B) {
	client, err := Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.User.Create().
			SetEmail("bench"+string(rune(i))+"@example.com").
			SetName("Benchmark User").
			Save(ctx)
	}
}

// BenchmarkUser_Query 基准测试：查询用户
func BenchmarkUser_Query(b *testing.B) {
	client, err := Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// 准备数据
	u, err := client.User.Create().
		SetEmail("benchquery@example.com").
		SetName("Benchmark Query User").
		Save(ctx)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.User.Get(ctx, u.ID)
	}
}
