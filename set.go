package main

// The keys of this map represent a Set, i.e. no duplicate values.
type Set map[string]int

// Returns true if the set contains value, otherwise it returns false.
func (s *Set) Contains(value string) bool {
	_, found := (*s)[value]
	if found {
		return true
	}
	return false
}
