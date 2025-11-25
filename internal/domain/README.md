# Domain Layer (领域层)

Clean Architecture 的领域层，包含核心业务逻辑。

## 结构

```
domain/
├── user/          # 用户领域
├── order/         # 订单领域
└── product/       # 产品领域
```

## 规则

- ✅ 不依赖任何外部框架
- ✅ 只包含业务逻辑
- ✅ 定义接口，不包含实现
- ❌ 不能导入 infrastructure 或 interfaces 层
