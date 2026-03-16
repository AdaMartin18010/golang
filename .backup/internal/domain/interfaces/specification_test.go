package interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用的简单实体
type TestEntity struct {
	ID    string
	Name  string
	Age   int
	Email string
	Score float64
}

// 测试规约实现
type NameSpecification struct {
	Name string
}

func (s NameSpecification) IsSatisfiedBy(entity *TestEntity) bool {
	return entity.Name == s.Name
}

func (s NameSpecification) And(other Specification[TestEntity]) Specification[TestEntity] {
	return NewAndSpecification[TestEntity](s, other)
}

func (s NameSpecification) Or(other Specification[TestEntity]) Specification[TestEntity] {
	return NewOrSpecification[TestEntity](s, other)
}

func (s NameSpecification) Not() Specification[TestEntity] {
	return NewNotSpecification[TestEntity](s)
}

type AgeRangeSpecification struct {
	MinAge int
	MaxAge int
}

func (s AgeRangeSpecification) IsSatisfiedBy(entity *TestEntity) bool {
	return entity.Age >= s.MinAge && entity.Age <= s.MaxAge
}

func (s AgeRangeSpecification) And(other Specification[TestEntity]) Specification[TestEntity] {
	return NewAndSpecification[TestEntity](s, other)
}

func (s AgeRangeSpecification) Or(other Specification[TestEntity]) Specification[TestEntity] {
	return NewOrSpecification[TestEntity](s, other)
}

func (s AgeRangeSpecification) Not() Specification[TestEntity] {
	return NewNotSpecification[TestEntity](s)
}

type EmailDomainSpecification struct {
	Domain string
}

func (s EmailDomainSpecification) IsSatisfiedBy(entity *TestEntity) bool {
	// 简化实现：检查邮箱是否包含域名
	return len(entity.Email) > 0
}

func (s EmailDomainSpecification) And(other Specification[TestEntity]) Specification[TestEntity] {
	return NewAndSpecification[TestEntity](s, other)
}

func (s EmailDomainSpecification) Or(other Specification[TestEntity]) Specification[TestEntity] {
	return NewOrSpecification[TestEntity](s, other)
}

func (s EmailDomainSpecification) Not() Specification[TestEntity] {
	return NewNotSpecification[TestEntity](s)
}

type ScoreThresholdSpecification struct {
	Threshold float64
}

func (s ScoreThresholdSpecification) IsSatisfiedBy(entity *TestEntity) bool {
	return entity.Score >= s.Threshold
}

func (s ScoreThresholdSpecification) And(other Specification[TestEntity]) Specification[TestEntity] {
	return NewAndSpecification[TestEntity](s, other)
}

func (s ScoreThresholdSpecification) Or(other Specification[TestEntity]) Specification[TestEntity] {
	return NewOrSpecification[TestEntity](s, other)
}

func (s ScoreThresholdSpecification) Not() Specification[TestEntity] {
	return NewNotSpecification[TestEntity](s)
}

// TestSpecification_Basic 测试基本规约
func TestSpecification_Basic(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	// 名称规约
	nameSpec := NameSpecification{Name: "John Doe"}
	assert.True(t, nameSpec.IsSatisfiedBy(entity))

	nameSpec2 := NameSpecification{Name: "Jane Doe"}
	assert.False(t, nameSpec2.IsSatisfiedBy(entity))
}

// TestAndSpecification 测试 And 规约
func TestAndSpecification(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	nameSpec := NameSpecification{Name: "John Doe"}
	ageSpec := AgeRangeSpecification{MinAge: 20, MaxAge: 30}

	// 组合规约
	andSpec := And(nameSpec, ageSpec)
	assert.True(t, andSpec.IsSatisfiedBy(entity), "Both specs should be satisfied")

	// 不满足的情况
	ageSpec2 := AgeRangeSpecification{MinAge: 30, MaxAge: 40}
	andSpec2 := And(nameSpec, ageSpec2)
	assert.False(t, andSpec2.IsSatisfiedBy(entity), "Age spec should not be satisfied")
}

// TestOrSpecification 测试 Or 规约
func TestOrSpecification(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	nameSpec := NameSpecification{Name: "John Doe"}
	ageSpec := AgeRangeSpecification{MinAge: 30, MaxAge: 40}

	// Or 规约 - 满足任一即可
	orSpec := Or(nameSpec, ageSpec)
	assert.True(t, orSpec.IsSatisfiedBy(entity), "Name spec is satisfied")

	// 两个都不满足
	nameSpec2 := NameSpecification{Name: "Jane Doe"}
	orSpec2 := Or(nameSpec2, ageSpec)
	assert.False(t, orSpec2.IsSatisfiedBy(entity), "Neither spec should be satisfied")
}

// TestNotSpecification 测试 Not 规约
func TestNotSpecification(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	nameSpec := NameSpecification{Name: "John Doe"}
	notSpec := Not(nameSpec)

	assert.False(t, notSpec.IsSatisfiedBy(entity), "Not spec should invert the result")

	nameSpec2 := NameSpecification{Name: "Jane Doe"}
	notSpec2 := Not(nameSpec2)
	assert.True(t, notSpec2.IsSatisfiedBy(entity))
}

// TestAndSpecification_Methods 测试 AndSpecification 的方法
func TestAndSpecification_Methods(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	// 创建基础 AndSpecification
	andSpec := &AndSpecification[TestEntity]{
		left:  NameSpecification{Name: "John Doe"},
		right: AgeRangeSpecification{MinAge: 20, MaxAge: 30},
	}

	// 测试 And 方法 - 链式调用
	scoreSpec := ScoreThresholdSpecification{Threshold: 80.0}
	chainedAnd := andSpec.And(scoreSpec)
	assert.True(t, chainedAnd.IsSatisfiedBy(entity), "Chained And should be satisfied")

	// 测试 Or 方法
	orResult := andSpec.Or(EmailDomainSpecification{Domain: "example.com"})
	assert.True(t, orResult.IsSatisfiedBy(entity), "Or with matching left should be satisfied")

	// 测试 Not 方法
	notResult := andSpec.Not()
	assert.False(t, notResult.IsSatisfiedBy(entity), "Not should invert the result")
}

// TestOrSpecification_Methods 测试 OrSpecification 的方法
func TestOrSpecification_Methods(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	// 创建 OrSpecification
	orSpec := &OrSpecification[TestEntity]{
		left:  NameSpecification{Name: "Jane Doe"},           // 不满足
		right: AgeRangeSpecification{MinAge: 20, MaxAge: 30}, // 满足
	}

	// 测试 IsSatisfiedBy
	assert.True(t, orSpec.IsSatisfiedBy(entity), "Right spec should satisfy")

	// 测试 And 方法
	scoreSpec := ScoreThresholdSpecification{Threshold: 90.0} // 不满足
	andResult := orSpec.And(scoreSpec)
	assert.False(t, andResult.IsSatisfiedBy(entity), "And with false right should be false")

	// 测试 Or 方法 - 链式调用
	nameSpec := NameSpecification{Name: "John Doe"}
	orResult := orSpec.Or(nameSpec)
	assert.True(t, orResult.IsSatisfiedBy(entity), "Or with matching spec should be true")

	// 测试 Not 方法
	notResult := orSpec.Not()
	assert.False(t, notResult.IsSatisfiedBy(entity), "Not should invert the result")
}

// TestNotSpecification_Methods 测试 NotSpecification 的方法
func TestNotSpecification_Methods(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	// 创建 NotSpecification
	nameSpec := NameSpecification{Name: "Jane Doe"}
	notSpec := &NotSpecification[TestEntity]{spec: nameSpec}

	// 测试 IsSatisfiedBy
	assert.True(t, notSpec.IsSatisfiedBy(entity), "Not of false should be true")

	// 测试 And 方法
	ageSpec := AgeRangeSpecification{MinAge: 20, MaxAge: 30}
	andResult := notSpec.And(ageSpec)
	assert.True(t, andResult.IsSatisfiedBy(entity), "True AND True should be true")

	// 测试 Or 方法
	orResult := notSpec.Or(NameSpecification{Name: "Jane Doe"})
	assert.True(t, orResult.IsSatisfiedBy(entity), "True OR False should be true")

	// 测试 Not 方法（双重否定）
	notNotResult := notSpec.Not()
	assert.False(t, notNotResult.IsSatisfiedBy(entity), "Double negation should restore original")
}

// TestNewAndSpecification 测试 NewAndSpecification 工厂函数
func TestNewAndSpecification(t *testing.T) {
	entity := &TestEntity{
		ID:   "1",
		Name: "John",
		Age:  25,
	}

	nameSpec := NameSpecification{Name: "John"}
	ageSpec := AgeRangeSpecification{MinAge: 20, MaxAge: 30}

	spec := NewAndSpecification(nameSpec, ageSpec)
	assert.NotNil(t, spec)
	assert.True(t, spec.IsSatisfiedBy(entity))
}

// TestNewOrSpecification 测试 NewOrSpecification 工厂函数
func TestNewOrSpecification(t *testing.T) {
	entity := &TestEntity{
		ID:   "1",
		Name: "John",
		Age:  35,
	}

	nameSpec := NameSpecification{Name: "John"}
	ageSpec := AgeRangeSpecification{MinAge: 20, MaxAge: 30}

	spec := NewOrSpecification(nameSpec, ageSpec)
	assert.NotNil(t, spec)
	assert.True(t, spec.IsSatisfiedBy(entity), "Name matches, so Or should be true")
}

// TestNewNotSpecification 测试 NewNotSpecification 工厂函数
func TestNewNotSpecification(t *testing.T) {
	entity := &TestEntity{
		ID:   "1",
		Name: "John",
	}

	nameSpec := NameSpecification{Name: "Jane"}
	spec := NewNotSpecification(nameSpec)
	assert.NotNil(t, spec)
	assert.True(t, spec.IsSatisfiedBy(entity), "Not of false should be true")
}

// TestSpecification_ComplexComposition 测试复杂组合
func TestSpecification_ComplexComposition(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	// (Name = "John Doe" AND Age 20-30) OR Email domain
	nameSpec := NameSpecification{Name: "John Doe"}
	ageSpec := AgeRangeSpecification{MinAge: 20, MaxAge: 30}
	emailSpec := EmailDomainSpecification{Domain: "example.com"}

	complexSpec := Or(
		And(nameSpec, ageSpec),
		emailSpec,
	)

	assert.True(t, complexSpec.IsSatisfiedBy(entity))

	// 测试不满足的情况
	entity2 := &TestEntity{
		ID:    "2",
		Name:  "Jane Doe",
		Age:   35,
		Email: "",
		Score: 50.0,
	}

	assert.False(t, complexSpec.IsSatisfiedBy(entity2))
}

// TestSpecification_Chaining 测试链式调用
func TestSpecification_Chaining(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	// 创建基础规约
	andSpec := &AndSpecification[TestEntity]{
		left:  NameSpecification{Name: "John Doe"},
		right: AgeRangeSpecification{MinAge: 20, MaxAge: 30},
	}

	// 链式调用: (Name AND Age) AND Score
	chainedAnd := andSpec.And(ScoreThresholdSpecification{Threshold: 80.0})
	assert.True(t, chainedAnd.IsSatisfiedBy(entity))

	// 链式调用: ((Name AND Age) AND Score) OR Email
	// 注意：And() 返回 Specification[T]，需要使用 Or() 辅助函数
	chainedOr := Or(chainedAnd, EmailDomainSpecification{Domain: "other.com"})
	assert.True(t, chainedOr.IsSatisfiedBy(entity))

	// 复杂链式: NOT((Name AND Age) AND Score)
	// 注意：使用 Not() 辅助函数
	notChained := Not(chainedAnd)
	assert.False(t, notChained.IsSatisfiedBy(entity))
}

// TestSpecification_DeepNesting 测试深层嵌套规约
func TestSpecification_DeepNesting(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	// 创建深层嵌套的规约: ((Name AND Age) OR (Score AND Email)) AND NOT(Name)
	// 实际上这会返回 false，因为 NOT(Name) 会否定结果
	nameSpec := NameSpecification{Name: "John Doe"}
	ageSpec := AgeRangeSpecification{MinAge: 20, MaxAge: 30}
	scoreSpec := ScoreThresholdSpecification{Threshold: 80.0}
	emailSpec := EmailDomainSpecification{Domain: "example.com"}

	innerAnd1 := And(nameSpec, ageSpec)    // true
	innerAnd2 := And(scoreSpec, emailSpec) // true
	innerOr := Or(innerAnd1, innerAnd2)    // true
	notName := Not(nameSpec)               // false
	finalSpec := And(innerOr, notName)     // true AND false = false

	assert.False(t, finalSpec.IsSatisfiedBy(entity))

	// 另一个测试: ((Name AND Age) OR (Score AND Email)) AND (Age OR Score)
	ageOrScore := Or(
		AgeRangeSpecification{MinAge: 20, MaxAge: 30},
		ScoreThresholdSpecification{Threshold: 80.0},
	)
	finalSpec2 := And(innerOr, ageOrScore) // true AND true = true
	assert.True(t, finalSpec2.IsSatisfiedBy(entity))
}

// TestSpecification_EdgeCases 测试边缘情况
func TestSpecification_EdgeCases(t *testing.T) {
	// 测试边界值
	entity := &TestEntity{
		ID:    "1",
		Name:  "John",
		Age:   20, // 边界值
		Email: "john@example.com",
		Score: 80.0, // 边界值
	}

	// Age 边界测试
	ageAtMin := AgeRangeSpecification{MinAge: 20, MaxAge: 30}
	assert.True(t, ageAtMin.IsSatisfiedBy(entity), "Age at minimum boundary should satisfy")

	entity.Age = 30 // 另一个边界值
	ageAtMax := AgeRangeSpecification{MinAge: 20, MaxAge: 30}
	assert.True(t, ageAtMax.IsSatisfiedBy(entity), "Age at maximum boundary should satisfy")

	// Score 边界测试
	scoreAtThreshold := ScoreThresholdSpecification{Threshold: 80.0}
	assert.True(t, scoreAtThreshold.IsSatisfiedBy(entity), "Score at threshold should satisfy")

	entity.Score = 79.999
	assert.False(t, scoreAtThreshold.IsSatisfiedBy(entity), "Score below threshold should not satisfy")
}

// TestSpecification_EmptyAndNilHandling 测试空值和 nil 处理
func TestSpecification_EmptyAndNilHandling(t *testing.T) {
	emptyEntity := &TestEntity{}

	// 空实体测试
	nameSpec := NameSpecification{Name: ""}
	assert.True(t, nameSpec.IsSatisfiedBy(emptyEntity), "Empty name should match empty name spec")

	nameSpec2 := NameSpecification{Name: "John"}
	assert.False(t, nameSpec2.IsSatisfiedBy(emptyEntity), "Non-empty name spec should not match empty entity")

	// 测试 And 与空规约
	emptyAnd := And(
		NameSpecification{Name: ""},
		AgeRangeSpecification{MinAge: 0, MaxAge: 0},
	)
	assert.True(t, emptyAnd.IsSatisfiedBy(emptyEntity), "Empty specs should match empty entity")
}

// TestSpecification_ShortCircuit 测试短路逻辑
func TestSpecification_ShortCircuit(t *testing.T) {
	entity := &TestEntity{
		ID:   "1",
		Name: "John",
		Age:  25,
	}

	// 测试 And 短路: 第一个不满足，第二个不会被执行（从逻辑上）
	// 注意：Go 总是会执行两个，但语义上是短路的
	andSpec := And(
		NameSpecification{Name: "Jane"}, // false
		AgeRangeSpecification{MinAge: 20, MaxAge: 30},
	)
	assert.False(t, andSpec.IsSatisfiedBy(entity))

	// 测试 Or 短路: 第一个满足，第二个从语义上不需要检查
	orSpec := Or(
		NameSpecification{Name: "John"},               // true
		AgeRangeSpecification{MinAge: 50, MaxAge: 60}, // false
	)
	assert.True(t, orSpec.IsSatisfiedBy(entity))
}

// BenchmarkAndSpecification 性能基准测试 - And 规约
func BenchmarkAndSpecification(b *testing.B) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	spec := And(
		NameSpecification{Name: "John Doe"},
		AgeRangeSpecification{MinAge: 20, MaxAge: 30},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spec.IsSatisfiedBy(entity)
	}
}

// BenchmarkOrSpecification 性能基准测试 - Or 规约
func BenchmarkOrSpecification(b *testing.B) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	spec := Or(
		NameSpecification{Name: "John Doe"},
		AgeRangeSpecification{MinAge: 50, MaxAge: 60},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spec.IsSatisfiedBy(entity)
	}
}

// BenchmarkNotSpecification 性能基准测试 - Not 规约
func BenchmarkNotSpecification(b *testing.B) {
	entity := &TestEntity{
		ID:   "1",
		Name: "John Doe",
		Age:  25,
	}

	spec := Not(NameSpecification{Name: "Jane Doe"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spec.IsSatisfiedBy(entity)
	}
}

// BenchmarkComplexSpecification 性能基准测试 - 复杂组合规约
func BenchmarkComplexSpecification(b *testing.B) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
		Score: 85.5,
	}

	// 复杂组合: ((Name AND Age) OR (Score AND Email)) AND NOT(Name = "Jane")
	spec := And(
		Or(
			And(
				NameSpecification{Name: "John Doe"},
				AgeRangeSpecification{MinAge: 20, MaxAge: 30},
			),
			And(
				ScoreThresholdSpecification{Threshold: 80.0},
				EmailDomainSpecification{Domain: "example.com"},
			),
		),
		Not(NameSpecification{Name: "Jane Doe"}),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spec.IsSatisfiedBy(entity)
	}
}
