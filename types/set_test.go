package types

import (
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet()

	var h1, h2, h3 Hash

	s.Add(h1)
	s.Add(h1)
	s.Add(h2)

	if !s.IsEmpty() {
		t.Logf("0 item")
	}

	s.Clear()
	if s.IsEmpty() {
		t.Logf("0 item")
	}

	s.Add(h1)
	s.Add(h2)
	s.Add(h3)

	if s.Has(h2) {
		t.Logf("2 does exist")
	}

	s.Remove(h2)
	s.Remove(h3)
	t.Log("list of all items", s.List())
}

func TestMapSet(t *testing.T) {
	m := make(map[Hash]Set)

	k := NewHashString("key")
	v := NewHashString("value")

	if _, ok := m[k]; !ok {
		t.Log("Not nil!")
		m[k] = NewSet()
	}

	m[k].Add(v)

	if m[k].Len() != 1 {
		t.Fatal("Lenght is not 1!")
	}

	if !m[k].Has(v) {
		t.Fatal("Don't have value!")
	}

	v2 := m[k].List()[0]
	if v != v2 {
		t.Fatal("Value not equal!")
	}

	t.Log(v2, len(v2))
}
