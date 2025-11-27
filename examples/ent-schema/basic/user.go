// Ent Schema 定义示例 - User 实体
// 这是一个示例，展示如何定义 Ent Schema
// 用户应该在自己的项目中定义自己的业务 Schema
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// User 用户实体（示例）
type User struct {
	ent.Schema
}

// Fields 定义用户字段
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("email").
			Unique().
			NotEmpty().
			MaxLen(255),
		field.String("name").
			NotEmpty().
			MaxLen(100),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Indexes 定义索引
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("created_at"),
	}
}
