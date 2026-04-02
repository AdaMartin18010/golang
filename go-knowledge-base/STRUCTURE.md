# Directory Structure Guide

> **Version**: 1.0 S-Level
> **Created**: 2026-04-02
> **Status**: Active
> **Applies to**: Entire Knowledge Base

---

## Table of Contents

1. [Structure Overview](#structure-overview)
2. [Root Level Files](#root-level-files)
3. [Dimension Directories](#dimension-directories)
4. [Supporting Directories](#supporting-directories)
5. [File Naming Conventions](#file-naming-conventions)
6. [Navigation Structure](#navigation-structure)
7. [Cross-Reference Structure](#cross-reference-structure)
8. [Maintenance Guidelines](#maintenance-guidelines)

---

## Structure Overview

### Directory Hierarchy

```
go-knowledge-base/
│
├── 📄 Root Level Documentation (11 files)
│   ├── README.md                 # Main overview
│   ├── INDEX.md                  # Complete document index
│   ├── CONTRIBUTING.md           # Contribution guidelines
│   ├── CHANGELOG.md              # Version history
│   ├── GOALS.md                  # Project goals
│   ├── METHODOLOGY.md            # Documentation methodology
│   ├── STRUCTURE.md              # This file
│   ├── QUALITY-STANDARDS.md      # S/A/B/C definitions
│   ├── TEMPLATES.md              # Document templates
│   ├── GLOSSARY.md               # Terminology
│   ├── REFERENCES.md             # Bibliography
│   ├── FAQ.md                    # FAQ
│   ├── VISUAL-TEMPLATES.md       # Visualization standards
│   ├── ARCHITECTURE.md           # System architecture
│   ├── ROADMAP.md                # Development roadmap
│   ├── CROSS-REFERENCES.md       # Cross-reference matrix
│   └── STATUS.md                 # Completion status
│
├── 📁 01-Formal-Theory/          # Dimension 1: Formal Theory
├── 📁 02-Language-Design/        # Dimension 2: Language Design
├── 📁 03-Engineering-CloudNative/ # Dimension 3: Engineering
├── 📁 04-Technology-Stack/       # Dimension 4: Technology Stack
├── 📁 05-Application-Domains/    # Dimension 5: Application Domains
│
├── 📁 indices/                   # Navigation indexes
├── 📁 learning-paths/            # Curated learning paths
├── 📁 examples/                  # Code examples
└── 📁 scripts/                   # Automation tools
```

### Structural Principles

```
┌─────────────────────────────────────────────────────────────────┐
│                  STRUCTURAL PRINCIPLES                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. HIERARCHY     → Logical nesting by abstraction level        │
│  2. CONSISTENCY   → Same structure across similar categories    │
│  3. DISCOVERABILITY → Easy to find related content              │
│  4. SCALABILITY   → Structure accommodates growth               │
│  5. STABILITY     → Core structure changes rarely               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Root Level Files

### Core Documentation Files

| File | Purpose | Size Target | Quality |
|------|---------|-------------|---------|
| `README.md` | Main entry point, overview | >15KB | S-Level |
| `INDEX.md` | Complete document listing | >15KB | S-Level |
| `CONTRIBUTING.md` | How to contribute | >15KB | S-Level |
| `CHANGELOG.md` | Version history | >15KB | S-Level |

### Meta Documentation Files

| File | Purpose | Size Target | Quality |
|------|---------|-------------|---------|
| `GOALS.md` | Project objectives | >15KB | S-Level |
| `METHODOLOGY.md` | Documentation approach | >15KB | S-Level |
| `STRUCTURE.md` | This guide | >15KB | S-Level |
| `QUALITY-STANDARDS.md` | Quality definitions | >15KB | S-Level |
| `TEMPLATES.md` | Document templates | >15KB | S-Level |

### Reference Files

| File | Purpose | Size Target | Quality |
|------|---------|-------------|---------|
| `GLOSSARY.md` | Terminology definitions | >15KB | S-Level |
| `REFERENCES.md` | Bibliography | >15KB | S-Level |
| `FAQ.md` | Common questions | >15KB | S-Level |
| `VISUAL-TEMPLATES.md` | Visualization guide | >15KB | S-Level |

### Status & Planning Files

| File | Purpose | Size Target | Quality |
|------|---------|-------------|---------|
| `ARCHITECTURE.md` | System architecture | >10KB | A-Level |
| `ROADMAP.md` | Development plan | >10KB | A-Level |
| `CROSS-REFERENCES.md` | Link matrix | >10KB | A-Level |
| `STATUS.md` | Completion status | >10KB | A-Level |

### Root Level Structure Diagram

```
Root Level Organization:

┌─────────────────────────────────────────────────────────────────┐
│                      ROOT LEVEL FILES                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Entry Points          Meta Docs         Reference              │
│  ├─ README.md          ├─ GOALS.md       ├─ GLOSSARY.md        │
│  ├─ INDEX.md           ├─ METHODOLOGY.md ├─ REFERENCES.md       │
│  ├─ QUICK-START.md     ├─ STRUCTURE.md   ├─ FAQ.md              │
│  ├─ CONTRIBUTING.md    ├─ QUALITY-...    ├─ VISUAL-TEMPLATES.md │
│  └─ CHANGELOG.md       └─ TEMPLATES.md   │                      │
│                                           │                      │
│  Planning                  Status         │                      │
│  ├─ ARCHITECTURE.md        ├─ STATUS.md   │                      │
│  ├─ ROADMAP.md             ├─ COMPLETION-...                    │
│  └─ CROSS-REFERENCES.md    └─ PROGRESS-...                      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Dimension Directories

### Dimension 1: 01-Formal-Theory/

**Purpose**: Mathematical foundations and formal semantics

```
01-Formal-Theory/
├── README.md                      # Dimension overview
│
├── 01-Semantics/                  # Semantics theory
│   ├── README.md
│   ├── 01-Operational-Semantics.md
│   ├── 02-Denotational-Semantics.md
│   ├── 03-Axiomatic-Semantics.md
│   └── 04-Featherweight-Go.md
│
├── 02-Type-Theory/                # Type systems
│   ├── README.md
│   ├── 01-Structural-Typing.md
│   ├── 02-Interface-Types.md
│   ├── 03-Generics-Theory/
│   │   ├── README.md
│   │   ├── 01-F-Bounded-Polymorphism.md
│   │   └── 02-Type-Sets.md
│   └── 04-Subtyping.md
│
├── 03-Concurrency-Models/         # Concurrency theory
│   ├── README.md
│   ├── 01-CSP-Theory.md
│   └── 02-Go-Concurrency-Semantics.md
│
├── 04-Memory-Models/              # Memory models
│   ├── README.md
│   ├── 01-Happens-Before.md
│   └── 02-DRF-SC.md
│
├── 05-Category-Theory/            # Category theory
│   ├── README.md
│   └── 01-Functors.md
│
└── FT-001 to FT-022               # Individual topic files
    ├── FT-001-Go-Memory-Model-Formal-Specification.md
    ├── FT-002-Raft-Consensus-Formal.md
    └── ... (22 files total)
```

### Dimension 2: 02-Language-Design/

**Purpose**: Go language internals and design

```
02-Language-Design/
├── README.md                      # Dimension overview
│
├── 01-Design-Philosophy/          # Design principles
│   ├── README.md
│   ├── 01-Simplicity.md
│   ├── 02-Composition.md
│   ├── 03-Explicitness.md
│   └── 04-Orthogonality.md
│
├── 02-Language-Features/          # Language mechanisms
│   ├── README.md
│   ├── 01-Type-System.md
│   ├── 02-Interfaces.md
│   ├── 03-Goroutines.md
│   ├── 04-Channels.md
│   └── ... (20 files total)
│
├── 03-Evolution/                  # Language evolution
│   ├── README.md
│   ├── 01-Go1-to-Go115.md
│   ├── 02-Go116-to-Go120.md
│   ├── 03-Go121-to-Go124.md
│   └── 04-Go125-to-Go126.md
│
├── 04-Comparison/                 # Language comparisons
│   ├── README.md
│   ├── vs-Cpp.md
│   ├── vs-Java.md
│   └── vs-Rust.md
│
└── LD-001 to LD-015               # Deep dive files
    ├── LD-001-Go-Memory-Model-Formal.md
    └── ... (15 files total)
```

### Dimension 3: 03-Engineering-CloudNative/

**Purpose**: Engineering practices and cloud-native patterns

```
03-Engineering-CloudNative/
├── README.md                      # Dimension overview
│
├── 01-Methodology/                # Engineering methodology
│   ├── README.md
│   ├── 01-Clean-Code.md
│   ├── 02-Design-Patterns.md
│   ├── 03-Testing-Strategies.md
│   └── ... (8 files total)
│
├── 02-Cloud-Native/               # Cloud-native patterns
│   ├── README.md
│   ├── 05-Scheduled-Tasks/        # Task scheduling
│   │   ├── 110-Task-Resource-Quota-Management.md
│   │   ├── 111-Task-Event-Sourcing-Implementation.md
│   │   └── ... (11 files)
│   ├── 07-Graceful-Shutdown.md
│   └── ... (60+ files)
│
├── 03-Performance/                # Performance engineering
│   ├── README.md
│   ├── 01-Profiling.md
│   └── ... (8 files)
│
├── 04-Security/                   # Security practices
│   ├── README.md
│   └── ... (9 files)
│
└── EC-001 to EC-121               # Individual pattern files
    ├── EC-001-Microservices.md
    ├── EC-002-Retry-Pattern.md
    └── ... (121 files total)
```

### Dimension 4: 04-Technology-Stack/

**Purpose**: Technology deep-dives

```
04-Technology-Stack/
├── README.md                      # Dimension overview
│
├── 01-Core-Library/               # Go standard library
│   ├── README.md
│   ├── 01-Standard-Library-Overview.md
│   ├── 02-IO-Package.md
│   └── ... (15 files)
│
├── 02-Database/                   # Database technologies
│   ├── README.md
│   ├── 01-Database-Connectivity.md
│   ├── 02-ORM-GORM.md
│   └── ... (13 files)
│
├── 03-Network/                    # Network technologies
│   ├── README.md
│   ├── 01-Gin-Framework.md
│   └── ... (13 files)
│
├── 04-Development-Tools/          # Development tools
│   ├── README.md
│   └── ... (10 files)
│
└── TS-001 to TS-017               # Technology-specific files
    └── ... (17 files)
```

### Dimension 5: 05-Application-Domains/

**Purpose**: Real-world application patterns

```
05-Application-Domains/
├── README.md                      # Dimension overview
│
├── 01-Backend-Development/        # Backend patterns
│   ├── README.md
│   ├── 01-RESTful-API.md
│   └── ... (14 files)
│
├── 02-Cloud-Infrastructure/       # Infrastructure patterns
│   ├── README.md
│   └── ... (11 files)
│
├── 03-DevOps-Tools/               # DevOps patterns
│   ├── README.md
│   └── ... (14 files)
│
└── AD-001 to AD-016               # Domain-specific files
    └── ... (16 files)
```

### Dimension Structure Comparison

```
┌─────────────────────────────────────────────────────────────────┐
│                DIMENSION STRUCTURE PATTERN                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  All dimensions follow this structure:                          │
│                                                                  │
│  Dimension/                                                      │
│  ├── README.md           ← Dimension overview                   │
│  ├── 01-Category/        ← Organized by topic                   │
│  │   ├── README.md       ← Category overview                    │
│  │   └── 01-Topic.md     ← Individual documents                 │
│  ├── 02-Category/                                               │
│  └── XX-Prefix-###       ← Flat files (legacy/comprehensive)    │
│                                                                  │
│  Prefixes:                                                       │
│  • FT-### = Formal Theory                                        │
│  • LD-### = Language Design                                      │
│  • EC-### = Engineering Cloud-Native                             │
│  • TS-### = Technology Stack                                     │
│  • AD-### = Application Domains                                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Supporting Directories

### indices/ - Navigation Indexes

**Purpose**: Multiple entry points for finding content

```
indices/
├── README.md              # Index overview
├── by-topic.md            # Topic-based organization
└── by-difficulty.md       # Difficulty-based organization
```

**Index Structure**:

```markdown
# by-topic.md structure

## Category

### Subcategory
- [Document Title](../path/to/document.md) - Brief description
- [Document Title](../path/to/document.md) - Brief description

## Category
...
```

### learning-paths/ - Curated Paths

**Purpose**: Guided learning sequences

```
learning-paths/
├── README.md              # Learning paths overview
├── go-specialist.md       # Go deep dive path
├── backend-engineer.md    # Backend career path
├── distributed-systems-engineer.md
└── cloud-native-engineer.md
```

**Path Structure**:

```markdown
# [Role] Learning Path

## Phase 1: Foundations (Weeks 1-4)
1. [Document](../path) - Why this matters
2. [Document](../path) - Key concepts
...

## Phase 2: Core Skills (Weeks 5-12)
...
```

### examples/ - Code Examples

**Purpose**: Runnable, tested code demonstrations

```
examples/
├── README.md              # Examples overview
├── task-scheduler/        # Complete project example
│   ├── README.md
│   ├── go.mod
│   ├── main.go
│   └── ... (full project)
└── saga/                  # Pattern implementation
    ├── README.md
    ├── go.mod
    └── ... (full project)
```

**Example Requirements**:

- Must be complete, runnable projects
- Include go.mod for dependencies
- Have comprehensive README
- Include tests
- Document how to run

### scripts/ - Automation

**Purpose**: Maintenance and validation scripts

```
scripts/
├── README.md              # Scripts documentation
├── validate-links.sh      # Link checking
├── check-size.sh          # Size validation
├── generate-index.sh      # Index generation
└── quality-report.sh      # Quality metrics
```

---

## File Naming Conventions

### General Conventions

```
┌─────────────────────────────────────────────────────────────────┐
│                  NAMING CONVENTIONS                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. Use kebab-case (lowercase with hyphens)                     │
│     ✓ my-document.md                                            │
│     ✗ my_document.md, MyDocument.md, myDocument.md              │
│                                                                  │
│  2. Be descriptive but concise                                  │
│     ✓ context-cancellation-patterns.md                          │
│     ✗ patterns.md, context-cancellation-and-timeout-patterns... │
│                                                                  │
│  3. Use consistent terminology                                  │
│     ✓ goroutine, channel, interface                             │
│     ✗ go-routine, channels, interfaces (plural inconsistency)   │
│                                                                  │
│  4. Include category prefix for flat files                      │
│     ✓ EC-001-Circuit-Breaker-Pattern.md                         │
│     ✗ Circuit-Breaker-Pattern.md                                │
│                                                                  │
│  5. Use numbered prefixes for ordering                          │
│     ✓ 01-Introduction.md, 02-Concepts.md                        │
│     ✗ Introduction.md, Concepts.md                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Prefix Conventions

| Prefix | Dimension | Example |
|--------|-----------|---------|
| `FT-###` | Formal Theory | `FT-001-Go-Memory-Model-Formal.md` |
| `LD-###` | Language Design | `LD-001-Go-Type-System-Formal.md` |
| `EC-###` | Engineering | `EC-001-Microservices-Patterns.md` |
| `TS-###` | Technology Stack | `TS-001-PostgreSQL-Transaction-Internals.md` |
| `AD-###` | Application Domains | `AD-001-DDD-Strategic-Patterns.md` |

### Special Files

| File | Purpose | Location |
|------|---------|----------|
| `README.md` | Directory overview | Every directory |
| `go.mod` | Go module definition | `examples/*/` |
| `_test.go` | Test files | `examples/*/` |

---

## Navigation Structure

### Navigation Hierarchy

```
Navigation Flow:

User Entry
    │
    ├───► README.md (Overview)
    │       │
    │       ├───► INDEX.md (All documents)
    │       │
    │       ├───► QUICK-START.md (Fast start)
    │       │
    │       └───► [Dimension]/README.md
    │               │
    │               ├───► [Category]/README.md
    │               │       │
    │               │       └───► [Document].md
    │               │
    │               └───► [XX-###-Document].md
    │
    ├───► indices/by-topic.md
    │
    ├───► indices/by-difficulty.md
    │
    └───► learning-paths/[path].md
```

### Breadcrumb Pattern

Every document should enable navigation:

```markdown
<!-- At top of document -->
[Home](../README.md) > [Dimension](../01-Formal-Theory/) > [Category](../01-Formal-Theory/01-Semantics/) > Current Document

<!-- In content -->
For background, see [Prerequisite Topic](./prerequisite.md).

<!-- At bottom -->
## Next Steps
- [Related Topic 1](./related-1.md)
- [Related Topic 2](./related-2.md)
```

---

## Cross-Reference Structure

### Cross-Reference Types

| Type | Purpose | Example |
|------|---------|---------|
| **Prerequisite** | Required before reading | "Before reading this, understand [Goroutines](../)" |
| **Related** | Connected concepts | "See also [Channels](../)" |
| **Parent** | Broader category | "Part of [Concurrency Patterns](../)" |
| **Child** | Specific details | "For details, see [Select Statement](../)" |
| **Example** | Practical demonstration | "Example implementation in [examples/task-scheduler/](../)" |

### Cross-Reference Maintenance

```
CROSS-REFERENCES.md structure:

| Document | Prerequisites | Related | Examples | Quality |
|----------|--------------|---------|----------|---------|
| EC-007   | EC-001, LD-002 | EC-008, EC-009 | task-scheduler | S |
| FT-002   | FT-001 | FT-003, EC-008 | - | S |
```

---

## Maintenance Guidelines

### Structural Changes

**When to Modify Structure**:

| Change Type | Approval Required | Migration |
|-------------|-------------------|-----------|
| New category | Maintainer team | Automatic |
| Rename directory | Maintainer team | Redirects |
| Move files | Content lead | Update links |
| New dimension | Steering committee | Major version |

### Directory Growth Management

```
Growth Thresholds:

Category Size        Action
─────────────        ──────
> 50 files      →    Split into subcategories
> 20 subdirs    →    Consider restructuring
> 3 levels deep →    Review hierarchy
```

### Validation Checklist

**For New Directories**:

- [ ] README.md present
- [ ] Parent README updated
- [ ] INDEX.md updated
- [ ] Navigation links work
- [ ] Naming follows conventions
- [ ] Added to CROSS-REFERENCES.md

**For Structural Changes**:

- [ ] All links updated
- [ ] Redirects configured (if needed)
- [ ] CHANGELOG.md updated
- [ ] No broken references

---

## Document History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | 2026-04-02 | Initial S-level structure document | Knowledge Base Team |

---

*For the overall architecture, see [ARCHITECTURE.md](./ARCHITECTURE.md). For navigation, see [INDEX.md](./INDEX.md).*
