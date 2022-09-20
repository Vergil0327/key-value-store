package provider

import (
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
			want: &FIFO{queue: make([]*entryFIFO, 0), cache: make(map[string]*entryFIFO)},
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
		queue []*entryFIFO
		cache map[string]*entryFIFO
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
		queue []*entryFIFO
		cache map[string]*entryFIFO
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CacheEntry
	}{
		{
			name: "it_should_get_from_cache",
			fields: fields{
				queue: []*entryFIFO{{key: "key1", value: "value"}},
				cache: map[string]*entryFIFO{"key1": {key: "key1", value: "value"}},
			},
			args: args{key: "key1"},
			want: &entryFIFO{key: "key1", value: "value"},
		},
		{
			name: "it_should_get_nil",
			fields: fields{
				queue: []*entryFIFO{{key: "key1", value: "value"}},
				cache: map[string]*entryFIFO{"key1": {key: "key1", value: "value"}},
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
		queue []*entryFIFO
		cache map[string]*entryFIFO
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
			fields: fields{queue: make([]*entryFIFO, 0), cache: make(map[string]*entryFIFO)},
			args:   args{key: "key1", value: &entryFIFO{key: "key1", value: "value1"}},
		},
		{
			name:   "it_should_add_entry_to_queue",
			fields: fields{queue: make([]*entryFIFO, 0), cache: make(map[string]*entryFIFO)},
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
			if v, ok := f.cache[tt.args.key]; !ok || !reflect.DeepEqual(v, want) {
				t.Errorf("After FIFO.Set(): got %v, want %v", v, want)
			}
			if v := f.queue[0]; !reflect.DeepEqual(v, want) {
				t.Errorf("After FIFO.Set(): got %v, want %v", v, want)
			}
		})
	}
}

func TestFIFO_Peek(t *testing.T) {
	type fields struct {
		queue []*entryFIFO
		cache map[string]*entryFIFO
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CacheEntry
	}{
		{
			name: "it_should_peek_from_cache",
			fields: fields{
				queue: []*entryFIFO{{key: "key1", value: "value"}},
				cache: map[string]*entryFIFO{"key1": {key: "key1", value: "value"}},
			},
			args: args{key: "key1"},
			want: &entryFIFO{key: "key1", value: "value"},
		},
		{
			name: "it_should_peek_nil",
			fields: fields{
				queue: []*entryFIFO{{key: "key1", value: "value"}},
				cache: map[string]*entryFIFO{"key1": {key: "key1", value: "value"}},
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
		queue []*entryFIFO
		cache map[string]*entryFIFO
	}
	tests := []struct {
		name        string
		fields      fields
		wantEvicted uint
	}{
		{
			name:        "it_should_evict_nothing",
			fields:      fields{queue: make([]*entryFIFO, 0), cache: make(map[string]*entryFIFO)},
			wantEvicted: 0,
		},
		{
			name:        "it_should_evict_successfully",
			fields:      fields{queue: []*entryFIFO{{key: "6bytes", value: "6bytes"}}, cache: make(map[string]*entryFIFO)},
			wantEvicted: 12,
		},
		{
			name: "it_should_evict_left_most",
			fields: fields{queue: []*entryFIFO{
				{key: "6bytes", value: "6bytes"},
				{key: "hello", value: "world"},
			}, cache: make(map[string]*entryFIFO)},
			wantEvicted: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FIFO{
				queue: tt.fields.queue,
				cache: tt.fields.cache,
			}
			original := len(f.queue)

			if gotEvicted := f.Evict(); gotEvicted != tt.wantEvicted {
				t.Errorf("FIFO.Evict() = %v, want %v", gotEvicted, tt.wantEvicted)
				if want := int(math.Max(0, float64(original))); len(f.queue) != want {
					t.Errorf("After FIFO.Evict() queue's length = %v, want %v", len(f.queue), want)
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
				v, ok := f.cache["key1"]
				return ok && len(f.queue) == 1 && reflect.DeepEqual(ent, v) && reflect.DeepEqual(ent, f.queue[0]) && ent.Size() == 10
			},
		},
		{
			name: "set_key2",
			check: func() bool {
				ent := f.NewEntry("key2", "value2")
				f.Set("key2", ent)
				v, ok := f.cache["key2"]
				return ok && len(f.queue) == 2 && reflect.DeepEqual(ent, v) && reflect.DeepEqual(ent, f.queue[1])
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
				originalLen := len(f.queue)
				oriSize := len(f.cache)
				evictEnt := f.queue[0]
				evicted := f.Evict()
				_, ok := f.cache[evictEnt.key]

				return evicted == 10 && len(f.queue) == originalLen-1 && len(f.cache) == oriSize-1 && !ok
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
