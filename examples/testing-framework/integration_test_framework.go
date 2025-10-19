package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TestSuite 测试套件
type TestSuite struct {
	Name        string
	Description string
	Tests       []Test
	Setup       func() error
	Teardown    func() error
	Timeout     time.Duration
	Parallel    bool
	mu          sync.RWMutex
}

// Test 单个测试
type Test struct {
	Name        string
	Description string
	Run         func(ctx context.Context) error
	Timeout     time.Duration
	Retries     int
	Required    bool
	Tags        []string
}

// TestResult 测试结果
type TestResult struct {
	Test      *Test
	Status    TestStatus
	Error     error
	Duration  time.Duration
	Retries   int
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// TestStatus 测试状态
type TestStatus string

const (
	TestStatusPassed  TestStatus = "passed"
	TestStatusFailed  TestStatus = "failed"
	TestStatusSkipped TestStatus = "skipped"
	TestStatusTimeout TestStatus = "timeout"
)

// TestEnvironment 测试环境
type TestEnvironment struct {
	Name         string
	Config       map[string]interface{}
	Resources    map[string]interface{}
	CleanupFuncs []func() error
	mu           sync.RWMutex
}

// TestExecutor 测试执行器
type TestExecutor struct {
	suites       map[string]*TestSuite
	environments map[string]*TestEnvironment
	results      []TestResult
	config       *TestConfig
	mu           sync.RWMutex
}

// TestConfig 测试配置
type TestConfig struct {
	DefaultTimeout time.Duration
	MaxRetries     int
	Parallel       bool
	MaxWorkers     int
	ReportFormat   string
	OutputDir      string
}

// NewTestSuite 创建新的测试套件
func NewTestSuite(name, description string) *TestSuite {
	return &TestSuite{
		Name:        name,
		Description: description,
		Tests:       make([]Test, 0),
		Timeout:     30 * time.Second,
		Parallel:    false,
	}
}

// AddTest 添加测试到套件
func (ts *TestSuite) AddTest(test Test) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.Tests = append(ts.Tests, test)
}

// Run 运行测试套件
func (ts *TestSuite) Run(ctx context.Context, executor *TestExecutor) ([]TestResult, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	var results []TestResult

	// 执行Setup
	if ts.Setup != nil {
		if err := ts.Setup(); err != nil {
			return nil, fmt.Errorf("setup failed: %w", err)
		}
	}

	// 确保Teardown被执行
	defer func() {
		if ts.Teardown != nil {
			ts.Teardown()
		}
	}()

	// 运行测试
	if ts.Parallel {
		results = ts.runParallel(ctx, executor)
	} else {
		results = ts.runSequential(ctx, executor)
	}

	return results, nil
}

// runSequential 顺序运行测试
func (ts *TestSuite) runSequential(ctx context.Context, executor *TestExecutor) []TestResult {
	var results []TestResult

	for _, test := range ts.Tests {
		result := executor.runTest(ctx, &test)
		results = append(results, result)

		// 如果测试失败且是必需的，停止执行
		if result.Status == TestStatusFailed && test.Required {
			break
		}
	}

	return results
}

// runParallel 并行运行测试
func (ts *TestSuite) runParallel(ctx context.Context, executor *TestExecutor) []TestResult {
	var results []TestResult
	var wg sync.WaitGroup
	resultChan := make(chan TestResult, len(ts.Tests))

	// 启动工作协程
	for _, test := range ts.Tests {
		wg.Add(1)
		go func(t Test) {
			defer wg.Done()
			result := executor.runTest(ctx, &t)
			resultChan <- result
		}(test)
	}

	// 等待所有测试完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// NewTestExecutor 创建测试执行器
func NewTestExecutor(config *TestConfig) *TestExecutor {
	if config == nil {
		config = &TestConfig{
			DefaultTimeout: 30 * time.Second,
			MaxRetries:     3,
			Parallel:       false,
			MaxWorkers:     4,
			ReportFormat:   "json",
			OutputDir:      "./test-results",
		}
	}

	return &TestExecutor{
		suites:       make(map[string]*TestSuite),
		environments: make(map[string]*TestEnvironment),
		results:      make([]TestResult, 0),
		config:       config,
	}
}

// RegisterSuite 注册测试套件
func (te *TestExecutor) RegisterSuite(suite *TestSuite) {
	te.mu.Lock()
	defer te.mu.Unlock()
	te.suites[suite.Name] = suite
}

// RegisterEnvironment 注册测试环境
func (te *TestExecutor) RegisterEnvironment(env *TestEnvironment) {
	te.mu.Lock()
	defer te.mu.Unlock()
	te.environments[env.Name] = env
}

// RunSuite 运行指定套件
func (te *TestExecutor) RunSuite(ctx context.Context, suiteName string) ([]TestResult, error) {
	te.mu.RLock()
	suite, exists := te.suites[suiteName]
	te.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("test suite '%s' not found", suiteName)
	}

	results, err := suite.Run(ctx, te)
	if err != nil {
		return nil, err
	}

	// 保存结果
	te.mu.Lock()
	te.results = append(te.results, results...)
	te.mu.Unlock()

	return results, nil
}

// RunAllSuites 运行所有套件
func (te *TestExecutor) RunAllSuites(ctx context.Context) (map[string][]TestResult, error) {
	te.mu.RLock()
	suiteNames := make([]string, 0, len(te.suites))
	for name := range te.suites {
		suiteNames = append(suiteNames, name)
	}
	te.mu.RUnlock()

	allResults := make(map[string][]TestResult)

	for _, suiteName := range suiteNames {
		results, err := te.RunSuite(ctx, suiteName)
		if err != nil {
			return allResults, fmt.Errorf("failed to run suite '%s': %w", suiteName, err)
		}
		allResults[suiteName] = results
	}

	return allResults, nil
}

// runTest 运行单个测试
func (te *TestExecutor) runTest(ctx context.Context, test *Test) TestResult {
	result := TestResult{
		Test:      test,
		Status:    TestStatusFailed,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	// 设置超时
	timeout := test.Timeout
	if timeout == 0 {
		timeout = te.config.DefaultTimeout
	}

	testCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 执行测试，支持重试
	var err error
	for attempt := 0; attempt <= test.Retries; attempt++ {
		start := time.Now()
		err = test.Run(testCtx)
		result.Duration = time.Since(start)
		result.Retries = attempt

		if err == nil {
			result.Status = TestStatusPassed
			break
		}

		// 如果不是最后一次尝试，等待后重试
		if attempt < test.Retries {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	if err != nil {
		result.Error = err
		if testCtx.Err() == context.DeadlineExceeded {
			result.Status = TestStatusTimeout
		} else {
			result.Status = TestStatusFailed
		}
	}

	return result
}

// GetResults 获取测试结果
func (te *TestExecutor) GetResults() []TestResult {
	te.mu.RLock()
	defer te.mu.RUnlock()

	results := make([]TestResult, len(te.results))
	copy(results, te.results)
	return results
}

// GetResultsByStatus 按状态获取结果
func (te *TestExecutor) GetResultsByStatus(status TestStatus) []TestResult {
	te.mu.RLock()
	defer te.mu.RUnlock()

	var filtered []TestResult
	for _, result := range te.results {
		if result.Status == status {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// GetTestSummary 获取测试摘要
func (te *TestExecutor) GetTestSummary() TestSummary {
	te.mu.RLock()
	defer te.mu.RUnlock()

	summary := TestSummary{
		Total:    len(te.results),
		Passed:   0,
		Failed:   0,
		Skipped:  0,
		Timeout:  0,
		Duration: 0,
	}

	for _, result := range te.results {
		summary.Duration += result.Duration
		switch result.Status {
		case TestStatusPassed:
			summary.Passed++
		case TestStatusFailed:
			summary.Failed++
		case TestStatusSkipped:
			summary.Skipped++
		case TestStatusTimeout:
			summary.Timeout++
		}
	}

	return summary
}

// TestSummary 测试摘要
type TestSummary struct {
	Total    int           `json:"total"`
	Passed   int           `json:"passed"`
	Failed   int           `json:"failed"`
	Skipped  int           `json:"skipped"`
	Timeout  int           `json:"timeout"`
	Duration time.Duration `json:"duration"`
}

// NewTestEnvironment 创建测试环境
func NewTestEnvironment(name string) *TestEnvironment {
	return &TestEnvironment{
		Name:         name,
		Config:       make(map[string]interface{}),
		Resources:    make(map[string]interface{}),
		CleanupFuncs: make([]func() error, 0),
	}
}

// SetConfig 设置环境配置
func (te *TestEnvironment) SetConfig(key string, value interface{}) {
	te.mu.Lock()
	defer te.mu.Unlock()
	te.Config[key] = value
}

// GetConfig 获取环境配置
func (te *TestEnvironment) GetConfig(key string) (interface{}, bool) {
	te.mu.RLock()
	defer te.mu.RUnlock()
	value, exists := te.Config[key]
	return value, exists
}

// AddResource 添加资源
func (te *TestEnvironment) AddResource(name string, resource interface{}) {
	te.mu.Lock()
	defer te.mu.Unlock()
	te.Resources[name] = resource
}

// GetResource 获取资源
func (te *TestEnvironment) GetResource(name string) (interface{}, bool) {
	te.mu.RLock()
	defer te.mu.RUnlock()
	resource, exists := te.Resources[name]
	return resource, exists
}

// AddCleanup 添加清理函数
func (te *TestEnvironment) AddCleanup(cleanup func() error) {
	te.mu.Lock()
	defer te.mu.Unlock()
	te.CleanupFuncs = append(te.CleanupFuncs, cleanup)
}

// Cleanup 执行清理
func (te *TestEnvironment) Cleanup() error {
	te.mu.Lock()
	defer te.mu.Unlock()

	var lastErr error
	for _, cleanup := range te.CleanupFuncs {
		if err := cleanup(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
