package keyvaluestore

import (
	"sync"
	"testing"
)

func TestCache_calculateSize(t *testing.T) {
	type fields struct {
		size uint
	}
	type args struct {
		items []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint
	}{
		{
			name: "it_should_work",
			args: args{items: []any{"key", "value"}},
			want: 8,
		},
		{
			name: "size_should_be_sum_of_number_of_bytes",
			args: args{items: []any{"key0123456789", "hello world"}},
			want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cache{
				mux:  sync.Mutex{},
				size: tt.fields.size,
			}
			if got := c.calculateSize(tt.args.items); got != tt.want {
				t.Errorf("Cache.calculateSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
