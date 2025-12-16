package base

import (
	"testing"
	"time"
)

func TestBaseEvent_GetID(t *testing.T) {
	event := &BaseEvent{
		ID: "test-id-123",
	}

	if got := event.GetID(); got != "test-id-123" {
		t.Errorf("GetID() = %s, esperava test-id-123", got)
	}
}

func TestBaseEvent_GetCategory(t *testing.T) {
	event := &BaseEvent{
		Category: EventCategory("EVENT_CATEGORY_ORDER"),
	}

	if got := event.GetCategory(); string(got) != "EVENT_CATEGORY_ORDER" {
		t.Errorf("GetCategory() = %s, esperava EVENT_CATEGORY_ORDER", got)
	}
}

func TestBaseEvent_GetSubType(t *testing.T) {
	event := &BaseEvent{
		Sub: EventSub("EVENT_SUB_ORDER_STATUS"),
	}

	if got := event.GetSubType(); string(got) != "EVENT_SUB_ORDER_STATUS" {
		t.Errorf("GetSubType() = %s, esperava EVENT_SUB_ORDER_STATUS", got)
	}
}

func TestBaseEvent_GetDeviceID(t *testing.T) {
	event := &BaseEvent{
		Attributes: Attributes{
			Device: &Device{
				ID: "device-123",
			},
		},
	}

	if got := event.GetDeviceID(); got != "device-123" {
		t.Errorf("GetDeviceID() = %s, esperava device-123", got)
	}
}

func TestBaseEvent_GetDeviceID_Nil(t *testing.T) {
	event := &BaseEvent{
		Attributes: Attributes{
			Device: nil,
		},
	}

	if got := event.GetDeviceID(); got != "" {
		t.Errorf("GetDeviceID() = %s, esperava string vazia", got)
	}
}

func TestBaseEvent_GetCreatedAt(t *testing.T) {
	now := time.Now()
	event := &BaseEvent{
		CreatedAt: now,
	}

	if got := event.GetCreatedAt(); !got.Equal(now) {
		t.Errorf("GetCreatedAt() = %v, esperava %v", got, now)
	}
}
