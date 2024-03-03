package main

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

var acceptedDeviceTypes = []string{"computer", "laptop", "mobile phone", "smart watch",
	"tablet", "television", "vehicle"}
var acceptedManufacturer = []string{"Apple", "Chevorlet", "Dell", "Ford", "GM", "Google", "HP",
	"Hisense", "Huawei", "Hyundai", "IBM", "KIA", "LG", "Microsoft", "Motorola", "Nissan", "Nokia",
	"OnePlus", "Panasonic", "Samsung", "Sony", "TCL", "Toyota", "Vizio"}


// Parse will parse the csv file and return a DeviceData struct and will handle error/logging
func ParseRecord(r [][]string) []DeviceData {

	var d []DeviceData

	for i, record := range r {
		if i%1000 == 0 {
			Logger.WriteLogs()
			time.Sleep(50 * time.Millisecond)
		}
		invalidRecord := strings.Join(record, ",")
        serial, device_type, manufacturer := "", "", ""

		if len(record) < 4 {
			msg := fmt.Sprintf("Invalid Record: missing fields [%s]\n", invalidRecord)
			Logger.AddWarn(Message{msg, time.Now()})
			InvalidRecordCount++
			continue
		}

		for i := 1; i < len(record); i++ {
			if strings.Contains(record[i], "'") {
				record[i] = strings.ReplaceAll(record[i], "'", "")
			}
			if strings.HasPrefix(record[i], "SN-") {
				serial = record[i]
			} else if slices.Contains(acceptedDeviceTypes, record[i]) {
				device_type = record[i]
			} else if slices.Contains(acceptedManufacturer, record[i]) {
				manufacturer = record[i]
			}
		}

        if strings.Compare(serial, "") == 0 {
            msg := fmt.Sprintf("Invalid Record: serial_number missing [%s]\n", invalidRecord)
            Logger.AddWarn(Message{msg, time.Now()})
            InvalidRecordCount++
            continue
        } else if len(serial) != 67 {
            msg := fmt.Sprintf("Invalid Record: serial_number invalid length [%s]\n", invalidRecord)
            Logger.AddWarn(Message{msg, time.Now()})
            InvalidRecordCount++
            continue
        } 
        mu.Lock()
        if slices.Contains(SerialNumbers, serial) {
            mu.Unlock()
            msg := fmt.Sprintf("Invalid Record: serial_number already exists [%s]\n", invalidRecord)
            Logger.AddWarn(Message{msg, time.Now()})
            InvalidRecordCount++
            continue
        }
        mu.Unlock()

		mu.Lock()
		SerialNumbers = append(SerialNumbers, record[3])
		mu.Unlock()

		d = append(d, DeviceData{
			device_type:   device_type,
			manufacturer:  manufacturer,
			serial_number: serial,
		})
	}
	return d
}
