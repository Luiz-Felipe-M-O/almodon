package repository

type Index[K, V comparable] map[K]V

func (i Index[K, V]) Get(key K) (V, bool) {
	val, ok := i[key]
	return val, ok
}

func (i Index[K, V]) Set(key K, value V) {
	i[key] = value
}

func (i Index[K, V]) Remove(key K) {
	delete(i, key)
}

type SliceIndex[K, V comparable] map[K][]V

func (s SliceIndex[K, V]) Get(key K) ([]V, bool) {
	values, ok := s[key]
	return values, ok
}

func (s SliceIndex[K, V]) Add(key K, value V) {
	s[key] = append(s[key], value)
}

func (s SliceIndex[K, V]) Swap(key K, old, new V) {
	values, ok := s[key]
	if !ok {
		return
	}

	for i, v := range values {
		if v == old {
			s[key][i] = new
			break
		}
	}
}

func (s SliceIndex[K, V]) Del(key K, value V) {
	values, ok := s[key]
	if !ok {
		return
	}

	for i, v := range values {
		if v == value {
			s[key] = append(values[:i], values[i+1:]...)
			break
		}
	}

	if len(s[key]) == 0 {
		delete(s, key)
	}
}
