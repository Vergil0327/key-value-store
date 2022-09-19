package provider

import "reflect"

// return size in bytes
func calculateSize(items []any) uint {
	var size uint
	for _, item := range items {
		switch v := item.(type) {
		case string:
			size += uint(len(v))
		default:
			size += uint(reflect.TypeOf(item).Size())
		}
	}
	return size
}
