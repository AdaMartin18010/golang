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
}

// 测试规约实现
type NameSpecification struct {
	Name string
}

func (s NameSpecification) IsSatisfiedBy(entity *TestEntity) bool {
	return entity.Name == s.Name
}

type AgeRangeSpecification struct {
	MinAge int
	MaxAge int
}

func (s AgeRangeSpecification) IsSatisfiedBy(entity *TestEntity) bool {
	return entity.Age >= s.MinAge && entity.Age <= s.MaxAge
}

type EmailDomainSpecification struct {
	Domain string
}

func (s EmailDomainSpecification) IsSatisfiedBy(entity *TestEntity) bool {
	// 简化实现
	return len(entity.Email) > 0
}

// 测试基本规约
func TestSpecification_Basic(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
	}

	// 名称规约
	nameSpec := NameSpecification{Name: "John Doe"}
	assert.True(t, nameSpec.IsSatisfiedBy(entity))

	nameSpec2 := NameSpecification{Name: "Jane Doe"}
	assert.False(t, nameSpec2.IsSatisfiedBy(entity))
}

// 测试 And 规约
func TestAndSpecification(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
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

// 测试 Or 规约
func TestOrSpecification(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
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

// 测试 Not 规约
func TestNotSpecification(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
	}

	nameSpec := NameSpecification{Name: "John Doe"}
	notSpec := Not(nameSpec)

	assert.False(t, notSpec.IsSatisfiedBy(entity), "Not spec should invert the result")

	nameSpec2 := NameSpecification{Name: "Jane Doe"}
	notSpec2 := Not(nameSpec2)
	assert.True(t, notSpec2.IsSatisfiedBy(entity))
}

// 测试复杂组合
func TestSpecification_ComplexComposition(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
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
}

// 测试链式调用
func TestSpecification_Chaining(t *testing.T) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
	}

	spec := &AndSpecification[TestEntity]{
		left:  NameSpecification{Name: "John Doe"},
		right: AgeRangeSpecification{MinAge: 20, MaxAge: 30},
	}

	// 链式调用
	finalSpec := spec.And(EmailDomainSpecification{Domain: "example.com"})
	assert.True(t, finalSpec.IsSatisfiedBy(entity))
}

// 性能测试
func BenchmarkAndSpecification(b *testing.B) {
	entity := &TestEntity{
		ID:    "1",
		Name:  "John Doe",
		Age:   25,
		Email: "john@example.com",
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
