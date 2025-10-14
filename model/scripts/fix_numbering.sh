#!/bin/bash

# Model目录文档序号修正脚本
# 用于批量修正所有markdown文档的序号错误

echo "🚀 开始修正Model目录文档序号..."

# 定义需要修正的模式
declare -A patterns=(
    ["^# 1 1 1 1 1 1 1"]="# "
    ["^## 1 1 1 1 1 1 1"]="## 1. "
    ["^### 1 1 1 1 1 1 1"]="### 1.1 "
    ["^#### 1 1 1 1 1 1 1"]="#### 1.1.1 "
    ["^## 9 9 9 9 9 9 9"]="## 2. "
    ["^### 9 9 9 9 9 9 9"]="### 2.1 "
    ["^## 13 13 13 13 13 13 13"]="## 3. "
    ["^### 13 13 13 13 13 13 13"]="### 3.1 "
    ["^## 14 14 14 14 14 14 14"]="## 4. "
    ["^### 14 14 14 14 14 14 14"]="### 4.1 "
    ["^## 15 15 15 15 15 15 15"]="## 5. "
    ["^### 15 15 15 15 15 15 15"]="### 5.1 "
    ["^## 7 7 7 7 7 7 7"]="## 6. "
    ["^### 7 7 7 7 7 7 7"]="### 6.1 "
    ["^## 8 8 8 8 8 8 8"]="## 7. "
    ["^### 8 8 8 8 8 8 8"]="### 7.1 "
    ["^## 11 11 11 11 11 11 11"]="## 8. "
    ["^### 11 11 11 11 11 11 11"]="### 8.1 "
    ["^## 12 12 12 12 12 12 12"]="## 9. "
    ["^### 12 12 12 12 12 12 12"]="### 9.1 "
)

# 统计变量
total_files=0
fixed_files=0
skipped_files=0

# 遍历所有markdown文件
find model/ -name "*.md" -type f | while read file; do
    total_files=$((total_files + 1))
    echo "📝 处理文件: $file"
    
    # 检查文件是否需要修正
    needs_fix=false
    for pattern in "${!patterns[@]}"; do
        if grep -q "$pattern" "$file"; then
            needs_fix=true
            break
        fi
    done
    
    if [ "$needs_fix" = true ]; then
        echo "  🔧 发现序号错误，开始修正..."
        
        # 创建备份
        cp "$file" "$file.bak"
        
        # 应用所有修正规则
        for pattern in "${!patterns[@]}"; do
            sed -i "s|$pattern|${patterns[$pattern]}|g" "$file"
        done
        
        # 验证修正结果
        if grep -q "1 1 1 1 1 1 1\|9 9 9 9 9 9 9\|13 13 13 13 13 13 13" "$file"; then
            echo "  ⚠️  仍有序号错误，请手动检查"
        else
            echo "  ✅ 序号修正完成"
            fixed_files=$((fixed_files + 1))
            # 删除备份文件
            rm "$file.bak"
        fi
    else
        echo "  ✅ 序号格式正确，跳过"
        skipped_files=$((skipped_files + 1))
    fi
    
    echo ""
done

echo "📊 修正完成统计:"
echo "  总文件数: $total_files"
echo "  修正文件数: $fixed_files"
echo "  跳过文件数: $skipped_files"
echo ""

# 检查是否还有序号错误
echo "🔍 检查剩余序号错误..."
remaining_errors=$(find model/ -name "*.md" -exec grep -l "1 1 1 1 1 1 1\|9 9 9 9 9 9 9\|13 13 13 13 13 13 13" {} \; 2>/dev/null | wc -l)

if [ "$remaining_errors" -gt 0 ]; then
    echo "⚠️  仍有 $remaining_errors 个文件存在序号错误，需要手动处理:"
    find model/ -name "*.md" -exec grep -l "1 1 1 1 1 1 1\|9 9 9 9 9 9 9\|13 13 13 13 13 13 13" {} \; 2>/dev/null
else
    echo "🎉 所有文档序号修正完成！"
fi

echo ""
echo "✨ 序号修正脚本执行完成"
