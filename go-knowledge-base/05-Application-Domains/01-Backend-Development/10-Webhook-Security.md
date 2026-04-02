# Webhook 安全实践

> **分类**: 成熟应用领域  
> **标签**: #webhook #security #signature

---

## 签名验证

### HMAC-SHA256 验证

```go
func VerifyWebhookSignature(payload []byte, signature string, secret string) error {
    // 提取签名算法和值
    parts := strings.SplitN(signature, "=", 2)
    if len(parts) != 2 {
        return errors.New("invalid signature format")
    }
    
    algo, sigValue := parts[0], parts[1]
    if algo != "sha256" {
        return errors.New("unsupported algorithm")
    }
    
    // 计算 HMAC
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(payload)
    expectedSig := hex.EncodeToString(mac.Sum(nil))
    
    // 常量时间比较
    if !hmac.Equal([]byte(sigValue), []byte(expectedSig)) {
        return errors.New("signature mismatch")
    }
    
    return nil
}
```

### 中间件实现

```go
func WebhookAuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        signature := c.GetHeader("X-Webhook-Signature")
        if signature == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "missing signature"})
            return
        }
        
        body, _ := io.ReadAll(c.Request.Body)
        c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
        
        if err := VerifyWebhookSignature(body, signature, secret); err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid signature"})
            return
        }
        
        c.Next()
    }
}
```

---

## 重放攻击防护

```go
func VerifyTimestamp(timestamp string, tolerance time.Duration) error {
    ts, err := strconv.ParseInt(timestamp, 10, 64)
    if err != nil {
        return err
    }
    
    eventTime := time.Unix(ts, 0)
    now := time.Now()
    
    if now.Sub(eventTime) > tolerance {
        return errors.New("timestamp too old")
    }
    
    if eventTime.After(now.Add(time.Minute)) {
        return errors.New("timestamp in future")
    }
    
    return nil
}
```

---

## 幂等性处理

```go
type WebhookProcessor struct {
    processed cache.Cache  // 使用 Redis 等
}

func (p *WebhookProcessor) Process(ctx context.Context, event WebhookEvent) error {
    // 检查是否已处理
    key := fmt.Sprintf("webhook:%s", event.ID)
    if exists, _ := p.processed.Exists(key); exists {
        return nil  // 已处理，直接返回
    }
    
    // 处理事件
    if err := p.handleEvent(event); err != nil {
        return err
    }
    
    // 标记为已处理
    p.processed.Set(key, true, 24*time.Hour)
    
    return nil
}
```

---

## 完整示例

```go
func WebhookHandler(c *gin.Context) {
    // 1. 验证时间戳
    timestamp := c.GetHeader("X-Webhook-Timestamp")
    if err := VerifyTimestamp(timestamp, 5*time.Minute); err != nil {
        c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 2. 验证签名
    signature := c.GetHeader("X-Webhook-Signature")
    body, _ := io.ReadAll(c.Request.Body)
    
    if err := VerifyWebhookSignature(body, signature, webhookSecret); err != nil {
        c.AbortWithStatusJSON(401, gin.H{"error": "invalid signature"})
        return
    }
    
    // 3. 解析事件
    var event WebhookEvent
    if err := json.Unmarshal(body, &event); err != nil {
        c.AbortWithStatusJSON(400, gin.H{"error": "invalid JSON"})
        return
    }
    
    // 4. 幂等处理
    if err := processor.Process(c.Request.Context(), event); err != nil {
        c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"status": "ok"})
}
```

---

## 安全建议

1. **使用 HTTPS**
2. **验证签名**
3. **检查时间戳**
4. **实现幂等性**
5. **限制请求大小**
6. **使用 IP 白名单**
7. **记录审计日志**
