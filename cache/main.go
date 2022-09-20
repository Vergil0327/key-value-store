package cache

import (
	"kvstore/cache/provider"
	"sync"
)

type Cache struct {
	mux        sync.Mutex
	upperbound uint
	size       uint
	provider   provider.CacheProvider
}

type Config struct {
	// maximum storage size
	StorageSize uint

	// can change cache policy by using different provider
	Cacher provider.CacheProvider
}

// Constructs cache with given configuration
func NewCache(config Config) (*Cache, error) {
	if config.StorageSize <= 0 {
		return nil, ErrStorageLimit
	}

	if config.Cacher == nil {
		config.Cacher = provider.NewFIFO() // default cache provider
	}

	return new(Cache).init(config), nil
}

func (c *Cache) init(config Config) *Cache {
	c.upperbound = config.StorageSize
	c.provider = config.Cacher

	return c
}

// Looks up value from cache
func (c *Cache) Get(key string) any {
	c.mux.Lock()
	defer c.mux.Unlock()
	ent := c.provider.Get(key)
	if ent != nil {
		return ent.Value()
	}

	return nil
}

// Set value to cache
func (c *Cache) Set(key string, value any) uint {
	c.mux.Lock()
	defer c.mux.Unlock()

	ent := c.provider.NewEntry(key, value)
	c.provider.Set(key, ent)
	c.size += ent.Size()

	var evicted uint
	for c.size > c.upperbound {
		v := c.provider.Evict()
		evicted += v
		c.size -= v
	}

	return evicted
}

// Looks up value from cache without updating.
// ex. `Peek` will NOT update recently used one in LRU
func (c *Cache) Peek(key string) any {
	c.mux.Lock()
	defer c.mux.Unlock()

	ent := c.provider.Peek(key)
	if ent != nil {
		return ent.Value()
	}

	return nil
}

// Get current storage size
func (c *Cache) Size() uint {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.size
}
