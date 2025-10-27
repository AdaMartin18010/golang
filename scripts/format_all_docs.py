#!/usr/bin/env python3
"""
批量格式化文档脚本
为所有新创建的文档添加目录和序号
"""

import os
import re
from pathlib import Path

# 需要更新的所有文档列表
DOCS_TO_UPDATE = [
    # testing (剩余6个，01已完成，02已完成)
    "docs/practices/testing/03-集成测试.md",
    "docs/practices/testing/04-性能测试.md",
    "docs/practices/testing/05-测试覆盖率.md",
    "docs/practices/testing/06-Mock与Stub.md",
    "docs/practices/testing/07-测试最佳实践.md",
    "docs/practices/testing/08-常见问题与技巧.md",
    
    # deployment (7个)
    "docs/practices/deployment/01-部署概览.md",
    "docs/practices/deployment/02-Docker部署.md",
    "docs/practices/deployment/03-Kubernetes部署.md",
    "docs/practices/deployment/04-CI-CD流程.md",
    "docs/practices/deployment/05-监控与日志.md",
    "docs/practices/deployment/06-滚动更新.md",
    "docs/practices/deployment/07-生产环境最佳实践.md",
    
    # distributed (6个)
    "docs/advanced/distributed/01-分布式系统基础.md",
    "docs/advanced/distributed/02-服务注册与发现.md",
    "docs/advanced/distributed/03-分布式一致性.md",
    "docs/advanced/distributed/04-分布式锁.md",
    "docs/advanced/distributed/05-分布式事务.md",
    "docs/advanced/distributed/06-负载均衡.md",
    
    # security (6个)
    "docs/advanced/security/01-Web安全基础.md",
    "docs/advanced/security/02-身份认证.md",
    "docs/advanced/security/03-授权机制.md",
    "docs/advanced/security/04-数据保护.md",
    "docs/advanced/security/05-安全审计.md",
    "docs/advanced/security/06-最佳实践.md",
    
    # concurrency (5个)
    "docs/fundamentals/concurrency/01-并发基础概念.md",
    "docs/fundamentals/concurrency/02-Goroutine深入.md",
    "docs/fundamentals/concurrency/03-Channel深入.md",
    "docs/fundamentals/concurrency/04-Context应用.md",
    "docs/fundamentals/concurrency/05-并发模式.md",
    
    # stdlib (1个)
    "docs/fundamentals/stdlib/01-核心包概览.md",
    
    # ai-ml (6个)
    "docs/advanced/ai-ml/01-Go与AI集成.md",
    "docs/advanced/ai-ml/02-机器学习库.md",
    "docs/advanced/ai-ml/03-深度学习框架.md",
    "docs/advanced/ai-ml/04-模型推理.md",
    "docs/advanced/ai-ml/05-数据处理.md",
    "docs/advanced/ai-ml/06-实战案例.md",
    
    # modern-web (6个)
    "docs/advanced/modern-web/01-现代Web框架.md",
    "docs/advanced/modern-web/02-实时通信.md",
    "docs/advanced/modern-web/03-GraphQL.md",
    "docs/advanced/modern-web/04-微服务网关.md",
    "docs/advanced/modern-web/05-服务网格.md",
    "docs/advanced/modern-web/06-云原生实践.md",
    
    # templates (6个)
    "docs/projects/templates/01-项目结构模板.md",
    "docs/projects/templates/02-微服务模板.md",
    "docs/projects/templates/03-Web应用模板.md",
    "docs/projects/templates/04-CLI工具模板.md",
    "docs/projects/templates/05-库项目模板.md",
    "docs/projects/templates/06-快速开始指南.md",
    
    # reference/api (4个)
    "docs/reference/api/01-核心API参考.md",
    "docs/reference/api/02-标准库API.md",
    "docs/reference/api/03-常用第三方库.md",
    "docs/reference/api/04-API设计指南.md",
    
    # reference/guides (3个)
    "docs/reference/guides/01-学习路线图.md",
    "docs/reference/guides/02-资源汇总.md",
    "docs/reference/guides/03-常见问题.md",
]

def extract_sections(content):
    """提取文档中的章节"""
    lines = content.split('\n')
    sections = []
    h2_count = 0
    h3_counts = {}
    
    for i, line in enumerate(lines):
        # 匹配二级标题 ##
        h2_match = re.match(r'^##\s+(.+)$', line)
        if h2_match:
            title = h2_match.group(1)
            # 排除目录
            if '目录' not in title and 'TOC' not in title.upper():
                h2_count += 1
                h3_counts[h2_count] = 0
                emoji = extract_emoji(title)
                clean_title = title.replace(emoji, '').strip()
                sections.append({
                    'level': 2,
                    'line': i,
                    'number': str(h2_count),
                    'emoji': emoji,
                    'title': clean_title,
                    'original': line
                })
                continue
        
        # 匹配三级标题 ###
        h3_match = re.match(r'^###\s+(.+)$', line)
        if h3_match and h2_count > 0:
            title = h3_match.group(1)
            h3_counts[h2_count] += 1
            emoji = extract_emoji(title)
            clean_title = title.replace(emoji, '').strip()
            # 移除开头的数字（如果有）
            clean_title = re.sub(r'^\d+\.\s*', '', clean_title)
            clean_title = re.sub(r'^\d+\.\d+\s*', '', clean_title)
            sections.append({
                'level': 3,
                'line': i,
                'number': f"{h2_count}.{h3_counts[h2_count]}",
                'emoji': emoji,
                'title': clean_title,
                'original': line
            })
    
    return sections

def extract_emoji(text):
    """提取emoji"""
    # 常见的emoji模式
    emoji_pattern = r'^([^\w\s]+)\s+'
    match = re.match(emoji_pattern, text)
    if match:
        return match.group(1)
    return ''

def generate_toc(sections):
    """生成目录"""
    toc_lines = ['## 📋 目录', '']
    
    for section in sections:
        indent = '  ' * (section['level'] - 2)
        anchor_title = section['title'].lower().replace(' ', '-')
        # 移除特殊字符
        anchor_title = re.sub(r'[^\w\-\u4e00-\u9fa5]', '', anchor_title)
        
        emoji_part = f" {section['emoji']}" if section['emoji'] else ""
        anchor = f"#{section['number']}-{emoji_part}-{anchor_title}".replace(' ', '')
        
        toc_line = f"{indent}- [{section['number']}.{emoji_part} {section['title']}]({anchor})"
        toc_lines.append(toc_line)
    
    toc_lines.extend(['', '---', ''])
    return '\n'.join(toc_lines)

def add_numbering(content, sections):
    """添加序号到章节标题"""
    lines = content.split('\n')
    
    for section in sections:
        if section['level'] == 2:
            lines[section['line']] = f"## {section['number']}. {section['emoji']} {section['title']}"
        elif section['level'] == 3:
            lines[section['line']] = f"### {section['number']} {section['title']}"
    
    return '\n'.join(lines)

def insert_toc(content, toc):
    """插入目录"""
    lines = content.split('\n')
    
    # 找到第一个## 之前插入
    for i, line in enumerate(lines):
        if line.startswith('## ') and '目录' not in line:
            lines.insert(i, toc)
            return '\n'.join(lines)
    
    return content

def format_document(file_path):
    """格式化单个文档"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # 如果已经有目录，先删除
        content = re.sub(r'## 📋 目录.*?---\n', '', content, flags=re.DOTALL)
        
        # 提取章节
        sections = extract_sections(content)
        
        if not sections:
            print(f"  ⚠️  没有找到章节")
            return False
        
        # 生成目录
        toc = generate_toc(sections)
        
        # 添加序号
        content = add_numbering(content, sections)
        
        # 插入目录
        content = insert_toc(content, toc)
        
        # 写回文件
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        
        return True
    
    except Exception as e:
        print(f"  ❌ 错误: {e}")
        return False

def main():
    print("=" * 60)
    print("批量格式化文档")
    print("=" * 60)
    print(f"\n总计: {len(DOCS_TO_UPDATE)} 个文档\n")
    
    success = 0
    failed = 0
    skipped = 0
    
    for i, doc_path in enumerate(DOCS_TO_UPDATE, 1):
        print(f"[{i}/{len(DOCS_TO_UPDATE)}] {doc_path}")
        
        if not os.path.exists(doc_path):
            print(f"  ⏭️  跳过(文件不存在)")
            skipped += 1
            continue
        
        if format_document(doc_path):
            print(f"  ✅ 完成")
            success += 1
        else:
            failed += 1
    
    print("\n" + "=" * 60)
    print(f"完成! 成功: {success}, 失败: {failed}, 跳过: {skipped}")
    print("=" * 60)

if __name__ == '__main__':
    main()

