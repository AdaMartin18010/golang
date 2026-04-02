# EC-099: Kubernetes CronJob 深度分析 (Kubernetes CronJob Deep Dive)

> **维度**: Engineering CloudNative
> **级别**: S (20+ KB)
> **标签**: #kubernetes #cronjob #controller #source-analysis
> **相关**: EC-007, EC-008, EC-109

---

## 整合说明

本文档合并了以下文档：

- `59-Kubernetes-CronJob-Controller-Deep-Dive.md` (19 KB)
- `68-Kubernetes-CronJob-V2-Controller.md` (26 KB)
- `114-Task-K8s-CronJob-Controller-Analysis.md` (11 KB)

---

## 架构概览

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        Kubernetes CronJob Controller                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Informer ──► SyncHandler ──► JobControl ──► API Server ──► etcd       │
│      │            │              │                                      │
│      │            │              └── 创建/删除/管理 Jobs                │
│      │            │                                                     │
│      │            └── 处理 CronJob 调度逻辑                             │
│      │                                                                 │
│      └── 监视 CronJob/Job/Pod 变更                                      │
│                                                                          │
│  Key Components:                                                         │
│  - CronJobController: 主控制器循环                                       │
│  - syncOne: 单个 CronJob 同步                                            │
│  - getNextScheduleTime: 计算下次执行时间                                 │
│  - adoptOrphanJobs: 处理孤儿 Job                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## V1 vs V2 控制器对比

| 特性 | V1 (旧) | V2 (默认 v1.21+) |
|------|---------|-----------------|
| 实现 | 单文件 | 模块化 |
| 性能 | 全量列表 | 增量同步 |
| 并发策略 | 简单 | Forbid/Allow/Replace |
| 时区支持 | 无 | 有 (v1.24+) |
| 历史清理 | 固定 | 可配置 |

---

## V2 控制器源码分析

```go
// 基于 Kubernetes v1.28 pkg/controller/cronjob/cronjob_controllerv2.go

// ControllerV2 CronJob控制器 V2版本
type ControllerV2 struct {
 kubeClient clientset.Interface
 recorder   record.EventRecorder
 cjLister   batchv1listers.CronJobLister
 jobLister  batchv1listers.JobLister
 queue      workqueue.RateLimitingInterface
 jobControl jobControlInterface
}

// sync 同步 CronJob
func (jm *ControllerV2) sync(ctx context.Context, cronJobKey string) error {
 ns, name, err := cache.SplitMetaNamespaceKey(cronJobKey)
 if err != nil {
  return err
 }

 cronJob, err := jm.cjLister.CronJobs(ns).Get(name)
 if errors.IsNotFound(err) {
  return nil
 }
 if err != nil {
  return err
 }

 jobs, err := jm.getJobsToBeReconciled(ctx, cronJob)
 if err != nil {
  return err
 }

 return jm.syncCronJob(ctx, cronJob, jobs)
}

// syncCronJob 核心同步逻辑
func (jm *ControllerV2) syncCronJob(ctx context.Context, cronJob *batchv1.CronJob,
 jobs []*batchv1.Job) error {

 // 1. 检查是否被挂起
 if cronJob.Spec.Suspend != nil && *cronJob.Spec.Suspend {
  jm.recorder.Eventf(cronJob, corev1.EventTypeNormal, "Suspended",
   "CronJob is suspended, skipping execution")
  return nil
 }

 // 2. 解析 Cron 表达式
 sched, err := cron.ParseStandard(cronJob.Spec.Schedule)
 if err != nil {
  jm.recorder.Eventf(cronJob, corev1.EventTypeWarning, "InvalidSchedule",
   "Unparseable schedule: %s", err)
  return fmt.Errorf("unparseable schedule: %s", err)
 }

 // 3. 计算下次执行时间
 now := jm.now()
 scheduledTime, missedSchedule, err := getNextScheduleTime(*cronJob, now, sched)
 if err != nil {
  return err
 }

 // 4. 处理错过的调度
 if missedSchedule > 0 {
  jm.recorder.Eventf(cronJob, corev1.EventTypeWarning, "MissedSchedule",
   "Missed %d scheduled times", missedSchedule)

  // 检查截止时间
  if cronJob.Spec.StartingDeadlineSeconds != nil {
   deadline := now.Add(-time.Duration(*cronJob.Spec.StartingDeadlineSeconds) * time.Second)
   if scheduledTime.Before(deadline) {
    jm.recorder.Eventf(cronJob, corev1.EventTypeWarning, "MissedDeadline",
     "Missed scheduled time to start within deadline")
    return nil
   }
  }

  // 检查并发策略
  if cronJob.Spec.ConcurrencyPolicy == batchv1.ForbidConcurrent && len(jobs) > 0 {
   jm.recorder.Eventf(cronJob, corev1.EventTypeNormal, "JobAlreadyActive",
    "Not starting job because prior execution is running and concurrency policy is Forbid")
   return nil
  }

  if cronJob.Spec.ConcurrencyPolicy == batchv1.ReplaceConcurrent {
   for _, j := range jobs {
    if j.Status.Active > 0 {
     jm.deleteJob(ctx, cronJob, j)
    }
   }
  }
 }

 // 5. 创建新 Job
 if scheduledTime != nil {
  jobReq, err := getJobFromTemplate2(cronJob, *scheduledTime)
  if err != nil {
   return err
  }

  jobResp, err := jm.kubeClient.BatchV1().Jobs(cronJob.Namespace).Create(ctx, jobReq, metav1.CreateOptions{})
  if err != nil {
   jm.recorder.Eventf(cronJob, corev1.EventTypeWarning, "FailedCreate",
    "Error creating job: %v", err)
   return err
  }

  jm.recorder.Eventf(cronJob, corev1.EventTypeNormal, "SuccessfulCreate",
   "Created job %v", jobResp.Name)
 }

 // 6. 清理完成的 Jobs
 return jm.cleanupFinishedJobs(ctx, cronJob, jobs)
}

// getNextScheduleTime 计算下次调度时间
func getNextScheduleTime(cj batchv1.CronJob, now time.Time, sched cron.Schedule)
 (*time.Time, int64, error) {

 var earliestTime time.Time
 if cj.Status.LastScheduleTime != nil {
  earliestTime = cj.Status.LastScheduleTime.Time
 } else {
  earliestTime = cj.ObjectMeta.CreationTimestamp.Time
 }

 if cj.Spec.StartingDeadlineSeconds != nil {
  schedulingDeadline := now.Add(-time.Second * time.Duration(*cj.Spec.StartingDeadlineSeconds))
  if schedulingDeadline.After(earliestTime) {
   earliestTime = schedulingDeadline
  }
 }

 if earliestTime.After(now) {
  return nil, 0, nil
 }

 missedSchedules := int64(0)
 mostRecentTime := sched.Next(earliestTime)

 for mostRecentTime.Before(now) {
  missedSchedules++
  if missedSchedules > 100 {
   return nil, 0, fmt.Errorf("too many missed start times (> 100)")
  }
  mostRecentTime = sched.Next(mostRecentTime)
 }

 return &mostRecentTime, missedSchedules, nil
}

// cleanupFinishedJobs 清理完成的 Jobs
func (jm *ControllerV2) cleanupFinishedJobs(ctx context.Context, cj *batchv1.CronJob,
 jobs []*batchv1.Job) error {

 if cj.Spec.SuccessfulJobHistoryLimit == nil && cj.Spec.FailedJobHistoryLimit == nil {
  return nil
 }

 successfulJobs := []*batchv1.Job{}
 failedJobs := []*batchv1.Job{}

 for _, job := range jobs {
  if isJobFinished(job) {
   if job.Status.Succeeded > 0 {
    successfulJobs = append(successfulJobs, job)
   } else {
    failedJobs = append(failedJobs, job)
   }
  }
 }

 successfulLimit := int32(3)
 if cj.Spec.SuccessfulJobHistoryLimit != nil {
  successfulLimit = *cj.Spec.SuccessfulJobHistoryLimit
 }

 failedLimit := int32(1)
 if cj.Spec.FailedJobHistoryLimit != nil {
  failedLimit = *cj.Spec.FailedJobHistoryLimit
 }

 if int32(len(successfulJobs)) > successfulLimit {
  sort.Sort(byJobStartTime(successfulJobs))
  for i := int32(0); i < int32(len(successfulJobs))-successfulLimit; i++ {
   jm.deleteJob(ctx, cj, successfulJobs[i])
  }
 }

 if int32(len(failedJobs)) > failedLimit {
  sort.Sort(byJobStartTime(failedJobs))
  for i := int32(0); i < int32(len(failedJobs))-failedLimit; i++ {
   jm.deleteJob(ctx, cj, failedJobs[i])
  }
 }

 return nil
}
```

---

## 时区处理（v1.24+）

```go
func getTimezone(cj *batchv1.CronJob) (*time.Location, error) {
 if cj.Spec.TimeZone == nil {
  return time.UTC, nil
 }

 tz := *cj.Spec.TimeZone
 if tz == "Local" {
  return nil, fmt.Errorf("Local timezone not supported")
 }

 loc, err := time.LoadLocation(tz)
 if err != nil {
  return nil, fmt.Errorf("invalid timezone: %s", tz)
 }

 return loc, nil
}
```

---

## 关键设计决策

| 设计点 | 决策 | 原因 |
|--------|------|------|
| 错过调度处理 | 统计错过的次数（上限100） | 防止惊群效应 |
| 并发策略 | Forbid/Allow/Replace | 满足不同场景需求 |
| 历史限制 | 成功/失败分别限制 | 避免存储无限增长 |
| 时区支持 | 显式声明（禁止Local） | 避免分布式时区问题 |
| 截止时间 | startingDeadlineSeconds | 跳过过期的调度 |

---

## 性能优化

```go
// 指数退避重试
jm.queue = workqueue.NewRateLimitingQueue(
 workqueue.NewItemExponentialFailureRateLimiter(5*time.Second, 5*time.Minute),
)

// 事件过滤（避免不必要的同步）
func (jm *ControllerV2) updateCronJob(old, cur interface{}) {
 oldCJ := old.(*batchv1.CronJob)
 curCJ := cur.(*batchv1.CronJob)

 if !reflect.DeepEqual(oldCJ.Spec, curCJ.Spec) {
  jm.enqueueController(curCJ)
 }
}
```
