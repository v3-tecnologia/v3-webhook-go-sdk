package processing

import (
	"fmt"
	"strings"

	eventsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/events"
	locationv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/location"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type resolvedEvent struct {
	envelope proto.Message
	payload  proto.Message
}

func resolveFromEventData(data *eventsv1.EventData) (*resolvedEvent, error) {
	if data == nil {
		return nil, nil
	}

	var container proto.Message
	switch {
	case data.GetTripEvent() != nil:
		container = data.GetTripEvent()
	case data.GetStandaloneEvent() != nil:
		container = data.GetStandaloneEvent()
	default:
		return nil, fmt.Errorf("protocol violation: event data missing trip_event or standalone_event")
	}

	return resolveContainer(container)
}

func resolveContainer(container proto.Message) (*resolvedEvent, error) {
	if container == nil {
		return nil, fmt.Errorf("protocol violation: empty event container")
	}

	envelope := extractGroupEnvelope(container)
	if envelope == nil {
		return nil, fmt.Errorf("protocol violation: event_group oneof is unset on %s", proto.MessageName(container))
	}

	payload := extractOneOfMessage(envelope)
	if payload == nil {
		return nil, fmt.Errorf("protocol violation: event oneof is unset on %s", proto.MessageName(envelope))
	}

	return &resolvedEvent{
		envelope: envelope,
		payload:  payload,
	}, nil
}

func extractGroupEnvelope(container proto.Message) proto.Message {
	m := container.ProtoReflect()
	oneofs := m.Descriptor().Oneofs()
	for i := 0; i < oneofs.Len(); i++ {
		od := oneofs.Get(i)
		if od.Name() != "event_group" {
			continue
		}
		fd := m.WhichOneof(od)
		if fd == nil || fd.Kind() != protoreflect.MessageKind {
			continue
		}
		return m.Get(fd).Message().Interface()
	}
	return nil
}

func extractOneOfMessage(msg proto.Message) proto.Message {
	m := msg.ProtoReflect()
	oneofs := m.Descriptor().Oneofs()
	for i := 0; i < oneofs.Len(); i++ {
		od := oneofs.Get(i)
		fd := m.WhichOneof(od)
		if fd == nil || fd.Kind() != protoreflect.MessageKind {
			continue
		}
		return m.Get(fd).Message().Interface()
	}
	return nil
}

func extractLocation(msg proto.Message) *locationv1.Location {
	if msg == nil {
		return nil
	}
	if loc, ok := msg.(*locationv1.Location); ok {
		return loc
	}

	m := msg.ProtoReflect()
	fields := m.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if fd.IsMap() || fd.IsList() || fd.Kind() != protoreflect.MessageKind || !m.Has(fd) {
			continue
		}
		nested := m.Get(fd).Message().Interface()
		if fd.Name() == "location" {
			if loc, ok := nested.(*locationv1.Location); ok {
				return loc
			}
		}
		if loc := extractLocation(nested); loc != nil {
			return loc
		}
	}
	return nil
}

func messageKey(msg proto.Message) string {
	if msg == nil {
		return ""
	}
	return string(proto.MessageName(msg))
}

func normalizeWhitespace(s string) string {
	return strings.TrimSpace(s)
}
