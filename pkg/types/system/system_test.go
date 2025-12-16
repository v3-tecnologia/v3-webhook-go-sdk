package system

import (
	"testing"
	"time"

	"go-eventlib/pkg/types/base"
)

func TestSystemEvent_GetSystemData(t *testing.T) {
	systemData := map[string]interface{}{
		"id":         "system-123",
		"event_name": "UPLOAD",
		"timestamp":  time.Now(),
	}

	standalone := map[string]interface{}{
		"event_group_name": "SYSTEM",
		"system":           systemData,
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:        "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategorySystem,
		Sub:        base.EventSubSystemUpload,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetSystemData()
	if got == nil {
		t.Error("GetSystemData() retornou nil, esperava SystemEventData")
		return
	}

	if got.ID != "system-123" {
		t.Errorf("GetSystemData().ID = %s, esperava system-123", got.ID)
	}

	if got.EventName != "UPLOAD" {
		t.Errorf("GetSystemData().EventName = %s, esperava UPLOAD", got.EventName)
	}
}

func TestSystemEvent_GetUploadData(t *testing.T) {
	uploadData := map[string]interface{}{
		"name": "upload-test",
		"files": []interface{}{},
	}

	systemData := map[string]interface{}{
		"id":         "system-123",
		"event_name": "UPLOAD",
		"timestamp":  time.Now(),
		"upload":     uploadData,
	}

	standalone := map[string]interface{}{
		"event_group_name": "SYSTEM",
		"system":           systemData,
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:        "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategorySystem,
		Sub:        base.EventSubSystemUpload,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	got := event.GetUploadData()
	if got == nil {
		t.Error("GetUploadData() retornou nil, esperava UploadEvent")
		return
	}

	if got.Name != "upload-test" {
		t.Errorf("GetUploadData().Name = %s, esperava upload-test", got.Name)
	}
}

func TestSystemEvent_GetSystemData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:        "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategorySystem,
		Sub:        base.EventSubSystemUpload,
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetSystemData(); got != nil {
		t.Error("GetSystemData() retornou SystemEventData, esperava nil")
	}
}

func TestSystemEvent_GetUploadData_Nil(t *testing.T) {
	baseEvent := &base.BaseEvent{
		ID:        "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategorySystem,
		Sub:        base.EventSubSystemUpload,
		Attributes: base.Attributes{Data: nil},
	}

	event := New(baseEvent)

	if got := event.GetUploadData(); got != nil {
		t.Error("GetUploadData() retornou UploadEvent, esperava nil")
	}
}

func TestSystemEvent_GetSystemData_NoSystem(t *testing.T) {
	standalone := map[string]interface{}{
		"event_group_name": "SYSTEM",
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:        "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategorySystem,
		Sub:        base.EventSubSystemUpload,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetSystemData(); got != nil {
		t.Error("GetSystemData() retornou SystemEventData, esperava nil")
	}
}

func TestSystemEvent_GetSystemData_InvalidJSON(t *testing.T) {
	standalone := map[string]interface{}{
		"event_group_name": "SYSTEM",
		"system":          make(chan int),
	}

	data := &base.Data{
		StandaloneEvent: standalone,
	}

	baseEvent := &base.BaseEvent{
		ID:        "event-123",
		Status:     base.StatusReceived,
		CreatedAt:  time.Now(),
		Type:       base.EventTypeGeneral,
		Category:   base.EventCategorySystem,
		Sub:        base.EventSubSystemUpload,
		Attributes: base.Attributes{Data: data},
	}

	event := New(baseEvent)

	if got := event.GetSystemData(); got != nil {
		t.Error("GetSystemData() retornou SystemEventData, esperava nil para JSON inválido")
	}
}
