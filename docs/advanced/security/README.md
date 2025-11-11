# Go安全实践

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Go安全实践](#go安全实践)
  - [📋 目录](#-目录)
  - [📚 核心内容](#-核心内容)
  - [🚀 JWT认证示例](#-jwt认证示例)
  - [🔒 密码加密](#-密码加密)
  - [📖 系统文档](#-系统文档)

---

## 📚 核心内容

1. **[Web安全基础](./01-Web安全基础.md)** ⭐⭐⭐⭐⭐
2. **[身份认证](./02-身份认证.md)** ⭐⭐⭐⭐⭐
3. **[授权机制](./03-授权机制.md)** ⭐⭐⭐⭐
4. **[数据保护](./04-数据保护.md)** ⭐⭐⭐⭐⭐
5. **[安全审计](./05-安全审计.md)** ⭐⭐⭐⭐
6. **[最佳实践](./06-最佳实践.md)** ⭐⭐⭐⭐⭐

---

## 🚀 JWT认证示例

```go
import "github.com/golang-jwt/jwt/v5"

token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": 123,
    "exp": time.Now().Add(24 * time.Hour).Unix(),
})
tokenString, _ := token.SignedString([]byte("secret"))
```

---

## 🔒 密码加密

```go
import "golang.org/x/crypto/bcrypt"

hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
err := bcrypt.CompareHashAndPassword(hashed, []byte("password"))
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
