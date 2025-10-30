# 🤝 Contributing Guide

> **Thank you for considering contributing to this project!**
> This guide will help you get started with contributing to our Go 1.23+ knowledge system.

**Languages**: [中文](CONTRIBUTING.md) | **English**

---

## 📋 Table of Contents

- [Code of Conduct](#-code-of-conduct)
- [How Can I Contribute?](#-how-can-i-contribute)
- [Development Setup](#-development-setup)
- [Code Standards](#-code-standards)
- [Testing Requirements](#-testing-requirements)
- [Pull Request Process](#-pull-request-process)
- [Documentation Guidelines](#-documentation-guidelines)
- [Community](#-community)

---

## 📜 Code of Conduct

This project adheres to a Code of Conduct. By participating, you are expected to uphold this code. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

**Core Principles**:

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Collaborate openly

---

## 🎯 How Can I Contribute?

### 1. 🐛 Reporting Bugs

**Before submitting**:

- Check existing [issues](https://github.com/AdaMartin18010/golang/issues)
- Verify the bug exists in the latest version
- Collect relevant information

**Bug Report Should Include**:

- Clear, descriptive title
- Steps to reproduce
- Expected vs actual behavior
- Go version (`go version`)
- OS and architecture
- Code sample (if applicable)

**Use Template**: [Bug Report Template](.github/ISSUE_TEMPLATE/bug_report.md)

---

### 2. 💡 Suggesting Features

**Before suggesting**:

- Check [existing features](README_EN.md)
- Review [open issues](https://github.com/AdaMartin18010/golang/issues)
- Consider if it fits project scope

**Feature Request Should Include**:

- Clear, descriptive title
- Problem statement
- Proposed solution
- Alternative solutions
- Use cases and examples

**Use Template**: [Feature Request Template](.github/ISSUE_TEMPLATE/feature_request.md)

---

### 3. 📖 Improving Documentation

**Documentation improvements are always welcome!**

**Types of contributions**:

- Fix typos and grammar
- Improve clarity and examples
- Add missing documentation
- Translate to other languages
- Update outdated information

**Documentation locations**:

- Main docs: `docs/`
- Examples: `EXAMPLES_EN.md`
- Guides: `QUICK_START_EN.md`, `README_EN.md`
- API docs: Inline code comments

---

### 4. 💻 Contributing Code

**Types of code contributions**:

- Bug fixes
- New features
- Performance improvements
- Test coverage improvements
- Refactoring

**Before coding**:

1. Open an issue to discuss your idea
2. Wait for maintainer approval
3. Fork the repository
4. Create a feature branch

---

## 🛠️ Development Setup

### Prerequisites

```bash
# Required
go version  # 1.23+ (1.24 or 1.25 recommended)
git version  # 2.0+

# Recommended
golangci-lint version  # Latest
```

### Setup Steps

```bash
# 1. Fork the repository
# Click "Fork" on GitHub

# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/golang.git
cd golang

# 3. Add upstream remote
git remote add upstream https://github.com/AdaMartin18010/golang.git

# 4. Create a branch
git checkout -b feature/amazing-feature

# 5. Verify setup
go mod download
go build ./...
go test ./...
```

### Development Tools

**Required**:

- Go 1.23+ toolchain
- Git

**Recommended**:

- golangci-lint (linting)
- govulncheck (security)
- gosec (security)

**Install Tools** (optional but recommended):

```bash
# golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

---

## 📏 Code Standards

### Go Code Style

**Follow**:

- Official Go [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- Project-specific conventions (see below)

### Formatting

```bash
# Format code
go fmt ./...

# Run vet
go vet ./...

# Run linter (if installed)
golangci-lint run
```

**Required**:

- All code must pass `go fmt`
- All code must pass `go vet`
- Zero linter warnings (if golangci-lint available)

### Naming Conventions

**Packages**:

```go
// Good
package user
package httputil

// Bad
package UserPackage
package http_util
```

**Functions/Methods**:

```go
// Exported: PascalCase
func ProcessData() {}
func (a *Agent) MakeDecision() {}

// Unexported: camelCase
func processInternal() {}
func (a *Agent) calculateScore() {}
```

**Variables**:

```go
// Short-lived: short names
for i, v := range items {}

// Long-lived: descriptive names
var requestTimeout time.Duration
var maxConcurrentWorkers int
```

**Constants**:

```go
// Exported
const MaxRetries = 3
const DefaultTimeout = 5 * time.Second

// Unexported
const bufferSize = 1024
```

### Error Handling

**Always handle errors**:

```go
// Good
result, err := DoSomething()
if err != nil {
    return fmt.Errorf("do something: %w", err)
}

// Bad
result, _ := DoSomething()  // Never ignore errors!
```

**Error messages**:

- Start with lowercase
- No punctuation at end
- Use `%w` for error wrapping
- Provide context

```go
// Good
return fmt.Errorf("failed to connect to database: %w", err)

// Bad
return fmt.Errorf("Error connecting to database.")  // Capitalized, punctuation
```

### Comments

**Package comments**:

```go
// Package agent provides AI agent functionality.
// It implements the base agent, decision engine, and learning engine.
package agent
```

**Function comments**:

```go
// ProcessTask processes a task and returns the result.
// It returns an error if the task cannot be processed.
func ProcessTask(task Task) (Result, error) {
    // ...
}
```

**Exported names must have comments**:

```go
// Agent represents an AI agent.
type Agent interface {
    // Process processes an input and returns an output.
    Process(input Input) (Output, error)
}
```

---

## 🧪 Testing Requirements

### Test Coverage

**Requirements**:

- All new code must have tests
- Aim for 60%+ coverage
- Critical paths must have 80%+ coverage

**Run tests**:

```bash
# All tests
go test -v ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# With race detection
go test -race ./...
```

### Test Structure

**Test file naming**:

```bash
# Implementation
agent.go

# Tests
agent_test.go
```

**Test function naming**:

```go
// Test functions
func TestProcessTask(t *testing.T) {}
func TestProcessTask_WithError(t *testing.T) {}

// Benchmark functions
func BenchmarkProcessTask(b *testing.B) {}

// Example functions
func ExampleAgent_Process() {}
```

### Writing Tests

**Table-driven tests** (preferred):

```go
func TestProcessTask(t *testing.T) {
    tests := []struct {
        name    string
        input   Input
        want    Output
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   Input{ID: "1", Data: "test"},
            want:    Output{ID: "1", Result: "processed"},
            wantErr: false,
        },
        {
            name:    "invalid input",
            input:   Input{},
            want:    Output{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessTask(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessTask() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ProcessTask() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Subtests**:

```go
func TestAgent(t *testing.T) {
    t.Run("process valid input", func(t *testing.T) {
        // Test code
    })

    t.Run("process invalid input", func(t *testing.T) {
        // Test code
    })
}
```

**Benchmarks**:

```go
func BenchmarkProcessTask(b *testing.B) {
    input := Input{ID: "test", Data: "data"}
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        ProcessTask(input)
    }
}
```

---

## 🔄 Pull Request Process

### Before Submitting

**Checklist**:

- [ ] Code follows project standards
- [ ] All tests pass (`go test -v ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] Documentation updated (if needed)
- [ ] CHANGELOG updated (for features/fixes)
- [ ] Commit messages are clear

### Commit Messages

**Format**:

```text
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style (formatting)
- `refactor`: Code refactoring
- `test`: Add/update tests
- `chore`: Maintenance tasks

**Examples**:

```bash
feat(agent): add decision consensus mechanism

- Implement multi-agent consensus
- Add voting system
- Update tests

Closes #123

---

fix(learning): correct reward calculation

The reward was calculated incorrectly for negative feedback.
This commit fixes the calculation logic.

Fixes #456

---

docs(readme): update installation instructions

Update Go version requirement from 1.22 to 1.23.
Add instructions for Windows users.
```

### Creating Pull Request

1. **Push to your fork**:

    ```bash
    git push origin feature/amazing-feature
    ```

2. **Open PR on GitHub**:
   - Go to original repository
   - Click "New Pull Request"
   - Select your fork and branch
   - Fill in PR template

3. **PR Description Should Include**:
   - What changes were made
   - Why changes were needed
   - How to test changes
   - Related issues (if any)
   - Screenshots (if UI changes)

**Use Template**: [Pull Request Template](.github/PULL_REQUEST_TEMPLATE.md)

### Review Process

**What to expect**:

1. Automated checks run (CI/CD)
2. Maintainer reviews code
3. Feedback and requested changes
4. Approval and merge

**Response time**:

- Initial response: 1-3 days
- Review cycles: 2-5 days
- Total time: 1-2 weeks (average)

### After Approval

**Maintainers will**:

- Merge your PR
- Update CHANGELOG
- Close related issues
- Thank you! 🎉

---

## 📚 Documentation Guidelines

### Documentation Structure

```text
docs/
├── INDEX.md                 # Documentation index
├── 01-Go语言基础/            # Go basics
├── 02-Go语言现代化/          # Go 1.23+ features
│   ├── 08-智能化架构集成/     # AI-Agent
│   ├── 10-建立完整测试体系/   # Testing
│   ├── 13-Go-1.23工具链增强/ # Toolchain
│   └── 14-Go-1.23并发和网络/ # Concurrency
└── ...
```

### Writing Documentation

**Markdown format**:

- Use clear headings (H1-H6)
- Include code examples
- Add diagrams (if helpful)
- Link to related docs

**Code examples**:

````markdown
```go
// Always include imports
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```
````

**Best practices**:

- Be clear and concise
- Use examples liberally
- Explain "why", not just "what"
- Keep it up-to-date

### Bilingual Support

**English contributions**:

- Create `_EN.md` version
- Keep structure same as Chinese version
- Update both when making changes

**Example**:

- Chinese: `README.md`
- English: `README_EN.md`

---

## 🌍 Community

### Communication Channels

- **GitHub Issues**: Bug reports, feature requests
- **GitHub Discussions**: General questions, ideas
- **Pull Requests**: Code contributions

### Getting Help

**Before asking**:

1. Search [existing issues](https://github.com/AdaMartin18010/golang/issues)
2. Read [FAQ](FAQ.md)
3. Check [documentation](docs/INDEX.md)

**When asking**:

- Be specific and clear
- Provide context and examples
- Be patient and respectful

### Recognition

**Contributors will be**:

- Listed in README (major contributions)
- Mentioned in CHANGELOG
- Thanked in release notes

---

## 🎓 Learning Resources

### For New Contributors

**Start here**:

1. [Quick Start Guide](QUICK_START_EN.md)
2. [Examples Showcase](EXAMPLES_EN.md)
3. [Documentation Index](docs/INDEX.md)

**Good first issues**:

- Look for `good first issue` label
- Start with documentation
- Fix typos or improve examples

### For Experienced Contributors

**Challenging tasks**:

- Look for `help wanted` label
- New features
- Performance improvements
- Complex refactoring

---

## 📞 Contact

### Maintainers

- Primary: GitHub [@AdaMartin18010](https://github.com/AdaMartin18010)

### Links

- **Repository**: [github.com/AdaMartin18010/golang](https://github.com/AdaMartin18010/golang)
- **Issues**: [github.com/AdaMartin18010/golang/issues](https://github.com/AdaMartin18010/golang/issues)
- **Discussions**: [github.com/AdaMartin18010/golang/discussions](https://github.com/AdaMartin18010/golang/discussions)

---

## 🙏 Thank You

**Your contributions make this project better!**

Every contribution, no matter how small, is valuable:

- 🐛 Bug reports help us improve quality
- 💡 Feature suggestions guide our roadmap
- 📖 Documentation improvements help users
- 💻 Code contributions advance the project

**We appreciate your time and effort!** 🎉

---

<div align="center">

## 🤝 Ready to Contribute?

**Choose your path and get started!**

[🐛 Report Bug](https://github.com/AdaMartin18010/golang/issues/new?template=bug_report.md) • [💡 Suggest Feature](https://github.com/AdaMartin18010/golang/issues/new?template=feature_request.md) • [💻 Submit PR](https://github.com/AdaMartin18010/golang/compare)

---

**Last Updated**: October 19, 2025
**Languages**: [中文](CONTRIBUTING.md) | **English**

Made with ❤️ by Contributors

</div>
