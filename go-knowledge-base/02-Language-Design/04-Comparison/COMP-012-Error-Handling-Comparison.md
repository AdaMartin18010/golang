# Error Handling Comparison: Go vs Other Languages

## Executive Summary

Error handling is a fundamental aspect of software reliability. Go's explicit error returns contrast with exceptions in Java/Python, Result types in Rust, and optionals in Swift. This document compares approaches across languages with code examples and best practices.

---

## Table of Contents

- [Error Handling Comparison: Go vs Other Languages](#error-handling-comparison-go-vs-other-languages)
  - [Executive Summary](#executive-summary)
  - [Table of Contents](#table-of-contents)
  - [Go: Explicit Error Returns](#go-explicit-error-returns)
  - [Rust: Result Type](#rust-result-type)
  - [Java: Exceptions](#java-exceptions)
  - [Python: Exceptions](#python-exceptions)
  - [Swift: Optionals and Throws](#swift-optionals-and-throws)
  - [TypeScript/JavaScript](#typescriptjavascript)
  - [Comparison Matrix](#comparison-matrix)
  - [Best Practices](#best-practices)
    - [Go Error Handling Best Practices](#go-error-handling-best-practices)

---

## Go: Explicit Error Returns

Go treats errors as values, forcing explicit handling:

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

// Basic error return
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

func validateUser(name string, age int) error {
    if name == "" {
        return &ValidationError{Field: "name", Message: "cannot be empty"}
    }
    if age < 0 || age > 150 {
        return &ValidationError{Field: "age", Message: "must be between 0 and 150"}
    }
    return nil
}

// Error wrapping (Go 1.13+)
func readConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config from %s: %w", path, err)
    }

    config, err := parseConfig(data)
    if err != nil {
        return nil, fmt.Errorf("parsing config: %w", err)
    }

    return config, nil
}

// Error checking with errors.Is
func handleError(err error) {
    if err == nil {
        return
    }

    // Check specific error
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("File not found")
        return
    }

    // Check error type with errors.As
    var valErr *ValidationError
    if errors.As(err, &valErr) {
        fmt.Printf("Validation failed: %s\n", valErr.Error())
        return
    }

    fmt.Printf("Unexpected error: %v\n", err)
}

// Sentinel errors
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid credentials")

func findUser(id string) (*User, error) {
    // ...
    if user == nil {
        return nil, ErrUserNotFound
    }
    return user, nil
}

// Must pattern for must-not-fail scenarios
func mustOpen(path string) *os.File {
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    return f
}

// Structured error handling
type Result[T any] struct {
    Value T
    Error error
}

func (r Result[T]) Unwrap() (T, error) {
    return r.Value, r.Error
}

func (r Result[T]) Or(defaultValue T) T {
    if r.Error != nil {
        return defaultValue
    }
    return r.Value
}
```

---

## Rust: Result Type

Rust uses Result<T, E> for explicit error handling:

```rust
use std::fs;
use std::io;

// Result type
fn divide(a: f64, b: f64) -> Result<f64, String> {
    if b == 0.0 {
        Err("division by zero".to_string())
    } else {
        Ok(a / b)
    }
}

// Custom error type
#[derive(Debug)]
struct ValidationError {
    field: String,
    message: String,
}

impl std::fmt::Display for ValidationError {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "validation error for {}: {}", self.field, self.message)
    }
}

impl std::error::Error for ValidationError {}

fn validate_user(name: &str, age: i32) -> Result<(), ValidationError> {
    if name.is_empty() {
        return Err(ValidationError {
            field: "name".to_string(),
            message: "cannot be empty".to_string(),
        });
    }
    if age < 0 || age > 150 {
        return Err(ValidationError {
            field: "age".to_string(),
            message: "must be between 0 and 150".to_string(),
        });
    }
    Ok(())
}

// The ? operator for propagation
fn read_config(path: &str) -> Result<Config, Box<dyn std::error::Error>> {
    let data = fs::read_to_string(path)?;
    let config = parse_config(&data)?;
    Ok(config)
}

// unwrap and expect for quick prototyping
let result = divide(10.0, 2.0).unwrap();
let result = divide(10.0, 2.0).expect("division should work");

// unwrap_or for defaults
let result = divide(10.0, 0.0).unwrap_or(f64::INFINITY);

// match on Result
match divide(10.0, 0.0) {
    Ok(result) => println!("Result: {}", result),
    Err(e) => println!("Error: {}", e),
}

// if let for single case
if let Ok(result) = divide(10.0, 2.0) {
    println!("Result: {}", result);
}

// map and map_err for transforming
let result = divide(10.0, 2.0)
    .map(|r| r * 2.0)
    .map_err(|e| format!("Calculation failed: {}", e));

// Result chaining with and_then
let result = read_file("config.txt")
    .and_then(|content| parse_config(&content))
    .and_then(|config| validate_config(config));

// thiserror for derive macros
use thiserror::Error;

#[derive(Error, Debug)]
enum AppError {
    #[error("IO error: {0}")]
    Io(#[from] io::Error),

    #[error("Parse error: {0}")]
    Parse(String),

    #[error("Validation error: {message}")]
    Validation { field: String, message: String },
}
```

---

## Java: Exceptions

Java uses checked and unchecked exceptions:

```java
// Checked exception
public class ConfigNotFoundException extends Exception {
    public ConfigNotFoundException(String message) {
        super(message);
    }

    public ConfigNotFoundException(String message, Throwable cause) {
        super(message, cause);
    }
}

// Unchecked exception
public class ValidationException extends RuntimeException {
    private final String field;

    public ValidationException(String field, String message) {
        super(message);
        this.field = field;
    }

    public String getField() { return field; }
}

public class ConfigService {
    // Throws checked exception
    public Config loadConfig(String path) throws ConfigNotFoundException {
        try {
            String content = Files.readString(Path.of(path));
            return parseConfig(content);
        } catch (IOException e) {
            throw new ConfigNotFoundException("Failed to load config: " + path, e);
        }
    }

    // Throws unchecked exception
    public void validateUser(String name, int age) {
        if (name == null || name.isEmpty()) {
            throw new ValidationException("name", "cannot be empty");
        }
        if (age < 0 || age > 150) {
            throw new ValidationException("age", "must be between 0 and 150");
        }
    }

    // Try-catch-finally
    public void processFile(String path) {
        FileInputStream fis = null;
        try {
            fis = new FileInputStream(path);
            process(fis);
        } catch (FileNotFoundException e) {
            logger.error("File not found: {}", path);
        } catch (IOException e) {
            logger.error("IO error", e);
        } finally {
            if (fis != null) {
                try {
                    fis.close();
                } catch (IOException e) {
                    logger.warn("Failed to close file", e);
                }
            }
        }
    }

    // Try-with-resources (Java 7+)
    public void processFileModern(String path) {
        try (FileInputStream fis = new FileInputStream(path);
             BufferedReader reader = new BufferedReader(new InputStreamReader(fis))) {
            process(reader);
        } catch (IOException e) {
            logger.error("IO error", e);
        }
    }

    // Optional for null safety
    public Optional<User> findUser(String id) {
        User user = repository.findById(id);
        return Optional.ofNullable(user);
    }

    // Usage
    public void displayUser(String id) {
        findUser(id)
            .ifPresentOrElse(
                user -> System.out.println(user.getName()),
                () -> System.out.println("User not found")
            );
    }
}

// Result-like pattern with Vavr library
import io.vavr.control.Either;
import io.vavr.control.Try;

public Either<Error, Config> loadConfigSafe(String path) {
    return Try.of(() -> Files.readString(Path.of(path)))
        .mapTry(this::parseConfig)
        .toEither()
        .mapLeft(Throwable::getMessage);
}
```

---

## Python: Exceptions

Python uses exceptions with try/except:

```python
# Custom exceptions
class ValidationError(Exception):
    def __init__(self, field: str, message: str):
        self.field = field
        self.message = message
        super().__init__(f"Validation error for {field}: {message}")

class ConfigNotFoundError(FileNotFoundError):
    pass

# Exception raising
def divide(a: float, b: float) -> float:
    if b == 0:
        raise ValueError("division by zero")
    return a / b

def validate_user(name: str, age: int) -> None:
    if not name:
        raise ValidationError("name", "cannot be empty")
    if not 0 <= age <= 150:
        raise ValidationError("age", "must be between 0 and 150")

# Exception handling
def load_config(path: str) -> dict:
    try:
        with open(path, 'r') as f:
            content = f.read()
        return parse_config(content)
    except FileNotFoundError as e:
        raise ConfigNotFoundError(f"Config not found: {path}") from e
    except json.JSONDecodeError as e:
        raise ValueError(f"Invalid JSON in {path}") from e
    except Exception as e:
        logger.exception("Unexpected error loading config")
        raise

# Multiple except blocks
def process_data(data: dict):
    try:
        result = calculate(data)
    except KeyError as e:
        logger.error(f"Missing key: {e}")
        return None
    except ValueError as e:
        logger.error(f"Invalid value: {e}")
        return default_value
    except Exception as e:
        logger.exception("Unexpected error")
        raise

# Else and finally
def read_file_safe(path: str) -> str:
    try:
        f = open(path, 'r')
    except FileNotFoundError:
        return ""
    else:
        content = f.read()
        f.close()
        return content
    finally:
        print("Operation completed")

# Context manager (with statement)
def read_file_modern(path: str) -> str:
    try:
        with open(path, 'r') as f:
            return f.read()
    except FileNotFoundError:
        return ""

# Result pattern with returns-Result library
from returns.result import Result, Success, Failure

def divide_safe(a: float, b: float) -> Result[float, str]:
    if b == 0:
        return Failure("division by zero")
    return Success(a / b)

# Usage
result = divide_safe(10, 0)
match result:
    case Success(value):
        print(f"Result: {value}")
    case Failure(error):
        print(f"Error: {error}")
```

---

## Swift: Optionals and Throws

Swift combines optionals with throwing functions:

```swift
// Optional type
func findUser(id: String) -> User? {
    // Returns nil if not found
    return database[id]
}

// Force unwrap (dangerous!)
let user = findUser(id: "123")!

// Optional binding
if let user = findUser(id: "123") {
    print(user.name)
} else {
    print("User not found")
}

// Guard statement
func processUser(id: String) {
    guard let user = findUser(id: id) else {
        print("User not found")
        return
    }
    // user is non-optional here
    print(user.name)
}

// Nil-coalescing
let name = findUser(id: "123")?.name ?? "Anonymous"

// Optional chaining
let street = findUser(id: "123")?.address?.street

// Throwing functions
enum ValidationError: Error {
    case invalidName
    case invalidAge(Int)
}

func validateUser(name: String, age: Int) throws {
    guard !name.isEmpty else {
        throw ValidationError.invalidName
    }
    guard age >= 0 && age <= 150 else {
        throw ValidationError.invalidAge(age)
    }
}

// Do-catch
do {
    try validateUser(name: "John", age: 25)
    print("Valid")
} catch ValidationError.invalidName {
    print("Invalid name")
} catch ValidationError.invalidAge(let age) {
    print("Invalid age: \(age)")
} catch {
    print("Unknown error: \(error)")
}

// Try?
let result = try? validateUser(name: "", age: 25)  // Returns Void? or nil on error

// Try!
try! validateUser(name: "John", age: 25)  // Crashes on error

// Result type
func divide(a: Double, b: Double) -> Result<Double, Error> {
    if b == 0 {
        return .failure(DivisionError.divideByZero)
    }
    return .success(a / b)
}

// Using Result
let result = divide(a: 10, b: 0)
switch result {
case .success(let value):
    print("Result: \(value)")
case .failure(let error):
    print("Error: \(error)")
}

// Result with map
let doubled = divide(a: 10, b: 2)
    .map { $0 * 2 }
    .mapError { error in
        CustomError.calculationFailed(error)
    }
```

---

## TypeScript/JavaScript

Modern TypeScript/JavaScript uses various approaches:

```typescript
// Traditional try-catch
function divide(a: number, b: number): number {
    if (b === 0) {
        throw new Error("division by zero");
    }
    return a / b;
}

try {
    const result = divide(10, 0);
} catch (error) {
    if (error instanceof Error) {
        console.error(error.message);
    }
}

// Result type pattern
type Result<T, E = Error> =
    | { success: true; data: T }
    | { success: false; error: E };

function divideSafe(a: number, b: number): Result<number, string> {
    if (b === 0) {
        return { success: false, error: "division by zero" };
    }
    return { success: true, data: a / b };
}

// Using Result
const result = divideSafe(10, 0);
if (result.success) {
    console.log(result.data);
} else {
    console.error(result.error);
}

// neverthrow library
import { ok, err, Result } from 'neverthrow';

function parseNumber(input: string): Result<number, string> {
    const parsed = parseFloat(input);
    if (isNaN(parsed)) {
        return err("Invalid number");
    }
    return ok(parsed);
}

// Chaining
const result = parseNumber("10")
    .map(n => n * 2)
    .andThen(n => divideSafe(n, 0));

// Async error handling
async function fetchUser(id: string): Promise<Result<User, Error>> {
    try {
        const response = await fetch(`/api/users/${id}`);
        if (!response.ok) {
            return err(new Error(`HTTP ${response.status}`));
        }
        const user = await response.json();
        return ok(user);
    } catch (error) {
        return err(error as Error);
    }
}

// Optional chaining and nullish coalescing
const userName = user?.profile?.name ?? "Anonymous";

// Assertion functions
function assertDefined<T>(value: T | undefined | null): asserts value is T {
    if (value === undefined || value === null) {
        throw new Error("Value is not defined");
    }
}
```

---

## Comparison Matrix

| Aspect | Go | Rust | Java | Python | Swift | TS/JS |
|--------|-----|------|------|--------|-------|-------|
| Explicit | ✓✓ | ✓✓ | Partial | ✗ | Partial | Variable |
| Compile-time | ✓ | ✓ | Checked ex | ✗ | ✓ | ✗ |
| Ergonomics | Good | Excellent | Poor | Good | Excellent | Good |
| Performance | Zero cost | Zero cost | Costly | Costly | Zero cost | Variable |
| Debugging | Easy | Easy | Hard | Hard | Easy | Variable |
| Boilerplate | Medium | Low | Medium | Low | Low | Low |

---

## Best Practices

### Go Error Handling Best Practices

```go
// 1. Check errors immediately
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close()

// 2. Wrap errors with context
if err := doSomething(); err != nil {
    return fmt.Errorf("doing something: %w", err)
}

// 3. Use sentinel errors for specific cases
if errors.Is(err, ErrNotFound) {
    // handle not found
}

// 4. Define custom error types for rich information
// 5. Don't ignore errors with _
// 6. Use panic only for unrecoverable errors
```

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Size: ~24KB*
