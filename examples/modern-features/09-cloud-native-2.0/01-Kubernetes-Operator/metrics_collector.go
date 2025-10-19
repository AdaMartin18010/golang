package kubernetes_operator

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsCollector 指标收集器
type MetricsCollector struct {
	// 应用相关指标
	applicationsTotal   prometheus.Counter
	applicationsRunning prometheus.Gauge
	applicationsFailed  prometheus.Gauge
	applicationsScaling prometheus.Gauge

	// 调和相关指标
	reconcileTotal       prometheus.Counter
	reconcileErrors      prometheus.Counter
	reconcileDuration    prometheus.Histogram
	reconcileQueueLength prometheus.Gauge

	// 资源相关指标
	deploymentsTotal prometheus.Counter
	servicesTotal    prometheus.Counter
	hpasTotal        prometheus.Counter
	podsTotal        prometheus.Gauge
	podsReady        prometheus.Gauge
	podsFailed       prometheus.Gauge

	// 性能相关指标
	resourceUsageCPU       *prometheus.GaugeVec
	resourceUsageMemory    *prometheus.GaugeVec
	resourceRequestsCPU    *prometheus.GaugeVec
	resourceRequestsMemory *prometheus.GaugeVec
	resourceLimitsCPU      *prometheus.GaugeVec
	resourceLimitsMemory   *prometheus.GaugeVec

	// 网络相关指标
	networkRequestsTotal   prometheus.Counter
	networkRequestsLatency prometheus.Histogram
	networkErrorsTotal     prometheus.Counter

	// 存储相关指标
	storageUsage    *prometheus.GaugeVec
	storageCapacity *prometheus.GaugeVec
	storageIOPS     *prometheus.CounterVec

	// 安全相关指标
	securityVulnerabilities *prometheus.GaugeVec
	complianceViolations    *prometheus.GaugeVec
	certificateExpiryDays   *prometheus.GaugeVec

	// 事件相关指标
	eventsTotal    prometheus.Counter
	eventsByType   *prometheus.CounterVec
	eventsByReason *prometheus.CounterVec

	// 内部状态
	mu                sync.RWMutex
	startTime         time.Time
	lastReconcileTime time.Time
	reconcileCount    int64
	errorCount        int64
}

// NewMetricsCollector 创建新的指标收集器
func NewMetricsCollector() *MetricsCollector {
	mc := &MetricsCollector{
		startTime: time.Now(),
	}

	// 初始化Prometheus指标
	mc.initializeMetrics()

	return mc
}

// initializeMetrics 初始化指标
func (mc *MetricsCollector) initializeMetrics() {
	// 应用相关指标
	mc.applicationsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "application_operator_applications_total",
		Help: "Total number of applications managed by the operator",
	})

	mc.applicationsRunning = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "application_operator_applications_running",
		Help: "Number of applications currently running",
	})

	mc.applicationsFailed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "application_operator_applications_failed",
		Help: "Number of applications in failed state",
	})

	mc.applicationsScaling = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "application_operator_applications_scaling",
		Help: "Number of applications currently scaling",
	})

	// 调和相关指标
	mc.reconcileTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "application_operator_reconcile_total",
		Help: "Total number of reconcile operations",
	})

	mc.reconcileErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "application_operator_reconcile_errors_total",
		Help: "Total number of reconcile errors",
	})

	mc.reconcileDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "application_operator_reconcile_duration_seconds",
		Help:    "Duration of reconcile operations",
		Buckets: prometheus.DefBuckets,
	})

	mc.reconcileQueueLength = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "application_operator_reconcile_queue_length",
		Help: "Current length of the reconcile queue",
	})

	// 资源相关指标
	mc.deploymentsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "application_operator_deployments_total",
		Help: "Total number of deployments created",
	})

	mc.servicesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "application_operator_services_total",
		Help: "Total number of services created",
	})

	mc.hpasTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "application_operator_hpas_total",
		Help: "Total number of HPAs created",
	})

	mc.podsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "application_operator_pods_total",
		Help: "Total number of pods managed",
	})

	mc.podsReady = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "application_operator_pods_ready",
		Help: "Number of pods in ready state",
	})

	mc.podsFailed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "application_operator_pods_failed",
		Help: "Number of pods in failed state",
	})

	// 性能相关指标
	mc.resourceUsageCPU = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_resource_usage_cpu_cores",
		Help: "CPU usage in cores",
	}, []string{"application", "namespace", "pod"})

	mc.resourceUsageMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_resource_usage_memory_bytes",
		Help: "Memory usage in bytes",
	}, []string{"application", "namespace", "pod"})

	mc.resourceRequestsCPU = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_resource_requests_cpu_cores",
		Help: "CPU requests in cores",
	}, []string{"application", "namespace"})

	mc.resourceRequestsMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_resource_requests_memory_bytes",
		Help: "Memory requests in bytes",
	}, []string{"application", "namespace"})

	mc.resourceLimitsCPU = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_resource_limits_cpu_cores",
		Help: "CPU limits in cores",
	}, []string{"application", "namespace"})

	mc.resourceLimitsMemory = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_resource_limits_memory_bytes",
		Help: "Memory limits in bytes",
	}, []string{"application", "namespace"})

	// 网络相关指标
	mc.networkRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "application_operator_network_requests_total",
		Help: "Total number of network requests",
	})

	mc.networkRequestsLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "application_operator_network_requests_latency_seconds",
		Help:    "Network request latency",
		Buckets: prometheus.DefBuckets,
	})

	mc.networkErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "application_operator_network_errors_total",
		Help: "Total number of network errors",
	})

	// 存储相关指标
	mc.storageUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_storage_usage_bytes",
		Help: "Storage usage in bytes",
	}, []string{"application", "namespace", "pvc"})

	mc.storageCapacity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_storage_capacity_bytes",
		Help: "Storage capacity in bytes",
	}, []string{"application", "namespace", "pvc"})

	mc.storageIOPS = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "application_operator_storage_iops_total",
		Help: "Storage I/O operations",
	}, []string{"application", "namespace", "pvc", "operation"})

	// 安全相关指标
	mc.securityVulnerabilities = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_security_vulnerabilities",
		Help: "Number of security vulnerabilities",
	}, []string{"application", "namespace", "severity"})

	mc.complianceViolations = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_compliance_violations",
		Help: "Number of compliance violations",
	}, []string{"application", "namespace", "standard"})

	mc.certificateExpiryDays = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "application_operator_certificate_expiry_days",
		Help: "Days until certificate expires",
	}, []string{"application", "namespace", "certificate"})

	// 事件相关指标
	mc.eventsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "application_operator_events_total",
		Help: "Total number of events recorded",
	})

	mc.eventsByType = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "application_operator_events_by_type_total",
		Help: "Events by type",
	}, []string{"application", "namespace", "type"})

	mc.eventsByReason = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "application_operator_events_by_reason_total",
		Help: "Events by reason",
	}, []string{"application", "namespace", "reason"})
}

// RecordReconcileStart 记录调和开始
func (mc *MetricsCollector) RecordReconcileStart(applicationName string) {
	atomic.AddInt64(&mc.reconcileCount, 1)
	mc.reconcileTotal.Inc()
	mc.mu.Lock()
	mc.lastReconcileTime = time.Now()
	mc.mu.Unlock()
}

// RecordReconcileSuccess 记录调和成功
func (mc *MetricsCollector) RecordReconcileSuccess(applicationName string) {
	duration := time.Since(mc.lastReconcileTime).Seconds()
	mc.reconcileDuration.Observe(duration)
}

// RecordReconcileError 记录调和错误
func (mc *MetricsCollector) RecordReconcileError(applicationName string) {
	atomic.AddInt64(&mc.errorCount, 1)
	mc.reconcileErrors.Inc()
}

// RecordApplicationCreated 记录应用创建
func (mc *MetricsCollector) RecordApplicationCreated(applicationName string) {
	mc.applicationsTotal.Inc()
	mc.applicationsRunning.Inc()
}

// RecordApplicationDeleted 记录应用删除
func (mc *MetricsCollector) RecordApplicationDeleted(applicationName string) {
	mc.applicationsRunning.Dec()
}

// RecordApplicationFailed 记录应用失败
func (mc *MetricsCollector) RecordApplicationFailed(applicationName string) {
	mc.applicationsRunning.Dec()
	mc.applicationsFailed.Inc()
}

// RecordApplicationRecovered 记录应用恢复
func (mc *MetricsCollector) RecordApplicationRecovered(applicationName string) {
	mc.applicationsFailed.Dec()
	mc.applicationsRunning.Inc()
}

// RecordApplicationScaling 记录应用扩缩容
func (mc *MetricsCollector) RecordApplicationScaling(applicationName string, scaling bool) {
	if scaling {
		mc.applicationsScaling.Inc()
	} else {
		mc.applicationsScaling.Dec()
	}
}

// RecordDeploymentCreated 记录Deployment创建
func (mc *MetricsCollector) RecordDeploymentCreated(applicationName, namespace string) {
	mc.deploymentsTotal.Inc()
}

// RecordServiceCreated 记录Service创建
func (mc *MetricsCollector) RecordServiceCreated(applicationName, namespace string) {
	mc.servicesTotal.Inc()
}

// RecordHPACreated 记录HPA创建
func (mc *MetricsCollector) RecordHPACreated(applicationName, namespace string) {
	mc.hpasTotal.Inc()
}

// RecordPodStatus 记录Pod状态
func (mc *MetricsCollector) RecordPodStatus(applicationName, namespace string, total, ready, failed int) {
	mc.podsTotal.Set(float64(total))
	mc.podsReady.Set(float64(ready))
	mc.podsFailed.Set(float64(failed))
}

// RecordResourceUsage 记录资源使用情况
func (mc *MetricsCollector) RecordResourceUsage(applicationName, namespace, podName string, cpuCores, memoryBytes float64) {
	mc.resourceUsageCPU.WithLabelValues(applicationName, namespace, podName).Set(cpuCores)
	mc.resourceUsageMemory.WithLabelValues(applicationName, namespace, podName).Set(memoryBytes)
}

// RecordResourceRequests 记录资源请求
func (mc *MetricsCollector) RecordResourceRequests(applicationName, namespace string, cpuCores, memoryBytes float64) {
	mc.resourceRequestsCPU.WithLabelValues(applicationName, namespace).Set(cpuCores)
	mc.resourceRequestsMemory.WithLabelValues(applicationName, namespace).Set(memoryBytes)
}

// RecordResourceLimits 记录资源限制
func (mc *MetricsCollector) RecordResourceLimits(applicationName, namespace string, cpuCores, memoryBytes float64) {
	mc.resourceLimitsCPU.WithLabelValues(applicationName, namespace).Set(cpuCores)
	mc.resourceLimitsMemory.WithLabelValues(applicationName, namespace).Set(memoryBytes)
}

// RecordNetworkRequest 记录网络请求
func (mc *MetricsCollector) RecordNetworkRequest(applicationName, namespace string, latency time.Duration, success bool) {
	mc.networkRequestsTotal.Inc()
	mc.networkRequestsLatency.Observe(latency.Seconds())

	if !success {
		mc.networkErrorsTotal.Inc()
	}
}

// RecordStorageUsage 记录存储使用情况
func (mc *MetricsCollector) RecordStorageUsage(applicationName, namespace, pvcName string, usageBytes, capacityBytes float64) {
	mc.storageUsage.WithLabelValues(applicationName, namespace, pvcName).Set(usageBytes)
	mc.storageCapacity.WithLabelValues(applicationName, namespace, pvcName).Set(capacityBytes)
}

// RecordStorageIOPS 记录存储I/O操作
func (mc *MetricsCollector) RecordStorageIOPS(applicationName, namespace, pvcName, operation string) {
	mc.storageIOPS.WithLabelValues(applicationName, namespace, pvcName, operation).Inc()
}

// RecordSecurityVulnerability 记录安全漏洞
func (mc *MetricsCollector) RecordSecurityVulnerability(applicationName, namespace, severity string, count int) {
	mc.securityVulnerabilities.WithLabelValues(applicationName, namespace, severity).Set(float64(count))
}

// RecordComplianceViolation 记录合规违规
func (mc *MetricsCollector) RecordComplianceViolation(applicationName, namespace, standard string, count int) {
	mc.complianceViolations.WithLabelValues(applicationName, namespace, standard).Set(float64(count))
}

// RecordCertificateExpiry 记录证书过期时间
func (mc *MetricsCollector) RecordCertificateExpiry(applicationName, namespace, certificateName string, daysUntilExpiry int) {
	mc.certificateExpiryDays.WithLabelValues(applicationName, namespace, certificateName).Set(float64(daysUntilExpiry))
}

// RecordEvent 记录事件
func (mc *MetricsCollector) RecordEvent(applicationName, namespace, eventType, reason string) {
	mc.eventsTotal.Inc()
	mc.eventsByType.WithLabelValues(applicationName, namespace, eventType).Inc()
	mc.eventsByReason.WithLabelValues(applicationName, namespace, reason).Inc()
}

// SetReconcileQueueLength 设置调和队列长度
func (mc *MetricsCollector) SetReconcileQueueLength(length int) {
	mc.reconcileQueueLength.Set(float64(length))
}

// GetStats 获取统计信息
func (mc *MetricsCollector) GetStats() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	return map[string]interface{}{
		"uptime_seconds":      time.Since(mc.startTime).Seconds(),
		"reconcile_count":     atomic.LoadInt64(&mc.reconcileCount),
		"error_count":         atomic.LoadInt64(&mc.errorCount),
		"last_reconcile_time": mc.lastReconcileTime,
		"start_time":          mc.startTime,
	}
}

// Reset 重置指标
func (mc *MetricsCollector) Reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// 重置计数器
	atomic.StoreInt64(&mc.reconcileCount, 0)
	atomic.StoreInt64(&mc.errorCount, 0)

	// 重置时间
	mc.startTime = time.Now()
	mc.lastReconcileTime = time.Time{}
}
