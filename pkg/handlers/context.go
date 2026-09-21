package handlers

import (
	"context"
	"fmt"

	eventsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/events"
	locationv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/location"
	ordersv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/orders"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EventContext carries shared metadata and optional persistence helpers for a handler invocation.
type EventContext struct {
	ID          string
	HasMedia    bool
	PayloadKind EventPayloadKind
	Status      eventsv1.Status
	CreatedAt   *timestamppb.Timestamp
	Type        eventsv1.EventType
	Category    eventsv1.EventCategory
	Sub         eventsv1.EventSub
	Device      *ordersv1.Device
	Location    *locationv1.Location

	writer EventWriter
	reader EventReader
}

// NewEventContext builds an EventContext used by the processor when invoking handlers.
func NewEventContext(
	id string,
	hasMedia bool,
	status eventsv1.Status,
	createdAt *timestamppb.Timestamp,
	eventType eventsv1.EventType,
	category eventsv1.EventCategory,
	sub eventsv1.EventSub,
	device *ordersv1.Device,
	location *locationv1.Location,
	payloadKind EventPayloadKind,
	writer EventWriter,
	reader EventReader,
) EventContext {
	return EventContext{
		ID:          id,
		HasMedia:    hasMedia,
		Status:      status,
		CreatedAt:   createdAt,
		Type:        eventType,
		Category:    category,
		Sub:         sub,
		Device:      device,
		Location:    location,
		PayloadKind: payloadKind,
		writer:      writer,
		reader:      reader,
	}
}

// Save persists evt through the configured EventWriter.
func (c EventContext) Save(ctx context.Context, evt proto.Message) error {
	if c.writer == nil {
		return fmt.Errorf("no EventWriter configured")
	}
	return c.writer.Save(ctx, c, evt)
}

// GetEventByID loads a previously saved event by id into dest.
func (c EventContext) GetEventByID(ctx context.Context, id string, dest proto.Message) (bool, error) {
	if c.reader == nil {
		return false, fmt.Errorf("no EventReader configured")
	}
	return c.reader.GetEventByID(ctx, id, dest)
}

// GetEvents returns up to max events of dest's protobuf type.
func (c EventContext) GetEvents(ctx context.Context, dest proto.Message, max int) ([]proto.Message, error) {
	if c.reader == nil {
		return nil, fmt.Errorf("no EventReader configured")
	}
	return c.reader.GetEvents(ctx, dest, max)
}

// GetRootEvent resolves a parent/root event for child using the configured reader.
func (c EventContext) GetRootEvent(ctx context.Context, child proto.Message, dest proto.Message) (bool, error) {
	if c.reader == nil {
		return false, fmt.Errorf("no EventReader configured")
	}
	return c.reader.GetRootEvent(ctx, child, dest)
}
