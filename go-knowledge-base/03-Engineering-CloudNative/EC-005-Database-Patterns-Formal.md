# EC-005: 数据库访问模式的形式化 (Database Access Patterns: Formalization)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #database #patterns #repository #unit-of-work #caching
> **权威来源**:
>
> - [Patterns of Enterprise Application Architecture](https://martinfowler.com/books/eaa.html) - Martin Fowler (2002)
> - [Database Internals](https://www.oreilly.com/library/view/database-internals/9781492043401/) - Alex Petrov (2019)

---

## 1. 数据访问模式

### 1.1 Repository 模式

**定义 1.1 (Repository)**
$$\text{Repository}: \text{DomainObject} \leftrightarrow \text{Database}$$

**操作**:

- Add, Remove, Get, GetAll, Find

### 1.2 Unit of Work

**定义 1.2 (工作单元)**
$$\text{UoW} = \langle \text{new}, \text{dirty}, \text{deleted} \rangle$$

**提交**:
$$\text{Commit}() = \text{INSERT}(new) \circ \text{UPDATE}(dirty) \circ \text{DELETE}(deleted)$$

---

## 2. 缓存模式

### 2.1 缓存策略

| 模式 | 写操作 | 一致性 | 适用 |
|------|--------|--------|------|
| Cache-Aside | 写 DB | 低 | 读多写少 |
| Write-Through | 写 Cache+DB | 高 | 写重要 |
| Write-Behind | 写 Cache (异步写DB) | 低 | 写性能 |
| Read-Through | 自动加载 | 中 | 通用 |

---

## 3. 多元表征

### 3.1 数据流图

```
Application
    │
    ├──► Repository
    │       │
    │       ├──► Cache
    │       │       │
    │       │       └── Miss? ──► Database
    │       │
    │       └──► Database
    │
    └──► Unit of Work
            │
            └──► Commit (事务)
```

---

**质量评级**: S (15KB)
