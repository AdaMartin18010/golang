# TS-DT-008: Go Workspaces (Go 1.18+)

> **维度**: Technology Stack > Development Tools
> **级别**: S (16+ KB)
> **标签**: #go-workspaces #go-modules #multi-module #development
> **权威来源**:
>
> - [Go Workspaces Tutorial](https://go.dev/doc/tutorial/workspaces) - Go team
> - [Workspace Mode](https://go.dev/ref/mod#workspaces) - Go modules reference

---

## 1. Workspace Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Go Workspace Architecture                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Project Structure:                                                          │
│  /myproject/                                                                 │
│  ├── go.work              # Workspace file                                   │
│  ├── go.work.sum          # Workspace checksums                              │
│  ├── api/                 # Module 1                                         │
│  │   ├── go.mod           # module github.com/example/api                    │
│  │   └── api.go                                                            │
│  ├── service/             # Module 2                                         │
│  │   ├── go.mod           # module github.com/example/service                │
│  │   └── service.go                                                         │
│  ├── common/              # Module 3 (shared library)                        │
│  │   ├── go.mod           # module github.com/example/common                 │
│  │   └── common.go                                                          │
│  └── client/              # Module 4                                         │
│      ├── go.mod           # module github.com/example/client                 │
│      └── client.go                                                          │
│                                                                              │
│  go.work file:                                                               │
│  go 1.21                                                                     │
│                                                                              │
│  use (                                                                       │
│      ./api                                                                   │
│      ./service                                                               │
│      ./common                                                                │
│      ./client                                                                │
│  )                                                                           │
│                                                                              │
│  replace (                                                                   │
│      github.com/example/api => ./api                                         │
│      github.com/example/service => ./service                                 │
│      github.com/example/common => ./common                                   │
│      github.com/example/client => ./client                                   │
│  )                                                                           │
│                                                                              │
│  Benefits:                                                                   │
│  - Work on multiple modules simultaneously                                   │
│  - Changes in one module immediately visible to others                       │
│  - No need to publish to test changes                                        │
│  - Atomic commits across modules                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Workspace Commands

```bash
# Initialize workspace in current directory
go work init

# Initialize workspace with specific modules
go work init ./api ./service ./common

# Add module to workspace
go work use ./client
go work use ./new-module

# Remove module from workspace
go work edit -dropuse=./old-module

# View workspace status
go work sync

# Build all modules in workspace
go build ./...

# Test all modules in workspace
go test ./...

# List workspace modules
go list -m all

# Tidy workspace
go work sync

# Vendor dependencies for workspace
go work vendor
```

---

## 3. Workspace Use Cases

```
Use Case 1: Multi-Module Development
- Main application depends on library
- Both in same repository
- Changes to library immediately available

Use Case 2: Dependency Override
- Need to test with local fork of dependency
- Override without modifying go.mod

Use Case 3: Large Projects
- Monorepo with multiple services
- Each service is a module
- Shared libraries in same repo
```

---

## 4. Best Practices

```
Workspace Best Practices:
□ Don't commit go.work to version control (usually)
□ Use go.work for local development
□ Each module should have its own go.mod
□ Use replace directives for local paths
□ Keep modules loosely coupled
□ Document workspace setup in README
```
