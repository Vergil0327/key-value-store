package keyvaluestore

import "testing"

func Test_calculateSize(t *testing.T) {
	type args struct {
		items []any
	}
	tests := []struct {
		name string
		args args
		want uint
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
			if got := calculateSize(tt.args.items); got != tt.want {
				t.Errorf("calculateSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
