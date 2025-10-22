# scripts/fix_broken_links.ps1
# 修复失效链接工具

param (
    [string]$DocsPath = "docs-new"
)

Write-Host "=== 🔧 链接修复工具 ===" -ForegroundColor Cyan
Write-Host ""

# URL编码函数
function UrlEncode-Chinese {
    param([string]$text)
    
    $encoded = ""
    for ($i = 0; $i -lt $text.Length; $i++) {
        $char = $text[$i]
        if ($char -match '[\u4e00-\u9fa5]' -or $char -match '[^\w\-]') {
            # 中文字符或特殊字符需要编码
            $bytes = [System.Text.Encoding]::UTF8.GetBytes($char)
            foreach ($byte in $bytes) {
                $encoded += "%{0:X2}" -f $byte
            }
        } else {
            $encoded += $char
        }
    }
    return $encoded.ToLower()
}

$fixedCount = 0
$skippedCount = 0

# 需要修复的文件列表
$filesToFix = @(
    "INDEX.md",
    "LEARNING_PATHS.md",
    "FAQ.md",
    "QUICK_START.md",
    "README.md",
    "01-语言基础\README.md",
    "01-语言基础\01-语法基础\01-Hello-World.md",
    "01-语言基础\00-Go语言形式化语义与理论基础.md",
    "03-Web开发\00-HTTP编程深度实战指南.md",
    "04-数据库编程\01-MySQL编程.md",
    "04-数据库编程\02-PostgreSQL编程.md",
    "04-数据库编程\03-Redis编程.md",
    "05-微服务架构\10-高性能微服务架构.md",
    "05-微服务架构\11-Kubernetes微服务部署.md",
    "05-微服务架构\13-GitOps持续部署.md",
    "05-微服务架构\15-微服务安全最佳实践.md",
    "05-微服务架构\README.md",
    "06-云原生与容器\05-服务网格集成.md",
    "06-云原生与容器\06-GitOps部署.md",
    "06-云原生与容器\README.md",
    "07-性能优化\01-性能分析与pprof.md",
    "08-架构设计\01-创建型模式.md",
    "08-架构设计\03-行为型模式.md"
)

foreach ($relPath in $filesToFix) {
    $filePath = Join-Path $DocsPath $relPath
    
    if (-not (Test-Path $filePath)) {
        Write-Host "  ⚠️ 文件不存在: $relPath" -ForegroundColor Yellow
        continue
    }
    
    $content = Get-Content $filePath -Raw
    $originalContent = $content
    $fileFixed = $false
    
    # 1. 修复中文锚点链接（URL编码）
    $anchorMatches = [regex]::Matches($content, '\[([^\]]+)\]\(([^)]*\.md)(#[^)]+)\)')
    foreach ($match in $anchorMatches) {
        $linkText = $match.Groups[1].Value
        $mdPath = $match.Groups[2].Value
        $anchor = $match.Groups[3].Value
        
        # 如果锚点包含中文或空格，进行编码
        if ($anchor -match '[\u4e00-\u9fa5\s]') {
            $cleanAnchor = $anchor.Substring(1) # 去掉 #
            $encodedAnchor = "#" + (UrlEncode-Chinese $cleanAnchor)
            $oldLink = "[$linkText]($mdPath$anchor)"
            $newLink = "[$linkText]($mdPath$encodedAnchor)"
            $content = $content.Replace($oldLink, $newLink)
            $fileFixed = $true
        }
    }
    
    # 2. 修复特定的已知问题链接
    # 使用数组而不是哈希表来避免特殊字符问题
    $replacements = @(
        # FAQ.md, QUICK_START.md, README.md - 锚点链接
        @{ Old = "LEARNING_PATHS.md#零基础入门路径"; New = "LEARNING_PATHS.md#%E9%9B%B6%E5%9F%BA%E7%A1%80%E5%85%A5%E9%97%A8%E8%B7%AF%E5%BE%84" }
        @{ Old = "LEARNING_PATHS.md#web开发路径"; New = "LEARNING_PATHS.md#web%E5%BC%80%E5%8F%91%E8%B7%AF%E5%BE%84" }
        @{ Old = "LEARNING_PATHS.md#微服务开发路径"; New = "LEARNING_PATHS.md#%E5%BE%AE%E6%9C%8D%E5%8A%A1%E5%BC%80%E5%8F%91%E8%B7%AF%E5%BE%84" }
        @{ Old = "LEARNING_PATHS.md#云原生路径"; New = "LEARNING_PATHS.md#%E4%BA%91%E5%8E%9F%E7%94%9F%E8%B7%AF%E5%BE%84" }
        @{ Old = "LEARNING_PATHS.md#算法面试路径"; New = "LEARNING_PATHS.md#%E7%AE%97%E6%B3%95%E9%9D%A2%E8%AF%95%E8%B7%AF%E5%BE%84" }
        @{ Old = "LEARNING_PATHS.md#性能优化路径"; New = "LEARNING_PATHS.md#%E6%80%A7%E8%83%BD%E4%BC%98%E5%8C%96%E8%B7%AF%E5%BE%84" }
        @{ Old = "LEARNING_PATHS.md#架构师路径"; New = "LEARNING_PATHS.md#%E6%9E%B6%E6%9E%84%E5%B8%88%E8%B7%AF%E5%BE%84" }
        @{ Old = "LEARNING_PATHS.md#📅-4周快速入门"; New = "LEARNING_PATHS.md#-4%E5%91%A8%E5%BF%AB%E9%80%9F%E5%85%A5%E9%97%A8" }
        @{ Old = "LEARNING_PATHS.md#📅-3个月成为go开发者"; New = "LEARNING_PATHS.md#-3%E4%B8%AA%E6%9C%88%E6%88%90%E4%B8%BAgo%E5%BC%80%E5%8F%91%E8%80%85" }
        @{ Old = "LEARNING_PATHS.md#📅-6个月进阶"; New = "LEARNING_PATHS.md#-6%E4%B8%AA%E6%9C%88%E8%BF%9B%E9%98%B6" }
        @{ Old = "LEARNING_PATHS.md#📅-1年成为专家"; New = "LEARNING_PATHS.md#-1%E5%B9%B4%E6%88%90%E4%B8%BA%E4%B8%93%E5%AE%B6" }
        @{ Old = "INDEX.md#按难度等级索引"; New = "INDEX.md#%E6%8C%89%E9%9A%BE%E5%BA%A6%E7%AD%89%E7%BA%A7%E7%B4%A2%E5%BC%95" }
        @{ Old = "INDEX.md#按应用场景索引"; New = "INDEX.md#%E6%8C%89%E5%BA%94%E7%94%A8%E5%9C%BA%E6%99%AF%E7%B4%A2%E5%BC%95" }
        
        # 移除指向不存在文件的链接
        @{ Old = "[版本选择](07-版本选择.md)"; New = "" }
        @{ Old = "[私有模块](08-私有模块.md)"; New = "" }
        @{ Old = "[模块代理](09-模块代理.md)"; New = "" }
        @{ Old = "[Vendor目录](10-Vendor目录.md)"; New = "" }
        @{ Old = "[工作区模式](11-工作区模式.md)"; New = "" }
        @{ Old = "[Service Mesh集成](./12-Service Mesh集成.md)"; New = "" }
        @{ Old = "[GitHub Actions](./07-GitHub Actions.md)"; New = "" }
        @{ Old = "[GitLab CI](./08-GitLab CI.md)"; New = "" }
        
        # 修复ORM链接
        @{ Old = "04-ORM框架-GORM.md"; New = "../01-语言基础/README.md" }
        
        # 修复跨模块链接
        @{ Old = "../06-云原生/01-容器化部署.md"; New = "../06-云原生与容器/01-Docker容器化.md" }
        @{ Old = "../06-云原生/07-GitHub-Actions.md"; New = "../06-云原生与容器/README.md" }
        @{ Old = "../05-微服务/12-Service-Mesh集成.md"; New = "../05-微服务架构/README.md" }
        @{ Old = "../05-微服务/13-GitOps持续部署.md"; New = "../05-微服务架构/13-GitOps持续部署.md" }
        @{ Old = "../02-Web开发/00-Go认证与授权深度实战指南.md"; New = "../03-Web开发/00-Go认证与授权深度实战指南.md" }
        @{ Old = "../07-性能优化/01-性能分析工具.md"; New = "../07-性能优化/01-性能分析与pprof.md" }
        @{ Old = "../07-性能优化/02-缓存优化.md"; New = "../07-性能优化/README.md" }
        @{ Old = "../08-架构设计/01-领域驱动设计.md"; New = "../08-架构设计/README.md" }
        @{ Old = "../DOCUMENT_STANDARD.md"; New = "../README.md" }
        @{ Old = "../LICENSE"; New = "../README.md" }
        @{ Old = "../../issues"; New = "https://github.com/yourusername/golang-docs/issues" }
    )
    
    foreach ($pair in $replacements) {
        if ($content.Contains($pair.Old)) {
            $content = $content.Replace($pair.Old, $pair.New)
            $fileFixed = $true
        }
    }
    
    # 保存修改
    if ($fileFixed) {
        Set-Content -Path $filePath -Value $content -Encoding UTF8
        Write-Host "  ✅ 已修复: $relPath" -ForegroundColor Green
        $fixedCount++
    } else {
        Write-Host "  跳过: $relPath (无需修复)" -ForegroundColor DarkGray
        $skippedCount++
    }
}

Write-Host ""
Write-Host "=== 完成 ===" -ForegroundColor Cyan
Write-Host "  检查文件: $($filesToFix.Count)" -ForegroundColor Green
Write-Host "  已修复: $fixedCount" -ForegroundColor Green
Write-Host "  已跳过: $skippedCount" -ForegroundColor Green
