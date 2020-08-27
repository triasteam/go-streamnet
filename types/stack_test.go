package types

import "testing"

func TestStack(t *testing.T) {
	s := NewStack()

	h1 := NewHashString("h1")
	h2 := NewHashString("h2")

	s.Push(h1)
	s.Push(h2)

	if h2 != s.Peek() {
		t.Fatal("Peek error!")
	}

	h := s.Pop()
	if h != h2 {
		t.Fatal("Pop error!")
	}

	if s.length != 1 {
		t.Fatal("Length != 1")
	}
}
