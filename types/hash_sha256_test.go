package types

import (
	"log"
	"testing"
)

func TestSha256(t *testing.T) {
	data := "hello, world"
	h := Sha256([]byte(data))
	log.Printf("lenght = %d\n", len(h))
	log.Printf("%s", h)

	data2 := "hello, world"
	h2 := Sha256([]byte(data2))
	if h != h2 {
		log.Fatal("Hash not equal!")
	}
}
