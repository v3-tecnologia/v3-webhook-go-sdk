// Package persistence provides in-memory EventWriter/EventReader implementations for tests and examples.
package persistence

import (
	"context"
	"fmt"
	"sync"

	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/handlers"
	"google.golang.org/protobuf/proto"
)

type storedEvent struct {
	id      string
	payload []byte
	msgType string
}

// InMemoryStore is a concurrency-safe store shared by writer and reader.
type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string]storedEvent
}

// NewInMemoryStore creates an empty store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string]storedEvent)}
}

// InMemoryWriter implements handlers.EventWriter against an InMemoryStore.
type InMemoryWriter struct {
	store *InMemoryStore
}

// NewInMemoryWriter creates a writer bound to store.
func NewInMemoryWriter(store *InMemoryStore) *InMemoryWriter {
	return &InMemoryWriter{store: store}
}

// Save marshals and stores evt keyed by eventContext.ID.
func (w *InMemoryWriter) Save(_ context.Context, eventContext handlers.EventContext, evt proto.Message) error {
	bytes, err := proto.Marshal(evt)
	if err != nil {
		return err
	}
	w.store.mu.Lock()
	defer w.store.mu.Unlock()
	w.store.data[eventContext.ID] = storedEvent{
		id:      eventContext.ID,
		payload: bytes,
		msgType: string(evt.ProtoReflect().Descriptor().FullName()),
	}
	return nil
}

// InMemoryReader implements handlers.EventReader against an InMemoryStore.
type InMemoryReader struct {
	store *InMemoryStore
}

// NewInMemoryReader creates a reader bound to store.
func NewInMemoryReader(store *InMemoryStore) *InMemoryReader {
	return &InMemoryReader{store: store}
}

// GetEventByID unmarshals a stored event into dest when types match.
func (r *InMemoryReader) GetEventByID(_ context.Context, id string, dest proto.Message) (bool, error) {
	r.store.mu.RLock()
	stored, ok := r.store.data[id]
	r.store.mu.RUnlock()
	if !ok {
		return false, nil
	}
	if string(dest.ProtoReflect().Descriptor().FullName()) != stored.msgType {
		return false, nil
	}
	if err := proto.Unmarshal(stored.payload, dest); err != nil {
		return false, err
	}
	return true, nil
}

// GetEvents returns up to max stored events matching dest's protobuf type.
func (r *InMemoryReader) GetEvents(_ context.Context, dest proto.Message, max int) ([]proto.Message, error) {
	if max <= 0 {
		max = 10
	}
	wantType := string(dest.ProtoReflect().Descriptor().FullName())
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	out := make([]proto.Message, 0, max)
	for _, stored := range r.store.data {
		if stored.msgType != wantType {
			continue
		}
		msg := proto.Clone(dest)
		if err := proto.Unmarshal(stored.payload, msg); err != nil {
			return nil, err
		}
		out = append(out, msg)
		if len(out) >= max {
			break
		}
	}
	return out, nil
}

// GetRootEvent looks up a parent event using child's source_id field when present.
func (r *InMemoryReader) GetRootEvent(ctx context.Context, child proto.Message, dest proto.Message) (bool, error) {
	m := child.ProtoReflect()
	fd := m.Descriptor().Fields().ByName("source_id")
	if fd == nil {
		return false, nil
	}
	sourceID := m.Get(fd).String()
	if sourceID == "" {
		return false, nil
	}
	ok, err := r.GetEventByID(ctx, sourceID, dest)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("root event not found for source_id=%s", sourceID)
	}
	return true, nil
}
