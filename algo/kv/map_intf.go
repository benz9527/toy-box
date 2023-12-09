package kv

import (
	"io"
)

type SafeStoreKeyFilterFunc func(key string) bool

func defaultAllKeysFilter(key string) bool {
	return true
}

type IThreadSafeStore[T any] interface {
	io.Closer
	AddOrUpdate(key string, obj T)
	Replace(items map[string]T)
	Delete(key string)
	Get(key string) (item T, exists bool)
	ListKeys(filters ...SafeStoreKeyFilterFunc) []string
	ListValues(keys ...string) (items []T)
}
