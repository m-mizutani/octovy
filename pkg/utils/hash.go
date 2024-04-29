package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/m-mizutani/octovy/pkg/domain/types"
)

func HashBranch(branch string) types.FSDocumentID {
	h := sha256.New()
	h.Write([]byte(branch))
	v := hex.EncodeToString(h.Sum(nil))
	return types.FSDocumentID(v)
}
