package abac

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConditionFunc(t *testing.T) {
	ctx := context.Background()
	req := Request{}

	// 测试返回 true 的条件函数
	trueFunc := ConditionFunc(func(ctx context.Context, req Request) (bool, error) {
		return true, nil
	})

	result, err := trueFunc.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	// 测试返回 false 的条件函数
	falseFunc := ConditionFunc(func(ctx context.Context, req Request) (bool, error) {
		return false, nil
	})

	result, err = falseFunc.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestConditionFunc_String(t *testing.T) {
	fn := ConditionFunc(func(ctx context.Context, req Request) (bool, error) {
		return true, nil
	})
	assert.Equal(t, "custom condition function", fn.String())
}

func TestEqualsCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			ID:    "user-123",
			Roles: []string{"admin"},
			Attributes: map[string]interface{}{
				"level": 5,
			},
		},
	}

	tests := []struct {
		name     string
		attr     string
		value    interface{}
		expected bool
	}{
		{"subject.id equals", "subject.id", "user-123", true},
		{"subject.id not equals", "subject.id", "user-456", false},
		{"subject.attributes.level equals", "subject.attributes.level", 5, true},
		{"nonexistent attribute", "subject.attributes.nonexistent", "value", false},
		{"invalid path", "invalid.path.here", "value", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Eq(tt.attr, tt.value)
			result, err := cond.Evaluate(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotEqualsCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{ID: "user-123"},
	}

	cond := Ne("subject.id", "user-456")
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = Ne("subject.id", "user-123")
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestGreaterThanCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			Attributes: map[string]interface{}{
				"level": 5,
			},
		},
	}

	tests := []struct {
		name     string
		attr     string
		value    interface{}
		expected bool
	}{
		{"5 > 3", "subject.attributes.level", 3, true},
		{"5 > 5", "subject.attributes.level", 5, false},
		{"5 > 10", "subject.attributes.level", 10, false},
		{"nonexistent", "subject.attributes.nonexistent", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Gt(tt.attr, tt.value)
			result, err := cond.Evaluate(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGreaterThanOrEqualCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			Attributes: map[string]interface{}{
				"level": 5,
			},
		},
	}

	cond := Gte("subject.attributes.level", 5)
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = Gte("subject.attributes.level", 6)
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestLessThanCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			Attributes: map[string]interface{}{
				"level": 5,
			},
		},
	}

	cond := Lt("subject.attributes.level", 10)
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = Lt("subject.attributes.level", 5)
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestLessThanOrEqualCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			Attributes: map[string]interface{}{
				"level": 5,
			},
		},
	}

	cond := Lte("subject.attributes.level", 5)
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = Lte("subject.attributes.level", 4)
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestInCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			ID: "user-123",
			Attributes: map[string]interface{}{
				"level": 5,
			},
		},
	}

	tests := []struct {
		name     string
		attr     string
		values   interface{}
		expected bool
	}{
		{"id in slice", "subject.id", []string{"user-123", "user-456"}, true},
		{"id not in slice", "subject.id", []string{"user-789", "user-456"}, false},
		{"level in ints", "subject.attributes.level", []int{1, 3, 5, 7}, true},
		{"nonexistent", "subject.attributes.nonexistent", []string{"a", "b"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := In(tt.attr, tt.values)
			result, err := cond.Evaluate(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotInCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{ID: "user-123"},
	}

	cond := NotIn("subject.id", []string{"user-456", "user-789"})
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = NotIn("subject.id", []string{"user-123", "user-456"})
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestContainsCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			Attributes: map[string]interface{}{
				"tags": []string{"admin", "active", "verified"},
				"bio":  "Software engineer",
			},
		},
	}

	tests := []struct {
		name     string
		attr     string
		value    interface{}
		expected bool
	}{
		{"string in slice", "subject.attributes.tags", "admin", true},
		{"string not in slice", "subject.attributes.tags", "inactive", false},
		{"substring", "subject.attributes.bio", "engineer", true},
		{"substring not found", "subject.attributes.bio", "doctor", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Contains(tt.attr, tt.value)
			result, err := cond.Evaluate(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchesCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			ID: "user-123",
			Attributes: map[string]interface{}{
				"email": "test@example.com",
			},
		},
	}

	tests := []struct {
		name    string
		attr    string
		pattern string
		want    bool
		wantErr bool
	}{
		{"email pattern", "subject.attributes.email", `^[a-z]+@example\.com$`, true, false},
		{"id numeric", "subject.id", `^user-\d+$`, true, false},
		{"no match", "subject.id", `^admin-`, false, false},
		{"invalid regex", "subject.id", `[invalid`, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Matches(tt.attr, tt.pattern)
			result, err := cond.Evaluate(ctx, req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestStartsWithCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{ID: "user-123"},
	}

	cond := StartsWith("subject.id", "user-")
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = StartsWith("subject.id", "admin-")
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestEndsWithCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{ID: "user-123"},
	}

	cond := EndsWith("subject.id", "-123")
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = EndsWith("subject.id", "-456")
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestAllOfCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			ID:    "user-123",
			Roles: []string{"admin"},
		},
	}

	tests := []struct {
		name       string
		conditions []Condition
		expected   bool
	}{
		{
			"all true",
			[]Condition{Eq("subject.id", "user-123"), Eq("subject.id", "user-123")},
			true,
		},
		{
			"one false",
			[]Condition{Eq("subject.id", "user-123"), Eq("subject.id", "user-456")},
			false,
		},
		{
			"empty",
			[]Condition{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := AllOf(tt.conditions...)
			result, err := cond.Evaluate(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnyOfCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{ID: "user-123"},
	}

	tests := []struct {
		name       string
		conditions []Condition
		expected   bool
	}{
		{
			"one true",
			[]Condition{Eq("subject.id", "user-123"), Eq("subject.id", "user-456")},
			true,
		},
		{
			"all false",
			[]Condition{Eq("subject.id", "user-456"), Eq("subject.id", "user-789")},
			false,
		},
		{
			"empty",
			[]Condition{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := AnyOf(tt.conditions...)
			result, err := cond.Evaluate(ctx, req)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{ID: "user-123"},
	}

	cond := NotCondition(Eq("subject.id", "user-456"))
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = NotCondition(Eq("subject.id", "user-123"))
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestBetweenCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			Attributes: map[string]interface{}{
				"age": 25,
			},
		},
	}

	cond := Between("subject.attributes.age", 18, 65)
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = Between("subject.attributes.age", 30, 65)
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestExistsCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			ID: "user-123",
			Attributes: map[string]interface{}{
				"level": 5,
			},
		},
	}

	cond := Exists("subject.id")
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	cond = Exists("subject.attributes.nonexistent")
	result, err = cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestEmptyCondition(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		req      Request
		attr     string
		expected bool
	}{
		{
			"nil attribute",
			Request{Subject: Subject{ID: ""}},
			"subject.id",
			true,
		},
		{
			"empty string",
			Request{Subject: Subject{ID: ""}},
			"subject.id",
			true,
		},
		{
			"non-empty",
			Request{Subject: Subject{ID: "user-123"}},
			"subject.id",
			false,
		},
		{
			"nonexistent",
			Request{Subject: Subject{}},
			"subject.attributes.nonexistent",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Empty(tt.attr)
			result, err := cond.Evaluate(ctx, tt.req)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRuleFromCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{ID: "user-123"},
	}

	cond := Eq("subject.id", "user-123")
	rule := RuleFromCondition(cond)

	result, err := rule.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)

	// 测试 String 方法
	assert.Equal(t, cond.String(), rule.String())
}

func TestOrCondition(t *testing.T) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{ID: "user-123"},
	}

	cond := OrCondition(
		Eq("subject.id", "user-456"),
		Eq("subject.id", "user-123"),
	)
	result, err := cond.Evaluate(ctx, req)
	require.NoError(t, err)
	assert.True(t, result)
}

func TestConditionString(t *testing.T) {
	assert.Equal(t, "subject.id == user-123", Eq("subject.id", "user-123").String())
	assert.Equal(t, "subject.level > 5", Gt("subject.level", 5).String())
	assert.Equal(t, "subject.level < 10", Lt("subject.level", 10).String())
	assert.Equal(t, "subject.id IN [user-123 user-456]", In("subject.id", []string{"user-123", "user-456"}).String())
	assert.Contains(t, Contains("subject.tags", "admin").String(), "CONTAINS")
	assert.Contains(t, Matches("subject.email", `.*@example.com`).String(), "MATCHES")
	assert.Contains(t, StartsWith("subject.id", "user").String(), "STARTSWITH")
	assert.Contains(t, EndsWith("subject.id", "123").String(), "ENDSWITH")
	assert.Contains(t, AllOf(Eq("a", "b")).String(), "ALLOF")
	assert.Contains(t, AnyOf(Eq("a", "b")).String(), "ANYOF")
	assert.Contains(t, NotCondition(Eq("a", "b")).String(), "NOT")
	assert.Contains(t, Exists("subject.id").String(), "EXISTS")
	assert.Contains(t, Empty("subject.id").String(), "EMPTY")
}

func BenchmarkCondition_Evaluate(b *testing.B) {
	ctx := context.Background()
	req := Request{
		Subject: Subject{
			ID:    "user-123",
			Roles: []string{"admin", "user"},
			Attributes: map[string]interface{}{
				"level": 5,
			},
		},
	}

	cond := AllOf(
		Eq("subject.id", "user-123"),
		Gt("subject.attributes.level", 3),
		In("subject.roles", []string{"admin", "moderator"}),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cond.Evaluate(ctx, req)
	}
}
