package types

import (
	"sync"
)

// Set implement 'set' struct.
type Set struct {
<<<<<<< HEAD
	m map[Hash]bool
	sync.RWMutex
}

// NewSet creates a new set.
func NewSet() Set {
	return Set{
		m: map[Hash]bool{},
=======
	m map[int]bool
	sync.RWMutex
}

// New creates a new set.
func New() *Set {
	return &Set{
		m: map[int]bool{},
>>>>>>> 688dc7a... implement type 'Set'
	}
}

// Add adds new item to set.
<<<<<<< HEAD
func (s Set) Add(item Hash) {
=======
func (s *Set) Add(item int) {
>>>>>>> 688dc7a... implement type 'Set'
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

<<<<<<< HEAD
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
=======
// Remove deletes an item from set.
func (s *Set) Remove(item int) {
>>>>>>> 688dc7a... implement type 'Set'
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

// Has checks whether an item is in set.
<<<<<<< HEAD
func (s Set) Has(item Hash) bool {
=======
func (s *Set) Has(item int) bool {
>>>>>>> 688dc7a... implement type 'Set'
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// Len returns the length of set.
<<<<<<< HEAD
func (s Set) Len() int {
=======
func (s *Set) Len() int {
>>>>>>> 688dc7a... implement type 'Set'
	return len(s.List())
}

// Clear removes all items from set.
<<<<<<< HEAD
func (s Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[Hash]bool{}
}

// IsEmpty checks whether the set is empty.
func (s Set) IsEmpty() bool {
=======
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[int]bool{}
}

// IsEmpty checks whether the set is empty.
func (s *Set) IsEmpty() bool {
>>>>>>> 688dc7a... implement type 'Set'
	if s.Len() == 0 {
		return true
	}
	return false
}

// List returns all items in the set.
<<<<<<< HEAD
func (s Set) List() []Hash {
	s.RLock()
	defer s.RUnlock()
	list := []Hash{}
=======
func (s *Set) List() []int {
	s.RLock()
	defer s.RUnlock()
	list := []int{}
>>>>>>> 688dc7a... implement type 'Set'
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
