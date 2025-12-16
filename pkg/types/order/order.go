package order

import (
	"go-eventlib/pkg/types/base"
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
