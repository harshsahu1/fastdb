package engine

import (
	"bytes"
	"hash/crc32"
)

type Engine struct {
	shards     	[]*Shard
	numShards  	uint32
	ps 			*PubSubManager
}

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

func (e *Engine) Set(key string, value []byte) {
	shard := e.getShard(key)
	shard.set(key, value)

	if e.ps.HasSubscribers(key) {
		msg := bytes.Join([][]byte{[]byte("SET"), []byte(key), value}, []byte(" "))
		e.ps.Publish(key, msg)
	}
}

func (e *Engine) Get(key string) ([]byte, bool) {
	shard := e.getShard(key)
	val, ok := shard.get(key)
	return val, ok
}

func (e *Engine) Delete(key string) {
	shard := e.getShard(key)
	shard.delete(key)

	if e.ps.HasSubscribers(key) {
		msg := bytes.Join([][]byte{[]byte("DEL"), []byte(key)}, []byte(" "))
		e.ps.Publish(key, msg)
	}
}

func (e *Engine) getShard(key string) *Shard {
	hash := crc32.ChecksumIEEE([]byte(key))
	return e.shards[hash%uint32(e.numShards)]
}

func (e *Engine) PubSub() *PubSubManager {
	return e.ps
}