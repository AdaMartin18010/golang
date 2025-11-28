package system

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// KubernetesInfo Kubernetes 信息
type KubernetesInfo struct {
	PodName      string
	PodNamespace string
	NodeName     string
	PodIP        string
	HostIP       string
	ContainerID  string
	Labels       map[string]string
	Annotations  map[string]string
}

// KubernetesMonitor Kubernetes 监控器
type KubernetesMonitor struct {
	meter        metric.Meter
	info         KubernetesInfo
	attributes   []attribute.KeyValue
	enabled      bool
}

// NewKubernetesMonitor 创建 Kubernetes 监控器
func NewKubernetesMonitor(meter metric.Meter) (*KubernetesMonitor, error) {
	info, err := detectKubernetesInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to detect Kubernetes info: %w", err)
	}

	attrs := []attribute.KeyValue{
		attribute.String("k8s.pod.name", info.PodName),
		attribute.String("k8s.pod.namespace", info.PodNamespace),
		attribute.String("k8s.node.name", info.NodeName),
		attribute.String("k8s.pod.ip", info.PodIP),
		attribute.String("k8s.host.ip", info.HostIP),
	}

	// 添加标签
	for k, v := range info.Labels {
		attrs = append(attrs, attribute.String("k8s.label."+k, v))
	}

	// 添加注解（可选，通常较多，只添加关键注解）
	for k, v := range info.Annotations {
		if strings.HasPrefix(k, "prometheus.io/") || strings.HasPrefix(k, "opentelemetry.io/") {
			attrs = append(attrs, attribute.String("k8s.annotation."+k, v))
		}
	}

	return &KubernetesMonitor{
		meter:      meter,
		info:       info,
		attributes: attrs,
		enabled:    info.PodName != "",
	}, nil
}

// GetInfo 获取 Kubernetes 信息
func (km *KubernetesMonitor) GetInfo() KubernetesInfo {
	return km.info
}

// GetAttributes 获取属性
func (km *KubernetesMonitor) GetAttributes() []attribute.KeyValue {
	return km.attributes
}

// IsEnabled 检查是否启用
func (km *KubernetesMonitor) IsEnabled() bool {
	return km.enabled
}

// RecordKubernetesMetrics 记录 Kubernetes 指标
func (km *KubernetesMonitor) RecordKubernetesMetrics(ctx context.Context) error {
	if !km.enabled {
		return nil
	}

	// 创建 Kubernetes 信息指标
	k8sInfoGauge, err := km.meter.Int64ObservableGauge(
		"k8s.pod.info",
		metric.WithDescription("Kubernetes pod information (always 1, attributes contain the info)"),
	)
	if err != nil {
		return err
	}

	// 注册回调
	_, err = km.meter.RegisterCallback(
		func(ctx context.Context, obs metric.Observer) error {
			obs.ObserveInt64(k8sInfoGauge, 1, metric.WithAttributes(km.attributes...))
			return nil
		},
		k8sInfoGauge,
	)
	return err
}

// detectKubernetesInfo 检测 Kubernetes 信息
func detectKubernetesInfo() (KubernetesInfo, error) {
	info := KubernetesInfo{}

	// 从环境变量读取
	info.PodName = os.Getenv("POD_NAME")
	if info.PodName == "" {
		info.PodName = os.Getenv("HOSTNAME") // 通常 Pod 名就是 HOSTNAME
	}
	info.PodNamespace = os.Getenv("POD_NAMESPACE")
	info.NodeName = os.Getenv("NODE_NAME")
	info.PodIP = os.Getenv("POD_IP")
	info.HostIP = os.Getenv("HOST_IP")

	// 从文件读取（如果环境变量不存在）
	if info.PodName == "" {
		if data, err := os.ReadFile("/etc/hostname"); err == nil {
			info.PodName = strings.TrimSpace(string(data))
		}
	}

	// 读取 Pod 信息文件（如果存在）
	// Kubernetes 会将 Pod 信息挂载到 /etc/podinfo/
	podInfoDir := "/etc/podinfo"
	if _, err := os.Stat(podInfoDir); err == nil {
		// 读取 labels
		if data, err := os.ReadFile(filepath.Join(podInfoDir, "labels")); err == nil {
			info.Labels = parseKeyValueFile(string(data))
		}

		// 读取 annotations
		if data, err := os.ReadFile(filepath.Join(podInfoDir, "annotations")); err == nil {
			info.Annotations = parseKeyValueFile(string(data))
		}
	}

	// 从 Downward API 读取（如果配置了）
	// 通常挂载在 /etc/podinfo/
	if info.PodNamespace == "" {
		if data, err := os.ReadFile(filepath.Join(podInfoDir, "namespace")); err == nil {
			info.PodNamespace = strings.TrimSpace(string(data))
		}
	}

	// 检测容器 ID
	info.ContainerID, _ = detectContainer()

	return info, nil
}

// parseKeyValueFile 解析键值对文件
func parseKeyValueFile(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.Trim(parts[0], `"`)
			value := strings.Trim(parts[1], `"`)
			result[key] = value
		}
	}
	return result
}
