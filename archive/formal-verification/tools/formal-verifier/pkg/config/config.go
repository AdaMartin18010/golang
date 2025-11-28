// Package config 提供配置管理功能
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 表示FV工具的完整配置
type Config struct {
	// 项目配置
	Project ProjectConfig `yaml:"project"`

	// 分析配置
	Analysis AnalysisConfig `yaml:"analysis"`

	// 报告配置
	Report ReportConfig `yaml:"report"`

	// 规则配置
	Rules RulesConfig `yaml:"rules"`

	// 输出配置
	Output OutputConfig `yaml:"output"`
}

// ProjectConfig 项目相关配置
type ProjectConfig struct {
	// 项目根目录
	RootDir string `yaml:"root_dir"`

	// 是否递归扫描
	Recursive bool `yaml:"recursive"`

	// 是否包含测试文件
	IncludeTests bool `yaml:"include_tests"`

	// 排除模式列表
	ExcludePatterns []string `yaml:"exclude_patterns"`

	// Go模块路径
	ModulePath string `yaml:"module_path"`
}

// AnalysisConfig 分析相关配置
type AnalysisConfig struct {
	// 是否启用并发分析
	Concurrent bool `yaml:"concurrent"`

	// 并发工作线程数（0表示使用CPU核心数）
	Workers int `yaml:"workers"`

	// 是否启用详细模式
	Verbose bool `yaml:"verbose"`

	// 超时时间（秒）
	Timeout int `yaml:"timeout"`

	// 最大文件大小（KB）
	MaxFileSize int `yaml:"max_file_size"`
}

// ReportConfig 报告相关配置
type ReportConfig struct {
	// 报告格式（text, html, json, markdown）
	Format string `yaml:"format"`

	// 输出文件路径
	OutputPath string `yaml:"output_path"`

	// 是否在浏览器中打开HTML报告
	OpenInBrowser bool `yaml:"open_in_browser"`

	// 报告标题
	Title string `yaml:"title"`

	// 报告作者
	Author string `yaml:"author"`
}

// RulesConfig 规则相关配置
type RulesConfig struct {
	// 并发规则
	Concurrency ConcurrencyRules `yaml:"concurrency"`

	// 类型规则
	Type TypeRules `yaml:"type"`

	// 复杂度规则
	Complexity ComplexityRules `yaml:"complexity"`

	// 性能规则
	Performance PerformanceRules `yaml:"performance"`
}

// ConcurrencyRules 并发检查规则
type ConcurrencyRules struct {
	// 是否启用
	Enabled bool `yaml:"enabled"`

	// 检查goroutine泄漏
	CheckGoroutineLeak bool `yaml:"check_goroutine_leak"`

	// 检查数据竞争
	CheckDataRace bool `yaml:"check_data_race"`

	// 检查死锁
	CheckDeadlock bool `yaml:"check_deadlock"`

	// 检查channel使用
	CheckChannel bool `yaml:"check_channel"`
}

// TypeRules 类型检查规则
type TypeRules struct {
	// 是否启用
	Enabled bool `yaml:"enabled"`

	// 检查nil指针
	CheckNilPointer bool `yaml:"check_nil_pointer"`

	// 检查类型断言
	CheckTypeAssertion bool `yaml:"check_type_assertion"`

	// 检查接口实现
	CheckInterface bool `yaml:"check_interface"`
}

// ComplexityRules 复杂度检查规则
type ComplexityRules struct {
	// 是否启用
	Enabled bool `yaml:"enabled"`

	// 圈复杂度阈值
	CyclomaticThreshold int `yaml:"cyclomatic_threshold"`

	// 认知复杂度阈值
	CognitiveThreshold int `yaml:"cognitive_threshold"`

	// 函数行数阈值
	MaxFunctionLines int `yaml:"max_function_lines"`

	// 参数数量阈值
	MaxParameters int `yaml:"max_parameters"`
}

// PerformanceRules 性能检查规则
type PerformanceRules struct {
	// 是否启用
	Enabled bool `yaml:"enabled"`

	// 检查内存分配
	CheckAllocation bool `yaml:"check_allocation"`

	// 检查字符串拼接
	CheckStringConcat bool `yaml:"check_string_concat"`

	// 检查循环优化
	CheckLoopOptimization bool `yaml:"check_loop_optimization"`
}

// OutputConfig 输出相关配置
type OutputConfig struct {
	// 是否使用彩色输出
	ColorOutput bool `yaml:"color_output"`

	// 是否显示进度条
	ShowProgress bool `yaml:"show_progress"`

	// 是否显示统计信息
	ShowStats bool `yaml:"show_stats"`

	// 失败时退出码
	FailOnError bool `yaml:"fail_on_error"`

	// 最小质量分数（低于此分数视为失败）
	MinQualityScore int `yaml:"min_quality_score"`
}

// Default 返回默认配置
func Default() *Config {
	return &Config{
		Project: ProjectConfig{
			RootDir:         ".",
			Recursive:       true,
			IncludeTests:    false,
			ExcludePatterns: []string{"vendor", "testdata", ".git"},
			ModulePath:      "",
		},
		Analysis: AnalysisConfig{
			Concurrent:  true,
			Workers:     0, // 使用CPU核心数
			Verbose:     false,
			Timeout:     300,  // 5分钟
			MaxFileSize: 1024, // 1MB
		},
		Report: ReportConfig{
			Format:        "text",
			OutputPath:    "",
			OpenInBrowser: false,
			Title:         "形式化验证报告",
			Author:        "",
		},
		Rules: RulesConfig{
			Concurrency: ConcurrencyRules{
				Enabled:            true,
				CheckGoroutineLeak: true,
				CheckDataRace:      true,
				CheckDeadlock:      true,
				CheckChannel:       true,
			},
			Type: TypeRules{
				Enabled:            true,
				CheckNilPointer:    true,
				CheckTypeAssertion: true,
				CheckInterface:     true,
			},
			Complexity: ComplexityRules{
				Enabled:             true,
				CyclomaticThreshold: 10,
				CognitiveThreshold:  15,
				MaxFunctionLines:    50,
				MaxParameters:       5,
			},
			Performance: PerformanceRules{
				Enabled:               true,
				CheckAllocation:       true,
				CheckStringConcat:     true,
				CheckLoopOptimization: true,
			},
		},
		Output: OutputConfig{
			ColorOutput:     true,
			ShowProgress:    true,
			ShowStats:       true,
			FailOnError:     false,
			MinQualityScore: 0,
		},
	}
}

// Load 从YAML文件加载配置
func Load(path string) (*Config, error) {
	// 读取文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析YAML
	config := Default()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}

// LoadOrDefault 加载配置文件，如果不存在则返回默认配置
func LoadOrDefault(path string) *Config {
	if path == "" {
		return Default()
	}

	config, err := Load(path)
	if err != nil {
		// 配置文件不存在或无效，返回默认配置
		return Default()
	}

	return config
}

// Save 保存配置到YAML文件
func (c *Config) Save(path string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 序列化为YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate 验证配置有效性
func (c *Config) Validate() error {
	// 验证项目配置
	if c.Project.RootDir == "" {
		return fmt.Errorf("project.root_dir cannot be empty")
	}

	// 验证分析配置
	if c.Analysis.Workers < 0 {
		return fmt.Errorf("analysis.workers must be >= 0")
	}
	if c.Analysis.Timeout < 0 {
		return fmt.Errorf("analysis.timeout must be >= 0")
	}
	if c.Analysis.MaxFileSize < 0 {
		return fmt.Errorf("analysis.max_file_size must be >= 0")
	}

	// 验证报告配置
	validFormats := map[string]bool{
		"text":     true,
		"html":     true,
		"json":     true,
		"markdown": true,
	}
	if !validFormats[c.Report.Format] {
		return fmt.Errorf("invalid report.format: %s (must be text, html, json, or markdown)", c.Report.Format)
	}

	// 验证复杂度阈值
	if c.Rules.Complexity.CyclomaticThreshold < 1 {
		return fmt.Errorf("rules.complexity.cyclomatic_threshold must be >= 1")
	}
	if c.Rules.Complexity.CognitiveThreshold < 1 {
		return fmt.Errorf("rules.complexity.cognitive_threshold must be >= 1")
	}
	if c.Rules.Complexity.MaxFunctionLines < 1 {
		return fmt.Errorf("rules.complexity.max_function_lines must be >= 1")
	}
	if c.Rules.Complexity.MaxParameters < 0 {
		return fmt.Errorf("rules.complexity.max_parameters must be >= 0")
	}

	// 验证输出配置
	if c.Output.MinQualityScore < 0 || c.Output.MinQualityScore > 100 {
		return fmt.Errorf("output.min_quality_score must be between 0 and 100")
	}

	return nil
}

// Merge 合并配置（other中的非零值会覆盖当前配置）
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	// 合并项目配置
	if other.Project.RootDir != "" && other.Project.RootDir != "." {
		c.Project.RootDir = other.Project.RootDir
	}
	if len(other.Project.ExcludePatterns) > 0 {
		c.Project.ExcludePatterns = other.Project.ExcludePatterns
	}
	if other.Project.ModulePath != "" {
		c.Project.ModulePath = other.Project.ModulePath
	}

	// 合并分析配置
	if other.Analysis.Workers > 0 {
		c.Analysis.Workers = other.Analysis.Workers
	}
	if other.Analysis.Timeout > 0 {
		c.Analysis.Timeout = other.Analysis.Timeout
	}
	if other.Analysis.MaxFileSize > 0 {
		c.Analysis.MaxFileSize = other.Analysis.MaxFileSize
	}

	// 合并报告配置
	if other.Report.Format != "" && other.Report.Format != "text" {
		c.Report.Format = other.Report.Format
	}
	if other.Report.OutputPath != "" {
		c.Report.OutputPath = other.Report.OutputPath
	}
	if other.Report.Title != "" && other.Report.Title != "形式化验证报告" {
		c.Report.Title = other.Report.Title
	}

	// 合并复杂度规则
	if other.Rules.Complexity.CyclomaticThreshold > 0 {
		c.Rules.Complexity.CyclomaticThreshold = other.Rules.Complexity.CyclomaticThreshold
	}
	if other.Rules.Complexity.CognitiveThreshold > 0 {
		c.Rules.Complexity.CognitiveThreshold = other.Rules.Complexity.CognitiveThreshold
	}
	if other.Rules.Complexity.MaxFunctionLines > 0 {
		c.Rules.Complexity.MaxFunctionLines = other.Rules.Complexity.MaxFunctionLines
	}

	// 合并输出配置
	if other.Output.MinQualityScore > 0 {
		c.Output.MinQualityScore = other.Output.MinQualityScore
	}
}
