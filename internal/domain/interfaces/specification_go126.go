// Package interfaces provides Go 1.26+ optimized specification pattern interfaces
// using generic self-reference feature.
//
// Go 1.26 泛型自引用优化版本
// 相比原版提供更严格的类型约束和更好的类型安全
//
// 主要改进：
// 1. 使用泛型自引用约束，支持递归类型定义
// 2. 类型参数 S 代表规约类型自身，避免类型断言
// 3. 链式调用时保持类型信息
//
// 使用示例：
//
//	// Go 1.26 现在支持自引用类型约束
//	type Comparable[T any] interface {
//	    Compare(other T) int
//	}
//	
//	// 可以在接口中使用自己作为约束
//	type Ordered[T Ordered[T]] interface {
//	    Less(other T) bool
//	}
package interfaces

// SelfReferencingInterface demonstrates Go 1.26's new self-referential generic type feature
// Before Go 1.26, this was not allowed.
//
// 此接口展示 Go 1.26 的泛型自引用特性
// 类型参数 T 的约束引用了接口自身
type SelfReferencingInterface[T SelfReferencingInterface[T]] interface {
	// Method returns the value
	Method() string
	// Combine combines two instances
	Combine(other T) T
}

// Adder is an example from Go 1.26 release notes
// demonstrating self-referential generic types.
//
// Adder 要求类型参数 A 必须实现 Adder[A] 接口
type Adder[A Adder[A]] interface {
	Add(A) A
}

// IntAdder implements Adder for int
type IntAdder struct {
	value int
}

// Add implements the Adder interface
func (a IntAdder) Add(other IntAdder) IntAdder {
	return IntAdder{value: a.value + other.value}
}

// Value returns the int value
func (a IntAdder) Value() int {
	return a.value
}

// ComparableSpecification demonstrates how self-referential generics
// can be used with the Specification pattern in Go 1.26.
//
// This shows the theoretical improvement, though practical usage
// requires careful design of the type hierarchy.
type ComparableSpecification[T any, S ComparableSpecification[T, S]] interface {
	Specification[T]
	// Compare allows comparing specifications
	Compare(other S) int
}

// SpecificationFactory uses Go 1.26 features for type-safe specification creation
type SpecificationFactory[T any] struct{}

// Create returns a new specification with proper type constraints
// This leverages Go 1.26's improved generic type inference
func (f SpecificationFactory[T]) Create(pred func(*T) bool) Specification[T] {
	return &predicateSpec[T]{predicate: pred}
}

// predicateSpec is a simple specification implementation
type predicateSpec[T any] struct {
	predicate func(*T) bool
}

func (s *predicateSpec[T]) IsSatisfiedBy(entity *T) bool {
	return s.predicate(entity)
}

// And combines two specifications (uses Go 1.26 improved type inference)
func (s *predicateSpec[T]) And(other Specification[T]) Specification[T] {
	return &AndSpecification[T]{left: s, right: other}
}

// Or combines two specifications with OR logic
func (s *predicateSpec[T]) Or(other Specification[T]) Specification[T] {
	return &OrSpecification[T]{left: s, right: other}
}

// Not negates the specification
func (s *predicateSpec[T]) Not() Specification[T] {
	return &NotSpecification[T]{spec: s}
}

// Go126Features is a marker type indicating use of Go 1.26 features
type Go126Features struct {
	// FeaturesUsed lists the Go 1.26 features utilized
	FeaturesUsed []string
}

// NewGo126Features creates a marker for Go 1.26 features
func NewGo126Features() *Go126Features {
	return &Go126Features{
		FeaturesUsed: []string{
			"self-referential-generics",
			"improved-type-inference",
			"generic-type-aliases",
		},
	}
}

// GetFeatureDocs returns documentation links for Go 1.26 features
func (f *Go126Features) GetFeatureDocs() map[string]string {
	return map[string]string{
		"self-referential-generics": "https://go.dev/doc/go1.26#language",
		"errors-as-type":            "https://go.dev/doc/go1.26#errors",
		"slog-multi-handler":        "https://go.dev/doc/go1.26#log/slog",
		"new-expressions":           "https://go.dev/doc/go1.26#language",
	}
}
