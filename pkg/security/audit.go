package security

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrAuditLogNotFound 审计日志未找到
	ErrAuditLogNotFound = errors.New("audit log not found")
)

// AuditLogger 审计日志记录器
type AuditLogger struct {
	store AuditLogStore
	mu    sync.RWMutex
}

// AuditLog 审计日志
type AuditLog struct {
	ID          string
	Timestamp   time.Time
	UserID      string
	Action      string
	Resource    string
	ResourceID  string
	Result      AuditResult
	IPAddress   string
	UserAgent   string
	RequestID   string
	Details     map[string]interface{}
	Metadata    map[string]string
}

// AuditResult 审计结果
type AuditResult string

const (
	// AuditResultSuccess 成功
	AuditResultSuccess AuditResult = "success"
	// AuditResultFailure 失败
	AuditResultFailure AuditResult = "failure"
	// AuditResultDenied 拒绝
	AuditResultDenied AuditResult = "denied"
)

// AuditLogStore 审计日志存储接口
type AuditLogStore interface {
	Save(ctx context.Context, log *AuditLog) error
	Get(ctx context.Context, logID string) (*AuditLog, error)
	Query(ctx context.Context, filter *AuditLogFilter) ([]*AuditLog, error)
	Delete(ctx context.Context, logID string) error
}

// AuditLogFilter 审计日志过滤器
type AuditLogFilter struct {
	UserID     string
	Action     string
	Resource   string
	Result     AuditResult
	StartTime  *time.Time
	EndTime    *time.Time
	Limit      int
	Offset     int
}

// NewAuditLogger 创建审计日志记录器
func NewAuditLogger(store AuditLogStore) *AuditLogger {
	return &AuditLogger{
		store: store,
	}
}

// Log 记录审计日志
func (l *AuditLogger) Log(ctx context.Context, log *AuditLog) error {
	if log.ID == "" {
		log.ID = generateAuditLogID()
	}

	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	if log.Details == nil {
		log.Details = make(map[string]interface{})
	}

	if log.Metadata == nil {
		log.Metadata = make(map[string]string)
	}

	return l.store.Save(ctx, log)
}

// LogAction 记录操作审计日志
func (l *AuditLogger) LogAction(ctx context.Context, userID, action, resource, resourceID string, result AuditResult, details map[string]interface{}) error {
	log := &AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Result:     result,
		Details:    details,
	}

	return l.Log(ctx, log)
}

// LogAccess 记录访问审计日志
func (l *AuditLogger) LogAccess(ctx context.Context, userID, resource, resourceID string, result AuditResult, ipAddress, userAgent string) error {
	log := &AuditLog{
		UserID:     userID,
		Action:     "access",
		Resource:   resource,
		ResourceID: resourceID,
		Result:     result,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	return l.Log(ctx, log)
}

// LogSecurity 记录安全事件审计日志
func (l *AuditLogger) LogSecurity(ctx context.Context, userID, event string, details map[string]interface{}) error {
	log := &AuditLog{
		UserID:   userID,
		Action:   "security",
		Resource: "security",
		Result:   AuditResultFailure,
		Details:  details,
	}

	if details == nil {
		log.Details = make(map[string]interface{})
	}
	log.Details["event"] = event

	return l.Log(ctx, log)
}

// GetLog 获取审计日志
func (l *AuditLogger) GetLog(ctx context.Context, logID string) (*AuditLog, error) {
	return l.store.Get(ctx, logID)
}

// QueryLogs 查询审计日志
func (l *AuditLogger) QueryLogs(ctx context.Context, filter *AuditLogFilter) ([]*AuditLog, error) {
	return l.store.Query(ctx, filter)
}

// DeleteLog 删除审计日志
func (l *AuditLogger) DeleteLog(ctx context.Context, logID string) error {
	return l.store.Delete(ctx, logID)
}

// generateAuditLogID 生成审计日志 ID
func generateAuditLogID() string {
	return fmt.Sprintf("audit_%d_%d", time.Now().UnixNano(), time.Now().Unix())
}

// MemoryAuditLogStore 内存审计日志存储（用于测试）
type MemoryAuditLogStore struct {
	logs map[string]*AuditLog
	mu   sync.RWMutex
}

// NewMemoryAuditLogStore 创建内存审计日志存储
func NewMemoryAuditLogStore() *MemoryAuditLogStore {
	return &MemoryAuditLogStore{
		logs: make(map[string]*AuditLog),
	}
}

// Save 保存审计日志
func (s *MemoryAuditLogStore) Save(ctx context.Context, log *AuditLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs[log.ID] = log
	return nil
}

// Get 获取审计日志
func (s *MemoryAuditLogStore) Get(ctx context.Context, logID string) (*AuditLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log, exists := s.logs[logID]
	if !exists {
		return nil, ErrAuditLogNotFound
	}

	return log, nil
}

// Query 查询审计日志
func (s *MemoryAuditLogStore) Query(ctx context.Context, filter *AuditLogFilter) ([]*AuditLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*AuditLog

	for _, log := range s.logs {
		// 应用过滤器
		if filter != nil {
			if filter.UserID != "" && log.UserID != filter.UserID {
				continue
			}
			if filter.Action != "" && log.Action != filter.Action {
				continue
			}
			if filter.Resource != "" && log.Resource != filter.Resource {
				continue
			}
			if filter.Result != "" && log.Result != filter.Result {
				continue
			}
			if filter.StartTime != nil && log.Timestamp.Before(*filter.StartTime) {
				continue
			}
			if filter.EndTime != nil && log.Timestamp.After(*filter.EndTime) {
				continue
			}
		}

		results = append(results, log)
	}

	// 应用分页
	if filter != nil {
		if filter.Offset > 0 && filter.Offset < len(results) {
			results = results[filter.Offset:]
		}
		if filter.Limit > 0 && filter.Limit < len(results) {
			results = results[:filter.Limit]
		}
	}

	return results, nil
}

// Delete 删除审计日志
func (s *MemoryAuditLogStore) Delete(ctx context.Context, logID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.logs, logID)
	return nil
}

// AuditLogExporter 审计日志导出器
type AuditLogExporter struct {
	store AuditLogStore
}

// NewAuditLogExporter 创建审计日志导出器
func NewAuditLogExporter(store AuditLogStore) *AuditLogExporter {
	return &AuditLogExporter{
		store: store,
	}
}

// ExportJSON 导出为 JSON 格式
func (e *AuditLogExporter) ExportJSON(ctx context.Context, filter *AuditLogFilter) ([]byte, error) {
	logs, err := e.store.Query(ctx, filter)
	if err != nil {
		return nil, err
	}

	return json.Marshal(logs)
}

// ExportCSV 导出为 CSV 格式
func (e *AuditLogExporter) ExportCSV(ctx context.Context, filter *AuditLogFilter) ([]byte, error) {
	logs, err := e.store.Query(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 简单的 CSV 实现
	csv := "ID,Timestamp,UserID,Action,Resource,ResourceID,Result,IPAddress\n"
	for _, log := range logs {
		csv += fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n",
			log.ID,
			log.Timestamp.Format(time.RFC3339),
			log.UserID,
			log.Action,
			log.Resource,
			log.ResourceID,
			log.Result,
			log.IPAddress,
		)
	}

	return []byte(csv), nil
}
