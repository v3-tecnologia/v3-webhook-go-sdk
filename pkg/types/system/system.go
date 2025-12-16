package system

import (
	"encoding/json"
	"time"

	"go-eventlib/pkg/types/base"
)

type SystemEventData struct {
	ID        string       `json:"id"`
	EventName string       `json:"event_name"`
	Upload    *UploadEvent `json:"upload,omitempty"`
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

type Event struct {
	*base.BaseEvent
}

func New(baseEvent *base.BaseEvent) *Event {
	return &Event{BaseEvent: baseEvent}
}

func (e *Event) GetSystemData() *SystemEventData {
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

func (e *Event) GetUploadData() *UploadEvent {
	if sys := e.GetSystemData(); sys != nil {
		return sys.Upload
	}
	return nil
}
