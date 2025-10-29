# Go 1.22 HTTP 路由增强详解

> **难度**: ⭐⭐⭐
> **标签**: #http #routing #ServeMux

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3

---

**版本**: v1.0  
**更新日期**: 2025-10-29  
**适用于**: Go 1.25.3

---

## 📋 目录

- [📋 概述](#概述)
- [🚀 新特性详解](#新特性详解)
  - [1. HTTP 方法匹配](#1.-http-方法匹配)
    - [基本语法](#基本语法)
    - [方法不匹配处理](#方法不匹配处理)
  - [2. 路径参数（通配符）](#2.-路径参数通配符)
    - [单段通配符 `{name}`](#单段通配符-name)
    - [多段通配符 `{path...}`](#多段通配符-path...)
  - [3. 路由优先级](#3.-路由优先级)
    - [优先级示例](#优先级示例)
- [🔍 详细特性](#详细特性)
  - [PathValue() 方法](#pathvalue-方法)
  - [组合方法和通配符](#组合方法和通配符)
- [🏗️ 实战案例](#实战案例)
  - [案例 1: RESTful API 服务器](#案例-1-restful-api-服务器)
    - [测试 API](#测试-api)
  - [案例 2: 文件服务器](#案例-2-文件服务器)
- [💡 最佳实践](#最佳实践)
  - [1. 使用常量定义路由](#1.-使用常量定义路由)
  - [2. 参数验证](#2.-参数验证)
  - [3. 中间件模式](#3.-中间件模式)
- [📚 扩展阅读](#扩展阅读)

## 📋 概述

Go 1.22 对 `net/http` 包的 `ServeMux` 进行了重大增强，引入了现代化的路由功能：

- ✅ **HTTP 方法匹配**（`GET`, `POST`, `PUT`, `DELETE` 等）
- ✅ **路径参数**（通配符 `{id}`, `{name}` 等）
- ✅ **路由优先级**（精确匹配优先）
- ✅ **向后兼容**（不影响现有代码）

这些改进让 Go 的标准库路由功能接近第三方框架（如 Gin, Echo），**无需额外依赖**。

---

## 🚀 新特性详解

### 1. HTTP 方法匹配

#### 基本语法

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    // 指定 HTTP 方法
    mux.HandleFunc("GET /users", listUsers)
    mux.HandleFunc("POST /users", createUser)
    mux.HandleFunc("PUT /users/{id}", updateUser)
    mux.HandleFunc("DELETE /users/{id}", deleteUser)
    
    http.ListenAndServe(":8080", mux)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "List all users")
}

func createUser(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Create user")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Update user: %s\n", id)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Delete user: %s\n", id)
}
```

#### 方法不匹配处理

```go
package main

import (
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    // 只允许 GET 方法
    mux.HandleFunc("GET /api/data", getData)
    
    http.ListenAndServe(":8080", mux)
}

func getData(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Data"))
}

// 测试:
// curl -X GET http://localhost:8080/api/data
// 输出: Data

// curl -X POST http://localhost:8080/api/data
// 输出: 405 Method Not Allowed
```

---

### 2. 路径参数（通配符）

#### 单段通配符 `{name}`

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    // 单段通配符
    mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
        id := r.PathValue("id")
        fmt.Fprintf(w, "User ID: %s\n", id)
    })
    
    mux.HandleFunc("GET /posts/{postID}/comments/{commentID}", func(w http.ResponseWriter, r *http.Request) {
        postID := r.PathValue("postID")
        commentID := r.PathValue("commentID")
        fmt.Fprintf(w, "Post: %s, Comment: %s\n", postID, commentID)
    })
    
    http.ListenAndServe(":8080", mux)
}

// 测试:
// curl http://localhost:8080/users/123
// 输出: User ID: 123

// curl http://localhost:8080/posts/456/comments/789
// 输出: Post: 456, Comment: 789
```

#### 多段通配符 `{path...}`

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    // 多段通配符（贪婪匹配）
    mux.HandleFunc("GET /files/{path...}", func(w http.ResponseWriter, r *http.Request) {
        path := r.PathValue("path")
        fmt.Fprintf(w, "File path: %s\n", path)
    })
    
    http.ListenAndServe(":8080", mux)
}

// 测试:
// curl http://localhost:8080/files/docs/api/v1/users.md
// 输出: File path: docs/api/v1/users.md

// curl http://localhost:8080/files/image.png
// 输出: File path: image.png
```

---

### 3. 路由优先级

Go 1.22 的路由按以下优先级匹配（从高到低）：

1. **精确匹配**
2. **最长路径匹配**
3. **带参数的路径**
4. **通配符路径**

#### 优先级示例

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    // 优先级 1: 精确匹配（最高）
    mux.HandleFunc("GET /users/admin", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Admin user (exact match)")
    })
    
    // 优先级 2: 带参数
    mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
        id := r.PathValue("id")
        fmt.Fprintf(w, "User ID: %s (parameterized)\n", id)
    })
    
    // 优先级 3: 前缀匹配（最低）
    mux.HandleFunc("GET /users/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Users prefix (lowest priority)")
    })
    
    http.ListenAndServe(":8080", mux)
}

// 测试:
// curl http://localhost:8080/users/admin
// 输出: Admin user (exact match)

// curl http://localhost:8080/users/123
// 输出: User ID: 123 (parameterized)

// curl http://localhost:8080/users/list/all
// 输出: Users prefix (lowest priority)
```

---

## 🔍 详细特性

### PathValue() 方法

`r.PathValue(name string)` 用于获取路径参数：

```go
package main

import (
    "encoding/json"
    "net/http"
)

type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("GET /api/v1/users/{userID}/posts/{postID}", func(w http.ResponseWriter, r *http.Request) {
        // 获取多个路径参数
        userID := r.PathValue("userID")
        postID := r.PathValue("postID")
        
        // 参数验证
        if userID == "" || postID == "" {
            http.Error(w, "Invalid parameters", http.StatusBadRequest)
            return
        }
        
        // 返回 JSON
        response := map[string]string{
            "user_id": userID,
            "post_id": postID,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })
    
    http.ListenAndServe(":8080", mux)
}
```

### 组合方法和通配符

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    // RESTful API 示例
    
    // 列表和创建
    mux.HandleFunc("GET /api/resources", listResources)
    mux.HandleFunc("POST /api/resources", createResource)
    
    // 单个资源操作
    mux.HandleFunc("GET /api/resources/{id}", getResource)
    mux.HandleFunc("PUT /api/resources/{id}", updateResource)
    mux.HandleFunc("PATCH /api/resources/{id}", patchResource)
    mux.HandleFunc("DELETE /api/resources/{id}", deleteResource)
    
    // 嵌套资源
    mux.HandleFunc("GET /api/resources/{id}/items", listResourceItems)
    mux.HandleFunc("POST /api/resources/{id}/items", createResourceItem)
    
    http.ListenAndServe(":8080", mux)
}

func listResources(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "List all resources")
}

func createResource(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Create resource")
}

func getResource(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Get resource: %s\n", id)
}

func updateResource(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Update resource: %s\n", id)
}

func patchResource(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Patch resource: %s\n", id)
}

func deleteResource(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Delete resource: %s\n", id)
}

func listResourceItems(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "List items for resource: %s\n", id)
}

func createResourceItem(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Create item for resource: %s\n", id)
}
```

---

## 🏗️ 实战案例

### 案例 1: RESTful API 服务器

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "sync"
)

// 数据模型
type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

// 内存存储
type TodoStore struct {
    mu    sync.RWMutex
    todos map[int]*Todo
    nextID int
}

func NewTodoStore() *TodoStore {
    return &TodoStore{
        todos: make(map[int]*Todo),
        nextID: 1,
    }
}

func (s *TodoStore) List() []*Todo {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    result := make([]*Todo, 0, len(s.todos))
    for _, todo := range s.todos {
        result = append(result, todo)
    }
    return result
}

func (s *TodoStore) Get(id int) (*Todo, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    todo, exists := s.todos[id]
    return todo, exists
}

func (s *TodoStore) Create(title string) *Todo {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    todo := &Todo{
        ID:        s.nextID,
        Title:     title,
        Completed: false,
    }
    s.todos[s.nextID] = todo
    s.nextID++
    return todo
}

func (s *TodoStore) Update(id int, title string, completed bool) (*Todo, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    todo, exists := s.todos[id]
    if !exists {
        return nil, false
    }
    
    todo.Title = title
    todo.Completed = completed
    return todo, true
}

func (s *TodoStore) Delete(id int) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    _, exists := s.todos[id]
    if exists {
        delete(s.todos, id)
    }
    return exists
}

// HTTP Handlers
type TodoAPI struct {
    store *TodoStore
}

func NewTodoAPI() *TodoAPI {
    return &TodoAPI{
        store: NewTodoStore(),
    }
}

func (api *TodoAPI) ListTodos(w http.ResponseWriter, r *http.Request) {
    todos := api.store.List()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func (api *TodoAPI) GetTodo(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    todo, exists := api.store.Get(id)
    if !exists {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func (api *TodoAPI) CreateTodo(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Title string `json:"title"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if input.Title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
    }
    
    todo := api.store.Create(input.Title)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func (api *TodoAPI) UpdateTodo(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    var input struct {
        Title     string `json:"title"`
        Completed bool   `json:"completed"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    todo, exists := api.store.Update(id, input.Title, input.Completed)
    if !exists {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func (api *TodoAPI) DeleteTodo(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    if !api.store.Delete(id) {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func main() {
    api := NewTodoAPI()
    
    mux := http.NewServeMux()
    
    // 路由配置
    mux.HandleFunc("GET /api/todos", api.ListTodos)
    mux.HandleFunc("GET /api/todos/{id}", api.GetTodo)
    mux.HandleFunc("POST /api/todos", api.CreateTodo)
    mux.HandleFunc("PUT /api/todos/{id}", api.UpdateTodo)
    mux.HandleFunc("DELETE /api/todos/{id}", api.DeleteTodo)
    
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", mux)
}
```

#### 测试 API

```bash
# 创建 todo
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go 1.22"}'

# 列出所有 todos
curl http://localhost:8080/api/todos

# 获取单个 todo
curl http://localhost:8080/api/todos/1

# 更新 todo
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go 1.22","completed":true}'

# 删除 todo
curl -X DELETE http://localhost:8080/api/todos/1
```

---

### 案例 2: 文件服务器

```go
package main

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
)

func main() {
    mux := http.NewServeMux()
    
    // 静态文件服务
    mux.HandleFunc("GET /static/{path...}", func(w http.ResponseWriter, r *http.Request) {
        path := r.PathValue("path")
        filePath := filepath.Join("./static", path)
        
        // 安全检查
        if !isPathSafe(filePath) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        
        // 检查文件存在
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            http.Error(w, "File not found", http.StatusNotFound)
            return
        }
        
        http.ServeFile(w, r, filePath)
    })
    
    fmt.Println("File server running on :8080")
    http.ListenAndServe(":8080", mux)
}

func isPathSafe(path string) bool {
    // 防止路径遍历攻击
    cleanPath := filepath.Clean(path)
    return filepath.HasPrefix(cleanPath, "./static")
}
```

---

## 💡 最佳实践

### 1. 使用常量定义路由

```go
package main

const (
    RouteUsers       = "GET /api/users"
    RouteUserByID    = "GET /api/users/{id}"
    RouteCreateUser  = "POST /api/users"
    RouteUpdateUser  = "PUT /api/users/{id}"
    RouteDeleteUser  = "DELETE /api/users/{id}"
)

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc(RouteUsers, listUsers)
    mux.HandleFunc(RouteUserByID, getUser)
    mux.HandleFunc(RouteCreateUser, createUser)
    mux.HandleFunc(RouteUpdateUser, updateUser)
    mux.HandleFunc(RouteDeleteUser, deleteUser)
}
```

### 2. 参数验证

```go
func getUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    
    // 验证参数
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    // 处理请求...
}
```

### 3. 中间件模式

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("%s %s\n", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /api/users", listUsers)
    
    // 应用中间件
    handler := loggingMiddleware(mux)
    http.ListenAndServe(":8080", handler)
}
```

---

## 📚 扩展阅读

- [Go 1.22 Release Notes - HTTP Routing](https://go.dev/doc/go1.22#net/http)
- [net/http ServeMux 文档](https://pkg.go.dev/net/http#ServeMux)
- [Go Blog: Routing Enhancements](https://go.dev/blog/routing-enhancements)

---

**文档维护者**: Go Documentation Team  
**最后更新**: 2025-10-29  
**文档状态**: ✅ 完成  
**适用版本**: Go 1.22+
