package cache

import (
	"kvstore/cache/provider"
	"reflect"
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
			if got := c.init(tt.args.config); !reflect.DeepEqual(got, tt.want) {
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
		// {
		// 	name: "it_should_construct_cache_with_given_provider",
		// 	args: args{
		// 		config: Config{StorageSize: 100, Cacher: provider.NewLRU()},
		// 	},
		// 	want: &Cache{
		// 		upperbound:     100,
		// 		provider: provider.NewLRU(),
		// 	},
		// 	wantErr: false,
		// },
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
		mux      sync.Mutex
		size     uint
		provider provider.CacheProvider
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   provider.SizeGetter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				mux:      tt.fields.mux,
				size:     tt.fields.size,
				provider: tt.fields.provider,
			}
			if got := c.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
