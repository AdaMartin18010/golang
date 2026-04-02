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
