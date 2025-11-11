# Web应用模板

**版本**: v1.0
**更新日期**: 2025-10-29
**适用于**: Go 1.25.3

---

## 📋 目录

- [Web应用模板](#web应用模板)
  - [📋 目录](#-目录)
  - [1. 📖 Web应用结构](#1--web应用结构)
  - [🎯 模板渲染](#-模板渲染)
  - [🔐 会话管理](#-会话管理)
  - [📚 相关资源](#-相关资源)

---

## 1. 📖 Web应用结构

```
webapp/
├── cmd/web/
│   └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   └── templates/
├── static/
│   ├── css/
│   ├── js/
│   └── img/
├── templates/
│   ├── base.html
│   └── pages/
└── go.mod
```

---

## 🎯 模板渲染

```go
// internal/handlers/handlers.go
package handlers

import (
    "html/template"
    "net/http"
)

type App struct {
    templates *template.Template
}

func New() *App {
    tmpl := template.Must(template.ParseGlob("templates/*.html"))
    return &App{templates: tmpl}
}

func (app *App) Home(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "Title": "Home Page",
        "User":  "John Doe",
    }

    app.templates.ExecuteTemplate(w, "base.html", data)
}
```

```html
<!-- templates/base.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <h1>Welcome, {{.User}}!</h1>
    {{template "content" .}}
</body>
</html>
```

---

## 🔐 会话管理

```go
import "github.com/gorilla/sessions"

var store = sessions.NewCookieStore([]byte("secret-key"))

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session")
    session.Values["user_id"] = "123"
    session.Save(r, w)
}
```

---

## 📚 相关资源

- [Go Web Examples](https://gowebexamples.com/)

**下一步**: [04-CLI工具模板](./04-CLI工具模板.md)

---

**最后更新**: 2025-10-29
