// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// 本文件包含 ABAC 模块的单元测试。
package abac

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestSubject 测试 Subject 结构
type TestSubject struct {
	ID           string
	Roles        []string
	Department   string
	Organization string
	Attributes   map[string]interface{}
}

// TestResource 测试 Resource 结构
type TestResource struct {
	Type       string
	Owner      string
	ID         string
	Path       string
	Attributes map[string]interface{}
}

// TestAction 测试 Action 结构
type TestAction struct {
	Name       string
	Type       string
	Attributes map[string]interface{}
}

// TestEnvironment 测试 Environment 结构
type TestEnvironment struct {
	Time         int64
	Location     string
	DeviceType   string
	Connection   string
	Timezone     string
	RequestID    string
	Attributes   map[string]interface{}
}

// createTestSubject 创建测试主体
func createTestSubject() Subject {
	return Subject{
		ID:           "user-001",
		Roles:        []string{"user", "developer"},
		Department:   "engineering",
		Organization: "acme-corp",
		Attributes: map[string]interface{}{
			"clearance_level": 5,
			"location":        "office",
		},
	}
}

// createTestResource 创建测试资源
func createTestResource() Resource {
	return Resource{
		Type:   "document",
		Owner:  "user-001",
		ID:     "doc-001",
		Path:   "/documents/doc-001",
		Attributes: map[string]interface{}{
			"sensitivity_level": 3,
			"category":          "technical",
		},
	}
}

// createTestAction 创建测试操作
func createTestAction(name string) Action {
	actionType := "read"
	if name == "write" || name == "create" || name == "update" || name == "delete" {
		actionType = "write"
	}
	return Action{
		Name: name,
		Type: actionType,
	}
}

// createTestEnvironment 创建测试环境
func createTestEnvironment() Environment {
	return Environment{
		Time:       time.Now().Unix(),
		Location:   "192.168.1.100",
		DeviceType: "desktop",
		Connection: "internal",
		Timezone:   "Asia/Shanghai",
		RequestID:  "req-001",
	}
}

// createTestRequest 创建测试请求
func createTestRequest() Request {
	return Request{
		Subject:     createTestSubject(),
		Resource:    createTestResource(),
		Action:      createTestAction("read"),
		Environment: createTestEnvironment(),
	}
}

// ===== Subject 测试 =====

func TestSubjectHasRole(t *testing.T) {
	subject := createTestSubject()

	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"has existing role", "user", true},
		{"has another existing role", "developer", true},
		{"has non-existing role", "admin", false},
		{"case insensitive", "USER", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := subject.HasRole(tt.role)
			if result != tt.expected {
				t.Errorf("HasRole(%s) = %v, want %v", tt.role, result, tt.expected)
			}
		})
	}
}

func TestSubjectHasAnyRole(t *testing.T) {
	subject := createTestSubject()

	tests := []struct {
		name     string
		roles    []string
		expected bool
	}{
		{"has any of existing roles", []string{"admin", "user"}, true},
		{"has none of roles", []string{"admin", "manager"}, false},
		{"empty roles", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := subject.HasAnyRole(tt.roles...)
			if result != tt.expected {
				t.Errorf("HasAnyRole(%v) = %v, want %v", tt.roles, result, tt.expected)
			}
		})
	}
}

func TestSubjectHasAllRoles(t *testing.T) {
	subject := createTestSubject()

	tests := []struct {
		name     string
		roles    []string
		expected bool
	}{
		{"has all roles", []string{"user", "developer"}, true},
		{"missing one role", []string{"user", "admin"}, false},
		{"empty roles", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := subject.HasAllRoles(tt.roles...)
			if result != tt.expected {
				t.Errorf("HasAllRoles(%v) = %v, want %v", tt.roles, result, tt.expected)
			}
		})
	}
}

func TestSubjectGetAttribute(t *testing.T) {
	subject := createTestSubject()

	tests := []struct {
		name          string
		key           string
		expectedValue interface{}
		expectedFound bool
	}{
		{"get id", "id", "user-001", true},
		{"get department", "department", "engineering", true},
		{"get custom attribute", "clearance_level", 5, true},
		{"get non-existing attribute", "non_existing", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, found := subject.GetAttribute(tt.key)
			if found != tt.expectedFound {
				t.Errorf("GetAttribute(%s) found = %v, want %v", tt.key, found, tt.expectedFound)
			}
			if found && !CompareValues(value, tt.expectedValue) {
				t.Errorf("GetAttribute(%s) value = %v, want %v", tt.key, value, tt.expectedValue)
			}
		})
	}
}

// ===== Resource 测试 =====

func TestResourceIsOwnedBy(t *testing.T) {
	resource := createTestResource()

	tests := []struct {
		name      string
		subjectID string
		expected  bool
	}{
		{"is owner", "user-001", true},
		{"not owner", "user-002", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resource.IsOwnedBy(tt.subjectID)
			if result != tt.expected {
				t.Errorf("IsOwnedBy(%s) = %v, want %v", tt.subjectID, result, tt.expected)
			}
		})
	}
}

func TestResourceGetAttribute(t *testing.T) {
	resource := createTestResource()

	tests := []struct {
		name          string
		key           string
		expectedValue interface{}
		expectedFound bool
	}{
		{"get id", "id", "doc-001", true},
		{"get type", "type", "document", true},
		{"get owner", "owner", "user-001", true},
		{"get custom attribute", "sensitivity_level", 3, true},
		{"get non-existing attribute", "non_existing", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, found := resource.GetAttribute(tt.key)
			if found != tt.expectedFound {
				t.Errorf("GetAttribute(%s) found = %v, want %v", tt.key, found, tt.expectedFound)
			}
			if found && !CompareValues(value, tt.expectedValue) {
				t.Errorf("GetAttribute(%s) value = %v, want %v", tt.key, value, tt.expectedValue)
			}
		})
	}
}

// ===== Action 测试 =====

func TestActionIsRead(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		expected bool
	}{
		{"read action", "read", true},
		{"get action", "get", true},
		{"view action", "view", true},
		{"list action", "list", true},
		{"write action", "write", false},
		{"create action", "create", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := createTestAction(tt.action)
			action.Name = tt.action
			result := action.IsRead()
			if result != tt.expected {
				t.Errorf("IsRead() for action %s = %v, want %v", tt.action, result, tt.expected)
			}
		})
	}
}

func TestActionIsWrite(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		expected bool
	}{
		{"write action", "write", true},
		{"create action", "create", true},
		{"update action", "update", true},
		{"delete action", "delete", true},
		{"read action", "read", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := createTestAction(tt.action)
			action.Name = tt.action
			result := action.IsWrite()
			if result != tt.expected {
				t.Errorf("IsWrite() for action %s = %v, want %v", tt.action, result, tt.expected)
			}
		})
	}
}

// ===== Rule 测试 =====

func TestSubjectHasRoleRule(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"has user role", "user", true},
		{"has developer role", "developer", true},
		{"does not have admin role", "admin", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := SubjectHasRole(tt.role)
			result, err := rule.Evaluate(ctx, req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Evaluate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAndRule(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	tests := []struct {
		name     string
		rules    []Rule
		expected bool
	}{
		{
			name:     "all rules match",
			rules:    []Rule{SubjectHasRole("user"), ResourceTypeIs("document")},
			expected: true,
		},
		{
			name:     "one rule does not match",
			rules:    []Rule{SubjectHasRole("admin"), ResourceTypeIs("document")},
			expected: false,
		},
		{
			name:     "no rules",
			rules:    []Rule{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := And(tt.rules...)
			result, err := rule.Evaluate(ctx, req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Evaluate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOrRule(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	tests := []struct {
		name     string
		rules    []Rule
		expected bool
	}{
		{
			name:     "one rule matches",
			rules:    []Rule{SubjectHasRole("admin"), SubjectHasRole("user")},
			expected: true,
		},
		{
			name:     "no rules match",
			rules:    []Rule{SubjectHasRole("admin"), SubjectHasRole("manager")},
			expected: false,
		},
		{
			name:     "all rules match",
			rules:    []Rule{SubjectHasRole("user"), SubjectHasRole("developer")},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := Or(tt.rules...)
			result, err := rule.Evaluate(ctx, req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Evaluate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNotRule(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	tests := []struct {
		name     string
		rule     Rule
		expected bool
	}{
		{
			name:     "negate true",
			rule:     SubjectHasRole("user"),
			expected: false,
		},
		{
			name:     "negate false",
			rule:     SubjectHasRole("admin"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := Not(tt.rule)
			result, err := rule.Evaluate(ctx, req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Evaluate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// ===== Engine 测试 =====

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}

	stats := engine.GetStats()
	if stats.TotalPolicies != 0 {
		t.Errorf("expected 0 policies, got %d", stats.TotalPolicies)
	}
}

func TestEngineAddPolicy(t *testing.T) {
	engine := NewEngine()

	policy := Policy{
		ID:      "policy-001",
		Name:    "Test Policy",
		Effect:  Allow,
		Rules:   SubjectHasRole("admin"),
		Enabled: true,
	}

	err := engine.AddPolicy(policy)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// 添加重复策略应该失败
	err = engine.AddPolicy(policy)
	if err == nil {
		t.Error("expected error for duplicate policy")
	}

	// 验证策略数量
	stats := engine.GetStats()
	if stats.TotalPolicies != 1 {
		t.Errorf("expected 1 policy, got %d", stats.TotalPolicies)
	}
}

func TestEngineGetPolicy(t *testing.T) {
	engine := NewEngine()

	policy := Policy{
		ID:      "policy-001",
		Name:    "Test Policy",
		Effect:  Allow,
		Rules:   SubjectHasRole("admin"),
		Enabled: true,
	}

	engine.AddPolicy(policy)

	// 获取存在的策略
	retrieved, err := engine.GetPolicy("policy-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved.ID != "policy-001" {
		t.Errorf("expected policy ID policy-001, got %s", retrieved.ID)
	}

	// 获取不存在的策略
	_, err = engine.GetPolicy("non-existing")
	if err == nil {
		t.Error("expected error for non-existing policy")
	}
}

func TestEngineRemovePolicy(t *testing.T) {
	engine := NewEngine()

	policy := Policy{
		ID:      "policy-001",
		Name:    "Test Policy",
		Effect:  Allow,
		Rules:   SubjectHasRole("admin"),
		Enabled: true,
	}

	engine.AddPolicy(policy)

	// 删除存在的策略
	err := engine.RemovePolicy("policy-001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// 删除不存在的策略
	err = engine.RemovePolicy("policy-001")
	if err == nil {
		t.Error("expected error for non-existing policy")
	}

	// 验证策略数量
	stats := engine.GetStats()
	if stats.TotalPolicies != 0 {
		t.Errorf("expected 0 policies, got %d", stats.TotalPolicies)
	}
}

func TestEngineEvaluate(t *testing.T) {
	ctx := context.Background()
	engine := NewEngine()

	// 添加允许策略
	allowPolicy := Policy{
		ID:       "policy-allow",
		Name:     "Allow user access",
		Priority: 100,
		Effect:   Allow,
		Rules:    SubjectHasRole("user"),
		Enabled:  true,
	}
	engine.AddPolicy(allowPolicy)

	// 添加拒绝策略
	denyPolicy := Policy{
		ID:       "policy-deny",
		Name:     "Deny banned users",
		Priority: 200, // 更高优先级
		Effect:   Deny,
		Rules:    SubjectHasRole("banned"),
		Enabled:  true,
	}
	engine.AddPolicy(denyPolicy)

	tests := []struct {
		name           string
		subject        Subject
		expectedAllowed bool
	}{
		{
			name: "allow user",
			subject: Subject{
				ID:    "user-001",
				Roles: []string{"user"},
			},
			expectedAllowed: true,
		},
		{
			name: "deny banned user",
			subject: Subject{
				ID:    "user-002",
				Roles: []string{"user", "banned"},
			},
			expectedAllowed: false,
		},
		{
			name: "default deny",
			subject: Subject{
				ID:    "user-003",
				Roles: []string{"guest"},
			},
			expectedAllowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{
				Subject:     tt.subject,
				Resource:    createTestResource(),
				Action:      createTestAction("read"),
				Environment: createTestEnvironment(),
			}

			result := engine.Evaluate(ctx, req)
			if result.IsAllowed() != tt.expectedAllowed {
				t.Errorf("IsAllowed() = %v, want %v, reason: %s", result.IsAllowed(), tt.expectedAllowed, result.Reason)
			}
		})
	}
}

func TestEngineEvaluateDisabledPolicy(t *testing.T) {
	ctx := context.Background()
	engine := NewEngine()

	// 添加禁用的策略
	policy := Policy{
		ID:      "policy-001",
		Name:    "Disabled Policy",
		Effect:  Allow,
		Rules:   SubjectHasRole("user"),
		Enabled: false,
	}
	engine.AddPolicy(policy)

	req := createTestRequest()
	result := engine.Evaluate(ctx, req)

	// 策略被禁用，应该返回默认拒绝
	if result.IsAllowed() {
		t.Error("expected denied when policy is disabled")
	}
}

func TestEngineBatchEvaluate(t *testing.T) {
	ctx := context.Background()
	engine := NewEngine()

	// 添加策略
	policy := Policy{
		ID:      "policy-001",
		Name:    "Allow user access",
		Effect:  Allow,
		Rules:   SubjectHasRole("user"),
		Enabled: true,
	}
	engine.AddPolicy(policy)

	requests := []Request{
		{
			Subject:     Subject{ID: "user-001", Roles: []string{"user"}},
			Resource:    createTestResource(),
			Action:      createTestAction("read"),
			Environment: createTestEnvironment(),
		},
		{
			Subject:     Subject{ID: "user-002", Roles: []string{"admin"}},
			Resource:    createTestResource(),
			Action:      createTestAction("read"),
			Environment: createTestEnvironment(),
		},
	}

	results := engine.BatchEvaluate(ctx, requests)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if !results[0].IsAllowed() {
		t.Error("expected first request to be allowed")
	}

	if results[1].IsAllowed() {
		t.Error("expected second request to be denied")
	}
}

func TestEngineListPolicies(t *testing.T) {
	engine := NewEngine()

	// 添加多个策略
	policies := []Policy{
		{ID: "policy-001", Name: "Low Priority", Priority: 10, Effect: Allow, Rules: AlwaysAllow(), Enabled: true},
		{ID: "policy-002", Name: "High Priority", Priority: 100, Effect: Deny, Rules: AlwaysDeny(), Enabled: true},
		{ID: "policy-003", Name: "Medium Priority", Priority: 50, Effect: Allow, Rules: AlwaysAllow(), Enabled: true},
	}

	for _, policy := range policies {
		engine.AddPolicy(policy)
	}

	list := engine.ListPolicies()
	if len(list) != 3 {
		t.Errorf("expected 3 policies, got %d", len(list))
	}

	// 验证按优先级降序排序
	if list[0].Priority != 100 {
		t.Errorf("expected highest priority first, got %d", list[0].Priority)
	}
	if list[2].Priority != 10 {
		t.Errorf("expected lowest priority last, got %d", list[2].Priority)
	}
}

// ===== Condition 测试 =====

func TestEqualsCondition(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	tests := []struct {
		name     string
		attr     string
		value    interface{}
		expected bool
	}{
		{"subject department equals", "subject.department", "engineering", true},
		{"subject department not equals", "subject.department", "sales", false},
		{"resource type equals", "resource.type", "document", true},
		{"action name equals", "action.name", "read", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Eq(tt.attr, tt.value)
			result, err := cond.Evaluate(ctx, req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Evaluate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGreaterThanCondition(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	tests := []struct {
		name     string
		attr     string
		value    interface{}
		expected bool
	}{
		{"clearance greater than", "subject.clearance_level", 3, true},
		{"clearance equal", "subject.clearance_level", 5, false},
		{"clearance less than", "subject.clearance_level", 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Gt(tt.attr, tt.value)
			result, err := cond.Evaluate(ctx, req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Evaluate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInCondition(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	// 测试属性值在集合中 - 使用 ContainsValue 的方式
	cond := In("action.name", []string{"read", "write", "delete"})
	result, err := cond.Evaluate(ctx, req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !result {
		t.Error("expected true for matching action")
	}
}

// ===== Cache 测试 =====

func TestCacheBasicOperations(t *testing.T) {
	config := CacheConfig{
		MaxSize:    100,
		DefaultTTL: time.Minute,
	}
	cache := NewCache(config)

	// Set and Get
	cache.Set("key1", "value1", 0)
	value, found := cache.Get("key1")
	if !found {
		t.Error("expected to find key1")
	}
	if value != "value1" {
		t.Errorf("expected value1, got %v", value)
	}

	// Get non-existing key
	_, found = cache.Get("non-existing")
	if found {
		t.Error("expected not to find non-existing key")
	}

	// Delete
	cache.Delete("key1")
	_, found = cache.Get("key1")
	if found {
		t.Error("expected key1 to be deleted")
	}
}

func TestCacheExpiration(t *testing.T) {
	config := CacheConfig{
		MaxSize:    100,
		DefaultTTL: time.Millisecond * 50,
	}
	cache := NewCache(config)

	// Set with short TTL
	cache.Set("key1", "value1", time.Millisecond*50)

	// Should be found immediately
	_, found := cache.Get("key1")
	if !found {
		t.Error("expected to find key1 immediately")
	}

	// Wait for expiration
	time.Sleep(time.Millisecond * 100)

	// Should be expired
	_, found = cache.Get("key1")
	if found {
		t.Error("expected key1 to be expired")
	}
}

func TestCacheStats(t *testing.T) {
	config := CacheConfig{
		MaxSize:    100,
		DefaultTTL: time.Minute,
	}
	cache := NewCache(config)

	// Add items
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	// Access items
	cache.Get("key1") // hit
	cache.Get("key1") // hit
	cache.Get("key2") // hit
	cache.Get("key3") // miss

	stats := cache.Stats()
	if stats.Hits != 3 {
		t.Errorf("expected 3 hits, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("expected 1 miss, got %d", stats.Misses)
	}
	if stats.Size != 2 {
		t.Errorf("expected size 2, got %d", stats.Size)
	}
}

func TestCacheLRUEviction(t *testing.T) {
	config := CacheConfig{
		MaxSize:    3,
		DefaultTTL: time.Minute,
	}
	cache := NewCache(config)

	// Add items up to capacity
	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)
	cache.Set("key3", "value3", 0)

	// Access key1 to make it recently used
	cache.Get("key1")

	// Add new item to trigger eviction
	cache.Set("key4", "value4", 0)

	// key2 should be evicted (least recently used)
	_, found := cache.Get("key2")
	if found {
		t.Error("expected key2 to be evicted")
	}

	// key1 should still exist
	_, found = cache.Get("key1")
	if !found {
		t.Error("expected key1 to still exist")
	}
}

// ===== 工具函数测试 =====

func TestCompareValues(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected bool
	}{
		{"both nil", nil, nil, true},
		{"one nil", nil, "value", false},
		{"equal strings", "hello", "hello", true},
		{"different strings", "hello", "world", false},
		{"equal ints", 42, 42, true},
		{"int and float64", 42, 42.0, true},
		{"equal bools", true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareValues(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("CompareValues(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestContainsValue(t *testing.T) {
	tests := []struct {
		name     string
		slice    interface{}
		value    interface{}
		expected bool
	}{
		{"contains string", []string{"a", "b", "c"}, "b", true},
		{"not contains", []string{"a", "b", "c"}, "d", false},
		{"contains int", []int{1, 2, 3}, 2, true},
		{"nil slice", nil, "value", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsValue(tt.slice, tt.value)
			if result != tt.expected {
				t.Errorf("ContainsValue(%v, %v) = %v, want %v", tt.slice, tt.value, result, tt.expected)
			}
		})
	}
}

// ===== 集成测试 =====

func TestABACIntegration(t *testing.T) {
	ctx := context.Background()

	// 创建引擎
	engine := NewEngine(
		WithDefaultEffect(Deny),
	)

	// 添加业务策略
	// 1. 管理员可以执行任何操作
	adminPolicy := Policy{
		ID:          "policy-admin",
		Name:        "Admin full access",
		Description: "管理员拥有所有权限",
		Priority:    1000,
		Effect:      Allow,
		Rules:       SubjectHasRole("admin"),
		Enabled:     true,
	}
	engine.AddPolicy(adminPolicy)

	// 2. 文档所有者可以编辑自己的文档
	ownerPolicy := Policy{
		ID:          "policy-owner",
		Name:        "Owner can edit own documents",
		Description: "文档所有者可以编辑自己的文档",
		Priority:    500,
		Effect:      Allow,
		Rules: And(
			ResourceTypeIs("document"),
			ActionIn("read", "update", "delete"),
			SubjectIsOwner(),
		),
		Enabled: true,
	}
	engine.AddPolicy(ownerPolicy)

	// 3. 工程部成员可以读取敏感级别 <= 3 的文档
	engReadPolicy := Policy{
		ID:          "policy-eng-read",
		Name:        "Engineering can read low sensitivity docs",
		Description: "工程部成员可以读取敏感度 <= 3 的文档",
		Priority:    300,
		Effect:      Allow,
		Rules: And(
			SubjectDepartmentIs("engineering"),
			ResourceTypeIs("document"),
			ActionIsRead(),
			ResourceSensitivityLevelLte(3),
		),
		Enabled: true,
	}
	engine.AddPolicy(engReadPolicy)

	// 4. 拒绝被封禁的用户
	bannedPolicy := Policy{
		ID:          "policy-banned",
		Name:        "Deny banned users",
		Description: "拒绝被封禁的用户",
		Priority:    2000, // 最高优先级
		Effect:      Deny,
		Rules:       SubjectHasRole("banned"),
		Enabled:     true,
	}
	engine.AddPolicy(bannedPolicy)

	// 测试用例
	testCases := []struct {
		name           string
		subject        Subject
		resource       Resource
		action         Action
		expectedAllowed bool
	}{
		{
			name: "admin can do anything",
			subject: Subject{
				ID:    "admin-001",
				Roles: []string{"admin"},
			},
			resource: Resource{Type: "document", ID: "doc-001", Owner: "user-001"},
			action:   Action{Name: "delete"},
			expectedAllowed: true,
		},
		{
			name: "owner can edit own document",
			subject: Subject{
				ID:    "user-001",
				Roles: []string{"user"},
			},
			resource: Resource{Type: "document", ID: "doc-001", Owner: "user-001"},
			action:   Action{Name: "update"},
			expectedAllowed: true,
		},
		{
			name: "non-owner cannot edit others document",
			subject: Subject{
				ID:    "user-002",
				Roles: []string{"user"},
			},
			resource: Resource{Type: "document", ID: "doc-001", Owner: "user-001"},
			action:   Action{Name: "update"},
			expectedAllowed: false,
		},
		{
			name: "engineering can read low sensitivity doc",
			subject: Subject{
				ID:         "user-003",
				Roles:      []string{"user"},
				Department: "engineering",
			},
			resource: Resource{
				Type: "document",
				ID:   "doc-002",
				Attributes: map[string]interface{}{
					"sensitivity_level": 2,
				},
			},
			action: Action{Name: "read", Type: "read"},
			expectedAllowed: true,
		},
		{
			name: "engineering cannot read high sensitivity doc",
			subject: Subject{
				ID:         "user-003",
				Roles:      []string{"user"},
				Department: "engineering",
			},
			resource: Resource{
				Type: "document",
				ID:   "doc-003",
				Attributes: map[string]interface{}{
					"sensitivity_level": 5,
				},
			},
			action: Action{Name: "read"},
			expectedAllowed: false,
		},
		{
			name: "banned user denied even if admin",
			subject: Subject{
				ID:    "user-004",
				Roles: []string{"admin", "banned"},
			},
			resource: Resource{Type: "document", ID: "doc-001"},
			action:   Action{Name: "read"},
			expectedAllowed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := Request{
				Subject:     tc.subject,
				Resource:    tc.resource,
				Action:      tc.action,
				Environment: createTestEnvironment(),
			}

			result := engine.Evaluate(ctx, req)
			if result.IsAllowed() != tc.expectedAllowed {
				t.Errorf("IsAllowed() = %v, want %v, reason: %s", 
					result.IsAllowed(), tc.expectedAllowed, result.Reason)
			}
		})
	}
}

// BenchmarkEngineEvaluate 基准测试
func BenchmarkEngineEvaluate(b *testing.B) {
	ctx := context.Background()
	engine := NewEngine()

	// 添加多个策略
	for i := 0; i < 100; i++ {
		policy := Policy{
			ID:       fmt.Sprintf("policy-%d", i),
			Name:     fmt.Sprintf("Policy %d", i),
			Priority: i,
			Effect:   Allow,
			Rules:    SubjectHasRole(fmt.Sprintf("role-%d", i)),
			Enabled:  true,
		}
		engine.AddPolicy(policy)
	}

	req := createTestRequest()
	req.Subject.Roles = []string{"role-99"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Evaluate(ctx, req)
	}
}

func BenchmarkEngineEvaluateWithCache(b *testing.B) {
	ctx := context.Background()
	engine := NewEngine(
		WithCache(DefaultCacheConfig()),
	)

	policy := Policy{
		ID:      "policy-001",
		Name:    "Test Policy",
		Effect:  Allow,
		Rules:   SubjectHasRole("user"),
		Enabled: true,
	}
	engine.AddPolicy(policy)

	req := createTestRequest()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Evaluate(ctx, req)
	}
}
