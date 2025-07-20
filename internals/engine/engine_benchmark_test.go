package engine

import (
	"strconv"
	"testing"
)

func setupEngine() *Engine {
	return New(16)
}

func BenchmarkSet(b *testing.B) {
	engine := setupEngine()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Set("key"+strconv.Itoa(i), []byte("value"))
	}
}

func BenchmarkGet(b *testing.B) {
	engine := setupEngine()
	for i := 0; i < b.N; i++ {
		engine.Set("key"+strconv.Itoa(i), []byte("value"))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Get("key" + strconv.Itoa(i))
	}
}

func BenchmarkSetParallel(b *testing.B) {
	engine := setupEngine()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			engine.Set("pkey"+strconv.Itoa(i), []byte("value"))
			i++
		}
	})
}

func BenchmarkPubSubSet(b *testing.B) {
	engine := setupEngine()
	sub := engine.PubSub().Subscribe("benchkey", "test-client")

	go func() {
		for range sub.Chan {
			// simulate receiving
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Set("benchkey", []byte("val"+strconv.Itoa(i)))
	}
}
