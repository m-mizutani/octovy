package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.Int("user_id"),
		field.String("token").Sensitive().Immutable().NotEmpty(),
		field.Int64("created_at").Immutable(),
		field.Int64("expires_at").Immutable(),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("login", User.Type).Unique(),
	}
}
