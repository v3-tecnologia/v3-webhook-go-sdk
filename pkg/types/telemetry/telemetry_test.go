package telemetry

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
)

func TestTelemetryEvent_GetTelemetryData(t *testing.T) {
	telemetryData := map[string]interface{}{
		"id":        "telemetry-123",
		"status":    "IGNITION_STATUS_ON",
		"timestamp": time.Now(),
		"metrics": map[string]interface{}{
			"battery1": map[string]interface{}{
				"component": "battery",
				"status":    "OK",
				"voltage":   12.5,
			},
		},
	}

	data := &base.Data{
		Telemetry: telemetryData,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_TELEMETRY"),
		Sub:        base.EventSub("EVENT_SUB_TELEMETRY_BATTERY"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetTelemetryData()
	if got == nil {
		t.Error("GetTelemetryData() retornou nil, esperava Telemetry")
		return
	}

	if got.ID != "telemetry-123" {
		t.Errorf("GetTelemetryData().ID = %s, esperava telemetry-123", got.ID)
	}

	if got.Status != IgnitionStatusOn {
		t.Errorf("GetTelemetryData().Status = %s, esperava IGNITION_STATUS_ON", got.Status)
	}
}

func TestTelemetryEvent_GetBatteryMetrics(t *testing.T) {
	telemetryData := map[string]interface{}{
		"id":        "telemetry-123",
		"status":    "IGNITION_STATUS_ON",
		"timestamp": time.Now(),
		"metrics": map[string]interface{}{
			"battery1": map[string]interface{}{
				"component": "battery",
				"status":    "OK",
				"voltage":   12.5,
			},
		},
	}

	data := &base.Data{
		Telemetry: telemetryData,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_TELEMETRY"),
		Sub:        base.EventSub("EVENT_SUB_TELEMETRY_BATTERY"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetBatteryMetrics()
	if got == nil {
		t.Error("GetBatteryMetrics() retornou nil, esperava map[string]*BatteryMetric")
		return
	}

	battery, ok := got["battery1"]
	if !ok {
		t.Error("GetBatteryMetrics() não contém battery1")
		return
	}

	if battery.Component != "battery" {
		t.Errorf("GetBatteryMetrics()[\"battery1\"].Component = %s, esperava battery", battery.Component)
	}

	if battery.Voltage != 12.5 {
		t.Errorf("GetBatteryMetrics()[\"battery1\"].Voltage = %f, esperava 12.5", battery.Voltage)
	}
}

func TestTelemetryEvent_GetTelemetryData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_TELEMETRY"),
		Sub:        base.EventSub("EVENT_SUB_TELEMETRY_BATTERY"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetTelemetryData(); got != nil {
		t.Error("GetTelemetryData() retornou Telemetry, esperava nil")
	}
}

func TestTelemetryEvent_GetBatteryMetrics_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_TELEMETRY"),
		Sub:        base.EventSub("EVENT_SUB_TELEMETRY_BATTERY"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetBatteryMetrics(); got != nil {
		t.Error("GetBatteryMetrics() retornou map, esperava nil")
	}
}

func TestTelemetryEvent_GetTelemetryData_InvalidJSON(t *testing.T) {
	data := &base.Data{
		Telemetry: make(chan int),
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_TELEMETRY"),
		Sub:        base.EventSub("EVENT_SUB_TELEMETRY_BATTERY"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetTelemetryData(); got != nil {
		t.Error("GetTelemetryData() retornou Telemetry, esperava nil para JSON inválido")
	}
}
