# 安全最佳实践指南

> **版本**: v1.0.0
> **更新日期**: 2025-01-XX
> **状态**: ✅ 生产就绪

---

## 📋 概述

本文档提供了使用本项目安全功能的最佳实践指南，帮助开发者正确配置和使用安全功能。

---

## 🔐 认证和授权

### OAuth2/OIDC

#### 最佳实践

1. **使用 HTTPS**
   ```go
   // 始终在生产环境使用 HTTPS
   config := oauth2.DefaultServerConfig()
   // 配置 TLS
   ```

2. **令牌过期时间**
   ```go
   config := oauth2.DefaultServerConfig()
   config.AccessTokenLifetime = 1 * time.Hour  // 访问令牌：1小时
   config.RefreshTokenLifetime = 30 * 24 * time.Hour  // 刷新令牌：30天
   ```

3. **作用域验证**
   ```go
   // 始终验证请求的作用域
   if !hasRequiredScope(token.Scope, requiredScope) {
       return ErrInsufficientScope
   }
   ```

4. **客户端密钥管理**
   ```go
   // 使用强随机密钥
   // 定期轮换客户端密钥
   // 使用环境变量或密钥管理服务存储密钥
   ```

---

## 🔒 数据加密

### AES-256 加密

#### 最佳实践

1. **密钥管理**
   ```go
   // 使用密钥管理服务（如 HashiCorp Vault）
   // 不要硬编码密钥
   // 定期轮换密钥
   ```

2. **字段级加密**
   ```go
   encryptor, _ := security.NewAES256EncryptorFromString(os.Getenv("ENCRYPTION_KEY"))
   fieldEncryptor := security.NewFieldEncryptor(encryptor)
   
   // 加密敏感字段
   encryptedEmail, _ := fieldEncryptor.EncryptField(user.Email)
   ```

3. **密钥轮换**
   ```go
   // 定期轮换加密密钥
   // 使用密钥版本管理
   // 支持多版本密钥同时解密
   ```

---

## 🛡️ 防护措施

### CSRF 防护

#### 最佳实践

1. **令牌生成**
   ```go
   csrf := security.NewCSRFProtection(security.DefaultCSRFConfig())
   token, _ := csrf.GenerateToken(sessionID)
   ```

2. **令牌验证**
   ```go
   // 在表单中包含 CSRF 令牌
   // 在 AJAX 请求的 Header 中包含令牌
   err := csrf.ValidateToken(sessionID, token)
   ```

3. **配置**
   ```go
   // 设置合理的过期时间（24小时）
   // 使用安全的随机数生成器
   ```

### XSS 防护

#### 最佳实践

1. **输入清理**
   ```go
   xss := security.NewXSSProtection()
   sanitized := xss.Sanitize(userInput)
   ```

2. **输出转义**
   ```go
   // 在模板中自动转义
   // 使用 HTML 转义函数
   escaped := xss.EscapeHTML(userInput)
   ```

3. **内容安全策略**
   ```go
   // 配置 CSP 头部
   config := security.DefaultSecurityHeadersConfig()
   config.CSP = "default-src 'self'; script-src 'self' 'unsafe-inline'"
   ```

### SQL 注入防护

#### 最佳实践

1. **参数化查询**
   ```go
   // 始终使用参数化查询
   db.Query("SELECT * FROM users WHERE id = $1", userID)
   ```

2. **输入验证**
   ```go
   sqlProtection := security.NewSQLInjectionProtection(true)
   err := sqlProtection.ValidateInput(userInput)
   ```

3. **最小权限原则**
   ```go
   // 数据库用户使用最小权限
   // 避免使用超级用户
   ```

---

## 🔑 密码安全

### 密码哈希

#### 最佳实践

1. **使用 Argon2id**
   ```go
   hasher := security.NewPasswordHasher(security.DefaultPasswordHashConfig())
   hash, _ := hasher.Hash(password)
   ```

2. **密码验证**
   ```go
   validator := security.NewPasswordValidator(security.DefaultPasswordValidatorConfig())
   err := validator.Validate(password)
   ```

3. **密码策略**
   ```go
   // 最小长度：8 字符
   // 要求大小写字母、数字
   // 可选：特殊字符
   // 禁止常见弱密码
   ```

---

## 📊 审计日志

### 审计日志记录

#### 最佳实践

1. **记录所有安全事件**
   ```go
   logger := security.NewAuditLogger(store)
   logger.LogSecurity(ctx, userID, "failed_login", map[string]interface{}{
       "attempts": 3,
       "ip": "192.168.1.1",
   })
   ```

2. **日志保留**
   ```go
   // 配置合理的保留时间（90天）
   // 定期归档旧日志
   // 使用加密存储敏感日志
   ```

3. **日志查询**
   ```go
   // 支持按用户、时间范围、操作类型查询
   filter := &security.AuditLogFilter{
       UserID: "user-123",
       StartTime: &startTime,
       EndTime: &endTime,
   }
   logs, _ := logger.QueryLogs(ctx, filter)
   ```

---

## 🚦 速率限制

### 速率限制配置

#### 最佳实践

1. **多级别限制**
   ```go
   // IP 级别：防止 DDoS
   ipLimiter := security.NewIPRateLimiter(security.RateLimiterConfig{
       Limit:  100,
       Window: 1 * time.Minute,
   })
   
   // 用户级别：防止滥用
   userLimiter := security.NewUserRateLimiter(security.RateLimiterConfig{
       Limit:  1000,
       Window: 1 * time.Hour,
   })
   ```

2. **端点级别限制**
   ```go
   // 敏感端点使用更严格的限制
   endpointLimiter := security.NewEndpointRateLimiter(security.RateLimiterConfig{
       Limit:  10,
       Window: 1 * time.Minute,
   })
   ```

---

## 📁 文件上传安全

### 文件上传验证

#### 最佳实践

1. **文件类型验证**
   ```go
   validator := security.NewFileUploadValidator(security.DefaultFileUploadConfig())
   err := validator.ValidateFile(filename, contentType, size, content)
   ```

2. **文件大小限制**
   ```go
   // 设置合理的文件大小限制（10MB）
   // 根据文件类型设置不同限制
   ```

3. **文件内容扫描**
   ```go
   // 验证文件头（魔数）
   // 扫描恶意内容
   // 使用病毒扫描（可选）
   ```

---

## 🔐 会话管理

### 会话安全

#### 最佳实践

1. **会话配置**
   ```go
   config := security.DefaultSessionConfig()
   config.DefaultTTL = 24 * time.Hour
   sm := security.NewSessionManager(config)
   ```

2. **会话安全**
   ```go
   // 使用 HTTPS 传输会话 ID
   // 设置 HttpOnly Cookie
   // 设置 Secure Cookie（仅 HTTPS）
   // 定期刷新会话
   ```

3. **会话过期**
   ```go
   // 设置合理的过期时间
   // 实现自动过期清理
   // 支持手动撤销会话
   ```

---

## 🌐 HTTPS/TLS

### TLS 配置

#### 最佳实践

1. **TLS 版本**
   ```go
   // 最低 TLS 1.2
   // 推荐 TLS 1.3
   config := security.TLSConfig{
       MinVersion: tls.VersionTLS12,
       MaxVersion: tls.VersionTLS13,
   }
   ```

2. **密码套件**
   ```go
   // 使用强密码套件
   // 禁用弱密码套件
   // 优先使用 ECDHE
   ```

3. **证书管理**
   ```go
   // 使用有效的 SSL 证书
   // 定期更新证书
   // 监控证书过期
   ```

---

## 🛡️ 安全头部

### HTTP 安全头部

#### 最佳实践

1. **配置安全头部**
   ```go
   headers := security.NewSecurityHeaders(security.DefaultSecurityHeadersConfig())
   router.Use(headers.Middleware)
   ```

2. **内容安全策略**
   ```go
   // 配置严格的 CSP
   config.CSP = "default-src 'self'; script-src 'self' 'unsafe-inline'"
   ```

3. **HSTS**
   ```go
   // 启用 HSTS（仅 HTTPS）
   config.HSTS = "max-age=31536000; includeSubDomains; preload"
   ```

---

## 🔍 输入验证

### 输入验证和清理

#### 最佳实践

1. **验证所有输入**
   ```go
   validator := security.NewInputValidator(security.InputValidatorConfig{
       MinLength: 5,
       MaxLength: 100,
       Pattern:   "^[a-zA-Z0-9]+$",
   })
   err := validator.Validate(userInput)
   ```

2. **清理用户输入**
   ```go
   sanitizer := security.NewStringSanitizer()
   cleaned := sanitizer.Sanitize(userInput)
   ```

3. **类型验证**
   ```go
   emailValidator := security.NewEmailValidator()
   err := emailValidator.ValidateEmail(email)
   ```

---

## 📝 安全配置

### 配置管理

#### 最佳实践

1. **使用环境变量**
   ```go
   // 敏感配置使用环境变量
   encryptionKey := os.Getenv("ENCRYPTION_KEY")
   ```

2. **配置验证**
   ```go
   config := security.DefaultSecurityConfig()
   if err := config.Validate(); err != nil {
       log.Fatal(err)
   }
   ```

3. **配置热重载**
   ```go
   // 支持配置热重载（无需重启）
   // 验证新配置后再应用
   ```

---

## 🚨 安全事件响应

### 事件处理

#### 最佳实践

1. **监控安全事件**
   ```go
   // 监控失败登录尝试
   // 监控异常访问模式
   // 监控权限提升
   ```

2. **自动响应**
   ```go
   // 自动锁定账户（多次失败登录）
   // 自动触发告警
   // 自动记录审计日志
   ```

3. **告警通知**
   ```go
   // 配置告警规则
   // 发送告警通知
   // 记录告警历史
   ```

---

## 📚 参考资源

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/)

---

**最后更新**: 2025-01-XX

