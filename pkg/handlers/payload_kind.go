// Package handlers defines the handler contract used by the webhook processor:
// EventContext metadata, EventHandlingResult, and optional persistence ports.
package handlers

import (
	eventsv1 "github.com/v3-tecnologia/protocol-cloud/proto/gen/domain/v1/events"
	"google.golang.org/protobuf/proto"
)

// EventPayloadKind classifies the concrete payload domain passed to a handler.
// Values map to protocol envelope messages under domain.events.v1.
type EventPayloadKind string

const (
	PayloadKindDms                  EventPayloadKind = "dms"
	PayloadKindOrder                EventPayloadKind = "order"
	PayloadKindConnection           EventPayloadKind = "connection"
	PayloadKindVision               EventPayloadKind = "vision"
	PayloadKindHardware             EventPayloadKind = "hardware"
	PayloadKindSystem               EventPayloadKind = "system"
	PayloadKindTelemetry            EventPayloadKind = "telemetry"
	PayloadKindAlert                EventPayloadKind = "alert"
	PayloadKindDriverBehavior       EventPayloadKind = "driver_behavior"
	PayloadKindVehicle              EventPayloadKind = "vehicle"
	PayloadKindDriverIdentification EventPayloadKind = "driver_identification"
	PayloadKindGeofence             EventPayloadKind = "geofence"
	PayloadKindAdas                 EventPayloadKind = "adas"
	PayloadKindUnknown              EventPayloadKind = "unknown"
)

// PayloadKindFromEnvelope derives the kind from the protocol envelope message.
func PayloadKindFromEnvelope(envelope proto.Message) EventPayloadKind {
	switch envelope.(type) {
	case *eventsv1.Dms:
		return PayloadKindDms
	case *eventsv1.Alert:
		return PayloadKindAlert
	case *eventsv1.Vision:
		return PayloadKindVision
	case *eventsv1.Connection:
		return PayloadKindConnection
	case *eventsv1.Health:
		return PayloadKindHardware
	case *eventsv1.Telemetry:
		return PayloadKindTelemetry
	case *eventsv1.DriverBehavior:
		return PayloadKindDriverBehavior
	case *eventsv1.Vehicle:
		return PayloadKindVehicle
	case *eventsv1.DriverIdentification:
		return PayloadKindDriverIdentification
	case *eventsv1.Geofence:
		return PayloadKindGeofence
	case *eventsv1.Adas:
		return PayloadKindAdas
	case *eventsv1.System:
		return PayloadKindSystem
	case *eventsv1.OrderStatus:
		return PayloadKindOrder
	default:
		return PayloadKindUnknown
	}
}
