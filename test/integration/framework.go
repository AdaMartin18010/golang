package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
)

// TestFramework 集成测试框架
type TestFramework struct {
	db        *sql.DB
	dbURL     string
	cleanup   []func() error
	ctx       context.Context
	cancel    context.CancelFunc
}

// TestFrameworkConfig 集成测试框架配置
type TestFrameworkConfig struct {
	DatabaseURL string
	SetupTimeout time.Duration
	TeardownTimeout time.Duration
}

// DefaultTestFrameworkConfig 默认集成测试框架配置
func DefaultTestFrameworkConfig() TestFrameworkConfig {
	return TestFrameworkConfig{
		DatabaseURL:     os.Getenv("TEST_DATABASE_URL"),
		SetupTimeout:    30 * time.Second,
		TeardownTimeout: 30 * time.Second,
	}
}

// NewTestFramework 创建集成测试框架
func NewTestFramework(config TestFrameworkConfig) (*TestFramework, error) {
	ctx, cancel := context.WithCancel(context.Background())

	tf := &TestFramework{
		dbURL:  config.DatabaseURL,
		ctx:    ctx,
		cancel: cancel,
		cleanup: make([]func() error, 0),
	}

	// 初始化数据库连接
	if config.DatabaseURL != "" {
		if err := tf.initDatabase(config.DatabaseURL); err != nil {
			return nil, fmt.Errorf("failed to init database: %w", err)
		}
	}

	return tf, nil
}

// initDatabase 初始化数据库
func (tf *TestFramework) initDatabase(dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	tf.db = db

	// 注册清理函数
	tf.RegisterCleanup(func() error {
		return db.Close()
	})

	return nil
}

// GetDB 获取数据库连接
func (tf *TestFramework) GetDB() *sql.DB {
	return tf.db
}

// RegisterCleanup 注册清理函数
func (tf *TestFramework) RegisterCleanup(fn func() error) {
	tf.cleanup = append(tf.cleanup, fn)
}

// Setup 设置测试环境
func (tf *TestFramework) Setup() error {
	// 运行数据库迁移（如果需要）
	if tf.db != nil {
		if err := tf.runMigrations(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return nil
}

// Teardown 清理测试环境
func (tf *TestFramework) Teardown() error {
	// 执行所有清理函数（逆序）
	for i := len(tf.cleanup) - 1; i >= 0; i-- {
		if err := tf.cleanup[i](); err != nil {
			return fmt.Errorf("cleanup failed: %w", err)
		}
	}

	return nil
}

// runMigrations 运行数据库迁移
func (tf *TestFramework) runMigrations() error {
	// 这里可以集成数据库迁移工具（如 migrate）
	// 当前为占位实现
	return nil
}

// CleanupDatabase 清理数据库
func (tf *TestFramework) CleanupDatabase() error {
	if tf.db == nil {
		return nil
	}

	// 删除所有表（用于测试清理）
	tables := []string{
		"oauth2_tokens",
		"oauth2_clients",
		"oauth2_authorization_codes",
	}

	for _, table := range tables {
		_, err := tf.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// 表可能不存在，忽略错误
			continue
		}
	}

	return nil
}

// ResetDatabase 重置数据库
func (tf *TestFramework) ResetDatabase() error {
	if err := tf.CleanupDatabase(); err != nil {
		return err
	}

	return tf.runMigrations()
}

// WithTransaction 在事务中执行函数
func (tf *TestFramework) WithTransaction(fn func(*sql.Tx) error) error {
	if tf.db == nil {
		return fmt.Errorf("database not initialized")
	}

	tx, err := tf.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// Shutdown 关闭测试框架
func (tf *TestFramework) Shutdown() error {
	tf.cancel()
	return tf.Teardown()
}

// WaitForService 等待服务就绪
func (tf *TestFramework) WaitForService(url string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for service")
		case <-ticker.C:
			// 简单的 HTTP 检查（可以使用 http.Get）
			cmd := exec.CommandContext(ctx, "curl", "-f", url)
			if err := cmd.Run(); err == nil {
				return nil
			}
		}
	}
}

// WaitForDatabase 等待数据库就绪
func (tf *TestFramework) WaitForDatabase(dbURL string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for database")
		case <-ticker.C:
			db, err := sql.Open("postgres", dbURL)
			if err != nil {
				continue
			}

			if err := db.PingContext(ctx); err == nil {
				db.Close()
				return nil
			}

			db.Close()
		}
	}
}
