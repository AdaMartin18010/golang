# Generics Comparison: Go vs Other Languages

## Executive Summary

Generics (parametric polymorphism) allow writing reusable code that works with multiple types. This document compares Go's generics (introduced in 1.18) with implementations in Rust, Java, C++, TypeScript, and C#, analyzing syntax, constraints, type inference, and performance characteristics.

---

## Table of Contents

- [Generics Comparison: Go vs Other Languages](#generics-comparison-go-vs-other-languages)
  - [Executive Summary](#executive-summary)
  - [Table of Contents](#table-of-contents)
  - [Go Generics](#go-generics)
    - [Basic Syntax](#basic-syntax)
    - [Advanced Patterns](#advanced-patterns)
    - [Limitations](#limitations)
  - [Rust Generics](#rust-generics)
  - [Java Generics](#java-generics)
  - [C++ Templates](#c-templates)
  - [TypeScript Generics](#typescript-generics)
  - [C# Generics](#c-generics)
  - [Feature Matrix](#feature-matrix)
  - [Performance Comparison](#performance-comparison)
  - [Summary](#summary)

---

## Go Generics

Go 1.18 introduced generics with a focus on simplicity:

### Basic Syntax

```go
package main

import "fmt"

// Generic function
func Min[T comparable](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Generic type
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// Type constraints
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
    ~float32 | ~float64
}

func Sum[T Number](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}

// Using constraints package
import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// Type sets with approximation (~)
type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Custom type with underlying int works
type MyInt int

func main() {
    // Type inference
    m := Min(3, 5)  // T inferred as int

    // Explicit type parameter
    m2 := Min[float64](3.14, 2.71)

    // Generic types
    stack := Stack[string]{}
    stack.Push("hello")
    item, ok := stack.Pop()

    fmt.Println(m, m2, item, ok)
}
```

### Advanced Patterns

```go
// Generic interfaces
type Container[T any] interface {
    Get() T
    Set(T)
}

// Generic methods on non-generic types
type Processor struct{}

func (p Processor) Process[T any](item T) T {
    return item
}

// Multiple type parameters
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Generic channels
func Merge[T any](ch1, ch2 <-chan T) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for ch1 != nil || ch2 != nil {
            select {
            case v, ok := <-ch1:
                if !ok {
                    ch1 = nil
                    continue
                }
                out <- v
            case v, ok := <-ch2:
                if !ok {
                    ch2 = nil
                    continue
                }
                out <- v
            }
        }
    }()
    return out
}

// Generic result type
type Result[T any] struct {
    Value T
    Error error
}

func (r Result[T]) OrDefault(defaultValue T) T {
    if r.Error != nil {
        return defaultValue
    }
    return r.Value
}
```

### Limitations

```go
// Cannot use type parameters in method receivers
type MyType[T any] struct{}

// func (m MyType[T]) Method[U any]() {}  // ERROR!

// Cannot specialize methods
type Container[T any] struct {
    value T
}

// Cannot have different implementations for different types
// All types share the same compiled code (GC shape stenciling)

// No variadic type parameters
// func Call[T... any](f func(T...), args...T)  // Not supported
```

---

## Rust Generics

Rust offers powerful generics with trait bounds:

```rust
// Basic generic function
fn min<T: Ord>(a: T, b: T) -> T {
    if a < b { a } else { b }
}

// Generic struct
struct Stack<T> {
    items: Vec<T>,
}

impl<T> Stack<T> {
    fn new() -> Self {
        Stack { items: Vec::new() }
    }

    fn push(&mut self, item: T) {
        self.items.push(item);
    }

    fn pop(&mut self) -> Option<T> {
        self.items.pop()
    }
}

// Trait bounds
use std::ops::Add;

fn sum<T: Add<Output = T> + Default>(values: &[T]) -> T {
    values.iter().fold(T::default(), |acc, x| acc + x)
}

// Multiple trait bounds
use std::fmt::Display;

fn print_and_clone<T: Display + Clone>(item: &T) -> T {
    println!("{}", item);
    item.clone()
}

// Where clauses for complex bounds
fn complex_function<T, U>(t: T, u: U) -> i32
where
    T: Display + Clone,
    U: Clone + Default,
{
    0
}

// Associated types in traits
trait Container {
    type Item;
    fn get(&self) -> Option<&Self::Item>;
    fn set(&mut self, item: Self::Item);
}

struct BoxContainer<T> {
    value: Option<T>,
}

impl<T> Container for BoxContainer<T> {
    type Item = T;

    fn get(&self) -> Option<&T> {
        self.value.as_ref()
    }

    fn set(&mut self, item: T) {
        self.value = Some(item);
    }
}

// Const generics
struct Array<T, const N: usize> {
    data: [T; N],
}

impl<T: Default + Copy, const N: usize> Default for Array<T, N> {
    fn default() -> Self {
        Array { data: [T::default(); N] }
    }
}

// Lifetime parameters with generics
struct RefContainer<'a, T> {
    value: &'a T,
}

impl<'a, T> RefContainer<'a, T> {
    fn new(value: &'a T) -> Self {
        RefContainer { value }
    }
}
```

---

## Java Generics

Java uses type erasure for backward compatibility:

```java
// Generic class
public class Stack<T> {
    private List<T> items = new ArrayList<>();

    public void push(T item) {
        items.add(item);
    }

    public T pop() {
        if (items.isEmpty()) {
            throw new EmptyStackException();
        }
        return items.remove(items.size() - 1);
    }

    public boolean isEmpty() {
        return items.isEmpty();
    }
}

// Bounded type parameters
public class NumberStack<T extends Number> {
    private List<T> items = new ArrayList<>();

    public double sum() {
        return items.stream()
            .mapToDouble(Number::doubleValue)
            .sum();
    }
}

// Multiple bounds
interface Drawable {
    void draw();
}

class BoundedExample<T extends Number & Drawable> {
    // T must extend Number AND implement Drawable
}

// Wildcards
public static void printList(List<?> list) {
    for (Object elem : list) {
        System.out.print(elem + " ");
    }
}

// PECS: Producer Extends, Consumer Super
public static void copy(List<? extends T> src, List<? super T> dest) {
    for (T item : src) {
        dest.add(item);
    }
}

// Generic methods
public static <T> T getFirst(List<T> list) {
    return list.isEmpty() ? null : list.get(0);
}

// Type inference with diamond operator
Stack<String> stack = new Stack<>();  // <> infers String

// Generic interfaces
interface Container<T> {
    T get();
    void set(T value);
}

// Erasure limitations
public class ErasureExample {
    // At runtime, both are just List
    public void method(List<String> list) {}  // Erased to List
    // public void method(List<Integer> list) {}  // ERROR: same signature after erasure
}
```

---

## C++ Templates

C++ templates are compile-time code generation:

```cpp
#include <vector>
#include <algorithm>

// Function template
template<typename T>
T min(T a, T b) {
    return (a < b) ? a : b;
}

// Class template
template<typename T>
class Stack {
    std::vector<T> items;

public:
    void push(const T& item) {
        items.push_back(item);
    }

    void push(T&& item) {
        items.push_back(std::move(item));
    }

    T pop() {
        if (items.empty()) {
            throw std::runtime_error("Stack is empty");
        }
        T item = std::move(items.back());
        items.pop_back();
        return item;
    }

    bool empty() const {
        return items.empty();
    }
};

// Template specialization
template<>
class Stack<bool> {
    // Optimized implementation for bool
    std::vector<unsigned char> data;

public:
    void push(bool value) {
        // Bit-packing implementation
    }
};

// Concepts (C++20)
template<typename T>
concept Numeric = std::is_arithmetic_v<T>;

template<Numeric T>
T sum(const std::vector<T>& values) {
    T result = 0;
    for (const auto& v : values) {
        result += v;
    }
    return result;
}

// Variadic templates
template<typename T>
T sum(T value) {
    return value;
}

template<typename T, typename... Args>
T sum(T first, Args... args) {
    return first + sum(args...);
}

// SFINAE (Substitution Failure Is Not An Error)
template<typename T>
typename std::enable_if_t<std::is_integral_v<T>, bool>
is_even(T value) {
    return value % 2 == 0;
}

// Template metaprogramming
template<int N>
struct Factorial {
    static const int value = N * Factorial<N - 1>::value;
};

template<>
struct Factorial<0> {
    static const int value = 1;
};

constexpr int fact5 = Factorial<5>::value;  // 120, computed at compile time
```

---

## TypeScript Generics

TypeScript offers structural generics with powerful inference:

```typescript
// Generic function
function min<T>(a: T, b: T): T {
    return a < b ? a : b;
}

// Generic class
class Stack<T> {
    private items: T[] = [];

    push(item: T): void {
        this.items.push(item);
    }

    pop(): T | undefined {
        return this.items.pop();
    }

    peek(): T | undefined {
        return this.items[this.items.length - 1];
    }
}

// Generic constraints
interface HasLength {
    length: number;
}

function logLength<T extends HasLength>(arg: T): T {
    console.log(arg.length);
    return arg;
}

// Keyof constraint
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
    return obj[key];
}

// Mapped types
type Readonly<T> = {
    readonly [P in keyof T]: T[P];
};

type Partial<T> = {
    [P in keyof T]?: T[P];
};

type Pick<T, K extends keyof T> = {
    [P in K]: T[P];
};

// Conditional types
type IsString<T> = T extends string ? true : false;

type ReturnType<T> = T extends (...args: any[]) => infer R ? R : never;

// Utility function with inference
async function fetchJSON<T>(url: string): Promise<T> {
    const response = await fetch(url);
    return response.json() as Promise<T>;
}

// Usage with inferred type
interface User {
    id: number;
    name: string;
}

const user = await fetchJSON<User>("/api/user/1");

// Generic React component
interface ListProps<T> {
    items: T[];
    renderItem: (item: T) => React.ReactNode;
}

function List<T>({ items, renderItem }: ListProps<T>) {
    return (
        <ul>
            {items.map((item, index) => (
                <li key={index}>{renderItem(item)}</li>
            ))}
        </ul>
    );
}
```

---

## C# Generics

C# offers reified generics with runtime type preservation:

```csharp
// Generic class
public class Stack<T>
{
    private List<T> items = new List<T>();

    public void Push(T item) => items.Add(item);

    public T Pop()
    {
        if (items.Count == 0)
            throw new InvalidOperationException("Stack is empty");

        var item = items[^1];
        items.RemoveAt(items.Count - 1);
        return item;
    }
}

// Constraints
public class GenericConstraints<T> where T : class, new()
{
    public T CreateInstance() => new T();
}

// Multiple type parameters with constraints
public class Container<TKey, TValue>
    where TKey : notnull
    where TValue : class
{
    private Dictionary<TKey, TValue> items = new();

    public void Add(TKey key, TValue value) => items.Add(key, value);

    public TValue Get(TKey key) => items[key];
}

// Generic methods
public static T[] Slice<T>(T[] array, int start, int length)
{
    T[] result = new T[length];
    Array.Copy(array, start, result, 0, length);
    return result;
}

// Extension methods with generics
public static class LinqExtensions
{
    public static T? FirstOrNull<T>(this IEnumerable<T> source)
        where T : struct
    {
        foreach (var item in source)
            return item;
        return null;
    }
}

// Generic interfaces
public interface IContainer<T>
{
    T Value { get; set; }
}

// Covariance and contravariance
public interface IProducer<out T>  // Covariant
{
    T Produce();
}

public interface IConsumer<in T>   // Contravariant
{
    void Consume(T item);
}

// Pattern matching with generics
public static string Describe<T>(T value) => value switch
{
    int i => $"Integer: {i}",
    string s => $"String: {s}",
    null => "Null",
    _ => $"Other: {value}"
};

// Records with generics
public record Result<T>(bool Success, T Value, string? Error);

// Usage
var result = new Result<int>(true, 42, null);
```

---

## Feature Matrix

| Feature | Go | Rust | Java | C++ | TypeScript | C# |
|---------|-----|------|------|-----|------------|-----|
| Type Parameters | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Type Constraints | ~ | ✓ | ✓ | Concepts | ✓ | ✓ |
| Type Inference | ✓ | ✓ | ✓ | Partial | ✓ | ✓ |
| Specialization | ✗ | ✓ | ✗ | ✓ | ✗ | ✗ |
| Const Generics | ✗ | ✓ | ✗ | ✓ | ✗ | Partial |
| Higher-Kinded | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Variance | ✗ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Metaprogramming | ✗ | Macros | ✗ | Templates | ✗ | ✗ |
| Runtime Types | ✓ | ✓ | Erased | ✓ | Erased | ✓ |
| Compile-time Eval | ✗ | Const | ✗ | Constexpr | ✗ | Limited |

---

## Performance Comparison

| Language | Implementation | Overhead | Notes |
|----------|---------------|----------|-------|
| Go | GC Shape Stenciling | Low | Limited instantiations |
| Rust | Monomorphization | Zero | Code bloat possible |
| Java | Type Erasure | Boxing | No primitives in generics |
| C++ | Template Expansion | Zero | Compile-time cost |
| TypeScript | Erasure | None | Compile-time only |
| C# | Reified Generics | Low | Runtime type info |

---

## Summary

**Go Generics:**

- Simple and easy to understand
- Good for common use cases
- Limited advanced features
- Fast compilation

**Recommendation:**

- Use Go generics for: Containers, algorithms, result types
- Consider interfaces for: Behavior abstraction, loose coupling

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Size: ~22KB*
