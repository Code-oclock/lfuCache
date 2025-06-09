// Package linkedlist implements a doubly linked list.
package linkedlist

type ListInterface[T any] interface {
	PushFront(v T) *Element[T]
	PushBack(v T) *Element[T]
	InsertAfter(v T, mark *Element[T]) *Element[T]
	InsertBefore(v T, mark *Element[T]) *Element[T]
	Remove(el *Element[T]) T
	Len() int
	Front() *Element[T]
	Back() *Element[T]
	NewIteratorNext() func() (*Element[T], bool)
	NewIteratorPrev() func() (*Element[T], bool)
}

// Element is an element of type T of a linked list.
type Element[T any] struct {
	// next and prev - pointers to the next and previous element, respectively
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front())
	next, prev *Element[T]
	// The value of type T stored with this element.
	Value T
}

// Next returns the next list element
func (e *Element[T]) Next() *Element[T] {
	return e.next
}

// Prev returns the previous list element
func (e *Element[T]) Prev() *Element[T] {
	return e.prev
}

// list represents a doubly linked list.
type list[T any] struct {
	// sentinel list element, only &root, root.prev, and root.next are used
	root Element[T]
	// current list length excluding (this) sentinel element
	len int
}

// NewList returns an initialized list.
func NewList[T any]() ListInterface[T] {
	var list list[T]
	list.root.next = &list.root
	list.root.prev = &list.root
	list.len = 0
	return &list
}

// nodeConnection connects two nodes
func (l *list[T]) nodeConnection(first, second *Element[T]) {
	first.next = second
	second.prev = first
}

// insert inserts value v after dst element, increments l.len, and returns el
func (l *list[T]) insert(v T, dst *Element[T]) *Element[T] {
	// It's create an element and insert it
	el := &Element[T]{Value: v}
	l.nodeConnection(el, dst.next)
	l.nodeConnection(dst, el)
	l.len++
	return el
}

// remove removes el from its list, decrements l.len
func (l *list[T]) remove(el *Element[T]) {
	if el.prev == nil || el.next == nil {
		return
	}
	l.nodeConnection(el.prev, el.next)
	el.next = nil
	el.prev = nil
	l.len--
}

// Root returns root of list
// The complexity is O(1).
func (l *list[T]) Root() *Element[T] {
	return &l.root
}

// Len returns len of list
// The complexity is O(1).
func (l *list[T]) Len() int {
	return l.len
}

// Front returns first element of list or nil
// The complexity is O(1).
func (l *list[T]) Front() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns last element of list or nil
// The complexity is O(1).
func (l *list[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// PushFront insert v after root
// The complexity is O(1).
func (l *list[T]) PushFront(v T) *Element[T] {
	return l.insert(v, &l.root)
}

// PushBack insert v before root
// The complexity is O(1).
func (l *list[T]) PushBack(v T) *Element[T] {
	return l.insert(v, l.root.prev)
}

// InsertAfter insert v after element mark
// The complexity is O(1).
func (l *list[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	return l.insert(v, mark)
}

// InsertBefore insert v before element mark
// The complexity is O(1).
func (l *list[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	return l.insert(v, mark.prev)
}

// Remove removes e from l
// The element must not be nil
// The complexity is O(1).
func (l *list[T]) Remove(el *Element[T]) T {
	l.remove(el)
	return el.Value
}

// NewIteratorNext returns a func,
// which contains iterator, which go forward.
func (l *list[T]) NewIteratorNext() func() (*Element[T], bool) {
	current := l.Front()
	return func() (*Element[T], bool) {
		if current == nil || current == &l.root { // if list doesn't contain element
			return nil, false
		}
		el := current
		current = current.next
		return el, true
	}
}

// NewIteratorPrev returns a func,
// which contains iterator, which go back.
func (l *list[T]) NewIteratorPrev() func() (*Element[T], bool) {
	current := l.Back()
	return func() (*Element[T], bool) {
		if current == nil || current == &l.root { // if list doesn't contain element
			return nil, false
		}
		el := current
		current = current.prev
		return el, true
	}
}
