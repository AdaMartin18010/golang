package kubernetes_operator

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceManager Kubernetes资源管理器
type ResourceManager struct {
	client client.Client
}

// NewResourceManager 创建资源管理器
func NewResourceManager(client client.Client) *ResourceManager {
	return &ResourceManager{
		client: client,
	}
}

// createPersistentVolumeClaims 创建持久化卷声明
func (rm *ResourceManager) createPersistentVolumeClaims(ctx context.Context, app *Application) error {
	for _, pvcSpec := range app.Spec.Storage.PersistentVolumeClaims {
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pvcSpec.Name,
				Namespace: app.Namespace,
				Labels:    getLabels(app),
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(app, app.GroupVersionKind()),
				},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes:      pvcSpec.AccessModes,
				StorageClassName: pvcSpec.StorageClassName,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: pvcSpec.Size,
					},
				},
			},
		}

		if err := rm.client.Create(ctx, pvc); err != nil && !errors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create PVC %s: %w", pvcSpec.Name, err)
		}
	}
	return nil
}

// createOrUpdateDeployment 创建或更新Deployment
func (rm *ResourceManager) createOrUpdateDeployment(ctx context.Context, app *Application) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels:    getLabels(app),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, app.GroupVersionKind()),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(app.Spec.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: getLabels(app),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getLabels(app),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  app.Name,
							Image: app.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: app.Spec.Port,
								},
							},
							Env:             app.Spec.Environment,
							Resources:       app.Spec.Resources,
							LivenessProbe:   app.Spec.HealthCheck.LivenessProbe,
							ReadinessProbe:  app.Spec.HealthCheck.ReadinessProbe,
							StartupProbe:    app.Spec.HealthCheck.StartupProbe,
							SecurityContext: app.Spec.Security.SecurityContext,
						},
					},
					ServiceAccountName: app.Spec.Security.ServiceAccount,
					ImagePullSecrets:  app.Spec.Security.ImagePullSecrets,
					SecurityContext:   app.Spec.Security.PodSecurityContext,
				},
			},
		},
	}

	// 添加存储卷
	if len(app.Spec.Storage.PersistentVolumeClaims) > 0 {
		var volumes []corev1.Volume
		var volumeMounts []corev1.VolumeMount

		for i, pvcSpec := range app.Spec.Storage.PersistentVolumeClaims {
			volumeName := fmt.Sprintf("storage-%d", i)
			volumes = append(volumes, corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcSpec.Name,
					},
				},
			})
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      volumeName,
				MountPath: fmt.Sprintf("/data/%d", i),
			})
		}

		deployment.Spec.Template.Spec.Volumes = volumes
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMounts
	}

	existing := &appsv1.Deployment{}
	err := rm.client.Get(ctx, types.NamespacedName{
		Namespace: app.Namespace,
		Name:      app.Name,
	}, existing)

	if err != nil {
		if errors.IsNotFound(err) {
			return rm.client.Create(ctx, deployment)
		}
		return fmt.Errorf("failed to get existing deployment: %w", err)
	}

	// 更新现有Deployment
	existing.Spec = deployment.Spec
	return rm.client.Update(ctx, existing)
}

// createOrUpdateService 创建或更新Service
func (rm *ResourceManager) createOrUpdateService(ctx context.Context, app *Application) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels:    getLabels(app),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, app.GroupVersionKind()),
			},
		},
		Spec: corev1.ServiceSpec{
			Type: app.Spec.Network.ServiceType,
			Selector: map[string]string{
				"app": app.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       app.Spec.Port,
					TargetPort: intstr.FromInt(int(app.Spec.Port)),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	// 配置负载均衡器
	if app.Spec.Network.ServiceType == corev1.ServiceTypeLoadBalancer && app.Spec.Network.LoadBalancer != nil {
		service.Spec.LoadBalancerSourceRanges = app.Spec.Network.LoadBalancer.SourceRanges
		service.Spec.ExternalIPs = app.Spec.Network.LoadBalancer.ExternalIPs
	}

	existing := &corev1.Service{}
	err := rm.client.Get(ctx, types.NamespacedName{
		Namespace: app.Namespace,
		Name:      app.Name,
	}, existing)

	if err != nil {
		if errors.IsNotFound(err) {
			return rm.client.Create(ctx, service)
		}
		return fmt.Errorf("failed to get existing service: %w", err)
	}

	// 更新现有Service
	existing.Spec = service.Spec
	return rm.client.Update(ctx, existing)
}

// createOrUpdateHPA 创建或更新HPA
func (rm *ResourceManager) createOrUpdateHPA(ctx context.Context, app *Application) error {
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels:    getLabels(app),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, app.GroupVersionKind()),
			},
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       app.Name,
			},
			MinReplicas: pointer.Int32(app.Spec.Scaling.MinReplicas),
			MaxReplicas: app.Spec.Scaling.MaxReplicas,
			Metrics:     app.Spec.Scaling.CustomMetrics,
		},
	}

	// 添加CPU和内存指标
	if app.Spec.Scaling.TargetCPUUtilizationPercentage > 0 {
		hpa.Spec.Metrics = append(hpa.Spec.Metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: corev1.ResourceCPU,
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: pointer.Int32(app.Spec.Scaling.TargetCPUUtilizationPercentage),
				},
			},
		})
	}

	if app.Spec.Scaling.TargetMemoryUtilizationPercentage > 0 {
		hpa.Spec.Metrics = append(hpa.Spec.Metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: corev1.ResourceMemory,
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: pointer.Int32(app.Spec.Scaling.TargetMemoryUtilizationPercentage),
				},
			},
		})
	}

	existing := &autoscalingv2.HorizontalPodAutoscaler{}
	err := rm.client.Get(ctx, types.NamespacedName{
		Namespace: app.Namespace,
		Name:      app.Name,
	}, existing)

	if err != nil {
		if errors.IsNotFound(err) {
			return rm.client.Create(ctx, hpa)
		}
		return fmt.Errorf("failed to get existing HPA: %w", err)
	}

	// 更新现有HPA
	existing.Spec = hpa.Spec
	return rm.client.Update(ctx, existing)
}

// cleanupResources 清理资源
func (rm *ResourceManager) cleanupResources(ctx context.Context, namespacedName types.NamespacedName) error {
	// 删除Deployment
	deployment := &appsv1.Deployment{}
	if err := rm.client.Get(ctx, namespacedName, deployment); err == nil {
		if err := rm.client.Delete(ctx, deployment); err != nil {
			return fmt.Errorf("failed to delete deployment: %w", err)
		}
	}

	// 删除Service
	service := &corev1.Service{}
	if err := rm.client.Get(ctx, namespacedName, service); err == nil {
		if err := rm.client.Delete(ctx, service); err != nil {
			return fmt.Errorf("failed to delete service: %w", err)
		}
	}

	// 删除HPA
	hpa := &autoscalingv2.HorizontalPodAutoscaler{}
	if err := rm.client.Get(ctx, namespacedName, hpa); err == nil {
		if err := rm.client.Delete(ctx, hpa); err != nil {
			return fmt.Errorf("failed to delete HPA: %w", err)
		}
	}

	// 删除持久化卷声明
	pvcList := &corev1.PersistentVolumeClaimList{}
	if err := rm.client.List(ctx, pvcList, client.InNamespace(namespacedName.Namespace)); err == nil {
		for _, pvc := range pvcList.Items {
			if pvc.Labels["app"] == namespacedName.Name {
				if err := rm.client.Delete(ctx, &pvc); err != nil {
					return fmt.Errorf("failed to delete PVC %s: %w", pvc.Name, err)
				}
			}
		}
	}

	return nil
}

// getLabels 获取标签
func getLabels(app *Application) map[string]string {
	return map[string]string{
		"app": app.Name,
		"managed-by": "application-controller",
	}
}

// 为ApplicationController添加资源管理方法
func (ac *ApplicationController) createPersistentVolumeClaims(ctx context.Context, app *Application) error {
	rm := NewResourceManager(ac.client)
	return rm.createPersistentVolumeClaims(ctx, app)
}

func (ac *ApplicationController) createOrUpdateDeployment(ctx context.Context, app *Application) error {
	rm := NewResourceManager(ac.client)
	return rm.createOrUpdateDeployment(ctx, app)
}

func (ac *ApplicationController) createOrUpdateService(ctx context.Context, app *Application) error {
	rm := NewResourceManager(ac.client)
	return rm.createOrUpdateService(ctx, app)
}

func (ac *ApplicationController) createOrUpdateHPA(ctx context.Context, app *Application) error {
	rm := NewResourceManager(ac.client)
	return rm.createOrUpdateHPA(ctx, app)
}

func (ac *ApplicationController) cleanupResources(ctx context.Context, namespacedName types.NamespacedName) error {
	rm := NewResourceManager(ac.client)
	return rm.cleanupResources(ctx, namespacedName)
}
