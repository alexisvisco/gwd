package utils

import (
	"encoding/json"
)

// StringSet is a map from a string to an empty struct. We choose empty structs because they have a size
// of zero and make it fairly clear that this shouldn't be treated as a map.
type StringSet map[string]struct{}

// NewStringSet creates a new stringset with the specified strings in it
func NewStringSet(strings ...string) StringSet {
	set := make(map[string]struct{}, len(strings))
	for _, str := range strings {
		set[str] = struct{}{}
	}
	return set
}

// FromList converts a list of strings a StringSet
func FromList(strings []string) StringSet {
	return NewStringSet(strings...)
}

// FromInterfaceList converts a list of interfaces that are known to be strings into a StringSet
func FromInterfaceList(strings []interface{}) StringSet {
	set := make(map[string]struct{}, len(strings))
	for _, str := range strings {
		set[str.(string)] = struct{}{}
	}
	return set
}

// ToList converts a StringSet to a list of strings
func (s StringSet) ToList() []string {
	returnList := make([]string, 0, len(s))
	for key, _ := range s {
		returnList = append(returnList, key)
	}
	return returnList
}

// Clone copies a string set to a new string set
func (s StringSet) Clone() StringSet {
	returnSet := make(map[string]struct{}, len(s))
	for key, value := range s {
		returnSet[key] = value
	}
	return returnSet
}

// Intersect returns a new StringSet with the intersection of all the elements in both sets
func (s StringSet) Intersect(s2 StringSet) StringSet {
	return setOperation(s, s2, true)
}

// Minus returns a new StringSet with the
func (s StringSet) Minus(s2 StringSet) StringSet {
	return setOperation(s, s2, false)
}

// setOperation is a helper method to either intersect or subtract sets
func setOperation(s1, s2 StringSet, wantElemsInSet2 bool) map[string]struct{} {
	resultSet := make(map[string]struct{})
	for key, _ := range s1 {
		if _, ok := s2[key]; ok == wantElemsInSet2 {
			resultSet[key] = struct{}{}
		}
	}
	return resultSet
}

// AddSet adds all the elements in a string set to the operand set.
func (s StringSet) AddSet(newValues StringSet) {
	for newValue, _ := range newValues {
		s[newValue] = struct{}{}
	}
}

// AddAll adds all the elements in a string slice to the operand set.
func (s StringSet) AddAll(newValues []string) {
	for _, newValue := range newValues {
		s[newValue] = struct{}{}
	}
}

// Equals returns true if two string sets have exactly the same elements
func (s StringSet) Equals(s2 StringSet) bool {
	if len(s) != len(s2) {
		return false
	}
	for key, _ := range s {
		if _, ok := s2[key]; !ok {
			return false
		}
	}
	return true
}

// Add adds an element to a string set
func (s StringSet) Add(str string) {
	s[str] = struct{}{}
}

// Remove removes an element from a string set
func (s StringSet) Remove(str string) {
	delete(s, str)
}

// Contains returns true if a stringset contains the specified string
func (s StringSet) Contains(str string) bool {
	_, ok := s[str]
	return ok
}

// Partition takes in two string slices and returns a tuple with (strings only in the first set,
// strings in both sets, strings only in the second set). It is a utility function that uses the
// set implementation
func Partition(s1, s2 []string) (only1 []string, both []string, only2 []string) {
	set1 := FromList(s1)
	set2 := FromList(s2)
	return set1.Minus(set2).ToList(), set1.Intersect(set2).ToList(), set2.Minus(set1).ToList()
}

// MarshalJSON marshals a string set to JSON
func (s StringSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToList())
}
