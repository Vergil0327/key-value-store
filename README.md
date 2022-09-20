

# In-Memory Cache For KeyValueStore

## About

An in-memory cache that

- support configuring the maximum storage size

  ```go
  cache := NewCache(Config{StorageSize: 100})
  ```

- support both read-through and write-through caching strategies.

  see [example.go](example.go)
  ```go
  ```

- support configuring cache replacement policy

  ```go
  /* FIFO */
  NewCache(Config{
		StorageSize: 100,
		Cacher:      provider.NewFIFO(),
	})

  /* LRU */
	NewCache(Config{
		StorageSize: 100,
		Cacher:      provider.NewLRU(),
	})

  /*
  or use custom policy by implementing `CacheProvider` interface

  see cache/provider folder for more details
  */

  type CacheProvider interface {
	Get(key string) CacheEntry
	Set(key string, value CacheEntry)
	Peek(key string) CacheEntry
	Evict() (evicted uint)
	NewEntry(key string, value any) CacheEntry
  }
  ```

- support these methods:
  - Get: looks up value from cache
  - Set: set value to cache
  - Peek: looks up value from cache without updating. ex. it won't change most recently used one
  - Size: get current storage size

## Getting started

*required golang version `1.19`*

### Installation

- go to official website [download and install](https://go.dev/doc/install)

### Testing

``` bash
# under root folder
$ make test

# or
$ go test ./...
```