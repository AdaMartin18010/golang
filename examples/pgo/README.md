# PGO 示例

运行：

```bash
# 运行基线
go run ./examples/pgo

# 按需采集 CPU profile（若你的服务暴露了 pprof）
# 将 cpu.pprof 重命名为 default.pgo 并放置于 main 包目录
# 再次构建/运行，Go 工具链会使用 PGO 优化
```
