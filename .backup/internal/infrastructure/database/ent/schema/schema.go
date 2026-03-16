package schema

import (
	"entgo.io/ent"
)

// 导出所有 schema
// 注意：框架不提供具体的 Schema 定义
// 用户需要在自己的项目中定义 Ent Schema
var (
	Schemas = []ent.Interface{
		&User{},
	}
)
