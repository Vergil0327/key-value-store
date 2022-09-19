package keyvaluestore

type Cacher interface {
	Get(key string) any
	Set(key string, value any)
}

type Cache struct {
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
