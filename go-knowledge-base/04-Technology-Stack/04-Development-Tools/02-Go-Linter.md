# TS-DT-002: Go Linting and Static Analysis

> **维度**: Technology Stack > Development Tools
> **级别**: S (16+ KB)
> **标签**: #golangci-lint #static-analysis #code-quality #linting
> **权威来源**:
>
> - [golangci-lint](https://golangci-lint.run/) - Official docs
> - [Go Vet](https://golang.org/cmd/vet/) - Go standard tool
> - [Static Analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) - Go analysis framework

---

## 1. Go Linting Ecosystem

### 1.1 Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Go Linting Ecosystem                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  golangci-lint (Meta-linter)                                                │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐ │   │
│  │  │   go vet    │  │   errcheck  │  │   staticcheck│  │   revive   │ │   │
│  │  │             │  │             │  │              │  │            │ │   │
│  │  │ - Std tool  │  │ - Unchecked │  │ - Advanced   │  │ - Style    │ │   │
│  │  │ - Built-in  │  │   errors    │  │   analysis   │  │   guide    │ │   │
│  │  │   checks    │  │             │  │ - SA* rules  │  │ - Config   │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘ │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐ │   │
│  │  │   gosimple  │  │   structlint│  │   ineffassign│  │   gocritic │ │   │
│  │  │             │  │             │  │              │  │            │ │   │
│  │  │ - Simplify  │  │ - Struct    │  │ - Detect     │  │ - Opinion  │ │   │
│  │  │   code      │  │   tags      │  │   ineffect.  │  │   ated     │ │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘ │   │
│  │                                                                      │   │
│  │  + 50+ more linters...                                              │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  Configuration: .golangci.yml                                               │
│  Parallel execution for speed                                               │
│  Cache for incremental analysis                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. golangci-lint Configuration

### 2.1 Complete Configuration

```yaml
# .golangci.yml
run:
  # Timeout for analysis
  timeout: 5m

  # Include test files
  tests: true

  # Which files to skip
  skip-files:
    - ".*\\.pb\\.go$"
    - ".*_gen\\.go$"

  # Which dirs to skip
  skip-dirs:
    - vendor
    - third_party

  # Build tags
  build-tags:
    - integration

# Output configuration
output:
  # Format: line-number|json|colored-line-number|tab|checkstyle|code-climate|junit-xml|github-actions
  format: colored-line-number

  # Print lines of code with issue
  print-issued-lines: true

  # Print linter name in the end of issue text
  print-linter-name: true

  # Make issues output unique by line
  uniq-by-line: true

# Linter settings
linters-settings:
  # Errcheck - unchecked errors
  errcheck:
    # Report about not checking of errors in type assertions
    check-type-assertions: true
    # Report about assignment of errors to blank identifier
    check-blank: true
    # List of functions to exclude
    exclude-functions:
      - (*os.File).Close

  # Revive - fast, configurable, extensible linter
  revive:
    ignore-generated-header: true
    severity: warning
    rules:
      - name: blank-imports
      - name: context-as-argument
      - name: context-keys-type
      - name: dot-imports
      - name: error-return
      - name: error-strings
      - name: error-naming
      - name: exported
      - name: if-return
      - name: increment-decrement
      - name: var-naming
      - name: var-declaration
      - name: package-comments
      - name: range
      - name: receiver-naming
      - name: time-naming
      - name: unexported-return
      - name: indent-error-flow
      - name: errorf
      - name: empty-block
      - name: superfluous-else
      - name: unused-parameter
      - name: unreachable-code
      - name: redefines-builtin-id

  # Staticcheck - comprehensive static analysis
  staticcheck:
    # SA*: static analysis checks
    # ST*: style checks
    # S*: simplification checks
    checks: ["all"]

  # Gocritic - style critique
  gocritic:
    enabled-tags:
      - performance
      - style
      - experimental
    disabled-checks:
      - wrapperFunc
      - dupImport

  # Goimports - import formatting
  goimports:
    local-prefixes: github.com/mycompany/myproject

  # Go fmt
  gofmt:
    simplify: true

  # Go vet
  govet:
    check-shadowing: true
    enable-all: true

  # Cyclomatic complexity
  cyclop:
    max-complexity: 15
    package-average: 10

  # Function length
  funlen:
    lines: 100
    statements: 50

  # Line length
  lll:
    line-length: 120

  # Maligned - struct field alignment
  maligned:
    suggest-new: true

  # Nestif - nested if statements
  nestif:
    min-complexity: 5

  # Gosec - security checker
  gosec:
    excludes:
      - G204
    config:
      G306: "0600"
      G101:
        pattern: "(?i)example"

  # Prealloc - slice preallocation
  prealloc:
    simple: true
    range-loops: true
    for-loops: true

# Which linters to enable
linters:
  enable:
    - bodyclose      # Checks whether HTTP response body is closed
    - deadcode       # Finds unused code
    - depguard       # Go linter that checks if package imports are in a list of acceptable packages
    - dogsled        # Checks assignments with too many blank identifiers
    - dupl           # Tool for code clone detection
    - errcheck       # Errcheck is a program for checking for unchecked errors
    - exhaustive     # Check exhaustiveness of enum switch statements
    - funlen         # Tool for detection of long functions
    - gochecknoinits # Checks that no init functions are present in Go code
    - goconst        # Finds repeated strings that could be replaced by a constant
    - gocritic       # Provides diagnostics that check for bugs, performance and style issues
    - gocyclo        # Computes and checks the cyclomatic complexity of functions
    - gofmt          # Gofmt checks whether code was gofmt-ed
    - goimports      # In addition to fixing imports, goimports also formats your code
    - golint         # Golint differs from gofmt. Gofmt reformats Go source code, whereas golint prints out style mistakes
    - gomnd          # An analyzer to detect magic numbers
    - goprintffuncname # Checks that printf-like functions are named with f at the end
    - gosec          # Inspects source code for security problems
    - gosimple       # Linter for Go source code that specializes in simplifying a code
    - govet          # Vet examines Go source code and reports suspicious constructs
    - ineffassign    # Detects when assignments to existing variables are not used
    - interfacer     # Linter that suggests narrower interface types
    - lll            # Reports long lines
    - misspell       # Finds commonly misspelled English words in comments
    - nakedret       # Finds naked returns in functions greater than a specified function length
    - noctx          # Noctx finds sending http request without context.Context
    - nolintlint     # Reports ill-formed or insufficient nolint directives
    - rowserrcheck   # Checks whether Err of rows is checked successfully
    - scopelint      # Scopelint checks for unpinned variables in go programs
    - staticcheck    # Staticcheck is a go vet on steroids, applying a ton of static analysis checks
    - structcheck    # Finds unused struct fields
    - stylecheck     # Stylecheck is a replacement for golint
    - typecheck      # Like the front-end of a Go compiler, parses and type-checks Go code
    - unconvert      # Remove unnecessary type conversions
    - unparam        # Reports unused function parameters
    - unused         # Checks Go code for unused constants, variables, functions and types
    - varcheck       # Finds unused global variables and constants
    - whitespace     # Tool for detection of leading and trailing whitespace
    - wsl            # Whitespace Linter - Forces you to use empty lines!

  disable:
    - maligned  # Deprecated
    - prealloc  # Can be noisy

  # Run only fast linters from enabled linters set (first run won't be fast)
  fast: false

# Issues configuration
issues:
  # List of regexps of issue texts to exclude
  exclude:
    - "Error return value of .((os\\.)?std(out|err)\\..*|.*Close|.*Flush|os\\.Remove(All)?|.*print(f|ln)?|os\\.(Un)?Setenv). is not checked"
    - "exported (type|method|function) (.+) should have comment or be unexported"
    - "ST1000: at least one file in a package should have a package comment"

  # Excluding configuration per-path, per-linter, per-text and per-source
  exclude-rules:
    # Exclude some linters from running on tests files
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec
        - lll

    # Exclude known linter issues
    - text: "weak cryptographic primitive"
      linters:
        - gosec

    # Exclude shadow checking in test files
    - path: _test\.go
      text: "shadow"
      linters:
        - govet

  # Show only new issues: if there are unstaged changes or untracked files,
  # only those changes are analyzed, else only changes in HEAD~ are analyzed.
  new: false

  # Maximum issues count per one linter. Set to 0 to disable.
  max-issues-per-linter: 0

  # Maximum count of issues with the same text. Set to 0 to disable.
  max-same-issues: 0

  # Fix found issues (if it's supported by the linter)
  fix: false
```

---

## 3. Running Linters

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all linters
golangci-lint run

# Run on specific directory
golangci-lint run ./pkg/...

# Run specific linters only
golangci-lint run --no-config --disable-all -E errcheck -E gosimple

# Run with fast linters only
golangci-lint run --fast

# Show all issues (not grouped by file)
golangci-lint run --out-format=line-number

# Generate report
golangci-lint run --out-format=json > report.json
golangci-lint run --out-format=checkstyle > checkstyle.xml

# Fix issues automatically (if supported)
golangci-lint run --fix

# Run in CI (new issues only compared to main)
golangci-lint run --new-from-rev=main

# Cache control
golangci-lint cache status
golangci-lint cache clean
```

---

## 4. Custom Linters

### 4.1 Writing a Custom Analyzer

```go
package customlint

import (
    "go/ast"
    "go/types"

    "golang.org/x/tools/go/analysis"
    "golang.org/x/tools/go/analysis/passes/inspect"
    "golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
    Name:     "nocontext",
    Doc:      "Check for functions that should take context.Context",
    Requires: []*analysis.Analyzer{inspect.Analyzer},
    Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
    inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

    nodeFilter := []ast.Node{
        (*ast.FuncDecl)(nil),
        (*ast.FuncLit)(nil),
    }

    inspect.Preorder(nodeFilter, func(n ast.Node) {
        switch fn := n.(type) {
        case *ast.FuncDecl:
            checkFunction(pass, fn.Name.Name, fn.Type, fn.Body)
        case *ast.FuncLit:
            checkFunction(pass, "<anonymous>", fn.Type, fn.Body)
        }
    })

    return nil, nil
}

func checkFunction(pass *analysis.Pass, name string, typ *ast.FuncType, body *ast.BlockStmt) {
    // Check if function makes network calls but doesn't accept context
    if hasNetworkCall(body) && !hasContextParam(typ) {
        pass.Reportf(typ.Pos(), "function %s makes network calls but does not accept context.Context", name)
    }
}

func hasNetworkCall(body *ast.BlockStmt) bool {
    // Implementation to detect http.Get, sql.Query, etc.
    // ...
    return false
}

func hasContextParam(typ *ast.FuncType) bool {
    if typ.Params == nil {
        return false
    }
    for _, param := range typ.Params.List {
        if sel, ok := param.Type.(*ast.SelectorExpr); ok {
            if ident, ok := sel.X.(*ast.Ident); ok {
                if ident.Name == "context" && sel.Sel.Name == "Context" {
                    return true
                }
            }
        }
    }
    return false
}
```

---

## 5. Checklist

```
Linting Checklist:
□ golangci-lint configured (.golangci.yml)
□ All critical linters enabled (errcheck, govet, staticcheck)
□ Security checks enabled (gosec)
□ Complexity limits set (cyclop, funlen)
□ Pre-commit hooks configured
□ CI pipeline includes linting
□ No linting errors on main branch
□ Regular linting runs (daily/weekly)
□ Team trained on linter rules
□ Documentation for custom rules
```
