// Package abac provides Attribute-Based Access Control (ABAC) implementation.
//
// evaluator_test.go 包含评估器条件的单元测试
package abac

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestConditionFunc 测试条件函数
func TestConditionFunc(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	cond := ConditionFunc(func(ctx context.Context, req Request) (bool, error) {
		return true, nil
	})

	result, err := cond.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.True(t, result)

	// 测试 String 方法
	str := cond.String()
	assert.Equal(t, "custom condition function", str)
}

// TestEq_String 测试等于条件的 String 方法
func TestEq_String(t *testing.T) {
	cond := Eq("subject.department", "engineering")
	str := cond.String()
	assert.Equal(t, "subject.department == engineering", str)
}

// TestGt_String 测试大于条件的 String 方法
func TestGt_String(t *testing.T) {
	cond := Gt("subject.age", 18)
	str := cond.String()
	assert.Equal(t, "subject.age > 18", str)
}

// TestLt_String 测试小于条件的 String 方法
func TestLt_String(t *testing.T) {
	cond := Lt("subject.age", 65)
	str := cond.String()
	assert.Equal(t, "subject.age < 65", str)
}

// TestIn_String 测试包含于条件的 String 方法
func TestIn_String(t *testing.T) {
	cond := In("subject.role", []string{"admin", "user"})
	str := cond.String()
	assert.Contains(t, str, "IN")
	assert.Contains(t, str, "subject.role")
}

// TestContains_String 测试包含条件的 String 方法
func TestContains_String(t *testing.T) {
	cond := Contains("subject.email", "@example.com")
	str := cond.String()
	assert.Contains(t, str, "CONTAINS")
}

// TestMatches_InvalidRegex 测试无效正则表达式
func TestMatches_InvalidRegex(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	cond := Matches("subject.email", `[`)
	_, err := cond.Evaluate(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid regex pattern")
}

// TestStartsWith_String 测试前缀匹配条件的 String 方法
func TestStartsWith_String(t *testing.T) {
	cond := StartsWith("resource.path", "/api/v1")
	str := cond.String()
	assert.Contains(t, str, "STARTSWITH")
}

// TestEndsWith_String 测试后缀匹配条件的 String 方法
func TestEndsWith_String(t *testing.T) {
	cond := EndsWith("resource.path", ".pdf")
	str := cond.String()
	assert.Contains(t, str, "ENDSWITH")
}

// TestAllOf_String 测试 AllOf 条件的 String 方法
func TestAllOf_String(t *testing.T) {
	cond := AllOf(Eq("a", 1), Eq("b", 2))
	str := cond.String()
	assert.Contains(t, str, "ALLOF")
}

// TestAnyOf_String 测试 AnyOf 条件的 String 方法
func TestAnyOf_String(t *testing.T) {
	cond := AnyOf(Eq("a", 1), Eq("b", 2))
	str := cond.String()
	assert.Contains(t, str, "ANYOF")
}

// TestNotCondition_String 测试非条件的 String 方法
func TestNotCondition_String(t *testing.T) {
	cond := NotCondition(Eq("a", 1))
	str := cond.String()
	assert.Contains(t, str, "NOT")
}

// TestBetween 测试范围条件
func TestBetween(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	cond := Between("subject.clearance_level", 3, 7)
	result, err := cond.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.True(t, result) // clearance_level is 5

	cond = Between("subject.clearance_level", 1, 3)
	result, err = cond.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.False(t, result)
}

// TestExists_String 测试存在条件的 String 方法
func TestExists_String(t *testing.T) {
	cond := Exists("subject.department")
	str := cond.String()
	assert.Contains(t, str, "EXISTS")
}

// TestEmpty_String 测试空条件的 String 方法
func TestEmpty_String(t *testing.T) {
	cond := Empty("resource.owner")
	str := cond.String()
	assert.Contains(t, str, "EMPTY")
}

// TestRuleFromCondition 测试从条件创建规则
func TestRuleFromCondition(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	cond := Eq("subject.department", "engineering")
	rule := RuleFromCondition(cond)

	result, err := rule.Evaluate(ctx, req)
	assert.NoError(t, err)
	assert.True(t, result)

	str := rule.String()
	assert.NotEmpty(t, str)

	typ := rule.Type()
	assert.Equal(t, "condition", typ)
}

// TestTimeBetween_InvalidTime 测试无效时间格式
func TestTimeBetween_InvalidTime(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	cond := TimeBetween("environment.time", "invalid", "18:00")
	_, err := cond.Evaluate(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start time format")

	cond = TimeBetween("environment.time", "09:00", "invalid")
	_, err = cond.Evaluate(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid end time format")
}

// TestTimeBetween_String 测试时间范围条件的 String 方法
func TestTimeBetween_String(t *testing.T) {
	cond := TimeBetween("environment.time", "09:00", "18:00")
	str := cond.String()
	assert.Contains(t, str, "BETWEEN")
}

// TestIPInRange_InvalidCIDR 测试无效 CIDR
func TestIPInRange_InvalidCIDR(t *testing.T) {
	ctx := context.Background()
	req := createTestRequest()

	cond := IPInRange("environment.location", "invalid")
	_, err := cond.Evaluate(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid CIDR format")
}

// TestIPInRange_String 测试 IP 范围条件的 String 方法
func TestIPInRange_String(t *testing.T) {
	cond := IPInRange("environment.location", "192.168.0.0/16")
	str := cond.String()
	assert.Contains(t, str, "IN")
}

// TestWeekdayIn_String 测试星期条件的 String 方法
func TestWeekdayIn_String(t *testing.T) {
	cond := WeekdayIn("environment.time", []time.Weekday{time.Monday, time.Tuesday})
	str := cond.String()
	assert.Contains(t, str, "WEEKDAY IN")
}

// TestAndRule_Type 测试 And 规则的 Type 方法
func TestAndRule_Type(t *testing.T) {
	rule := And(SubjectHasRole("admin"))
	assert.Equal(t, "and", rule.Type())
}

// TestOrRule_Type 测试 Or 规则的 Type 方法
func TestOrRule_Type(t *testing.T) {
	rule := Or(SubjectHasRole("admin"))
	assert.Equal(t, "or", rule.Type())
}

// TestNotRule_Type 测试 Not 规则的 Type 方法
func TestNotRule_Type(t *testing.T) {
	rule := Not(AlwaysAllow())
	assert.Equal(t, "not", rule.Type())
}
