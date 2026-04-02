# DNS 解析

> **分类**: 开源技术堆栈

---

## 基础解析

```go
import "net"

// 解析域名
addrs, err := net.LookupHost("google.com")
// ["142.250.80.46", ...]

// 解析 IP
names, err := net.LookupAddr("8.8.8.8")

// CNAME
CNAME, err := net.LookupCNAME("www.google.com")

// SRV 记录
_, addrs, err := net.LookupSRV("xmpp-server", "tcp", "google.com")
```

---

## 自定义 Resolver

```go
resolver := &net.Resolver{
    PreferGo: true,  // 使用 Go 实现而非系统
    Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
        d := net.Dialer{
            Timeout: 5 * time.Second,
        }
        return d.DialContext(ctx, "udp", "8.8.8.8:53")
    },
}

// 使用自定义解析器
addrs, err := resolver.LookupHost(ctx, "google.com")
```

---

## DNS 缓存

```go
type DNSCache struct {
    cache map[string]*cacheEntry
    mu    sync.RWMutex
    ttl   time.Duration
}

type cacheEntry struct {
    addrs     []string
    expiresAt time.Time
}

func (c *DNSCache) Lookup(ctx context.Context, host string) ([]string, error) {
    c.mu.RLock()
    entry, ok := c.cache[host]
    c.mu.RUnlock()
    
    if ok && time.Now().Before(entry.expiresAt) {
        return entry.addrs, nil
    }
    
    // 查询
    addrs, err := net.DefaultResolver.LookupHost(ctx, host)
    if err != nil {
        return nil, err
    }
    
    // 缓存
    c.mu.Lock()
    c.cache[host] = &cacheEntry{
        addrs:     addrs,
        expiresAt: time.Now().Add(c.ttl),
    }
    c.mu.Unlock()
    
    return addrs, nil
}
```

---

## 连接池优化

```go
transport := &http.Transport{
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
    }).DialContext,
    ForceAttemptHTTP2:     true,
    MaxIdleConns:          100,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
}

client := &http.Client{Transport: transport}
```
