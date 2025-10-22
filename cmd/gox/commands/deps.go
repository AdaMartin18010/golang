package commands

import (
	"fmt"
	"os"
	"os/exec"
)

// DepsCommand 依赖管理命令
func DepsCommand(args []string) error {
	if len(args) == 0 {
		return showDeps()
	}

	action := args[0]

	switch action {
	case "list":
		return showDeps()
	case "tidy":
		return tidyDeps()
	case "update":
		return updateDeps()
	case "verify":
		return verifyDeps()
	case "graph":
		return graphDeps()
	default:
		return fmt.Errorf("未知操作: %s", action)
	}
}

func showDeps() error {
	fmt.Println("📦 依赖列表...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	cmd := exec.Command("go", "list", "-m", "all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func tidyDeps() error {
	fmt.Println("🧹 整理依赖...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("依赖整理失败: %w", err)
	}

	fmt.Println("✅ 依赖整理完成")
	return nil
}

func updateDeps() error {
	fmt.Println("⬆️  更新依赖...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	cmd := exec.Command("go", "get", "-u", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("依赖更新失败: %w", err)
	}

	// 自动tidy
	return tidyDeps()
}

func verifyDeps() error {
	fmt.Println("✅ 验证依赖...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	cmd := exec.Command("go", "mod", "verify")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("依赖验证失败: %w", err)
	}

	fmt.Println("✅ 依赖验证通过")
	return nil
}

func graphDeps() error {
	fmt.Println("📊 依赖关系图...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	cmd := exec.Command("go", "mod", "graph")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
