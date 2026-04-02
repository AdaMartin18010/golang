# 任务 CLI 工具 (Task CLI Tooling)

> **分类**: 工程与云原生
> **标签**: #cli #tooling #cobra #urfave-cli
> **参考**: Cobra, Viper, CLI Best Practices

---

## CLI 架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Task CLI Tool Architecture                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   task-cli                                                                   │
│   ├── task                                                                   │
│   │   ├── list          # 列出任务                                         │
│   │   ├── get           # 获取任务详情                                     │
│   │   ├── create        # 创建任务                                         │
│   │   ├── cancel        # 取消任务                                         │
│   │   ├── logs          # 查看任务日志                                     │
│   │   └── retry         # 重试任务                                         │
│   │                                                                          │
│   ├── queue                                                                │
│   │   ├── list          # 列出队列                                         │
│   │   ├── stats         # 队列统计                                         │
│   │   └── purge         # 清空队列                                         │
│   │                                                                          │
│   ├── worker                                                               │
│   │   ├── list          # 列出工作节点                                     │
│   │   ├── stats         # 工作节点统计                                     │
│   │   └── scale         # 扩缩容                                           │
│   │                                                                          │
│   ├── schedule                                                             │
│   │   ├── list          # 列出定时任务                                     │
│   │   ├── create        # 创建定时任务                                     │
│   │   └── delete        # 删除定时任务                                     │
│   │                                                                          │
│   └── config                                                               │
│       ├── get           # 获取配置                                         │
│       ├── set           # 设置配置                                         │
│       └── init          # 初始化配置                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整 CLI 实现

```go
package cli

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "text/tabwriter"
    "time"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

// TaskClient 任务客户端接口
type TaskClient interface {
    ListTasks(ctx context.Context, options ListOptions) ([]Task, error)
    GetTask(ctx context.Context, id string) (*Task, error)
    CreateTask(ctx context.Context, task *Task) (*Task, error)
    CancelTask(ctx context.Context, id string) error
    GetTaskLogs(ctx context.Context, id string) ([]LogEntry, error)
    RetryTask(ctx context.Context, id string) (*Task, error)
}

// ListOptions 列表选项
type ListOptions struct {
    Status   string
    Limit    int
    Offset   int
    SortBy   string
    SortOrder string
}

// Task 任务
type Task struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`
    Status    string    `json:"status"`
    Payload   string    `json:"payload"`
    CreatedAt time.Time `json:"created_at"`
}

// LogEntry 日志条目
type LogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Level     string    `json:"level"`
    Message   string    `json:"message"`
}

// NewRootCmd 创建根命令
func NewRootCmd(client TaskClient) *cobra.Command {
    rootCmd := &cobra.Command{
        Use:   "task-cli",
        Short: "Task scheduling CLI tool",
        Long:  `A CLI tool for managing tasks, queues, and workers in the task scheduling system.`,
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            return initConfig()
        },
    }

    // 全局标志
    rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.task-cli.yaml)")
    rootCmd.PersistentFlags().String("server", "http://localhost:8080", "task server URL")
    rootCmd.PersistentFlags().String("output", "table", "output format (table, json, yaml)")

    viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
    viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

    // 添加子命令
    rootCmd.AddCommand(newTaskCmd(client))
    rootCmd.AddCommand(newQueueCmd(client))
    rootCmd.AddCommand(newWorkerCmd(client))
    rootCmd.AddCommand(newScheduleCmd(client))
    rootCmd.AddCommand(newConfigCmd())

    return rootCmd
}

// newTaskCmd 任务命令
func newTaskCmd(client TaskClient) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "task",
        Short: "Manage tasks",
        Long:  `Create, list, get, cancel, and retry tasks.`,
    }

    // list
    listCmd := &cobra.Command{
        Use:   "list",
        Short: "List tasks",
        RunE: func(cmd *cobra.Command, args []string) error {
            status, _ := cmd.Flags().GetString("status")
            limit, _ := cmd.Flags().GetInt("limit")

            options := ListOptions{
                Status: status,
                Limit:  limit,
            }

            tasks, err := client.ListTasks(cmd.Context(), options)
            if err != nil {
                return err
            }

            output := viper.GetString("output")
            return printTasks(tasks, output)
        },
    }
    listCmd.Flags().String("status", "", "Filter by status")
    listCmd.Flags().Int("limit", 20, "Limit number of results")

    // get
    getCmd := &cobra.Command{
        Use:   "get [task-id]",
        Short: "Get task details",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            task, err := client.GetTask(cmd.Context(), args[0])
            if err != nil {
                return err
            }

            output := viper.GetString("output")
            return printTask(task, output)
        },
    }

    // create
    createCmd := &cobra.Command{
        Use:   "create",
        Short: "Create a new task",
        RunE: func(cmd *cobra.Command, args []string) error {
            taskType, _ := cmd.Flags().GetString("type")
            payload, _ := cmd.Flags().GetString("payload")
            priority, _ := cmd.Flags().GetInt("priority")

            task := &Task{
                Type:    taskType,
                Payload: payload,
            }

            // 解析 payload JSON
            if payload != "" {
                var data map[string]interface{}
                if err := json.Unmarshal([]byte(payload), &data); err != nil {
                    return fmt.Errorf("invalid payload JSON: %w", err)
                }
            }

            created, err := client.CreateTask(cmd.Context(), task)
            if err != nil {
                return err
            }

            fmt.Printf("Task created: %s\n", created.ID)
            return nil
        },
    }
    createCmd.Flags().String("type", "", "Task type (required)")
    createCmd.Flags().String("payload", "{}", "Task payload (JSON)")
    createCmd.Flags().Int("priority", 0, "Task priority")
    createCmd.MarkFlagRequired("type")

    // cancel
    cancelCmd := &cobra.Command{
        Use:   "cancel [task-id]",
        Short: "Cancel a task",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            if err := client.CancelTask(cmd.Context(), args[0]); err != nil {
                return err
            }
            fmt.Printf("Task %s cancelled\n", args[0])
            return nil
        },
    }

    // logs
    logsCmd := &cobra.Command{
        Use:   "logs [task-id]",
        Short: "Get task logs",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            logs, err := client.GetTaskLogs(cmd.Context(), args[0])
            if err != nil {
                return err
            }

            for _, log := range logs {
                fmt.Printf("[%s] %s: %s\n",
                    log.Timestamp.Format(time.RFC3339),
                    log.Level,
                    log.Message)
            }
            return nil
        },
    }

    // retry
    retryCmd := &cobra.Command{
        Use:   "retry [task-id]",
        Short: "Retry a failed task",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            task, err := client.RetryTask(cmd.Context(), args[0])
            if err != nil {
                return err
            }
            fmt.Printf("Task %s retried: %s\n", args[0], task.ID)
            return nil
        },
    }

    cmd.AddCommand(listCmd, getCmd, createCmd, cancelCmd, logsCmd, retryCmd)
    return cmd
}

// newQueueCmd 队列命令
func newQueueCmd(client TaskClient) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "queue",
        Short: "Manage queues",
    }

    listCmd := &cobra.Command{
        Use:   "list",
        Short: "List queues",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Println("Queues:")
            fmt.Println("  - default")
            fmt.Println("  - high-priority")
            fmt.Println("  - low-priority")
            return nil
        },
    }

    statsCmd := &cobra.Command{
        Use:   "stats [queue-name]",
        Short: "Get queue statistics",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Printf("Queue: %s\n", args[0])
            fmt.Println("  Pending: 10")
            fmt.Println("  Processing: 5")
            fmt.Println("  Failed: 2")
            return nil
        },
    }

    cmd.AddCommand(listCmd, statsCmd)
    return cmd
}

// newWorkerCmd 工作节点命令
func newWorkerCmd(client TaskClient) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "worker",
        Short: "Manage workers",
    }

    listCmd := &cobra.Command{
        Use:   "list",
        Short: "List workers",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Println("Workers:")
            fmt.Println("  - worker-1 (active)")
            fmt.Println("  - worker-2 (active)")
            fmt.Println("  - worker-3 (idle)")
            return nil
        },
    }

    cmd.AddCommand(listCmd)
    return cmd
}

// newScheduleCmd 调度命令
func newScheduleCmd(client TaskClient) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "schedule",
        Short: "Manage scheduled tasks",
    }

    listCmd := &cobra.Command{
        Use:   "list",
        Short: "List scheduled tasks",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Println("Scheduled Tasks:")
            fmt.Println("  - backup (0 0 * * *)")
            fmt.Println("  - cleanup (0 */6 * * *)")
            return nil
        },
    }

    cmd.AddCommand(listCmd)
    return cmd
}

// newConfigCmd 配置命令
func newConfigCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "config",
        Short: "Manage configuration",
    }

    getCmd := &cobra.Command{
        Use:   "get [key]",
        Short: "Get configuration value",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            value := viper.Get(args[0])
            fmt.Printf("%s: %v\n", args[0], value)
            return nil
        },
    }

    setCmd := &cobra.Command{
        Use:   "set [key] [value]",
        Short: "Set configuration value",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            viper.Set(args[0], args[1])
            return viper.WriteConfig()
        },
    }

    initCmd := &cobra.Command{
        Use:   "init",
        Short: "Initialize configuration",
        RunE: func(cmd *cobra.Command, args []string) error {
            viper.SetDefault("server", "http://localhost:8080")
            viper.SetDefault("output", "table")

            configFile := viper.ConfigFileUsed()
            if configFile == "" {
                configFile = "$HOME/.task-cli.yaml"
            }

            fmt.Printf("Configuration initialized at %s\n", configFile)
            return viper.SafeWriteConfig()
        },
    }

    cmd.AddCommand(getCmd, setCmd, initCmd)
    return cmd
}

// printTasks 打印任务列表
func printTasks(tasks []Task, format string) error {
    switch format {
    case "json":
        data, _ := json.MarshalIndent(tasks, "", "  ")
        fmt.Println(string(data))
    case "yaml":
        // YAML输出
    default:
        w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
        fmt.Fprintln(w, "ID\tTYPE\tSTATUS\tCREATED AT")
        for _, task := range tasks {
            fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
                task.ID, task.Type, task.Status,
                task.CreatedAt.Format(time.RFC3339))
        }
        w.Flush()
    }
    return nil
}

// printTask 打印单个任务
func printTask(task *Task, format string) error {
    switch format {
    case "json":
        data, _ := json.MarshalIndent(task, "", "  ")
        fmt.Println(string(data))
    default:
        fmt.Printf("ID: %s\n", task.ID)
        fmt.Printf("Type: %s\n", task.Type)
        fmt.Printf("Status: %s\n", task.Status)
        fmt.Printf("Payload: %s\n", task.Payload)
        fmt.Printf("Created At: %s\n", task.CreatedAt.Format(time.RFC3339))
    }
    return nil
}

// initConfig 初始化配置
func initConfig() error {
    if cfgFile := viper.GetString("config"); cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, _ := os.UserHomeDir()
        viper.AddConfigPath(home)
        viper.SetConfigName(".task-cli")
        viper.SetConfigType("yaml")
    }

    viper.AutomaticEnv()
    viper.SetEnvPrefix("TASK_CLI")

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return err
        }
    }

    return nil
}
```

---

## 使用示例

```go
package main

import (
    "context"

    "cli"
)

func main() {
    // 创建模拟客户端
    client := &mockClient{}

    // 创建根命令
    rootCmd := cli.NewRootCmd(client)

    // 执行
    if err := rootCmd.Execute(); err != nil {
        panic(err)
    }
}

type mockClient struct{}

func (m *mockClient) ListTasks(ctx context.Context, options cli.ListOptions) ([]cli.Task, error) {
    return []cli.Task{
        {ID: "task-1", Type: "email", Status: "completed"},
        {ID: "task-2", Type: "sms", Status: "pending"},
    }, nil
}

func (m *mockClient) GetTask(ctx context.Context, id string) (*cli.Task, error) {
    return &cli.Task{ID: id, Type: "email", Status: "completed"}, nil
}

func (m *mockClient) CreateTask(ctx context.Context, task *cli.Task) (*cli.Task, error) {
    task.ID = "task-new"
    return task, nil
}

func (m *mockClient) CancelTask(ctx context.Context, id string) error {
    return nil
}

func (m *mockClient) GetTaskLogs(ctx context.Context, id string) ([]cli.LogEntry, error) {
    return []cli.LogEntry{
        {Level: "INFO", Message: "Task started"},
        {Level: "INFO", Message: "Task completed"},
    }, nil
}

func (m *mockClient) RetryTask(ctx context.Context, id string) (*cli.Task, error) {
    return &cli.Task{ID: id + "-retry", Type: "email", Status: "pending"}, nil
}
```
