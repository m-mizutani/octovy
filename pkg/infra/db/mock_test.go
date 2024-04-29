package db_test

import (
	"testing"

	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/infra/db"
)

func TestMock(t *testing.T) {
	mock := db.NewMock()
	interfaces.FirestoreClientTest(t, mock)
}
