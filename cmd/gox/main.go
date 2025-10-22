package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const version = "v1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "quality", "q":
		runQuality(args)
	case "test", "t":
		runTest(args)
	case "stats", "s":
		runStats(args)
	case "format", "f":
		runFormat(args)
	case "docs", "d":
		runDocs(args)
	case "migrate", "m":
		runMigrate(args)
	case "verify", "v":
		runVerify(args)
	case "help", "h", "--help", "-h":
		printHelp()
	case "version", "--version", "-v":
		printVersion()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
gox - Golang项目管理工具

使用方式:
  gox <command> [options]

常用命令:
  quality, q     代码质量检查
  test, t        运行测试并生成报告
  stats, s       项目统计分析
  format, f      代码格式化
  docs, d        文档处理
  migrate, m     项目迁移
  verify, v      结构验证
  help, h        显示帮助信息
  version        显示版本信息

示例:
  gox quality           运行质量检查
  gox test              运行所有测试
  gox stats             查看项目统计
  gox format --check    检查代码格式

详细帮助: gox help
`)
}

func printHelp() {
	fmt.Println(`
gox - Golang项目管理工具 v1.0.0

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 quality (q) - 代码质量检查
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  执行完整的代码质量检查，包括:
  ✅ go fmt   - 代码格式检查
  ✅ go vet   - 静态分析
  ✅ go build - 编译检查
  ✅ go test  - 测试运行

  使用:
    gox quality           完整检查
    gox quality --fast    快速检查(跳过测试)
    gox q                 简写形式

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🧪 test (t) - 测试统计
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  运行测试并生成统计报告

  使用:
    gox test              运行所有测试
    gox test --coverage   生成覆盖率报告
    gox test --verbose    详细输出
    gox t                 简写形式

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📊 stats (s) - 项目统计
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  分析并展示项目统计信息

  使用:
    gox stats             显示项目统计
    gox stats --detail    详细统计
    gox s                 简写形式

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

💻 format (f) - 代码格式化
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  格式化Go代码

  使用:
    gox format            格式化所有代码
    gox format --check    只检查不格式化
    gox f                 简写形式

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📝 docs (d) - 文档处理
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  文档生成和处理

  使用:
    gox docs toc          生成文档目录
    gox docs links        检查文档链接
    gox docs format       格式化文档
    gox d                 简写形式

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔄 migrate (m) - 项目迁移
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  执行Workspace迁移

  使用:
    gox migrate --dry-run 预览迁移
    gox migrate           执行迁移
    gox m                 简写形式

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ verify (v) - 结构验证
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  验证项目结构

  使用:
    gox verify            验证项目结构
    gox verify workspace  验证Workspace配置
    gox v                 简写形式

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

全局选项:
  --verbose, -v         详细输出
  --quiet, -q           安静模式
  --help, -h            显示帮助

示例:
  gox quality           运行质量检查
  gox test --coverage   运行测试并生成覆盖率
  gox stats --detail    查看详细统计
  gox format --check    检查代码格式
  gox verify workspace  验证Workspace配置

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`)
}

func printVersion() {
	fmt.Printf("gox version %s\n", version)
	fmt.Printf("go version %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// runQuality 运行代码质量检查
func runQuality(args []string) {
	fmt.Println("🔍 开始代码质量检查...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fast := contains(args, "--fast")

	// 1. go fmt检查
	fmt.Println("\n📋 检查代码格式 (go fmt)...")
	if err := runCommand("go", "fmt", "./..."); err != nil {
		fmt.Printf("❌ 格式检查失败: %v\n", err)
	} else {
		fmt.Println("✅ 格式检查通过")
	}

	// 2. go vet检查
	fmt.Println("\n📋 静态分析 (go vet)...")
	if err := runCommand("go", "vet", "./..."); err != nil {
		fmt.Printf("⚠️  发现警告: %v\n", err)
	} else {
		fmt.Println("✅ 静态分析通过")
	}

	// 3. 编译检查
	fmt.Println("\n📋 编译检查...")
	if err := runCommand("go", "build", "./..."); err != nil {
		fmt.Printf("❌ 编译失败: %v\n", err)
	} else {
		fmt.Println("✅ 编译成功")
	}

	// 4. 测试运行 (如果不是fast模式)
	if !fast {
		fmt.Println("\n📋 运行测试...")
		if err := runCommand("go", "test", "./..."); err != nil {
			fmt.Printf("❌ 测试失败: %v\n", err)
		} else {
			fmt.Println("✅ 测试通过")
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ 质量检查完成！")
}

// runTest 运行测试
func runTest(args []string) {
	fmt.Println("🧪 运行测试...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	coverage := contains(args, "--coverage")
	verbose := contains(args, "--verbose")

	testArgs := []string{"test"}

	if verbose {
		testArgs = append(testArgs, "-v")
	}

	if coverage {
		testArgs = append(testArgs, "-cover", "-coverprofile=coverage.out")
	}

	testArgs = append(testArgs, "./...")

	if err := runCommand("go", testArgs...); err != nil {
		fmt.Printf("❌ 测试失败: %v\n", err)
		os.Exit(1)
	}

	if coverage {
		fmt.Println("\n📊 生成覆盖率报告...")
		if err := runCommand("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html"); err != nil {
			fmt.Printf("⚠️  无法生成HTML报告: %v\n", err)
		} else {
			fmt.Println("✅ 覆盖率报告已生成: coverage.html")
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ 测试完成！")
}

// runStats 运行项目统计
func runStats(args []string) {
	fmt.Println("📊 项目统计分析...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 检查project_stats工具是否存在
	statsPath := filepath.Join("scripts", "project_stats")
	if _, err := os.Stat(statsPath); err == nil {
		fmt.Println("\n使用Go工具进行统计...")
		if err := runCommandInDir(statsPath, "go", "run", "main.go"); err != nil {
			fmt.Printf("⚠️  统计失败: %v\n", err)
		}
	} else {
		// 简单统计
		fmt.Println("\n📁 文件统计:")
		runCommand("find", ".", "-type", "f", "-name", "*.go", "|", "wc", "-l")

		fmt.Println("\n📝 文档统计:")
		runCommand("find", ".", "-type", "f", "-name", "*.md", "|", "wc", "-l")
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ 统计完成！")
}

// runFormat 格式化代码
func runFormat(args []string) {
	check := contains(args, "--check")

	if check {
		fmt.Println("🔍 检查代码格式...")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		if err := runCommand("gofmt", "-l", "."); err != nil {
			fmt.Printf("❌ 发现格式问题: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ 代码格式检查通过")
	} else {
		fmt.Println("💻 格式化代码...")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		if err := runCommand("go", "fmt", "./..."); err != nil {
			fmt.Printf("❌ 格式化失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ 代码格式化完成")
	}
}

// runDocs 文档处理
func runDocs(args []string) {
	if len(args) == 0 {
		fmt.Println("请指定文档操作: toc, links, format")
		return
	}

	operation := args[0]
	fmt.Printf("📝 文档处理: %s\n", operation)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	switch operation {
	case "toc":
		fmt.Println("生成文档目录...")
		// TODO: 实现TOC生成
		fmt.Println("⏳ 功能开发中...")
	case "links":
		fmt.Println("检查文档链接...")
		// TODO: 实现链接检查
		fmt.Println("⏳ 功能开发中...")
	case "format":
		fmt.Println("格式化文档...")
		// TODO: 实现文档格式化
		fmt.Println("⏳ 功能开发中...")
	default:
		fmt.Printf("未知操作: %s\n", operation)
	}
}

// runMigrate 运行迁移
func runMigrate(args []string) {
	dryRun := contains(args, "--dry-run")

	fmt.Println("🔄 Workspace迁移...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	scriptPath := filepath.Join("scripts", "migrate-to-workspace.ps1")

	if dryRun {
		fmt.Println("预览模式 (不会实际修改文件)")
		if err := runPowerShell(scriptPath, "-DryRun"); err != nil {
			fmt.Printf("❌ 迁移预览失败: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("⚠️  即将执行实际迁移，请确认已备份！")
		fmt.Print("继续? [y/N]: ")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("❌ 迁移已取消")
			return
		}

		if err := runPowerShell(scriptPath); err != nil {
			fmt.Printf("❌ 迁移失败: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("✅ 迁移完成！")
}

// runVerify 验证项目结构
func runVerify(args []string) {
	target := "structure"
	if len(args) > 0 {
		target = args[0]
	}

	fmt.Printf("✅ 验证: %s\n", target)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	switch target {
	case "structure":
		scriptPath := filepath.Join("scripts", "verify_structure.ps1")
		if err := runPowerShell(scriptPath); err != nil {
			fmt.Printf("❌ 验证失败: %v\n", err)
			os.Exit(1)
		}
	case "workspace":
		scriptPath := filepath.Join("scripts", "verify-workspace.ps1")
		if err := runPowerShell(scriptPath); err != nil {
			fmt.Printf("❌ 验证失败: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("未知验证目标: %s\n", target)
		os.Exit(1)
	}

	fmt.Println("✅ 验证通过！")
}

// 辅助函数

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommandInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runPowerShell(scriptPath string, args ...string) error {
	psArgs := []string{"-ExecutionPolicy", "Bypass", "-File", scriptPath}
	psArgs = append(psArgs, args...)
	return runCommand("powershell", psArgs...)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
