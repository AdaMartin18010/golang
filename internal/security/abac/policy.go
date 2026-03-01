// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// 本文件定义了 ABAC 的策略结构、效果和规则。
//
// 策略（Policy）是 ABAC 的核心概念，它定义了在何种条件下允许或拒绝访问。
// 每个策略包含：
//   - 名称和描述
//   - 效果（Allow 或 Deny）
//   - 优先级（数字越大优先级越高）
//   - 规则（Rule）：决定策略是否适用的条件
//
// 策略评估流程：
//  1. 按优先级排序策略
//  2. 依次评估每个策略的规则
//  3. 规则匹配时，返回策略的效果
//  4. 如果没有匹配的策略，返回默认效果（通常为 Deny）
//
// 使用示例：
//
//	policy := abac.Policy{
//	    ID:          "policy-001",
//	    Name:        "Allow document owners to edit",
//	    Description: "文档所有者可以编辑自己的文档",
//	    Priority:    100,
//	    Effect:      abac.Allow,
//	    Rules: abac.And(
//	        abac.SubjectHasRole("user"),
//	        abac.ResourceTypeIs("document"),
//	        abac.ActionIs("edit"),
//	        abac.SubjectIsOwner(),
//	    ),
//	    Enabled: true,
//	}
//
//	if err := policy.Validate(); err != nil {
//	    log.Fatal(err)
//	}
package abac

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Effect 表示策略的效果类型
type Effect int

const (
	// Deny 拒绝访问
	Deny Effect = iota
	// Allow 允许访问
	Allow
)

// String 返回效果的字符串表示
func (e Effect) String() string {
	switch e {
	case Allow:
		return "Allow"
	case Deny:
		return "Deny"
	default:
		return "Unknown"
	}
}

// MarshalJSON 实现 json.Marshaler 接口
func (e Effect) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (e *Effect) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "allow":
		*e = Allow
	case "deny":
		*e = Deny
	default:
		return fmt.Errorf("unknown effect: %s", s)
	}
	return nil
}

// Policy 表示一个访问控制策略
type Policy struct {
	ID          string `json:"id"`
	Version     string `json:"version"` // 策略版本，用于版本控制
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"` // 优先级，数字越大优先级越高
	Effect      Effect `json:"effect"`
	Rules       Rule   `json:"-"`        // 规则（不直接序列化）
	RulesJSON   string `json:"rules"`    // 规则的 JSON 表示
	Enabled     bool   `json:"enabled"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	CreatedBy   string `json:"created_by"`
}

// Validate 验证策略的有效性
//
// 验证规则：
//   - ID 不能为空
//   - Name 不能为空
//   - Rules 不能为空
//
// 返回：
//   - error: 验证失败时返回错误
func (p Policy) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("policy ID is required")
	}
	if p.Name == "" {
		return fmt.Errorf("policy name is required")
	}
	if p.Rules == nil && p.RulesJSON == "" {
		return fmt.Errorf("policy rules are required")
	}
	return nil
}

// Match 检查策略是否匹配给定的请求
//
// 参数：
//   - ctx: 上下文
//   - req: 访问请求
//
// 返回：
//   - bool: 如果策略规则匹配返回 true
//   - error: 评估过程中的错误
func (p Policy) Match(ctx context.Context, req Request) (bool, error) {
	if !p.Enabled {
		return false, nil
	}
	if p.Rules == nil {
		return false, fmt.Errorf("policy rules are nil")
	}
	return p.Rules.Evaluate(ctx, req)
}

// Rule 是策略规则的接口
//
// 所有规则类型都实现此接口
type Rule interface {
	// Evaluate 评估规则
	//
	// 参数：
	//   - ctx: 上下文
	//   - req: 访问请求
	//
	// 返回：
	//   - bool: 如果规则匹配返回 true
	//   - error: 评估过程中的错误
	Evaluate(ctx context.Context, req Request) (bool, error)

	// String 返回规则的可读描述
	String() string

	// Type 返回规则类型
	Type() string
}

// RuleFunc 是一个函数类型，实现 Rule 接口
type RuleFunc struct {
	fn   func(ctx context.Context, req Request) (bool, error)
	desc string
	typ  string
}

// Evaluate 实现 Rule 接口
func (f RuleFunc) Evaluate(ctx context.Context, req Request) (bool, error) {
	return f.fn(ctx, req)
}

// String 返回规则描述
func (f RuleFunc) String() string {
	return f.desc
}

// Type 返回规则类型
func (f RuleFunc) Type() string {
	if f.typ == "" {
		return "custom"
	}
	return f.typ
}

// NewRuleFunc 创建新的规则函数
//
// 参数：
//   - desc: 规则描述
//   - fn: 规则评估函数
//
// 返回：
//   - Rule: 规则接口
func NewRuleFunc(desc string, fn func(ctx context.Context, req Request) (bool, error)) Rule {
	return RuleFunc{fn: fn, desc: desc}
}

// compoundRule 是复合规则的基结构
type compoundRule struct {
	op     string
	rules  []Rule
	typeName string
}

// And 创建逻辑与规则
//
// 所有子规则都必须匹配
//
// 示例：
//
//	rule := abac.And(
//	    abac.SubjectHasRole("manager"),
//	    abac.ResourceTypeIs("document"),
//	)
func And(rules ...Rule) Rule {
	return &andRule{compoundRule{op: "AND", rules: rules, typeName: "and"}}
}

// Or 创建逻辑或规则
//
// 至少一个子规则必须匹配
//
// 示例：
//
//	rule := abac.Or(
//	    abac.SubjectHasRole("admin"),
//	    abac.SubjectHasRole("manager"),
//	)
func Or(rules ...Rule) Rule {
	return &orRule{compoundRule{op: "OR", rules: rules, typeName: "or"}}
}

// Not 创建逻辑非规则
//
// 子规则必须不匹配
//
// 示例：
//
//	rule := abac.Not(abac.SubjectHasRole("banned"))
func Not(rule Rule) Rule {
	return &notRule{rule: rule}
}

// andRule 实现逻辑与
type andRule struct {
	compoundRule
}

// Evaluate 实现 Rule 接口
func (r *andRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	for _, rule := range r.rules {
		match, err := rule.Evaluate(ctx, req)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

// String 返回规则描述
func (r *andRule) String() string {
	parts := make([]string, len(r.rules))
	for i, rule := range r.rules {
		parts[i] = rule.String()
	}
	return fmt.Sprintf("AND(%s)", strings.Join(parts, ", "))
}

// Type 返回规则类型
func (r *andRule) Type() string {
	return "and"
}

// orRule 实现逻辑或
type orRule struct {
	compoundRule
}

// Evaluate 实现 Rule 接口
func (r *orRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	for _, rule := range r.rules {
		match, err := rule.Evaluate(ctx, req)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

// String 返回规则描述
func (r *orRule) String() string {
	parts := make([]string, len(r.rules))
	for i, rule := range r.rules {
		parts[i] = rule.String()
	}
	return fmt.Sprintf("OR(%s)", strings.Join(parts, ", "))
}

// Type 返回规则类型
func (r *orRule) Type() string {
	return "or"
}

// notRule 实现逻辑非
type notRule struct {
	rule Rule
}

// Evaluate 实现 Rule 接口
func (r *notRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	match, err := r.rule.Evaluate(ctx, req)
	if err != nil {
		return false, err
	}
	return !match, nil
}

// String 返回规则描述
func (r *notRule) String() string {
	return fmt.Sprintf("NOT(%s)", r.rule.String())
}

// Type 返回规则类型
func (r *notRule) Type() string {
	return "not"
}

// ===== 预定义的规则构建函数 =====

// SubjectHasRole 创建检查主体角色的规则
//
// 示例：
//
//	rule := abac.SubjectHasRole("admin")
func SubjectHasRole(role string) Rule {
	return &subjectRoleRule{role: role}
}

// subjectRoleRule 检查主体角色
type subjectRoleRule struct {
	role string
}

func (r *subjectRoleRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return req.Subject.HasRole(r.role), nil
}

func (r *subjectRoleRule) String() string {
	return fmt.Sprintf("SubjectHasRole(%s)", r.role)
}

func (r *subjectRoleRule) Type() string {
	return "subject_role"
}

// SubjectHasAnyRole 创建检查主体拥有任意角色的规则
//
// 示例：
//
//	rule := abac.SubjectHasAnyRole("admin", "manager")
func SubjectHasAnyRole(roles ...string) Rule {
	return &subjectAnyRoleRule{roles: roles}
}

// subjectAnyRoleRule 检查主体是否有任意指定角色
type subjectAnyRoleRule struct {
	roles []string
}

func (r *subjectAnyRoleRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return req.Subject.HasAnyRole(r.roles...), nil
}

func (r *subjectAnyRoleRule) String() string {
	return fmt.Sprintf("SubjectHasAnyRole(%s)", strings.Join(r.roles, ", "))
}

func (r *subjectAnyRoleRule) Type() string {
	return "subject_any_role"
}

// SubjectHasAllRoles 创建检查主体拥有所有角色的规则
//
// 示例：
//
//	rule := abac.SubjectHasAllRoles("user", "verified")
func SubjectHasAllRoles(roles ...string) Rule {
	return &subjectAllRolesRule{roles: roles}
}

// subjectAllRolesRule 检查主体是否有所有指定角色
type subjectAllRolesRule struct {
	roles []string
}

func (r *subjectAllRolesRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return req.Subject.HasAllRoles(r.roles...), nil
}

func (r *subjectAllRolesRule) String() string {
	return fmt.Sprintf("SubjectHasAllRoles(%s)", strings.Join(r.roles, ", "))
}

func (r *subjectAllRolesRule) Type() string {
	return "subject_all_roles"
}

// SubjectAttributeEquals 创建检查主体属性等于指定值的规则
//
// 示例：
//
//	rule := abac.SubjectAttributeEquals("department", "engineering")
func SubjectAttributeEquals(key string, value interface{}) Rule {
	return &attributeEqualsRule{
		accessor: func(req Request) AttributeAccessor { return req.Subject },
		key:      key,
		value:    value,
		typeName: "Subject",
	}
}

// SubjectDepartmentIs 创建检查主体部门的规则
//
// 示例：
//
//	rule := abac.SubjectDepartmentIs("engineering")
func SubjectDepartmentIs(department string) Rule {
	return SubjectAttributeEquals("department", department)
}

// SubjectClearanceLevelGte 创建检查主体安全级别大于等于指定值的规则
//
// 示例：
//
//	rule := abac.SubjectClearanceLevelGte(5)
func SubjectClearanceLevelGte(level int) Rule {
	return &attributeGteRule{
		accessor: func(req Request) AttributeAccessor { return req.Subject },
		key:      "clearance_level",
		value:    level,
		typeName: "Subject",
	}
}

// ResourceTypeIs 创建检查资源类型的规则
//
// 示例：
//
//	rule := abac.ResourceTypeIs("document")
func ResourceTypeIs(resourceType string) Rule {
	return &resourceTypeRule{resourceType: resourceType}
}

// resourceTypeRule 检查资源类型
type resourceTypeRule struct {
	resourceType string
}

func (r *resourceTypeRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return strings.EqualFold(req.Resource.Type, r.resourceType), nil
}

func (r *resourceTypeRule) String() string {
	return fmt.Sprintf("ResourceTypeIs(%s)", r.resourceType)
}

func (r *resourceTypeRule) Type() string {
	return "resource_type"
}

// ResourceAttributeEquals 创建检查资源属性等于指定值的规则
//
// 示例：
//
//	rule := abac.ResourceAttributeEquals("sensitivity", "confidential")
func ResourceAttributeEquals(key string, value interface{}) Rule {
	return &attributeEqualsRule{
		accessor: func(req Request) AttributeAccessor { return req.Resource },
		key:      key,
		value:    value,
		typeName: "Resource",
	}
}

// ResourceSensitivityLevelLte 创建检查资源敏感度小于等于指定值的规则
//
// 示例：
//
//	rule := abac.ResourceSensitivityLevelLte(3)
func ResourceSensitivityLevelLte(level int) Rule {
	return &attributeLteRule{
		accessor: func(req Request) AttributeAccessor { return req.Resource },
		key:      "sensitivity_level",
		value:    level,
		typeName: "Resource",
	}
}

// SubjectIsOwner 创建检查主体是否是资源所有者的规则
//
// 示例：
//
//	rule := abac.SubjectIsOwner()
func SubjectIsOwner() Rule {
	return &subjectOwnerRule{}
}

// subjectOwnerRule 检查主体是否是资源所有者
type subjectOwnerRule struct{}

func (r *subjectOwnerRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return req.Resource.IsOwnedBy(req.Subject.ID), nil
}

func (r *subjectOwnerRule) String() string {
	return "SubjectIsOwner()"
}

func (r *subjectOwnerRule) Type() string {
	return "subject_owner"
}

// ActionIs 创建检查操作名称的规则
//
// 示例：
//
//	rule := abac.ActionIs("edit")
func ActionIs(action string) Rule {
	return &actionIsRule{action: action}
}

// actionIsRule 检查操作名称
type actionIsRule struct {
	action string
}

func (r *actionIsRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return req.Action.Equals(r.action), nil
}

func (r *actionIsRule) String() string {
	return fmt.Sprintf("ActionIs(%s)", r.action)
}

func (r *actionIsRule) Type() string {
	return "action_is"
}

// ActionIn 创建检查操作是否在指定列表中的规则
//
// 示例：
//
//	rule := abac.ActionIn("create", "update", "delete")
func ActionIn(actions ...string) Rule {
	return &actionInRule{actions: actions}
}

// actionInRule 检查操作是否在列表中
type actionInRule struct {
	actions []string
}

func (r *actionInRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	for _, action := range r.actions {
		if req.Action.Equals(action) {
			return true, nil
		}
	}
	return false, nil
}

func (r *actionInRule) String() string {
	return fmt.Sprintf("ActionIn(%s)", strings.Join(r.actions, ", "))
}

func (r *actionInRule) Type() string {
	return "action_in"
}

// ActionIsRead 创建检查是否为读操作的规则
//
// 示例：
//
//	rule := abac.ActionIsRead()
func ActionIsRead() Rule {
	return &actionReadRule{}
}

// actionReadRule 检查是否为读操作
type actionReadRule struct{}

func (r *actionReadRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return req.Action.IsRead(), nil
}

func (r *actionReadRule) String() string {
	return "ActionIsRead()"
}

func (r *actionReadRule) Type() string {
	return "action_read"
}

// ActionIsWrite 创建检查是否为写操作的规则
//
// 示例：
//
//	rule := abac.ActionIsWrite()
func ActionIsWrite() Rule {
	return &actionWriteRule{}
}

// actionWriteRule 检查是否为写操作
type actionWriteRule struct{}

func (r *actionWriteRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return req.Action.IsWrite(), nil
}

func (r *actionWriteRule) String() string {
	return "ActionIsWrite()"
}

func (r *actionWriteRule) Type() string {
	return "action_write"
}

// EnvironmentAttributeEquals 创建检查环境属性等于指定值的规则
//
// 示例：
//
//	rule := abac.EnvironmentAttributeEquals("location", "office")
func EnvironmentAttributeEquals(key string, value interface{}) Rule {
	return &attributeEqualsRule{
		accessor: func(req Request) AttributeAccessor { return req.Environment },
		key:      key,
		value:    value,
		typeName: "Environment",
	}
}

// EnvironmentIsInternalNetwork 创建检查是否来自内部网络的规则
//
// 示例：
//
//	rule := abac.EnvironmentIsInternalNetwork()
func EnvironmentIsInternalNetwork() Rule {
	return &environmentInternalRule{}
}

// environmentInternalRule 检查是否来自内部网络
type environmentInternalRule struct{}

func (r *environmentInternalRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return req.Environment.IsInternalNetwork(), nil
}

func (r *environmentInternalRule) String() string {
	return "EnvironmentIsInternalNetwork()"
}

func (r *environmentInternalRule) Type() string {
	return "env_internal"
}

// attributeEqualsRule 通用属性相等规则
type attributeEqualsRule struct {
	accessor func(Request) AttributeAccessor
	key      string
	value    interface{}
	typeName string
}

func (r *attributeEqualsRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	accessor := r.accessor(req)
	actual, exists := accessor.GetAttribute(r.key)
	if !exists {
		return false, nil
	}
	return CompareValues(actual, r.value), nil
}

func (r *attributeEqualsRule) String() string {
	return fmt.Sprintf("%sAttributeEquals(%s, %v)", r.typeName, r.key, r.value)
}

func (r *attributeEqualsRule) Type() string {
	return "attr_equals"
}

// attributeGteRule 通用属性大于等于规则
type attributeGteRule struct {
	accessor func(Request) AttributeAccessor
	key      string
	value    interface{}
	typeName string
}

func (r *attributeGteRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	accessor := r.accessor(req)
	actual, exists := accessor.GetAttribute(r.key)
	if !exists {
		return false, nil
	}
	result, err := GreaterThan(actual, r.value)
	if err != nil {
		return false, err
	}
	// Check for equality as well
	if !result {
		return CompareValues(actual, r.value), nil
	}
	return true, nil
}

func (r *attributeGteRule) String() string {
	return fmt.Sprintf("%sAttributeGte(%s, %v)", r.typeName, r.key, r.value)
}

func (r *attributeGteRule) Type() string {
	return "attr_gte"
}

// attributeLteRule 通用属性小于等于规则
type attributeLteRule struct {
	accessor func(Request) AttributeAccessor
	key      string
	value    interface{}
	typeName string
}

func (r *attributeLteRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	accessor := r.accessor(req)
	actual, exists := accessor.GetAttribute(r.key)
	if !exists {
		return false, nil
	}
	// Check if actual < r.value (less than)
	result, err := LessThan(actual, r.value)
	if err != nil {
		return false, err
	}
	// If less than, it satisfies <= condition
	if result {
		return true, nil
	}
	// Otherwise, check for equality
	return CompareValues(actual, r.value), nil
}

func (r *attributeLteRule) String() string {
	return fmt.Sprintf("%sAttributeLte(%s, %v)", r.typeName, r.key, r.value)
}

func (r *attributeLteRule) Type() string {
	return "attr_lte"
}

// AlwaysAllow 创建一个总是允许的规则
func AlwaysAllow() Rule {
	return &alwaysRule{result: true}
}

// AlwaysDeny 创建一个总是拒绝的规则
func AlwaysDeny() Rule {
	return &alwaysRule{result: false}
}

// alwaysRule 总是返回固定结果的规则
type alwaysRule struct {
	result bool
}

func (r *alwaysRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return r.result, nil
}

func (r *alwaysRule) String() string {
	if r.result {
		return "AlwaysAllow()"
	}
	return "AlwaysDeny()"
}

func (r *alwaysRule) Type() string {
	return "always"
}
