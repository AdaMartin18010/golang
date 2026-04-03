# Container Best Practices

> **分类**: 工程与云原生
> **标签**: #containers #docker #security #optimization #production
> **参考**: Docker Security, CIS Benchmarks, NIST SP 800-190

---

## 1. Formal Definition

### 1.1 Containerization Fundamentals

Containerization is an operating system-level virtualization method where the kernel allows the existence of multiple isolated user-space instances. Containers package application code with its dependencies, enabling consistent execution across different environments.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Container Architecture                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   TRADITIONAL DEPLOYMENT          CONTAINERIZED DEPLOYMENT                  │
│   ━━━━━━━━━━━━━━━━━━━━━━━         ━━━━━━━━━━━━━━━━━━━━━━━━━                 │
│                                                                             │
│   ┌─────────────────┐             ┌─────────────────────────────────────┐   │
│   │  Application A  │             │  ┌─────────┐ ┌─────────┐ ┌────────┐ │   │
│   ├─────────────────┤             │  │  App A  │ │  App B  │ │ App C  │ │   │
│   │  Dependencies   │             │  │ + Libs  │ │ + Libs  │ │+ Libs  │ │   │
│   ├─────────────────┤             │  │ + Bin   │ │ + Bin   │ │+ Bin   │ │   │
│   │  Operating      │             │  └────┬────┘ └────┬────┘ └───┬────┘ │   │
│   │  System         │             │       └───────────┴──────────┘      │   │
│   ├─────────────────┤             │         Container Runtime           │   │
│   │  Hardware       │             │              (Docker/containerd)     │   │
│   └─────────────────┘             ├─────────────────────────────────────┤   │
│                                   │         Host Operating System        │   │
│   ┌─────────────────┐             ├─────────────────────────────────────┤   │
│   │  Application B  │             │              Hardware                │   │
│   ├─────────────────┤             └─────────────────────────────────────┘   │
│   │  Dependencies   │                                                       │
│   ├─────────────────┤             ADVANTAGES:                               │
│   │  Operating      │             • Isolation between applications          │
│   │  System         │             • Consistent environment                  │
│   ├─────────────────┤             • Resource efficiency                     │
│   │  Hardware       │             • Rapid deployment                        │
│   └─────────────────┘             • Scalability                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Container Security Layers

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Container Security Layers                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  LAYER 7: Application Security                                              │
│  ├─ Secure code practices                                                   │
│  ├─ Dependency scanning                                                       │
│  ├─ Secrets management                                                        │
│  └─ Application-level auth                                                    │
│                                                                             │
│  LAYER 6: Container Image Security                                            │
│  ├─ Minimal base images                                                       │
│  ├─ Image scanning                                                            │
│  ├─ Signed images                                                             │
│  └─ Immutable tags                                                            │
│                                                                             │
│  LAYER 5: Container Runtime Security                                          │
│  ├─ Capability dropping                                                       │
│  ├─ Seccomp profiles                                                          │
│  ├─ AppArmor/SELinux                                                          │
│  └─ Resource limits                                                           │
│                                                                             │
│  LAYER 4: Orchestration Security                                              │
│  ├─ Pod security policies / Pod Security Standards                            │
│  ├─ Network policies                                                          │
│  ├─ RBAC                                                                      │
│  └─ Service mesh                                                              │
│                                                                             │
│  LAYER 3: Host Security                                                       │
│  ├─ OS hardening                                                              │
│  ├─ Kernel security modules                                                   │
│  ├─ Audit logging                                                             │
│  └─ Minimal host OS                                                           │
│                                                                             │
│  LAYER 2: Network Security                                                    │
│  ├─ TLS everywhere                                                            │
│  ├─ Network segmentation                                                      │
│  ├─ Ingress/Egress controls                                                   │
│  └─ Service discovery security                                                │
│                                                                             │
│  LAYER 1: Infrastructure Security                                             │
│  ├─ Cloud IAM                                                                 │
│  ├─ Node security                                                             │
│  ├─ Storage encryption                                                        │
│  └─ Backup and recovery                                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Implementation Patterns in Go

### 2.1 Secure Dockerfile Patterns

```dockerfile
# ============================================================================
# Pattern 1: Multi-Stage Build with Minimal Base Image
# ============================================================================

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for build
RUN adduser -D -g '' appuser

WORKDIR /build

# Copy and download dependencies (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o /build/server \
    ./cmd/server

# ============================================================================
# Production stage
# ============================================================================
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy binary
COPY --from=builder /build/server /server

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/server", "-health-check"] || exit 1

# Run
ENTRYPOINT ["/server"]

# ============================================================================
# Pattern 2: Distroless Base Image
# ============================================================================
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /build/server /server

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/server"]

# ============================================================================
# Pattern 3: Alpine with Security Hardening
# ============================================================================
FROM alpine:3.18

# Update and install security patches
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary
COPY --from=builder /build/server /app/server

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Read-only root filesystem
EXPOSE 8080

ENTRYPOINT ["/app/server"]
```

### 2.2 Go Application Containerization

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config holds application configuration
type Config struct {
    Port        string
    MetricsPort string
    LogLevel    string
}

func main() {
    cfg := &Config{
        Port:        getEnv("PORT", "8080"),
        MetricsPort: getEnv("METRICS_PORT", "9090"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
    }

    // Check if health check mode
    if len(os.Args) > 1 && os.Args[1] == "-health-check" {
        os.Exit(performHealthCheck(cfg.Port))
    }

    // Create application
    app := NewApplication(cfg)

    // Start servers
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Start main server
    go func() {
        if err := app.Start(ctx); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Start metrics server
    go func() {
        metricsMux := http.NewServeMux()
        metricsMux.Handle("/metrics", promhttp.Handler())

        metricsServer := &http.Server{
            Addr:    ":" + cfg.MetricsPort,
            Handler: metricsMux,
        }

        log.Printf("Metrics server starting on port %s", cfg.MetricsPort)
        if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("Metrics server error: %v", err)
        }
    }()

    // Wait for shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    sig := <-sigChan
    log.Printf("Received signal: %v", sig)

    // Graceful shutdown
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()

    if err := app.Shutdown(shutdownCtx); err != nil {
        log.Printf("Shutdown error: %v", err)
    }

    log.Println("Server stopped gracefully")
}

// Application represents the containerized application
type Application struct {
    cfg    *Config
    server *http.Server
}

// NewApplication creates a new application
func NewApplication(cfg *Config) *Application {
    mux := http.NewServeMux()

    // Health endpoints
    mux.HandleFunc("/health/live", handleLiveness)
    mux.HandleFunc("/health/ready", handleReadiness)

    // Application endpoints
    mux.HandleFunc("/", handleRequest)

    server := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    return &Application{
        cfg:    cfg,
        server: server,
    }
}

// Start starts the application
func (a *Application) Start(ctx context.Context) error {
    log.Printf("Server starting on port %s", a.cfg.Port)
    return a.server.ListenAndServe()
}

// Shutdown gracefully shuts down the application
func (a *Application) Shutdown(ctx context.Context) error {
    return a.server.Shutdown(ctx)
}

// Health check handlers
func handleLiveness(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"alive"}`))
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
    // Check dependencies (DB, cache, etc.)
    if !checkDependencies() {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte(`{"status":"not ready"}`))
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ready"}`))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok","service":"api"}`))
}

func checkDependencies() bool {
    // Implement dependency checks
    return true
}

func performHealthCheck(port string) int {
    resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health/live", port))
    if err != nil || resp.StatusCode != http.StatusOK {
        return 1
    }
    return 0
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 2.3 Container Security Scanner

```go
package security

import (
    "context"
    "encoding/json"
    "fmt"
    "os/exec"
    "strings"
)

// ImageScanner scans container images for vulnerabilities
type ImageScanner struct {
    scannerPath string
    timeout     int
}

// ScanResult contains scan results
type ScanResult struct {
    Image       string          `json:"image"`
    ScanTime    string          `json:"scan_time"`
    Vulnerabilities []Vulnerability `json:"vulnerabilities"`
    Summary     ScanSummary     `json:"summary"`
}

// Vulnerability represents a single vulnerability
type Vulnerability struct {
    ID          string   `json:"id"`
    Package     string   `json:"package"`
    Severity    string   `json:"severity"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    FixedVersion string  `json:"fixed_version,omitempty"`
    CVSS        float64  `json:"cvss,omitempty"`
}

// ScanSummary provides vulnerability counts
type ScanSummary struct {
    Total        int `json:"total"`
    Critical     int `json:"critical"`
    High         int `json:"high"`
    Medium       int `json:"medium"`
    Low          int `json:"low"`
    Fixable      int `json:"fixable"`
}

// NewImageScanner creates a new image scanner
func NewImageScanner(scannerPath string) *ImageScanner {
    return &ImageScanner{
        scannerPath: scannerPath,
        timeout:     300,
    }
}

// ScanImage scans a container image
func (s *ImageScanner) ScanImage(ctx context.Context, image string) (*ScanResult, error) {
    // Use Trivy as the scanner
    cmd := exec.CommandContext(ctx, "trivy", "image",
        "--format", "json",
        "--severity", "CRITICAL,HIGH,MEDIUM,LOW",
        image,
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("scan failed: %w", err)
    }

    return s.parseTrivyOutput(image, output)
}

// ScanDockerfile scans a Dockerfile
func (s *ImageScanner) ScanDockerfile(ctx context.Context, dockerfilePath string) ([]string, error) {
    issues := make([]string, 0)

    // Check for common security issues in Dockerfile
    checks := []struct {
        name        string
        pattern     string
        description string
    }{
        {
            name:        "ROOT_USER",
            pattern:     "USER root",
            description: "Container should not run as root",
        },
        {
            name:        "LATEST_TAG",
            pattern:     ":latest",
            description: "Avoid using 'latest' tag",
        },
        {
            name:        "ADD_COMMAND",
            pattern:     "ADD ",
            description: "Prefer COPY over ADD",
        },
        {
            name:        "NO_HEALTHCHECK",
            pattern:     "HEALTHCHECK",
            description: "Missing HEALTHCHECK instruction",
            // This is inverted - we check for absence
        },
        {
            name:        "SUDO_USAGE",
            pattern:     "sudo",
            description: "Avoid using sudo in containers",
        },
        {
            name:        "SSH_KEYS",
            pattern:     ".ssh",
            description: "Potential SSH keys in image",
        },
        {
            name:        "SECRETS_ENV",
            pattern:     "ENV.*SECRET",
            description: "Secrets should not be in ENV",
        },
        {
            name:        "PASSWORD_IN_ENV",
            pattern:     "ENV.*PASSWORD",
            description: "Password should not be in ENV",
        },
    }

    content, err := exec.CommandContext(ctx, "cat", dockerfilePath).Output()
    if err != nil {
        return nil, err
    }

    contentStr := string(content)
    hasHealthcheck := strings.Contains(contentStr, "HEALTHCHECK")

    for _, check := range checks {
        if check.name == "NO_HEALTHCHECK" {
            if !hasHealthcheck {
                issues = append(issues, fmt.Sprintf("SECURITY: %s - %s", check.name, check.description))
            }
            continue
        }

        if strings.Contains(contentStr, check.pattern) {
            issues = append(issues, fmt.Sprintf("SECURITY: %s - %s", check.name, check.description))
        }
    }

    // Check for secrets with gitleaks or similar
    // This is a simplified check
    secretPatterns := []string{
        "AKIA[0-9A-Z]{16}",  // AWS Access Key
        "ghp_[a-zA-Z0-9]{36}", // GitHub Personal Access Token
        "[a-zA-Z0-9_-]*:[a-zA-Z0-9_-]+@github.com", // GitHub URL with credentials
    }

    for _, pattern := range secretPatterns {
        cmd := exec.CommandContext(ctx, "grep", "-E", pattern, dockerfilePath)
        if output, _ := cmd.Output(); len(output) > 0 {
            issues = append(issues, fmt.Sprintf("SECURITY: POTENTIAL_SECRET - Found potential secret matching pattern: %s", pattern))
        }
    }

    return issues, nil
}

// parseTrivyOutput parses Trivy JSON output
func (s *ImageScanner) parseTrivyOutput(image string, output []byte) (*ScanResult, error) {
    var trivyResult struct {
        Results []struct {
            Target          string `json:"Target"`
            Vulnerabilities []struct {
                VulnerabilityID string `json:"VulnerabilityID"`
                PkgName         string `json:"PkgName"`
                Title           string `json:"Title"`
                Description     string `json:"Description"`
                Severity        string `json:"Severity"`
                FixedVersion    string `json:"FixedVersion"`
                CVSS            struct {
                    Nvd struct {
                        V3Score float64 `json:"V3Score"`
                    } `json:"nvd"`
                } `json:"CVSS"`
            } `json:"Vulnerabilities"`
        } `json:"Results"`
    }

    if err := json.Unmarshal(output, &trivyResult); err != nil {
        return nil, err
    }

    result := &ScanResult{
        Image:           image,
        Vulnerabilities: make([]Vulnerability, 0),
    }

    severityCount := map[string]int{
        "CRITICAL": 0,
        "HIGH":     0,
        "MEDIUM":   0,
        "LOW":      0,
    }

    for _, r := range trivyResult.Results {
        for _, v := range r.Vulnerabilities {
            vuln := Vulnerability{
                ID:          v.VulnerabilityID,
                Package:     v.PkgName,
                Severity:    v.Severity,
                Title:       v.Title,
                Description: v.Description,
                FixedVersion: v.FixedVersion,
                CVSS:        v.CVSS.Nvd.V3Score,
            }

            result.Vulnerabilities = append(result.Vulnerabilities, vuln)
            severityCount[v.Severity]++

            if v.FixedVersion != "" {
                result.Summary.Fixable++
            }
        }
    }

    result.Summary.Total = len(result.Vulnerabilities)
    result.Summary.Critical = severityCount["CRITICAL"]
    result.Summary.High = severityCount["HIGH"]
    result.Summary.Medium = severityCount["MEDIUM"]
    result.Summary.Low = severityCount["LOW"]

    return result, nil
}

// HasCriticalVulnerabilities checks if scan has critical vulnerabilities
func (r *ScanResult) HasCriticalVulnerabilities() bool {
    return r.Summary.Critical > 0
}

// HasHighVulnerabilities checks if scan has high vulnerabilities
func (r *ScanResult) HasHighVulnerabilities() bool {
    return r.Summary.High > 0
}

// GenerateReport generates a human-readable report
func (r *ScanResult) GenerateReport() string {
    var sb strings.Builder

    sb.WriteString(fmt.Sprintf("Container Image Security Scan Report\n"))
    sb.WriteString(fmt.Sprintf("====================================\n\n"))
    sb.WriteString(fmt.Sprintf("Image: %s\n", r.Image))
    sb.WriteString(fmt.Sprintf("Scan Time: %s\n\n", r.ScanTime))

    sb.WriteString(fmt.Sprintf("Summary:\n"))
    sb.WriteString(fmt.Sprintf("  Total: %d\n", r.Summary.Total))
    sb.WriteString(fmt.Sprintf("  Critical: %d\n", r.Summary.Critical))
    sb.WriteString(fmt.Sprintf("  High: %d\n", r.Summary.High))
    sb.WriteString(fmt.Sprintf("  Medium: %d\n", r.Summary.Medium))
    sb.WriteString(fmt.Sprintf("  Low: %d\n", r.Summary.Low))
    sb.WriteString(fmt.Sprintf("  Fixable: %d\n\n", r.Summary.Fixable))

    if len(r.Vulnerabilities) > 0 {
        sb.WriteString(fmt.Sprintf("Vulnerabilities:\n"))
        for _, v := range r.Vulnerabilities {
            sb.WriteString(fmt.Sprintf("  [%s] %s - %s\n", v.Severity, v.ID, v.Package))
            if v.FixedVersion != "" {
                sb.WriteString(fmt.Sprintf("    Fixed in: %s\n", v.FixedVersion))
            }
        }
    }

    return sb.String()
}
```

### 2.4 Container Runtime Security

```go
package security

import (
    "context"
    "fmt"
    "os/exec"
    "strings"
)

// RuntimeSecurityConfig defines container runtime security settings
type RuntimeSecurityConfig struct {
    // User configuration
    User                string
    Group               string
    NoNewPrivileges     bool

    // Capabilities
    DropAllCapabilities bool
    AddCapabilities     []string

    // Security options
    ReadOnlyRootFS      bool
    SeccompProfile      string
    AppArmorProfile     string
    SELinuxOptions      []string

    // Resource limits
    MemoryLimit         string
    CPUQuota            int
    CPUPeriod           int

    // Network
    NetworkMode         string
    DisableNetworking   bool

    // Filesystem
    TmpfsMounts         map[string]string
    VolumeMounts        map[string]string
}

// SecurityProfileGenerator generates Docker/Kubernetes security profiles
type SecurityProfileGenerator struct{}

// NewSecurityProfileGenerator creates a new generator
func NewSecurityProfileGenerator() *SecurityProfileGenerator {
    return &SecurityProfileGenerator{}
}

// GenerateDockerFlags generates Docker run flags
func (g *SecurityProfileGenerator) GenerateDockerFlags(config *RuntimeSecurityConfig) []string {
    flags := make([]string, 0)

    // User
    if config.User != "" {
        flags = append(flags, "--user", config.User)
    }

    // No new privileges
    if config.NoNewPrivileges {
        flags = append(flags, "--security-opt", "no-new-privileges:true")
    }

    // Capabilities
    if config.DropAllCapabilities {
        flags = append(flags, "--cap-drop", "ALL")
    }
    for _, cap := range config.AddCapabilities {
        flags = append(flags, "--cap-add", cap)
    }

    // Read-only root filesystem
    if config.ReadOnlyRootFS {
        flags = append(flags, "--read-only")
    }

    // Seccomp
    if config.SeccompProfile != "" {
        flags = append(flags, "--security-opt", fmt.Sprintf("seccomp=%s", config.SeccompProfile))
    }

    // AppArmor
    if config.AppArmorProfile != "" {
        flags = append(flags, "--security-opt", fmt.Sprintf("apparmor=%s", config.AppArmorProfile))
    }

    // SELinux
    for _, opt := range config.SELinuxOptions {
        flags = append(flags, "--security-opt", fmt.Sprintf("label=%s", opt))
    }

    // Memory limits
    if config.MemoryLimit != "" {
        flags = append(flags, "--memory", config.MemoryLimit)
    }

    // CPU limits
    if config.CPUQuota > 0 {
        flags = append(flags, "--cpu-quota", fmt.Sprintf("%d", config.CPUQuota))
    }
    if config.CPUPeriod > 0 {
        flags = append(flags, "--cpu-period", fmt.Sprintf("%d", config.CPUPeriod))
    }

    // Network
    if config.NetworkMode != "" {
        flags = append(flags, "--network", config.NetworkMode)
    }
    if config.DisableNetworking {
        flags = append(flags, "--network", "none")
    }

    // Tmpfs mounts
    for target, options := range config.TmpfsMounts {
        mount := fmt.Sprintf("--tmpfs=%s", target)
        if options != "" {
            mount = fmt.Sprintf("--tmpfs=%s:%s", target, options)
        }
        flags = append(flags, mount)
    }

    // Volume mounts
    for source, target := range config.VolumeMounts {
        flags = append(flags, "-v", fmt.Sprintf("%s:%s:ro", source, target))
    }

    return flags
}

// GenerateKubernetesSecurityContext generates Kubernetes security context
func (g *SecurityProfileGenerator) GenerateKubernetesSecurityContext(config *RuntimeSecurityConfig) map[string]interface{} {
    securityContext := make(map[string]interface{})

    // User
    if config.User != "" {
        var uid int
        fmt.Sscanf(config.User, "%d", &uid)
        securityContext["runAsUser"] = uid
        securityContext["runAsNonRoot"] = true
    }

    // Privileges
    if config.NoNewPrivileges {
        securityContext["allowPrivilegeEscalation"] = false
    }

    // Capabilities
    if config.DropAllCapabilities {
        securityContext["capabilities"] = map[string][]string{
            "drop": {"ALL"},
        }
        if len(config.AddCapabilities) > 0 {
            securityContext["capabilities"].(map[string][]string)["add"] = config.AddCapabilities
        }
    }

    // Read-only root filesystem
    if config.ReadOnlyRootFS {
        securityContext["readOnlyRootFilesystem"] = true
    }

    // Seccomp
    if config.SeccompProfile != "" {
        securityContext["seccompProfile"] = map[string]string{
            "type": "Localhost",
            "localhostProfile": config.SeccompProfile,
        }
    }

    // SELinux
    if len(config.SELinuxOptions) > 0 {
        securityContext["seLinuxOptions"] = parseSELinuxOptions(config.SELinuxOptions)
    }

    return securityContext
}

// parseSELinuxOptions parses SELinux options
func parseSELinuxOptions(options []string) map[string]string {
    result := make(map[string]string)
    for _, opt := range options {
        parts := strings.SplitN(opt, ":", 2)
        if len(parts) == 2 {
            result[parts[0]] = parts[1]
        }
    }
    return result
}

// DefaultSecurityConfig returns a secure default configuration
func DefaultSecurityConfig() *RuntimeSecurityConfig {
    return &RuntimeSecurityConfig{
        User:                "1000:1000",
        NoNewPrivileges:     true,
        DropAllCapabilities: true,
        AddCapabilities:     []string{},
        ReadOnlyRootFS:      true,
        SeccompProfile:      "/var/lib/kubelet/seccomp/default.json",
        NetworkMode:         "bridge",
        TmpfsMounts: map[string]string{
            "/tmp":     "rw,noexec,nosuid,size=100m",
            "/var/tmp": "rw,noexec,nosuid,size=100m",
        },
    }
}

// ValidateContainer validates a running container's security
func ValidateContainer(ctx context.Context, containerID string) ([]string, error) {
    issues := make([]string, 0)

    // Check if running as root
    cmd := exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.Config.User}}", containerID)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    user := strings.TrimSpace(string(output))
    if user == "" || user == "root" || user == "0" {
        issues = append(issues, "Container is running as root")
    }

    // Check capabilities
    cmd = exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.HostConfig.CapDrop}}", containerID)
    output, _ = cmd.Output()
    if !strings.Contains(string(output), "ALL") {
        issues = append(issues, "Container has not dropped all capabilities")
    }

    // Check if privileged
    cmd = exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.HostConfig.Privileged}}", containerID)
    output, _ = cmd.Output()
    if strings.TrimSpace(string(output)) == "true" {
        issues = append(issues, "Container is running in privileged mode")
    }

    // Check read-only root fs
    cmd = exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.HostConfig.ReadonlyRootfs}}", containerID)
    output, _ = cmd.Output()
    if strings.TrimSpace(string(output)) != "true" {
        issues = append(issues, "Container root filesystem is not read-only")
    }

    return issues, nil
}
```

---

## 3. Production-Ready Configurations

### 3.1 Kubernetes Pod Security

```yaml
# pod-security.yaml
apiVersion: v1
kind: Pod
metadata:
  name: secure-app
  labels:
    app: secure-app
spec:
  # Use specific service account (not default)
  serviceAccountName: secure-app-sa
  automountServiceAccountToken: false

  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    runAsGroup: 1000
    fsGroup: 1000
    seccompProfile:
      type: RuntimeDefault

  containers:
  - name: app
    image: myapp:v1.0.0
    imagePullPolicy: Always

    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL

    resources:
      requests:
        memory: "128Mi"
        cpu: "100m"
      limits:
        memory: "256Mi"
        cpu: "500m"

    # Read-only root filesystem requires writable directories as volumes
    volumeMounts:
    - name: tmp
      mountPath: /tmp
    - name: cache
      mountPath: /app/cache

    livenessProbe:
      httpGet:
        path: /health/live
        port: 8080
      initialDelaySeconds: 10
      periodSeconds: 10

    readinessProbe:
      httpGet:
        path: /health/ready
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 5

    # Environment variables (no secrets here)
    env:
    - name: LOG_LEVEL
      value: "info"
    - name: PORT
      value: "8080"

    # Secrets via secret references
    envFrom:
    - secretRef:
        name: app-secrets
        optional: false

  volumes:
  - name: tmp
    emptyDir:
      sizeLimit: 100Mi
  - name: cache
    emptyDir:
      sizeLimit: 500Mi

  # Node affinity for dedicated nodes
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: node-type
            operator: In
            values:
            - general

  # Anti-affinity for HA
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - secure-app
        topologyKey: kubernetes.io/hostname

---
# Pod Security Standard (PSS) - Restricted
apiVersion: v1
kind: Namespace
metadata:
  name: restricted-ns
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

### 3.2 BuildKit Configuration

```yaml
# buildkitd.toml
# BuildKit daemon configuration for secure builds

[worker.oci]
  enabled = true
  snapshotter = "overlayfs"

  # Security options
  no-process-sandbox = false

  # GC policy
  [worker.oci.gcpolicy]
    [[worker.oci.gcpolicy.rules]]
      filter = ["type==source.local", "type==exec.cachemount"]
      keepDuration = 3600
      keepBytes = 10737418240  # 10GB

    [[worker.oci.gcpolicy.rules]]
      all = true
      keepBytes = 53687091200  # 50GB

[registry]
  # Configure registry mirrors for faster pulls
  [registry."docker.io"]
    mirrors = ["mirror.gcr.io"]

  # Insecure registries (for development only)
  # [registry."insecure-registry:5000"]
  #   http = true

  # TLS configuration
  [registry."secure-registry.example.com"]
    ca = ["/etc/buildkit/certs/ca.crt"]
    [[registry."secure-registry.example.com".keypair]]
      key = "/etc/buildkit/certs/client.key"
      cert = "/etc/buildkit/certs/client.crt"
```

---

## 4. Security Considerations

### 4.1 Container Security Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     Container Security Checklist                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  IMAGE BUILD                                                                │
│  [ ] Use minimal base images (scratch, distroless, alpine)                  │
│  [ ] Pin image tags to specific versions (no 'latest')                      │
│  [ ] Multi-stage builds to reduce attack surface                            │
│  [ ] Scan images for vulnerabilities in CI/CD                               │
│  [ ] Sign images with Cosign/Notary                                         │
│  [ ] No secrets in images (use build secrets)                               │
│  [ ] Remove build dependencies from final image                             │
│  [ ] Use non-root user in Dockerfile                                        │
│                                                                             │
│  RUNTIME SECURITY                                                           │
│  [ ] Run containers as non-root user                                        │
│  [ ] Drop all capabilities, add only required ones                          │
│  [ ] Enable read-only root filesystem                                       │
│  [ ] Disable privilege escalation                                           │
│  [ ] Use seccomp and AppArmor/SELinux profiles                              │
│  [ ] Set resource limits (CPU, memory, disk)                                │
│  [ ] Use specific user/group IDs                                            │
│  [ ] No privileged containers in production                                 │
│                                                                             │
│  NETWORK SECURITY                                                           │
│  [ ] Use specific network modes (avoid host network)                        │
│  [ ] Implement network policies                                             │
│  [ ] TLS for all inter-service communication                                │
│  [ ] Expose only necessary ports                                            │
│  [ ] Service mesh for zero-trust networking                                 │
│                                                                             │
│  SECRETS MANAGEMENT                                                         │
│  [ ] No secrets in environment variables                                    │
│  [ ] Use secret management solutions (Vault, Sealed Secrets)                │
│  [ ] Rotate secrets regularly                                               │
│  [ ] Encrypt secrets at rest                                                │
│  [ ] Audit secret access                                                    │
│                                                                             │
│  MONITORING & AUDITING                                                      │
│  [ ] Container runtime security monitoring (Falco)                          │
│  [ ] Audit container events                                                 │
│  [ ] Log all container activities                                           │
│  [ ] Alert on security policy violations                                    │
│  [ ] Regular security scans of running containers                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Compliance Requirements

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Container Compliance Requirements                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  NIST SP 800-190                     CIS Benchmarks                         │
│  ━━━━━━━━━━━━━━━━                    ━━━━━━━━━━━━━━                         │
│                                                                             │
│  • Image vulnerabilities (4.1)       • 4.1 - Create user for container      │
│  • Image configuration (4.2)         • 4.3 - No sensitive data in ENV       │
│  • Registry security (4.3)           • 4.6 - Health check instruction       │
│  • Orchestrator vulnerabilities (4.4)• 5.1 - SELinux/AppArmor options       │
│  • Container vulnerabilities (4.5)   • 5.2 - Linux kernel capabilities      │
│  • Host OS vulnerabilities (4.6)     • 5.4 - Privileged containers          │
│                                                                             │
│  PCI DSS                             SOC 2                                  │
│  ━━━━━━━━                            ━━━━━━━                                │
│                                                                             │
│  Req 1.3 - Network segmentation      CC6.1 - Logical access                 │
│  Req 2.1 - Vendor defaults           CC6.6 - Encryption                     │
│  Req 6.2 - Security patches          CC6.8 - Security infrastructure        │
│  Req 10.2 - Audit trails             CC7.1 - Detection                      │
│  Req 11.2.1 - Vulnerability scans    CC7.2 - Monitoring                     │
│                                                                             │
│  HIPAA                                                                      │
│  ━━━━━━                                                                     │
│                                                                             │
│  §164.312(a)(1) - Access control                                            │
│  §164.312(a)(2)(iv) - Encryption                                            │
│  §164.312(b) - Audit controls                                               │
│  §164.312(e)(1) - Transmission security                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. Decision Matrices

### 6.1 Base Image Selection Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Base Image Selection Matrix                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Criteria         │  SCRATCH   │  DISTROLESS │  ALPINE    │  UBUNTU/DEBIAN│
├───────────────────┼────────────┼─────────────┼────────────┼───────────────│
│  Size             │  ★★★★★     │  ★★★★★      │  ★★★★☆     │  ★★☆☆☆       │
│  (Small is good)  │  ~0MB      │  ~2MB       │  ~5MB      │  ~100MB+     │
├───────────────────┼────────────┼─────────────┼────────────┼───────────────│
│  Security         │  ★★★★★     │  ★★★★★      │  ★★★★☆     │  ★★★☆☆       │
│  (Less attack     │  Minimal   │  Minimal    │  Musl libc │  Full OS     │
│   surface)        │  surface   │  surface    │            │              │
├───────────────────┼────────────┼─────────────┼────────────┼───────────────│
│  Compatibility    │  ★★☆☆☆     │  ★★★☆☆      │  ★★★★☆     │  ★★★★★       │
│  (Apps work)      │  Static    │  Static     │  Good      │  Excellent   │
│                   │  only      │  only       │            │              │
├───────────────────┼────────────┼─────────────┼────────────┼───────────────│
│  Debuggability    │  ★☆☆☆☆     │  ★☆☆☆☆      │  ★★★☆☆     │  ★★★★★       │
│  (Shell access)   │  No shell  │  No shell   │  BusyBox   │  Full tools  │
├───────────────────┼────────────┼─────────────┼────────────┼───────────────│
│  Build Complexity │  ★★★☆☆     │  ★★★☆☆      │  ★★★★☆     │  ★★★★★       │
│  (Easy to use)    │  Complex   │  Complex    │  Easy      │  Very Easy   │
├───────────────────┼────────────┼─────────────┼────────────┼───────────────│
│  WHEN TO USE      │  Go, Rust  │  Go, Java   │  Multi-lang│  Dev, Debug  │
│                   │  static    │  static     │  static    │  legacy apps │
│                                                                             │
│  Recommendation:                                                            │
│  • Production services: Distroless or Scratch (static languages)            │
│  • Multi-language support: Alpine (with security updates)                   │
│  • Development/ debugging: Ubuntu/Debian                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Security Hardening Decision Matrix

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Security Hardening Decision Matrix                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Container Use Case            │  Security Level  │  Key Controls          │
├────────────────────────────────┼──────────────────┼────────────────────────│
│  Public-facing web service     │  Maximum         │  • Non-root user       │
│                                │                  │  • Read-only root fs   │
│                                │                  │  • Drop ALL caps       │
│                                │                  │  • Seccomp + AppArmor  │
│                                │                  │  • Network policies    │
│                                │                  │  • Pod Security Std    │
├────────────────────────────────┼──────────────────┼────────────────────────│
│  Internal microservice         │  High            │  • Non-root user       │
│                                │                  │  • Read-only root fs   │
│                                │                  │  • Drop ALL caps       │
│                                │                  │  • Seccomp profile     │
│                                │                  │  • Network policies    │
├────────────────────────────────┼──────────────────┼────────────────────────│
│  Batch job / CronJob           │  High            │  • Non-root user       │
│                                │                  │  • Drop ALL caps       │
│                                │                  │  • No network (if ok)  │
│                                │                  │  • Resource limits     │
├────────────────────────────────┼──────────────────┼────────────────────────│
│  Development / Testing         │  Medium          │  • Non-root user       │
│                                │                  │  • No privileged       │
│                                │                  │  • Resource limits     │
├────────────────────────────────┼──────────────────┼────────────────────────│
│  Monitoring / Observability    │  Medium          │  • Non-root where      │
│  agents                        │                  │    possible            │
│                                │                  │  • Minimal caps        │
│                                │                  │  • Read-only root      │
├────────────────────────────────┼──────────────────┼────────────────────────│
│  System daemon (node-exporter) │  Low-Medium      │  • Required caps only  │
│                                │                  │  • Host namespace      │
│                                │                  │    restrictions        │
│                                │                  │  • Minimal privileges  │
│                                                                             │
│  NEVER IN PRODUCTION:                                                       │
│  • Privileged containers (unless absolutely necessary)                      │
│  • Host network namespace (unless required)                                 │
│  • Host PID namespace (unless required)                                     │
│  • Running as root (unless absolutely necessary)                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Container Best Practices Summary                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  BUILD                                                                      │
│  ✓ Use multi-stage builds to minimize final image size                      │
│  ✓ Choose minimal base images (scratch, distroless, alpine)                 │
│  ✓ Pin versions for all dependencies and base images                        │
│  ✓ Scan images for vulnerabilities in CI/CD pipeline                        │
│  ✓ Sign images before pushing to registry                                   │
│  ✓ Use BuildKit for advanced features and security                          │
│  ✓ Layer caching: order commands by change frequency                        │
│  ✓ Remove build tools and dependencies from production                      │
│  ✓ Never embed secrets in images                                            │
│                                                                             │
│  RUNTIME                                                                    │
│  ✓ Always run as non-root user                                              │
│  ✓ Drop all capabilities, add only required ones                            │
│  ✓ Enable read-only root filesystem                                         │
│  ✓ Set resource limits (CPU, memory, disk)                                  │
│  ✓ Use seccomp and AppArmor/SELinux profiles                                │
│  ✓ Disable privilege escalation                                             │
│  ✓ Use specific user/group IDs (don't rely on names)                        │
│  ✓ Mount tmpfs for writable directories                                     │
│                                                                             │
│  SECURITY                                                                   │
│  ✓ Implement Pod Security Standards                                         │
│  ✓ Use network policies for traffic segmentation                            │
│  ✓ Enable container runtime security monitoring                             │
│  ✓ Audit all container activities                                           │
│  ✓ Regular vulnerability scanning                                           │
│  ✓ Secrets management outside containers                                    │
│  ✓ TLS for all inter-service communication                                  │
│                                                                             │
│  OPERATIONS                                                                 │
│  ✓ Health checks for all containers                                         │
│  ✓ Graceful shutdown handling (SIGTERM)                                     │
│  ✓ Structured logging                                                       │
│  ✓ Observability (metrics, traces)                                          │
│  ✓ Resource monitoring and alerting                                         │
│  ✓ Image update strategy (rolling updates)                                  │
│  ✓ Disaster recovery testing                                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## References

1. NIST SP 800-190 - Application Container Security Guide
2. CIS Docker Benchmark
3. OWASP Container Security Verification Standard
4. Docker Security Documentation
5. Kubernetes Security Best Practices
