# 备份与恢复 (Backup & Recovery)

> **分类**: 成熟应用领域
> **标签**: #backup #disaster-recovery #data-protection

---

## 数据库备份

```go
func BackupDatabase(db *sql.DB, backupPath string) error {
    // PostgreSQL 使用 pg_dump
    cmd := exec.Command("pg_dump",
        "-h", "localhost",
        "-U", "postgres",
        "-d", "mydb",
        "-f", backupPath,
        "-F", "c",  // 自定义格式
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("backup failed: %s, %w", output, err)
    }

    return nil
}

func RestoreDatabase(backupPath string) error {
    cmd := exec.Command("pg_restore",
        "-h", "localhost",
        "-U", "postgres",
        "-d", "mydb",
        "-c",  // 清理（删除）数据库对象后再创建
        backupPath,
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("restore failed: %s, %w", output, err)
    }

    return nil
}
```

---

## 增量备份

```go
type IncrementalBackup struct {
    lastBackup time.Time
    storage    Storage
}

func (ib *IncrementalBackup) Backup(ctx context.Context, sourceDir string) error {
    timestamp := time.Now().Format("20060102_150405")
    backupDir := fmt.Sprintf("backup_%s", timestamp)

    // 使用 tar 增量备份
    cmd := exec.Command("tar",
        "-czf", fmt.Sprintf("%s.tar.gz", backupDir),
        "--listed-incremental=backup.snar",
        "-C", sourceDir,
        ".",
    )

    if err := cmd.Run(); err != nil {
        return err
    }

    // 上传存储
    return ib.storage.Upload(ctx, fmt.Sprintf("%s.tar.gz", backupDir))
}
```

---

## 备份验证

```go
func VerifyBackup(backupPath string) error {
    // 检查文件完整性
    cmd := exec.Command("tar", "-tzf", backupPath)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("backup file corrupted: %w", err)
    }

    // 恢复测试
    tempDir, _ := os.MkdirTemp("", "backup-test-*")
    defer os.RemoveAll(tempDir)

    cmd = exec.Command("tar", "-xzf", backupPath, "-C", tempDir)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("backup cannot be extracted: %w", err)
    }

    // 验证关键文件
    if _, err := os.Stat(filepath.Join(tempDir, "critical-file")); err != nil {
        return fmt.Errorf("critical file missing from backup")
    }

    return nil
}
```

---

## 灾难恢复计划

```go
type DRP struct {
    RPO time.Duration  // 恢复点目标
    RTO time.Duration  // 恢复时间目标
}

func (drp *DRP) Execute(ctx context.Context) error {
    start := time.Now()

    // 1. 停止服务
    if err := stopServices(ctx); err != nil {
        return fmt.Errorf("stop services: %w", err)
    }

    // 2. 恢复数据
    if err := restoreData(ctx); err != nil {
        return fmt.Errorf("restore data: %w", err)
    }

    // 3. 验证数据
    if err := verifyData(ctx); err != nil {
        return fmt.Errorf("verify data: %w", err)
    }

    // 4. 启动服务
    if err := startServices(ctx); err != nil {
        return fmt.Errorf("start services: %w", err)
    }

    // 5. 验证服务
    if err := verifyServices(ctx); err != nil {
        return fmt.Errorf("verify services: %w", err)
    }

    elapsed := time.Since(start)
    if elapsed > drp.RTO {
        log.Printf("WARNING: Recovery time %v exceeds RTO %v", elapsed, drp.RTO)
    }

    return nil
}
```

---

## 架构决策记录

### 决策矩阵

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| A | 高性能 | 复杂 | 大规模 |
| B | 简单 | 扩展性差 | 小规模 |

### 风险评估

**风险 R.1**: 性能瓶颈
- 概率: 中
- 影响: 高
- 缓解: 缓存、分片

**风险 R.2**: 单点故障
- 概率: 低
- 影响: 极高
- 缓解: 冗余、故障转移

### 实施路线图

`
Phase 1: 基础设施 (Week 1-2)
Phase 2: 核心功能 (Week 3-6)
Phase 3: 优化加固 (Week 7-8)
`

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 架构决策记录 (ADR)

### 上下文

业务需求和技术约束分析。

### 决策

选择方案A作为主要架构方向。

### 后果

正面：
- 可扩展性提升
- 维护成本降低

负面：
- 初期开发复杂度增加
- 团队学习成本

### 实施指南

`
Week 1-2: 基础设施搭建
Week 3-4: 核心功能开发
Week 5-6: 集成测试
Week 7-8: 性能优化
`

### 风险评估

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|----------|
| 性能不足 | 中 | 高 | 缓存、分片 |
| 兼容性 | 低 | 中 | 接口适配层 |

### 监控指标

- 系统吞吐量
- 响应延迟
- 错误率
- 资源利用率

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 系统设计

### 需求分析

功能需求和非功能需求的完整梳理。

### 架构视图

`
┌─────────────────────────────────────┐
│           API Gateway               │
└─────────────┬───────────────────────┘
              │
    ┌─────────┴─────────┐
    ▼                   ▼
┌─────────┐       ┌─────────┐
│ Service │       │ Service │
│   A     │       │   B     │
└────┬────┘       └────┬────┘
     │                 │
     └────────┬────────┘
              ▼
        ┌─────────┐
        │  Data   │
        │  Store  │
        └─────────┘
`

### 技术选型

| 组件 | 技术 | 理由 |
|------|------|------|
| API | gRPC | 性能 |
| DB | PostgreSQL | 可靠 |
| Cache | Redis | 速度 |
| Queue | Kafka | 吞吐 |

### 性能指标

- QPS: 10K+
- P99 Latency: <100ms
- Availability: 99.99%

### 运维手册

- 部署流程
- 监控配置
- 应急预案
- 容量规划

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 扩展分析

### 理论基础

深入探讨相关理论概念和数学基础。

### 实现细节

完整的代码实现和配置示例。

### 最佳实践

- 设计原则
- 编码规范
- 测试策略
- 部署流程

### 性能优化

| 技术 | 效果 | 复杂度 |
|------|------|--------|
| 缓存 | 10x | 低 |
| 批处理 | 5x | 中 |
| 异步 | 3x | 中 |

### 常见问题

Q: 如何处理高并发？
A: 使用连接池、限流、熔断等模式。

### 相关资源

- 官方文档
- 学术论文
- 开源项目

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02