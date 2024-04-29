package db

import (
	"context"
	"encoding/json"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/domain/interfaces"
	"github.com/m-mizutani/octovy/pkg/domain/types"
)

type MockCollection struct {
	Docs map[types.FSDocumentID]*MockDoc
}

type MockDoc struct {
	Value       any
	Collections map[types.FSCollectionID]*MockCollection
}

type Mock struct {
	Collections map[types.FSCollectionID]*MockCollection
}

// Get implements interfaces.Firestore.
func (m *Mock) Get(ctx context.Context, value any, docRefs ...types.FireStoreRef) error {
	var target *MockDoc

	collections := m.Collections
	// get data from mock
	for _, docRef := range docRefs {
		if collections == nil {
			return nil
		}
		collection, ok := collections[docRef.CollectionID]
		if !ok {
			return nil
		}

		doc, ok := collection.Docs[docRef.DocumentID]
		if !ok {
			return nil
		}

		target = doc
		collections = doc.Collections
	}

	if target.Value == nil {
		return nil
	}

	raw, err := json.Marshal(target.Value)
	if err != nil {
		return goerr.Wrap(err)
	}
	if err := json.Unmarshal(raw, value); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}

// Put implements interfaces.Firestore.
func (m *Mock) Put(ctx context.Context, value any, docRefs ...types.FireStoreRef) error {
	collections := m.Collections
	var target *MockDoc

	// get or create data from mock
	for _, docRef := range docRefs {
		collection, ok := collections[docRef.CollectionID]
		if !ok {
			collection = &MockCollection{Docs: map[types.FSDocumentID]*MockDoc{}}
			collections[docRef.CollectionID] = collection
		}

		doc, ok := collection.Docs[docRef.DocumentID]
		if !ok {
			doc = &MockDoc{Collections: map[types.FSCollectionID]*MockCollection{}}
			collection.Docs[docRef.DocumentID] = doc
		}

		target = doc
		if target.Collections == nil {
			target.Collections = map[types.FSCollectionID]*MockCollection{}
		}
		collections = target.Collections
	}

	target.Value = value

	return nil
}

func NewMock() *Mock {
	return &Mock{
		Collections: map[types.FSCollectionID]*MockCollection{},
	}
}

var _ interfaces.Firestore = &Mock{}
