package utils

import (
	"crypto/sha256"

	"github.com/m-mizutani/octovy/pkg/domain/types"
)

func HashBranch(branch string) types.FSDocumentID {
	h := sha256.New()
	h.Write([]byte(branch))
	return types.FSDocumentID(h.Sum(nil))
}
