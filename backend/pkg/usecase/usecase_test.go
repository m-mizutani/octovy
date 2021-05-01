package usecase_test

import (
	"testing"

	"github.com/m-mizutani/octovy/backend/pkg/infra"
	"github.com/m-mizutani/octovy/backend/pkg/infra/db"
)

func newTestTable(t *testing.T) infra.DBClient {
	tableName := "dynamodb-test"

	client, err := db.NewDynamoClientLocal("ap-northeast-1", tableName)
	if err != nil {
		panic("Failed to use local DynamoDB: " + err.Error())
	}

	t.Log("Created table name: ", client.(*db.DynamoClient).TableName())

	t.Cleanup(func() {
		if t.Failed() {
			return // Failed test table is not deleted
		}

		if err := client.Close(); err != nil {
			panic("Failed to delete test table: " + err.Error())
		}
	})
	return client
}
