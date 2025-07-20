package engine

import (
	"sync"
)

type Shard struct {
	mu sync.RWMutex
	data map[string][]byte
}

func newShard() *Shard {
	return &Shard{
		data: make(map[string][]byte),
	}
}

func (s *Shard) set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s * Shard) get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *Shard) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}