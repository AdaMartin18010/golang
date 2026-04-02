# 文本模板 (Text Template)

> **分类**: 开源技术堆栈  
> **标签**: #template #text #code-generation

---

## 基础模板

```go
import "text/template"

const tmplStr = `
Hello, {{.Name}}!
You have {{.Count}} new messages.
`

type Data struct {
    Name  string
    Count int
}

func main() {
    tmpl, err := template.New("hello").Parse(tmplStr)
    if err != nil {
        log.Fatal(err)
    }
    
    data := Data{Name: "John", Count: 5}
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(buf.String())
}
```

---

## 条件与循环

```go
const emailTmpl = `
Subject: {{.Subject}}

Dear {{.User.Name}},

{{if .HasOrders}}
Your recent orders:
{{range .Orders}}
- {{.Product}}: ${{.Price}}
{{end}}
Total: ${{.Total}}
{{else}}
You have no recent orders.
{{end}}

{{if .VIP}}Thank you for being a VIP customer!{{end}}
`
```

---

## 函数

```go
funcMap := template.FuncMap{
    "upper": strings.ToUpper,
    "lower": strings.ToLower,
    "formatDate": func(t time.Time) string {
        return t.Format("2006-01-02")
    },
    "json": func(v interface{}) string {
        b, _ := json.Marshal(v)
        return string(b)
    },
}

tmpl, _ := template.New("email").Funcs(funcMap).Parse(emailTmpl)
```

---

## 代码生成

```go
const modelTmpl = `package {{.Package}}

{{range .Imports}}import "{{.}}"
{{end}}

type {{.Name}} struct {
{{range .Fields}}    {{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}}\"`" + `
{{end}}
}

{{range .Methods}}
func (m *{{$.Name}}) {{.Name}}({{.Params}}) {{.Return}} {
    {{.Body}}
}
{{end}}
`

type ModelData struct {
    Package string
    Imports []string
    Name    string
    Fields  []Field
    Methods []Method
}

func GenerateModel(data ModelData) (string, error) {
    tmpl, _ := template.New("model").Parse(modelTmpl)
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    
    // 格式化代码
    formatted, err := format.Source(buf.Bytes())
    if err != nil {
        return buf.String(), err
    }
    
    return string(formatted), nil
}
```

---

## 模板继承

```go
// base.tmpl
{{define "base"}}
<!DOCTYPE html>
<html>
<head>{{template "head" .}}</head>
<body>{{template "body" .}}</body>
</html>
{{end}}

// page.tmpl
{{template "base" .}}

{{define "head"}}
<title>{{.Title}}</title>
{{end}}

{{define "body"}}
<h1>{{.Heading}}</h1>
<p>{{.Content}}</p>
{{end}}
```

---

## 管道

```go
const pipeTmpl = `
{{.Name | upper}}
{{.CreatedAt | formatDate}}
{{.Price | printf "%.2f"}}
{{.Items | len}}
`

// 嵌套管道
{{.Description | html | js}}
```
