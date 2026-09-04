# 检查仓库内失效的相对 Markdown 链接
# 扫描逻辑对齐 scripts/tmp/rescan_deadlinks.py：
#   - 链接正则: [锚文本](目标.md)（目标中不含空白与 # 之前的部分，含 URL 解码）
#   - 跳过 http/https 外链
#   - 跳过代码围栏 ``` 内的链接
#   - 跳过顶层目录 .git / node_modules / archive 下的所有 .md
# 退出码: 存在死链时 Exit 1，否则 Exit 0

param(
    [string]$Root = (Split-Path -Parent $PSScriptRoot)
)

$ErrorActionPreference = "Stop"

$linkRe = [regex]'\[([^\]]*)\]\(([^)#\s]+\.md)'
$skipTopDirs = @('.git', 'node_modules', 'archive')

function Test-InFence([string]$text, [int]$pos) {
    $count = 0
    $idx = 0
    while (($idx = $text.IndexOf('```', $idx)) -ge 0 -and $idx -lt $pos) {
        $count++
        $idx += 3
    }
    return ($count % 2 -eq 1)
}

# 收集 .md 文件（在根目录层剪枝 skipTopDirs，与 rescan 脚本一致）
$mdFiles = New-Object System.Collections.Generic.List[string]
$dirQueue = New-Object System.Collections.Generic.Queue[string]
$dirQueue.Enqueue($Root)
while ($dirQueue.Count -gt 0) {
    $dir = $dirQueue.Dequeue()
    foreach ($sub in [System.IO.Directory]::EnumerateDirectories($dir)) {
        $relTop = $sub.Substring($Root.Length).TrimStart('\','/').Split([char[]]@('\','/'))[0]
        if ($skipTopDirs -contains $relTop) { continue }
        $dirQueue.Enqueue($sub)
    }
    foreach ($f in [System.IO.Directory]::EnumerateFiles($dir, '*.md')) {
        $mdFiles.Add($f)
    }
}

Write-Host "扫描 $($mdFiles.Count) 个 Markdown 文件（跳过 $($skipTopDirs -join ', ')）..." -ForegroundColor Cyan

$dead = @{}      # target -> [ "file:line  锚文本" ]
$scanned = 0

foreach ($path in $mdFiles) {
    $rel = $path.Substring($Root.Length).TrimStart('\','/').Replace('\','/')
    $text = [System.IO.File]::ReadAllText($path)
    $scanned++

    foreach ($m in $linkRe.Matches($text)) {
        $url = $m.Groups[2].Value
        if ($url -match '^(https?://)') { continue }
        if (Test-InFence $text $m.Index) { continue }

        $dec = [System.Uri]::UnescapeDataString($url)
        $target = [System.IO.Path]::GetFullPath([System.IO.Path]::Combine([System.IO.Path]::GetDirectoryName($path), $dec))
        if (-not (Test-Path -LiteralPath $target)) {
            $line = ($text.Substring(0, $m.Index) -split "`n").Count
            $snippet = $m.Groups[1].Value
            if ($snippet.Length -gt 40) { $snippet = $snippet.Substring(0, 40) + '…' }
            if (-not $dead.ContainsKey($dec)) { $dead[$dec] = @() }
            $dead[$dec] += "$rel`:$line  [$snippet]"
        }
    }
}

$linksTotalDead = 0
foreach ($v in $dead.Values) { $linksTotalDead += $v.Count }

Write-Host ""
Write-Host "TOTAL_DEAD = $($dead.Count)"
Write-Host "LINKS_TOTAL_DEAD = $linksTotalDead"

if ($dead.Count -gt 0) {
    Write-Host ""
    Write-Host "死链明细（目标 <- 来源）:" -ForegroundColor Red
    foreach ($tgt in ($dead.Keys | Sort-Object)) {
        Write-Host ""
        Write-Host "  → $tgt" -ForegroundColor Yellow
        foreach ($src in $dead[$tgt]) {
            Write-Host "    $src" -ForegroundColor Gray
        }
    }
    Write-Host ""
    Write-Host "❌ 发现 $($dead.Count) 个死链目标 / $linksTotalDead 条死链" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "✅ 未发现死链（扫描 $scanned 个文件）" -ForegroundColor Green
exit 0
