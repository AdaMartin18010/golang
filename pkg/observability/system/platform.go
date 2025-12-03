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
	KubernetesNS    string // Kubernetes Namespace
	Virtualization  string
	CloudProvider   string // AWS/GCP/Azure/Alibaba/Tencent
	CloudRegion     string
	CloudZone       string
	CPUs            int
	MemoryTotal     uint64 // Total memory in bytes
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
	pod, node, ns := detectKubernetes()
	info.KubernetesPod = pod
	info.KubernetesNode = node
	info.KubernetesNS = ns

	// 检测云环境
	provider, region, zone := detectCloudProvider()
	info.CloudProvider = provider
	info.CloudRegion = region
	info.CloudZone = zone

	// 检测虚拟化
	info.Virtualization = detectVirtualization()

	// 检测内存
	info.MemoryTotal = detectTotalMemory()

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
func detectKubernetes() (pod, node, namespace string) {
	// 从环境变量读取（推荐的 Downward API 方式）
	pod = os.Getenv("POD_NAME")
	if pod == "" {
		pod = os.Getenv("HOSTNAME") // 通常 Pod 名就是 HOSTNAME
	}

	node = os.Getenv("NODE_NAME")
	namespace = os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = os.Getenv("NAMESPACE")
	}

	// 尝试从文件读取（如果存在）
	if pod == "" {
		if data, err := os.ReadFile("/etc/hostname"); err == nil {
			pod = strings.TrimSpace(string(data))
		}
	}

	// 从 service account 读取 namespace
	if namespace == "" {
		if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
			namespace = strings.TrimSpace(string(data))
		}
	}

	// 检查是否在 Kubernetes 中
	if pod != "" {
		// 检查是否有 Kubernetes 相关的挂载点或文件
		// 1. 检查 service account token
		if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount"); err == nil {
			return pod, node, namespace
		}
		// 2. 检查 podinfo 目录
		if _, err := os.Stat("/etc/podinfo"); err == nil {
			return pod, node, namespace
		}
		// 3. 检查 KUBERNETES_SERVICE_HOST 环境变量
		if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
			return pod, node, namespace
		}
	}

	return "", "", ""
}

// detectCloudProvider 检测云服务提供商
func detectCloudProvider() (provider, region, zone string) {
	// AWS 检测
	if os.Getenv("AWS_EXECUTION_ENV") != "" || os.Getenv("AWS_REGION") != "" {
		provider = "aws"
		region = os.Getenv("AWS_REGION")
		zone = os.Getenv("AWS_AVAILABILITY_ZONE")
		
		// 尝试从 EC2 metadata 读取
		if region == "" {
			// 实际实现需要调用 EC2 metadata API
			// http://169.254.169.254/latest/meta-data/placement/region
		}
		return provider, region, zone
	}

	// GCP 检测
	if os.Getenv("GOOGLE_CLOUD_PROJECT") != "" || os.Getenv("GCP_PROJECT") != "" {
		provider = "gcp"
		region = os.Getenv("GOOGLE_CLOUD_REGION")
		zone = os.Getenv("GOOGLE_CLOUD_ZONE")
		
		// 尝试从 GCE metadata 读取
		if region == "" {
			// 实际实现需要调用 GCE metadata API
			// http://metadata.google.internal/computeMetadata/v1/instance/zone
		}
		return provider, region, zone
	}

	// Azure 检测
	if os.Getenv("WEBSITE_SITE_NAME") != "" || os.Getenv("AZURE_CLIENT_ID") != "" {
		provider = "azure"
		region = os.Getenv("AZURE_REGION")
		// Azure 使用 region 而不是 zone
		return provider, region, ""
	}

	// Alibaba Cloud 检测
	if os.Getenv("ALIBABA_CLOUD_REGION_ID") != "" {
		provider = "alibaba"
		region = os.Getenv("ALIBABA_CLOUD_REGION_ID")
		zone = os.Getenv("ALIBABA_CLOUD_ZONE_ID")
		return provider, region, zone
	}

	// Tencent Cloud 检测
	if os.Getenv("TENCENTCLOUD_REGION") != "" {
		provider = "tencent"
		region = os.Getenv("TENCENTCLOUD_REGION")
		zone = os.Getenv("TENCENTCLOUD_ZONE")
		return provider, region, zone
	}

	// 通用云环境变量
	if provider := os.Getenv("CLOUD_PROVIDER"); provider != "" {
		return provider, os.Getenv("CLOUD_REGION"), os.Getenv("CLOUD_ZONE")
	}

	return "", "", ""
}

// detectTotalMemory 检测总内存
func detectTotalMemory() uint64 {
	// Linux: 读取 /proc/meminfo
	if runtime.GOOS == "linux" {
		if data, err := os.ReadFile("/proc/meminfo"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						var kb uint64
						fmt.Sscanf(fields[1], "%d", &kb)
						return kb * 1024 // KB to Bytes
					}
				}
			}
		}
	}
	// 其他操作系统可以使用 gopsutil 库
	return 0
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
