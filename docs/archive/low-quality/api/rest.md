# REST API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **Content-Type**: `application/json`

## 用户 API

### 创建用户

```http
POST /api/v1/users
Content-Type: application/json

{
  "email": "user@example.com",
  "name": "User Name"
}
```

**响应**:

```json
{
  "code": 201,
  "message": "success",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "User Name",
    "created_at": "2025-11-11T10:00:00Z",
    "updated_at": "2025-11-11T10:00:00Z"
  }
}
```

### 获取用户

```http
GET /api/v1/users/{id}
```

### 更新用户

```http
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "name": "Updated Name",
  "email": "updated@example.com"
}
```

### 删除用户

```http
DELETE /api/v1/users/{id}
```

### 列出用户

```http
GET /api/v1/users?limit=10&offset=0
```

---

## 错误响应格式

```json
{
  "code": 400,
  "message": "error",
  "error": {
    "code": "INVALID_INPUT",
    "message": "Invalid request body"
  }
}
```

## 状态码

- `200` - 成功
- `201` - 创建成功
- `204` - 删除成功
- `400` - 请求错误
- `404` - 资源不存在
- `409` - 资源冲突
- `500` - 服务器错误
