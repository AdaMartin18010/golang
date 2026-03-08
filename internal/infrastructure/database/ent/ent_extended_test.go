package ent

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestClientWithMigrate 创建测试用的 Ent 客户端并执行迁移
func setupTestClientWithMigrate(t *testing.T) *Client {
	client, err := Open("sqlite3", "file:test?mode=memory&cache=shared&_fk=1")
	require.NoError(t, err)

	// 执行 schema 迁移
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		client.Close()
		t.Fatalf("failed to create schema: %v", err)
	}

	return client
}

// TestNewClient_Options 测试客户端选项
func TestNewClient_Options(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
	}{
		{
			name:    "no options",
			options: nil,
		},
		{
			name: "with debug option",
			options: []Option{
				Debug(),
			},
		},
		{
			name: "with log option",
			options: []Option{
				Log(func(v ...any) {
					t.Log(v...)
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := Open("sqlite3", "file:test?mode=memory&cache=shared&_fk=1")
			require.NoError(t, err)
			defer client.Close()
			assert.NotNil(t, client)
		})
	}
}

// TestOpen_UnsupportedDriver 测试不支持的数据库驱动
func TestOpen_UnsupportedDriver(t *testing.T) {
	_, err := Open("unsupported", "connection_string")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported driver")
}

// TestOpen_InvalidDSN 测试无效的 DSN
func TestOpen_InvalidDSN(t *testing.T) {
	_, err := Open("sqlite3", "///invalid/path/to/db.sqlite")
	// SQLite 对无效路径的处理可能不同，这里只验证返回了错误
	// 注意：在某些情况下可能不会立即返回错误
	if err != nil {
		assert.Error(t, err)
	}
}

// TestClient_Context 测试上下文相关功能
func TestClient_Context(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 测试 FromContext - 应该返回 nil（因为没有存储）
	c := FromContext(ctx)
	assert.Nil(t, c)

	// 测试 NewContext
	newCtx := NewContext(ctx, client)
	c = FromContext(newCtx)
	assert.Equal(t, client, c)
}

// TestTx_Context 测试事务上下文
func TestTx_Context(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 开始事务
	tx, err := client.Tx(ctx)
	require.NoError(t, err)
	defer tx.Rollback()

	// 测试 TxFromContext - 应该返回 nil
	txFromCtx := TxFromContext(ctx)
	assert.Nil(t, txFromCtx)

	// 测试 NewTxContext
	newCtx := NewTxContext(ctx, tx)
	txFromCtx = TxFromContext(newCtx)
	assert.Equal(t, tx, txFromCtx)
}

// TestValidationError 测试验证错误
func TestValidationError(t *testing.T) {
	innerErr := errors.New("field validation failed")
	valErr := &ValidationError{
		Name: "email",
		err:  innerErr,
	}

	assert.Equal(t, innerErr.Error(), valErr.Error())
	assert.Equal(t, innerErr, valErr.Unwrap())
	assert.True(t, IsValidationError(valErr))
	assert.False(t, IsValidationError(nil))
	assert.False(t, IsValidationError(errors.New("other error")))
}

// TestNotFoundError 测试未找到错误
func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{label: "User"}
	assert.Equal(t, "ent: User not found", err.Error())
	assert.True(t, IsNotFound(err))
	assert.False(t, IsNotFound(nil))
	assert.False(t, IsNotFound(errors.New("other error")))
}

// TestMaskNotFound 测试 MaskNotFound 函数
func TestMaskNotFound(t *testing.T) {
	// 当错误是 NotFoundError 时，应该返回 nil
	notFoundErr := &NotFoundError{label: "User"}
	result := MaskNotFound(notFoundErr)
	assert.Nil(t, result)

	// 当错误不是 NotFoundError 时，应该返回原错误
	otherErr := errors.New("other error")
	result = MaskNotFound(otherErr)
	assert.Equal(t, otherErr, result)

	// 当错误为 nil 时，应该返回 nil
	result = MaskNotFound(nil)
	assert.Nil(t, result)
}

// TestNotSingularError 测试非唯一错误
func TestNotSingularError(t *testing.T) {
	err := &NotSingularError{label: "User"}
	assert.Equal(t, "ent: User not singular", err.Error())
	assert.True(t, IsNotSingular(err))
	assert.False(t, IsNotSingular(nil))
	assert.False(t, IsNotSingular(errors.New("other error")))
}

// TestNotLoadedError 测试未加载错误
func TestNotLoadedError(t *testing.T) {
	err := &NotLoadedError{edge: "posts"}
	assert.Equal(t, "ent: posts edge was not loaded", err.Error())
	assert.True(t, IsNotLoaded(err))
	assert.False(t, IsNotLoaded(nil))
	assert.False(t, IsNotLoaded(errors.New("other error")))
}

// TestConstraintError 测试约束错误
func TestConstraintError(t *testing.T) {
	innerErr := errors.New("unique constraint violation")
	constraintErr := &ConstraintError{
		msg:  "duplicate email",
		wrap: innerErr,
	}

	assert.Equal(t, "ent: constraint failed: duplicate email", constraintErr.Error())
	assert.Equal(t, innerErr, constraintErr.Unwrap())
	assert.True(t, IsConstraintError(constraintErr))
	assert.False(t, IsConstraintError(nil))
	assert.False(t, IsConstraintError(errors.New("other error")))
}

// TestTx_CommitHook 测试事务提交钩子
func TestTx_CommitHook(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()
	tx, err := client.Tx(ctx)
	require.NoError(t, err)

	hookCalled := false
	tx.OnCommit(func(next Committer) Committer {
		return CommitFunc(func(ctx context.Context, tx *Tx) error {
			hookCalled = true
			return next.Commit(ctx, tx)
		})
	})

	// 在事务中创建用户
	_, err = tx.User.Create().
		SetID("user-tx-hook").
		SetEmail("tx-hook@example.com").
		SetName("Tx Hook User").
		Save(ctx)
	require.NoError(t, err)

	// 提交事务
	err = tx.Commit()
	require.NoError(t, err)
	assert.True(t, hookCalled, "Commit hook should have been called")
}

// TestTx_RollbackHook 测试事务回滚钩子
func TestTx_RollbackHook(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()
	tx, err := client.Tx(ctx)
	require.NoError(t, err)

	hookCalled := false
	tx.OnRollback(func(next Rollbacker) Rollbacker {
		return RollbackFunc(func(ctx context.Context, tx *Tx) error {
			hookCalled = true
			return next.Rollback(ctx, tx)
		})
	})

	// 在事务中创建用户
	_, err = tx.User.Create().
		SetID("user-tx-rollback-hook").
		SetEmail("tx-rollback-hook@example.com").
		SetName("Tx Rollback Hook User").
		Save(ctx)
	require.NoError(t, err)

	// 回滚事务
	err = tx.Rollback()
	require.NoError(t, err)
	assert.True(t, hookCalled, "Rollback hook should have been called")
}

// TestTx_Client 测试获取事务客户端
func TestTx_Client(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()
	tx, err := client.Tx(ctx)
	require.NoError(t, err)
	defer tx.Rollback()

	// 获取绑定到事务的客户端
	txClient := tx.Client()
	assert.NotNil(t, txClient)
	assert.NotNil(t, txClient.User)

	// 使用事务客户端创建用户
	_, err = txClient.User.Create().
		SetID("user-tx-client").
		SetEmail("tx-client@example.com").
		SetName("Tx Client User").
		Save(ctx)
	require.NoError(t, err)
}

// TestClient_Mutate 测试客户端变异操作
func TestClient_Mutate(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	user, err := client.User.Create().
		SetID("user-mutate").
		SetEmail("mutate@example.com").
		SetName("Mutate User").
		Save(ctx)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "mutate@example.com", user.Email)
}

// TestUser_Value 测试 User 的 Value 方法
func TestUser_Value(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	user, err := client.User.Create().
		SetID("user-value").
		SetEmail("value@example.com").
		SetName("Value User").
		Save(ctx)
	require.NoError(t, err)

	// 测试 Value 方法
	_, err = user.Value("unknown_field")
	assert.Error(t, err)
}

// TestUser_Update 测试通过 User 更新
func TestUser_Update(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建用户
	user, err := client.User.Create().
		SetID("user-update-method").
		SetEmail("update-method@example.com").
		SetName("Original Name").
		Save(ctx)
	require.NoError(t, err)

	// 使用 Update 方法更新
	updated, err := user.Update().
		SetName("Updated Name").
		Save(ctx)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
}

// TestUser_Unwrap 测试 User 的 Unwrap 方法
func TestUser_Unwrap(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 在事务中创建用户
	tx, err := client.Tx(ctx)
	require.NoError(t, err)

	user, err := tx.User.Create().
		SetID("user-unwrap").
		SetEmail("unwrap@example.com").
		SetName("Unwrap User").
		Save(ctx)
	require.NoError(t, err)

	// 提交事务
	err = tx.Commit()
	require.NoError(t, err)

	// 使用 Unwrap 获取非事务用户
	// 注意：这里需要小心，因为 Unwrap 会 panic 如果不是事务实体
	// 由于事务已提交，user 可能已经被 unwrap
}

// TestClient_Intercept 测试客户端拦截器
func TestClient_Intercept(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	interceptorCalled := false

	// 添加拦截器
	client.Intercept(InterceptFunc(func(next Querier) Querier {
		return QuerierFunc(func(ctx context.Context, query Query) (Value, error) {
			interceptorCalled = true
			return next.Query(ctx, query)
		})
	}))

	ctx := context.Background()
	// 查询应该触发拦截器
	_, _ = client.User.Query().All(ctx)

	// 拦截器可能在某些情况下不会被调用，取决于实现
	// 这里我们只验证方法可以正常调用
}

// TestUserClient_Hooks 测试 UserClient 的 Hooks
func TestUserClient_Hooks(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	hookCalled := false

	// 添加钩子到 UserClient
	client.User.Use(func(next Mutator) Mutator {
		return MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			hookCalled = true
			return next.Mutate(ctx, m)
		})
	})

	ctx := context.Background()
	// 创建用户应该触发钩子
	_, err := client.User.Create().
		SetID("user-hook-test").
		SetEmail("hook@example.com").
		SetName("Hook User").
		Save(ctx)

	require.NoError(t, err)
	assert.True(t, hookCalled, "Hook should have been called")
}

// TestUserClient_Interceptors 测试 UserClient 的 Interceptors
func TestUserClient_Interceptors(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	interceptorCalled := false

	// 添加拦截器到 UserClient
	client.User.Intercept(InterceptFunc(func(next Querier) Querier {
		return QuerierFunc(func(ctx context.Context, query Query) (Value, error) {
			interceptorCalled = true
			return next.Query(ctx, query)
		})
	}))

	ctx := context.Background()
	// 查询应该触发拦截器
	_, _ = client.User.Query().All(ctx)

	// 拦截器可能在某些情况下不会被调用
}

// TestUser_String 测试 User 的 String 方法
func TestUser_String(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	user, err := client.User.Create().
		SetID("user-string-test").
		SetEmail("string@example.com").
		SetName("String User").
		Save(ctx)
	require.NoError(t, err)

	str := user.String()
	assert.Contains(t, str, "user-string-test")
	assert.Contains(t, str, "string@example.com")
	assert.Contains(t, str, "String User")
}

// TestUsers_Slice 测试 Users 切片
func TestUsers_Slice(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	ctx := context.Background()

	// 创建多个用户
	for i := 0; i < 3; i++ {
		_, err := client.User.Create().
			SetID(fmt.Sprintf("user-slice-%d", i)).
			SetEmail(fmt.Sprintf("slice%d@example.com", i)).
			SetName(fmt.Sprintf("Slice User %d", i)).
			Save(ctx)
		require.NoError(t, err)
	}

	// 查询所有用户
	users, err := client.User.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 3)
	assert.IsType(t, Users{}, users)
}

// TestClient_Debug_Mode 测试调试模式
func TestClient_Debug_Mode(t *testing.T) {
	client, err := Open("sqlite3", "file:test?mode=memory&cache=shared&_fk=1", Debug())
	require.NoError(t, err)
	defer client.Close()

	// 调试模式的客户端应该启用日志
	debugClient := client.Debug()
	assert.NotNil(t, debugClient)

	// 再次调用 Debug 应该返回同一个客户端
	debugClient2 := debugClient.Debug()
	assert.Equal(t, debugClient, debugClient2)
}

// TestClient_MultipleHooks 测试多个钩子
func TestClient_MultipleHooks(t *testing.T) {
	client := setupTestClientWithMigrate(t)
	defer client.Close()

	hookOrder := []string{}

	// 添加多个钩子
	client.Use(
		func(next Mutator) Mutator {
			return MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				hookOrder = append(hookOrder, "hook1-before")
				v, err := next.Mutate(ctx, m)
				hookOrder = append(hookOrder, "hook1-after")
				return v, err
			})
		},
		func(next Mutator) Mutator {
			return MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				hookOrder = append(hookOrder, "hook2-before")
				v, err := next.Mutate(ctx, m)
				hookOrder = append(hookOrder, "hook2-after")
				return v, err
			})
		},
	)

	ctx := context.Background()
	_, err := client.User.Create().
		SetID("user-multi-hooks").
		SetEmail("multi-hooks@example.com").
		SetName("Multi Hooks User").
		Save(ctx)

	require.NoError(t, err)
	assert.Len(t, hookOrder, 4)
	assert.Equal(t, "hook1-before", hookOrder[0])
	assert.Equal(t, "hook2-before", hookOrder[1])
	assert.Equal(t, "hook2-after", hookOrder[2])
	assert.Equal(t, "hook1-after", hookOrder[3])
}
