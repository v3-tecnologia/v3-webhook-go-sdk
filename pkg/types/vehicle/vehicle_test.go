package vehicle

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
	"go-eventlib/pkg/types/telemetry"
)

func TestVehicleEvent_GetVehicleTelemetry(t *testing.T) {
	telemetryData := map[string]interface{}{
		"id":        "telemetry-123",
		"status":    "IGNITION_STATUS_ON",
		"timestamp": time.Now(),
	}

	data := &base.Data{
		Telemetry: telemetryData,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVehicle,
		Sub:        base.EventSubTelemetryIgnition,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetVehicleTelemetry()
	if got == nil {
		t.Error("GetVehicleTelemetry() retornou nil, esperava Telemetry")
		return
	}

	if got.ID != "telemetry-123" {
		t.Errorf("GetVehicleTelemetry().ID = %s, esperava telemetry-123", got.ID)
	}

	if got.Status != telemetry.IgnitionStatusOn {
		t.Errorf("GetVehicleTelemetry().Status = %s, esperava IGNITION_STATUS_ON", got.Status)
	}
}

func TestVehicleEvent_GetVehicleTelemetry_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVehicle,
		Sub:        base.EventSubTelemetryIgnition,
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetVehicleTelemetry(); got != nil {
		t.Error("GetVehicleTelemetry() retornou Telemetry, esperava nil")
	}
}
