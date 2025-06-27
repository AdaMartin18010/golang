# Mock与依赖注入

## 📚 **理论分析**

- Mock用于隔离外部依赖，便于单元测试。
- 依赖注入（DI）提升代码可测试性与解耦性。
- Go常用接口+手动注入，或用第三方Mock库（如gomock、testify/mock）。

## 🛠️ **主流Mock方案**

- 手动实现接口Mock
- 使用gomock自动生成Mock
- 使用testify/mock简化Mock

## 💻 **代码示例**

### **手动Mock接口**

```go
type DB interface {
    Get(key string) string
}
type mockDB struct{}
func (m *mockDB) Get(key string) string { return "mock" }
func TestQuery(t *testing.T) {
    db := &mockDB{}
    got := db.Get("id")
    if got != "mock" {
        t.Errorf("want mock, got %s", got)
    }
}
```

### **gomock用法**

```bash
go install github.com/golang/mock/mockgen@latest
mockgen -source=db.go -destination=mock_db.go -package=yourpkg
```

## 🎯 **最佳实践**

- 依赖均用接口抽象，便于Mock
- Mock只用于单元测试，集成测试用真实依赖
- Mock行为应可配置，覆盖边界与异常

## 🔍 **常见问题**

- Q: Mock和Stub区别？
  A: Mock可校验调用行为，Stub只返回固定值
- Q: 依赖注入框架推荐？
  A: Go多用手动注入，少用复杂框架

## 📚 **扩展阅读**

- [Go Mock实战](https://geektutu.com/post/hpg-golang-mock.html)
- [gomock官方文档](https://github.com/golang/mock)
- [testify/mock文档](https://pkg.go.dev/github.com/stretchr/testify/mock)

---

**文档维护者**: AI Assistant  
**最后更新**: 2024年6月27日  
**文档状态**: 完成
