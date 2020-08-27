package types

type Object interface{}

type Node struct {
	Data Hash
	Next *Node
}

type List struct {
	headNode *Node
}

func (l *List) IsEmpty() bool {
	if l.headNode == nil {
		return true
	} else {
		return false
	}
}

func (this *List) Length() int {
	cur := this.headNode
	count := 0

	for cur != nil {
		count++
		cur = cur.Next
	}
	return count
}

func (l *List) Add(data Hash) *Node {
	node := &Node{Data: data}
	node.Next = l.headNode
	l.headNode = node
	return node
}

func (l *List) Append(data Hash) {
	node := &Node{Data: data}
	if l.IsEmpty() {
		l.headNode = node
	} else {
		cur := l.headNode
		for cur.Next != nil {
			cur = cur.Next
		}
		cur.Next = node
	}
}

func (l *List) Insert(index int, data Hash) {
	if index < 0 {
		l.Add(data)
	} else if index > l.Length() {
		l.Add(data)
	} else {
		pre := l.headNode
		count := 0
		for count < index-1 {
			pre = pre.Next
			count++
		}
		node := &Node{Data: data}
		node.Next = pre.Next
		pre.Next = node
	}
}

func (l *List) Remove(data Hash) {
	pre := l.headNode
	if pre.Data == data {
		l.headNode = pre.Next
	} else {
		for pre.Next != nil {
			if pre.Next.Data == data {
				pre.Next = pre.Next.Next
			} else {
				pre = pre.Next
			}
		}
	}
}

func (l *List) RemoveAtIndex(index int) (ret Hash) {
	pre := l.headNode
	if index <= 0 {
		ret = pre.Data
		l.headNode = pre.Next
	} else if index > l.Length() {
		ret = NewHash(nil)
	} else {
		count := 0
		for count != index-1 && pre.Next != nil {
			count++
			pre = pre.Next
		}
		ret = pre.Next.Data
		pre.Next = pre.Next.Next
	}
	return
}

func (l *List) Index(index int) (ret Hash) {
	pre := l.headNode
	if index <= 0 {
		ret = pre.Data
		l.headNode = pre.Next
	} else if index > l.Length() {
		ret = NewHash(nil)
	} else {
		count := 0
		for count != index-1 && pre.Next != nil {
			count++
			pre = pre.Next
		}
		ret = pre.Next.Data
	}
	return
}

func (l *List) GetLast() Hash {
	return l.Index(l.Length() - 1)
}

func (l *List) Contain(data Hash) bool {
	cur := l.headNode
	for cur != nil {
		if cur.Data == data {
			return true
		}
		cur = cur.Next
	}
	return false
}

func (l *List) Show() {
}
