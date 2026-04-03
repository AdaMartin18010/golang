#!/usr/bin/env pwsh
#Requires -Version 5.1

<#
.SYNOPSIS
    Automatic Index Generator for Go Knowledge Base

.DESCRIPTION
    This script scans all markdown files in the knowledge base and generates
    cross-referenced indices including by-topic, by-tag, by-date, and complete-map.

.PARAMETER KnowledgeBasePath
    Path to the knowledge base directory. Defaults to parent of script directory.

.EXAMPLE
    .\generate-index.ps1
    Generates indices using default path

.EXAMPLE
    .\generate-index.ps1 -KnowledgeBasePath C:\path\to\go-knowledge-base
    Generates indices for specified path
#>

[CmdletBinding()]
param(
    [Parameter()]
    [string]$KnowledgeBasePath = (Split-Path -Parent $PSScriptRoot)
)

# Ensure path is absolute and normalized
$KBPath = (Resolve-Path $KnowledgeBasePath).Path
$IndicesDir = Join-Path $KBPath "indices"

# Console colors
$Colors = @{
    Info = "Cyan"
    Success = "Green"
    Warn = "Yellow"
    Error = "Red"
}

function Write-Log {
    param([string]$Message, [string]$Level = "Info")
    Write-Host "[$Level] $Message" -ForegroundColor $Colors[$Level]
}

# Ensure directories exist
function Initialize-Directories {
    if (-not (Test-Path $IndicesDir)) {
        New-Item -ItemType Directory -Path $IndicesDir -Force | Out-Null
    }
    Write-Log "Indices directory: $IndicesDir"
}

# Extract metadata from markdown file
function Get-DocumentMetadata {
    param([string]$FilePath)
    
    $relPath = $FilePath.Substring($KBPath.Length + 1).Replace('\', '/')
    $filename = Split-Path $FilePath -Leaf
    $content = Get-Content $FilePath -TotalCount 50 -ErrorAction SilentlyContinue
    $header = $content -join "`n"
    
    # Extract title
    $title = if ($content[0] -match '^#\s+(.+)$') { 
        $matches[1].Substring(0, [Math]::Min(100, $matches[1].Length))
    } else { 
        $filename 
    }
    
    # Extract dimension
    $dimension = if ($header -match '\*\*维度\*\*[：:]\s*([^|]+)') {
        $matches[1].Trim()
    } elseif ($header -match '>\s*\*\*维度\*\*\s*[:|]\s*([^|)]+)') {
        $matches[1].Trim()
    } else { "" }
    
    # Extract category
    $category = if ($header -match '\*\*分类\*\*[：:]\s*([^|]+)') {
        $matches[1].Trim()
    } else { "" }
    
    # Extract level
    $level = ""
    if ($header -match '\*\*级别\*\*[：:]\s*([SABC])') {
        $level = $matches[1]
    } elseif ($header -match '\*\*难度\*\*[：:]\s*(初级|中级|高级|Beginner|Intermediate|Advanced|Expert)') {
        switch ($matches[1]) {
            "初级" { $level = "B" }
            "中级" { $level = "A" }
            "高级" { $level = "S" }
            "Beginner" { $level = "B" }
            "Intermediate" { $level = "A" }
            "Advanced" { $level = "S" }
            "Expert" { $level = "S" }
            default { $level = "B" }
        }
    } elseif ($header -match 'S-Level|S级') {
        $level = "S"
    } elseif ($header -match 'A-Level|A级') {
        $level = "A"
    } elseif ($header -match 'B-Level|B级') {
        $level = "B"
    }
    
    # Extract tags
    $tags = @()
    if ($header -match '\*\*标签\*\*[：:]\s*([^|]+)') {
        $tags = $matches[1].Split(',') | ForEach-Object { $_.Trim().TrimStart('#') } | Where-Object { $_ }
    } elseif ($header -match '(#[a-z-]+)') {
        $tagMatches = $header | Select-String -Pattern '#([a-z-]+)' -AllMatches
        $tags = $tagMatches.Matches | ForEach-Object { $_.Groups[1].Value } | Select-Object -Unique
    }
    
    # Extract date
    $date = ""
    $datePatterns = @(
        '\*\*最后更新\*\*[：:]\s*([0-9]{4}-[0-9]{2}-[0-9]{2})',
        '\*\*完成日期\*\*[：:]\s*([0-9]{4}-[0-9]{2}-[0-9]{2})',
        '\*\*Created\*\*[：:]\s*([0-9]{4}-[0-9]{2}-[0-9]{2})'
    )
    foreach ($pattern in $datePatterns) {
        if ($header -match $pattern) {
            $date = $matches[1]
            break
        }
    }
    if (-not $date) {
        $date = (Get-Item $FilePath).LastWriteTime.ToString("yyyy-MM-dd")
    }
    
    # Infer dimension from path if not found
    if (-not $dimension) {
        $dimension = Get-DimensionFromPath -Path $relPath
    }
    
    # Infer level from file size if not found
    if (-not $level) {
        $level = Get-LevelFromSize -FilePath $FilePath
    }
    
    return [PSCustomObject]@{
        Title = $title
        Dimension = $dimension
        Category = $category
        Level = $level
        Tags = $tags
        Date = $date
        Path = $relPath
    }
}

function Get-DimensionFromPath {
    param([string]$Path)
    
    switch -Regex ($Path) {
        "^01-Formal-Theory" { return "Formal Theory" }
        "^02-Language-Design" { return "Language Design" }
        "^03-Engineering-CloudNative" { return "Engineering & Cloud Native" }
        "^04-Technology-Stack" { return "Technology Stack" }
        "^05-Application-Domains" { return "Application Domains" }
        "^examples" { return "Examples" }
        "^indices" { return "Indices" }
        "^learning-paths" { return "Learning Paths" }
        default { return "Other" }
    }
}

function Get-LevelFromSize {
    param([string]$FilePath)
    
    $size = (Get-Item $FilePath).Length
    if ($size -gt 15000) { return "S" }
    elseif ($size -gt 10000) { return "A" }
    elseif ($size -gt 5000) { return "B" }
    else { return "C" }
}

# Scan all markdown files
function Invoke-FileScan {
    Write-Log "Scanning markdown files in $KBPath..."
    
    $excludedDirs = @('indices', 'scripts', '.git')
    $files = Get-ChildItem -Path $KBPath -Filter "*.md" -Recurse | Where-Object {
        $file = $_
        -not ($excludedDirs | Where-Object { $file.FullName -like "*\$_\*" -or $file.FullName -like "*\$($_)\*" })
    }
    
    $documents = @()
    $count = 0
    
    foreach ($file in $files) {
        try {
            $metadata = Get-DocumentMetadata -FilePath $file.FullName
            $documents += $metadata
            
            $count++
            if ($count % 100 -eq 0) {
                Write-Host -NoNewline "`r[INFO] Scanned $count files..." -ForegroundColor $Colors.Info
            }
        }
        catch {
            Write-Log "Error scanning $($file.Name): $_" -Level "Warn"
        }
    }
    
    Write-Host ""
    Write-Log "Scanned $count markdown files" -Level "Success"
    
    return $documents
}

# Generate by-topic.md
function Export-ByTopicIndex {
    param([array]$Documents)
    
    $output = Join-Path $IndicesDir "by-topic.md"
    Write-Log "Generating by-topic.md..."
    
    $sb = New-Object System.Text.StringBuilder
    [void]$sb.AppendLine("# Go Knowledge Base - Topic-Based Cross Reference")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("> **Version**: Auto-generated")
    [void]$sb.AppendLine("> **Last Updated**: $(Get-Date -Format 'yyyy-MM-dd')")
    [void]$sb.AppendLine("> **Purpose**: Find documents by topic, concept, or technology")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("## 📑 Topics Overview")
    [void]$sb.AppendLine("")
    
    # Build topic index
    $topics = @{}
    foreach ($doc in $Documents) {
        # Add dimension as topic
        if ($doc.Dimension) {
            if (-not $topics[$doc.Dimension]) { $topics[$doc.Dimension] = @() }
            $topics[$doc.Dimension] += $doc
        }
        
        # Add category as topic
        if ($doc.Category) {
            if (-not $topics[$doc.Category]) { $topics[$doc.Category] = @() }
            $topics[$doc.Category] += $doc
        }
        
        # Add tags as topics
        foreach ($tag in $doc.Tags) {
            $cleanTag = $tag.Trim().ToLower()
            if ($cleanTag -and -not $topics[$cleanTag]) { $topics[$cleanTag] = @() }
            if ($cleanTag) { $topics[$cleanTag] += $doc }
        }
    }
    
    # Sort topics alphabetically
    $sortedTopics = $topics.Keys | Sort-Object
    
    # Group by first letter
    $currentLetter = ""
    foreach ($topic in $sortedTopics) {
        $firstLetter = $topic.Substring(0, 1).ToUpper()
        if ($firstLetter -ne $currentLetter) {
            $currentLetter = $firstLetter
            [void]$sb.AppendLine("")
            [void]$sb.AppendLine("### $currentLetter")
            [void]$sb.AppendLine("")
        }
        
        $docs = $topics[$topic] | Select-Object -First 5
        [void]$sb.AppendLine("- **$topic**")
        
        foreach ($doc in $docs) {
            $title = $doc.Title -replace '\[', '\[' -replace '\]', '\]'
            [void]$sb.AppendLine("  - [$title](../$($doc.Path))")
        }
        
        if ($topics[$topic].Count -gt 5) {
            [void]$sb.AppendLine("  - *(and $($topics[$topic].Count - 5) more...)*")
        }
        [void]$sb.AppendLine("")
    }
    
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("*This index is automatically generated. Run ``./scripts/generate-index.ps1`` to update.*")
    
    $sb.ToString() | Out-File -FilePath $output -Encoding UTF8
    $lineCount = (Get-Content $output).Count
    Write-Log "Generated $output ($lineCount lines)" -Level "Success"
}

# Generate by-tag.md
function Export-ByTagIndex {
    param([array]$Documents)
    
    $output = Join-Path $IndicesDir "by-tag.md"
    Write-Log "Generating by-tag.md..."
    
    $sb = New-Object System.Text.StringBuilder
    [void]$sb.AppendLine("# Go Knowledge Base - Tag Index")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("> **Version**: Auto-generated")
    [void]$sb.AppendLine("> **Last Updated**: $(Get-Date -Format 'yyyy-MM-dd')")
    [void]$sb.AppendLine("> **Purpose**: Find documents by tags")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("## 🏷️ Tag Index")
    [void]$sb.AppendLine("")
    
    # Collect tags
    $tagDocs = @{}
    foreach ($doc in $Documents) {
        foreach ($tag in $doc.Tags) {
            $cleanTag = $tag.Trim().ToLower()
            if ($cleanTag) {
                if (-not $tagDocs[$cleanTag]) { $tagDocs[$cleanTag] = @() }
                $tagDocs[$cleanTag] += $doc
            }
        }
        
        # Also add dimension as tag
        if ($doc.Dimension) {
            $dimTag = $doc.Dimension.ToLower() -replace '[^a-z0-9]', '-' -replace '-+', '-' -replace '^-|-$'
            if (-not $tagDocs[$dimTag]) { $tagDocs[$dimTag] = @() }
            $tagDocs[$dimTag] += $doc
        }
    }
    
    # Sort tags
    $sortedTags = $tagDocs.Keys | Sort-Object
    
    foreach ($tag in $sortedTags) {
        $count = $tagDocs[$tag].Count
        [void]$sb.AppendLine("### ` #$tag ` ($count documents)")
        [void]$sb.AppendLine("")
        
        $uniqueDocs = $tagDocs[$tag] | Sort-Object -Property Path -Unique
        foreach ($doc in $uniqueDocs) {
            $title = $doc.Title -replace '\[', '\[' -replace '\]', '\]'
            [void]$sb.AppendLine("- [$title](../$($doc.Path))")
        }
        [void]$sb.AppendLine("")
    }
    
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("*This index is automatically generated. Run ``./scripts/generate-index.ps1`` to update.*")
    
    $sb.ToString() | Out-File -FilePath $output -Encoding UTF8
    $lineCount = (Get-Content $output).Count
    Write-Log "Generated $output ($lineCount lines)" -Level "Success"
}

# Generate by-date.md
function Export-ByDateIndex {
    param([array]$Documents)
    
    $output = Join-Path $IndicesDir "by-date.md"
    Write-Log "Generating by-date.md..."
    
    $sb = New-Object System.Text.StringBuilder
    [void]$sb.AppendLine("# Go Knowledge Base - Chronological Index")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("> **Version**: Auto-generated")
    [void]$sb.AppendLine("> **Last Updated**: $(Get-Date -Format 'yyyy-MM-dd')")
    [void]$sb.AppendLine("> **Purpose**: Find documents by creation/update date")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("## 📅 Documents by Date")
    [void]$sb.AppendLine("")
    
    # Group by date
    $dateDocs = @{}
    foreach ($doc in $Documents) {
        $normDate = $doc.Date
        if ($normDate -match '^[0-9]{4}-[0-9]{2}-[0-9]{2}$') {
            if (-not $dateDocs[$normDate]) { $dateDocs[$normDate] = @() }
            $dateDocs[$normDate] += $doc
        }
    }
    
    # Sort dates in reverse chronological order
    $sortedDates = $dateDocs.Keys | Sort-Object -Descending
    
    # Group by month
    $currentMonth = ""
    foreach ($date in $sortedDates) {
        $month = $date.Substring(0, 7)
        if ($month -ne $currentMonth) {
            $currentMonth = $month
            [void]$sb.AppendLine("")
            [void]$sb.AppendLine("### $currentMonth")
            [void]$sb.AppendLine("")
            [void]$sb.AppendLine("| Date | Document | Dimension | Level |")
            [void]$sb.AppendLine("|------|----------|-----------|-------|")
        }
        
        foreach ($doc in $dateDocs[$date]) {
            $title = $doc.Title -replace '\|', '\|' -replace '\[', '\[' -replace '\]', '\]'
            [void]$sb.AppendLine("| $date | [$title](../$($doc.Path)) | $($doc.Dimension) | $($doc.Level) |")
        }
    }
    
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("## 📊 Monthly Statistics")
    [void]$sb.AppendLine("")
    
    $monthCounts = @{}
    foreach ($date in $sortedDates) {
        $month = $date.Substring(0, 7)
        if (-not $monthCounts[$month]) { $monthCounts[$month] = 0 }
        $monthCounts[$month] += $dateDocs[$date].Count
    }
    
    foreach ($month in ($monthCounts.Keys | Sort-Object -Descending)) {
        [void]$sb.AppendLine("- **$month**: $($monthCounts[$month]) documents")
    }
    
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("*This index is automatically generated. Run ``./scripts/generate-index.ps1`` to update.*")
    
    $sb.ToString() | Out-File -FilePath $output -Encoding UTF8
    $lineCount = (Get-Content $output).Count
    Write-Log "Generated $output ($lineCount lines)" -Level "Success"
}

# Generate complete-map.md
function Export-CompleteMapIndex {
    param([array]$Documents)
    
    $output = Join-Path $IndicesDir "complete-map.md"
    Write-Log "Generating complete-map.md..."
    
    $totalDocs = $Documents.Count
    $dimStats = $Documents | Group-Object -Property Dimension | Sort-Object -Property Count -Descending
    $levelStats = $Documents | Group-Object -Property Level | Sort-Object -Property Count -Descending
    
    $sb = New-Object System.Text.StringBuilder
    [void]$sb.AppendLine("# Go Knowledge Base - Complete Document Map")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("> **Version**: Auto-generated")
    [void]$sb.AppendLine("> **Last Updated**: $(Get-Date -Format 'yyyy-MM-dd')")
    [void]$sb.AppendLine("> **Purpose**: Complete inventory of all knowledge base documents")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("## 📊 Statistics")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("| Metric | Value |")
    [void]$sb.AppendLine("|--------|-------|")
    [void]$sb.AppendLine("| **Total Documents** | $totalDocs |")
    [void]$sb.AppendLine("| **Last Updated** | $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') |")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("### By Dimension")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("| Dimension | Count |")
    [void]$sb.AppendLine("|-----------|-------|")
    
    foreach ($stat in $dimStats) {
        $name = if ($stat.Name) { $stat.Name } else { "Other" }
        [void]$sb.AppendLine("| $name | $($stat.Count) |")
    }
    
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("### By Level")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("| Level | Count | Description |")
    [void]$sb.AppendLine("|-------|-------|-------------|")
    
    foreach ($stat in $levelStats) {
        $desc = switch ($stat.Name) {
            "S" { "Expert" }
            "A" { "Advanced" }
            "B" { "Intermediate" }
            "C" { "Basic" }
            default { "Unknown" }
        }
        [void]$sb.AppendLine("| $($stat.Name) | $($stat.Count) | $desc |")
    }
    
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("## 🗂️ Document Directory")
    [void]$sb.AppendLine("")
    
    # Group by dimension
    $dimOrder = @("Formal Theory", "Language Design", "Engineering & Cloud Native", "Technology Stack", "Application Domains", "Examples", "Learning Paths", "Other")
    
    foreach ($dim in $dimOrder) {
        $dimDocs = $Documents | Where-Object { $_.Dimension -eq $dim } | Sort-Object -Property Date -Descending
        
        if ($dimDocs) {
            [void]$sb.AppendLine("")
            [void]$sb.AppendLine("### $dim ($($dimDocs.Count) documents)")
            [void]$sb.AppendLine("")
            [void]$sb.AppendLine("| Document | Category | Level | Date | Path |")
            [void]$sb.AppendLine("|----------|----------|-------|------|------|")
            
            foreach ($doc in $dimDocs) {
                $title = $doc.Title -replace '\|', '\|' -replace '\[', '\[' -replace '\]', '\]'
                $cat = if ($doc.Category) { $doc.Category } else { "-" }
                [void]$sb.AppendLine("| [$title](../$($doc.Path)) | $cat | $($doc.Level) | $($doc.Date) | $(Split-Path $doc.Path -Leaf) |")
            }
        }
    }
    
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("## 🔍 Quick Reference")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("### Document ID Prefixes")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("| Prefix | Dimension |")
    [void]$sb.AppendLine("|--------|-----------|")
    [void]$sb.AppendLine("| FT-* | Formal Theory |")
    [void]$sb.AppendLine("| LD-* | Language Design |")
    [void]$sb.AppendLine("| EC-* | Engineering & Cloud Native |")
    [void]$sb.AppendLine("| TS-* | Technology Stack |")
    [void]$sb.AppendLine("| AD-* | Application Domains |")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("### Level Definitions")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("| Level | Description | Target Audience |")
    [void]$sb.AppendLine("|-------|-------------|-----------------|")
    [void]$sb.AppendLine("| S | Expert/S-Level | Principal Engineers, Researchers |")
    [void]$sb.AppendLine("| A | Advanced/A-Level | Senior Engineers |")
    [void]$sb.AppendLine("| B | Intermediate/B-Level | Mid-level Engineers |")
    [void]$sb.AppendLine("| C | Basic/C-Level | Junior Engineers |")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("---")
    [void]$sb.AppendLine("")
    [void]$sb.AppendLine("*This index is automatically generated. Run ``./scripts/generate-index.ps1`` to update.*")
    
    $sb.ToString() | Out-File -FilePath $output -Encoding UTF8
    $lineCount = (Get-Content $output).Count
    Write-Log "Generated $output ($lineCount lines)" -Level "Success"
}

# Update indices README
function Update-IndicesReadme {
    $output = Join-Path $IndicesDir "README.md"
    Write-Log "Updating indices/README.md..."
    
    $content = @"
# Go Knowledge Base - Indices

> **Version**: Auto-generated
> **Last Updated**: $(Get-Date -Format 'yyyy-MM-dd')

This directory contains auto-generated indices for the Go Knowledge Base.

## Available Indices

| Index | Description | File |
|-------|-------------|------|
| **By Topic** | Documents organized by topic/concept | [by-topic.md](./by-topic.md) |
| **By Tag** | Documents organized by tags | [by-tag.md](./by-tag.md) |
| **By Date** | Chronological listing of documents | [by-date.md](./by-date.md) |
| **Complete Map** | Full inventory with statistics | [complete-map.md](./complete-map.md) |
| **By Difficulty** | Learning paths by experience level | [by-difficulty.md](./by-difficulty.md) |
| **Cross Reference** | Topic relationships and pathways | [cross-reference.md](./cross-reference.md) |

## How to Update

Run the index generator script from the knowledge base root:

```powershell
.\scripts\generate-index.ps1
```

Or from anywhere with the path:

```powershell
.\path\to\scripts\generate-index.ps1 C:\path\to\go-knowledge-base
```

## Index Details

### by-topic.md
Alphabetical index of topics with associated documents. Topics are extracted from:
- Document dimensions (维度)
- Document categories (分类)
- Document tags (标签)

### by-tag.md
Tag cloud style index with document counts per tag.

### by-date.md
Chronological view of document creation/updates, grouped by month.

### complete-map.md
Comprehensive document inventory including:
- Complete statistics
- All documents sorted by dimension
- Quick reference guides

---

*This README is automatically updated when indices are regenerated.*
"@
    
    $content | Out-File -FilePath $output -Encoding UTF8
    Write-Log "Updated $output" -Level "Success"
}

# Main execution
function Main {
    Write-Log "Starting index generation..."
    Write-Log "Knowledge Base Path: $KBPath"
    
    if (-not (Test-Path $KBPath)) {
        Write-Log "Knowledge base path does not exist: $KBPath" -Level "Error"
        exit 1
    }
    
    # Ensure directories exist
    Initialize-Directories
    
    # Scan files
    $documents = Invoke-FileScan
    
    if ($documents.Count -eq 0) {
        Write-Log "No documents found to index!" -Level "Error"
        exit 1
    }
    
    # Generate indices
    Export-ByTopicIndex -Documents $documents
    Export-ByTagIndex -Documents $documents
    Export-ByDateIndex -Documents $documents
    Export-CompleteMapIndex -Documents $documents
    Update-IndicesReadme
    
    # Show statistics
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  Index Generation Complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "Total Documents: $($documents.Count)" -ForegroundColor Cyan
    Write-Host "Indices Generated:" -ForegroundColor Cyan
    Write-Host "  - by-topic.md" -ForegroundColor White
    Write-Host "  - by-tag.md" -ForegroundColor White
    Write-Host "  - by-date.md" -ForegroundColor White
    Write-Host "  - complete-map.md" -ForegroundColor White
    Write-Host "========================================" -ForegroundColor Green
    
    Write-Log "Index generation complete!" -Level "Success"
    Write-Log "Indices location: $IndicesDir"
}

# Run main
Main
