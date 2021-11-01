package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Severity holds the schema definition for the Severity entity.
type Severity struct {
	ent.Schema
}

// Fields of the Severity.
func (Severity) Fields() []ent.Field {
	return []ent.Field{
		field.String("label").NotEmpty().Unique(),
	}
}

// Edges of the Severity.
func (Severity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("vulnerabilities", Vulnerability.Type).Ref("sev"),
	}
}
