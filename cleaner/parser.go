package cleaner

import (
	"aswe-importer/models"
	"fmt"
	"slices"
	"strings"
	"time"
)

var (
	acceptedDeviceTypes = map[string]int{"computer": 1, "laptop": 2, "mobile phone": 3, "smart watch": 4,
		"tablet": 5, "television": 6, "vehicle": 7}
	acceptedManufacturer = map[string]int{"Apple": 1, "Chevorlet": 2, "Dell": 3, "Ford": 4, "GM": 5, "Google": 6, "HP": 7,
		"Hisense": 8, "Huawei": 9, "Hyundai": 10, "IBM": 11, "KIA": 12, "LG": 13, "Microsoft": 14, "Motorola": 15, "Nissan": 16,
		"Nokia": 17, "OnePlus": 18, "Panasonic": 19, "Samsung": 20, "Sony": 21, "TCL": 22, "Toyota": 23, "Vizio": 24}
)

// Parse will parse the csv file and return a DeviceData struct and will handle error/logging
func ParseDirtyRecord(r [][]string) ([]models.DeviceData, models.InvalidError) {

	var d []models.DeviceData
	var IE models.InvalidError

	for _, record := range r {
		invalidRecord := strings.Join(record, ",")
		line_number, serial, manufacturer, device_type := "", "", "", ""

		if len(record) < 4 {
			msg := fmt.Sprintf("Invalid Record: missing fields [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			Logger.AddBadData(models.Message{Message: invalidRecord})
			IE.MissingFields++
			continue
		}
		// NOTE: this should not error since we know every line has a line number...
		line_number = record[0]

		for i := 1; i < len(record); i++ {
			if strings.Contains(record[i], "'") {
				record[i] = strings.ReplaceAll(record[i], "'", "")
			}
			if strings.HasPrefix(record[i], "SN-") {
				serial = record[i]
			} else if _, ok := acceptedDeviceTypes[record[i]]; ok {
				device_type = record[i]
			} else if _, ok := acceptedManufacturer[record[i]]; ok {
				manufacturer = record[i]
			}
		}

		// handle empty device_type
		if strings.Compare(device_type, "") == 0 {
			msg := fmt.Sprintf("Invalid Record: device_type missing [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			Logger.AddBadData(models.Message{Message: invalidRecord})
			IE.DeviceTypeMissing++
			continue
		}

		// handle empty manufacturer
		if strings.Compare(manufacturer, "") == 0 {
			msg := fmt.Sprintf("Invalid Record: manufacturer missing [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			Logger.AddBadData(models.Message{Message: invalidRecord})
			IE.ManufacturerMissing++
			continue
		}

		// handle all serial number errors
		if strings.Compare(serial, "") == 0 {
			msg := fmt.Sprintf("Invalid Record: serial_number missing [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			Logger.AddBadData(models.Message{Message: invalidRecord})
			IE.SerialNumberMissing++
			continue
		} else if len(serial) != 67 {
			msg := fmt.Sprintf("Invalid Record: serial_number invalid length [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			Logger.AddBadData(models.Message{Message: invalidRecord})
			IE.SerialNumberLength++
			continue
		}
		mu.Lock()
		if slices.Contains(SerialNumbers, serial) {
			mu.Unlock()
			msg := fmt.Sprintf("Invalid Record: serial_number already exists [%s]\n", invalidRecord)
			Logger.AddErr(models.Message{Message: msg, Time: time.Now()})
			Logger.AddBadData(models.Message{Message: invalidRecord})
			IE.SerialNumberExists++
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
	return d, IE
}
