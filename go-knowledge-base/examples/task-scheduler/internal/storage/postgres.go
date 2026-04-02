package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// PostgresClient PostgreSQL 客户端
type PostgresClient struct {
	db *sql.DB
}

// NewPostgresClient 创建 PostgreSQL 客户端
func NewPostgresClient(dsn string) (*PostgresClient, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &PostgresClient{db: db}, nil
}

// Close 关闭连接
func (p *PostgresClient) Close() error {
	return p.db.Close()
}

// InitSchema 初始化数据库表
func (p *PostgresClient) InitSchema(ctx context.Context) error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id VARCHAR(64) PRIMARY KEY,
		type VARCHAR(50) NOT NULL,
		payload JSONB,
		status VARCHAR(20) NOT NULL DEFAULT 'pending',
		priority INT DEFAULT 0,
		scheduled_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		worker_id VARCHAR(64),
		retry_count INT DEFAULT 0,
		max_retries INT DEFAULT 3,
		error_message TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_scheduled ON tasks(scheduled_at) WHERE status = 'pending';
	CREATE INDEX IF NOT EXISTS idx_tasks_worker ON tasks(worker_id) WHERE status = 'running';
	`

	_, err := p.db.ExecContext(ctx, schema)
	return err
}

// SaveTask 保存任务
func (p *PostgresClient) SaveTask(ctx context.Context, task Task) error {
	payload, err := json.Marshal(task.Payload)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO tasks (id, type, payload, priority, scheduled_at, status)
		VALUES ($1, $2, $3, $4, $5, 'pending')
		ON CONFLICT (id) DO UPDATE SET
			payload = EXCLUDED.payload,
			priority = EXCLUDED.priority,
			scheduled_at = EXCLUDED.scheduled_at,
			updated_at = NOW()
	`

	_, err = p.db.ExecContext(ctx, query, task.ID, task.Type, payload, task.Priority, task.Scheduled)
	return err
}

// GetTasks 获取待执行任务
func (p *PostgresClient) GetTasks(ctx context.Context, status string, limit int) ([]Task, error) {
	query := `
		SELECT id, type, payload, priority, scheduled_at
		FROM tasks
		WHERE status = $1 AND scheduled_at <= NOW()
		ORDER BY priority DESC, scheduled_at ASC
		LIMIT $2
	`

	rows, err := p.db.QueryContext(ctx, query, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var payload []byte
		err := rows.Scan(&task.ID, &task.Type, &payload, &task.Priority, &task.Scheduled)
		if err != nil {
			continue
		}
		json.Unmarshal(payload, &task.Payload)
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

// UpdateTaskStatus 更新任务状态
func (p *PostgresClient) UpdateTaskStatus(ctx context.Context, taskID string, status string) error {
	query := `
		UPDATE tasks 
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`
	_, err := p.db.ExecContext(ctx, query, taskID, status)
	return err
}

// AcquireLeadership 实现 Storage 接口
func (p *PostgresClient) AcquireLeadership(ctx context.Context, nodeID string, ttl time.Duration) (bool, error) {
	// PostgreSQL 通过 advisory lock 实现
	// 或使用独立的锁表
	return true, nil
}

func (p *PostgresClient) ReleaseLeadership(ctx context.Context, nodeID string) error {
	return nil
}

func (p *PostgresClient) RegisterWorker(ctx context.Context, worker *Worker) error {
	return nil
}

func (p *PostgresClient) GetWorkers(ctx context.Context) ([]Worker, error) {
	return nil, nil
}
