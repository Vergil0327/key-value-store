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

func (c *Cache) Get(key string) any {
	c.mux.Lock()
	defer c.mux.Unlock()
	ent := c.provider.Get(key)

	return ent.Value()
}

func (c *Cache) Set(key string, value any) uint {
	c.mux.Lock()
	defer c.mux.Unlock()

	ent := c.provider.NewEntry(key, value)
	c.provider.Set(key, ent)
	c.size += ent.Size()

	var evicted uint
	for c.size > c.upperbound {
		evicted += c.provider.Evict()
	}

	return evicted
}

func (c *Cache) Peek(key string) any {
	c.mux.Lock()
	defer c.mux.Unlock()

	ent := c.provider.Peek(key)
	return ent.Value()
}

func (c *Cache) Size() uint {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.size
}

/* TODO:
1. Support configuring the maximum storage size
	- storage size is defined as the sum of the number of bytes of all keys and values
	- This definition is intentionally simplified so don't count the size of other on-disk or in-memory data structures and metadata

2. Support both read-through and write-through caching strategies
3. Support configuring the cache replacement policy. Your implementation should support FIFO and provide the flexibility for adding another policy such as LRU in the future
4. Support get and set methods
5. thread-safe
6. unit tests
7. README
*/
