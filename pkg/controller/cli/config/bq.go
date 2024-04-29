package config

import (
	"context"
	"log/slog"

	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/bq"
	"github.com/urfave/cli/v2"
)

type BigQuery struct {
	projectID types.GoogleProjectID
	datasetID types.BQDatasetID
	tableID   types.BQTableID
}

func (x *BigQuery) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "bigquery-project-id",
			Usage:       "BigQuery project ID",
			Category:    "BigQuery",
			Destination: (*string)(&x.projectID),
			EnvVars:     []string{"OCTOVY_BIGQUERY_PROJECT_ID"},
		},
		&cli.StringFlag{
			Name:        "bigquery-dataset-id",
			Usage:       "BigQuery dataset ID",
			Category:    "BigQuery",
			Destination: (*string)(&x.datasetID),
			EnvVars:     []string{"OCTOVY_BIGQUERY_DATASET_ID"},
		},
		&cli.StringFlag{
			Name:        "bigquery-table-id",
			Usage:       "BigQuery table ID",
			Category:    "BigQuery",
			Destination: (*string)(&x.tableID),
			EnvVars:     []string{"OCTOVY_BIGQUERY_TABLE_ID"},
			Value:       "scans",
		},
	}
}

func (x *BigQuery) TableID() types.BQTableID {
	return x.tableID
}

func (x *BigQuery) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("ProjectID", x.projectID),
		slog.Any("DatasetID", x.datasetID),
	)
}

func (x *BigQuery) NewClient(ctx context.Context) (interfaces.BigQuery, error) {
	if x.projectID == "" && x.datasetID == "" {
		return nil, nil
	}
	return bq.New(ctx, x.projectID, x.datasetID)
}
