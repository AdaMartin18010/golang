// Package openapi 提供 OpenAPI/Swagger UI 集成
//
// 设计原理：
// 1. 提供 Swagger UI 集成，方便查看和测试 API
// 2. 支持从 OpenAPI 规范文件动态加载
// 3. 支持开发和生产环境的不同配置
//
// 架构位置：
// - Interfaces Layer (internal/interfaces/http/openapi/)
// - 用于 HTTP 服务器集成 Swagger UI
//
// 使用场景：
// 1. 开发环境：查看 API 文档和测试 API
// 2. 测试环境：API 文档展示
// 3. 生产环境：可选的 API 文档展示（需要认证）
package openapi

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed swagger-ui/*
var swaggerUIFiles embed.FS

// SwaggerConfig Swagger UI 配置
type SwaggerConfig struct {
	// OpenAPISpecPath OpenAPI 规范文件路径
	OpenAPISpecPath string
	// Title 文档标题
	Title string
	// EnableUI 是否启用 UI（生产环境可禁用）
	EnableUI bool
	// RequireAuth 是否需要认证（生产环境建议启用）
	RequireAuth bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() SwaggerConfig {
	return SwaggerConfig{
		OpenAPISpecPath: "api/openapi/openapi.yaml",
		Title:           "API Documentation",
		EnableUI:        true,
		RequireAuth:     false,
	}
}

// Handler 返回 Swagger UI 处理器
//
// 功能说明：
// - 提供 Swagger UI 界面
// - 从 OpenAPI 规范文件加载 API 定义
// - 支持自定义配置
//
// 参数：
// - config: Swagger 配置
//
// 返回：
// - http.Handler: HTTP 处理器
func Handler(config SwaggerConfig) http.Handler {
	if !config.EnableUI {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Swagger UI is disabled", http.StatusNotFound)
		})
	}

	mux := http.NewServeMux()

	// Swagger UI 静态文件
	mux.Handle("/swagger-ui/", http.FileServer(http.FS(swaggerUIFiles)))

	// OpenAPI 规范文件
	mux.HandleFunc("/swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		if config.RequireAuth {
			// TODO: 添加认证检查
		}

		// 读取 OpenAPI 规范文件
		specPath := config.OpenAPISpecPath
		if !filepath.IsAbs(specPath) {
			// 相对路径，从项目根目录查找
			cwd, _ := os.Getwd()
			specPath = filepath.Join(cwd, specPath)
		}

		data, err := os.ReadFile(specPath)
		if err != nil {
			http.Error(w, "Failed to read OpenAPI spec", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/yaml")
		w.Write(data)
	})

	// Swagger UI 页面
	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if config.RequireAuth {
			// TODO: 添加认证检查
		}

		// 渲染 Swagger UI 页面
		tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <link rel="stylesheet" type="text/css" href="/swagger-ui/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="/swagger-ui/swagger-ui-bundle.js"></script>
    <script src="/swagger-ui/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger/openapi.yaml",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>
`
		t, err := template.New("swagger").Parse(tmpl)
		if err != nil {
			http.Error(w, "Failed to render Swagger UI", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, map[string]string{
			"Title": config.Title,
		})
	})

	return mux
}
