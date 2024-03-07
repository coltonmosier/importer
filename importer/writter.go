package importer

import (
	"aswe-importer/models"
	"log"
)

var (
	acceptedDeviceTypes = map[string]int{"computer": 1, "laptop": 2, "mobile phone": 3, "smart watch": 4,
		"tablet": 5, "television": 6, "vehicle": 7}
	acceptedManufacturer = map[string]int{"Apple": 1, "Chevorlet": 2, "Dell": 3, "Ford": 4, "GM": 5, "Google": 6, "HP": 7,
		"Hisense": 8, "Huawei": 9, "Hyundai": 10, "IBM": 11, "KIA": 12, "LG": 13, "Microsoft": 14, "Motorola": 15, "Nissan": 16,
		"Nokia": 17, "OnePlus": 18, "Panasonic": 19, "Samsung": 20, "Sony": 21, "TCL": 22, "Toyota": 23, "Vizio": 24}
)

// WriteDeviceData writes the device data to the database
func WriteDeviceData(data []models.DeviceData) {
	defer Wg.Done()
	stmt, err := db.Prepare("INSERT INTO serial_numbers(device_type, manufacturer, serial_number) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatalf("%v: creating prepared statement\n", err)
	}

	for _, d := range data {
        dt := acceptedDeviceTypes[d.Device_type]
        m := acceptedManufacturer[d.Manufacturer]
		_, err := stmt.Exec(dt, m, d.Serial_number)
		if err != nil {
			log.Fatalf("%v: executing prepared statement %v\n", err, d)
		}
	}
	stmt.Close()
}
