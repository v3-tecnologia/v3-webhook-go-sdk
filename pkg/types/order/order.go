package order

import (
	"github.com/v3-tecnologia/v3-webhook-go-sdk/pkg/types/base"
)

type Event struct {
	*base.BaseEvent
}

func New(baseEvent *base.BaseEvent) *Event {
	return &Event{BaseEvent: baseEvent}
}

func (e *Event) GetOrder() *base.Order {
	if e.Attributes.Order != nil {
		return e.Attributes.Order
	}
	return nil
}
