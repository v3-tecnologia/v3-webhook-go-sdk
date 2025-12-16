package hardware

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
)

func TestHardwareEvent_GetSystemEventData(t *testing.T) {
	systemData := map[string]interface{}{
		"id":         "system-123",
		"event_name": "RESTART",
		"timestamp":  time.Now(),
	}

	standalone := map[string]interface{}{
		"event_group_name": "SYSTEM",
		"system":           systemData,
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_SYSTEM_UPLOAD"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetSystemEventData()
	if got == nil {
		t.Error("GetSystemEventData() retornou nil, esperava SystemEventData")
		return
	}

	if got.ID != "system-123" {
		t.Errorf("GetSystemEventData().ID = %s, esperava system-123", got.ID)
	}

	if got.EventName != "RESTART" {
		t.Errorf("GetSystemEventData().EventName = %s, esperava RESTART", got.EventName)
	}
}

func TestHardwareEvent_GetAlertEventData(t *testing.T) {
	alertData := map[string]interface{}{
		"id":         "alert-123",
		"event_name": "SD_CARD_MOUNTED",
		"timestamp":  time.Now(),
		"sd_card_mounted": map[string]interface{}{
			"name": "SD_CARD_MOUNTED",
		},
	}

	standalone := map[string]interface{}{
		"event_group_name": "ALERT",
		"alert":            alertData,
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetAlertEventData()
	if got == nil {
		t.Error("GetAlertEventData() retornou nil, esperava AlertEventData")
		return
	}

	if got.ID != "alert-123" {
		t.Errorf("GetAlertEventData().ID = %s, esperava alert-123", got.ID)
	}

	if got.EventName != "SD_CARD_MOUNTED" {
		t.Errorf("GetAlertEventData().EventName = %s, esperava SD_CARD_MOUNTED", got.EventName)
	}
}

func TestHardwareEvent_GetEventName(t *testing.T) {
	alertData := map[string]interface{}{
		"id":         "alert-123",
		"event_name": "SD_CARD_MOUNTED",
		"timestamp":  time.Now(),
	}

	standalone := map[string]interface{}{
		"event_group_name": "ALERT",
		"alert":            alertData,
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "SD_CARD_MOUNTED" {
		t.Errorf("GetEventName() = %s, esperava SD_CARD_MOUNTED", got)
	}
}

func TestHardwareEvent_GetEventName_FromSystem(t *testing.T) {
	systemData := map[string]interface{}{
		"id":         "system-123",
		"event_name": "RESTART",
		"timestamp":  time.Now(),
	}

	standalone := map[string]interface{}{
		"event_group_name": "SYSTEM",
		"system":           systemData,
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_SYSTEM_UPLOAD"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "RESTART" {
		t.Errorf("GetEventName() = %s, esperava RESTART", got)
	}
}

func TestHardwareEvent_GetEventName_Empty(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_SYSTEM_UPLOAD"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "" {
		t.Errorf("GetEventName() = %s, esperava string vazia", got)
	}
}

func TestHardwareEvent_GetSystemEventData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_SYSTEM_UPLOAD"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetSystemEventData(); got != nil {
		t.Error("GetSystemEventData() retornou SystemEventData, esperava nil")
	}
}

func TestHardwareEvent_GetAlertEventData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetAlertEventData(); got != nil {
		t.Error("GetAlertEventData() retornou AlertEventData, esperava nil")
	}
}

func TestHardwareEvent_GetSystemEventData_NoSystem(t *testing.T) {
	standalone := map[string]interface{}{
		"event_group_name": "SYSTEM",
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_SYSTEM_UPLOAD"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetSystemEventData(); got != nil {
		t.Error("GetSystemEventData() retornou SystemEventData, esperava nil")
	}
}

func TestHardwareEvent_GetAlertEventData_NoAlert(t *testing.T) {
	standalone := map[string]interface{}{
		"event_group_name": "ALERT",
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetAlertEventData(); got != nil {
		t.Error("GetAlertEventData() retornou AlertEventData, esperava nil")
	}
}

func TestHardwareEvent_GetSystemEventData_InvalidJSON(t *testing.T) {
	standalone := map[string]interface{}{
		"event_group_name": "SYSTEM",
		"system":           make(chan int),
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_SYSTEM_UPLOAD"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetSystemEventData(); got != nil {
		t.Error("GetSystemEventData() retornou SystemEventData, esperava nil para JSON inválido")
	}
}

func TestHardwareEvent_GetAlertEventData_InvalidJSON(t *testing.T) {
	standalone := map[string]interface{}{
		"event_group_name": "ALERT",
		"alert":            make(chan int),
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_HEALTH"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetAlertEventData(); got != nil {
		t.Error("GetAlertEventData() retornou AlertEventData, esperava nil para JSON inválido")
	}
}
