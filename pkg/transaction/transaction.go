package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrTransactionNotFound 事务未找到
	ErrTransactionNotFound = errors.New("transaction not found")
	// ErrTransactionAlreadyCommitted 事务已提交
	ErrTransactionAlreadyCommitted = errors.New("transaction already committed")
	// ErrTransactionAlreadyRolledBack 事务已回滚
	ErrTransactionAlreadyRolledBack = errors.New("transaction already rolled back")
)

// Transaction 事务接口
type Transaction interface {
	// Commit 提交事务
	Commit() error
	// Rollback 回滚事务
	Rollback() error
	// GetTx 获取底层事务对象
	GetTx() interface{}
}

// Manager 事务管理器接口
type Manager interface {
	// Begin 开始事务
	Begin(ctx context.Context) (Transaction, error)
	// Get 获取当前事务
	Get(ctx context.Context) (Transaction, error)
	// Commit 提交当前事务
	Commit(ctx context.Context) error
	// Rollback 回滚当前事务
	Rollback(ctx context.Context) error
	// WithTransaction 在事务中执行函数
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

// SQLTransaction SQL事务实现
type SQLTransaction struct {
	tx     *sql.Tx
	committed   bool
	rolledBack  bool
	mu          sync.Mutex
}

// NewSQLTransaction 创建SQL事务
func NewSQLTransaction(tx *sql.Tx) *SQLTransaction {
	return &SQLTransaction{
		tx: tx,
	}
}

// Commit 提交事务
func (t *SQLTransaction) Commit() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.committed {
		return ErrTransactionAlreadyCommitted
	}
	if t.rolledBack {
		return ErrTransactionAlreadyRolledBack
	}

	err := t.tx.Commit()
	if err == nil {
		t.committed = true
	}
	return err
}

// Rollback 回滚事务
func (t *SQLTransaction) Rollback() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.committed {
		return ErrTransactionAlreadyCommitted
	}
	if t.rolledBack {
		return nil // 已经回滚，直接返回
	}

	err := t.tx.Rollback()
	if err == nil {
		t.rolledBack = true
	}
	return err
}

// GetTx 获取底层事务对象
func (t *SQLTransaction) GetTx() interface{} {
	return t.tx
}

// SQLTransactionManager SQL事务管理器
type SQLTransactionManager struct {
	db *sql.DB
}

// NewSQLTransactionManager 创建SQL事务管理器
func NewSQLTransactionManager(db *sql.DB) *SQLTransactionManager {
	return &SQLTransactionManager{
		db: db,
	}
}

// Begin 开始事务
func (m *SQLTransactionManager) Begin(ctx context.Context) (Transaction, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return NewSQLTransaction(tx), nil
}

// Get 获取当前事务（从context中）
func (m *SQLTransactionManager) Get(ctx context.Context) (Transaction, error) {
	tx, ok := ctx.Value(transactionKey{}).(Transaction)
	if !ok {
		return nil, ErrTransactionNotFound
	}
	return tx, nil
}

// Commit 提交当前事务
func (m *SQLTransactionManager) Commit(ctx context.Context) error {
	tx, err := m.Get(ctx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Rollback 回滚当前事务
func (m *SQLTransactionManager) Rollback(ctx context.Context) error {
	tx, err := m.Get(ctx)
	if err != nil {
		return err
	}
	return tx.Rollback()
}

// WithTransaction 在事务中执行函数
func (m *SQLTransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.Begin(ctx)
	if err != nil {
		return err
	}

	// 将事务添加到context
	ctx = context.WithValue(ctx, transactionKey{}, tx)

	// 执行函数
	err = fn(ctx)
	if err != nil {
		// 发生错误，回滚事务
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction rollback error: %w (original error: %v)", rollbackErr, err)
		}
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %w", err)
	}

	return nil
}

// transactionKey context key类型
type transactionKey struct{}

// WithTransaction 在事务中执行函数（辅助函数）
func WithTransaction(ctx context.Context, manager Manager, fn func(context.Context) error) error {
	return manager.WithTransaction(ctx, fn)
}

// GetTransaction 从context获取事务
func GetTransaction(ctx context.Context) (Transaction, bool) {
	tx, ok := ctx.Value(transactionKey{}).(Transaction)
	return tx, ok
}

// GetSQLTx 从context获取SQL事务
func GetSQLTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := GetTransaction(ctx)
	if !ok {
		return nil, false
	}
	sqlTx, ok := tx.GetTx().(*sql.Tx)
	return sqlTx, ok
}
