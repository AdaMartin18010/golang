package report

import (
	"fmt"
	"html/template"
	"os"

	"github.com/your-org/formal-verifier/pkg/project"
)

// HTMLReport HTML格式报告生成器
type HTMLReport struct {
	result *project.AnalysisResult
}

// NewHTMLReport 创建HTML报告生成器
func NewHTMLReport(result *project.AnalysisResult) *HTMLReport {
	return &HTMLReport{result: result}
}

// Generate 生成HTML报告
func (h *HTMLReport) Generate(output string) error {
	// 准备模板数据
	data := h.prepareTemplateData()

	// 解析模板
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"severityClass": severityClass,
		"categoryIcon":  categoryIcon,
		"scoreClass":    scoreClass,
		"scoreLabel":    scoreLabel,
	}).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// 创建输出文件
	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// 渲染模板
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// prepareTemplateData 准备模板数据
func (h *HTMLReport) prepareTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"Stats":       h.result.Stats,
		"Issues":      h.result.Issues,
		"Errors":      h.result.GetIssuesBySeverity("error"),
		"Warnings":    h.result.GetIssuesBySeverity("warning"),
		"Infos":       h.result.GetIssuesBySeverity("info"),
		"Summary":     h.result.Summary,
		"HasErrors":   h.result.HasErrors(),
		"HasWarnings": h.result.HasWarnings(),
	}
}

// severityClass 返回严重程度对应的CSS类
func severityClass(severity string) string {
	switch severity {
	case "error":
		return "severity-error"
	case "warning":
		return "severity-warning"
	case "info":
		return "severity-info"
	default:
		return ""
	}
}

// categoryIcon 返回类别对应的图标
func categoryIcon(category string) string {
	icons := map[string]string{
		"concurrency":  "⚡",
		"type":         "🔤",
		"dataflow":     "📊",
		"optimization": "⚙️",
		"complexity":   "🔀",
		"syntax":       "📝",
	}

	if icon, ok := icons[category]; ok {
		return icon
	}
	return "📌"
}

// scoreClass 返回质量评分对应的CSS类
func scoreClass(score int) string {
	if score >= 90 {
		return "score-excellent"
	} else if score >= 70 {
		return "score-good"
	} else if score >= 50 {
		return "score-fair"
	}
	return "score-poor"
}

// scoreLabel 返回质量评分对应的标签
func scoreLabel(score int) string {
	if score >= 90 {
		return "优秀"
	} else if score >= 70 {
		return "良好"
	} else if score >= 50 {
		return "一般"
	}
	return "较差"
}

// htmlTemplate HTML报告模板
const htmlTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Formal Verifier - 项目分析报告</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }

        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
        }

        .header p {
            font-size: 1.1em;
            opacity: 0.9;
        }

        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            padding: 40px;
            background: #f8f9fa;
        }

        .stat-card {
            background: white;
            padding: 25px;
            border-radius: 10px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            text-align: center;
            transition: transform 0.3s ease;
        }

        .stat-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
        }

        .stat-card h3 {
            font-size: 0.9em;
            color: #666;
            margin-bottom: 10px;
            text-transform: uppercase;
            letter-spacing: 1px;
        }

        .stat-card .value {
            font-size: 2.5em;
            font-weight: bold;
            color: #667eea;
        }

        .score-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }

        .score-card .value {
            color: white;
            font-size: 3em;
        }

        .score-excellent { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
        .score-good { background: linear-gradient(135deg, #3498db 0%, #2980b9 100%); }
        .score-fair { background: linear-gradient(135deg, #f39c12 0%, #e67e22 100%); }
        .score-poor { background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%); }

        .content {
            padding: 40px;
        }

        .section {
            margin-bottom: 40px;
        }

        .section h2 {
            font-size: 1.8em;
            color: #667eea;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 3px solid #667eea;
        }

        .issue-list {
            display: grid;
            gap: 15px;
        }

        .issue-card {
            background: white;
            border-left: 4px solid #ddd;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            transition: all 0.3s ease;
        }

        .issue-card:hover {
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
            transform: translateX(5px);
        }

        .severity-error {
            border-left-color: #e74c3c;
            background: #fff5f5;
        }

        .severity-warning {
            border-left-color: #f39c12;
            background: #fffbf5;
        }

        .severity-info {
            border-left-color: #3498db;
            background: #f5f9ff;
        }

        .issue-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 10px;
        }

        .issue-icon {
            font-size: 1.5em;
            margin-right: 10px;
        }

        .issue-location {
            font-family: 'Courier New', monospace;
            font-size: 0.9em;
            color: #666;
            background: #f8f9fa;
            padding: 4px 8px;
            border-radius: 4px;
        }

        .issue-category {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.85em;
            font-weight: 600;
            background: #667eea;
            color: white;
        }

        .issue-message {
            font-size: 1.05em;
            margin: 15px 0;
            color: #333;
        }

        .issue-suggestion {
            background: #f0f7ff;
            border-left: 3px solid #3498db;
            padding: 12px 15px;
            margin-top: 15px;
            border-radius: 4px;
        }

        .issue-suggestion::before {
            content: '💡 建议: ';
            font-weight: bold;
            color: #3498db;
        }

        .chart {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }

        .chart-item {
            background: white;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }

        .chart-item h4 {
            color: #667eea;
            margin-bottom: 15px;
        }

        .progress-bar {
            background: #e9ecef;
            height: 30px;
            border-radius: 15px;
            overflow: hidden;
            margin: 10px 0;
        }

        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            font-size: 0.85em;
            transition: width 0.3s ease;
        }

        .footer {
            background: #f8f9fa;
            padding: 30px;
            text-align: center;
            color: #666;
            border-top: 1px solid #dee2e6;
        }

        .footer p {
            margin: 5px 0;
        }

        .footer a {
            color: #667eea;
            text-decoration: none;
            font-weight: 600;
        }

        .footer a:hover {
            text-decoration: underline;
        }

        @media print {
            body {
                background: white;
                padding: 0;
            }

            .container {
                box-shadow: none;
            }

            .issue-card:hover {
                transform: none;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 Go Formal Verifier</h1>
            <p>项目形式化分析报告</p>
        </div>

        <div class="summary">
            <div class="stat-card">
                <h3>文件数</h3>
                <div class="value">{{.Stats.TotalFiles}}</div>
            </div>
            <div class="stat-card">
                <h3>代码行数</h3>
                <div class="value">{{.Stats.TotalLines}}</div>
            </div>
            <div class="stat-card">
                <h3>总问题</h3>
                <div class="value">{{.Stats.TotalIssues}}</div>
            </div>
            <div class="stat-card score-card {{scoreClass .Stats.QualityScore}}">
                <h3>质量评分</h3>
                <div class="value">{{.Stats.QualityScore}}</div>
                <p>{{scoreLabel .Stats.QualityScore}}</p>
            </div>
        </div>

        <div class="content">
            <div class="section">
                <h2>📊 统计信息</h2>
                <div class="chart">
                    <div class="chart-item">
                        <h4>按严重程度分类</h4>
                        {{if gt .Stats.ErrorCount 0}}
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: 100%; background: #e74c3c;">
                                ❌ 错误: {{.Stats.ErrorCount}}
                            </div>
                        </div>
                        {{end}}
                        {{if gt .Stats.WarningCount 0}}
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: 100%; background: #f39c12;">
                                ⚠️  警告: {{.Stats.WarningCount}}
                            </div>
                        </div>
                        {{end}}
                        {{if gt .Stats.InfoCount 0}}
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: 100%; background: #3498db;">
                                ℹ️  提示: {{.Stats.InfoCount}}
                            </div>
                        </div>
                        {{end}}
                    </div>

                    <div class="chart-item">
                        <h4>按问题类别分类</h4>
                        {{if gt .Stats.ConcurrencyIssues 0}}
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: 100%; background: #667eea;">
                                ⚡ 并发: {{.Stats.ConcurrencyIssues}}
                            </div>
                        </div>
                        {{end}}
                        {{if gt .Stats.TypeIssues 0}}
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: 100%; background: #764ba2;">
                                🔤 类型: {{.Stats.TypeIssues}}
                            </div>
                        </div>
                        {{end}}
                        {{if gt .Stats.DataFlowIssues 0}}
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: 100%; background: #11998e;">
                                📊 数据流: {{.Stats.DataFlowIssues}}
                            </div>
                        </div>
                        {{end}}
                        {{if gt .Stats.OptimizationIssues 0}}
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: 100%; background: #f39c12;">
                                ⚙️  优化: {{.Stats.OptimizationIssues}}
                            </div>
                        </div>
                        {{end}}
                    </div>
                </div>
            </div>

            {{if .HasErrors}}
            <div class="section">
                <h2>❌ 错误</h2>
                <div class="issue-list">
                    {{range .Errors}}
                    <div class="issue-card {{severityClass .Severity}}">
                        <div class="issue-header">
                            <span class="issue-icon">{{categoryIcon .Category}}</span>
                            <span class="issue-category">{{.Category}}</span>
                            <span class="issue-location">{{.File}}:{{.Line}}:{{.Column}}</span>
                        </div>
                        <div class="issue-message">{{.Message}}</div>
                        {{if .Suggestion}}
                        <div class="issue-suggestion">{{.Suggestion}}</div>
                        {{end}}
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}

            {{if .HasWarnings}}
            <div class="section">
                <h2>⚠️ 警告</h2>
                <div class="issue-list">
                    {{range .Warnings}}
                    <div class="issue-card {{severityClass .Severity}}">
                        <div class="issue-header">
                            <span class="issue-icon">{{categoryIcon .Category}}</span>
                            <span class="issue-category">{{.Category}}</span>
                            <span class="issue-location">{{.File}}:{{.Line}}:{{.Column}}</span>
                        </div>
                        <div class="issue-message">{{.Message}}</div>
                        {{if .Suggestion}}
                        <div class="issue-suggestion">{{.Suggestion}}</div>
                        {{end}}
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}

            {{if eq .Stats.TotalIssues 0}}
            <div class="section">
                <div style="text-align: center; padding: 60px 20px;">
                    <div style="font-size: 4em; margin-bottom: 20px;">✅</div>
                    <h2 style="color: #11998e; border: none;">代码质量优秀！</h2>
                    <p style="color: #666; font-size: 1.1em; margin-top: 10px;">没有发现任何问题</p>
                </div>
            </div>
            {{end}}
        </div>

        <div class="footer">
            <p><strong>Go Formal Verifier</strong> - 基于 Go 1.25.3 形式化理论体系</p>
            <p>文档位置: docs/01-语言基础/00-Go-1.25.3形式化理论体系/</p>
            <p><a href="https://github.com/your-org/formal-verifier">GitHub</a> | 
               <a href="https://github.com/your-org/formal-verifier/docs">文档</a></p>
        </div>
    </div>
</body>
</html>`
