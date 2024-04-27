package interfaces

import (
	"context"

	"cloud.google.com/go/bigquery"

	"github.com/m-mizutani/octovy/pkg/domain/types"
)

type BigQuery interface {
	Insert(ctx context.Context, tableID types.BQTableID, schema bigquery.Schema, data any) error

	GetMetadata(ctx context.Context, table types.BQTableID) (*bigquery.TableMetadata, error)
	UpdateTable(ctx context.Context, table types.BQTableID, md bigquery.TableMetadataToUpdate, eTag string) error
	CreateTable(ctx context.Context, table types.BQTableID, md *bigquery.TableMetadata) error
}

type Database interface {
	Put(ctx context.Context, colID types.FSCollectionID, docID types.FSDocumentID, value any) error
	Get(ctx context.Context, colID types.FSCollectionID, docID types.FSDocumentID, value any) error
}
