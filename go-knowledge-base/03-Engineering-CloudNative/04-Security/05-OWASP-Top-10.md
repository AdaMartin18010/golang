# OWASP Top 10 for Go

> **分类**: 工程与云原生

---

## A01: 访问控制失效

```go
// ❌ 错误：没有权限检查
func GetUserData(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    data := db.GetUser(userID)
    json.NewEncoder(w).Encode(data)
}

// ✅ 正确：验证权限
func GetUserData(w http.ResponseWriter, r *http.Request) {
    currentUser := GetCurrentUser(r)
    targetID := r.URL.Query().Get("id")

    if !currentUser.CanAccess(targetID) {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    data := db.GetUser(targetID)
    json.NewEncoder(w).Encode(data)
}
```

---

## A02: 敏感数据泄露

```go
// ❌ 错误：明文存储密码
func StorePassword(password string) {
    db.Exec("INSERT users (password) VALUES (?)", password)
}

// ✅ 正确：使用 bcrypt
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
```

---

## A03: 注入攻击

```go
// ❌ SQL 注入
query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
db.Query(query)

// ✅ 参数化查询
db.Query("SELECT * FROM users WHERE name = ?", name)
```

---

## A05: 安全配置错误

```go
// ❌ 默认配置不安全
server := &http.Server{
    Addr: ":8080",
}

// ✅ 安全配置
server := &http.Server{
    Addr:         ":8080",
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
    },
}
```

---

## A06: 易受攻击组件

```bash
# 扫描依赖漏洞
govulncheck ./...
snyk test
```

---

## A07: 身份识别与认证失效

```go
// ❌ 弱会话管理
func Login(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
        Name:  "session",
        Value: userID,  // 可预测
    })
}

// ✅ 安全会话
func Login(w http.ResponseWriter, r *http.Request) {
    sessionID := generateSecureSessionID()
    storeSession(sessionID, userID)

    http.SetCookie(w, &http.Cookie{
        Name:     "session",
        Value:    sessionID,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        MaxAge:   3600,
    })
}
```
