# flag 包

> **分类**: 开源技术堆栈

---

## 基本用法

```go
import "flag"

func main() {
    // 定义参数
    host := flag.String("host", "localhost", "Server host")
    port := flag.Int("port", 8080, "Server port")
    debug := flag.Bool("debug", false, "Enable debug mode")
    
    // 解析
    flag.Parse()
    
    fmt.Printf("Server: %s:%d\n", *host, *port)
}
```

---

## 运行

```bash
go run main.go -host=0.0.0.0 -port=3000 -debug

# 简写
go run main.go -h 0.0.0.0 -p 3000
```

---

## 自定义类型

```go
type DurationFlag time.Duration

func (d *DurationFlag) String() string {
    return time.Duration(*d).String()
}

func (d *DurationFlag) Set(s string) error {
    duration, err := time.ParseDuration(s)
    if err != nil {
        return err
    }
    *d = DurationFlag(duration)
    return nil
}

var timeout DurationFlag

func init() {
    flag.Var(&timeout, "timeout", "Request timeout")
}
```

---

## 子命令

```go
func main() {
    if len(os.Args) < 2 {
        fmt.Println("Expected 'serve' or 'migrate' subcommand")
        os.Exit(1)
    }
    
    switch os.Args[1] {
    case "serve":
        serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
        port := serveCmd.Int("port", 8080, "Port")
        serveCmd.Parse(os.Args[2:])
        runServer(*port)
        
    case "migrate":
        migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
        direction := migrateCmd.String("direction", "up", "up or down")
        migrateCmd.Parse(os.Args[2:])
        runMigrate(*direction)
    }
}
```

---

## 与 Cobra 对比

| 特性 | flag | Cobra |
|------|------|-------|
| 复杂度 | 简单 | 复杂 |
| 子命令 | 手动 | 内置 |
| 帮助信息 | 基础 | 丰富 |
| 配置集成 | 无 | Viper |

**推荐**: 简单工具用 flag，复杂 CLI 用 Cobra。
