package types

import (
	"sync"
)

// Set implement 'set' struct.
type Set struct {
	m map[Hash]bool
	sync.RWMutex
}

// NewSet creates a new set.
func NewSet() Set {
	return Set{
		m: map[Hash]bool{},
	}
}

// Add adds new item to set.
func (s Set) Add(item Hash) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

// AddAll adds all items to set.
func (s Set) AddAll(items Set) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items.List() {
		s.m[item] = true
	}
}

// Remove deletes an item from set.
func (s Set) Remove(item Hash) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

// Has checks whether an item is in set.
func (s Set) Has(item Hash) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// Len returns the length of set.
func (s Set) Len() int {
	return len(s.List())
}

// Clear removes all items from set.
func (s Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[Hash]bool{}
}

// IsEmpty checks whether the set is empty.
func (s Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

// List returns all items in the set.
func (s Set) List() []Hash {
	s.RLock()
	defer s.RUnlock()
	list := []Hash{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
