package types

import (
	"encoding/hex"
	"math/rand"
	"time"
)

// HashLen is the length of type 'Hash'.
const HashLen = 32

// Hash is the hash type.
type Hash [HashLen]byte

// NilHash is nil with type 'Hash'
var NilHash Hash

// NewHash transvers byte slice to Hash.
func NewHash(s []byte) Hash {
	var h Hash

	if s == nil {
		return h
	}

	l := len(s)
	if l == 0 {
		return h
	}

	if l > HashLen {
		copy(h[:HashLen], s[:HashLen])
	} else {
		copy(h[:l], s[:l])
	}

	return h
}

// NewHashString transvers string to Hash.
func NewHashString(s string) Hash {
	return NewHash([]byte(s))
}

func NewHashHex(hexString string) Hash {
	str, err := hex.DecodeString(hexString)
	if err != nil {
		return NilHash
	}

	var h Hash
	copy(h[:], str[:HashLen])

	return h
}

func RandomHash() Hash {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, HashLen)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return Sha256(b)
}

func (h Hash) Bytes() []byte {
	var b = make([]byte, HashLen, HashLen)
	copy(b[:], h[:])

	return b
}

func (h Hash) String() string {
	var b = make([]byte, HashLen, HashLen)
	copy(b[:], h[:])

	return hex.EncodeToString(b[:])
}
