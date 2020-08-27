package types

import (
	"crypto/sha256"
)

func Sha256(bytes []byte) Hash {
	hasher := sha256.New()
	hasher.Write(bytes)
	sum := hasher.Sum(nil)

	var hash Hash
	copy(hash[:], sum[:4])
	return hash
}
