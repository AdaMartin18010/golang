# Contributing to Go Formal Verification

First off, thank you for considering contributing to the Go Formal Verification project! It's people like you that make this project better.

## 🌟 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Coding Guidelines](#coding-guidelines)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Project Structure](#project-structure)
- [Testing Guidelines](#testing-guidelines)
- [Documentation Guidelines](#documentation-guidelines)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### 🐛 Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include as many details as possible:

- Use a clear and descriptive title
- Describe the exact steps to reproduce the problem
- Provide specific examples
- Describe the behavior you observed and what you expected
- Include screenshots if relevant
- Note your environment (OS, Go version, tool version)

Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.yml) when filing issues.

### ✨ Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- Use a clear and descriptive title
- Provide a detailed description of the proposed feature
- Explain why this enhancement would be useful
- List examples of how the feature would be used

Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.yml).

### 📝 Improving Documentation

Documentation improvements are always welcome! This includes:

- Fixing typos or grammatical errors
- Adding examples or clarifications
- Translating documentation
- Writing tutorials or blog posts

### 💻 Contributing Code

#### Good First Issues

Look for issues labeled `good first issue` - these are great for newcomers!

#### Areas for Contribution

1. **Formal Verifier**
   - New verification algorithms
   - Performance improvements
   - Bug fixes
   - Test coverage

2. **Pattern Generator**
   - New concurrency patterns
   - Pattern combinations
   - Code generation improvements
   - Template enhancements

3. **Examples**
   - Real-world use cases
   - Best practices demonstrations
   - Performance benchmarks

4. **Tools**
   - IDE plugins
   - Web UI
   - CLI improvements

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional but recommended)

### Setting Up Your Development Environment

1. **Fork the repository**

   Click the 'Fork' button on GitHub.

2. **Clone your fork**

   ```bash
   git clone https://github.com/YOUR-USERNAME/golang-formal-verification.git
   cd golang-formal-verification
   ```

3. **Add upstream remote**

   ```bash
   git remote add upstream https://github.com/ORIGINAL-OWNER/golang-formal-verification.git
   ```

4. **Install dependencies**

   ```bash
   # For Formal Verifier
   cd tools/formal-verifier
   go mod download

   # For Pattern Generator
   cd ../concurrency-pattern-generator
   go mod download
   ```

5. **Build the tools**

   ```bash
   # Formal Verifier
   cd tools/formal-verifier
   go build ./cmd/fv

   # Pattern Generator
   cd ../concurrency-pattern-generator
   go build ./cmd/cpg
   ```

6. **Run tests**

   ```bash
   # Test all
   go test ./...

   # Test with coverage
   go test -cover ./...

   # Test with race detector
   go test -race ./...
   ```

## Coding Guidelines

### Go Style

We follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go).

### Key Principles

1. **Clarity over Cleverness**
   - Write clear, readable code
   - Add comments for complex logic
   - Use meaningful variable names

2. **Error Handling**
   - Always check errors
   - Provide context in error messages
   - Use `fmt.Errorf` with `%w` for error wrapping

3. **Testing**
   - Write tests for new features
   - Maintain or improve test coverage
   - Use table-driven tests when appropriate

### Code Formatting

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run all checks
make lint  # if Makefile is available
```

### Naming Conventions

- **Files**: Use lowercase with underscores (`data_flow.go`)
- **Packages**: Use lowercase, single word if possible
- **Exported names**: Use `CamelCase`
- **Private names**: Use `camelCase`
- **Constants**: Use `CamelCase` or `ALL_CAPS` for package-level
- **Interfaces**: Use `er` suffix when appropriate (`Reader`, `Writer`)

## Commit Guidelines

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```text
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes

### Examples

```text
feat(verifier): add deadlock detection algorithm

Implement a new deadlock detection algorithm based on
resource allocation graphs. This improves detection
accuracy by 25% compared to the previous approach.

Closes #123
```

```text
fix(generator): correct mutex pattern generation

The previous implementation had a race condition in the
mutex unlock logic. This fix ensures proper synchronization.

Fixes #456
```

### Scope

Use appropriate scope tags:

- `verifier`: Formal Verifier
- `generator`: Pattern Generator
- `docs`: Documentation
- `examples`: Example code
- `ci`: CI/CD
- `test`: Tests

## Pull Request Process

### Before Submitting

1. **Update your branch**

   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run tests**

   ```bash
   go test ./...
   go test -race ./...
   ```

3. **Check formatting**

   ```bash
   go fmt ./...
   golangci-lint run
   ```

4. **Update documentation**
   - Update README if needed
   - Add/update code comments
   - Update CHANGELOG.md

### Submitting the PR

1. **Push to your fork**

   ```bash
   git push origin your-branch-name
   ```

2. **Create Pull Request**
   - Use the PR template
   - Fill in all sections
   - Link related issues

3. **Description Guidelines**
   - Describe what changes you made
   - Explain why you made these changes
   - Include screenshots if UI changes
   - List any breaking changes

### Review Process

1. **Automated Checks**
   - CI tests must pass
   - Code coverage should not decrease
   - Linting must pass

2. **Code Review**
   - At least one approval required
   - Address review comments
   - Update PR as needed

3. **Merging**
   - Squash commits before merging
   - Use meaningful commit message
   - Delete branch after merge

## Project Structure

```text
.
├── docs/                          # Documentation
│   ├── 01-语言基础/              # Core theory
│   └── 04-高级特性/              # Advanced topics
├── tools/
│   ├── formal-verifier/          # Verification tool
│   │   ├── cmd/fv/              # CLI entry
│   │   ├── pkg/                 # Core packages
│   │   │   ├── cfg/            # Control flow
│   │   │   ├── ssa/            # SSA
│   │   │   ├── concurrency/    # Concurrency checks
│   │   │   └── types/          # Type verification
│   │   └── README.md
│   └── concurrency-pattern-generator/  # Pattern generator
│       ├── cmd/cpg/             # CLI entry
│       ├── pkg/
│       │   ├── generator/       # Code generator
│       │   └── patterns/        # Pattern implementations
│       └── README.md
├── examples/                     # Example projects
├── scripts/                      # Utility scripts
├── .github/                      # GitHub config
│   ├── workflows/               # CI/CD
│   └── ISSUE_TEMPLATE/          # Issue templates
├── CONTRIBUTING.md              # This file
├── CODE_OF_CONDUCT.md           # Code of conduct
└── README.md                    # Main README
```

## Testing Guidelines

### Unit Tests

- Place tests in `*_test.go` files
- Use `testing` package
- Follow table-driven test pattern
- Test edge cases and error conditions

```go
func TestWorkerPool(t *testing.T) {
    tests := []struct {
        name     string
        workers  int
        jobs     int
        expected int
    }{
        {"basic", 5, 10, 10},
        {"edge", 0, 10, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

- Test component interactions
- Use realistic test data
- Clean up resources in tests

### Benchmarks

- Add benchmarks for performance-critical code
- Use `testing.B`
- Include in PR description

```go
func BenchmarkVerifier(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Benchmark code
    }
}
```

### Coverage

- Aim for 80%+ coverage
- Focus on critical paths
- Don't sacrifice quality for coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation Guidelines

### Code Comments

```go
// Package cfg implements control flow graph generation.
package cfg

// Node represents a node in the control flow graph.
// Each node corresponds to a statement or expression in the source code.
type Node struct {
    // ID is the unique identifier for this node
    ID int

    // Stmt is the AST node this CFG node represents
    Stmt ast.Stmt
}

// NewGraph creates a new control flow graph from the given function.
// It returns an error if the function body is invalid.
func NewGraph(fn *ast.FuncDecl) (*Graph, error) {
    // Implementation
}
```

### Documentation Files

- Use Markdown format
- Include code examples
- Add diagrams when helpful
- Keep examples up-to-date

### Commit Documentation

- Document user-facing changes in CHANGELOG.md
- Update README for new features
- Add examples for new APIs

## Getting Help

- 💬 **GitHub Discussions**: Ask questions and share ideas
- 🐛 **GitHub Issues**: Report bugs or request features
- 📧 **Email**: <team@go-formal-verification.org>
- 📝 **Documentation**: Check the `/docs` directory

## Recognition

Contributors will be recognized in:

- README.md
- Release notes
- Project documentation

Thank you for contributing! 🎉

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
