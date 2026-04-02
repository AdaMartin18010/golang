# ORM - GORM

> **分类**: 开源技术堆栈

---

## 快速开始

```go
import "gorm.io/gorm"
import "gorm.io/driver/mysql"

db, err := gorm.Open(mysql.Open("user:pass@/dbname"), &gorm.Config{})
```

---

## 模型定义

```go
type User struct {
    gorm.Model
    Name     string `gorm:"size:255"`
    Email    string `gorm:"uniqueIndex"`
    Age      int    `gorm:"default:0"`
    Profile  Profile `gorm:"foreignKey:UserID"`
}

type Profile struct {
    gorm.Model
    UserID uint
    Bio    string
}
```

---

## CRUD 操作

### 创建

```go
// 创建
user := User{Name: "Alice", Email: "alice@example.com"}
result := db.Create(&user)

// 批量创建
users := []User{{Name: "A"}, {Name: "B"}}
db.Create(&users)
```

### 查询

```go
// 单条
var user User
db.First(&user, 1)                    // 主键
db.First(&user, "name = ?", "Alice")  // 条件

// 多条
var users []User
db.Where("age > ?", 18).Find(&users)
db.Where("name IN ?", []string{"A", "B"}).Find(&users)
```

### 更新

```go
// 更新单列
db.Model(&user).Update("age", 20)

// 更新多列
db.Model(&user).Updates(User{Name: "Bob", Age: 25})

// 条件更新
db.Model(&User{}).Where("age > ?", 18).Update("status", "adult")
```

### 删除

```go
// 软删除（默认）
db.Delete(&user)

// 硬删除
db.Unscoped().Delete(&user)
```

---

## 关联

```go
// 预加载
db.Preload("Profile").Find(&users)

// Joins
db.Joins("Profile").Find(&users)
```

---

## 迁移

```go
// 自动迁移
db.AutoMigrate(&User{}, &Profile{})

// 创建表
db.Migrator().CreateTable(&User{})
```
