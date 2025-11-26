package transaction

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQLTransaction_Commit(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	sqlTx := NewSQLTransaction(tx)
	err = sqlTx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	// 再次提交应该返回错误
	err = sqlTx.Commit()
	if err != ErrTransactionAlreadyCommitted {
		t.Errorf("Expected ErrTransactionAlreadyCommitted, got %v", err)
	}
}

func TestSQLTransaction_Rollback(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	sqlTx := NewSQLTransaction(tx)
	err = sqlTx.Rollback()
	if err != nil {
		t.Fatalf("Failed to rollback transaction: %v", err)
	}

	// 再次回滚应该返回nil（幂等）
	err = sqlTx.Rollback()
	if err != nil {
		t.Errorf("Expected nil on second rollback, got %v", err)
	}
}

func TestSQLTransactionManager_Begin(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	manager := NewSQLTransactionManager(db)
	tx, err := manager.Begin(context.Background())
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	if tx == nil {
		t.Error("Expected transaction, got nil")
	}
}

func TestSQLTransactionManager_WithTransaction(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 创建表
	_, err = db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	manager := NewSQLTransactionManager(db)

	// 成功的事务
	err = manager.WithTransaction(context.Background(), func(ctx context.Context) error {
		tx, ok := GetSQLTx(ctx)
		if !ok {
			return errors.New("transaction not found in context")
		}

		_, err := tx.Exec("INSERT INTO test (name) VALUES (?)", "test")
		return err
	})

	if err != nil {
		t.Fatalf("Transaction should succeed, got error: %v", err)
	}

	// 验证数据已提交
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM test").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}
}

func TestSQLTransactionManager_WithTransaction_Rollback(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 创建表
	_, err = db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	manager := NewSQLTransactionManager(db)

	// 失败的事务（应该回滚）
	err = manager.WithTransaction(context.Background(), func(ctx context.Context) error {
		tx, ok := GetSQLTx(ctx)
		if !ok {
			return errors.New("transaction not found in context")
		}

		_, err := tx.Exec("INSERT INTO test (name) VALUES (?)", "test")
		if err != nil {
			return err
		}

		// 返回错误，触发回滚
		return errors.New("intentional error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// 验证数据已回滚
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM test").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 rows (rolled back), got %d", count)
	}
}

func TestGetTransaction(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	manager := NewSQLTransactionManager(db)

	ctx := context.Background()
	_, ok := GetTransaction(ctx)
	if ok {
		t.Error("Expected no transaction in context")
	}

	// 在事务中
	err = manager.WithTransaction(ctx, func(ctx context.Context) error {
		tx, ok := GetTransaction(ctx)
		if !ok {
			return errors.New("transaction not found")
		}
		if tx == nil {
			return errors.New("transaction is nil")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
}
