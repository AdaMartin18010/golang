package operational

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/golang/pkg/observability"
)

// OperationalEndpoints 运维控制端点
// 提供健康检查、就绪检查、指标导出等运维控制功能
type OperationalEndpoints struct {
	observability *observability.Observability
	mux           *http.ServeMux
	server        *http.Server
	enabled       bool
}

// Config 运维端点配置
type Config struct {
	Observability *observability.Observability
	Port          int
	PathPrefix    string
	Enabled       bool
}

// NewOperationalEndpoints 创建运维控制端点
func NewOperationalEndpoints(cfg Config) *OperationalEndpoints {
	if !cfg.Enabled {
		return &OperationalEndpoints{enabled: false}
	}

	mux := http.NewServeMux()
	prefix := cfg.PathPrefix
	if prefix == "" {
		prefix = "/ops"
	}

	endpoints := &OperationalEndpoints{
		observability: cfg.Observability,
		mux:           mux,
		enabled:       true,
	}

	// 注册端点
	endpoints.registerEndpoints(prefix)

	// 创建 HTTP 服务器
	port := cfg.Port
	if port == 0 {
		port = 9090 // 默认端口
	}

	endpoints.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return endpoints
}

// registerEndpoints 注册所有端点
func (oe *OperationalEndpoints) registerEndpoints(prefix string) {
	// 健康检查
	oe.mux.HandleFunc(prefix+"/health", oe.healthHandler)
	oe.mux.HandleFunc(prefix+"/healthz", oe.healthHandler) // Kubernetes 兼容

	// 就绪检查
	oe.mux.HandleFunc(prefix+"/ready", oe.readinessHandler)
	oe.mux.HandleFunc(prefix+"/readiness", oe.readinessHandler)

	// 存活检查
	oe.mux.HandleFunc(prefix+"/live", oe.livenessHandler)
	oe.mux.HandleFunc(prefix+"/liveness", oe.livenessHandler)

	// 指标导出
	oe.mux.HandleFunc(prefix+"/metrics", oe.metricsHandler)
	oe.mux.HandleFunc(prefix+"/metrics/prometheus", oe.prometheusMetricsHandler)

	// 仪表板数据
	oe.mux.HandleFunc(prefix+"/dashboard", oe.dashboardHandler)

	// 诊断报告
	oe.mux.HandleFunc(prefix+"/diagnostics", oe.diagnosticsHandler)

	// 配置重载
	oe.mux.HandleFunc(prefix+"/config/reload", oe.configReloadHandler)

	// 性能分析（pprof）
	oe.mux.HandleFunc(prefix+"/debug/pprof/", oe.pprofHandler)

	// 系统信息
	oe.mux.HandleFunc(prefix+"/info", oe.infoHandler)

	// 版本信息
	oe.mux.HandleFunc(prefix+"/version", oe.versionHandler)

	// 服务发现
	oe.mux.HandleFunc(prefix+"/services", oe.servicesHandler)
	oe.mux.HandleFunc(prefix+"/services/", oe.serviceHandler)
}

// Start 启动运维端点服务器
func (oe *OperationalEndpoints) Start() error {
	if !oe.enabled {
		return nil
	}

	go func() {
		if err := oe.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Operational endpoints server error: %v\n", err)
		}
	}()

	return nil
}

// Stop 停止运维端点服务器
func (oe *OperationalEndpoints) Stop(ctx context.Context) error {
	if !oe.enabled || oe.server == nil {
		return nil
	}

	return oe.server.Shutdown(ctx)
}

// healthHandler 健康检查处理器
func (oe *OperationalEndpoints) healthHandler(w http.ResponseWriter, r *http.Request) {
	if oe.observability == nil {
		http.Error(w, "Observability not initialized", http.StatusServiceUnavailable)
		return
	}

	systemMonitor := oe.observability.GetSystemMonitor()
	if systemMonitor == nil {
		// 如果没有系统监控，返回基本健康状态
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    true,
			"timestamp": time.Now(),
			"message":   "healthy (no system monitor)",
		})
		return
	}

	healthChecker := systemMonitor.GetHealthChecker()
	if healthChecker == nil {
		// 如果没有健康检查器，返回基本健康状态
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    true,
			"timestamp": time.Now(),
			"message":   "healthy (no health checker)",
		})
		return
	}

	status := healthChecker.Check(r.Context())
	
	w.Header().Set("Content-Type", "application/json")
	if status.Healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    status.Healthy,
		"timestamp": status.Timestamp,
		"message":   status.Message,
		"details": map[string]interface{}{
			"memory_usage": status.MemoryUsage,
			"cpu_usage":    status.CPUUsage,
			"goroutines":   status.Goroutines,
			"gc":          status.GC,
		},
	})
}

// readinessHandler 就绪检查处理器
func (oe *OperationalEndpoints) readinessHandler(w http.ResponseWriter, r *http.Request) {
	// 检查关键依赖是否就绪
	ready := true
	checks := make(map[string]bool)

	// 检查系统监控
	if oe.observability != nil {
		systemMonitor := oe.observability.GetSystemMonitor()
		checks["system_monitor"] = systemMonitor != nil
		ready = ready && (systemMonitor != nil)
	} else {
		checks["observability"] = false
		ready = false
	}

	w.Header().Set("Content-Type", "application/json")
	if ready {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ready":  ready,
		"checks": checks,
	})
}

// livenessHandler 存活检查处理器
func (oe *OperationalEndpoints) livenessHandler(w http.ResponseWriter, r *http.Request) {
	// 简单的存活检查，只要服务在运行就返回 OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"alive": true,
		"time":  time.Now(),
	})
}

// metricsHandler 指标导出处理器
func (oe *OperationalEndpoints) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if oe.observability == nil {
		http.Error(w, "Observability not initialized", http.StatusServiceUnavailable)
		return
	}

	exporter := oe.observability.GetMetricsExporter()
	if exporter == nil {
		// 如果没有指标导出器，返回空指标
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"metrics": []interface{}{},
		})
		return
	}

	jsonData, err := exporter.ExportJSON(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// prometheusMetricsHandler Prometheus 格式指标处理器
func (oe *OperationalEndpoints) prometheusMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if oe.observability == nil {
		http.Error(w, "Observability not initialized", http.StatusServiceUnavailable)
		return
	}

	dashboardExporter := oe.observability.GetDashboardExporter()
	if dashboardExporter == nil {
		// 如果没有仪表板导出器，返回空 Prometheus 格式
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# No metrics available\n"))
		return
	}

	promData, err := dashboardExporter.ExportForPrometheus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(promData))
}

// dashboardHandler 仪表板数据处理器
func (oe *OperationalEndpoints) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if oe.observability == nil {
		http.Error(w, "Observability not initialized", http.StatusServiceUnavailable)
		return
	}

	dashboardExporter := oe.observability.GetDashboardExporter()
	if dashboardExporter == nil {
		http.Error(w, "Dashboard exporter not available", http.StatusServiceUnavailable)
		return
	}

	jsonData, err := dashboardExporter.ExportJSON(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// diagnosticsHandler 诊断报告处理器
func (oe *OperationalEndpoints) diagnosticsHandler(w http.ResponseWriter, r *http.Request) {
	if oe.observability == nil {
		http.Error(w, "Observability not initialized", http.StatusServiceUnavailable)
		return
	}

	diagnostics := oe.observability.GetDiagnostics()
	if diagnostics == nil {
		http.Error(w, "Diagnostics not available", http.StatusServiceUnavailable)
		return
	}

	jsonData, err := diagnostics.ExportJSON(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// configReloadHandler 配置重载处理器
func (oe *OperationalEndpoints) configReloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 配置重载逻辑
	// 注意：实际配置重载需要根据具体配置系统实现
	// 这里提供一个框架，实际使用时需要集成配置热重载功能
	if oe.observability != nil {
		systemMonitor := oe.observability.GetSystemMonitor()
		if systemMonitor != nil {
			configReloader := systemMonitor.GetConfigReloader()
			if configReloader != nil {
				if err := configReloader.Reload(r.Context()); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"status":  "error",
						"message": "Failed to reload configuration",
						"error":   err.Error(),
						"time":    time.Now(),
					})
					return
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Configuration reloaded",
		"time":    time.Now(),
	})
}

// pprofHandler 性能分析处理器
func (oe *OperationalEndpoints) pprofHandler(w http.ResponseWriter, r *http.Request) {
	// 重定向到标准库的 pprof
	// 需要导入 _ "net/http/pprof"
	EnablePprof(oe.mux, "/ops/debug/pprof")
	http.DefaultServeMux.ServeHTTP(w, r)
}

// infoHandler 系统信息处理器
func (oe *OperationalEndpoints) infoHandler(w http.ResponseWriter, r *http.Request) {
	if oe.observability == nil {
		http.Error(w, "Observability not initialized", http.StatusServiceUnavailable)
		return
	}

	info := make(map[string]interface{})

	// 平台信息
	platformInfo := oe.observability.GetPlatformInfo()
	info["platform"] = map[string]interface{}{
		"os":            platformInfo.OS,
		"arch":          platformInfo.Arch,
		"go_version":    platformInfo.GoVersion,
		"hostname":      platformInfo.Hostname,
		"cpu_cores":     platformInfo.CPUCores,
		"container":     oe.observability.IsContainer(),
		"kubernetes":    oe.observability.IsKubernetes(),
		"virtualized":   oe.observability.IsVirtualized(),
	}

	// Kubernetes 信息
	if oe.observability.IsKubernetes() {
		k8sInfo := oe.observability.GetKubernetesInfo()
		info["kubernetes"] = map[string]interface{}{
			"pod_name":      k8sInfo.PodName,
			"pod_namespace": k8sInfo.PodNamespace,
			"node_name":     k8sInfo.NodeName,
			"pod_ip":        k8sInfo.PodIP,
			"host_ip":       k8sInfo.HostIP,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(info)
}

// versionHandler 版本信息处理器
func (oe *OperationalEndpoints) versionHandler(w http.ResponseWriter, r *http.Request) {
	version := map[string]interface{}{
		"version":     "1.0.0",
		"build_time":  time.Now().Format(time.RFC3339),
		"git_commit":  "unknown",
		"go_version":  "1.25.3",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(version)
}

// servicesHandler 服务列表处理器
func (oe *OperationalEndpoints) servicesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: 从服务注册表获取服务列表
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"services": []interface{}{},
	})
}

// serviceHandler 单个服务信息处理器
func (oe *OperationalEndpoints) serviceHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: 从服务注册表获取单个服务信息
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": nil,
	})
}
