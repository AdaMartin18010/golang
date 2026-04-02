# Kubernetes CronJob Controller V2 深度解析

> **分类**: 工程与云原生
> **标签**: #kubernetes #cronjob #controller #v2
> **参考**: `k8s.io/kubernetes/pkg/controller/cronjob`

---

## CronJob Controller V2 架构概览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Kubernetes CronJob Controller V2                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Controller Loop                                  │   │
│  │                                                                      │   │
│  │  ┌──────────┐     ┌──────────┐     ┌──────────┐                    │   │
│  │  │  Informer│────►│   Queue  │────►│  Worker  │                    │   │
│  │  │  (Watch) │     │          │     │          │                    │   │
│  │  └──────────┘     └──────────┘     └────┬─────┘                    │   │
│  │                                         │                          │   │
│  │                                         ▼                          │   │
│  │  ┌─────────────────────────────────────────────────────────────┐  │   │
│  │  │              syncCronJob() - 核心同步逻辑                    │  │   │
│  │  │                                                              │  │   │
│  │  │  1. 获取 CronJob 和关联 Job 列表                              │  │   │
│  │  │  2. 计算需要执行的调度时间                                     │  │   │
│  │  │  3. 并发控制（StartingDeadlineSeconds）                        │  │   │
│  │  │  4. 并发策略处理（Allow/Forbid/Replace）                       │  │   │
│  │  │  5. 创建 Job 资源                                             │  │   │
│  │  │  6. 清理历史 Job（Successful/Failed Job History）               │  │   │
│  │  │  7. 更新 CronJob 状态（LastScheduleTime, Active）              │  │   │
│  │  │                                                              │  │   │
│  │  └─────────────────────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    Key Data Structures                               │   │
│  │                                                                      │   │
│  │  cronJobController:                                                  │   │
│  │    - kubeClient: Kubernetes 客户端                                  │   │
│  │    - cronJobLister: CronJob 缓存列表                                │   │
│  │    - jobLister: Job 缓存列表                                        │   │
│  │    - queue: 工作队列                                                 │   │
│  │    - recorder: 事件记录器                                            │   │
│  │    - syncHandler: 同步处理函数                                       │   │
│  │                                                                      │   │
│  │  cronJobStore: map[types.UID]*CronJob                               │   │
│  │  jobStore:    map[types.UID]*Job                                    │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 核心数据结构

```go
// pkg/controller/cronjob/cronjob_controllerv2.go (简化版)

// ControllerV2 CronJob Controller V2 实现
type ControllerV2 struct {
    kubeClient clientset.Interface

    // Lister 缓存
    cronJobLister batchv1listers.CronJobLister
    cronJobsSynced cache.InformerSynced
    jobLister      batchv1listers.JobLister
    jobsSynced     cache.InformerSynced

    // 工作队列
    queue workqueue.RateLimitingInterface

    // 事件记录
    recorder record.EventRecorder

    // 时区支持
    now func() time.Time
}

// CronJob 关键字段解析
// CronJobSpec 定义了定时任务规范
type CronJobSpec struct {
    // Schedule: Cron 表达式
    // 格式: "分 时 日 月 周" (标准 Unix cron)
    // 特殊扩展:
    //   - "@yearly", "@monthly", "@weekly", "@daily", "@hourly"
    //   - "TZ=Asia/Shanghai 0 9 * * *"  # 时区支持
    Schedule string

    // TimeZone: 时区名称 ("Asia/Shanghai", "America/New_York")
    // 使用 IANA Time Zone database 名称
    TimeZone *string

    // ConcurrencyPolicy: 并发策略
    // - Allow (默认): 允许并发执行
    // - Forbid: 禁止并发，跳过新调度
    // - Replace: 取消旧任务，执行新任务
    ConcurrencyPolicy ConcurrencyPolicy

    // StartingDeadlineSeconds: 任务启动截止时间
    // 如果任务错过了超过此时间的调度，则跳过
    StartingDeadlineSeconds *int64

    // Suspend: 暂停 CronJob
    Suspend *bool

    // JobTemplate: Job 模板
    // 定义创建的 Job 的规范
    JobTemplate JobTemplateSpec

    // SuccessfulJobHistoryLimit: 保留的成功 Job 数量 (默认 3)
    SuccessfulJobHistoryLimit *int32

    // FailedJobHistoryLimit: 保留的失败 Job 数量 (默认 1)
    FailedJobHistoryLimit *int32
}

// CronJobStatus 状态字段
type CronJobStatus struct {
    // Active: 当前正在运行的 Job 列表
    Active []corev1.ObjectReference

    // LastScheduleTime: 上次调度时间
    LastScheduleTime *metav1.Time

    // LastSuccessfulTime: 上次成功完成时间
    LastSuccessfulTime *metav1.Time
}
```

---

## 调度计算实现

```go
package cronjob

import (
    "time"

    "github.com/robfig/cron/v3"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CronJobScheduleDeadline 调度截止时间检查
type CronJobScheduleDeadline struct {
    StartingDeadlineSeconds *int64
    LastScheduleTime        *metav1.Time
    Now                     time.Time
}

// IsOverdue 检查调度是否已过期
func (d *CronJobScheduleDeadline) IsOverdue(scheduledTime time.Time) bool {
    // 如果没有设置 StartingDeadlineSeconds，不过期
    if d.StartingDeadlineSeconds == nil {
        return false
    }

    deadline := *d.StartingDeadlineSeconds

    // 计算从调度时间到当前时间的差值
    elapsed := int64(d.Now.Sub(scheduledTime).Seconds())

    // 如果超过截止时间，则过期
    return elapsed > deadline
}

// GetNextScheduleTimes 计算下次调度时间
func GetNextScheduleTimes(
    cronJob *CronJob,
    now time.Time,
) (times []time.Time, err error) {
    // 解析 cron 表达式
    sched, err := parseCronSchedule(cronJob)
    if err != nil {
        return nil, err
    }

    // 获取上次调度时间
    var earliestTime time.Time
    if cronJob.Status.LastScheduleTime != nil {
        earliestTime = cronJob.Status.LastScheduleTime.Time
    } else {
        // 如果没有上次调度，使用 CreationTimestamp
        earliestTime = cronJob.CreationTimestamp.Time
    }

    // 如果设置了 StartingDeadlineSeconds，考虑截止时间
    if cronJob.Spec.StartingDeadlineSeconds != nil {
        deadline := now.Add(-time.Duration(*cronJob.Spec.StartingDeadlineSeconds) * time.Second)
        if deadline.After(earliestTime) {
            earliestTime = deadline
        }
    }

    // 计算错过的调度时间
    var missedRun, nextRun time.Time
    for t := sched.Next(earliestTime); ; t = sched.Next(t) {
        if t.IsZero() {
            break
        }

        // 只考虑早于现在的时间
        if t.After(now) {
            nextRun = t
            break
        }

        // 检查是否过期
        deadline := &CronJobScheduleDeadline{
            StartingDeadlineSeconds: cronJob.Spec.StartingDeadlineSeconds,
            LastScheduleTime:        cronJob.Status.LastScheduleTime,
            Now:                     now,
        }

        if deadline.IsOverdue(t) {
            // 跳过过期调度
            continue
        }

        missedRun = t
        times = append(times, t)

        // 限制单次处理的错过调度数量
        if len(times) >= 100 {
            break
        }
    }

    return times, nil
}

// parseCronSchedule 解析 cron 表达式，支持时区
func parseCronSchedule(cronJob *CronJob) (cron.Schedule, error) {
    schedule := cronJob.Spec.Schedule

    // 创建解析器
    parser := cron.NewParser(
        cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
    )

    // 如果设置了 TimeZone，使用时区
    if cronJob.Spec.TimeZone != nil && *cronJob.Spec.TimeZone != "" {
        loc, err := time.LoadLocation(*cronJob.Spec.TimeZone)
        if err != nil {
            return nil, err
        }

        return parser.Parse(schedule)
    }

    return parser.Parse(schedule)
}

// StandardCronExpressions 常用 Cron 表达式示例
var StandardCronExpressions = map[string]string{
    "EveryMinute":     "*/1 * * * *",
    "Every5Minutes":   "*/5 * * * *",
    "Every15Minutes":  "*/15 * * * *",
    "EveryHour":       "0 * * * *",
    "EveryDayAt9AM":   "0 9 * * *",
    "EveryMonday":     "0 0 * * 1",
    "EveryMonth1st":   "0 0 1 * *",
}
```

---

## 并发策略实现

```go
package cronjob

import (
    "context"
    "fmt"

    batchv1 "k8s.io/api/batch/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConcurrencyPolicy 并发策略类型
type ConcurrencyPolicy string

const (
    // AllowConcurrent 允许并发执行
    AllowConcurrent ConcurrencyPolicy = "Allow"

    // ForbidConcurrent 禁止并发，跳过新调度
    ForbidConcurrent ConcurrencyPolicy = "Forbid"

    // ReplaceConcurrent 取消旧任务，执行新任务
    ReplaceConcurrent ConcurrencyPolicy = "Replace"
)

// ConcurrencyController 并发控制
type ConcurrencyController struct {
    jobControl JobControlInterface
}

// ApplyConcurrencyPolicy 应用并发策略
func (c *ConcurrencyController) ApplyConcurrencyPolicy(
    ctx context.Context,
    cronJob *CronJob,
    activeJobs []*batchv1.Job,
    scheduledTime time.Time,
) (shouldCreate bool, jobsToDelete []*batchv1.Job, err error) {
    policy := cronJob.Spec.ConcurrencyPolicy

    switch policy {
    case AllowConcurrent, "":
        // 允许并发，直接创建
        return true, nil, nil

    case ForbidConcurrent:
        // 如果有活跃 Job，跳过
        if len(activeJobs) > 0 {
            // 记录事件：跳过调度
            return false, nil, nil
        }
        return true, nil, nil

    case ReplaceConcurrent:
        // 取消所有活跃 Job
        if len(activeJobs) > 0 {
            return true, activeJobs, nil
        }
        return true, nil, nil

    default:
        return false, nil, fmt.Errorf("invalid concurrency policy: %s", policy)
    }
}

// DeleteJobs 批量删除 Job
func (c *ConcurrencyController) DeleteJobs(
    ctx context.Context,
    jobs []*batchv1.Job,
    cronJob *CronJob,
    recorder record.EventRecorder,
) error {
    for _, job := range jobs {
        // 设置删除策略
        propagationPolicy := metav1.DeletePropagationBackground

        err := c.jobControl.DeleteJob(ctx, job.Namespace, job.Name, &propagationPolicy)
        if err != nil {
            recorder.Eventf(cronJob, corev1.EventTypeWarning, "FailedDelete",
                "Error deleting job %s: %v", job.Name, err)
            continue
        }

        recorder.Eventf(cronJob, corev1.EventTypeNormal, "SuccessfulDelete",
            "Deleted job %s", job.Name)
    }

    return nil
}

// 并发策略决策流程
/*
┌─────────────────────────────────────────────────────────────┐
│                  Concurrency Policy Flow                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────┐                                             │
│  │ ActiveJobs? │                                             │
│  └──────┬──────┘                                             │
│         │                                                    │
│     ┌───┴───┐                                                │
│     │       │                                                │
│    Yes      No ───────► CREATE JOB                           │
│     │                                                        │
│     ▼                                                        │
│  ┌─────────────────┐                                         │
│  │ ConcurrencyPolicy│                                        │
│  └────────┬────────┘                                         │
│           │                                                  │
│     ┌─────┼─────┐                                            │
│     │     │     │                                            │
│  Allow Forbid Replace                                        │
│     │     │     │                                            │
│     ▼     ▼     ▼                                            │
│  CREATE  SKIP  DELETE &                                      │
│             CREATE                                           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
*/
```

---

## Job 创建实现

```go
package cronjob

import (
    "context"
    "fmt"

    batchv1 "k8s.io/api/batch/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/utils/pointer"
)

// JobCreator 负责从 CronJob 创建 Job
type JobCreator struct {
    jobControl JobControlInterface
}

// CreateJobFromCronJob 从 CronJob 创建 Job
func (c *JobCreator) CreateJobFromCronJob(
    ctx context.Context,
    cronJob *CronJob,
    scheduledTime time.Time,
) (*batchv1.Job, error) {
    // 深拷贝 Job 模板
    jobSpec := *cronJob.Spec.JobTemplate.DeepCopy()

    // 生成 Job 名称
    // 格式: {cronjob-name}-{timestamp}
    name := fmt.Sprintf("%s-%d", cronJob.Name, scheduledTime.Unix())

    // 设置标签
    if jobSpec.Labels == nil {
        jobSpec.Labels = make(map[string]string)
    }
    jobSpec.Labels["app.kubernetes.io/managed-by"] = "cronjob-controller"
    jobSpec.Labels["cronjob-name"] = cronJob.Name

    // 设置注解
    if jobSpec.Annotations == nil {
        jobSpec.Annotations = make(map[string]string)
    }
    jobSpec.Annotations["cronjob.kubernetes.io/scheduled-timestamp"] =
        scheduledTime.Format(time.RFC3339)

    // 设置 OwnerReference，建立父子关系
    ownerRef := metav1.OwnerReference{
        APIVersion:         batchv1.SchemeGroupVersion.String(),
        Kind:               "CronJob",
        Name:               cronJob.Name,
        UID:                cronJob.UID,
        Controller:         pointer.Bool(true),
        BlockOwnerDeletion: pointer.Bool(true),
    }
    jobSpec.OwnerReferences = append(jobSpec.OwnerReferences, ownerRef)

    // 创建 Job
    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:            name,
            Namespace:       cronJob.Namespace,
            Labels:          jobSpec.Labels,
            Annotations:     jobSpec.Annotations,
            OwnerReferences: jobSpec.OwnerReferences,
        },
        Spec: jobSpec.Spec,
    }

    createdJob, err := c.jobControl.CreateJob(cronJob.Namespace, job)
    if err != nil {
        return nil, err
    }

    return createdJob, nil
}

// GetActiveJobs 获取活跃的 Job 列表
func GetActiveJobs(
    cronJob *CronJob,
    jobLister batchv1listers.JobLister,
) ([]*batchv1.Job, []*batchv1.Job, error) {
    // 通过 Label 选择器查找关联的 Job
    selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
        MatchLabels: map[string]string{
            "cronjob-name": cronJob.Name,
        },
    })
    if err != nil {
        return nil, nil, err
    }

    jobs, err := jobLister.Jobs(cronJob.Namespace).List(selector)
    if err != nil {
        return nil, nil, err
    }

    var activeJobs, successfulJobs, failedJobs []*batchv1.Job

    for _, job := range jobs {
        // 检查 Job 是否由当前 CronJob 拥有
        if !metav1.IsControlledBy(job, cronJob) {
            continue
        }

        // 根据 Job 状态分类
        if IsJobActive(job) {
            activeJobs = append(activeJobs, job)
        } else if IsJobSucceeded(job) {
            successfulJobs = append(successfulJobs, job)
        } else if IsJobFailed(job) {
            failedJobs = append(failedJobs, job)
        }
    }

    return activeJobs, append(successfulJobs, failedJobs...), nil
}

// IsJobActive 检查 Job 是否活跃
func IsJobActive(job *batchv1.Job) bool {
    return job.Status.Active > 0 ||
        (job.Status.StartTime != nil && job.Status.CompletionTime == nil)
}

// IsJobSucceeded 检查 Job 是否成功
func IsJobSucceeded(job *batchv1.Job) bool {
    for _, cond := range job.Status.Conditions {
        if cond.Type == batchv1.JobComplete && cond.Status == "True" {
            return true
        }
    }
    return false
}

// IsJobFailed 检查 Job 是否失败
func IsJobFailed(job *batchv1.Job) bool {
    for _, cond := range job.Status.Conditions {
        if cond.Type == batchv1.JobFailed && cond.Status == "True" {
            return true
        }
    }
    return false
}
```

---

## 历史清理实现

```go
package cronjob

import (
    "context"
    "sort"

    batchv1 "k8s.io/api/batch/v1"
)

// HistoryCleaner 历史 Job 清理器
type HistoryCleaner struct {
    jobControl JobControlInterface
}

// CleanupHistory 清理历史 Job
func (c *HistoryCleaner) CleanupHistory(
    ctx context.Context,
    cronJob *CronJob,
    successfulJobs []*batchv1.Job,
    failedJobs []*batchv1.Job,
) error {
    // 获取限制
    successfulLimit := int32(3)
    if cronJob.Spec.SuccessfulJobHistoryLimit != nil {
        successfulLimit = *cronJob.Spec.SuccessfulJobHistoryLimit
    }

    failedLimit := int32(1)
    if cronJob.Spec.FailedJobHistoryLimit != nil {
        failedLimit = *cronJob.Spec.FailedJobHistoryLimit
    }

    // 清理成功的 Job
    if err := c.cleanupSuccessfulJobs(ctx, successfulJobs, int(successfulLimit)); err != nil {
        return err
    }

    // 清理失败的 Job
    if err := c.cleanupFailedJobs(ctx, failedJobs, int(failedLimit)); err != nil {
        return err
    }

    return nil
}

func (c *HistoryCleaner) cleanupSuccessfulJobs(
    ctx context.Context,
    jobs []*batchv1.Job,
    limit int,
) error {
    if len(jobs) <= limit {
        return nil
n    }

    // 按完成时间排序（最近的在前）
    sort.Slice(jobs, func(i, j int) bool {
        if jobs[i].Status.CompletionTime == nil {
            return false
        }
        if jobs[j].Status.CompletionTime == nil {
            return true
        }
        return jobs[i].Status.CompletionTime.After(jobs[j].Status.CompletionTime.Time)
    })

    // 删除超出限制的 Job
    for i := limit; i < len(jobs); i++ {
        c.jobControl.DeleteJob(ctx, jobs[i].Namespace, jobs[i].Name, nil)
    }

    return nil
}

func (c *HistoryCleaner) cleanupFailedJobs(
    ctx context.Context,
    jobs []*batchv1.Job,
    limit int,
) error {
    if len(jobs) <= limit {
        return nil
    }

    // 按失败时间排序
    sort.Slice(jobs, func(i, j int) bool {
        timeI := getJobConditionTime(jobs[i], batchv1.JobFailed)
        timeJ := getJobConditionTime(jobs[j], batchv1.JobFailed)
        return timeI.After(timeJ)
    })

    // 删除超出限制的 Job
    for i := limit; i < len(jobs); i++ {
        c.jobControl.DeleteJob(ctx, jobs[i].Namespace, jobs[i].Name, nil)
    }

    return nil
}

func getJobConditionTime(job *batchv1.Job, condType batchv1.JobConditionType) time.Time {
    for _, cond := range job.Status.Conditions {
        if cond.Type == condType {
            return cond.LastTransitionTime.Time
        }
    }
    return time.Time{}
}
```

---

## 完整同步流程

```go
package cronjob

import (
    "context"
    "fmt"
    "time"

    batchv1 "k8s.io/api/batch/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SyncCronJob CronJob 完整同步逻辑
func (c *ControllerV2) SyncCronJob(ctx context.Context, cronJob *CronJob) error {
    now := c.now()

    // 1. 检查是否暂停
    if cronJob.Spec.Suspend != nil && *cronJob.Spec.Suspend {
        return nil
    }

    // 2. 获取所有关联的 Job
    activeJobs, inactiveJobs, err := GetActiveJobs(cronJob, c.jobLister)
    if err != nil {
        return err
    }

    // 3. 计算错过的调度时间
    missedSchedules, err := GetNextScheduleTimes(cronJob, now)
    if err != nil {
        c.recorder.Eventf(cronJob, corev1.EventTypeWarning, "InvalidSchedule",
            "Failed to parse schedule: %v", err)
        return err
    }

    // 4. 更新状态中的活跃 Job 列表
    cronJob.Status.Active = nil
    for _, job := range activeJobs {
        jobRef, err := ref.GetReference(scheme.Scheme, job)
        if err != nil {
            continue
        }
        cronJob.Status.Active = append(cronJob.Status.Active, *jobRef)
    }

    // 5. 处理错过的调度
    for _, scheduledTime := range missedSchedules {
        // 应用并发策略
        shouldCreate, jobsToDelete, err := c.concurrencyController.ApplyConcurrencyPolicy(
            ctx, cronJob, activeJobs, scheduledTime)
        if err != nil {
            return err
        }

        // 删除需要取消的 Job
        if len(jobsToDelete) > 0 {
            c.concurrencyController.DeleteJobs(ctx, jobsToDelete, cronJob, c.recorder)
        }

        if !shouldCreate {
            continue
        }

        // 创建 Job
        job, err := c.jobCreator.CreateJobFromCronJob(ctx, cronJob, scheduledTime)
        if err != nil {
            c.recorder.Eventf(cronJob, corev1.EventTypeWarning, "FailedCreate",
                "Error creating job: %v", err)
            continue
        }

        // 添加到活跃列表
        jobRef, _ := ref.GetReference(scheme.Scheme, job)
        cronJob.Status.Active = append(cronJob.Status.Active, *jobRef)

        // 记录事件
        c.recorder.Eventf(cronJob, corev1.EventTypeNormal, "SuccessfulCreate",
            "Created job %s", job.Name)
    }

    // 6. 更新最后调度时间
    if len(missedSchedules) > 0 {
        cronJob.Status.LastScheduleTime = &metav1.Time{Time: missedSchedules[len(missedSchedules)-1]}
    }

    // 7. 清理历史 Job
    successfulJobs, failedJobs := classifyJobs(inactiveJobs)
    if err := c.historyCleaner.CleanupHistory(ctx, cronJob, successfulJobs, failedJobs); err != nil {
        return err
    }

    // 8. 更新 CronJob 状态
    if _, err := c.kubeClient.BatchV1().CronJobs(cronJob.Namespace).UpdateStatus(
        ctx, cronJob, metav1.UpdateOptions{}); err != nil {
        return err
    }

    return nil
}

func classifyJobs(jobs []*batchv1.Job) (successful, failed []*batchv1.Job) {
    for _, job := range jobs {
        if IsJobSucceeded(job) {
            successful = append(successful, job)
        } else if IsJobFailed(job) {
            failed = append(failed, job)
        }
    }
    return
}
```
