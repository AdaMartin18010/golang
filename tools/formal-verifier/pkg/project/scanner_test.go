package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewScanner(t *testing.T) {
	scanner := NewScanner(".")
	if scanner == nil {
		t.Fatal("NewScanner returned nil")
	}

	if scanner.RootDir != "." {
		t.Errorf("Expected RootDir '.', got '%s'", scanner.RootDir)
	}

	if !scanner.Recursive {
		t.Error("Expected Recursive to be true by default")
	}

	if scanner.IncludeTests {
		t.Error("Expected IncludeTests to be false by default")
	}

	if len(scanner.ExcludePatterns) == 0 {
		t.Error("Expected default exclude patterns")
	}
}

func TestScannerWithOptions(t *testing.T) {
	scanner := NewScanner(".").
		WithRecursive(false).
		WithIncludeTests(true).
		WithExcludePatterns([]string{"custom/*"})

	if scanner.Recursive {
		t.Error("Expected Recursive to be false")
	}

	if !scanner.IncludeTests {
		t.Error("Expected IncludeTests to be true")
	}

	if len(scanner.ExcludePatterns) != 1 || scanner.ExcludePatterns[0] != "custom/*" {
		t.Errorf("Expected custom exclude patterns, got %v", scanner.ExcludePatterns)
	}
}

func TestIsGoFile(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		includeTests bool
		want         bool
	}{
		{"go file", "main.go", false, true},
		{"test file excluded", "main_test.go", false, false},
		{"test file included", "main_test.go", true, true},
		{"non-go file", "main.txt", false, false},
		{"hidden go file", ".hidden.go", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(".").WithIncludeTests(tt.includeTests)
			got := scanner.isGoFile(tt.path)
			if got != tt.want {
				t.Errorf("isGoFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestShouldExclude(t *testing.T) {
	scanner := NewScanner("/project").
		WithExcludePatterns([]string{"vendor/*", "testdata/*", ".git/*"})

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"vendor dir", "/project/vendor/pkg/mod.go", true},
		{"testdata dir", "/project/pkg/testdata/test.go", true},
		{"git dir", "/project/.git/config", true},
		{"normal file", "/project/main.go", false},
		{"sub package", "/project/pkg/util/util.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scanner.shouldExclude(tt.path)
			if got != tt.want {
				t.Errorf("shouldExclude(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestScan(t *testing.T) {
	// 创建临时测试目录
	tmpDir := t.TempDir()

	// 创建测试文件
	testFiles := map[string]string{
		"main.go":           "package main\n\nfunc main() {}\n",
		"util/util.go":      "package util\n\nfunc Helper() {}\n",
		"util/util_test.go": "package util\n\nimport \"testing\"\n",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("scan without tests", func(t *testing.T) {
		scanner := NewScanner(tmpDir).WithIncludeTests(false)
		result, err := scanner.Scan()
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}

		if result.Stats.TotalFiles != 2 {
			t.Errorf("Expected 2 files, got %d", result.Stats.TotalFiles)
		}

		if len(result.Files) != 2 {
			t.Errorf("Expected 2 files in result, got %d", len(result.Files))
		}
	})

	t.Run("scan with tests", func(t *testing.T) {
		scanner := NewScanner(tmpDir).WithIncludeTests(true)
		result, err := scanner.Scan()
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}

		if result.Stats.TotalFiles != 3 {
			t.Errorf("Expected 3 files, got %d", result.Stats.TotalFiles)
		}

		if len(result.Files) != 3 {
			t.Errorf("Expected 3 files in result, got %d", len(result.Files))
		}
	})

	t.Run("non-recursive scan", func(t *testing.T) {
		scanner := NewScanner(tmpDir).WithRecursive(false)
		result, err := scanner.Scan()
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}

		if result.Stats.TotalFiles != 1 {
			t.Errorf("Expected 1 file (only root level), got %d", result.Stats.TotalFiles)
		}
	})
}

func TestScanNonExistentDirectory(t *testing.T) {
	scanner := NewScanner("/non/existent/directory")
	_, err := scanner.Scan()
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestFindGoModFiles(t *testing.T) {
	// 创建临时测试目录
	tmpDir := t.TempDir()

	// 创建go.mod文件
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewScanner(tmpDir)
	goModFiles, err := scanner.FindGoModFiles()
	if err != nil {
		t.Fatalf("FindGoModFiles failed: %v", err)
	}

	if len(goModFiles) != 1 {
		t.Errorf("Expected 1 go.mod file, got %d", len(goModFiles))
	}

	if goModFiles[0] != goModPath {
		t.Errorf("Expected go.mod at %s, got %s", goModPath, goModFiles[0])
	}
}

func TestGetProjectInfo(t *testing.T) {
	// 创建临时测试目录
	tmpDir := t.TempDir()

	// 创建测试文件
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	mainGoPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte("package main\n\nfunc main() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewScanner(tmpDir)
	info, err := scanner.GetProjectInfo()
	if err != nil {
		t.Fatalf("GetProjectInfo failed: %v", err)
	}

	if !info.HasGoMod {
		t.Error("Expected HasGoMod to be true")
	}

	if info.GoModPath != goModPath {
		t.Errorf("Expected GoModPath %s, got %s", goModPath, info.GoModPath)
	}

	if info.TotalFiles != 1 {
		t.Errorf("Expected 1 file, got %d", info.TotalFiles)
	}

	if info.TotalLines == 0 {
		t.Error("Expected non-zero line count")
	}
}

func TestScanWithFilter(t *testing.T) {
	// 创建临时测试目录
	tmpDir := t.TempDir()

	// 创建测试文件
	testFiles := []string{"main.go", "util.go", "test.go"}
	for _, file := range testFiles {
		path := filepath.Join(tmpDir, file)
		if err := os.WriteFile(path, []byte("package main\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	scanner := NewScanner(tmpDir)

	// 使用自定义过滤器，只包含main.go和util.go
	result, err := scanner.ScanWithFilter(func(path string, info os.FileInfo) bool {
		if info.IsDir() {
			return true
		}
		name := filepath.Base(path)
		return name == "main.go" || name == "util.go"
	})

	if err != nil {
		t.Fatalf("ScanWithFilter failed: %v", err)
	}

	if result.Stats.TotalFiles != 2 {
		t.Errorf("Expected 2 files, got %d", result.Stats.TotalFiles)
	}
}

func BenchmarkScan(b *testing.B) {
	// 创建临时测试目录
	tmpDir := b.TempDir()

	// 创建100个测试文件
	for i := 0; i < 100; i++ {
		path := filepath.Join(tmpDir, filepath.Join("pkg", "file"+string(rune(i))+".go"))
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			b.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("package pkg\n\nfunc F() {}\n"), 0644); err != nil {
			b.Fatal(err)
		}
	}

	scanner := NewScanner(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.Scan()
		if err != nil {
			b.Fatal(err)
		}
	}
}
