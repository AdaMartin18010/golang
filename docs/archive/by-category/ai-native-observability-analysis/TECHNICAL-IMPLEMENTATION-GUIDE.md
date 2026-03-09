# 技术实施指南

> 详细的技术实现方案、代码模板与最佳实践

---

## 一、项目结构

### 1.1 目录组织

```
ai-native-observability-analysis/
├── cmd/                          # 可执行命令
│   ├── analyzer/                 # 静态分析工具
│   │   └── main.go
│   ├── import/                   # 知识图谱导入工具
│   │   └── main.go
│   ├── mcp-server/               # MCP服务
│   │   └── main.go
│   └── nl-query/                 # 自然语言查询服务
│       └── main.go
│
├── pkg/                          # 核心库
│   ├── analyzer/                 # 分析器
│   │   ├── golang/               # Go代码分析
│   │   │   ├── ast.go
│   │   │   ├── interface.go
│   │   │   └── config.go
│   │   └── extractor.go
│   │
│   ├── graph/                    # 知识图谱客户端
│   │   ├── client.go
│   │   ├── models.go
│   │   └── queries.go
│   │
│   ├── vector/                   # 向量检索客户端
│   │   ├── client.go
│   │   └── embeddings.go
│   │
│   ├── mcp/                      # MCP协议实现
│   │   ├── server.go
│   │   └── tools.go
│   │
│   ├── nlp/                      # 自然语言处理
│   │   ├── parser.go
│   │   └── generator.go
│   │
│   └── llm/                      # LLM客户端封装
│       ├── client.go
│       └── prompts.go
│
├── internal/                     # 内部实现
│   ├── config/                   # 配置管理
│   ├── utils/                    # 工具函数
│   └── validators/               # 验证器
│
├── projects/                     # 项目分析配置
│   ├── otel-collector.yaml
│   ├── prometheus.yaml
│   └── templates/                # 分析模板
│
├── schema/                       # 知识图谱Schema
│   ├── entities.cypher
│   ├── relations.cypher
│   └── migrations/
│
├── data/                         # 分析数据输出
│   └── .gitignore
│
├── tests/                        # 测试
│   ├── integration/
│   ├── e2e/
│   └── fixtures/
│
├── docs/                         # 文档
│   └── (已创建的分析文档)
│
├── scripts/                      # 脚本
│   ├── setup.sh
│   ├── analyze.sh
│   └── deploy.sh
│
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

### 1.2 Go Module配置

```go
// go.mod
module github.com/yourorg/observability-analyzer

go 1.21

require (
    // Neo4j驱动
    github.com/neo4j/neo4j-go-driver/v5 v5.17.0

    // Weaviate客户端
    github.com/weaviate/weaviate-go-client/v4 v4.13.0

    // MCP协议
    github.com/metoro-io/mcp-golang v0.8.0

    // LLM客户端
    github.com/sashabaranov/go-openai v1.20.0
    github.com/anthropics/anthropic-sdk-go v0.1.0

    // Go分析工具
    golang.org/x/tools v0.19.0

    // Web框架
    github.com/gin-gonic/gin v1.9.1

    // 配置管理
    github.com/spf13/viper v1.18.2

    // 日志
    github.com/rs/zerolog v1.32.0

    // 测试
    github.com/stretchr/testify v1.9.0
    github.com/testcontainers/testcontainers-go v0.29.0
)
```

---

## 二、核心组件实现

### 2.1 Go AST分析器

```go
// pkg/analyzer/golang/ast.go
package golang

import (
    "go/ast"
    "go/parser"
    "go/token"
    "path/filepath"
    "strings"
)

// Analyzer Go代码分析器
type Analyzer struct {
    fset *token.FileSet
}

func NewAnalyzer() *Analyzer {
    return &Analyzer{
        fset: token.NewFileSet(),
    }
}

// AnalysisResult 分析结果
type AnalysisResult struct {
    Packages    []PackageInfo
    Interfaces  []InterfaceInfo
    Functions   []FunctionInfo
    Configs     []ConfigInfo
}

type PackageInfo struct {
    Path        string
    Name        string
    Imports     []string
    Files       []string
}

type InterfaceInfo struct {
    Package     string
    File        string
    Name        string
    Methods     []MethodInfo
    Doc         string
    Line        int
}

type MethodInfo struct {
    Name        string
    Signature   string
    Params      []ParamInfo
    Returns     []string
    Doc         string
}

type ParamInfo struct {
    Name        string
    Type        string
}

type FunctionInfo struct {
    Package     string
    File        string
    Name        string
    Signature   string
    Doc         string
    Line        int
}

// AnalyzeDirectory 分析整个目录
func (a *Analyzer) AnalyzeDirectory(root string) (*AnalysisResult, error) {
    result := &AnalysisResult{
        Packages:   make([]PackageInfo, 0),
        Interfaces: make([]InterfaceInfo, 0),
        Functions:  make([]FunctionInfo, 0),
        Configs:    make([]ConfigInfo, 0),
    }

    // 收集所有Go文件
    files, err := a.collectGoFiles(root)
    if err != nil {
        return nil, err
    }

    // 按包分组
    packages := make(map[string]*PackageInfo)

    for _, file := range files {
        f, err := parser.ParseFile(a.fset, file, nil, parser.ParseComments)
        if err != nil {
            continue // 跳过解析失败的文件
        }

        // 获取包信息
        dir := filepath.Dir(file)
        if _, ok := packages[dir]; !ok {
            pkg := &PackageInfo{
                Path:    dir,
                Name:    f.Name.Name,
                Imports: make([]string, 0),
                Files:   make([]string, 0),
            }

            // 提取imports
            for _, imp := range f.Imports {
                path := strings.Trim(imp.Path.Value, `"`)
                pkg.Imports = append(pkg.Imports, path)
            }

            packages[dir] = pkg
            result.Packages = append(result.Packages, *pkg)
        }

        packages[dir].Files = append(packages[dir].Files, file)

        // 分析文件内容
        a.analyzeFile(f, file, packages[dir].Name, result)
    }

    return result, nil
}

// collectGoFiles 收集所有Go文件
func (a *Analyzer) collectGoFiles(root string) ([]string, error) {
    var files []string

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // 跳过vendor和测试文件
        if strings.Contains(path, "vendor/") ||
           strings.Contains(path, "_test.go") ||
           strings.Contains(path, ".git/") {
            return nil
        }

        if !info.IsDir() && strings.HasSuffix(path, ".go") {
            files = append(files, path)
        }

        return nil
    })

    return files, err
}

// analyzeFile 分析单个文件
func (a *Analyzer) analyzeFile(f *ast.File, filename, pkgName string, result *AnalysisResult) {
    // 遍历AST
    ast.Inspect(f, func(n ast.Node) bool {
        switch x := n.(type) {
        case *ast.TypeSpec:
            // 分析类型定义
            a.analyzeTypeSpec(x, f, filename, pkgName, result)

        case *ast.FuncDecl:
            // 分析函数定义
            a.analyzeFuncDecl(x, filename, pkgName, result)
        }

        return true
    })
}

// analyzeTypeSpec 分析类型定义
func (a *Analyzer) analyzeTypeSpec(spec *ast.TypeSpec, f *ast.File, filename, pkgName string, result *AnalysisResult) {
    // 检查是否是接口类型
    if iface, ok := spec.Type.(*ast.InterfaceType); ok {
        interfaceInfo := InterfaceInfo{
            Package: pkgName,
            File:    filename,
            Name:    spec.Name.Name,
            Methods: make([]MethodInfo, 0),
            Line:    a.fset.Position(spec.Pos()).Line,
        }

        // 获取文档注释
        if spec.Doc != nil {
            interfaceInfo.Doc = spec.Doc.Text()
        }

        // 分析方法
        for _, method := range iface.Methods.List {
            if fnType, ok := method.Type.(*ast.FuncType); ok {
                for _, name := range method.Names {
                    methodInfo := MethodInfo{
                        Name:      name.Name,
                        Doc:       method.Doc.Text(),
                    }

                    // 提取签名
                    methodInfo.Signature = a.extractSignature(fnType)

                    // 提取参数
                    methodInfo.Params = a.extractParams(fnType)

                    interfaceInfo.Methods = append(interfaceInfo.Methods, methodInfo)
                }
            }
        }

        result.Interfaces = append(result.Interfaces, interfaceInfo)
    }

    // 检查是否是配置结构体
    if _, ok := spec.Type.(*ast.StructType); ok {
        // 检查结构体名是否包含Config
        if strings.Contains(spec.Name.Name, "Config") {
            a.analyzeConfigStruct(spec, f, filename, pkgName, result)
        }
    }
}

// analyzeFuncDecl 分析函数定义
func (a *Analyzer) analyzeFuncDecl(fn *ast.FuncDecl, filename, pkgName string, result *AnalysisResult) {
    // 跳过私有函数
    if !ast.IsExported(fn.Name.Name) {
        return
    }

    funcInfo := FunctionInfo{
        Package:   pkgName,
        File:      filename,
        Name:      fn.Name.Name,
        Line:      a.fset.Position(fn.Pos()).Line,
    }

    // 获取文档
    if fn.Doc != nil {
        funcInfo.Doc = fn.Doc.Text()
    }

    // 提取签名
    if fn.Type != nil {
        funcInfo.Signature = a.extractSignature(fn.Type)
    }

    result.Functions = append(result.Functions, funcInfo)
}

// extractSignature 提取函数签名
func (a *Analyzer) extractSignature(ft *ast.FuncType) string {
    var params, returns []string

    if ft.Params != nil {
        for _, param := range ft.Params.List {
            paramType := a.typeToString(param.Type)
            if len(param.Names) > 0 {
                for _, name := range param.Names {
                    params = append(params, name.Name+" "+paramType)
                }
            } else {
                params = append(params, paramType)
            }
        }
    }

    if ft.Results != nil {
        for _, result := range ft.Results.List {
            returns = append(returns, a.typeToString(result.Type))
        }
    }

    sig := "func(" + strings.Join(params, ", ") + ")"
    if len(returns) > 0 {
        if len(returns) == 1 {
            sig += " " + returns[0]
        } else {
            sig += " (" + strings.Join(returns, ", ") + ")"
        }
    }

    return sig
}

// extractParams 提取参数信息
func (a *Analyzer) extractParams(ft *ast.FuncType) []ParamInfo {
    var params []ParamInfo

    if ft.Params != nil {
        for _, param := range ft.Params.List {
            paramType := a.typeToString(param.Type)
            if len(param.Names) > 0 {
                for _, name := range param.Names {
                    params = append(params, ParamInfo{
                        Name: name.Name,
                        Type: paramType,
                    })
                }
            }
        }
    }

    return params
}

// typeToString 类型转字符串（简化版）
func (a *Analyzer) typeToString(expr ast.Expr) string {
    switch t := expr.(type) {
    case *ast.Ident:
        return t.Name
    case *ast.StarExpr:
        return "*" + a.typeToString(t.X)
    case *ast.ArrayType:
        return "[]" + a.typeToString(t.Elt)
    case *ast.SelectorExpr:
        return a.typeToString(t.X) + "." + t.Sel.Name
    case *ast.InterfaceType:
        return "interface{}"
    case *ast.FuncType:
        return "func"
    default:
        return "unknown"
    }
}
```

### 2.2 知识图谱客户端

```go
// pkg/graph/client.go
package graph

import (
    "context"
    "fmt"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Client Neo4j客户端
type Client struct {
    driver neo4j.DriverWithContext
    dbName string
}

// NewClient 创建新客户端
func NewClient(uri, username, password, dbName string) (*Client, error) {
    auth := neo4j.BasicAuth(username, password, "")

    driver, err := neo4j.NewDriverWithContext(uri, auth)
    if err != nil {
        return nil, fmt.Errorf("failed to create driver: %w", err)
    }

    // 验证连接
    ctx := context.Background()
    if err := driver.VerifyConnectivity(ctx); err != nil {
        return nil, fmt.Errorf("failed to verify connectivity: %w", err)
    }

    return &Client{
        driver: driver,
        dbName: dbName,
    }, nil
}

// Close 关闭连接
func (c *Client) Close(ctx context.Context) error {
    return c.driver.Close(ctx)
}

// CreateProject 创建项目节点
func (c *Client) CreateProject(ctx context.Context, project Project) error {
    session := c.driver.NewSession(ctx, neo4j.SessionConfig{
        DatabaseName: c.dbName,
    })
    defer session.Close(ctx)

    query := `
        CREATE (p:Project {
            id: $id,
            name: $name,
            description: $description,
            repository: $repository,
            license: $license,
            primary_language: $primary_language,
            ai_summary: $ai_summary,
            created_at: datetime()
        })
        RETURN p
    `

    params := map[string]any{
        "id":               project.ID,
        "name":             project.Name,
        "description":      project.Description,
        "repository":       project.Repository,
        "license":          project.License,
        "primary_language": project.PrimaryLanguage,
        "ai_summary":       project.AISummary,
    }

    _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
        result, err := tx.Run(ctx, query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume(ctx)
    })

    return err
}

// CreateComponent 创建组件节点
func (c *Client) CreateComponent(ctx context.Context, component Component) error {
    session := c.driver.NewSession(ctx, neo4j.SessionConfig{
        DatabaseName: c.dbName,
    })
    defer session.Close(ctx)

    query := `
        MATCH (p:Project {id: $project_id})
        CREATE (c:Component {
            id: $id,
            name: $name,
            type: $type,
            purpose: $purpose,
            stability: $stability,
            semantic_description: $semantic_description,
            source_location: $source_location
        })
        CREATE (c)-[:BELONGS_TO]->(p)
        RETURN c
    `

    params := map[string]any{
        "project_id":           component.ProjectID,
        "id":                   component.ID,
        "name":                 component.Name,
        "type":                 component.Type,
        "purpose":              component.Purpose,
        "stability":            component.Stability,
        "semantic_description": component.SemanticDescription,
        "source_location":      component.SourceLocation,
    }

    _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
        result, err := tx.Run(ctx, query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume(ctx)
    })

    return err
}

// CreateInterface 创建接口节点
func (c *Client) CreateInterface(ctx context.Context, iface Interface) error {
    session := c.driver.NewSession(ctx, neo4j.SessionConfig{
        DatabaseName: c.dbName,
    })
    defer session.Close(ctx)

    query := `
        MATCH (comp:Component {id: $component_id})
        CREATE (i:Interface {
            id: $id,
            name: $name,
            type: $type,
            signature: $signature,
            semantic_description: $semantic_description,
            example_usage: $example_usage,
            time_complexity: $time_complexity,
            space_complexity: $space_complexity
        })
        CREATE (comp)-[:IMPLEMENTS]->(i)
        RETURN i
    `

    signatureJSON, _ := json.Marshal(iface.Signature)

    params := map[string]any{
        "component_id":         iface.ComponentID,
        "id":                   iface.ID,
        "name":                 iface.Name,
        "type":                 iface.Type,
        "signature":            string(signatureJSON),
        "semantic_description": iface.SemanticDescription,
        "example_usage":        iface.ExampleUsage,
        "time_complexity":      iface.TimeComplexity,
        "space_complexity":     iface.SpaceComplexity,
    }

    _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
        result, err := tx.Run(ctx, query, params)
        if err != nil {
            return nil, err
        }
        return result.Consume(ctx)
    })

    return err
}

// QueryProject 查询项目信息
func (c *Client) QueryProject(ctx context.Context, projectID string) (*Project, error) {
    session := c.driver.NewSession(ctx, neo4j.SessionConfig{
        DatabaseName: c.dbName,
    })
    defer session.Close(ctx)

    query := `
        MATCH (p:Project {id: $id})
        RETURN p {
            .*,
            components: [(p)<-[:BELONGS_TO]-(c:Component) | c {.*}]
        } as project
    `

    result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
        result, err := tx.Run(ctx, query, map[string]any{"id": projectID})
        if err != nil {
            return nil, err
        }

        if result.Next(ctx) {
            record := result.Record()
            projectMap, _ := record.Get("project")
            return projectMap, nil
        }

        return nil, nil
    })

    if err != nil {
        return nil, err
    }

    if result == nil {
        return nil, fmt.Errorf("project not found: %s", projectID)
    }

    // 转换为Project结构体
    projectMap := result.(map[string]any)
    project := &Project{
        ID:              getString(projectMap, "id"),
        Name:            getString(projectMap, "name"),
        Description:     getString(projectMap, "description"),
        Repository:      getString(projectMap, "repository"),
        License:         getString(projectMap, "license"),
        PrimaryLanguage: getString(projectMap, "primary_language"),
        AISummary:       getString(projectMap, "ai_summary"),
    }

    return project, nil
}

// SearchComponents 搜索组件
func (c *Client) SearchComponents(ctx context.Context, query string, types []string) ([]Component, error) {
    session := c.driver.NewSession(ctx, neo4j.SessionConfig{
        DatabaseName: c.dbName,
    })
    defer session.Close(ctx)

    cypher := `
        MATCH (c:Component)
        WHERE c.name CONTAINS $query OR c.semantic_description CONTAINS $query
    `

    params := map[string]any{
        "query": query,
    }

    if len(types) > 0 {
        cypher += ` AND c.type IN $types`
        params["types"] = types
    }

    cypher += ` RETURN c LIMIT 20`

    result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
        result, err := tx.Run(ctx, cypher, params)
        if err != nil {
            return nil, err
        }

        var components []Component
        for result.Next(ctx) {
            record := result.Record()
            node, _ := record.Get("c")
            nodeMap := node.(neo4j.Node).Props

            components = append(components, Component{
                ID:                  getString(nodeMap, "id"),
                Name:                getString(nodeMap, "name"),
                Type:                getString(nodeMap, "type"),
                Purpose:             getString(nodeMap, "purpose"),
                SemanticDescription: getString(nodeMap, "semantic_description"),
            })
        }

        return components, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]Component), nil
}

// Helper函数
func getString(m map[string]any, key string) string {
    if v, ok := m[key]; ok {
        if s, ok := v.(string); ok {
            return s
        }
    }
    return ""
}
```

### 2.3 MCP服务实现

```go
// pkg/mcp/server.go
package mcp

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/yourorg/observability-analyzer/pkg/graph"
    mcpgolang "github.com/metoro-io/mcp-golang"
)

// Server MCP服务器
type Server struct {
    graphClient *graph.Client
    engine      *gin.Engine
    mcpServer   *mcpgolang.Server
}

// NewServer 创建MCP服务器
func NewServer(graphClient *graph.Client) *Server {
    engine := gin.Default()

    s := &Server{
        graphClient: graphClient,
        engine:      engine,
    }

    // 初始化MCP服务器
    s.mcpServer = mcpgolang.NewServer(
        mcpgolang.WithName("observability-analyzer"),
        mcpgolang.WithVersion("0.1.0"),
    )

    // 注册工具
    s.registerTools()

    // 设置路由
    engine.POST("/mcp", s.handleMCPRequest)
    engine.GET("/health", s.healthCheck)

    return s
}

// Run 启动服务器
func (s *Server) Run(addr string) error {
    return s.engine.Run(addr)
}

// handleMCPRequest 处理MCP请求
func (s *Server) handleMcpRequest(c *gin.Context) {
    var request mcpgolang.JSONRPCRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    response := s.mcpServer.HandleRequest(request)
    c.JSON(200, response)
}

// healthCheck 健康检查
func (s *Server) healthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "healthy",
        "version": "0.1.0",
    })
}

// registerTools 注册MCP工具
func (s *Server) registerTools() {
    // 工具1: 查询项目信息
    s.mcpServer.RegisterTool("query_project",
        "Query detailed information about a project by ID or name",
        func(args QueryProjectArgs) (*mcpgolang.ToolResponse, error) {
            ctx := context.Background()

            project, err := s.graphClient.QueryProject(ctx, args.ProjectID)
            if err != nil {
                return nil, fmt.Errorf("failed to query project: %w", err)
            }

            content := fmt.Sprintf(
                "Project: %s\nDescription: %s\nAI Summary: %s\nRepository: %s\nLicense: %s",
                project.Name,
                project.Description,
                project.AISummary,
                project.Repository,
                project.License,
            )

            return mcpgolang.NewToolResponse(
                mcpgolang.NewTextContent(content),
            ), nil
        })

    // 工具2: 搜索组件
    s.mcpServer.RegisterTool("search_components",
        "Search for components by name, type, or capability",
        func(args SearchComponentsArgs) (*mcpgolang.ToolResponse, error) {
            ctx := context.Background()

            components, err := s.graphClient.SearchComponents(ctx, args.Query, args.Types)
            if err != nil {
                return nil, fmt.Errorf("failed to search components: %w", err)
            }

            var content string
            if len(components) == 0 {
                content = "No components found matching your query."
            } else {
                content = fmt.Sprintf("Found %d components:\n\n", len(components))
                for _, comp := range components {
                    content += fmt.Sprintf(
                        "- %s (%s): %s\n",
                        comp.Name,
                        comp.Type,
                        comp.SemanticDescription,
                    )
                }
            }

            return mcpgolang.NewToolResponse(
                mcpgolang.NewTextContent(content),
            ), nil
        })

    // 工具3: 获取组件详情
    s.mcpServer.RegisterTool("get_component",
        "Get detailed information about a specific component including its interfaces and configurations",
        func(args GetComponentArgs) (*mcpgolang.ToolResponse, error) {
            ctx := context.Background()

            component, err := s.graphClient.QueryComponent(ctx, args.ComponentID)
            if err != nil {
                return nil, fmt.Errorf("failed to get component: %w", err)
            }

            content := fmt.Sprintf(
                "Component: %s\nType: %s\nPurpose: %s\n\nDescription: %s\n\nInterfaces:\n",
                component.Name,
                component.Type,
                component.Purpose,
                component.SemanticDescription,
            )

            for _, iface := range component.Interfaces {
                content += fmt.Sprintf("- %s: %s\n", iface.Name, iface.SemanticDescription)
            }

            content += "\nConfigurations:\n"
            for _, cfg := range component.Configurations {
                content += fmt.Sprintf("- %s: %s (default: %v)\n", cfg.Key, cfg.Description, cfg.DefaultValue)
            }

            return mcpgolang.NewToolResponse(
                mcpgolang.NewTextContent(content),
            ), nil
        })

    // 工具4: 解释指标
    s.mcpServer.RegisterTool("explain_metric",
        "Explain the meaning and usage of a runtime metric",
        func(args ExplainMetricArgs) (*mcpgolang.ToolResponse, error) {
            ctx := context.Background()

            metric, err := s.graphClient.QueryMetric(ctx, args.MetricName)
            if err != nil {
                return nil, fmt.Errorf("failed to get metric: %w", err)
            }

            content := fmt.Sprintf(
                "Metric: %s\nType: %s\n\nDescription: %s\n\nInterpretation:\n%s\n\nAffected by:\n",
                metric.Name,
                metric.MetricType,
                metric.SemanticDescription,
                metric.Interpretation,
            )

            for _, driver := range metric.ConfigurationDrivers {
                content += fmt.Sprintf("- %s\n", driver)
            }

            return mcpgolang.NewToolResponse(
                mcpgolang.NewTextContent(content),
            ), nil
        })

    // 工具5: 生成配置
    s.mcpServer.RegisterTool("generate_config",
        "Generate configuration for a specific use case",
        func(args GenerateConfigArgs) (*mcpgolang.ToolResponse, error) {
            config, err := s.generateConfig(args)
            if err != nil {
                return nil, fmt.Errorf("failed to generate config: %w", err)
            }

            return mcpgolang.NewToolResponse(
                mcpgolang.NewTextContent(config),
            ), nil
        })
}

// 参数结构体
type QueryProjectArgs struct {
    ProjectID string `json:"project_id" description:"The unique identifier of the project"`
}

type SearchComponentsArgs struct {
    Query string   `json:"query" description:"Search query string"`
    Types []string `json:"types,omitempty" description:"Filter by component types (e.g., receiver, processor, exporter)"`
}

type GetComponentArgs struct {
    ComponentID string `json:"component_id" description:"The unique identifier of the component"`
}

type ExplainMetricArgs struct {
    MetricName string `json:"metric_name" description:"The name of the metric to explain"`
}

type GenerateConfigArgs struct {
    UseCase    string            `json:"use_case" description:"The use case (e.g., high-throughput-traces)"`
    Project    string            `json:"project" description:"Target project (e.g., otel-collector)"`
    Constraints map[string]string `json:"constraints,omitempty" description:"Constraints like memory_limit, cpu_limit"`
}

// generateConfig 生成配置
func (s *Server) generateConfig(args GenerateConfigArgs) (string, error) {
    // 基于知识图谱和使用场景生成配置
    // 这是一个简化实现，实际应该查询知识图谱并基于最佳实践生成

    config := fmt.Sprintf(`# Generated configuration for %s
# Use case: %s

receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 1s
    send_batch_size: 1024

exporters:
  otlp:
    endpoint: localhost:4317
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
`, args.Project, args.UseCase)

    return config, nil
}
```

---

## 三、测试策略

### 3.1 单元测试

```go
// pkg/graph/client_test.go
package graph

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func TestClient_CreateAndQueryProject(t *testing.T) {
    ctx := context.Background()

    // 启动Neo4j测试容器
    req := testcontainers.ContainerRequest{
        Image:        "neo4j:5.15-community",
        ExposedPorts: []string{"7687/tcp"},
        Env: map[string]string{
            "NEO4J_AUTH": "neo4j/testpass",
        },
        WaitingFor: wait.ForListeningPort("7687/tcp"),
    }

    neo4jC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err)
    defer neo4jC.Terminate(ctx)

    // 获取连接信息
    host, err := neo4jC.Host(ctx)
    require.NoError(t, err)

    port, err := neo4jC.MappedPort(ctx, "7687")
    require.NoError(t, err)

    // 创建客户端
    client, err := NewClient(
        fmt.Sprintf("bolt://%s:%s", host, port.Port()),
        "neo4j",
        "testpass",
        "neo4j",
    )
    require.NoError(t, err)
    defer client.Close(ctx)

    // 测试创建项目
    project := Project{
        ID:              "test-project",
        Name:            "Test Project",
        Description:     "A test project",
        Repository:      "https://github.com/test/project",
        License:         "MIT",
        PrimaryLanguage: "Go",
        AISummary:       "This is a test project for unit testing",
    }

    err = client.CreateProject(ctx, project)
    require.NoError(t, err)

    // 测试查询项目
    retrieved, err := client.QueryProject(ctx, "test-project")
    require.NoError(t, err)

    assert.Equal(t, project.Name, retrieved.Name)
    assert.Equal(t, project.Description, retrieved.Description)
    assert.Equal(t, project.AISummary, retrieved.AISummary)
}
```

### 3.2 集成测试

```go
// tests/integration/end_to_end_test.go
package integration

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/stretchr/testify/suite"
    "github.com/yourorg/observability-analyzer/pkg/graph"
    "github.com/yourorg/observability-analyzer/pkg/mcp"
)

type IntegrationSuite struct {
    suite.Suite
    graphClient *graph.Client
    mcpServer   *mcp.Server
}

func (s *IntegrationSuite) SetupSuite() {
    ctx := context.Background()

    // 连接测试数据库
    var err error
    s.graphClient, err = graph.NewClient(
        "bolt://localhost:7687",
        "neo4j",
        "password",
        "neo4j",
    )
    s.Require().NoError(err)

    // 创建MCP服务器
    s.mcpServer = mcp.NewServer(s.graphClient)

    // 等待服务器就绪
    time.Sleep(1 * time.Second)
}

func (s *IntegrationSuite) TearDownSuite() {
    if s.graphClient != nil {
        s.graphClient.Close(context.Background())
    }
}

func (s *IntegrationSuite) TestQueryProject() {
    // 首先创建测试数据
    ctx := context.Background()
    project := graph.Project{
        ID:              "integration-test",
        Name:            "Integration Test Project",
        Description:     "For integration testing",
        Repository:      "https://github.com/test/integration",
        AISummary:       "A project for testing the MCP API",
    }

    err := s.graphClient.CreateProject(ctx, project)
    s.Require().NoError(err)

    // 构造MCP请求
    reqBody := map[string]interface{}{
        "jsonrpc": "2.0",
        "id":      1,
        "method":  "tools/call",
        "params": map[string]interface{}{
            "name": "query_project",
            "arguments": map[string]string{
                "project_id": "integration-test",
            },
        },
    }

    jsonBody, _ := json.Marshal(reqBody)

    // 发送请求
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/mcp", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")

    // 这里需要实际调用Gin路由
    // 简化处理，实际应该通过HTTP客户端访问运行的服务器

    s.T().Log("Query project test passed")
}

func TestIntegrationSuite(t *testing.T) {
    suite.Run(t, new(IntegrationSuite))
}
```

---

## 四、部署配置

### 4.1 Dockerfile

```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o analyzer ./cmd/analyzer
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mcp-server ./cmd/mcp-server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o import-tool ./cmd/import

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 复制二进制文件
COPY --from=builder /app/analyzer .
COPY --from=builder /app/mcp-server .
COPY --from=builder /app/import-tool .

# 暴露端口
EXPOSE 8080

# 默认运行MCP服务器
CMD ["./mcp-server"]
```

### 4.2 Kubernetes部署

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: observability-analyzer
  labels:
    app: observability-analyzer
spec:
  replicas: 2
  selector:
    matchLabels:
      app: observability-analyzer
  template:
    metadata:
      labels:
        app: observability-analyzer
    spec:
      containers:
      - name: mcp-server
        image: yourorg/observability-analyzer:latest
        ports:
        - containerPort: 8080
        env:
        - name: NEO4J_URI
          value: "bolt://neo4j:7687"
        - name: NEO4J_USER
          valueFrom:
            secretKeyRef:
              name: neo4j-credentials
              key: username
        - name: NEO4J_PASSWORD
          valueFrom:
            secretKeyRef:
              name: neo4j-credentials
              key: password
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: observability-analyzer
spec:
  selector:
    app: observability-analyzer
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

---

## 五、最佳实践

### 5.1 代码规范

1. **错误处理**
   - 使用`fmt.Errorf("context: %w", err)`包装错误
   - 在边界处（API、存储）记录错误日志
   - 不要忽略错误

2. **上下文管理**
   - 所有I/O操作接受`context.Context`
   - 设置合理的超时
   - 传递上下文而非创建新上下文

3. **资源管理**
   - 使用`defer`关闭资源
   - 连接池化（Neo4j驱动已内置）
   - 限制并发数

### 5.2 性能优化

1. **查询优化**
   - 为常用查询创建索引
   - 限制返回结果数
   - 使用参数化查询防止注入

2. **缓存策略**
   - 热点数据内存缓存
   - LLM响应缓存
   - 向量索引缓存

### 5.3 安全考虑

1. **访问控制**
   - Neo4j启用认证
   - API密钥管理
   - 审计日志

2. **数据保护**
   - 敏感数据脱敏
   - 传输加密
   - 最小权限原则

---

*本指南将持续更新，反映最新实现模式。*
