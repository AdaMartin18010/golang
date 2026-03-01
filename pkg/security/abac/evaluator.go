// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// 本文件定义了策略评估器，负责评估复杂的条件表达式。
//
// 评估器支持的条件类型：
//   - 相等/不等比较
//   - 包含/不包含检查
//   - 大于/小于比较
//   - 字符串包含
//   - 正则表达式匹配
//   - 自定义函数
//
// 使用示例：
//
//	// 创建相等条件
//	condition := abac.Eq("subject.department", "engineering")
//
//	// 创建复合条件
//	condition := abac.AllOf(
//	    abac.Eq("subject.clearance_level", 5),
//	    abac.Gt("resource.sensitivity", 3),
//	    abac.In("action.name", []string{"read", "write"}),
//	)
//
//	// 评估条件
//	result, err := condition.Evaluate(ctx, request)

package abac

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Condition 表示评估条件接口
//
// 所有条件类型都实现此接口
type Condition interface {
	// Evaluate 评估条件
	//
	// 参数：
	//   - ctx: 上下文
	//   - req: 访问请求
	//
	// 返回：
	//   - bool: 如果条件满足返回 true
	//   - error: 评估过程中的错误
	Evaluate(ctx context.Context, req Request) (bool, error)
	
	// String 返回条件的可读描述
	String() string
}

// ConditionFunc 是一个函数类型，实现 Condition 接口
type ConditionFunc func(ctx context.Context, req Request) (bool, error)

// Evaluate 实现 Condition 接口
func (f ConditionFunc) Evaluate(ctx context.Context, req Request) (bool, error) {
	return f(ctx, req)
}

// String 返回条件描述
func (f ConditionFunc) String() string {
	return "custom condition function"
}

// ===== 基本比较条件 =====

// Eq 创建等于条件
//
// 检查属性值是否等于指定值
//
// 示例：
//
//	condition := abac.Eq("subject.department", "engineering")
func Eq(attribute string, value interface{}) Condition {
	return &equalsCondition{attribute: attribute, expected: value}
}

// equalsCondition 等于条件
type equalsCondition struct {
	attribute string
	expected  interface{}
}

func (c *equalsCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	actual, err := c.resolveAttribute(req)
	if err != nil {
		return false, nil // 属性不存在视为不匹配
	}
	return CompareValues(actual, c.expected), nil
}

func (c *equalsCondition) String() string {
	return fmt.Sprintf("%s == %v", c.attribute, c.expected)
}

func (c *equalsCondition) resolveAttribute(req Request) (interface{}, error) {
	return resolveRequestAttribute(req, c.attribute)
}

// Ne 创建不等于条件
//
// 检查属性值是否不等于指定值
//
// 示例：
//
//	condition := abac.Ne("subject.status", "banned")
func Ne(attribute string, value interface{}) Condition {
	return NotCondition(Eq(attribute, value))
}

// Gt 创建大于条件
//
// 检查属性值是否大于指定值
//
// 示例：
//
//	condition := abac.Gt("subject.clearance_level", 3)
func Gt(attribute string, value interface{}) Condition {
	return &greaterThanCondition{attribute: attribute, threshold: value}
}

// greaterThanCondition 大于条件
type greaterThanCondition struct {
	attribute string
	threshold interface{}
}

func (c *greaterThanCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	actual, err := resolveRequestAttribute(req, c.attribute)
	if err != nil {
		return false, nil
	}
	result, err := GreaterThan(actual, c.threshold)
	if err != nil {
		return false, fmt.Errorf("cannot compare values: %w", err)
	}
	return result, nil
}

func (c *greaterThanCondition) String() string {
	return fmt.Sprintf("%s > %v", c.attribute, c.threshold)
}

// Gte 创建大于等于条件
//
// 检查属性值是否大于等于指定值
//
// 示例：
//
//	condition := abac.Gte("subject.clearance_level", 3)
func Gte(attribute string, value interface{}) Condition {
	return OrCondition(
		Gt(attribute, value),
		Eq(attribute, value),
	)
}

// Lt 创建小于条件
//
// 检查属性值是否小于指定值
//
// 示例：
//
//	condition := abac.Lt("subject.failed_attempts", 5)
func Lt(attribute string, value interface{}) Condition {
	return &lessThanCondition{attribute: attribute, threshold: value}
}

// lessThanCondition 小于条件
type lessThanCondition struct {
	attribute string
	threshold interface{}
}

func (c *lessThanCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	actual, err := resolveRequestAttribute(req, c.attribute)
	if err != nil {
		return false, nil
	}
	result, err := LessThan(actual, c.threshold)
	if err != nil {
		return false, fmt.Errorf("cannot compare values: %w", err)
	}
	return result, nil
}

func (c *lessThanCondition) String() string {
	return fmt.Sprintf("%s < %v", c.attribute, c.threshold)
}

// Lte 创建小于等于条件
//
// 检查属性值是否小于等于指定值
//
// 示例：
//
//	condition := abac.Lte("subject.failed_attempts", 5)
func Lte(attribute string, value interface{}) Condition {
	return OrCondition(
		Lt(attribute, value),
		Eq(attribute, value),
	)
}

// ===== 集合条件 =====

// In 创建包含于条件
//
// 检查属性值是否在指定集合中
//
// 示例：
//
//	condition := abac.In("action.name", []string{"read", "write"})
func In(attribute string, values interface{}) Condition {
	return &inCondition{attribute: attribute, values: values}
}

// inCondition 包含于条件
type inCondition struct {
	attribute string
	values    interface{}
}

func (c *inCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	actual, err := resolveRequestAttribute(req, c.attribute)
	if err != nil {
		return false, nil
	}
	return ContainsValue(c.values, actual), nil
}

func (c *inCondition) String() string {
	return fmt.Sprintf("%s IN %v", c.attribute, c.values)
}

// NotIn 创建不包含于条件
//
// 检查属性值是否不在指定集合中
//
// 示例：
//
//	condition := abac.NotIn("subject.status", []string{"banned", "suspended"})
func NotIn(attribute string, values interface{}) Condition {
	return NotCondition(In(attribute, values))
}

// Contains 创建包含条件
//
// 检查属性值（字符串或切片）是否包含指定值
//
// 示例：
//
//	condition := abac.Contains("resource.tags", "confidential")
func Contains(attribute string, value interface{}) Condition {
	return &containsCondition{attribute: attribute, searchValue: value}
}

// containsCondition 包含条件
type containsCondition struct {
	attribute   string
	searchValue interface{}
}

func (c *containsCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	actual, err := resolveRequestAttribute(req, c.attribute)
	if err != nil {
		return false, nil
	}

	// 如果是字符串，检查子串
	if str, ok := actual.(string); ok {
		searchStr := fmt.Sprintf("%v", c.searchValue)
		return strings.Contains(str, searchStr), nil
	}

	// 否则检查切片包含
	return ContainsValue(actual, c.searchValue), nil
}

func (c *containsCondition) String() string {
	return fmt.Sprintf("%s CONTAINS %v", c.attribute, c.searchValue)
}

// Matches 创建正则匹配条件
//
// 检查属性值是否匹配正则表达式
//
// 示例：
//
//	condition := abac.Matches("subject.email", `^[a-zA-Z0-9._%+-]+@company\.com$`)
func Matches(attribute string, pattern string) Condition {
	return &matchesCondition{attribute: attribute, pattern: pattern}
}

// matchesCondition 正则匹配条件
type matchesCondition struct {
	attribute string
	pattern   string
	regexp    *regexp.Regexp
}

func (c *matchesCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	// 延迟编译正则表达式
	if c.regexp == nil {
		re, err := regexp.Compile(c.pattern)
		if err != nil {
			return false, fmt.Errorf("invalid regex pattern: %w", err)
		}
		c.regexp = re
	}

	actual, err := resolveRequestAttribute(req, c.attribute)
	if err != nil {
		return false, nil
	}

	str := fmt.Sprintf("%v", actual)
	return c.regexp.MatchString(str), nil
}

func (c *matchesCondition) String() string {
	return fmt.Sprintf("%s MATCHES /%s/", c.attribute, c.pattern)
}

// StartsWith 创建前缀匹配条件
//
// 检查属性值是否以指定前缀开头
//
// 示例：
//
//	condition := abac.StartsWith("resource.path", "/api/v1/")
func StartsWith(attribute string, prefix string) Condition {
	return &startsWithCondition{attribute: attribute, prefix: prefix}
}

// startsWithCondition 前缀匹配条件
type startsWithCondition struct {
	attribute string
	prefix    string
}

func (c *startsWithCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	actual, err := resolveRequestAttribute(req, c.attribute)
	if err != nil {
		return false, nil
	}

	str := fmt.Sprintf("%v", actual)
	return strings.HasPrefix(str, c.prefix), nil
}

func (c *startsWithCondition) String() string {
	return fmt.Sprintf("%s STARTSWITH %q", c.attribute, c.prefix)
}

// EndsWith 创建后缀匹配条件
//
// 检查属性值是否以指定后缀结尾
//
// 示例：
//
//	condition := abac.EndsWith("resource.name", ".pdf")
func EndsWith(attribute string, suffix string) Condition {
	return &endsWithCondition{attribute: attribute, suffix: suffix}
}

// endsWithCondition 后缀匹配条件
type endsWithCondition struct {
	attribute string
	suffix    string
}

func (c *endsWithCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	actual, err := resolveRequestAttribute(req, c.attribute)
	if err != nil {
		return false, nil
	}

	str := fmt.Sprintf("%v", actual)
	return strings.HasSuffix(str, c.suffix), nil
}

func (c *endsWithCondition) String() string {
	return fmt.Sprintf("%s ENDSWITH %q", c.attribute, c.suffix)
}

// ===== 复合条件 =====

// AllOf 创建所有条件都必须满足的条件
//
// 等同于逻辑与
//
// 示例：
//
//	condition := abac.AllOf(
//	    abac.Eq("subject.role", "admin"),
//	    abac.Gt("subject.clearance", 3),
//	)
func AllOf(conditions ...Condition) Condition {
	return &allOfCondition{conditions: conditions}
}

// allOfCondition 所有条件都必须满足
type allOfCondition struct {
	conditions []Condition
}

func (c *allOfCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	for _, cond := range c.conditions {
		result, err := cond.Evaluate(ctx, req)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

func (c *allOfCondition) String() string {
	parts := make([]string, len(c.conditions))
	for i, cond := range c.conditions {
		parts[i] = cond.String()
	}
	return fmt.Sprintf("ALLOF(%s)", strings.Join(parts, " AND "))
}

// AnyOf 创建至少一个条件必须满足的条件
//
// 等同于逻辑或
//
// 示例：
//
//	condition := abac.AnyOf(
//	    abac.Eq("subject.role", "admin"),
//	    abac.Eq("subject.role", "manager"),
//	)
func AnyOf(conditions ...Condition) Condition {
	return &anyOfCondition{conditions: conditions}
}

// anyOfCondition 至少一个条件必须满足
type anyOfCondition struct {
	conditions []Condition
}

func (c *anyOfCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	for _, cond := range c.conditions {
		result, err := cond.Evaluate(ctx, req)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

func (c *anyOfCondition) String() string {
	parts := make([]string, len(c.conditions))
	for i, cond := range c.conditions {
		parts[i] = cond.String()
	}
	return fmt.Sprintf("ANYOF(%s)", strings.Join(parts, " OR "))
}

// NotCondition 创建条件的逻辑非
//
// 示例：
//
//	condition := abac.NotCondition(abac.Eq("subject.status", "banned"))
func NotCondition(condition Condition) Condition {
	return &notCondition{condition: condition}
}

// notCondition 逻辑非条件
type notCondition struct {
	condition Condition
}

func (c *notCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	result, err := c.condition.Evaluate(ctx, req)
	if err != nil {
		return false, err
	}
	return !result, nil
}

func (c *notCondition) String() string {
	return fmt.Sprintf("NOT(%s)", c.condition.String())
}

// OrCondition 创建两个条件的逻辑或
//
// 示例：
//
//	condition := abac.OrCondition(
//	    abac.Eq("subject.role", "admin"),
//	    abac.Eq("subject.role", "manager"),
//	)
func OrCondition(a, b Condition) Condition {
	return AnyOf(a, b)
}

// Between 创建范围条件
//
// 检查属性值是否在指定范围内（包含边界）
//
// 示例：
//
//	condition := abac.Between("subject.age", 18, 65)
func Between(attribute string, min, max interface{}) Condition {
	return AllOf(
		Gte(attribute, min),
		Lte(attribute, max),
	)
}

// Exists 创建属性存在条件
//
// 检查属性是否存在（不为 nil 或空）
//
// 示例：
//
//	condition := abac.Exists("subject.department")
func Exists(attribute string) Condition {
	return &existsCondition{attribute: attribute}
}

// existsCondition 属性存在条件
type existsCondition struct {
	attribute string
}

func (c *existsCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	_, err := resolveRequestAttribute(req, c.attribute)
	return err == nil, nil
}

func (c *existsCondition) String() string {
	return fmt.Sprintf("EXISTS(%s)", c.attribute)
}

// Empty 创建属性为空条件
//
// 检查属性是否为空（nil、空字符串或空切片）
//
// 示例：
//
//	condition := abac.Empty("resource.owner")
func Empty(attribute string) Condition {
	return &emptyCondition{attribute: attribute}
}

// emptyCondition 属性为空条件
type emptyCondition struct {
	attribute string
}

func (c *emptyCondition) Evaluate(ctx context.Context, req Request) (bool, error) {
	value, err := resolveRequestAttribute(req, c.attribute)
	if err != nil {
		return true, nil // 不存在视为空
	}

	if value == nil {
		return true, nil
	}

	str := fmt.Sprintf("%v", value)
	return str == "" || str == "[]" || str == "{}", nil
}

func (c *emptyCondition) String() string {
	return fmt.Sprintf("EMPTY(%s)", c.attribute)
}

// resolveRequestAttribute 从请求中解析属性值
//
// 支持属性路径：
//   - subject.id, subject.department, subject.attributes.xxx
//   - resource.id, resource.type, resource.owner, resource.attributes.xxx
//   - action.name, action.attributes.xxx
//   - environment.time, environment.location, environment.attributes.xxx
func resolveRequestAttribute(req Request, path string) (interface{}, error) {
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

// RuleFromCondition 将 Condition 转换为 Rule
//
// 用于在策略规则中使用复杂条件
func RuleFromCondition(condition Condition) Rule {
	return &conditionRule{condition: condition}
}

// conditionRule 包装 Condition 为 Rule
type conditionRule struct {
	condition Condition
}

func (r *conditionRule) Evaluate(ctx context.Context, req Request) (bool, error) {
	return r.condition.Evaluate(ctx, req)
}

func (r *conditionRule) String() string {
	return r.condition.String()
}
