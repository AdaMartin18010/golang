# 密钥管理 (Secrets Management)

> **分类**: 工程与云原生
> **标签**: #security #secrets #vault

---

## 环境变量

### 基本使用

```go
import "github.com/joho/godotenv"

// 加载 .env 文件
godotenv.Load()

dbPassword := os.Getenv("DB_PASSWORD")
apiKey := os.Getenv("API_KEY")
```

### 验证必需变量

```go
func requireEnv(key string) string {
    value := os.Getenv(key)
    if value == "" {
        log.Fatalf("required environment variable %s is not set", key)
    }
    return value
}
```

---

## HashiCorp Vault

### 客户端初始化

```go
import "github.com/hashicorp/vault/api"

config := api.DefaultConfig()
config.Address = "http://localhost:8200"

client, err := api.NewClient(config)
if err != nil {
    log.Fatal(err)
}

client.SetToken("your-token")
```

### 读取密钥

```go
// 读取 KV v2 密钥
secret, err := client.KVv2("secret").Get(context.Background(), "myapp/database")
if err != nil {
    log.Fatal(err)
}

password := secret.Data["password"].(string)
```

### 动态凭据

```go
// 获取动态数据库凭据
dbCreds, err := client.Logical().Read("database/creds/my-role")
if err != nil {
    log.Fatal(err)
}

username := dbCreds.Data["username"].(string)
password := dbCreds.Data["password"].(string)
```

---

## Kubernetes Secrets

### 读取 Secret

```go
import "k8s.io/client-go/kubernetes"

clientset, _ := kubernetes.NewForConfig(config)

secret, err := clientset.CoreV1().Secrets("default").Get(ctx, "my-secret", metav1.GetOptions{})
if err != nil {
    log.Fatal(err)
}

password := string(secret.Data["password"])
```

---

## 加密配置

### 使用 sops

```go
// 解密配置文件
import "github.com/mozilla/sops/v3/decrypt"

 decrypted, err := decrypt.File("config.enc.yaml", "yaml")
 if err != nil {
     log.Fatal(err)
 }

 var config Config
 yaml.Unmarshal(decrypted, &config)
```

---

## 最佳实践

1. **不在代码中硬编码密钥**
2. **使用密钥管理服务**
3. **定期轮换密钥**
4. **最小权限原则**
5. **审计密钥访问日志**
