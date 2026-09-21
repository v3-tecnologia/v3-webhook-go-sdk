package handlers

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// EventWriter persists protobuf events. Implemented by the SDK consumer.
type EventWriter interface {
	Save(ctx context.Context, eventContext EventContext, evt proto.Message) error
}

// EventReader reads persisted protobuf events and relationships. Implemented by the SDK consumer.
type EventReader interface {
	GetEventByID(ctx context.Context, id string, dest proto.Message) (bool, error)
	GetEvents(ctx context.Context, dest proto.Message, max int) ([]proto.Message, error)
	GetRootEvent(ctx context.Context, child proto.Message, dest proto.Message) (bool, error)
}
