package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RepoLabel holds the schema definition for the RepoLabel entity.
type RepoLabel struct {
	ent.Schema
}

// Fields of the RepoLabel.
func (RepoLabel) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("description"),
		field.String("color"),
	}
}

// Edges of the RepoLabel.
func (RepoLabel) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repos", Repository.Type).Ref("labels"),
	}
}
