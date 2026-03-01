// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// 本文件包含属性相关功能的单元测试

package abac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSubject_HasRole 测试角色检查
func TestSubject_HasRole(t *testing.T) {
	subject := Subject{
		ID:    "user1",
		Roles: []string{"admin", "user", "manager"},
	}

	tests := []struct {
		role string
		want bool
	}{
		{"admin", true},
		{"user", true},
		{"manager", true},
		{"superadmin", false},
		{"ADMIN", true},    // 大小写不敏感
		{"User", true},     // 大小写不敏感
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			assert.Equal(t, tt.want, subject.HasRole(tt.role))
		})
	}
}

// TestSubject_HasAnyRole 测试任意角色检查
func TestSubject_HasAnyRole(t *testing.T) {
	subject := Subject{
		ID:    "user1",
		Roles: []string{"admin", "user"},
	}

	tests := []struct {
		name  string
		roles []string
		want  bool
	}{
		{"has one", []string{"superadmin", "admin"}, true},
		{"has none", []string{"superadmin", "manager"}, false},
		{"empty roles", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, subject.HasAnyRole(tt.roles...))
		})
	}
}

// TestSubject_HasAllRoles 测试所有角色检查
func TestSubject_HasAllRoles(t *testing.T) {
	subject := Subject{
		ID:    "user1",
		Roles: []string{"admin", "user", "manager"},
	}

	tests := []struct {
		name  string
		roles []string
		want  bool
	}{
		{"has all", []string{"admin", "user"}, true},
		{"has all 2", []string{"admin", "user", "manager"}, true},
		{"missing one", []string{"admin", "superadmin"}, false},
		{"empty roles", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, subject.HasAllRoles(tt.roles...))
		})
	}
}

// TestSubject_GetAttribute 测试获取属性
func TestSubject_GetAttribute(t *testing.T) {
	subject := Subject{
		ID:         "user1",
		Department: "engineering",
		Attributes: map[string]interface{}{
			"clearance_level": 5,
			"location":       "office",
		},
	}

	tests := []struct {
		key       string
		wantValue interface{}
		wantExist bool
	}{
		{"id", "user1", true},
		{"department", "engineering", true},
		{"clearance_level", 5, true},
		{"location", "office", true},
		{"nonexistent", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			val, exists := subject.GetAttribute(tt.key)
			assert.Equal(t, tt.wantExist, exists)
			if tt.wantExist {
				assert.Equal(t, tt.wantValue, val)
			}
		})
	}
}

// TestResource_GetAttribute 测试资源属性获取
func TestResource_GetAttribute(t *testing.T) {
	resource := Resource{
		Type:  "document",
		Owner: "user1",
		ID:    "doc123",
		Attributes: map[string]interface{}{
			"sensitivity": "confidential",
		},
	}

	tests := []struct {
		key       string
		wantValue interface{}
		wantExist bool
	}{
		{"id", "doc123", true},
		{"type", "document", true},
		{"owner", "user1", true},
		{"sensitivity", "confidential", true},
		{"nonexistent", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			val, exists := resource.GetAttribute(tt.key)
			assert.Equal(t, tt.wantExist, exists)
			if tt.wantExist {
				assert.Equal(t, tt.wantValue, val)
			}
		})
	}
}

// TestResource_IsOwnedBy 测试资源所有者检查
func TestResource_IsOwnedBy(t *testing.T) {
	resource := Resource{
		Type:  "document",
		Owner: "user1",
	}

	assert.True(t, resource.IsOwnedBy("user1"))
	assert.False(t, resource.IsOwnedBy("user2"))
}

// TestAction_GetAttribute 测试操作属性获取
func TestAction_GetAttribute(t *testing.T) {
	action := Action{
		Name: "edit",
		Attributes: map[string]interface{}{
			"severity": "high",
		},
	}

	// 测试 name 属性
	val, exists := action.GetAttribute("name")
	assert.True(t, exists)
	assert.Equal(t, "edit", val)

	// 测试自定义属性
	val, exists = action.GetAttribute("severity")
	assert.True(t, exists)
	assert.Equal(t, "high", val)

	// 测试不存在的属性
	_, exists = action.GetAttribute("nonexistent")
	assert.False(t, exists)
}

// TestAction_Equals 测试操作相等检查
func TestAction_Equals(t *testing.T) {
	action := Action{Name: "edit"}

	assert.True(t, action.Equals("edit"))
	assert.True(t, action.Equals("EDIT")) // 大小写不敏感
	assert.False(t, action.Equals("read"))
}

// TestAction_IsRead 测试读操作检查
func TestAction_IsRead(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		isRead  bool
		isWrite bool
	}{
		{"read", "read", true, false},
		{"READ", "READ", true, false},
		{"get", "get", true, false},
		{"view", "view", true, false},
		{"list", "list", true, false},
		{"write", "write", false, true},
		{"create", "create", false, true},
		{"update", "update", false, true},
		{"delete", "delete", false, true},
		{"custom", "custom", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := Action{Name: tt.action}
			assert.Equal(t, tt.isRead, action.IsRead())
			assert.Equal(t, tt.isWrite, action.IsWrite())
		})
	}
}

// TestEnvironment_GetAttribute 测试环境属性获取
func TestEnvironment_GetAttribute(t *testing.T) {
	env := Environment{
		Time:       1234567890,
		Location:   "office",
		DeviceType: "desktop",
		Attributes: map[string]interface{}{
			"ip_address": "192.168.1.1",
		},
	}

	tests := []struct {
		key       string
		wantValue interface{}
		wantExist bool
	}{
		{"time", int64(1234567890), true},
		{"location", "office", true},
		{"device_type", "desktop", true},
		{"ip_address", "192.168.1.1", true},
		{"nonexistent", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			val, exists := env.GetAttribute(tt.key)
			assert.Equal(t, tt.wantExist, exists)
			if tt.wantExist {
				assert.Equal(t, tt.wantValue, val)
			}
		})
	}
}

// TestCompareValues 测试值比较
func TestCompareValues(t *testing.T) {
	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want bool
	}{
		{"same int", 5, 5, true},
		{"different int", 5, 10, false},
		{"same string", "hello", "hello", true},
		{"different string", "hello", "world", false},
		{"int and int64", int(5), int64(5), true},
		{"int and float64", int(5), float64(5), true},
		{"both nil", nil, nil, true},
		{"one nil", 5, nil, false},
		{"same bool", true, true, true},
		{"different bool", true, false, false},
		{"same slice", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"different slice", []int{1, 2, 3}, []int{1, 2, 4}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CompareValues(tt.a, tt.b))
		})
	}
}

// TestContainsValue 测试包含检查
func TestContainsValue(t *testing.T) {
	tests := []struct {
		name  string
		slice interface{}
		value interface{}
		want  bool
	}{
		{"int in slice", []int{1, 2, 3, 4, 5}, 3, true},
		{"int not in slice", []int{1, 2, 3, 4, 5}, 10, false},
		{"string in slice", []string{"a", "b", "c"}, "b", true},
		{"string not in slice", []string{"a", "b", "c"}, "d", false},
		{"empty slice", []int{}, 1, false},
		{"nil slice", nil, 1, false},
		{"type conversion", []int{1, 2, 3}, int64(2), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ContainsValue(tt.slice, tt.value))
		})
	}
}

// TestGreaterThan 测试大于比较
func TestGreaterThan(t *testing.T) {
	tests := []struct {
		name    string
		a       interface{}
		b       interface{}
		want    bool
		wantErr bool
	}{
		{"5 > 3", 5, 3, true, false},
		{"3 > 5", 3, 5, false, false},
		{"5 > 5", 5, 5, false, false},
		{"int and int64", int(5), int64(3), true, false},
		{"float and int", float64(5.5), 5, true, false},
		{"string comparison", "b", "a", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GreaterThan(tt.a, tt.b)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestLessThan 测试小于比较
func TestLessThan(t *testing.T) {
	tests := []struct {
		name    string
		a       interface{}
		b       interface{}
		want    bool
		wantErr bool
	}{
		{"3 < 5", 3, 5, true, false},
		{"5 < 3", 5, 3, false, false},
		{"5 < 5", 5, 5, false, false},
		{"int and int64", int(3), int64(5), true, false},
		{"string comparison", "a", "b", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LessThan(tt.a, tt.b)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestResolveRequestAttribute 测试请求属性解析
func TestResolveRequestAttribute(t *testing.T) {
	req := Request{
		Subject: Subject{
			ID:    "user1",
			Roles: []string{"admin"},
			Attributes: map[string]interface{}{
				"clearance": 5,
			},
		},
		Resource: Resource{
			Type:  "document",
			Owner: "user1",
		},
		Action: Action{
			Name: "read",
		},
		Environment: Environment{
			Location: "office",
		},
	}

	tests := []struct {
		path    string
		want    interface{}
		wantErr bool
	}{
		{"subject.id", "user1", false},
		{"resource.type", "document", false},
		{"action.name", "read", false},
		{"environment.location", "office", false},
		{"subject.attributes.clearance", 5, false},
		{"invalid.path.here", nil, true},
		{"unknown.field", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, err := resolveRequestAttribute(req, tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestSubject_EmptyRoles 测试空角色处理
func TestSubject_EmptyRoles(t *testing.T) {
	subject := Subject{
		ID:    "user1",
		Roles: []string{},
	}

	assert.False(t, subject.HasRole("admin"))
	assert.False(t, subject.HasAnyRole("admin", "user"))
	assert.True(t, subject.HasAllRoles()) // 空角色列表时返回 true
}

// TestResource_EmptyAttributes 测试空属性处理
func TestResource_EmptyAttributes(t *testing.T) {
	resource := Resource{
		Type: "document",
		// Attributes 为 nil
	}

	val, exists := resource.GetAttribute("nonexistent")
	assert.False(t, exists)
	assert.Nil(t, val)
}
