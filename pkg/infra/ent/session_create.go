// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/octovy/pkg/infra/ent/session"
	"github.com/m-mizutani/octovy/pkg/infra/ent/user"
)

// SessionCreate is the builder for creating a Session entity.
type SessionCreate struct {
	config
	mutation *SessionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetUserID sets the "user_id" field.
func (sc *SessionCreate) SetUserID(i int) *SessionCreate {
	sc.mutation.SetUserID(i)
	return sc
}

// SetToken sets the "token" field.
func (sc *SessionCreate) SetToken(s string) *SessionCreate {
	sc.mutation.SetToken(s)
	return sc
}

// SetCreatedAt sets the "created_at" field.
func (sc *SessionCreate) SetCreatedAt(i int64) *SessionCreate {
	sc.mutation.SetCreatedAt(i)
	return sc
}

// SetExpiresAt sets the "expires_at" field.
func (sc *SessionCreate) SetExpiresAt(i int64) *SessionCreate {
	sc.mutation.SetExpiresAt(i)
	return sc
}

// SetID sets the "id" field.
func (sc *SessionCreate) SetID(s string) *SessionCreate {
	sc.mutation.SetID(s)
	return sc
}

// SetLoginID sets the "login" edge to the User entity by ID.
func (sc *SessionCreate) SetLoginID(id int) *SessionCreate {
	sc.mutation.SetLoginID(id)
	return sc
}

// SetNillableLoginID sets the "login" edge to the User entity by ID if the given value is not nil.
func (sc *SessionCreate) SetNillableLoginID(id *int) *SessionCreate {
	if id != nil {
		sc = sc.SetLoginID(*id)
	}
	return sc
}

// SetLogin sets the "login" edge to the User entity.
func (sc *SessionCreate) SetLogin(u *User) *SessionCreate {
	return sc.SetLoginID(u.ID)
}

// Mutation returns the SessionMutation object of the builder.
func (sc *SessionCreate) Mutation() *SessionMutation {
	return sc.mutation
}

// Save creates the Session in the database.
func (sc *SessionCreate) Save(ctx context.Context) (*Session, error) {
	var (
		err  error
		node *Session
	)
	if len(sc.hooks) == 0 {
		if err = sc.check(); err != nil {
			return nil, err
		}
		node, err = sc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SessionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = sc.check(); err != nil {
				return nil, err
			}
			sc.mutation = mutation
			if node, err = sc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(sc.hooks) - 1; i >= 0; i-- {
			if sc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = sc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SessionCreate) SaveX(ctx context.Context) *Session {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SessionCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SessionCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SessionCreate) check() error {
	if _, ok := sc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "user_id"`)}
	}
	if _, ok := sc.mutation.Token(); !ok {
		return &ValidationError{Name: "token", err: errors.New(`ent: missing required field "token"`)}
	}
	if v, ok := sc.mutation.Token(); ok {
		if err := session.TokenValidator(v); err != nil {
			return &ValidationError{Name: "token", err: fmt.Errorf(`ent: validator failed for field "token": %w`, err)}
		}
	}
	if _, ok := sc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "created_at"`)}
	}
	if _, ok := sc.mutation.ExpiresAt(); !ok {
		return &ValidationError{Name: "expires_at", err: errors.New(`ent: missing required field "expires_at"`)}
	}
	return nil
}

func (sc *SessionCreate) sqlSave(ctx context.Context) (*Session, error) {
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		_node.ID = _spec.ID.Value.(string)
	}
	return _node, nil
}

func (sc *SessionCreate) createSpec() (*Session, *sqlgraph.CreateSpec) {
	var (
		_node = &Session{config: sc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: session.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: session.FieldID,
			},
		}
	)
	_spec.OnConflict = sc.conflict
	if id, ok := sc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := sc.mutation.UserID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: session.FieldUserID,
		})
		_node.UserID = value
	}
	if value, ok := sc.mutation.Token(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: session.FieldToken,
		})
		_node.Token = value
	}
	if value, ok := sc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: session.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := sc.mutation.ExpiresAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt64,
			Value:  value,
			Column: session.FieldExpiresAt,
		})
		_node.ExpiresAt = value
	}
	if nodes := sc.mutation.LoginIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   session.LoginTable,
			Columns: []string{session.LoginColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.session_login = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Session.Create().
//		SetUserID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SessionUpsert) {
//			SetUserID(v+v).
//		}).
//		Exec(ctx)
//
func (sc *SessionCreate) OnConflict(opts ...sql.ConflictOption) *SessionUpsertOne {
	sc.conflict = opts
	return &SessionUpsertOne{
		create: sc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Session.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (sc *SessionCreate) OnConflictColumns(columns ...string) *SessionUpsertOne {
	sc.conflict = append(sc.conflict, sql.ConflictColumns(columns...))
	return &SessionUpsertOne{
		create: sc,
	}
}

type (
	// SessionUpsertOne is the builder for "upsert"-ing
	//  one Session node.
	SessionUpsertOne struct {
		create *SessionCreate
	}

	// SessionUpsert is the "OnConflict" setter.
	SessionUpsert struct {
		*sql.UpdateSet
	}
)

// SetUserID sets the "user_id" field.
func (u *SessionUpsert) SetUserID(v int) *SessionUpsert {
	u.Set(session.FieldUserID, v)
	return u
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *SessionUpsert) UpdateUserID() *SessionUpsert {
	u.SetExcluded(session.FieldUserID)
	return u
}

// SetToken sets the "token" field.
func (u *SessionUpsert) SetToken(v string) *SessionUpsert {
	u.Set(session.FieldToken, v)
	return u
}

// UpdateToken sets the "token" field to the value that was provided on create.
func (u *SessionUpsert) UpdateToken() *SessionUpsert {
	u.SetExcluded(session.FieldToken)
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *SessionUpsert) SetCreatedAt(v int64) *SessionUpsert {
	u.Set(session.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *SessionUpsert) UpdateCreatedAt() *SessionUpsert {
	u.SetExcluded(session.FieldCreatedAt)
	return u
}

// SetExpiresAt sets the "expires_at" field.
func (u *SessionUpsert) SetExpiresAt(v int64) *SessionUpsert {
	u.Set(session.FieldExpiresAt, v)
	return u
}

// UpdateExpiresAt sets the "expires_at" field to the value that was provided on create.
func (u *SessionUpsert) UpdateExpiresAt() *SessionUpsert {
	u.SetExcluded(session.FieldExpiresAt)
	return u
}

// UpdateNewValues updates the fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Session.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(session.FieldID)
//			}),
//		).
//		Exec(ctx)
//
func (u *SessionUpsertOne) UpdateNewValues() *SessionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(session.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//  client.Session.Create().
//      OnConflict(sql.ResolveWithIgnore()).
//      Exec(ctx)
//
func (u *SessionUpsertOne) Ignore() *SessionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SessionUpsertOne) DoNothing() *SessionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SessionCreate.OnConflict
// documentation for more info.
func (u *SessionUpsertOne) Update(set func(*SessionUpsert)) *SessionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SessionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserID sets the "user_id" field.
func (u *SessionUpsertOne) SetUserID(v int) *SessionUpsertOne {
	return u.Update(func(s *SessionUpsert) {
		s.SetUserID(v)
	})
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *SessionUpsertOne) UpdateUserID() *SessionUpsertOne {
	return u.Update(func(s *SessionUpsert) {
		s.UpdateUserID()
	})
}

// SetToken sets the "token" field.
func (u *SessionUpsertOne) SetToken(v string) *SessionUpsertOne {
	return u.Update(func(s *SessionUpsert) {
		s.SetToken(v)
	})
}

// UpdateToken sets the "token" field to the value that was provided on create.
func (u *SessionUpsertOne) UpdateToken() *SessionUpsertOne {
	return u.Update(func(s *SessionUpsert) {
		s.UpdateToken()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *SessionUpsertOne) SetCreatedAt(v int64) *SessionUpsertOne {
	return u.Update(func(s *SessionUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *SessionUpsertOne) UpdateCreatedAt() *SessionUpsertOne {
	return u.Update(func(s *SessionUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetExpiresAt sets the "expires_at" field.
func (u *SessionUpsertOne) SetExpiresAt(v int64) *SessionUpsertOne {
	return u.Update(func(s *SessionUpsert) {
		s.SetExpiresAt(v)
	})
}

// UpdateExpiresAt sets the "expires_at" field to the value that was provided on create.
func (u *SessionUpsertOne) UpdateExpiresAt() *SessionUpsertOne {
	return u.Update(func(s *SessionUpsert) {
		s.UpdateExpiresAt()
	})
}

// Exec executes the query.
func (u *SessionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SessionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SessionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SessionUpsertOne) ID(ctx context.Context) (id string, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: SessionUpsertOne.ID is not supported by MySQL driver. Use SessionUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SessionUpsertOne) IDX(ctx context.Context) string {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SessionCreateBulk is the builder for creating many Session entities in bulk.
type SessionCreateBulk struct {
	config
	builders []*SessionCreate
	conflict []sql.ConflictOption
}

// Save creates the Session entities in the database.
func (scb *SessionCreateBulk) Save(ctx context.Context) ([]*Session, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Session, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SessionMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = scb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SessionCreateBulk) SaveX(ctx context.Context) []*Session {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SessionCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SessionCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Session.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SessionUpsert) {
//			SetUserID(v+v).
//		}).
//		Exec(ctx)
//
func (scb *SessionCreateBulk) OnConflict(opts ...sql.ConflictOption) *SessionUpsertBulk {
	scb.conflict = opts
	return &SessionUpsertBulk{
		create: scb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Session.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
//
func (scb *SessionCreateBulk) OnConflictColumns(columns ...string) *SessionUpsertBulk {
	scb.conflict = append(scb.conflict, sql.ConflictColumns(columns...))
	return &SessionUpsertBulk{
		create: scb,
	}
}

// SessionUpsertBulk is the builder for "upsert"-ing
// a bulk of Session nodes.
type SessionUpsertBulk struct {
	create *SessionCreateBulk
}

// UpdateNewValues updates the fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Session.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(session.FieldID)
//			}),
//		).
//		Exec(ctx)
//
func (u *SessionUpsertBulk) UpdateNewValues() *SessionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(session.FieldID)
				return
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Session.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
//
func (u *SessionUpsertBulk) Ignore() *SessionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SessionUpsertBulk) DoNothing() *SessionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SessionCreateBulk.OnConflict
// documentation for more info.
func (u *SessionUpsertBulk) Update(set func(*SessionUpsert)) *SessionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SessionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUserID sets the "user_id" field.
func (u *SessionUpsertBulk) SetUserID(v int) *SessionUpsertBulk {
	return u.Update(func(s *SessionUpsert) {
		s.SetUserID(v)
	})
}

// UpdateUserID sets the "user_id" field to the value that was provided on create.
func (u *SessionUpsertBulk) UpdateUserID() *SessionUpsertBulk {
	return u.Update(func(s *SessionUpsert) {
		s.UpdateUserID()
	})
}

// SetToken sets the "token" field.
func (u *SessionUpsertBulk) SetToken(v string) *SessionUpsertBulk {
	return u.Update(func(s *SessionUpsert) {
		s.SetToken(v)
	})
}

// UpdateToken sets the "token" field to the value that was provided on create.
func (u *SessionUpsertBulk) UpdateToken() *SessionUpsertBulk {
	return u.Update(func(s *SessionUpsert) {
		s.UpdateToken()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *SessionUpsertBulk) SetCreatedAt(v int64) *SessionUpsertBulk {
	return u.Update(func(s *SessionUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *SessionUpsertBulk) UpdateCreatedAt() *SessionUpsertBulk {
	return u.Update(func(s *SessionUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetExpiresAt sets the "expires_at" field.
func (u *SessionUpsertBulk) SetExpiresAt(v int64) *SessionUpsertBulk {
	return u.Update(func(s *SessionUpsert) {
		s.SetExpiresAt(v)
	})
}

// UpdateExpiresAt sets the "expires_at" field to the value that was provided on create.
func (u *SessionUpsertBulk) UpdateExpiresAt() *SessionUpsertBulk {
	return u.Update(func(s *SessionUpsert) {
		s.UpdateExpiresAt()
	})
}

// Exec executes the query.
func (u *SessionUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the SessionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SessionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SessionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}