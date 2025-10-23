// cpg (Concurrency Pattern Generator) - Go并发模式代码生成器
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/your-org/concurrency-pattern-generator/pkg/generator"
	"github.com/your-org/concurrency-pattern-generator/pkg/patterns"
)

const version = "v1.0.0"

var (
	// 命令标志
	patternFlag  = flag.String("pattern", "", "Pattern type to generate (e.g., worker-pool, fan-in)")
	outputFlag   = flag.String("output", "", "Output file path")
	packageFlag  = flag.String("package", "main", "Package name")
	workersFlag  = flag.Int("workers", 10, "Number of workers (for worker-pool)")
	bufferFlag   = flag.Int("buffer", 0, "Buffer size for channels")
	fanoutFlag   = flag.Int("fanout", 5, "Fan-out N")
	listFlag     = flag.Bool("list", false, "List all available patterns")
	categoryFlag = flag.String("category", "", "List patterns by category")
	versionFlag  = flag.Bool("version", false, "Show version")
	helpFlag     = flag.Bool("help", false, "Show help")
)

func main() {
	flag.Parse()

	// 处理版本标志
	if *versionFlag {
		fmt.Printf("cpg (Concurrency Pattern Generator) %s\n", version)
		fmt.Println("Based on CSP formal verification")
		return
	}

	// 处理帮助标志
	if *helpFlag {
		printHelp()
		return
	}

	// 列出所有模式
	if *listFlag {
		listPatterns("")
		return
	}

	// 按类别列出模式
	if *categoryFlag != "" {
		listPatterns(*categoryFlag)
		return
	}

	// 生成模式
	if *patternFlag == "" {
		fmt.Println("❌ Error: pattern is required")
		fmt.Println("Use --help for usage information")
		os.Exit(1)
	}

	if err := generatePattern(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Pattern generated successfully!")
}

func generatePattern() error {
	// 验证模式类型
	patternType := generator.PatternType(*patternFlag)

	// 准备配置
	config := &generator.Config{
		PatternType: patternType,
		PackageName: *packageFlag,
		OutputFile:  *outputFlag,
		NumWorkers:  *workersFlag,
		BufferSize:  *bufferFlag,
		FanOutN:     *fanoutFlag,
	}

	// 生成代码
	var code string
	var err error

	switch patternType {
	case generator.WorkerPoolPattern:
		data := map[string]interface{}{
			"PackageName": config.PackageName,
			"NumWorkers":  config.NumWorkers,
		}
		code = patterns.GenerateWorkerPool(data)

	case generator.FanInPattern:
		data := map[string]interface{}{
			"PackageName": config.PackageName,
		}
		code = patterns.GenerateFanIn(data)

	case generator.FanOutPattern:
		data := map[string]interface{}{
			"PackageName": config.PackageName,
			"FanOutN":     config.FanOutN,
		}
		code = patterns.GenerateFanOut(data)

	case generator.PipelinePattern:
		data := map[string]interface{}{
			"PackageName": config.PackageName,
		}
		code = patterns.GeneratePipeline(data)

	case generator.GeneratorPattern:
		data := map[string]interface{}{
			"PackageName": config.PackageName,
		}
		code = patterns.GenerateGenerator(data)

	// Sync Patterns (使用简化版本)
	case generator.MutexPattern:
		code = patterns.GenerateMutexSimple(config.PackageName)

	case generator.RWMutexPattern:
		code = patterns.GenerateRWMutexSimple(config.PackageName)

	case generator.WaitGroupPattern:
		code = patterns.GenerateWaitGroupSimple(config.PackageName)

	case generator.OncePattern:
		code = patterns.GenerateOnceSimple(config.PackageName)

	case generator.CondPattern:
		code = patterns.GenerateCondSimple(config.PackageName)

	case generator.SemaphorePattern:
		code = patterns.GenerateSemaphoreSimple(config.PackageName)

	case generator.BarrierPattern:
		code = patterns.GenerateBarrierSimple(config.PackageName)

	case generator.CountDownLatchPattern:
		code = patterns.GenerateCountDownLatchSimple(config.PackageName)

	default:
		return fmt.Errorf("pattern not yet implemented: %s", patternType)
	}

	if err != nil {
		return err
	}

	// 输出代码
	if *outputFlag != "" {
		if err := os.WriteFile(*outputFlag, []byte(code), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("📝 Generated: %s\n", *outputFlag)
		fmt.Printf("📊 Pattern: %s\n", patternType)
		fmt.Printf("📏 Lines: %d\n", countLines(code))
	} else {
		fmt.Println(code)
	}

	return nil
}

func listPatterns(category string) {
	fmt.Println("🎯 Available Concurrency Patterns\n")

	categories := map[string]string{
		"":         "All Patterns",
		"classic":  "Classic Patterns",
		"sync":     "Synchronization Patterns",
		"control":  "Control Flow Patterns",
		"dataflow": "Data Flow Patterns",
		"advanced": "Advanced Patterns",
	}

	if title, ok := categories[category]; ok {
		fmt.Printf("📚 %s:\n\n", title)
	}

	var patterns []generator.PatternType
	if category == "" {
		patterns = generator.GetAllPatterns()
	} else {
		patterns = generator.GetPatternsByCategory(category)
	}

	// 按类别分组
	groups := make(map[string][]generator.PatternType)

	for _, p := range patterns {
		cat := getPatternCategory(p)
		groups[cat] = append(groups[cat], p)
	}

	// 按顺序输出
	catOrder := []string{"Classic", "Sync", "Control Flow", "Data Flow", "Advanced"}
	for _, cat := range catOrder {
		if pats, ok := groups[cat]; ok && len(pats) > 0 {
			fmt.Printf("  %s:\n", cat)
			for _, p := range pats {
				fmt.Printf("    - %s\n", p)
			}
			fmt.Println()
		}
	}

	fmt.Println("💡 Usage:")
	fmt.Println("   cpg --pattern <pattern-name> --output <file.go>")
	fmt.Println()
	fmt.Println("📖 Examples:")
	fmt.Println("   cpg --pattern worker-pool --workers 10 --output pool.go")
	fmt.Println("   cpg --pattern fan-in --output fanin.go")
	fmt.Println("   cpg --pattern pipeline --output pipeline.go")
}

func getPatternCategory(p generator.PatternType) string {
	classic := []generator.PatternType{
		generator.WorkerPoolPattern, generator.FanInPattern,
		generator.FanOutPattern, generator.PipelinePattern,
		generator.GeneratorPattern,
	}
	sync := []generator.PatternType{
		generator.MutexPattern, generator.RWMutexPattern,
		generator.WaitGroupPattern, generator.OncePattern,
		generator.CondPattern, generator.SemaphorePattern,
		generator.BarrierPattern, generator.CountDownLatchPattern,
	}
	control := []generator.PatternType{
		generator.ContextCancelPattern, generator.ContextTimeoutPattern,
		generator.ContextValuePattern, generator.GracefulShutdownPattern,
		generator.RateLimitingPattern,
	}
	dataflow := []generator.PatternType{
		generator.ProducerConsumerPattern, generator.BufferedChannelPattern,
		generator.UnbufferedChannelPattern, generator.SelectPattern,
		generator.ForSelectLoopPattern, generator.DoneChannelPattern,
		generator.ErrorChannelPattern,
	}

	for _, pat := range classic {
		if pat == p {
			return "Classic"
		}
	}
	for _, pat := range sync {
		if pat == p {
			return "Sync"
		}
	}
	for _, pat := range control {
		if pat == p {
			return "Control Flow"
		}
	}
	for _, pat := range dataflow {
		if pat == p {
			return "Data Flow"
		}
	}
	return "Advanced"
}

func printHelp() {
	fmt.Println("🎯 cpg - Concurrency Pattern Generator")
	fmt.Printf("Version: %s\n\n", version)

	fmt.Println("Usage:")
	fmt.Println("  cpg [flags]")
	fmt.Println()

	fmt.Println("Flags:")
	fmt.Println("  --pattern string     Pattern type to generate (required)")
	fmt.Println("  --output string      Output file path (default: stdout)")
	fmt.Println("  --package string     Package name (default: main)")
	fmt.Println("  --workers int        Number of workers for worker-pool (default: 10)")
	fmt.Println("  --buffer int         Buffer size for channels (default: 0)")
	fmt.Println("  --fanout int         Fan-out N (default: 5)")
	fmt.Println("  --list              List all available patterns")
	fmt.Println("  --category string    List patterns by category")
	fmt.Println("  --version           Show version")
	fmt.Println("  --help              Show this help message")
	fmt.Println()

	fmt.Println("Examples:")
	fmt.Println("  # Generate worker pool pattern")
	fmt.Println("  cpg --pattern worker-pool --workers 10 --output pool.go")
	fmt.Println()
	fmt.Println("  # Generate fan-in pattern")
	fmt.Println("  cpg --pattern fan-in --output fanin.go")
	fmt.Println()
	fmt.Println("  # List all patterns")
	fmt.Println("  cpg --list")
	fmt.Println()
	fmt.Println("  # List patterns by category")
	fmt.Println("  cpg --category classic")
	fmt.Println()

	fmt.Println("Pattern Categories:")
	fmt.Println("  classic   - Classic patterns (worker-pool, fan-in, fan-out, pipeline, generator)")
	fmt.Println("  sync      - Synchronization patterns")
	fmt.Println("  control   - Control flow patterns")
	fmt.Println("  dataflow  - Data flow patterns")
	fmt.Println("  advanced  - Advanced patterns")
	fmt.Println()

	fmt.Println("Theory:")
	fmt.Println("  All patterns are based on CSP (Communicating Sequential Processes)")
	fmt.Println("  formal verification and include:")
	fmt.Println("  - CSP process definitions")
	fmt.Println("  - Deadlock freedom proofs")
	fmt.Println("  - Data race analysis")
	fmt.Println("  - Happens-before relations")
	fmt.Println()

	fmt.Println("For more information, see:")
	fmt.Println("  https://github.com/your-org/concurrency-pattern-generator")
}

func countLines(s string) int {
	return len(strings.Split(s, "\n"))
}

// runInteractive 交互式模式（TODO）
func runInteractive(ctx context.Context) error {
	// TODO: 实现交互式模式
	fmt.Println("🎨 Interactive mode (Coming soon...)")
	return nil
}

// runBatch 批量生成模式（TODO）
func runBatch(ctx context.Context, configFile string) error {
	// TODO: 实现批量生成
	fmt.Println("📦 Batch mode (Coming soon...)")
	return nil
}
