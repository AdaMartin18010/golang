package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIntAdder_Add 测试 IntAdder 的 Add 方法
func TestIntAdder_Add(t *testing.T) {
	tests := []struct {
		name     string
		a        IntAdder
		b        IntAdder
		expected int
	}{
		{
			name:     "add positive numbers",
			a:        IntAdder{value: 5},
			b:        IntAdder{value: 3},
			expected: 8,
		},
		{
			name:     "add negative numbers",
			a:        IntAdder{value: -5},
			b:        IntAdder{value: -3},
			expected: -8,
		},
		{
			name:     "add mixed numbers",
			a:        IntAdder{value: 5},
			b:        IntAdder{value: -3},
			expected: 2,
		},
		{
			name:     "add zero",
			a:        IntAdder{value: 5},
			b:        IntAdder{value: 0},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Add(tt.b)
			assert.Equal(t, tt.expected, result.Value())
		})
	}
}

// TestIntAdder_Value 测试 IntAdder 的 Value 方法
func TestIntAdder_Value(t *testing.T) {
	tests := []struct {
		name     string
		adder    IntAdder
		expected int
	}{
		{
			name:     "positive value",
			adder:    IntAdder{value: 100},
			expected: 100,
		},
		{
			name:     "negative value",
			adder:    IntAdder{value: -50},
			expected: -50,
		},
		{
			name:     "zero value",
			adder:    IntAdder{value: 0},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.adder.Value())
		})
	}
}

// TestSpecificationFactory_Create 测试 SpecificationFactory 的 Create 方法
func TestSpecificationFactory_Create(t *testing.T) {
	factory := SpecificationFactory[TestEntity]{}

	tests := []struct {
		name       string
		predicate  func(*TestEntity) bool
		entity     *TestEntity
		shouldPass bool
	}{
		{
			name:       "predicate returns true",
			predicate:  func(e *TestEntity) bool { return e.Age >= 18 },
			entity:     &TestEntity{ID: "1", Name: "Adult", Age: 25},
			shouldPass: true,
		},
		{
			name:       "predicate returns false",
			predicate:  func(e *TestEntity) bool { return e.Age >= 18 },
			entity:     &TestEntity{ID: "2", Name: "Child", Age: 10},
			shouldPass: false,
		},
		{
			name:       "name predicate true",
			predicate:  func(e *TestEntity) bool { return e.Name == "John" },
			entity:     &TestEntity{ID: "3", Name: "John", Age: 30},
			shouldPass: true,
		},
		{
			name:       "name predicate false",
			predicate:  func(e *TestEntity) bool { return e.Name == "John" },
			entity:     &TestEntity{ID: "4", Name: "Jane", Age: 30},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := factory.Create(tt.predicate)
			assert.NotNil(t, spec)
			assert.Equal(t, tt.shouldPass, spec.IsSatisfiedBy(tt.entity))
		})
	}
}

// TestPredicateSpec_And 测试 predicateSpec 的 And 方法
func TestPredicateSpec_And(t *testing.T) {
	factory := SpecificationFactory[TestEntity]{}

	ageSpec := factory.Create(func(e *TestEntity) bool { return e.Age >= 18 })
	nameSpec := factory.Create(func(e *TestEntity) bool { return e.Name == "John" })

	combined := ageSpec.And(nameSpec)

	tests := []struct {
		name       string
		entity     *TestEntity
		shouldPass bool
	}{
		{
			name:       "both conditions met",
			entity:     &TestEntity{ID: "1", Name: "John", Age: 25},
			shouldPass: true,
		},
		{
			name:       "only age condition met",
			entity:     &TestEntity{ID: "2", Name: "Jane", Age: 25},
			shouldPass: false,
		},
		{
			name:       "only name condition met",
			entity:     &TestEntity{ID: "3", Name: "John", Age: 16},
			shouldPass: false,
		},
		{
			name:       "neither condition met",
			entity:     &TestEntity{ID: "4", Name: "Jane", Age: 16},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.shouldPass, combined.IsSatisfiedBy(tt.entity))
		})
	}
}

// TestPredicateSpec_Or 测试 predicateSpec 的 Or 方法
func TestPredicateSpec_Or(t *testing.T) {
	factory := SpecificationFactory[TestEntity]{}

	ageSpec := factory.Create(func(e *TestEntity) bool { return e.Age >= 18 })
	nameSpec := factory.Create(func(e *TestEntity) bool { return e.Name == "John" })

	combined := ageSpec.Or(nameSpec)

	tests := []struct {
		name       string
		entity     *TestEntity
		shouldPass bool
	}{
		{
			name:       "both conditions met",
			entity:     &TestEntity{ID: "1", Name: "John", Age: 25},
			shouldPass: true,
		},
		{
			name:       "only age condition met",
			entity:     &TestEntity{ID: "2", Name: "Jane", Age: 25},
			shouldPass: true,
		},
		{
			name:       "only name condition met",
			entity:     &TestEntity{ID: "3", Name: "John", Age: 16},
			shouldPass: true,
		},
		{
			name:       "neither condition met",
			entity:     &TestEntity{ID: "4", Name: "Jane", Age: 16},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.shouldPass, combined.IsSatisfiedBy(tt.entity))
		})
	}
}

// TestPredicateSpec_Not 测试 predicateSpec 的 Not 方法
func TestPredicateSpec_Not(t *testing.T) {
	factory := SpecificationFactory[TestEntity]{}

	ageSpec := factory.Create(func(e *TestEntity) bool { return e.Age >= 18 })
	notAgeSpec := ageSpec.Not()

	tests := []struct {
		name       string
		entity     *TestEntity
		shouldPass bool
	}{
		{
			name:       "original true, not false",
			entity:     &TestEntity{ID: "1", Name: "Adult", Age: 25},
			shouldPass: false,
		},
		{
			name:       "original false, not true",
			entity:     &TestEntity{ID: "2", Name: "Child", Age: 10},
			shouldPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.shouldPass, notAgeSpec.IsSatisfiedBy(tt.entity))
		})
	}
}

// TestNewGo126Features 测试 NewGo126Features 构造函数
func TestNewGo126Features(t *testing.T) {
	features := NewGo126Features()

	assert.NotNil(t, features)
	assert.NotNil(t, features.FeaturesUsed)
	assert.Len(t, features.FeaturesUsed, 3)
	assert.Contains(t, features.FeaturesUsed, "self-referential-generics")
	assert.Contains(t, features.FeaturesUsed, "improved-type-inference")
	assert.Contains(t, features.FeaturesUsed, "generic-type-aliases")
}

// TestGo126Features_GetFeatureDocs 测试 GetFeatureDocs 方法
func TestGo126Features_GetFeatureDocs(t *testing.T) {
	features := NewGo126Features()
	docs := features.GetFeatureDocs()

	assert.NotNil(t, docs)
	assert.Len(t, docs, 4)
	assert.Contains(t, docs, "self-referential-generics")
	assert.Contains(t, docs, "errors-as-type")
	assert.Contains(t, docs, "slog-multi-handler")
	assert.Contains(t, docs, "new-expressions")

	// 验证 URL 格式
	for _, url := range docs {
		assert.NotEmpty(t, url)
		assert.Contains(t, url, "https://go.dev/doc/go1.26")
	}
}

// TestGo126Features_MultipleInstances 测试创建多个实例的独立性
func TestGo126Features_MultipleInstances(t *testing.T) {
	f1 := NewGo126Features()
	f2 := NewGo126Features()

	// 验证它们是独立的实例
	assert.NotSame(t, f1, f2)
	// 验证底层数组不同（通过修改f1不影响f2来验证）
	assert.Len(t, f1.FeaturesUsed, 3)
	assert.Len(t, f2.FeaturesUsed, 3)

	// 修改 f1 不影响 f2
	f1.FeaturesUsed = append(f1.FeaturesUsed, "extra-feature")
	assert.Len(t, f1.FeaturesUsed, 4)
	assert.Len(t, f2.FeaturesUsed, 3)
}

// BenchmarkSpecificationFactory_Create 基准测试：SpecificationFactory
func BenchmarkSpecificationFactory_Create(b *testing.B) {
	factory := SpecificationFactory[TestEntity]{}
	predicate := func(e *TestEntity) bool { return e.Age >= 18 }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = factory.Create(predicate)
	}
}

// BenchmarkIntAdder_Add 基准测试：IntAdder
func BenchmarkIntAdder_Add(b *testing.B) {
	a := IntAdder{value: 100}
	c := IntAdder{value: 50}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.Add(c)
	}
}
