package cache

import "errors"

var (
	ErrStorageLimit = errors.New("storage size should be greater than zero")
)
