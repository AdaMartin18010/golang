# 批量重命名 Go 1.25.3 文档到 Go 1.26
# 策略: 创建新的1.26版本，移动旧版本到archive

$files = @(
    # advanced
    @{Old="docs/advanced/00-Go-1.25.3新兴技术应用-2025.md"; New="docs/advanced/00-Go-1.26新兴技术应用-2026.md"},
    @{Old="docs/advanced/20-Go-1.25.3消息队列与异步处理完整实战.md"; New="docs/advanced/20-Go-1.26消息队列与异步处理完整实战.md"},
    @{Old="docs/advanced/21-Go-1.25.3分布式缓存架构完整实战.md"; New="docs/advanced/21-Go-1.26分布式缓存架构完整实战.md"},
    @{Old="docs/advanced/22-Go-1.25.3安全加固与认证授权完整实战.md"; New="docs/advanced/22-Go-1.26安全加固与认证授权完整实战.md"},
    @{Old="docs/advanced/23-Go-1.25.3分布式追踪与可观测性完整实战.md"; New="docs/advanced/23-Go-1.26分布式追踪与可观测性完整实战.md"},
    @{Old="docs/advanced/24-Go-1.25.3流量控制与限流完整实战.md"; New="docs/advanced/24-Go-1.26流量控制与限流完整实战.md"},
    @{Old="docs/advanced/25-Go-1.25.3API网关完整实战.md"; New="docs/advanced/25-Go-1.26API网关完整实战.md"},
    @{Old="docs/advanced/26-Go-1.25.3分布式事务完整实战.md"; New="docs/advanced/26-Go-1.26分布式事务完整实战.md"},
    @{Old="docs/advanced/27-Go-1.25.3配置中心与服务治理完整实战.md"; New="docs/advanced/27-Go-1.26配置中心与服务治理完整实战.md"},
    @{Old="docs/advanced/28-Go-1.25.3服务网格与高级流量治理完整实战.md"; New="docs/advanced/28-Go-1.26服务网格与高级流量治理完整实战.md"},
    @{Old="docs/advanced/29-Go-1.25.3事件溯源与CQRS完整实战.md"; New="docs/advanced/29-Go-1.26事件溯源与CQRS完整实战.md"},
    @{Old="docs/advanced/30-Go-1.25.3实时数据处理与流计算完整实战.md"; New="docs/advanced/30-Go-1.26实时数据处理与流计算完整实战.md"},
    @{Old="docs/advanced/31-Go-1.25.3GraphQL现代API完整实战.md"; New="docs/advanced/31-Go-1.26GraphQL现代API完整实战.md"},
    @{Old="docs/advanced/32-Go-1.25.3Serverless与FaaS完整实战.md"; New="docs/advanced/32-Go-1.26Serverless与FaaS完整实战.md"},
    @{Old="docs/advanced/33-Go-1.25.3AI与机器学习集成完整实战.md"; New="docs/advanced/33-Go-1.26AI与机器学习集成完整实战.md"},
    @{Old="docs/advanced/34-Go-1.25.3WebAssembly完整实战.md"; New="docs/advanced/34-Go-1.26WebAssembly完整实战.md"},
    @{Old="docs/advanced/35-Go-1.25.3高级DevOps完整实战.md"; New="docs/advanced/35-Go-1.26高级DevOps完整实战.md"},
    @{Old="docs/advanced/36-Go-1.25.3边缘计算与IoT完整实战.md"; New="docs/advanced/36-Go-1.26边缘计算与IoT完整实战.md"},
    @{Old="docs/advanced/37-Go-1.25.3区块链与Web3完整实战.md"; New="docs/advanced/37-Go-1.26区块链与Web3完整实战.md"},
    @{Old="docs/advanced/38-Go-1.25.3搜索引擎与全文检索完整实战.md"; New="docs/advanced/38-Go-1.26搜索引擎与全文检索完整实战.md"},
    @{Old="docs/advanced/architecture/00-Go-1.25.3编程设计模式与最佳实践-2025.md"; New="docs/advanced/architecture/00-Go-1.26编程设计模式与最佳实践-2026.md"},
    @{Old="docs/advanced/performance/08-Go-1.25.3性能优化完整实战.md"; New="docs/advanced/performance/08-Go-1.26性能优化完整实战.md"},
    # development
    @{Old="docs/development/cloud-native/00-Go-1.25.3云原生生态全景-2025.md"; New="docs/development/cloud-native/00-Go-1.26云原生生态全景-2026.md"},
    @{Old="docs/development/cloud-native/10-Go-1.25.3云原生部署完整实战.md"; New="docs/development/cloud-native/10-Go-1.26云原生部署完整实战.md"},
    @{Old="docs/development/database/04-Go-1.25.3数据库编程完整实战.md"; New="docs/development/database/04-Go-1.26数据库编程完整实战.md"},
    @{Old="docs/development/microservices/09-Go 1.25.1微服务优化.md"; New="docs/development/microservices/09-Go-1.26微服务优化.md"},
    @{Old="docs/development/microservices/10-Go-1.25.3微服务架构完整实战.md"; New="docs/development/microservices/10-Go-1.26微服务架构完整实战.md"},
    @{Old="docs/development/web/23-Go-1.25.3现代Web服务完整项目.md"; New="docs/development/web/23-Go-1.26现代Web服务完整项目.md"},
    # fundamentals
    @{Old="docs/fundamentals/data-structures/05-Go-1.25.3泛型数据结构实战.md"; New="docs/fundamentals/data-structures/05-Go-1.26泛型数据结构实战.md"},
    @{Old="docs/fundamentals/language/00-Go-1.25.3语言语义模型更新.md"; New="docs/fundamentals/language/00-Go-1.26语言语义模型更新.md"},
    @{Old="docs/fundamentals/language/02-并发编程/08-Go-1.25.3并发编程完整实战.md"; New="docs/fundamentals/language/02-并发编程/08-Go-1.26并发编程完整实战.md"},
    @{Old="docs/fundamentals/language/03-模块管理/07-Go-Workspace完整指南-Go1.25.3.md"; New="docs/fundamentals/language/03-模块管理/07-Go-Workspace完整指南-Go1.26.md"},
    # practices
    @{Old="docs/practices/engineering/05-Go-1.25.3测试工程完整实战.md"; New="docs/practices/engineering/05-Go-1.26测试工程完整实战.md"}
)

$archiveDir = "docs/archive/go125-versions"

Write-Host "=== 批量处理文档重命名 ===" -ForegroundColor Cyan
Write-Host "总计: $($files.Count) 个文件" -ForegroundColor Yellow
Write-Host ""

$success = 0
$failed = 0

foreach ($file in $files) {
    $oldPath = $file.Old
    $newPath = $file.New
    
    if (Test-Path $oldPath) {
        try {
            # 1. 读取旧文件内容
            $content = Get-Content $oldPath -Raw -Encoding UTF8
            
            # 2. 更新内容中的版本号
            $newContent = $content -replace "Go 1\.25\.3", "Go 1.26" -replace "Go1\.25\.3", "Go1.26" -replace "适用于: Go 1\.26", "适用于: Go 1.26" -replace "2025-11-11", "2026-03-07"
            
            # 3. 写入新文件
            $newContent | Set-Content $newPath -Encoding UTF8 -NoNewline
            
            # 4. 移动旧文件到archive
            $fileName = Split-Path $oldPath -Leaf
            $archivePath = Join-Path $archiveDir $fileName
            Move-Item $oldPath $archivePath -Force
            
            Write-Host "✓ $fileName" -ForegroundColor Green
            $success++
        } catch {
            Write-Host "✗ $fileName - $_" -ForegroundColor Red
            $failed++
        }
    } else {
        Write-Host "⚠ 不存在: $oldPath" -ForegroundColor Yellow
    }
}

# 处理子目录中的文件
$subDirs = @(
    "docs/fundamentals/language/00-Go-1.25.3核心机制完整解析"
)

foreach ($dir in $subDirs) {
    if (Test-Path $dir) {
        $parentDir = Split-Path $dir -Parent
        $newDirName = (Split-Path $dir -Leaf) -replace "1\.25\.3", "1.26"
        $newDir = Join-Path $parentDir $newDirName
        
        # 重命名目录
        Rename-Item $dir $newDir -Force
        Write-Host "✓ 目录重命名: $(Split-Path $dir -Leaf)" -ForegroundColor Green
        
        # 重命名目录内的文件
        Get-ChildItem $newDir -File | Where-Object { $_.Name -match "1\.25\.3" } | ForEach-Object {
            $newFileName = $_.Name -replace "1\.25\.3", "1.26"
            Rename-Item $_.FullName (Join-Path $newDir $newFileName) -Force
            Write-Host "  ✓ $($_.Name)" -ForegroundColor Green
        }
    }
}

Write-Host ""
Write-Host "完成: $success 成功, $failed 失败" -ForegroundColor Cyan
