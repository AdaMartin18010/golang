# TS-CL-007: Go encoding/json Package

> **维度**: Technology Stack > Core Library
> **级别**: S (16+ KB)
> **标签**: #golang #json #serialization #encoding #marshaling
> **权威来源**:
>
> - [encoding/json Package](https://golang.org/pkg/encoding/json/) - Go standard library

---

## 1. JSON Basics

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

// Struct with JSON tags
type Person struct {
    ID        int      `json:"id"`
    Name      string   `json:"name"`
    Email     string   `json:"email,omitempty"`  // Omit if empty
    Age       int      `json:"age"`
    IsActive  bool     `json:"is_active"`
    Tags      []string `json:"tags,omitempty"`
    Address   Address  `json:"address"`
    CreatedAt string   `json:"-"`  // Ignore this field
    Password  string   `json:"-"`  // Never serialize password
}

type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
}

// Marshal (struct to JSON)
func marshalExample() {
    p := Person{
        ID:       1,
        Name:     "John Doe",
        Email:    "john@example.com",
        Age:      30,
        IsActive: true,
        Tags:     []string{"developer", "golang"},
        Address: Address{
            Street:  "123 Main St",
            City:    "San Francisco",
            Country: "USA",
        },
    }

    // Marshal to JSON
    data, err := json.Marshal(p)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(data))
    // {"id":1,"name":"John Doe","email":"john@example.com","age":30,"is_active":true,"tags":["developer","golang"],"address":{"street":"123 Main St","city":"San Francisco","country":"USA"}}

    // Marshal with indentation
    prettyJSON, _ := json.MarshalIndent(p, "", "  ")
    fmt.Println(string(prettyJSON))
}

// Unmarshal (JSON to struct)
func unmarshalExample() {
    jsonStr := `{
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "age": 30,
        "is_active": true,
        "tags": ["developer", "golang"],
        "address": {
            "street": "123 Main St",
            "city": "San Francisco",
            "country": "USA"
        }
    }`

    var p Person
    if err := json.Unmarshal([]byte(jsonStr), &p); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", p)
}

// Raw message
func rawMessageExample() {
    type Message struct {
        Type string          `json:"type"`
        Data json.RawMessage `json:"data"` // Delay parsing
    }

    jsonStr := `{"type": "user", "data": {"name": "John", "age": 30}}`

    var msg Message
    json.Unmarshal([]byte(jsonStr), &msg)

    // Parse data based on type
    switch msg.Type {
    case "user":
        var user Person
        json.Unmarshal(msg.Data, &user)
        fmt.Printf("User: %+v\n", user)
    }
}

// Dynamic JSON with map
func dynamicJSONExample() {
    jsonStr := `{"name": "John", "age": 30, "unknown_field": "value"}`

    var result map[string]interface{}
    if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
        log.Fatal(err)
    }

    // Access fields
    name := result["name"].(string)
    age := result["age"].(float64) // JSON numbers are float64

    fmt.Printf("Name: %s, Age: %f\n", name, age)
}

// Custom MarshalJSON/UnmarshalJSON
type CustomTime struct {
    time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
    return json.Marshal(ct.Format("2006-01-02"))
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }

    parsed, err := time.Parse("2006-01-02", s)
    if err != nil {
        return err
    }

    ct.Time = parsed
    return nil
}

// Streaming JSON
type LargeDataset struct {
    Items []Item `json:"items"`
}

type Item struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func streamEncoding() {
    items := []Item{
        {ID: 1, Name: "Item 1"},
        {ID: 2, Name: "Item 2"},
        // ... many items
    }

    // Stream to file
    file, _ := os.Create("output.json")
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")

    for _, item := range items {
        if err := encoder.Encode(item); err != nil {
            log.Fatal(err)
        }
    }
}

func streamDecoding() {
    file, _ := os.Open("input.json")
    defer file.Close()

    decoder := json.NewDecoder(file)

    for decoder.More() {
        var item Item
        if err := decoder.Decode(&item); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Item: %+v\n", item)
    }
}
```

---

## 2. JSON Best Practices

```
JSON Best Practices:
□ Use struct tags for field names
□ Use omitempty for optional fields
□ Use - to exclude sensitive fields
□ Handle errors properly
□ Use json.RawMessage for delayed parsing
□ Use streaming for large datasets
□ Define custom types for special formats
□ Validate JSON before parsing
```

---

## 3. Checklist

```
JSON Checklist:
□ Struct tags defined
□ Sensitive fields excluded
□ Optional fields use omitempty
□ Error handling implemented
□ Types match JSON schema
□ Custom marshaling if needed
□ Streaming for large data
```
