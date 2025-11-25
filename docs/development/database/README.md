# Go数据库开发

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go数据库开发](#go数据库开发)
  - [📋 目录](#-目录)
  - [📚 核心内容](#-核心内容)
  - [🚀 GORM示例](#-gorm示例)
  - [📖 系统文档](#-系统文档)

---

## 📚 核心内容

1. **SQL数据库**: MySQL, PostgreSQL
2. **NoSQL数据库**: MongoDB, Redis
3. **ORM框架**: GORM
4. **数据库设计**
5. **性能优化**

---

## 🚀 GORM示例

```go
type User struct {
    ID   uint
    Name string
    Age  int
}

db.Create(&User{Name: "Alice", Age: 25})
db.First(&user, 1)
db.Model(&user).Update("Age", 26)
```

---

## 📖 系统文档
