//go:generate go run github.com/99designs/gqlgen generate

package graph

import (
	"sync"

	"github.com/jsr-probitas/dockerfiles/echo-graphql/graph/model"
)

// Resolver is the root resolver for all GraphQL operations
type Resolver struct {
	mu              sync.RWMutex
	messages        map[string]*model.Message
	nextID          int
	messageChannels []chan *model.Message
}

// NewResolver creates a new resolver instance
func NewResolver() *Resolver {
	return &Resolver{
		messages: make(map[string]*model.Message),
		nextID:   1,
	}
}

// Subscribe adds a channel to receive message events
func (r *Resolver) Subscribe() chan *model.Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan *model.Message, 1)
	r.messageChannels = append(r.messageChannels, ch)
	return ch
}

// Unsubscribe removes a channel from message events
func (r *Resolver) Unsubscribe(ch chan *model.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, c := range r.messageChannels {
		if c == ch {
			r.messageChannels = append(r.messageChannels[:i], r.messageChannels[i+1:]...)
			close(ch)
			return
		}
	}
}

// Broadcast sends a message to all subscribers
func (r *Resolver) Broadcast(msg *model.Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ch := range r.messageChannels {
		select {
		case ch <- msg:
		default:
		}
	}
}
