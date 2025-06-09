package lfu

import (
	"errors"
	"iter"
	"lfucache/internal/linkedlist"
)

var ErrKeyNotFound = errors.New("key not found")

const DefaultCapacity = 5

// Cache
// O(capacity) memory
type Cache[K comparable, V any] interface {
	// Get returns the value of the key if the key exists in the cache,
	// otherwise, returns ErrKeyNotFound.
	//
	// O(1)
	Get(key K) (V, error)

	// Put updates the value of the key if present, or inserts the key if not already present.
	//
	// When the cache reaches its capacity, it should invalidate and remove the least frequently used key
	// before inserting a new item. For this problem, when there is a tie
	// (i.e., two or more keys with the same frequency), the least recently used key would be invalidated.
	//
	// O(1)
	Put(key K, value V)

	// All returns the iterator in descending order of frequency.
	// If two or more keys have tmented")he same frequency, the most recently used key will be listed first.
	//
	// O(capacity)
	All() iter.Seq2[K, V]

	// Size returns the cache size.
	//
	// O(1)
	Size() int

	// Capacity returns the cache capacity.
	//
	// O(1)
	Capacity() int

	// GetKeyFrequency returns the element's frequency if the key exists in the cache,
	// otherwise, returns ErrKeyNotFound.
	//
	// O(1)
	GetKeyFrequency(key K) (int, error)
}

// node is an element, which contains
// key and value from lfu
type node[K comparable, V any] struct {
	key   K
	value V
	freq  int
}

// cacheImpl represents LFU cache implementation
type cacheImpl[K comparable, V any] struct {
	freqToList linkedlist.ListInterface[linkedlist.ListInterface[node[K, V]]]
	freqToElem map[int]*linkedlist.Element[linkedlist.ListInterface[node[K, V]]]
	keyToElem  map[K]*linkedlist.Element[node[K, V]]
	capacity   int
	minFreq    int
}

// New initializes the cache with the given capacity.
// If no capacity is provided, the cache will use DefaultCapacity.
func New[K comparable, V any](capacity ...int) *cacheImpl[K, V] {
	cap := DefaultCapacity
	if len(capacity) > 0 {
		if cap = capacity[0]; cap < 0 {
			panic("Capacity must be a positive integer")
		}
	}
	return &cacheImpl[K, V]{
		freqToList: linkedlist.NewList[linkedlist.ListInterface[node[K, V]]](),                   // List of nodes with frequencies that store lists of nodes with elements
		freqToElem: make(map[int]*linkedlist.Element[linkedlist.ListInterface[node[K, V]]], cap), // Map from frequency to node with frequency
		keyToElem:  make(map[K]*linkedlist.Element[node[K, V]], cap),                             // Map from the key to the node with the element
		capacity:   cap,
		minFreq:    0,
	}
}

// clearNodes clears node, which contained a current element,
// and checks, that next node exists
func (l *cacheImpl[K, V]) clearNodes(frequency int) {
	if _, ok := l.freqToElem[frequency+1]; !ok {
		list := linkedlist.NewList[node[K, V]]()
		l.freqToElem[frequency+1] = l.freqToList.InsertAfter(list, l.freqToElem[frequency])
	}

	if l.freqToElem[frequency].Value.Len() == 0 {
		l.freqToList.Remove(l.freqToElem[frequency])
		delete(l.freqToElem, frequency)
		if frequency == l.minFreq {
			l.minFreq++
		}
	}
}

// incrementFreq get a current element,
// increment its frequency and checks nodes
func (l *cacheImpl[K, V]) incrementFreq(key K) {
	elem := l.keyToElem[key]
	curFreq := elem.Value.freq
	l.freqToElem[curFreq].Value.Remove(elem)
	l.clearNodes(curFreq)
	elem.Value.freq++
	l.keyToElem[key] = l.freqToElem[curFreq+1].Value.PushFront(elem.Value)
}

func (l *cacheImpl[K, V]) Get(key K) (V, error) {
	if _, ok := l.keyToElem[key]; !ok {
		var zeroValue V
		return zeroValue, ErrKeyNotFound
	}
	l.incrementFreq(key)
	return l.keyToElem[key].Value.value, nil
}

// Remove least frequently used item from cache
func (l *cacheImpl[K, V]) extractLatest() {
	list := l.freqToElem[l.minFreq]
	elem := list.Value.Remove(list.Value.Back())
	delete(l.keyToElem, elem.key)
}

func (l *cacheImpl[K, V]) Put(key K, value V) {
	if l.capacity == 0 {
		return
	}

	// if an element exists, we must increase its frequency
	if elem, ok := l.keyToElem[key]; ok {
		elem.Value.value = value
		l.incrementFreq(key)
		return
	}

	if l.Size() == l.capacity {
		l.extractLatest()
	}

	// Create new node
	l.minFreq = 1
	n := node[K, V]{
		key:   key,
		value: value,
		freq:  l.minFreq,
	}

	// Check if there is a node with the next frequency in the list
	if _, ok := l.freqToElem[l.minFreq]; !ok {
		list := linkedlist.NewList[node[K, V]]()
		l.freqToElem[l.minFreq] = l.freqToList.PushFront(list)
	}
	l.keyToElem[key] = l.freqToElem[l.minFreq].Value.PushFront(n)
}

func (l *cacheImpl[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		iter1 := l.freqToList.NewIteratorPrev()
		for nodeFreq, ok := iter1(); ok; nodeFreq, ok = iter1() {
			iter2 := nodeFreq.Value.NewIteratorNext()
			for node, skp := iter2(); skp; node, skp = iter2() {
				if !yield(node.Value.key, node.Value.value) {
					return
				}
			}
		}
	}
}

// Size return size of cache
func (l *cacheImpl[K, V]) Size() int {
	return len(l.keyToElem)
}

func (l *cacheImpl[K, V]) Capacity() int {
	return l.capacity
}

func (l *cacheImpl[K, V]) GetKeyFrequency(key K) (int, error) {
	// l.Get(key)
	if el, ok := l.keyToElem[key]; ok {
		return el.Value.freq, nil
	}
	return -1, ErrKeyNotFound
}

// [...] <-> [...] <-> [...] <-> [...] <-> [...] <-> [...]
//   |         |         |
// [...]     [...]     [...]
