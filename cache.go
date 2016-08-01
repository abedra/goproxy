package main

import "container/list"

type LRU struct {
	Maximum int
	entries *list.List
	cache   map[interface{}]*list.Element
}

type Key interface{}

type entry struct {
	key   Key
	value interface{}
}

func NewCache(max int) *LRU {
	return &LRU{
		Maximum: max,
		entries: list.New(),
		cache:   make(map[interface{}]*list.Element),
	}
}

func (lru *LRU) Add(key Key, value interface{}) {
	if lru.cache == nil {
		lru.cache = make(map[interface{}]*list.Element)
		lru.entries = list.New()
	}

	if ent, ok := lru.cache[key]; ok {
		lru.entries.MoveToFront(ent)
		ent.Value.(*entry).value = value
		return
	}

	element := lru.entries.PushFront(&entry{key, value})
	lru.cache[key] = element

	if lru.Maximum != 0 && lru.entries.Len() > lru.Maximum {
		lru.RemoveOldest()
	}
}

func (lru *LRU) Get(key Key) (value interface{}, ok bool) {
	if lru.cache == nil {
		return
	}

	if element, hit := lru.cache[key]; hit {
		lru.entries.MoveToFront(element)
		return element.Value.(*entry).value, true
	}

	return
}

func (lru *LRU) Remove(key Key) {
	if lru.cache == nil {
		return
	}

	if element, hit := lru.cache[key]; hit {
		lru.removeElement(element)
	}
}

func (lru *LRU) RemoveOldest() {
	if lru.cache == nil {
		return
	}

	element := lru.entries.Back()

	if element != nil {
		lru.removeElement(element)
	}
}

func (lru *LRU) removeElement(element *list.Element) {
	lru.entries.Remove(element)
	kv := element.Value.(*entry)
	delete(lru.cache, kv.key)
}

func (lru *LRU) Len() int {
	if lru.cache == nil {
		return 0
	}

	return lru.entries.Len()
}
