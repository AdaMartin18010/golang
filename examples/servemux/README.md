# servemux

`net/http` ServeMux 路由模式测试驱动示例（Go 1.22+ 的增强路由模式：通配符、HTTP 方法、路径变量等）。

本目录只有 `servemux_test.go`，无 main 程序，通过测试用例演示用法。

## 运行

```bash
cd examples/servemux
GOWORK=off go test ./...
```

本目录无独立 go.mod，测试归入上级 `examples` 模块（`example.com/golang-examples/servemux`）。
