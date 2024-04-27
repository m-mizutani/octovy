package bq_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/m-mizutani/bqs"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/model"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/bq"
	"github.com/m-mizutani/octovy/pkg/utils"
)

func TestClient(t *testing.T) {
	projectID := utils.LoadEnv(t, "TEST_BIGQUERY_PROJECT_ID")
	datasetID := utils.LoadEnv(t, "TEST_BIGQUERY_DATASET_ID")

	ctx := context.Background()

	tblName := types.BQTableID(time.Now().Format("insert_test_20060102_150405"))
	client, err := bq.New(ctx, types.GoogleProjectID(projectID), types.BQDatasetID(datasetID))
	gt.NoError(t, err)

	var baseSchema bigquery.Schema

	t.Run("Create base table at first", func(t *testing.T) {
		var scan model.Scan
		baseSchema = gt.R1(bqs.Infer(scan)).NoError(t)
		gt.NoError(t, err)

		gt.NoError(t, client.CreateTable(ctx, tblName, &bigquery.TableMetadata{
			Name:   tblName.String(),
			Schema: baseSchema,
		}))
	})

	t.Run("Insert record", func(t *testing.T) {
		var scan model.Scan
		utils.LoadJson(t, "testdata/data.json", &scan.Report)
		dataSchema := gt.R1(bqs.Infer(scan)).NoError(t)
		mergedSchema := gt.R1(bqs.Merge(baseSchema, dataSchema)).NoError(t)

		md := gt.R1(client.GetMetadata(ctx, tblName)).NoError(t)
		gt.False(t, bqs.Equal(mergedSchema, baseSchema))
		gt.NoError(t, client.UpdateTable(ctx, tblName, bigquery.TableMetadataToUpdate{
			Schema: mergedSchema,
		}, md.ETag)).Must()

		record := model.ScanRawRecord{
			Scan:      scan,
			Timestamp: scan.Timestamp.UnixMicro(),
		}
		gt.NoError(t, client.Insert(ctx, tblName, mergedSchema, record))
	})
}
