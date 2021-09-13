package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Branch holds the schema definition for the Branch entity.
type Branch struct {
	ent.Schema
}

// Fields of the Branch.
func (Branch) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").Unique(),
		field.String("owner"),
		field.String("repo_name"),
		field.String("name"),
	}
}

// Edges of the Branch.
func (Branch) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("scan", Scan.Type),
	}
}
