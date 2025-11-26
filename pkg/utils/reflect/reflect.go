package reflect

import (
	"fmt"
	"reflect"
)

// GetType 获取值的类型名称
func GetType(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// GetKind 获取值的Kind
func GetKind(v interface{}) reflect.Kind {
	return reflect.TypeOf(v).Kind()
}

// IsNil 检查值是否为nil
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// IsZero 检查值是否为零值
func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	return rv.IsZero()
}

// IsPointer 检查值是否为指针
func IsPointer(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Ptr
}

// IsSlice 检查值是否为切片
func IsSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

// IsMap 检查值是否为映射
func IsMap(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}

// IsStruct 检查值是否为结构体
func IsStruct(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

// IsInterface 检查值是否为接口
func IsInterface(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Interface
}

// IsFunc 检查值是否为函数
func IsFunc(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}

// IsChan 检查值是否为通道
func IsChan(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Chan
}

// Dereference 解引用指针，如果不是指针则返回原值
func Dereference(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		return rv.Elem().Interface()
	}
	return v
}

// GetField 获取结构体字段的值
func GetField(v interface{}, fieldName string) (interface{}, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("value is not a struct")
	}
	field := rv.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}
	return field.Interface(), nil
}

// SetField 设置结构体字段的值
func SetField(v interface{}, fieldName string, value interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be a pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("value is not a struct")
	}
	field := rv.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}
	fieldValue := reflect.ValueOf(value)
	if fieldValue.Type() != field.Type() {
		return fmt.Errorf("type mismatch: expected %s, got %s", field.Type(), fieldValue.Type())
	}
	field.Set(fieldValue)
	return nil
}

// HasField 检查结构体是否有指定字段
func HasField(v interface{}, fieldName string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	field := rv.FieldByName(fieldName)
	return field.IsValid()
}

// GetFieldNames 获取结构体的所有字段名
func GetFieldNames(v interface{}) []string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}
	t := rv.Type()
	fields := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fields[i] = t.Field(i).Name
	}
	return fields
}

// GetFieldTags 获取结构体字段的标签
func GetFieldTags(v interface{}, fieldName string) (map[string]string, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("value is not a struct")
	}
	t := rv.Type()
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}
	tags := make(map[string]string)
	for _, tag := range []string{"json", "xml", "yaml", "db", "form", "validate"} {
		if tagValue := field.Tag.Get(tag); tagValue != "" {
			tags[tag] = tagValue
		}
	}
	return tags, nil
}

// CallMethod 调用方法
func CallMethod(v interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	rv := reflect.ValueOf(v)
	method := rv.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found", methodName)
	}
	methodType := method.Type()
	if methodType.NumIn() != len(args) {
		return nil, fmt.Errorf("method %s expects %d arguments, got %d", methodName, methodType.NumIn(), len(args))
	}
	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValues[i] = reflect.ValueOf(arg)
	}
	results := method.Call(argValues)
	resultValues := make([]interface{}, len(results))
	for i, result := range results {
		resultValues[i] = result.Interface()
	}
	return resultValues, nil
}

// HasMethod 检查值是否有指定方法
func HasMethod(v interface{}, methodName string) bool {
	rv := reflect.ValueOf(v)
	method := rv.MethodByName(methodName)
	return method.IsValid()
}

// GetMethodNames 获取值的所有方法名
func GetMethodNames(v interface{}) []string {
	rv := reflect.ValueOf(v)
	t := rv.Type()
	methods := make([]string, t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		methods[i] = t.Method(i).Name
	}
	return methods
}

// NewInstance 创建类型的新实例
func NewInstance(v interface{}) interface{} {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}

// NewSlice 创建切片的新实例
func NewSlice(v interface{}, length, capacity int) interface{} {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		return reflect.MakeSlice(t, length, capacity).Interface()
	}
	return nil
}

// NewMap 创建映射的新实例
func NewMap(keyType, valueType interface{}) interface{} {
	kt := reflect.TypeOf(keyType)
	vt := reflect.TypeOf(valueType)
	if kt.Kind() == reflect.Ptr {
		kt = kt.Elem()
	}
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	mapType := reflect.MapOf(kt, vt)
	return reflect.MakeMap(mapType).Interface()
}

// Convert 转换值的类型
func Convert(v interface{}, targetType interface{}) (interface{}, error) {
	rv := reflect.ValueOf(v)
	tt := reflect.TypeOf(targetType)
	if tt.Kind() == reflect.Ptr {
		tt = tt.Elem()
	}
	if !rv.CanConvert(tt) {
		return nil, fmt.Errorf("cannot convert %s to %s", rv.Type(), tt)
	}
	return rv.Convert(tt).Interface(), nil
}

// DeepEqual 深度比较两个值是否相等
func DeepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// Copy 深度复制值
func Copy(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	return rv.Interface()
}

// GetSliceElement 获取切片元素
func GetSliceElement(slice interface{}, index int) (interface{}, error) {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value is not a slice")
	}
	if index < 0 || index >= rv.Len() {
		return nil, fmt.Errorf("index out of range")
	}
	return rv.Index(index).Interface(), nil
}

// SetSliceElement 设置切片元素
func SetSliceElement(slice interface{}, index int, value interface{}) error {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("value is not a slice")
	}
	if index < 0 || index >= rv.Len() {
		return fmt.Errorf("index out of range")
	}
	element := rv.Index(index)
	if !element.CanSet() {
		return fmt.Errorf("element cannot be set")
	}
	elementValue := reflect.ValueOf(value)
	if elementValue.Type() != element.Type() {
		return fmt.Errorf("type mismatch")
	}
	element.Set(elementValue)
	return nil
}

// GetMapValue 获取映射的值
func GetMapValue(m interface{}, key interface{}) (interface{}, bool) {
	rv := reflect.ValueOf(m)
	if rv.Kind() != reflect.Map {
		return nil, false
	}
	keyValue := reflect.ValueOf(key)
	if !keyValue.Type().AssignableTo(rv.Type().Key()) {
		return nil, false
	}
	value := rv.MapIndex(keyValue)
	if !value.IsValid() {
		return nil, false
	}
	return value.Interface(), true
}

// SetMapValue 设置映射的值
func SetMapValue(m interface{}, key, value interface{}) error {
	rv := reflect.ValueOf(m)
	if rv.Kind() != reflect.Map {
		return fmt.Errorf("value is not a map")
	}
	keyValue := reflect.ValueOf(key)
	valueValue := reflect.ValueOf(value)
	if !keyValue.Type().AssignableTo(rv.Type().Key()) {
		return fmt.Errorf("key type mismatch")
	}
	if !valueValue.Type().AssignableTo(rv.Type().Elem()) {
		return fmt.Errorf("value type mismatch")
	}
	rv.SetMapIndex(keyValue, valueValue)
	return nil
}

// GetLength 获取值的长度（切片、映射、字符串、数组）
func GetLength(v interface{}) (int, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Array:
		return rv.Len(), nil
	default:
		return 0, fmt.Errorf("value does not have a length")
	}
}

// GetCapacity 获取值的容量（切片、数组）
func GetCapacity(v interface{}) (int, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return rv.Cap(), nil
	default:
		return 0, fmt.Errorf("value does not have a capacity")
	}
}

// IsAssignable 检查值是否可以赋值给目标类型
func IsAssignable(value, targetType interface{}) bool {
	rv := reflect.ValueOf(value)
	tt := reflect.TypeOf(targetType)
	if tt.Kind() == reflect.Ptr {
		tt = tt.Elem()
	}
	return rv.Type().AssignableTo(tt)
}

// IsConvertible 检查值是否可以转换为目标类型
func IsConvertible(value, targetType interface{}) bool {
	rv := reflect.ValueOf(value)
	tt := reflect.TypeOf(targetType)
	if tt.Kind() == reflect.Ptr {
		tt = tt.Elem()
	}
	return rv.Type().ConvertibleTo(tt)
}
