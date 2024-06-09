package config

import (
	"context"
	"log/slog"

	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/infra/cs"
	"github.com/urfave/cli/v2"
)

type CloudStorage struct {
	bucket string
	prefix string
}

func (x *CloudStorage) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "cloud-storage-bucket",
			Usage:       "Cloud Storage bucket name",
			Category:    "Cloud Storage",
			Destination: &x.bucket,
			EnvVars:     []string{"OCTOVY_CLOUD_STORAGE_BUCKET"},
		},
		&cli.StringFlag{
			Name:        "cloud-storage-prefix",
			Usage:       "Cloud Storage prefix",
			Category:    "Cloud Storage",
			Destination: &x.prefix,
			EnvVars:     []string{"OCTOVY_CLOUD_STORAGE_PREFIX"},
		},
	}
}

func (x *CloudStorage) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("Bucket", x.bucket),
		slog.Any("Prefix", x.prefix),
	)
}

func (x *CloudStorage) NewClient(ctx context.Context) (interfaces.Storage, error) {
	if x.bucket == "" {
		return nil, nil
	}

	var options []cs.Option
	if x.prefix != "" {
		options = append(options, cs.WithPrefix(x.prefix))
	}

	return cs.New(ctx, x.bucket, options...)
}
