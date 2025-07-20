package engine

import (
	"strconv"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	db := New(256, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := []byte("value")
		db.Set(key, value)
	}
}

func BenchmarkGet(b *testing.B) {
	db := New(256, nil)

	// Preload with some keys
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		value := []byte("value")
		db.Set(key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		db.Get(key)
	}
}
