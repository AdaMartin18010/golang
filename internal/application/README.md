# Application Layer (应用层)

Clean Architecture 的应用层，包含用例编排。

## 结构

```
application/
├── user/          # 用户用例
├── order/         # 订单用例
└── product/       # 产品用例
```

## 规则

- ✅ 只能导入 domain 层
- ✅ 协调领域对象
- ✅ 包含用例逻辑
- ❌ 不能导入 infrastructure 或 interfaces 层
