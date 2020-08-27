package types

import "bytes"

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

func (h Hash) String() string {
	var b = make([]byte, HashLen, HashLen)
	copy(b[:], h[:])

	index := bytes.IndexByte(b, 0)
	return string(b[:index])
}
