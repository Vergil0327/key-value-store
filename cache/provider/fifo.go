package provider

type FIFO struct {
	queue []*entryFIFO
	cache map[string]*entryFIFO
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
		queue: make([]*entryFIFO, 0),
		cache: make(map[string]*entryFIFO),
	}
}

func (f *FIFO) NewEntry(key string, value any) CacheEntry {
	return &entryFIFO{key: key, value: value}
}

func (f *FIFO) Get(key string) CacheEntry {
	if v, ok := f.cache[key]; ok {
		return v
	}

	return nil
}

func (f *FIFO) Set(key string, value CacheEntry) {
	ent := &entryFIFO{key: key, value: value.Value()}

	if _, ok := f.cache[key]; !ok {
		f.queue = append(f.queue, ent)
	}
	f.cache[key] = ent
}

func (f *FIFO) Peek(key string) CacheEntry {
	if v, ok := f.cache[key]; ok {
		return v
	}

	return nil
}

func (f *FIFO) Evict() (evicted uint) {
	if len(f.queue) == 0 {
		return 0
	}

	item := f.queue[0]
	f.queue = f.queue[1:]
	delete(f.cache, item.key)
	return item.Size()
}
