// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// 属性访问控制 (ABAC) 实现，支持基于属性的细粒度权限控制。
// ABAC 通过评估主体（Subject）、资源（Resource）、操作（Action）和环境（Environment）
// 的属性来决定访问是否被允许。
//
// 核心概念：
//   - Subject: 访问主体（用户、服务等）的属性
//   - Resource: 被访问资源的属性
//   - Action: 操作的属性
//   - Environment: 环境上下文属性
//
// 使用示例：
//
//	subject := abac.Subject{
//	    ID:       "user1",
//	    Roles:    []string{"manager"},
//	    Department: "engineering",
//	    Attributes: map[string]interface{}{
//	        "clearance_level": 5,
//	        "location":       "office",
//	    },
//	}
//
//	resource := abac.Resource{
//	    Type:   "document",
//	    Owner:  "user1",
//	    ID:     "doc123",
//	    Attributes: map[string]interface{}{
//	        "sensitivity": "confidential",
//	        "department":  "engineering",
//	    },
//	}
package abac

import (
	"fmt"
	"reflect"
	"strings"
)

// Subject 表示访问主体的属性
//
// 包含用户或服务的身份信息、角色、部门以及自定义属性
type Subject struct {
	ID         string                 `json:"id"`
	Roles      []string               `json:"roles"`
	Department string                 `json:"department"`
	Attributes map[string]interface{} `json:"attributes"`
}

// GetAttribute 获取主体的属性值
//
// 参数：
//   - key: 属性键名
//
// 返回：
//   - value: 属性值
//   - exists: 属性是否存在
func (s Subject) GetAttribute(key string) (interface{}, bool) {
	switch key {
	case "id":
		return s.ID, true
	case "department":
		return s.Department, true
	default:
		if s.Attributes != nil {
			val, exists := s.Attributes[key]
			return val, exists
		}
		return nil, false
	}
}

// HasRole 检查主体是否拥有指定角色
//
// 参数：
//   - role: 角色名称
//
// 返回：
//   - bool: 如果拥有该角色返回 true
func (s Subject) HasRole(role string) bool {
	for _, r := range s.Roles {
		if strings.EqualFold(r, role) {
			return true
		}
	}
	return false
}

// HasAnyRole 检查主体是否拥有任意一个指定角色
//
// 参数：
//   - roles: 角色列表
//
// 返回：
//   - bool: 如果拥有至少一个角色返回 true
func (s Subject) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if s.HasRole(role) {
			return true
		}
	}
	return false
}

// HasAllRoles 检查主体是否拥有所有指定角色
//
// 参数：
//   - roles: 角色列表
//
// 返回：
//   - bool: 如果拥有所有角色返回 true
func (s Subject) HasAllRoles(roles ...string) bool {
	for _, role := range roles {
		if !s.HasRole(role) {
			return false
		}
	}
	return true
}

// Resource 表示被访问资源的属性
//
// 包含资源的类型、所有者、标识以及自定义属性
type Resource struct {
	Type       string                 `json:"type"`
	Owner      string                 `json:"owner"`
	ID         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

// GetAttribute 获取资源的属性值
//
// 参数：
//   - key: 属性键名
//
// 返回：
//   - value: 属性值
//   - exists: 属性是否存在
func (r Resource) GetAttribute(key string) (interface{}, bool) {
	switch key {
	case "id":
		return r.ID, true
	case "type":
		return r.Type, true
	case "owner":
		return r.Owner, true
	default:
		if r.Attributes != nil {
			val, exists := r.Attributes[key]
			return val, exists
		}
		return nil, false
	}
}

// IsOwnedBy 检查资源是否属于指定主体
//
// 参数：
//   - subjectID: 主体ID
//
// 返回：
//   - bool: 如果资源属于该主体返回 true
func (r Resource) IsOwnedBy(subjectID string) bool {
	return r.Owner == subjectID
}

// Action 表示操作的属性
//
// 包含操作名称和自定义属性
type Action struct {
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"attributes"`
}

// GetAttribute 获取操作的属性值
//
// 参数：
//   - key: 属性键名
//
// 返回：
//   - value: 属性值
//   - exists: 属性是否存在
func (a Action) GetAttribute(key string) (interface{}, bool) {
	switch key {
	case "name":
		return a.Name, true
	default:
		if a.Attributes != nil {
			val, exists := a.Attributes[key]
			return val, exists
		}
		return nil, false
	}
}

// Equals 检查操作是否等于指定名称
//
// 参数：
//   - name: 操作名称
//
// 返回：
//   - bool: 如果匹配返回 true
func (a Action) Equals(name string) bool {
	return strings.EqualFold(a.Name, name)
}

// IsRead 检查是否为读操作
//
// 返回：
//   - bool: 如果是读操作返回 true
func (a Action) IsRead() bool {
	return strings.EqualFold(a.Name, "read") ||
		strings.EqualFold(a.Name, "get") ||
		strings.EqualFold(a.Name, "view") ||
		strings.EqualFold(a.Name, "list")
}

// IsWrite 检查是否为写操作
//
// 返回：
//   - bool: 如果是写操作返回 true
func (a Action) IsWrite() bool {
	return strings.EqualFold(a.Name, "write") ||
		strings.EqualFold(a.Name, "create") ||
		strings.EqualFold(a.Name, "update") ||
		strings.EqualFold(a.Name, "delete")
}

// Environment 表示环境上下文属性
//
// 包含时间、位置、设备信息等环境因素
type Environment struct {
	Time       int64                  `json:"time"`       // Unix timestamp
	Location   string                 `json:"location"`   // e.g., "192.168.1.1", "office"
	DeviceType string                 `json:"device_type"` // e.g., "mobile", "desktop"
	Attributes map[string]interface{} `json:"attributes"`
}

// GetAttribute 获取环境的属性值
//
// 参数：
//   - key: 属性键名
//
// 返回：
//   - value: 属性值
//   - exists: 属性是否存在
func (e Environment) GetAttribute(key string) (interface{}, bool) {
	switch key {
	case "time":
		return e.Time, true
	case "location":
		return e.Location, true
	case "device_type":
		return e.DeviceType, true
	default:
		if e.Attributes != nil {
			val, exists := e.Attributes[key]
			return val, exists
		}
		return nil, false
	}
}

// AttributeAccessor 定义了获取属性的通用接口
type AttributeAccessor interface {
	GetAttribute(key string) (interface{}, bool)
}

// AttributeResolver 用于解析和比较属性值的辅助函数集合

// ResolveAttribute 从访问器中解析属性值
//
// 参数：
//   - accessor: 属性访问器
//   - key: 属性键名，支持点号表示法如 "subject.id"
//
// 返回：
//   - value: 属性值
//   - exists: 属性是否存在
func ResolveAttribute(accessor AttributeAccessor, key string) (interface{}, bool) {
	// 处理点号表示法
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		return accessor.GetAttribute(key)
	}

	// 递归解析嵌套属性
	value, exists := accessor.GetAttribute(parts[0])
	if !exists {
		return nil, false
	}

	// 对剩余部分进行解析
	for i := 1; i < len(parts); i++ {
		value = resolveNestedField(value, parts[i])
		if value == nil {
			return nil, false
		}
	}

	return value, true
}

// resolveNestedField 解析嵌套字段
func resolveNestedField(obj interface{}, field string) interface{} {
	if obj == nil {
		return nil
	}

	// 如果是 map，直接查找
	if m, ok := obj.(map[string]interface{}); ok {
		return m[field]
	}

	// 使用反射获取字段
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	// 尝试通过字段名获取
	fieldValue := val.FieldByName(field)
	if fieldValue.IsValid() {
		return fieldValue.Interface()
	}

	// 尝试通过 JSON tag 获取
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		tag := f.Tag.Get("json")
		if tag != "" {
			tagParts := strings.Split(tag, ",")
			if tagParts[0] == field {
				return val.Field(i).Interface()
			}
		}
	}

	return nil
}

// CompareValues 比较两个值是否相等
//
// 支持类型：
//   - 基本类型（int, float64, string, bool）
//   - 切片类型
//   - 实现了 Equal 接口的类型
func CompareValues(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// 尝试转换为相同类型后比较
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	// 如果类型相同，直接比较
	if aVal.Type() == bVal.Type() {
		return reflect.DeepEqual(a, b)
	}

	// 尝试数值比较
	if aNum, ok := toFloat64(a); ok {
		if bNum, ok := toFloat64(b); ok {
			return aNum == bNum
		}
	}

	// 字符串比较
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// toFloat64 尝试将值转换为 float64
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}

// ContainsValue 检查切片是否包含指定值
//
// 参数：
//   - slice: 切片
//   - value: 要查找的值
//
// 返回：
//   - bool: 如果包含返回 true
func ContainsValue(slice interface{}, value interface{}) bool {
	if slice == nil {
		return false
	}

	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice && s.Kind() != reflect.Array {
		return CompareValues(slice, value)
	}

	for i := 0; i < s.Len(); i++ {
		if CompareValues(s.Index(i).Interface(), value) {
			return true
		}
	}
	return false
}

// GreaterThan 比较两个值（a > b）
//
// 参数：
//   - a: 第一个值
//   - b: 第二个值
//
// 返回：
//   - bool: 如果 a > b 返回 true
//   - error: 如果无法比较返回错误
func GreaterThan(a, b interface{}) (bool, error) {
	aNum, ok1 := toFloat64(a)
	bNum, ok2 := toFloat64(b)

	if !ok1 || !ok2 {
		// 尝试字符串比较
		aStr := fmt.Sprintf("%v", a)
		bStr := fmt.Sprintf("%v", b)
		return aStr > bStr, nil
	}

	return aNum > bNum, nil
}

// LessThan 比较两个值（a < b）
//
// 参数：
//   - a: 第一个值
//   - b: 第二个值
//
// 返回：
//   - bool: 如果 a < b 返回 true
//   - error: 如果无法比较返回错误
func LessThan(a, b interface{}) (bool, error) {
	aNum, ok1 := toFloat64(a)
	bNum, ok2 := toFloat64(b)

	if !ok1 || !ok2 {
		// 尝试字符串比较
		aStr := fmt.Sprintf("%v", a)
		bStr := fmt.Sprintf("%v", b)
		return aStr < bStr, nil
	}

	return aNum < bNum, nil
}
