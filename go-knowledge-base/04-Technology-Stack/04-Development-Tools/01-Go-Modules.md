# TS-DT-001: Go Modules - Dependency Management

> **维度**: Technology Stack > Development Tools
> **级别**: S (16+ KB)
> **标签**: #go-modules #dependency-management #semver #vendoring
> **权威来源**:
>
> - [Go Modules Reference](https://go.dev/ref/mod) - Go Team
> - [Go Modules Wiki](https://github.com/golang/go/wiki/Modules) - Go Wiki
> - [Semantic Versioning](https://semver.org/) - Semver spec

---

## 1. Go Modules Architecture

### 1.1 Module System Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Go Modules Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Module Resolution Graph:                                                    │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │  myapp (main module)                                                 │   │
│  │  ├── github.com/gin-gonic/gin v1.9.1                                │   │
│  │  │   ├── github.com/bytedance/sonic v1.9.1                          │   │
│  │  │   └── github.com/gin-contrib/sse v0.1.0                          │   │
│  │  ├── github.com/go-redis/redis/v9 v9.0.5                            │   │
│  │  │   └── github.com/cespare/xxhash/v2 v2.2.0                        │   │
│  │  └── github.com/stretchr/testify v1.8.4                             │   │
│  │       ├── github.com/davecgh/go-spew v1.1.1                         │   │
│  │       └── github.com/pmezard/go-difflib v1.0.0                      │   │
│  │                                                                      │   │
│  │  Minimum Version Selection (MVS):                                    │   │
│  │  - Finds minimum versions that satisfy all requirements             │   │
│  │  - Deterministic and reproducible builds                             │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  File Structure:                                                             │
│  myproject/                                                                  │
│  ├── go.mod          # Module definition and dependencies                  │
│  ├── go.sum          # Cryptographic checksums                             │
│  ├── vendor/         # Vendored dependencies (optional)                    │
│  └── internal/       # Private packages                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 go.mod Structure

```go
// Module declaration
module github.com/example/myapp

// Go version required
go 1.21

// Require direct dependencies
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/go-redis/redis/v9 v9.0.5
    github.com/stretchr/testify v1.8.4
)

// Indirect dependencies (dependencies of dependencies)
require (
    github.com/bytedance/sonic v1.9.1 // indirect
    github.com/cespare/xxhash/v2 v2.2.0 // indirect
)

// Replace directive (local development or forks)
replace github.com/example/mypackage => ../mypackage

replace github.com/original/package => github.com/fork/package v1.0.0

// Exclude problematic versions
exclude github.com/problematic/package v1.2.3

// Retract (mark versions as unusable)
retract (
    v1.5.0 // Published accidentally
    v1.4.0 // Contains security vulnerability
)
```

---

## 2. Module Operations

### 2.1 Common Commands

```bash
# Initialize new module
go mod init github.com/user/project

# Download dependencies
go mod download

# Tidy module - add missing, remove unused
go mod tidy

# Verify dependencies
go mod verify

# Vendor dependencies
go mod vendor

# Build with vendor directory
go build -mod=vendor

# List module dependencies
go list -m all
go list -m -versions github.com/gin-gonic/gin
go list -m -json github.com/gin-gonic/gin

# Graph dependencies
go mod graph

# Edit go.mod directly
go mod edit -require=github.com/pkg/errors@v0.9.1
go mod edit -replace=github.com/original=../fork
go mod edit -dropreplace=github.com/original
go mod edit -go=1.21

# Clean module cache
go clean -modcache

# Get specific version
go get github.com/gin-gonic/gin@v1.9.1
go get github.com/gin-gonic/gin@latest
go get github.com/gin-gonic/gin@master  # Latest commit

# Update all dependencies
go get -u ./...
go get -u=patch ./...  # Only patch updates

# Why is this module needed?
go mod why github.com/pkg/errors
go mod why -m github.com/pkg/errors
```

### 2.2 Version Selection

```
Semantic Versioning in Go:

┌────────────────────────────────────────────────────────────────┐
│ Version Format: vMAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]        │
│                                                                │
│ v1.2.3          - Release version                              │
│ v1.2.3-alpha    - Pre-release                                  │
│ v1.2.3-alpha.1  - Pre-release with identifier                  │
│ v0.0.0-20200101000000-abcdef123456 - Pseudo-version            │
│                                                                │
│ Major version compatibility:                                   │
│ - v0.x.x: Unstable, breaking changes allowed                   │
│ - v1.x.x: Stable, no breaking changes                          │
│ - v2.x.x: Breaking changes from v1.x.x                         │
│   → Must use /v2 in module path                                │
│   → github.com/user/pkg/v2                                     │
│                                                                │
└────────────────────────────────────────────────────────────────┘

Pseudo-Version Format:
v0.0.0-yyyymmddhhmmss-abcdefabcdef
  │     │                │
  │     │                └── Commit hash (12 chars)
  │     └── Commit timestamp (UTC)
  └── Version prefix (v0.0.0 or last tag)
```

### 2.3 Minimal Version Selection

```
MVS Algorithm:

Given requirements:
- Main module requires A@v1.1.0
- A@v1.1.0 requires B@v1.2.0
- A@v1.1.0 requires C@v1.3.0
- B@v1.2.0 requires C@v1.4.0

MVS selects:
- A@v1.1.0 (explicit)
- B@v1.2.0 (from A)
- C@v1.4.0 (max of requirements: max(v1.3.0, v1.4.0))

Why minimum?
- Ensures reproducibility
- Avoids unnecessary upgrades
- Only upgrades when required

Upgrading strategy:
- Explicit: go get package@version
- Patch: go get -u=patch
- Minor: go get -u
- Major: go get package/v2@version
```

---

## 3. Workspace Mode (Go 1.18+)

### 3.1 Multi-Module Development

```go
// go.work file
// Placed in parent directory of multiple modules

go 1.21

use (
    ./api
    ./service
    ./common
    ./client
)

replace (
    github.com/example/api => ./api
    github.com/example/service => ./service
    github.com/example/common => ./common
)
```

```bash
# Initialize workspace
go work init ./api ./service ./common

# Add module to workspace
go work use ./new-module

# Edit workspace
go work edit -use=./another-module

# Build in workspace mode
go build ./...

# Run tests across all modules
go test ./...
```

---

## 4. Private Modules

### 4.1 Configuration

```bash
# Configure GOPRIVATE
# Prevents fetching from public proxy

# Single private host
export GOPRIVATE=github.com/mycompany

# Multiple patterns
export GOPRIVATE="github.com/mycompany,*.internal.company.com"

# Disable proxy for all (except public)
export GOPRIVATE="*"
export GOPROXY="proxy.golang.org,direct"

# Configure Git for private repos
# ~/.gitconfig
[url "ssh://git@github.com/"]
    insteadOf = https://github.com/

# Or use .netrc
# ~/.netrc
machine github.com
login username
password token
```

---

## 5. Vendoring

```bash
# Create vendor directory
go mod vendor

# Build with vendor
go build -mod=vendor

# Verify vendor matches go.mod
go mod vendor -v

# In CI/CD
# Set -mod=vendor to ensure reproducible builds
go test -mod=vendor ./...

# Vendor directory structure
vendor/
├── modules.txt          # List of vendored modules
├── github.com/
│   ├── gin-gonic/
│   │   └── gin/
│   └── other/
│       └── module/
└── golang.org/
    └── x/
        └── tools/
```

---

## 6. Best Practices

```
Module Best Practices:
□ Use semantic versioning tags
□ Keep go.mod and go.sum in version control
□ Run go mod tidy before commits
□ Use go.work for multi-module projects
□ Set GOPRIVATE for company code
□ Use -mod=vendor in CI/CD
□ Pin to specific versions, not latest
□ Document breaking changes
□ Test with go test ./... before release
□ Use retract for bad releases
```

---

## 7. Checklist

```
Go Modules Checklist:
□ go.mod properly initialized
□ All dependencies properly versioned
□ go.sum committed to version control
□ GOPRIVATE configured for internal repos
□ vendor directory (if required)
□ No replace directives in production
□ go mod tidy run regularly
□ No unused dependencies
□ Semantic versioning tags pushed
□ go.work for local development (not committed)
```
