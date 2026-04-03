# Type System Comparison: Go vs Other Languages

## Executive Summary

Type systems define how languages categorize and check data types, affecting safety, expressiveness, and performance. This document compares Go's simple static typing with advanced systems in Haskell, TypeScript, Rust, and others.

---

## Table of Contents

- [Type System Comparison: Go vs Other Languages](#type-system-comparison-go-vs-other-languages)
  - [Executive Summary](#executive-summary)
  - [Table of Contents](#table-of-contents)
  - [Go: Simple Static Typing](#go-simple-static-typing)
  - [Haskell: Advanced Static](#haskell-advanced-static)
  - [TypeScript: Gradual Typing](#typescript-gradual-typing)
  - [Rust: Ownership + Types](#rust-ownership--types)
  - [Python: Dynamic + Optional](#python-dynamic--optional)
  - [Java: Nominal Subtyping](#java-nominal-subtyping)
  - [Scala: Advanced OOP + FP](#scala-advanced-oop--fp)
  - [Feature Matrix](#feature-matrix)

---

## Go: Simple Static Typing

Go uses a simple, pragmatic type system focused on clarity:

```go
package main

// Basic types
var i int = 42
var f float64 = 3.14
var s string = "hello"
var b bool = true

// Struct types
type Person struct {
    Name string
    Age  int
}

// Interface types (structural subtyping)
type Greeter interface {
    Greet() string
}

// Type aliases
type UserID = int64

// Defined types (distinct from underlying)
type Miles int64
type Kilometers int64

func (m Miles) ToKilometers() Kilometers {
    return Kilometers(float64(m) * 1.60934)
}

// Generics (Go 1.18+)
type Container[T any] struct {
    Value T
}

func (c Container[T]) Get() T {
    return c.Value
}

// Type constraints
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64
}

func Min[T Number](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Interface embedding
type ReadWriter interface {
    Reader
    Writer
}

type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Type assertions
var i interface{} = "hello"
s := i.(string)  // Panic if wrong type
s, ok := i.(string)  // Safe assertion

// Type switches
func describe(i interface{}) string {
    switch v := i.(type) {
    case int:
        return fmt.Sprintf("int: %d", v)
    case string:
        return fmt.Sprintf("string: %s", v)
    case Person:
        return fmt.Sprintf("Person: %s", v.Name)
    default:
        return fmt.Sprintf("unknown: %T", v)
    }
}

// Type inference
x := 42  // Inferred as int
y := 3.14  // Inferred as float64
z := "hello"  // Inferred as string

// Empty interface (any)
var anything any = 42
anything = "hello"
anything = struct{ X int }{42}
```

---

## Haskell: Advanced Static

Haskell features one of the most sophisticated type systems:

```haskell
-- Type declarations
x :: Int
x = 42

name :: String
name = "Hello"

-- Parametric polymorphism (generics)
identity :: a -> a
identity x = x

-- Multiple type variables
pair :: a -> b -> (a, b)
pair x y = (x, y)

-- Type classes (ad-hoc polymorphism)
class Describable a where
    describe :: a -> String

instance Describable Int where
    describe n = "Integer: " ++ show n

instance Describable Bool where
    describe b = if b then "True" else "False"

-- Constraints
showAndDescribe :: (Show a, Describable a) => a -> String
showAndDescribe x = show x ++ " - " ++ describe x

-- Algebraic Data Types
data Color = Red | Green | Blue
data Maybe a = Nothing | Just a
data Either a b = Left a | Right b

-- Record types
data Person = Person
    { personName :: String
    , personAge :: Int
    } deriving (Show, Eq)

-- Recursive types
data Tree a = Leaf a | Node (Tree a) (Tree a)

-- Type synonyms
type String = [Char]
type Point = (Double, Double)

-- Newtypes (distinct from underlying)
newtype UserID = UserID Int
deriving (Show, Eq)

-- Phantom types
data Temperature a = Temperature Double

data Celsius
data Fahrenheit

toCelsius :: Temperature Fahrenheit -> Temperature Celsius
toCelsius (Temperature f) = Temperature ((f - 32) * 5 / 9)

-- GADTs (Generalized Algebraic Data Types)
data Expr a where
    IntLit :: Int -> Expr Int
    BoolLit :: Bool -> Expr Bool
    Add :: Expr Int -> Expr Int -> Expr Int
    If :: Expr Bool -> Expr a -> Expr a -> Expr a

eval :: Expr a -> a
eval (IntLit n) = n
eval (BoolLit b) = b
eval (Add e1 e2) = eval e1 + eval e2
eval (If cond then_ else_) = if eval cond then eval then_ else eval else_

-- Type families
type family If (c :: Bool) (t :: *) (f :: *) :: * where
    If 'True t f = t
    If 'False t f = f

-- Higher-kinded types
class Functor f where
    fmap :: (a -> b) -> f a -> f b

instance Functor Maybe where
    fmap _ Nothing = Nothing
    fmap f (Just x) = Just (f x)

instance Functor [] where
    fmap = map

-- Monads
class Monad m where
    return :: a -> m a
    (>>=) :: m a -> (a -> m b) -> m b

instance Monad Maybe where
    return = Just
    Nothing >>= _ = Nothing
    Just x >>= f = f x

-- Type-level programming
{-# LANGUAGE DataKinds #-}
{-# LANGUAGE TypeFamilies #-}

data Nat = Z | S Nat

type family Add (n :: Nat) (m :: Nat) :: Nat where
    Add 'Z m = m
    Add ('S n) m = 'S (Add n m)

-- Dependent types (with singletons)
data SNat (n :: Nat) where
    SZ :: SNat 'Z
    SS :: SNat n -> SNat ('S n)
```

---

## TypeScript: Gradual Typing

TypeScript adds static types to JavaScript:

```typescript
// Basic types
let num: number = 42;
let str: string = "hello";
let bool: boolean = true;
let arr: number[] = [1, 2, 3];
let tuple: [string, number] = ["age", 25];

// Object types
interface Person {
    name: string;
    age: number;
    email?: string;  // Optional
    readonly id: number;  // Readonly
}

// Union types
type Status = "pending" | "approved" | "rejected";
type ID = string | number;

// Intersection types
type Employee = Person & {
    employeeId: string;
    department: string;
};

// Type aliases
type Point = {
    x: number;
    y: number;
};

type UserID = string;

// Generics
function identity<T>(arg: T): T {
    return arg;
}

interface Container<T> {
    value: T;
    get(): T;
}

// Generic constraints
interface HasLength {
    length: number;
}

function logLength<T extends HasLength>(arg: T): T {
    console.log(arg.length);
    return arg;
}

// keyof operator
type PersonKeys = keyof Person;  // "name" | "age" | "email" | "id"

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

type Record<K extends keyof any, T> = {
    [P in K]: T;
};

// Conditional types
type IsString<T> = T extends string ? true : false;

type NonNullable<T> = T extends null | undefined ? never : T;

type ReturnType<T> = T extends (...args: any[]) => infer R ? R : never;

type Parameters<T> = T extends (...args: infer P) => any ? P : never;

// Template literal types
type EventName<T extends string> = `on${Capitalize<T>}`;
// EventName<"click"> = "onClick"

// Variadic tuple types
type Concat<T extends readonly unknown[], U extends readonly unknown[]> = [...T, ...U];

// Structural typing
interface Named {
    name: string;
}

class Person {
    name: string;
    constructor(name: string) {
        this.name = name;
    }
}

let p: Named = new Person("John");  // OK - structural match

// Discriminated unions
type Shape =
    | { kind: "circle"; radius: number }
    | { kind: "rectangle"; width: number; height: number }
    | { kind: "square"; side: number };

function area(shape: Shape): number {
    switch (shape.kind) {
        case "circle":
            return Math.PI * shape.radius ** 2;
        case "rectangle":
            return shape.width * shape.height;
        case "square":
            return shape.side ** 2;
    }
}

// Type guards
function isString(x: unknown): x is string {
    return typeof x === "string";
}

// Assertion functions
function assertDefined<T>(value: T | undefined | null): asserts value is T {
    if (value === undefined || value === null) {
        throw new Error("Value is not defined");
    }
}
```

---

## Rust: Ownership + Types

Rust combines ownership with a powerful type system:

```rust
// Basic types
let i: i32 = 42;
let f: f64 = 3.14;
let s: String = String::from("hello");
let b: bool = true;

// Struct types
struct Person {
    name: String,
    age: u32,
}

// Tuple structs
struct Point(f64, f64);

// Unit structs
struct Empty;

// Generic types
struct Container<T> {
    value: T,
}

impl<T> Container<T> {
    fn new(value: T) -> Self {
        Container { value }
    }
}

// Traits (interfaces)
trait Drawable {
    fn draw(&self);
    fn bounds(&self) -> Rect;
}

struct Circle {
    radius: f64,
}

impl Drawable for Circle {
    fn draw(&self) {
        println!("Drawing circle");
    }

    fn bounds(&self) -> Rect {
        // ...
    }
}

// Trait bounds
fn draw_all<T: Drawable>(items: &[T]) {
    for item in items {
        item.draw();
    }
}

// Multiple bounds
use std::fmt::Display;

fn print_and_clone<T: Display + Clone>(item: &T) -> T {
    println!("{}", item);
    item.clone()
}

// Where clauses
fn complex<T, U>(t: T, u: U) -> i32
where
    T: Display + Clone,
    U: Clone + Default,
{
    0
}

// Associated types
trait Iterator {
    type Item;
    fn next(&mut self) -> Option<Self::Item>;
}

// Lifetimes
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() { x } else { y }
}

// Lifetime elision
fn first_word(s: &str) -> &str {
    &s[..s.find(' ').unwrap_or(s.len())]
}

// Smart pointers
use std::rc::Rc;
use std::sync::Arc;
use std::cell::RefCell;

// Enums with data
enum Message {
    Quit,
    Move { x: i32, y: i32 },
    Write(String),
    ChangeColor(i32, i32, i32),
}

// Option and Result
fn divide(a: f64, b: f64) -> Option<f64> {
    if b == 0.0 {
        None
    } else {
        Some(a / b)
    }
}

fn may_fail() -> Result<String, Error> {
    Ok("success".to_string())
}

// Const generics
struct Array<T, const N: usize> {
    data: [T; N],
}

// Type inference
let x = 42;  // i32
let y = vec![1, 2, 3];  // Vec<i32>
let z = "hello";  // &str

// Never type
fn diverges() -> ! {
    panic!("Never returns");
}

// Phantom types
use std::marker::PhantomData;

struct Tag<T>(PhantomData<T>);

struct Kilograms;
struct Pounds;

fn kilograms_to_pounds(_: Tag<Kilograms>) -> Tag<Pounds> {
    Tag(PhantomData)
}
```

---

## Python: Dynamic + Optional

```python
# Dynamic typing (runtime)
x = 42  # int
x = "hello"  # now str
x = [1, 2, 3]  # now list

# Type hints (optional, static check)
from typing import (
    List, Dict, Optional, Union, Callable,
    TypeVar, Generic, Protocol
)

def greet(name: str) -> str:
    return f"Hello, {name}"

# Generic functions
T = TypeVar('T')

def first(items: List[T]) -> Optional[T]:
    return items[0] if items else None

# Generic classes
class Container(Generic[T]):
    def __init__(self, value: T) -> None:
        self.value = value

    def get(self) -> T:
        return self.value

# Union types
def process(value: Union[int, str]) -> str:
    if isinstance(value, int):
        return str(value)
    return value

# Python 3.10+ syntax
def process_new(value: int | str) -> str:
    ...

# Optional
def find_user(user_id: int) -> Optional[User]:
    ...

# Python 3.10+
def find_user_new(user_id: int) -> User | None:
    ...

# Protocols (structural subtyping)
class Drawable(Protocol):
    def draw(self) -> None: ...

def render(item: Drawable) -> None:
    item.draw()

# Callable
Handler = Callable[[int, str], bool]

# Type aliases
Vector = List[float]
Matrix = List[Vector]
UserDict = Dict[int, User]

# Final
from typing import Final
MAX_SIZE: Final[int] = 100

# Literal types
from typing import Literal
Direction = Literal["north", "south", "east", "west"]

def move(direction: Direction) -> None:
    ...

# TypedDict
from typing import TypedDict

class UserDict(TypedDict):
    name: str
    age: int
    email: Optional[str]

# Runtime check
from typing import get_type_hints
hints = get_type_hints(greet)
```

---

## Java: Nominal Subtyping

```java
// Class-based nominal typing
public class Person {
    private final String name;
    private final int age;

    public Person(String name, int age) {
        this.name = name;
        this.age = age;
    }

    // Getters...
}

// Interface inheritance
public interface Drawable {
    void draw();
    default void print() {
        System.out.println("Printing...");
    }
}

public class Circle implements Drawable {
    @Override
    public void draw() {
        System.out.println("Drawing circle");
    }
}

// Generics
public class Container<T> {
    private T value;

    public Container(T value) {
        this.value = value;
    }

    public T get() {
        return value;
    }
}

// Bounded type parameters
public static <T extends Comparable<T>> T max(T a, T b) {
    return a.compareTo(b) > 0 ? a : b;
}

// Wildcards
public static void printList(List<?> list) {
    for (Object item : list) {
        System.out.println(item);
    }
}

// Variance: Producer extends, Consumer super
public static <T> void copy(List<? extends T> src, List<? super T> dest) {
    for (T item : src) {
        dest.add(item);
    }
}

// Records (Java 16+)
public record Point(int x, int y) {}

// Sealed classes (Java 17+)
public abstract sealed class Shape
    permits Circle, Rectangle, Square {
}

public final class Circle extends Shape {
    private final double radius;
    public Circle(double radius) { this.radius = radius; }
}

// Pattern matching (Java 17+)
public static double area(Shape shape) {
    if (shape instanceof Circle c) {
        return Math.PI * c.radius() * c.radius();
    }
    // ...
    throw new IllegalArgumentException();
}

// Switch expressions
public static String describe(DayOfWeek day) {
    return switch (day) {
        case MONDAY -> "Start of week";
        case FRIDAY -> "End of week";
        case SATURDAY, SUNDAY -> "Weekend";
        default -> "Midweek";
    };
}
```

---

## Scala: Advanced OOP + FP

```scala
// Object-oriented
class Person(val name: String, val age: Int) {
  def greet(): String = s"Hello, I'm $name"
}

// Case classes (value types)
case class User(id: Long, name: String, email: String)

// Pattern matching
val user = User(1, "John", "john@example.com")
user match {
  case User(_, name, _) if name.startsWith("J") => println("Starts with J")
  case User(1, _, _) => println("First user")
  case _ => println("Other")
}

// Traits (mixins)
trait Logger {
  def log(msg: String): Unit = println(s"[LOG] $msg")
}

trait Validatable {
  def validate(): Boolean
}

class Service extends Logger with Validatable {
  def validate(): Boolean = true
}

// Generics with variance
class Container[+A](value: A) {  // Covariant
  def get: A = value
}

class Writer[-A] {  // Contravariant
  def write(value: A): Unit = ()
}

// Higher-kinded types
trait Functor[F[_]] {
  def map[A, B](fa: F[A])(f: A => B): F[B]
}

// Implicits (Scala 2) / Given (Scala 3)
given Functor[List] with {
  def map[A, B](fa: List[A])(f: A => B): List[B] = fa.map(f)
}

// Type classes
trait Show[A] {
  def show(a: A): String
}

object Show {
  given Show[Int] with {
    def show(n: Int): String = n.toString
  }
}

// Extension methods
extension (s: String)
  def greet: String = s"Hello, $s!"

// Opaque types (Scala 3)
opaque type UserId = Long
object UserId {
  def apply(value: Long): UserId = value
  extension (id: UserId) def value: Long = id
}

// Union and intersection types (Scala 3)
type ErrorOr[T] = T | Error
type CloneableResetable = Cloneable & Resetable
```

---

## Feature Matrix

| Feature | Go | Haskell | TypeScript | Rust | Python | Java | Scala |
|---------|-----|---------|------------|------|--------|------|-------|
| Static Checking | Yes | Yes | Optional | Yes | Optional | Yes | Yes |
| Inference | Yes | Yes | Yes | Yes | N/A | Partial | Yes |
| Generics | Yes | Yes | Yes | Yes | Yes | Yes | Yes |
| Higher-Kinded | No | Yes | Yes | No | No | No | Yes |
| Dependent Types | No | Partial | No | Partial | No | No | No |
| Variance | No | Yes | Yes | No | Yes | Yes | Yes |
| Subtyping | Structural | No | Structural | No | Duck | Nominal | Nominal |
| Type Classes | Interfaces | Yes | No | Traits | Protocols | No | Yes |
| Null Safety | Partial | Yes | Partial | Yes | Optional | Partial | Partial |

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Size: ~23KB*
