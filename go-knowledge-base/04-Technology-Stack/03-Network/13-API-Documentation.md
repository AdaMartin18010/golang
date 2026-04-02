# TS-NET-013: API Documentation Best Practices

> **维度**: Technology Stack > Network
> **级别**: S (16+ KB)
> **标签**: #api-documentation #openapi #rest #best-practices
> **权威来源**:
>
> - [API Documentation Best Practices](https://swagger.io/resources/articles/best-practices-in-api-documentation/) - Swagger

---

## 1. API Documentation Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       API Documentation Components                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. Overview Section                                                         │
│     - API purpose and value proposition                                      │
│     - Base URL and environment details                                       │
│     - Authentication requirements                                            │
│     - Rate limiting information                                              │
│                                                                              │
│  2. Getting Started                                                          │
│     - Quick start guide                                                      │
│     - First API call example                                                 │
│     - SDKs and client libraries                                              │
│                                                                              │
│  3. Authentication                                                           │
│     - Authentication methods                                                 │
│     - Token acquisition                                                      │
│     - Security best practices                                                │
│                                                                              │
│  4. API Reference                                                            │
│     - Endpoint descriptions                                                  │
│     - Request/response schemas                                               │
│     - Error codes                                                            │
│     - Code examples in multiple languages                                    │
│                                                                              │
│  5. Guides and Tutorials                                                     │
│     - Common use cases                                                       │
│     - Step-by-step tutorials                                                 │
│     - Best practices                                                         │
│                                                                              │
│  6. Changelog                                                                │
│     - Version history                                                        │
│     - Breaking changes                                                       │
│     - Deprecation notices                                                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. REST API Design Documentation

```markdown
# User API Documentation

## Base URL

```

Production: <https://api.example.com/v1>
Staging: <https://api-staging.example.com/v1>

```

## Authentication

All API requests require an API key passed in the Authorization header:

```

Authorization: Bearer YOUR_API_KEY

```

Obtain your API key from the [developer dashboard](https://dashboard.example.com).

## Rate Limiting

- 1000 requests per hour per API key
- Rate limit headers included in all responses:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Unix timestamp when limit resets

## Endpoints

### List Users

```http
GET /users
```

Returns a list of users.

#### Query Parameters

| Parameter | Type    | Required | Default | Description                |
|-----------|---------|----------|---------|----------------------------|
| limit     | integer | No       | 10      | Number of results per page |
| offset    | integer | No       | 0       | Offset for pagination      |
| sort      | string  | No       | id      | Sort field                 |
| order     | string  | No       | asc     | Sort order (asc/desc)      |

#### Response

```json
{
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "total": 100,
    "limit": 10,
    "offset": 0,
    "has_more": true
  }
}
```

#### Example Request

```bash
curl -X GET "https://api.example.com/v1/users?limit=10" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### Get User

```http
GET /users/{id}
```

Returns a specific user by ID.

#### Path Parameters

| Parameter | Type    | Required | Description    |
|-----------|---------|----------|----------------|
| id        | integer | Yes      | User ID        |

#### Response

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Create User

```http
POST /users
```

Creates a new user.

#### Request Body

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "age": 28
}
```

#### Validation Rules

- `name`: Required, 2-100 characters
- `email`: Required, valid email format
- `age`: Optional, 0-150

#### Response

```json
{
  "id": 2,
  "name": "Jane Doe",
  "email": "jane@example.com",
  "created_at": "2024-01-16T08:00:00Z"
}
```

## Error Handling

### Error Response Format

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

### Error Codes

| Code                  | Status | Description                    |
|----------------------|--------|--------------------------------|
| INVALID_REQUEST      | 400    | Malformed request              |
| VALIDATION_ERROR     | 400    | Input validation failed        |
| UNAUTHORIZED         | 401    | Authentication required        |
| FORBIDDEN            | 403    | Insufficient permissions       |
| NOT_FOUND            | 404    | Resource not found             |
| RATE_LIMIT_EXCEEDED  | 429    | Too many requests              |
| INTERNAL_ERROR       | 500    | Server error                   |

## SDKs and Libraries

- [JavaScript/TypeScript](https://github.com/example/js-sdk)
- [Python](https://github.com/example/python-sdk)
- [Go](https://github.com/example/go-sdk)
- [Java](https://github.com/example/java-sdk)

## Changelog

### v1.1.0 (2024-01-15)

- Added pagination support to List Users endpoint
- Added `sort` and `order` query parameters

### v1.0.0 (2024-01-01)

- Initial release

```

---

## 3. Best Practices

```

API Documentation Best Practices:

1. Keep it up to date
   - Update docs with every API change
   - Version your documentation
   - Use automated tools (Swagger/OpenAPI)

2. Be comprehensive
   - Document all endpoints
   - Include all parameters
   - Provide complete examples
   - Explain error scenarios

3. Make it accessible
   - Clear navigation
   - Search functionality
   - Multiple code examples
   - Interactive try-it feature

4. Use consistent formatting
   - Standard response formats
   - Consistent naming conventions
   - Clear error messages

5. Include practical examples
   - Real-world use cases
   - Complete request/response cycles
   - Common integration patterns

```

---

## 4. Checklist

```

API Documentation Checklist:
□ Overview and purpose clear
□ Base URLs documented
□ Authentication explained
□ All endpoints documented
□ Request/response examples
□ Error codes documented
□ Rate limits specified
□ SDKs and tools listed
□ Changelog maintained
□ Code examples in multiple languages

```
