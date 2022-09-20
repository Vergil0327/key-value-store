package keyvaluestore

import (
	"kvstore/cache"
	"kvstore/cache/provider"
)

type KeyValueStore struct {
	store map[string]string
	cache *cache.Cache
}

func New() KeyValueStore {
	cache, err := cache.NewCache(cache.Config{
		StorageSize: 100,
		Cacher:      provider.NewFIFO(),
	})
	if err != nil {
		panic(err)
	}

	return KeyValueStore{
		store: make(map[string]string),
		cache: cache,
	}
}

func (kvs *KeyValueStore) Get(key string) string {
	v := kvs.cache.Get(key)
	if v != nil {
		return v.(string)
	}

	if v := kvs.store[key]; v != "" {
		// Use read-through strategy
		kvs.cache.Set(key, v)

		return v
	}

	return ""
}

func (kvs *KeyValueStore) Set(key string, value string) {
	kvs.store[key] = value

	// Use write-through strategy
	kvs.cache.Set(key, value)
}
