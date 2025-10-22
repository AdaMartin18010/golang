# PowerShell Script: 创建新的文档目录结构
# 用途: 自动创建重构后的13个主题模块目录
# 版本: v1.0
# 日期: 2025-10-22

param(
    [string]$BaseDir = "docs-new",
    [switch]$DryRun
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Golang 文档结构创建脚本" -ForegroundColor Cyan
Write-Host "  版本: v1.0" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

if ($DryRun) {
    Write-Host "[DRY RUN 模式] 仅显示将要创建的目录" -ForegroundColor Yellow
    Write-Host ""
}

# 定义新结构
$structure = @(
    # 01-语言基础
    "01-语言基础/01-语法基础",
    "01-语言基础/02-并发编程",
    "01-语言基础/03-模块管理",
    "01-语言基础/04-标准库精讲",
    
    # 02-数据结构与算法
    "02-数据结构与算法/01-基础数据结构",
    "02-数据结构与算法/02-算法实现",
    "02-数据结构与算法/03-Go特有实现",
    
    # 03-Web开发
    "03-Web开发/01-HTTP编程",
    "03-Web开发/02-Web框架",
    "03-Web开发/03-Web进阶",
    "03-Web开发/04-RESTful-API",
    "03-Web开发/05-认证与授权",
    "03-Web开发/06-安全与性能",
    
    # 04-数据库编程
    "04-数据库编程/01-SQL数据库",
    "04-数据库编程/02-NoSQL数据库",
    "04-数据库编程/03-ORM框架",
    "04-数据库编程/04-数据库高级",
    
    # 05-微服务架构
    "05-微服务架构/01-微服务基础",
    "05-微服务架构/02-RPC框架",
    "05-微服务架构/03-服务发现",
    "05-微服务架构/04-API网关",
    "05-微服务架构/05-服务网格",
    "05-微服务架构/06-消息队列",
    "05-微服务架构/07-微服务高级",
    
    # 06-云原生
    "06-云原生/01-容器化",
    "06-云原生/02-Kubernetes",
    "06-云原生/03-Service-Mesh",
    "06-云原生/04-云原生生态",
    
    # 07-性能优化
    "07-性能优化/01-性能分析",
    "07-性能优化/02-内存优化",
    "07-性能优化/03-并发优化",
    "07-性能优化/04-其他优化",
    "07-性能优化/05-常见问题",
    
    # 08-架构设计
    "08-架构设计/01-设计模式",
    "08-架构设计/02-架构模式",
    "08-架构设计/03-领域驱动设计",
    "08-架构设计/04-架构实战",
    
    # 09-工程实践
    "09-工程实践/01-测试",
    "09-工程实践/02-代码质量",
    "09-工程实践/03-CI-CD",
    "09-工程实践/04-监控与可观测性",
    "09-工程实践/05-文档与规范",
    
    # 10-Go版本特性
    "10-Go版本特性/01-Go-1.21特性",
    "10-Go版本特性/02-Go-1.22特性",
    "10-Go版本特性/03-Go-1.23特性",
    "10-Go版本特性/04-Go-1.24特性",
    "10-Go版本特性/05-Go-1.25特性",
    "10-Go版本特性/06-实验性特性",
    "10-Go版本特性/07-版本迁移指南",
    
    # 11-高级专题
    "11-高级专题/01-底层原理",
    "11-高级专题/02-汇编与CGO",
    "11-高级专题/03-安全",
    "11-高级专题/04-分布式系统",
    "11-高级专题/05-WebAssembly",
    "11-高级专题/06-新兴技术",
    
    # 12-行业应用
    "12-行业应用/01-金融科技",
    "12-行业应用/02-游戏开发",
    "12-行业应用/03-物联网",
    "12-行业应用/04-人工智能",
    "12-行业应用/05-大数据",
    "12-行业应用/06-电商系统",
    "12-行业应用/07-其他行业",
    
    # 13-参考资料
    "13-参考资料"
)

$createdCount = 0
$totalCount = $structure.Count

# README模板
$readmeTemplate = @"
# {0}

> 📚 **简介**: 本模块文档待补充

## 📋 内容列表

待添加...

## 🎯 学习目标

- 目标1
- 目标2
- 目标3

## 📚 参考资料

待添加...

---

**维护者**: Documentation Team  
**最后更新**: 2025年10月22日  
**状态**: 待完善
"@

# 创建目录
foreach ($dir in $structure) {
    $fullPath = Join-Path $BaseDir $dir
    
    if ($DryRun) {
        Write-Host "[DRY RUN] 将创建: $fullPath" -ForegroundColor Yellow
    } else {
        # 创建目录
        if (!(Test-Path $fullPath)) {
            New-Item -ItemType Directory -Path $fullPath -Force | Out-Null
            Write-Host "✅ 已创建: $fullPath" -ForegroundColor Green
            $createdCount++
        } else {
            Write-Host "⏭️  已存在: $fullPath" -ForegroundColor Gray
        }
        
        # 创建README.md
        $readmePath = Join-Path $fullPath "README.md"
        if (!(Test-Path $readmePath)) {
            $dirName = Split-Path $dir -Leaf
            $readmeContent = $readmeTemplate -f $dirName
            $readmeContent | Out-File -FilePath $readmePath -Encoding UTF8
            Write-Host "   📄 已创建README: $readmePath" -ForegroundColor Cyan
        }
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "✅ DRY RUN 完成! 共 $totalCount 个目录待创建" -ForegroundColor Yellow
} else {
    Write-Host "✅ 结构创建完成!" -ForegroundColor Green
    Write-Host "   新创建: $createdCount 个目录" -ForegroundColor Green
    Write-Host "   总目录: $totalCount 个" -ForegroundColor Cyan
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "📝 下一步操作:" -ForegroundColor Yellow
Write-Host "   1. 运行 migrate_documents.ps1 迁移文档" -ForegroundColor White
Write-Host "   2. 运行 fix_links.ps1 修复链接" -ForegroundColor White
Write-Host "   3. 运行 check_quality.ps1 质量检查" -ForegroundColor White
Write-Host ""

