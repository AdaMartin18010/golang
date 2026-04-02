# 内容协商 (Content Negotiation)

> **分类**: 成熟应用领域  
> **标签**: #content-negotiation #api #rest

---

## Accept Header 解析

```go
func ParseAccept(header string) []MediaType {
    var types []MediaType
    
    for _, part := range strings.Split(header, ",") {
        part = strings.TrimSpace(part)
        
        mediaType, q := part, 1.0
        if idx := strings.Index(part, ";"); idx != -1 {
            mediaType = strings.TrimSpace(part[:idx])
            params := part[idx+1:]
            
            // 解析 q 值
            if strings.Contains(params, "q=") {
                qStr := strings.TrimPrefix(params[strings.Index(params, "q="):], "q=")
                q, _ = strconv.ParseFloat(strings.TrimSpace(qStr), 64)
            }
        }
        
        types = append(types, MediaType{
            Type:    mediaType,
            Quality: q,
        })
    }
    
    // 按 q 值排序
    sort.Slice(types, func(i, j int) bool {
        return types[i].Quality > types[j].Quality
    })
    
    return types
}
```

---

## 响应格式协商

```go
func ContentNegotiation() gin.HandlerFunc {
    return func(c *gin.Context) {
        accept := c.GetHeader("Accept")
        if accept == "" {
            accept = "application/json"  // 默认
        }
        
        mediaTypes := ParseAccept(accept)
        
        // 查找最佳匹配
        supported := []string{"application/json", "application/xml", "text/csv"}
        
        for _, mt := range mediaTypes {
            for _, s := range supported {
                if match(mt.Type, s) {
                    c.Set("response_format", s)
                    c.Next()
                    return
                }
            }
        }
        
        c.AbortWithStatus(406)  // Not Acceptable
    }
}

func match(accept, supported string) bool {
    // application/json 匹配 application/*
    if strings.HasSuffix(accept, "/*") {
        prefix := strings.TrimSuffix(accept, "/*")
        return strings.HasPrefix(supported, prefix)
    }
    return accept == supported
}
```

---

## 多格式响应

```go
func RespondWithFormat(c *gin.Context, data interface{}) {
    format := c.GetString("response_format")
    
    switch format {
    case "application/xml":
        c.XML(200, data)
    case "text/csv":
        respondCSV(c, data)
    default:
        c.JSON(200, data)
    }
}

func respondCSV(c *gin.Context, data interface{}) {
    c.Header("Content-Type", "text/csv")
    c.Header("Content-Disposition", "attachment; filename=data.csv")
    
    writer := csv.NewWriter(c.Writer)
    defer writer.Flush()
    
    // 写入数据...
}
```

---

## 语言协商

```go
func LanguageNegotiation(supported []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        acceptLang := c.GetHeader("Accept-Language")
        if acceptLang == "" {
            c.Set("lang", supported[0])
            c.Next()
            return
        }
        
        // 解析语言偏好
        langs := parseLanguages(acceptLang)
        
        // 查找匹配
        for _, l := range langs {
            for _, s := range supported {
                if matchLanguage(l.Tag, s) {
                    c.Set("lang", s)
                    c.Next()
                    return
                }
            }
        }
        
        c.Set("lang", supported[0])
        c.Next()
    }
}

func matchLanguage(accepted, supported string) bool {
    // en-US 匹配 en
    if strings.HasPrefix(accepted, supported) {
        return true
    }
    return accepted == supported
}
```
