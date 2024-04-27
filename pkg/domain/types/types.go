package types

import "github.com/google/uuid"

type (
	ScanID string

	RequestID string
)

func NewScanID() ScanID         { return ScanID(uuid.NewString()) }
func (x ScanID) String() string { return string(x) }

func NewRequestID() RequestID      { return RequestID(uuid.NewString()) }
func (x RequestID) String() string { return string(x) }

type (
	GoogleProjectID string

	BQDatasetID string
	BQTableID   string

	FSDatabaseID   string
	FSCollectionID string
	FSDocumentID   string
)

func (x GoogleProjectID) String() string { return string(x) }
func (x BQDatasetID) String() string     { return string(x) }
func (x BQTableID) String() string       { return string(x) }
func (x FSDatabaseID) String() string    { return string(x) }
func (x FSCollectionID) String() string  { return string(x) }
func (x FSDocumentID) String() string    { return string(x) }
