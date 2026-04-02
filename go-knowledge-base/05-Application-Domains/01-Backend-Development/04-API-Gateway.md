# API 网关

> **分类**: 成熟应用领域

---

## 功能

- 路由
- 认证
- 限流
- 负载均衡

---

## 实现

```go
// 反向代理
type Proxy struct {
    targets []*url.URL
    current uint32
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    idx := atomic.AddUint32(&p.current, 1)
    target := p.targets[idx%uint32(len(p.targets))]
    
    proxy := httputil.NewSingleHostReverseProxy(target)
    proxy.ServeHTTP(w, r)
}
```

---

## 限流

```go
// 令牌桶
func NewLimiter(rate int) *Limiter {
    return &Limiter{
        tokens: make(chan struct{}, rate),
    }
}
```
