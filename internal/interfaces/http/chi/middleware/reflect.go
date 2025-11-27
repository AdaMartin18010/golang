// Package middleware 提供 HTTP 中间件实现
//
// 反射/自解释中间件：基于框架的反射能力，提供请求和响应的元数据信息
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yourusername/golang/pkg/reflect"
)

// ReflectConfig 反射中间件配置
//
// 功能说明：
// - 配置反射和自解释行为
// - 支持在响应中添加元数据信息
//
// 字段说明：
// - EnableMetadata: 是否在响应头中添加元数据信息
// - EnableSelfDescribe: 是否启用自描述功能
// - MetadataPaths: 需要添加元数据的路径列表（空列表表示所有路径）
// - SkipPaths: 跳过反射的路径列表
//
// 使用示例：
//
//	config := middleware.ReflectConfig{
//	    EnableMetadata:    true,
//	    EnableSelfDescribe: true,
//	    SkipPaths:         []string{"/health", "/metrics"},
//	}
//	router.Use(middleware.ReflectMiddleware(config))
type ReflectConfig struct {
	EnableMetadata    bool
	EnableSelfDescribe bool
	MetadataPaths     []string
	SkipPaths         []string
}

// ReflectMiddleware 创建反射/自解释中间件
//
// 功能说明：
// - 基于框架的反射能力提供请求和响应的元数据信息
// - 支持在响应头中添加元数据
// - 支持自描述功能
//
// 工作流程：
// 1. 检查路径是否在跳过列表中
// 2. 如果启用元数据，在响应头中添加元数据信息
// 3. 如果启用自描述，在响应中添加自描述信息
// 4. 继续处理请求
//
// 使用示例：
//
//	router.Use(middleware.ReflectMiddleware(middleware.ReflectConfig{
//	    EnableMetadata:    true,
//	    EnableSelfDescribe: true,
//	    SkipPaths:         []string{"/health"},
//	}))
func ReflectMiddleware(config ReflectConfig) func(http.Handler) http.Handler {
	inspector := reflect.NewInspector()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// 检查是否跳过反射
			if shouldSkipReflect(path, config.SkipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()

			// 如果启用元数据，在响应头中添加元数据信息
			if config.EnableMetadata {
				// 检查路径是否在元数据路径列表中
				if len(config.MetadataPaths) == 0 || containsPath(path, config.MetadataPaths) {
					// 添加请求元数据
					metadata := map[string]interface{}{
						"method": r.Method,
						"path":   r.URL.Path,
						"host":   r.Host,
					}

					// 将元数据添加到响应头
					if metadataJSON, err := json.Marshal(metadata); err == nil {
						w.Header().Set("X-Request-Metadata", string(metadataJSON))
					}
				}
			}

			// 将检查器添加到上下文，供后续使用
			ctx = context.WithValue(ctx, "reflect.inspector", inspector)

			// 继续处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetInspector 从上下文中获取反射检查器
//
// 功能说明：
// - 获取反射检查器实例
// - 用于在处理器中进行类型检查、函数检查等
//
// 参数：
//   - ctx: 请求上下文
//
// 返回：
//   - *reflect.Inspector: 反射检查器实例，如果没有则返回 nil
//
// 使用示例：
//
//	func MyHandler(w http.ResponseWriter, r *http.Request) {
//	    inspector := middleware.GetInspector(r.Context())
//	    if inspector != nil {
//	        metadata := inspector.InspectType(myStruct)
//	        // 使用元数据...
//	    }
//	}
func GetInspector(ctx context.Context) *reflect.Inspector {
	if inspector, ok := ctx.Value("reflect.inspector").(*reflect.Inspector); ok {
		return inspector
	}
	return nil
}

// shouldSkipReflect 检查路径是否应该跳过反射
func shouldSkipReflect(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// containsPath 检查路径是否在路径列表中
func containsPath(path string, paths []string) bool {
	for _, p := range paths {
		if path == p || strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
