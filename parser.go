package main

import (
	"slices"
	"strings"
)


var acceptedDeviceTypes = []string{"computer", "laptop", "mobile phone", "smart watch",
	"tablet", "television", "vehicle"}
var acceptedManufacturer = []string{"Apple", "Chevorlet", "Dell", "Ford", "GM", "Google", "HP",
	"Hisense", "Huawei", "Hyundai", "IBM", "KIA", "LG", "Microsoft", "Motorola", "Nissan", "Nokia",
	"OnePlus", "Panasonic", "Samsung", "Sony", "TCL", "Toyota", "Vizio"}
var serialNumbers = []string{}

// Parse will parse the csv file and return a DeviceData struct and will handle error/logging
func ParseRecord(record []string) DeviceData {

	invalidRecord := strings.Join(record, ",")

    // strip single quotes
    for i := 1; i < len(record); i++ {
        if strings.Contains(record[i], "'") {
			record[i] = strings.ReplaceAll(record[i], "'", "")
        }
    }

	if len(record) < 4 {
		WarnLog.Printf("Invalid Record: missing fields [%s]\n", invalidRecord)
		InvalidRecordCount++
		return DeviceData{}
	}

	if len(record) > 4 {
		WarnLog.Printf("Invalid Record: too many fields [%s]\n", invalidRecord)
		InvalidRecordCount++
		return DeviceData{}
	}

	if !slices.Contains(acceptedDeviceTypes, record[1]) {
		WarnLog.Printf("Invalid Record: device_type invalid [%s]\n", invalidRecord)
		InvalidRecordCount++
		return DeviceData{}
	}

	if !slices.Contains(acceptedManufacturer, record[2]) {
		WarnLog.Printf("Invalid Record: manufacturer invalid [%s]\n", invalidRecord)
		InvalidRecordCount++
		return DeviceData{}
	}

    if slices.Contains(serialNumbers, record[3]) {
        WarnLog.Printf("Invalid Record: serial_number already exists [%s]\n", invalidRecord)
        InvalidRecordCount++
        return DeviceData{}
    }

	if !strings.HasPrefix(record[3], "SN-") {
		WarnLog.Printf("Invalid Record: serial_number invalid or in wrong position [%s]\n",
			invalidRecord)
		InvalidRecordCount++
		return DeviceData{}
	}

    if len(record[3]) != 67 {
        WarnLog.Printf("Invalid Record: serial_number invalid length [%s]\n", invalidRecord)
        InvalidRecordCount++
        return DeviceData{}
    }

    serialNumbers = append(serialNumbers, record[3])

	return DeviceData{
		device_type:   record[1],
		manufacturer:  record[2],
		serial_number: record[3],
	}
}
