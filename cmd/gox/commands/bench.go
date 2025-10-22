package commands

import (
	"fmt"
	"os"
	"os/exec"
)

// BenchCommand 基准测试命令
func BenchCommand(args []string) error {
	fmt.Println("⚡ 运行基准测试...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	benchArgs := []string{"test", "-bench=.", "-benchmem"}

	// 解析参数
	for _, arg := range args {
		switch arg {
		case "--cpu":
			benchArgs = append(benchArgs, "-cpu=1,2,4,8")
		case "--count":
			benchArgs = append(benchArgs, "-count=5")
		case "--time":
			benchArgs = append(benchArgs, "-benchtime=10s")
		default:
			benchArgs = append(benchArgs, arg)
		}
	}

	benchArgs = append(benchArgs, "./...")

	cmd := exec.Command("go", benchArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("基准测试失败: %w", err)
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ 基准测试完成！")
	return nil
}
