package connection

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
)

func TestConnectionEvent_GetConnectionData(t *testing.T) {
	wifiConn := &WifiConnection{
		Name:   "TestWiFi",
		Status: "CONNECTED",
	}

	connData := map[string]interface{}{
		"id":              "conn-123",
		"event_name":      "WIFI_CONNECTED",
		"wifi_connection": wifiConn,
		"timestamp":       time.Now(),
	}

	standalone := map[string]interface{}{
		"event_group_name": "CONNECTION",
		"connection":       connData,
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_CONNECTION"),
		Sub:        base.EventSub("EVENT_SUB_CONNECTION_STATUS_CHANGED"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetConnectionData()
	if got == nil {
		t.Error("GetConnectionData() retornou nil, esperava ConnectionEventData")
		return
	}

	if got.ID != "conn-123" {
		t.Errorf("GetConnectionData().ID = %s, esperava conn-123", got.ID)
	}

	if got.EventName != "WIFI_CONNECTED" {
		t.Errorf("GetConnectionData().EventName = %s, esperava WIFI_CONNECTED", got.EventName)
	}
}

func TestConnectionEvent_GetWifiConnection(t *testing.T) {
	wifiConn := &WifiConnection{
		Name:   "TestWiFi",
		Status: "CONNECTED",
	}

	connData := map[string]interface{}{
		"id":              "conn-123",
		"event_name":      "WIFI_CONNECTED",
		"wifi_connection": wifiConn,
		"timestamp":       time.Now(),
	}

	standalone := map[string]interface{}{
		"event_group_name": "CONNECTION",
		"connection":       connData,
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_CONNECTION"),
		Sub:        base.EventSub("EVENT_SUB_CONNECTION_STATUS_CHANGED"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetWifiConnection()
	if got == nil {
		t.Error("GetWifiConnection() retornou nil, esperava WifiConnection")
		return
	}

	if got.Name != "TestWiFi" {
		t.Errorf("GetWifiConnection().Name = %s, esperava TestWiFi", got.Name)
	}

	if got.Status != "CONNECTED" {
		t.Errorf("GetWifiConnection().Status = %s, esperava CONNECTED", got.Status)
	}
}

func TestConnectionEvent_GetSimCard(t *testing.T) {
	simCard := &SimCard{
		Name:   "SIM_CARD",
		Status: "PRESENT",
	}

	connData := map[string]interface{}{
		"id":         "conn-123",
		"event_name": "SIMCARD",
		"sim_card":   simCard,
		"timestamp":  time.Now(),
	}

	standalone := map[string]interface{}{
		"event_group_name": "CONNECTION",
		"connection":       connData,
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_CONNECTION"),
		Sub:        base.EventSub("EVENT_SUB_CONNECTION_STATUS_CHANGED"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetSimCard()
	if got == nil {
		t.Error("GetSimCard() retornou nil, esperava SimCard")
		return
	}

	if got.Name != "SIM_CARD" {
		t.Errorf("GetSimCard().Name = %s, esperava SIM_CARD", got.Name)
	}

	if got.Status != "PRESENT" {
		t.Errorf("GetSimCard().Status = %s, esperava PRESENT", got.Status)
	}
}

func TestConnectionEvent_GetConnectionData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_CONNECTION"),
		Sub:        base.EventSub("EVENT_SUB_CONNECTION_STATUS_CHANGED"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetConnectionData(); got != nil {
		t.Error("GetConnectionData() retornou ConnectionEventData, esperava nil")
	}
}

func TestConnectionEvent_GetWifiConnection_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_CONNECTION"),
		Sub:        base.EventSub("EVENT_SUB_CONNECTION_STATUS_CHANGED"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetWifiConnection(); got != nil {
		t.Error("GetWifiConnection() retornou WifiConnection, esperava nil")
	}
}

func TestConnectionEvent_GetSimCard_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_CONNECTION"),
		Sub:        base.EventSub("EVENT_SUB_CONNECTION_STATUS_CHANGED"),
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetSimCard(); got != nil {
		t.Error("GetSimCard() retornou SimCard, esperava nil")
	}
}

func TestConnectionEvent_GetConnectionData_NoConnection(t *testing.T) {
	standalone := map[string]interface{}{
		"event_group_name": "CONNECTION",
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_CONNECTION"),
		Sub:        base.EventSub("EVENT_SUB_CONNECTION_STATUS_CHANGED"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetConnectionData(); got != nil {
		t.Error("GetConnectionData() retornou ConnectionEventData, esperava nil")
	}
}

func TestConnectionEvent_GetConnectionData_InvalidJSON(t *testing.T) {
	standalone := map[string]interface{}{
		"event_group_name": "CONNECTION",
		"connection":       make(chan int),
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:         "event-123",
		Status:     base.EventStatus("STATUS_RECEIVED"),
		CreatedAt:  time.Now(),
		Type:       base.EventType("EVENT_TYPE_GENERAL"),
		Category:   base.EventCategory("EVENT_CATEGORY_CONNECTION"),
		Sub:        base.EventSub("EVENT_SUB_CONNECTION_STATUS_CHANGED"),
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetConnectionData(); got != nil {
		t.Error("GetConnectionData() retornou ConnectionEventData, esperava nil para JSON inválido")
	}
}
