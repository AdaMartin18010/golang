# 任务文档生成器 (Task Documentation Generator)

> **分类**: 工程与云原生  
> **标签**: #documentation #code-generation #openapi

---

## 从代码生成文档

```go
type TaskDocGenerator struct {
    registry TaskRegistry
}

func (tdg *TaskDocGenerator) GenerateAll() (*TaskDocumentation, error) {
    tasks := tdg.registry.ListTasks()
    
    doc := &TaskDocumentation{
        Version:   "1.0.0",
        Generated: time.Now(),
        Tasks:     []TaskDoc{},
    }
    
    for _, task := range tasks {
        taskDoc := tdg.generateTaskDoc(task)
        doc.Tasks = append(doc.Tasks, taskDoc)
    }
    
    return doc, nil
}

func (tdg *TaskDocGenerator) generateTaskDoc(task TaskDefinition) TaskDoc {
    // 提取结构体字段作为参数文档
    params := tdg.extractParams(task.PayloadType)
    
    // 提取错误码
    errors := tdg.extractErrors(task.Handler)
    
    return TaskDoc{
        Name:        task.Name,
        Type:        task.Type,
        Description: tdg.extractDescription(task.Handler),
        Parameters:  params,
        Returns:     tdg.extractReturns(task.Handler),
        Errors:      errors,
        Examples:    tdg.generateExamples(task),
        Timeout:     task.DefaultTimeout,
        Retries:     task.DefaultRetries,
    }
}

func (tdg *TaskDocGenerator) extractParams(t reflect.Type) []ParamDoc {
    var params []ParamDoc
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        
        param := ParamDoc{
            Name:        field.Name,
            Type:        field.Type.String(),
            Required:    tdg.isRequired(field),
            Description: field.Tag.Get("doc"),
            Default:     field.Tag.Get("default"),
        }
        
        params = append(params, param)
    }
    
    return params
}
```

---

## OpenAPI 生成

```go
func (tdg *TaskDocGenerator) GenerateOpenAPI() (*OpenAPI, error) {
    tasks, _ := tdg.GenerateAll()
    
    api := &OpenAPI{
        OpenAPI: "3.0.0",
        Info: Info{
            Title:       "Task API",
            Version:     "1.0.0",
            Description: "Auto-generated task API documentation",
        },
        Paths: make(map[string]PathItem),
    }
    
    // 为每个任务生成端点
    for _, task := range tasks.Tasks {
        path := fmt.Sprintf("/tasks/%s", task.Type)
        
        api.Paths[path] = PathItem{
            Post: &Operation{
                Summary:     task.Description,
                Description: fmt.Sprintf("Execute %s task", task.Name),
                RequestBody: &RequestBody{
                    Content: map[string]MediaType{
                        "application/json": {
                            Schema: tdg.schemaFromParams(task.Parameters),
                        },
                    },
                },
                Responses: map[string]Response{
                    "200": {
                        Description: "Task submitted successfully",
                        Content: map[string]MediaType{
                            "application/json": {
                                Schema: Schema{Ref: "#/components/schemas/TaskResponse"},
                            },
                        },
                    },
                },
            },
        }
    }
    
    return api, nil
}
```

---

## Markdown 文档生成

```go
func (tdg *TaskDocGenerator) GenerateMarkdown(w io.Writer) error {
    tasks, _ := tdg.GenerateAll()
    
    tmpl := template.Must(template.New("docs").Parse(docTemplate))
    
    return tmpl.Execute(w, tasks)
}

const docTemplate = `# Task Documentation

> Generated at {{ .Generated }}

## 目录

{{ range .Tasks }}- [{{ .Name }}](#{{ .Type }})
{{ end }}

---

{{ range .Tasks }}
## {{ .Name }} {#{{ .Type }} }

{{ .Description }}

### 参数

| 名称 | 类型 | 必需 | 描述 | 默认值 |
|------|------|------|------|--------|
{{ range .Parameters }}| {{ .Name }} | {{ .Type }} | {{ .Required }} | {{ .Description }} | {{ .Default }} |
{{ end }}

### 返回值

{{ range .Returns }}- **{{ .Name }}** ({{ .Type }}): {{ .Description }}
{{ end }}

### 错误码

{{ range .Errors }}- **{{ .Code }}**: {{ .Description }}
{{ end }}

### 示例

{{ .Examples }}

---
{{ end }}
`
```

---

## 运行时文档

```go
// 内嵌文档处理器
type TaskDocHandler struct {
    generator *TaskDocGenerator
}

func (tdh *TaskDocHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    format := r.URL.Query().Get("format")
    
    switch format {
    case "openapi", "swagger":
        doc, _ := tdh.generator.GenerateOpenAPI()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(doc)
        
    case "markdown":
        w.Header().Set("Content-Type", "text/markdown")
        tdh.generator.GenerateMarkdown(w)
        
    default:
        doc, _ := tdh.generator.GenerateAll()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(doc)
    }
}
```
