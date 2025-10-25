package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	config := Default()

	// 验证默认值
	if config.Project.RootDir != "." {
		t.Errorf("Expected default root dir '.', got %s", config.Project.RootDir)
	}

	if !config.Project.Recursive {
		t.Error("Expected recursive to be true by default")
	}

	if config.Analysis.Timeout != 300 {
		t.Errorf("Expected default timeout 300, got %d", config.Analysis.Timeout)
	}

	if config.Report.Format != "text" {
		t.Errorf("Expected default format 'text', got %s", config.Report.Format)
	}

	if !config.Rules.Concurrency.Enabled {
		t.Error("Expected concurrency rules to be enabled by default")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name:      "Valid default config",
			config:    Default(),
			wantError: false,
		},
		{
			name: "Empty root dir",
			config: &Config{
				Project: ProjectConfig{RootDir: ""},
			},
			wantError: true,
		},
		{
			name: "Negative workers",
			config: func() *Config {
				c := Default()
				c.Analysis.Workers = -1
				return c
			}(),
			wantError: true,
		},
		{
			name: "Invalid report format",
			config: func() *Config {
				c := Default()
				c.Report.Format = "invalid"
				return c
			}(),
			wantError: true,
		},
		{
			name: "Invalid cyclomatic threshold",
			config: func() *Config {
				c := Default()
				c.Rules.Complexity.CyclomaticThreshold = 0
				return c
			}(),
			wantError: true,
		},
		{
			name: "Invalid quality score",
			config: func() *Config {
				c := Default()
				c.Output.MinQualityScore = 150
				return c
			}(),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestLoadAndSave(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// 创建测试配置
	config := Default()
	config.Project.RootDir = "/test/path"
	config.Analysis.Workers = 4
	config.Report.Format = "json"
	config.Rules.Complexity.CyclomaticThreshold = 15

	// 保存配置
	if err := config.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// 加载配置
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 验证配置内容
	if loaded.Project.RootDir != "/test/path" {
		t.Errorf("Expected root dir '/test/path', got %s", loaded.Project.RootDir)
	}
	if loaded.Analysis.Workers != 4 {
		t.Errorf("Expected 4 workers, got %d", loaded.Analysis.Workers)
	}
	if loaded.Report.Format != "json" {
		t.Errorf("Expected format 'json', got %s", loaded.Report.Format)
	}
	if loaded.Rules.Complexity.CyclomaticThreshold != 15 {
		t.Errorf("Expected threshold 15, got %d", loaded.Rules.Complexity.CyclomaticThreshold)
	}
}

func TestLoadOrDefault(t *testing.T) {
	// 测试空路径
	config := LoadOrDefault("")
	if config == nil {
		t.Fatal("Expected default config, got nil")
	}

	// 测试不存在的文件
	config = LoadOrDefault("/nonexistent/config.yaml")
	if config == nil {
		t.Fatal("Expected default config for nonexistent file, got nil")
	}

	// 测试有效文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	testConfig := Default()
	testConfig.Project.RootDir = "/custom/path"
	if err := testConfig.Save(configPath); err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}

	loaded := LoadOrDefault(configPath)
	if loaded.Project.RootDir != "/custom/path" {
		t.Errorf("Expected custom path, got %s", loaded.Project.RootDir)
	}
}

func TestMerge(t *testing.T) {
	base := Default()
	other := &Config{
		Project: ProjectConfig{
			RootDir: "/new/path",
		},
		Analysis: AnalysisConfig{
			Workers: 8,
		},
		Report: ReportConfig{
			Format: "html",
		},
		Rules: RulesConfig{
			Complexity: ComplexityRules{
				CyclomaticThreshold: 20,
			},
		},
	}

	base.Merge(other)

	if base.Project.RootDir != "/new/path" {
		t.Errorf("Expected merged root dir '/new/path', got %s", base.Project.RootDir)
	}
	if base.Analysis.Workers != 8 {
		t.Errorf("Expected merged workers 8, got %d", base.Analysis.Workers)
	}
	if base.Report.Format != "html" {
		t.Errorf("Expected merged format 'html', got %s", base.Report.Format)
	}
	if base.Rules.Complexity.CyclomaticThreshold != 20 {
		t.Errorf("Expected merged threshold 20, got %d", base.Rules.Complexity.CyclomaticThreshold)
	}

	// 验证其他值未被覆盖
	if !base.Project.Recursive {
		t.Error("Expected recursive to remain true")
	}
}

func TestMergeWithNil(t *testing.T) {
	config := Default()
	originalRootDir := config.Project.RootDir

	config.Merge(nil)

	if config.Project.RootDir != originalRootDir {
		t.Error("Config should not change when merging with nil")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	// 写入无效的YAML
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content: [[["), 0644); err != nil {
		t.Fatalf("Failed to write invalid yaml: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error when loading invalid YAML, got nil")
	}
}

func TestComplexityRulesValidation(t *testing.T) {
	tests := []struct {
		name      string
		threshold int
		wantError bool
	}{
		{"Valid threshold", 10, false},
		{"Zero threshold", 0, true},
		{"Negative threshold", -5, true},
		{"High threshold", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Default()
			config.Rules.Complexity.CyclomaticThreshold = tt.threshold
			err := config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestReportFormatValidation(t *testing.T) {
	formats := []struct {
		format    string
		wantError bool
	}{
		{"text", false},
		{"html", false},
		{"json", false},
		{"markdown", false},
		{"xml", true},
		{"pdf", true},
		{"", true},
	}

	for _, tt := range formats {
		t.Run(tt.format, func(t *testing.T) {
			config := Default()
			config.Report.Format = tt.format
			err := config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() with format '%s' error = %v, wantError %v", tt.format, err, tt.wantError)
			}
		})
	}
}
