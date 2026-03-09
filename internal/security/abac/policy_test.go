// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// policy_test.go 包含策略相关的单元测试
package abac

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEffect_String 测试 Effect 字符串表示
func TestEffect_String(t *testing.T) {
	tests := []struct {
		name     string
		effect   Effect
		expected string
	}{
		{"Allow", Allow, "Allow"},
		{"Deny", Deny, "Deny"},
		{"Unknown", Effect(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.effect.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestEffect_MarshalJSON 测试 Effect JSON 序列化
func TestEffect_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		effect   Effect
		expected string
	}{
		{"Allow", Allow, `"Allow"`},
		{"Deny", Deny, `"Deny"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.effect.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

// TestEffect_UnmarshalJSON 测试 Effect JSON 反序列化
func TestEffect_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected Effect
		wantErr  bool
	}{
		{"allow lowercase", `"allow"`, Allow, false},
		{"Allow uppercase", `"Allow"`, Allow, false},
		{"deny lowercase", `"deny"`, Deny, false},
		{"Deny uppercase", `"Deny"`, Deny, false},
		{"unknown effect", `"unknown"`, Deny, true},
		{"invalid json", `invalid`, Deny, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var effect Effect
			err := effect.UnmarshalJSON([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, effect)
			}
		})
	}
}

// TestPolicy_Validate 测试策略验证
func TestPolicy_Validate(t *testing.T) {
	tests := []struct {
		name    string
		policy  Policy
		wantErr bool
	}{
		{
			name:    "valid policy",
			policy:  Policy{ID: "p1", Name: "Test", Effect: Allow, Rules: AlwaysAllow()},
			wantErr: false,
		},
		{
			name:    "missing ID",
			policy:  Policy{ID: "", Name: "Test", Effect: Allow, Rules: AlwaysAllow()},
			wantErr: true,
		},
		{
			name:    "missing Name",
			policy:  Policy{ID: "p1", Name: "", Effect: Allow, Rules: AlwaysAllow()},
			wantErr: true,
		},
		{
			name:    "missing Rules and RulesJSON",
			policy:  Policy{ID: "p1", Name: "Test", Effect: Allow},
			wantErr: true,
		},
		{
			name:    "with RulesJSON only",
			policy:  Policy{ID: "p1", Name: "Test", Effect: Allow, RulesJSON: `{"type":"always"}`},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestPolicy_Match_Disabled 测试禁用策略的匹配
func TestPolicy_Match_Disabled(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	policy := Policy{
		ID:      "disabled-policy",
		Name:    "Disabled Policy",
		Effect:  Allow,
		Rules:   AlwaysAllow(),
		Enabled: false,
	}

	matched, err := policy.Match(ctx, req)
	assert.NoError(t, err)
	assert.False(t, matched)
}

// TestPolicy_Match_NilRules 测试 nil 规则的匹配
func TestPolicy_Match_NilRules(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	policy := Policy{
		ID:      "nil-rules-policy",
		Name:    "Nil Rules Policy",
		Effect:  Allow,
		Rules:   nil,
		Enabled: true,
	}

	matched, err := policy.Match(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rules are nil")
	assert.False(t, matched)
}

// TestNewRuleFunc 测试创建规则函数
func TestNewRuleFunc(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	rule := NewRuleFunc("test rule", func(ctx context.Context, req Request) (bool, error) {
		return true, nil
	})

	result, err := rule.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.True(t, result)
	assert.Equal(t, "test rule", rule.String())
	assert.Equal(t, "custom", rule.Type())
}

// TestRuleFunc_WithType 测试带类型的规则函数
func TestRuleFunc_WithType(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	rule := RuleFunc{
		fn:   func(ctx context.Context, req Request) (bool, error) { return true, nil },
		desc: "typed rule",
		typ:  "special",
	}

	result, err := rule.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.True(t, result)
	assert.Equal(t, "typed rule", rule.String())
	assert.Equal(t, "special", rule.Type())
}

// TestAndRule_String 测试 And 规则的 String 方法
func TestAndRule_String(t *testing.T) {
	rule := And(SubjectHasRole("admin"), ResourceTypeIs("document"))
	str := rule.String()
	assert.Contains(t, str, "AND")
}

// TestOrRule_String 测试 Or 规则的 String 方法
func TestOrRule_String(t *testing.T) {
	rule := Or(SubjectHasRole("admin"), SubjectHasRole("user"))
	str := rule.String()
	assert.Contains(t, str, "OR")
}

// TestNotRule_String 测试 Not 规则的 String 方法
func TestNotRule_String(t *testing.T) {
	rule := Not(SubjectHasRole("admin"))
	str := rule.String()
	assert.Contains(t, str, "NOT")
}

// TestSubjectHasAnyRole_String 测试任意角色检查的 String 方法
func TestSubjectHasAnyRole_String(t *testing.T) {
	rule := SubjectHasAnyRole("admin", "manager", "user")
	str := rule.String()
	assert.Contains(t, str, "SubjectHasAnyRole")
}

// TestSubjectHasAllRoles_String 测试所有角色检查的 String 方法
func TestSubjectHasAllRoles_String(t *testing.T) {
	rule := SubjectHasAllRoles("user", "developer")
	str := rule.String()
	assert.Contains(t, str, "SubjectHasAllRoles")
}

// TestSubjectDepartmentIs 测试部门检查
func TestSubjectDepartmentIs(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	rule := SubjectDepartmentIs("engineering")
	result, err := rule.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.True(t, result)
}

// TestSubjectClearanceLevelGte 测试安全级别检查
func TestSubjectClearanceLevelGte(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	rule := SubjectClearanceLevelGte(3)
	result, err := rule.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.True(t, result) // clearance_level is 5
}

// TestResourceSensitivityLevelLte 测试资源敏感度检查
func TestResourceSensitivityLevelLte(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	rule := ResourceSensitivityLevelLte(5)
	result, err := rule.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.True(t, result) // sensitivity_level is 3
}

// TestActionIn_String 测试操作列表检查的 String 方法
func TestActionIn_String(t *testing.T) {
	rule := ActionIn("read", "write", "delete")
	str := rule.String()
	assert.Contains(t, str, "ActionIn")
}

// TestActionIsRead_String 测试读操作检查的 String 方法
func TestActionIsRead_String(t *testing.T) {
	rule := ActionIsRead()
	str := rule.String()
	assert.Equal(t, "ActionIsRead()", str)
}

// TestActionIsWrite_String 测试写操作检查的 String 方法
func TestActionIsWrite_String(t *testing.T) {
	rule := ActionIsWrite()
	str := rule.String()
	assert.Equal(t, "ActionIsWrite()", str)
}

// TestEnvironmentIsInternalNetwork_String 测试内部网络检查的 String 方法
func TestEnvironmentIsInternalNetwork_String(t *testing.T) {
	rule := EnvironmentIsInternalNetwork()
	str := rule.String()
	assert.Equal(t, "EnvironmentIsInternalNetwork()", str)
}

// TestAlwaysAllow_String 测试总是允许规则的 String 方法
func TestAlwaysAllow_String(t *testing.T) {
	rule := AlwaysAllow()
	str := rule.String()
	assert.Equal(t, "AlwaysAllow()", str)
}

// TestAlwaysDeny_String 测试总是拒绝规则的 String 方法
func TestAlwaysDeny_String(t *testing.T) {
	rule := AlwaysDeny()
	str := rule.String()
	assert.Equal(t, "AlwaysDeny()", str)
}

// TestAndRule_Empty 测试空的 And 规则
func TestAndRule_Empty(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	rule := And()
	result, err := rule.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.True(t, result)
}

// TestOrRule_Empty 测试空的 Or 规则
func TestOrRule_Empty(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	rule := Or()
	result, err := rule.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.False(t, result)
}

// TestPolicy_JSONSerialization 测试策略 JSON 序列化
func TestPolicy_JSONSerialization(t *testing.T) {
	policy := Policy{
		ID:     "policy-1",
		Name:   "Test Policy",
		Effect: Allow,
	}

	data, err := json.Marshal(policy)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"id":"policy-1"`)
	assert.Contains(t, string(data), `"effect":"Allow"`)
}
