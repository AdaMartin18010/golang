# A4: 迁移指南

> **层级**: 应用层 (Application)
> **地位**: 从 Go 1.25 到 1.26 的迁移指南
> **依赖**: A1-A3

---

## 迁移步骤

### 1. 环境准备

```bash
# 更新Go版本
go version  # 确保 >= 1.26

# 更新go.mod
go mod edit -go=1.26
go mod tidy
```

### 2. 运行自动化工具

```bash
# 自动现代化代码
go fix ./...

# 查看变更
go fix -diff ./...
```

### 3. 手动优化

```go
// 迁移前
config := &Config{
    Timeout: &[]int{30}[0],  // 丑陋
}

// 迁移后
config := &Config{
    Timeout: new(30),  // 简洁
}
```

---

## 兼容性检查清单

- [ ] 代码编译通过
- [ ] 所有测试通过
- [ ] 基准测试无性能回退
- [ ] 生产配置更新

---

**参考层**: [R1-概念图谱](../reference/R1-concept-graph.md)
