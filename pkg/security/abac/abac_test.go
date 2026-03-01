// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// 本文件包含 ABAC 引擎的单元测试

package abac

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewEngine 测试创建 ABAC 引擎
func TestNewEngine(t *testing.T) {
	tests := []struct {
		name string
		opts []EngineOption
		want Effect
	}{
		{
			name: "default engine",
			opts: nil,
			want: Deny,
		},
		{
			name: "with allow default",
			opts: []EngineOption{WithDefaultEffect(Allow)},
			want: Allow,
		},
		{
			name: "with deny default",
			opts: []EngineOption{WithDefaultEffect(Deny)},
			want: Deny,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine(tt.opts...)
			assert.NotNil(t, engine)
			assert.Equal(t, tt.want, engine.defaultEffect)
		})
	}
}

// TestEngine_AddPolicy 测试添加策略
func TestEngine_AddPolicy(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name    string
		policy  Policy
		wantErr bool
	}{
		{
			name: "valid policy",
			policy: Policy{
				ID:      "policy-1",
				Name:    "Test Policy",
				Rules:   AlwaysAllow(),
				Enabled: true,
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			policy: Policy{
				ID:      "",
				Name:    "Test Policy",
				Rules:   AlwaysAllow(),
				Enabled: true,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			policy: Policy{
				ID:      "policy-2",
				Name:    "",
				Rules:   AlwaysAllow(),
				Enabled: true,
			},
			wantErr: true,
		},
		{
			name: "nil rules",
			policy: Policy{
				ID:      "policy-3",
				Name:    "Test Policy",
				Rules:   nil,
				Enabled: true,
			},
			wantErr: true,
		},
		{
			name: "duplicate ID",
			policy: Policy{
				ID:      "policy-1", // 重复的 ID
				Name:    "Another Policy",
				Rules:   AlwaysAllow(),
				Enabled: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.AddPolicy(tt.policy)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestEngine_UpdatePolicy 测试更新策略
func TestEngine_UpdatePolicy(t *testing.T) {
	engine := NewEngine()

	// 添加初始策略
	policy := Policy{
		ID:      "policy-1",
		Name:    "Original Policy",
		Rules:   AlwaysAllow(),
		Enabled: true,
	}
	require.NoError(t, engine.AddPolicy(policy))

	// 更新策略
	updatedPolicy := Policy{
		ID:      "policy-1",
		Name:    "Updated Policy",
		Rules:   AlwaysDeny(),
		Enabled: true,
	}
	err := engine.UpdatePolicy(updatedPolicy)
	assert.NoError(t, err)

	// 验证更新
	retrieved, err := engine.GetPolicy("policy-1")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Policy", retrieved.Name)

	// 更新不存在的策略
	err = engine.UpdatePolicy(Policy{
		ID:      "non-existent",
		Name:    "Policy",
		Rules:   AlwaysAllow(),
		Enabled: true,
	})
	assert.Error(t, err)
}

// TestEngine_RemovePolicy 测试删除策略
func TestEngine_RemovePolicy(t *testing.T) {
	engine := NewEngine()

	// 添加策略
	policy := Policy{
		ID:      "policy-1",
		Name:    "Test Policy",
		Rules:   AlwaysAllow(),
		Enabled: true,
	}
	require.NoError(t, engine.AddPolicy(policy))

	// 删除策略
	err := engine.RemovePolicy("policy-1")
	assert.NoError(t, err)

	// 验证删除
	_, err = engine.GetPolicy("policy-1")
	assert.Error(t, err)

	// 删除不存在的策略
	err = engine.RemovePolicy("non-existent")
	assert.Error(t, err)
}

// TestEngine_GetPolicy 测试获取策略
func TestEngine_GetPolicy(t *testing.T) {
	engine := NewEngine()

	policy := Policy{
		ID:      "policy-1",
		Name:    "Test Policy",
		Rules:   AlwaysAllow(),
		Enabled: true,
	}
	require.NoError(t, engine.AddPolicy(policy))

	// 获取存在的策略
	retrieved, err := engine.GetPolicy("policy-1")
	assert.NoError(t, err)
	assert.Equal(t, policy.ID, retrieved.ID)
	assert.Equal(t, policy.Name, retrieved.Name)

	// 获取不存在的策略
	_, err = engine.GetPolicy("non-existent")
	assert.Error(t, err)
}

// TestEngine_ListPolicies 测试列出策略
func TestEngine_ListPolicies(t *testing.T) {
	engine := NewEngine()

	// 添加多个策略
	policies := []Policy{
		{ID: "policy-1", Name: "Policy 1", Rules: AlwaysAllow(), Enabled: true, Priority: 10},
		{ID: "policy-2", Name: "Policy 2", Rules: AlwaysAllow(), Enabled: true, Priority: 20},
		{ID: "policy-3", Name: "Policy 3", Rules: AlwaysAllow(), Enabled: true, Priority: 5},
	}

	for _, p := range policies {
		require.NoError(t, engine.AddPolicy(p))
	}

	// 列出策略（应该按优先级排序）
	list := engine.ListPolicies()
	assert.Len(t, list, 3)
	assert.Equal(t, "policy-2", list[0].ID) // 优先级 20
	assert.Equal(t, "policy-1", list[1].ID) // 优先级 10
	assert.Equal(t, "policy-3", list[2].ID) // 优先级 5
}

// TestEngine_ClearPolicies 测试清空策略
func TestEngine_ClearPolicies(t *testing.T) {
	engine := NewEngine()

	// 添加策略
	require.NoError(t, engine.AddPolicy(Policy{
		ID:      "policy-1",
		Name:    "Test Policy",
		Rules:   AlwaysAllow(),
		Enabled: true,
	}))

	// 清空
	engine.ClearPolicies()

	// 验证
	list := engine.ListPolicies()
	assert.Len(t, list, 0)
}

// TestEngine_Evaluate 测试策略评估
func TestEngine_Evaluate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		policies []Policy
		request  Request
		want     Effect
	}{
		{
			name: "allow by matching policy",
			policies: []Policy{
				{
					ID:       "allow-admin",
					Name:     "Allow Admin",
					Effect:   Allow,
					Rules:    SubjectHasRole("admin"),
					Enabled:  true,
					Priority: 100,
				},
			},
			request: Request{
				Subject: Subject{ID: "user1", Roles: []string{"admin"}},
			},
			want: Allow,
		},
		{
			name: "deny by matching policy",
			policies: []Policy{
				{
					ID:       "deny-banned",
					Name:     "Deny Banned Users",
					Effect:   Deny,
					Rules:    SubjectHasRole("banned"),
					Enabled:  true,
					Priority: 100,
				},
			},
			request: Request{
				Subject: Subject{ID: "user1", Roles: []string{"banned"}},
			},
			want: Deny,
		},
		{
			name: "no matching policy - default deny",
			policies: []Policy{
				{
					ID:       "allow-admin",
					Name:     "Allow Admin",
					Effect:   Allow,
					Rules:    SubjectHasRole("admin"),
					Enabled:  true,
					Priority: 100,
				},
			},
			request: Request{
				Subject: Subject{ID: "user1", Roles: []string{"user"}},
			},
			want: Deny,
		},
		{
			name: "disabled policy not matched",
			policies: []Policy{
				{
					ID:       "allow-admin",
					Name:     "Allow Admin",
					Effect:   Allow,
					Rules:    SubjectHasRole("admin"),
					Enabled:  false, // 禁用
					Priority: 100,
				},
			},
			request: Request{
				Subject: Subject{ID: "user1", Roles: []string{"admin"}},
			},
			want: Deny, // 默认拒绝
		},
		{
			name: "higher priority wins",
			policies: []Policy{
				{
					ID:       "deny-all",
					Name:     "Deny All",
					Effect:   Deny,
					Rules:    AlwaysAllow(),
					Enabled:  true,
					Priority: 200, // 更高优先级
				},
				{
					ID:       "allow-admin",
					Name:     "Allow Admin",
					Effect:   Allow,
					Rules:    SubjectHasRole("admin"),
					Enabled:  true,
					Priority: 100,
				},
			},
			request: Request{
				Subject: Subject{ID: "user1", Roles: []string{"admin"}},
			},
			want: Deny, // 高优先级的 Deny 策略先匹配
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine()

			for _, p := range tt.policies {
				require.NoError(t, engine.AddPolicy(p))
			}

			result := engine.Evaluate(ctx, tt.request)
			assert.Equal(t, tt.want, result.Decision)
		})
	}
}

// TestEngine_EvaluateWithReason 测试带原因的评估
func TestEngine_EvaluateWithReason(t *testing.T) {
	ctx := context.Background()
	engine := NewEngine()

	// 添加策略
	require.NoError(t, engine.AddPolicy(Policy{
		ID:       "policy-1",
		Name:     "Allow Admin",
		Effect:   Allow,
		Rules:    SubjectHasRole("admin"),
		Enabled:  true,
		Priority: 100,
	}))

	require.NoError(t, engine.AddPolicy(Policy{
		ID:       "policy-2",
		Name:     "Allow User",
		Effect:   Allow,
		Rules:    SubjectHasRole("user"),
		Enabled:  true,
		Priority: 50,
	}))

	// 测试匹配
	result, reasons := engine.EvaluateWithReason(ctx, Request{
		Subject: Subject{ID: "user1", Roles: []string{"admin"}},
	})

	assert.True(t, result.Allowed)
	assert.Equal(t, "policy-1", result.MatchedPolicy.ID)
	assert.True(t, reasons["policy-1"])
	assert.False(t, reasons["policy-2"])
}

// TestEngine_BatchEvaluate 测试批量评估
func TestEngine_BatchEvaluate(t *testing.T) {
	ctx := context.Background()
	engine := NewEngine()

	// 添加策略
	require.NoError(t, engine.AddPolicy(Policy{
		ID:       "allow-admin",
		Name:     "Allow Admin",
		Effect:   Allow,
		Rules:    SubjectHasRole("admin"),
		Enabled:  true,
		Priority: 100,
	}))

	requests := []Request{
		{Subject: Subject{ID: "user1", Roles: []string{"admin"}}},
		{Subject: Subject{ID: "user2", Roles: []string{"user"}}},
		{Subject: Subject{ID: "user3", Roles: []string{"admin"}}},
	}

	results := engine.BatchEvaluate(ctx, requests)

	assert.Len(t, results, 3)
	assert.True(t, results[0].Allowed)  // admin
	assert.False(t, results[1].Allowed) // user
	assert.True(t, results[2].Allowed)  // admin
}

// TestEngine_GetStats 测试统计信息
func TestEngine_GetStats(t *testing.T) {
	engine := NewEngine()

	// 初始状态
	stats := engine.GetStats()
	assert.Equal(t, 0, stats.TotalPolicies)
	assert.Equal(t, 0, stats.EnabledPolicies)

	// 添加策略
	require.NoError(t, engine.AddPolicy(Policy{
		ID:       "policy-1",
		Name:     "Policy 1",
		Rules:    AlwaysAllow(),
		Enabled:  true,
		Priority: 100,
	}))

	require.NoError(t, engine.AddPolicy(Policy{
		ID:       "policy-2",
		Name:     "Policy 2",
		Rules:    AlwaysAllow(),
		Enabled:  false, // 禁用
		Priority: 50,
	}))

	stats = engine.GetStats()
	assert.Equal(t, 2, stats.TotalPolicies)
	assert.Equal(t, 1, stats.EnabledPolicies)
}

// TestEngine_Validate 测试引擎验证
func TestEngine_Validate(t *testing.T) {
	engine := NewEngine()

	// 空引擎验证通过
	assert.NoError(t, engine.Validate())

	// 添加有效策略
	require.NoError(t, engine.AddPolicy(Policy{
		ID:      "valid",
		Name:    "Valid Policy",
		Rules:   AlwaysAllow(),
		Enabled: true,
	}))
	assert.NoError(t, engine.Validate())
}

// TestEngine_Clone 测试克隆引擎
func TestEngine_Clone(t *testing.T) {
	engine := NewEngine(WithDefaultEffect(Allow))

	require.NoError(t, engine.AddPolicy(Policy{
		ID:      "policy-1",
		Name:    "Test Policy",
		Rules:   AlwaysAllow(),
		Enabled: true,
	}))

	// 克隆
	clone := engine.Clone()

	// 验证克隆
	assert.Equal(t, engine.defaultEffect, clone.defaultEffect)
	assert.Len(t, clone.ListPolicies(), 1)

	// 修改原始引擎不应影响克隆
	engine.ClearPolicies()
	assert.Len(t, engine.ListPolicies(), 0)
	assert.Len(t, clone.ListPolicies(), 1)
}

// TestEngine_WithPolicyChangeCallback 测试策略变更回调
func TestEngine_WithPolicyChangeCallback(t *testing.T) {
	var events []PolicyChangeEvent

	callback := func(event PolicyChangeEvent) {
		events = append(events, event)
	}

	engine := NewEngine(WithPolicyChangeCallback(callback))

	// 添加策略
	policy := Policy{
		ID:      "policy-1",
		Name:    "Test Policy",
		Rules:   AlwaysAllow(),
		Enabled: true,
	}
	require.NoError(t, engine.AddPolicy(policy))

	// 更新策略
	policy.Name = "Updated Policy"
	require.NoError(t, engine.UpdatePolicy(policy))

	// 删除策略
	require.NoError(t, engine.RemovePolicy("policy-1"))

	// 验证回调
	assert.Len(t, events, 3)
	assert.Equal(t, "add", events[0].Type)
	assert.Equal(t, "update", events[1].Type)
	assert.Equal(t, "remove", events[2].Type)
}

// TestResult 测试结果结构
func TestResult(t *testing.T) {
	tests := []struct {
		name    string
		result  Result
		allowed bool
		denied  bool
	}{
		{
			name:    "allowed result",
			result:  Result{Allowed: true, Decision: Allow},
			allowed: true,
			denied:  false,
		},
		{
			name:    "denied result",
			result:  Result{Allowed: false, Decision: Deny},
			allowed: false,
			denied:  true,
		},
		{
			name:    "denied by default",
			result:  Result{Allowed: false, Decision: Deny},
			allowed: false,
			denied:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.allowed, tt.result.IsAllowed())
			assert.Equal(t, tt.denied, tt.result.IsDenied())
		})
	}
}

// BenchmarkEngine_Evaluate 评估性能基准测试
func BenchmarkEngine_Evaluate(b *testing.B) {
	ctx := context.Background()
	engine := NewEngine()

	// 添加多个策略
	for i := 0; i < 10; i++ {
		require.NoError(b, engine.AddPolicy(Policy{
			ID:       fmt.Sprintf("policy-%d", i),
			Name:     fmt.Sprintf("Policy %d", i),
			Effect:   Allow,
			Rules:    SubjectHasRole(fmt.Sprintf("role-%d", i)),
			Enabled:  true,
			Priority: i * 10,
		}))
	}

	req := Request{
		Subject: Subject{ID: "user1", Roles: []string{"role-5"}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Evaluate(ctx, req)
	}
}

// TestComplexScenario 测试复杂场景
func TestComplexScenario(t *testing.T) {
	ctx := context.Background()
	engine := NewEngine()

	// 场景：企业文档管理系统
	// 1. 管理员可以执行所有操作
	// 2. 部门经理可以编辑本部门文档
	// 3. 普通用户可以查看文档
	// 4. 文档所有者可以编辑自己的文档
	// 5. 禁止被禁用的用户

	// 策略 1: 禁止被禁用的用户（最高优先级）
	require.NoError(t, engine.AddPolicy(Policy{
		ID:          "deny-banned",
		Name:        "Deny Banned Users",
		Description: "禁止被禁用的用户访问",
		Effect:      Deny,
		Priority:    1000,
		Rules:       SubjectHasRole("banned"),
		Enabled:     true,
	}))

	// 策略 2: 管理员拥有所有权限
	require.NoError(t, engine.AddPolicy(Policy{
		ID:          "allow-admin",
		Name:        "Allow Admin All",
		Description: "管理员可以执行所有操作",
		Effect:      Allow,
		Priority:    900,
		Rules:       SubjectHasRole("admin"),
		Enabled:     true,
	}))

	// 策略 3: 文档所有者可以编辑
	require.NoError(t, engine.AddPolicy(Policy{
		ID:          "allow-owner-edit",
		Name:        "Allow Owner Edit",
		Description: "文档所有者可以编辑自己的文档",
		Effect:      Allow,
		Priority:    800,
		Rules: And(
			ResourceTypeIs("document"),
			ActionIn("edit", "update", "delete"),
			SubjectIsOwner(),
		),
		Enabled: true,
	}))

	// 策略 4: 部门经理可以编辑本部门文档
	require.NoError(t, engine.AddPolicy(Policy{
		ID:          "allow-manager-department",
		Name:        "Allow Manager Edit Department Documents",
		Description: "部门经理可以编辑本部门文档",
		Effect:      Allow,
		Priority:    700,
		Rules: And(
			SubjectHasRole("manager"),
			ResourceTypeIs("document"),
			ActionIn("edit", "update"),
			RuleFromCondition(Eq("subject.department", "resource.attributes.department")),
		),
		Enabled: true,
	}))

	// 策略 5: 普通用户可以查看
	require.NoError(t, engine.AddPolicy(Policy{
		ID:          "allow-user-read",
		Name:        "Allow User Read",
		Description: "普通用户可以查看文档",
		Effect:      Allow,
		Priority:    100,
		Rules: And(
			SubjectHasAnyRole("user", "manager", "admin"),
			ResourceTypeIs("document"),
			ActionIs("read"),
		),
		Enabled: true,
	}))

	tests := []struct {
		name    string
		request Request
		allowed bool
	}{
		{
			name: "admin can do anything",
			request: Request{
				Subject:  Subject{ID: "admin1", Roles: []string{"admin"}},
				Resource: Resource{Type: "document", Owner: "other"},
				Action:   Action{Name: "delete"},
			},
			allowed: true,
		},
		{
			name: "banned user denied",
			request: Request{
				Subject:  Subject{ID: "banned1", Roles: []string{"banned", "admin"}},
				Resource: Resource{Type: "document"},
				Action:   Action{Name: "read"},
			},
			allowed: false,
		},
		{
			name: "owner can edit own document",
			request: Request{
				Subject:  Subject{ID: "user1", Roles: []string{"user"}},
				Resource: Resource{Type: "document", Owner: "user1"},
				Action:   Action{Name: "edit"},
			},
			allowed: true,
		},
		{
			name: "user can read document",
			request: Request{
				Subject:  Subject{ID: "user1", Roles: []string{"user"}},
				Resource: Resource{Type: "document", Owner: "other"},
				Action:   Action{Name: "read"},
			},
			allowed: true,
		},
		{
			name: "user cannot edit others document",
			request: Request{
				Subject:  Subject{ID: "user1", Roles: []string{"user"}},
				Resource: Resource{Type: "document", Owner: "other"},
				Action:   Action{Name: "edit"},
			},
			allowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Evaluate(ctx, tt.request)
			assert.Equal(t, tt.allowed, result.Allowed, "Reason: %s", result.Reason)
		})
	}
}

// TestEnvironmentBasedAccess 测试基于环境的访问控制
func TestEnvironmentBasedAccess(t *testing.T) {
	ctx := context.Background()
	engine := NewEngine()

	// 策略：仅允许在工作时间访问敏感资源
	require.NoError(t, engine.AddPolicy(Policy{
		ID:          "business-hours-only",
		Name:        "Business Hours Only",
		Description: "仅允许在工作时间访问敏感资源",
		Effect:      Allow,
		Priority:    100,
		Rules: And(
			SubjectHasAnyRole("manager", "admin"),
			ResourceAttributeEquals("sensitivity", "high"),
			RuleFromCondition(And(
				Gte("environment.time", parseTime("09:00")),
				Lte("environment.time", parseTime("18:00")),
			)),
		),
		Enabled: true,
	}))

	now := time.Now()
	today9AM := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	today6PM := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())

	tests := []struct {
		name    string
		request Request
		allowed bool
	}{
		{
			name: "access during business hours",
			request: Request{
				Subject:  Subject{ID: "manager1", Roles: []string{"manager"}},
				Resource: Resource{Type: "document", Attributes: map[string]interface{}{"sensitivity": "high"}},
				Action:      Action{Name: "read"},
				Environment: Environment{Time: today9AM.Add(2 * time.Hour).Unix()}, // 11:00 AM
			},
			allowed: true,
		},
		{
			name: "access outside business hours",
			request: Request{
				Subject:  Subject{ID: "manager1", Roles: []string{"manager"}},
				Resource: Resource{Type: "document", Attributes: map[string]interface{}{"sensitivity": "high"}},
				Action:      Action{Name: "read"},
				Environment: Environment{Time: today6PM.Add(2 * time.Hour).Unix()}, // 20:00 PM
			},
			allowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Evaluate(ctx, tt.request)
			assert.Equal(t, tt.allowed, result.Allowed)
		})
	}
}

// parseTime 辅助函数：解析时间字符串为 Unix 时间戳
func parseTime(timeStr string) int64 {
	t, _ := time.Parse("15:04", timeStr)
	return t.Unix()
}


