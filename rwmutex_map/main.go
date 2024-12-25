package main

import "sync"

type SafeMap[K comparable, V any] struct {
	l sync.RWMutex
	m map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{m: make(map[K]V)}
}

func (s *SafeMap[K, V]) Set(key K, value V) {
	s.l.Lock()
	defer s.l.Unlock()
	s.m[key] = value
}

func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	value, ok := s.m[key]
	return value, ok
}

func (s *SafeMap[K, V]) Delete(key K) {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.m, key)
}

func (s *SafeMap[K, V]) Len() int {
	return len(s.m)
}
