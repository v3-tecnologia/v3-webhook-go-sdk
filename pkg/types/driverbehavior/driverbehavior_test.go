package driverbehavior

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
)

func TestDriverBehaviorEvent_GetDriverBehaviorData(t *testing.T) {
	dbData := map[string]interface{}{
		"id":         "db-123",
		"event_name": "HARSH_ACCELERATION",
		"timestamp":  time.Now(),
		"acceleration_harsh": map[string]interface{}{
			"name": "HARSH_ACCELERATION",
		},
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DRIVER_BEHAVIOR",
		"driver_behavior":  dbData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryDriverBehavior,
		Sub:        base.EventSubDriverBehaviorAdvanced,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetDriverBehaviorData()
	if got == nil {
		t.Error("GetDriverBehaviorData() retornou nil, esperava DriverBehaviorEventData")
		return
	}

	if got.ID != "db-123" {
		t.Errorf("GetDriverBehaviorData().ID = %s, esperava db-123", got.ID)
	}

	if got.EventName != "HARSH_ACCELERATION" {
		t.Errorf("GetDriverBehaviorData().EventName = %s, esperava HARSH_ACCELERATION", got.EventName)
	}
}

func TestDriverBehaviorEvent_GetHarshAccelerationData(t *testing.T) {
	dbData := map[string]interface{}{
		"id":         "db-123",
		"event_name": "HARSH_ACCELERATION",
		"timestamp":  time.Now(),
		"acceleration_harsh": map[string]interface{}{
			"name": "HARSH_ACCELERATION",
		},
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DRIVER_BEHAVIOR",
		"driver_behavior":  dbData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryDriverBehavior,
		Sub:        base.EventSubDriverBehaviorAdvanced,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetHarshAccelerationData()
	if got == nil {
		t.Error("GetHarshAccelerationData() retornou nil, esperava map[string]interface{}")
		return
	}

	if got["name"] != "HARSH_ACCELERATION" {
		t.Errorf("GetHarshAccelerationData()[\"name\"] = %v, esperava HARSH_ACCELERATION", got["name"])
	}
}

func TestDriverBehaviorEvent_GetEventName(t *testing.T) {
	dbData := map[string]interface{}{
		"id":         "db-123",
		"event_name": "HARSH_ACCELERATION",
		"timestamp":  time.Now(),
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DRIVER_BEHAVIOR",
		"driver_behavior":  dbData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryDriverBehavior,
		Sub:        base.EventSubDriverBehaviorAdvanced,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "HARSH_ACCELERATION" {
		t.Errorf("GetEventName() = %s, esperava HARSH_ACCELERATION", got)
	}
}

func TestDriverBehaviorEvent_GetHarshBrakingData(t *testing.T) {
	dbData := map[string]interface{}{
		"id":         "db-123",
		"event_name": "HARSH_BRAKING",
		"timestamp":  time.Now(),
		"braking_harsh": map[string]interface{}{
			"name": "HARSH_BRAKING",
		},
	}

	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DRIVER_BEHAVIOR",
		"driver_behavior":  dbData,
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryDriverBehavior,
		Sub:        base.EventSubDriverBehaviorAdvanced,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetHarshBrakingData()
	if got == nil {
		t.Error("GetHarshBrakingData() retornou nil, esperava map[string]interface{}")
		return
	}

	if got["name"] != "HARSH_BRAKING" {
		t.Errorf("GetHarshBrakingData()[\"name\"] = %v, esperava HARSH_BRAKING", got["name"])
	}
}

func TestDriverBehaviorEvent_GetHarshBrakingData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryDriverBehavior,
		Sub:        base.EventSubDriverBehaviorAdvanced,
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetHarshBrakingData(); got != nil {
		t.Error("GetHarshBrakingData() retornou dados, esperava nil")
	}
}

func TestDriverBehaviorEvent_GetEventName_Empty(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryDriverBehavior,
		Sub:        base.EventSubDriverBehaviorAdvanced,
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "" {
		t.Errorf("GetEventName() = %s, esperava string vazia", got)
	}
}

func TestDriverBehaviorEvent_GetDriverBehaviorData_NoDriverBehavior(t *testing.T) {
	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DRIVER_BEHAVIOR",
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryDriverBehavior,
		Sub:        base.EventSubDriverBehaviorAdvanced,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetDriverBehaviorData(); got != nil {
		t.Error("GetDriverBehaviorData() retornou DriverBehaviorEventData, esperava nil")
	}
}

func TestDriverBehaviorEvent_GetDriverBehaviorData_InvalidJSON(t *testing.T) {
	tripEvent := map[string]interface{}{
		"trip_id":          "trip-123",
		"event_group_name": "DRIVER_BEHAVIOR",
		"driver_behavior":  make(chan int),
	}

	data := &base.Data{
		TripEvent: tripEvent,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategoryDriverBehavior,
		Sub:        base.EventSubDriverBehaviorAdvanced,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetDriverBehaviorData(); got != nil {
		t.Error("GetDriverBehaviorData() retornou DriverBehaviorEventData, esperava nil para JSON inválido")
	}
}
