# Go数据库开发

Go数据库开发完整指南，涵盖SQL、NoSQL、ORM和数据库设计。

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

- [知识图谱](./00-知识图谱.md)
- [对比矩阵](./00-对比矩阵.md)
- [概念定义体系](./00-概念定义体系.md)

---

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3
