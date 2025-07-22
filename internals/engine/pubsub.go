package engine

import (
	"sync"
)

type Subscriber struct {
	ID   string
	Chan chan []byte
}

type PubSubManager struct {
	mu 			sync.RWMutex
	subscribers map[string][]*Subscriber
}

func NewPubSubManager() *PubSubManager {
	return &PubSubManager{
		subscribers: make(map[string][]*Subscriber),
	}
}

func (ps *PubSubManager) Subscribe(key string, id string) *Subscriber {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	sub := &Subscriber{
		ID:   id,
		Chan: make(chan []byte, 64), // buffered to avoid blocking
	}
	ps.subscribers[key] = append(ps.subscribers[key], sub)
	return sub
}


func (ps *PubSubManager) Unsubscribe(key string, id string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subs := ps.subscribers[key]
	for i, sub := range subs {
		if sub.ID == id {
			ps.subscribers[key] = append(subs[:i], subs[i+1:]...)
			close(sub.Chan)
			break
		}
	}
}

func (ps *PubSubManager) HasSubscribers(key string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return len(ps.subscribers[key]) > 0
}

func (ps *PubSubManager) Publish(key string, message []byte) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, sub := range ps.subscribers[key] {
		select {
		case sub.Chan <- message:
		default:
			// drop message if channel is full to avoid blocking
		}
	}
}