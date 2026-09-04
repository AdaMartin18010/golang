# Go 1.27 新特性示例

本目录演示 **Go 1.27 / 1.27.1** 的正式新特性，需要 Go ≥ 1.27。

权威参考：[Go 1.27 Release Notes](https://go.dev/doc/go1.27)

概念权威页：泛型方法 → [`LD-037`](../../go-knowledge-base/02-Language-Design/LD-037-Go-1.27-Generic-Methods.md)；版本演进 → [`03-Evolution/07`](../../go-knowledge-base/02-Language-Design/03-Evolution/07-Go126-to-Go127.md)

## 运行

```bash
cd examples/go127-features
go run .        # 运行全部演示
go test ./...   # 运行测试（含 synctest 气泡演示）
```

## 演示清单

| 特性 | 文件位置 | 说明 |
|------|---------|------|
| 泛型方法 | `main.go` Feature 1 | 方法可声明自己的类型参数（proposal #77273）；接口方法仍不允许。含标准库 `math/rand/v2.Rand.N[Int]` |
| `uuid` 新包 | `main.go` Feature 2 | 顶层新包：`NewV7` / `NewV4` / `Parse` / `MustParse` / `Max` |
| `strings.CutLast` / `bytes.CutLast` | `main.go` Feature 3 | 按最后一次出现位置切分 |
| `url.URL.Clone` / `url.Values.Clone` | `main.go` Feature 4 | 深拷贝辅助函数 |
| `encoding/json/v2` | `main.go` Feature 5 | GA。Marshal/Unmarshal + Options（`AllowDuplicateNames`、`RejectUnknownMembers`）；v1 `encoding/json` 已改由 v2 实现 |
| `httptest.NewTestServer` + `testing/synctest` | `synctest_test.go` | 内存网络测试服务器 + 测试气泡，无需真实端口与 sleep |
| container-aware GOMAXPROCS | `main.go` Feature 7 | Go 1.27 起默认感知 cgroup 配额（`GODEBUG=containermaxprocs=0` 关闭） |
| goroutineleak profile | `main.go` Feature 8 | Go 1.27 GA，端点 `/debug/pprof/goroutineleak` |
| `crypto/mldsa` | `advanced.go` Feature 9 | FIPS 204 后量子签名（ML-DSA-44/65/87） |
| `math/big.Int.Divide` | `advanced.go` Feature 10 | 四种舍入模式（EUCLID/_floor/_ceil/_trunc）的整数除法 |
| `hash/maphash.Hasher[T]` | `advanced.go` Feature 11 | 类型化哈希器 / `ComparableHasher[T]` |
| `simd` / `archsimd` | `advanced.go` Feature 12 | 实验性 SIMD（需 `GOEXPERIMENT=simd`，无实验环境时优雅降级） |
| `testing/cryptotest` | `cryptotest_test.go` | 标准库提供的加密实现正确性测试框架（以 AES-GCM 为例） |

## 其他 1.27 重要变更（未在 demo 中展开）

- 小对象分配提速（<80B，最高 30%，`GOEXPERIMENT=nosizespecializedmalloc` 可关，1.28 移除开关）
- `go test` 默认启用 `stdversion` vet 检查
- `go fix` 新 modernizer：`atomictypes`、`embedlit`、`slicesbackward`、`unsafefuncs`
- 新包：`crypto/mldsa`（FIPS 204 后量子签名）、`encoding/json/jsontext`、`simd`/`archsimd`（实验）
- 新 API：`math/big.Int.Divide`、`hash/maphash.Hasher[T]`、`go/types.Hasher`、`testing/synctest.Sleep`、`net/http.Server.MaxHeaderValueCount` 等
- 行为变化：HTTP/1 Response.Body 关闭自动 drain、HTTP/2 RFC 9218 客户端优先级默认尊重、`flate` 输出字节可能变化、Unicode 15 → 17
- 平台：macOS 最低 13 Ventura；`asynctimerchan` GODEBUG 永久移除

详见 [docs/01-Go-1.27-Comprehensive-Knowledge-System-2026.md](../../docs/go127-comprehensive-guide/01-Go-1.27-Comprehensive-Knowledge-System-2026.md)。
