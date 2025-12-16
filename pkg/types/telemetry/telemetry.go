package telemetry

import (
	"encoding/json"
	"time"

	"go-eventlib/pkg/types/base"
)

type IgnitionStatus string

const (
	IgnitionStatusOn  IgnitionStatus = "IGNITION_STATUS_ON"
	IgnitionStatusOff IgnitionStatus = "IGNITION_STATUS_OFF"
)

type Telemetry struct {
	ID              string                    `json:"id"`
	Status          IgnitionStatus            `json:"status"`
	Hardware        interface{}               `json:"hardware,omitempty"`
	FirmwareVersion interface{}               `json:"firmware_version,omitempty"`
	Connection      interface{}               `json:"connection,omitempty"`
	Connectivity    interface{}               `json:"connectivity,omitempty"`
	Metrics         map[string]*BatteryMetric `json:"metrics,omitempty"`
	Timestamp       time.Time                 `json:"timestamp"`
}

type BatteryMetric struct {
	Component string  `json:"component"`
	Status    string  `json:"status"`
	Voltage   float64 `json:"voltage,omitempty"`
}

type TripTelemetry struct {
	ID        string         `json:"id"`
	EventName string         `json:"event_name,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
	Ignition  *IgnitionEvent `json:"ignition,omitempty"`
}

type IgnitionEvent struct {
	Name     string      `json:"name"`
	Status   string      `json:"status"`
	Location interface{} `json:"location,omitempty"`
}

type Event struct {
	*base.BaseEvent
}

func New(baseEvent *base.BaseEvent) *Event {
	return &Event{BaseEvent: baseEvent}
}

func (e *Event) GetTelemetryData() *Telemetry {
	if e.Attributes.Data == nil || e.Attributes.Data.Telemetry == nil {
		return nil
	}

	data, err := json.Marshal(e.Attributes.Data.Telemetry)
	if err != nil {
		return nil
	}

	var telemetry Telemetry
	if err := json.Unmarshal(data, &telemetry); err != nil {
		return nil
	}

	return &telemetry
}

func (e *Event) GetBatteryMetrics() map[string]*BatteryMetric {
	if telemetry := e.GetTelemetryData(); telemetry != nil {
		return telemetry.Metrics
	}
	return nil
}
