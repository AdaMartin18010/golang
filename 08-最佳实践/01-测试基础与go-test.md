# 测试基础与go test

## 📚 **理论分析**

- 测试是保障代码质量的核心手段，Go内置强大的测试框架，无需第三方依赖。
- go test 支持单元测试、基准测试、示例测试，集成简单，执行高效。

## 🛠️ **go test基本用法**

- 测试文件以`_test.go`结尾，测试函数以`Test`开头
- 运行所有测试：

  ```bash
  go test
  ```

- 运行指定测试：

  ```bash
  go test -run TestFuncName
  ```

- 显示详细输出：

  ```bash
  go test -v
  ```

## 💻 **代码示例**

### **基本单元测试**

```go
// file: math_test.go
package math
import "testing"
func TestAdd(t *testing.T) {
    got := Add(1, 2)
    want := 3
    if got != want {
        t.Errorf("Add(1,2) = %d; want %d", got, want)
    }
}
```

### **测试失败输出**

- t.Error/t.Errorf：记录失败但继续
- t.Fatal/t.Fatalf：记录失败并终止当前测试

## 🎯 **最佳实践**

- 每个包都应有测试，覆盖核心逻辑
- 用表驱动法组织测试用例，提升可维护性
- 测试用例应独立、可复现
- 测试命名清晰，便于定位

## 🔍 **常见问题**

- Q: go test如何只测试某个函数？
  A: 用`-run`参数指定正则
- Q: 测试文件必须和源码同包吗？
  A: 推荐同包，便于测试私有函数

## 📚 **扩展阅读**

- [Go官方测试文档](https://golang.org/pkg/testing/)
- [Go测试实战](https://geektutu.com/post/hpg-golang-unit-test.html)

---

**文档维护者**: AI Assistant  
**最后更新**: 2024年6月27日  
**文档状态**: 完成
