package common

import (
	"testing"
)

func TestLocation(t *testing.T) {
	loc := &Location{
		Method: "GNSS",
		Coordinates: &Coordinates{
			Latitude:  -23.5505,
			Longitude: -46.6333,
			Altitude:  760.0,
			Speed:     60.0,
		},
		GNSS: &GNSS{
			GNSSClass:  "GPS",
			Satellites: 12,
			Fixed:      true,
			GPSHeading: 45.0,
			HDOP:       1.2,
			VDOP:       1.5,
		},
		Connectivity: &Connectivity{
			SSID:           "TestWiFi",
			SignalStrength: -50,
		},
		Fix: &Fix{
			Timestamp:          1234567890,
			LastTimestampOfFix: 1234567890,
		},
	}

	if loc.Method != "GNSS" {
		t.Errorf("Location.Method = %s, esperava GNSS", loc.Method)
	}

	if loc.Coordinates == nil {
		t.Error("Location.Coordinates é nil")
	} else {
		if loc.Coordinates.Latitude != -23.5505 {
			t.Errorf("Coordinates.Latitude = %f, esperava -23.5505", loc.Coordinates.Latitude)
		}
		if loc.Coordinates.Longitude != -46.6333 {
			t.Errorf("Coordinates.Longitude = %f, esperava -46.6333", loc.Coordinates.Longitude)
		}
		if loc.Coordinates.Altitude != 760.0 {
			t.Errorf("Coordinates.Altitude = %f, esperava 760.0", loc.Coordinates.Altitude)
		}
		if loc.Coordinates.Speed != 60.0 {
			t.Errorf("Coordinates.Speed = %f, esperava 60.0", loc.Coordinates.Speed)
		}
	}

	if loc.GNSS == nil {
		t.Error("Location.GNSS é nil")
	} else {
		if loc.GNSS.GNSSClass != "GPS" {
			t.Errorf("GNSS.GNSSClass = %s, esperava GPS", loc.GNSS.GNSSClass)
		}
		if loc.GNSS.Satellites != 12 {
			t.Errorf("GNSS.Satellites = %d, esperava 12", loc.GNSS.Satellites)
		}
		if !loc.GNSS.Fixed {
			t.Error("GNSS.Fixed = false, esperava true")
		}
		if loc.GNSS.GPSHeading != 45.0 {
			t.Errorf("GNSS.GPSHeading = %f, esperava 45.0", loc.GNSS.GPSHeading)
		}
		if loc.GNSS.HDOP != 1.2 {
			t.Errorf("GNSS.HDOP = %f, esperava 1.2", loc.GNSS.HDOP)
		}
		if loc.GNSS.VDOP != 1.5 {
			t.Errorf("GNSS.VDOP = %f, esperava 1.5", loc.GNSS.VDOP)
		}
	}

	if loc.Connectivity == nil {
		t.Error("Location.Connectivity é nil")
	} else {
		if loc.Connectivity.SSID != "TestWiFi" {
			t.Errorf("Connectivity.SSID = %s, esperava TestWiFi", loc.Connectivity.SSID)
		}
		if loc.Connectivity.SignalStrength != -50 {
			t.Errorf("Connectivity.SignalStrength = %d, esperava -50", loc.Connectivity.SignalStrength)
		}
	}

	if loc.Fix == nil {
		t.Error("Location.Fix é nil")
	} else {
		if loc.Fix.Timestamp != 1234567890 {
			t.Errorf("Fix.Timestamp = %d, esperava 1234567890", loc.Fix.Timestamp)
		}
		if loc.Fix.LastTimestampOfFix != 1234567890 {
			t.Errorf("Fix.LastTimestampOfFix = %d, esperava 1234567890", loc.Fix.LastTimestampOfFix)
		}
	}
}

func TestLocation_NilFields(t *testing.T) {
	loc := &Location{
		Method: "GNSS",
	}

	if loc.Coordinates != nil {
		t.Error("Location.Coordinates deveria ser nil")
	}
	if loc.GNSS != nil {
		t.Error("Location.GNSS deveria ser nil")
	}
	if loc.Connectivity != nil {
		t.Error("Location.Connectivity deveria ser nil")
	}
	if loc.Fix != nil {
		t.Error("Location.Fix deveria ser nil")
	}
}

func TestCoordinates(t *testing.T) {
	coords := &Coordinates{
		Latitude:  -23.5505,
		Longitude: -46.6333,
		Altitude:  760.0,
		Speed:     60.0,
	}

	if coords.Latitude != -23.5505 {
		t.Errorf("Coordinates.Latitude = %f, esperava -23.5505", coords.Latitude)
	}
	if coords.Longitude != -46.6333 {
		t.Errorf("Coordinates.Longitude = %f, esperava -46.6333", coords.Longitude)
	}
	if coords.Altitude != 760.0 {
		t.Errorf("Coordinates.Altitude = %f, esperava 760.0", coords.Altitude)
	}
	if coords.Speed != 60.0 {
		t.Errorf("Coordinates.Speed = %f, esperava 60.0", coords.Speed)
	}
}

func TestGNSS(t *testing.T) {
	gnss := &GNSS{
		GNSSClass:  "GPS",
		Satellites: 12,
		Fixed:      true,
		GPSHeading: 45.0,
		HDOP:       1.2,
		VDOP:       1.5,
	}

	if gnss.GNSSClass != "GPS" {
		t.Errorf("GNSS.GNSSClass = %s, esperava GPS", gnss.GNSSClass)
	}
	if gnss.Satellites != 12 {
		t.Errorf("GNSS.Satellites = %d, esperava 12", gnss.Satellites)
	}
	if !gnss.Fixed {
		t.Error("GNSS.Fixed = false, esperava true")
	}
	if gnss.GPSHeading != 45.0 {
		t.Errorf("GNSS.GPSHeading = %f, esperava 45.0", gnss.GPSHeading)
	}
	if gnss.HDOP != 1.2 {
		t.Errorf("GNSS.HDOP = %f, esperava 1.2", gnss.HDOP)
	}
	if gnss.VDOP != 1.5 {
		t.Errorf("GNSS.VDOP = %f, esperava 1.5", gnss.VDOP)
	}
}

func TestConnectivity(t *testing.T) {
	conn := &Connectivity{
		SSID:           "TestWiFi",
		SignalStrength: -50,
	}

	if conn.SSID != "TestWiFi" {
		t.Errorf("Connectivity.SSID = %s, esperava TestWiFi", conn.SSID)
	}
	if conn.SignalStrength != -50 {
		t.Errorf("Connectivity.SignalStrength = %d, esperava -50", conn.SignalStrength)
	}
}

func TestFix(t *testing.T) {
	fix := &Fix{
		Timestamp:          1234567890,
		LastTimestampOfFix: 1234567890,
	}

	if fix.Timestamp != 1234567890 {
		t.Errorf("Fix.Timestamp = %d, esperava 1234567890", fix.Timestamp)
	}
	if fix.LastTimestampOfFix != 1234567890 {
		t.Errorf("Fix.LastTimestampOfFix = %d, esperava 1234567890", fix.LastTimestampOfFix)
	}
}
