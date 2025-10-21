#!/usr/bin/env python3
"""
代码格式统一修复脚本
功能：统一代码块格式、Mermaid图表格式、移除多余空行
"""

import re
import os
import sys
from pathlib import Path
from typing import Tuple
from datetime import datetime

# 统计变量
stats = {
    'files_scanned': 0,
    'files_modified': 0,
    'issues': {
        'trailing_empty_lines_in_code_blocks': 0,
        'excessive_newlines': 0,
        'trailing_spaces': 0,
        'missing_language_tags': 0,
        'mermaid_format': 0,
    }
}

def log(message: str, level: str = "INFO"):
    """打印日志"""
    timestamp = datetime.now().strftime("%H:%M:%S")
    colors = {
        "ERROR": "\033[91m",
        "WARN": "\033[93m",
        "SUCCESS": "\033[92m",
        "INFO": "\033[97m",
    }
    reset = "\033[0m"
    color = colors.get(level, colors["INFO"])
    print(f"{color}[{timestamp}] {level}: {message}{reset}")

def fix_code_blocks(content: str) -> Tuple[str, bool]:
    """修复代码块格式"""
    modified = False
    lines = content.split('\n')
    result = []
    in_code_block = False
    code_block_lang = ""
    i = 0
    
    while i < len(lines):
        line = lines[i]
        
        # 检测代码块开始
        match = re.match(r'^```(\w*)(.*)$', line)
        if match and not in_code_block:
            lang = match.group(1)
            extra = match.group(2).strip()
            
            # 尝试推断缺失的语言标记
            if not lang and i + 1 < len(lines):
                next_line = lines[i + 1]
                if re.match(r'^(package|func|import|type|var|const)\s', next_line):
                    lang = "go"
                    modified = True
                    stats['issues']['missing_language_tags'] += 1
                elif re.match(r'^(\$|#|cd|ls|git|npm|go\s)', next_line):
                    lang = "bash"
                    modified = True
                    stats['issues']['missing_language_tags'] += 1
            
            # 统一格式：```language
            result.append(f"```{lang}")
            if extra:
                modified = True
            
            in_code_block = True
            code_block_lang = lang
            i += 1
            continue
        
        # 检测代码块结束
        if line.strip() == '```' and in_code_block:
            # 移除代码块末尾的空行
            while result and not result[-1].strip():
                result.pop()
                modified = True
                stats['issues']['trailing_empty_lines_in_code_blocks'] += 1
            
            result.append('```')
            in_code_block = False
            code_block_lang = ""
            
            # 代码块后应该有一个空行
            if i + 1 < len(lines) and lines[i + 1].strip():
                result.append('')
                modified = True
            
            i += 1
            continue
        
        # 在代码块内或外，移除行尾空格
        trimmed = line.rstrip()
        if trimmed != line:
            modified = True
            stats['issues']['trailing_spaces'] += 1
        result.append(trimmed)
        i += 1
    
    return '\n'.join(result), modified

def fix_mermaid_diagrams(content: str) -> Tuple[str, bool]:
    """修复Mermaid图表格式"""
    modified = False
    
    # 统一Mermaid开始标记
    original = content
    content = re.sub(r'^```mermaid\s+.*$', '```mermaid', content, flags=re.MULTILINE)
    if content != original:
        modified = True
        stats['issues']['mermaid_format'] += 1
    
    return content, modified

def fix_excessive_newlines(content: str) -> Tuple[str, bool]:
    """修复连续空行（3个以上 → 2个）"""
    modified = False
    
    # 修复：连续3个以上空行 → 2个空行
    original = content
    content = re.sub(r'\n\s*\n\s*\n\s*\n+', '\n\n\n', content)
    if content != original:
        modified = True
        stats['issues']['excessive_newlines'] += 1
    
    return content, modified

def process_markdown_file(filepath: Path, dry_run: bool = False) -> None:
    """处理单个Markdown文件"""
    stats['files_scanned'] += 1
    
    try:
        # 读取文件
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        if not content:
            return
        
        original_content = content
        file_modified = False
        
        # 1. 修复代码块格式
        content, mod1 = fix_code_blocks(content)
        if mod1:
            file_modified = True
        
        # 2. 修复Mermaid图表格式
        content, mod2 = fix_mermaid_diagrams(content)
        if mod2:
            file_modified = True
        
        # 3. 修复连续空行
        content, mod3 = fix_excessive_newlines(content)
        if mod3:
            file_modified = True
        
        # 4. 确保文件末尾只有一个换行符
        content = content.rstrip() + '\n'
        
        # 保存文件
        if file_modified or content != original_content:
            if not dry_run:
                with open(filepath, 'w', encoding='utf-8', newline='\n') as f:
                    f.write(content)
                log(f"Modified: {filepath.name}", "SUCCESS")
            else:
                log(f"Would modify: {filepath.name}", "WARN")
            stats['files_modified'] += 1
    
    except Exception as e:
        log(f"Error processing {filepath}: {e}", "ERROR")

def main():
    """主程序"""
    import argparse
    
    parser = argparse.ArgumentParser(description='代码格式统一修复脚本')
    parser.add_argument('target_dir', nargs='?', default='docs', help='目标目录')
    parser.add_argument('--dry-run', action='store_true', help='Dry Run模式（预览）')
    args = parser.parse_args()
    
    target_dir = Path(args.target_dir)
    dry_run = args.dry_run
    
    log("=" * 50, "INFO")
    log("代码格式统一修复脚本", "INFO")
    log("=" * 50, "INFO")
    log(f"目标目录: {target_dir}", "INFO")
    log(f"模式: {'Dry Run (预览)' if dry_run else '实际修改'}", "INFO")
    log("=" * 50, "INFO")
    
    # 获取所有Markdown文件
    md_files = list(target_dir.rglob('*.md'))
    log(f"找到 {len(md_files)} 个Markdown文件", "INFO")
    print()
    
    # 处理每个文件
    for filepath in md_files:
        process_markdown_file(filepath, dry_run)
    
    # 输出统计报告
    print()
    log("=" * 50, "SUCCESS")
    log("修复完成！统计报告：", "SUCCESS")
    log("=" * 50, "INFO")
    log(f"扫描文件数: {stats['files_scanned']}", "INFO")
    log(f"修改文件数: {stats['files_modified']}", "INFO")
    log("-" * 50, "INFO")
    log("修复问题统计：", "INFO")
    log(f"缺失语言标记: {stats['issues']['missing_language_tags']}", "INFO")
    log(f"代码块尾部空行: {stats['issues']['trailing_empty_lines_in_code_blocks']}", "INFO")
    log(f"Mermaid格式: {stats['issues']['mermaid_format']}", "INFO")
    log(f"连续空行修复: {stats['issues']['excessive_newlines']}", "INFO")
    log(f"行尾空格移除: {stats['issues']['trailing_spaces']}", "INFO")
    log("=" * 50, "INFO")
    
    if dry_run:
        print()
        log("这是Dry Run模式，未实际修改文件", "WARN")
        log("移除 --dry-run 参数以执行实际修改", "WARN")

if __name__ == '__main__':
    main()

