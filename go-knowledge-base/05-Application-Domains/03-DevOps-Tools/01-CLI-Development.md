# CLI Application Development in Go

> **Dimension**: Application Domains
> **Level**: S (18+ KB)
> **Tags**: #cli #cobra #urfave #command #flag #terminal

---

## 1. Domain Requirements Analysis

### 1.1 CLI vs GUI vs API Comparison

| Aspect | CLI | GUI | Web API |
|--------|-----|-----|---------|
| Automation | Excellent (scripts) | Poor | Good (curl/scripts) |
| Learning Curve | Moderate | Low | High |
| Remote Access | SSH | VNC/RDP | HTTP |
| Batch Operations | Native | Limited | Programmatic |
| Integration | Shell pipes | Clipboard | HTTP clients |
| Output Format | Text/JSON/CSV | Visual | JSON/XML |
| Use Case | DevOps, automation | End users | Services |

### 1.2 CLI Application Types

| Type | Examples | Characteristics |
|------|----------|-----------------|
| System Tools | kubectl, docker, git | Complex hierarchies, plugins |
| Build Tools | make, npm, go | Task-oriented, dependency management |
| Dev Tools | linters, formatters, generators | Single-purpose, pipe-friendly |
| Cloud CLIs | aws, gcloud, az | Service integration, config management |
| Utilities | jq, curl, wget | UNIX philosophy, composable |

---

## 2. Architecture Formalization

### 2.1 CLI Architecture Patterns

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      CLI Application Architecture                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Entry Point                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  func main() {                                                      │    │
│  │      cmd := rootCommand()                                           │    │
│  │      if err := cmd.Execute(); err != nil {                          │    │
│  │          fmt.Fprintln(os.Stderr, err)                               │    │
│  │          os.Exit(1)                                                 │    │
│  │      }                                                              │    │
│  │  }                                                                  │    │
│  └─────────────────────┬───────────────────────────────────────────────┘    │
│                        │                                                     │
│                        ▼                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Command Layer                                    │    │
│  │                                                                     │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │    │
│  │  │  Parse   │  │ Validate │  │  Route   │  │    Execute       │   │    │
│  │  │  Flags   │  │  Input   │  │ Command  │  │                  │   │    │
│  │  │  Args    │──│          │──│          │──│  Call Handler    │   │    │
│  │  │          │  │          │  │          │  │                  │   │    │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                        │                                                     │
│                        ▼                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Handler Layer                                    │    │
│  │                                                                     │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │    │
│  │  │  Business    │  │  External    │  │   Output     │              │    │
│  │  │  Logic       │  │  API Calls   │  │   Formatting │              │    │
│  │  │              │  │              │  │              │              │    │
│  │  │ - Process    │  │ - HTTP       │  │ - Table      │              │    │
│  │  │ - Transform  │  │ - Database   │  │ - JSON       │              │    │
│  │  │ - Validate   │  │ - File       │  │ - YAML       │              │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘              │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                        │                                                     │
│                        ▼                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    Output Layer                                     │    │
│  │                                                                     │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │    │
│  │  │  Stdout  │  │  Stderr  │  │  Exit    │  │  Logs    │            │    │
│  │  │          │  │          │  │  Code    │  │          │            │    │
│  │  │ Success  │  │ Errors   │  │ 0=success│  │ Debug    │            │    │
│  │  │ Data     │  │ Warnings │  │ 1=error  │  │ Info     │            │    │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Command Structure Patterns

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Command Structure Patterns                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Pattern 1: Flat Structure (Simple Tools)                                    │
│  ════════════════════════════════════════                                    │
│                                                                              │
│  myapp [flags] <arg1> <arg2>                                                 │
│                                                                              │
│  Examples: grep, cat, ls                                                     │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Pattern 2: Verb-Noun Structure (AWS Style)                                  │
│  ══════════════════════════════════════════                                  │
│                                                                              │
│  myapp <verb> <noun> [flags]                                                 │
│                                                                              │
│  myapp create bucket                                                         │
│  myapp delete bucket                                                         │
│  myapp list buckets                                                          │
│  myapp get bucket                                                            │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Pattern 3: Noun-Verb Structure (kubectl Style)                              │
│  ══════════════════════════════════════════════                              │
│                                                                              │
│  myapp <noun> <verb> [flags]                                                 │
│                                                                              │
│  myapp bucket create                                                         │
│  myapp bucket delete                                                         │
│  myapp bucket list                                                           │
│  myapp bucket get                                                            │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Pattern 4: Contextual Structure (git Style)                                 │
│  ════════════════════════════════════════════                                │
│                                                                              │
│  myapp <command> [subcommand] [flags]                                        │
│                                                                              │
│  myapp remote add origin <url>                                               │
│  myapp remote remove origin                                                  │
│  myapp remote list                                                           │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Pattern 5: Plugin Architecture (kubectl, helm)                              │
│  ══════════════════════════════════════════════                              │
│                                                                              │
│  myapp <plugin-command> [flags]                                              │
│                                                                              │
│  myapp-mycluster create                                                      │
│  myapp-mycluster list                                                        │
│                                                                              │
│  (Binary named myapp-mycluster in PATH)                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Scalability and Performance Considerations

### 3.1 CLI Performance Optimization

| Technique | Before | After | Benefit |
|-----------|--------|-------|---------|
| Concurrent operations | Sequential | Parallel | 5-10x faster |
| Lazy loading | Load all | On-demand | Lower memory |
| Streaming | Buffer all | Process chunks | Handle large files |
| Connection reuse | New each time | HTTP keep-alive | Faster API calls |
| Progress indicators | No feedback | Progress bars | Better UX |

### 3.2 Concurrent Operations

```go
package cmd

import (
    "context"
    "sync"
    "golang.org/x/sync/errgroup"
)

// ParallelExecutor executes operations in parallel
type ParallelExecutor struct {
    maxConcurrency int
}

func NewParallelExecutor(maxConcurrency int) *ParallelExecutor {
    return &ParallelExecutor{maxConcurrency: maxConcurrency}
}

// Execute runs operations concurrently with error handling
func (pe *ParallelExecutor) Execute(ctx context.Context, operations []Operation) error {
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(pe.maxConcurrency)

    for _, op := range operations {
        op := op // capture range variable
        g.Go(func() error {
            return op.Execute(ctx)
        })
    }

    return g.Wait()
}

// ProgressReporter reports progress for long-running operations
type ProgressReporter struct {
    total     int
    completed int32
    mu        sync.RWMutex
    bar       *progressbar.ProgressBar
}

func NewProgressReporter(total int) *ProgressReporter {
    return &ProgressReporter{
        total: total,
        bar: progressbar.NewOptions(total,
            progressbar.OptionEnableColorCodes(true),
            progressbar.OptionShowBytes(false),
            progressbar.OptionSetWidth(50),
            progressbar.OptionSetDescription("Processing..."),
            progressbar.OptionSetTheme(progressbar.Theme{
                Saucer:        "[green]█[reset]",
                SaucerHead:    "[green]>[reset]",
                SaucerPadding: " ",
                BarStart:      "[",
                BarEnd:        "]",
            }),
        ),
    }
}

func (pr *ProgressReporter) Increment() {
    pr.bar.Add(1)
}

func (pr *ProgressReporter) Finish() {
    pr.bar.Finish()
}
```

---

## 4. Technology Stack Recommendations

### 4.1 CLI Framework Comparison

| Framework | Maturity | Features | Learning Curve | Best For |
|-----------|----------|----------|----------------|----------|
| Cobra | ★★★★★ | ★★★★★ | Medium | Large projects, kubectl-style |
| urfave/cli | ★★★★★ | ★★★★☆ | Low | Simple to medium projects |
| Kingpin | ★★★★☆ | ★★★★☆ | Low | Flag-heavy applications |
| Flaggy | ★★★☆☆ | ★★★☆☆ | Very Low | Minimal dependencies |
| CLI (std lib) | ★★★★★ | ★★☆☆☆ | Low | Learning, minimal tools |

### 4.2 Recommended Stack

| Component | Library | Purpose |
|-----------|---------|---------|
| Framework | Cobra | Command structure |
| Config | Viper | Configuration management |
| Colors | fatih/color | Terminal colors |
| Tables | olekukonko/tablewriter | Tabular output |
| Progress | schollz/progressbar | Progress indicators |
| Prompts | AlecAivazis/survey | Interactive prompts |
| Spinner | briandowns/spinner | Loading indicators |
| TUI | charmbracelet/bubbletea | Rich terminal UI |

---

## 5. Case Studies

### 5.1 kubectl Design

**Scale:** Industry standard for Kubernetes management

**Key Decisions:**

- Noun-verb structure: `kubectl get pods`
- Context-aware: kubeconfig for cluster selection
- Output formats: `-o json|yaml|wide|custom-columns`
- Plugin system: `kubectl-<plugin>` binaries

**Lessons:**

- Consistent flag naming (`-n` for namespace everywhere)
- Helpful error messages with suggestions
- Autocompletion generation

### 5.2 Docker CLI Evolution

**Evolution Path:**

- v1.x: Monolithic binary
- v1.12+: Client-server architecture
- v19.03+: Context support
- Present: Docker Compose integration

**Key Features:**

- Plugin support (docker-buildx)
- Context switching
- JSON output for automation

---

## 6. Go Implementation Examples

### 6.1 Complete Cobra Application

```go
package cmd

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "github.com/fatih/color"
    "github.com/olekukonko/tablewriter"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    cfgFile    string
    outputMode string
    verbose    bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "A CLI application for managing resources",
    Long: `MyApp is a comprehensive CLI tool for resource management.

It supports multiple output formats, parallel operations, and
provides a great user experience with progress indicators.`,
    Version: "1.0.0",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // Initialize logging based on verbose flag
        if verbose {
            // Set debug logging
        }
    },
}

// Execute runs the root command
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        color.Red("Error: %v", err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    // Global flags
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.myapp.yaml)")
    rootCmd.PersistentFlags().StringVarP(&outputMode, "output", "o", "table", "Output format: table|json|yaml")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

    viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
    viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)

        viper.AddConfigPath(home)
        viper.AddConfigPath(".")
        viper.SetConfigName(".myapp")
        viper.SetConfigType("yaml")
    }

    viper.AutomaticEnv()
    viper.SetEnvPrefix("MYAPP")

    if err := viper.ReadInConfig(); err == nil {
        if verbose {
            fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
        }
    }
}

// listCmd lists resources
var listCmd = &cobra.Command{
    Use:     "list [resource-type]",
    Aliases: []string{"ls", "l"},
    Short:   "List resources",
    Long:    `List all resources of a specific type with optional filtering.`,
    Example: `  # List all users
  myapp list users

  # List users with filter
  myapp list users --status active

  # List in JSON format
  myapp list users -o json`,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        resourceType := args[0]

        // Get filter flags
        status, _ := cmd.Flags().GetString("status")
        limit, _ := cmd.Flags().GetInt("limit")

        // Fetch resources
        ctx := context.Background()
        resources, err := fetchResources(ctx, resourceType, status, limit)
        if err != nil {
            return fmt.Errorf("failed to list %s: %w", resourceType, err)
        }

        // Output based on format
        switch outputMode {
        case "json":
            return outputJSON(resources)
        case "yaml":
            return outputYAML(resources)
        case "table":
            return outputTable(resources)
        default:
            return fmt.Errorf("unknown output format: %s", outputMode)
        }
    },
}

func init() {
    rootCmd.AddCommand(listCmd)

    listCmd.Flags().StringP("status", "s", "", "Filter by status")
    listCmd.Flags().IntP("limit", "l", 50, "Maximum number of results")
    listCmd.Flags().StringSliceP("columns", "c", []string{}, "Columns to display")
}

// createCmd creates a new resource
var createCmd = &cobra.Command{
    Use:   "create [resource-type] [name]",
    Short: "Create a new resource",
    Long:  `Create a new resource with the specified name and configuration.`,
    Example: `  # Create a user interactively
  myapp create user john

  # Create with flags
  myapp create user john --email john@example.com --role admin`,
    Args: cobra.MinimumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        resourceType := args[0]
        name := ""
        if len(args) > 1 {
            name = args[1]
        }

        // Interactive mode if no name provided
        if name == "" {
            var err error
            name, err = promptForName(resourceType)
            if err != nil {
                return err
            }
        }

        // Collect additional fields interactively or from flags
        fields, err := collectFields(cmd, resourceType)
        if err != nil {
            return err
        }

        // Create resource
        ctx := context.Background()
        resource, err := createResource(ctx, resourceType, name, fields)
        if err != nil {
            return fmt.Errorf("failed to create %s: %w", resourceType, err)
        }

        color.Green("✓ Successfully created %s: %s", resourceType, name)

        if verbose {
            return outputJSON(resource)
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(createCmd)

    createCmd.Flags().String("email", "", "Email address")
    createCmd.Flags().String("role", "user", "User role")
    createCmd.Flags().Bool("interactive", false, "Interactive mode")
}

// deleteCmd deletes resources
var deleteCmd = &cobra.Command{
    Use:     "delete [resource-type] [name]",
    Aliases: []string{"del", "rm"},
    Short:   "Delete resources",
    Long:    `Delete one or more resources. Use with caution - this cannot be undone.`,
    Example: `  # Delete a user
  myapp delete user john

  # Delete multiple users
  myapp delete user john jane bob

  # Delete with force (skip confirmation)
  myapp delete user john --force`,
    Args: cobra.MinimumNArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        resourceType := args[0]
        names := args[1:]

        force, _ := cmd.Flags().GetBool("force")

        // Confirm deletion
        if !force {
            confirmed, err := confirmDeletion(resourceType, names)
            if err != nil || !confirmed {
                return err
            }
        }

        // Delete resources
        ctx := context.Background()
        var failed []string

        for _, name := range names {
            if err := deleteResource(ctx, resourceType, name); err != nil {
                color.Red("✗ Failed to delete %s: %v", name, err)
                failed = append(failed, name)
            } else {
                color.Green("✓ Deleted %s: %s", resourceType, name)
            }
        }

        if len(failed) > 0 {
            return fmt.Errorf("failed to delete: %v", failed)
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(deleteCmd)
    deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

// Helper functions for output formatting
func outputTable(resources []Resource) error {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Name", "Status", "Created", "ID"})
    table.SetAutoWrapText(false)
    table.SetAutoFormatHeaders(true)
    table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
    table.SetAlignment(tablewriter.ALIGN_LEFT)
    table.SetCenterSeparator("")
    table.SetColumnSeparator("")
    table.SetRowSeparator("")
    table.SetHeaderLine(false)
    table.SetBorder(false)
    table.SetTablePadding("\t")
    table.SetNoWhiteSpace(true)

    for _, r := range resources {
        table.Append([]string{
            r.Name,
            r.Status,
            r.CreatedAt.Format("2006-01-02"),
            r.ID,
        })
    }

    table.Render()
    return nil
}

func outputJSON(resources []Resource) error {
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")
    return encoder.Encode(resources)
}

func outputYAML(resources []Resource) error {
    data, err := yaml.Marshal(resources)
    if err != nil {
        return err
    }
    fmt.Print(string(data))
    return nil
}
```

### 6.2 Interactive Prompts

```go
package cmd

import (
    "github.com/AlecAivazis/survey/v2"
)

// promptForName prompts user for resource name
func promptForName(resourceType string) (string, error) {
    var name string
    prompt := &survey.Input{
        Message: fmt.Sprintf("Enter %s name:", resourceType),
        Help:    "Name must be unique and contain only lowercase letters, numbers, and hyphens",
    }

    validator := survey.Required

    err := survey.AskOne(prompt, &name, survey.WithValidator(validator))
    return name, err
}

// confirmDeletion prompts for deletion confirmation
func confirmDeletion(resourceType string, names []string) (bool, error) {
    var confirmed bool

    message := fmt.Sprintf("Are you sure you want to delete %d %s(s): %v?",
        len(names), resourceType, names)

    prompt := &survey.Confirm{
        Message: message,
        Default: false,
    }

    err := survey.AskOne(prompt, &confirmed)
    return confirmed, err
}

// collectFields collects resource fields interactively
func collectFields(cmd *cobra.Command, resourceType string) (map[string]string, error) {
    fields := make(map[string]string)

    switch resourceType {
    case "user":
        // Email
        if email, _ := cmd.Flags().GetString("email"); email != "" {
            fields["email"] = email
        } else {
            var email string
            survey.AskOne(&survey.Input{
                Message: "Email:",
            }, &email, survey.WithValidator(survey.Required))
            fields["email"] = email
        }

        // Role selection
        if role, _ := cmd.Flags().GetString("role"); role != "" {
            fields["role"] = role
        } else {
            var role string
            survey.AskOne(&survey.Select{
                Message: "Select role:",
                Options: []string{"admin", "editor", "viewer"},
                Default: "viewer",
            }, &role)
            fields["role"] = role
        }

    case "project":
        var description string
        survey.AskOne(&survey.Multiline{
            Message: "Description:",
        }, &description)
        fields["description"] = description

        var public bool
        survey.AskOne(&survey.Confirm{
            Message: "Make public?",
            Default: false,
        }, &public)
        fields["public"] = fmt.Sprintf("%v", public)
    }

    return fields, nil
}
```

### 6.3 Configuration Management

```go
package config

import (
    "os"
    "path/filepath"

    "github.com/spf13/viper"
)

// Config holds CLI configuration
type Config struct {
    API       APIConfig       `mapstructure:"api"`
    Output    OutputConfig    `mapstructure:"output"`
    Defaults  DefaultsConfig  `mapstructure:"defaults"`
}

type APIConfig struct {
    Endpoint string `mapstructure:"endpoint"`
    Token    string `mapstructure:"token"`
    Timeout  int    `mapstructure:"timeout"`
}

type OutputConfig struct {
    Format string `mapstructure:"format"`
    Color  bool   `mapstructure:"color"`
}

type DefaultsConfig struct {
    Namespace string `mapstructure:"namespace"`
    Region    string `mapstructure:"region"`
}

// LoadConfig loads configuration from file and environment
func LoadConfig() (*Config, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    // Set defaults
    viper.SetDefault("api.timeout", 30)
    viper.SetDefault("output.format", "table")
    viper.SetDefault("output.color", true)

    // Config file paths
    viper.AddConfigPath(filepath.Join(home, ".config", "myapp"))
    viper.AddConfigPath(home)
    viper.AddConfigPath(".")
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    // Environment variables
    viper.SetEnvPrefix("MYAPP")
    viper.AutomaticEnv()

    // Read config file
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config) error {
    home, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    configDir := filepath.Join(home, ".config", "myapp")
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return err
    }

    viper.Set("api", config.API)
    viper.Set("output", config.Output)
    viper.Set("defaults", config.Defaults)

    configFile := filepath.Join(configDir, "config.yaml")
    return viper.WriteConfigAs(configFile)
}
```

---

## 7. Visual Representations

### 7.1 CLI Command Hierarchy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        CLI Command Hierarchy                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  myapp                                                                       │
│  ├── config                    # Configuration management                    │
│  │   ├── get                   # Get configuration value                     │
│  │   ├── set                   # Set configuration value                     │
│  │   ├── list                  # List all configurations                    │
│  │   └── delete                # Delete configuration                       │
│  │                                                                           │
│  ├── auth                      # Authentication                              │
│  │   ├── login                 # Authenticate user                           │
│  │   ├── logout                # Logout user                                 │
│  │   ├── status                # Check auth status                           │
│  │   └── token                 # Manage API tokens                           │
│  │       ├── create                                                          │
│  │       ├── list                                                            │
│  │       └── revoke                                                          │
│  │                                                                           │
│  ├── user                      # User management                             │
│  │   ├── list                  # List users                                  │
│  │   ├── create                # Create user                                 │
│  │   ├── get                   # Get user details                            │
│  │   ├── update                # Update user                                 │
│  │   └── delete                # Delete user                                 │
│  │                                                                           │
│  ├── project                   # Project management                          │
│  │   ├── list                                                              │
│  │   ├── create                                                            │
│  │   ├── get                                                               │
│  │   ├── update                                                            │
│  │   ├── delete                                                            │
│  │   └── members                # Project member management                 │
│  │       ├── add                                                            │
│  │       ├── list                                                           │
│  │       └── remove                                                         │
│  │                                                                           │
│  └── completion                # Shell completion                            │
│      ├── bash                                                              │
│      ├── zsh                                                               │
│      └── fish                                                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 CLI Execution Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        CLI Execution Flow                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Shell Input                                                                 │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  1. Command Parsing                                                  │    │
│  │     - Tokenize input                                                │    │
│  │     - Identify command path                                         │    │
│  │     - Extract flags and arguments                                   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  2. Configuration Loading                                            │    │
│  │     - Load config file (~/.myapp/config.yaml)                       │    │
│  │     - Override with environment variables                           │    │
│  │     - Override with CLI flags                                       │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  3. Pre-execution Hooks                                              │    │
│  │     - Validate authentication                                       │    │
│  │     - Check API connectivity                                        │    │
│  │     - Set up logging                                                │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  4. Command Execution                                                │    │
│  │     - Validate arguments                                            │    │
│  │     - Execute business logic                                        │    │
│  │     - Handle API calls                                              │    │
│  │     - Process results                                               │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  5. Output Formatting                                                │    │
│  │     - Format as table/json/yaml                                     │    │
│  │     - Apply colors if terminal                                      │    │
│  │     - Paginate if needed                                            │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│       │                                                                      │
│       ▼                                                                      │
│  ┌──────────────────────┬──────────────────────┐                          │
│  │     Success          │       Error          │                          │
│  │  ┌────────────────┐  │  ┌────────────────┐  │                          │
│  │  │ Output results │  │  │ Print error    │  │                          │
│  │  │ to stdout      │  │  │ to stderr      │  │                          │
│  │  │ Exit code: 0   │  │  │ Exit code: 1   │  │                          │
│  │  └────────────────┘  │  └────────────────┘  │                          │
│  └──────────────────────┴──────────────────────┘                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.3 Plugin Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Plugin Architecture                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Core CLI (myapp)                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                                                                     │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐  │    │
│  │  │  Core        │  │   Plugin     │  │      Command Router      │  │    │
│  │  │  Commands    │  │   Manager    │  │                          │  │    │
│  │  │              │  │              │  │  - Built-in commands     │  │    │
│  │  │  - list      │  │  - Discover  │  │  - Plugin commands       │  │    │
│  │  │  - create    │  │  - Load      │  │  - Help generation       │  │    │
│  │  │  - delete    │  │  - Execute   │  │                          │  │    │
│  │  └──────────────┘  └──────────────┘  └──────────────────────────┘  │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│              ┌─────────────────────┼─────────────────────┐                   │
│              │                     │                     │                   │
│              ▼                     ▼                     ▼                   │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐          │
│  │  myapp-aws      │    │  myapp-gcp      │    │  myapp-custom   │          │
│  │                 │    │                 │    │                 │          │
│  │  AWS-specific   │    │  GCP-specific   │    │  Custom         │          │
│  │  commands       │    │  commands       │    │  extensions     │          │
│  │                 │    │                 │    │                 │          │
│  │  $ myapp aws    │    │  $ myapp gcp    │    │  $ myapp custom │          │
│  │     ec2 list    │    │     compute list│    │     workflow    │          │
│  │     s3 sync     │    │     storage ls  │    │     deploy      │          │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘          │
│                                                                              │
│  Plugin Discovery:                                                           │
│  - Scan $PATH for binaries named myapp-*                                     │
│  - Parse plugin metadata (version, description)                              │
│  - Load plugin manifest for command structure                                │
│                                                                              │
│  Plugin Execution:                                                           │
│  - Invoke plugin binary with context (config, auth)                          │
│  - Pass through stdin/stdout/stderr                                          │
│  - Handle exit codes uniformly                                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Security Requirements

### 8.1 CLI Security Checklist

| Category | Requirement | Implementation |
|----------|-------------|----------------|
| Credentials | Secure storage | OS keyring integration |
| Tokens | Encryption at rest | AES-256 encryption |
| API Calls | TLS only | Force TLS 1.3 |
| Input | Validation | Sanitize all inputs |
| Output | No secrets | Mask sensitive data |
| Logging | No credentials | Redact tokens/passwords |

### 8.2 Secure Credential Storage

```go
package auth

import (
    "fmt"

    "github.com/zalando/go-keyring"
)

const serviceName = "myapp-cli"

// SecureStorage handles credential storage
type SecureStorage struct{}

// StoreToken securely stores API token
func (s *SecureStorage) StoreToken(username, token string) error {
    return keyring.Set(serviceName, username, token)
}

// GetToken retrieves stored API token
func (s *SecureStorage) GetToken(username string) (string, error) {
    return keyring.Get(serviceName, username)
}

// DeleteToken removes stored token
func (s *SecureStorage) DeleteToken(username string) error {
    return keyring.Delete(serviceName, username)
}

// TokenMasker masks tokens in output
func TokenMasker(token string) string {
    if len(token) <= 8 {
        return "****"
    }
    return token[:4] + "****" + token[len(token)-4:]
}
```

---

**Quality Rating**: S (18+ KB)
**Last Updated**: 2026-04-02
