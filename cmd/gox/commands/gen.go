package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// GenCommand 代码生成命令
func GenCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定生成类型: handler, model, service, test")
	}

	genType := args[0]
	name := ""
	if len(args) > 1 {
		name = args[1]
	}

	fmt.Printf("🔨 生成 %s: %s\n", genType, name)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	switch genType {
	case "handler":
		return genHandler(name)
	case "model":
		return genModel(name)
	case "service":
		return genService(name)
	case "test":
		return genTest(name)
	case "middleware":
		return genMiddleware(name)
	default:
		return fmt.Errorf("未知的生成类型: %s", genType)
	}
}

func genHandler(name string) error {
	if name == "" {
		return fmt.Errorf("请指定handler名称")
	}

	// 确保首字母大写
	name = strings.Title(name)

	tmpl := `package handlers

import (
	"encoding/json"
	"net/http"
)

// {{.Name}}Handler {{.Name}}处理器
type {{.Name}}Handler struct {
	// TODO: 添加依赖
}

// New{{.Name}}Handler 创建{{.Name}}处理器
func New{{.Name}}Handler() *{{.Name}}Handler {
	return &{{.Name}}Handler{}
}

// Handle{{.Name}} 处理{{.Name}}请求
func (h *{{.Name}}Handler) Handle{{.Name}}(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现业务逻辑
	response := map[string]interface{}{
		"status":  "success",
		"message": "{{.Name}} handler",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Validate{{.Name}} 验证{{.Name}}请求
func (h *{{.Name}}Handler) Validate{{.Name}}(r *http.Request) error {
	// TODO: 实现验证逻辑
	return nil
}
`

	return generateFromTemplate(tmpl, name, fmt.Sprintf("%s_handler.go", strings.ToLower(name)))
}

func genModel(name string) error {
	if name == "" {
		return fmt.Errorf("请指定model名称")
	}

	name = strings.Title(name)

	tmpl := `package models

import (
	"time"
)

// {{.Name}} {{.Name}}模型
type {{.Name}} struct {
	ID        int64     ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
	// TODO: 添加字段
}

// Validate 验证{{.Name}}
func (m *{{.Name}}) Validate() error {
	// TODO: 实现验证逻辑
	return nil
}

// BeforeCreate 创建前钩子
func (m *{{.Name}}) BeforeCreate() error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前钩子
func (m *{{.Name}}) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}
`

	return generateFromTemplate(tmpl, name, fmt.Sprintf("%s.go", strings.ToLower(name)))
}

func genService(name string) error {
	if name == "" {
		return fmt.Errorf("请指定service名称")
	}

	name = strings.Title(name)

	tmpl := `package services

import (
	"context"
	"fmt"
)

// {{.Name}}Service {{.Name}}服务
type {{.Name}}Service struct {
	// TODO: 添加依赖
}

// New{{.Name}}Service 创建{{.Name}}服务
func New{{.Name}}Service() *{{.Name}}Service {
	return &{{.Name}}Service{}
}

// Get{{.Name}} 获取{{.Name}}
func (s *{{.Name}}Service) Get{{.Name}}(ctx context.Context, id int64) error {
	// TODO: 实现获取逻辑
	return nil
}

// Create{{.Name}} 创建{{.Name}}
func (s *{{.Name}}Service) Create{{.Name}}(ctx context.Context) error {
	// TODO: 实现创建逻辑
	return nil
}

// Update{{.Name}} 更新{{.Name}}
func (s *{{.Name}}Service) Update{{.Name}}(ctx context.Context, id int64) error {
	// TODO: 实现更新逻辑
	return nil
}

// Delete{{.Name}} 删除{{.Name}}
func (s *{{.Name}}Service) Delete{{.Name}}(ctx context.Context, id int64) error {
	// TODO: 实现删除逻辑
	return nil
}
`

	return generateFromTemplate(tmpl, name, fmt.Sprintf("%s_service.go", strings.ToLower(name)))
}

func genTest(name string) error {
	if name == "" {
		return fmt.Errorf("请指定测试文件名")
	}

	name = strings.Title(name)

	tmpl := `package tests

import (
	"testing"
)

// Test{{.Name}} 测试{{.Name}}
func Test{{.Name}}(t *testing.T) {
	// TODO: 实现测试
	t.Run("基础测试", func(t *testing.T) {
		// Arrange
		
		// Act
		
		// Assert
		if false {
			t.Error("测试失败")
		}
	})
}

// Benchmark{{.Name}} 基准测试{{.Name}}
func Benchmark{{.Name}}(b *testing.B) {
	// Setup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// TODO: 实现基准测试
	}
}
`

	return generateFromTemplate(tmpl, name, fmt.Sprintf("%s_test.go", strings.ToLower(name)))
}

func genMiddleware(name string) error {
	if name == "" {
		return fmt.Errorf("请指定middleware名称")
	}

	name = strings.Title(name)

	tmpl := `package middleware

import (
	"net/http"
)

// {{.Name}}Middleware {{.Name}}中间件
func {{.Name}}Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: 前置处理
		
		// 调用下一个处理器
		next.ServeHTTP(w, r)
		
		// TODO: 后置处理
	})
}
`

	return generateFromTemplate(tmpl, name, fmt.Sprintf("%s_middleware.go", strings.ToLower(name)))
}

func generateFromTemplate(tmplStr, name, filename string) error {
	// 创建模板
	tmpl, err := template.New("gen").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("模板解析失败: %w", err)
	}

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 执行模板
	data := struct{ Name string }{Name: name}
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("模板执行失败: %w", err)
	}

	absPath, _ := filepath.Abs(filename)
	fmt.Printf("✅ 生成成功: %s\n", absPath)
	return nil
}
