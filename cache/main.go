package cache

import (
	"sync"
)

type Cacher interface {
	Get(key string) CacheEntry
	Set(key string, value CacheEntry)
	Peek(key string) CacheEntry
	Size() uint
}

type CacheEntry interface {
	Size() uint
}

type Cache struct {
	mux      sync.Mutex
	size     uint
	provider Cacher
}

type Config struct {
	// maximum storage size
	Size uint

	// TODO: Support configuring the cache replacement policy. Your implementation should support FIFO and provide the flexibility for adding another policy such as LRU in the future
	Cache Cacher

	// TODO: Support both read-through and write-through caching strategies
	Strategy any
}

// Constructs cache with given configuration
func NewCache(config Config) (*Cache, error) {
	if config.Size <= 0 {
		return nil, ErrStorageLimit
	}

	if config.Cache == nil {
		config.Cache = NewFIFO() // default cache provider
	}

	return new(Cache).init(config.Size, config.Cache), nil
}

func (c *Cache) init(size uint, provider Cacher) *Cache {
	c.size = size
	c.provider = provider

	return c
}

func (c *Cache) Get(key string) CacheEntry {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.provider.Get(key)
}

func (c *Cache) Set(key string, value CacheEntry) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.provider.Set(key, value)

	// handle evic & update size
	c.size += c.provider.Size()
}

func (c *Cache) Peek(key string) CacheEntry {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.provider.Peek(key)
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
