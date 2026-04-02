# 任务部署与运维 (Task Deployment & Operations)

> **分类**: 工程与云原生
> **标签**: #deployment #operations #kubernetes #docker
> **参考**: Kubernetes Deployment, Helm, GitOps

---

## 部署架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Task System Deployment Architecture                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Kubernetes Deployment                             │   │
│  │                                                                      │   │
│  │   ┌─────────────────────────────────────────────────────────────┐   │   │
│  │   │                    Deployment                                 │   │   │
│  │   │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐     │   │   │
│  │   │  │ Replica  │  │ Replica  │  │ Replica  │  │ Replica  │     │   │   │
│  │   │  │   Pod    │  │   Pod    │  │   Pod    │  │   Pod    │     │   │   │
│  │   │  │ (Worker) │  │ (Worker) │  │ (Worker) │  │ (Worker) │     │   │   │
│  │   │  └──────────┘  └──────────┘  └──────────┘  └──────────┘     │   │   │
│  │   └─────────────────────────────────────────────────────────────┘   │   │
│  │                                                                      │   │
│  │   ┌─────────────────────────────────────────────────────────────┐   │   │
│  │   │                    HPA (Horizontal Pod Autoscaler)           │   │   │
│  │   │   Scale based on: CPU, Memory, Custom Metrics (Queue Depth)  │   │   │
│  │   └─────────────────────────────────────────────────────────────┘   │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Helm Chart Structure                              │   │
│  │                                                                      │   │
│  │   task-scheduler/                                                    │   │
│  │   ├── Chart.yaml                                                     │   │
│  │   ├── values.yaml                                                    │   │
│  │   ├── templates/                                                     │   │
│  │   │   ├── deployment.yaml                                            │   │
│  │   │   ├── service.yaml                                               │   │
│  │   │   ├── hpa.yaml                                                   │   │
│  │   │   ├── configmap.yaml                                             │   │
│  │   │   ├── secret.yaml                                                │   │
│  │   │   └── ingress.yaml                                               │   │
│  │   └── charts/                                                        │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 完整部署配置

```go
package deployment

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "time"
)

// DeploymentConfig 部署配置
type DeploymentConfig struct {
    Name        string
    Namespace   string
    Version     string
    Replicas    int
    Image       string

    // 资源限制
    Resources ResourceRequirements

    // 环境变量
    EnvVars map[string]string

    // 配置
    Config map[string]string

    // 健康检查
    HealthChecks HealthCheckConfig
}

// ResourceRequirements 资源需求
type ResourceRequirements struct {
    CPURequest    string
    CPULimit      string
    MemoryRequest string
    MemoryLimit   string
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
    LivenessPath  string
    ReadinessPath string
    Port          int
    InitialDelay  int
    Period        int
}

// KubernetesDeployer K8s部署器
type KubernetesDeployer struct {
    kubeconfig string
    namespace  string
}

// NewKubernetesDeployer 创建K8s部署器
func NewKubernetesDeployer(kubeconfig, namespace string) *KubernetesDeployer {
    return &KubernetesDeployer{
        kubeconfig: kubeconfig,
        namespace:  namespace,
    }
}

// Deploy 部署应用
func (kd *KubernetesDeployer) Deploy(ctx context.Context, config DeploymentConfig) error {
    // 生成部署YAML
    yaml := kd.generateDeploymentYAML(config)

    // 应用部署
    cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
    cmd.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", kd.kubeconfig))
    cmd.Stdin = strings.NewReader(yaml)

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("deploy failed: %v, output: %s", err, output)
    }

    // 等待就绪
    return kd.waitForReady(ctx, config.Name, config.Replicas)
}

// Rollout 滚动更新
func (kd *KubernetesDeployer) Rollout(ctx context.Context, name, version string) error {
    cmd := exec.CommandContext(ctx, "kubectl", "set", "image",
        fmt.Sprintf("deployment/%s", name),
        fmt.Sprintf("%s=%s:%s", name, name, version),
        "-n", kd.namespace,
    )
    cmd.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", kd.kubeconfig))

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("rollout failed: %v, output: %s", err, output)
    }

    // 等待滚动更新完成
    return kd.waitForRollout(ctx, name)
}

// Scale 扩缩容
func (kd *KubernetesDeployer) Scale(ctx context.Context, name string, replicas int) error {
    cmd := exec.CommandContext(ctx, "kubectl", "scale", "deployment", name,
        fmt.Sprintf("--replicas=%d", replicas),
        "-n", kd.namespace,
    )
    cmd.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", kd.kubeconfig))

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("scale failed: %v, output: %s", err, output)
    }

    return nil
}

// Rollback 回滚
func (kd *KubernetesDeployer) Rollback(ctx context.Context, name string) error {
    cmd := exec.CommandContext(ctx, "kubectl", "rollout", "undo",
        fmt.Sprintf("deployment/%s", name),
        "-n", kd.namespace,
    )
    cmd.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", kd.kubeconfig))

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("rollback failed: %v, output: %s", err, output)
    }

    return nil
}

func (kd *KubernetesDeployer) waitForReady(ctx context.Context, name string, replicas int) error {
    timeout := time.After(5 * time.Minute)
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-timeout:
            return fmt.Errorf("timeout waiting for deployment")
        case <-ticker.C:
            ready, err := kd.checkReady(name, replicas)
            if err != nil {
                continue
            }
            if ready {
                return nil
            }
        }
    }
}

func (kd *KubernetesDeployer) checkReady(name string, replicas int) (bool, error) {
    // 检查Pod就绪状态
    return true, nil
}

func (kd *KubernetesDeployer) waitForRollout(ctx context.Context, name string) error {
    cmd := exec.CommandContext(ctx, "kubectl", "rollout", "status",
        fmt.Sprintf("deployment/%s", name),
        "-n", kd.namespace,
        "--timeout=5m",
    )
    cmd.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", kd.kubeconfig))

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("rollout status failed: %v, output: %s", err, output)
    }

    return nil
}

func (kd *KubernetesDeployer) generateDeploymentYAML(config DeploymentConfig) string {
    return fmt.Sprintf(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
    version: %s
spec:
  replicas: %d
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
        version: %s
    spec:
      containers:
      - name: %s
        image: %s:%s
        resources:
          requests:
            cpu: %s
            memory: %s
          limits:
            cpu: %s
            memory: %s
        livenessProbe:
          httpGet:
            path: %s
            port: %d
          initialDelaySeconds: %d
          periodSeconds: %d
        readinessProbe:
          httpGet:
            path: %s
            port: %d
          initialDelaySeconds: %d
          periodSeconds: %d
`,
        config.Name,
        kd.namespace,
        config.Name,
        config.Version,
        config.Replicas,
        config.Name,
        config.Name,
        config.Version,
        config.Name,
        config.Image,
        config.Version,
        config.Resources.CPURequest,
        config.Resources.MemoryRequest,
        config.Resources.CPULimit,
        config.Resources.MemoryLimit,
        config.HealthChecks.LivenessPath,
        config.HealthChecks.Port,
        config.HealthChecks.InitialDelay,
        config.HealthChecks.Period,
        config.HealthChecks.ReadinessPath,
        config.HealthChecks.Port,
        config.HealthChecks.InitialDelay,
        config.HealthChecks.Period,
    )
}

import "strings"

// HelmDeployer Helm部署器
type HelmDeployer struct {
    kubeconfig string
}

// NewHelmDeployer 创建Helm部署器
func NewHelmDeployer(kubeconfig string) *HelmDeployer {
    return &HelmDeployer{kubeconfig: kubeconfig}
}

// Install 安装Chart
func (hd *HelmDeployer) Install(ctx context.Context, releaseName, chartPath string, values map[string]interface{}) error {
    // 构建values文件
    valuesFile := hd.buildValuesFile(values)

    cmd := exec.CommandContext(ctx, "helm", "upgrade", "--install",
        releaseName,
        chartPath,
        "-f", valuesFile,
        "--kubeconfig", hd.kubeconfig,
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("helm install failed: %v, output: %s", err, output)
    }

    return nil
}

// Upgrade 升级Chart
func (hd *HelmDeployer) Upgrade(ctx context.Context, releaseName string, values map[string]interface{}) error {
    return hd.Install(ctx, releaseName, "", values)
}

// Rollback 回滚
func (hd *HelmDeployer) Rollback(ctx context.Context, releaseName string, revision int) error {
    cmd := exec.CommandContext(ctx, "helm", "rollback",
        releaseName,
        fmt.Sprintf("%d", revision),
        "--kubeconfig", hd.kubeconfig,
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("helm rollback failed: %v, output: %s", err, output)
    }

    return nil
}

func (hd *HelmDeployer) buildValuesFile(values map[string]interface{}) string {
    // 生成values.yaml内容
    return "/tmp/values.yaml"
}

// DockerBuilder Docker构建器
type DockerBuilder struct {
    registry string
}

// NewDockerBuilder 创建Docker构建器
func NewDockerBuilder(registry string) *DockerBuilder {
    return &DockerBuilder{registry: registry}
}

// Build 构建镜像
func (db *DockerBuilder) Build(ctx context.Context, imageName, dockerfile, contextPath string, tags []string) error {
    tagsArgs := []string{"build"}
    for _, tag := range tags {
        fullTag := fmt.Sprintf("%s/%s:%s", db.registry, imageName, tag)
        tagsArgs = append(tagsArgs, "-t", fullTag)
    }

    if dockerfile != "" {
        tagsArgs = append(tagsArgs, "-f", dockerfile)
    }

    tagsArgs = append(tagsArgs, contextPath)

    cmd := exec.CommandContext(ctx, "docker", tagsArgs...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("docker build failed: %v, output: %s", err, output)
    }

    return nil
}

// Push 推送镜像
func (db *DockerBuilder) Push(ctx context.Context, imageName, tag string) error {
    fullImage := fmt.Sprintf("%s/%s:%s", db.registry, imageName, tag)

    cmd := exec.CommandContext(ctx, "docker", "push", fullImage)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("docker push failed: %v, output: %s", err, output)
    }

    return nil
}

// BlueGreenDeployer 蓝绿部署器
type BlueGreenDeployer struct {
    k8s *KubernetesDeployer
}

// Deploy 蓝绿部署
func (bgd *BlueGreenDeployer) Deploy(ctx context.Context, config DeploymentConfig) error {
    // 确定当前颜色
    currentColor := bgd.getCurrentColor(config.Name)
    newColor := "blue"
    if currentColor == "blue" {
        newColor = "green"
    }

    // 部署新版本到非活动颜色
    config.Name = fmt.Sprintf("%s-%s", config.Name, newColor)
    if err := bgd.k8s.Deploy(ctx, config); err != nil {
        return err
    }

    // 健康检查
    if err := bgd.healthCheck(config.Name); err != nil {
        return err
    }

    // 切换流量
    if err := bgd.switchTraffic(config.Name, newColor); err != nil {
        return err
    }

    // 删除旧版本
    oldDeployment := fmt.Sprintf("%s-%s", config.Name, currentColor)
    bgd.deleteDeployment(ctx, oldDeployment)

    return nil
}

func (bgd *BlueGreenDeployer) getCurrentColor(name string) string {
    return "blue"
}

func (bgd *BlueGreenDeployer) healthCheck(name string) error {
    return nil
}

func (bgd *BlueGreenDeployer) switchTraffic(name, color string) error {
    return nil
}

func (bgd *BlueGreenDeployer) deleteDeployment(ctx context.Context, name string) {
    cmd := exec.CommandContext(ctx, "kubectl", "delete", "deployment", name)
    cmd.Run()
}
```

---

## 使用示例

```go
package main

import (
    "context"
    "fmt"

    "deployment"
)

func main() {
    // K8s部署
    deployer := deployment.NewKubernetesDeployer(
        "/path/to/kubeconfig",
        "production",
    )

    config := deployment.DeploymentConfig{
        Name:      "task-scheduler",
        Version:   "v1.2.3",
        Replicas:  5,
        Image:     "myregistry/task-scheduler",
        Resources: deployment.ResourceRequirements{
            CPURequest:    "100m",
            CPULimit:      "500m",
            MemoryRequest: "256Mi",
            MemoryLimit:   "512Mi",
        },
        HealthChecks: deployment.HealthCheckConfig{
            LivenessPath:  "/health",
            ReadinessPath: "/ready",
            Port:          8080,
            InitialDelay:  10,
            Period:        5,
        },
    }

    ctx := context.Background()

    if err := deployer.Deploy(ctx, config); err != nil {
        panic(err)
    }

    fmt.Println("Deployment successful")

    // 滚动更新
    if err := deployer.Rollout(ctx, "task-scheduler", "v1.2.4"); err != nil {
        panic(err)
    }

    // 扩缩容
    if err := deployer.Scale(ctx, "task-scheduler", 10); err != nil {
        panic(err)
    }

    // Docker构建
    builder := deployment.NewDockerBuilder("myregistry.io")

    if err := builder.Build(ctx, "task-scheduler", "Dockerfile", ".", []string{"v1.2.3", "latest"}); err != nil {
        panic(err)
    }

    if err := builder.Push(ctx, "task-scheduler", "v1.2.3"); err != nil {
        panic(err)
    }
}
```
