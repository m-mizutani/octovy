package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashBranch(branch string) string {
	h := sha256.New()
	h.Write([]byte(branch))
	v := hex.EncodeToString(h.Sum(nil))
	return v
}
