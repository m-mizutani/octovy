package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Repository holds the schema definition for the Repository entity.
type Repository struct {
	ent.Schema
}

// Fields of the Repository.
func (Repository) Fields() []ent.Field {
	return []ent.Field{
		field.String("owner"),
		field.String("name"),
		field.Int64("install_id").Optional(),
		field.String("url").Optional(),
		field.String("avatar_url").Optional().Nillable(),
		field.String("default_branch").Optional().Nillable(),
	}
}

// Edges of the Repository.
func (Repository) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("scan", Scan.Type).Comment("All scan results"),
		edge.To("main", Scan.Type).Comment("Scan results for default branch"),
		edge.To("latest", Scan.Type).Unique().Comment("A latest Scan result for default branch"),
		edge.To("status", VulnStatusIndex.Type),
		edge.To("labels", RepoLabel.Type),
	}
}

func (Repository) Index() []ent.Index {
	return []ent.Index{
		index.Fields("owner", "name").Unique(),
	}
}
