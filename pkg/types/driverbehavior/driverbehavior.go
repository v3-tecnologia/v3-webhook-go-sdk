package driverbehavior

import (
	"encoding/json"
	"time"

	"go-eventlib/pkg/types/base"
)

type DriverBehaviorEventData struct {
	ID                 string                 `json:"id"`
	EventName          string                 `json:"event_name"`
	Timestamp          time.Time              `json:"timestamp"`
	AccelerationHarsh  map[string]interface{} `json:"acceleration_harsh,omitempty"`
	BrakingHarsh       map[string]interface{} `json:"braking_harsh,omitempty"`
	MaxSpeedFault      map[string]interface{} `json:"max_speed_fault,omitempty"`
	NormalSpeedReturn  map[string]interface{} `json:"normal_speed_return,omitempty"`
	PersistentMaxSpeed map[string]interface{} `json:"persistent_max_speed,omitempty"`
	SharpTurn          map[string]interface{} `json:"sharp_turn,omitempty"`
	StartOvertaking    map[string]interface{} `json:"start_overtaking,omitempty"`
}

type Event struct {
	*base.BaseEvent
}

func New(baseEvent *base.BaseEvent) *Event {
	return &Event{BaseEvent: baseEvent}
}

func (e *Event) GetDriverBehaviorData() *DriverBehaviorEventData {
	if e.Attributes.Data == nil || e.Attributes.Data.TripEvent == nil {
		return nil
	}

	data, err := json.Marshal(e.Attributes.Data.TripEvent)
	if err != nil {
		return nil
	}

	var tripEvent map[string]interface{}
	if err := json.Unmarshal(data, &tripEvent); err != nil {
		return nil
	}

	dbData, ok := tripEvent["driver_behavior"]
	if !ok {
		return nil
	}

	dbBytes, err := json.Marshal(dbData)
	if err != nil {
		return nil
	}

	var db DriverBehaviorEventData
	if err := json.Unmarshal(dbBytes, &db); err != nil {
		return nil
	}

	return &db
}

func (e *Event) GetHarshAccelerationData() map[string]interface{} {
	if db := e.GetDriverBehaviorData(); db != nil {
		return db.AccelerationHarsh
	}
	return nil
}

func (e *Event) GetHarshBrakingData() map[string]interface{} {
	if db := e.GetDriverBehaviorData(); db != nil {
		return db.BrakingHarsh
	}
	return nil
}

func (e *Event) GetEventName() string {
	if db := e.GetDriverBehaviorData(); db != nil {
		return db.EventName
	}
	return ""
}
