package security

import (
	"context"
	"testing"
	"time"
)

func TestAuditLogger_Log(t *testing.T) {
	store := NewMemoryAuditLogStore()
	logger := NewAuditLogger(store)

	ctx := context.Background()

	log := &AuditLog{
		UserID:     "user-123",
		Action:     "create",
		Resource:   "user",
		ResourceID: "user-456",
		Result:     AuditResultSuccess,
		Details: map[string]interface{}{
			"name": "Test User",
		},
	}

	if err := logger.Log(ctx, log); err != nil {
		t.Fatalf("Failed to log: %v", err)
	}

	if log.ID == "" {
		t.Error("Log ID should be generated")
	}

	if log.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
}

func TestAuditLogger_LogAction(t *testing.T) {
	store := NewMemoryAuditLogStore()
	logger := NewAuditLogger(store)

	ctx := context.Background()

	err := logger.LogAction(ctx, "user-123", "update", "user", "user-456", AuditResultSuccess, map[string]interface{}{
		"field": "email",
		"value": "new@example.com",
	})

	if err != nil {
		t.Fatalf("Failed to log action: %v", err)
	}
}

func TestAuditLogger_LogAccess(t *testing.T) {
	store := NewMemoryAuditLogStore()
	logger := NewAuditLogger(store)

	ctx := context.Background()

	err := logger.LogAccess(ctx, "user-123", "api", "endpoint-1", AuditResultSuccess, "192.168.1.1", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("Failed to log access: %v", err)
	}
}

func TestAuditLogger_LogSecurity(t *testing.T) {
	store := NewMemoryAuditLogStore()
	logger := NewAuditLogger(store)

	ctx := context.Background()

	err := logger.LogSecurity(ctx, "user-123", "failed_login", map[string]interface{}{
		"attempts": 3,
		"reason":   "invalid_password",
	})

	if err != nil {
		t.Fatalf("Failed to log security event: %v", err)
	}
}

func TestAuditLogger_QueryLogs(t *testing.T) {
	store := NewMemoryAuditLogStore()
	logger := NewAuditLogger(store)

	ctx := context.Background()

	// 记录多个日志
	logger.LogAction(ctx, "user-1", "create", "user", "user-1", AuditResultSuccess, nil)
	logger.LogAction(ctx, "user-1", "update", "user", "user-1", AuditResultSuccess, nil)
	logger.LogAction(ctx, "user-2", "create", "user", "user-2", AuditResultSuccess, nil)

	// 查询特定用户的日志
	filter := &AuditLogFilter{
		UserID: "user-1",
	}

	logs, err := logger.QueryLogs(ctx, filter)
	if err != nil {
		t.Fatalf("Failed to query logs: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}

	// 查询特定操作的日志
	filter = &AuditLogFilter{
		Action: "create",
	}

	logs, err = logger.QueryLogs(ctx, filter)
	if err != nil {
		t.Fatalf("Failed to query logs: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}
}

func TestAuditLogger_QueryLogsWithTimeRange(t *testing.T) {
	store := NewMemoryAuditLogStore()
	logger := NewAuditLogger(store)

	ctx := context.Background()

	// 记录日志
	logger.LogAction(ctx, "user-1", "create", "user", "user-1", AuditResultSuccess, nil)

	time.Sleep(10 * time.Millisecond)

	startTime := time.Now()

	time.Sleep(10 * time.Millisecond)

	logger.LogAction(ctx, "user-2", "create", "user", "user-2", AuditResultSuccess, nil)

	endTime := time.Now()

	// 查询时间范围内的日志
	filter := &AuditLogFilter{
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	logs, err := logger.QueryLogs(ctx, filter)
	if err != nil {
		t.Fatalf("Failed to query logs: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}
}

func TestAuditLogExporter_ExportJSON(t *testing.T) {
	store := NewMemoryAuditLogStore()
	logger := NewAuditLogger(store)

	ctx := context.Background()

	// 记录日志
	logger.LogAction(ctx, "user-1", "create", "user", "user-1", AuditResultSuccess, nil)

	// 导出
	exporter := NewAuditLogExporter(store)
	jsonData, err := exporter.ExportJSON(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to export JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON data should not be empty")
	}
}

func TestAuditLogExporter_ExportCSV(t *testing.T) {
	store := NewMemoryAuditLogStore()
	logger := NewAuditLogger(store)

	ctx := context.Background()

	// 记录日志
	logger.LogAction(ctx, "user-1", "create", "user", "user-1", AuditResultSuccess, nil)

	// 导出
	exporter := NewAuditLogExporter(store)
	csvData, err := exporter.ExportCSV(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to export CSV: %v", err)
	}

	if len(csvData) == 0 {
		t.Error("CSV data should not be empty")
	}
}
