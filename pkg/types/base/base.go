package base

import (
	"time"
)

type EventStatus string

type EventType string

type EventCategory string

type EventSub string

type Event interface {
	GetID() string
	GetCategory() EventCategory
	GetSubType() EventSub
	GetDeviceID() string
	GetCreatedAt() time.Time
}

type BaseEvent struct {
	ID         string        `json:"id"`
	Status     EventStatus   `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	Type       EventType     `json:"type"`
	Category   EventCategory `json:"category"`
	Sub        EventSub      `json:"sub"`
	Attributes Attributes    `json:"attributes"`
}

func (e *BaseEvent) GetID() string              { return e.ID }
func (e *BaseEvent) GetCategory() EventCategory { return e.Category }
func (e *BaseEvent) GetSubType() EventSub       { return e.Sub }
func (e *BaseEvent) GetDeviceID() string {
	if e.Attributes.Device != nil {
		return e.Attributes.Device.ID
	}
	return ""
}
func (e *BaseEvent) GetCreatedAt() time.Time { return e.CreatedAt }

type Device struct {
	ID            string   `json:"id"`
	CorrelationID string   `json:"correlation_id"`
	UID           string   `json:"uid"`
	AccountID     string   `json:"account_id"`
	Orders        []string `json:"orders,omitempty"`
}

type Attributes struct {
	Device *Device `json:"device,omitempty"`
	Data   *Data   `json:"data,omitempty"`
	Order  *Order  `json:"order,omitempty"`
}

type Order struct {
	ID            string      `json:"id"`
	CorrelationID string      `json:"correlation_id"`
	Group         OrderGroup  `json:"group"`
	Status        OrderStatus `json:"status"`
	Type          OrderType   `json:"type"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderStatus string

type OrderGroup string

type OrderType string

type Data struct {
	Telemetry       interface{} `json:"telemetry,omitempty"`
	GroupName       string      `json:"group_name,omitempty"`
	StandaloneEvent interface{} `json:"standalone_event,omitempty"`
	TripEvent       interface{} `json:"trip_event,omitempty"`
}
