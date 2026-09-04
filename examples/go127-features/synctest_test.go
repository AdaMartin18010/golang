package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"testing/synctest"
	"uuid"
)

// TestSynctestWithNewTestServer 演示 Go 1.27 的 httptest.NewTestServer
// 与 testing/synctest 配合：在测试气泡内使用内存网络，
// synctest.Wait 会阻塞直到气泡内所有 goroutine 都处于持久阻塞状态，
// 从而无需 sleep 等待后台 goroutine 完成。
func TestSynctestWithNewTestServer(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Go 1.27 新增：NewTestServer 基于 testing.TB 自动管理生命周期，
		// 默认使用内存网络实现，无需监听真实端口。
		srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "pong")
		}))

		var body string
		go func() {
			resp, err := srv.Client().Get(srv.URL + "/ping")
			if err != nil {
				t.Errorf("GET: %v", err)
				return
			}
			defer resp.Body.Close()
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("ReadAll: %v", err)
				return
			}
			body = string(b)
		}()

		// 等待气泡内所有 goroutine 持久阻塞（即后台请求已完成），无需 time.Sleep。
		synctest.Wait()

		if body != "pong\n" {
			t.Errorf("body = %q, want %q", body, "pong\n")
		}
	})
}

// TestSynctestSleep 演示 synctest.Sleep：气泡内时间被虚拟化，
// Sleep 只会在所有 goroutine 阻塞后才推进，真实时间几乎不流逝。
func TestSynctestSleep(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		done := false
		go func() {
			// 气泡内的 time.Sleep 是虚拟的
			done = true
		}()
		synctest.Wait()
		if !done {
			t.Error("goroutine did not run after synctest.Wait")
		}
	})
}

// TestGenericMethod 断言自定义类型的泛型方法（Go 1.27 语言特性）。
func TestGenericMethod(t *testing.T) {
	s := Stack[int]{}
	s.Push(10)
	s.Push(20)

	got := s.Drain(func(v int) string { return fmt.Sprint(v * 2) })
	want := []string{"20", "40"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("Drain = %v, want %v", got, want)
	}

	// 标准库泛型方法 math/rand/v2.Rand.N[Int]
	r := rand.New(rand.NewPCG(1, 2))
	if n := r.N(int8(10)); n < 0 || n >= 10 {
		t.Errorf("r.N(int8(10)) = %d, want in [0,10)", n)
	}
}

// TestJSONv2Options 断言 json/v2 的 Options 行为。
func TestJSONv2Options(t *testing.T) {
	// v2 默认拒绝重复字段名
	var m map[string]string
	if err := json.Unmarshal([]byte(`{"a":"1","a":"2"}`), &m); err == nil {
		t.Error("expected duplicate name error by default")
	}

	// jsontext.AllowDuplicateNames(true) 恢复宽松行为
	if err := json.Unmarshal([]byte(`{"a":"1","a":"2"}`), &m, jsontext.AllowDuplicateNames(true)); err != nil {
		t.Fatalf("AllowDuplicateNames: %v", err)
	}
	if m["a"] != "2" {
		t.Errorf(`m["a"] = %q, want "2" (last wins)`, m["a"])
	}

	// RejectUnknownMembers 严格校验
	type strict struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(`{"name":"x","oops":1}`), &strict{},
		json.RejectUnknownMembers(true)); err == nil {
		t.Error("expected unknown member error with RejectUnknownMembers")
	}
}

// TestUUID 断言 uuid 包的基本行为。
func TestUUID(t *testing.T) {
	v7 := uuid.NewV7()
	if v7[6]>>4 != 7 {
		t.Errorf("NewV7 version = %d, want 7", v7[6]>>4)
	}
	back, err := uuid.Parse(v7.String())
	if err != nil || back != v7 {
		t.Errorf("Parse round-trip: %v, %v", back, err)
	}
	if _, err := uuid.Parse("bogus"); err == nil {
		t.Error("expected error parsing bogus uuid")
	}
}

// TestCutLast 断言 strings.CutLast 行为。
func TestCutLast(t *testing.T) {
	before, after, found := strings.CutLast("a/b/c", "/")
	if !found || before != "a/b" || after != "c" {
		t.Errorf("CutLast = %q, %q, %t", before, after, found)
	}
	if _, _, found := strings.CutLast("abc", "/"); found {
		t.Error("CutLast should report found=false without separator")
	}
}

// TestURLClone 断言 Clone 返回深拷贝。
func TestURLClone(t *testing.T) {
	u, err := url.Parse("https://example.com/x?a=1")
	if err != nil {
		t.Fatal(err)
	}
	c := u.Clone()
	c.Query().Set("a", "2") // 修改副本的 query map 不应影响原值（Clone 深拷贝）
	if u.Query().Get("a") != "1" {
		t.Error("URL.Clone should deep-copy the query map")
	}
}
