package common

type ReadOnlyMap[K comparable, V any] struct {
	data map[K]V
}

func NewReadOnlyMap[K comparable, V any](data map[K]V) *ReadOnlyMap[K, V] {
	return &ReadOnlyMap[K, V]{data: data}
}

func (m *ReadOnlyMap[K, V]) Get(key K) (V, bool) {
	value, ok := m.data[key]
	return value, ok
}

func (m *ReadOnlyMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys
}

func (m *ReadOnlyMap[K, V]) Len() int {
	return len(m.data)
}
