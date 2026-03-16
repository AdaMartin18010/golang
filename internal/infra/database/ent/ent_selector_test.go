package ent

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/golang/internal/infra/database/ent/user"
)

// TestUser_Query_Only 测试 Only 查询
func TestUser_Query_Only(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	created, err := client.User.Create().
		SetID("user-only").
		SetEmail("only@example.com").
		SetName("Only User").
		Save(ctx)
	require.NoError(t, err)

	// 使用 Only 查询
	user, err := client.User.Query().
		Where(user.ID(created.ID)).
		Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, created.ID, user.ID)
}

// TestUser_Query_Only_NotFound 测试 Only 查询不存在
func TestUser_Query_Only_NotFound(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 查询不存在的用户
	_, err := client.User.Query().
		Where(user.ID("nonexistent-id")).
		Only(ctx)
	assert.Error(t, err)
	assert.True(t, IsNotFound(err))
}

// TestUser_Query_OnlyX 测试 OnlyX 查询
func TestUser_Query_OnlyX(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	created, err := client.User.Create().
		SetID("user-onlyx").
		SetEmail("onlyx@example.com").
		SetName("OnlyX User").
		Save(ctx)
	require.NoError(t, err)

	// 使用 OnlyX 查询
	user := client.User.Query().
		Where(user.ID(created.ID)).
		OnlyX(ctx)
	assert.Equal(t, created.ID, user.ID)
}

// TestUser_Query_OnlyX_NotFound 测试 OnlyX 查询不存在时 panic
func TestUser_Query_OnlyX_NotFound(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// OnlyX 应该 panic
	assert.Panics(t, func() {
		client.User.Query().
			Where(user.ID("nonexistent-id")).
			OnlyX(ctx)
	})
}

// TestUser_Query_First 测试 First 查询
func TestUser_Query_First(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	_, err := client.User.Create().
		SetID("user-first").
		SetEmail("first@example.com").
		SetName("First User").
		Save(ctx)
	require.NoError(t, err)

	// 使用 First 查询
	user, err := client.User.Query().
		First(ctx)
	require.NoError(t, err)
	assert.NotNil(t, user)
}

// TestUser_Query_First_NotFound 测试 First 查询不存在
func TestUser_Query_First_NotFound(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 查询不存在的用户
	_, err := client.User.Query().
		First(ctx)
	assert.Error(t, err)
	assert.True(t, IsNotFound(err))
}

// TestUser_Query_FirstX 测试 FirstX 查询
func TestUser_Query_FirstX(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	_, err := client.User.Create().
		SetID("user-firstx").
		SetEmail("firstx@example.com").
		SetName("FirstX User").
		Save(ctx)
	require.NoError(t, err)

	// 使用 FirstX 查询
	user := client.User.Query().
		FirstX(ctx)
	assert.NotNil(t, user)
}

// TestUser_Query_All 测试 All 查询
func TestUser_Query_All(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建多个用户
	for i := 0; i < 3; i++ {
		_, err := client.User.Create().
			SetID(fmt.Sprintf("user-all-%d", i)).
			SetEmail(fmt.Sprintf("all%d@example.com", i)).
			SetName(fmt.Sprintf("All User %d", i)).
			Save(ctx)
		require.NoError(t, err)
	}

	// 查询所有用户
	users, err := client.User.Query().
		All(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 3)
}

// TestUser_Query_Count 测试 Count 查询
func TestUser_Query_Count(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 初始计数
	count, err := client.User.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// 创建用户
	for i := 0; i < 5; i++ {
		_, err := client.User.Create().
			SetID(fmt.Sprintf("user-count-%d", i)).
			SetEmail(fmt.Sprintf("count%d@example.com", i)).
			SetName(fmt.Sprintf("Count User %d", i)).
			Save(ctx)
		require.NoError(t, err)
	}

	// 再次计数
	count, err = client.User.Query().Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 5, count)
}

// TestUser_Query_Exist 测试 Exist 查询
func TestUser_Query_Exist(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 检查空表
	exists, err := client.User.Query().
		Where(user.ID("nonexistent")).
		Exist(ctx)
	require.NoError(t, err)
	assert.False(t, exists)

	// 创建用户
	_, err = client.User.Create().
		SetID("user-exist").
		SetEmail("exist@example.com").
		SetName("Exist User").
		Save(ctx)
	require.NoError(t, err)

	// 检查存在
	exists, err = client.User.Query().
		Where(user.ID("user-exist")).
		Exist(ctx)
	require.NoError(t, err)
	assert.True(t, exists)
}

// TestUser_Query_LimitOffset 测试 Limit 和 Offset
func TestUser_Query_LimitOffset(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建多个用户
	for i := 0; i < 10; i++ {
		_, err := client.User.Create().
			SetID(fmt.Sprintf("user-page-%d", i)).
			SetEmail(fmt.Sprintf("page%d@example.com", i)).
			SetName(fmt.Sprintf("Page User %d", i)).
			Save(ctx)
		require.NoError(t, err)
	}

	// 测试 Limit
	users, err := client.User.Query().
		Limit(5).
		All(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 5)

	// 测试 Limit 和 Offset
	users, err = client.User.Query().
		Limit(3).
		Offset(3).
		All(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 3)
}

// TestUser_Query_Order 测试排序
func TestUser_Query_Order(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	_, err := client.User.Create().
		SetID("user-order-a").
		SetEmail("a@example.com").
		SetName("A User").
		Save(ctx)
	require.NoError(t, err)

	_, err = client.User.Create().
		SetID("user-order-b").
		SetEmail("b@example.com").
		SetName("B User").
		Save(ctx)
	require.NoError(t, err)

	// 测试升序
	users, err := client.User.Query().
		Order(Asc(user.FieldID)).
		All(ctx)
	require.NoError(t, err)
	assert.True(t, len(users) >= 2)

	// 测试降序
	users, err = client.User.Query().
		Order(Desc(user.FieldID)).
		All(ctx)
	require.NoError(t, err)
	assert.True(t, len(users) >= 2)
}

// TestUser_Query_GroupBy 测试 GroupBy
func TestUser_Query_GroupBy(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	_, err := client.User.Create().
		SetID("user-group-1").
		SetEmail("group@example.com").
		SetName("Group User").
		Save(ctx)
	require.NoError(t, err)

	// 测试 GroupBy - 只是验证不会 panic
	_, _ = client.User.Query().
		GroupBy(user.FieldName).
		Strings(ctx)
}

// TestUser_Query_Select 测试 Select
func TestUser_Query_Select(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	_, err := client.User.Create().
		SetID("user-select").
		SetEmail("select@example.com").
		SetName("Select User").
		Save(ctx)
	require.NoError(t, err)

	// 测试 Select - 选择特定字段
	users, err := client.User.Query().
		Select(user.FieldID, user.FieldEmail).
		All(ctx)
	require.NoError(t, err)
	assert.True(t, len(users) > 0)
}

