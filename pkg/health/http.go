package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// HTTPHandler 健康检查 HTTP 处理器
type HTTPHandler struct {
	checker *HealthChecker
}

// NewHTTPHandler 创建健康检查 HTTP 处理器
func NewHTTPHandler(checker *HealthChecker) *HTTPHandler {
	return &HTTPHandler{
		checker: checker,
	}
}

// LivenessHandler 存活探针处理器
// Kubernetes liveness probe: 检查应用是否存活
// 返回 200 OK 表示应用正在运行
func (h *HTTPHandler) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "alive",
			"timestamp": time.Now().Unix(),
		})
	}
}

// ReadinessHandler 就绪探针处理器
// Kubernetes readiness probe: 检查应用是否准备好接收流量
// 返回 200 OK 表示应用已就绪，503 表示未就绪
func (h *HTTPHandler) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		results := h.checker.Check(ctx)
		overallStatus := h.checker.OverallStatus(ctx)

		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"status":    string(overallStatus),
			"timestamp": time.Now().Unix(),
			"checks":    results,
		}

		statusCode := http.StatusOK
		if overallStatus == StatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		} else if overallStatus == StatusDegraded {
			statusCode = http.StatusOK // 降级但仍可服务
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}
}

// HealthHandler 综合健康检查处理器
// 返回详细的健康状态信息
func (h *HTTPHandler) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		results := h.checker.Check(ctx)
		overallStatus := h.checker.OverallStatus(ctx)

		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"status":    string(overallStatus),
			"timestamp": time.Now().Unix(),
			"checks":    results,
		}

		statusCode := http.StatusOK
		if overallStatus == StatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}
}

// StartupHandler 启动探针处理器
// Kubernetes startup probe: 检查应用是否启动完成
// 返回 200 OK 表示应用已启动
func (h *HTTPHandler) StartupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// 启动探针通常只检查关键依赖是否就绪
		results := h.checker.Check(ctx)
		overallStatus := h.checker.OverallStatus(ctx)

		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"status":    string(overallStatus),
			"timestamp": time.Now().Unix(),
			"checks":    results,
		}

		statusCode := http.StatusOK
		if overallStatus == StatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}
}
