package keyvaluestore

import (
	"container/list"
	"errors"
	"reflect"
	"sync"
)

type Cacher interface {
	Get(key string) any
	Set(key string, value any)
}

type Cache struct {
	mux  sync.Mutex
	size uint

	values   map[string]*entry
	evitList *list.List
}

type entry struct {
	k    string
	v    any
	size uint
}

type Config struct {
	// maximum storage size
	Size uint

	// TODO: Support configuring the cache replacement policy. Your implementation should support FIFO and provide the flexibility for adding another policy such as LRU in the future
	Policy any

	// TODO: Support both read-through and write-through caching strategies
	Strategy any
}

var (
	ErrStorageLimit = errors.New("storage size should be greater than zero")
)

// Constructs cache with given configuration
func NewCache(config Config) (*Cache, error) {
	if config.Size <= 0 {
		return nil, ErrStorageLimit
	}

	return new(Cache).init(config.Size), nil
}

func (c *Cache) init(size uint) *Cache {
	c.size = size
	c.evitList = list.New()
	c.values = make(map[string]*entry)
	return c
}

// TODO: thread-safe
func (c *Cache) Get(key string) any {
	c.mux.Lock()
	defer c.mux.Unlock()

	if v, ok := c.values[key]; ok {
		return v
	}

	return nil
}

// TODO: thread-safe
func (c *Cache) Set(key string, value any) {
	c.mux.Lock()
	defer c.mux.Unlock()

	size := c.calculateSize([]any{key, value})
	ent := &entry{k: key, v: value, size: size}

	c.values[key] = ent
	c.evitList.PushFront(ent)
	c.size += size
}

// return size in bytes
func (c *Cache) calculateSize(items []any) uint {
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
