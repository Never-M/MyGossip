package set

import (
	"fmt"
	"github.com/String-Reconciliation-Ditributed-System/RCDS_GO/pkg/lib/algorithm"
	"reflect"
)

type Set map[interface{}]interface{}

// Create a new set
func New() *Set {
	return &Set{}
}

// Find the difference between two sets (s - set)
func (s *Set) Difference(set *Set) *Set {
	n := make(Set)

	for k := range *s {
		if v, exists := (*set)[k]; !exists {
			n[k] = v
		}
	}

	return &n
}

// Call f for each item in the set
func (s *Set) Do(f func(interface{})) {
	for k := range *s {
		f(k)
	}
}

// Test to see whether or not the element is in the set
func (s *Set) Has(key interface{}) bool {
	key = convertToHashableElement(key)
	_, exists := (*s)[key]
	return exists
}

func (s *Set) Get(key interface{}) interface{} {
	key = convertToHashableElement(key)
	return (*s)[key]
}

// Add an element to the set, val can be struct{}
func (s *Set) Insert(key, val interface{}) {
	if val == nil {
		val = struct{}{}
	}
	key = convertToHashableElement(key)
	(*s)[key] = val
}

func (s *Set) InsertKey(key interface{}) {
	key = convertToHashableElement(key)
	(*s)[key] = struct{}{}
}

// Find the intersection of two sets
func (s *Set) Intersection(otherSet *Set) *Set {
	n := make(Set)

	for k := range *s {
		if val, exists := (*otherSet)[k]; exists {
			n[k] = val
		}
	}

	return &n
}

// Return the number of items in the set
func (s *Set) Len() int {
	return len(*s)
}

// Test whether or not this set is a proper subset of "set"
func (s *Set) ProperSubsetOf(set *Set) bool {
	return s.SubsetOf(set) && s.Len() < set.Len()
}

// Remove an element from the set
func (s *Set) Remove(key interface{}) {
	key = convertToHashableElement(key)
	delete(*s, key)
}

// Test whether or not this set is a subset of "set"
func (s *Set) SubsetOf(set *Set) bool {
	if s.Len() > set.Len() {
		return false
	}
	for k := range *s {
		if _, exists := (*set)[k]; !exists {
			return false
		}
	}
	return true
}

// Find the union of two sets
func (s *Set) Union(set *Set) *Set {
	n := make(Set)

	for k, v := range *s {
		n[k] = v
	}
	for k, v := range *set {
		n[k] = v
	}

	return &n
}

// Digest is the xor sum of the entire set's hash.
func (s *Set) GetDigest() (uint64, error) {
	var sum uint64
	for k, v := range *s {
		e, err := algorithm.HashString(fmt.Sprint(k) + fmt.Sprint(v)).ToUint64()
		if err != nil {
			return 0, err
		}
		sum = sum ^ e
	}
	return sum, nil
}

// convertToHashableElement includes all methods converting element into a hashable type:
// 1. convert []byte to string.
func convertToHashableElement(element interface{}) interface{} {
	if reflect.TypeOf(element).String() == reflect.TypeOf([]byte{}).String() {
		element = string(element.([]byte))
	}
	return element
}
