package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// VulnStatusIndex holds the schema definition for the VulnStatusIndex entity.
type VulnStatusIndex struct {
	ent.Schema
}

// Fields of the VulnStatusIndex.
func (VulnStatusIndex) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Immutable().Unique(),
	}
}

// Edges of the VulnStatusIndex.
func (VulnStatusIndex) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("latest", VulnStatus.Type).Unique(),
		edge.To("status", VulnStatus.Type),
	}
}
