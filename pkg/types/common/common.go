package common

type Location struct {
	Method       string        `json:"method"`
	Coordinates  *Coordinates  `json:"coordinates,omitempty"`
	GNSS         *GNSS         `json:"gnss,omitempty"`
	Connectivity *Connectivity `json:"connectivity,omitempty"`
	Fix          *Fix          `json:"fix,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude,omitempty"`
	Speed     float64 `json:"speed,omitempty"`
}

type GNSS struct {
	GNSSClass  string  `json:"gnss_class"`
	Satellites int     `json:"satellites"`
	Fixed      bool    `json:"fixed"`
	GPSHeading float64 `json:"gps_heading,omitempty"`
	HDOP       float64 `json:"hdop,omitempty"`
	VDOP       float64 `json:"vdop,omitempty"`
}

type Fix struct {
	Timestamp          int64 `json:"timestamp"`
	LastTimestampOfFix int64 `json:"last_timestamp_of_fix"`
}

type Connectivity struct {
	SSID           string `json:"ssid,omitempty"`
	SignalStrength int    `json:"signal_strength,omitempty"`
}
