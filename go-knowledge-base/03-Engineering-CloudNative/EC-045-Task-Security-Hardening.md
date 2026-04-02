# 任务安全加固 (Task Security Hardening)

> **分类**: 工程与云原生
> **标签**: #security #hardening #isolation

---

## 代码注入防护

```go
// 安全的任务处理器注册
type SecureRegistry struct {
    allowedTypes map[string]TaskHandler
    sanitizer    *InputSanitizer
}

func (sr *SecureRegistry) Register(taskType string, handler TaskHandler) error {
    // 验证类型名
    if !isValidTaskType(taskType) {
        return fmt.Errorf("invalid task type: %s", taskType)
    }

    sr.allowedTypes[taskType] = handler
    return nil
}

func (sr *SecureRegistry) Execute(ctx context.Context, task *Task) error {
    // 验证任务类型
    handler, ok := sr.allowedTypes[task.Type]
    if !ok {
        return fmt.Errorf("unknown task type: %s", task.Type)
    }

    // 净化输入
    sanitized, err := sr.sanitizer.Sanitize(task)
    if err != nil {
        return fmt.Errorf("input sanitization failed: %w", err)
    }

    return handler.Handle(ctx, sanitized)
}

// 输入净化
type InputSanitizer struct {
    maxPayloadSize int64
    forbiddenPatterns []string
}

func (is *InputSanitizer) Sanitize(task *Task) (*Task, error) {
    // 检查 payload 大小
    if int64(len(task.Payload)) > is.maxPayloadSize {
        return nil, fmt.Errorf("payload too large: %d > %d",
            len(task.Payload), is.maxPayloadSize)
    }

    // 检查危险模式
    payloadStr := string(task.Payload)
    for _, pattern := range is.forbiddenPatterns {
        if strings.Contains(payloadStr, pattern) {
            return nil, fmt.Errorf("forbidden pattern found: %s", pattern)
        }
    }

    // 验证 JSON 结构
    var v interface{}
    if err := json.Unmarshal(task.Payload, &v); err != nil {
        return nil, fmt.Errorf("invalid json: %w", err)
    }

    return task, nil
}
```

---

## 沙箱执行

```go
// gVisor 沙箱执行器
type SandboxExecutor struct {
    runtime runsc.Runtime
    config  *SandboxConfig
}

type SandboxConfig struct {
    MaxCPU     float64
    MaxMemory  int64
    MaxDisk    int64
    Network    NetworkPolicy
    Seccomp    string
}

func (se *SandboxExecutor) Execute(ctx context.Context, task *Task) (*Result, error) {
    // 创建沙箱配置
    conf := &runsc.Config{
        Root:       "/var/run/runsc",
        Platform:   "ptrace",
        Network:    runsc.NetworkNone, // 默认无网络
        Seccomp:    se.config.Seccomp,
        Debug:      false,
    }

    // 配置资源限制
    spec := &specs.Spec{
        Linux: &specs.Linux{
            Resources: &specs.LinuxResources{
                CPU: &specs.LinuxCPU{
                    Quota:  int64(se.config.MaxCPU * 100000),
                    Period: 100000,
                },
                Memory: &specs.LinuxMemory{
                    Limit: &se.config.MaxMemory,
                },
            },
        },
    }

    // 启动沙箱
    container, err := se.runtime.Create(ctx, task.ID, spec, conf)
    if err != nil {
        return nil, fmt.Errorf("create sandbox: %w", err)
    }
    defer container.Destroy()

    // 在沙箱中执行
    return se.runInSandbox(ctx, container, task)
}
```

---

## 审计日志

```go
type SecurityAuditor struct {
    store AuditStore
}

func (sa *SecurityAuditor) LogEvent(ctx context.Context, event SecurityEvent) error {
    auditLog := AuditLog{
        Timestamp:   time.Now(),
        EventType:   event.Type,
        TaskID:      event.TaskID,
        UserID:      event.UserID,
        IPAddress:   event.IPAddress,
        Action:      event.Action,
        Result:      event.Result,
        Details:     event.Details,
    }

    // 敏感操作立即告警
    if event.IsSensitive() {
        sa.alert(auditLog)
    }

    return sa.store.Save(ctx, auditLog)
}

func (sa *SecurityAuditor) LogTaskExecution(ctx context.Context, task *Task, result *Result) error {
    return sa.LogEvent(ctx, SecurityEvent{
        Type:     "task_execution",
        TaskID:   task.ID,
        UserID:   getUserID(ctx),
        Action:   "execute",
        Result:   result.Status,
        Details: map[string]interface{}{
            "task_type": task.Type,
            "duration":  result.Duration,
        },
    })
}

// 异常检测
func (sa *SecurityAuditor) DetectAnomalies(ctx context.Context, window time.Duration) ([]Anomaly, error) {
    events, _ := sa.store.GetEvents(ctx, time.Now().Add(-window), time.Now())

    var anomalies []Anomaly

    // 检测异常模式
    if sa.detectRateSpike(events) {
        anomalies = append(anomalies, Anomaly{
            Type:        "rate_spike",
            Severity:    "high",
            Description: "Unusual task execution rate detected",
        })
    }

    if sa.detectRepeatedFailures(events) {
        anomalies = append(anomalies, Anomaly{
            Type:        "repeated_failures",
            Severity:    "medium",
            Description: "Multiple task failures from same source",
        })
    }

    return anomalies, nil
}
```

---

## 访问控制

```go
type TaskAuthorization struct {
    policy PolicyEngine
}

func (ta *TaskAuthorization) Authorize(ctx context.Context, user User, action Action, resource Resource) error {
    // 检查权限
    allowed, err := ta.policy.Evaluate(ctx, user, action, resource)
    if err != nil {
        return fmt.Errorf("policy evaluation failed: %w", err)
    }

    if !allowed {
        return &AuthorizationError{
            User:   user.ID,
            Action: action.String(),
            Resource: resource.String(),
        }
    }

    return nil
}

// RBAC 实现
type RBACPolicy struct {
    roles map[string]Role
    userRoles map[string][]string  // userID -> roleIDs
}

func (rp *RBACPolicy) Evaluate(ctx context.Context, user User, action Action, resource Resource) (bool, error) {
    userRoles := rp.userRoles[user.ID]

    for _, roleID := range userRoles {
        role := rp.roles[roleID]

        for _, permission := range role.Permissions {
            if permission.Resource == resource.Type &&
               (permission.Action == action.Type || permission.Action == "*") {
                return true, nil
            }
        }
    }

    return false, nil
}
```
