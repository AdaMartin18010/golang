package database

import "errors"

var (
	// ErrUnsupportedDriver 不支持的数据库驱动
	ErrUnsupportedDriver = errors.New("unsupported database driver")

	// ErrConnectionFailed 连接失败
	ErrConnectionFailed = errors.New("database connection failed")

	// ErrPingFailed Ping 失败
	ErrPingFailed = errors.New("database ping failed")

	// ErrTransactionNotStarted 事务未开始
	ErrTransactionNotStarted = errors.New("transaction not started")

	// ErrTransactionAlreadyCommitted 事务已提交
	ErrTransactionAlreadyCommitted = errors.New("transaction already committed")

	// ErrTransactionAlreadyRolledBack 事务已回滚
	ErrTransactionAlreadyRolledBack = errors.New("transaction already rolled back")
)
