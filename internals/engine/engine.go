package engine

import "hash/crc32"

// ChangeHook is a callback to notify on data updates (e.g. pubsub system).
type ChangeHook func(key string, value []byte, op string)

type Engine struct {
	shards     []*Shard
	numShards  uint32
	onChange   ChangeHook
}

// New creates a new sharded engine with the given number of shards.
func New(numShards uint32, hook ChangeHook) *Engine {
	shards := make([]*Shard, numShards)
	for i := uint32(0); i < numShards; i++ {
		shards[i] = newShard()
	}
	return &Engine{
		shards:    shards,
		numShards: numShards,
		onChange:  hook,
	}
}

// Set sets a key to a value.
func (e *Engine) Set(key string, value []byte) {
	shard := e.getShard(key)
	shard.set(key, value)
	if e.onChange != nil {
		e.onChange(key, value, "set")
	}
}

// Get retrieves a key's value.
func (e *Engine) Get(key string) ([]byte, bool) {
	return e.getShard(key).get(key)
}

// Delete removes a key.
func (e *Engine) Delete(key string) {
	shard := e.getShard(key)
	shard.delete(key)
	if e.onChange != nil {
		e.onChange(key, nil, "delete")
	}
}

// getShard returns the shard for a given key.
func (e *Engine) getShard(key string) *Shard {
	hash := crc32.ChecksumIEEE([]byte(key))
	return e.shards[hash%uint32(e.numShards)]
}