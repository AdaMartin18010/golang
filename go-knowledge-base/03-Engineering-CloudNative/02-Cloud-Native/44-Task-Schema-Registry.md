# 任务 Schema 注册中心 (Task Schema Registry)

> **分类**: 工程与云原生  
> **标签**: #schema-registry #validation #compatibility

---

## Schema 定义

```go
// 任务类型 Schema
type TaskSchema struct {
    ID          string          `json:"id"`
    Type        string          `json:"type"`
    Version     int             `json:"version"`
    Schema      json.RawMessage `json:"schema"`
    Description string          `json:"description"`
    CreatedAt   time.Time       `json:"created_at"`
    CreatedBy   string          `json:"created_by"`
    Status      string          `json:"status"` // active, deprecated
}

// JSON Schema 验证器
type SchemaValidator struct {
    registry SchemaRegistry
    compiler *jsonschema.Compiler
}

func (sv *SchemaValidator) Validate(taskType string, version int, payload []byte) error {
    schema, err := sv.registry.GetSchema(taskType, version)
    if err != nil {
        return fmt.Errorf("schema not found: %w", err)
    }
    
    // 编译 schema
    s, err := sv.compiler.Compile(string(schema.Schema))
    if err != nil {
        return fmt.Errorf("invalid schema: %w", err)
    }
    
    // 验证 payload
    var v interface{}
    if err := json.Unmarshal(payload, &v); err != nil {
        return fmt.Errorf("invalid payload json: %w", err)
    }
    
    if err := s.Validate(v); err != nil {
        return &ValidationError{Errors: formatValidationErrors(err)}
    }
    
    return nil
}
```

---

## Schema 版本管理

```go
type SchemaRegistry struct {
    store SchemaStore
}

func (sr *SchemaRegistry) RegisterSchema(ctx context.Context, schema TaskSchema) error {
    // 检查兼容性
    if schema.Version > 1 {
        prevSchema, _ := sr.store.GetSchema(ctx, schema.Type, schema.Version-1)
        if prevSchema != nil {
            if err := sr.checkCompatibility(prevSchema, &schema); err != nil {
                return fmt.Errorf("schema incompatible: %w", err)
            }
        }
    }
    
    // 保存 schema
    return sr.store.SaveSchema(ctx, schema)
}

func (sr *SchemaRegistry) checkCompatibility(old, new *TaskSchema) error {
    oldFields := extractFields(old.Schema)
    newFields := extractFields(new.Schema)
    
    // 检查破坏性变更
    for name, field := range oldFields {
        newField, exists := newFields[name]
        if !exists && field.Required {
            return fmt.Errorf("required field %s removed", name)
        }
        
        if exists && field.Type != newField.Type {
            return fmt.Errorf("field %s type changed from %s to %s", 
                name, field.Type, newField.Type)
        }
    }
    
    return nil
}

func (sr *SchemaRegistry) GetLatestSchema(ctx context.Context, taskType string) (*TaskSchema, error) {
    schemas, _ := sr.store.ListSchemas(ctx, taskType)
    if len(schemas) == 0 {
        return nil, fmt.Errorf("no schema found for type: %s", taskType)
    }
    
    // 返回最新版本
    latest := schemas[0]
    for _, s := range schemas {
        if s.Version > latest.Version {
            latest = s
        }
    }
    
    return &latest, nil
}
```

---

## 兼容性检查

```go
type CompatibilityChecker struct {
    level CompatibilityLevel
}

type CompatibilityLevel int

const (
    Backward CompatibilityLevel = iota  // 消费者可旧代码读新数据
    Forward                              // 消费者可新代码读旧数据
    Full                                 // 双向兼容
    None                                 // 不保证兼容
)

func (cc *CompatibilityChecker) Check(old, new *TaskSchema) error {
    switch cc.level {
    case Backward:
        return cc.checkBackwardCompatible(old, new)
    case Forward:
        return cc.checkForwardCompatible(old, new)
    case Full:
        if err := cc.checkBackwardCompatible(old, new); err != nil {
            return err
        }
        return cc.checkForwardCompatible(old, new)
    default:
        return nil
    }
}

func (cc *CompatibilityChecker) checkBackwardCompatible(old, new *TaskSchema) error {
    // 新 schema 必须能被旧代码读取
    // 1. 不能删除必填字段
    // 2. 不能更改字段类型
    // 3. 可以添加可选字段
    
    var oldDef, newDef SchemaDefinition
    json.Unmarshal(old.Schema, &oldDef)
    json.Unmarshal(new.Schema, &newDef)
    
    for name, prop := range oldDef.Properties {
        newProp, exists := newDef.Properties[name]
        if !exists && contains(oldDef.Required, name) {
            return fmt.Errorf("removing required field %s breaks backward compatibility", name)
        }
        
        if exists && !typesCompatible(prop.Type, newProp.Type) {
            return fmt.Errorf("changing type of %s breaks backward compatibility", name)
        }
    }
    
    return nil
}
```

---

## Schema 演进

```go
type SchemaEvolution struct {
    registry SchemaRegistry
}

func (se *SchemaEvolution) Evolve(ctx context.Context, taskType string, changes SchemaChanges) error {
    current, _ := se.registry.GetLatestSchema(ctx, taskType)
    
    // 应用变更
    newSchema := se.applyChanges(current, changes)
    newSchema.Version = current.Version + 1
    
    // 验证新 schema
    if err := se.validateSchema(newSchema); err != nil {
        return err
    }
    
    // 检查兼容性
    if err := se.registry.checkCompatibility(current, newSchema); err != nil {
        return err
    }
    
    // 注册新版本
    return se.registry.RegisterSchema(ctx, *newSchema)
}

func (se *SchemaEvolution) applyChanges(schema *TaskSchema, changes SchemaChanges) *TaskSchema {
    var def SchemaDefinition
    json.Unmarshal(schema.Schema, &def)
    
    // 添加字段
    for _, field := range changes.AddFields {
        def.Properties[field.Name] = field
        if field.Required {
            def.Required = append(def.Required, field.Name)
        }
    }
    
    // 移除字段
    for _, name := range changes.RemoveFields {
        delete(def.Properties, name)
        def.Required = removeString(def.Required, name)
    }
    
    // 修改字段
    for name, updates := range changes.ModifyFields {
        if prop, exists := def.Properties[name]; exists {
            if updates.Type != "" {
                prop.Type = updates.Type
            }
            if updates.Description != "" {
                prop.Description = updates.Description
            }
            def.Properties[name] = prop
        }
    }
    
    newSchema := *schema
    newSchema.Schema = mustMarshal(def)
    return &newSchema
}
```
