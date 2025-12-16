package dms

import (
	"encoding/json"
	"time"

	"go-eventlib/pkg/types/base"
)

type DMSEventData struct {
	ID              string                 `json:"id"`
	EventName       string                 `json:"event_name"`
	Timestamp       time.Time              `json:"timestamp"`
	Drowsiness      map[string]interface{} `json:"drowsiness,omitempty"`
	Drinking        map[string]interface{} `json:"drinking,omitempty"`
	Eating          map[string]interface{} `json:"eating,omitempty"`
	EyeClosure      map[string]interface{} `json:"eye_closure,omitempty"`
	GazeDistraction map[string]interface{} `json:"gaze_distraction,omitempty"`
	GazeFixation    map[string]interface{} `json:"gaze_fixation,omitempty"`
	OnPhone         map[string]interface{} `json:"on_phone,omitempty"`
	PoseDistraction map[string]interface{} `json:"pose_distraction_pitch,omitempty"`
	Smoking         map[string]interface{} `json:"smoking,omitempty"`
	Yawning         map[string]interface{} `json:"yawning,omitempty"`
}

type Event struct {
	*base.BaseEvent
}

func New(baseEvent *base.BaseEvent) *Event {
	return &Event{BaseEvent: baseEvent}
}

func (e *Event) GetDMSData() *DMSEventData {
	if e.Attributes.Data == nil || e.Attributes.Data.TripEvent == nil {
		return nil
	}

	data, err := json.Marshal(e.Attributes.Data.TripEvent)
	if err != nil {
		return nil
	}

	var tripEvent map[string]interface{}
	if err := json.Unmarshal(data, &tripEvent); err != nil {
		return nil
	}

	dmsData, ok := tripEvent["dms"]
	if !ok {
		return nil
	}

	dmsBytes, err := json.Marshal(dmsData)
	if err != nil {
		return nil
	}

	var dms DMSEventData
	if err := json.Unmarshal(dmsBytes, &dms); err != nil {
		return nil
	}

	return &dms
}

func (e *Event) GetDrowsinessData() map[string]interface{} {
	if dms := e.GetDMSData(); dms != nil {
		return dms.Drowsiness
	}
	return nil
}

func (e *Event) GetDrinkingData() map[string]interface{} {
	if dms := e.GetDMSData(); dms != nil {
		return dms.Drinking
	}
	return nil
}

func (e *Event) GetEventName() string {
	if dms := e.GetDMSData(); dms != nil {
		return dms.EventName
	}
	return ""
}
