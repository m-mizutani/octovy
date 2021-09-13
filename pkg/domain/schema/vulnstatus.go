package schema

import (
	"entgo.io/ent"
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
		field.String("id").NotEmpty().Immutable().Unique(),
		field.Enum("status").GoType(types.VulnStatusType("")),
		field.String("vuln_id"),
		field.Int64("expires_at"),
		field.Int64("created_at"),
	}
}

// Edges of the VulnStatus.
func (VulnStatus) Edges() []ent.Edge {
	return nil
}
