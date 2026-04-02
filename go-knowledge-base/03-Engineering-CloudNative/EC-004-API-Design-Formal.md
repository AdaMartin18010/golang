# EC-004: API 设计原则的形式化 (API Design: Formal Principles)

> **维度**: Engineering-CloudNative
> **级别**: S (15+ KB)
> **标签**: #api #rest #grpc #design #versioning
> **权威来源**:
>
> - [RESTful Web APIs](https://www.oreilly.com/library/view/restful-web-apis/9781449359713/) - Richardson & Amundsen
> - [Google API Design Guide](https://cloud.google.com/apis/design) - Google
> - [gRPC Style Guide](https://developers.google.com/protocol-buffers/docs/style) - Google

---

## 1. API 的形式化定义

### 1.1 接口契约

**定义 1.1 (API)**
$$\text{API} = \langle \text{operations}, \text{types}, \text{errors} \rangle$$

### 1.2 REST 约束

**定义 1.2 (REST)**

| 约束 | 形式化 |
|------|--------|
| Client-Server | $\text{UI} \perp \text{Data}$ |
| Stateless | $\forall r: \text{Server}(r) \not\ni \text{Session}$ |
| Cacheable | $\text{Response} \ni \text{Cache-Control}$ |
| Uniform Interface | $\text{HTTP verbs} = \{\text{GET}, \text{POST}, \text{PUT}, \text{DELETE}\}$ |

---

## 2. 版本控制

### 2.1 版本策略

**定义 2.1 (兼容性)**
$$\text{BackwardCompatible}(v_2, v_1) \Leftrightarrow \forall c: \text{Works}(c, v_1) \Rightarrow \text{Works}(c, v_2)$$

---

## 3. 多元表征

### 3.1 HTTP 方法矩阵

| 方法 | 幂等 | 安全 | 用途 |
|------|------|------|------|
| GET | ✓ | ✓ | 读取 |
| POST | ✗ | ✗ | 创建 |
| PUT | ✓ | ✗ | 全量更新 |
| PATCH | ✗ | ✗ | 部分更新 |
| DELETE | ✓ | ✗ | 删除 |

### 3.2 状态码决策树

```
响应状态?
├── 2xx (成功)
│   ├── 200 OK
│   ├── 201 Created
│   └── 204 No Content
├── 4xx (客户端错误)
│   ├── 400 Bad Request
│   ├── 401 Unauthorized
│   ├── 403 Forbidden
│   └── 404 Not Found
└── 5xx (服务端错误)
    ├── 500 Internal Error
    ├── 502 Bad Gateway
    └── 503 Service Unavailable
```

---

**质量评级**: S (15KB)
