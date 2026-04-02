# 知识跟踪目录

本目录包含自动生成的跟踪报告。

---

## 文件说明

| 文件 | 内容 | 更新频率 |
|------|------|----------|
| `go-releases.md` | Go 版本发布跟踪 | 每日 |
| `papers.md` | 学术论文跟踪 | 每周 |
| `weekly-digest-*.md` | 每周技术动态 | 每周 |
| `ecosystem-report-*.md` | 生态季度报告 | 每季 |

---

## 自动化工具

跟踪报告由以下工具自动生成：

- **Go 版本跟踪**: `scripts/knowledge-tracker/track_go_releases.py`
- **论文跟踪**: `scripts/knowledge-tracker/track_papers.py`
- **CI/CD**: `.github/workflows/knowledge-tracker.yml`

---

## 手动更新

如需手动更新：

```bash
# 安装依赖
pip install -r scripts/knowledge-tracker/requirements.txt

# 运行跟踪脚本
python scripts/knowledge-tracker/track_go_releases.py
python scripts/knowledge-tracker/track_papers.py
```

---

## 报告使用

跟踪报告用于：

1. **了解最新动态**: 快速浏览本周更新
2. **决策支持**: 评估是否升级/采用新技术
3. **内容灵感**: 发现值得深入分析的主题
4. **历史存档**: 追溯技术演进历史

---

*最后更新: [自动生成]*
