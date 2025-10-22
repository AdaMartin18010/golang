package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(coverageCmd)
}

var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "Generate test coverage report",
	Long:  `Generates test coverage report for all packages and displays summary.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📊 Generating coverage report...")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		
		packages := []string{
			"./pkg/agent/...",
			"./pkg/concurrency/...",
			"./pkg/http3/...",
			"./pkg/memory/...",
		}
		
		fmt.Println("\n🧪 Running tests with coverage...")
		
		for _, pkg := range packages {
			fmt.Printf("\n📦 Coverage for %s\n", pkg)
			
			coverageFile := "coverage-" + extractPkgName(pkg) + ".out"
			testArgs := []string{"test", "-cover", "-coverprofile=" + coverageFile, pkg}
			
			execCmd := exec.Command("go", testArgs...)
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			
			if err := execCmd.Run(); err != nil {
				fmt.Printf("⚠️  Warning: Coverage generation failed for %s\n", pkg)
				continue
			}
			
			// 显示覆盖率详情
			if _, err := os.Stat(coverageFile); err == nil {
				funcCmd := exec.Command("go", "tool", "cover", "-func="+coverageFile)
				funcCmd.Stdout = os.Stdout
				funcCmd.Stderr = os.Stderr
				funcCmd.Run()
				
				// 清理临时文件
				os.Remove(coverageFile)
			}
		}
		
		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("✅ Coverage report generated!")
		fmt.Println("\n💡 Tip: Use 'go test -coverprofile=coverage.out ./...' for detailed report")
		fmt.Println("   Then: 'go tool cover -html=coverage.out' to view in browser")
	},
}

func extractPkgName(pkg string) string {
	// 从 "./pkg/agent/..." 提取 "agent"
	parts := []string{}
	for _, part := range []rune(pkg) {
		if part != '.' && part != '/' {
			parts = append(parts, string(part))
		}
	}
	
	result := ""
	for i := 0; i < len(parts) && i < 10; i++ {
		result += parts[i]
	}
	
	if result == "" {
		result = "coverage"
	}
	
	return result
}

