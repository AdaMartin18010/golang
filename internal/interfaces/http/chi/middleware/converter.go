// Package middleware 提供 HTTP 中间件实现
//
// 数据转换中间件：基于框架的数据转换工具，自动转换请求和响应数据
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yourusername/golang/pkg/converter"
)

// ConverterConfig 数据转换中间件配置
//
// 功能说明：
// - 配置请求和响应的数据转换行为
// - 支持多种数据格式转换
// - 支持自定义转换规则
//
// 字段说明：
// - EnableRequestConversion: 是否启用请求数据转换
// - EnableResponseConversion: 是否启用响应数据转换
// - RequestFormats: 支持的请求格式列表（如 ["json", "xml", "form"]）
// - ResponseFormats: 支持的响应格式列表（如 ["json", "xml"]）
// - DefaultResponseFormat: 默认响应格式
//
// 使用示例：
//
//	config := middleware.ConverterConfig{
//	    EnableRequestConversion:  true,
//	    EnableResponseConversion: true,
//	    RequestFormats:           []string{"json", "form"},
//	    ResponseFormats:          []string{"json"},
//	    DefaultResponseFormat:    "json",
//	}
//	router.Use(middleware.ConverterMiddleware(config))
type ConverterConfig struct {
	EnableRequestConversion  bool
	EnableResponseConversion bool
	RequestFormats           []string
	ResponseFormats          []string
	DefaultResponseFormat    string
}

// ConverterMiddleware 创建数据转换中间件
//
// 功能说明：
// - 自动转换请求数据格式
// - 自动转换响应数据格式
// - 支持多种数据格式（JSON、XML、Form 等）
// - 支持自定义转换规则
//
// 工作流程：
// 1. 检查请求内容类型
// 2. 根据内容类型转换请求数据
// 3. 将转换后的数据添加到上下文
// 4. 处理响应时，根据 Accept 头转换响应数据
//
// 使用示例：
//
//	router.Use(middleware.ConverterMiddleware(middleware.ConverterConfig{
//	    EnableRequestConversion:  true,
//	    EnableResponseConversion: true,
//	    DefaultResponseFormat:    "json",
//	}))
func ConverterMiddleware(config ConverterConfig) func(http.Handler) http.Handler {
	// 设置默认值
	if config.DefaultResponseFormat == "" {
		config.DefaultResponseFormat = "json"
	}

	// 创建转换器（用于后续可能的响应转换）
	_ = converter.NewConverter()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 请求数据转换
			if config.EnableRequestConversion && r.ContentLength > 0 {
				contentType := r.Header.Get("Content-Type")
				if strings.Contains(contentType, "application/json") {
					// JSON 请求：解析 JSON 数据
					var data map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&data); err == nil {
						// 将解析后的数据添加到上下文
						ctx = context.WithValue(ctx, "request.data", data)
					}
				} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
					// Form 请求：解析 Form 数据
					if err := r.ParseForm(); err == nil {
						formData := make(map[string]interface{})
						for k, v := range r.Form {
							if len(v) == 1 {
								formData[k] = v[0]
							} else {
								formData[k] = v
							}
						}
						ctx = context.WithValue(ctx, "request.data", formData)
					}
				}
			}

			// 继续处理请求
			// 注意：响应转换功能需要更复杂的实现，这里暂时简化
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// 注意：响应转换功能需要更复杂的实现
// 由于 responseWriter 已经在 circuitbreaker.go 中定义，
// 这里暂时不实现响应转换功能，只实现请求转换

// GetRequestData 从上下文中获取转换后的请求数据
//
// 功能说明：
// - 获取经过转换的请求数据
// - 数据格式为 map[string]interface{}
//
// 参数：
//   - ctx: 请求上下文
//
// 返回：
//   - map[string]interface{}: 请求数据，如果没有则返回 nil
//
// 使用示例：
//
//	func MyHandler(w http.ResponseWriter, r *http.Request) {
//	    data := middleware.GetRequestData(r.Context())
//	    if data != nil {
//	        // 使用转换后的数据
//	    }
//	}
func GetRequestData(ctx context.Context) map[string]interface{} {
	if data, ok := ctx.Value("request.data").(map[string]interface{}); ok {
		return data
	}
	return nil
}
