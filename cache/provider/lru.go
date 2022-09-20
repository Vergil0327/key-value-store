package provider

import (
	"container/list"
)

type LRU struct {
	list  *list.List
	cache map[string]*list.Element
}

type entryLRU struct {
	key   string
	value any
}

func (ent entryLRU) Key() any {
	return ent.key
}
func (ent entryLRU) Value() any {
	return ent.value
}
func (ent entryLRU) Size() uint {
	return calculateSize([]any{ent.key, ent.value})
}

func NewLRU() *LRU {
	return &LRU{
		list:  list.New(),
		cache: make(map[string]*list.Element),
	}
}

func (f *LRU) NewEntry(key string, value any) CacheEntry {
	return &entryLRU{key: key, value: value}
}

// Looks up value from cache
func (f *LRU) Get(key string) CacheEntry {
	if ent, ok := f.cache[key]; ok {
		f.list.MoveToFront(ent)
		return ent.Value.(*entryLRU)
	}

	return nil
}

// Set value to cache
func (f *LRU) Set(key string, value CacheEntry) {
	if ent, ok := f.cache[key]; ok {
		ent.Value.(*entryLRU).value = value.Value()
		f.list.MoveToFront(ent)
		return
	}

	ent := f.NewEntry(key, value.Value())
	el := f.list.PushFront(ent)
	f.cache[key] = el
}

// Looks up value from cache without updating recently used record
func (f *LRU) Peek(key string) CacheEntry {
	if ent, ok := f.cache[key]; ok {
		return ent.Value.(*entryLRU)
	}

	return nil
}

// remove oldest one from cache
func (f *LRU) Evict() (evicted uint) {
	if f.list.Len() == 0 {
		return 0
	}

	oldest := f.list.Back()
	if oldest == nil {
		return 0
	}

	ent := f.list.Remove(oldest).(*entryLRU)
	delete(f.cache, ent.key)
	return ent.Size()
}
