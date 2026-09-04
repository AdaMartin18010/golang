import os
import shutil
from pathlib import Path

archive_dir = Path('e:/_src/golang/docs/archive/low-quality')
docs_dir = Path('e:/_src/golang/docs')
report = []

# 确保归档目录存在
archive_dir.mkdir(parents=True, exist_ok=True)

# 扫描所有markdown文件
placeholder_keywords = ['待补充', 'TODO', 'FIXME', 'placeholder', 'PLACEHOLDER']
placeholder_files = []
small_files = []

for root, dirs, files in os.walk(docs_dir):
    # 跳过archive目录
    if 'archive' in root:
        continue
    for file in files:
        if file.endswith('.md'):
            filepath = Path(root) / file
            try:
                content = filepath.read_text(encoding='utf-8', errors='ignore')
                # 检查是否包含占位符
                has_placeholder = any(kw in content for kw in placeholder_keywords)
                if has_placeholder:
                    placeholder_files.append(filepath)
                # 检查是否小于2KB
                elif filepath.stat().st_size < 2048:
                    small_files.append(filepath)
            except:
                pass

print(f'找到包含占位符的文件: {len(placeholder_files)} 个')
print(f'找到小于2KB的文件: {len(small_files)} 个')

# 合并文件列表并去重
all_files = list(dict.fromkeys(placeholder_files + small_files))
print(f'总共需要移动: {len(all_files)} 个文件')

# 移动文件并记录
for filepath in all_files:
    try:
        # 计算相对路径
        rel_path = filepath.relative_to(docs_dir)
        target_path = archive_dir / rel_path
        
        # 创建目标目录
        target_path.parent.mkdir(parents=True, exist_ok=True)
        
        # 移动文件
        shutil.move(str(filepath), str(target_path))
        
        reason = '包含占位符' if filepath in placeholder_files else '小于2KB'
        size_kb = filepath.stat().st_size / 1024
        report.append(f'{rel_path}|{size_kb:.2f}|{reason}')
        print(f'已移动: {rel_path}')
    except Exception as e:
        print(f'移动失败 {filepath}: {e}')

# 保存报告
report_path = archive_dir / 'cleanup-report.txt'
with open(report_path, 'w', encoding='utf-8') as f:
    f.write('文档清理报告\n')
    f.write('=' * 50 + '\n\n')
    f.write(f'清理时间: 2026-03-08 17:30:00\n')
    f.write(f'总共移动文件: {len(report)} 个\n\n')
    f.write('文件名|大小(KB)|原因\n')
    f.write('-' * 50 + '\n')
    for line in report:
        f.write(line + '\n')

print(f'\n报告已保存到: {report_path}')
