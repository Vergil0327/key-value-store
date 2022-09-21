package provider

import "container/list"

type FIFO struct {
	queue *list.List
	cache map[string]*list.Element // Value: *entryFIFO
}

type entryFIFO struct {
	key   string
	value any
}

func (ent entryFIFO) Key() any {
	return ent.key
}
func (ent entryFIFO) Value() any {
	return ent.value
}
func (ent entryFIFO) Size() uint {
	return calculateSize([]any{ent.key, ent.value})
}

func NewFIFO() *FIFO {
	return &FIFO{
		queue: list.New(),
		cache: make(map[string]*list.Element),
	}
}

func (f *FIFO) NewEntry(key string, value any) CacheEntry {
	return &entryFIFO{key: key, value: value}
}

func (f *FIFO) Get(key string) CacheEntry {
	if el, ok := f.cache[key]; ok {
		return el.Value.(*entryFIFO)
	}

	return nil
}

func (f *FIFO) Set(key string, value CacheEntry) (storage uint) {
	if el, ok := f.cache[key]; ok {
		el.Value.(*entryFIFO).value = value.Value()
		return
	}

	ent := &entryFIFO{key: key, value: value.Value()}
	el := f.queue.PushBack(ent)
	f.cache[key] = el
	storage = ent.Size()

	return storage
}

func (f *FIFO) Peek(key string) CacheEntry {
	if el, ok := f.cache[key]; ok {
		return el.Value.(*entryFIFO)
	}

	return nil
}

func (f *FIFO) Evict() (evicted uint) {
	if f.queue.Len() == 0 {
		return 0
	}

	el := f.queue.Front()
	if el == nil {
		return 0
	}
	f.queue.Remove(el)

	ent := el.Value.(*entryFIFO)
	delete(f.cache, ent.key)
	return ent.Size()
}
