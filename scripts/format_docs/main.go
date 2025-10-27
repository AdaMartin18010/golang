package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// 需要更新的文档列表
var docsToUpdate = []string{
	// testing
	"docs/practices/testing/02-表格驱动测试.md",
	"docs/practices/testing/03-集成测试.md",
	"docs/practices/testing/04-性能测试.md",
	"docs/practices/testing/05-测试覆盖率.md",
	"docs/practices/testing/06-Mock与Stub.md",
	"docs/practices/testing/07-测试最佳实践.md",
	"docs/practices/testing/08-常见问题与技巧.md",
	// deployment
	"docs/practices/deployment/01-部署概览.md",
	"docs/practices/deployment/02-Docker部署.md",
	"docs/practices/deployment/03-Kubernetes部署.md",
	"docs/practices/deployment/04-CI-CD流程.md",
	"docs/practices/deployment/05-监控与日志.md",
	"docs/practices/deployment/06-滚动更新.md",
	"docs/practices/deployment/07-生产环境最佳实践.md",
	// distributed
	"docs/advanced/distributed/01-分布式系统基础.md",
	"docs/advanced/distributed/02-服务注册与发现.md",
	"docs/advanced/distributed/03-分布式一致性.md",
	"docs/advanced/distributed/04-分布式锁.md",
	"docs/advanced/distributed/05-分布式事务.md",
	"docs/advanced/distributed/06-负载均衡.md",
	// security
	"docs/advanced/security/01-Web安全基础.md",
	"docs/advanced/security/02-身份认证.md",
	"docs/advanced/security/03-授权机制.md",
	"docs/advanced/security/04-数据保护.md",
	"docs/advanced/security/05-安全审计.md",
	"docs/advanced/security/06-最佳实践.md",
	// 其他目录...
}

func main() {
	fmt.Println("开始批量格式化文档...")
	fmt.Printf("总计: %d 个文档\n\n", len(docsToUpdate))

	success := 0
	failed := 0

	for i, docPath := range docsToUpdate {
		fmt.Printf("[%d/%d] 处理: %s\n", i+1, len(docsToUpdate), docPath)

		if err := formatDocument(docPath); err != nil {
			fmt.Printf("  ❌ 失败: %v\n", err)
			failed++
		} else {
			fmt.Printf("  ✅ 完成\n")
			success++
		}
	}

	fmt.Printf("\n完成! 成功: %d, 失败: %d\n", success, failed)
}

func formatDocument(filePath string) error {
	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	// 解析文档结构
	sections := parseDocument(lines)

	// 生成目录
	toc := generateTOC(sections)

	// 添加序号
	formatted := addNumbering(lines, sections)

	// 插入目录
	final := insertTOC(formatted, toc)

	// 写回文件
	return os.WriteFile(filePath, []byte(strings.Join(final, "\n")), 0644)
}

type Section struct {
	Level  int
	Title  string
	Line   int
	Number string
	Emoji  string
}

func parseDocument(lines []string) []Section {
	sections := []Section{}
	h2Pattern := regexp.MustCompile(`^##\s+(.+)$`)
	h3Pattern := regexp.MustCompile(`^###\s+(.+)$`)

	h2Count := 0
	h3Counts := make(map[int]int)

	for i, line := range lines {
		if matches := h2Pattern.FindStringSubmatch(line); matches != nil {
			h2Count++
			h3Counts[h2Count] = 0

			title := matches[1]
			emoji := extractEmoji(title)
			cleanTitle := strings.TrimSpace(strings.TrimPrefix(title, emoji))

			sections = append(sections, Section{
				Level:  2,
				Title:  cleanTitle,
				Line:   i,
				Number: fmt.Sprintf("%d", h2Count),
				Emoji:  emoji,
			})
		} else if matches := h3Pattern.FindStringSubmatch(line); matches != nil {
			if h2Count > 0 {
				h3Counts[h2Count]++

				title := matches[1]
				emoji := extractEmoji(title)
				cleanTitle := strings.TrimSpace(strings.TrimPrefix(title, emoji))

				sections = append(sections, Section{
					Level:  3,
					Title:  cleanTitle,
					Line:   i,
					Number: fmt.Sprintf("%d.%d", h2Count, h3Counts[h2Count]),
					Emoji:  emoji,
				})
			}
		}
	}

	return sections
}

func extractEmoji(text string) string {
	emojiPattern := regexp.MustCompile(`^([^\w\s]+)\s+`)
	if matches := emojiPattern.FindStringSubmatch(text); matches != nil {
		return matches[1]
	}
	return ""
}

func generateTOC(sections []Section) []string {
	toc := []string{
		"## 📋 目录",
		"",
	}

	for _, s := range sections {
		indent := strings.Repeat("  ", s.Level-2)
		anchor := strings.ToLower(strings.ReplaceAll(s.Title, " ", "-"))
		anchor = regexp.MustCompile(`[^\w\-一-龥]`).ReplaceAllString(anchor, "")

		line := fmt.Sprintf("%s- [%s. %s %s](#%s-%s-%s)",
			indent, s.Number, s.Emoji, s.Title, s.Number, s.Emoji, anchor)
		toc = append(toc, line)
	}

	toc = append(toc, "", "---", "")
	return toc
}

func addNumbering(lines []string, sections []Section) []string {
	result := make([]string, len(lines))
	copy(result, lines)

	for _, s := range sections {
		if s.Level == 2 {
			result[s.Line] = fmt.Sprintf("## %s. %s %s", s.Number, s.Emoji, s.Title)
		} else if s.Level == 3 {
			result[s.Line] = fmt.Sprintf("### %s %s", s.Number, s.Title)
		}
	}

	return result
}

func insertTOC(lines []string, toc []string) []string {
	// 找到第一个 ## 之前插入TOC
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") && !strings.Contains(line, "目录") {
			result := make([]string, 0, len(lines)+len(toc))
			result = append(result, lines[:i]...)
			result = append(result, toc...)
			result = append(result, lines[i:]...)
			return result
		}
	}
	return lines
}
