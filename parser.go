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

	for _, record := range r {
		invalidRecord := strings.Join(record, ",")

		// strip single quotes
		for i := 1; i < len(record); i++ {
			if strings.Contains(record[i], "'") {
				record[i] = strings.ReplaceAll(record[i], "'", "")
			}
		}

		if len(record) < 4 {
            fmt.Println("too short")
			msg := fmt.Sprintf("Invalid Record: missing fields [%s]\n", invalidRecord)
			Logger.AddWarn(Message{Message: msg, Time: time.Now()})
			InvalidRecordCount++
			continue
		}

		if len(record) > 4 {
            fmt.Println("too long")
			msg := fmt.Sprintf("Invalid Record: too many fields [%s]\n", invalidRecord)
			Logger.AddWarn(Message{Message: msg, Time: time.Now()})
			InvalidRecordCount++
			continue
		}

		if !slices.Contains(acceptedDeviceTypes, record[1]) {
            fmt.Println("device wrong")
			msg := fmt.Sprintf("Invalid Record: device_type invalid [%s]\n", invalidRecord)
			Logger.AddWarn(Message{Message: msg, Time: time.Now()})
			InvalidRecordCount++
			continue
		}

		if !slices.Contains(acceptedManufacturer, record[2]) {
            fmt.Println("manu wrong")
			msg := fmt.Sprintf("Invalid Record: manufacturer invalid [%s]\n", invalidRecord)
			Logger.AddWarn(Message{Message: msg, Time: time.Now()})
			InvalidRecordCount++
			continue
		}

		mu.Lock()
		if slices.Contains(SerialNumbers, record[3]) {
            fmt.Println("SN exists")
			msg := fmt.Sprintf("Invalid Record: serial_number already exists [%s]\n", invalidRecord)
			Logger.AddWarn(Message{Message: msg, Time: time.Now()})
			InvalidRecordCount++
			continue
		}
		mu.Unlock()

		if !strings.HasPrefix(record[3], "SN-") {
            fmt.Println("SN wrong")
			msg := fmt.Sprintf("Invalid Record: serial_number invalid or in wrong position [%s]\n", invalidRecord)
			Logger.AddWarn(Message{Message: msg, Time: time.Now()})
			InvalidRecordCount++
			continue
		}

		if len(record[3]) != 67 {
            fmt.Println("SN len")
			msg := fmt.Sprintf("Invalid Record: serial_number invalid length [%s]\n", invalidRecord)
			Logger.AddWarn(Message{Message: msg, Time: time.Now()})
			InvalidRecordCount++
			continue
		}

		mu.Lock()
		SerialNumbers = append(SerialNumbers, record[3])
		mu.Unlock()

		d = append(d, DeviceData{
			device_type:   record[1],
			manufacturer:  record[2],
			serial_number: record[3],
		})
        fmt.Println("Added to d:", len(d), "failed:", InvalidRecordCount)
	}
	return d
}
