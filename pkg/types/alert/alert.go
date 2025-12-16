package alert

import (
	"encoding/json"
	"time"

	"go-eventlib/pkg/types/base"
)

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

func (e *Event) GetAlertLevel() string {
	switch e.Sub {
	case base.EventSubAlertCritical:
		return "critical"
	case base.EventSubAlertWarning:
		return "warning"
	case base.EventSubAlertInfo:
		return "info"
	default:
		return "unknown"
	}
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
	return ""
}
