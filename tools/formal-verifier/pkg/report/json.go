package report

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/formal-verifier/pkg/project"
)

// JSONReport JSON格式报告生成器
type JSONReport struct {
	result *project.AnalysisResult
}

// NewJSONReport 创建JSON报告生成器
func NewJSONReport(result *project.AnalysisResult) *JSONReport {
	return &JSONReport{result: result}
}

// Generate 生成JSON报告
func (j *JSONReport) Generate(output string) error {
	// 准备JSON数据
	data := map[string]interface{}{
		"summary":  j.result.Summary,
		"stats":    j.result.Stats,
		"issues":   j.result.Issues,
		"errors":   j.result.GetIssuesBySeverity("error"),
		"warnings": j.result.GetIssuesBySeverity("warning"),
		"infos":    j.result.GetIssuesBySeverity("info"),
	}

	// 创建输出文件
	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// 编码JSON（格式化输出）
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
