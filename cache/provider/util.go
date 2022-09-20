package provider

import "reflect"

type CacheProvider interface {
	Get(key string) SizeGetter
	Set(key string, value SizeGetter)
	Peek(key string) SizeGetter
	Evict() (evicted uint)
	Size() uint
}

// return cache entry size
type SizeGetter interface {
	Size() uint
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
