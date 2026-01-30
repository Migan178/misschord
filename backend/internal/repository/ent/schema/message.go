package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.Int("author_id"),
		field.Int("room_id"),
		field.String("message"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).
			Ref("messages").
			Unique().
			Required().
			Field("author_id"),
		edge.From("room", Room.Type).
			Ref("messages").
			Unique().
			Required().
			Field("room_id"),
	}
}

func (Message) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("room_id"),
	}
}
