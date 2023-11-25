package container

import (
	"io"
	"log/slog"
	"reflect"
	"strings"
	"sync"
)

type threadSafeMap[T any] struct {
	lock  sync.RWMutex
	items map[string]T
}

func (t *threadSafeMap[T]) AddOrUpdate(key string, obj T) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.items[key] = obj
}

func (t *threadSafeMap[T]) Replace(items map[string]T) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.items = items
}

func (t *threadSafeMap[T]) Delete(key string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, exists := t.items[key]; exists {
		delete(t.items, key)
	}
}

func (t *threadSafeMap[T]) Get(key string) (item T, exists bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	item, exists = t.items[key]
	return
}

func (t *threadSafeMap[T]) ListKeys(filters ...SafeStoreKeyFilterFunc) []string {

	realFilters := make([]SafeStoreKeyFilterFunc, 0, len(filters))
	for _, filter := range filters {
		if filter != nil {
			realFilters = append(realFilters, filter)
		}
	}
	if len(realFilters) == 0 {
		realFilters = append(realFilters, defaultAllKeysFilter)
	}

	t.lock.RLock()
	defer t.lock.RUnlock()

	keys := make([]string, 0, len(t.items))
	for key := range t.items {
		for _, filter := range realFilters {
			if filter(key) {
				keys = append(keys, key)
				break
			}
		}
	}
	return keys
}

func (t *threadSafeMap[T]) ListValues(keys ...string) (items []T) {
	realKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		if len(strings.TrimSpace(key)) > 0 {
			realKeys = append(realKeys, key)
		}
	}

	contains := func(keys []string, key string) bool {
		for _, k := range keys {
			if k == key {
				return true
			}
		}
		return false
	}

	t.lock.RLock()
	defer t.lock.RUnlock()
	values := make([]T, 0, len(t.items))
	for key, item := range t.items {
		i := item
		if len(realKeys) > 0 && contains(realKeys, key) {
			values = append(values, i)
		} else if len(realKeys) == 0 {
			values = append(values, i)
		}
	}
	return values
}

func (t *threadSafeMap[T]) Close() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	for _, item := range t.items {
		if item == nil {
			continue
		}

		typ := reflect.TypeOf(item)
		if typ.Implements(reflect.TypeOf((*io.Closer)(nil)).Elem()) {
			vals := reflect.ValueOf(item).MethodByName("Close").Call([]reflect.Value{})
			if len(vals) > 0 && !vals[0].IsNil() {
				intf := vals[0].Elem().Interface()
				switch intf.(type) {
				case error:
					err := intf.(error)
					slog.Error("Close info", "error", err)
				}
			}
		}
	}

	t.items = nil
	return nil
}

func NewThreadSafeMap[T any]() IThreadSafeStore[T] {
	return &threadSafeMap[T]{items: make(map[string]T)}
}
