package framework

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CoverageReporter 覆盖率报告器
type CoverageReporter struct {
	outputDir string
	format    string
}

// CoverageReport 覆盖率报告
type CoverageReport struct {
	TotalCoverage float64
	Packages      []PackageCoverage
	Summary       CoverageSummary
}

// PackageCoverage 包覆盖率
type PackageCoverage struct {
	Package     string
	Coverage    float64
	Statements  int
	Covered     int
	Uncovered   int
	Files       []FileCoverage
}

// FileCoverage 文件覆盖率
type FileCoverage struct {
	File      string
	Coverage  float64
	Statements int
	Covered   int
	Uncovered int
}

// CoverageSummary 覆盖率摘要
type CoverageSummary struct {
	TotalStatements int
	TotalCovered    int
	TotalUncovered  int
	AverageCoverage float64
}

// NewCoverageReporter 创建覆盖率报告器
func NewCoverageReporter(outputDir, format string) *CoverageReporter {
	if outputDir == "" {
		outputDir = "coverage"
	}
	if format == "" {
		format = "html"
	}

	return &CoverageReporter{
		outputDir: outputDir,
		format:    format,
	}
}

// Generate 生成覆盖率报告
func (r *CoverageReporter) Generate(packages []string) (*CoverageReport, error) {
	// 创建输出目录
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成覆盖率文件
	coverageFile := filepath.Join(r.outputDir, "coverage.out")
	if err := r.generateCoverageFile(packages, coverageFile); err != nil {
		return nil, fmt.Errorf("failed to generate coverage file: %w", err)
	}

	// 解析覆盖率数据
	report, err := r.parseCoverageFile(coverageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coverage file: %w", err)
	}

	// 生成报告
	if err := r.generateReport(coverageFile, report); err != nil {
		return nil, fmt.Errorf("failed to generate report: %w", err)
	}

	return report, nil
}

// generateCoverageFile 生成覆盖率文件
func (r *CoverageReporter) generateCoverageFile(packages []string, outputFile string) error {
	args := []string{"test", "-coverprofile", outputFile, "-covermode=atomic"}
	args = append(args, packages...)

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// parseCoverageFile 解析覆盖率文件
func (r *CoverageReporter) parseCoverageFile(coverageFile string) (*CoverageReport, error) {
	// 使用 go tool cover 解析
	cmd := exec.Command("go", "tool", "cover", "-func", coverageFile)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to parse coverage: %w", err)
	}

	return r.parseCoverageOutput(string(output))
}

// parseCoverageOutput 解析覆盖率输出
func (r *CoverageReporter) parseCoverageOutput(output string) (*CoverageReport, error) {
	report := &CoverageReport{
		Packages: make([]PackageCoverage, 0),
	}

	lines := strings.Split(output, "\n")
	packageMap := make(map[string]*PackageCoverage)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "total:") {
			if strings.HasPrefix(line, "total:") {
				// 解析总覆盖率
				var totalCoverage float64
				fmt.Sscanf(line, "total:%*s%f%%", &totalCoverage)
				report.TotalCoverage = totalCoverage
			}
			continue
		}

		// 解析包覆盖率
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		filePath := parts[0]
		packageName := r.extractPackage(filePath)
		coverageStr := parts[len(parts)-1]
		coverageStr = strings.TrimSuffix(coverageStr, "%")

		var coverage float64
		var statements, covered int
		fmt.Sscanf(line, "%s %d %d %f%%", &filePath, &statements, &covered, &coverage)

		pkg, exists := packageMap[packageName]
		if !exists {
			pkg = &PackageCoverage{
				Package: packageName,
				Files:   make([]FileCoverage, 0),
			}
			packageMap[packageName] = pkg
		}

		fileCoverage := FileCoverage{
			File:       filePath,
			Coverage:   coverage,
			Statements: statements,
			Covered:    covered,
			Uncovered:  statements - covered,
		}

		pkg.Files = append(pkg.Files, fileCoverage)
		pkg.Statements += statements
		pkg.Covered += covered
		pkg.Uncovered += (statements - covered)
	}

	// 计算包覆盖率
	for _, pkg := range packageMap {
		if pkg.Statements > 0 {
			pkg.Coverage = float64(pkg.Covered) / float64(pkg.Statements) * 100
		}
		report.Packages = append(report.Packages, *pkg)
	}

	// 计算摘要
	report.Summary = r.calculateSummary(report.Packages)

	return report, nil
}

// extractPackage 提取包名
func (r *CoverageReporter) extractPackage(filePath string) string {
	parts := strings.Split(filePath, "/")
	if len(parts) < 2 {
		return "unknown"
	}

	// 查找 pkg 或 internal 目录
	for i, part := range parts {
		if part == "pkg" || part == "internal" {
			if i+1 < len(parts) {
				return strings.Join(parts[:i+2], "/")
			}
		}
	}

	return filepath.Dir(filePath)
}

// calculateSummary 计算摘要
func (r *CoverageReporter) calculateSummary(packages []PackageCoverage) CoverageSummary {
	summary := CoverageSummary{}

	for _, pkg := range packages {
		summary.TotalStatements += pkg.Statements
		summary.TotalCovered += pkg.Covered
		summary.TotalUncovered += pkg.Uncovered
	}

	if summary.TotalStatements > 0 {
		summary.AverageCoverage = float64(summary.TotalCovered) / float64(summary.TotalStatements) * 100
	}

	return summary
}

// generateReport 生成报告
func (r *CoverageReporter) generateReport(coverageFile string, report *CoverageReport) error {
	switch r.format {
	case "html":
		return r.generateHTMLReport(coverageFile)
	case "json":
		return r.generateJSONReport(report)
	case "text":
		return r.generateTextReport(report)
	default:
		return fmt.Errorf("unsupported format: %s", r.format)
	}
}

// generateHTMLReport 生成 HTML 报告
func (r *CoverageReporter) generateHTMLReport(coverageFile string) error {
	outputFile := filepath.Join(r.outputDir, "coverage.html")
	cmd := exec.Command("go", "tool", "cover", "-html", coverageFile, "-o", outputFile)
	return cmd.Run()
}

// generateJSONReport 生成 JSON 报告
func (r *CoverageReporter) generateJSONReport(report *CoverageReport) error {
	outputFile := filepath.Join(r.outputDir, "coverage.json")
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	return os.WriteFile(outputFile, data, 0644)
}

// generateTextReport 生成文本报告
func (r *CoverageReporter) generateTextReport(report *CoverageReport) error {
	outputFile := filepath.Join(r.outputDir, "coverage.txt")
	var builder strings.Builder

	builder.WriteString("Coverage Report\n")
	builder.WriteString("===============\n\n")
	builder.WriteString(fmt.Sprintf("Total Coverage: %.2f%%\n\n", report.TotalCoverage))

	builder.WriteString("Summary:\n")
	builder.WriteString(fmt.Sprintf("  Total Statements: %d\n", report.Summary.TotalStatements))
	builder.WriteString(fmt.Sprintf("  Total Covered: %d\n", report.Summary.TotalCovered))
	builder.WriteString(fmt.Sprintf("  Total Uncovered: %d\n", report.Summary.TotalUncovered))
	builder.WriteString(fmt.Sprintf("  Average Coverage: %.2f%%\n\n", report.Summary.AverageCoverage))

	builder.WriteString("Packages:\n")
	for _, pkg := range report.Packages {
		builder.WriteString(fmt.Sprintf("  %s: %.2f%% (%d/%d)\n",
			pkg.Package, pkg.Coverage, pkg.Covered, pkg.Statements))
	}

	return os.WriteFile(outputFile, []byte(builder.String()), 0644)
}

// CheckThreshold 检查覆盖率阈值
func (r *CoverageReporter) CheckThreshold(report *CoverageReport, threshold float64) error {
	if report.TotalCoverage < threshold {
		return fmt.Errorf("coverage %.2f%% is below threshold %.2f%%",
			report.TotalCoverage, threshold)
	}
	return nil
}
