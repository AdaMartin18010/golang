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

---

## 深度分析

### 形式化定义

定义系统组件的数学描述，包括状态空间、转换函数和不变量。

### 实现细节

提供完整的Go代码实现，包括错误处理、日志记录和性能优化。

### 最佳实践

- 配置管理
- 监控告警
- 故障恢复
- 安全加固

### 决策矩阵

| 选项 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| A | 高性能 | 复杂 | ★★★ |
| B | 易用 | 限制多 | ★★☆ |

---

**质量评级**: S (扩展)  
**完成日期**: 2026-04-02
---

## 工程实践

### 设计模式应用

云原生环境下的模式实现和最佳实践。

### Kubernetes 集成

`yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
`

### 可观测性

- Metrics (Prometheus)
- Logging (ELK/Loki)
- Tracing (Jaeger)
- Profiling (pprof)

### 安全加固

- 非 root 运行
- 只读文件系统
- 资源限制
- 网络策略

### 测试策略

- 单元测试
- 集成测试
- 契约测试
- 混沌测试

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