package models

import (
	"fmt"
	"time"
)


type DeviceData struct {
    Line_number int
    Device_type  string
    Manufacturer string
    Serial_number string
}

func (d DeviceData) String() []string {
    s := []string{fmt.Sprintf("%v", d.Line_number), d.Device_type, d.Manufacturer, d.Serial_number}
    return s
}


type Message struct {
	Message string
	Time    time.Time
}
