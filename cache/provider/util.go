package provider

import "reflect"

type CacheProvider interface {
	Get(key string) CacheEntry
	Set(key string, value CacheEntry) uint
	Peek(key string) CacheEntry
	Evict() (evicted uint)
	NewEntry(key string, value any) CacheEntry
}

// return cache entry size
type SizeGetter interface {
	Size() uint
}

type CacheEntry interface {
	SizeGetter
	Value() any
	Key() any
}

// sum of item's size in bytes
func calculateSize(items []any) uint {
	var size uint
	for _, item := range items {
		switch v := item.(type) {
		case string:
			size += uint(len(v))
		default:
			size += uint(reflect.TypeOf(item).Size())
		}
	}
	return size
}
