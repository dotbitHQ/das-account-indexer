package crypto

import (
	"crypto/sha256"
)

// Hash implements crypto sha256 hash algorithm.
func Hash(s []byte) []byte {
	h := sha256.New()
	// Hash.Write never returns an error per godoc
	h.Write(s)
	return h.Sum(nil)
}
