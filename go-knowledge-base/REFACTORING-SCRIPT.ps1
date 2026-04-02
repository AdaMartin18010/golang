# 知识库重构脚本
# 执行文档重命名和映射

$mapping = @{
    # 01-Microservices 系列
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/01-Microservices.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-001-Microservices.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/05-Context-Management.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-005-Context-Management.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/06-Distributed-Tracing.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-006-Distributed-Tracing.md"
    
    # 计划任务系列 (07, 08 已合并)
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/09-Job-Scheduling.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-009-Job-Scheduling.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/10-Async-Task-Queue.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-010-Async-Task-Queue.md"
    
    # 11-99 系列
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/11-Context-Cancellation-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-011-Context-Cancellation-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/12-State-Machine-Workflow.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-012-State-Machine-Workflow.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/13-Concurrent-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-013-Concurrent-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/14-Health-Checks.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-014-Health-Checks.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/15-Resource-Limits.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-015-Resource-Limits.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/16-Service-Discovery.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-016-Service-Discovery.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/17-Scheduled-Task-Framework.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-017-Scheduled-Task-Framework.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/18-Context-Propagation-Framework.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-018-Context-Propagation-Framework.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/19-Task-Execution-Engine.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-019-Task-Execution-Engine.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/20-Distributed-Cron.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-020-Distributed-Cron.md"
}

foreach ($oldPath in $mapping.Keys) {
    $newPath = $mapping[$oldPath]
    if (Test-Path $oldPath) {
        Move-Item -Path $oldPath -Destination $newPath -Force
        Write-Host "Renamed: $oldPath -> $newPath"
    }
}

Write-Host "Refactoring completed!"
