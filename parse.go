package main

import "strings"

// Parse will parse the csv file and return a DeviceData struct and will handle error/logging
func ParseRecord(record []string) DeviceData {

	if len(record) < 3 {
		WarnLog.Println("Invalid record: ", record)
		WarnLog.Println("Record missing fields")
		InvalidRecordCount++
		return DeviceData{}
	}

	if len(record) > 3 {
		WarnLog.Println("Invalid record: ", record)
		WarnLog.Println("Record has too many fields")
		InvalidRecordCount++
		return DeviceData{}
	}

	if !strings.HasPrefix(record[2], "SN-") {
		WarnLog.Println("Invalid record: ", record)
		WarnLog.Println("Record serial_number is invalid or in wrong position")
		InvalidRecordCount++
		return DeviceData{}
	}

	return DeviceData{
		device_type:   record[0],
		manufacturer:  record[1],
		serial_number: record[2],
	}
}
