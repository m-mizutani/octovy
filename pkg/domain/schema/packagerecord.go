package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PackageRecord holds the schema definition for the PackageRecord entity.
type PackageRecord struct {
	ent.Schema
}

// Fields of the PackageRecord.
func (PackageRecord) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").Immutable(),
		field.String("source").Immutable(),
		field.String("name").Immutable(),
		field.String("version").Immutable(),
		field.Strings("vuln_ids"),
	}
}

// Edges of the PackageRecord.
func (PackageRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("scan", Scan.Type).Ref("packages"),
		edge.To("vulnerabilities", Vulnerability.Type),
	}
}
