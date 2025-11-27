package reflect

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Metadata 程序元数据
// 提供程序的自解释能力，包括类型信息、函数信息、结构信息等
type Metadata struct {
	Type     TypeMetadata     `json:"type"`
	Function FunctionMetadata `json:"function"`
	Struct   StructMetadata   `json:"struct"`
}

// TypeMetadata 类型元数据
type TypeMetadata struct {
	Name    string                 `json:"name"`
	Kind    string                 `json:"kind"`
	Package string                 `json:"package"`
	Methods []MethodMetadata       `json:"methods,omitempty"`
	Fields  []FieldMetadata        `json:"fields,omitempty"`
	Tags    map[string]interface{} `json:"tags,omitempty"`
}

// FunctionMetadata 函数元数据
type FunctionMetadata struct {
	Name     string   `json:"name"`
	Package  string   `json:"package"`
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Inputs   []string `json:"inputs,omitempty"`
	Outputs  []string `json:"outputs,omitempty"`
	Variadic bool     `json:"variadic"`
}

// StructMetadata 结构体元数据
type StructMetadata struct {
	Name   string          `json:"name"`
	Fields []FieldMetadata `json:"fields"`
	Tags   map[string]map[string]string `json:"tags,omitempty"`
}

// MethodMetadata 方法元数据
type MethodMetadata struct {
	Name     string   `json:"name"`
	Inputs   []string `json:"inputs,omitempty"`
	Outputs  []string `json:"outputs,omitempty"`
	Exported bool     `json:"exported"`
}

// FieldMetadata 字段元数据
type FieldMetadata struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Tag      string                 `json:"tag,omitempty"`
	Exported bool                   `json:"exported"`
	Tags     map[string]interface{} `json:"tags,omitempty"`
}

// Inspector 反射检查器
type Inspector struct{}

// NewInspector 创建反射检查器
func NewInspector() *Inspector {
	return &Inspector{}
}

// InspectType 检查类型
func (i *Inspector) InspectType(v interface{}) TypeMetadata {
	rt := reflect.TypeOf(v)
	if rt == nil {
		return TypeMetadata{}
	}

	// 处理指针类型
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	metadata := TypeMetadata{
		Name:    rt.Name(),
		Kind:    rt.Kind().String(),
		Package: rt.PkgPath(),
		Methods: make([]MethodMetadata, 0),
		Fields:  make([]FieldMetadata, 0),
		Tags:    make(map[string]interface{}),
	}

	// 获取方法
	for j := 0; j < rt.NumMethod(); j++ {
		method := rt.Method(j)
		metadata.Methods = append(metadata.Methods, MethodMetadata{
			Name:     method.Name,
			Inputs:   i.getTypeNames(method.Type.In),
			Outputs:  i.getTypeNames(method.Type.Out),
			Exported: method.IsExported(),
		})
	}

	// 获取字段（仅结构体）
	if rt.Kind() == reflect.Struct {
		for j := 0; j < rt.NumField(); j++ {
			field := rt.Field(j)
			fieldMeta := FieldMetadata{
				Name:     field.Name,
				Type:     field.Type.String(),
				Tag:      string(field.Tag),
				Exported: field.IsExported(),
				Tags:     i.parseTags(field.Tag),
			}
			metadata.Fields = append(metadata.Fields, fieldMeta)
		}
	}

	return metadata
}

// InspectFunction 检查函数
func (i *Inspector) InspectFunction(fn interface{}) FunctionMetadata {
	rv := reflect.ValueOf(fn)
	if rv.Kind() != reflect.Func {
		return FunctionMetadata{}
	}

	rt := rv.Type()
	pc := runtime.FuncForPC(rv.Pointer())
	if pc == nil {
		return FunctionMetadata{}
	}

	file, line := pc.FileLine(pc.Entry())
	name := pc.Name()

	// 解析包名和函数名
	parts := strings.Split(name, ".")
	var pkgName, funcName string
	if len(parts) >= 2 {
		pkgName = strings.Join(parts[:len(parts)-1], ".")
		funcName = parts[len(parts)-1]
	} else {
		funcName = name
	}

	metadata := FunctionMetadata{
		Name:     funcName,
		Package:  pkgName,
		File:     file,
		Line:     line,
		Inputs:   make([]string, 0),
		Outputs:  make([]string, 0),
		Variadic: rt.IsVariadic(),
	}

	// 获取输入参数
	for j := 0; j < rt.NumIn(); j++ {
		metadata.Inputs = append(metadata.Inputs, rt.In(j).String())
	}

	// 获取输出参数
	for j := 0; j < rt.NumOut(); j++ {
		metadata.Outputs = append(metadata.Outputs, rt.Out(j).String())
	}

	return metadata
}

// InspectStruct 检查结构体
func (i *Inspector) InspectStruct(v interface{}) StructMetadata {
	rt := reflect.TypeOf(v)
	if rt == nil {
		return StructMetadata{}
	}

	// 处理指针类型
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		return StructMetadata{}
	}

	metadata := StructMetadata{
		Name:   rt.Name(),
		Fields: make([]FieldMetadata, 0),
		Tags:   make(map[string]map[string]string),
	}

	for j := 0; j < rt.NumField(); j++ {
		field := rt.Field(j)
		fieldMeta := FieldMetadata{
			Name:     field.Name,
			Type:     field.Type.String(),
			Tag:      string(field.Tag),
			Exported: field.IsExported(),
			Tags:     i.parseTags(field.Tag),
		}
		metadata.Fields = append(metadata.Fields, fieldMeta)

		// 解析标签
		if field.Tag != "" {
			tags := make(map[string]string)
			for _, tag := range []string{"json", "xml", "yaml", "db", "validate"} {
				if value := field.Tag.Get(tag); value != "" {
					tags[tag] = value
				}
			}
			if len(tags) > 0 {
				metadata.Tags[field.Name] = tags
			}
		}
	}

	return metadata
}

// GetTypeName 获取类型名称
func (i *Inspector) GetTypeName(v interface{}) string {
	rt := reflect.TypeOf(v)
	if rt == nil {
		return "nil"
	}

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	return rt.String()
}

// GetPackageName 获取包名
func (i *Inspector) GetPackageName(v interface{}) string {
	rt := reflect.TypeOf(v)
	if rt == nil {
		return ""
	}

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	return rt.PkgPath()
}

// getTypeNames 获取类型名称列表
func (i *Inspector) getTypeNames(types []reflect.Type) []string {
	result := make([]string, len(types))
	for j, t := range types {
		result[j] = t.String()
	}
	return result
}

// parseTags 解析标签
func (i *Inspector) parseTags(tag reflect.StructTag) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range []string{"json", "xml", "yaml", "db", "validate", "form"} {
		if value := tag.Get(key); value != "" {
			result[key] = value
		}
	}
	return result
}

// Describe 描述对象
// 返回对象的完整描述信息
func (i *Inspector) Describe(v interface{}) string {
	rt := reflect.TypeOf(v)
	if rt == nil {
		return "nil"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Type: %s\n", rt.String()))
	builder.WriteString(fmt.Sprintf("Kind: %s\n", rt.Kind().String()))

	if rt.PkgPath() != "" {
		builder.WriteString(fmt.Sprintf("Package: %s\n", rt.PkgPath()))
	}

	if rt.Kind() == reflect.Struct {
		builder.WriteString("\nFields:\n")
		for j := 0; j < rt.NumField(); j++ {
			field := rt.Field(j)
			builder.WriteString(fmt.Sprintf("  %s: %s\n", field.Name, field.Type.String()))
		}
	}

	return builder.String()
}
