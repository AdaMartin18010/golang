package framework

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

// Factory 测试数据工厂
type Factory struct {
	seed int64
}

// NewFactory 创建测试数据工厂
func NewFactory() *Factory {
	return &Factory{
		seed: time.Now().UnixNano(),
	}
}

// NewFactoryWithSeed 使用指定种子创建工厂
func NewFactoryWithSeed(seed int64) *Factory {
	return &Factory{
		seed: seed,
	}
}

// String 生成随机字符串
func (f *Factory) String(length int) string {
	if length <= 0 {
		length = 10
	}

	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// Int 生成随机整数
func (f *Factory) Int(min, max int) int {
	if min >= max {
		return min
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(n.Int64()) + min
}

// Int64 生成随机 int64
func (f *Factory) Int64(min, max int64) int64 {
	if min >= max {
		return min
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	return n.Int64() + min
}

// Float64 生成随机浮点数
func (f *Factory) Float64(min, max float64) float64 {
	if min >= max {
		return min
	}

	rangeVal := max - min
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	ratio := float64(n.Int64()) / 1000000.0
	return min + ratio*rangeVal
}

// Bool 生成随机布尔值
func (f *Factory) Bool() bool {
	return f.Int(0, 2) == 1
}

// Email 生成随机邮箱
func (f *Factory) Email() string {
	return fmt.Sprintf("%s@%s.com", f.String(8), f.String(6))
}

// Phone 生成随机手机号
func (f *Factory) Phone() string {
	prefixes := []string{"138", "139", "150", "151", "152", "188", "189"}
	prefix := prefixes[f.Int(0, len(prefixes))]
	return fmt.Sprintf("%s%08d", prefix, f.Int(0, 99999999))
}

// URL 生成随机 URL
func (f *Factory) URL() string {
	domains := []string{"example", "test", "demo", "sample"}
	domain := domains[f.Int(0, len(domains))]
	return fmt.Sprintf("https://%s-%s.com/%s", domain, f.String(6), f.String(8))
}

// UUID 生成随机 UUID（简化版）
func (f *Factory) UUID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		f.String(8),
		f.String(4),
		f.String(4),
		f.String(4),
		f.String(12),
	)
}

// Time 生成随机时间
func (f *Factory) Time(start, end time.Time) time.Time {
	if start.After(end) {
		start, end = end, start
	}

	diff := end.Sub(start)
	seconds := f.Int64(0, int64(diff.Seconds()))
	return start.Add(time.Duration(seconds) * time.Second)
}

// TimeInRange 生成指定范围内的随机时间
func (f *Factory) TimeInRange(daysAgo, daysFromNow int) time.Time {
	now := time.Now()
	start := now.AddDate(0, 0, -daysAgo)
	end := now.AddDate(0, 0, daysFromNow)
	return f.Time(start, end)
}

// Date 生成随机日期
func (f *Factory) Date() time.Time {
	return f.TimeInRange(365, 0)
}

// FutureDate 生成未来日期
func (f *Factory) FutureDate(days int) time.Time {
	return f.TimeInRange(0, days)
}

// PastDate 生成过去日期
func (f *Factory) PastDate(days int) time.Time {
	return f.TimeInRange(days, 0)
}

// Slice 生成随机切片
func (f *Factory) Slice(length int, generator func() interface{}) []interface{} {
	result := make([]interface{}, length)
	for i := 0; i < length; i++ {
		result[i] = generator()
	}
	return result
}

// StringSlice 生成随机字符串切片
func (f *Factory) StringSlice(length int) []string {
	result := make([]string, length)
	for i := 0; i < length; i++ {
		result[i] = f.String(10)
	}
	return result
}

// IntSlice 生成随机整数切片
func (f *Factory) IntSlice(length, min, max int) []int {
	result := make([]int, length)
	for i := 0; i < length; i++ {
		result[i] = f.Int(min, max)
	}
	return result
}

// Map 生成随机 map
func (f *Factory) Map(length int, keyGenerator, valueGenerator func() interface{}) map[interface{}]interface{} {
	result := make(map[interface{}]interface{})
	for i := 0; i < length; i++ {
		key := keyGenerator()
		value := valueGenerator()
		result[key] = value
	}
	return result
}

// StringMap 生成随机字符串 map
func (f *Factory) StringMap(length int) map[string]string {
	result := make(map[string]string)
	for i := 0; i < length; i++ {
		result[f.String(8)] = f.String(10)
	}
	return result
}

// UserFactory 用户数据工厂
type UserFactory struct {
	*Factory
}

// NewUserFactory 创建用户数据工厂
func NewUserFactory() *UserFactory {
	return &UserFactory{
		Factory: NewFactory(),
	}
}

// User 生成用户数据
func (f *UserFactory) User() map[string]interface{} {
	return map[string]interface{}{
		"id":        f.UUID(),
		"email":     f.Email(),
		"phone":     f.Phone(),
		"name":      f.String(10),
		"created_at": f.Date(),
		"updated_at": f.Date(),
		"active":    f.Bool(),
	}
}

// OAuth2ClientFactory OAuth2 客户端工厂
type OAuth2ClientFactory struct {
	*Factory
}

// NewOAuth2ClientFactory 创建 OAuth2 客户端工厂
func NewOAuth2ClientFactory() *OAuth2ClientFactory {
	return &OAuth2ClientFactory{
		Factory: NewFactory(),
	}
}

// Client 生成 OAuth2 客户端数据
func (f *OAuth2ClientFactory) Client() map[string]interface{} {
	return map[string]interface{}{
		"id":           f.UUID(),
		"secret":       f.String(32),
		"redirect_uris": f.StringSlice(f.Int(1, 5)),
		"grant_types":  f.StringSlice(f.Int(1, 3)),
		"scopes":       f.StringSlice(f.Int(1, 5)),
		"created_at":   f.Date(),
	}
}

// TokenFactory 令牌工厂
type TokenFactory struct {
	*Factory
}

// NewTokenFactory 创建令牌工厂
func NewTokenFactory() *TokenFactory {
	return &TokenFactory{
		Factory: NewFactory(),
	}
}

// Token 生成令牌数据
func (f *TokenFactory) Token() map[string]interface{} {
	return map[string]interface{}{
		"access_token":  f.String(40),
		"refresh_token": f.String(40),
		"token_type":    "Bearer",
		"expires_in":    f.Int(3600, 7200),
		"scope":         f.String(20),
		"created_at":    time.Now(),
		"expires_at":    time.Now().Add(time.Duration(f.Int(3600, 7200)) * time.Second),
	}
}

// AuditLogFactory 审计日志工厂
type AuditLogFactory struct {
	*Factory
}

// NewAuditLogFactory 创建审计日志工厂
func NewAuditLogFactory() *AuditLogFactory {
	return &AuditLogFactory{
		Factory: NewFactory(),
	}
}

// AuditLog 生成审计日志数据
func (f *AuditLogFactory) AuditLog() map[string]interface{} {
	actions := []string{"create", "update", "delete", "read", "access"}
	resources := []string{"user", "order", "product", "api", "system"}

	return map[string]interface{}{
		"id":          f.UUID(),
		"user_id":     f.UUID(),
		"action":      actions[f.Int(0, len(actions))],
		"resource":    resources[f.Int(0, len(resources))],
		"resource_id": f.UUID(),
		"ip_address":  fmt.Sprintf("%d.%d.%d.%d", f.Int(1, 255), f.Int(1, 255), f.Int(1, 255), f.Int(1, 255)),
		"user_agent":  f.String(50),
		"timestamp":   time.Now(),
	}
}
