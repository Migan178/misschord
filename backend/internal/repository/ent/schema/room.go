package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Room holds the schema definition for the Room entity.
type Room struct {
	ent.Schema
}

// Fields of the Room.
func (Room) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("room_type").
			Values("DM", "CHANNEL"),
		field.String("dm_key").
			Optional().
			Nillable().
			Unique(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Room.
func (Room) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("members", User.Type).
			Ref("rooms"),
		edge.To("messages", Message.Type),
	}
}
