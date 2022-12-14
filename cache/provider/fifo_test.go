package provider

import (
	"container/list"
	"math"
	"reflect"
	"testing"
)

func TestEntryFIFO_Key(t *testing.T) {
	type fields struct {
		key   string
		value any
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name:   "it_should_return_entry_key",
			fields: fields{key: "whatever"},
			want:   "whatever",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ent := entryFIFO{
				key:   tt.fields.key,
				value: tt.fields.value,
			}
			if got := ent.Key(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EntryFIFO.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryFIFO_Value(t *testing.T) {
	type fields struct {
		key   string
		value any
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name:   "it_should_return_entry_value",
			fields: fields{value: "whatever"},
			want:   "whatever",
		},
		{
			name:   "it_should_return_any_entry_value_you_store",
			fields: fields{value: 1000},
			want:   1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ent := entryFIFO{
				key:   tt.fields.key,
				value: tt.fields.value,
			}
			if got := ent.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EntryFIFO.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryFIFO_Size(t *testing.T) {
	type fields struct {
		key   string
		value any
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name: "it_should_return_entry_storage_size",
			fields: fields{
				key:   "6bytes",
				value: "should_be_18_bytes",
			},
			want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ent := entryFIFO{
				key:   tt.fields.key,
				value: tt.fields.value,
			}
			if got := ent.Size(); got != tt.want {
				t.Errorf("EntryFIFO.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFIFO(t *testing.T) {
	tests := []struct {
		name string
		want *FIFO
	}{
		{
			name: "it_should_construct_fifo_correctly",
			want: &FIFO{queue: list.New(), cache: make(map[string]*list.Element)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFIFO(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFIFO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFIFO_NewEntry(t *testing.T) {
	type fields struct {
		queue *list.List
		cache map[string]*list.Element
	}
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CacheEntry
	}{
		{
			name: "it_should_construct_entry",
			args: args{key: "KEY", value: "VALUE"},
			want: &entryFIFO{key: "KEY", value: "VALUE"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FIFO{
				queue: tt.fields.queue,
				cache: tt.fields.cache,
			}
			if got := f.NewEntry(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FIFO.NewEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFIFO_Get(t *testing.T) {
	type fields struct {
		queue *list.List
		cache map[string]*list.Element
	}
	type args struct {
		key string
	}
	que1 := list.New()
	ent1 := &entryFIFO{key: "key1", value: "value"}
	el1 := que1.PushBack(ent1)
	m1 := map[string]*list.Element{}
	m1[ent1.key] = el1

	que2 := list.New()
	ent2 := &entryFIFO{key: "key1", value: "value"}
	el2 := que2.PushBack(ent2)
	m2 := map[string]*list.Element{}
	m2[ent2.key] = el2

	tests := []struct {
		name   string
		fields fields
		args   args
		want   CacheEntry
	}{
		{
			name: "it_should_get_from_cache",
			fields: fields{
				queue: que1,
				cache: m1,
			},
			args: args{key: "key1"},
			want: &entryFIFO{key: "key1", value: "value"},
		},
		{
			name: "it_should_get_nil",
			fields: fields{
				queue: que2,
				cache: m2,
			},
			args: args{key: "key2"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FIFO{
				queue: tt.fields.queue,
				cache: tt.fields.cache,
			}
			if got := f.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FIFO.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFIFO_Set(t *testing.T) {
	type fields struct {
		queue *list.List
		cache map[string]*list.Element
	}
	type args struct {
		key   string
		value CacheEntry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "it_should_add_entry_to_cache",
			fields: fields{queue: list.New(), cache: make(map[string]*list.Element)},
			args:   args{key: "key1", value: &entryFIFO{key: "key1", value: "value1"}},
		},
		{
			name:   "it_should_add_entry_to_queue",
			fields: fields{queue: list.New(), cache: make(map[string]*list.Element)},
			args:   args{key: "key2", value: &entryFIFO{key: "key2", value: "value2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FIFO{
				queue: tt.fields.queue,
				cache: tt.fields.cache,
			}
			f.Set(tt.args.key, tt.args.value)
			want := &entryFIFO{key: tt.args.key, value: tt.args.value.Value()}
			if v, ok := f.cache[tt.args.key]; !ok || !reflect.DeepEqual(v.Value, want) {
				t.Errorf("After FIFO.Set(): got %v, want %v", v.Value, want)
			}
			if v := f.queue.Front(); !reflect.DeepEqual(v.Value, want) {
				t.Errorf("After FIFO.Set(): got %v, want %v", v.Value, want)
			}
		})
	}
}

func TestFIFO_Peek(t *testing.T) {
	type fields struct {
		queue *list.List
		cache map[string]*list.Element
	}
	type args struct {
		key string
	}

	que1 := list.New()
	ent1 := &entryFIFO{key: "key1", value: "value"}
	el1 := que1.PushBack(ent1)
	m1 := map[string]*list.Element{}
	m1[ent1.key] = el1

	que2 := list.New()
	ent2 := &entryFIFO{key: "key1", value: "value"}
	el2 := que2.PushBack(ent2)
	m2 := map[string]*list.Element{}
	m2[ent2.key] = el2

	tests := []struct {
		name   string
		fields fields
		args   args
		want   CacheEntry
	}{
		{
			name: "it_should_peek_from_cache",
			fields: fields{
				queue: que1,
				cache: m1,
			},
			args: args{key: "key1"},
			want: &entryFIFO{key: "key1", value: "value"},
		},
		{
			name: "it_should_peek_nil",
			fields: fields{
				queue: que2,
				cache: m2,
			},
			args: args{key: "key2"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FIFO{
				queue: tt.fields.queue,
				cache: tt.fields.cache,
			}
			if got := f.Peek(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FIFO.Peek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFIFO_Evict(t *testing.T) {
	type fields struct {
		queue *list.List
		cache map[string]*list.Element
	}

	que1 := list.New()
	ent1 := &entryFIFO{key: "6bytes", value: "6bytes"}
	que1.PushBack(ent1)

	que2 := list.New()
	ent2a := &entryFIFO{key: "6bytes", value: "6bytes"}
	ent2b := &entryFIFO{key: "key1", value: "value"}
	que2.PushBack(ent2a)
	que2.PushBack(ent2b)

	tests := []struct {
		name        string
		fields      fields
		wantEvicted uint
	}{
		{
			name:        "it_should_evict_nothing",
			fields:      fields{queue: list.New(), cache: make(map[string]*list.Element)},
			wantEvicted: 0,
		},
		{
			name:        "it_should_evict_successfully",
			fields:      fields{queue: que1, cache: make(map[string]*list.Element)},
			wantEvicted: 12,
		},
		{
			name:        "it_should_evict_left_most",
			fields:      fields{queue: que2, cache: make(map[string]*list.Element)},
			wantEvicted: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FIFO{
				queue: tt.fields.queue,
				cache: tt.fields.cache,
			}
			original := f.queue.Len()

			if gotEvicted := f.Evict(); gotEvicted != tt.wantEvicted {
				t.Errorf("FIFO.Evict() = %v, want %v", gotEvicted, tt.wantEvicted)
				if want := int(math.Max(0, float64(original))); f.queue.Len() != want {
					t.Errorf("After FIFO.Evict() queue's length = %v, want %v", f.queue.Len(), want)
				}
			}
		})
	}
}

func TestFIFO(t *testing.T) {
	f := NewFIFO()

	tests := []struct {
		name  string
		check func() bool
	}{
		{
			name: "get_key1",
			check: func() bool {
				got := f.Get("key1")
				return got == nil
			},
		},
		{
			name: "peek_key1",
			check: func() bool {
				got := f.Peek("key1")
				return got == nil
			},
		},
		{
			name: "set_key1",
			check: func() bool {
				ent := f.NewEntry("key1", "value1")
				f.Set("key1", ent)
				el, ok := f.cache["key1"]
				return ok && f.queue.Len() == 1 && reflect.DeepEqual(ent, el.Value) && reflect.DeepEqual(ent, f.queue.Back().Value) && ent.Size() == 10
			},
		},
		{
			name: "set_key2",
			check: func() bool {
				ent := f.NewEntry("key2", "value2")
				f.Set("key2", ent)
				el, ok := f.cache["key2"]
				return ok && f.queue.Len() == 2 && reflect.DeepEqual(ent, el.Value) && reflect.DeepEqual(ent, f.queue.Back().Value)
			},
		},
		{
			name: "get_key1",
			check: func() bool {
				ent := f.NewEntry("key1", "value1")
				got := f.Get("key1")
				return reflect.DeepEqual(ent, got)
			},
		},
		{
			name: "peek_key1",
			check: func() bool {
				ent := f.NewEntry("key1", "value1")
				got := f.Peek("key1")
				return reflect.DeepEqual(ent, got)
			},
		},
		{
			name: "get_key2",
			check: func() bool {
				ent := f.NewEntry("key2", "value2")
				got := f.Get("key2")
				return reflect.DeepEqual(ent, got)
			},
		},
		{
			name: "evict",
			check: func() bool {
				originalLen := f.queue.Len()
				oriSize := len(f.cache)
				evictEl := f.queue.Front()
				ent := evictEl.Value.(*entryFIFO)
				evicted := f.Evict()
				_, ok := f.cache[ent.key]

				return evicted == 10 && f.queue.Len() == originalLen-1 && len(f.cache) == oriSize-1 && !ok
			},
		},
		{
			name: "get_key1",
			check: func() bool {
				got := f.Get("key1")
				return got == nil
			},
		},
		{
			name: "get_key2",
			check: func() bool {
				ent := f.NewEntry("key2", "value2")
				got := f.Get("key2")
				return reflect.DeepEqual(ent, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Errorf("FIFO operations failed")
			}
		})
	}
}
