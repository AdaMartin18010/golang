#!/bin/bash

# Model目录文档格式检查脚本
# 用于检查所有markdown文档的格式规范

echo "🔍 开始检查Model目录文档格式..."

# 统计变量
total_files=0
format_ok_files=0
format_error_files=0
missing_toc_files=0
unclosed_code_files=0

# 检查单个文件的格式
check_file() {
    local file="$1"
    local has_errors=false
    
    echo "📝 检查文件: $file"
    
    # 检查标题格式
    if grep -q "^# 1 1 1 1 1 1 1\|^## 9 9 9 9 9 9 9\|^## 13 13 13 13 13 13 13" "$file"; then
        echo "  ❌ 发现格式错误的标题"
        has_errors=true
    fi
    
    # 检查TOC格式
    if ! grep -q "<!-- TOC START -->" "$file"; then
        echo "  ⚠️  缺少TOC"
        missing_toc_files=$((missing_toc_files + 1))
    fi
    
    # 检查代码块
    local code_blocks=$(grep -c "```" "$file" 2>/dev/null || echo 0)
    if [ "$code_blocks" -gt 0 ] && [ $((code_blocks % 2)) -ne 0 ]; then
        echo "  ❌ 发现未闭合的代码块"
        unclosed_code_files=$((unclosed_code_files + 1))
        has_errors=true
    fi
    
    # 检查空行问题
    if grep -q "^$" "$file" && grep -A1 -B1 "^$" "$file" | grep -q "^$.*^$"; then
        echo "  ⚠️  发现多余空行"
    fi
    
    # 检查链接格式
    local broken_links=$(grep -c "\[.*\]()" "$file" 2>/dev/null || echo 0)
    if [ "$broken_links" -gt 0 ]; then
        echo "  ⚠️  发现 $broken_links 个空链接"
    fi
    
    if [ "$has_errors" = true ]; then
        format_error_files=$((format_error_files + 1))
        echo "  ❌ 格式检查失败"
    else
        format_ok_files=$((format_ok_files + 1))
        echo "  ✅ 格式检查通过"
    fi
    
    echo ""
}

# 遍历所有markdown文件
find model/ -name "*.md" -type f | while read file; do
    total_files=$((total_files + 1))
    check_file "$file"
done

echo "📊 格式检查完成统计:"
echo "  总文件数: $total_files"
echo "  格式正确: $format_ok_files"
echo "  格式错误: $format_error_files"
echo "  缺少TOC: $missing_toc_files"
echo "  未闭合代码块: $unclosed_code_files"
echo ""

# 生成详细报告
echo "📋 详细检查报告:"
echo ""

# 检查序号错误的文件
echo "🔢 序号错误文件:"
find model/ -name "*.md" -exec grep -l "1 1 1 1 1 1 1\|9 9 9 9 9 9 9\|13 13 13 13 13 13 13" {} \; 2>/dev/null | head -10
if [ $(find model/ -name "*.md" -exec grep -l "1 1 1 1 1 1 1\|9 9 9 9 9 9 9\|13 13 13 13 13 13 13" {} \; 2>/dev/null | wc -l) -gt 10 ]; then
    echo "  ... 还有更多文件"
fi
echo ""

# 检查缺少TOC的文件
echo "📑 缺少TOC的文件:"
find model/ -name "*.md" -not -exec grep -q "<!-- TOC START -->" {} \; -print | head -10
if [ $(find model/ -name "*.md" -not -exec grep -q "<!-- TOC START -->" {} \; -print | wc -l) -gt 10 ]; then
    echo "  ... 还有更多文件"
fi
echo ""

# 检查内容质量
echo "📄 内容质量检查:"
find model/ -name "*.md" -type f | while read file; do
    local lines=$(wc -l < "$file" 2>/dev/null || echo 0)
    local code_blocks=$(grep -c "```" "$file" 2>/dev/null || echo 0)
    local links=$(grep -c "\[.*\](.*)" "$file" 2>/dev/null || echo 0)
    
    if [ "$lines" -lt 50 ]; then
        echo "  ⚠️  $file: 文档过短 ($lines 行)"
    elif [ "$code_blocks" -eq 0 ] && [[ "$file" == *"README.md" ]]; then
        echo "  ⚠️  $file: 缺少代码示例"
    elif [ "$links" -eq 0 ] && [[ "$file" == *"README.md" ]]; then
        echo "  ⚠️  $file: 缺少内部链接"
    fi
done

echo ""
echo "✨ 格式检查脚本执行完成"
