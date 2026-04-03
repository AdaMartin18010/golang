# Runbooks Documentation

> **分类**: 工程与云原生
> **标签**: #runbooks #documentation #procedures #operations #playbooks
> **参考**: Google SRE, AWS Well-Architected, Azure Operations

---

## 1. Formal Definition

### 1.1 What is a Runbook?

A runbook (or playbook) is a documented set of procedures and operations that guide engineers through routine maintenance tasks, troubleshooting steps, and incident response. Runbooks codify operational knowledge, reduce cognitive load during incidents, and ensure consistent execution of procedures across team members.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Runbook Hierarchy                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                    STANDARD OPERATING PROCEDURES                     │   │
│   │  (High-level operational standards and principles)                   │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                         RUNBOOKS                                     │   │
│   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │   │
│   │  │  Alert      │  │  Routine    │  │  Troubleshoot│  │  Emergency  │ │   │
│   │  │  Response   │  │  Maintenance│  │  Procedures  │  │  Procedures │ │   │
│   │  │             │  │             │  │              │  │             │ │   │
│   │  │ • CPU High  │  │ • Database  │  │ • Network   │  │ • Failover  │ │   │
│   │  │ • Disk Full │  │   Backup    │  │   Issues    │  │ • Rollback  │ │   │
│   │  │ • 5xx Errors│  │ • Log Rotate│  │ • Auth Prob │  │ • Shutdown  │ │   │
│   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                      PROCEDURES & SCRIPTS                            │   │
│   │  (Step-by-step commands, scripts, and verification steps)            │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│   Key Principles:                                                           │
│   • Executable: Can be followed step-by-step without expert knowledge      │
│   • Verified: Tested and validated regularly                               │
│   • Versioned: Tracked in source control                                   │
│   • Accessible: Available when needed (even during outages)                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Runbook Types Classification

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Runbook Types                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  BY FREQUENCY                    BY TRIGGER                  BY COMPLEXITY  │
│  ───────────────────             ───────────────────         ─────────────  │
│                                                                             │
│  ┌───────────────┐               ┌───────────────┐           ┌───────────┐ │
│  │  ROUTINE      │               │  ALERT-BASED  │           │  SIMPLE   │ │
│  │  Daily/Weekly │               │               │           │           │ │
│  │               │               │ • Threshold   │           │ < 5 steps │ │
│  │ • Log rotate  │               │   breaches    │           │ No decisions│
│  │ • Health      │               │ • Error rate  │           │           │ │
│  │   checks      │               │ • Latency     │           │ Verify &  │ │
│  │ • Backup      │               │ • Capacity    │           │ execute   │ │
│  │   verification│               │               │           │           │ │
│  └───────────────┘               └───────────────┘           └───────────┘ │
│                                                                             │
│  ┌───────────────┐               ┌───────────────┐           ┌───────────┐ │
│  │  PERIODIC     │               │  SCHEDULED    │           │  MODERATE │ │
│  │  Monthly/     │               │               │           │           │ │
│  │  Quarterly    │               │ • Maintenance │           │ 5-15 steps│ │
│  │               │               │   windows     │           │ Decision  │ │
│  │ • Cert        │               │ • Release     │           │   points  │ │
│  │   rotation    │               │   deploy      │           │           │ │
│  │ • Access      │               │ • Compliance  │           │ Requires  │ │
│  │   review      │               │   audits      │           │ judgment  │ │
│  └───────────────┘               └───────────────┘           └───────────┘ │
│                                                                             │
│  ┌───────────────┐               ┌───────────────┐           ┌───────────┐ │
│  │  AD-HOC       │               │  MANUAL       │           │  COMPLEX  │ │
│  │  As needed    │               │               │           │           │ │
│  │               │               │ • On-demand   │           │ > 15 steps│ │
│  │ • Incident    │               │ • Customer    │           │ Multiple  │ │
│  │   response    │               │   request     │           │   phases  │ │
│  │ • Troubleshoot│               │ • Security    │           │           │ │
│  │ • Escalation  │               │   incident    │           │ Requires  │ │
│  └───────────────┘               └───────────────┘           │ expertise │ │
│                                                              └───────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns in Go

### 2.1 Executable Runbook Engine

```go
package runbook

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "sync"
    "time"
)

// Runbook represents an executable runbook
type Runbook struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Version     string            `json:"version"`
    Category    string            `json:"category"`

    // Metadata
    Author      string    `json:"author"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
    ReviewedBy  string    `json:"reviewed_by,omitempty"`

    // Execution
    Steps       []Step            `json:"steps"`
    Variables   map[string]string `json:"variables,omitempty"`
    Preconditions []Precondition  `json:"preconditions,omitempty"`

    // Safety
    AutoApprove bool              `json:"auto_approve"`
    DangerLevel string            `json:"danger_level"` // low, medium, high, critical
    Rollback    *RollbackProcedure `json:"rollback,omitempty"`
}

// Step represents a single runbook step
type Step struct {
    ID          string            `json:"id"`
    Title       string            `json:"title"`
    Description string            `json:"description"`
    Type        StepType          `json:"type"`

    // Command execution
    Command     string            `json:"command,omitempty"`
    Args        []string          `json:"args,omitempty"`
    WorkingDir  string            `json:"working_dir,omitempty"`
    Environment map[string]string `json:"environment,omitempty"`

    // Verification
    Verification *Verification    `json:"verification,omitempty"`

    // Control flow
    OnSuccess   string            `json:"on_success,omitempty"` // next step ID
    OnFailure   string            `json:"on_failure,omitempty"` // next step ID or "rollback"
    Condition   string            `json:"condition,omitempty"`  // conditional execution

    // Safety
    RequiresConfirmation bool     `json:"requires_confirmation"`
    Timeout     time.Duration     `json:"timeout,omitempty"`
    DangerLevel string            `json:"danger_level,omitempty"`
}

// StepType defines the type of step
type StepType string

const (
    StepTypeCommand     StepType = "command"
    StepTypeScript      StepType = "script"
    StepTypePrompt      StepType = "prompt"
    StepTypeDecision    StepType = "decision"
    StepTypeVerification StepType = "verification"
    StepTypeNotification StepType = "notification"
)

// Verification defines verification criteria
type Verification struct {
    Type        string `json:"type"` // exit_code, output_contains, metric_threshold
    Expected    string `json:"expected,omitempty"`
    RetryCount  int    `json:"retry_count,omitempty"`
    RetryDelay  time.Duration `json:"retry_delay,omitempty"`
}

// Precondition defines prerequisites for runbook execution
type Precondition struct {
    Name        string `json:"name"`
    Check       string `json:"check"`
    AutoFix     string `json:"auto_fix,omitempty"`
    Required    bool   `json:"required"`
}

// RollbackProcedure defines how to undo changes
type RollbackProcedure struct {
    Steps []Step `json:"steps"`
    TriggerConditions []string `json:"trigger_conditions,omitempty"`
}

// ExecutionEngine executes runbooks
type ExecutionEngine struct {
    executor CommandExecutor
    logger   ExecutionLogger
    notifier Notifier

    activeExecutions map[string]*Execution
    mu               sync.RWMutex
}

// CommandExecutor executes commands
type CommandExecutor interface {
    Execute(ctx context.Context, command string, args []string, env map[string]string, timeout time.Duration) (*CommandResult, error)
}

// CommandResult contains command execution results
type CommandResult struct {
    ExitCode int
    Stdout   string
    Stderr   string
    Duration time.Duration
}

// ExecutionLogger logs execution details
type ExecutionLogger interface {
    LogStepStart(executionID, stepID string)
    LogStepComplete(executionID, stepID string, result *StepResult)
    LogExecutionComplete(executionID string, status ExecutionStatus)
}

// Notifier sends notifications
type Notifier interface {
    NotifyStepComplete(step Step, result *StepResult)
    NotifyExecutionComplete(runbook *Runbook, status ExecutionStatus)
}

// Execution represents a runbook execution
type Execution struct {
    ID        string
    RunbookID string
    StartedAt time.Time
    Status    ExecutionStatus
    Variables map[string]string
    Results   map[string]*StepResult
    CurrentStep int
    mu        sync.RWMutex
}

// ExecutionStatus represents execution status
type ExecutionStatus string

const (
    ExecutionStatusPending   ExecutionStatus = "pending"
    ExecutionStatusRunning   ExecutionStatus = "running"
    ExecutionStatusPaused    ExecutionStatus = "paused"
    ExecutionStatusCompleted ExecutionStatus = "completed"
    ExecutionStatusFailed    ExecutionStatus = "failed"
    ExecutionStatusRolledBack ExecutionStatus = "rolled_back"
)

// StepResult represents step execution result
type StepResult struct {
    StepID    string
    StartedAt time.Time
    EndedAt   time.Time
    Status    string // success, failure, skipped
    Output    string
    Error     string
}

// NewExecutionEngine creates a new execution engine
func NewExecutionEngine(executor CommandExecutor, logger ExecutionLogger, notifier Notifier) *ExecutionEngine {
    return &ExecutionEngine{
        executor:         executor,
        logger:           logger,
        notifier:         notifier,
        activeExecutions: make(map[string]*Execution),
    }
}

// ExecuteRunbook executes a runbook
func (e *ExecutionEngine) ExecuteRunbook(ctx context.Context, runbook *Runbook, variables map[string]string) (*Execution, error) {
    execution := &Execution{
        ID:        generateExecutionID(),
        RunbookID: runbook.ID,
        StartedAt: time.Now(),
        Status:    ExecutionStatusPending,
        Variables: mergeVariables(runbook.Variables, variables),
        Results:   make(map[string]*StepResult),
    }

    e.mu.Lock()
    e.activeExecutions[execution.ID] = execution
    e.mu.Unlock()

    // Check preconditions
    if err := e.checkPreconditions(ctx, runbook.Preconditions, execution); err != nil {
        execution.Status = ExecutionStatusFailed
        return execution, fmt.Errorf("precondition check failed: %w", err)
    }

    execution.Status = ExecutionStatusRunning

    // Execute steps
    for i, step := range runbook.Steps {
        execution.mu.Lock()
        execution.CurrentStep = i
        execution.mu.Unlock()

        result, err := e.executeStep(ctx, &step, execution)
        execution.Results[step.ID] = result

        if err != nil {
            e.logger.LogStepComplete(execution.ID, step.ID, result)

            // Check if rollback needed
            if step.OnFailure == "rollback" && runbook.Rollback != nil {
                e.rollback(ctx, runbook.Rollback, execution)
                execution.Status = ExecutionStatusRolledBack
            } else {
                execution.Status = ExecutionStatusFailed
            }

            e.logger.LogExecutionComplete(execution.ID, execution.Status)
            return execution, err
        }

        e.logger.LogStepComplete(execution.ID, step.ID, result)
        e.notifier.NotifyStepComplete(step, result)
    }

    execution.Status = ExecutionStatusCompleted
    e.logger.LogExecutionComplete(execution.ID, execution.Status)
    e.notifier.NotifyExecutionComplete(runbook, execution.Status)

    return execution, nil
}

// executeStep executes a single step
func (e *ExecutionEngine) executeStep(ctx context.Context, step *Step, execution *Execution) (*StepResult, error) {
    e.logger.LogStepStart(execution.ID, step.ID)

    result := &StepResult{
        StepID:    step.ID,
        StartedAt: time.Now(),
    }

    // Check condition
    if step.Condition != "" && !e.evaluateCondition(step.Condition, execution) {
        result.Status = "skipped"
        result.EndedAt = time.Now()
        return result, nil
    }

    // Handle confirmation
    if step.RequiresConfirmation {
        if !e.requestConfirmation(step) {
            result.Status = "cancelled"
            result.EndedAt = time.Now()
            return result, fmt.Errorf("step cancelled by user")
        }
    }

    var err error

    switch step.Type {
    case StepTypeCommand:
        err = e.executeCommandStep(ctx, step, execution, result)
    case StepTypePrompt:
        err = e.executePromptStep(step, execution, result)
    case StepTypeDecision:
        err = e.executeDecisionStep(step, execution, result)
    default:
        err = fmt.Errorf("unknown step type: %s", step.Type)
    }

    result.EndedAt = time.Now()

    if err != nil {
        result.Status = "failure"
        result.Error = err.Error()
        return result, err
    }

    result.Status = "success"
    return result, nil
}

// executeCommandStep executes a command step
func (e *ExecutionEngine) executeCommandStep(ctx context.Context, step *Step, execution *Execution, result *StepResult) error {
    // Substitute variables
    command := e.substituteVariables(step.Command, execution.Variables)
    args := make([]string, len(step.Args))
    for i, arg := range step.Args {
        args[i] = e.substituteVariables(arg, execution.Variables)
    }

    env := make(map[string]string)
    for k, v := range step.Environment {
        env[k] = e.substituteVariables(v, execution.Variables)
    }

    timeout := step.Timeout
    if timeout == 0 {
        timeout = 5 * time.Minute
    }

    cmdResult, err := e.executor.Execute(ctx, command, args, env, timeout)

    result.Output = cmdResult.Stdout
    if cmdResult.Stderr != "" {
        result.Output += "\nSTDERR: " + cmdResult.Stderr
    }

    if err != nil {
        return err
    }

    // Verify result
    if step.Verification != nil {
        if err := e.verifyResult(cmdResult, step.Verification); err != nil {
            return err
        }
    }

    return nil
}

// substituteVariables replaces variable placeholders
func (e *ExecutionEngine) substituteVariables(input string, variables map[string]string) string {
    result := input
    for key, value := range variables {
        placeholder := fmt.Sprintf("{{%s}}", key)
        result = strings.ReplaceAll(result, placeholder, value)
    }
    return result
}

// evaluateCondition evaluates a condition expression
func (e *ExecutionEngine) evaluateCondition(condition string, execution *Execution) bool {
    // Simplified condition evaluation
    // In production, use a proper expression evaluator
    return true
}

// requestConfirmation requests user confirmation
func (e *ExecutionEngine) requestConfirmation(step *Step) bool {
    fmt.Printf("\n⚠️  Step requires confirmation: %s\n", step.Title)
    fmt.Printf("Description: %s\n", step.Description)
    if step.DangerLevel != "" {
        fmt.Printf("Danger Level: %s\n", step.DangerLevel)
    }
    fmt.Print("\nProceed? (yes/no): ")

    var response string
    fmt.Scanln(&response)

    return strings.ToLower(response) == "yes"
}

// checkPreconditions checks runbook preconditions
func (e *ExecutionEngine) checkPreconditions(ctx context.Context, preconditions []Precondition, execution *Execution) error {
    for _, pre := range preconditions {
        // Execute check
        cmdResult, err := e.executor.Execute(ctx, "sh", []string{"-c", pre.Check}, nil, 30*time.Second)

        if err != nil || cmdResult.ExitCode != 0 {
            if pre.AutoFix != "" {
                // Try auto-fix
                _, fixErr := e.executor.Execute(ctx, "sh", []string{"-c", pre.AutoFix}, nil, 60*time.Second)
                if fixErr != nil && pre.Required {
                    return fmt.Errorf("precondition '%s' failed and auto-fix failed: %v", pre.Name, fixErr)
                }

                // Re-check
                cmdResult, err = e.executor.Execute(ctx, "sh", []string{"-c", pre.Check}, nil, 30*time.Second)
                if err != nil || cmdResult.ExitCode != 0 {
                    if pre.Required {
                        return fmt.Errorf("precondition '%s' failed after auto-fix", pre.Name)
                    }
                }
            } else if pre.Required {
                return fmt.Errorf("precondition '%s' failed: %s", pre.Name, cmdResult.Stderr)
            }
        }
    }

    return nil
}

// rollback executes rollback procedure
func (e *ExecutionEngine) rollback(ctx context.Context, rollback *RollbackProcedure, execution *Execution) {
    fmt.Println("\n⚠️  Executing rollback procedure...")

    for _, step := range rollback.Steps {
        result := &StepResult{
            StepID:    step.ID,
            StartedAt: time.Now(),
        }

        e.logger.LogStepStart(execution.ID, "rollback-"+step.ID)

        err := e.executeCommandStep(ctx, &step, execution, result)
        result.EndedAt = time.Now()

        if err != nil {
            result.Status = "failure"
            result.Error = err.Error()
            fmt.Printf("Rollback step failed: %v\n", err)
        } else {
            result.Status = "success"
        }

        e.logger.LogStepComplete(execution.ID, "rollback-"+step.ID, result)
    }
}

// verifyResult verifies command result
func (e *ExecutionEngine) verifyResult(result *CommandResult, verification *Verification) error {
    switch verification.Type {
    case "exit_code":
        if verification.Expected != "" {
            expectedExitCode := 0
            fmt.Sscanf(verification.Expected, "%d", &expectedExitCode)
            if result.ExitCode != expectedExitCode {
                return fmt.Errorf("exit code verification failed: expected %d, got %d", expectedExitCode, result.ExitCode)
            }
        }
    case "output_contains":
        if !strings.Contains(result.Stdout, verification.Expected) {
            return fmt.Errorf("output verification failed: expected to contain '%s'", verification.Expected)
        }
    }

    return nil
}

// Helper functions
func generateExecutionID() string {
    return fmt.Sprintf("exec-%d", time.Now().UnixNano())
}

func mergeVariables(base, override map[string]string) map[string]string {
    result := make(map[string]string)
    for k, v := range base {
        result[k] = v
    }
    for k, v := range override {
        result[k] = v
    }
    return result
}

func (e *ExecutionEngine) executePromptStep(step *Step, execution *Execution, result *StepResult) error {
    // Implementation for prompt steps
    return nil
}

func (e *ExecutionEngine) executeDecisionStep(step *Step, execution *Execution, result *StepResult) error {
    // Implementation for decision steps
    return nil
}
```

### 2.2 Runbook Validation Framework

```go
package runbook

import (
    "context"
    "fmt"
    "regexp"
    "strings"
)

// Validator validates runbooks
type Validator struct {
    rules []ValidationRule
}

// ValidationRule defines a validation rule
type ValidationRule struct {
    ID          string
    Name        string
    Description string
    Severity    string // error, warning
    Check       func(*Runbook) []ValidationIssue
}

// ValidationIssue represents a validation issue
type ValidationIssue struct {
    RuleID      string
    Severity    string
    Message     string
    Location    string
    Suggestion  string
}

// ValidationResult contains validation results
type ValidationResult struct {
    Valid   bool
    Issues  []ValidationIssue
    Errors  int
    Warnings int
}

// NewValidator creates a new validator with default rules
func NewValidator() *Validator {
    v := &Validator{
        rules: make([]ValidationRule, 0),
    }

    // Add default rules
    v.AddRule(ValidationRule{
        ID:          "R001",
        Name:        "Required Fields",
        Description: "Runbook must have required fields",
        Severity:    "error",
        Check:       checkRequiredFields,
    })

    v.AddRule(ValidationRule{
        ID:          "R002",
        Name:        "Step IDs",
        Description: "All steps must have unique IDs",
        Severity:    "error",
        Check:       checkUniqueStepIDs,
    })

    v.AddRule(ValidationRule{
        ID:          "R003",
        Name:        "Danger Level",
        Description: "High danger steps must have confirmation",
        Severity:    "error",
        Check:       checkDangerLevelConfirmation,
    })

    v.AddRule(ValidationRule{
        ID:          "R004",
        Name:        "Timeout Specified",
        Description: "Command steps should have timeout",
        Severity:    "warning",
        Check:       checkTimeoutSpecified,
    })

    v.AddRule(ValidationRule{
        ID:          "R005",
        Name:        "Verification Steps",
        Description: "Destructive operations should have verification",
        Severity:    "warning",
        Check:       checkVerificationSteps,
    })

    v.AddRule(ValidationRule{
        ID:          "R006",
        Name:        "Rollback Procedure",
        Description: "High danger runbooks should have rollback",
        Severity:    "warning",
        Check:       checkRollbackProcedure,
    })

    v.AddRule(ValidationRule{
        ID:          "R007",
        Name:        "Variable Substitution",
        Description: "Variables should be properly formatted",
        Severity:    "error",
        Check:       checkVariableFormat,
    })

    v.AddRule(ValidationRule{
        ID:          "R008",
        Name:        "Last Reviewed",
        Description: "Runbook should be reviewed within last 90 days",
        Severity:    "warning",
        Check:       checkReviewDate,
    })

    return v
}

// AddRule adds a validation rule
func (v *Validator) AddRule(rule ValidationRule) {
    v.rules = append(v.rules, rule)
}

// Validate validates a runbook
func (v *Validator) Validate(runbook *Runbook) *ValidationResult {
    result := &ValidationResult{
        Valid:  true,
        Issues: make([]ValidationIssue, 0),
    }

    for _, rule := range v.rules {
        issues := rule.Check(runbook)
        for _, issue := range issues {
            issue.RuleID = rule.ID
            if issue.Severity == "" {
                issue.Severity = rule.Severity
            }
            result.Issues = append(result.Issues, issue)

            if issue.Severity == "error" {
                result.Errors++
                result.Valid = false
            } else {
                result.Warnings++
            }
        }
    }

    return result
}

// Validation check functions

func checkRequiredFields(r *Runbook) []ValidationIssue {
    issues := make([]ValidationIssue, 0)

    if r.ID == "" {
        issues = append(issues, ValidationIssue{
            Severity:   "error",
            Message:    "Runbook ID is required",
            Location:   "id",
            Suggestion: "Add a unique identifier for the runbook",
        })
    }

    if r.Name == "" {
        issues = append(issues, ValidationIssue{
            Severity:   "error",
            Message:    "Runbook name is required",
            Location:   "name",
            Suggestion: "Add a descriptive name",
        })
    }

    if r.Description == "" {
        issues = append(issues, ValidationIssue{
            Severity:   "error",
            Message:    "Description is required",
            Location:   "description",
            Suggestion: "Add a description of when to use this runbook",
        })
    }

    if len(r.Steps) == 0 {
        issues = append(issues, ValidationIssue{
            Severity:   "error",
            Message:    "Runbook must have at least one step",
            Location:   "steps",
            Suggestion: "Add execution steps",
        })
    }

    return issues
}

func checkUniqueStepIDs(r *Runbook) []ValidationIssue {
    issues := make([]ValidationIssue, 0)
    seen := make(map[string]int)

    for i, step := range r.Steps {
        if step.ID == "" {
            issues = append(issues, ValidationIssue{
                Severity:   "error",
                Message:    fmt.Sprintf("Step %d is missing ID", i+1),
                Location:   fmt.Sprintf("steps[%d].id", i),
                Suggestion: "Add a unique step ID",
            })
            continue
        }

        if firstSeen, exists := seen[step.ID]; exists {
            issues = append(issues, ValidationIssue{
                Severity:   "error",
                Message:    fmt.Sprintf("Duplicate step ID '%s'", step.ID),
                Location:   fmt.Sprintf("steps[%d].id", i),
                Suggestion: fmt.Sprintf("Step ID already used at index %d", firstSeen),
            })
        }
        seen[step.ID] = i
    }

    return issues
}

func checkDangerLevelConfirmation(r *Runbook) []ValidationIssue {
    issues := make([]ValidationIssue, 0)

    highDangerLevels := map[string]bool{
        "high":     true,
        "critical": true,
    }

    for i, step := range r.Steps {
        if highDangerLevels[step.DangerLevel] && !step.RequiresConfirmation {
            issues = append(issues, ValidationIssue{
                Severity:   "error",
                Message:    fmt.Sprintf("Step '%s' has danger level '%s' but no confirmation", step.Title, step.DangerLevel),
                Location:   fmt.Sprintf("steps[%d].requires_confirmation", i),
                Suggestion: "Add requires_confirmation: true for safety",
            })
        }
    }

    return issues
}

func checkTimeoutSpecified(r *Runbook) []ValidationIssue {
    issues := make([]ValidationIssue, 0)

    for i, step := range r.Steps {
        if step.Type == StepTypeCommand && step.Timeout == 0 {
            issues = append(issues, ValidationIssue{
                Severity:   "warning",
                Message:    fmt.Sprintf("Step '%s' has no timeout specified", step.Title),
                Location:   fmt.Sprintf("steps[%d].timeout", i),
                Suggestion: "Add a timeout to prevent hung processes",
            })
        }
    }

    return issues
}

func checkVerificationSteps(r *Runbook) []ValidationIssue {
    issues := make([]ValidationIssue, 0)

    dangerousCommands := []string{"rm", "drop", "delete", "shutdown", "restart"}

    for i, step := range r.Steps {
        if step.Type != StepTypeCommand {
            continue
        }

        cmdLower := strings.ToLower(step.Command + " " + strings.Join(step.Args, " "))

        for _, dangerous := range dangerousCommands {
            if strings.Contains(cmdLower, dangerous) && step.Verification == nil {
                issues = append(issues, ValidationIssue{
                    Severity:   "warning",
                    Message:    fmt.Sprintf("Step '%s' contains potentially dangerous command but no verification", step.Title),
                    Location:   fmt.Sprintf("steps[%d].verification", i),
                    Suggestion: "Add verification to ensure command succeeded",
                })
                break
            }
        }
    }

    return issues
}

func checkRollbackProcedure(r *Runbook) []ValidationIssue {
    issues := make([]ValidationIssue, 0)

    if r.DangerLevel == "critical" && r.Rollback == nil {
        issues = append(issues, ValidationIssue{
            Severity:   "warning",
            Message:    "Critical danger runbook should have rollback procedure",
            Location:   "rollback",
            Suggestion: "Add rollback steps for safety",
        })
    }

    return issues
}

func checkVariableFormat(r *Runbook) []ValidationIssue {
    issues := make([]ValidationIssue, 0)
    variablePattern := regexp.MustCompile(`\{\{\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*\}\}`)

    for i, step := range r.Steps {
        // Check command
        matches := variablePattern.FindAllString(step.Command, -1)
        for _, match := range matches {
            varName := strings.Trim(match, "{}")
            if _, exists := r.Variables[varName]; !exists {
                issues = append(issues, ValidationIssue{
                    Severity:   "error",
                    Message:    fmt.Sprintf("Undefined variable '%s' in step command", varName),
                    Location:   fmt.Sprintf("steps[%d].command", i),
                    Suggestion: fmt.Sprintf("Define variable '%s' or check spelling", varName),
                })
            }
        }
    }

    return issues
}

func checkReviewDate(r *Runbook) []ValidationIssue {
    issues := make([]ValidationIssue, 0)

    // This would check actual dates
    // For now, just check if reviewed_at is set
    if r.ReviewedAt == nil {
        issues = append(issues, ValidationIssue{
            Severity:   "warning",
            Message:    "Runbook has not been reviewed",
            Location:   "reviewed_at",
            Suggestion: "Schedule a review to ensure accuracy",
        })
    }

    return issues
}
```

---

## 3. Production-Ready Configurations

### 3.1 Kubernetes Runbook Operator

```yaml
# runbook-crd.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: runbooks.operations.example.com
spec:
  group: operations.example.com
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              required:
                - name
                - steps
              properties:
                name:
                  type: string
                description:
                  type: string
                category:
                  type: string
                  enum:
                    - alert-response
                    - routine-maintenance
                    - troubleshooting
                    - emergency
                severity:
                  type: string
                  enum:
                    - low
                    - medium
                    - high
                    - critical
                triggers:
                  type: array
                  items:
                    type: object
                    properties:
                      type:
                        type: string
                        enum:
                          - alert
                          - schedule
                          - manual
                      condition:
                        type: string
                steps:
                  type: array
                  items:
                    type: object
                    required:
                      - name
                      - action
                    properties:
                      name:
                        type: string
                      description:
                        type: string
                      action:
                        type: string
                      args:
                        type: array
                        items:
                          type: string
                      timeout:
                        type: string
                      verification:
                        type: object
                        properties:
                          type:
                            type: string
                          expected:
                            type: string
                      requiresConfirmation:
                        type: boolean
                      rollbackAction:
                        type: string
                rollback:
                  type: object
                  properties:
                    steps:
                      type: array
                      items:
                        type: object
            status:
              type: object
              properties:
                lastExecuted:
                  type: string
                  format: date-time
                executionCount:
                  type: integer
                lastResult:
                  type: string
      additionalPrinterColumns:
        - name: Category
          type: string
          jsonPath: .spec.category
        - name: Severity
          type: string
          jsonPath: .spec.severity
        - name: Last Run
          type: string
          jsonPath: .status.lastExecuted

---
# Example Runbook
apiVersion: operations.example.com/v1
kind: Runbook
metadata:
  name: high-cpu-response
  namespace: operations
spec:
  name: "High CPU Usage Response"
  description: "Respond to high CPU usage alerts"
  category: alert-response
  severity: high
  triggers:
    - type: alert
      condition: "cpu_usage > 80"
  steps:
    - name: identify-process
      description: "Identify the process consuming CPU"
      action: "kubectl top pods"
      args:
        - "-n"
        - "{{ namespace }}"
      timeout: "30s"

    - name: check-recent-deployments
      description: "Check for recent deployments"
      action: "kubectl rollout history"
      args:
        - "deployment/{{ deployment }}"
        - "-n"
        - "{{ namespace }}"

    - name: collect-logs
      description: "Collect logs from affected pods"
      action: "kubectl logs"
      args:
        - "-l"
        - "app={{ app }}"
        - "-n"
        - "{{ namespace }}"
        - "--tail=100"

    - name: scale-up
      description: "Scale up the deployment"
      action: "kubectl scale"
      args:
        - "deployment/{{ deployment }}"
        - "--replicas={{ target_replicas }}"
        - "-n"
        - "{{ namespace }}"
      requiresConfirmation: true
      verification:
        type: "deployment_ready"
        expected: "{{ target_replicas }}"

    - name: verify-resolution
      description: "Verify CPU usage has decreased"
      action: "kubectl top pods"
      args:
        - "-l"
        - "app={{ app }}"
        - "-n"
        - "{{ namespace }}"
      timeout: "2m"
      verification:
        type: "cpu_threshold"
        expected: "< 70"

  rollback:
    steps:
      - name: rollback-scale
        action: "kubectl scale"
        args:
          - "deployment/{{ deployment }}"
          - "--replicas={{ original_replicas }}"
          - "-n"
          - "{{ namespace }}"
```

---

## 4. Security Considerations

### 4.1 Runbook Security Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Runbook Security Matrix                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ASPECT              │  RISK                  │  MITIGATION                 │
├──────────────────────┼────────────────────────┼─────────────────────────────│
│  Secrets in runbooks │  Credential exposure   │  • Use secret references    │
│                      │                        │  • Never hardcode passwords │
│                      │                        │  • Vault integration        │
│  ────────────────────┼────────────────────────┼─────────────────────────────│
│  Privilege escalation│  Unauthorized access   │  • RBAC for runbook access  │
│                      │                        │  • Just-in-time access      │
│                      │                        │  • Audit all executions     │
│  ────────────────────┼────────────────────────┼─────────────────────────────│
│  Malicious commands  │  System damage         │  • Code review for runbooks │
│                      │                        │  • Approval for dangerous   │
│                      │                        │  • Dry-run mode             │
│  ────────────────────┼────────────────────────┼─────────────────────────────│
│  Data exposure       │  Sensitive data leak   │  • Sanitize outputs         │
│                      │                        │  • No PII in runbooks       │
│                      │                        │  • Audit logging            │
│  ────────────────────┼────────────────────────┼─────────────────────────────│
│  Execution in        │  Unintended prod       │  • Environment confirmation │
│  wrong environment   │  changes               │  • Pre-execution checks     │
│                      │                        │  • Multi-env validation     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Compliance Requirements

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Runbook Compliance Requirements                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  SOC 2 TYPE II                                                              │
│  ├─ CC6.1: Logical access - Runbook access controlled by RBAC              │
│  ├─ CC7.2: System monitoring - Runbooks include monitoring procedures      │
│  ├─ CC7.5: Incident recovery - Documented recovery procedures              │
│  └─ CC7.1: Detection - Runbooks for security event detection               │
│                                                                             │
│  ISO 27001                                                                  │
│  ├─ A.12.1: Operational procedures - Documented operational procedures     │
│  ├─ A.16.1: Incident management - Incident response runbooks               │
│  ├─ A.17.1: Continuity - Business continuity procedures                    │
│  └─ A.12.4: Logging - Log review and retention procedures                  │
│                                                                             │
│  HIPAA                                                                      │
│  ├─ §164.308(a)(7): Contingency plan - Data backup and recovery            │
│  ├─ §164.312(a)(2)(ii): Emergency access - Emergency access procedures     │
│  └─ §164.308(a)(1)(ii)(D): Information access - Access management          │
│                                                                             │
│  PCI DSS                                                                    │
│  ├─ Req 12.10.1: IR procedures - Incident response procedures              │
│  ├─ Req 12.10.4: 24/7 availability - Coverage and escalation               │
│  └─ Req 12.11.1: Quarterly IR testing - Test procedures documented         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Decision Matrices

### 6.1 Runbook Type Selection Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Runbook Type Selection Matrix                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Situation                        │  Recommended Type    │  Example         │
├───────────────────────────────────┼──────────────────────┼──────────────────│
│  Alert fires with known solution  │  Alert Response      │  Disk full       │
│  ─────────────────────────────────┼──────────────────────┼──────────────────│
│  Recurring scheduled task         │  Routine Maintenance │  Backup verify   │
│  ─────────────────────────────────┼──────────────────────┼──────────────────│
│  Unknown issue, investigation     │  Troubleshooting     │  Debug latency   │
│  needed                           │                      │                  │
│  ─────────────────────────────────┼──────────────────────┼──────────────────│
│  Service down, immediate action   │  Emergency Procedure │  Failover        │
│  needed                           │                      │                  │
│  ─────────────────────────────────┼──────────────────────┼──────────────────│
│  Complex multi-step process       │  Workflow            │  Deployment      │
│  ─────────────────────────────────┼──────────────────────┼──────────────────│
│  Compliance requirement           │  Compliance          │  Access review   │
│  ─────────────────────────────────┼──────────────────────┼──────────────────│
│  Security incident                │  Security IR         │  Breach response │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Automation Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Automation Decision Matrix                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Frequency     │  Risk      │  Complexity  │  Recommendation               │
├────────────────┼────────────┼──────────────┼───────────────────────────────│
│  > 10x/day     │  Low       │  Low         │  Full automation              │
│  ──────────────┼────────────┼──────────────┼───────────────────────────────│
│  > 10x/day     │  Low       │  Medium      │  Semi-automation with verify  │
│  ──────────────┼────────────┼──────────────┼───────────────────────────────│
│  > 10x/day     │  Any       │  High        │  Keep manual, optimize steps  │
│  ──────────────┼────────────┼──────────────┼───────────────────────────────│
│  1-10x/day     │  Low       │  Low/Medium  │  Full automation              │
│  ──────────────┼────────────┼──────────────┼───────────────────────────────│
│  1-10x/day     │  Medium    │  Any         │  Semi-automation              │
│  ──────────────┼────────────┼──────────────┼───────────────────────────────│
│  1-10x/day     │  High      │  Any         │  Manual with strong guardrails│
│  ──────────────┼────────────┼──────────────┼───────────────────────────────│
│  Weekly        │  Low       │  Any         │  Full automation              │
│  ──────────────┼────────────┼──────────────┼───────────────────────────────│
│  Weekly        │  Medium+   │  Any         │  Manual with checklist        │
│  ──────────────┼────────────┼──────────────┼───────────────────────────────│
│  Monthly       │  Any       │  Any         │  Manual with detailed runbook │
│                                                                             │
│  Risk Assessment:                                                          │
│  • Low: Read-only operations, reversible changes                           │
│  • Medium: Data modification, service restart                              │
│  • High: Data deletion, infrastructure changes, customer impact            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.3 Escalation Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Runbook Escalation Decision Matrix                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Condition                            │  Action                             │
├───────────────────────────────────────┼─────────────────────────────────────│
│  Step fails with known workaround     │  Attempt workaround, continue       │
│  ─────────────────────────────────────┼─────────────────────────────────────│
│  Step fails, no workaround            │  Execute rollback, escalate         │
│  ─────────────────────────────────────┼─────────────────────────────────────│
│  Verification fails                   │  Retry 2x, then escalate            │
│  ─────────────────────────────────────┼─────────────────────────────────────│
│  Timeout exceeded                     │  Check status, conditional continue │
│  ─────────────────────────────────────┼─────────────────────────────────────│
│  Customer impact detected             │  Immediately escalate to IC         │
│  ─────────────────────────────────────┼─────────────────────────────────────│
│  Security implication discovered      │  Pause, notify security team        │
│  ─────────────────────────────────────┼─────────────────────────────────────│
│  Data loss risk identified            │  Stop immediately, escalate         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Runbook Best Practices Summary                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  WRITING RUNBOOKS                                                           │
│  ✓ Start with clear prerequisites and preconditions                         │
│  ✓ Use simple, direct language                                              │
│  ✓ Include exact commands, not descriptions                                 │
│  ✓ Add expected output or results for each step                             │
│  ✓ Include rollback procedures for destructive operations                   │
│  ✓ Add verification steps after changes                                     │
│  ✓ Document known failure modes and workarounds                             │
│  ✓ Include estimated time for each step                                     │
│  ✓ Add links to relevant dashboards and logs                                │
│  ✓ Version control all runbooks                                             │
│                                                                             │
│  VALIDATION                                                                 │
│  ✓ Test runbooks in non-production first                                    │
│  ✓ Dry-run mode for destructive operations                                  │
│  ✓ Regular review schedule (quarterly minimum)                              │
│  ✓ Update after each incident where runbook was used                        │
│  ✓ Peer review for new runbooks                                             │
│  ✓ Automated validation where possible                                      │
│                                                                             │
│  ORGANIZATION                                                               │
│  ✓ Standard naming conventions                                              │
│  ✓ Clear categorization (alert, maintenance, emergency)                     │
│  ✓ Searchable and discoverable                                              │
│  ✓ Indexed by service/component                                             │
│  ✓ Cross-reference related runbooks                                         │
│  ✓ Maintenance ownership assigned                                           │
│                                                                             │
│  EXECUTION                                                                  │
│  ✓ Confirmation for destructive operations                                  │
│  ✓ Audit logging of all executions                                          │
│  ✓ Environment checks before execution                                      │
│  ✓ Progress tracking during execution                                       │
│  ✓ Easy access during outages (offline copies)                              │
│  ✓ Integration with incident management                                     │
│                                                                             │
│  CONTINUOUS IMPROVEMENT                                                     │
│  ✓ Track runbook effectiveness                                              │
│  ✓ Measure MTTR with and without runbooks                                   │
│  ✓ Collect feedback from users                                              │
│  ✓ Analyze failure patterns                                                 │
│  ✓ Automate steps that are frequently executed                              │
│  ✓ Share learnings across teams                                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. Google SRE Book - Operational Readiness
2. AWS Well-Architected Framework - Operational Excellence
3. Azure Operations Guide - Runbooks
4. PagerDuty - Incident Response Automation
5. ITIL 4 - Service Operation
