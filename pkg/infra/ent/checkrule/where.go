// Code generated by entc, DO NOT EDIT.

package checkrule

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CheckResult applies equality check predicate on the "check_result" field. It's identical to CheckResultEQ.
func CheckResult(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCheckResult), vc))
	})
}

// CheckResultEQ applies the EQ predicate on the "check_result" field.
func CheckResultEQ(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCheckResult), vc))
	})
}

// CheckResultNEQ applies the NEQ predicate on the "check_result" field.
func CheckResultNEQ(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCheckResult), vc))
	})
}

// CheckResultIn applies the In predicate on the "check_result" field.
func CheckResultIn(vs ...types.GitHubCheckResult) predicate.CheckRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = string(vs[i])
	}
	return predicate.CheckRule(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCheckResult), v...))
	})
}

// CheckResultNotIn applies the NotIn predicate on the "check_result" field.
func CheckResultNotIn(vs ...types.GitHubCheckResult) predicate.CheckRule {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = string(vs[i])
	}
	return predicate.CheckRule(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCheckResult), v...))
	})
}

// CheckResultGT applies the GT predicate on the "check_result" field.
func CheckResultGT(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCheckResult), vc))
	})
}

// CheckResultGTE applies the GTE predicate on the "check_result" field.
func CheckResultGTE(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCheckResult), vc))
	})
}

// CheckResultLT applies the LT predicate on the "check_result" field.
func CheckResultLT(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCheckResult), vc))
	})
}

// CheckResultLTE applies the LTE predicate on the "check_result" field.
func CheckResultLTE(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCheckResult), vc))
	})
}

// CheckResultContains applies the Contains predicate on the "check_result" field.
func CheckResultContains(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCheckResult), vc))
	})
}

// CheckResultHasPrefix applies the HasPrefix predicate on the "check_result" field.
func CheckResultHasPrefix(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCheckResult), vc))
	})
}

// CheckResultHasSuffix applies the HasSuffix predicate on the "check_result" field.
func CheckResultHasSuffix(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCheckResult), vc))
	})
}

// CheckResultEqualFold applies the EqualFold predicate on the "check_result" field.
func CheckResultEqualFold(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCheckResult), vc))
	})
}

// CheckResultContainsFold applies the ContainsFold predicate on the "check_result" field.
func CheckResultContainsFold(v types.GitHubCheckResult) predicate.CheckRule {
	vc := string(v)
	return predicate.CheckRule(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCheckResult), vc))
	})
}

// HasSeverity applies the HasEdge predicate on the "severity" edge.
func HasSeverity() predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SeverityTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SeverityTable, SeverityColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSeverityWith applies the HasEdge predicate on the "severity" edge with a given conditions (other predicates).
func HasSeverityWith(preds ...predicate.Severity) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SeverityInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SeverityTable, SeverityColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.CheckRule) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.CheckRule) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.CheckRule) predicate.CheckRule {
	return predicate.CheckRule(func(s *sql.Selector) {
		p(s.Not())
	})
}
