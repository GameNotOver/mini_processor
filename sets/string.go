package sets

import (
	"encoding/json"
	"fmt"
)

// String .
type String map[string]struct{}

// NewString .
func NewString(items ...string) String {
	sets := String{}
	sets.Add(items...)
	return sets
}

// Add .
func (s String) Add(items ...string) {
	for _, item := range items {
		s.add(item)
	}
}

func (s String) add(item string) {
	s[item] = struct{}{}
}

// Remove .
func (s String) Remove(items ...string) {
	for _, item := range items {
		s.remove(item)
	}
}

func (s String) remove(item string) {
	_, ok := s[item]
	if ok {
		delete(s, item)
	}
}

// Contains .
func (s String) Contains(item string) bool {
	_, ok := s[item]
	return ok
}

// Sub .
func (s String) Sub(ss String) String {
	out := String{}
	for item := range s {
		if !ss.Contains(item) {
			out.add(item)
		}
	}
	return out
}

// Union .
func (s String) Union(ss String) String {
	out := String{}
	for item := range s {
		out.add(item)
	}
	for item := range ss {
		out.add(item)
	}
	return out
}

// Intersect .
func (s String) Intersect(ss String) String {
	return s.Sub(s.Sub(ss))
}

// Members .
func (s String) Members() []string {
	items := make([]string, 0, len(s))
	for item := range s {
		items = append(items, item)
	}
	return items
}

// Equal .
func (s String) Equal(ss String) bool {
	for item := range s {
		if !ss.Contains(item) {
			return false
		}
	}
	for item := range ss {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

func (s String) String() string {
	return fmt.Sprintf("set%v", s.Members())
}

func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Members())
}

func (s *String) UnmarshalJSON(data []byte) error {
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	*s = NewString(list...)
	return nil
}

func (s String) MarshalYAML() (interface{}, error) {
	return s.Members(), nil
}

func (s *String) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var list []string
	if err := unmarshal(&list); err != nil {
		return err
	}
	*s = NewString(list...)
	return nil
}
