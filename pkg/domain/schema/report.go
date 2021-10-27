package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Report holds the schema definition for the Report entity.
type Report struct {
	ent.Schema
}

// Fields of the Report.
func (Report) Fields() []ent.Field {
	return []ent.Field{
		field.String("scanner").Immutable(),
		field.String("resource_type"),
		field.String("resource_name"),
		field.Int64("scanned_at"),
		field.Int64("requested_at").Optional(),
	}
}

// Edges of the Report.
func (Report) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("objects", Object.Type),
	}
}
