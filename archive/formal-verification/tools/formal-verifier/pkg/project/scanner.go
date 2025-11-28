package project

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Scanner 扫描Go项目中的所有源文件
//
// 功能：
//   - 递归扫描目录
//   - 过滤.go文件（排除测试文件和vendor）
//   - 解析Go源文件
//   - 收集项目统计信息
type Scanner struct {
	// RootDir 项目根目录
	RootDir string

	// Recursive 是否递归扫描子目录
	Recursive bool

	// ExcludePatterns 排除模式（glob patterns）
	ExcludePatterns []string

	// IncludeTests 是否包含测试文件
	IncludeTests bool

	// FileSet Go语法树文件集
	FileSet *token.FileSet

	// 统计信息
	TotalFiles    int
	TotalLines    int
	TotalFuncs    int
	ExcludedFiles int
}

// ScanResult 扫描结果
type ScanResult struct {
	// Files 所有扫描到的文件路径
	Files []string

	// FileASTs 文件路径到AST的映射
	FileASTs map[string]*token.File

	// Stats 统计信息
	Stats *Stats
}

// Stats 项目统计信息
type Stats struct {
	TotalFiles    int
	TotalLines    int
	TotalPackages int
	TotalFuncs    int
	ExcludedFiles int
}

// NewScanner 创建新的项目扫描器
func NewScanner(rootDir string) *Scanner {
	return &Scanner{
		RootDir:         rootDir,
		Recursive:       true,
		ExcludePatterns: defaultExcludePatterns(),
		IncludeTests:    false,
		FileSet:         token.NewFileSet(),
	}
}

// defaultExcludePatterns 返回默认的排除模式
func defaultExcludePatterns() []string {
	return []string{
		"vendor/*",
		"node_modules/*",
		".git/*",
		"**/testdata/*",
	}
}

// WithRecursive 设置是否递归扫描
func (s *Scanner) WithRecursive(recursive bool) *Scanner {
	s.Recursive = recursive
	return s
}

// WithExcludePatterns 设置排除模式
func (s *Scanner) WithExcludePatterns(patterns []string) *Scanner {
	s.ExcludePatterns = patterns
	return s
}

// WithIncludeTests 设置是否包含测试文件
func (s *Scanner) WithIncludeTests(includeTests bool) *Scanner {
	s.IncludeTests = includeTests
	return s
}

// Scan 扫描项目
func (s *Scanner) Scan() (*ScanResult, error) {
	result := &ScanResult{
		Files:    make([]string, 0),
		FileASTs: make(map[string]*token.File),
		Stats:    &Stats{},
	}

	// 检查根目录是否存在
	if _, err := os.Stat(s.RootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", s.RootDir)
	}

	// 扫描文件
	err := s.scanDirectory(s.RootDir, result)
	if err != nil {
		return nil, err
	}

	// 更新统计信息
	result.Stats.TotalFiles = s.TotalFiles
	result.Stats.TotalLines = s.TotalLines
	result.Stats.TotalFuncs = s.TotalFuncs
	result.Stats.ExcludedFiles = s.ExcludedFiles

	return result, nil
}

// scanDirectory 递归扫描目录
func (s *Scanner) scanDirectory(dir string, result *ScanResult) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		// 检查是否应该排除
		if s.shouldExclude(path) {
			s.ExcludedFiles++
			continue
		}

		if entry.IsDir() {
			// 递归扫描子目录
			if s.Recursive {
				if err := s.scanDirectory(path, result); err != nil {
					return err
				}
			}
		} else {
			// 处理文件
			if s.isGoFile(path) {
				if err := s.processFile(path, result); err != nil {
					// 记录错误但继续扫描
					fmt.Fprintf(os.Stderr, "Warning: failed to process %s: %v\n", path, err)
				}
			}
		}
	}

	return nil
}

// isGoFile 判断是否是Go源文件
func (s *Scanner) isGoFile(path string) bool {
	// 检查扩展名
	if !strings.HasSuffix(path, ".go") {
		return false
	}

	// 检查是否是测试文件
	if !s.IncludeTests && strings.HasSuffix(path, "_test.go") {
		return false
	}

	return true
}

// shouldExclude 判断路径是否应该被排除
func (s *Scanner) shouldExclude(path string) bool {
	relPath, err := filepath.Rel(s.RootDir, path)
	if err != nil {
		return false
	}

	for _, pattern := range s.ExcludePatterns {
		matched, err := filepath.Match(pattern, relPath)
		if err == nil && matched {
			return true
		}

		// 检查是否匹配目录
		if strings.Contains(relPath, strings.TrimSuffix(pattern, "/*")) {
			return true
		}
	}

	return false
}

// processFile 处理单个Go源文件
func (s *Scanner) processFile(path string, result *ScanResult) error {
	// 添加文件到结果（即使解析失败也要记录）
	result.Files = append(result.Files, path)
	s.TotalFiles++

	// 解析文件
	f, err := parser.ParseFile(s.FileSet, path, nil, parser.ParseComments)
	if err != nil {
		// 解析失败，返回错误但文件已被记录
		return fmt.Errorf("failed to parse file: %w", err)
	}

	// 解析成功，记录AST和更新统计
	result.FileASTs[path] = s.FileSet.File(f.Pos())
	s.TotalFuncs += len(f.Decls)

	// 计算行数
	if tokenFile := s.FileSet.File(f.Pos()); tokenFile != nil {
		s.TotalLines += tokenFile.LineCount()
	}

	return nil
}

// ScanWithFilter 使用自定义过滤器扫描
func (s *Scanner) ScanWithFilter(filter func(path string, info fs.FileInfo) bool) (*ScanResult, error) {
	result := &ScanResult{
		Files:    make([]string, 0),
		FileASTs: make(map[string]*token.File),
		Stats:    &Stats{},
	}

	err := filepath.Walk(s.RootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 应用自定义过滤器
		if !filter(path, info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 检查默认排除规则
		if s.shouldExclude(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() && s.isGoFile(path) {
			if err := s.processFile(path, result); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to process %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 更新统计信息
	result.Stats.TotalFiles = s.TotalFiles
	result.Stats.TotalLines = s.TotalLines
	result.Stats.TotalFuncs = s.TotalFuncs
	result.Stats.ExcludedFiles = s.ExcludedFiles

	return result, nil
}

// FindGoModFiles 查找项目中的go.mod文件
func (s *Scanner) FindGoModFiles() ([]string, error) {
	var goModFiles []string

	err := filepath.Walk(s.RootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "go.mod" {
			goModFiles = append(goModFiles, path)
		}

		return nil
	})

	return goModFiles, err
}

// GetProjectInfo 获取项目基本信息
func (s *Scanner) GetProjectInfo() (*ProjectInfo, error) {
	info := &ProjectInfo{
		RootDir: s.RootDir,
	}

	// 查找go.mod
	goModFiles, err := s.FindGoModFiles()
	if err != nil {
		return nil, err
	}

	if len(goModFiles) > 0 {
		info.HasGoMod = true
		info.GoModPath = goModFiles[0]
	}

	// 扫描文件统计
	result, err := s.Scan()
	if err != nil {
		return nil, err
	}

	info.TotalFiles = result.Stats.TotalFiles
	info.TotalLines = result.Stats.TotalLines

	return info, nil
}

// ProjectInfo 项目基本信息
type ProjectInfo struct {
	RootDir    string
	HasGoMod   bool
	GoModPath  string
	TotalFiles int
	TotalLines int
}
