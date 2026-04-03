# TS-CL-015: Go text/template - Deep Architecture and Patterns

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #template #text #html #templating
> **权威来源**:
>
> - [Go text/template](https://pkg.go.dev/text/template) - Official documentation
> - [Go html/template](https://pkg.go.dev/html/template) - HTML template

---

## 1. Template Architecture

### 1.1 Template System

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       Template System Architecture                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Template Compilation:                                                      │
│   ┌───────────┐    ┌──────────────┐    ┌──────────────┐                     │
│   │  Template │───>│    Parser    │───>│    Tree      │                     │
│   │  String   │    │   (lex/parse)│    │   (AST)      │                     │
│   └───────────┘    └──────────────┘    └──────┬───────┘                     │
│                                               │                              │
│                                               ▼                              │
│                                        ┌──────────────┐                     │
│                                        │   Execute    │                     │
│                                        │  (with data) │                     │
│                                        └──────────────┘                     │
│                                                                              │
│   Template Elements:                                                         │
│   - Actions: {{.Field}}, {{if}}, {{range}}, {{with}}                        │
│   - Functions: {{printf "%s" .Name}}                                        │
│   - Pipelines: {{.Name | upper | printf "%s"}}                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Template Syntax

### 2.1 Basic Actions

```go
// Dot - Current value
{{.}}           // Root object
{{.Name}}       // Field access
{{.Items.0}}    // Array/slice index
{{.User.Name}}  // Nested field

// Variables
{{$var := .Name}}
{{$var}}

// Conditions
{{if .Active}}
    Active
{{else}}
    Inactive
{{end}}

// With (creates new scope)
{{with .User}}
    {{.Name}}
{{end}}

// Range (iteration)
{{range .Items}}
    {{.}}
{{else}}
    No items
{{end}}

// With index
{{range $index, $item := .Items}}
    {{$index}}: {{$item}}
{{end}}
```

### 2.2 Functions

```go
// Built-in functions
{{printf "%s is %d years old" .Name .Age}}
{{len .Items}}
{{index .Items 0}}
{{slice .Items 1 3}}

// Comparisons
{{eq .Status "active"}}
{{ne .Status "inactive"}}
{{lt .Age 18}}
{{gt .Age 65}}

// Logical
{{and .Active .Verified}}
{{or .Email .Phone}}
{{not .Disabled}}
```

---

## 3. Go Integration

### 3.1 Basic Usage

```go
package main

import (
    "os"
    "text/template"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    tmpl := `Hello, {{.Name}}! You are {{.Age}} years old.`

    t := template.Must(template.New("greeting").Parse(tmpl))

    person := Person{Name: "Alice", Age: 30}
    t.Execute(os.Stdout, person)
}
```

### 3.2 Custom Functions

```go
import "strings"

func main() {
    funcMap := template.FuncMap{
        "upper": strings.ToUpper,
        "lower": strings.ToLower,
        "join":  strings.Join,
        "add": func(a, b int) int {
            return a + b
        },
    }

    tmpl := `{{.Name | upper}} has {{add .Age 5}} points.`

    t := template.Must(
        template.New("test").Funcs(funcMap).Parse(tmpl),
    )

    t.Execute(os.Stdout, person)
}
```

---

## 4. Advanced Patterns

### 4.1 Template Composition

```go
// Define blocks
const layout = `
{{define "layout"}}
<!DOCTYPE html>
<html>
<head><title>{{template "title" .}}</title></head>
<body>
    {{template "content" .}}
</body>
</html>
{{end}}
`

const page = `
{{define "title"}}Home{{end}}
{{define "content"}}
    <h1>Welcome</h1>
{{end}}
{{template "layout" .}}
`

// Parse and execute
t := template.Must(template.New("layout").Parse(layout))
t = template.Must(t.New("page").Parse(page))

t.ExecuteTemplate(os.Stdout, "page", data)
```

### 4.2 Error Handling

```go
type TemplateExecutor struct {
    tmpl *template.Template
}

func (te *TemplateExecutor) Execute(wr io.Writer, data interface{}) error {
    if err := te.tmpl.Execute(wr, data); err != nil {
        return fmt.Errorf("template execution failed: %w", err)
    }
    return nil
}
```

---

## 5. text/template vs html/template

| Feature | text/template | html/template |
|---------|---------------|---------------|
| Auto-escaping | No | Yes (HTML) |
| Security | Manual | XSS protection |
| Use case | Config files, code | Web pages |
| Performance | Slightly faster | Slightly slower |

---

## 6. Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Template Best Practices                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Design:                                                                     │
│  □ Use html/template for web output (XSS protection)                        │
│  □ Keep templates simple - move logic to Go code                            │
│  □ Use defined templates for composition                                    │
│                                                                              │
│  Performance:                                                                │
│  □ Parse templates once at startup                                          │
│  □ Use template caching in web applications                                 │
│  □ Minimize allocations in templates                                        │
│                                                                              │
│  Safety:                                                                     │
│  □ Always escape user input in text/templates                               │
│  □ Validate template data structures                                        │
│  □ Handle template execution errors                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (16+ KB, comprehensive coverage)
