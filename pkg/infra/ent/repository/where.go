// Code generated by entc, DO NOT EDIT.

package repository

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/m-mizutani/octovy/pkg/infra/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
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
func IDGT(id int) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Owner applies equality check predicate on the "owner" field. It's identical to OwnerEQ.
func Owner(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOwner), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// InstallID applies equality check predicate on the "install_id" field. It's identical to InstallIDEQ.
func InstallID(v int64) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldInstallID), v))
	})
}

// URL applies equality check predicate on the "url" field. It's identical to URLEQ.
func URL(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldURL), v))
	})
}

// AvatarURL applies equality check predicate on the "avatar_url" field. It's identical to AvatarURLEQ.
func AvatarURL(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAvatarURL), v))
	})
}

// DefaultBranch applies equality check predicate on the "default_branch" field. It's identical to DefaultBranchEQ.
func DefaultBranch(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDefaultBranch), v))
	})
}

// OwnerEQ applies the EQ predicate on the "owner" field.
func OwnerEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOwner), v))
	})
}

// OwnerNEQ applies the NEQ predicate on the "owner" field.
func OwnerNEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldOwner), v))
	})
}

// OwnerIn applies the In predicate on the "owner" field.
func OwnerIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldOwner), v...))
	})
}

// OwnerNotIn applies the NotIn predicate on the "owner" field.
func OwnerNotIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldOwner), v...))
	})
}

// OwnerGT applies the GT predicate on the "owner" field.
func OwnerGT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldOwner), v))
	})
}

// OwnerGTE applies the GTE predicate on the "owner" field.
func OwnerGTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldOwner), v))
	})
}

// OwnerLT applies the LT predicate on the "owner" field.
func OwnerLT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldOwner), v))
	})
}

// OwnerLTE applies the LTE predicate on the "owner" field.
func OwnerLTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldOwner), v))
	})
}

// OwnerContains applies the Contains predicate on the "owner" field.
func OwnerContains(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldOwner), v))
	})
}

// OwnerHasPrefix applies the HasPrefix predicate on the "owner" field.
func OwnerHasPrefix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldOwner), v))
	})
}

// OwnerHasSuffix applies the HasSuffix predicate on the "owner" field.
func OwnerHasSuffix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldOwner), v))
	})
}

// OwnerEqualFold applies the EqualFold predicate on the "owner" field.
func OwnerEqualFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldOwner), v))
	})
}

// OwnerContainsFold applies the ContainsFold predicate on the "owner" field.
func OwnerContainsFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldOwner), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// InstallIDEQ applies the EQ predicate on the "install_id" field.
func InstallIDEQ(v int64) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldInstallID), v))
	})
}

// InstallIDNEQ applies the NEQ predicate on the "install_id" field.
func InstallIDNEQ(v int64) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldInstallID), v))
	})
}

// InstallIDIn applies the In predicate on the "install_id" field.
func InstallIDIn(vs ...int64) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldInstallID), v...))
	})
}

// InstallIDNotIn applies the NotIn predicate on the "install_id" field.
func InstallIDNotIn(vs ...int64) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldInstallID), v...))
	})
}

// InstallIDGT applies the GT predicate on the "install_id" field.
func InstallIDGT(v int64) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldInstallID), v))
	})
}

// InstallIDGTE applies the GTE predicate on the "install_id" field.
func InstallIDGTE(v int64) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldInstallID), v))
	})
}

// InstallIDLT applies the LT predicate on the "install_id" field.
func InstallIDLT(v int64) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldInstallID), v))
	})
}

// InstallIDLTE applies the LTE predicate on the "install_id" field.
func InstallIDLTE(v int64) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldInstallID), v))
	})
}

// InstallIDIsNil applies the IsNil predicate on the "install_id" field.
func InstallIDIsNil() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldInstallID)))
	})
}

// InstallIDNotNil applies the NotNil predicate on the "install_id" field.
func InstallIDNotNil() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldInstallID)))
	})
}

// URLEQ applies the EQ predicate on the "url" field.
func URLEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldURL), v))
	})
}

// URLNEQ applies the NEQ predicate on the "url" field.
func URLNEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldURL), v))
	})
}

// URLIn applies the In predicate on the "url" field.
func URLIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldURL), v...))
	})
}

// URLNotIn applies the NotIn predicate on the "url" field.
func URLNotIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldURL), v...))
	})
}

// URLGT applies the GT predicate on the "url" field.
func URLGT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldURL), v))
	})
}

// URLGTE applies the GTE predicate on the "url" field.
func URLGTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldURL), v))
	})
}

// URLLT applies the LT predicate on the "url" field.
func URLLT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldURL), v))
	})
}

// URLLTE applies the LTE predicate on the "url" field.
func URLLTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldURL), v))
	})
}

// URLContains applies the Contains predicate on the "url" field.
func URLContains(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldURL), v))
	})
}

// URLHasPrefix applies the HasPrefix predicate on the "url" field.
func URLHasPrefix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldURL), v))
	})
}

// URLHasSuffix applies the HasSuffix predicate on the "url" field.
func URLHasSuffix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldURL), v))
	})
}

// URLIsNil applies the IsNil predicate on the "url" field.
func URLIsNil() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldURL)))
	})
}

// URLNotNil applies the NotNil predicate on the "url" field.
func URLNotNil() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldURL)))
	})
}

// URLEqualFold applies the EqualFold predicate on the "url" field.
func URLEqualFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldURL), v))
	})
}

// URLContainsFold applies the ContainsFold predicate on the "url" field.
func URLContainsFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldURL), v))
	})
}

// AvatarURLEQ applies the EQ predicate on the "avatar_url" field.
func AvatarURLEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLNEQ applies the NEQ predicate on the "avatar_url" field.
func AvatarURLNEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLIn applies the In predicate on the "avatar_url" field.
func AvatarURLIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldAvatarURL), v...))
	})
}

// AvatarURLNotIn applies the NotIn predicate on the "avatar_url" field.
func AvatarURLNotIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldAvatarURL), v...))
	})
}

// AvatarURLGT applies the GT predicate on the "avatar_url" field.
func AvatarURLGT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLGTE applies the GTE predicate on the "avatar_url" field.
func AvatarURLGTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLLT applies the LT predicate on the "avatar_url" field.
func AvatarURLLT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLLTE applies the LTE predicate on the "avatar_url" field.
func AvatarURLLTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLContains applies the Contains predicate on the "avatar_url" field.
func AvatarURLContains(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLHasPrefix applies the HasPrefix predicate on the "avatar_url" field.
func AvatarURLHasPrefix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLHasSuffix applies the HasSuffix predicate on the "avatar_url" field.
func AvatarURLHasSuffix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLIsNil applies the IsNil predicate on the "avatar_url" field.
func AvatarURLIsNil() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldAvatarURL)))
	})
}

// AvatarURLNotNil applies the NotNil predicate on the "avatar_url" field.
func AvatarURLNotNil() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldAvatarURL)))
	})
}

// AvatarURLEqualFold applies the EqualFold predicate on the "avatar_url" field.
func AvatarURLEqualFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldAvatarURL), v))
	})
}

// AvatarURLContainsFold applies the ContainsFold predicate on the "avatar_url" field.
func AvatarURLContainsFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldAvatarURL), v))
	})
}

// DefaultBranchEQ applies the EQ predicate on the "default_branch" field.
func DefaultBranchEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchNEQ applies the NEQ predicate on the "default_branch" field.
func DefaultBranchNEQ(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchIn applies the In predicate on the "default_branch" field.
func DefaultBranchIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldDefaultBranch), v...))
	})
}

// DefaultBranchNotIn applies the NotIn predicate on the "default_branch" field.
func DefaultBranchNotIn(vs ...string) predicate.Repository {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Repository(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldDefaultBranch), v...))
	})
}

// DefaultBranchGT applies the GT predicate on the "default_branch" field.
func DefaultBranchGT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchGTE applies the GTE predicate on the "default_branch" field.
func DefaultBranchGTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchLT applies the LT predicate on the "default_branch" field.
func DefaultBranchLT(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchLTE applies the LTE predicate on the "default_branch" field.
func DefaultBranchLTE(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchContains applies the Contains predicate on the "default_branch" field.
func DefaultBranchContains(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchHasPrefix applies the HasPrefix predicate on the "default_branch" field.
func DefaultBranchHasPrefix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchHasSuffix applies the HasSuffix predicate on the "default_branch" field.
func DefaultBranchHasSuffix(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchIsNil applies the IsNil predicate on the "default_branch" field.
func DefaultBranchIsNil() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldDefaultBranch)))
	})
}

// DefaultBranchNotNil applies the NotNil predicate on the "default_branch" field.
func DefaultBranchNotNil() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldDefaultBranch)))
	})
}

// DefaultBranchEqualFold applies the EqualFold predicate on the "default_branch" field.
func DefaultBranchEqualFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldDefaultBranch), v))
	})
}

// DefaultBranchContainsFold applies the ContainsFold predicate on the "default_branch" field.
func DefaultBranchContainsFold(v string) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldDefaultBranch), v))
	})
}

// HasScan applies the HasEdge predicate on the "scan" edge.
func HasScan() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ScanTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ScanTable, ScanPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasScanWith applies the HasEdge predicate on the "scan" edge with a given conditions (other predicates).
func HasScanWith(preds ...predicate.Scan) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ScanInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ScanTable, ScanPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasMain applies the HasEdge predicate on the "main" edge.
func HasMain() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(MainTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MainTable, MainColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMainWith applies the HasEdge predicate on the "main" edge with a given conditions (other predicates).
func HasMainWith(preds ...predicate.Scan) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(MainInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MainTable, MainColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLatest applies the HasEdge predicate on the "latest" edge.
func HasLatest() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LatestTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LatestTable, LatestColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLatestWith applies the HasEdge predicate on the "latest" edge with a given conditions (other predicates).
func HasLatestWith(preds ...predicate.Scan) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LatestInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, LatestTable, LatestColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasReport applies the HasEdge predicate on the "report" edge.
func HasReport() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReportTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ReportTable, ReportPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasReportWith applies the HasEdge predicate on the "report" edge with a given conditions (other predicates).
func HasReportWith(preds ...predicate.Report) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ReportInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ReportTable, ReportPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasLatestReport applies the HasEdge predicate on the "latest_report" edge.
func HasLatestReport() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LatestReportTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, LatestReportTable, LatestReportColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLatestReportWith applies the HasEdge predicate on the "latest_report" edge with a given conditions (other predicates).
func HasLatestReportWith(preds ...predicate.Report) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(LatestReportInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, LatestReportTable, LatestReportColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasStatus applies the HasEdge predicate on the "status" edge.
func HasStatus() predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(StatusTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, StatusTable, StatusColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStatusWith applies the HasEdge predicate on the "status" edge with a given conditions (other predicates).
func HasStatusWith(preds ...predicate.VulnStatusIndex) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(StatusInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, StatusTable, StatusColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Repository) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Repository) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
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
func Not(p predicate.Repository) predicate.Repository {
	return predicate.Repository(func(s *sql.Selector) {
		p(s.Not())
	})
}
