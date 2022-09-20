package provider

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

func (cache *FIFO) Get(key string) SizeGetter {
	return &EntryFIFO{}
}
func (cache *FIFO) Set(key string, value SizeGetter) {

}
func (cache *FIFO) Peek(key string) SizeGetter {
	return &EntryFIFO{}
}
func (cache *FIFO) Evict() uint {
	return 0
}
func (cache *FIFO) Size() uint {
	return 0
}
