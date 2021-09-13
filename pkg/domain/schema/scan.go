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
		field.String("commit_id").Immutable(),
		field.Int64("requested_at"),
		field.Int64("scanned_at").Optional(),
		field.Int64("check_id"),
		field.String("pull_request_target").Optional(),
	}
}

// Edges of the Scan.
func (Scan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("target", Branch.Type).Ref("scan"),
		edge.To("packages", PackageRecord.Type),
	}
}
