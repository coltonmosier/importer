package models

import (
	"time"
)

type DeviceData struct {
	Line_number   int
	Device_type   string
	Manufacturer  string
	Serial_number string
}

type Message struct {
	Message string
	Time    time.Time
}
