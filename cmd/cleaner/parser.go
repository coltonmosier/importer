package cleaner

import (
	"fmt"
	"importer/models"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	acceptedDeviceTypes = []string{"computer", "laptop", "mobile phone", "smart watch",
		"tablet", "television", "vehicle"}
	acceptedManufacturer = []string{"Apple", "Chevorlet", "Dell", "Ford", "GM", "Google", "HP",
		"Hisense", "Huawei", "Hyundai", "IBM", "KIA", "LG", "Microsoft", "Motorola", "Nissan",
		"Nokia", "OnePlus", "Panasonic", "Samsung", "Sony", "TCL", "Toyota", "Vizio"}
)

// Parse will parse the csv file and return a DeviceData struct and will handle error/logging
func ParseDirtyRecord(r [][]string) ([]models.DeviceData, int) {
	invalidRecordCount := 0

	var d []models.DeviceData

	for i, record := range r {
		if i%1000 == 0 {
			Logger.WriteLogs()
		}
		invalidRecord := strings.Join(record, ",")
		serial, device_type, manufacturer := "", "", ""

		if len(record) < 4 {
			msg := fmt.Sprintf("Invalid Record: missing fields [%s]\n", invalidRecord)
			Logger.AddWarn(models.Message{Message: msg, Time: time.Now()})
			invalidRecordCount++
			continue
		}
		// NOTE: this should not error since we know every line has a line number...
		line_number, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
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

		// handle empty device_type
		if strings.Compare(device_type, "") == 0 {
			msg := fmt.Sprintf("Invalid Record: device_type missing [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			invalidRecordCount++
			continue
		}

		// handle empty manufacturer
		if strings.Compare(manufacturer, "") == 0 {
			msg := fmt.Sprintf("Invalid Record: manufacturer missing [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			invalidRecordCount++
			continue
		}

		// handle all serial number errors
		if strings.Compare(serial, "") == 0 {
			msg := fmt.Sprintf("Invalid Record: serial_number missing [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			invalidRecordCount++
			continue
		} else if len(serial) != 67 {
			msg := fmt.Sprintf("Invalid Record: serial_number invalid length [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			invalidRecordCount++
			continue
		}
		mu.Lock()
		if slices.Contains(SerialNumbers, serial) {
			mu.Unlock()
			msg := fmt.Sprintf("Invalid Record: serial_number already exists [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			invalidRecordCount++
			continue
		}
		mu.Unlock()

		mu.Lock()
		SerialNumbers = append(SerialNumbers, record[3])
		mu.Unlock()

		d = append(d, models.DeviceData{
			Line_number:   line_number,
			Device_type:   device_type,
			Manufacturer:  manufacturer,
			Serial_number: serial,
		})
	}
	return d, invalidRecordCount
}
