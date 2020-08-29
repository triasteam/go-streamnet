package types

import (
	"log"
	"testing"
)

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

func TestNewHashHex(t *testing.T) {
	str := "a6d3735d908f57647e2bee133643b6163c3888a245c2dea9980693d0d617d48e"
	h := NewHashHex(str)
	log.Print(h)
	if h[0] != 0xa6 || h[HashLen-1] != 0x8e {
		log.Fatal("Hex decode failed!")
	}
}
