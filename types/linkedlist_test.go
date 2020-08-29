package types

import (
	"log"
	"testing"
)

/*
func TestList(t *testing.T) {
	l := List{}

	if !l.IsEmpty() {
		t.Fatal("Not empty!")
	}

	l.Append(1)
	l.Append(2)
	l.Append("a")
	l.Append("b")

	if l.IsEmpty() {
		t.Fatal("Is empty!!!")
	}
	if l.Length() != 4 {
		t.Fatal("Length != 4")
	}

	if !l.Contain(1) {
		t.Fatal("Don't have 1")
	}

	l.Remove(1)
	if l.Contain(1) {
		t.Fatal("Have 1")
	}

	v := l.RemoveAtIndex(0)
	t.Log(v)
	if l.Contain(2) {
		t.Fatal("Delete 2 error!")
	}
	if int(v) != 2 {
		t.Fatal("Not equal to 2")
	}
}
*/

func TestList(t *testing.T) {
	h := NewHashString("nihao")
	j := NewHashString("hello")

	l := List{}

	if !l.IsEmpty() {
		t.Fatal("Not empty!")
	}

	l.Append(h)
	l.Append(j)

	if l.IsEmpty() {
		t.Fatal("Is empty!!!")
	}
	if l.Length() != 2 {
		t.Fatal("Length != 4")
	}

	if !l.Contain(h) {
		t.Fatal("Don't have h!")
	}

	l.Remove(h)
	if l.Contain(h) {
		t.Fatal("Have h!")
	}

	v := l.RemoveAtIndex(0)
	t.Log(v)
	if l.Contain(j) {
		t.Fatal("Delete j error!")
	}
	if v != j {
		t.Fatal("Not equal: v != j")
	}

	if !l.IsEmpty() {
		t.Fatal("Not empty!")
	}
}

func TestAppendAndIndex(t *testing.T) {
	l := List{}
	trunk := Sha256([]byte("StreamNet_Trunk"))
	branch := Sha256([]byte("StreamNet_Branch"))

	l.Append(trunk)
	l.Append(branch)

	tr1 := l.Index(0)
	br1 := l.Index(1)
	log.Println(tr1, br1)

	if trunk != tr1 || branch != br1 {
		log.Fatal("Add and get not same!")
	}
}
