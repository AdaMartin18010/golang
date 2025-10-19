package kubernetes_operator

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Application 自定义资源定义
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ApplicationSpec   `json:"spec,omitempty"`
	Status            ApplicationStatus `json:"status,omitempty"`
}

// ApplicationSpec 应用规格
type ApplicationSpec struct {
	Replicas    int32                `json:"replicas"`
	Image       string               `json:"image"`
	Port        int32                `json:"port"`
	Environment map[string]string    `json:"environment,omitempty"`
	Resources   ResourceRequirements `json:"resources,omitempty"`
	HealthCheck HealthCheck          `json:"healthCheck,omitempty"`
	Scaling     ScalingPolicy        `json:"scaling,omitempty"`
}

// ApplicationStatus 应用状态
type ApplicationStatus struct {
	Phase      string             `json:"phase"`
	Replicas   int32              `json:"replicas"`
	Ready      int32              `json:"ready"`
	Available  int32              `json:"available"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ResourceRequirements 资源需求
type ResourceRequirements struct {
	Requests ResourceList `json:"requests,omitempty"`
	Limits   ResourceList `json:"limits,omitempty"`
}

// ResourceList 资源列表
type ResourceList struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// HealthCheck 健康检查
type HealthCheck struct {
	Path                string `json:"path"`
	InitialDelaySeconds int32  `json:"initialDelaySeconds"`
	PeriodSeconds       int32  `json:"periodSeconds"`
	TimeoutSeconds      int32  `json:"timeoutSeconds"`
	FailureThreshold    int32  `json:"failureThreshold"`
	SuccessThreshold    int32  `json:"successThreshold"`
}

// ScalingPolicy 扩缩容策略
type ScalingPolicy struct {
	MinReplicas                       int32 `json:"minReplicas"`
	MaxReplicas                       int32 `json:"maxReplicas"`
	TargetCPUUtilizationPercentage    int32 `json:"targetCPUUtilizationPercentage"`
	TargetMemoryUtilizationPercentage int32 `json:"targetMemoryUtilizationPercentage"`
}

// ApplicationList 应用列表
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

// ApplicationController 应用控制器
type ApplicationController struct {
	client   client.Client
	scheme   *runtime.Scheme
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
	recorder *EventRecorder
	metrics  *MetricsCollector
}

// NewApplicationController 创建新的应用控制器
func NewApplicationController(mgr manager.Manager) (*ApplicationController, error) {
	// 创建工作队列
	queue := workqueue.NewNamedRateLimitingQueue(
		workqueue.DefaultControllerRateLimiter(),
		"application-controller",
	)

	// 创建事件记录器
	recorder := NewEventRecorder(mgr.GetEventRecorderFor("application-controller"))

	// 创建指标收集器
	metrics := NewMetricsCollector()

	controller := &ApplicationController{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		queue:    queue,
		recorder: recorder,
		metrics:  metrics,
	}

	// 设置事件处理器
	if err := controller.setupEventHandlers(mgr); err != nil {
		return nil, err
	}

	return controller, nil
}

// setupEventHandlers 设置事件处理器
func (ac *ApplicationController) setupEventHandlers(mgr manager.Manager) error {
	// 创建控制器
	c, err := controller.New("application-controller", mgr, controller.Options{
		Reconciler: ac,
	})
	if err != nil {
		return err
	}

	// 设置事件源
	if err := c.Watch(
		&source.Kind{Type: &Application{}},
		&handler.EnqueueRequestForObject{},
		predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				ac.recorder.Event(e.Object, "Normal", "Created", "Application created")
				return true
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				ac.recorder.Event(e.ObjectNew, "Normal", "Updated", "Application updated")
				return true
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				ac.recorder.Event(e.Object, "Normal", "Deleted", "Application deleted")
				return true
			},
		},
	); err != nil {
		return err
	}

	return nil
}

// Reconcile 调和逻辑
func (ac *ApplicationController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	// 获取应用实例
	app := &Application{}
	if err := ac.client.Get(ctx, req.NamespacedName, app); err != nil {
		// 应用不存在，可能是被删除了
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// 记录调和开始
	ac.metrics.RecordReconcileStart(app.Name)

	// 执行调和逻辑
	result, err := ac.reconcileApplication(ctx, app)
	if err != nil {
		ac.metrics.RecordReconcileError(app.Name)
		ac.recorder.Event(app, "Warning", "ReconcileError", err.Error())
		return reconcile.Result{}, err
	}

	// 记录调和成功
	ac.metrics.RecordReconcileSuccess(app.Name)

	return result, nil
}

// reconcileApplication 调和应用
func (ac *ApplicationController) reconcileApplication(ctx context.Context, app *Application) (reconcile.Result, error) {
	// 检查应用是否被标记为删除
	if app.DeletionTimestamp != nil {
		return ac.handleDeletion(ctx, app)
	}

	// 确保最终器存在
	if err := ac.ensureFinalizer(ctx, app); err != nil {
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
	// 更新状态为创建中
	if app.Status.Phase != "Creating" {
		app.Status.Phase = "Creating"
		app.Status.Conditions = append(app.Status.Conditions, metav1.Condition{
			Type:               "Creating",
			Status:             "True",
			LastTransitionTime: metav1.Now(),
			Reason:             "ApplicationCreating",
			Message:            "Application is being created",
		})
		if err := ac.client.Status().Update(ctx, app); err != nil {
			return reconcile.Result{}, err
		}
	}

	// 创建Deployment
	if err := ac.createDeployment(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 创建Service
	if err := ac.createService(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 创建HPA（如果配置了自动扩缩容）
	if app.Spec.Scaling.MaxReplicas > app.Spec.Scaling.MinReplicas {
		if err := ac.createHPA(ctx, app); err != nil {
			return reconcile.Result{}, err
		}
	}

	// 更新状态为运行中
	app.Status.Phase = "Running"
	app.Status.Replicas = app.Spec.Replicas
	app.Status.Ready = 0
	app.Status.Available = 0
	app.Status.Conditions = append(app.Status.Conditions, metav1.Condition{
		Type:               "Running",
		Status:             "True",
		LastTransitionTime: metav1.Now(),
		Reason:             "ApplicationRunning",
		Message:            "Application is running",
	})

	if err := ac.client.Status().Update(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 记录创建成功
	ac.recorder.Event(app, "Normal", "Created", "Application created successfully")
	ac.metrics.RecordApplicationCreated(app.Name)

	return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

// handleRunning 处理运行中的应用
func (ac *ApplicationController) handleRunning(ctx context.Context, app *Application) (reconcile.Result, error) {
	// 检查Deployment状态
	deployment, err := ac.getDeployment(ctx, app)
	if err != nil {
		return reconcile.Result{}, err
	}

	// 更新应用状态
	app.Status.Replicas = deployment.Status.Replicas
	app.Status.Ready = deployment.Status.ReadyReplicas
	app.Status.Available = deployment.Status.AvailableReplicas

	// 检查是否需要扩缩容
	if app.Spec.Scaling.MaxReplicas > app.Spec.Scaling.MinReplicas {
		if err := ac.checkScaling(ctx, app); err != nil {
			return reconcile.Result{}, err
		}
	}

	// 检查健康状态
	if err := ac.checkHealth(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	if err := ac.client.Status().Update(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

// handleScaling 处理扩缩容
func (ac *ApplicationController) handleScaling(ctx context.Context, app *Application) (reconcile.Result, error) {
	// 更新HPA
	if err := ac.updateHPA(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 更新状态
	app.Status.Phase = "Running"
	app.Status.Conditions = append(app.Status.Conditions, metav1.Condition{
		Type:               "Scaled",
		Status:             "True",
		LastTransitionTime: metav1.Now(),
		Reason:             "ApplicationScaled",
		Message:            "Application scaling completed",
	})

	if err := ac.client.Status().Update(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	ac.recorder.Event(app, "Normal", "Scaled", "Application scaled successfully")

	return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

// handleUpdating 处理应用更新
func (ac *ApplicationController) handleUpdating(ctx context.Context, app *Application) (reconcile.Result, error) {
	// 更新Deployment
	if err := ac.updateDeployment(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 检查更新是否完成
	deployment, err := ac.getDeployment(ctx, app)
	if err != nil {
		return reconcile.Result{}, err
	}

	if deployment.Status.UpdatedReplicas == deployment.Status.Replicas &&
		deployment.Status.AvailableReplicas == deployment.Status.Replicas {
		// 更新完成
		app.Status.Phase = "Running"
		app.Status.Conditions = append(app.Status.Conditions, metav1.Condition{
			Type:               "Updated",
			Status:             "True",
			LastTransitionTime: metav1.Now(),
			Reason:             "ApplicationUpdated",
			Message:            "Application updated successfully",
		})

		if err := ac.client.Status().Update(ctx, app); err != nil {
			return reconcile.Result{}, err
		}

		ac.recorder.Event(app, "Normal", "Updated", "Application updated successfully")
	}

	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

// handleFailure 处理失败状态
func (ac *ApplicationController) handleFailure(ctx context.Context, app *Application) (reconcile.Result, error) {
	// 检查是否可以恢复
	if err := ac.checkRecovery(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 尝试恢复
	if err := ac.attemptRecovery(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: 60 * time.Second}, nil
}

// handleDeletion 处理删除
func (ac *ApplicationController) handleDeletion(ctx context.Context, app *Application) (reconcile.Result, error) {
	// 删除相关资源
	if err := ac.deleteResources(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	// 移除最终器
	if err := ac.removeFinalizer(ctx, app); err != nil {
		return reconcile.Result{}, err
	}

	ac.recorder.Event(app, "Normal", "Deleted", "Application deleted successfully")
	ac.metrics.RecordApplicationDeleted(app.Name)

	return reconcile.Result{}, nil
}

// 辅助方法

func (ac *ApplicationController) ensureFinalizer(ctx context.Context, app *Application) error {
	if !containsString(app.Finalizers, "application.finalizers.example.com") {
		app.Finalizers = append(app.Finalizers, "application.finalizers.example.com")
		return ac.client.Update(ctx, app)
	}
	return nil
}

func (ac *ApplicationController) removeFinalizer(ctx context.Context, app *Application) error {
	app.Finalizers = removeString(app.Finalizers, "application.finalizers.example.com")
	return ac.client.Update(ctx, app)
}

func (ac *ApplicationController) createDeployment(ctx context.Context, app *Application) error {
	// 实现Deployment创建逻辑
	return nil
}

func (ac *ApplicationController) createService(ctx context.Context, app *Application) error {
	// 实现Service创建逻辑
	return nil
}

func (ac *ApplicationController) createHPA(ctx context.Context, app *Application) error {
	// 实现HPA创建逻辑
	return nil
}

func (ac *ApplicationController) getDeployment(ctx context.Context, app *Application) (interface{}, error) {
	// 实现获取Deployment逻辑
	return nil, nil
}

func (ac *ApplicationController) checkScaling(ctx context.Context, app *Application) error {
	// 实现扩缩容检查逻辑
	return nil
}

func (ac *ApplicationController) checkHealth(ctx context.Context, app *Application) error {
	// 实现健康检查逻辑
	return nil
}

func (ac *ApplicationController) updateHPA(ctx context.Context, app *Application) error {
	// 实现HPA更新逻辑
	return nil
}

func (ac *ApplicationController) updateDeployment(ctx context.Context, app *Application) error {
	// 实现Deployment更新逻辑
	return nil
}

func (ac *ApplicationController) checkRecovery(ctx context.Context, app *Application) error {
	// 实现恢复检查逻辑
	return nil
}

func (ac *ApplicationController) attemptRecovery(ctx context.Context, app *Application) error {
	// 实现恢复尝试逻辑
	return nil
}

func (ac *ApplicationController) deleteResources(ctx context.Context, app *Application) error {
	// 实现资源删除逻辑
	return nil
}

// 工具函数
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}
