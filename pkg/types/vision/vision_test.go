package vision

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
)

func TestVisionEvent_GetVisionEventData(t *testing.T) {
	visionData := map[string]interface{}{
		"id":         "vision-123",
		"event_name": "FACE_DETECTED",
		"timestamp":  time.Now(),
		"face_detected": map[string]interface{}{
			"name": "FACE_DETECTED",
		},
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "VISION",
		"vision":           visionData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVision,
		Sub:        base.EventSubVisionBasic,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetVisionEventData()
	if got == nil {
		t.Error("GetVisionEventData() retornou nil, esperava VisionEventData")
		return
	}

	if got.ID != "vision-123" {
		t.Errorf("GetVisionEventData().ID = %s, esperava vision-123", got.ID)
	}

	if got.EventName != "FACE_DETECTED" {
		t.Errorf("GetVisionEventData().EventName = %s, esperava FACE_DETECTED", got.EventName)
	}
}

func TestVisionEvent_GetEventName(t *testing.T) {
	visionData := map[string]interface{}{
		"id":         "vision-123",
		"event_name": "FACE_DETECTED",
		"timestamp":  time.Now(),
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "VISION",
		"vision":           visionData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVision,
		Sub:        base.EventSubVisionBasic,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "FACE_DETECTED" {
		t.Errorf("GetEventName() = %s, esperava FACE_DETECTED", got)
	}
}

func TestVisionEvent_GetFaceDetectedData(t *testing.T) {
	visionData := map[string]interface{}{
		"id":         "vision-123",
		"event_name": "FACE_DETECTED",
		"timestamp":  time.Now(),
		"face_detected": map[string]interface{}{
			"name": "FACE_DETECTED",
			"state": map[string]interface{}{
				"pending": map[string]interface{}{
					"reason": "AWAITING_INFERENCE",
				},
			},
		},
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "VISION",
		"vision":           visionData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVision,
		Sub:        base.EventSubVisionBasic,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetFaceDetectedData()
	if got == nil {
		t.Error("GetFaceDetectedData() retornou nil, esperava map[string]interface{}")
		return
	}

	if got["name"] != "FACE_DETECTED" {
		t.Errorf("GetFaceDetectedData()[\"name\"] = %v, esperava FACE_DETECTED", got["name"])
	}
}

func TestVisionEvent_GetFaceLostData(t *testing.T) {
	visionData := map[string]interface{}{
		"id":         "vision-123",
		"event_name": "FACE_LOST",
		"timestamp":  time.Now(),
		"face_lost": map[string]interface{}{
			"name": "FACE_LOST",
		},
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "VISION",
		"vision":           visionData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVision,
		Sub:        base.EventSubVisionBasic,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetFaceLostData()
	if got == nil {
		t.Error("GetFaceLostData() retornou nil, esperava map[string]interface{}")
		return
	}

	if got["name"] != "FACE_LOST" {
		t.Errorf("GetFaceLostData()[\"name\"] = %v, esperava FACE_LOST", got["name"])
	}
}

func TestVisionEvent_GetFaceDetectedData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVision,
		Sub:        base.EventSubVisionBasic,
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetFaceDetectedData(); got != nil {
		t.Error("GetFaceDetectedData() retornou dados, esperava nil")
	}
}

func TestVisionEvent_GetEventName_Empty(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVision,
		Sub:        base.EventSubVisionBasic,
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "" {
		t.Errorf("GetEventName() = %s, esperava string vazia", got)
	}
}

func TestVisionEvent_GetVisionEventData_NoVision(t *testing.T) {
	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "VISION",
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVision,
		Sub:        base.EventSubVisionBasic,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetVisionEventData(); got != nil {
		t.Error("GetVisionEventData() retornou VisionEventData, esperava nil")
	}
}

func TestVisionEvent_GetVisionEventData_InvalidJSON(t *testing.T) {
	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "VISION",
		"vision":           make(chan int),
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryVision,
		Sub:        base.EventSubVisionBasic,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetVisionEventData(); got != nil {
		t.Error("GetVisionEventData() retornou VisionEventData, esperava nil para JSON inválido")
	}
}
