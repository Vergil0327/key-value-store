package cache

import (
	"fmt"
	"kvstore/cache/provider"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

func TestCache_init(t *testing.T) {
	type fields struct {
		mux      sync.Mutex
		size     uint
		provider provider.CacheProvider
	}
	type args struct {
		config Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Cache
	}{
		{
			name:   "it_should_init_cache_struct_correctly",
			fields: fields{},
			args:   args{config: Config{StorageSize: 50, Cacher: provider.NewFIFO()}},
			want:   &Cache{upperbound: 50, provider: provider.NewFIFO()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				mux:      tt.fields.mux,
				size:     tt.fields.size,
				provider: tt.fields.provider,
			}
			if got := c.initialize(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCache(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Cache
		wantErr bool
	}{
		{
			name: "it_should_construct_cache_with_default_provider",
			args: args{
				config: Config{StorageSize: 100},
			},
			want: &Cache{
				upperbound: 100,
				provider:   provider.NewFIFO(),
			},
			wantErr: false,
		},
		{
			name: "it_should_construct_cache_with_given_provider",
			args: args{
				config: Config{StorageSize: 100, Cacher: provider.NewLRU()},
			},
			want: &Cache{
				upperbound: 100,
				provider:   provider.NewLRU(),
			},
			wantErr: false,
		},
		{
			name: "storage_size_should_be_greater_than_0",
			args: args{
				config: Config{StorageSize: 0},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCache(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Get(t *testing.T) {
	type fields struct {
		mux        sync.Mutex
		upperbound uint
		size       uint
		provider   provider.CacheProvider
	}
	type args struct {
		key string
	}

	mock := provider.NewFIFO()
	mock.Set("key1", mock.NewEntry("key1", "value"))
	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
	}{
		{
			name:   "it_should_get_nothing",
			fields: fields{upperbound: 50, provider: provider.NewFIFO()},
			args:   args{key: "key1"},
			want:   nil,
		},
		{
			name:   "it_should_lookup_from_cache",
			fields: fields{upperbound: 50, provider: mock},
			args:   args{key: "key1"},
			want:   "value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				mux:        tt.fields.mux,
				upperbound: tt.fields.upperbound,
				size:       tt.fields.size,
				provider:   tt.fields.provider,
			}
			if got := c.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Set(t *testing.T) {
	type fields struct {
		mux        sync.Mutex
		upperbound uint
		size       uint
		provider   provider.CacheProvider
	}
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint
	}{
		{
			name:   "it_should_set_value_to_cache",
			fields: fields{upperbound: 50, provider: provider.NewLRU()},
			args:   args{key: "key1", value: "value1"},
			want:   0,
		},
		{
			name:   "it_should_evict_value",
			fields: fields{upperbound: 5, provider: provider.NewLRU()},
			args:   args{key: "key1", value: "value1"},
			want:   10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				mux:        tt.fields.mux,
				upperbound: tt.fields.upperbound,
				size:       tt.fields.size,
				provider:   tt.fields.provider,
			}
			if got := c.Set(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("Cache.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Peek(t *testing.T) {
	type fields struct {
		mux        sync.Mutex
		upperbound uint
		size       uint
		provider   provider.CacheProvider
	}
	type args struct {
		key string
	}

	mock := provider.NewFIFO()
	mock.Set("key1", mock.NewEntry("key1", "value"))

	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
	}{
		{
			name:   "it_should_peek_at_nothing",
			fields: fields{upperbound: 50, provider: provider.NewFIFO()},
			args:   args{key: "key1"},
			want:   nil,
		},
		{
			name:   "it_should_peek_in_cache",
			fields: fields{upperbound: 50, provider: mock},
			args:   args{key: "key1"},
			want:   "value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				mux:        tt.fields.mux,
				upperbound: tt.fields.upperbound,
				size:       tt.fields.size,
				provider:   tt.fields.provider,
			}
			if got := c.Peek(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Peek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Size(t *testing.T) {
	type fields struct {
		mux        sync.Mutex
		upperbound uint
		size       uint
		provider   provider.CacheProvider
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name:   "it_should_return_current_storage_size",
			fields: fields{size: 100},
			want:   100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				mux:        tt.fields.mux,
				upperbound: tt.fields.upperbound,
				size:       tt.fields.size,
				provider:   tt.fields.provider,
			}
			if got := c.Size(); got != tt.want {
				t.Errorf("Cache.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTreadSafe_SetNoEvict(t *testing.T) {
	c, err := NewCache(Config{StorageSize: 6780, Cacher: provider.NewLRU()})
	if err != nil {
		t.Errorf("Error occured:%s", err.Error())
	}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(idx int) {
			c.Set("k"+strconv.Itoa(idx), strconv.Itoa(idx))
			wg.Done()
		}(i)
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(idx int) {
			c.Get("k" + strconv.Itoa(idx))
			c.Peek("k" + strconv.Itoa(idx))
			wg.Done()
		}(i)
	}

	wg.Wait()
	var want uint = 6780
	if got := c.Size(); got != want {
		t.Errorf("Cache.Size() = %v, want %v", got, want)
	}
}

func TestTreadSafe_SetWithEvict(t *testing.T) {
	c, err := NewCache(Config{StorageSize: 3000, Cacher: provider.NewLRU()})
	if err != nil {
		t.Errorf("Error occured:%s", err.Error())
	}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(idx int) {
			c.Set(fmt.Sprintf("%03d", idx), "1")
			wg.Done()
		}(i)
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(idx int) {
			c.Get("k" + strconv.Itoa(idx))
			c.Peek("k" + strconv.Itoa(idx))
			wg.Done()
		}(i)
	}

	wg.Wait()
	var want uint = 3000
	if got := c.Size(); got != want {
		t.Errorf("Cache.Size() = %v, want %v", got, want)
	}
}
