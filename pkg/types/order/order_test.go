package order

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
)

func TestOrderEvent_GetOrder(t *testing.T) {
	order := &base.Order{
		ID:            "order-123",
		CorrelationID: "corr-123",
		Group:         base.OrderGroup("ORDER_GROUP_CONFIG"),
		Status:        base.OrderStatus("ORDER_STATUS_ACK"),
		Type:          base.OrderType("CONFIG"),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_ORDER"),
		Category:   base.EventCategory("EVENT_CATEGORY_ORDER"),
		Sub:        base.EventSub("EVENT_SUB_ORDER_STATUS"),
		Attributes: base.Attributes{Order: order},
	}

	event := New(baseEvent)

	if got := event.GetOrder(); got == nil {
		t.Error("GetOrder() retornou nil, esperava Order")
	} else if got.ID != "order-123" {
		t.Errorf("GetOrder().ID = %s, esperava order-123", got.ID)
	}
}

func TestOrderEvent_GetOrder_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_ORDER"),
		Category:   base.EventCategory("EVENT_CATEGORY_ORDER"),
		Sub:        base.EventSub("EVENT_SUB_ORDER_STATUS"),
		Attributes: base.Attributes{Order: nil},
	}

	event := New(baseEvent)

	if got := event.GetOrder(); got != nil {
		t.Error("GetOrder() retornou Order, esperava nil")
	}
}
