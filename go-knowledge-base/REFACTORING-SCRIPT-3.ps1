# 第三阶段重构脚本 - 处理 51-100 系列

$mapping = @{
    # 51-99 系列 (跳过已合并的 59, 68 -> EC-099; 58, 69 -> EC-100)
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/51-Task-Context-Propagation-Advanced.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-051-Task-Context-Propagation-Advanced.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/52-Task-Context-Cancellation-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-052-Task-Context-Cancellation-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/53-Task-Context-Value-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-053-Task-Context-Value-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/54-Task-Context-Propagation-Standards.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-054-Task-Context-Propagation-Standards.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/55-Task-Context-Propagation-Best-Practices.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-055-Task-Context-Propagation-Best-Practices.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/56-Task-Distributed-Tracing-Deep-Dive.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-056-Task-Distributed-Tracing-Deep-Dive.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/57-ETCD-Distributed-Task-Scheduler.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-057-ETCD-Distributed-Task-Scheduler.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/60-OpenTelemetry-Distributed-Tracing-Production.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-060-OpenTelemetry-Distributed-Tracing-Production.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/61-Task-Queue-Implementation-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-061-Task-Queue-Implementation-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/62-Distributed-Task-Scheduler-Architecture.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-062-Distributed-Task-Scheduler-Architecture.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/63-Task-State-Machine-Implementation.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-063-Task-State-Machine-Implementation.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/64-Context-Management-Production-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-064-Context-Management-Production-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/65-Database-Transaction-Isolation-MVCC.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-065-Database-Transaction-Isolation-MVCC.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/66-Context-Propagation-Implementation.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-066-Context-Propagation-Implementation.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/67-Distributed-Task-Scheduler-Production.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-067-Distributed-Task-Scheduler-Production.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/70-OpenTelemetry-W3C-Trace-Context.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-070-OpenTelemetry-W3C-Trace-Context.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/71-etcd-Distributed-Coordination.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-071-etcd-Distributed-Coordination.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/72-Task-Queue-Implementation.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-072-Task-Queue-Implementation.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/73-Worker-Pool-Dynamic-Scaling.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-073-Worker-Pool-Dynamic-Scaling.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/74-Context-Aware-Logging.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-074-Context-Aware-Logging.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/75-Retry-Backoff-Circuit-Breaker.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-075-Retry-Backoff-Circuit-Breaker.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/76-DAG-Task-Dependencies.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-076-DAG-Task-Dependencies.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/77-State-Machine-Task-Execution.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-077-State-Machine-Task-Execution.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/78-Rate-Limiting-Throttling.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-078-Rate-Limiting-Throttling.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/79-Graceful-Shutdown-Implementation.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-079-Graceful-Shutdown-Implementation.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/80-Observability-Metrics-Integration.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-080-Observability-Metrics-Integration.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/81-Task-Execution-Lifecycle-Management.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-081-Task-Execution-Lifecycle-Management.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/82-Distributed-Task-Sharding.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-082-Distributed-Task-Sharding.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/83-Task-Execution-Timeout-Control.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-083-Task-Execution-Timeout-Control.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/84-Cancellation-Propagation-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-084-Cancellation-Propagation-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/85-Resource-Management-Scheduling.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-085-Resource-Management-Scheduling.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/86-Health-Check-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-086-Health-Check-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/87-Async-Task-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-087-Async-Task-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/88-Delayed-Task-Scheduling.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-088-Delayed-Task-Scheduling.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/89-Task-Priority-Queue.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-089-Task-Priority-Queue.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/90-Task-Compensation-Saga-Pattern.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-090-Task-Compensation-Saga-Pattern.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/91-Distributed-Lock-Implementation.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-091-Distributed-Lock-Implementation.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/92-Task-Event-Sourcing-Persistence.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-092-Task-Event-Sourcing-Persistence.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/93-Multi-Tenancy-Task-Isolation.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-093-Multi-Tenancy-Task-Isolation.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/94-Task-Debugging-Diagnostics.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-094-Task-Debugging-Diagnostics.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/95-Task-Testing-Strategies.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-095-Task-Testing-Strategies.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/96-Task-Deployment-Operations.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-096-Task-Deployment-Operations.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/97-Task-CLI-Tooling.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-097-Task-CLI-Tooling.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/99-Task-System-Architecture-Overview.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-099-Task-System-Architecture-Overview.md"
}

foreach ($oldPath in $mapping.Keys) {
    $newPath = $mapping[$oldPath]
    if (Test-Path $oldPath) {
        Move-Item -Path $oldPath -Destination $newPath -Force
        Write-Host "Renamed: $(Split-Path $oldPath -Leaf)"
    }
}

Write-Host "Phase 3 refactoring completed!"
