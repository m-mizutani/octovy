package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

// VulnStatus holds the schema definition for the VulnStatus entity.
type VulnStatus struct {
	ent.Schema
}

// Fields of the VulnStatus.
func (VulnStatus) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").GoType(types.VulnStatusType("")),
		field.String("source"),
		field.String("pkg_name"),
		field.String("pkg_type"),
		field.String("vuln_id"),
		field.Int64("expires_at"),
		field.Int64("created_at"),
		field.String("comment"),
	}
}

// Edges of the VulnStatus.
func (VulnStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("author", User.Type).Unique(),
	}
}
