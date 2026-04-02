# 第二阶段重构脚本 - 处理 21-50 和 51-100 系列

$mapping = @{
    # 21-50 系列
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/21-Task-Queue-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-021-Task-Queue-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/23-Task-Dependency-Management.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-023-Task-Dependency-Management.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/24-Task-State-Machine.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-024-Task-State-Machine.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/25-Task-Compensation.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-025-Task-Compensation.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/27-Task-Versioning.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-027-Task-Versioning.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/28-Task-Data-Consistency.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-028-Task-Data-Consistency.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/29-Task-Failure-Recovery.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-029-Task-Failure-Recovery.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/30-Task-Rate-Limiting.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-030-Task-Rate-Limiting.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/31-Task-Scheduling-Strategies.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-031-Task-Scheduling-Strategies.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/33-Task-Batch-Processing.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-033-Task-Batch-Processing.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/34-Task-Event-Sourcing.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-034-Task-Event-Sourcing.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/35-Task-Multi-Tenancy.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-035-Task-Multi-Tenancy.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/36-Task-Debugging-Diagnostics.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-036-Task-Debugging-Diagnostics.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/37-Task-Testing-Strategies.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-037-Task-Testing-Strategies.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/38-Task-Documentation-Generator.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-038-Task-Documentation-Generator.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/39-Task-Migration-Guide.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-039-Task-Migration-Guide.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/40-Task-Configuration-Management.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-040-Task-Configuration-Management.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/41-Task-CLI-Tooling.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-041-Task-CLI-Tooling.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/43-Task-API-Design.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-043-Task-API-Design.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/44-Task-Schema-Registry.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-044-Task-Schema-Registry.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/45-Task-Security-Hardening.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-045-Task-Security-Hardening.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/46-Task-Performance-Tuning.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-046-Task-Performance-Tuning.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/47-Task-Deployment-Operations.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-047-Task-Deployment-Operations.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/48-Task-Case-Studies.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-048-Task-Case-Studies.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/49-Task-Integration-Patterns.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-049-Task-Integration-Patterns.md"
    "go-knowledge-base/03-Engineering-CloudNative/02-Cloud-Native/50-Task-Future-Trends.md" = "go-knowledge-base/03-Engineering-CloudNative/EC-050-Task-Future-Trends.md"
}

foreach ($oldPath in $mapping.Keys) {
    $newPath = $mapping[$oldPath]
    if (Test-Path $oldPath) {
        Move-Item -Path $oldPath -Destination $newPath -Force
        Write-Host "Renamed: $(Split-Path $oldPath -Leaf) -> $(Split-Path $newPath -Leaf)"
    }
}

Write-Host "Phase 2 refactoring completed!"
