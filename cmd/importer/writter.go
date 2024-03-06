package importer

import (
	"importer/models"
	"log"
)

// WriteDeviceData writes the device data to the database
func WriteDeviceData(dChan <-chan []models.DeviceData) {
	defer Wg.Done()
	data := <-dChan
	stmt, err := db.Prepare("INSERT INTO devices(device_type, manufacturer, serial_number) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatalf("%v: creating prepared statement\n", err)
	}

	for _, d := range data {
		_, err := stmt.Exec(d.Device_type, d.Manufacturer, d.Serial_number)
		if err != nil {
			log.Fatalf("%v: executing prepared statement %v\n", err, d)
		}
	}
	stmt.Close()
}
