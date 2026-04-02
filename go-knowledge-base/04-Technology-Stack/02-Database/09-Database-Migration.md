# 数据库迁移 (Database Migration)

> **分类**: 开源技术堆栈

---

## golang-migrate

```bash
# 安装
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 创建迁移
migrate create -ext sql -dir migrations -seq create_users_table
```

### 迁移文件

```sql
-- 001_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 001_create_users_table.down.sql
DROP TABLE users;
```

### Go 代码执行

```go
import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

m, err := migrate.New(
    "file://migrations",
    "postgres://user:pass@localhost/db?sslmode=disable",
)
if err != nil {
    log.Fatal(err)
}

// 执行迁移
if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    log.Fatal(err)
}
```

---

## 最佳实践

1. **版本控制**: 迁移文件纳入 Git
2. **不可变**: 已执行的迁移不要修改
3. **可回滚**: 每个 up 对应一个 down
4. **幂等性**: 迁移可重复执行不出错
5. **事务**: 每个迁移在事务中执行

---

## GORM 自动迁移

```go
db.AutoMigrate(&User{}, &Order{})

// 仅创建表，不删除
db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
```
