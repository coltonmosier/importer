package models

import (
	"time"
)

type DeviceData struct {
	Line_number   string
	Device_type   string
	Manufacturer  string
	Serial_number string
}

type Message struct {
	Message string
	Time    time.Time
}

type InvalidError struct {
    ContainQuotes       int
    LongFields          int
	MissingFields       int
	DeviceTypeMissing   int
	ManufacturerMissing int
	SerialNumberMissing int
	SerialNumberLength  int
	SerialNumberExists  int
}
