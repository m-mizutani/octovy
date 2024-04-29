package db_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/m-mizutani/octovy/pkg/utils"
)

func TestFirestore(t *testing.T) {
	projectID := utils.LoadEnv(t, "TEST_FIRESTORE_PROJECT_ID")
	dbID := utils.LoadEnv(t, "TEST_FIRESTORE_DATABASE_ID")

	client, err := db.New(context.Background(), types.GoogleProjectID(projectID), types.FSDatabaseID(dbID))
	gt.NoError(t, err)

	interfaces.FirestoreClientTest(t, client)
}
