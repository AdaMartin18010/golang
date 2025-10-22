package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// DoctorCommand 健康检查命令
func DoctorCommand(args []string) error {
	fmt.Println("🏥 系统健康检查...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	issues := 0

	// 1. Go版本检查
	fmt.Println("\n📋 Go环境检查")
	fmt.Println("─────────────────────────────────────────")
	if err := checkCommand("go", "version"); err != nil {
		fmt.Println("❌ Go未安装或不在PATH中")
		issues++
	} else {
		fmt.Printf("✅ Go版本: %s\n", runtime.Version())
		fmt.Printf("   GOOS: %s, GOARCH: %s\n", runtime.GOOS, runtime.GOARCH)
	}

	// 2. Git检查
	fmt.Println("\n📋 Git检查")
	fmt.Println("─────────────────────────────────────────")
	if err := checkCommand("git", "version"); err != nil {
		fmt.Println("❌ Git未安装")
		issues++
	} else {
		fmt.Println("✅ Git已安装")
	}

	// 3. 项目结构检查
	fmt.Println("\n📋 项目结构检查")
	fmt.Println("─────────────────────────────────────────")
	requiredFiles := []string{"go.mod", "go.work", "README.md"}
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); err == nil {
			fmt.Printf("✅ %s 存在\n", file)
		} else {
			fmt.Printf("⚠️  %s 不存在\n", file)
		}
	}

	requiredDirs := []string{"pkg", "cmd", "docs"}
	for _, dir := range requiredDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			fmt.Printf("✅ %s/ 目录存在\n", dir)
		} else {
			fmt.Printf("⚠️  %s/ 目录不存在\n", dir)
		}
	}

	// 4. Go模块检查
	fmt.Println("\n📋 Go模块检查")
	fmt.Println("─────────────────────────────────────────")
	if err := runQuietCommand("go", "mod", "verify"); err != nil {
		fmt.Println("⚠️  Go模块验证失败")
		issues++
	} else {
		fmt.Println("✅ Go模块验证通过")
	}

	// 5. 工具链检查
	fmt.Println("\n📋 开发工具检查")
	fmt.Println("─────────────────────────────────────────")

	tools := map[string][]string{
		"gofmt":         {"gofmt", "-h"},
		"go vet":        {"go", "vet", "-h"},
		"golangci-lint": {"golangci-lint", "--version"},
		"gopls":         {"gopls", "version"},
	}

	for name, cmd := range tools {
		if err := checkCommand(cmd[0], cmd[1:]...); err != nil {
			fmt.Printf("⚠️  %s 未安装\n", name)
		} else {
			fmt.Printf("✅ %s 已安装\n", name)
		}
	}

	// 6. 编译检查
	fmt.Println("\n📋 编译检查")
	fmt.Println("─────────────────────────────────────────")
	if err := runQuietCommand("go", "build", "./..."); err != nil {
		fmt.Println("❌ 编译失败")
		issues++
	} else {
		fmt.Println("✅ 编译成功")
	}

	// 7. 测试检查
	fmt.Println("\n📋 测试检查")
	fmt.Println("─────────────────────────────────────────")
	if err := runQuietCommand("go", "test", "-short", "./..."); err != nil {
		fmt.Println("⚠️  部分测试失败")
	} else {
		fmt.Println("✅ 测试通过")
	}

	// 总结
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if issues == 0 {
		fmt.Println("✅ 系统健康状态良好！")
	} else {
		fmt.Printf("⚠️  发现 %d 个问题，请检查\n", issues)
	}

	return nil
}

func checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func runQuietCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
