package utils

type StringSet interface {
	Add(elements ...string) StringSet
	Remove(element string) StringSet
	Contains(element string) bool
	Length() int
	Keys() []string
}

func NewSet(element ...string) StringSet {
	set := &stringSet{
		set: make(map[string]struct{}),
	}

	for _, el := range element {
		set.Add(el)
	}

	return set
}

type stringSet struct {
	set map[string]struct{}
}

func (s *stringSet) Add(elements ...string) StringSet {
	for _, element := range elements {
		s.set[element] = struct{}{}
	}
	return s
}

func (s *stringSet) Remove(element string) StringSet {
	delete(s.set, element)
	return s
}

func (s *stringSet) Contains(element string) bool {
	_, ok := s.set[element]
	return ok
}

func (s *stringSet) Length() int {
	return len(s.set)
}

func (s *stringSet) Keys() []string {
	keys := make([]string, s.Length())
	index := 0
	for key := range s.set {
		keys[index] = key
		index++
	}

	return keys
}
