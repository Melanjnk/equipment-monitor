package controller

import (
	"maps"
	"slices"
)

type StringSet struct {
	data map[string]struct{}
}

func NewStringSet(values ...string) StringSet {
	set := StringSet{data: make(map[string]struct{})}
	for _, value := range(values) {
		set.Include(value)
	}
	return set
}

func (set *StringSet) Include(value string) {
	set.data[value] = struct{}{}
}

func (set *StringSet) Exclude(value string) {
	delete (set.data, value)
}

func (set *StringSet) ExcludeMultiply(values ...string) {
	for _, value := range values {
		set.Exclude(value)
	}
}

func (set StringSet) IsEmpty() bool {
	return len(set.data) == 0
}

func (set StringSet) Contains(value string) bool {
	_, contains := set.data[value]
	return contains
}

func (set *StringSet) Slice() []string {
	return slices.Collect(maps.Keys(set.data))
}
