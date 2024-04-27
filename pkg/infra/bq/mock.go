package bq

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

type Mock struct {
	FnCreateTable func(ctx context.Context, table types.BQTableID, md *bigquery.TableMetadata) error
	FnGetMetadata func(ctx context.Context, table types.BQTableID) (*bigquery.TableMetadata, error)
	FnInsert      func(ctx context.Context, tableID types.BQTableID, schema bigquery.Schema, data any) error
	FnUpdateTable func(ctx context.Context, table types.BQTableID, md bigquery.TableMetadataToUpdate, eTag string) error
}

// CreateTable implements interfaces.BigQuery.
func (m *Mock) CreateTable(ctx context.Context, table types.BQTableID, md *bigquery.TableMetadata) error {
	return m.FnCreateTable(ctx, table, md)
}

// GetMetadata implements interfaces.BigQuery.
func (m *Mock) GetMetadata(ctx context.Context, table types.BQTableID) (*bigquery.TableMetadata, error) {
	return m.FnGetMetadata(ctx, table)
}

// Insert implements interfaces.BigQuery.
func (m *Mock) Insert(ctx context.Context, tableID types.BQTableID, schema bigquery.Schema, data any) error {
	return m.FnInsert(ctx, tableID, schema, data)
}

// UpdateTable implements interfaces.BigQuery.
func (m *Mock) UpdateTable(ctx context.Context, table types.BQTableID, md bigquery.TableMetadataToUpdate, eTag string) error {
	return m.FnUpdateTable(ctx, table, md, eTag)
}

var _ interfaces.BigQuery = &Mock{}
