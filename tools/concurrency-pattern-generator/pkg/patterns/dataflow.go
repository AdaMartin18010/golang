// Package patterns - 数据流模式
package patterns

import "fmt"

// GenerateProducerConsumer 生成生产者消费者模式
func GenerateProducerConsumer(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"context\"\n\n// Producer 生产者\nfunc Producer(ctx context.Context, n int) <-chan int {\n\tch := make(chan int)\n\tgo func() {\n\t\tdefer close(ch)\n\t\tfor i := 0; i < n; i++ {\n\t\t\tselect {\n\t\t\tcase ch <- i:\n\t\t\tcase <-ctx.Done():\n\t\t\t\treturn\n\t\t\t}\n\t\t}\n\t}()\n\treturn ch\n}\n\n// Consumer 消费者\nfunc Consumer(ctx context.Context, ch <-chan int) {\n\tfor {\n\t\tselect {\n\t\tcase val, ok := <-ch:\n\t\t\tif !ok {\n\t\t\t\treturn\n\t\t\t}\n\t\t\t// Process val\n\t\t\t_ = val\n\t\tcase <-ctx.Done():\n\t\t\treturn\n\t\t}\n\t}\n}\n", pkg)
}

// GenerateBufferedChannel 生成缓冲channel模式
func GenerateBufferedChannel(pkg string) string {
	return fmt.Sprintf("package %s\n\n// BufferedChannel 缓冲channel示例\nfunc BufferedChannel() {\n\t// 创建容量为10的缓冲channel\n\tch := make(chan int, 10)\n\t\n\t// 非阻塞发送（当缓冲区未满时）\n\tfor i := 0; i < 10; i++ {\n\t\tch <- i\n\t}\n\t\n\t// 接收数据\n\tfor i := 0; i < 10; i++ {\n\t\t<-ch\n\t}\n}\n", pkg)
}

// GenerateUnbufferedChannel 生成无缓冲channel模式
func GenerateUnbufferedChannel(pkg string) string {
	return fmt.Sprintf("package %s\n\n// UnbufferedChannel 无缓冲channel示例\nfunc UnbufferedChannel() {\n\tch := make(chan int)\n\t\n\tgo func() {\n\t\tch <- 42 // 阻塞直到接收\n\t}()\n\t\n\tval := <-ch // 阻塞直到发送\n\t_ = val\n}\n", pkg)
}

// GenerateSelectPattern 生成select模式
func GenerateSelectPattern(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"time\"\n\n// SelectPattern select多路复用\nfunc SelectPattern(ch1, ch2 <-chan int, done <-chan struct{}) {\n\tfor {\n\t\tselect {\n\t\tcase v1 := <-ch1:\n\t\t\t// Handle ch1\n\t\t\t_ = v1\n\t\tcase v2 := <-ch2:\n\t\t\t// Handle ch2\n\t\t\t_ = v2\n\t\tcase <-done:\n\t\t\treturn\n\t\tcase <-time.After(time.Second):\n\t\t\t// Timeout\n\t\t\treturn\n\t\t}\n\t}\n}\n", pkg)
}

// GenerateForSelectLoop 生成for-select循环模式
func GenerateForSelectLoop(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"time\"\n\n// ForSelectLoop for-select事件循环\nfunc ForSelectLoop(input <-chan int, done <-chan struct{}) {\n\tfor {\n\t\tselect {\n\t\tcase val, ok := <-input:\n\t\t\tif !ok {\n\t\t\t\treturn\n\t\t\t}\n\t\t\t// Process val\n\t\t\t_ = val\n\t\tcase <-done:\n\t\t\treturn\n\t\t}\n\t}\n}\n", pkg)
}

// GenerateDoneChannel 生成done channel模式
func GenerateDoneChannel(pkg string) string {
	return fmt.Sprintf("package %s\n\n// DoneChannel done channel模式\nfunc DoneChannel() {\n\tdone := make(chan struct{})\n\t\n\tgo func() {\n\t\t// Do work\n\t\tfor {\n\t\t\tselect {\n\t\t\tcase <-done:\n\t\t\t\treturn\n\t\t\tdefault:\n\t\t\t\t// Continue working\n\t\t\t}\n\t\t}\n\t}()\n\t\n\t// Signal done\n\tclose(done)\n}\n", pkg)
}

// GenerateErrorChannel 生成error channel模式
func GenerateErrorChannel(pkg string) string {
	return fmt.Sprintf("package %s\n\n// ErrorChannel error channel模式\nfunc ErrorChannel() error {\n\terrChan := make(chan error, 1)\n\t\n\tgo func() {\n\t\t// Do work\n\t\terrChan <- nil // or error\n\t}()\n\t\n\treturn <-errChan\n}\n\n// MultiError 收集多个错误\nfunc MultiError(n int) []error {\n\terrChan := make(chan error, n)\n\t\n\tfor i := 0; i < n; i++ {\n\t\tgo func() {\n\t\t\terrChan <- nil\n\t\t}()\n\t}\n\t\n\terrors := make([]error, 0, n)\n\tfor i := 0; i < n; i++ {\n\t\tif err := <-errChan; err != nil {\n\t\t\terrors = append(errors, err)\n\t\t}\n\t}\n\treturn errors\n}\n", pkg)
}
