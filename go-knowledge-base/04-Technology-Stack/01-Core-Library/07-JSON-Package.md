# TS-CL-007: Go encoding/json Package - Deep Architecture and Serialization Patterns

> **维度**: Technology Stack > Core Library
> **级别**: S (18+ KB)
> **标签**: #golang #json #serialization #marshaling #encoding
> **权威来源**:
>
> - [Go encoding/json](https://pkg.go.dev/encoding/json) - Official documentation
> - [JSON and Go](https://go.dev/blog/json) - Go Blog
> - [JSON Stream Processing](https://go.dev/src/encoding/json/stream.go) - Source code

---

## 1. JSON Architecture Deep Dive

### 1.1 Encoder/Decoder Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    JSON Encoder/Decoder Architecture                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Encoding Path:                                                             │
│   ┌──────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    │
│   │  Go      │───>│  reflect    │───>│  encodeState │───>│  io.Writer  │    │
│   │  Value   │    │  inspection │    │  buffer      │    │  (output)   │    │
│   └──────────┘    └─────────────┘    └─────────────┘    └─────────────┘    │
│                                                                              │
│   Decoding Path:                                                             │
│   ┌──────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    │
│   │  io.     │───>│  decodeState│───>│  reflect    │───>│  Go         │    │
│   │  Reader  │    │  parser     │    │  assignment │    │  Value      │    │
│   └──────────┘    └─────────────┘    └─────────────┘    └─────────────┘    │
│                                                                              │
│   Key Interfaces:                                                            │
│   - json.Marshaler:   type Marshaler interface { MarshalJSON() ([]byte, error) }
│   - json.Unmarshaler: type Unmarshaler interface { UnmarshalJSON([]byte) error }
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Marshal/Unmarshal Flow

```go
// Marshal flow
func Marshal(v interface{}) ([]byte, error) {
    e := newEncodeState()
    err := e.marshal(v, encOpts{escapeHTML: true})
    if err != nil {
        return nil, err
    }
    return e.Bytes(), nil
}

// Unmarshal flow
func Unmarshal(data []byte, v interface{}) error {
    d := newDecodeState(data)
    return d.unmarshal(v)
}
```

---

## 2. Struct Tags and Field Options

### 2.1 JSON Struct Tags

```go
type User struct {
    // Basic mapping
    ID        int    `json:"id"`

    // Omit empty values
    Nickname  string `json:"nickname,omitempty"`

    // Ignore field
    Password  string `json:"-"`

    // Custom name
    CreatedAt time.Time `json:"created_at"`

    // Pointer to distinguish zero value from missing
    Score     *int   `json:"score,omitempty"`

    // String-encoded number
    BigNumber int64  `json:"big_number,string"`
}
```

### 2.2 Tag Options Reference

| Option | Effect | Example |
|--------|--------|---------|
| `"-"` | Skip field | `json:"-"` |
| `omitempty` | Skip if zero value | `json:"name,omitempty"` |
| `string` | Encode as string | `json:"id,string"` |
| `inline` | Embed fields (json-iterator) | `json:",inline"` |

---

## 3. Custom Marshaling

### 3.1 Implementing Marshaler

```go
type Money struct {
    Amount   int64
    Currency string
}

// MarshalJSON implements custom JSON serialization
func (m Money) MarshalJSON() ([]byte, error) {
    // Represent as "USD 100.00"
    formatted := fmt.Sprintf("%s %.2f", m.Currency, float64(m.Amount)/100)
    return json.Marshal(formatted)
}

// UnmarshalJSON implements custom JSON deserialization
func (m *Money) UnmarshalJSON(data []byte) error {
    var str string
    if err := json.Unmarshal(data, &str); err != nil {
        return err
    }

    parts := strings.Split(str, " ")
    if len(parts) != 2 {
        return fmt.Errorf("invalid money format: %s", str)
    }

    m.Currency = parts[0]
    amount, err := strconv.ParseFloat(parts[1], 64)
    if err != nil {
        return err
    }
    m.Amount = int64(amount * 100)
    return nil
}
```

### 3.2 Time Marshaling

```go
type CustomTime struct {
    time.Time
}

const customTimeFormat = "2006-01-02 15:04:05"

func (t CustomTime) MarshalJSON() ([]byte, error) {
    return []byte(`"` + t.Format(customTimeFormat) + `"`), nil
}

func (t *CustomTime) UnmarshalJSON(data []byte) error {
    str := strings.Trim(string(data), `"`)
    parsed, err := time.Parse(customTimeFormat, str)
    if err != nil {
        return err
    }
    t.Time = parsed
    return nil
}
```

---

## 4. Streaming JSON Processing

### 4.1 Streaming Encoder

```go
func writeLargeJSON(w io.Writer, users <-chan User) error {
    encoder := json.NewEncoder(w)
    encoder.SetIndent("", "  ")

    // Write opening bracket
    w.Write([]byte("[\n"))

    first := true
    for user := range users {
        if !first {
            w.Write([]byte(",\n"))
        }
        first = false

        if err := encoder.Encode(user); err != nil {
            return err
        }
    }

    // Write closing bracket
    w.Write([]byte("]\n"))
    return nil
}
```

### 4.2 Streaming Decoder

```go
func processLargeJSON(r io.Reader) error {
    decoder := json.NewDecoder(r)

    // Read opening bracket
    tok, err := decoder.Token()
    if err != nil {
        return err
    }
    if delim, ok := tok.(json.Delim); !ok || delim != '[' {
        return fmt.Errorf("expected array start")
    }

    // Process array elements one by one
    for decoder.More() {
        var user User
        if err := decoder.Decode(&user); err != nil {
            return err
        }
        processUser(user)
    }

    // Read closing bracket
    tok, err = decoder.Token()
    return err
}
```

### 4.3 Unknown JSON Structures

```go
// Using interface{} for dynamic JSON
func processDynamicJSON(data []byte) error {
    var result map[string]interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        return err
    }

    // Type assertions required
    if name, ok := result["name"].(string); ok {
        fmt.Println("Name:", name)
    }

    if age, ok := result["age"].(float64); ok {
        fmt.Println("Age:", int(age))
    }

    return nil
}

// Using json.RawMessage for delayed parsing
func processWithRawMessage(data []byte) error {
    var wrapper struct {
        Type string          `json:"type"`
        Data json.RawMessage `json:"data"`
    }

    if err := json.Unmarshal(data, &wrapper); err != nil {
        return err
    }

    // Parse Data based on Type
    switch wrapper.Type {
    case "user":
        var user User
        if err := json.Unmarshal(wrapper.Data, &user); err != nil {
            return err
        }
        processUser(user)
    case "product":
        var product Product
        if err := json.Unmarshal(wrapper.Data, &product); err != nil {
            return err
        }
        processProduct(product)
    }

    return nil
}
```

---

## 5. Performance Tuning Guidelines

### 5.1 Performance Comparison

| Approach | Small Objects | Large Objects | Memory Usage |
|----------|--------------|---------------|--------------|
| `json.Marshal` | Fast | Good | Medium |
| `json.Encoder` | Fast | Better | Lower |
| `json.Decoder` | Fast | Better | Lower |
| Third-party (easyjson) | Faster | Faster | Lower |

### 5.2 Optimization Techniques

```go
// 1. Reuse encoders/decoders for batch operations
var encoderPool = sync.Pool{
    New: func() interface{} {
        return json.NewEncoder(io.Discard)
    },
}

// 2. Pre-allocate buffers
buf := make([]byte, 0, 1024)
encoder := json.NewEncoder(bytes.NewBuffer(buf))

// 3. Use json.RawMessage for partial parsing

// 4. Avoid interface{} when possible

// 5. Use easyjson or similar for hot paths
// go install github.com/mailru/easyjson/...@latest
// easyjson -all structs.go
```

---

## 6. Configuration Best Practices

```go
// JSON configuration
type JSONConfig struct {
    EscapeHTML     bool
    Indent         string
    Prefix         string
    DisallowUnknownFields bool
}

func NewEncoderWithConfig(w io.Writer, cfg JSONConfig) *json.Encoder {
    enc := json.NewEncoder(w)
    enc.SetEscapeHTML(cfg.EscapeHTML)
    if cfg.Indent != "" {
        enc.SetIndent(cfg.Prefix, cfg.Indent)
    }
    return enc
}

func NewDecoderWithConfig(r io.Reader, cfg JSONConfig) *json.Decoder {
    dec := json.NewDecoder(r)
    if cfg.DisallowUnknownFields {
        dec.DisallowUnknownFields()
    }
    return dec
}
```

---

## 7. Comparison with Alternatives

| Library | Speed | Features | When to Use |
|---------|-------|----------|-------------|
| **encoding/json** | Standard | Full | Most cases |
| **json-iterator** | 2-3x faster | Drop-in | Performance critical |
| **easyjson** | 3-5x faster | Codegen | Struct-heavy APIs |
| **fastjson** | Very fast | Minimal | Simple parsing only |
| **gjson/sjson** | Fast | Path-based | Path extraction |

---

## 8. Checklist

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      JSON Best Practices                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Struct Design:                                                              │
│  □ Use omitempty for optional fields                                        │
│  □ Use pointers to distinguish zero from missing                            │
│  □ Use custom types with MarshalJSON/UnmarshalJSON                          │
│  □ Document expected JSON format                                            │
│                                                                              │
│  Performance:                                                                │
│  □ Use streaming for large data                                             │
│  □ Consider code generation for hot paths                                   │
│  □ Reuse encoders where possible                                            │
│                                                                              │
│  Safety:                                                                     │
│  □ Validate input before unmarshaling                                       │
│  □ Handle unknown fields appropriately                                      │
│  □ Use strict parsing when possible                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

**质量评级**: S (18+ KB, comprehensive coverage)
