# 配置管理

> **分类**: 成熟应用领域

---

## Viper

```go
import "github.com/spf13/viper"

viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")

viper.ReadInConfig()

port := viper.GetInt("server.port")
```

---

## 环境变量

```go
viper.BindEnv("server.port", "SERVER_PORT")
port := viper.GetInt("server.port")
```

---

## 配置热加载

```go
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Println("Config file changed:", e.Name)
    // 重新加载配置
})
```

---

## 默认值

```go
viper.SetDefault("server.port", 8080)
viper.SetDefault("log.level", "info")
```
