package engine

import "hash/crc32"

type Engine struct {
	shards     	[]*Shard
	numShards  	uint32
	ps 			*PubSubManager
}

// New creates a new sharded engine with the given number of shards.
func New(numShards uint32) *Engine {
	shards := make([]*Shard, numShards)
	for i := uint32(0); i < numShards; i++ {
		shards[i] = newShard()
	}
	return &Engine{
		shards:    	shards,
		numShards: 	numShards,
		ps: 		NewPubSubManager(),
	}
}

// Set sets a key to a value.
func (e *Engine) Set(key string, value []byte) {
	shard := e.getShard(key)
	shard.set(key, value)

	e.ps.Publish(key, []byte("SET " + key + " " + string(value)))
}

// Get retrieves a key's value.
func (e *Engine) Get(key string) ([]byte, bool) {
	shard := e.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	val, ok := shard.get(key)
	return val, ok
}

// Delete removes a key.
func (e *Engine) Delete(key string) {
	shard := e.getShard(key)
	delete(shard.data, key)

	e.ps.Publish(key, []byte("DEL " + key))
}

// getShard returns the shard for a given key.
func (e *Engine) getShard(key string) *Shard {
	hash := crc32.ChecksumIEEE([]byte(key))
	return e.shards[hash%uint32(e.numShards)]
}

// Access to PubSub system
func (e *Engine) PubSub() *PubSubManager {
	return e.ps
}