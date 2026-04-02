# 平台工程 (Platform Engineering)

> **分类**: 成熟应用领域
> **标签**: #platform-engineering #developer-experience #internal-platform

---

## 内部开发者平台 (IDP)

### 平台架构

```
┌─────────────────────────────────────┐
│         Developer Portal            │
│  (Backstage / Port / Cortex)        │
└─────────────┬───────────────────────┘
              │
    ┌─────────┴─────────┐
    ↓                   ↓
┌──────────┐      ┌──────────┐
│ Platform │      │  Self-   │
│  APIs    │      │ Service  │
└────┬─────┘      └────┬─────┘
     │                 │
     └────────┬────────┘
              ↓
    ┌─────────────────────┐
    │  Infrastructure     │
    │  (K8s / Cloud)      │
    └─────────────────────┘
```

---

## Backstage 集成

### 实体描述

```yaml
# catalog-info.yaml
apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: my-service
  description: User management service
  tags:
    - go
    - microservice
  annotations:
    github.com/project-slug: org/my-service
    argocd/app-name: my-service
    grafana/dashboard-selector: "title = 'My Service'"
spec:
  type: service
  lifecycle: production
  owner: team-platform
  system: user-management
  dependsOn:
    - resource:postgres-db
    - component:auth-service
```

### Go 模板

```go
// 生成 Backstage 配置
type Component struct {
    APIVersion string   `yaml:"apiVersion"`
    Kind       string   `yaml:"kind"`
    Metadata   Metadata `yaml:"metadata"`
    Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
    Name        string            `yaml:"name"`
    Description string            `yaml:"description"`
    Tags        []string          `yaml:"tags"`
    Annotations map[string]string `yaml:"annotations"`
}

func GenerateCatalog(serviceName, description string) ([]byte, error) {
    component := Component{
        APIVersion: "backstage.io/v1alpha1",
        Kind:       "Component",
        Metadata: Metadata{
            Name:        serviceName,
            Description: description,
            Tags:        []string{"go", "microservice"},
            Annotations: map[string]string{
                "github.com/project-slug": fmt.Sprintf("org/%s", serviceName),
            },
        },
        Spec: Spec{
            Type:     "service",
            Lifecycle: "production",
            Owner:    "team-platform",
        },
    }

    return yaml.Marshal(component)
}
```

---

## 自助服务 API

```go
// 环境创建 API
type EnvironmentAPI struct {
    k8sClient kubernetes.Interface
    tfClient  terraform.Client
}

func (api *EnvironmentAPI) CreateEnvironment(ctx context.Context, req CreateEnvRequest) (*Environment, error) {
    // 1. 验证请求
    if err := api.validate(req); err != nil {
        return nil, err
    }

    // 2. 创建命名空间
    ns := &corev1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name: req.Name,
            Labels: map[string]string{
                "team":      req.Team,
                "environment": req.Type,
            },
        },
    }
    _, err := api.k8sClient.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
    if err != nil {
        return nil, err
    }

    // 3. 应用 Terraform 模块
    outputs, err := api.tfClient.Apply(ctx, terraform.ApplyRequest{
        Module:   "environments/standard",
        Vars: map[string]string{
            "namespace": req.Name,
            "team":      req.Team,
        },
    })

    // 4. 返回环境信息
    return &Environment{
        Name:      req.Name,
        Namespace: req.Name,
        Endpoints: outputs["endpoints"],
    }, nil
}
```

---

## 黄金路径 (Golden Paths)

```go
// 服务模板
type ServiceTemplate struct {
    Name        string
    Description string
    Tags        []string
    Parameters  []Parameter
}

var GoServiceTemplate = ServiceTemplate{
    Name:        "go-microservice",
    Description: "Standard Go microservice",
    Tags:        []string{"go", "grpc", "kubernetes"},
    Parameters: []Parameter{
        {
            Name:     "service_name",
            Type:     "string",
            Required: true,
        },
        {
            Name:    "enable_grpc",
            Type:    "boolean",
            Default: true,
        },
        {
            Name:    "enable_http",
            Type:    "boolean",
            Default: true,
        },
    },
}

// 脚手架生成
func (t *ServiceTemplate) Generate(params map[string]interface{}) (*Scaffold, error) {
    scaffold := &Scaffold{
        Files: make(map[string]string),
    }

    // 生成 main.go
    scaffold.Files["cmd/server/main.go"] = generateMainGo(params)

    // 生成 Dockerfile
    scaffold.Files["Dockerfile"] = generateDockerfile(params)

    // 生成 K8s 配置
    scaffold.Files["k8s/deployment.yaml"] = generateDeployment(params)

    // 生成 CI/CD
    scaffold.Files[".github/workflows/ci.yaml"] = generateCI(params)

    return scaffold, nil
}
```
