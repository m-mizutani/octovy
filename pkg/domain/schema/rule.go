package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Rule holds the schema definition for the Rule entity.
type Rule struct {
	ent.Schema
}

// Fields of the Rule.
func (Rule) Fields() []ent.Field {
	return []ent.Field{
		field.String("action"),
	}
}

// Edges of the Rule.
func (Rule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("severity", Severity.Type).Unique(),
	}
}
