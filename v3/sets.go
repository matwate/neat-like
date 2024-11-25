package main

// This i will at some point export it to a separate package since it will be a generic set

type Set[T comparable] struct {
	container map[T]bool
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{container: make(map[T]bool)}
}

func (s Set[T]) Add(value T) {
	s.container[value] = true
}

func (s Set[T]) Remove(value T) {
	delete(s.container, value)
}

func (s Set[T]) Contains(value T) bool {
	_, ok := s.container[value]
	return ok
}

func (s Set[T]) Len() int {
	return len(s.container)
}

func (s Set[T]) Get(v T) T {
	if s.Contains(v) {
		return v
	} else {
		panic("Value not found")
	}
}

func (s Set[T]) Front() T {
	for k := range s.container {
		return k
	}
	panic("Empty set")
}

func (s Set[T]) Back() T {
	var last T
	for k := range s.container {
		last = k
	}
	return last
}

func (s Set[T]) Pop_back() T {
	last := s.Back()
	s.Remove(last)
	return last
}

// Union and Intersection will do later.
