package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Object holds the schema definition for the Object entity.
type Object struct {
	ent.Schema
}

// Fields of the Object.
func (Object) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").Immutable().Comment("Identifiable key of the object in same resource"),
		field.String("name").Immutable().Comment("Human readable name of the object"),
		field.String("description").Optional(),
		field.String("version").Optional(),
	}
}

// Edges of the Object.
func (Object) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vulnerabilities", Vulnerability.Type),
		edge.From("report", Report.Type).Ref("objects"),
	}
}
