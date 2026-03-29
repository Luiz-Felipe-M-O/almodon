package store

type Index[K comparable] map[K]int

func (i Index[K]) Get(key K) (int, bool) {
	val, ok := i[key]
	return val, ok
}

func (i Index[K]) Set(key K, index int) {
	i[key] = index
}

func (i Index[K]) Del(key K) {
	delete(i, key)
}

type MultiIndex[K comparable] map[K][]int

func (s MultiIndex[K]) Get(key K) ([]int, bool) {
	values, ok := s[key]
	return values, ok
}

func (s MultiIndex[K]) Add(key K, index int) {
	s[key] = append(s[key], index)
}

func (s MultiIndex[K]) Swap(key K, old, new int) {
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

func (s MultiIndex[K]) Del(key K, value int) {
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
