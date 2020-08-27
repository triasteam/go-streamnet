package types

import (
	"sync"
)

type (
	Stack struct {
		top    *node
		length int
		lock   *sync.RWMutex
	}
	node struct {
		value Hash
		prev  *node
	}
)

// Create a new stack
func NewStack() *Stack {
	return &Stack{nil, 0, &sync.RWMutex{}}
}

// Return the number of items in the stack
func (this *Stack) Len() int {
	return this.length
}

func (this *Stack) Empty() bool {
	return this.length == 0
}

// View the top item on the stack
func (this *Stack) Peek() Hash {
	if this.length == 0 {
		return NewHash(nil)
	}
	return this.top.value
}

// Pop the top item of the stack and return it
func (this *Stack) Pop() Hash {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.length == 0 {
		return NewHash(nil)
	}
	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push a value onto the top of the stack
func (this *Stack) Push(value Hash) {
	this.lock.Lock()
	defer this.lock.Unlock()
	n := &node{value, this.top}
	this.top = n
	this.length++
}
