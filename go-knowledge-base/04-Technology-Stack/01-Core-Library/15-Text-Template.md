# TS-CL-015: Go text/template Package

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #template #text-processing #code-generation
> **权威来源**:
>
> - [text/template Package](https://golang.org/pkg/text/template/) - Go standard library
> - [Go Templates](https://golang.org/pkg/html/template/) - HTML template

---

## 1. Template Basics

```go
package main

import (
    "bytes"
    "os"
    "strings"
    "text/template"
)

// Basic template execution
func basicTemplate() {
    // Define template
    tmpl := `Hello, {{.Name}}! You are {{.Age}} years old.`

    // Parse template
    t := template.Must(template.New("greeting").Parse(tmpl))

    // Execute with data
    data := struct {
        Name string
        Age  int
    }{
        Name: "John",
        Age:  30,
    }

    var buf bytes.Buffer
    if err := t.Execute(&buf, data); err != nil {
        panic(err)
    }

    println(buf.String()) // Hello, John! You are 30 years old.
}

// Template with conditionals
func conditionalTemplate() {
    tmpl := `
{{if .LoggedIn}}
Welcome back, {{.Username}}!
{{else}}
Please log in.
{{end}}
`
    t := template.Must(template.New("auth").Parse(tmpl))

    // Logged in user
    t.Execute(os.Stdout, struct {
        LoggedIn bool
        Username string
    }{true, "john"})

    // Guest
    t.Execute(os.Stdout, struct {
        LoggedIn bool
        Username string
    }{false, ""})
}

// Template with loops
func loopTemplate() {
    tmpl := `
Users:
{{range .}}
- {{.Name}} ({{.Email}})
{{end}}
`
    t := template.Must(template.New("users").Parse(tmpl))

    users := []struct {
        Name  string
        Email string
    }{
        {"Alice", "alice@example.com"},
        {"Bob", "bob@example.com"},
    }

    t.Execute(os.Stdout, users)
}

// Template with functions
func templateWithFunctions() {
    funcMap := template.FuncMap{
        "upper": strings.ToUpper,
        "lower": strings.ToLower,
    }

    tmpl := `{{.Name | upper}} - {{.Name | lower}}`
    t := template.Must(template.New("funcs").Funcs(funcMap).Parse(tmpl))

    t.Execute(os.Stdout, struct{ Name string }{"John Doe"})
    // Output: JOHN DOE - john doe
}
```

---

## 2. Advanced Templates

```go
// Named templates
func namedTemplates() {
    tmpl := `
{{define "header"}}
=== Header ===
{{end}}

{{define "footer"}}
=== Footer ===
{{end}}

{{template "header"}}
Content here
{{template "footer"}}
`
    t := template.Must(template.New("main").Parse(tmpl))
    t.Execute(os.Stdout, nil)
}

// Template inheritance with blocks
func blockTemplates() {
    base := `
{{define "base"}}
Header
{{block "content" .}}Default content{{end}}
Footer
{{end}}
`

    child := `
{{define "content"}}
Custom content for {{.Name}}
{{end}}
`

    t := template.New("base")
    template.Must(t.Parse(base))
    template.Must(t.Parse(child))

    t.ExecuteTemplate(os.Stdout, "base", struct{ Name string }{"John"})
}

// Nested templates with data passing
func nestedTemplates() {
    tmpl := `
{{define "list"}}
{{range .}}
{{template "item" .}}
{{end}}
{{end}}

{{define "item"}}
- {{.Name}}: {{.Value}}
{{end}}

{{template "list" .Items}}
`
    t := template.Must(template.New("nested").Parse(tmpl))

    data := struct {
        Items []struct{ Name, Value string }
    }{
        Items: []struct{ Name, Value string }{
            {"Key1", "Value1"},
            {"Key2", "Value2"},
        },
    }

    t.Execute(os.Stdout, data)
}

// Custom data structures
func customDataTemplate() {
    type Address struct {
        Street  string
        City    string
        Country string
    }

    type Person struct {
        Name    string
        Address Address
    }

    tmpl := `
Name: {{.Name}}
Address:
  Street: {{.Address.Street}}
  City: {{.Address.City}}
  Country: {{.Address.Country}}
`
    t := template.Must(template.New("person").Parse(tmpl))

    p := Person{
        Name: "John",
        Address: Address{
            Street:  "123 Main St",
            City:    "San Francisco",
            Country: "USA",
        },
    }

    t.Execute(os.Stdout, p)
}
```

---

## 3. Template Functions

```go
// Built-in functions
func builtInFunctions() {
    tmpl := `
{{/* Comments */}}
{{/* This is a comment */}}

{{/* Variables */}}
{{$name := .Name}}
{{$age := .Age}}
Name: {{$name}}
Age: {{$age}}

{{/* Comparison */}}
{{if eq .Name "John"}}Hello John!{{end}}
{{if ne .Name "Jane"}}Not Jane{{end}}
{{if gt .Age 18}}Adult{{end}}
{{if lt .Age 18}}Minor{{end}}
{{if ge .Age 21}}Can drink{{end}}
{{if le .Age 65}}Not retired{{end}}

{{/* Logical operators */}}
{{if and .Active .Verified}}Active and verified{{end}}
{{if or .Admin .Moderator}}Has permissions{{end}}
{{if not .Banned}}Not banned{{end}}

{{/* Default value */}}
{{.Nickname | default "No nickname"}}

{{/* Index */}}
{{index .Items 0}}
{{index .Map "key"}}

{{/* Length */}}
Total items: {{len .Items}}

{{/* Printf */}}
{{printf "%.2f" .Price}}

{{/* HTML escaping (html/template only) */}}
{{.HTMLContent}}
`
}

// Custom function map
func customFunctions() {
    funcMap := template.FuncMap{
        "add": func(a, b int) int {
            return a + b
        },
        "subtract": func(a, b int) int {
            return a - b
        },
        "formatDate": func(t time.Time) string {
            return t.Format("2006-01-02")
        },
        "join": strings.Join,
    }

    tmpl := `
Sum: {{add .A .B}}
Difference: {{subtract .A .B}}
Date: {{formatDate .Date}}
Tags: {{join .Tags ", "}}
`

    t := template.Must(template.New("custom").Funcs(funcMap).Parse(tmpl))

    data := struct {
        A     int
        B     int
        Date  time.Time
        Tags  []string
    }{
        A:     10,
        B:     5,
        Date:  time.Now(),
        Tags:  []string{"go", "template", "example"},
    }

    t.Execute(os.Stdout, data)
}
```

---

## 4. Template Best Practices

```go
// 1. Use template.ParseGlob for multiple templates
func parseMultiple() {
    t := template.Must(template.ParseGlob("templates/*.tmpl"))
    t.ExecuteTemplate(os.Stdout, "header.tmpl", nil)
}

// 2. Use template cache
var tmplCache = template.Must(template.New("").ParseGlob("templates/*"))

func renderWithCache(name string, data interface{}) (string, error) {
    var buf bytes.Buffer
    if err := tmplCache.ExecuteTemplate(&buf, name, data); err != nil {
        return "", err
    }
    return buf.String(), nil
}

// 3. Error handling
func safeTemplateExecution(tmpl *template.Template, data interface{}) (string, error) {
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("template execution failed: %w", err)
    }
    return buf.String(), nil
}

// 4. Validate templates at startup
func validateTemplates(pattern string) (*template.Template, error) {
    t, err := template.ParseGlob(pattern)
    if err != nil {
        return nil, fmt.Errorf("failed to parse templates: %w", err)
    }
    return t, nil
}
```

---

## 5. Checklist

```
Template Checklist:
□ Use html/template for HTML (XSS protection)
□ Use text/template for non-HTML
□ Cache parsed templates
□ Validate templates at startup
□ Handle template execution errors
□ Use functions for complex logic
□ Keep templates simple
□ Document custom functions
```
