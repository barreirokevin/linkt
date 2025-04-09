package main

import "fmt"

// A key-value pair.
type MapEntry[K any, V any] struct {
	key   *K
	value *V
}

// Creates a MapEntry with the given key-value pair.
func NewMapEntry[K any, V any](key *K, value *V) *MapEntry[K, V] {
	return &MapEntry[K, V]{key: key, value: value}
}

// Returns the key of the key-value pair.
func (e MapEntry[K, V]) GetKey() *K {
	return e.key
}

// Returns the value of the key-value pair.
func (e MapEntry[K, V]) GetValue() *V {
	return e.value
}

// Changes the key in this MapEntry to the given key.
func (e MapEntry[K, V]) SetKey(key *K) {
	e.key = key
}

// Changes the value associated to the key of this MapEntry and
// returns the old value.
func (e MapEntry[K, V]) SetValue(value *V) *V {
	old := e.value
	e.value = value
	return old
}

// Returns a string representation of MapEntry. This is useful
// for debugging.
func (e MapEntry[K, V]) ToString() string {
	return fmt.Sprintf("key: %+v, value: %+v\n", *e.key, *e.value)
}
