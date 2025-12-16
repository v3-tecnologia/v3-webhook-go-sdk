package hardware

import (
	"encoding/json"
	"time"

	"go-eventlib/pkg/types/base"
)

type Hardware struct {
	Model           *HardwareModel   `json:"model,omitempty"`
	FirmwareVersion *FirmwareVersion `json:"firmware_version,omitempty"`
	PID             interface{}      `json:"pid,omitempty"`
	Uptime          interface{}      `json:"uptime,omitempty"`
}

type HardwareModel struct {
	Name       string                 `json:"name,omitempty"`
	Vendor     string                 `json:"vendor,omitempty"`
	Version    map[string]interface{} `json:"version,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

type FirmwareVersion struct {
	Name    string                 `json:"name,omitempty"`
	Version map[string]interface{} `json:"version,omitempty"`
}

type SystemEventData struct {
	ID        string       `json:"id"`
	EventName string       `json:"event_name"`
	Upload    *UploadEvent  `json:"upload,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

type UploadEvent struct {
	Name     string      `json:"name"`
	Files    []*FileInfo `json:"files,omitempty"`
	Location interface{} `json:"location,omitempty"`
}

type FileInfo struct {
	ID        string                 `json:"id"`
	SourceID  string                 `json:"source_id"`
	URL       string                 `json:"url"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type AlertEventData struct {
	ID              string                 `json:"id"`
	EventName       string                 `json:"event_name"`
	Timestamp       time.Time              `json:"timestamp"`
	SDCardMounted   map[string]interface{} `json:"sd_card_mounted,omitempty"`
	SDCardUnmounted map[string]interface{} `json:"sd_card_unmounted,omitempty"`
	SimCardInserted map[string]interface{} `json:"sim_card_inserted,omitempty"`
	SimCardRemoved  map[string]interface{} `json:"sim_card_removed,omitempty"`
}

type Event struct {
	*base.BaseEvent
}

func New(baseEvent *base.BaseEvent) *Event {
	return &Event{BaseEvent: baseEvent}
}

func (e *Event) GetSystemEventData() *SystemEventData {
	if e.Attributes.Data == nil || e.Attributes.Data.StandaloneEvent == nil {
		return nil
	}

	data, err := json.Marshal(e.Attributes.Data.StandaloneEvent)
	if err != nil {
		return nil
	}

	var standalone map[string]interface{}
	if err := json.Unmarshal(data, &standalone); err != nil {
		return nil
	}

	systemData, ok := standalone["system"]
	if !ok {
		return nil
	}

	systemBytes, err := json.Marshal(systemData)
	if err != nil {
		return nil
	}

	var system SystemEventData
	if err := json.Unmarshal(systemBytes, &system); err != nil {
		return nil
	}

	return &system
}

func (e *Event) GetAlertEventData() *AlertEventData {
	if e.Attributes.Data == nil || e.Attributes.Data.StandaloneEvent == nil {
		return nil
	}

	data, err := json.Marshal(e.Attributes.Data.StandaloneEvent)
	if err != nil {
		return nil
	}

	var standalone map[string]interface{}
	if err := json.Unmarshal(data, &standalone); err != nil {
		return nil
	}

	alertData, ok := standalone["alert"]
	if !ok {
		return nil
	}

	alertBytes, err := json.Marshal(alertData)
	if err != nil {
		return nil
	}

	var alert AlertEventData
	if err := json.Unmarshal(alertBytes, &alert); err != nil {
		return nil
	}

	return &alert
}

func (e *Event) GetEventName() string {
	if alert := e.GetAlertEventData(); alert != nil {
		return alert.EventName
	}
	if system := e.GetSystemEventData(); system != nil {
		return system.EventName
	}
	return ""
}
