package provider

import (
	"container/list"
	"reflect"
	"testing"
)

func Test_entryLRU_Key(t *testing.T) {
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
			ent := entryLRU{
				key:   tt.fields.key,
				value: tt.fields.value,
			}
			if got := ent.Key(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("entryLRU.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryLRU_Value(t *testing.T) {
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
			ent := entryLRU{
				key:   tt.fields.key,
				value: tt.fields.value,
			}
			if got := ent.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("entryLRU.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryLRU_Size(t *testing.T) {
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
			ent := entryLRU{
				key:   tt.fields.key,
				value: tt.fields.value,
			}
			if got := ent.Size(); got != tt.want {
				t.Errorf("entryLRU.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLRU(t *testing.T) {
	tests := []struct {
		name string
		want *LRU
	}{
		{
			name: "it_should_construct_lru_correctly",
			want: &LRU{list: list.New(), cache: make(map[string]*list.Element)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLRU(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLRU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLRU_NewEntry(t *testing.T) {
	type fields struct {
		list  *list.List
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
			want: &entryLRU{key: "KEY", value: "VALUE"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &LRU{
				list:  tt.fields.list,
				cache: tt.fields.cache,
			}
			if got := f.NewEntry(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LRU.NewEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLRU_Get(t *testing.T) {
	type fields struct {
		list  *list.List
		cache map[string]*list.Element
	}
	type args struct {
		key string
	}
	ent := &entryLRU{key: "key1", value: "value"}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CacheEntry
	}{
		{
			name: "it_should_get_from_cache",
			fields: fields{
				list:  list.New(),
				cache: map[string]*list.Element{"key1": {Value: ent}},
			},
			args: args{key: "key1"},
			want: &entryLRU{key: "key1", value: "value"},
		},
		{
			name: "it_should_get_nil",
			fields: fields{
				list:  list.New(),
				cache: make(map[string]*list.Element),
			},
			args: args{key: "key2"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &LRU{
				list:  tt.fields.list,
				cache: tt.fields.cache,
			}
			if got := f.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LRU.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLRU_Set(t *testing.T) {
	type fields struct {
		list  *list.List
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
			fields: fields{list: list.New(), cache: make(map[string]*list.Element)},
			args:   args{key: "key1", value: &entryLRU{key: "key1", value: "value1"}},
		},
		{
			name:   "it_should_add_entry_to_list",
			fields: fields{list: list.New(), cache: make(map[string]*list.Element)},
			args:   args{key: "key2", value: &entryLRU{key: "key2", value: "value2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &LRU{
				list:  tt.fields.list,
				cache: tt.fields.cache,
			}
			f.Set(tt.args.key, tt.args.value)
			want := &entryLRU{key: tt.args.key, value: tt.args.value.Value()}
			if v, ok := f.cache[tt.args.key].Value.(*entryLRU); !ok || !reflect.DeepEqual(v, want) {
				t.Errorf("After LRU.Set(): got %v, want %v", v, want)

				if el := f.list.Front(); !reflect.DeepEqual(el, v) {
					t.Errorf("After LRU.Set(): got %v, want %v", el, v)
				}
			}
		})
	}
}

func TestLRU_Peek(t *testing.T) {
	type fields struct {
		list  *list.List
		cache map[string]*list.Element
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
			name: "it_should_peek_in_cache",
			fields: fields{
				list:  list.New(),
				cache: map[string]*list.Element{"key1": {Value: &entryLRU{key: "key1", value: "value"}}},
			},
			args: args{key: "key1"},
			want: &entryLRU{key: "key1", value: "value"},
		},
		{
			name: "it_should_peek_at_nil",
			fields: fields{
				list:  list.New(),
				cache: make(map[string]*list.Element),
			},
			args: args{key: "key2"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &LRU{
				list:  tt.fields.list,
				cache: tt.fields.cache,
			}
			if got := f.Peek(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LRU.Peek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLRU_Evict(t *testing.T) {
	type fields struct {
		list  *list.List
		cache map[string]*list.Element
	}

	l := list.New()
	m := map[string]*list.Element{}
	ent := &entryLRU{key: "key1", value: "value"}
	el := l.PushFront(ent)
	m[ent.key] = el

	l2 := list.New()
	m2 := map[string]*list.Element{}
	ent1 := &entryLRU{key: "key1", value: "total16bytes"}
	el1 := l2.PushFront(ent1)
	ent2 := &entryLRU{key: "hello", value: "world"}
	el2 := l2.PushFront(ent2)
	m2[ent1.key] = el1
	m2[ent2.key] = el2

	tests := []struct {
		name        string
		fields      fields
		wantEvicted uint
	}{
		{
			name:        "it_should_evict_nothing",
			fields:      fields{list: list.New(), cache: make(map[string]*list.Element)},
			wantEvicted: 0,
		},
		{
			name:        "it_should_evict_successfully",
			fields:      fields{list: l, cache: m},
			wantEvicted: 9,
		},
		{
			name:        "it_should_evict_oldest_one",
			fields:      fields{list: l2, cache: m2},
			wantEvicted: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &LRU{
				list:  tt.fields.list,
				cache: tt.fields.cache,
			}
			if gotEvicted := f.Evict(); gotEvicted != tt.wantEvicted {
				t.Errorf("LRU.Evict() = %v, want %v", gotEvicted, tt.wantEvicted)
			}
		})
	}
}

func TestLRU(t *testing.T) {
	f := NewLRU()

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
				want := f.NewEntry("key1", "value1")

				f.Set("key1", want)
				el, ok := f.cache["key1"]
				return ok && f.list.Len() == 1 && reflect.DeepEqual(want, el.Value) && reflect.DeepEqual(want, f.list.Front().Value) && want.Size() == 10
			},
		},
		{
			name: "set_key2",
			check: func() bool {
				want := f.NewEntry("key2", "long_value2")

				f.Set("key2", want)
				el, ok := f.cache["key2"]
				return ok && f.list.Len() == 2 && reflect.DeepEqual(want, el.Value) && reflect.DeepEqual(want, f.list.Front().Value)
			},
		},
		{
			name: "get_key1",
			check: func() bool {
				want := f.NewEntry("key1", "value1")
				got := f.Get("key1")
				return reflect.DeepEqual(want, got)
			},
		},
		{
			name: "peek_key1",
			check: func() bool {
				want := f.NewEntry("key1", "value1")
				got := f.Peek("key1")
				return reflect.DeepEqual(want, got)
			},
		},
		{
			name: "get_key2",
			check: func() bool {
				want := f.NewEntry("key2", "long_value2")
				got := f.Get("key2")
				return reflect.DeepEqual(want, got)
			},
		},
		{
			name: "get_key2",
			check: func() bool {
				want := f.NewEntry("key2", "long_value2")
				got := f.Get("key2")
				return reflect.DeepEqual(want, got)
			},
		},
		{
			name: "evict",
			check: func() bool {
				originalLen := f.list.Len()
				oriSize := len(f.cache)
				evictEnt := f.list.Back()
				evicted := f.Evict()
				_, ok := f.cache[evictEnt.Value.(*entryLRU).key]

				// evict key1
				return evicted == 10 && f.list.Len() == originalLen-1 && len(f.cache) == oriSize-1 && !ok
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
				want := f.NewEntry("key2", "long_value2")
				got := f.Get("key2")
				return reflect.DeepEqual(want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Errorf("LRU operations failed")
			}
		})
	}
}
