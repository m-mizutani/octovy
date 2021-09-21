// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/m-mizutani/octovy/pkg/infra/ent/migrate"

	"github.com/m-mizutani/octovy/pkg/infra/ent/packagerecord"
	"github.com/m-mizutani/octovy/pkg/infra/ent/repository"
	"github.com/m-mizutani/octovy/pkg/infra/ent/scan"
	"github.com/m-mizutani/octovy/pkg/infra/ent/session"
	"github.com/m-mizutani/octovy/pkg/infra/ent/user"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnerability"
	"github.com/m-mizutani/octovy/pkg/infra/ent/vulnstatus"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// PackageRecord is the client for interacting with the PackageRecord builders.
	PackageRecord *PackageRecordClient
	// Repository is the client for interacting with the Repository builders.
	Repository *RepositoryClient
	// Scan is the client for interacting with the Scan builders.
	Scan *ScanClient
	// Session is the client for interacting with the Session builders.
	Session *SessionClient
	// User is the client for interacting with the User builders.
	User *UserClient
	// VulnStatus is the client for interacting with the VulnStatus builders.
	VulnStatus *VulnStatusClient
	// Vulnerability is the client for interacting with the Vulnerability builders.
	Vulnerability *VulnerabilityClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.PackageRecord = NewPackageRecordClient(c.config)
	c.Repository = NewRepositoryClient(c.config)
	c.Scan = NewScanClient(c.config)
	c.Session = NewSessionClient(c.config)
	c.User = NewUserClient(c.config)
	c.VulnStatus = NewVulnStatusClient(c.config)
	c.Vulnerability = NewVulnerabilityClient(c.config)
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:           ctx,
		config:        cfg,
		PackageRecord: NewPackageRecordClient(cfg),
		Repository:    NewRepositoryClient(cfg),
		Scan:          NewScanClient(cfg),
		Session:       NewSessionClient(cfg),
		User:          NewUserClient(cfg),
		VulnStatus:    NewVulnStatusClient(cfg),
		Vulnerability: NewVulnerabilityClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		config:        cfg,
		PackageRecord: NewPackageRecordClient(cfg),
		Repository:    NewRepositoryClient(cfg),
		Scan:          NewScanClient(cfg),
		Session:       NewSessionClient(cfg),
		User:          NewUserClient(cfg),
		VulnStatus:    NewVulnStatusClient(cfg),
		Vulnerability: NewVulnerabilityClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		PackageRecord.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.PackageRecord.Use(hooks...)
	c.Repository.Use(hooks...)
	c.Scan.Use(hooks...)
	c.Session.Use(hooks...)
	c.User.Use(hooks...)
	c.VulnStatus.Use(hooks...)
	c.Vulnerability.Use(hooks...)
}

// PackageRecordClient is a client for the PackageRecord schema.
type PackageRecordClient struct {
	config
}

// NewPackageRecordClient returns a client for the PackageRecord from the given config.
func NewPackageRecordClient(c config) *PackageRecordClient {
	return &PackageRecordClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `packagerecord.Hooks(f(g(h())))`.
func (c *PackageRecordClient) Use(hooks ...Hook) {
	c.hooks.PackageRecord = append(c.hooks.PackageRecord, hooks...)
}

// Create returns a create builder for PackageRecord.
func (c *PackageRecordClient) Create() *PackageRecordCreate {
	mutation := newPackageRecordMutation(c.config, OpCreate)
	return &PackageRecordCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of PackageRecord entities.
func (c *PackageRecordClient) CreateBulk(builders ...*PackageRecordCreate) *PackageRecordCreateBulk {
	return &PackageRecordCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for PackageRecord.
func (c *PackageRecordClient) Update() *PackageRecordUpdate {
	mutation := newPackageRecordMutation(c.config, OpUpdate)
	return &PackageRecordUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PackageRecordClient) UpdateOne(pr *PackageRecord) *PackageRecordUpdateOne {
	mutation := newPackageRecordMutation(c.config, OpUpdateOne, withPackageRecord(pr))
	return &PackageRecordUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *PackageRecordClient) UpdateOneID(id int) *PackageRecordUpdateOne {
	mutation := newPackageRecordMutation(c.config, OpUpdateOne, withPackageRecordID(id))
	return &PackageRecordUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for PackageRecord.
func (c *PackageRecordClient) Delete() *PackageRecordDelete {
	mutation := newPackageRecordMutation(c.config, OpDelete)
	return &PackageRecordDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *PackageRecordClient) DeleteOne(pr *PackageRecord) *PackageRecordDeleteOne {
	return c.DeleteOneID(pr.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *PackageRecordClient) DeleteOneID(id int) *PackageRecordDeleteOne {
	builder := c.Delete().Where(packagerecord.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PackageRecordDeleteOne{builder}
}

// Query returns a query builder for PackageRecord.
func (c *PackageRecordClient) Query() *PackageRecordQuery {
	return &PackageRecordQuery{
		config: c.config,
	}
}

// Get returns a PackageRecord entity by its id.
func (c *PackageRecordClient) Get(ctx context.Context, id int) (*PackageRecord, error) {
	return c.Query().Where(packagerecord.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PackageRecordClient) GetX(ctx context.Context, id int) *PackageRecord {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryScan queries the scan edge of a PackageRecord.
func (c *PackageRecordClient) QueryScan(pr *PackageRecord) *ScanQuery {
	query := &ScanQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(packagerecord.Table, packagerecord.FieldID, id),
			sqlgraph.To(scan.Table, scan.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, packagerecord.ScanTable, packagerecord.ScanPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryVulnerabilities queries the vulnerabilities edge of a PackageRecord.
func (c *PackageRecordClient) QueryVulnerabilities(pr *PackageRecord) *VulnerabilityQuery {
	query := &VulnerabilityQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(packagerecord.Table, packagerecord.FieldID, id),
			sqlgraph.To(vulnerability.Table, vulnerability.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, packagerecord.VulnerabilitiesTable, packagerecord.VulnerabilitiesPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *PackageRecordClient) Hooks() []Hook {
	return c.hooks.PackageRecord
}

// RepositoryClient is a client for the Repository schema.
type RepositoryClient struct {
	config
}

// NewRepositoryClient returns a client for the Repository from the given config.
func NewRepositoryClient(c config) *RepositoryClient {
	return &RepositoryClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `repository.Hooks(f(g(h())))`.
func (c *RepositoryClient) Use(hooks ...Hook) {
	c.hooks.Repository = append(c.hooks.Repository, hooks...)
}

// Create returns a create builder for Repository.
func (c *RepositoryClient) Create() *RepositoryCreate {
	mutation := newRepositoryMutation(c.config, OpCreate)
	return &RepositoryCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Repository entities.
func (c *RepositoryClient) CreateBulk(builders ...*RepositoryCreate) *RepositoryCreateBulk {
	return &RepositoryCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Repository.
func (c *RepositoryClient) Update() *RepositoryUpdate {
	mutation := newRepositoryMutation(c.config, OpUpdate)
	return &RepositoryUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *RepositoryClient) UpdateOne(r *Repository) *RepositoryUpdateOne {
	mutation := newRepositoryMutation(c.config, OpUpdateOne, withRepository(r))
	return &RepositoryUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *RepositoryClient) UpdateOneID(id int) *RepositoryUpdateOne {
	mutation := newRepositoryMutation(c.config, OpUpdateOne, withRepositoryID(id))
	return &RepositoryUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Repository.
func (c *RepositoryClient) Delete() *RepositoryDelete {
	mutation := newRepositoryMutation(c.config, OpDelete)
	return &RepositoryDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *RepositoryClient) DeleteOne(r *Repository) *RepositoryDeleteOne {
	return c.DeleteOneID(r.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *RepositoryClient) DeleteOneID(id int) *RepositoryDeleteOne {
	builder := c.Delete().Where(repository.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &RepositoryDeleteOne{builder}
}

// Query returns a query builder for Repository.
func (c *RepositoryClient) Query() *RepositoryQuery {
	return &RepositoryQuery{
		config: c.config,
	}
}

// Get returns a Repository entity by its id.
func (c *RepositoryClient) Get(ctx context.Context, id int) (*Repository, error) {
	return c.Query().Where(repository.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *RepositoryClient) GetX(ctx context.Context, id int) *Repository {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryScan queries the scan edge of a Repository.
func (c *RepositoryClient) QueryScan(r *Repository) *ScanQuery {
	query := &ScanQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(repository.Table, repository.FieldID, id),
			sqlgraph.To(scan.Table, scan.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, repository.ScanTable, repository.ScanPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryStatus queries the status edge of a Repository.
func (c *RepositoryClient) QueryStatus(r *Repository) *VulnStatusQuery {
	query := &VulnStatusQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(repository.Table, repository.FieldID, id),
			sqlgraph.To(vulnstatus.Table, vulnstatus.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, repository.StatusTable, repository.StatusColumn),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *RepositoryClient) Hooks() []Hook {
	return c.hooks.Repository
}

// ScanClient is a client for the Scan schema.
type ScanClient struct {
	config
}

// NewScanClient returns a client for the Scan from the given config.
func NewScanClient(c config) *ScanClient {
	return &ScanClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `scan.Hooks(f(g(h())))`.
func (c *ScanClient) Use(hooks ...Hook) {
	c.hooks.Scan = append(c.hooks.Scan, hooks...)
}

// Create returns a create builder for Scan.
func (c *ScanClient) Create() *ScanCreate {
	mutation := newScanMutation(c.config, OpCreate)
	return &ScanCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Scan entities.
func (c *ScanClient) CreateBulk(builders ...*ScanCreate) *ScanCreateBulk {
	return &ScanCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Scan.
func (c *ScanClient) Update() *ScanUpdate {
	mutation := newScanMutation(c.config, OpUpdate)
	return &ScanUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ScanClient) UpdateOne(s *Scan) *ScanUpdateOne {
	mutation := newScanMutation(c.config, OpUpdateOne, withScan(s))
	return &ScanUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ScanClient) UpdateOneID(id string) *ScanUpdateOne {
	mutation := newScanMutation(c.config, OpUpdateOne, withScanID(id))
	return &ScanUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Scan.
func (c *ScanClient) Delete() *ScanDelete {
	mutation := newScanMutation(c.config, OpDelete)
	return &ScanDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *ScanClient) DeleteOne(s *Scan) *ScanDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *ScanClient) DeleteOneID(id string) *ScanDeleteOne {
	builder := c.Delete().Where(scan.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ScanDeleteOne{builder}
}

// Query returns a query builder for Scan.
func (c *ScanClient) Query() *ScanQuery {
	return &ScanQuery{
		config: c.config,
	}
}

// Get returns a Scan entity by its id.
func (c *ScanClient) Get(ctx context.Context, id string) (*Scan, error) {
	return c.Query().Where(scan.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ScanClient) GetX(ctx context.Context, id string) *Scan {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryRepository queries the repository edge of a Scan.
func (c *ScanClient) QueryRepository(s *Scan) *RepositoryQuery {
	query := &RepositoryQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(scan.Table, scan.FieldID, id),
			sqlgraph.To(repository.Table, repository.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, scan.RepositoryTable, scan.RepositoryPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPackages queries the packages edge of a Scan.
func (c *ScanClient) QueryPackages(s *Scan) *PackageRecordQuery {
	query := &PackageRecordQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(scan.Table, scan.FieldID, id),
			sqlgraph.To(packagerecord.Table, packagerecord.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, scan.PackagesTable, scan.PackagesPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ScanClient) Hooks() []Hook {
	return c.hooks.Scan
}

// SessionClient is a client for the Session schema.
type SessionClient struct {
	config
}

// NewSessionClient returns a client for the Session from the given config.
func NewSessionClient(c config) *SessionClient {
	return &SessionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `session.Hooks(f(g(h())))`.
func (c *SessionClient) Use(hooks ...Hook) {
	c.hooks.Session = append(c.hooks.Session, hooks...)
}

// Create returns a create builder for Session.
func (c *SessionClient) Create() *SessionCreate {
	mutation := newSessionMutation(c.config, OpCreate)
	return &SessionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Session entities.
func (c *SessionClient) CreateBulk(builders ...*SessionCreate) *SessionCreateBulk {
	return &SessionCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Session.
func (c *SessionClient) Update() *SessionUpdate {
	mutation := newSessionMutation(c.config, OpUpdate)
	return &SessionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SessionClient) UpdateOne(s *Session) *SessionUpdateOne {
	mutation := newSessionMutation(c.config, OpUpdateOne, withSession(s))
	return &SessionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SessionClient) UpdateOneID(id int) *SessionUpdateOne {
	mutation := newSessionMutation(c.config, OpUpdateOne, withSessionID(id))
	return &SessionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Session.
func (c *SessionClient) Delete() *SessionDelete {
	mutation := newSessionMutation(c.config, OpDelete)
	return &SessionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *SessionClient) DeleteOne(s *Session) *SessionDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *SessionClient) DeleteOneID(id int) *SessionDeleteOne {
	builder := c.Delete().Where(session.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SessionDeleteOne{builder}
}

// Query returns a query builder for Session.
func (c *SessionClient) Query() *SessionQuery {
	return &SessionQuery{
		config: c.config,
	}
}

// Get returns a Session entity by its id.
func (c *SessionClient) Get(ctx context.Context, id int) (*Session, error) {
	return c.Query().Where(session.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SessionClient) GetX(ctx context.Context, id int) *Session {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryLogin queries the login edge of a Session.
func (c *SessionClient) QueryLogin(s *Session) *UserQuery {
	query := &UserQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(session.Table, session.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, session.LoginTable, session.LoginColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SessionClient) Hooks() []Hook {
	return c.hooks.Session
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...Hook) {
	c.hooks.User = append(c.hooks.User, hooks...)
}

// Create returns a create builder for User.
func (c *UserClient) Create() *UserCreate {
	mutation := newUserMutation(c.config, OpCreate)
	return &UserCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of User entities.
func (c *UserClient) CreateBulk(builders ...*UserCreate) *UserCreateBulk {
	return &UserCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	mutation := newUserMutation(c.config, OpUpdate)
	return &UserUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUser(u))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id string) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUserID(id))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	mutation := newUserMutation(c.config, OpDelete)
	return &UserDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *UserClient) DeleteOneID(id string) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Query returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{
		config: c.config,
	}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id string) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id string) *User {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryEditedStatus queries the edited_status edge of a User.
func (c *UserClient) QueryEditedStatus(u *User) *VulnStatusQuery {
	query := &VulnStatusQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := u.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(user.Table, user.FieldID, id),
			sqlgraph.To(vulnstatus.Table, vulnstatus.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, user.EditedStatusTable, user.EditedStatusColumn),
		)
		fromV = sqlgraph.Neighbors(u.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *UserClient) Hooks() []Hook {
	return c.hooks.User
}

// VulnStatusClient is a client for the VulnStatus schema.
type VulnStatusClient struct {
	config
}

// NewVulnStatusClient returns a client for the VulnStatus from the given config.
func NewVulnStatusClient(c config) *VulnStatusClient {
	return &VulnStatusClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `vulnstatus.Hooks(f(g(h())))`.
func (c *VulnStatusClient) Use(hooks ...Hook) {
	c.hooks.VulnStatus = append(c.hooks.VulnStatus, hooks...)
}

// Create returns a create builder for VulnStatus.
func (c *VulnStatusClient) Create() *VulnStatusCreate {
	mutation := newVulnStatusMutation(c.config, OpCreate)
	return &VulnStatusCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of VulnStatus entities.
func (c *VulnStatusClient) CreateBulk(builders ...*VulnStatusCreate) *VulnStatusCreateBulk {
	return &VulnStatusCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for VulnStatus.
func (c *VulnStatusClient) Update() *VulnStatusUpdate {
	mutation := newVulnStatusMutation(c.config, OpUpdate)
	return &VulnStatusUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *VulnStatusClient) UpdateOne(vs *VulnStatus) *VulnStatusUpdateOne {
	mutation := newVulnStatusMutation(c.config, OpUpdateOne, withVulnStatus(vs))
	return &VulnStatusUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *VulnStatusClient) UpdateOneID(id string) *VulnStatusUpdateOne {
	mutation := newVulnStatusMutation(c.config, OpUpdateOne, withVulnStatusID(id))
	return &VulnStatusUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for VulnStatus.
func (c *VulnStatusClient) Delete() *VulnStatusDelete {
	mutation := newVulnStatusMutation(c.config, OpDelete)
	return &VulnStatusDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *VulnStatusClient) DeleteOne(vs *VulnStatus) *VulnStatusDeleteOne {
	return c.DeleteOneID(vs.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *VulnStatusClient) DeleteOneID(id string) *VulnStatusDeleteOne {
	builder := c.Delete().Where(vulnstatus.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &VulnStatusDeleteOne{builder}
}

// Query returns a query builder for VulnStatus.
func (c *VulnStatusClient) Query() *VulnStatusQuery {
	return &VulnStatusQuery{
		config: c.config,
	}
}

// Get returns a VulnStatus entity by its id.
func (c *VulnStatusClient) Get(ctx context.Context, id string) (*VulnStatus, error) {
	return c.Query().Where(vulnstatus.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *VulnStatusClient) GetX(ctx context.Context, id string) *VulnStatus {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *VulnStatusClient) Hooks() []Hook {
	return c.hooks.VulnStatus
}

// VulnerabilityClient is a client for the Vulnerability schema.
type VulnerabilityClient struct {
	config
}

// NewVulnerabilityClient returns a client for the Vulnerability from the given config.
func NewVulnerabilityClient(c config) *VulnerabilityClient {
	return &VulnerabilityClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `vulnerability.Hooks(f(g(h())))`.
func (c *VulnerabilityClient) Use(hooks ...Hook) {
	c.hooks.Vulnerability = append(c.hooks.Vulnerability, hooks...)
}

// Create returns a create builder for Vulnerability.
func (c *VulnerabilityClient) Create() *VulnerabilityCreate {
	mutation := newVulnerabilityMutation(c.config, OpCreate)
	return &VulnerabilityCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Vulnerability entities.
func (c *VulnerabilityClient) CreateBulk(builders ...*VulnerabilityCreate) *VulnerabilityCreateBulk {
	return &VulnerabilityCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Vulnerability.
func (c *VulnerabilityClient) Update() *VulnerabilityUpdate {
	mutation := newVulnerabilityMutation(c.config, OpUpdate)
	return &VulnerabilityUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *VulnerabilityClient) UpdateOne(v *Vulnerability) *VulnerabilityUpdateOne {
	mutation := newVulnerabilityMutation(c.config, OpUpdateOne, withVulnerability(v))
	return &VulnerabilityUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *VulnerabilityClient) UpdateOneID(id string) *VulnerabilityUpdateOne {
	mutation := newVulnerabilityMutation(c.config, OpUpdateOne, withVulnerabilityID(id))
	return &VulnerabilityUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Vulnerability.
func (c *VulnerabilityClient) Delete() *VulnerabilityDelete {
	mutation := newVulnerabilityMutation(c.config, OpDelete)
	return &VulnerabilityDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *VulnerabilityClient) DeleteOne(v *Vulnerability) *VulnerabilityDeleteOne {
	return c.DeleteOneID(v.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *VulnerabilityClient) DeleteOneID(id string) *VulnerabilityDeleteOne {
	builder := c.Delete().Where(vulnerability.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &VulnerabilityDeleteOne{builder}
}

// Query returns a query builder for Vulnerability.
func (c *VulnerabilityClient) Query() *VulnerabilityQuery {
	return &VulnerabilityQuery{
		config: c.config,
	}
}

// Get returns a Vulnerability entity by its id.
func (c *VulnerabilityClient) Get(ctx context.Context, id string) (*Vulnerability, error) {
	return c.Query().Where(vulnerability.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *VulnerabilityClient) GetX(ctx context.Context, id string) *Vulnerability {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryPackages queries the packages edge of a Vulnerability.
func (c *VulnerabilityClient) QueryPackages(v *Vulnerability) *PackageRecordQuery {
	query := &PackageRecordQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := v.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(vulnerability.Table, vulnerability.FieldID, id),
			sqlgraph.To(packagerecord.Table, packagerecord.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, vulnerability.PackagesTable, vulnerability.PackagesPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(v.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryStatus queries the status edge of a Vulnerability.
func (c *VulnerabilityClient) QueryStatus(v *Vulnerability) *VulnStatusQuery {
	query := &VulnStatusQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := v.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(vulnerability.Table, vulnerability.FieldID, id),
			sqlgraph.To(vulnstatus.Table, vulnstatus.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, vulnerability.StatusTable, vulnerability.StatusColumn),
		)
		fromV = sqlgraph.Neighbors(v.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *VulnerabilityClient) Hooks() []Hook {
	return c.hooks.Vulnerability
}
