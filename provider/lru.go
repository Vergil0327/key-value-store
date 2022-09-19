package provider

import "kvstore/cache"

type LRU struct {
}

type EntryLRU struct {
	key, value string
}

func (ent EntryLRU) Size() uint {
	return calculateSize([]any{ent.key, ent.value})
}

func (cache *LRU) Get(key string) cache.CacheEntry {
	return nil
}
func (cache *LRU) Set(key string, value cache.CacheEntry) {

}
func (cache *LRU) Peek(key string) cache.CacheEntry {
	return nil
}
func (cache *LRU) Size() uint {
	return 0
}
