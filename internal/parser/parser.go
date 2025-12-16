package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-eventlib/pkg/types/base"
	"go-eventlib/pkg/types/connection"
	"go-eventlib/pkg/types/telemetry"
)

// ParseEvent parses a JSON byte slice into a BaseEvent
func ParseEvent(data []byte) (*base.BaseEvent, error) {
	var event base.BaseEvent

	// Custom unmarshaling to handle time parsing
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Parse created_at time
	if createdAtStr, ok := raw["created_at"].(string); ok {
		if createdAt, err := time.Parse(time.RFC3339Nano, createdAtStr); err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		} else {
			raw["created_at"] = createdAt
		}
	}

	// Convert back to JSON and unmarshal into struct
	processedData, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal processed data: %w", err)
	}

	if err := json.Unmarshal(processedData, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into BaseEvent: %w", err)
	}

	return &event, nil
}

// ParseEventFromString parses a JSON string into a BaseEvent
func ParseEventFromString(jsonStr string) (*base.BaseEvent, error) {
	return ParseEvent([]byte(jsonStr))
}

// ValidateEvent validates that an event has required fields
func ValidateEvent(event *base.BaseEvent) error {
	if event.ID == "" {
		return errors.New("event ID is required")
	}
	if event.Type == "" {
		return errors.New("event type is required")
	}
	if event.Category == "" {
		return errors.New("event category is required")
	}
	if event.CreatedAt.IsZero() {
		return errors.New("event created_at is required")
	}
	if event.Attributes.Device == nil {
		return errors.New("event device information is required")
	}
	return nil
}

// GetEventCategory returns the category of an event
func GetEventCategory(event *base.BaseEvent) base.EventCategory {
	return event.Category
}

// GetEventType returns the type of an event
func GetEventType(event *base.BaseEvent) base.EventType {
	return event.Type
}

// GetEventSubType returns the sub-type of an event
func GetEventSubType(event *base.BaseEvent) base.EventSub {
	return event.Sub
}

// IsOrderEvent checks if the event is an order-related event
func IsOrderEvent(event *base.BaseEvent) bool {
	return event.Category == base.EventCategoryOrder
}

// IsHardwareEvent checks if the event is a hardware-related event
func IsHardwareEvent(event *base.BaseEvent) bool {
	return event.Category == base.EventCategoryHardware
}

// IsConnectionEvent checks if the event is a connection-related event
func IsConnectionEvent(event *base.BaseEvent) bool {
	return event.Category == base.EventCategoryConnection
}

// IsVisionEvent checks if the event is a vision-related event
func IsVisionEvent(event *base.BaseEvent) bool {
	return event.Category == base.EventCategoryVision
}

// IsSystemEvent checks if the event is a system-related event
func IsSystemEvent(event *base.BaseEvent) bool {
	return event.Category == base.EventCategorySystem
}

// IsTelemetryEvent checks if the event is a telemetry-related event
func IsTelemetryEvent(event *base.BaseEvent) bool {
	return event.Category == base.EventCategoryTelemetry
}

// GetDeviceID returns the device ID from the event
func GetDeviceID(event *base.BaseEvent) string {
	if event.Attributes.Device != nil {
		return event.Attributes.Device.ID
	}
	return ""
}

// GetDeviceUID returns the device UID from the event
func GetDeviceUID(event *base.BaseEvent) string {
	if event.Attributes.Device != nil {
		return event.Attributes.Device.UID
	}
	return ""
}

// GetAccountID returns the account ID from the event
func GetAccountID(event *base.BaseEvent) string {
	if event.Attributes.Device != nil {
		return event.Attributes.Device.AccountID
	}
	return ""
}

// GetIgnitionStatus returns the ignition status if available
func GetIgnitionStatus(event *base.BaseEvent) (telemetry.IgnitionStatus, bool) {
	if event.Attributes.Data != nil && event.Attributes.Data.Telemetry != nil {
		if telemetryData, ok := event.Attributes.Data.Telemetry.(map[string]interface{}); ok {
			if statusStr, ok := telemetryData["status"].(string); ok {
				return telemetry.IgnitionStatus(statusStr), true
			}
		}
	}
	return "", false
}

// GetConnectionType returns the connection type if available
func GetConnectionType(event *base.BaseEvent) (connection.ConnectionType, bool) {
	if event.Attributes.Data != nil && event.Attributes.Data.Telemetry != nil {
		if telemetryData, ok := event.Attributes.Data.Telemetry.(map[string]interface{}); ok {
			if connData, ok := telemetryData["connection"].(map[string]interface{}); ok {
				if typeStr, ok := connData["type"].(string); ok {
					return connection.ConnectionType(typeStr), true
				}
			}
		}
	}
	return "", false
}

// GetEventTimestamp returns the creation timestamp of the event
func GetEventTimestamp(event *base.BaseEvent) time.Time {
	return event.CreatedAt
}

// GetTelemetryTimestamp returns the telemetry timestamp if available
func GetTelemetryTimestamp(event *base.BaseEvent) (time.Time, bool) {
	if event.Attributes.Data != nil && event.Attributes.Data.Telemetry != nil {
		if telemetryData, ok := event.Attributes.Data.Telemetry.(map[string]interface{}); ok {
			if timestampStr, ok := telemetryData["timestamp"].(string); ok {
				if timestamp, err := time.Parse(time.RFC3339Nano, timestampStr); err == nil {
					return timestamp, true
				}
			}
		}
	}
	return time.Time{}, false
}
