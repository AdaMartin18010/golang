# Web安全基础

**版本**: v1.0
**更新日期**: 2025-11-11
**适用于**: Go 1.25.3

---

## 📋 目录

- [Web安全基础](#web安全基础)
  - [📋 目录](#-目录)
  - [1. 📖 概念介绍](#1--概念介绍)
  - [2. 🎯 常见Web攻击](#2--常见web攻击)
    - [1. XSS (跨站脚本攻击)](#1-xss-跨站脚本攻击)
      - [反射型XSS](#反射型xss)
      - [存储型XSS](#存储型xss)
    - [2. CSRF (跨站请求伪造)](#2-csrf-跨站请求伪造)
    - [3. SQL注入](#3-sql注入)
    - [4. 路径遍历](#4-路径遍历)
  - [🔒 安全头部](#-安全头部)
  - [💡 输入验证](#-输入验证)
  - [⚠️ 敏感信息保护](#️-敏感信息保护)
  - [4. 📚 相关资源](#4--相关资源)

---

## 1. 📖 概念介绍

Web安全是保护Web应用免受各种攻击的实践。了解常见威胁和防御措施是构建安全应用的基础。

---

## 2. 🎯 常见Web攻击

### 1. XSS (跨站脚本攻击)

#### 反射型XSS

```go
// ❌ 不安全：直接输出用户输入
func badHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    fmt.Fprintf(w, "<h1>Hello %s</h1>", name)
    // 攻击: ?name=<script>alert('XSS')</script>
}

// ✅ 安全：HTML转义
import "html/template"

func safeHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    tmpl := template.Must(template.New("").Parse("<h1>Hello {{.}}</h1>"))
    tmpl.Execute(w, name)
}
```

#### 存储型XSS

```go
type Comment struct {
    Content string
}

// ✅ 存储前清理
import "github.com/microcosm-cc/bluemonday"

func sanitizeComment(content string) string {
    p := bluemonday.UGCPolicy()
    return p.Sanitize(content)
}

func saveComment(content string) {
    clean := sanitizeComment(content)
    db.Save(&Comment{Content: clean})
}
```

---

### 2. CSRF (跨站请求伪造)

```go
import (
    "crypto/rand"
    "encoding/base64"
    "github.com/gorilla/csrf"
)

// 生成CSRF令牌
func generateCSRFToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.StdEncoding.EncodeToString(b)
}

// 使用gorilla/csrf中间件
func main() {
    r := mux.NewRouter()

    // CSRF保护
    csrfMiddleware := csrf.Protect(
        []byte("32-byte-long-auth-key"),
        csrf.Secure(false), // 开发环境设为false
    )

    r.HandleFunc("/form", showForm)
    r.HandleFunc("/submit", submitForm).Methods("POST")

    http.ListenAndServe(":8080", csrfMiddleware(r))
}

// 在表单中包含CSRF令牌
func showForm(w http.ResponseWriter, r *http.Request) {
    token := csrf.Token(r)
    fmt.Fprintf(w, `
        <form method="POST" action="/submit">
            <input type="hidden" name="csrf_token" value="%s">
            <button type="submit">Submit</button>
        </form>
    `, token)
}
```

---

### 3. SQL注入

```go
// ❌ 不安全：字符串拼接
func badQuery(username string) {
    query := "SELECT * FROM users WHERE username = '" + username + "'"
    db.Query(query)
    // 攻击: username = "admin' OR '1'='1"
}

// ✅ 安全：预处理语句
func safeQuery(username string) {
    query := "SELECT * FROM users WHERE username = ?"
    db.Query(query, username)
}

// ✅ 使用ORM
type User struct {
    ID       uint
    Username string
}

func safeQueryWithORM(username string) {
    var user User
    db.Where("username = ?", username).First(&user)
}
```

---

### 4. 路径遍历

```go
import (
    "path/filepath"
    "strings"
)

// ❌ 不安全：直接使用用户输入
func badFileHandler(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")
    http.ServeFile(w, r, "/uploads/"+filename)
    // 攻击: ?file=../../etc/passwd
}

// ✅ 安全：路径验证
func safeFileHandler(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")

    // 清理路径
    cleaned := filepath.Clean(filename)

    // 检查是否包含..
    if strings.Contains(cleaned, "..") {
        http.Error(w, "Invalid filename", http.StatusBadRequest)
        return
    }

    // 构建完整路径
    fullPath := filepath.Join("/uploads", cleaned)

    // 验证路径前缀
    if !strings.HasPrefix(fullPath, "/uploads/") {
        http.Error(w, "Invalid path", http.StatusBadRequest)
        return
    }

    http.ServeFile(w, r, fullPath)
}
```

---

## 🔒 安全头部

```go
// 设置安全HTTP头部
func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 防止XSS
        w.Header().Set("X-XSS-Protection", "1; mode=block")

        // 防止点击劫持
        w.Header().Set("X-Frame-Options", "DENY")

        // 防止MIME类型嗅探
        w.Header().Set("X-Content-Type-Options", "nosniff")

        // 强制HTTPS
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

        // Content Security Policy
        w.Header().Set("Content-Security-Policy",
            "default-src 'self'; script-src 'self' 'unsafe-inline'")

        next.ServeHTTP(w, r)
    })
}

// 使用中间件
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handler)

    http.ListenAndServe(":8080", securityHeaders(mux))
}
```

---

## 💡 输入验证

```go
import (
    "regexp"
    "unicode/utf8"
)

// 验证电子邮件
func validateEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}

// 验证长度
func validateLength(s string, min, max int) bool {
    length := utf8.RuneCountInString(s)
    return length >= min && length <= max
}

// 完整的输入验证
type RegisterRequest struct {
    Username string
    Email    string
    Password string
}

func (r *RegisterRequest) Validate() error {
    // 用户名验证
    if !validateLength(r.Username, 3, 20) {
        return errors.New("username must be 3-20 characters")
    }

    // 邮箱验证
    if !validateEmail(r.Email) {
        return errors.New("invalid email format")
    }

    // 密码验证
    if !validateLength(r.Password, 8, 50) {
        return errors.New("password must be 8-50 characters")
    }

    return nil
}
```

---

## ⚠️ 敏感信息保护

```go
import "golang.org/x/crypto/bcrypt"

// ✅ 密码哈希
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// ✅ 密码验证
func checkPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// ❌ 不要在日志中记录敏感信息
func badLogging(username, password string) {
    log.Printf("Login attempt: %s/%s", username, password)
}

// ✅ 只记录必要信息
func goodLogging(username string) {
    log.Printf("Login attempt: %s", username)
}
```

---

## 4. 📚 相关资源
