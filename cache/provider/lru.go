package provider

type LRU struct {
}

type EntryLRU struct {
	key, value string
}

func (ent EntryLRU) Size() uint {
	return calculateSize([]any{ent.key, ent.value})
}

func NewLRU() *LRU {
	return &LRU{}
}

func (cache *LRU) Get(key string) SizeGetter {
	return &EntryLRU{}
}
func (cache *LRU) Set(key string, value SizeGetter) {

}
func (cache *LRU) Peek(key string) SizeGetter {
	return &EntryLRU{}
}
func (cache *LRU) Evict() uint {
	return 0
}
func (cache *LRU) Size() uint {
	return 0
}

// // EvictCallback is used to get a callback when a cache entry is evicted
// type EvictCallback func(key interface{}, value interface{})

// // LRU implements a non-thread safe fixed size LRU cache
// // Methods:
// // Purge()
// // Add(key, value interface{}) (evicted bool)
// // Get(key interface{}) (value interface{}, ok bool)
// // Contains(key interface{}) (ok bool)
// // Peek(key interface{}) (value interface{}, ok bool)
// // Remove(key interface{}) (present bool)
// // RemoveOldest() (key, value interface{}, ok bool)
// // GetOldest() (key, value interface{}, ok bool)
// // Keys() []interface{}
// // Len() int
// // Resize(size int) (evicted int)
// // removeOldest()
// // removeElement(e *list.Element)
// type LRU struct {
// 	size      int
// 	evictList *list.List
// 	items     map[interface{}]*list.Element
// 	onEvict   EvictCallback
// }

// // lruEntry is used to hold a value in the evictList
// type lruEntry struct {
// 	key   interface{}
// 	value interface{}
// }

// // NewLRU constructs an LRU of the given size
// func NewLRU(size int, onEvict EvictCallback) (*LRU, error) {
// 	if size <= 0 {
// 		return nil, errors.New("must provide a positive size")
// 	}
// 	c := &LRU{
// 		size:      size,
// 		evictList: list.New(),
// 		items:     make(map[interface{}]*list.Element),
// 		onEvict:   onEvict,
// 	}
// 	return c, nil
// }

// // Purge is used to completely clear the cache.
// func (c *LRU) Purge() {
// 	for k, v := range c.items {
// 		if c.onEvict != nil {
// 			c.onEvict(k, v.Value.(*lruEntry).value)
// 		}
// 		delete(c.items, k)
// 	}
// 	c.evictList.Init()
// }

// // Add adds a value to the cache.  Returns true if an eviction occurred.
// func (c *LRU) Add(key, value interface{}) (evicted bool) {
// 	// Check for existing item
// 	if ent, ok := c.items[key]; ok {
// 		c.evictList.MoveToFront(ent)
// 		ent.Value.(*lruEntry).value = value
// 		return false
// 	}

// 	// Add new item
// 	ent := &lruEntry{key, value}
// 	entry := c.evictList.PushFront(ent)
// 	c.items[key] = entry

// 	evict := c.evictList.Len() > c.size
// 	// Verify size not exceeded
// 	if evict {
// 		c.removeOldest()
// 	}
// 	return evict
// }

// // Get looks up a key's value from the cache.
// func (c *LRU) Get(key interface{}) (value interface{}, ok bool) {
// 	if ent, ok := c.items[key]; ok {
// 		c.evictList.MoveToFront(ent)
// 		if ent.Value.(*lruEntry) == nil {
// 			return nil, false
// 		}
// 		return ent.Value.(*lruEntry).value, true
// 	}
// 	return
// }

// // Contains checks if a key is in the cache, without updating the recent-ness
// // or deleting it for being stale.
// func (c *LRU) Contains(key interface{}) (ok bool) {
// 	_, ok = c.items[key]
// 	return ok
// }

// // Peek returns the key value (or undefined if not found) without updating
// // the "recently used"-ness of the key.
// func (c *LRU) Peek(key interface{}) (value interface{}, ok bool) {
// 	var ent *list.Element
// 	if ent, ok = c.items[key]; ok {
// 		return ent.Value.(*lruEntry).value, true
// 	}
// 	return nil, ok
// }

// // Remove removes the provided key from the cache, returning if the
// // key was contained.
// func (c *LRU) Remove(key interface{}) (present bool) {
// 	if ent, ok := c.items[key]; ok {
// 		c.removeElement(ent)
// 		return true
// 	}
// 	return false
// }

// // RemoveOldest removes the oldest item from the cache.
// func (c *LRU) RemoveOldest() (key, value interface{}, ok bool) {
// 	ent := c.evictList.Back()
// 	if ent != nil {
// 		c.removeElement(ent)
// 		kv := ent.Value.(*lruEntry)
// 		return kv.key, kv.value, true
// 	}
// 	return nil, nil, false
// }

// // GetOldest returns the oldest entry
// func (c *LRU) GetOldest() (key, value interface{}, ok bool) {
// 	ent := c.evictList.Back()
// 	if ent != nil {
// 		kv := ent.Value.(*lruEntry)
// 		return kv.key, kv.value, true
// 	}
// 	return nil, nil, false
// }

// // Keys returns a slice of the keys in the cache, from oldest to newest.
// func (c *LRU) Keys() []interface{} {
// 	keys := make([]interface{}, len(c.items))
// 	i := 0
// 	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
// 		keys[i] = ent.Value.(*lruEntry).key
// 		i++
// 	}
// 	return keys
// }

// // Len returns the number of items in the cache.
// func (c *LRU) Len() int {
// 	return c.evictList.Len()
// }

// // Resize changes the cache size.
// func (c *LRU) Resize(size int) (evicted int) {
// 	diff := c.Len() - size
// 	if diff < 0 {
// 		diff = 0
// 	}
// 	for i := 0; i < diff; i++ {
// 		c.removeOldest()
// 	}
// 	c.size = size
// 	return diff
// }

// // removeOldest removes the oldest item from the cache.
// func (c *LRU) removeOldest() {
// 	ent := c.evictList.Back()
// 	if ent != nil {
// 		c.removeElement(ent)
// 	}
// }

// // removeElement is used to remove a given list element from the cache
// func (c *LRU) removeElement(e *list.Element) {
// 	c.evictList.Remove(e)
// 	kv := e.Value.(*lruEntry)
// 	delete(c.items, kv.key)
// 	if c.onEvict != nil {
// 		c.onEvict(kv.key, kv.value)
// 	}
// }
