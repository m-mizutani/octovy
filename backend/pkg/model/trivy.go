package model

import (
	"encoding/json"

	"github.com/m-mizutani/goerr"
)

type AdvisoryData struct {
	VulnID string
	Data   []byte
}

func (x *AdvisoryData) Unmarshal(v interface{}) error {
	if err := json.Unmarshal(x.Data, v); err != nil {
		return goerr.Wrap(err, "Failed to unmarshal Advisory data").With("data", string(x.Data))
	}
	return nil
}
