package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

// CheckRule holds the schema definition for the CheckRule entity.
type CheckRule struct {
	ent.Schema
}

// Fields of the CheckRule.
func (CheckRule) Fields() []ent.Field {
	return []ent.Field{
		field.String("check_result").GoType(types.GitHubCheckResult("")),
	}
}

// Edges of the CheckRule.
func (CheckRule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("severity", Severity.Type).Unique(),
	}

}
