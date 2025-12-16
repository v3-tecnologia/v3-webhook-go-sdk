package base

import (
	"time"
)

type EventStatus string

const (
	StatusReceived  EventStatus = "STATUS_RECEIVED"
	StatusProcessed EventStatus = "STATUS_PROCESSED"
	StatusFailed    EventStatus = "STATUS_FAILED"
)

type EventType string

const (
	EventTypeOrder   EventType = "EVENT_TYPE_ORDER"
	EventTypeGeneral EventType = "EVENT_TYPE_GENERAL"
)

type EventCategory string

const (
	EventCategoryOrder          EventCategory = "EVENT_CATEGORY_ORDER"
	EventCategorySystem         EventCategory = "EVENT_CATEGORY_SYSTEM"
	EventCategoryHardware       EventCategory = "EVENT_CATEGORY_HARDWARE"
	EventCategoryConnection     EventCategory = "EVENT_CATEGORY_CONNECTION"
	EventCategoryVision         EventCategory = "EVENT_CATEGORY_VISION"
	EventCategoryTelemetry      EventCategory = "EVENT_CATEGORY_TELEMETRY"
	EventCategoryDMS            EventCategory = "EVENT_CATEGORY_DMS"
	EventCategoryDriverBehavior EventCategory = "EVENT_CATEGORY_DRIVER_BEHAVIOR"
	EventCategoryVehicle        EventCategory = "EVENT_CATEGORY_VEHICLE"
	EventCategoryAlert          EventCategory = "EVENT_CATEGORY_ALERT"
)

type EventSub string

const (
	EventSubOrderStatus EventSub = "EVENT_SUB_ORDER_STATUS"

	EventSubSystemUpload EventSub = "EVENT_SUB_SYSTEM_UPLOAD"

	EventSubConnectionStatusChanged EventSub = "EVENT_SUB_CONNECTION_STATUS_CHANGED"
	EventSubAlertCritical           EventSub = "EVENT_SUB_ALERT_CRITICAL"

	EventSubVisionBasic EventSub = "EVENT_SUB_VISION_BASIC"

	EventSubDMSBasic    EventSub = "EVENT_SUB_DMS_BASIC"
	EventSubDMSAdvanced EventSub = "EVENT_SUB_DMS_ADVANCED"

	EventSubDriverBehaviorAdvanced EventSub = "EVENT_SUB_DRIVER_BEHAVIOR_ADVANCED"

	EventSubAlertWarning EventSub = "EVENT_SUB_ALERT_WARNING"
	EventSubAlertInfo    EventSub = "EVENT_SUB_ALERT_INFO"

	EventSubTelemetryBattery  EventSub = "EVENT_SUB_TELEMETRY_BATTERY"
	EventSubTelemetryIgnition EventSub = "EVENT_SUB_TELEMETRY_IGNITION"
	EventSubTelemetryLocation EventSub = "EVENT_SUB_TELEMETRY_LOCATION"
)

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

const (
	OrderStatusAck    OrderStatus = "ORDER_STATUS_ACK"
	OrderStatusSent   OrderStatus = "ORDER_STATUS_SENT"
	OrderStatusFailed OrderStatus = "ORDER_STATUS_FAILED"
)

type OrderGroup string

const (
	OrderGroupConfig OrderGroup = "ORDER_GROUP_CONFIG"
)

type OrderType string

const (
	OrderTypeConfig OrderType = "CONFIG"
)

type Data struct {
	Telemetry       interface{} `json:"telemetry,omitempty"`
	GroupName       string      `json:"group_name,omitempty"`
	StandaloneEvent interface{} `json:"standalone_event,omitempty"`
	TripEvent       interface{} `json:"trip_event,omitempty"`
}
