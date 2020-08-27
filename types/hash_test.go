package types

import "testing"

func TestNewHash(t *testing.T) {
	h := NewHash(nil)
	t.Logf("%q - %d", h, len(h))
}

func TestHash(t *testing.T) {
	h1 := NewHash([]byte("hello, world"))
	h2 := NewHashString("hello, world")

	s1 := h1.String()
	s2 := h2.String()

	if s1 != s2 || len(s1) != 12 || len(s2) != 12 {
		t.Fatal("Error!")
	}
}

func TestHashEqual(t *testing.T) {
	h1 := NewHashString("hello")
	h2 := NewHashString("hello")
	if h1 != h2 {
		t.Fatal("Unequal")
	}
}
