package kubernetes_operator

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Application 自定义资源定义
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// ApplicationSpec 应用规格
type ApplicationSpec struct {
	// 副本数
	Replicas int32 `json:"replicas"`
	// 容器镜像
	Image string `json:"image"`
	// 端口
	Port int32 `json:"port"`
	// 环境变量
	Environment []corev1.EnvVar `json:"environment,omitempty"`
	// 资源限制
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// 健康检查
	HealthCheck HealthCheckSpec `json:"healthCheck,omitempty"`
	// 自动扩缩容配置
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// 存储配置
	Storage StorageSpec `json:"storage,omitempty"`
	// 网络配置
	Network NetworkSpec `json:"network,omitempty"`
	// 安全配置
	Security SecuritySpec `json:"security,omitempty"`
}

// HealthCheckSpec 健康检查规格
type HealthCheckSpec struct {
	// 存活探针
	LivenessProbe *corev1.Probe `json:"livenessProbe,omitempty"`
	// 就绪探针
	ReadinessProbe *corev1.Probe `json:"readinessProbe,omitempty"`
	// 启动探针
	StartupProbe *corev1.Probe `json:"startupProbe,omitempty"`
}

// ScalingSpec 扩缩容规格
type ScalingSpec struct {
	// 最小副本数
	MinReplicas int32 `json:"minReplicas,omitempty"`
	// 最大副本数
	MaxReplicas int32 `json:"maxReplicas,omitempty"`
	// CPU目标利用率
	TargetCPUUtilizationPercentage int32 `json:"targetCPUUtilizationPercentage,omitempty"`
	// 内存目标利用率
	TargetMemoryUtilizationPercentage int32 `json:"targetMemoryUtilizationPercentage,omitempty"`
	// 自定义指标
	CustomMetrics []autoscalingv2.MetricSpec `json:"customMetrics,omitempty"`
}

// StorageSpec 存储规格
type StorageSpec struct {
	// 持久化卷声明
	PersistentVolumeClaims []PersistentVolumeClaimSpec `json:"persistentVolumeClaims,omitempty"`
	// 临时存储
	EphemeralStorage *corev1.ResourceRequirements `json:"ephemeralStorage,omitempty"`
}

// PersistentVolumeClaimSpec 持久化卷声明规格
type PersistentVolumeClaimSpec struct {
	Name string `json:"name"`
	// 存储大小
	Size resource.Quantity `json:"size"`
	// 存储类
	StorageClassName *string `json:"storageClassName,omitempty"`
	// 访问模式
	AccessModes []corev1.PersistentVolumeAccessMode `json:"accessModes"`
}

// NetworkSpec 网络规格
type NetworkSpec struct {
	// 服务类型
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
	// 负载均衡器配置
	LoadBalancer *LoadBalancerSpec `json:"loadBalancer,omitempty"`
}

// LoadBalancerSpec 负载均衡器规格
type LoadBalancerSpec struct {
	// 源IP保持
	SourceRanges []string `json:"sourceRanges,omitempty"`
	// 外部IP
	ExternalIPs []string `json:"externalIPs,omitempty"`
}

// SecuritySpec 安全规格
type SecuritySpec struct {
	// 服务账户
	ServiceAccount string `json:"serviceAccount,omitempty"`
	// 安全上下文
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
	// Pod安全上下文
	PodSecurityContext *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`
	// 镜像拉取密钥
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// ApplicationStatus 应用状态
type ApplicationStatus struct {
	// 阶段
	Phase string `json:"phase,omitempty"`
	// 副本数
	Replicas int32 `json:"replicas"`
	// 就绪副本数
	ReadyReplicas int32 `json:"readyReplicas"`
	// 可用副本数
	AvailableReplicas int32 `json:"availableReplicas"`
	// 条件
	Conditions []ApplicationCondition `json:"conditions,omitempty"`
	// 服务端点
	ServiceEndpoints []string `json:"serviceEndpoints,omitempty"`
	// 最后更新时间
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}

// ApplicationCondition 应用条件
type ApplicationCondition struct {
	// 类型
	Type string `json:"type"`
	// 状态
	Status corev1.ConditionStatus `json:"status"`
	// 最后转换时间
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// 原因
	Reason string `json:"reason,omitempty"`
	// 消息
	Message string `json:"message,omitempty"`
}

// DeepCopyObject 实现 runtime.Object 接口
func (a *Application) DeepCopyObject() runtime.Object {
	if a == nil {
		return nil
	}
	out := new(Application)
	a.DeepCopyInto(out)
	return out
}

// DeepCopyInto 深度复制到目标对象
func (a *Application) DeepCopyInto(out *Application) {
	*out = *a
	out.TypeMeta = a.TypeMeta
	a.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	a.Spec.DeepCopyInto(&out.Spec)
	a.Status.DeepCopyInto(&out.Status)
}

// DeepCopy 创建深度副本
func (a *Application) DeepCopy() *Application {
	if a == nil {
		return nil
	}
	out := new(Application)
	a.DeepCopyInto(out)
	return out
}

// DeepCopyInto 深度复制应用规格
func (as *ApplicationSpec) DeepCopyInto(out *ApplicationSpec) {
	*out = *as
	if as.Environment != nil {
		out.Environment = make([]corev1.EnvVar, len(as.Environment))
		for i := range as.Environment {
			as.Environment[i].DeepCopyInto(&out.Environment[i])
		}
	}
	as.Resources.DeepCopyInto(&out.Resources)
	if as.HealthCheck.LivenessProbe != nil {
		out.HealthCheck.LivenessProbe = as.HealthCheck.LivenessProbe.DeepCopy()
	}
	if as.HealthCheck.ReadinessProbe != nil {
		out.HealthCheck.ReadinessProbe = as.HealthCheck.ReadinessProbe.DeepCopy()
	}
	if as.HealthCheck.StartupProbe != nil {
		out.HealthCheck.StartupProbe = as.HealthCheck.StartupProbe.DeepCopy()
	}
	if as.Scaling.CustomMetrics != nil {
		out.Scaling.CustomMetrics = make([]autoscalingv2.MetricSpec, len(as.Scaling.CustomMetrics))
		for i := range as.Scaling.CustomMetrics {
			as.Scaling.CustomMetrics[i].DeepCopyInto(&out.Scaling.CustomMetrics[i])
		}
	}
	if as.Storage.PersistentVolumeClaims != nil {
		out.Storage.PersistentVolumeClaims = make([]PersistentVolumeClaimSpec, len(as.Storage.PersistentVolumeClaims))
		copy(out.Storage.PersistentVolumeClaims, as.Storage.PersistentVolumeClaims)
	}
	if as.Storage.EphemeralStorage != nil {
		out.Storage.EphemeralStorage = new(corev1.ResourceRequirements)
		as.Storage.EphemeralStorage.DeepCopyInto(out.Storage.EphemeralStorage)
	}
	if as.Security.SecurityContext != nil {
		out.Security.SecurityContext = as.Security.SecurityContext.DeepCopy()
	}
	if as.Security.PodSecurityContext != nil {
		out.Security.PodSecurityContext = as.Security.PodSecurityContext.DeepCopy()
	}
	if as.Security.ImagePullSecrets != nil {
		out.Security.ImagePullSecrets = make([]corev1.LocalObjectReference, len(as.Security.ImagePullSecrets))
		copy(out.Security.ImagePullSecrets, as.Security.ImagePullSecrets)
	}
}

// DeepCopyInto 深度复制应用状态
func (as *ApplicationStatus) DeepCopyInto(out *ApplicationStatus) {
	*out = *as
	if as.Conditions != nil {
		out.Conditions = make([]ApplicationCondition, len(as.Conditions))
		copy(out.Conditions, as.Conditions)
	}
	if as.ServiceEndpoints != nil {
		out.ServiceEndpoints = make([]string, len(as.ServiceEndpoints))
		copy(out.ServiceEndpoints, as.ServiceEndpoints)
	}
}

// ApplicationController 应用控制器
type ApplicationController struct {
	client   client.Client
	scheme   *runtime.Scheme
	queue    workqueue.RateLimitingInterface
	recorder *EventRecorder
	metrics  *MetricsCollector
}

// NewApplicationController 创建应用控制器
func NewApplicationController(mgr manager.Manager) (*ApplicationController, error) {
	controller := &ApplicationController{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		queue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "application"),
		recorder: NewEventRecorder(mgr.GetEventRecorderFor("application-controller")),
		metrics:  NewMetricsCollector(),
	}

	// 设置informer
	informer, err := mgr.GetCache().GetInformer(context.Background(), &Application{})
	if err != nil {
		return nil, fmt.Errorf("failed to get informer: %w", err)
	}

	// 添加事件处理器
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.enqueueApplication,
		UpdateFunc: controller.updateApplication,
		DeleteFunc: controller.deleteApplication,
	})

	return controller, nil
}

// enqueueApplication 将应用加入队列
func (ac *ApplicationController) enqueueApplication(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	ac.queue.Add(key)
}

// updateApplication 更新应用
func (ac *ApplicationController) updateApplication(old, new interface{}) {
	oldApp := old.(*Application)
	newApp := new.(*Application)

	// 检查是否有重要变更
	if oldApp.Spec.Replicas != newApp.Spec.Replicas ||
		oldApp.Spec.Image != newApp.Spec.Image ||
		oldApp.Spec.Port != newApp.Spec.Port {
		ac.enqueueApplication(new)
	}
}

// deleteApplication 删除应用
func (ac *ApplicationController) deleteApplication(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	ac.queue.Add(key)
}

// Reconcile 调和逻辑
func (ac *ApplicationController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	ac.metrics.RecordReconcileStart(req.Name)

	// 获取应用实例
	app := &Application{}
	if err := ac.client.Get(ctx, req.NamespacedName, app); err != nil {
		if errors.IsNotFound(err) {
			// 应用已删除，清理相关资源
			return reconcile.Result{}, ac.cleanupResources(ctx, req.NamespacedName)
		}
		ac.metrics.RecordReconcileError(req.Name)
		ac.recorder.RecordReconcileError(app, err.Error())
		return reconcile.Result{}, err
	}

	// 执行调和逻辑
	result, err := ac.reconcileApplication(ctx, app)
	if err != nil {
		ac.metrics.RecordReconcileError(req.Name)
		ac.recorder.RecordReconcileError(app, err.Error())
	} else {
		ac.metrics.RecordReconcileSuccess(req.Name)
		ac.recorder.RecordReconcileSuccess(app)
	}

	return result, err
}

// reconcileApplication 调和应用
func (ac *ApplicationController) reconcileApplication(ctx context.Context, app *Application) (reconcile.Result, error) {
	// 检查应用是否被标记为删除
	if app.ObjectMeta.DeletionTimestamp != nil {
		return ac.handleDeletion(ctx, app)
	}

	// 添加finalizer
	if err := ac.addFinalizer(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 根据应用状态执行不同的调和逻辑
	switch app.Status.Phase {
	case "":
		// 新创建的应用
		return ac.handleCreation(ctx, app)
	case "Creating":
		return ac.handleCreation(ctx, app)
	case "Running":
		return ac.handleRunning(ctx, app)
	case "Scaling":
		return ac.handleScaling(ctx, app)
	case "Updating":
		return ac.handleUpdating(ctx, app)
	case "Failed":
		return ac.handleFailure(ctx, app)
	default:
		return reconcile.Result{}, fmt.Errorf("unknown application phase: %s", app.Status.Phase)
	}
}

// handleCreation 处理应用创建
func (ac *ApplicationController) handleCreation(ctx context.Context, app *Application) (reconcile.Result, error) {
	ac.recorder.RecordApplicationCreated(app)

	// 更新状态为创建中
	if app.Status.Phase != "Creating" {
		app.Status.Phase = "Creating"
		app.Status.LastUpdateTime = metav1.Now()
		if err := ac.client.Status().Update(ctx, app); err != nil {
			return reconcile.Result{}, err
		}
	}

	// 创建持久化卷声明
	if err := ac.createPersistentVolumeClaims(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 创建Deployment
	if err := ac.createOrUpdateDeployment(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 创建Service
	if err := ac.createOrUpdateService(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 创建HPA（如果配置了自动扩缩容）
	if app.Spec.Scaling.MaxReplicas > 0 {
		if err := ac.createOrUpdateHPA(ctx, app); err != nil {
			return reconcile.Result{}, err
		}
	}

	// 更新状态为运行中
	app.Status.Phase = "Running"
	app.Status.LastUpdateTime = metav1.Now()
	if err := ac.client.Status().Update(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	ac.recorder.RecordApplicationRunning(app)
	return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

// handleRunning 处理运行中的应用
func (ac *ApplicationController) handleRunning(ctx context.Context, app *Application) (reconcile.Result, error) {
	// 检查Deployment状态
	deployment := &appsv1.Deployment{}
	err := ac.client.Get(ctx, types.NamespacedName{
		Namespace: app.ObjectMeta.Namespace,
		Name:      app.ObjectMeta.Name,
	}, deployment)

	if err != nil {
		if errors.IsNotFound(err) {
			// Deployment不存在，重新创建
			app.Status.Phase = "Creating"
			return ac.handleCreation(ctx, app)
		}
		return reconcile.Result{}, err
	}

	// 更新应用状态
	app.Status.Replicas = deployment.Status.Replicas
	app.Status.ReadyReplicas = deployment.Status.ReadyReplicas
	app.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	app.Status.LastUpdateTime = metav1.Now()

	// 检查是否需要扩缩容
	if deployment.Spec.Replicas == nil || app.Spec.Replicas != *deployment.Spec.Replicas {
		app.Status.Phase = "Scaling"
		return ac.handleScaling(ctx, app)
	}

	// 检查是否需要更新
	if app.Spec.Image != deployment.Spec.Template.Spec.Containers[0].Image {
		app.Status.Phase = "Updating"
		return ac.handleUpdating(ctx, app)
	}

	// 检查健康状态
	if deployment.Status.AvailableReplicas < deployment.Status.Replicas {
		ac.recorder.RecordApplicationUnhealthy(app, "Not all replicas are available")
	} else {
		ac.recorder.RecordApplicationHealthy(app)
	}

	if err := ac.client.Status().Update(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

// handleScaling 处理扩缩容
func (ac *ApplicationController) handleScaling(ctx context.Context, app *Application) (reconcile.Result, error) {
	ac.recorder.RecordApplicationScaling(app, fmt.Sprintf("Scaling to %d replicas", app.Spec.Replicas))

	// 更新Deployment副本数
	deployment := &appsv1.Deployment{}
	err := ac.client.Get(ctx, types.NamespacedName{
		Namespace: app.ObjectMeta.Namespace,
		Name:      app.ObjectMeta.Name,
	}, deployment)

	if err != nil {
		return reconcile.Result{}, err
	}

	deployment.Spec.Replicas = pointer.Int32(app.Spec.Replicas)
	if err := ac.client.Update(ctx, deployment); err != nil {
		return reconcile.Result{}, err
	}

	// 等待扩缩容完成
	if deployment.Status.AvailableReplicas == app.Spec.Replicas {
		app.Status.Phase = "Running"
		ac.recorder.RecordApplicationScaled(app, fmt.Sprintf("Successfully scaled to %d replicas", app.Spec.Replicas))
	}

	if err := ac.client.Status().Update(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

// handleUpdating 处理应用更新
func (ac *ApplicationController) handleUpdating(ctx context.Context, app *Application) (reconcile.Result, error) {
	ac.recorder.RecordApplicationUpdating(app, fmt.Sprintf("Updating to image: %s", app.Spec.Image))

	// 更新Deployment镜像
	deployment := &appsv1.Deployment{}
	err := ac.client.Get(ctx, types.NamespacedName{
		Namespace: app.ObjectMeta.Namespace,
		Name:      app.ObjectMeta.Name,
	}, deployment)

	if err != nil {
		return reconcile.Result{}, err
	}

	deployment.Spec.Template.Spec.Containers[0].Image = app.Spec.Image
	if err := ac.client.Update(ctx, deployment); err != nil {
		return reconcile.Result{}, err
	}

	// 等待更新完成
	if deployment.Status.UpdatedReplicas == deployment.Status.Replicas &&
		deployment.Status.AvailableReplicas == deployment.Status.Replicas {
		app.Status.Phase = "Running"
		ac.recorder.RecordApplicationUpdated(app, "Successfully updated application")
	}

	if err := ac.client.Status().Update(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

// handleFailure 处理失败状态
func (ac *ApplicationController) handleFailure(ctx context.Context, app *Application) (reconcile.Result, error) {
	// 尝试恢复
	ac.recorder.RecordApplicationRecovery(app, "Attempting to recover from failure")

	// 重新创建资源
	return ac.handleCreation(ctx, app)
}

// handleDeletion 处理删除
func (ac *ApplicationController) handleDeletion(ctx context.Context, app *Application) (reconcile.Result, error) {
	ac.recorder.RecordApplicationDeleting(app)

	// 删除相关资源
	if err := ac.cleanupResources(ctx, types.NamespacedName{
		Namespace: app.ObjectMeta.Namespace,
		Name:      app.ObjectMeta.Name,
	}); err != nil {
		return reconcile.Result{}, err
	}

	// 移除finalizer
	if err := ac.removeFinalizer(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	ac.recorder.RecordApplicationDeleted(app)
	return reconcile.Result{}, nil
}

// 辅助方法
func (ac *ApplicationController) addFinalizer(ctx context.Context, app *Application) error {
	if !containsString(app.ObjectMeta.Finalizers, "application.finalizers.example.com") {
		app.ObjectMeta.Finalizers = append(app.ObjectMeta.Finalizers, "application.finalizers.example.com")
		return ac.client.Update(ctx, app)
	}
	return nil
}

func (ac *ApplicationController) removeFinalizer(ctx context.Context, app *Application) error {
	app.ObjectMeta.Finalizers = removeString(app.ObjectMeta.Finalizers, "application.finalizers.example.com")
	return ac.client.Update(ctx, app)
}

func (ac *ApplicationController) getLabels(app *Application) map[string]string {
	return map[string]string{
		"app":        app.ObjectMeta.Name,
		"managed-by": "application-controller",
	}
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) []string {
	for i, item := range slice {
		if item == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
