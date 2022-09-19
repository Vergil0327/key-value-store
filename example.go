package keyvaluestore

type KeyValueStore struct {
	store map[string]string
}

func New() KeyValueStore {
	return KeyValueStore{}
}

func (kvs *KeyValueStore) Get(key string) string {
	// if v, ok := cache[key]; ok {
	// 	return v
	// }

	if v := kvs.store[key]; v != "" {
		// read-through cache
		// cache[key] = v

		return v
	}

	return ""
}

func (kvs *KeyValueStore) Set(key string, value string) {
	kvs.store[key] = value

	// write-through cache
	// cache[key] = v
}
