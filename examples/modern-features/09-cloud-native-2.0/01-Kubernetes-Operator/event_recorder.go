package kubernetes_operator

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

// EventRecorder 事件记录器
type EventRecorder struct {
	recorder record.EventRecorder
}

// NewEventRecorder 创建新的事件记录器
func NewEventRecorder(recorder record.EventRecorder) *EventRecorder {
	return &EventRecorder{
		recorder: recorder,
	}
}

// Event 记录事件
func (er *EventRecorder) Event(object runtime.Object, eventType, reason, message string) {
	er.recorder.Event(object, eventType, reason, message)
}

// Eventf 格式化记录事件
func (er *EventRecorder) Eventf(object runtime.Object, eventType, reason, messageFmt string, args ...interface{}) {
	er.recorder.Eventf(object, eventType, reason, messageFmt, args...)
}

// AnnotatedEventf 带注解的事件记录
func (er *EventRecorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventType, reason, messageFmt string, args ...interface{}) {
	er.recorder.AnnotatedEventf(object, annotations, eventType, reason, messageFmt, args...)
}

// RecordApplicationCreated 记录应用创建事件
func (er *EventRecorder) RecordApplicationCreated(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "ApplicationCreated",
		fmt.Sprintf("Application %s created successfully", app.ObjectMeta.Name))
}

// RecordApplicationRunning 记录应用运行事件
func (er *EventRecorder) RecordApplicationRunning(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "ApplicationRunning",
		fmt.Sprintf("Application %s is now running", app.ObjectMeta.Name))
}

// RecordApplicationHealthy 记录应用健康事件
func (er *EventRecorder) RecordApplicationHealthy(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "ApplicationHealthy",
		fmt.Sprintf("Application %s is healthy", app.ObjectMeta.Name))
}

// RecordApplicationUnhealthy 记录应用不健康事件
func (er *EventRecorder) RecordApplicationUnhealthy(app *Application, message string) {
	er.Eventf(app, corev1.EventTypeWarning, "ApplicationUnhealthy",
		"Application %s is unhealthy: %s", app.ObjectMeta.Name, message)
}

// RecordApplicationUpdated 记录应用更新事件
func (er *EventRecorder) RecordApplicationUpdated(app *Application, message string) {
	er.Eventf(app, corev1.EventTypeNormal, "ApplicationUpdated",
		"Application %s updated: %s", app.ObjectMeta.Name, message)
}

// RecordApplicationUpdating 记录应用更新中事件
func (er *EventRecorder) RecordApplicationUpdating(app *Application, message string) {
	er.Eventf(app, corev1.EventTypeNormal, "ApplicationUpdating",
		"Application %s updating: %s", app.ObjectMeta.Name, message)
}

// RecordApplicationDeleted 记录应用删除事件
func (er *EventRecorder) RecordApplicationDeleted(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "ApplicationDeleted",
		fmt.Sprintf("Application %s deleted successfully", app.ObjectMeta.Name))
}

// RecordApplicationDeleting 记录应用删除中事件
func (er *EventRecorder) RecordApplicationDeleting(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "ApplicationDeleting",
		fmt.Sprintf("Application %s is being deleted", app.ObjectMeta.Name))
}

// RecordApplicationScaling 记录应用扩缩容中事件
func (er *EventRecorder) RecordApplicationScaling(app *Application, message string) {
	er.Eventf(app, corev1.EventTypeNormal, "ApplicationScaling",
		"Application %s scaling: %s", app.ObjectMeta.Name, message)
}

// RecordApplicationScaled 记录应用扩缩容完成事件
func (er *EventRecorder) RecordApplicationScaled(app *Application, message string) {
	er.Eventf(app, corev1.EventTypeNormal, "ApplicationScaled",
		"Application %s scaled: %s", app.ObjectMeta.Name, message)
}

// RecordApplicationFailed 记录应用失败事件
func (er *EventRecorder) RecordApplicationFailed(app *Application, reason, message string) {
	er.Eventf(app, corev1.EventTypeWarning, "ApplicationFailed",
		"Application %s failed: %s - %s", app.ObjectMeta.Name, reason, message)
}

// RecordApplicationRecovered 记录应用恢复事件
func (er *EventRecorder) RecordApplicationRecovered(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "ApplicationRecovered",
		fmt.Sprintf("Application %s recovered from failure", app.ObjectMeta.Name))
}

// RecordApplicationRecovery 记录应用恢复中事件
func (er *EventRecorder) RecordApplicationRecovery(app *Application, message string) {
	er.Eventf(app, corev1.EventTypeNormal, "ApplicationRecovery",
		"Application %s recovering: %s", app.ObjectMeta.Name, message)
}

// RecordDeploymentCreated 记录Deployment创建事件
func (er *EventRecorder) RecordDeploymentCreated(app *Application, deploymentName string) {
	er.Eventf(app, corev1.EventTypeNormal, "DeploymentCreated",
		"Deployment %s created for application %s", deploymentName, app.ObjectMeta.Name)
}

// RecordServiceCreated 记录Service创建事件
func (er *EventRecorder) RecordServiceCreated(app *Application, serviceName string) {
	er.Eventf(app, corev1.EventTypeNormal, "ServiceCreated",
		"Service %s created for application %s", serviceName, app.ObjectMeta.Name)
}

// RecordHPACreated 记录HPA创建事件
func (er *EventRecorder) RecordHPACreated(app *Application, hpaName string) {
	er.Eventf(app, corev1.EventTypeNormal, "HPACreated",
		"HPA %s created for application %s", hpaName, app.ObjectMeta.Name)
}

// RecordHealthCheckFailed 记录健康检查失败事件
func (er *EventRecorder) RecordHealthCheckFailed(app *Application, reason string) {
	er.Eventf(app, corev1.EventTypeWarning, "HealthCheckFailed",
		"Health check failed for application %s: %s", app.ObjectMeta.Name, reason)
}

// RecordHealthCheckPassed 记录健康检查通过事件
func (er *EventRecorder) RecordHealthCheckPassed(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "HealthCheckPassed",
		fmt.Sprintf("Health check passed for application %s", app.ObjectMeta.Name))
}

// RecordResourceQuotaExceeded 记录资源配额超限事件
func (er *EventRecorder) RecordResourceQuotaExceeded(app *Application, resource string) {
	er.Eventf(app, corev1.EventTypeWarning, "ResourceQuotaExceeded",
		"Resource quota exceeded for application %s: %s", app.ObjectMeta.Name, resource)
}

// RecordImagePullFailed 记录镜像拉取失败事件
func (er *EventRecorder) RecordImagePullFailed(app *Application, image string, reason string) {
	er.Eventf(app, corev1.EventTypeWarning, "ImagePullFailed",
		"Failed to pull image %s for application %s: %s", image, app.ObjectMeta.Name, reason)
}

// RecordPodScheduled 记录Pod调度事件
func (er *EventRecorder) RecordPodScheduled(app *Application, podName, nodeName string) {
	er.Eventf(app, corev1.EventTypeNormal, "PodScheduled",
		"Pod %s for application %s scheduled on node %s", podName, app.ObjectMeta.Name, nodeName)
}

// RecordPodFailed 记录Pod失败事件
func (er *EventRecorder) RecordPodFailed(app *Application, podName, reason string) {
	er.Eventf(app, corev1.EventTypeWarning, "PodFailed",
		"Pod %s for application %s failed: %s", podName, app.ObjectMeta.Name, reason)
}

// RecordReconcileError 记录调和错误事件
func (er *EventRecorder) RecordReconcileError(app *Application, error string) {
	er.Eventf(app, corev1.EventTypeWarning, "ReconcileError",
		"Reconcile error for application %s: %s", app.ObjectMeta.Name, error)
}

// RecordReconcileSuccess 记录调和成功事件
func (er *EventRecorder) RecordReconcileSuccess(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "ReconcileSuccess",
		fmt.Sprintf("Reconcile successful for application %s", app.ObjectMeta.Name))
}

// RecordBackupCreated 记录备份创建事件
func (er *EventRecorder) RecordBackupCreated(app *Application, backupName string) {
	er.Eventf(app, corev1.EventTypeNormal, "BackupCreated",
		"Backup %s created for application %s", backupName, app.ObjectMeta.Name)
}

// RecordRestoreCompleted 记录恢复完成事件
func (er *EventRecorder) RecordRestoreCompleted(app *Application, backupName string) {
	er.Eventf(app, corev1.EventTypeNormal, "RestoreCompleted",
		"Restore from backup %s completed for application %s", backupName, app.ObjectMeta.Name)
}

// RecordMigrationStarted 记录迁移开始事件
func (er *EventRecorder) RecordMigrationStarted(app *Application, targetNamespace string) {
	er.Eventf(app, corev1.EventTypeNormal, "MigrationStarted",
		"Migration of application %s to namespace %s started", app.ObjectMeta.Name, targetNamespace)
}

// RecordMigrationCompleted 记录迁移完成事件
func (er *EventRecorder) RecordMigrationCompleted(app *Application, targetNamespace string) {
	er.Eventf(app, corev1.EventTypeNormal, "MigrationCompleted",
		"Migration of application %s to namespace %s completed", app.ObjectMeta.Name, targetNamespace)
}

// RecordRollbackStarted 记录回滚开始事件
func (er *EventRecorder) RecordRollbackStarted(app *Application, version string) {
	er.Eventf(app, corev1.EventTypeNormal, "RollbackStarted",
		"Rollback of application %s to version %s started", app.ObjectMeta.Name, version)
}

// RecordRollbackCompleted 记录回滚完成事件
func (er *EventRecorder) RecordRollbackCompleted(app *Application, version string) {
	er.Eventf(app, corev1.EventTypeNormal, "RollbackCompleted",
		"Rollback of application %s to version %s completed", app.ObjectMeta.Name, version)
}

// RecordConfigMapUpdated 记录ConfigMap更新事件
func (er *EventRecorder) RecordConfigMapUpdated(app *Application, configMapName string) {
	er.Eventf(app, corev1.EventTypeNormal, "ConfigMapUpdated",
		"ConfigMap %s updated for application %s", configMapName, app.ObjectMeta.Name)
}

// RecordSecretUpdated 记录Secret更新事件
func (er *EventRecorder) RecordSecretUpdated(app *Application, secretName string) {
	er.Eventf(app, corev1.EventTypeNormal, "SecretUpdated",
		"Secret %s updated for application %s", secretName, app.ObjectMeta.Name)
}

// RecordNetworkPolicyApplied 记录网络策略应用事件
func (er *EventRecorder) RecordNetworkPolicyApplied(app *Application, policyName string) {
	er.Eventf(app, corev1.EventTypeNormal, "NetworkPolicyApplied",
		"Network policy %s applied to application %s", policyName, app.ObjectMeta.Name)
}

// RecordIngressCreated 记录Ingress创建事件
func (er *EventRecorder) RecordIngressCreated(app *Application, ingressName string) {
	er.Eventf(app, corev1.EventTypeNormal, "IngressCreated",
		"Ingress %s created for application %s", ingressName, app.ObjectMeta.Name)
}

// RecordCertificateIssued 记录证书签发事件
func (er *EventRecorder) RecordCertificateIssued(app *Application, certName string) {
	er.Eventf(app, corev1.EventTypeNormal, "CertificateIssued",
		"Certificate %s issued for application %s", certName, app.ObjectMeta.Name)
}

// RecordCertificateExpiring 记录证书即将过期事件
func (er *EventRecorder) RecordCertificateExpiring(app *Application, certName string, daysLeft int) {
	er.Eventf(app, corev1.EventTypeWarning, "CertificateExpiring",
		"Certificate %s for application %s expires in %d days", certName, app.ObjectMeta.Name, daysLeft)
}

// RecordCertificateExpired 记录证书过期事件
func (er *EventRecorder) RecordCertificateExpired(app *Application, certName string) {
	er.Eventf(app, corev1.EventTypeWarning, "CertificateExpired",
		"Certificate %s for application %s has expired", certName, app.ObjectMeta.Name)
}

// RecordStorageProvisioned 记录存储供应事件
func (er *EventRecorder) RecordStorageProvisioned(app *Application, pvcName string) {
	er.Eventf(app, corev1.EventTypeNormal, "StorageProvisioned",
		"Storage provisioned for application %s: %s", app.ObjectMeta.Name, pvcName)
}

// RecordStorageFailed 记录存储失败事件
func (er *EventRecorder) RecordStorageFailed(app *Application, pvcName, reason string) {
	er.Eventf(app, corev1.EventTypeWarning, "StorageFailed",
		"Storage failed for application %s: %s - %s", app.ObjectMeta.Name, pvcName, reason)
}

// RecordMonitoringEnabled 记录监控启用事件
func (er *EventRecorder) RecordMonitoringEnabled(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "MonitoringEnabled",
		fmt.Sprintf("Monitoring enabled for application %s", app.ObjectMeta.Name))
}

// RecordLoggingEnabled 记录日志启用事件
func (er *EventRecorder) RecordLoggingEnabled(app *Application) {
	er.Event(app, corev1.EventTypeNormal, "LoggingEnabled",
		fmt.Sprintf("Logging enabled for application %s", app.ObjectMeta.Name))
}

// RecordSecurityScanCompleted 记录安全扫描完成事件
func (er *EventRecorder) RecordSecurityScanCompleted(app *Application, vulnerabilities int) {
	if vulnerabilities == 0 {
		er.Event(app, corev1.EventTypeNormal, "SecurityScanCompleted",
			fmt.Sprintf("Security scan completed for application %s: no vulnerabilities found", app.ObjectMeta.Name))
	} else {
		er.Eventf(app, corev1.EventTypeWarning, "SecurityScanCompleted",
			"Security scan completed for application %s: %d vulnerabilities found", app.ObjectMeta.Name, vulnerabilities)
	}
}

// RecordComplianceCheckPassed 记录合规检查通过事件
func (er *EventRecorder) RecordComplianceCheckPassed(app *Application, standard string) {
	er.Eventf(app, corev1.EventTypeNormal, "ComplianceCheckPassed",
		"Compliance check passed for application %s: %s", app.ObjectMeta.Name, standard)
}

// RecordComplianceCheckFailed 记录合规检查失败事件
func (er *EventRecorder) RecordComplianceCheckFailed(app *Application, standard, reason string) {
	er.Eventf(app, corev1.EventTypeWarning, "ComplianceCheckFailed",
		"Compliance check failed for application %s: %s - %s", app.ObjectMeta.Name, standard, reason)
}
