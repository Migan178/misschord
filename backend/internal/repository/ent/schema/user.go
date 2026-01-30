package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("profile").Default("/defaults/default_profile.png"),
		field.String("handle").Unique(),
		field.String("email").Unique().Sensitive(),
		field.String("hashed_password").Sensitive(),
		field.String("description").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).StructTag(`json:"createdAt"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("messages", Message.Type),
		edge.To("rooms", Room.Type),
	}
}
