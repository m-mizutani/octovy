package db

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	client *firestore.Client
}

var _ interfaces.Firestore = (*Client)(nil)

func New(ctx context.Context, projectID types.GoogleProjectID, dbID types.FSDatabaseID) (*Client, error) {
	client, err := firestore.NewClientWithDatabase(ctx, string(projectID), string(dbID))
	if err != nil {
		return nil, goerr.Wrap(err, "failed to create Firestore client").With("projectID", projectID).With("dbID", dbID).With("err", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) toDocRef(docRefs []types.FireStoreRef) *firestore.DocumentRef {
	var ref *firestore.DocumentRef
	for _, docRef := range docRefs {
		if ref == nil {
			ref = c.client.Collection(string(docRef.CollectionID)).Doc(string(docRef.DocumentID))
		} else {
			ref = ref.Collection(string(docRef.CollectionID)).Doc(string(docRef.DocumentID))
		}
	}

	return ref
}

// Get implements interfaces.Database.
func (c *Client) Get(ctx context.Context, value any, docRefs ...types.FireStoreRef) error {
	ref := c.toDocRef(docRefs)

	doc, err := ref.Get(ctx)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				return nil
			}
		}
		return goerr.Wrap(err, "failed to get document").With("ref", docRefs).With("err", err)
	}

	if err := doc.DataTo(value); err != nil {
		return goerr.Wrap(err, "failed to convert document to value").With("docRefs", docRefs).With("err", err)
	}

	return nil
}

// Put implements interfaces.Database.
func (c *Client) Put(ctx context.Context, value any, docRefs ...types.FireStoreRef) error {
	ref := c.toDocRef(docRefs)

	if _, err := ref.Set(ctx, value); err != nil {
		return goerr.Wrap(err, "failed to set document").With("docRefs", docRefs).With("err", err)
	}

	return nil
}
