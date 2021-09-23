// Code generated by entc, DO NOT EDIT.

package scan

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
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
func IDNotIn(ids ...string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
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
func IDGT(id string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Branch applies equality check predicate on the "branch" field. It's identical to BranchEQ.
func Branch(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBranch), v))
	})
}

// CommitID applies equality check predicate on the "commit_id" field. It's identical to CommitIDEQ.
func CommitID(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCommitID), v))
	})
}

// RequestedAt applies equality check predicate on the "requested_at" field. It's identical to RequestedAtEQ.
func RequestedAt(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRequestedAt), v))
	})
}

// ScannedAt applies equality check predicate on the "scanned_at" field. It's identical to ScannedAtEQ.
func ScannedAt(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldScannedAt), v))
	})
}

// CheckID applies equality check predicate on the "check_id" field. It's identical to CheckIDEQ.
func CheckID(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCheckID), v))
	})
}

// PullRequestTarget applies equality check predicate on the "pull_request_target" field. It's identical to PullRequestTargetEQ.
func PullRequestTarget(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPullRequestTarget), v))
	})
}

// BranchEQ applies the EQ predicate on the "branch" field.
func BranchEQ(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBranch), v))
	})
}

// BranchNEQ applies the NEQ predicate on the "branch" field.
func BranchNEQ(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldBranch), v))
	})
}

// BranchIn applies the In predicate on the "branch" field.
func BranchIn(vs ...string) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldBranch), v...))
	})
}

// BranchNotIn applies the NotIn predicate on the "branch" field.
func BranchNotIn(vs ...string) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldBranch), v...))
	})
}

// BranchGT applies the GT predicate on the "branch" field.
func BranchGT(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldBranch), v))
	})
}

// BranchGTE applies the GTE predicate on the "branch" field.
func BranchGTE(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldBranch), v))
	})
}

// BranchLT applies the LT predicate on the "branch" field.
func BranchLT(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldBranch), v))
	})
}

// BranchLTE applies the LTE predicate on the "branch" field.
func BranchLTE(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldBranch), v))
	})
}

// BranchContains applies the Contains predicate on the "branch" field.
func BranchContains(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldBranch), v))
	})
}

// BranchHasPrefix applies the HasPrefix predicate on the "branch" field.
func BranchHasPrefix(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldBranch), v))
	})
}

// BranchHasSuffix applies the HasSuffix predicate on the "branch" field.
func BranchHasSuffix(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldBranch), v))
	})
}

// BranchEqualFold applies the EqualFold predicate on the "branch" field.
func BranchEqualFold(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldBranch), v))
	})
}

// BranchContainsFold applies the ContainsFold predicate on the "branch" field.
func BranchContainsFold(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldBranch), v))
	})
}

// CommitIDEQ applies the EQ predicate on the "commit_id" field.
func CommitIDEQ(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCommitID), v))
	})
}

// CommitIDNEQ applies the NEQ predicate on the "commit_id" field.
func CommitIDNEQ(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCommitID), v))
	})
}

// CommitIDIn applies the In predicate on the "commit_id" field.
func CommitIDIn(vs ...string) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCommitID), v...))
	})
}

// CommitIDNotIn applies the NotIn predicate on the "commit_id" field.
func CommitIDNotIn(vs ...string) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCommitID), v...))
	})
}

// CommitIDGT applies the GT predicate on the "commit_id" field.
func CommitIDGT(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCommitID), v))
	})
}

// CommitIDGTE applies the GTE predicate on the "commit_id" field.
func CommitIDGTE(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCommitID), v))
	})
}

// CommitIDLT applies the LT predicate on the "commit_id" field.
func CommitIDLT(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCommitID), v))
	})
}

// CommitIDLTE applies the LTE predicate on the "commit_id" field.
func CommitIDLTE(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCommitID), v))
	})
}

// CommitIDContains applies the Contains predicate on the "commit_id" field.
func CommitIDContains(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCommitID), v))
	})
}

// CommitIDHasPrefix applies the HasPrefix predicate on the "commit_id" field.
func CommitIDHasPrefix(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCommitID), v))
	})
}

// CommitIDHasSuffix applies the HasSuffix predicate on the "commit_id" field.
func CommitIDHasSuffix(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCommitID), v))
	})
}

// CommitIDEqualFold applies the EqualFold predicate on the "commit_id" field.
func CommitIDEqualFold(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCommitID), v))
	})
}

// CommitIDContainsFold applies the ContainsFold predicate on the "commit_id" field.
func CommitIDContainsFold(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCommitID), v))
	})
}

// RequestedAtEQ applies the EQ predicate on the "requested_at" field.
func RequestedAtEQ(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRequestedAt), v))
	})
}

// RequestedAtNEQ applies the NEQ predicate on the "requested_at" field.
func RequestedAtNEQ(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRequestedAt), v))
	})
}

// RequestedAtIn applies the In predicate on the "requested_at" field.
func RequestedAtIn(vs ...int64) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldRequestedAt), v...))
	})
}

// RequestedAtNotIn applies the NotIn predicate on the "requested_at" field.
func RequestedAtNotIn(vs ...int64) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldRequestedAt), v...))
	})
}

// RequestedAtGT applies the GT predicate on the "requested_at" field.
func RequestedAtGT(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRequestedAt), v))
	})
}

// RequestedAtGTE applies the GTE predicate on the "requested_at" field.
func RequestedAtGTE(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRequestedAt), v))
	})
}

// RequestedAtLT applies the LT predicate on the "requested_at" field.
func RequestedAtLT(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRequestedAt), v))
	})
}

// RequestedAtLTE applies the LTE predicate on the "requested_at" field.
func RequestedAtLTE(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRequestedAt), v))
	})
}

// ScannedAtEQ applies the EQ predicate on the "scanned_at" field.
func ScannedAtEQ(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldScannedAt), v))
	})
}

// ScannedAtNEQ applies the NEQ predicate on the "scanned_at" field.
func ScannedAtNEQ(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldScannedAt), v))
	})
}

// ScannedAtIn applies the In predicate on the "scanned_at" field.
func ScannedAtIn(vs ...int64) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldScannedAt), v...))
	})
}

// ScannedAtNotIn applies the NotIn predicate on the "scanned_at" field.
func ScannedAtNotIn(vs ...int64) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldScannedAt), v...))
	})
}

// ScannedAtGT applies the GT predicate on the "scanned_at" field.
func ScannedAtGT(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldScannedAt), v))
	})
}

// ScannedAtGTE applies the GTE predicate on the "scanned_at" field.
func ScannedAtGTE(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldScannedAt), v))
	})
}

// ScannedAtLT applies the LT predicate on the "scanned_at" field.
func ScannedAtLT(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldScannedAt), v))
	})
}

// ScannedAtLTE applies the LTE predicate on the "scanned_at" field.
func ScannedAtLTE(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldScannedAt), v))
	})
}

// CheckIDEQ applies the EQ predicate on the "check_id" field.
func CheckIDEQ(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCheckID), v))
	})
}

// CheckIDNEQ applies the NEQ predicate on the "check_id" field.
func CheckIDNEQ(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCheckID), v))
	})
}

// CheckIDIn applies the In predicate on the "check_id" field.
func CheckIDIn(vs ...int64) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCheckID), v...))
	})
}

// CheckIDNotIn applies the NotIn predicate on the "check_id" field.
func CheckIDNotIn(vs ...int64) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCheckID), v...))
	})
}

// CheckIDGT applies the GT predicate on the "check_id" field.
func CheckIDGT(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCheckID), v))
	})
}

// CheckIDGTE applies the GTE predicate on the "check_id" field.
func CheckIDGTE(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCheckID), v))
	})
}

// CheckIDLT applies the LT predicate on the "check_id" field.
func CheckIDLT(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCheckID), v))
	})
}

// CheckIDLTE applies the LTE predicate on the "check_id" field.
func CheckIDLTE(v int64) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCheckID), v))
	})
}

// CheckIDIsNil applies the IsNil predicate on the "check_id" field.
func CheckIDIsNil() predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldCheckID)))
	})
}

// CheckIDNotNil applies the NotNil predicate on the "check_id" field.
func CheckIDNotNil() predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldCheckID)))
	})
}

// PullRequestTargetEQ applies the EQ predicate on the "pull_request_target" field.
func PullRequestTargetEQ(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetNEQ applies the NEQ predicate on the "pull_request_target" field.
func PullRequestTargetNEQ(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetIn applies the In predicate on the "pull_request_target" field.
func PullRequestTargetIn(vs ...string) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldPullRequestTarget), v...))
	})
}

// PullRequestTargetNotIn applies the NotIn predicate on the "pull_request_target" field.
func PullRequestTargetNotIn(vs ...string) predicate.Scan {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Scan(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldPullRequestTarget), v...))
	})
}

// PullRequestTargetGT applies the GT predicate on the "pull_request_target" field.
func PullRequestTargetGT(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetGTE applies the GTE predicate on the "pull_request_target" field.
func PullRequestTargetGTE(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetLT applies the LT predicate on the "pull_request_target" field.
func PullRequestTargetLT(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetLTE applies the LTE predicate on the "pull_request_target" field.
func PullRequestTargetLTE(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetContains applies the Contains predicate on the "pull_request_target" field.
func PullRequestTargetContains(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetHasPrefix applies the HasPrefix predicate on the "pull_request_target" field.
func PullRequestTargetHasPrefix(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetHasSuffix applies the HasSuffix predicate on the "pull_request_target" field.
func PullRequestTargetHasSuffix(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetIsNil applies the IsNil predicate on the "pull_request_target" field.
func PullRequestTargetIsNil() predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldPullRequestTarget)))
	})
}

// PullRequestTargetNotNil applies the NotNil predicate on the "pull_request_target" field.
func PullRequestTargetNotNil() predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldPullRequestTarget)))
	})
}

// PullRequestTargetEqualFold applies the EqualFold predicate on the "pull_request_target" field.
func PullRequestTargetEqualFold(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPullRequestTarget), v))
	})
}

// PullRequestTargetContainsFold applies the ContainsFold predicate on the "pull_request_target" field.
func PullRequestTargetContainsFold(v string) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPullRequestTarget), v))
	})
}

// HasRepository applies the HasEdge predicate on the "repository" edge.
func HasRepository() predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(RepositoryTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, RepositoryTable, RepositoryPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRepositoryWith applies the HasEdge predicate on the "repository" edge with a given conditions (other predicates).
func HasRepositoryWith(preds ...predicate.Repository) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(RepositoryInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, RepositoryTable, RepositoryPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPackages applies the HasEdge predicate on the "packages" edge.
func HasPackages() predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PackagesTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, PackagesTable, PackagesPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPackagesWith applies the HasEdge predicate on the "packages" edge with a given conditions (other predicates).
func HasPackagesWith(preds ...predicate.PackageRecord) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PackagesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, PackagesTable, PackagesPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Scan) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Scan) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
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
func Not(p predicate.Scan) predicate.Scan {
	return predicate.Scan(func(s *sql.Selector) {
		p(s.Not())
	})
}