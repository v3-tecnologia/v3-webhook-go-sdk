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
		Group:         base.OrderGroupConfig,
		Status:        base.OrderStatusAck,
		Type:          base.OrderTypeConfig,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeOrder,
		Category:   base.EventCategoryOrder,
		Sub:        base.EventSubOrderStatus,
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
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeOrder,
		Category:   base.EventCategoryOrder,
		Sub:        base.EventSubOrderStatus,
		Attributes: base.Attributes{Order: nil},
	}

	event := New(baseEvent)

	if got := event.GetOrder(); got != nil {
		t.Error("GetOrder() retornou Order, esperava nil")
	}
}
