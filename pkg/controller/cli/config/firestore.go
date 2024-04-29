package config

import (
	"context"
	"log/slog"

	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"github.com/m-mizutani/octovy/pkg/infra/db"
	"github.com/urfave/cli/v2"
)

type Firestore struct {
	projectID  types.GoogleProjectID
	databaseID types.FSDatabaseID
}

func (x *Firestore) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "firestore-project-id",
			Usage:       "Firestore project ID",
			Category:    "Firestore",
			Destination: (*string)(&x.projectID),
			EnvVars:     []string{"OCTOVY_FIRESTORE_PROJECT_ID"},
		},
		&cli.StringFlag{
			Name:        "firestore-database-id",
			Usage:       "Firestore database ID",
			Category:    "Firestore",
			Destination: (*string)(&x.databaseID),
			EnvVars:     []string{"OCTOVY_FIRESTORE_DATABASE_ID"},
		},
	}
}

func (x *Firestore) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("ProjectID", x.projectID),
		slog.Any("DatabaseID", x.databaseID),
	)
}

func (x *Firestore) ProjectID() types.GoogleProjectID {
	return x.projectID
}

func (x *Firestore) DatabaseID() types.FSDatabaseID {
	return x.databaseID
}

func (x *Firestore) NewClient(ctx context.Context) (interfaces.Firestore, error) {
	if x.projectID == "" && x.databaseID == "" {
		return nil, nil
	}
	return db.New(ctx, x.projectID, x.databaseID)
}
