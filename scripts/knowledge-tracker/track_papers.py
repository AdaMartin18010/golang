#!/usr/bin/env python3
"""
学术论文跟踪器
自动搜索 Go 语言相关学术论文
"""

import arxiv
from datetime import datetime, timedelta

# 搜索关键词
KEYWORDS = [
    'Go language',
    'Golang semantics',
    'Go concurrency',
    'Featherweight Go',
    'Go type system',
    'Go memory model',
]

def fetch_recent_papers(days=7):
    """获取最近发表的论文"""
    papers = []
    cutoff_date = datetime.now() - timedelta(days=days)
    
    for keyword in KEYWORDS:
        print(f"Searching: {keyword}")
        
        search = arxiv.Search(
            query=keyword,
            max_results=5,
            sort_by=arxiv.SortCriterion.SubmittedDate
        )
        
        for result in search.results():
            if result.published > cutoff_date:
                papers.append({
                    'title': result.title,
                    'authors': [a.name for a in result.authors],
                    'published': result.published,
                    'summary': result.summary[:200] + "...",
                    'url': result.entry_id,
                    'keyword': keyword
                })
    
    return papers

def generate_report(papers):
    """生成论文跟踪报告"""
    if not papers:
        return "# 学术论文跟踪报告\n\n本周无新论文。\n"
    
    report = f"""# 学术论文跟踪报告

生成时间: {datetime.now().isoformat()}
本周新论文: {len(papers)} 篇

"""
    
    for i, paper in enumerate(papers, 1):
        report += f"""## {i}. {paper['title']}

- **作者**: {', '.join(paper['authors'])}
- **发布时间**: {paper['published'].strftime('%Y-%m-%d')}
- **关键词**: {paper['keyword']}
- **摘要**: {paper['summary']}
- **链接**: {paper['url']}

---

"""
    
    report += """
## 建议行动

- [ ] 阅读相关论文
- [ ] 评估对项目的影响
- [ ] 更新技术文档
- [ ] 分享社区讨论

---

*自动生成的学术跟踪报告*
"""
    
    return report

def main():
    papers = fetch_recent_papers()
    report = generate_report(papers)
    
    output_file = "docs/tracking/papers.md"
    with open(output_file, "w", encoding="utf-8") as f:
        f.write(report)
    
    print(f"Found {len(papers)} new papers")
    print(f"Report saved to {output_file}")

if __name__ == "__main__":
    main()
