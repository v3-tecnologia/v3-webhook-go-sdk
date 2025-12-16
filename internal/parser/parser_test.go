package parser

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
	"go-eventlib/pkg/types/connection"
	"go-eventlib/pkg/types/telemetry"
)

func TestParseEvent(t *testing.T) {
	eventJSON := `{
		"id": "01KCHNX0HKAAAZD4ZMTY35P8XJ",
		"status": "STATUS_RECEIVED",
		"created_at": "2025-12-15T18:55:59.748719972Z",
		"type": "EVENT_TYPE_ORDER",
		"category": "EVENT_CATEGORY_ORDER",
		"sub": "EVENT_SUB_ORDER_STATUS",
		"attributes": {
			"device": {
				"id": "01KBJ38J7DG4831CW327H1C908",
				"correlation_id": "01KBJ38J7D369CAFVGNWM6H2QX",
				"uid": "862798052131337",
				"account_id": "01GZXXCVVPEKM7E830XAMJKA14"
			},
			"order": {
				"id": "01KCHNWZGS5VXG8EY5SDZ9EVHW",
				"status": "ORDER_STATUS_ACK",
				"type": "CONFIG"
			}
		}
	}`

	event, err := ParseEventFromString(eventJSON)
	if err != nil {
		t.Fatalf("Erro ao fazer parse do evento: %v", err)
	}

	if event.ID != "01KCHNX0HKAAAZD4ZMTY35P8XJ" {
		t.Errorf("ID esperado: 01KCHNX0HKAAAZD4ZMTY35P8XJ, obtido: %s", event.ID)
	}

	if event.Category != base.EventCategoryOrder {
		t.Errorf("Categoria esperada: EVENT_CATEGORY_ORDER, obtida: %s", event.Category)
	}

	if event.Type != base.EventTypeOrder {
		t.Errorf("Tipo esperado: EVENT_TYPE_ORDER, obtido: %s", event.Type)
	}

	if event.Attributes.Order == nil {
		t.Error("Order não deveria ser nil")
	}

	if event.Attributes.Order.Status != base.OrderStatusAck {
		t.Errorf("Status da ordem esperado: ORDER_STATUS_ACK, obtido: %s", event.Attributes.Order.Status)
	}

	expectedTime, _ := time.Parse(time.RFC3339Nano, "2025-12-15T18:55:59.748719972Z")
	if !event.CreatedAt.Equal(expectedTime) {
		t.Errorf("Timestamp esperado: %v, obtido: %v", expectedTime, event.CreatedAt)
	}
}

func TestParseEventFromRealFiles(t *testing.T) {
	basePath := "../../test/events"
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		t.Skip("Pasta de eventos de exemplo não encontrada, pulando teste")
	}

	testFiles := []struct {
		path     string
		category base.EventCategory
	}{
		{"ack-events/ack-order-event.json", base.EventCategoryOrder},
		{"telemetry-events/telemetry-ignition.json", base.EventCategoryVehicle},
		{"dms-events/vision-drowsiness.json", base.EventCategoryDMS},
	}

	for _, tt := range testFiles {
		t.Run(tt.path, func(t *testing.T) {
			fullPath := filepath.Join(basePath, tt.path)
			data, err := os.ReadFile(fullPath)
			if err != nil {
				t.Skipf("Arquivo não encontrado: %s", fullPath)
			}

			event, err := ParseEvent(data)
			if err != nil {
				t.Fatalf("Erro ao fazer parse do evento %s: %v", tt.path, err)
			}

			if event.Category != tt.category {
				t.Errorf("Categoria esperada: %s, obtida: %s", tt.category, event.Category)
			}

			if event.ID == "" {
				t.Error("ID do evento não deveria ser vazio")
			}

			if event.Attributes.Device == nil {
				t.Error("Device não deveria ser nil")
			}
		})
	}
}

func TestValidateEvent(t *testing.T) {
	// Evento válido
	validEvent := &base.BaseEvent{
		ID:        "test-id",
		Type:      base.EventTypeOrder,
		Category:  base.EventCategoryOrder,
		CreatedAt: time.Now(),
		Attributes: base.Attributes{
			Device: &base.Device{
				ID: "device-001",
			},
		},
	}

	if err := ValidateEvent(validEvent); err != nil {
		t.Errorf("Evento válido falhou na validação: %v", err)
	}

	// Evento inválido - sem ID
	invalidEvent := &base.BaseEvent{
		Type:     base.EventTypeOrder,
		Category: base.EventCategoryOrder,
	}

	if err := ValidateEvent(invalidEvent); err == nil {
		t.Error("Evento inválido passou na validação")
	}
}

func TestEventTypeChecks(t *testing.T) {
	tests := []struct {
		name     string
		event    *base.BaseEvent
		checkFn  func(*base.BaseEvent) bool
		expected bool
	}{
		{"OrderEvent", &base.BaseEvent{Category: base.EventCategoryOrder}, IsOrderEvent, true},
		{"ConnectionEvent", &base.BaseEvent{Category: base.EventCategoryConnection}, IsConnectionEvent, true},
		{"VisionEvent", &base.BaseEvent{Category: base.EventCategoryVision}, IsVisionEvent, true},
		{"HardwareEvent", &base.BaseEvent{Category: base.EventCategoryHardware}, IsHardwareEvent, true},
		{"SystemEvent", &base.BaseEvent{Category: base.EventCategorySystem}, IsSystemEvent, true},
		{"TelemetryEvent", &base.BaseEvent{Category: base.EventCategoryTelemetry}, IsTelemetryEvent, true},
		{"OrderEvent wrong category", &base.BaseEvent{Category: base.EventCategoryConnection}, IsOrderEvent, false},
		{"DMS Event", &base.BaseEvent{Category: base.EventCategoryDMS}, IsVisionEvent, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.checkFn(tt.event)
			if result != tt.expected {
				t.Errorf("checkFn(%s) = %v, esperado %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGetDeviceID(t *testing.T) {
	event := &base.BaseEvent{
		Attributes: base.Attributes{
			Device: &base.Device{
				ID:        "device-001",
				UID:       "862798052131337",
				AccountID: "account-001",
			},
		},
	}

	deviceID := GetDeviceID(event)
	if deviceID != "device-001" {
		t.Errorf("Device ID esperado: device-001, obtido: %s", deviceID)
	}

	deviceUID := GetDeviceUID(event)
	if deviceUID != "862798052131337" {
		t.Errorf("Device UID esperado: 862798052131337, obtido: %s", deviceUID)
	}

	accountID := GetAccountID(event)
	if accountID != "account-001" {
		t.Errorf("Account ID esperado: account-001, obtido: %s", accountID)
	}

	eventNoDevice := &base.BaseEvent{}
	deviceID = GetDeviceID(eventNoDevice)
	if deviceID != "" {
		t.Errorf("Device ID de evento sem dispositivo deveria ser vazio, obtido: %s", deviceID)
	}
}

func TestGetTelemetryData(t *testing.T) {
	now := time.Now()
	event := &base.BaseEvent{
		Attributes: base.Attributes{
			Data: &base.Data{
				Telemetry: map[string]interface{}{
					"status":    "IGNITION_STATUS_ON",
					"timestamp": now.Format(time.RFC3339Nano),
					"connection": map[string]interface{}{
						"type": "CONNECTION_TYPE_WIFI",
					},
				},
			},
		},
	}

	ignitionStatus, ok := GetIgnitionStatus(event)
	if !ok {
		t.Error("GetIgnitionStatus deveria retornar true")
	}
	if ignitionStatus != telemetry.IgnitionStatusOn {
		t.Errorf("Ignition status esperado: IGNITION_STATUS_ON, obtido: %s", ignitionStatus)
	}

	connType, ok := GetConnectionType(event)
	if !ok {
		t.Error("GetConnectionType deveria retornar true")
	}
	if connType != connection.ConnectionTypeWiFi {
		t.Errorf("Connection type esperado: CONNECTION_TYPE_WIFI, obtido: %s", connType)
	}

	timestamp, ok := GetTelemetryTimestamp(event)
	if !ok {
		t.Error("GetTelemetryTimestamp deveria retornar true")
	}
	if timestamp.IsZero() {
		t.Error("Telemetry timestamp não deveria ser zero")
	}
}

func TestGetEventCategory(t *testing.T) {
	event := &base.BaseEvent{
		Category: base.EventCategoryOrder,
	}

	if got := GetEventCategory(event); got != base.EventCategoryOrder {
		t.Errorf("GetEventCategory() = %s, esperava EVENT_CATEGORY_ORDER", got)
	}
}

func TestGetEventType(t *testing.T) {
	event := &base.BaseEvent{
		Type: base.EventTypeOrder,
	}

	if got := GetEventType(event); got != base.EventTypeOrder {
		t.Errorf("GetEventType() = %s, esperava EVENT_TYPE_ORDER", got)
	}
}

func TestGetEventSubType(t *testing.T) {
	event := &base.BaseEvent{
		Sub: base.EventSubOrderStatus,
	}

	if got := GetEventSubType(event); got != base.EventSubOrderStatus {
		t.Errorf("GetEventSubType() = %s, esperava EVENT_SUB_ORDER_STATUS", got)
	}
}

func TestGetEventTimestamp(t *testing.T) {
	now := time.Now()
	event := &base.BaseEvent{
		CreatedAt: now,
	}

	if got := GetEventTimestamp(event); !got.Equal(now) {
		t.Errorf("GetEventTimestamp() = %v, esperava %v", got, now)
	}
}

func TestGetIgnitionStatus_False(t *testing.T) {
	event := &base.BaseEvent{
		Attributes: base.Attributes{
			Data: &base.Data{
				Telemetry: map[string]interface{}{},
			},
		},
	}

	_, ok := GetIgnitionStatus(event)
	if ok {
		t.Error("GetIgnitionStatus deveria retornar false quando status não está presente")
	}
}

func TestGetConnectionType_False(t *testing.T) {
	event := &base.BaseEvent{
		Attributes: base.Attributes{
			Data: &base.Data{
				Telemetry: map[string]interface{}{},
			},
		},
	}

	_, ok := GetConnectionType(event)
	if ok {
		t.Error("GetConnectionType deveria retornar false quando connection não está presente")
	}
}

func TestParseEvent_InvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json}`
	_, err := ParseEvent([]byte(invalidJSON))
	if err == nil {
		t.Error("ParseEvent deveria retornar erro para JSON inválido")
	}
}

func TestParseEvent_InvalidTimeFormat(t *testing.T) {
	invalidTimeJSON := `{
		"id": "test",
		"created_at": "invalid-time",
		"type": "EVENT_TYPE_ORDER",
		"category": "EVENT_CATEGORY_ORDER"
	}`
	_, err := ParseEvent([]byte(invalidTimeJSON))
	if err == nil {
		t.Error("ParseEvent deveria retornar erro para formato de tempo inválido")
	}
}

func TestValidateEvent_MissingFields(t *testing.T) {
	tests := []struct {
		name  string
		event *base.BaseEvent
	}{
		{"MissingID", &base.BaseEvent{Type: "TYPE", Category: "CATEGORY"}},
		{"MissingType", &base.BaseEvent{ID: "id", Category: "CATEGORY"}},
		{"MissingCategory", &base.BaseEvent{ID: "id", Type: "TYPE"}},
		{"MissingCreatedAt", &base.BaseEvent{ID: "id", Type: "TYPE", Category: "CATEGORY"}},
		{"MissingDevice", &base.BaseEvent{ID: "id", Type: "TYPE", Category: "CATEGORY", CreatedAt: time.Now()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateEvent(tt.event); err == nil {
				t.Errorf("ValidateEvent deveria retornar erro para %s", tt.name)
			}
		})
	}
}

func TestGetDeviceUID(t *testing.T) {
	event := &base.BaseEvent{
		Attributes: base.Attributes{
			Device: &base.Device{
				UID: "test-uid",
			},
		},
	}

	if got := GetDeviceUID(event); got != "test-uid" {
		t.Errorf("GetDeviceUID() = %s, esperava test-uid", got)
	}
}

func TestGetDeviceUID_Nil(t *testing.T) {
	event := &base.BaseEvent{
		Attributes: base.Attributes{
			Device: nil,
		},
	}

	if got := GetDeviceUID(event); got != "" {
		t.Errorf("GetDeviceUID() = %s, esperava string vazia", got)
	}
}

func TestGetAccountID(t *testing.T) {
	event := &base.BaseEvent{
		Attributes: base.Attributes{
			Device: &base.Device{
				AccountID: "test-account",
			},
		},
	}

	if got := GetAccountID(event); got != "test-account" {
		t.Errorf("GetAccountID() = %s, esperava test-account", got)
	}
}

func TestGetAccountID_Nil(t *testing.T) {
	event := &base.BaseEvent{
		Attributes: base.Attributes{
			Device: nil,
		},
	}

	if got := GetAccountID(event); got != "" {
		t.Errorf("GetAccountID() = %s, esperava string vazia", got)
	}
}

func TestGetTelemetryTimestamp_InvalidFormat(t *testing.T) {
	event := &base.BaseEvent{
		Attributes: base.Attributes{
			Data: &base.Data{
				Telemetry: map[string]interface{}{
					"timestamp": "invalid-time",
				},
			},
		},
	}

	_, ok := GetTelemetryTimestamp(event)
	if ok {
		t.Error("GetTelemetryTimestamp deveria retornar false para formato de tempo inválido")
	}
}
