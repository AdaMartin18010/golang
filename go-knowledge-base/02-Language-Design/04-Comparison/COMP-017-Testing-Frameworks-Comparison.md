# Testing Frameworks Comparison

## Executive Summary

Testing frameworks vary significantly across languages in terms of features, syntax, and capabilities. This document compares testing approaches in Go, Rust, JavaScript, Java, Python, and C# with code examples and best practices.

---

## Table of Contents

- [Testing Frameworks Comparison](#testing-frameworks-comparison)
  - [Executive Summary](#executive-summary)
  - [Table of Contents](#table-of-contents)
  - [Go: Built-in Testing](#go-built-in-testing)
  - [Rust: Built-in with Macros](#rust-built-in-with-macros)
  - [JavaScript: Jest/Vitest](#javascript-jestvitest)
  - [Java: JUnit 5](#java-junit-5)
  - [Python: pytest](#python-pytest)
  - [C#: xUnit](#c-xunit)
  - [Feature Comparison](#feature-comparison)

---

## Go: Built-in Testing

Go includes testing in the standard library:

```go
package mypackage

import (
    "errors"
    "testing"
)

// Unit test
func Add(a, b int) int {
    return a + b
}

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}

// Table-driven tests (idiomatic Go)
func TestAddTableDriven(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
        {"mixed", -2, 3, 1},
        {"zero", 0, 5, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

// Benchmark tests
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(100, 200)
    }
}

// Benchmark with different inputs
func BenchmarkAddVarious(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Add(size, size)
            }
        })
    }
}

// Example tests (documentation + test)
func ExampleAdd() {
    result := Add(2, 3)
    fmt.Println(result)
    // Output: 5
}

// Fuzz testing (Go 1.18+)
func FuzzAdd(f *testing.F) {
    f.Add(1, 2)  // Seed corpus
    f.Fuzz(func(t *testing.T, a, b int) {
        result := Add(a, b)
        expected := a + b
        if result != expected {
            t.Errorf("Add(%d, %d) = %d; want %d", a, b, result, expected)
        }
    })
}

// Parallel tests
func TestParallel(t *testing.T) {
    tests := []struct {
        name string
        id   int
    }{
        {"test1", 1},
        {"test2", 2},
        {"test3", 3},
    }

    for _, tt := range tests {
        tt := tt  // capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // Run in parallel
            process(tt.id)
        })
    }
}

// Test helpers
type testCase struct {
    input    string
    expected int
    wantErr  bool
}

func runTestCases(t *testing.T, cases []testCase, fn func(string) (int, error)) {
    for _, tc := range cases {
        t.Run(tc.input, func(t *testing.T) {
            result, err := fn(tc.input)
            if tc.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            if result != tc.expected {
                t.Errorf("got %d, want %d", result, tc.expected)
            }
        })
    }
}

// Mocking with interfaces
type UserRepository interface {
    FindByID(id int) (*User, error)
}

type mockUserRepository struct {
    users map[int]*User
}

func (m *mockUserRepository) FindByID(id int) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("not found")
    }
    return user, nil
}

func TestUserService(t *testing.T) {
    mock := &mockUserRepository{
        users: map[int]*User{
            1: {ID: 1, Name: "John"},
        },
    }

    service := NewUserService(mock)
    user, err := service.GetUser(1)

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "John" {
        t.Errorf("got %q, want %q", user.Name, "John")
    }
}

// Subtests setup/teardown
func TestWithSetup(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(t, db)

    t.Run("insert", func(t *testing.T) {
        // test insert
    })

    t.Run("update", func(t *testing.T) {
        // test update
    })
}

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    return db
}

func teardownTestDB(t *testing.T, db *sql.DB) {
    if err := db.Close(); err != nil {
        t.Errorf("failed to close db: %v", err)
    }
}
```

---

## Rust: Built-in with Macros

```rust
// Unit tests (in same file)
fn add(a: i32, b: i32) -> i32 {
    a + b
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_add() {
        assert_eq!(add(2, 3), 5);
    }

    #[test]
    fn test_add_negative() {
        assert_eq!(add(-2, -3), -5);
    }

    // Custom error message
    #[test]
    fn test_add_with_message() {
        let result = add(2, 3);
        assert!(result == 5, "Expected 5, got {}", result);
    }

    // Should panic
    #[test]
    #[should_panic(expected = "divide by zero")]
    fn test_divide_by_zero() {
        divide(10, 0);
    }

    // Ignored test
    #[test]
    #[ignore = "not implemented yet"]
    fn test_future_feature() {
        // ...
    }

    // Using rstest for parameterized tests
    use rstest::rstest;

    #[rstest]
    #[case(2, 3, 5)]
    #[case(-2, -3, -5)]
    #[case(0, 5, 5)]
    fn test_add_parameterized(#[case] a: i32, #[case] b: i32, #[case] expected: i32) {
        assert_eq!(add(a, b), expected);
    }
}

// Integration tests (tests/ directory)
// tests/integration_test.rs
use mycrate::add;

#[test]
fn test_add_integration() {
    assert_eq!(add(10, 20), 30);
}

// Benchmarks with criterion
use criterion::{black_box, criterion_group, criterion_main, Criterion};

fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("add 100 + 200", |b| {
        b.iter(|| add(black_box(100), black_box(200)))
    });
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);

// Property-based testing with proptest
use proptest::prelude::*;

proptest! {
    #[test]
    fn test_add_commutative(a: i32, b: i32) {
        assert_eq!(add(a, b), add(b, a));
    }

    #[test]
    fn test_add_associative(a: i32, b: i32, c: i32) {
        assert_eq!(add(add(a, b), c), add(a, add(b, c)));
    }
}

// Mocking with mockall
use mockall::automock;

#[automock]
pub trait UserRepository {
    fn find_by_id(&self, id: i64) -> Option<User>;
}

#[cfg(test)]
mod service_tests {
    use super::*;

    #[test]
    fn test_get_user() {
        let mut mock = MockUserRepository::new();
        mock.expect_find_by_id()
            .with(eq(1))
            .times(1)
            .returning(|_| Some(User { id: 1, name: "John".to_string() }));

        let service = UserService::new(mock);
        let user = service.get_user(1).unwrap();
        assert_eq!(user.name, "John");
    }
}

// Async tests
#[tokio::test]
async fn test_async_fetch() {
    let result = fetch_data().await;
    assert!(result.is_ok());
}

// Snapshot testing with insta
#[test]
fn test_snapshot() {
    let result = generate_complex_output();
    insta::assert_yaml_snapshot!(result);
}
```

---

## JavaScript: Jest/Vitest

```javascript
// Jest test example
import { describe, it, expect, beforeEach, jest } from '@jest/globals';
import { add, UserService } from './utils';

// Unit test
describe('add', () => {
    it('should add two positive numbers', () => {
        expect(add(2, 3)).toBe(5);
    });

    it('should add negative numbers', () => {
        expect(add(-2, -3)).toBe(-5);
    });

    // Table-driven
    it.each([
        [2, 3, 5],
        [-2, -3, -5],
        [0, 5, 5],
    ])('add(%i, %i) should return %i', (a, b, expected) => {
        expect(add(a, b)).toBe(expected);
    });
});

// Mocking
describe('UserService', () => {
    let service;
    let mockRepository;

    beforeEach(() => {
        mockRepository = {
            findById: jest.fn(),
            save: jest.fn(),
        };
        service = new UserService(mockRepository);
    });

    it('should return user when found', async () => {
        const user = { id: 1, name: 'John' };
        mockRepository.findById.mockResolvedValue(user);

        const result = await service.getUser(1);

        expect(result).toEqual(user);
        expect(mockRepository.findById).toHaveBeenCalledWith(1);
    });

    it('should throw when user not found', async () => {
        mockRepository.findById.mockResolvedValue(null);

        await expect(service.getUser(1)).rejects.toThrow('User not found');
    });

    // Spy
    it('should call repository once', async () => {
        const spy = jest.spyOn(mockRepository, 'findById');
        await service.getUser(1);
        expect(spy).toHaveBeenCalledTimes(1);
    });
});

// Snapshot testing
describe('Component', () => {
    it('renders correctly', () => {
        const tree = renderer.create(<MyComponent />).toJSON();
        expect(tree).toMatchSnapshot();
    });
});

// Coverage
// jest --coverage

// Vitest (faster alternative)
import { describe, it, expect, vi } from 'vitest';

describe('with vitest', () => {
    it('uses vi for mocking', () => {
        const mockFn = vi.fn();
        mockFn('arg');
        expect(mockFn).toHaveBeenCalledWith('arg');
    });
});
```

---

## Java: JUnit 5

```java
import org.junit.jupiter.api.*;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.*;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@DisplayName("Calculator Tests")
public class CalculatorTest {

    private Calculator calculator;

    @BeforeAll
    static void initAll() {
        // Runs once before all tests
    }

    @BeforeEach
    void init() {
        calculator = new Calculator();
    }

    @Test
    @DisplayName("Should add two numbers")
    void testAdd() {
        assertEquals(5, calculator.add(2, 3), "2 + 3 should equal 5");
    }

    @Test
    void testDivide() {
        assertThrows(ArithmeticException.class, () -> {
            calculator.divide(1, 0);
        });
    }

    // Grouped assertions
    @Test
    void testMultipleAssertions() {
        assertAll(" calculator operations",
            () -> assertEquals(5, calculator.add(2, 3)),
            () -> assertEquals(1, calculator.subtract(3, 2)),
            () -> assertEquals(6, calculator.multiply(2, 3))
        );
    }

    // Parameterized tests
    @ParameterizedTest
    @CsvSource({
        "1, 1, 2",
        "2, 3, 5",
        "10, 20, 30"
    })
    void testAddWithParameters(int a, int b, int expected) {
        assertEquals(expected, calculator.add(a, b));
    }

    @ParameterizedTest
    @ValueSource(strings = {"racecar", "radar", "level"})
    void testIsPalindrome(String candidate) {
        assertTrue(calculator.isPalindrome(candidate));
    }

    // Nested tests
    @Nested
    @DisplayName("When calculator is new")
    class WhenNew {
        @Test
        @DisplayName("display shows zero")
        void displayShowsZero() {
            assertEquals(0, calculator.getDisplay());
        }
    }

    // Mocking with Mockito
    @Test
    void testWithMock() {
        UserRepository mockRepo = mock(UserRepository.class);
        when(mockRepo.findById(1L)).thenReturn(new User(1L, "John"));

        UserService service = new UserService(mockRepo);
        User user = service.getUser(1L);

        assertEquals("John", user.getName());
        verify(mockRepo, times(1)).findById(1L);
    }

    // Assumptions
    @Test
    void testOnlyOnCI() {
        assumeTrue(System.getenv("CI") != null);
        // Test only runs on CI
    }

    // Conditional execution
    @Test
    @EnabledOnOs(OS.LINUX)
    void linuxOnlyTest() {
        // Only runs on Linux
    }

    @RepeatedTest(5)
    void repeatedTest() {
        // Runs 5 times
    }

    @Timeout(2)
    @Test
    void timeoutTest() {
        // Fails if takes more than 2 seconds
    }
}
```

---

## Python: pytest

```python
import pytest
from unittest.mock import Mock, patch, MagicMock

# Basic test
class TestCalculator:
    def test_add(self):
        calc = Calculator()
        assert calc.add(2, 3) == 5

    def test_add_negative(self):
        calc = Calculator()
        assert calc.add(-2, -3) == -5

    # Parametrized test
    @pytest.mark.parametrize("a,b,expected", [
        (2, 3, 5),
        (-2, -3, -5),
        (0, 5, 5),
    ])
    def test_add_parametrized(self, a, b, expected):
        calc = Calculator()
        assert calc.add(a, b) == expected

    # Fixture
    @pytest.fixture
    def calculator(self):
        return Calculator()

    def test_with_fixture(self, calculator):
        assert calculator.add(2, 3) == 5

    # Setup/teardown
    @pytest.fixture
    def db(self):
        db = create_test_db()
        yield db
        db.cleanup()

    # Mocking
    def test_with_mock(self):
        mock_repo = Mock()
        mock_repo.find_by_id.return_value = User(id=1, name="John")

        service = UserService(mock_repo)
        user = service.get_user(1)

        assert user.name == "John"
        mock_repo.find_by_id.assert_called_once_with(1)

    # Patch
    @patch('module.external_api')
    def test_with_patch(self, mock_api):
        mock_api.return_value = {"status": "ok"}
        result = call_external()
        assert result["status"] == "ok"

    # Exception testing
    def test_divide_by_zero(self):
        calc = Calculator()
        with pytest.raises(ZeroDivisionError):
            calc.divide(10, 0)

    # Async test
    @pytest.mark.asyncio
    async def test_async_fetch(self):
        result = await fetch_data()
        assert result is not None

    # Snapshot testing
    def test_snapshot(self, snapshot):
        result = generate_complex_data()
        assert result == snapshot

# Property-based testing with Hypothesis
from hypothesis import given, strategies as st

@given(st.integers(), st.integers())
def test_add_commutative(a, b):
    assert Calculator().add(a, b) == Calculator().add(b, a)

@given(st.lists(st.integers()))
def test_sum_matches_manual(data):
    assert Calculator().sum(data) == sum(data)

# Benchmark with pytest-benchmark
def test_benchmark(benchmark):
    calc = Calculator()
    result = benchmark(calc.add, 100, 200)
    assert result == 300
```

---

## C#: xUnit

```csharp
using Xunit;
using Moq;
using FluentAssertions;

public class CalculatorTests
{
    private readonly Calculator _calculator;

    public CalculatorTests()
    {
        _calculator = new Calculator();
    }

    [Fact]
    public void Add_TwoNumbers_ReturnsSum()
    {
        var result = _calculator.Add(2, 3);
        Assert.Equal(5, result);
    }

    // Using FluentAssertions
    [Fact]
    public void Add_WithFluentAssertions()
    {
        var result = _calculator.Add(2, 3);
        result.Should().Be(5);
    }

    // Theory (parameterized)
    [Theory]
    [InlineData(2, 3, 5)]
    [InlineData(-2, -3, -5)]
    [InlineData(0, 5, 5)]
    public void Add_VariousInputs_ReturnsExpected(int a, int b, int expected)
    {
        var result = _calculator.Add(a, b);
        Assert.Equal(expected, result);
    }

    // MemberData for complex parameters
    public static IEnumerable<object[]> AddData =>
        new List<object[]>
        {
            new object[] { 2, 3, 5 },
            new object[] { -2, -3, -5 },
        };

    [Theory]
    [MemberData(nameof(AddData))]
    public void Add_WithMemberData(int a, int b, int expected)
    {
        Assert.Equal(expected, _calculator.Add(a, b));
    }

    // ClassData
    public class AddTestData : IEnumerable<object[]>
    {
        public IEnumerator<object[]> GetEnumerator()
        {
            yield return new object[] { 2, 3, 5 };
            yield return new object[] { -2, -3, -5 };
        }

        IEnumerator IEnumerable.GetEnumerator() => GetEnumerator();
    }

    [Theory]
    [ClassData(typeof(AddTestData))]
    public void Add_WithClassData(int a, int b, int expected)
    {
        Assert.Equal(expected, _calculator.Add(a, b));
    }

    // Exception testing
    [Fact]
    public void Divide_ByZero_ThrowsException()
    {
        Assert.Throws<DivideByZeroException>(() =>
            _calculator.Divide(10, 0));
    }

    // Async test
    [Fact]
    public async Task FetchDataAsync_ReturnsData()
    {
        var result = await _service.FetchDataAsync();
        Assert.NotNull(result);
    }

    // Mocking with Moq
    [Fact]
    public void GetUser_WithMock_ReturnsUser()
    {
        var mockRepo = new Mock<IUserRepository>();
        mockRepo.Setup(r => r.FindById(1))
            .Returns(new User { Id = 1, Name = "John" });

        var service = new UserService(mockRepo.Object);
        var user = service.GetUser(1);

        Assert.Equal("John", user.Name);
        mockRepo.Verify(r => r.FindById(1), Times.Once);
    }

    // Collection fixture
    [Collection("Database")]
    public class DatabaseTests
    {
        private readonly DatabaseFixture _fixture;

        public DatabaseTests(DatabaseFixture fixture)
        {
            _fixture = fixture;
        }

        [Fact]
        public void TestWithDatabase()
        {
            // Uses shared database fixture
        }
    }

    // Skipping tests
    [Fact(Skip = "Not implemented yet")]
    public void FutureFeature()
    {
    }

    // Conditional fact
    [Fact]
    [Trait("Category", "Integration")]
    public void IntegrationTest()
    {
        // Only runs when explicitly selected
    }
}

// Collection fixture
[CollectionDefinition("Database")]
public class DatabaseCollection : ICollectionFixture<DatabaseFixture>
{
}

public class DatabaseFixture : IDisposable
{
    public DatabaseFixture()
    {
        // Setup database
    }

    public void Dispose()
    {
        // Cleanup
    }
}
```

---

## Feature Comparison

| Feature | Go | Rust | Jest | JUnit 5 | pytest | xUnit |
|---------|-----|------|------|---------|--------|-------|
| Built-in | Yes | Yes | No | No | No | No |
| Table-driven | Native | With rstest | it.each | @ParameterizedTest | @pytest.mark.parametrize | [Theory] |
| Mocking | Interfaces | mockall | jest.mock | Mockito | unittest.mock | Moq |
| Benchmarks | Built-in | criterion | Built-in | JMH | pytest-benchmark | BenchmarkDotNet |
| Coverage | Built-in | cargo-tarpaulin | Built-in | JaCoCo | pytest-cov | Coverlet |
| Parallel | t.Parallel() | Default | Yes | Yes | pytest-xdist | Yes |
| Property | go-fuzz | proptest | fast-check | jqwik | Hypothesis | FsCheck |
| Snapshot | No | insta | Built-in | No | syrupy | Verify |

---

*Document Version: 1.0*
*Last Updated: 2026-04-03*
*Size: ~21KB*
