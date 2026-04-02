# 安全通信 (Secure Communication)

> **分类**: 工程与云原生  
> **标签**: #tls #mtls #encryption

---

## TLS 配置

```go
func CreateTLSConfig() *tls.Config {
    return &tls.Config{
        MinVersion: tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        },
        PreferServerCipherSuites: true,
        CurvePreferences: []tls.CurveID{
            tls.X25519,
            tls.CurveP256,
        },
    }
}
```

---

## 证书验证

```go
func LoadCertificate(certFile, keyFile string) (tls.Certificate, error) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return tls.Certificate{}, err
    }
    
    // 验证证书
    cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
    if err != nil {
        return tls.Certificate{}, err
    }
    
    // 检查过期
    if time.Now().After(cert.Leaf.NotAfter) {
        return tls.Certificate{}, fmt.Errorf("certificate expired")
    }
    
    return cert, nil
}
```

---

## 双向 TLS

```go
func CreateMutualTLSConfig(caCert []byte) *tls.Config {
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    return &tls.Config{
        ClientCAs:  caCertPool,
        ClientAuth: tls.RequireAndVerifyClientCert,
        MinVersion: tls.VersionTLS12,
    }
}
```

---

## 证书轮换

```go
type CertManager struct {
    certPath string
    keyPath  string
    cert     *tls.Certificate
    mu       sync.RWMutex
}

func (cm *CertManager) StartWatching() {
    watcher, _ := fsnotify.NewWatcher()
    watcher.Add(cm.certPath)
    watcher.Add(cm.keyPath)
    
    go func() {
        for event := range watcher.Events {
            if event.Op&fsnotify.Write == fsnotify.Write {
                log.Println("Certificate changed, reloading...")
                cm.Reload()
            }
        }
    }()
}

func (cm *CertManager) Reload() error {
    cert, err := tls.LoadX509KeyPair(cm.certPath, cm.keyPath)
    if err != nil {
        return err
    }
    
    cm.mu.Lock()
    cm.cert = &cert
    cm.mu.Unlock()
    
    return nil
}

func (cm *CertManager) GetCertificate() *tls.Certificate {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    return cm.cert
}
```
