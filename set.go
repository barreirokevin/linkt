package main

// The keys of this map represent a Set, i.e. no duplicate values.
type Set[K comparable, V any] map[K]V

// Returns true if the set contains value, otherwise it returns false.
func (s *Set[K, V]) Contains(value K) bool {
	_, found := (*s)[value]
	if found {
		return true
	}
	return false
}
