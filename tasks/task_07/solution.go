package main

import (
	"container/list"
	"sync"
)

type LRU[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
}
type entry[K comparable, V any] struct {
	key   K
	value V
}
type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	ll       list.List
	items    map[K]*list.Element
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		ll:       list.List{},
		items:    make(map[K]*list.Element),
	}
}

func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.items[key]
	if !ok {
		return value, false
	}

	c.ll.MoveToFront(v)
	item := v.Value.(*entry[K, V])
	return item.value, ok
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.items[key]

	switch {
	case c.capacity <= 0:
		return
	case ok:
		item := v.Value.(*entry[K, V])
		item.value = value
		return
	case len(c.items) >= c.capacity:
		kick := c.ll.Back()
		item := kick.Value.(*entry[K, V])
		delete(c.items, item.key)
		c.ll.Remove(kick)
	}

	v = c.ll.PushFront(&entry[K, V]{
		key:   key,
		value: value,
	})
	c.items[key] = v
}
