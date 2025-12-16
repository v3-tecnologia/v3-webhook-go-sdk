package dms

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
)

func TestDMSEvent_GetDMSData(t *testing.T) {
	dmsData := map[string]interface{}{
		"id":         "dms-123",
		"event_name": "DROWSINESS",
		"timestamp":  time.Now(),
		"drowsiness": map[string]interface{}{
			"name":       "DROWSINESS",
			"confidence": 0.90,
		},
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DMS",
		"dms":              dmsData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_DMS"),
		Sub:        base.EventSub("EVENT_SUB_DMS_BASIC"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetDMSData()
	if got == nil {
		t.Error("GetDMSData() retornou nil, esperava DMSEventData")
		return
	}

	if got.ID != "dms-123" {
		t.Errorf("GetDMSData().ID = %s, esperava dms-123", got.ID)
	}

	if got.EventName != "DROWSINESS" {
		t.Errorf("GetDMSData().EventName = %s, esperava DROWSINESS", got.EventName)
	}
}

func TestDMSEvent_GetDrowsinessData(t *testing.T) {
	dmsData := map[string]interface{}{
		"id":         "dms-123",
		"event_name": "DROWSINESS",
		"timestamp":  time.Now(),
		"drowsiness": map[string]interface{}{
			"name":       "DROWSINESS",
			"confidence": 0.90,
		},
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DMS",
		"dms":              dmsData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_DMS"),
		Sub:        base.EventSub("EVENT_SUB_DMS_BASIC"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetDrowsinessData()
	if got == nil {
		t.Error("GetDrowsinessData() retornou nil, esperava map[string]interface{}")
		return
	}

	if got["name"] != "DROWSINESS" {
		t.Errorf("GetDrowsinessData()[\"name\"] = %v, esperava DROWSINESS", got["name"])
	}
}

func TestDMSEvent_GetEventName(t *testing.T) {
	dmsData := map[string]interface{}{
		"id":         "dms-123",
		"event_name": "DROWSINESS",
		"timestamp":  time.Now(),
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DMS",
		"dms":              dmsData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_DMS"),
		Sub:        base.EventSub("EVENT_SUB_DMS_BASIC"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "DROWSINESS" {
		t.Errorf("GetEventName() = %s, esperava DROWSINESS", got)
	}
}

func TestDMSEvent_GetDrinkingData(t *testing.T) {
	dmsData := map[string]interface{}{
		"id":         "dms-123",
		"event_name": "DRINKING",
		"timestamp":  time.Now(),
		"drinking": map[string]interface{}{
			"name":       "DRINKING",
			"confidence": 0.85,
		},
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DMS",
		"dms":              dmsData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_DMS"),
		Sub:        base.EventSub("EVENT_SUB_DMS_BASIC"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetDrinkingData()
	if got == nil {
		t.Error("GetDrinkingData() retornou nil, esperava map[string]interface{}")
		return
	}

	if got["name"] != "DRINKING" {
		t.Errorf("GetDrinkingData()[\"name\"] = %v, esperava DRINKING", got["name"])
	}
}

func TestDMSEvent_GetDrinkingData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_DMS"),
		Sub:        base.EventSub("EVENT_SUB_DMS_BASIC"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetDrinkingData(); got != nil {
		t.Error("GetDrinkingData() retornou dados, esperava nil")
	}
}

func TestDMSEvent_GetEventName_Empty(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_DMS"),
		Sub:        base.EventSub("EVENT_SUB_DMS_BASIC"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "" {
		t.Errorf("GetEventName() = %s, esperava string vazia", got)
	}
}

func TestDMSEvent_GetDMSData_NoDMS(t *testing.T) {
	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DMS",
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_DMS"),
		Sub:        base.EventSub("EVENT_SUB_DMS_BASIC"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetDMSData(); got != nil {
		t.Error("GetDMSData() retornou DMSEventData, esperava nil")
	}
}

func TestDMSEvent_GetDMSData_InvalidJSON(t *testing.T) {
	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DMS",
		"dms":              make(chan int),
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_DMS"),
		Sub:        base.EventSub("EVENT_SUB_DMS_BASIC"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetDMSData(); got != nil {
		t.Error("GetDMSData() retornou DMSEventData, esperava nil para JSON inválido")
	}
}
