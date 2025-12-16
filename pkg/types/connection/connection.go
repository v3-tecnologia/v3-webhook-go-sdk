package connection

import (
	"encoding/json"
	"time"

	"go-eventlib/pkg/types/base"
)

type ConnectionType string

const (
	ConnectionTypeWiFi     ConnectionType = "CONNECTION_TYPE_WIFI"
	ConnectionTypeCellular ConnectionType = "CONNECTION_TYPE_CELLULAR"
	ConnectionTypeGPS      ConnectionType = "CONNECTION_TYPE_GPS"
)

type Connection struct {
	Type           ConnectionType `json:"type"`
	Area           int            `json:"area"`
	CellID         int            `json:"cell_id"`
	MCC            string         `json:"mcc"`
	MNC            string         `json:"mnc"`
	IMEI           string         `json:"imei"`
	SignalStrength interface{}    `json:"signal_strength"`
}

type ConnectionEventData struct {
	ID             string          `json:"id"`
	EventName      string          `json:"event_name"`
	WifiConnection *WifiConnection `json:"wifi_connection,omitempty"`
	SimCard        *SimCard        `json:"sim_card,omitempty"`
	Timestamp      time.Time       `json:"timestamp"`
}

type WifiConnection struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SimCard struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Event struct {
	*base.BaseEvent
}

func New(baseEvent *base.BaseEvent) *Event {
	return &Event{BaseEvent: baseEvent}
}

func (e *Event) GetConnectionData() *ConnectionEventData {
	if e.Attributes.Data == nil || e.Attributes.Data.StandaloneEvent == nil {
		return nil
	}

	data, err := json.Marshal(e.Attributes.Data.StandaloneEvent)
	if err != nil {
		return nil
	}

	var standalone map[string]interface{}
	if err := json.Unmarshal(data, &standalone); err != nil {
		return nil
	}

	connData, ok := standalone["connection"]
	if !ok {
		return nil
	}

	connBytes, err := json.Marshal(connData)
	if err != nil {
		return nil
	}

	var conn ConnectionEventData
	if err := json.Unmarshal(connBytes, &conn); err != nil {
		return nil
	}

	return &conn
}

func (e *Event) GetWifiConnection() *WifiConnection {
	if conn := e.GetConnectionData(); conn != nil {
		return conn.WifiConnection
	}
	return nil
}

func (e *Event) GetSimCard() *SimCard {
	if conn := e.GetConnectionData(); conn != nil {
		return conn.SimCard
	}
	return nil
}
