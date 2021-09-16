package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Scan holds the schema definition for the Scan entity.
type Scan struct {
	ent.Schema
}

// Fields of the Scan.
func (Scan) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().Unique(),
		field.String("branch").Immutable(),
		field.String("commit_id").Immutable(),
		field.Int64("requested_at"),
		field.Int64("scanned_at"),
		field.Int64("check_id").Optional(),
		field.String("pull_request_target").Optional(),
	}
}

// Edges of the Scan.
func (Scan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository", Repository.Type).Ref("scan"),
		edge.To("packages", PackageRecord.Type),
	}
}
