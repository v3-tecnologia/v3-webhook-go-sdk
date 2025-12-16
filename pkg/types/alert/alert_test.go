package alert

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
)

func TestAlertEvent_GetAlertLevel(t *testing.T) {
	tests := []struct {
		name     string
		sub      base.EventSub
		expected string
	}{
		{"Critical", base.EventSub("EVENT_SUB_ALERT_CRITICAL"), "critical"},
		{"Warning", base.EventSub("EVENT_SUB_ALERT_WARNING"), "warning"},
		{"Info", base.EventSub("EVENT_SUB_ALERT_INFO"), "info"},
		{"Unknown", base.EventSub("EVENT_SUB_ORDER_STATUS"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseEvent := &base.BaseEvent{
				ID:         "event-123",
				Status:     base.EventStatus("STATUS_RECEIVED"),
				CreatedAt:  time.Now(),
				Type:       base.EventType("EVENT_TYPE_GENERAL"),
				Category:   base.EventCategory("EVENT_CATEGORY_ALERT"),
				Sub:        tt.sub,
				Attributes: base.Attributes{},
			}

			event := New(baseEvent)

			if got := event.GetAlertLevel(); got != tt.expected {
				t.Errorf("GetAlertLevel() = %s, esperava %s", got, tt.expected)
			}
		})
	}
}

func TestAlertEvent_GetAlertEventData(t *testing.T) {
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
		Category:   base.EventCategory("EVENT_CATEGORY_ALERT"),
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

func TestAlertEvent_GetEventName(t *testing.T) {
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
		Category:   base.EventCategory("EVENT_CATEGORY_ALERT"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "SD_CARD_MOUNTED" {
		t.Errorf("GetEventName() = %s, esperava SD_CARD_MOUNTED", got)
	}
}

func TestAlertEvent_GetAlertEventData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_ALERT"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetAlertEventData(); got != nil {
		t.Error("GetAlertEventData() retornou AlertEventData, esperava nil")
	}
}

func TestAlertEvent_GetEventName_Empty(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_ALERT"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetEventName(); got != "" {
		t.Errorf("GetEventName() = %s, esperava string vazia", got)
	}
}

func TestAlertEvent_GetAlertEventData_NoAlert(t *testing.T) {
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
		Category:   base.EventCategory("EVENT_CATEGORY_ALERT"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetAlertEventData(); got != nil {
		t.Error("GetAlertEventData() retornou AlertEventData, esperava nil")
	}
}

func TestAlertEvent_GetAlertEventData_InvalidJSON(t *testing.T) {
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
		Category:   base.EventCategory("EVENT_CATEGORY_ALERT"),
		Sub:        base.EventSub("EVENT_SUB_ALERT_INFO"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetAlertEventData(); got != nil {
		t.Error("GetAlertEventData() retornou AlertEventData, esperava nil para JSON inválido")
	}
}
