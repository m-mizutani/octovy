package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// AuthStateCache holds the schema definition for the AuthStateCache entity.
type AuthStateCache struct {
	ent.Schema
}

// Fields of the AuthStateCache.
func (AuthStateCache) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.Int64("expires_at"),
	}
}

// Edges of the AuthStateCache.
func (AuthStateCache) Edges() []ent.Edge {
	return nil
}
