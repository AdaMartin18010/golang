# 任务 CLI 工具 (Task CLI Tooling)

> **分类**: 工程与云原生  
> **标签**: #cli #tooling #automation

---

## 命令行工具设计

```go
package main

import (
    "github.com/spf13/cobra"
    "github.com/fatih/color"
)

var rootCmd = &cobra.Command{
    Use:   "taskctl",
    Short: "Task system control tool",
}

func main() {
    rootCmd.AddCommand(
        listCmd(),
        submitCmd(),
        statusCmd(),
        cancelCmd(),
        logsCmd(),
        statsCmd(),
        debugCmd(),
    )
    
    rootCmd.Execute()
}

func listCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "list",
        Short: "List tasks",
        Run: func(cmd *cobra.Command, args []string) {
            client := NewTaskClient()
            
            status, _ := cmd.Flags().GetString("status")
            limit, _ := cmd.Flags().GetInt("limit")
            
            tasks, err := client.List(cmd.Context(), ListOptions{
                Status: status,
                Limit:  limit,
            })
            if err != nil {
                color.Red("Error: %v", err)
                return
            }
            
            // 表格输出
            table := NewTable()
            table.AddHeader("ID", "Name", "Status", "Created")
            
            for _, task := range tasks {
                table.AddRow(
                    task.ID[:8],
                    task.Name,
                    coloredStatus(task.Status),
                    task.CreatedAt.Format("2006-01-02 15:04"),
                )
            }
            
            table.Print()
        },
    }
}

func coloredStatus(status string) string {
    switch status {
    case "completed":
        return color.GreenString(status)
    case "failed":
        return color.RedString(status)
    case "running":
        return color.YellowString(status)
    default:
        return status
    }
}
```

---

## 交互式命令

```go
func debugCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "debug <task-id>",
        Short: "Interactive debug mode",
        Run: func(cmd *cobra.Command, args []string) {
            if len(args) < 1 {
                color.Red("Usage: taskctl debug <task-id>")
                return
            }
            
            taskID := args[0]
            client := NewTaskClient()
            
            // 进入交互模式
            reader := bufio.NewReader(os.Stdin)
            
            fmt.Printf("Debug mode for task %s\n", taskID)
            fmt.Println("Commands: step, continue, break, vars, stack, quit")
            
            for {
                fmt.Print("> ")
                input, _ := reader.ReadString('\n')
                input = strings.TrimSpace(input)
                
                switch input {
                case "step":
                    resp, _ := client.DebugStep(cmd.Context(), taskID)
                    fmt.Println("Stepped to:", resp.CurrentStep)
                    
                case "vars":
                    vars, _ := client.DebugGetVars(cmd.Context(), taskID)
                    for k, v := range vars {
                        fmt.Printf("%s = %v\n", k, v)
                    }
                    
                case "stack":
                    stack, _ := client.DebugGetStack(cmd.Context(), taskID)
                    for i, frame := range stack {
                        fmt.Printf("%d: %s:%d %s\n", i, frame.File, frame.Line, frame.Function)
                    }
                    
                case "quit":
                    return
                }
            }
        },
    }
}
```

---

## 批量操作

```go
func bulkCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "bulk <file>",
        Short: "Bulk submit tasks from file",
        Run: func(cmd *cobra.Command, args []string) {
            file := args[0]
            data, _ := os.ReadFile(file)
            
            var tasks []TaskDefinition
            json.Unmarshal(data, &tasks)
            
            client := NewTaskClient()
            bar := progressbar.New(len(tasks))
            
            for _, task := range tasks {
                _, err := client.Submit(cmd.Context(), task)
                if err != nil {
                    color.Yellow("Failed to submit %s: %v", task.Name, err)
                }
                bar.Add(1)
            }
            
            color.Green("Bulk submit completed!")
        },
    }
}
```

---

## 自动补全

```go
func completionCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "completion [bash|zsh|fish]",
        Short: "Generate shell completion",
        Run: func(cmd *cobra.Command, args []string) {
            switch args[0] {
            case "bash":
                cmd.Root().GenBashCompletion(os.Stdout)
            case "zsh":
                cmd.Root().GenZshCompletion(os.Stdout)
            case "fish":
                cmd.Root().GenFishCompletion(os.Stdout, true)
            }
        },
    }
}

// 动态补全
type TaskIDCompletion struct {
    client TaskClient
}

func (tc *TaskIDCompletion) Complete(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    tasks, _ := tc.client.List(cmd.Context(), ListOptions{Limit: 100})
    
    var completions []string
    for _, task := range tasks {
        if strings.HasPrefix(task.ID, toComplete) {
            completions = append(completions, task.ID)
        }
    }
    
    return completions, cobra.ShellCompDirectiveNoFileComp
}
```
