package vehicle

import (
	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/types/base"
	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/types/telemetry"
)

type Event struct {
	*base.BaseEvent
}

func New(baseEvent *base.BaseEvent) *Event {
	return &Event{BaseEvent: baseEvent}
}

func (e *Event) GetVehicleTelemetry() *telemetry.Telemetry {
	if e.Attributes.Data == nil {
		return nil
	}

	telemetryEvent := telemetry.New(e.BaseEvent)
	return telemetryEvent.GetTelemetryData()
}
