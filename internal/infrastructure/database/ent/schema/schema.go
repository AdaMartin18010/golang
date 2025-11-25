package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
)

// 导出所有 schema
var (
	Schemas = []ent.Interface{
		&User{},
	}
)

// 实现 schema.Annotation 接口
func (User) Annotations() []schema.Annotation {
	return nil
}
