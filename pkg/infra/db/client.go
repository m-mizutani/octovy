package db

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

type Client struct {
	client *firestore.Client
}

var _ interfaces.Database = (*Client)(nil)

func New(ctx context.Context, projectID types.GoogleProjectID, dbID types.FSDatabaseID) (*Client, error) {
	client, err := firestore.NewClient(ctx, string(projectID))
	if err != nil {
		return nil, goerr.Wrap(err, "failed to create Firestore client").With("projectID", projectID).With("dbID", dbID).With("err", err)
	}

	return &Client{
		client: client,
	}, nil
}

// Get implements interfaces.Database.
func (c *Client) Get(ctx context.Context, colID types.FSCollectionID, docID types.FSDocumentID, value any) error {
	ref := c.client.Collection(string(colID)).Doc(string(docID))

	doc, err := ref.Get(ctx)
	if err != nil {
		return goerr.Wrap(err, "failed to get document").With("colID", colID).With("docID", docID).With("err", err)
	}

	if err := doc.DataTo(value); err != nil {
		return goerr.Wrap(err, "failed to convert document to value").With("colID", colID).With("docID", docID).With("err", err)
	}

	return nil
}

// Put implements interfaces.Database.
func (c *Client) Put(ctx context.Context, colID types.FSCollectionID, docID types.FSDocumentID, value any) error {
	ref := c.client.Collection(string(colID)).Doc(string(docID))

	if _, err := ref.Set(ctx, value); err != nil {
		return goerr.Wrap(err, "failed to set document").With("colID", colID).With("docID", docID).With("err", err)
	}

	return nil
}
