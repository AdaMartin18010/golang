// Package main demonstrates Go 1.27 new features
//
// 演示内容包括：
// 1. 泛型方法（含标准库 math/rand/v2.Rand.N[Int]）(Go 1.27)
// 2. 新标准库包 uuid：NewV7 / NewV4 / Parse (Go 1.27)
// 3. strings.CutLast / bytes.CutLast (Go 1.27)
// 4. net/url.URL.Clone 与 url.Values.Clone (Go 1.27)
// 5. encoding/json/v2：Marshal/Unmarshal 与 Options (Go 1.27 GA)
// 6. httptest.NewTestServer 与 testing/synctest (Go 1.27 / synctest 自 1.25 GA)
// 7. container-aware GOMAXPROCS（默认感知 cgroup）(Go 1.27)
// 8. goroutineleak profile（GA，/debug/pprof/goroutineleak）(Go 1.27)
package main

import (
	"bytes"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"uuid"
)

// ============================================================================
// Feature 1: Generic methods (Go 1.27)
// ============================================================================

// Stack is a simple stack with a generic method (Go 1.27).
// Go 1.27 之前，方法不允许有自己的类型参数；现在普通（非接口）类型的方法可以声明类型参数。
type Stack[T any] struct {
	items []T
}

// Push 是带接收者类型参数的方法的普通示例（T 来自类型声明，非方法级）。
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Drain 是 Go 1.27 新增的“方法级类型参数”示例：
// 方法自己声明了与接收者类型无关的类型参数 U。
// 注意：接口类型的方法仍然不允许类型参数（interface method type parameters 仍未接受），
// 即无法写出 type Foo interface { Bar[T any](x T) }。
func (s *Stack[T]) Drain[U any](convert func(T) U) []U {
	out := make([]U, 0, len(s.items))
	for _, v := range s.items {
		out = append(out, convert(v))
	}
	s.items = nil
	return out
}

// demonstrateGenericMethods shows Go 1.27 generic methods
func demonstrateGenericMethods() {
	fmt.Println("=== Feature 1: Generic methods ===")
	fmt.Println("Go 1.27: methods may now declare their own type parameters.")
	fmt.Println("Note: interface methods still cannot have type parameters.")
	fmt.Println()

	s := Stack[int]{}
	s.Push(1)
	s.Push(2)
	s.Push(3)

	// 方法级类型参数由实参推导（convert 为 func(int) string，故 U 推导为 string）
	strs := s.Drain(func(v int) string { return fmt.Sprintf("#%d", v) })
	fmt.Printf("Stack.Drain -> %v (stack now empty: %v)\n", strs, len(s.items) == 0)

	// 标准库 math/rand/v2 新增泛型方法 Rand.N[Int]：随机 [0, n) 的任意整数类型
	r := rand.New(rand.NewPCG(42, 1024))
	n32 := r.N(uint32(100))       // 返回 uint32
	n64 := r.N(int64(1_000_000))  // 返回 int64
	fmt.Printf("rand.N[uint32](100) = %d, rand.N[int64](1000000) = %d\n", n32, n64)
	fmt.Println()
}

// ============================================================================
// Feature 2: New uuid package (Go 1.27)
// ============================================================================

// demonstrateUUID shows the new top-level uuid package
func demonstrateUUID() {
	fmt.Println("=== Feature 2: uuid package ===")

	// UUID v7：时间戳 + 随机数，适合数据库主键（大致单调）
	id7 := uuid.NewV7()
	fmt.Printf("NewV7: %s (version %d)\n", id7, id7[6]>>4)

	// UUID v4：纯随机
	id4 := uuid.NewV4()
	fmt.Printf("NewV4: %s (version %d)\n", id4, id4[6]>>4)

	// Parse：解析标准字符串形式
	parsed, err := uuid.Parse(id7.String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parse: %s, equals original: %t\n", parsed, parsed == id7)

	// 非法输入会返回错误
	if _, err := uuid.Parse("not-a-uuid"); err != nil {
		fmt.Printf("Parse invalid -> error: %v\n", err)
	}
	fmt.Println()
}

// ============================================================================
// Feature 3: strings.CutLast / bytes.CutLast (Go 1.27)
// ============================================================================

// demonstrateCutLast shows CutLast for strings and bytes
func demonstrateCutLast() {
	fmt.Println("=== Feature 3: strings.CutLast / bytes.CutLast ===")

	// 按“最后一次”出现的位置切分（CutPrefix/CutSuffix 之外的最常见补全）
	path := "/usr/local/go/bin/go"
	dir, file, found := strings.CutLast(path, "/")
	fmt.Printf("strings.CutLast(%q, \"/\") = dir=%q file=%q found=%t\n", path, dir, file, found)

	host := "example.com"
	domain, tld, found := strings.CutLast(host, ".")
	fmt.Printf("strings.CutLast(%q, \".\") = domain=%q tld=%q found=%t\n", host, domain, tld, found)

	before, after, found := bytes.CutLast([]byte("key=value=extra"), []byte("="))
	fmt.Printf("bytes.CutLast -> before=%q after=%q found=%t\n", before, after, found)

	// 分隔符不存在时 found 为 false，before 为原串
	b, a, ok := strings.CutLast("no-separator-here", "=")
	fmt.Printf("not found: before=%q after=%q found=%t\n", b, a, ok)
	fmt.Println()
}

// ============================================================================
// Feature 4: net/url URL.Clone / Values.Clone (Go 1.27)
// ============================================================================

// demonstrateURLClone shows deep-copy helpers for url.URL and url.Values
func demonstrateURLClone() {
	fmt.Println("=== Feature 4: url.URL.Clone / url.Values.Clone ===")

	u, err := url.Parse("https://example.com/api/v1?limit=10&offset=0")
	if err != nil {
		panic(err)
	}

	// Clone 深拷贝：修改副本不影响原 URL
	clone := u.Clone()
	clone.Scheme = "http"
	clone.Path = "/changed"
	clone.RawQuery = "limit=99"
	fmt.Printf("original: %s\n", u)
	fmt.Printf("clone   : %s\n", clone)

	// Values.Clone 同样返回深拷贝副本
	vals := url.Values{"tag": {"go", "127"}}
	copied := vals.Clone()
	copied.Add("tag", "cloned")
	fmt.Printf("original Values: %v\n", vals["tag"])
	fmt.Printf("cloned   Values: %v\n", copied["tag"])
	fmt.Println()
}

// ============================================================================
// Feature 5: encoding/json/v2 (Go 1.27 GA)
// ============================================================================

// Product is a sample value for json/v2 marshaling.
type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

// demonstrateJSONv2 shows encoding/json/v2 marshal/unmarshal and options
func demonstrateJSONv2() {
	fmt.Println("=== Feature 5: encoding/json/v2 ===")
	fmt.Println("Note: Go 1.27 起旧的 encoding/json 已改由 v2 实现，行为兼容。")
	fmt.Println()

	p := Product{Name: "keyboard", Price: 299, Stock: 42}

	// 基本 Marshal / Unmarshal（import 路径为 encoding/json/v2）
	data, err := jsonv2.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Marshal: %s\n", data)

	var back Product
	if err := jsonv2.Unmarshal(data, &back); err != nil {
		panic(err)
	}
	fmt.Printf("Unmarshal round-trip equal: %t\n", back == p)

	// Options：jsonv2 通过 Option 参数配置语义行为。
	// 重复字段名属于语法层行为，由 jsontext 包提供：默认拒绝重复名（返回错误），
	// 可用 jsontext.AllowDuplicateNames(true) 放宽（对应旧 v1 的宽松行为）。
	dup := []byte(`{"name":"a","name":"b"}`)
	var m map[string]string
	if err := jsonv2.Unmarshal(dup, &m); err != nil {
		fmt.Printf("duplicate names rejected by default: %v\n", err)
	}
	if err := jsonv2.Unmarshal(dup, &m, jsontext.AllowDuplicateNames(true)); err != nil {
		panic(err)
	}
	fmt.Printf("with jsontext.AllowDuplicateNames(true): %v\n", m)

	// RejectUnknownMembers：严格模式，拒绝未知字段
	type Strict struct {
		Name string `json:"name"`
	}
	if err := jsonv2.Unmarshal([]byte(`{"name":"x","extra":1}`), &Strict{},
		jsonv2.RejectUnknownMembers(true)); err != nil {
		fmt.Printf("unknown members rejected: %v\n", err)
	}
	fmt.Println()
}

// ============================================================================
// Feature 6: httptest / synctest（详见 synctest_test.go）
// ============================================================================

// demonstrateHTTPServer shows httptest server basics.
// NewTestServer(t, handler) 是 Go 1.27 新增：基于 testing.TB 自动管理生命周期、
// 默认使用内存网络，配合 testing/synctest 可无真实网络测试并发代码。
// 该 API 只能在测试中使用，因此这里用经典的 httptest.NewServer 演示，
// synctest + NewTestServer 的完整用法见 synctest_test.go。
func demonstrateHTTPServer() {
	fmt.Println("=== Feature 6: httptest (see synctest_test.go for NewTestServer) ===")

	var mu sync.Mutex
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests++
		mu.Unlock()
		fmt.Fprintf(w, "hello from %s", r.URL.Path)
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/demo")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	fmt.Printf("GET %s/demo -> %d %q\n", srv.URL, resp.StatusCode, buf.String())
	fmt.Println()
}

// ============================================================================
// Feature 7: container-aware GOMAXPROCS (Go 1.27)
// ============================================================================

// demonstrateGOMAXPROCS prints GOMAXPROCS and explains the Go 1.27 default.
func demonstrateGOMAXPROCS() {
	fmt.Println("=== Feature 7: container-aware GOMAXPROCS ===")
	fmt.Printf("runtime.GOMAXPROCS(0) = %d\n", runtime.GOMAXPROCS(0))
	fmt.Println("Go 1.27 起运行时默认感知容器 CPU 配额（cgroup），")
	fmt.Println("GOMAXPROCS 默认取 min(宿主机核数, 容器配额) 而非宿主机核数。")
	fmt.Println("可用 GODEBUG=containermaxprocs=0 恢复旧行为。")
	fmt.Println()
}

// ============================================================================
// Feature 8: goroutineleak profile (Go 1.27 GA)
// ============================================================================

// demonstrateGoroutineLeak shows the new goroutineleak profile.
func demonstrateGoroutineLeak() {
	fmt.Println("=== Feature 8: goroutineleak profile (GA in Go 1.27) ===")

	// 1.27 起 goroutineleak profile 正式 GA：
	// 挂 net/http/pprof 后可通过 /debug/pprof/goroutineleak 查看“疑似泄漏”的 goroutine。
	// 与 /debug/pprof/goroutine?debug=1 的区别：只保留启动后一直阻塞且未结束的 goroutine。
	prof := pprof.Lookup("goroutineleak")
	fmt.Printf("pprof.Lookup(\"goroutineleak\") available: %t\n", prof != nil)
	fmt.Println("Endpoint: /debug/pprof/goroutineleak (via net/http/pprof)")
	fmt.Println()
}

// ============================================================================
// Main
// ============================================================================

func main() {
	fmt.Println("Go 1.27 Features Demo")
	fmt.Println("====================")

	demonstrateGenericMethods()
	demonstrateUUID()
	demonstrateCutLast()
	demonstrateURLClone()
	demonstrateJSONv2()
	demonstrateHTTPServer()
	demonstrateGOMAXPROCS()
	demonstrateGoroutineLeak()

	// synctest.Sleep / synctest.Wait 只能在 synctest 测试气泡内使用，
	// 完整示例见同目录 synctest_test.go（httptest.NewTestServer + synctest.Wait）。

	fmt.Println("All demonstrations completed!")
}
