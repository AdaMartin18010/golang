package testing_system

import (
	"context"
	"fmt"
	"time"
)

// ExampleTestSystem 测试体系使用示例
func ExampleTestSystem() {
	// 1. 创建测试执行器
	config := &TestConfig{
		DefaultTimeout: 30 * time.Second,
		MaxRetries:     3,
		Parallel:       true,
		MaxWorkers:     4,
		ReportFormat:   "json",
		OutputDir:      "./test-results",
	}
	
	executor := NewTestExecutor(config)
	
	// 2. 创建测试套件
	suite := NewTestSuite("API测试套件", "测试REST API的功能和性能")
	
	// 3. 添加测试用例
	suite.AddTest(Test{
		Name:        "用户注册测试",
		Description: "测试用户注册功能",
		Run: func(ctx context.Context) error {
			// 模拟用户注册测试
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		Timeout:  10 * time.Second,
		Retries:  2,
		Required: true,
		Tags:     []string{"api", "user", "registration"},
	})
	
	suite.AddTest(Test{
		Name:        "用户登录测试",
		Description: "测试用户登录功能",
		Run: func(ctx context.Context) error {
			// 模拟用户登录测试
			time.Sleep(50 * time.Millisecond)
			return nil
		},
		Timeout:  5 * time.Second,
		Retries:  1,
		Required: true,
		Tags:     []string{"api", "user", "login"},
	})
	
	suite.AddTest(Test{
		Name:        "数据查询测试",
		Description: "测试数据查询功能",
		Run: func(ctx context.Context) error {
			// 模拟数据查询测试
			time.Sleep(200 * time.Millisecond)
			return nil
		},
		Timeout:  15 * time.Second,
		Retries:  3,
		Required: false,
		Tags:     []string{"api", "data", "query"},
	})
	
	// 4. 注册测试套件
	executor.RegisterSuite(suite)
	
	// 5. 运行测试套件
	ctx := context.Background()
	results, err := executor.RunSuite(ctx, "API测试套件")
	if err != nil {
		fmt.Printf("测试执行失败: %v\n", err)
		return
	}
	
	// 6. 分析测试结果
	summary := executor.GetTestSummary()
	fmt.Printf("测试摘要:\n")
	fmt.Printf("  总数: %d\n", summary.Total)
	fmt.Printf("  通过: %d\n", summary.Passed)
	fmt.Printf("  失败: %d\n", summary.Failed)
	fmt.Printf("  跳过: %d\n", summary.Skipped)
	fmt.Printf("  超时: %d\n", summary.Timeout)
	fmt.Printf("  总耗时: %v\n", summary.Duration)
	
	// 7. 输出详细结果
	for _, result := range results {
		fmt.Printf("测试: %s - %s (耗时: %v, 重试: %d)\n", 
			result.Test.Name, result.Status, result.Duration, result.Retries)
		if result.Error != nil {
			fmt.Printf("  错误: %v\n", result.Error)
		}
	}
}

// ExamplePerformanceTesting 性能测试示例
func ExamplePerformanceTesting() {
	// 1. 创建性能监控器
	perfConfig := &PerformanceConfig{
		DefaultTimeout:        60 * time.Second,
		DefaultIterations:     100,
		DefaultWarmup:         10,
		RegressionThreshold:   0.1,
		OutputDir:             "./performance-results",
		ReportFormat:          "json",
		EnableProfiling:       true,
		ProfilingDir:          "./profiles",
	}
	
	monitor := NewPerformanceMonitor(perfConfig)
	
	// 2. 创建性能基准测试
	apiBenchmark := NewPerformanceBenchmark(
		"API响应时间测试",
		"测试API接口的响应时间性能",
		func(ctx context.Context) (BenchmarkResult, error) {
			start := time.Now()
			
			// 模拟API调用
			time.Sleep(10 * time.Millisecond)
			
			duration := time.Since(start)
			
			return BenchmarkResult{
				Name:       "API响应时间测试",
				Duration:   duration,
				Operations: 1,
				Throughput: 1.0 / duration.Seconds(),
				MemoryUsage: MemoryUsage{
					Allocated: 1024,
					Total:     2048,
					Heap:      512,
					Stack:     256,
					GC:        0,
				},
				CPUUsage: CPUUsage{
					UserTime:   duration / 2,
					SystemTime: duration / 4,
					IdleTime:   duration / 4,
					Usage:      75.0,
				},
				Iterations:  1,
				MinDuration: duration,
				MaxDuration: duration,
				AvgDuration: duration,
				StdDev:      0,
				Percentiles: map[int]time.Duration{
					50: duration,
					90: duration,
					95: duration,
					99: duration,
				},
				Timestamp: time.Now(),
				Metadata:  make(map[string]interface{}),
			}, nil
		},
	)
	
	// 3. 注册基准测试
	monitor.RegisterBenchmark(apiBenchmark)
	
	// 4. 运行性能测试
	ctx := context.Background()
	result, err := monitor.RunBenchmark(ctx, "API响应时间测试")
	if err != nil {
		fmt.Printf("性能测试失败: %v\n", err)
		return
	}
	
	// 5. 输出性能结果
	fmt.Printf("性能测试结果:\n")
	fmt.Printf("  测试名称: %s\n", result.Name)
	fmt.Printf("  平均耗时: %v\n", result.AvgDuration)
	fmt.Printf("  最小耗时: %v\n", result.MinDuration)
	fmt.Printf("  最大耗时: %v\n", result.MaxDuration)
	fmt.Printf("  标准差: %.2f\n", result.StdDev)
	fmt.Printf("  吞吐量: %.2f ops/s\n", result.Throughput)
	fmt.Printf("  内存分配: %d bytes\n", result.MemoryUsage.Allocated)
	fmt.Printf("  CPU使用率: %.2f%%\n", result.CPUUsage.Usage)
	
	// 6. 检查性能回归
	alerts := monitor.GetRegressionAlerts()
	if len(alerts) > 0 {
		fmt.Printf("检测到性能回归:\n")
		for _, alert := range alerts {
			fmt.Printf("  %s: %s (基准: %.2f, 当前: %.2f, 下降: %.2f%%)\n",
				alert.BenchmarkName, alert.Message,
				alert.Baseline, alert.Current, alert.Degradation*100)
		}
	} else {
		fmt.Printf("未检测到性能回归\n")
	}
}

// ExampleQualityDashboard 质量监控仪表板示例
func ExampleQualityDashboard() {
	// 1. 创建质量监控仪表板
	dashboardConfig := &DashboardConfig{
		Port:            8080,
		RefreshInterval: 30 * time.Second,
		RetentionPeriod: 24 * time.Hour,
		MaxDataPoints:   1000,
		EnableRealTime:  true,
		Theme:           "default",
	}
	
	dashboard := NewQualityDashboard(dashboardConfig)
	
	// 2. 启动仪表板
	ctx := context.Background()
	if err := dashboard.Start(ctx); err != nil {
		fmt.Printf("启动仪表板失败: %v\n", err)
		return
	}
	
	// 3. 获取监控数据
	metrics := dashboard.GetMetrics()
	fmt.Printf("当前指标数量: %d\n", len(metrics))
	
	alerts := dashboard.GetAlerts()
	fmt.Printf("当前告警数量: %d\n", len(alerts))
	
	charts := dashboard.GetCharts()
	fmt.Printf("当前图表数量: %d\n", len(charts))
	
	// 4. 创建示例图表
	visualizer := dashboard.visualizer
	chartData := []ChartDataPoint{
		{Label: "通过", Value: 85, Color: "#2ca02c"},
		{Label: "失败", Value: 10, Color: "#d62728"},
		{Label: "跳过", Value: 5, Color: "#ff7f0e"},
	}
	
	chart := visualizer.CreateChart("test-results", "测试结果分布", ChartTypePie, chartData)
	fmt.Printf("创建图表: %s (%s)\n", chart.Name, chart.Type)
	
	// 5. 模拟运行一段时间
	fmt.Printf("仪表板已启动，运行中...\n")
	time.Sleep(5 * time.Second)
	
	// 6. 停止仪表板
	dashboard.Stop()
	fmt.Printf("仪表板已停止\n")
}

// ExampleCompleteTestingWorkflow 完整测试工作流示例
func ExampleCompleteTestingWorkflow() {
	fmt.Printf("=== 完整测试工作流示例 ===\n")
	
	// 1. 集成测试
	fmt.Printf("\n1. 运行集成测试...\n")
	ExampleTestSystem()
	
	// 2. 性能测试
	fmt.Printf("\n2. 运行性能测试...\n")
	ExamplePerformanceTesting()
	
	// 3. 质量监控
	fmt.Printf("\n3. 启动质量监控...\n")
	ExampleQualityDashboard()
	
	fmt.Printf("\n=== 测试工作流完成 ===\n")
}

// ExampleTestEnvironment 测试环境管理示例
func ExampleTestEnvironment() {
	// 1. 创建测试环境
	env := NewTestEnvironment("测试环境")
	
	// 2. 设置环境配置
	env.SetConfig("database_url", "postgres://localhost:5432/testdb")
	env.SetConfig("api_base_url", "http://localhost:8080/api")
	env.SetConfig("timeout", 30)
	
	// 3. 添加资源
	env.AddResource("db_connection", "模拟数据库连接")
	env.AddResource("http_client", "模拟HTTP客户端")
	
	// 4. 添加清理函数
	env.AddCleanup(func() error {
		fmt.Printf("清理数据库连接...\n")
		return nil
	})
	
	env.AddCleanup(func() error {
		fmt.Printf("清理HTTP客户端...\n")
		return nil
	})
	
	// 5. 使用环境配置
	if dbURL, exists := env.GetConfig("database_url"); exists {
		fmt.Printf("数据库URL: %v\n", dbURL)
	}
	
	if client, exists := env.GetResource("http_client"); exists {
		fmt.Printf("HTTP客户端: %v\n", client)
	}
	
	// 6. 执行清理
	if err := env.Cleanup(); err != nil {
		fmt.Printf("清理失败: %v\n", err)
	}
}

// ExampleCustomTest 自定义测试示例
func ExampleCustomTest() {
	// 1. 创建测试执行器
	executor := NewTestExecutor(nil)
	
	// 2. 创建自定义测试套件
	suite := NewTestSuite("自定义测试套件", "演示自定义测试用例")
	
	// 3. 添加自定义测试
	suite.AddTest(Test{
		Name:        "自定义业务逻辑测试",
		Description: "测试特定的业务逻辑",
		Run: func(ctx context.Context) error {
			// 模拟业务逻辑测试
			fmt.Printf("执行自定义业务逻辑测试...\n")
			
			// 模拟一些业务验证
			if err := validateBusinessLogic(); err != nil {
				return fmt.Errorf("业务逻辑验证失败: %w", err)
			}
			
			fmt.Printf("业务逻辑测试通过\n")
			return nil
		},
		Timeout:  20 * time.Second,
		Retries:  2,
		Required: true,
		Tags:     []string{"business", "custom"},
	})
	
	// 4. 注册并运行
	executor.RegisterSuite(suite)
	
	ctx := context.Background()
	results, err := executor.RunSuite(ctx, "自定义测试套件")
	if err != nil {
		fmt.Printf("自定义测试失败: %v\n", err)
		return
	}
	
	// 5. 输出结果
	for _, result := range results {
		fmt.Printf("自定义测试结果: %s - %s\n", result.Test.Name, result.Status)
	}
}

// validateBusinessLogic 模拟业务逻辑验证
func validateBusinessLogic() error {
	// 模拟业务逻辑验证
	time.Sleep(100 * time.Millisecond)
	
	// 模拟随机失败（10%概率）
	if time.Now().UnixNano()%10 == 0 {
		return fmt.Errorf("业务逻辑验证失败")
	}
	
	return nil
}

// ExampleParallelTesting 并行测试示例
func ExampleParallelTesting() {
	// 1. 创建支持并行的测试执行器
	config := &TestConfig{
		DefaultTimeout: 30 * time.Second,
		MaxRetries:     1,
		Parallel:       true,
		MaxWorkers:     4,
	}
	
	executor := NewTestExecutor(config)
	
	// 2. 创建并行测试套件
	suite := NewTestSuite("并行测试套件", "演示并行测试执行")
	suite.Parallel = true
	
	// 3. 添加多个可以并行执行的测试
	for i := 1; i <= 5; i++ {
		testNum := i
		suite.AddTest(Test{
			Name:        fmt.Sprintf("并行测试-%d", testNum),
			Description: fmt.Sprintf("第%d个并行测试", testNum),
			Run: func(ctx context.Context) error {
				fmt.Printf("执行并行测试-%d...\n", testNum)
				time.Sleep(time.Duration(testNum*100) * time.Millisecond)
				fmt.Printf("并行测试-%d完成\n", testNum)
				return nil
			},
			Timeout:  10 * time.Second,
			Retries:  0,
			Required: false,
			Tags:     []string{"parallel", fmt.Sprintf("test-%d", testNum)},
		})
	}
	
	// 4. 注册并运行
	executor.RegisterSuite(suite)
	
	ctx := context.Background()
	start := time.Now()
	results, err := executor.RunSuite(ctx, "并行测试套件")
	duration := time.Since(start)
	
	if err != nil {
		fmt.Printf("并行测试失败: %v\n", err)
		return
	}
	
	// 5. 输出结果
	fmt.Printf("并行测试完成，总耗时: %v\n", duration)
	for _, result := range results {
		fmt.Printf("  %s: %s (耗时: %v)\n", result.Test.Name, result.Status, result.Duration)
	}
}
