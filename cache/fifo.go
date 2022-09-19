package cache

type FIFO struct {
}

type EntryFIFO struct {
	key, value string
}

func (ent EntryFIFO) Size() uint {
	return calculateSize([]any{ent.key, ent.value})
}

func NewFIFO() *FIFO {
	return &FIFO{}
}

func (cache *FIFO) Get(key string) CacheEntry {
	return nil
}
func (cache *FIFO) Set(key string, value CacheEntry) {

}
func (cache *FIFO) Peek(key string) CacheEntry {
	return nil
}
func (cache *FIFO) Size() uint {
	return 0
}
