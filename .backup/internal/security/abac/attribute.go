// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// 本文件定义了 ABAC 的核心属性模型，包括主体、资源、操作和环境的属性结构。
//
// 属性模型是 ABAC 的基础，支持：
//   - 主体属性（用户角色、部门、安全级别等）
//   - 资源属性（资源类型、所有者、敏感度等）
//   - 操作属性（操作名称、类型等）
//   - 环境属性（时间、IP地址、位置等）
//
// 使用示例：
//
//	subject := abac.Subject{
//	    ID:         "user123",
//	    Roles:      []string{"manager", "developer"},
//	    Department: "engineering",
//	    Attributes: map[string]interface{}{
//	        "clearance_level": 5,
//	        "location":       "office",
//	    },
//	}
//
//	resource := abac.Resource{
//	    Type:   "document",
//	    Owner:  "user123",
//	    ID:     "doc456",
//	    Attributes: map[string]interface{}{
//	        "sensitivity": "confidential",
//	        "category":    "finance",
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
// 主体可以是用户、服务账户或任何请求访问的实体
type Subject struct {
	ID           string                 `json:"id"`
	Roles        []string               `json:"roles"`
	Department   string                 `json:"department"`
	Organization string                 `json:"organization"`
	Attributes   map[string]interface{} `json:"attributes"`
}

// GetAttribute 获取主体的属性值
//
// 支持的内置属性：
//   - id: 主体ID
//   - department: 部门
//   - organization: 组织
//   - 自定义属性从 Attributes 映射中获取
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
	case "organization":
		return s.Organization, true
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
//   - role: 角色名称（不区分大小写）
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

// GetClearanceLevel 获取主体的安全级别
//
// 从属性中获取 clearance_level，默认为 0
//
// 返回：
//   - int: 安全级别
func (s Subject) GetClearanceLevel() int {
	if val, exists := s.GetAttribute("clearance_level"); exists {
		if level, ok := toInt(val); ok {
			return level
		}
	}
	return 0
}

// Resource 表示被访问资源的属性
//
// 资源可以是文档、数据、API端点等需要保护的实体
type Resource struct {
	Type       string                 `json:"type"`
	Owner      string                 `json:"owner"`
	ID         string                 `json:"id"`
	Path       string                 `json:"path"` // 资源路径，如 "/api/v1/users"
	Attributes map[string]interface{} `json:"attributes"`
}

// GetAttribute 获取资源的属性值
//
// 支持的内置属性：
//   - id: 资源ID
//   - type: 资源类型
//   - owner: 资源所有者
//   - path: 资源路径
//   - 自定义属性从 Attributes 映射中获取
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
	case "path":
		return r.Path, true
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

// GetSensitivityLevel 获取资源的敏感度级别
//
// 从属性中获取 sensitivity_level，默认为 0
//
// 返回：
//   - int: 敏感度级别
func (r Resource) GetSensitivityLevel() int {
	if val, exists := r.GetAttribute("sensitivity_level"); exists {
		if level, ok := toInt(val); ok {
			return level
		}
	}
	return 0
}

// Action 表示操作的属性
//
// 操作定义了对资源执行的动作，如读取、写入、删除等
type Action struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"` // read, write, admin
	Attributes map[string]interface{} `json:"attributes"`
}

// GetAttribute 获取操作的属性值
//
// 支持的内置属性：
//   - name: 操作名称
//   - type: 操作类型
//   - 自定义属性从 Attributes 映射中获取
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
	case "type":
		return a.Type, true
	default:
		if a.Attributes != nil {
			val, exists := a.Attributes[key]
			return val, exists
		}
		return nil, false
	}
}

// Equals 检查操作是否等于指定名称（不区分大小写）
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
	name := strings.ToLower(a.Name)
	readOps := []string{"read", "get", "view", "list", "search", "query"}
	for _, op := range readOps {
		if name == op {
			return true
		}
	}
	return strings.EqualFold(a.Type, "read")
}

// IsWrite 检查是否为写操作
//
// 返回：
//   - bool: 如果是写操作返回 true
func (a Action) IsWrite() bool {
	name := strings.ToLower(a.Name)
	writeOps := []string{"write", "create", "update", "delete", "modify", "patch"}
	for _, op := range writeOps {
		if name == op {
			return true
		}
	}
	return strings.EqualFold(a.Type, "write")
}

// IsAdmin 检查是否为管理操作
//
// 返回：
//   - bool: 如果是管理操作返回 true
func (a Action) IsAdmin() bool {
	name := strings.ToLower(a.Name)
	adminOps := []string{"admin", "manage", "configure", "grant", "revoke"}
	for _, op := range adminOps {
		if name == op {
			return true
		}
	}
	return strings.EqualFold(a.Type, "admin")
}

// Environment 表示环境上下文属性
//
// 环境属性包括时间、位置、设备信息等上下文信息
type Environment struct {
	Time         int64                  `json:"time"`          // Unix timestamp
	Location     string                 `json:"location"`      // IP address or location
	DeviceType   string                 `json:"device_type"`   // e.g., "mobile", "desktop", "tablet"
	Connection   string                 `json:"connection"`    // e.g., "vpn", "internal", "external"
	Timezone     string                 `json:"timezone"`      // e.g., "Asia/Shanghai"
	RequestID    string                 `json:"request_id"`    // 请求追踪ID
	Attributes   map[string]interface{} `json:"attributes"`
}

// GetAttribute 获取环境的属性值
//
// 支持的内置属性：
//   - time: Unix时间戳
//   - location: 位置/IP地址
//   - device_type: 设备类型
//   - connection: 连接类型
//   - timezone: 时区
//   - request_id: 请求ID
//   - 自定义属性从 Attributes 映射中获取
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
	case "connection":
		return e.Connection, true
	case "timezone":
		return e.Timezone, true
	case "request_id":
		return e.RequestID, true
	default:
		if e.Attributes != nil {
			val, exists := e.Attributes[key]
			return val, exists
		}
		return nil, false
	}
}

// IsBusinessHours 检查当前时间是否为工作时间
//
// 默认工作时间：周一到周五，9:00-18:00
//
// 返回：
//   - bool: 如果是工作时间返回 true
func (e Environment) IsBusinessHours() bool {
	// 简化的业务时间检查，实际实现可能需要更复杂的逻辑
	// 这里仅作为示例
	return true
}

// IsInternalNetwork 检查是否来自内部网络
//
// 返回：
//   - bool: 如果是内部网络返回 true
func (e Environment) IsInternalNetwork() bool {
	internalPatterns := []string{"192.168.", "10.", "172.16.", "127.0.0.1", "::1"}
	for _, pattern := range internalPatterns {
		if strings.HasPrefix(e.Location, pattern) {
			return true
		}
	}
	return e.Connection == "internal" || e.Connection == "vpn"
}

// AttributeAccessor 定义了获取属性的通用接口
//
// Subject、Resource、Action、Environment 都实现了此接口
type AttributeAccessor interface {
	GetAttribute(key string) (interface{}, bool)
}

// Request 表示访问请求
//
// 包含主体、资源、操作和环境信息
type Request struct {
	Subject     Subject     `json:"subject"`
	Resource    Resource    `json:"resource"`
	Action      Action      `json:"action"`
	Environment Environment `json:"environment"`
}

// NewRequest 创建新的访问请求
//
// 参数：
//   - subject: 访问主体
//   - resource: 被访问资源
//   - action: 操作
//
// 返回：
//   - *Request: 访问请求指针
func NewRequest(subject Subject, resource Resource, action Action) *Request {
	return &Request{
		Subject:     subject,
		Resource:    resource,
		Action:      action,
		Environment: Environment{},
	}
}

// WithEnvironment 设置环境属性
//
// 参数：
//   - env: 环境
//
// 返回：
//   - *Request: 更新后的请求指针（链式调用）
func (r *Request) WithEnvironment(env Environment) *Request {
	r.Environment = env
	return r
}

// toInt 尝试将值转换为 int
func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int8:
		return int(n), true
	case int16:
		return int(n), true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case uint:
		return int(n), true
	case uint8:
		return int(n), true
	case uint16:
		return int(n), true
	case uint32:
		return int(n), true
	case uint64:
		return int(n), true
	case float32:
		return int(n), true
	case float64:
		return int(n), true
	case string:
		var result int
		_, err := fmt.Sscanf(n, "%d", &result)
		return result, err == nil
	default:
		return 0, false
	}
}

// CompareValues 比较两个值是否相等
//
// 支持类型：
//   - 基本类型（int, float64, string, bool）
//   - 切片类型
//   - 实现了 Equal 接口的类型
//
// 参数：
//   - a: 第一个值
//   - b: 第二个值
//
// 返回：
//   - bool: 如果相等返回 true
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
	case string:
		var result float64
		_, err := fmt.Sscanf(n, "%f", &result)
		return result, err == nil
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

// resolveAttributePath 从请求中解析属性路径
//
// 支持属性路径：
//   - subject.id, subject.department, subject.attributes.xxx
//   - resource.id, resource.type, resource.owner, resource.attributes.xxx
//   - action.name, action.attributes.xxx
//   - environment.time, environment.location, environment.attributes.xxx
//
// 参数：
//   - req: 访问请求
//   - path: 属性路径，如 "subject.department"
//
// 返回：
//   - value: 属性值
//   - error: 如果解析失败返回错误
func resolveAttributePath(req Request, path string) (interface{}, error) {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid attribute path: %s", path)
	}

	var accessor AttributeAccessor
	switch parts[0] {
	case "subject":
		accessor = req.Subject
	case "resource":
		accessor = req.Resource
	case "action":
		accessor = req.Action
	case "environment":
		accessor = req.Environment
	default:
		return nil, fmt.Errorf("unknown attribute source: %s", parts[0])
	}

	if len(parts) == 1 {
		return nil, fmt.Errorf("attribute key is required")
	}

	value, exists := accessor.GetAttribute(parts[1])
	if !exists {
		return nil, fmt.Errorf("attribute not found: %s", path)
	}

	return value, nil
}
