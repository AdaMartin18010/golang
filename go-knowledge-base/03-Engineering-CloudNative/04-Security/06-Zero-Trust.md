# 零信任架构 (Zero Trust)

> **分类**: 工程与云原生  
> **标签**: #zerotrust #security #mTLS

---

## 核心原则

1. **永不信任，始终验证**
2. **最小权限访问**
3. **假设网络被攻破**
4. **持续验证和监控**

---

## mTLS 实现

### 服务端

```go
func main() {
    // 加载证书
    cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        log.Fatal(err)
    }
    
    // 加载客户端 CA
    caCert, _ := os.ReadFile("ca.crt")
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientCAs:    caCertPool,
        ClientAuth:   tls.RequireAndVerifyClientCert,
    }
    
    server := &http.Server{
        Addr:      ":8443",
        TLSConfig: tlsConfig,
    }
    
    log.Fatal(server.ListenAndServeTLS("", ""))
}
```

### 客户端

```go
func createMTLSClient() *http.Client {
    cert, _ := tls.LoadX509KeyPair("client.crt", "client.key")
    caCert, _ := os.ReadFile("ca.crt")
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
    }
    
    return &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig,
        },
    }
}
```

---

## SPIFFE/SPIRE

```go
import "github.com/spiffe/go-spiffe/v2/workloadapi"

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // 连接到 SPIRE Agent
    source, err := workloadapi.NewX509Source(
        ctx,
        workloadapi.WithClientOptions(
            workloadapi.WithAddr("unix:///tmp/spire-agent/public/api.sock"),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer source.Close()
    
    // 获取 SVID
    svid, err := source.GetX509SVID()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("SPIFFE ID: %s\n", svid.ID)
}
```

---

## 服务到服务认证

```go
type AuthInterceptor struct {
    verifier TokenVerifier
}

func (a *AuthInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // 提取客户端证书
    peer, ok := peer.FromContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "no peer info")
    }
    
    tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "no TLS info")
    }
    
    // 验证证书
    if len(tlsInfo.State.PeerCertificates) == 0 {
        return nil, status.Error(codes.Unauthenticated, "no client certificate")
    }
    
    clientCert := tlsInfo.State.PeerCertificates[0]
    
    // 检查 SPIFFE ID
    spiffeID, err := spiffeid.FromString(clientCert.URIs[0].String())
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid SPIFFE ID")
    }
    
    // 将身份信息添加到上下文
    ctx = WithSPIFFEID(ctx, spiffeID)
    
    return handler(ctx, req)
}
```

---

## 策略执行

```go
type PolicyEngine struct {
    policies []Policy
}

func (p *PolicyEngine) Authorize(ctx context.Context, request Request) error {
    identity := GetIdentity(ctx)
    resource := request.Resource
    action := request.Action
    
    for _, policy := range p.policies {
        if !policy.Matches(identity, resource, action) {
            continue
        }
        
        if policy.Effect == Deny {
            return fmt.Errorf("explicitly denied")
        }
        
        if policy.Effect == Allow {
            return nil
        }
    }
    
    return fmt.Errorf("no matching policy")
}
```
