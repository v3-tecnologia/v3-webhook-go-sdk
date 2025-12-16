package vision

import (
	"encoding/json"
	"time"

	"go-eventlib/pkg/types/base"
)

type VisionEventData struct {
	ID               string                 `json:"id"`
	EventName        string                 `json:"event_name"`
	Timestamp        time.Time              `json:"timestamp"`
	FaceDetected     map[string]interface{} `json:"face_detected,omitempty"`
	FaceLost         map[string]interface{} `json:"face_lost,omitempty"`
	FaceTracked      map[string]interface{} `json:"face_tracked,omitempty"`
	NoFaceDetected   map[string]interface{} `json:"no_face_detected,omitempty"`
	CameraObstructed map[string]interface{} `json:"camera_obstructed,omitempty"`
}

type Event struct {
	*base.BaseEvent
}

func New(baseEvent *base.BaseEvent) *Event {
	return &Event{BaseEvent: baseEvent}
}

func (e *Event) GetVisionEventData() *VisionEventData {
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

	visionData, ok := tripEvent["vision"]
	if !ok {
		return nil
	}

	visionBytes, err := json.Marshal(visionData)
	if err != nil {
		return nil
	}

	var vision VisionEventData
	if err := json.Unmarshal(visionBytes, &vision); err != nil {
		return nil
	}

	return &vision
}

func (e *Event) GetFaceDetectedData() map[string]interface{} {
	if vision := e.GetVisionEventData(); vision != nil {
		return vision.FaceDetected
	}
	return nil
}

func (e *Event) GetFaceLostData() map[string]interface{} {
	if vision := e.GetVisionEventData(); vision != nil {
		return vision.FaceLost
	}
	return nil
}

func (e *Event) GetEventName() string {
	if vision := e.GetVisionEventData(); vision != nil {
		return vision.EventName
	}
	return ""
}
