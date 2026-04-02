# Kubernetes CronJob Controller 源码深度解析 (Kubernetes CronJob Controller Deep Dive)

> **分类**: 工程与云原生
> **标签**: #kubernetes #cronjob #controller #source-code
> **参考**: k8s.io/kubernetes/pkg/controller/cronjob, Kubernetes 1.28+

---

## 架构概述

Kubernetes CronJob Controller 是一个控制平面组件，负责根据 Cron 表达式调度 Job 的创建。

```
┌─────────────────────────────────────────────────────────────────┐
│                     CronJob Controller                          │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Controller V2                        │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │    │
│  │  │   Informer   │  │  DelayingQueue│  │   Recorder   │  │    │
│  │  │ (CronJob/Job)│  │ (Next Schedule)│  │   (Events) │   │    │
│  │  └──────┬───────┘  └──────┬───────┘  └──────────────┘   │    │
│  │         │                 │                             │    │
│  │         └────────┬────────┘                             │    │
│  │                  │                                      │    │
│  │         ┌────────▼────────┐                             │    │
│  │         │   syncHandler   │                             │    │
│  │         │  • syncCronJob  │                             │    │
│  │         │  • enqueueCronJob│                            │    │
│  │         └─────────────────┘                             │    │
│  └─────────────────────────────────────────────────────────┘    │
│                              │                                  │
│                              ▼                                  │
│                    ┌─────────────────┐                          │
│                    │   kube-apiserver │                         │
│                    │  (Watch/Update)  │                         │
│                    └─────────────────┘                          │
└─────────────────────────────────────────────────────────────────┘
```

---

## Controller V2 实现

```go
// ControllerV2 CronJob 控制器 V2 版本实现
// 路径: k8s.io/kubernetes/pkg/controller/cronjob/cronjob_controllerv2.go

type ControllerV2 struct {
    // KubeClient 用于操作 API
    kubeClient clientset.Interface

    // Informers 用于监听资源变化
    cronJobLister batchv1listers.CronJobLister
    cronJobSynced cache.InformerSynced
    jobLister     batchv1listers.JobLister
    jobSynced     cache.InformerSynced

    // 工作队列
    queue workqueue.RateLimitingInterface

    // 延迟队列，用于调度下一次执行
    delayQueue workqueue.DelayingInterface

    // 记录事件
    recorder record.EventRecorder

    // 时钟，用于测试
    now func() time.Time
}

// NewControllerV2 创建控制器实例
func NewControllerV2(
    kubeClient clientset.Interface,
    cronJobInformer batchv1informers.CronJobInformer,
    jobInformer batchv1informers.JobInformer,
) *ControllerV2 {

    eventBroadcaster := record.NewBroadcaster()

    jm := &ControllerV2{
        kubeClient: kubeClient,
        queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "cronjob"),
        delayQueue: workqueue.NewNamedDelayingQueue("cronjob-delay"),
        recorder:   eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "cronjob-controller"}),
        now:        time.Now,
    }

    // 设置 CronJob Informer 回调
    cronJobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc:    jm.addCronJob,
        UpdateFunc: jm.updateCronJob,
        DeleteFunc: jm.deleteCronJob,
    })

    // 设置 Job Informer 回调
    jobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc:    jm.addJob,
        UpdateFunc: jm.updateJob,
        DeleteFunc: jm.deleteJob,
    })

    jm.cronJobLister = cronJobInformer.Lister()
    jm.cronJobSynced = cronJobInformer.Informer().HasSynced
    jm.jobLister = jobInformer.Lister()
    jm.jobSynced = jobInformer.Informer().HasSynced

    return jm
}

// Run 启动控制器
func (jm *ControllerV2) Run(workers int, stopCh <-chan struct{}) {
    defer utilruntime.HandleCrash()
    defer jm.queue.ShutDown()

    klog.InfoS("Starting CronJob Controller")

    // 等待缓存同步
    if !cache.WaitForNamedCacheSync("cronjob", stopCh, jm.cronJobSynced, jm.jobSynced) {
        return
    }

    // 启动工作线程
    for i := 0; i < workers; i++ {
        go wait.Until(jm.worker, time.Second, stopCh)
    }

    // 启动延迟队列处理
    go wait.Until(jm.delayWorker, time.Second, stopCh)

    <-stopCh
    klog.InfoS("Shutting down CronJob Controller")
}

// worker 处理队列中的 CronJob
func (jm *ControllerV2) worker() {
    for jm.processNextWorkItem() {
    }
}

func (jm *ControllerV2) processNextWorkItem() bool {
    key, quit := jm.queue.Get()
    if quit {
        return false
    }
    defer jm.queue.Done(key)

    err := jm.syncCronJob(key.(string))
    if err == nil {
        jm.queue.Forget(key)
        return true
    }

    // 失败重试
    utilruntime.HandleError(fmt.Errorf("sync cronjob failed: %v", err))
    jm.queue.AddRateLimited(key)
    return true
}
```

---

## 核心调度逻辑

```go
// syncCronJob 同步 CronJob 状态
func (jm *ControllerV2) syncCronJob(key string) error {
    ns, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        return err
    }

    // 获取 CronJob
    cronJob, err := jm.cronJobLister.CronJobs(ns).Get(name)
    if errors.IsNotFound(err) {
        return nil
    }
    if err != nil {
        return err
    }

    // 检查是否被暂停
    if cronJob.Spec.Suspend != nil && *cronJob.Spec.Suspend {
        klog.V(4).InfoS("CronJob suspended, skipping", "cronjob", klog.KObj(cronJob))
        return nil
    }

    // 获取关联的 Jobs
    jobs, err := jm.getJobsToReconcile(cronJob)
    if err != nil {
        return err
    }

    // 清理已完成的 Jobs
    if err := jm.cleanupFinishedJobs(cronJob, jobs); err != nil {
        return err
    }

    // 计算下一次调度时间
    nextSchedule, missedSchedules, err := getNextScheduleTime(cronJob, jm.now(), jm.recorder)
    if err != nil {
        return err
    }

    // 处理错过的调度
    for _, missedSchedule := range missedSchedules {
        // 检查 StartingDeadlineSeconds
        if cronJob.Spec.StartingDeadlineSeconds != nil {
            deadline := missedSchedule.Add(time.Duration(*cronJob.Spec.StartingDeadlineSeconds) * time.Second)
            if deadline.Before(jm.now()) {
                // 已过截止时间，记录警告
                jm.recorder.Eventf(cronJob, corev1.EventTypeWarning, "MissSchedule",
                    "Missed scheduled time to start a job: %s", missedSchedule.Format(time.RFC1123))
                continue
            }
        }

        // 检查并发策略
        if cronJob.Spec.ConcurrencyPolicy == batchv1.ForbidConcurrent && len(jobs.active) > 0 {
            jm.recorder.Eventf(cronJob, corev1.EventTypeWarning, "JobAlreadyActive",
                "Not starting job because prior execution is running and concurrency policy is Forbid")
            continue
        }

        if cronJob.Spec.ConcurrencyPolicy == batchv1.ReplaceConcurrent && len(jobs.active) > 0 {
            // 取消正在运行的 Job
            for _, job := range jobs.active {
                if err := jm.kubeClient.BatchV1().Jobs(job.Namespace).Delete(context.TODO(), job.Name, metav1.DeleteOptions{}); err != nil {
                    return err
                }
            }
        }

        // 创建新的 Job
        job, err := jm.createJob(cronJob, missedSchedule)
        if err != nil {
            return err
        }

        jm.recorder.Eventf(cronJob, corev1.EventTypeNormal, "SuccessfulCreate",
            "Created job %v", job.Name)
    }

    // 更新 CronJob 状态
    if err := jm.updateCronJobStatus(cronJob, nextSchedule); err != nil {
        return err
    }

    // 安排下一次调度
    if nextSchedule != nil {
        duration := nextSchedule.Sub(jm.now())
        if duration < 0 {
            duration = 0
        }
        jm.delayQueue.AddAfter(key, duration)
    }

    return nil
}

// createJob 创建 Job 对象
func (jm *ControllerV2) createJob(cronJob *batchv1.CronJob, scheduledTime time.Time) (*batchv1.Job, error) {
    // 生成 Job 名称
    scheduledTimeUnix := scheduledTime.Unix()
    name := fmt.Sprintf("%s-%d", cronJob.Name, scheduledTimeUnix)

    // 构建 Job 对象
    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: cronJob.Namespace,
            Labels:    cronJob.Spec.JobTemplate.Labels,
            Annotations: map[string]string{
                "cronjob.kubernetes.io/instantiate": "manual",
            },
            OwnerReferences: []metav1.OwnerReference{
                *metav1.NewControllerRef(cronJob, batchv1.SchemeGroupVersion.WithKind("CronJob")),
            },
        },
        Spec: cronJob.Spec.JobTemplate.Spec,
    }

    // 设置标签
    if job.Labels == nil {
        job.Labels = make(map[string]string)
    }
    job.Labels["cronjob-name"] = cronJob.Name
    job.Labels["scheduled-timestamp"] = strconv.FormatInt(scheduledTimeUnix, 10)

    // 创建 Job
    return jm.kubeClient.BatchV1().Jobs(cronJob.Namespace).Create(context.TODO(), job, metav1.CreateOptions{})
}
```

---

## Cron 表达式解析

```go
// getNextScheduleTime 计算下一次调度时间
func getNextScheduleTime(
    cronJob *batchv1.CronJob,
    now time.Time,
    recorder record.EventRecorderLogger,
) (*time.Time, []time.Time, error) {

    // 解析 Cron 表达式
    sched, err := cron.ParseStandard(cronJob.Spec.Schedule)
    if err != nil {
        recorder.Eventf(cronJob, corev1.EventTypeWarning, "InvalidSchedule",
            "Unparseable schedule: %s : %s", cronJob.Spec.Schedule, err)
        return nil, nil, fmt.Errorf("unparseable schedule: %s", err)
    }

    var (
        missedSchedules []time.Time
        nextSchedule    time.Time
    )

    // 获取上次调度时间
    lastScheduleTime := cronJob.Status.LastScheduleTime
    if lastScheduleTime == nil {
        // 首次调度，使用 CreationTimestamp
        lastScheduleTime = &cronJob.CreationTimestamp
    }

    // 计算错过的调度
    for t := sched.Next(lastScheduleTime.Time); t.Before(now); t = sched.Next(t) {
        missedSchedules = append(missedSchedules, t)

        // 限制错过的调度数量，防止创建大量 Job
        if len(missedSchedules) > 100 {
            recorder.Eventf(cronJob, corev1.EventTypeWarning, "TooManyMissedStarts",
                "Too many missed start times (> 100). Setting or lowering .spec.startingDeadlineSeconds or .spec.concurrencyPolicy are recommended.")
            break
        }
    }

    // 计算下一次调度
    next := sched.Next(now)
    nextSchedule = &next

    return nextSchedule, missedSchedules, nil
}

// Cron 表达式特殊语法支持
const (
    // @yearly (or @annually)    Run once a year at midnight of 1 January    0 0 1 1 *
    // @monthly                  Run once a month at midnight of the first day    0 0 1 * *
    // @weekly                   Run once a week at midnight on Sunday    0 0 * * 0
    // @daily (or @midnight)     Run once a day at midnight    0 0 * * *
    // @hourly                   Run once an hour at the beginning of the hour    0 * * * *
    // @every <duration>         Run every duration
)

// 使用 robfig/cron/v3 库解析
func parseCronSchedule(schedule string) (cron.Schedule, error) {
    parser := cron.NewParser(
        cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
    )
    return parser.Parse(schedule)
}
```

---

## Job 清理策略

```go
// cleanupFinishedJobs 清理已完成的 Job
func (jm *ControllerV2) cleanupFinishedJobs(cronJob *batchv1.CronJob, jobs *jobReconciliation) error {
    // 确定保留数量
    successfulJobsHistoryLimit := int32(3)
    if cronJob.Spec.SuccessfulJobsHistoryLimit != nil {
        successfulJobsHistoryLimit = *cronJob.Spec.SuccessfulJobsHistoryLimit
    }

    failedJobsHistoryLimit := int32(1)
    if cronJob.Spec.FailedJobsHistoryLimit != nil {
        failedJobsHistoryLimit = *cronJob.Spec.FailedJobsHistoryLimit
    }

    // 按完成时间排序
    sortJobsByCompletionTime(jobs.successful)
    sortJobsByCompletionTime(jobs.failed)

    // 删除超出的成功 Job
    for i := int(successfulJobsHistoryLimit); i < len(jobs.successful); i++ {
        job := jobs.successful[i]
        if err := jm.kubeClient.BatchV1().Jobs(job.Namespace).Delete(
            context.TODO(), job.Name, metav1.DeleteOptions{}); err != nil {
            return err
        }
        jm.recorder.Eventf(cronJob, corev1.EventTypeNormal, "SuccessfulDelete",
            "Deleted job %v", job.Name)
    }

    // 删除超出的失败 Job
    for i := int(failedJobsHistoryLimit); i < len(jobs.failed); i++ {
        job := jobs.failed[i]
        if err := jm.kubeClient.BatchV1().Jobs(job.Namespace).Delete(
            context.TODO(), job.Name, metav1.DeleteOptions{}); err != nil {
            return err
        }
        jm.recorder.Eventf(cronJob, corev1.EventTypeNormal, "SuccessfulDelete",
            "Deleted job %v", job.Name)
    }

    return nil
}

// jobReconciliation Job 分类
type jobReconciliation struct {
    active    []*batchv1.Job
    successful []*batchv1.Job
    failed    []*batchv1.Job
}

// getJobsToReconcile 获取需要协调的 Jobs
func (jm *ControllerV2) getJobsToReconcile(cronJob *batchv1.CronJob) (*jobReconciliation, error) {
    // 获取所有 Job
    jobList, err := jm.jobLister.Jobs(cronJob.Namespace).List(labels.SelectorFromSet(map[string]string{
        "cronjob-name": cronJob.Name,
    }))
    if err != nil {
        return nil, err
    }

    result := &jobReconciliation{}

    for _, job := range jobList {
        // 检查 OwnerReference
        if !metav1.IsControlledBy(job, cronJob) {
            continue
        }

        // 分类 Job
        if isJobFinished(job) {
            if job.Status.Succeeded > 0 {
                result.successful = append(result.successful, job)
            } else {
                result.failed = append(result.failed, job)
            }
        } else {
            result.active = append(result.active, job)
        }
    }

    return result, nil
}

func isJobFinished(job *batchv1.Job) bool {
    for _, c := range job.Status.Conditions {
        if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == corev1.ConditionTrue {
            return true
        }
    }
    return false
}
```

---

## 时区支持 (Kubernetes 1.27+)

```go
// Kubernetes 1.27+ 引入的时区支持

type CronJobSpec struct {
    // ... 其他字段

    // TimeZone 指定调度时使用的时区
    // 例如 "Asia/Shanghai", "America/New_York", "UTC"
    // +optional
    TimeZone *string `json:"timeZone,omitempty" protobuf:"bytes,9,opt,name=timeZone"`
}

// withTimeZone 应用时区到调度计算
func withTimeZone(cronJob *batchv1.CronJob, now time.Time) (*time.Location, time.Time, error) {
    if cronJob.Spec.TimeZone == nil {
        return time.UTC, now.UTC(), nil
    }

    loc, err := time.LoadLocation(*cronJob.Spec.TimeZone)
    if err != nil {
        return nil, time.Time{}, fmt.Errorf("invalid timeZone %s: %w", *cronJob.Spec.TimeZone, err)
    }

    return loc, now.In(loc), nil
}

// 使用时区的调度计算
func getNextScheduleTimeWithTZ(cronJob *batchv1.CronJob, now time.Time) (*time.Time, []time.Time, error) {
    loc, nowInTZ, err := withTimeZone(cronJob, now)
    if err != nil {
        return nil, nil, err
    }

    // 在指定时区解析 Cron 表达式
    locationAwareParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
    sched, err := locationAwareParser.Parse(cronJob.Spec.Schedule)
    if err != nil {
        return nil, nil, err
    }

    // 计算调度时间
    next := sched.Next(nowInTZ).In(loc)

    // ... 其他逻辑

    return &next, missedSchedules, nil
}
```

---

## 控制器性能优化

```go
// ControllerV2 的性能优化策略

// 1. 使用 DelayingQueue 代替轮询
// 传统的 V1 使用 10 秒轮询，V2 使用 DelayingQueue 实现事件驱动

// 2. 批量事件处理
func (jm *ControllerV2) addCronJob(obj interface{}) {
    key, err := cache.MetaNamespaceKeyFunc(obj)
    if err != nil {
        utilruntime.HandleError(err)
        return
    }
    jm.queue.Add(key)
}

// 3. 去重更新
func (jm *ControllerV2) updateCronJob(oldObj, newObj interface{}) {
    oldCJ := oldObj.(*batchv1.CronJob)
    newCJ := newObj.(*batchv1.CronJob)

    // 只有 Spec 变化或删除时才重新调度
    if oldCJ.Spec.Schedule != newCJ.Spec.Schedule ||
        !equality.Semantic.DeepEqual(oldCJ.Spec.JobTemplate, newCJ.Spec.JobTemplate) ||
        newCJ.DeletionTimestamp != nil {

        key, err := cache.MetaNamespaceKeyFunc(newObj)
        if err != nil {
            utilruntime.HandleError(err)
            return
        }
        jm.queue.Add(key)
    }
}

// 4. 限速队列防止雪崩
func newRateLimitingQueue() workqueue.RateLimitingInterface {
    return workqueue.NewNamedRateLimitingQueue(
        workqueue.NewMaxOfRateLimiter(
            workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
            &workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
        ),
        "cronjob",
    )
}

// 5. 内存优化：使用 Lister 而不是直接访问 API Server
// jm.cronJobLister.CronJobs(ns).Get(name) // 读取本地缓存
```
