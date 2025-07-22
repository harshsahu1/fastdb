// internals/engine/shard.go

package engine

import (
	"sync"
)

type Shard struct {
	data sync.Map
}

func newShard() *Shard {
	return &Shard{}
}

func (s *Shard) set(key string, value []byte) {
	s.data.Store(key, value)
}

func (s *Shard) get(key string) ([]byte, bool) {
	val, ok := s.data.Load(key)
	if !ok {
		return nil, false
	}
	return val.([]byte), true
}

func (s *Shard) delete(key string) {
	s.data.Delete(key)
}