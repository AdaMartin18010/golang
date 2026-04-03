# EC-091-Configuration-Management

> **Dimension**: 03-Engineering-CloudNative
> **Status**: S-Level
> **Created**: 2026-04-03
> **Version**: 2026 (Viper, Consul, etcd, Kubernetes ConfigMaps)
> **Size**: >20KB

---

## 1. 配置管理概述

### 1.1 十二要素应用配置原则

```
1. 配置与代码分离
2. 环境变量优先
3. 配置文件版本控制
4. 敏感信息加密
5. 配置热更新支持
```

### 1.2 配置层级

```
配置优先级 (从高到低):

1. 命令行参数
2. 环境变量
3. 配置文件
4. 默认值
```

---

## 2. Viper配置库

### 2.1 基础使用

```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Log      LogConfig
}

type ServerConfig struct {
    Port         int           `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func Load() (*Config, error) {
    // 设置默认值
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.read_timeout", "15s")
    viper.SetDefault("server.write_timeout", "15s")

    // 读取配置文件
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("/etc/myapp/")

    // 自动读取环境变量
    viper.SetEnvPrefix("MYAPP")
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // 读取配置
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// 热重载
func (c *Config) Watch() {
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Printf("Config file changed: %s", e.Name)

        var newCfg Config
        if err := viper.Unmarshal(&newCfg); err != nil {
            log.Printf("Failed to reload config: %v", err)
            return
        }

        // 原子更新
        atomic.StorePointer(&cfgPtr, unsafe.Pointer(&newCfg))
    })
}
```

### 2.2 多环境配置

```yaml
# config.yaml (默认)
server:
  port: 8080
  read_timeout: 15s
  write_timeout: 15s

database:
  host: localhost
  port: 5432
  name: myapp
  max_connections: 10

# config.production.yaml
server:
  port: 80

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  password: ${DB_PASSWORD}
  max_connections: 100

# config.development.yaml
server:
  port: 3000

database:
  host: localhost
  port: 5432
  name: myapp_dev
```

```go
// 根据环境加载
func LoadForEnv(env string) (*Config, error) {
    viper.SetConfigName(fmt.Sprintf("config.%s", env))

    // 先加载默认配置
    viper.SetConfigName("config")
    _ = viper.ReadInConfig()

    // 加载环境特定配置 (覆盖默认)
    viper.SetConfigName(fmt.Sprintf("config.%s", env))
    _ = viper.MergeInConfig()

    var cfg Config
    viper.Unmarshal(&cfg)
    return &cfg, nil
}
```

---

## 3. 环境变量处理

### 3.1 结构化解析

```go
// 使用envdecode
import "github.com/joeshaw/envdecode"

type Config struct {
    Server struct {
        Port         int           `env:"SERVER_PORT,default=8080"`
        ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT,default=15s"`
    }

    Database struct {
        Host     string `env:"DB_HOST,default=localhost"`
        Port     int    `env:"DB_PORT,default=5432"`
        User     string `env:"DB_USER,required"`
        Password string `env:"DB_PASSWORD,required"`
        Name     string `env:"DB_NAME,required"`
    }

    Features struct {
        EnableCache bool `env:"FEATURE_CACHE,default=true"`
        EnableRateLimit bool `env:"FEATURE_RATE_LIMIT,default=false"`
    }
}

func LoadFromEnv() (*Config, error) {
    var cfg Config
    if err := envdecode.StrictDecode(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

### 3.2 敏感信息管理

```go
// 使用HashiCorp Vault
type VaultProvider struct {
    client *vault.Client
}

func (v *VaultProvider) GetDatabaseCredentials() (username, password string, err error) {
    secret, err := v.client.Logical().Read("database/creds/myapp")
    if err != nil {
        return "", "", err
    }

    username = secret.Data["username"].(string)
    password = secret.Data["password"].(string)

    return username, password, nil
}

// 使用AWS Secrets Manager
func GetSecretFromAWS(secretName string) (map[string]string, error) {
    sess, _ := session.NewSession()
    svc := secretsmanager.New(sess)

    input := &secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretName),
    }

    result, err := svc.GetSecretValue(input)
    if err != nil {
        return nil, err
    }

    var secretData map[string]string
    json.Unmarshal([]byte(*result.SecretString), &secretData)

    return secretData, nil
}
```

---

## 4. 分布式配置中心

### 4.1 etcd集成

```go
import clientv3 "go.etcd.io/etcd/client/v3"

type EtcdConfig struct {
    client *clientv3.Client
    prefix string
}

func NewEtcdConfig(endpoints []string, prefix string) (*EtcdConfig, error) {
    client, err := clientv3.New(clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        return nil, err
    }

    return &EtcdConfig{
        client: client,
        prefix: prefix,
    }, nil
}

func (e *EtcdConfig) Get(key string) (string, error) {
    fullKey := e.prefix + "/" + key
    resp, err := e.client.Get(context.Background(), fullKey)
    if err != nil {
        return "", err
    }

    if len(resp.Kvs) == 0 {
        return "", ErrKeyNotFound
    }

    return string(resp.Kvs[0].Value), nil
}

func (e *EtcdConfig) Watch(key string, callback func(string)) {
    fullKey := e.prefix + "/" + key
    watchChan := e.client.Watch(context.Background(), fullKey)

    go func() {
        for resp := range watchChan {
            for _, ev := range resp.Events {
                if ev.Type == clientv3.EventTypePut {
                    callback(string(ev.Kv.Value))
                }
            }
        }
    }()
}

// 配置初始化
func (e *EtcdConfig) LoadAll() (map[string]string, error) {
    resp, err := e.client.Get(context.Background(), e.prefix, clientv3.WithPrefix())
    if err != nil {
        return nil, err
    }

    configs := make(map[string]string)
    for _, kv := range resp.Kvs {
        key := strings.TrimPrefix(string(kv.Key), e.prefix+"/")
        configs[key] = string(kv.Value)
    }

    return configs, nil
}
```

### 4.2 Consul集成

```go
import "github.com/hashicorp/consul/api"

type ConsulConfig struct {
    client *api.Client
    prefix string
}

func NewConsulConfig(addr, prefix string) (*ConsulConfig, error) {
    config := api.DefaultConfig()
    config.Address = addr

    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }

    return &ConsulConfig{
        client: client,
        prefix: prefix,
    }, nil
}

func (c *ConsulConfig) Get(key string) (string, error) {
    kv := c.client.KV()
    pair, _, err := kv.Get(c.prefix+"/"+key, nil)
    if err != nil {
        return "", err
    }

    if pair == nil {
        return "", ErrKeyNotFound
    }

    return string(pair.Value), nil
}

func (c *ConsulConfig) Watch(key string, callback func(string)) {
    go func() {
        var lastIndex uint64

        for {
            kv := c.client.KV()
            pair, meta, err := kv.Get(c.prefix+"/"+key, &api.QueryOptions{
                WaitIndex: lastIndex,
            })

            if err != nil {
                time.Sleep(5 * time.Second)
                continue
            }

            if meta.LastIndex != lastIndex && pair != nil {
                lastIndex = meta.LastIndex
                callback(string(pair.Value))
            }
        }
    }()
}
```

---

## 5. Kubernetes ConfigMaps和Secrets

### 5.1 配置映射

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  database.host: "postgres"
  database.port: "5432"
  log.level: "info"
  server.yaml: |
    port: 8080
    timeout: 30s
---
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: app
          envFrom:
            - configMapRef:
                name: app-config
          volumeMounts:
            - name: config
              mountPath: /etc/app
      volumes:
        - name: config
          configMap:
            name: app-config
```

```go
// 从Kubernetes加载
func LoadFromK8s() (*Config, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, err
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }

    // 读取ConfigMap
    cm, err := clientset.CoreV1().ConfigMaps("default").Get(context.Background(), "app-config", metav1.GetOptions{})
    if err != nil {
        return nil, err
    }

    cfg := &Config{
        Database: DatabaseConfig{
            Host: cm.Data["database.host"],
            Port: mustAtoi(cm.Data["database.port"]),
        },
    }

    return cfg, nil
}
```

### 5.2 Secret管理

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
stringData:
  db.password: "supersecret"
  api.key: "abc123xyz"
---
# 或使用Sealed Secrets加密
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  name: app-secrets
spec:
  encryptedData:
    db.password: AgByA0...  # 加密值
```

---

## 6. 配置验证

### 6.1 结构化验证

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

type Config struct {
    Server struct {
        Port         int           `validate:"required,min=1,max=65535"`
        ReadTimeout  time.Duration `validate:"required,min=1s,max=5m"`
    }

    Database struct {
        Host     string `validate:"required,hostname"`
        Port     int    `validate:"required,min=1,max=65535"`
        User     string `validate:"required"`
        Password string `validate:"required,min=8"`
        Name     string `validate:"required"`
    }

    RateLimit struct {
        RequestsPerSecond int `validate:"required,min=1,max=10000"`
        BurstSize         int `validate:"required,min=1"`
    }
}

func (c *Config) Validate() error {
    return validate.Struct(c)
}

// 自定义验证器
func init() {
    validate.RegisterValidation("duration_range", func(fl validator.FieldLevel) bool {
        duration := fl.Field().Interface().(time.Duration)
        min := time.Second
        max := time.Hour
        return duration >= min && duration <= max
    })
}
```

### 6.2 交叉验证

```go
func (c *Config) CrossValidate() error {
    // 验证数据库连接
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return fmt.Errorf("invalid database config: %w", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        return fmt.Errorf("cannot connect to database: %w", err)
    }

    // 验证Redis连接
    if c.Redis.Host != "" {
        client := redis.NewClient(&redis.Options{
            Addr: fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
        })
        defer client.Close()

        if err := client.Ping(context.Background()).Err(); err != nil {
            return fmt.Errorf("cannot connect to redis: %w", err)
        }
    }

    return nil
}
```

---

## 7. 最佳实践

### 7.1 配置设计原则

| 原则 | 说明 |
|------|------|
| 最小权限 | 只暴露必要配置 |
| 验证优先 | 启动时验证配置 |
| 文档化 | 所有配置项有说明 |
| 版本化 | 配置变更可追溯 |
| 热更新 | 支持不停机更新 |

### 7.2 安全清单

- [ ] 敏感信息加密存储
- [ ] 生产环境不使用默认配置
- [ ] 配置访问权限控制
- [ ] 配置变更审计日志
- [ ] 定期轮换敏感凭证

---

## 8. 参考文献

1. "The Twelve-Factor App" - Configuration
2. Viper Documentation
3. etcd Documentation
4. Consul Documentation
5. Kubernetes ConfigMaps Guide

---

*Last Updated: 2026-04-03*
