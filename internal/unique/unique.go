package unique

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Exists(key K) (V, bool) {
	value, exists := m[key]
	return value, exists
}

// Add adds the value for the given key, if it does not exist already.
//
// Returns whether the value was added.
func (m Map[K, V]) Add(key K, value V) bool {
	_, exists := m.Exists(key)
	if !exists {
		m[key] = value
	}

	return exists
}
