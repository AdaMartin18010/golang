# Go Fuzzing (模糊测试)

> **分类**: 开源技术堆栈  
> **标签**: #fuzzing #testing #security

---

## 基础 Fuzzing

```go
// fuzz_test.go
package mypkg

import "testing"

func FuzzParse(f *testing.F) {
    // 种子语料库
    f.Add("hello world")
    f.Add("12345")
    f.Add("")
    
    f.Fuzz(func(t *testing.T, input string) {
        result, err := Parse(input)
        if err != nil {
            // 解析错误是允许的
            return
        }
        
        // 验证不变式
        if result.Original != input {
            t.Errorf("Original mismatch: got %q, want %q", result.Original, input)
        }
    })
}
```

---

## 运行 Fuzzing

```bash
# 运行 fuzzing
go test -fuzz=FuzzParse

# 指定时间
go test -fuzz=FuzzParse -fuzztime=30s

# 指定并行度
go test -fuzz=FuzzParse -parallel=4

# 使用 corpus
go test -fuzz=FuzzParse -fuzzcachedir=./corpus
```

---

## 结构化 Fuzzing

```go
func FuzzProcessUser(f *testing.F) {
    f.Add("john", 25, "john@example.com")
    f.Add("jane", 30, "jane@example.com")
    
    f.Fuzz(func(t *testing.T, name string, age int, email string) {
        user := User{
            Name:  name,
            Age:   age,
            Email: email,
        }
        
        // 验证验证逻辑
        err := Validate(user)
        
        // 如果年龄为负，应该返回错误
        if age < 0 && err == nil {
            t.Error("expected error for negative age")
        }
    })
}
```

---

## 自定义语料库

```go
func FuzzJSONParser(f *testing.F) {
    // 从文件加载语料库
    corpus, _ := os.ReadDir("testdata/corpus")
    for _, entry := range corpus {
        data, _ := os.ReadFile(filepath.Join("testdata/corpus", entry.Name()))
        f.Add(data)
    }
    
    f.Fuzz(func(t *testing.T, data []byte) {
        var v interface{}
        if err := json.Unmarshal(data, &v); err != nil {
            return
        }
        
        // 重新序列化应该成功
        _, err := json.Marshal(v)
        if err != nil {
            t.Errorf("remarshal failed: %v", err)
        }
    })
}
```

---

## 发现崩溃

```bash
# 发现崩溃时会保存到 testdata/fuzz/FuzzXXX/

# 重现崩溃
go test -run=FuzzParse/testdata/fuzz/FuzzParse/crash-xxx

# 最小化崩溃
go test -fuzz=FuzzParse -minimize=30s
```

---

## 与 CI 集成

```yaml
# .github/workflows/fuzz.yml
name: Fuzzing
on: [push, pull_request]

jobs:
  fuzz:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.22'
    
    - name: Run Fuzzing
      run: go test -fuzz=. -fuzztime=5m
      continue-on-error: true
    
    - name: Upload Crashers
      uses: actions/upload-artifact@v3
      if: failure()
      with:
        name: fuzz-crashers
        path: testdata/fuzz/
```
