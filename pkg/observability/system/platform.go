package system

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// PlatformInfo 平台信息
type PlatformInfo struct {
	OS              string
	Arch            string
	GoVersion       string
	Hostname        string
	ContainerID     string
	ContainerName   string
	KubernetesPod   string
	KubernetesNode  string
	Virtualization  string
	CPUs            int
}

// PlatformMonitor 平台信息监控器
type PlatformMonitor struct {
	meter      metric.Meter
	info       PlatformInfo
	attributes []attribute.KeyValue
}

// NewPlatformMonitor 创建平台信息监控器
func NewPlatformMonitor(meter metric.Meter) (*PlatformMonitor, error) {
	info, err := detectPlatform()
	if err != nil {
		return nil, fmt.Errorf("failed to detect platform: %w", err)
	}

	attrs := []attribute.KeyValue{
		attribute.String("system.os", info.OS),
		attribute.String("system.arch", info.Arch),
		attribute.String("system.go_version", info.GoVersion),
		attribute.String("system.hostname", info.Hostname),
		attribute.Int("system.cpus", info.CPUs),
	}

	if info.ContainerID != "" {
		attrs = append(attrs, attribute.String("container.id", info.ContainerID))
	}
	if info.ContainerName != "" {
		attrs = append(attrs, attribute.String("container.name", info.ContainerName))
	}
	if info.KubernetesPod != "" {
		attrs = append(attrs, attribute.String("k8s.pod", info.KubernetesPod))
	}
	if info.KubernetesNode != "" {
		attrs = append(attrs, attribute.String("k8s.node", info.KubernetesNode))
	}
	if info.Virtualization != "" {
		attrs = append(attrs, attribute.String("system.virtualization", info.Virtualization))
	}

	return &PlatformMonitor{
		meter:      meter,
		info:       info,
		attributes: attrs,
	}, nil
}

// GetInfo 获取平台信息
func (pm *PlatformMonitor) GetInfo() PlatformInfo {
	return pm.info
}

// GetAttributes 获取属性（用于添加到指标和追踪）
func (pm *PlatformMonitor) GetAttributes() []attribute.KeyValue {
	return pm.attributes
}

// RecordPlatformMetrics 记录平台指标
func (pm *PlatformMonitor) RecordPlatformMetrics(ctx context.Context) error {
	// 创建平台信息指标
	platformInfoGauge, err := pm.meter.Int64ObservableGauge(
		"system.platform.info",
		metric.WithDescription("Platform information (always 1, attributes contain the info)"),
	)
	if err != nil {
		return err
	}

	// 注册回调
	_, err = pm.meter.RegisterCallback(
		func(ctx context.Context, obs metric.Observer) error {
			obs.ObserveInt64(platformInfoGauge, 1, metric.WithAttributes(pm.attributes...))
			return nil
		},
		platformInfoGauge,
	)
	return err
}

// detectPlatform 检测平台信息
func detectPlatform() (PlatformInfo, error) {
	info := PlatformInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
		CPUs:      runtime.NumCPU(),
	}

	// 获取主机名
	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}

	// 检测容器环境
	containerID, containerName := detectContainer()
	info.ContainerID = containerID
	info.ContainerName = containerName

	// 检测 Kubernetes
	pod, node := detectKubernetes()
	info.KubernetesPod = pod
	info.KubernetesNode = node

	// 检测虚拟化
	info.Virtualization = detectVirtualization()

	return info, nil
}

// detectContainer 检测容器环境
func detectContainer() (id, name string) {
	// 检测 Docker
	// 检查 /.dockerenv 文件
	if _, err := os.Stat("/.dockerenv"); err == nil {
		// 尝试读取容器 ID
		if data, err := os.ReadFile("/proc/self/cgroup"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.Contains(line, "docker") {
					parts := strings.Split(line, "/")
					if len(parts) > 0 {
						id = parts[len(parts)-1]
						if len(id) > 12 {
							id = id[:12] // Docker 容器 ID 前 12 位
						}
					}
					break
				}
			}
		}

		// 尝试从环境变量获取容器名
		name = os.Getenv("HOSTNAME")
		return id, name
	}

	// 检测 Kubernetes Pod
	if podName := os.Getenv("HOSTNAME"); podName != "" {
		// 可能是 Kubernetes Pod
		if _, err := os.Stat("/var/run/secrets/kubernetes.io"); err == nil {
			return podName, podName
		}
	}

	return "", ""
}

// detectKubernetes 检测 Kubernetes 环境
func detectKubernetes() (pod, node string) {
	// 从环境变量读取
	pod = os.Getenv("POD_NAME")
	if pod == "" {
		pod = os.Getenv("HOSTNAME") // 通常 Pod 名就是 HOSTNAME
	}

	node = os.Getenv("NODE_NAME")

	// 尝试从文件读取（如果存在）
	if pod == "" {
		if data, err := os.ReadFile("/etc/hostname"); err == nil {
			pod = strings.TrimSpace(string(data))
		}
	}

	// 检查是否在 Kubernetes 中
	if pod != "" {
		// 检查是否有 Kubernetes 相关的挂载点或文件
		// 1. 检查 service account token
		if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount"); err == nil {
			return pod, node
		}
		// 2. 检查 podinfo 目录
		if _, err := os.Stat("/etc/podinfo"); err == nil {
			return pod, node
		}
		// 3. 检查 KUBERNETES_SERVICE_HOST 环境变量
		if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
			return pod, node
		}
	}

	return "", ""
}

// detectVirtualization 检测虚拟化环境
func detectVirtualization() string {
	// 检测 Docker
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "docker"
	}

	// 检测 systemd-nspawn
	if _, err := os.Stat("/run/systemd/container"); err == nil {
		return "systemd-nspawn"
	}

	// 检测 LXC
	if _, err := os.Stat("/.lxc"); err == nil {
		return "lxc"
	}

	// 检测虚拟化（通过 /proc/cpuinfo 或 /sys/class/dmi/id）
	if data, err := os.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
		product := strings.ToLower(strings.TrimSpace(string(data)))
		if strings.Contains(product, "vmware") {
			return "vmware"
		}
		if strings.Contains(product, "virtualbox") {
			return "virtualbox"
		}
		if strings.Contains(product, "kvm") || strings.Contains(product, "qemu") {
			return "kvm"
		}
		if strings.Contains(product, "xen") {
			return "xen"
		}
	}

	// 检测云环境
	if vendor := os.Getenv("CLOUD_PROVIDER"); vendor != "" {
		return vendor
	}

	// 检测 AWS
	if os.Getenv("AWS_EXECUTION_ENV") != "" {
		return "aws"
	}

	// 检测 GCP
	if project := os.Getenv("GOOGLE_CLOUD_PROJECT"); project != "" {
		return "gcp"
	}

	// 检测 Azure
	if os.Getenv("WEBSITE_SITE_NAME") != "" {
		return "azure"
	}

	return "bare-metal"
}

// IsContainer 检查是否在容器中
func (pm *PlatformMonitor) IsContainer() bool {
	return pm.info.ContainerID != "" || pm.info.ContainerName != ""
}

// IsKubernetes 检查是否在 Kubernetes 中
func (pm *PlatformMonitor) IsKubernetes() bool {
	return pm.info.KubernetesPod != ""
}

// IsVirtualized 检查是否在虚拟化环境中
func (pm *PlatformMonitor) IsVirtualized() bool {
	return pm.info.Virtualization != "" && pm.info.Virtualization != "bare-metal"
}

// GetContainerInfo 获取容器信息
func (pm *PlatformMonitor) GetContainerInfo() (id, name string) {
	return pm.info.ContainerID, pm.info.ContainerName
}

// GetKubernetesInfo 获取 Kubernetes 信息
func (pm *PlatformMonitor) GetKubernetesInfo() (pod, node string) {
	return pm.info.KubernetesPod, pm.info.KubernetesNode
}

// GetVirtualization 获取虚拟化类型
func (pm *PlatformMonitor) GetVirtualization() string {
	return pm.info.Virtualization
}
