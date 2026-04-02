# CLI 工具开发

> **分类**: 成熟应用领域

---

## Cobra 框架

```go
import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
    Use:   "mytool",
    Short: "My CLI tool",
    Long:  `A longer description of my tool.`,
}

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the server",
    Run: func(cmd *cobra.Command, args []string) {
        port, _ := cmd.Flags().GetInt("port")
        fmt.Printf("Serving on port %d\n", port)
    },
}

func init() {
    serveCmd.Flags().IntP("port", "p", 8080, "Port to serve on")
    rootCmd.AddCommand(serveCmd)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

---

## 子命令

```go
var userCmd = &cobra.Command{
    Use:   "user",
    Short: "User management",
}

var userListCmd = &cobra.Command{
    Use:   "list",
    Short: "List users",
    Run: func(cmd *cobra.Command, args []string) {
        // 实现
    },
}

func init() {
    userCmd.AddCommand(userListCmd)
    rootCmd.AddCommand(userCmd)
}
```

---

## 配置

```go
import "github.com/spf13/viper"

func init() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

    if err := viper.ReadInConfig(); err != nil {
        log.Fatal(err)
    }
}
```
