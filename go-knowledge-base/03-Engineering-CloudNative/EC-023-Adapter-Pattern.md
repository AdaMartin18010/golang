# EC-023: Adapter Pattern

## Problem Formalization

### The Interface Mismatch Problem

In distributed systems, components often need to interact with incompatible interfaces—different protocols, data formats, or interaction models. The Adapter pattern bridges these incompatibilities without modifying existing code.

#### Problem Statement

Given:

- Target interface T required by client C
- Adaptee A with incompatible interface
- Conversion functions F = {f₁, f₂, ..., fₙ}

Find adapter D such that:

```
∀ operations o ∈ T:
    D.o() = F(A.equivalent_operations())

Constraints:
    - C remains unchanged (uses T)
    - A remains unchanged (provides its interface)
    - D maintains semantic equivalence
    - Performance overhead is acceptable
```

### Adapter Pattern Variants

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Adapter Pattern Variants                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. Object Adapter (Composition)                                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  ┌──────────┐         ┌──────────────┐      ┌──────────────┐   │   │
│  │  │  Client  │────────►│   Target     │◄─────│   Adapter    │   │   │
│  │  │          │         │   Interface  │      │              │   │   │
│  │  └──────────┘         └──────────────┘      │  ┌──────────┐ │   │   │
│  │                                             │  │  Adaptee │ │   │   │
│  │                                             │  │  (field) │ │   │   │
│  │                                             │  └──────────┘ │   │   │
│  │                                             └──────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  2. Class Adapter (Inheritance) - Not available in Go                   │
│                                                                         │
│  3. Two-Way Adapter                                                     │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  ┌──────────┐         ┌──────────────┐      ┌──────────────┐   │   │
│  │  │ Client A │◄───────►│   Adapter    │◄────►│  Client B    │   │   │
│  │  │ (FormatX)│         │ (Translates) │      │  (FormatY)   │   │   │
│  │  └──────────┘         └──────────────┘      └──────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Solution Architecture

### Protocol Adapter Example

```
┌─────────────────────────────────────────────────────────────────────────┐
│              REST to gRPC Adapter Architecture                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  REST Client                                                            │
│       │                                                                 │
│       │ GET /api/v1/users/123                                           │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  REST Adapter (HTTP Server)                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  HTTP Request Parser                                        │ │   │
│  │  │  • Path parsing: /users/{id}                                │ │   │
│  │  │  • Method mapping: GET → GetUser                            │ │   │
│  │  │  • Header extraction: Authorization                         │ │   │
│  │  │  • Query params: ?include=profile                           │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  │                              │                                    │   │
│  │                              ▼                                    │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  Protocol Translation                                       │ │   │
│  │  │  • HTTP/1.1 → HTTP/2                                        │ │   │
│  │  │  • JSON → Protocol Buffers                                  │ │   │
│  │  │  • REST semantics → gRPC methods                            │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  │                              │                                    │   │
│  └──────────────────────────────┼────────────────────────────────────┘   │
│                                 │ gRPC                                    │
│                                 ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  gRPC Service                                                    │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  service UserService {                                      │ │   │
│  │  │    rpc GetUser(GetUserRequest) returns (User);              │ │   │
│  │  │  }                                                          │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Message Format Adapter

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Avro to JSON Adapter                                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Kafka Topic: events.avro                                               │
│       │                                                                 │
│       │ Binary Avro with Schema ID                                      │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  Avro Adapter                                                    │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  1. Schema Resolution                                         │ │   │
│  │  │     • Extract schema ID from message                        │ │   │
│  │  │     • Fetch from Schema Registry                            │ │   │
│  │  │     • Cache schema locally                                  │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  │                              │                                    │   │
│  │                              ▼                                    │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  2. Deserialization                                         │ │   │
│  │  │     • Binary Avro → GenericRecord                           │ │   │
│  │  │     • Apply schema evolution rules                          │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  │                              │                                    │   │
│  │                              ▼                                    │   │
│  │  ┌─────────────────────────────────────────────────────────────┐ │   │
│  │  │  3. Transformation                                          │ │   │
│  │  │     • Avro types → JSON types                               │ │   │
│  │  │     • Byte arrays → Base64 strings                          │ │   │
│  │  │     • Logical types → ISO formats                           │ │   │
│  │  │       (timestamp-millis → "2024-01-15T10:30:00Z")           │ │   │
│  │  │       (decimal → "123.45")                                  │ │   │
│  │  └─────────────────────────────────────────────────────────────┘ │   │
│  │                              │                                    │   │
│  └──────────────────────────────┼────────────────────────────────────┘   │
│                                 │ JSON                                    │
│                                 ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  {                                                              │   │
│  │    "userId": "123",                                             │   │
│  │    "timestamp": "2024-01-15T10:30:00Z",                         │   │
│  │    "amount": "123.45"                                           │   │
│  │  }                                                              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Production-Ready Go Implementation

### Generic Adapter Framework

```go
// pkg/adapter/adapter.go
package adapter

import (
    "context"
    "fmt"
)

// Target is the interface clients expect
type Target interface {
    Request(ctx context.Context, req Request) (Response, error)
}

// Request represents a generic request
type Request struct {
    Method  string
    Path    string
    Headers map[string]string
    Body    []byte
}

// Response represents a generic response
type Response struct {
    StatusCode int
    Headers    map[string]string
    Body       []byte
}

// Adaptee is the interface we need to adapt
type Adaptee interface {
    SpecificRequest(input interface{}) (interface{}, error)
}

// Adapter implements the Target interface using Adaptee
type Adapter struct {
    adaptee    Adaptee
    translator Translator
    validator  Validator
}

// Translator converts between request/response formats
type Translator interface {
    ToAdapteeFormat(req Request) (interface{}, error)
    FromAdapteeFormat(result interface{}) (Response, error)
}

// Validator validates requests before translation
type Validator interface {
    Validate(req Request) error
}

func NewAdapter(adaptee Adaptee, translator Translator, validator Validator) *Adapter {
    return &Adapter{
        adaptee:    adaptee,
        translator: translator,
        validator:  validator,
    }
}

func (a *Adapter) Request(ctx context.Context, req Request) (Response, error) {
    // Validate
    if err := a.validator.Validate(req); err != nil {
        return Response{StatusCode: 400}, fmt.Errorf("validation failed: %w", err)
    }

    // Translate to adaptee format
    adapteeInput, err := a.translator.ToAdapteeFormat(req)
    if err != nil {
        return Response{StatusCode: 400}, fmt.Errorf("translation failed: %w", err)
    }

    // Call adaptee
    result, err := a.adaptee.SpecificRequest(adapteeInput)
    if err != nil {
        return Response{StatusCode: 500}, fmt.Errorf("adaptee error: %w", err)
    }

    // Translate response
    response, err := a.translator.FromAdapteeFormat(result)
    if err != nil {
        return Response{StatusCode: 500}, fmt.Errorf("response translation failed: %w", err)
    }

    return response, nil
}
```

### REST to gRPC Adapter

```go
// internal/adapter/restgrpc/adapter.go
package restgrpc

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/go-chi/chi/v5"
    "google.golang.org/grpc"
    "google.golang.org/protobuf/encoding/protojson"
    "google.golang.org/protobuf/types/known/emptypb"
)

// RESTToGRPCAdapter adapts REST HTTP requests to gRPC calls
type RESTToGRPCAdapter struct {
    router     *chi.Mux
    connections map[string]*grpc.ClientConn
    converters map[string]MethodConverter
}

// MethodConverter converts REST requests to gRPC
type MethodConverter interface {
    Convert(r *http.Request) (grpcMethod string, request interface{}, err error)
    ConvertResponse(grpcResp interface{}) (status int, body []byte, err error)
}

func NewRESTToGRPCAdapter() *RESTToGRPCAdapter {
    return &RESTToGRPCAdapter{
        router:      chi.NewRouter(),
        connections: make(map[string]*grpc.ClientConn),
        converters:  make(map[string]MethodConverter),
    }
}

func (a *RESTToGRPCAdapter) RegisterService(serviceName string, target string) error {
    conn, err := grpc.Dial(target, grpc.WithInsecure()) // Use TLS in production
    if err != nil {
        return fmt.Errorf("connecting to %s: %w", target, err)
    }

    a.connections[serviceName] = conn
    return nil
}

func (a *RESTToGRPCAdapter) RegisterConverter(path string, method string, converter MethodConverter) {
    key := fmt.Sprintf("%s %s", method, path)
    a.converters[key] = converter

    a.router.Method(method, path, a.handleRequest(serviceName, converter))
}

func (a *RESTToGRPCAdapter) handleRequest(serviceName string, converter MethodConverter) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        conn, ok := a.connections[serviceName]
        if !ok {
            http.Error(w, "Service not found", http.StatusNotFound)
            return
        }

        // Convert REST to gRPC
        grpcMethod, grpcReq, err := converter.Convert(r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Invoke gRPC
        ctx := r.Context()
        var grpcResp interface{}

        // Use grpc.Invoke for generic method calls
        err = conn.Invoke(ctx, grpcMethod, grpcReq, &grpcResp)
        if err != nil {
            status := grpcStatusFromError(err)
            http.Error(w, status.Message(), grpcCodeToHTTP(status.Code()))
            return
        }

        // Convert response
        httpStatus, body, err := converter.ConvertResponse(grpcResp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(httpStatus)
        w.Write(body)
    }
}

// UserServiceConverter example converter
type UserServiceConverter struct {
    client pb.UserServiceClient // Generated protobuf client
}

func (c *UserServiceConverter) Convert(r *http.Request) (string, interface{}, error) {
    path := r.URL.Path
    method := r.Method

    switch {
    case method == "GET" && strings.HasPrefix(path, "/users/"):
        // Extract user ID from path
        userID := chi.URLParam(r, "id")
        return "/pb.UserService/GetUser", &pb.GetUserRequest{Id: userID}, nil

    case method == "POST" && path == "/users":
        var req pb.CreateUserRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            return "", nil, err
        }
        return "/pb.UserService/CreateUser", &req, nil

    case method == "GET" && path == "/users":
        // Parse query params
        page := parseInt(r.URL.Query().Get("page"), 1)
        size := parseInt(r.URL.Query().Get("size"), 20)
        return "/pb.UserService/ListUsers", &pb.ListUsersRequest{
            Page: int32(page),
            Size: int32(size),
        }, nil

    default:
        return "", nil, fmt.Errorf("unsupported endpoint")
    }
}

func (c *UserServiceConverter) ConvertResponse(grpcResp interface{}) (int, []byte, error) {
    // Use protojson for proper JSON encoding
    body, err := protojson.Marshal(grpcResp.(proto.Message))
    if err != nil {
        return 0, nil, err
    }

    return http.StatusOK, body, nil
}

func grpcCodeToHTTP(code codes.Code) int {
    switch code {
    case codes.OK:
        return http.StatusOK
    case codes.InvalidArgument:
        return http.StatusBadRequest
    case codes.NotFound:
        return http.StatusNotFound
    case codes.PermissionDenied:
        return http.StatusForbidden
    case codes.Unauthenticated:
        return http.StatusUnauthorized
    case codes.AlreadyExists:
        return http.StatusConflict
    case codes.Unimplemented:
        return http.StatusNotImplemented
    case codes.Unavailable:
        return http.StatusServiceUnavailable
    case codes.DeadlineExceeded:
        return http.StatusGatewayTimeout
    default:
        return http.StatusInternalServerError
    }
}
```

### Database Protocol Adapter

```go
// internal/adapter/postgresmysql/adapter.go
package postgresmysql

import (
    "context"
    "database/sql"
    "fmt"
    "strings"

    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
)

// PostgresToMySQLAdapter allows PostgreSQL clients to connect to MySQL
type PostgresToMySQLAdapter struct {
    mysqlDB *sql.DB

    // Query translation
    translator *SQLTranslator

    // Type mapping
    typeMapper *TypeMapper
}

// SQLTranslator converts PostgreSQL syntax to MySQL
type SQLTranslator struct {
    rewriteRules []RewriteRule
}

type RewriteRule struct {
    Pattern     string
    Replacement string
    Matcher     func(query string) bool
}

func NewPostgresToMySQLAdapter(mysqlDSN string) (*PostgresToMySQLAdapter, error) {
    db, err := sql.Open("mysql", mysqlDSN)
    if err != nil {
        return nil, err
    }

    return &PostgresToMySQLAdapter{
        mysqlDB:    db,
        translator: NewSQLTranslator(),
        typeMapper: NewTypeMapper(),
    }, nil
}

func (a *PostgresToMySQLAdapter) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    // Translate PostgreSQL query to MySQL
    mysqlQuery, mysqlArgs, err := a.translator.Translate(query, args)
    if err != nil {
        return nil, fmt.Errorf("translation failed: %w", err)
    }

    return a.mysqlDB.QueryContext(ctx, mysqlQuery, mysqlArgs...)
}

func (a *PostgresToMySQLAdapter) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    mysqlQuery, mysqlArgs, err := a.translator.Translate(query, args)
    if err != nil {
        return nil, fmt.Errorf("translation failed: %w", err)
    }

    return a.mysqlDB.ExecContext(ctx, mysqlQuery, mysqlArgs...)
}

func NewSQLTranslator() *SQLTranslator {
    return &SQLTranslator{
        rewriteRules: []RewriteRule{
            // LIMIT/OFFSET syntax
            {
                Pattern: "LIMIT $1 OFFSET $2",
                Replacement: "LIMIT $2, $1",
            },
            // ILIKE to LIKE with lowercase
            {
                Pattern: "ILIKE",
                Replacement: "LIKE",
                Matcher: func(q string) bool {
                    return strings.Contains(q, "ILIKE")
                },
            },
            // RETURNING clause (MySQL uses LAST_INSERT_ID())
            {
                Pattern: "RETURNING *",
                Replacement: "",
            },
            // SERIAL to AUTO_INCREMENT
            {
                Pattern: "SERIAL",
                Replacement: "AUTO_INCREMENT",
            },
            // NOW() is the same
            // CURRENT_TIMESTAMP is the same
            // String concatenation
            {
                Pattern: "||",
                Replacement: "CONCAT()",
            },
        },
    }
}

func (t *SQLTranslator) Translate(pgQuery string, args []interface{}) (string, []interface{}, error) {
    mysqlQuery := pgQuery

    // Apply rewrite rules
    for _, rule := range t.rewriteRules {
        if rule.Matcher == nil || rule.Matcher(mysqlQuery) {
            mysqlQuery = strings.ReplaceAll(mysqlQuery, rule.Pattern, rule.Replacement)
        }
    }

    // Handle placeholder conversion ($1, $2 → ?, ?)
    mysqlQuery, mysqlArgs := convertPlaceholders(mysqlQuery, args)

    return mysqlQuery, mysqlArgs, nil
}

func convertPlaceholders(query string, args []interface{}) (string, []interface{}) {
    // PostgreSQL uses $1, $2; MySQL uses ?
    // Need to reorder args because MySQL uses positional ?

    var result strings.Builder
    argIndex := 0
    mysqlArgs := make([]interface{}, 0, len(args))

    for i := 0; i < len(query); i++ {
        if query[i] == '$' && i+1 < len(query) && isDigit(query[i+1]) {
            // Found placeholder
            j := i + 1
            for j < len(query) && isDigit(query[j]) {
                j++
            }

            // Extract index (1-based in PostgreSQL)
            idx := parseInt(query[i+1:j]) - 1
            if idx >= 0 && idx < len(args) {
                mysqlArgs = append(mysqlArgs, args[idx])
                argIndex++
            }

            result.WriteByte('?')
            i = j - 1
        } else {
            result.WriteByte(query[i])
        }
    }

    return result.String(), mysqlArgs
}

// TypeMapper handles PostgreSQL to MySQL type conversions
type TypeMapper struct {
    pgToMySQL map[string]string
}

func NewTypeMapper() *TypeMapper {
    return &TypeMapper{
        pgToMySQL: map[string]string{
            "serial":           "INT AUTO_INCREMENT",
            "bigserial":        "BIGINT AUTO_INCREMENT",
            "varchar":          "VARCHAR",
            "text":             "TEXT",
            "integer":          "INT",
            "bigint":           "BIGINT",
            "boolean":          "TINYINT(1)",
            "timestamp":        "DATETIME",
            "timestamptz":      "TIMESTAMP",
            "json":             "JSON",
            "jsonb":            "JSON",
            "uuid":             "CHAR(36)",
            "bytea":            "BLOB",
            "real":             "FLOAT",
            "double precision": "DOUBLE",
            "numeric":          "DECIMAL",
        },
    }
}
```

## Trade-off Analysis

### Adapter vs Alternative Approaches

| Approach | Flexibility | Performance | Maintenance | When to Use |
|----------|-------------|-------------|-------------|-------------|
| **Adapter Pattern** | High | Medium | Medium | Protocol/format bridging |
| **Direct Modification** | N/A | High | Low | When you control both sides |
| **Proxy with Translation** | Medium | High | Low | Simple transformations |
| **Message Queue** | High | Lower | High | Async, eventual consistency |

### Performance Considerations

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Adapter Overhead Analysis                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Operation Type          Latency Impact    CPU Impact    Memory         │
│  ─────────────────────────────────────────────────────────────────     │
│  Protocol Translation    1-5ms            Medium        Low             │
│  (HTTP↔gRPC)                                                            │
│                                                                         │
│  Serialization Change    0.5-2ms          Medium        Medium          │
│  (JSON↔Protobuf)                                                        │
│                                                                         │
│  Character Encoding      <0.1ms           Low           Low             │
│  (UTF-8↔ASCII)                                                          │
│                                                                         │
│  Data Format Parsing     1-10ms           High          High            │
│  (XML↔JSON)                                                             │
│                                                                         │
│  SQL Dialect Translation 0.5-2ms          Low           Low             │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Testing Strategies

### Adapter Testing

```go
// test/adapter/adapter_test.go
package adapter

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock implementations
type MockAdaptee struct {
    mock.Mock
}

func (m *MockAdaptee) SpecificRequest(input interface{}) (interface{}, error) {
    args := m.Called(input)
    return args.Get(0), args.Error(1)
}

func TestAdapterRequestTranslation(t *testing.T) {
    // Setup
    adaptee := new(MockAdaptee)
    translator := &testTranslator{}
    validator := &testValidator{}

    adapter := NewAdapter(adaptee, translator, validator)

    // Expectations
    adaptee.On("SpecificRequest", mock.Anything).Return(
        map[string]interface{}{"id": "123", "name": "Test"},
        nil,
    )

    // Execute
    req := Request{
        Method: "GET",
        Path:   "/users/123",
    }

    resp, err := adapter.Request(context.Background(), req)

    // Verify
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    adaptee.AssertExpectations(t)
}

func TestAdapterValidationFailure(t *testing.T) {
    validator := &failingValidator{error: fmt.Errorf("invalid request")}
    adapter := NewAdapter(nil, nil, validator)

    resp, err := adapter.Request(context.Background(), Request{})

    assert.Error(t, err)
    assert.Equal(t, 400, resp.StatusCode)
}

// Round-trip test for bidirectional adapters
func TestBidirectionalRoundTrip(t *testing.T) {
    // Test that converting A→B→A preserves data
    original := &pb.User{
        Id:    "123",
        Name:  "John",
        Email: "john@example.com",
    }

    // A → B
    jsonBytes, err := protojson.Marshal(original)
    require.NoError(t, err)

    // B → A
    restored := &pb.User{}
    err = protojson.Unmarshal(jsonBytes, restored)
    require.NoError(t, err)

    // Verify
    assert.True(t, proto.Equal(original, restored))
}

// Fuzz testing for robustness
func TestAdapterFuzz(t *testing.T) {
    adapter := setupTestAdapter()

    // Generate random inputs
    for i := 0; i < 1000; i++ {
        req := generateRandomRequest()

        // Should not panic
        _, _ = adapter.Request(context.Background(), req)
    }
}
```

## Summary

The Adapter Pattern provides:

1. **Interface Compatibility**: Bridge incompatible interfaces without changing existing code
2. **Protocol Translation**: Convert between different communication protocols
3. **Format Transformation**: Transform data between different formats
4. **Legacy Integration**: Connect modern systems with legacy components
5. **Reusability**: Reuse existing components in new contexts

Key implementation considerations:

- Performance overhead of translation
- Error handling across interface boundaries
- Stateful vs stateless adapters
- Caching frequently translated data
- Monitoring adapter health and performance
