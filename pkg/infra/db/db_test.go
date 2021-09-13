package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/infra/db"
)

func setupDB(t *testing.T) db.Interface {
	client := db.NewDBMock(t)
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning failed to close DB: %+v", err)
		}
	})

	return client
}
