// Package main implements the Go Formal Verifier command-line tool
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/your-org/formal-verifier/pkg/cfg"
	"github.com/your-org/formal-verifier/pkg/concurrency"
	"github.com/your-org/formal-verifier/pkg/dataflow"
	"github.com/your-org/formal-verifier/pkg/optimization"
	"github.com/your-org/formal-verifier/pkg/project"
	"github.com/your-org/formal-verifier/pkg/report"
	fvtypes "github.com/your-org/formal-verifier/pkg/types"
)

const version = "v1.0.0"

func main() {
	// 定义子命令
	analyzeCmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	cfgCmd := flag.NewFlagSet("cfg", flag.ExitOnError)
	concurrencyCmd := flag.NewFlagSet("concurrency", flag.ExitOnError)
	dataflowCmd := flag.NewFlagSet("dataflow", flag.ExitOnError)
	typesCmd := flag.NewFlagSet("types", flag.ExitOnError)
	optimizerCmd := flag.NewFlagSet("optimizer", flag.ExitOnError)

	// 项目分析命令参数
	analyzeDir := analyzeCmd.String("dir", ".", "项目根目录路径")
	analyzeRecursive := analyzeCmd.Bool("recursive", true, "递归扫描子目录")
	analyzeOutput := analyzeCmd.String("output", "", "输出文件路径 (留空输出到stdout)")
	analyzeFormat := analyzeCmd.String("format", "text", "输出格式: text, json, html, markdown")
	analyzeExclude := analyzeCmd.String("exclude", "", "排除模式，逗号分隔 (例如: vendor/*,testdata/*)")
	analyzeIncludeTests := analyzeCmd.Bool("include-tests", false, "包含测试文件")
	analyzeFailOnError := analyzeCmd.Bool("fail-on-error", false, "发现错误时以非零退出码退出")

	// CFG命令参数
	cfgFile := cfgCmd.String("file", "", "Go源文件路径")
	cfgOutput := cfgCmd.String("output", "cfg.dot", "输出文件路径")
	cfgSSA := cfgCmd.Bool("ssa", false, "启用SSA转换")
	cfgFormat := cfgCmd.String("format", "dot", "输出格式: dot, json, html")

	// 并发检查命令参数
	concFile := concurrencyCmd.String("file", "", "Go源文件路径")
	concCheck := concurrencyCmd.String("check", "all", "检查类型: all, goroutine-leak, deadlock, race")

	// 数据流分析命令参数
	dataflowFile := dataflowCmd.String("file", "", "Go源文件路径")
	dataflowAnalysis := dataflowCmd.String("analysis", "all", "分析类型: all, liveness, reaching, available")

	// 类型检查命令参数
	typesFile := typesCmd.String("file", "", "Go源文件路径")
	typesCheck := typesCmd.String("check", "all", "检查类型: all, safety, progress, preservation, constraints")

	// 优化分析命令参数
	optimizerFile := optimizerCmd.String("file", "", "Go源文件路径")
	optimizerCheck := optimizerCmd.String("check", "all", "检查类型: all, escape, inline, bce")

	// 检查命令
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// 处理特殊命令
	switch os.Args[1] {
	case "version", "--version", "-v":
		fmt.Printf("Go Formal Verifier %s\n", version)
		fmt.Println("基于 Go 1.25.3 形式化理论体系")
		return
	case "help", "--help", "-h":
		printUsage()
		return
	}

	// 解析子命令
	switch os.Args[1] {
	case "analyze":
		analyzeCmd.Parse(os.Args[2:])
		runProjectAnalysis(
			*analyzeDir,
			*analyzeRecursive,
			*analyzeOutput,
			*analyzeFormat,
			*analyzeExclude,
			*analyzeIncludeTests,
			*analyzeFailOnError,
		)

	case "cfg":
		cfgCmd.Parse(os.Args[2:])
		if *cfgFile == "" {
			fmt.Println("错误: 必须指定 --file 参数")
			cfgCmd.Usage()
			os.Exit(1)
		}
		runCFGAnalysis(*cfgFile, *cfgOutput, *cfgSSA, *cfgFormat)

	case "concurrency":
		concurrencyCmd.Parse(os.Args[2:])
		if *concFile == "" {
			fmt.Println("错误: 必须指定 --file 参数")
			concurrencyCmd.Usage()
			os.Exit(1)
		}
		runConcurrencyCheck(*concFile, *concCheck)

	case "dataflow":
		dataflowCmd.Parse(os.Args[2:])
		if *dataflowFile == "" {
			fmt.Println("错误: 必须指定 --file 参数")
			dataflowCmd.Usage()
			os.Exit(1)
		}
		runDataflowAnalysis(*dataflowFile, *dataflowAnalysis)

	case "types":
		typesCmd.Parse(os.Args[2:])
		if *typesFile == "" {
			fmt.Println("错误: 必须指定 --file 参数")
			typesCmd.Usage()
			os.Exit(1)
		}
		runTypesCheck(*typesFile, *typesCheck)

	case "optimizer":
		optimizerCmd.Parse(os.Args[2:])
		if *optimizerFile == "" {
			fmt.Println("错误: 必须指定 --file 参数")
			optimizerCmd.Usage()
			os.Exit(1)
		}
		runOptimizerCheck(*optimizerFile, *optimizerCheck)

	default:
		fmt.Printf("未知命令: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`Go Formal Verifier - Go形式化验证工具

用法:
  fv <command> [options]

命令:
  analyze      项目级分析 (NEW!)
  cfg          控制流图分析
  concurrency  并发安全检查
  dataflow     数据流分析
  types        类型系统验证
  optimizer    优化分析

  version      显示版本信息
  help         显示帮助信息

项目分析:
  fv analyze [options]
    --dir=<path>           项目根目录 (默认: .)
    --recursive            递归扫描子目录 (默认: true)
    --output=<path>        输出文件路径 (留空输出到stdout)
    --format=<fmt>         输出格式: text, json, html, markdown (默认: text)
    --exclude=<patterns>   排除模式，逗号分隔
    --include-tests        包含测试文件
    --fail-on-error        发现错误时以非零退出码退出

  示例:
    fv analyze --dir=./myproject
    fv analyze --dir=. --format=html --output=report.html
    fv analyze --dir=. --exclude="vendor/*,testdata/*"
    fv analyze --dir=. --fail-on-error  # 适用于CI/CD

CFG分析:
  fv cfg --file=<file> [options]
    --file=<path>      Go源文件路径 (必需)
    --output=<path>    输出文件路径 (默认: cfg.dot)
    --ssa              启用SSA转换
    --format=<fmt>     输出格式: dot, json, html (默认: dot)

  示例:
    fv cfg --file=main.go --output=cfg.dot --ssa
    dot -Tpng cfg.dot -o cfg.png

并发检查:
  fv concurrency --file=<file> --check=<type>
    --file=<path>      Go源文件路径 (必需)
    --check=<type>     检查类型: all, goroutine-leak, deadlock, race

  示例:
    fv concurrency --file=main.go --check=goroutine-leak
    fv concurrency --file=main.go --check=deadlock

数据流分析:
  fv dataflow --file=<file> --analysis=<type>
    --file=<path>      Go源文件路径 (必需)
    --analysis=<type>  分析类型: all, liveness, reaching, available

  示例:
    fv dataflow --file=main.go --analysis=liveness

类型验证:
  fv types --file=<file> --check=<type>
    --file=<path>      Go源文件路径 (必需)
    --check=<type>     检查类型: all, safety, progress, preservation, constraints

  示例:
    fv types --file=main.go --check=safety

优化分析:
  fv optimizer --file=<file> --check=<type>
    --file=<path>      Go源文件路径 (必需)
    --check=<type>     检查类型: all, escape, inline, bce

  示例:
    fv optimizer --file=main.go --check=escape

理论基础:
  本工具基于 Go 1.25.3 形式化理论体系:
  - 文档02: CSP并发模型与形式化证明
  - 文档03: Go类型系统形式化定义
  - 文档13: Go控制流形式化完整分析
  - 文档15: Go编译器优化形式化证明
  - 文档16: Go并发模式完整形式化分析

文档位置:
  docs/01-语言基础/00-Go-1.25.3形式化理论体系/

更多信息:
  https://github.com/your-org/formal-verifier
`)
}

func runCFGAnalysis(file, output string, enableSSA bool, format string) {
	fmt.Printf("🔍 CFG分析: %s\n", filepath.Base(file))
	fmt.Printf("   输出: %s\n", output)
	fmt.Printf("   SSA: %v\n", enableSSA)
	fmt.Printf("   格式: %s\n", format)
	fmt.Println()

	// Build CFG
	fmt.Println("📊 正在构造CFG...")
	builder := cfg.NewBuilder()
	cfgGraph, err := builder.BuildFromFile(file)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		os.Exit(1)
	}

	// Print statistics
	visualizer := cfg.NewVisualizer(cfgGraph)
	stats := visualizer.GetStats()

	fmt.Println("✅ CFG构造完成")
	fmt.Printf("   节点数: %d\n", stats.NodeCount)
	fmt.Printf("   边数: %d\n", stats.EdgeCount)
	fmt.Printf("   最大深度: %d\n", stats.MaxDepth)
	fmt.Printf("   循环数: %d\n", stats.LoopCount)
	fmt.Printf("   分支数: %d\n", stats.BranchCount)
	fmt.Println()

	// SSA transformation (if requested)
	if enableSSA {
		fmt.Println("🔄 正在进行SSA转换...")
		converter := cfg.NewSSAConverter(cfgGraph)
		ssaCFG, err := converter.Convert()
		if err != nil {
			fmt.Printf("❌ SSA转换失败: %v\n", err)
		} else {
			fmt.Println("✅ SSA转换完成")

			// Verify SSA property
			fmt.Println("🔍 验证SSA性质...")
			valid, errors := ssaCFG.VerifySSAProperty()
			if valid {
				fmt.Println("✅ SSA性质验证通过")
			} else {
				fmt.Printf("⚠️  发现 %d 个SSA性质违反:\n", len(errors))
				for _, err := range errors {
					fmt.Printf("   - %s\n", err)
				}
			}

			// Print SSA statistics
			fmt.Printf("   φ-函数总数: ")
			phiCount := 0
			for _, node := range ssaCFG.Nodes {
				phiCount += len(ssaCFG.SSANodes[node].PhiFunctions)
			}
			fmt.Printf("%d\n", phiCount)
		}
		fmt.Println()
	}

	// Export to requested format
	fmt.Printf("💾 正在导出到 %s 格式...\n", format)
	switch format {
	case "dot":
		err = visualizer.ExportDOT(output)
		if err != nil {
			fmt.Printf("❌ 导出失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ DOT文件已保存: %s\n", output)
		fmt.Println()
		fmt.Println("💡 使用 Graphviz 可视化:")
		fmt.Printf("   dot -Tpng %s -o %s.png\n", output, output[:len(output)-4])

	case "json":
		err = visualizer.ExportJSON(output)
		if err != nil {
			fmt.Printf("❌ 导出失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ JSON文件已保存: %s\n", output)

	case "html":
		err = visualizer.ExportHTML(output)
		if err != nil {
			fmt.Printf("❌ 导出失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ HTML文件已保存: %s\n", output)
		fmt.Println("💡 在浏览器中打开查看交互式可视化")

	default:
		fmt.Printf("❌ 未知格式: %s\n", format)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("📚 理论基础:")
	fmt.Println("   - 文档13: Go控制流形式化完整分析")
}

func runConcurrencyCheck(file, check string) {
	fmt.Printf("🔍 并发安全检查: %s\n", filepath.Base(file))
	fmt.Printf("   检查类型: %s\n", check)
	fmt.Println()

	// 创建并发分析器
	analyzer := concurrency.NewAnalyzer()

	// 分析文件
	fmt.Println("📊 正在分析...")
	err := analyzer.AnalyzeFile(file)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 分析完成")
	fmt.Println()

	// 执行指定的检查
	fmt.Println("📊 检查结果:")
	fmt.Println()

	goroutines := analyzer.GetGoroutines()
	channels := analyzer.GetChannels()
	dataRaces := analyzer.GetDataRaces()

	switch check {
	case "goroutine-leak":
		fmt.Println("   🔍 Goroutine泄露检查:")
		fmt.Println("   理论: Leak(g) ⟺ ¬CanExit(g) ∧ WaitedBy(g) = ∅")
		fmt.Println()

		leakCount := 0
		for _, g := range goroutines {
			if !g.CanExit && len(g.WaitedBy) == 0 {
				fmt.Printf("   ⚠️  泄露 (Goroutine #%d) at %s\n", g.ID, g.Position)
				leakCount++
			}
		}
		if leakCount == 0 {
			fmt.Println("   ✅ 未检测到goroutine泄露")
		} else {
			fmt.Printf("   ⚠️  检测到 %d 个潜在泄露\n", leakCount)
		}

	case "deadlock":
		fmt.Println("   🔍 Channel死锁检查:")
		fmt.Println("   理论: Deadlock(ch) ⟺ Unbuffered ∧ Sends > Receives")
		fmt.Println()

		deadlockCount := 0
		for name, ch := range channels {
			if (!ch.Buffered && len(ch.Sends) > 0 && len(ch.Receives) == 0) ||
				(ch.Buffered && len(ch.Sends) > ch.BufferSize+len(ch.Receives)) {
				fmt.Printf("   ⚠️  死锁 (%s): %d sends, %d receives\n",
					name, len(ch.Sends), len(ch.Receives))
				deadlockCount++
			}
		}
		if deadlockCount == 0 {
			fmt.Println("   ✅ 未检测到channel死锁")
		} else {
			fmt.Printf("   ⚠️  检测到 %d 个潜在死锁\n", deadlockCount)
		}

	case "race":
		fmt.Println("   🔍 数据竞争检查:")
		fmt.Println("   理论: DataRace(v) ⟺ ∃concurrent accesses ∧ ¬(a1 <HB a2)")
		fmt.Println()

		raceCount := 0
		for varName, info := range dataRaces {
			if info.IsRace {
				fmt.Printf("   ⚠️  数据竞争 (变量 %s): %d accesses\n",
					varName, len(info.Accesses))
				raceCount++
			}
		}
		if raceCount == 0 {
			fmt.Println("   ✅ 未检测到数据竞争")
		} else {
			fmt.Printf("   ⚠️  检测到 %d 个潜在数据竞争\n", raceCount)
		}

	case "all":
		// 完整报告
		report := analyzer.Report()
		fmt.Println(report)
		return

	default:
		fmt.Printf("❌ 未知检查类型: %s\n", check)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("📚 理论基础:")
	fmt.Println("   - 文档02: CSP并发模型与形式化证明")
	fmt.Println("   - 文档16: Go并发模式完整形式化分析")
	fmt.Println("   - Happens-Before关系: Go Memory Model")
}

func runDataflowAnalysis(file, analysis string) {
	fmt.Printf("🔍 数据流分析: %s\n", filepath.Base(file))
	fmt.Printf("   分析类型: %s\n", analysis)
	fmt.Println()

	// Build CFG first
	fmt.Println("📊 构造CFG...")
	builder := cfg.NewBuilder()
	cfgGraph, err := builder.BuildFromFile(file)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ CFG构造完成")
	fmt.Println()

	// Run data flow analysis
	switch analysis {
	case "liveness":
		fmt.Println("📊 活跃变量分析 (Liveness)")
		fmt.Println("   理论: OUT[n] = ⋃(s∈succ(n)) IN[s]")
		fmt.Println("        IN[n] = use[n] ∪ (OUT[n] - def[n])")
		fmt.Println()
		liveness := dataflow.NewLivenessAnalysis(cfgGraph)
		liveness.Run()

	case "reaching":
		fmt.Println("📊 可达定义分析 (Reaching Definitions)")
		fmt.Println("   理论: OUT[n] = gen[n] ∪ (IN[n] - kill[n])")
		fmt.Println("        IN[n] = ⋃(p∈pred(n)) OUT[p]")
		fmt.Println()
		reaching := dataflow.NewReachingDefinitionsAnalysis(cfgGraph)
		reaching.Run()

	case "available":
		fmt.Println("📊 可用表达式分析 (Available Expressions)")
		fmt.Println("   理论: OUT[n] = gen[n] ∪ (IN[n] - kill[n])")
		fmt.Println("        IN[n] = ⋂(p∈pred(n)) OUT[p]")
		fmt.Println()
		available := dataflow.NewAvailableExpressionsAnalysis(cfgGraph)
		available.Run()

	case "all":
		dataflow.RunAll(cfgGraph)

	default:
		fmt.Printf("❌ 未知分析类型: %s\n", analysis)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("📚 理论基础:")
	fmt.Println("   - 文档13: Go控制流形式化完整分析 (第4章)")
}

func runTypesCheck(file, check string) {
	fmt.Printf("🔍 类型系统验证: %s\n", filepath.Base(file))
	fmt.Printf("   检查类型: %s\n", check)
	fmt.Println()

	// 创建类型验证器
	verifier := fvtypes.NewVerifier()

	// 分析文件
	fmt.Println("📊 正在验证...")
	err := verifier.VerifyFile(file)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 验证完成")
	fmt.Println()

	// 执行指定的检查
	fmt.Println("📊 验证结果:")
	fmt.Println()

	progressErrors := verifier.GetProgressErrors()
	preservationErrors := verifier.GetPreservationErrors()
	constraintErrors := verifier.GetConstraintErrors()

	switch check {
	case "progress":
		fmt.Println("   🔍 Progress定理验证:")
		fmt.Println("   理论: ∀e, T. (⊢ e : T) ⟹ (value(e) ∨ ∃e'. e ↦ e')")
		fmt.Println()

		if len(progressErrors) == 0 {
			fmt.Println("   ✅ Progress定理: 验证通过")
		} else {
			fmt.Printf("   ⚠️  Progress定理: %d个违反\n", len(progressErrors))
			for i, err := range progressErrors {
				if i < 3 {
					fmt.Printf("      - %s: %s\n", err.Position, err.Message)
				}
			}
			if len(progressErrors) > 3 {
				fmt.Printf("      ... and %d more\n", len(progressErrors)-3)
			}
		}

	case "preservation":
		fmt.Println("   🔍 Preservation定理验证:")
		fmt.Println("   理论: ∀e, e', T. (⊢ e : T ∧ e ↦ e') ⟹ ⊢ e' : T")
		fmt.Println()

		if len(preservationErrors) == 0 {
			fmt.Println("   ✅ Preservation定理: 验证通过")
		} else {
			fmt.Printf("   ⚠️  Preservation定理: %d个违反\n", len(preservationErrors))
			for i, err := range preservationErrors {
				if i < 3 {
					fmt.Printf("      - %s: %s\n", err.Position, err.Message)
				}
			}
			if len(preservationErrors) > 3 {
				fmt.Printf("      ... and %d more\n", len(preservationErrors)-3)
			}
		}

	case "constraints":
		fmt.Println("   🔍 泛型约束验证:")
		fmt.Println("   理论: ∀T, C. (T : C) ⟹ satisfies(T, C)")
		fmt.Println()

		if len(constraintErrors) == 0 {
			fmt.Println("   ✅ 泛型约束: 验证通过")
		} else {
			fmt.Printf("   ⚠️  泛型约束: %d个违反\n", len(constraintErrors))
		}

	case "safety":
		fmt.Println("   🔍 类型安全性验证:")
		fmt.Println("   理论: Type Safety = Progress ∧ Preservation")
		fmt.Println()

		if verifier.IsSafe() {
			fmt.Println("   ✅ 类型安全性: 验证通过")
		} else {
			fmt.Println("   ⚠️  类型安全性: 存在违反")
			fmt.Printf("      - Progress errors: %d\n", len(progressErrors))
			fmt.Printf("      - Preservation errors: %d\n", len(preservationErrors))
		}

	case "all":
		// 完整报告
		report := verifier.Report()
		fmt.Println(report)
		return

	default:
		fmt.Printf("❌ 未知检查类型: %s\n", check)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("📚 理论基础:")
	fmt.Println("   - 文档03: Go类型系统形式化定义")
	fmt.Println("   - Progress: 良型项要么是值要么可以继续计算")
	fmt.Println("   - Preservation: 类型在计算过程中保持不变")
	fmt.Println("   - Type Safety: Progress ∧ Preservation")
}

func runOptimizerCheck(file, check string) {
	fmt.Printf("🔍 优化分析: %s\n", filepath.Base(file))
	fmt.Printf("   检查类型: %s\n", check)
	fmt.Println()

	// 创建优化分析器
	analyzer := optimization.NewAnalyzer()

	// 分析文件
	fmt.Println("📊 正在分析...")
	if err := analyzer.AnalyzeFile(file); err != nil {
		fmt.Printf("❌ 分析失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 分析完成")
	fmt.Println()

	// 根据check参数输出结果
	switch check {
	case "escape":
		printEscapeAnalysis(analyzer)
	case "inline":
		printInlineAnalysis(analyzer)
	case "bce":
		printBCEAnalysis(analyzer)
	case "all":
		// 输出完整报告
		fmt.Println(analyzer.Report())
	default:
		fmt.Printf("未知的检查类型: %s\n", check)
		os.Exit(1)
	}

	// 输出理论基础
	if check != "all" {
		fmt.Println()
		printOptimizationTheory(check)
	}
}

func printEscapeAnalysis(analyzer *optimization.OptimizerAnalyzer) {
	escapes := analyzer.GetEscapeAnalysis()

	fmt.Println("📊 逃逸分析结果:")
	fmt.Println()
	fmt.Println("   🔍 逃逸分析:")
	fmt.Println("   理论: obj escapes ⟺ ∃ reference outliving function")
	fmt.Println()

	if len(escapes) == 0 {
		fmt.Println("   ✅ 无逃逸对象")
		return
	}

	stackCount := 0
	heapCount := 0
	for _, info := range escapes {
		if info.EscapesTo == "stack" {
			stackCount++
		} else if info.EscapesTo == "heap" {
			heapCount++
		}
	}

	fmt.Printf("   📊 栈分配: %d 个 | 堆分配: %d 个\n\n", stackCount, heapCount)

	for i, info := range escapes {
		if i >= 15 {
			fmt.Printf("   ... and %d more\n", len(escapes)-15)
			break
		}

		icon := "⚠️ "
		if info.CanOptimize {
			icon = "✅"
		}

		fmt.Printf("   %s %s:%d: %s -> %s (%s)\n",
			icon,
			filepath.Base(info.Position.Filename),
			info.Position.Line,
			info.VarName,
			info.EscapesTo,
			info.Reason)
	}
}

func printInlineAnalysis(analyzer *optimization.OptimizerAnalyzer) {
	inlines := analyzer.GetInlineAnalysis()

	fmt.Println("📊 内联分析结果:")
	fmt.Println()
	fmt.Println("   🔍 函数内联分析:")
	fmt.Println("   理论: CanInline(f) ⟺ cost < 80 ∧ ¬recursive ∧ ¬complex")
	fmt.Println()

	if len(inlines) == 0 {
		fmt.Println("   ℹ️  无函数分析")
		return
	}

	canInlineCount := 0
	for _, info := range inlines {
		if info.CanInline {
			canInlineCount++
		}
	}

	fmt.Printf("   📊 可内联: %d 个 | 不可内联: %d 个\n\n",
		canInlineCount, len(inlines)-canInlineCount)

	for i, info := range inlines {
		if i >= 15 {
			fmt.Printf("   ... and %d more\n", len(inlines)-15)
			break
		}

		icon := "❌"
		if info.CanInline {
			icon = "✅"
		}

		fmt.Printf("   %s %s (cost: %d): %s\n",
			icon,
			info.FuncName,
			info.Cost,
			info.Reason)
	}
}

func printBCEAnalysis(analyzer *optimization.OptimizerAnalyzer) {
	bces := analyzer.GetBCEAnalysis()

	fmt.Println("📊 边界检查消除结果:")
	fmt.Println()
	fmt.Println("   🔍 边界检查消除 (BCE):")
	fmt.Println("   理论: CanEliminate(a[i]) ⟺ provable: 0 ≤ i < len(a)")
	fmt.Println()

	if len(bces) == 0 {
		fmt.Println("   ℹ️  无数组访问")
		return
	}

	eliminateCount := 0
	for _, info := range bces {
		if info.CanEliminate {
			eliminateCount++
		}
	}

	fmt.Printf("   📊 可消除: %d 个 | 不可消除: %d 个\n\n",
		eliminateCount, len(bces)-eliminateCount)

	for i, info := range bces {
		if i >= 15 {
			fmt.Printf("   ... and %d more\n", len(bces)-15)
			break
		}

		icon := "❌"
		if info.CanEliminate {
			icon = "✅"
		}

		fmt.Printf("   %s %s:%d: %s[%s] (%s)\n",
			icon,
			filepath.Base(info.Position.Filename),
			info.Position.Line,
			info.ArrayExpr,
			info.IndexExpr,
			info.Reason)
	}
}

func printOptimizationTheory(check string) {
	fmt.Println("📐 形式化理论基础:")

	switch check {
	case "escape":
		fmt.Println("   - 逃逸分析: obj escapes ⟺ ∃ reference outliving function")
		fmt.Println("   - 栈分配优化: lifetime(obj) ⊆ function lifetime")
	case "inline":
		fmt.Println("   - 内联分析: InlineCost < threshold ∧ ¬recursive")
		fmt.Println("   - 成本计算: Σ(instruction weights)")
	case "bce":
		fmt.Println("   - BCE: provable(0 ≤ i < len(a)) ⟹ eliminate check")
		fmt.Println("   - 安全性: eliminate ⟺ cannot panic")
	}

	fmt.Println()
	fmt.Println("📚 理论基础:")
	fmt.Println("   - 文档15: Go编译器优化形式化证明")
}

// runProjectAnalysis 执行项目级分析
func runProjectAnalysis(dir string, recursive bool, output, format, exclude string, includeTests, failOnError bool) {
	fmt.Printf("🔍 项目分析: %s\n", dir)
	fmt.Println()

	// 创建分析器
	analyzer := project.NewAnalyzer(dir)

	// 配置扫描器
	scanner := analyzer.Scanner
	scanner.WithRecursive(recursive)
	scanner.WithIncludeTests(includeTests)

	// 处理排除模式
	if exclude != "" {
		patterns := strings.Split(exclude, ",")
		for i := range patterns {
			patterns[i] = strings.TrimSpace(patterns[i])
		}
		scanner.WithExcludePatterns(patterns)
	}

	// 执行分析
	fmt.Println("📊 正在扫描和分析项目...")
	result, err := analyzer.Analyze()
	if err != nil {
		fmt.Printf("❌ 分析失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 分析完成")
	fmt.Println()

	// 根据格式输出结果
	switch format {
	case "text":
		outputTextReport(result, output)
	case "json":
		if output == "" {
			output = "analysis-report.json"
		}
		jsonReport := report.NewJSONReport(result)
		if err := jsonReport.Generate(output); err != nil {
			fmt.Printf("❌ 生成JSON报告失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ JSON报告已保存到: %s\n", output)
	case "html":
		if output == "" {
			output = "analysis-report.html"
		}
		htmlReport := report.NewHTMLReport(result)
		if err := htmlReport.Generate(output); err != nil {
			fmt.Printf("❌ 生成HTML报告失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ HTML报告已保存到: %s\n", output)
		absPath, _ := filepath.Abs(output)
		fmt.Printf("📊 在浏览器中打开: file://%s\n", absPath)
	case "markdown":
		if output == "" {
			output = "analysis-report.md"
		}
		mdReport := report.NewMarkdownReport(result)
		if err := mdReport.Generate(output); err != nil {
			fmt.Printf("❌ 生成Markdown报告失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Markdown报告已保存到: %s\n", output)
	default:
		fmt.Printf("未知的输出格式: %s\n", format)
		os.Exit(1)
	}

	// 根据failOnError参数决定退出码
	if failOnError && result.HasErrors() {
		os.Exit(1)
	}
}

// outputTextReport 输出文本格式报告
func outputTextReport(result *project.AnalysisResult, output string) {
	// 准备输出
	var writer *os.File
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			fmt.Printf("❌ 无法创建输出文件: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		writer = f
	} else {
		writer = os.Stdout
	}

	// 输出标题
	fmt.Fprintln(writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(writer, "       Go Formal Verifier - 项目分析报告")
	fmt.Fprintln(writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(writer)

	// 输出摘要
	fmt.Fprintln(writer, result.Summary)
	fmt.Fprintln(writer)

	// 输出统计
	fmt.Fprintln(writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(writer, "统计信息")
	fmt.Fprintln(writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintf(writer, "文件数: %d\n", result.Stats.TotalFiles)
	fmt.Fprintf(writer, "代码行数: %d\n", result.Stats.TotalLines)
	fmt.Fprintf(writer, "总问题数: %d\n", result.Stats.TotalIssues)
	fmt.Fprintf(writer, "  - 错误: %d\n", result.Stats.ErrorCount)
	fmt.Fprintf(writer, "  - 警告: %d\n", result.Stats.WarningCount)
	fmt.Fprintf(writer, "  - 提示: %d\n", result.Stats.InfoCount)
	fmt.Fprintln(writer)

	// 按类别统计
	fmt.Fprintln(writer, "问题分布:")
	fmt.Fprintf(writer, "  - 并发问题: %d\n", result.Stats.ConcurrencyIssues)
	fmt.Fprintf(writer, "  - 类型问题: %d\n", result.Stats.TypeIssues)
	fmt.Fprintf(writer, "  - 数据流问题: %d\n", result.Stats.DataFlowIssues)
	fmt.Fprintf(writer, "  - 优化建议: %d\n", result.Stats.OptimizationIssues)
	fmt.Fprintln(writer)

	// 如果没有问题，直接返回
	if result.Stats.TotalIssues == 0 {
		fmt.Fprintln(writer, "✅ 没有发现问题！代码质量优秀。")
		return
	}

	// 输出详细问题
	fmt.Fprintln(writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(writer, "问题详情")
	fmt.Fprintln(writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(writer)

	// 按严重程度输出
	if result.Stats.ErrorCount > 0 {
		fmt.Fprintln(writer, "❌ 错误:")
		errors := result.GetIssuesBySeverity("error")
		for _, issue := range errors {
			printIssue(writer, issue)
		}
		fmt.Fprintln(writer)
	}

	if result.Stats.WarningCount > 0 {
		fmt.Fprintln(writer, "⚠️  警告:")
		warnings := result.GetIssuesBySeverity("warning")
		// 限制显示数量
		maxDisplay := 20
		for i, issue := range warnings {
			if i >= maxDisplay {
				fmt.Fprintf(writer, "... 还有 %d 个警告（使用 --format=json 查看全部）\n", len(warnings)-maxDisplay)
				break
			}
			printIssue(writer, issue)
		}
		fmt.Fprintln(writer)
	}

	if result.Stats.InfoCount > 0 && result.Stats.InfoCount <= 10 {
		fmt.Fprintln(writer, "ℹ️  提示:")
		infos := result.GetIssuesBySeverity("info")
		for _, issue := range infos {
			printIssue(writer, issue)
		}
		fmt.Fprintln(writer)
	} else if result.Stats.InfoCount > 10 {
		fmt.Fprintf(writer, "ℹ️  提示: %d 个（使用 --format=json 查看详情）\n\n", result.Stats.InfoCount)
	}

	// 输出总结
	fmt.Fprintln(writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(writer, "分析总结")
	fmt.Fprintln(writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintf(writer, "质量评分: %d/100\n", result.Stats.QualityScore)

	if result.Stats.QualityScore >= 90 {
		fmt.Fprintln(writer, "✅ 代码质量优秀！")
	} else if result.Stats.QualityScore >= 70 {
		fmt.Fprintln(writer, "✓ 代码质量良好")
	} else if result.Stats.QualityScore >= 50 {
		fmt.Fprintln(writer, "⚠️  代码质量需要改进")
	} else {
		fmt.Fprintln(writer, "❌ 代码质量较差 - 需要立即关注")
	}

	if output != "" {
		fmt.Printf("✅ 报告已保存到: %s\n", output)
	}
}

// printIssue 打印单个问题
func printIssue(writer *os.File, issue project.Issue) {
	fmt.Fprintf(writer, "  [%s] %s:%d:%d\n",
		issue.Category,
		filepath.Base(issue.File),
		issue.Line,
		issue.Column)
	fmt.Fprintf(writer, "    %s\n", issue.Message)
	if issue.Suggestion != "" {
		fmt.Fprintf(writer, "    💡 建议: %s\n", issue.Suggestion)
	}
	fmt.Fprintln(writer)
}
