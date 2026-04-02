# 安全编码 (Secure Coding)

> **分类**: 工程与云原生

---

## 输入验证

```go
func ValidateEmail(email string) error {
    if email == "" {
        return errors.New("email required")
    }
    if !regexp.MustCompile(`^[\w.-]+@[\w.-]+\.\w+$`).MatchString(email) {
        return errors.New("invalid format")
    }
    return nil
}
```

---

## SQL 注入防护

```go
// ✅ 使用参数化查询
rows, err := db.Query("SELECT * FROM users WHERE id = ?", userID)

// ❌ 不要拼接 SQL
rows, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID))
```

---

## 密码安全

```go
import "golang.org/x/crypto/bcrypt"

// 哈希
hash, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

// 验证
err := bcrypt.CompareHashAndPassword(hash, password)
```

---

## 敏感信息

```go
// 不要硬编码密码
const dbPassword = "secret123"  // ❌

// 使用环境变量
password := os.Getenv("DB_PASSWORD")  // ✅
```
