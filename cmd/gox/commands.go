package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

// =============================================================================
// Gen Command - 代码生成
// =============================================================================

func runGen(args []string) {
	if len(args) == 0 {
		fmt.Println("❌ 请指定生成类型: handler, model, service, test")
		return
	}

	genType := args[0]
	name := ""
	if len(args) > 1 {
		name = args[1]
	}

	fmt.Printf("🔨 生成 %s: %s\n", genType, name)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	var err error
	switch genType {
	case "handler":
		err = genHandler(name)
	case "model":
		err = genModel(name)
	case "service":
		err = genService(name)
	case "test":
		err = genTest(name)
	case "middleware":
		err = genMiddleware(name)
	default:
		fmt.Printf("❌ 未知的生成类型: %s\n", genType)
		return
	}

	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n", err)
	}
}

func genHandler(name string) error {
	if name == "" {
		return fmt.Errorf("请指定handler名称")
	}

	name = strings.Title(name)
	tmpl := `package handlers

import (
	"encoding/json"
	"net/http"
)

// {{.Name}}Handler {{.Name}}处理器
type {{.Name}}Handler struct{}

// New{{.Name}}Handler 创建{{.Name}}处理器
func New{{.Name}}Handler() *{{.Name}}Handler {
	return &{{.Name}}Handler{}
}

// Handle{{.Name}} 处理{{.Name}}请求
func (h *{{.Name}}Handler) Handle{{.Name}}(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "success",
		"message": "{{.Name}} handler",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

import "time"

// {{.Name}} {{.Name}}模型
type {{.Name}} struct {
	ID        int64     ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// Validate 验证{{.Name}}
func (m *{{.Name}}) Validate() error {
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

import "context"

// {{.Name}}Service {{.Name}}服务
type {{.Name}}Service struct{}

// New{{.Name}}Service 创建{{.Name}}服务
func New{{.Name}}Service() *{{.Name}}Service {
	return &{{.Name}}Service{}
}

// Get{{.Name}} 获取{{.Name}}
func (s *{{.Name}}Service) Get{{.Name}}(ctx context.Context, id int64) error {
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

import "testing"

// Test{{.Name}} 测试{{.Name}}
func Test{{.Name}}(t *testing.T) {
	t.Run("基础测试", func(t *testing.T) {
		// TODO: 实现测试
		if false {
			t.Error("测试失败")
		}
	})
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

import "net/http"

// {{.Name}}Middleware {{.Name}}中间件
func {{.Name}}Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: 前置处理
		next.ServeHTTP(w, r)
		// TODO: 后置处理
	})
}
`
	return generateFromTemplate(tmpl, name, fmt.Sprintf("%s_middleware.go", strings.ToLower(name)))
}

func generateFromTemplate(tmplStr, name, filename string) error {
	tmpl, err := template.New("gen").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("模板解析失败: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	data := struct{ Name string }{Name: name}
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("模板执行失败: %w", err)
	}

	absPath, _ := filepath.Abs(filename)
	fmt.Printf("✅ 生成成功: %s\n", absPath)
	return nil
}

// =============================================================================
// Init Command - 项目初始化
// =============================================================================

func runInit(args []string) {
	projectName := "myproject"
	if len(args) > 0 {
		projectName = args[0]
	}

	fmt.Printf("🚀 初始化项目: %s\n", projectName)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	dirs := []string{
		"cmd/" + projectName,
		"pkg/handlers",
		"pkg/models",
		"internal/config",
		"api",
		"docs",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("❌ 创建目录失败: %v\n", err)
			return
		}
		fmt.Printf("✅ 创建目录: %s\n", dir)
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ 项目初始化完成!")
}

// =============================================================================
// Doctor Command - 健康检查
// =============================================================================

func runDoctor(args []string) {
	fmt.Println("🏥 系统健康检查...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Go版本检查
	fmt.Println("\n📋 Go环境检查")
	fmt.Printf("✅ Go版本: %s\n", runtime.Version())
	fmt.Printf("   GOOS: %s, GOARCH: %s\n", runtime.GOOS, runtime.GOARCH)

	// 项目结构检查
	fmt.Println("\n📋 项目结构检查")
	requiredFiles := []string{"go.mod", "go.work", "README.md"}
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); err == nil {
			fmt.Printf("✅ %s 存在\n", file)
		} else {
			fmt.Printf("⚠️  %s 不存在\n", file)
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ 健康检查完成！")
}

// =============================================================================
// Config Command - 配置管理
// =============================================================================

type Config struct {
	Project struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"project"`
}

func runConfig(args []string) {
	if len(args) == 0 {
		fmt.Println("⚙️  配置管理")
		fmt.Println("使用: gox config [init|list|get|set]")
		return
	}

	action := args[0]

	switch action {
	case "init":
		fmt.Println("✅ 配置初始化...")
		// 创建默认配置
		config := Config{}
		config.Project.Name = "myproject"
		config.Project.Version = "1.0.0"
		data, _ := json.MarshalIndent(config, "", "  ")
		os.WriteFile(".goxconfig.json", data, 0644)
		fmt.Println("✅ 配置已创建: .goxconfig.json")
	case "list":
		fmt.Println("⚙️  当前配置:")
		data, err := os.ReadFile(".goxconfig.json")
		if err != nil {
			fmt.Println("❌ 配置文件不存在")
			return
		}
		fmt.Println(string(data))
	default:
		fmt.Printf("未知操作: %s\n", action)
	}
}

// =============================================================================
// Bench Command - 基准测试
// =============================================================================

func runBench(args []string) {
	fmt.Println("⚡ 运行基准测试...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	benchArgs := []string{"test", "-bench=.", "-benchmem", "./..."}
	cmd := exec.Command("go", benchArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ 基准测试失败: %v\n", err)
		return
	}

	fmt.Println("\n✅ 基准测试完成！")
}

// =============================================================================
// Deps Command - 依赖管理
// =============================================================================

func runDeps(args []string) {
	action := "list"
	if len(args) > 0 {
		action = args[0]
	}

	fmt.Printf("📦 依赖管理: %s\n", action)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	var cmd *exec.Cmd
	switch action {
	case "list":
		cmd = exec.Command("go", "list", "-m", "all")
	case "tidy":
		cmd = exec.Command("go", "mod", "tidy")
	case "verify":
		cmd = exec.Command("go", "mod", "verify")
	case "update":
		cmd = exec.Command("go", "get", "-u", "./...")
	default:
		fmt.Printf("未知操作: %s\n", action)
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ 操作失败: %v\n", err)
		return
	}

	fmt.Println("\n✅ 操作完成！")
}
